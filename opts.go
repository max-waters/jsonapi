package jsonapi

import (
	"encoding/json"
	"fmt"
)

type marshalResourceOpts struct {
	resourceLinks     ResourceLinker
	relationshipLinks RelationshipLinker
}

type marshalResourceOpt func(opts marshalResourceOpts) marshalResourceOpts

// A ResourceIdentifier represents a JSON:API [resource identifier].
//
// [resource identifier]: https://jsonapi.org/format/#document-resource-identifier-objects
type ResourceIdentifier struct {
	Type string
	Id   string
}

func newResourceIdentifier(r resourceIdentifier) ResourceIdentifier {
	ri := ResourceIdentifier{
		Type: r.Type,
	}

	if r.Id != nil {
		if r.Id[0] == '"' {
			ri.Id = string(r.Id[1 : len(r.Id)-1])
		} else {
			ri.Id = string(r.Id)
		}
	}

	return ri
}

// The Link interface represents a JSON:API link,
// ie either a uri reference or an object. It has
// a single unexported method and so cannot be
// implemented by code outside of this package.
type Link interface {
	linkTag()
}

// A LinkObject represents a JSON:API [object link].
//
// [object link]: https://jsonapi.org/format/#auto-id--link-objects
type LinkObject struct {
	Href        string         `json:"href,omitempty"`
	DescribedBy Link           `json:"described_by,omitempty"`
	Title       string         `json:"title,omitempty"`
	Type        string         `json:"type,omitempty"`
	HrefLang    []string       `json:"hreflang,omitempty"`
	Meta        map[string]any `json:"meta,omitempty"`
}

func (LinkObject) linkTag() {}

func (l *LinkObject) UnmarshalJSON(data []byte) error {
	type alias LinkObject

	type proxy struct {
		alias
		DescribedBy json.RawMessage `json:"described_by,omitempty"`
	}

	a := proxy{
		alias: alias(*l),
	}

	if err := json.Unmarshal(data, &a); err != nil {
		return err
	}

	*l = LinkObject(a.alias)

	if len(a.DescribedBy) > 0 {
		switch a.DescribedBy[0] {
		case '"':
			linkUri := LinkUri{}
			if err := json.Unmarshal(a.DescribedBy, &linkUri); err != nil {
				return err
			}
			l.DescribedBy = linkUri
		case '{':
			linkObj := LinkObject{}
			if err := json.Unmarshal(a.DescribedBy, &linkObj); err != nil {
				return err
			}
			l.DescribedBy = linkObj
		default:
			return fmt.Errorf("cannot unmarshal \"described_by\", unexpected byte %d", a.DescribedBy)
		}
	}

	return nil
}

// A LinkUri represents a JSON:API [URI-reference link].
//
// [URI-reference link]: https://jsonapi.org/format/#document-links
type LinkUri struct {
	Uri string
}

func (LinkUri) linkTag() {}

func (l LinkUri) MarshalJSON() ([]byte, error) {
	return json.Marshal(l.Uri)
}

func (l *LinkUri) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &l.Uri)
}

// A ResourceLinker function should return all links for the
// supplied resource. Input r is the value passed to [MarshalResource],
// and id is the [ResourceIdentifier] extracted from r.
// Can be used with [WithResourceLinker] to set links that
// are not present in the struct passed to [MarshalResource].
//
// Example:
//
//	func MyResourceLinker(r any, id ResourceIdentifier) (map[string]Link, error) {
//		return map[string]Link{
//			"self": LinkUri{
//				Uri: fmt.Sprintf("https://example.com/%s/%s", id.Type, id.Id),
//			},
//		}, nil
//	}
type ResourceLinker func(r any, id ResourceIdentifier) (map[string]Link, error)

// A RelationshipLinker function should return all links for the given
// relationship, on the given resource. Input r is the value passed to [MarshalResource],
// id is the [ResourceIdentifier] extracted from r, rel is the relationship name,
// toOne indicates if the relationship is toOne or toMany, and data contains a
// [ResourceIdentifier] for each related resource.
// Can be used with [WithRelationshipLinker] to set links that
// are not present in the struct passed to [MarshalResource].
//
// Example:
//
//	func MyRelationshipLinker(r any, id ResourceIdentifier, rel string, toOne bool, data ...ResourceIdentifier) (map[string]Link, error) {
//		return map[string]Link{
//			"related": LinkUri{Uri: fmt.Sprintf("https://example.com/%s/%s/%s", id.Type, id.Id, rel)},
//		}, nil
//	}
type RelationshipLinker func(r any, id ResourceIdentifier, rel string, toOne bool, data ...ResourceIdentifier) (map[string]Link, error)

// WithResourceLinker returns a resource marshaling option that will
// set resource links to those returned by the suppled [ResourceLinker] function.
func WithResourceLinker(links ResourceLinker) marshalResourceOpt {
	return func(opts marshalResourceOpts) marshalResourceOpts {
		opts.resourceLinks = links
		return opts
	}
}

// WithRelationshipLinker returns a resource marshaling option that will
// set links on relationships to those returned by the
// the supplied [RelationshipLinker] function.
func WithRelationshipLinker(links RelationshipLinker) marshalResourceOpt {
	return func(opts marshalResourceOpts) marshalResourceOpts {
		opts.relationshipLinks = links
		return opts
	}
}

func applyMarshalOpts(a any, r resource, opts marshalResourceOpts) error {
	if opts.resourceLinks == nil && opts.relationshipLinks == nil {
		return nil
	}

	rscId := newResourceIdentifier(r.resourceIdentifier)

	if opts.resourceLinks != nil {
		links, err := opts.resourceLinks(a, rscId)
		if err != nil {
			return fmt.Errorf("adding resource links: %w", err)
		}

		for name, link := range links {
			l, err := json.Marshal(link)
			if err != nil {
				return &MarshalError{Field: name, Err: err}
			}
			r.Links[name] = l
		}
	}

	if opts.relationshipLinks != nil {
		for relName, rel := range r.ToOneRelationships {
			links, err := opts.relationshipLinks(a, rscId, relName, true, newResourceIdentifier(rel.Data))
			if err != nil {
				return fmt.Errorf("adding links for relationship %s: %w", relName, err)
			}

			rel.Links, err = formatMap(links)
			if err != nil {
				return err
			}
		}

		for relName, rel := range r.ToManyRelationships {
			data := make([]ResourceIdentifier, len(rel.Data))
			for i, d := range rel.Data {
				data[i] = newResourceIdentifier(d)
			}

			links, err := opts.relationshipLinks(a, rscId, relName, false, data...)
			if err != nil {
				return fmt.Errorf("adding links for relationship %s: %w", relName, err)
			}

			rel.Links, err = formatMap(links)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func formatMap[A any](m map[string]A) (map[string]json.RawMessage, error) {
	formatted := make(map[string]json.RawMessage, len(m))
	for k, v := range m {
		f, err := json.Marshal(v)
		if err != nil {
			return nil, &MarshalError{Field: k, Err: err}
		}
		formatted[k] = f
	}
	return formatted, nil
}
