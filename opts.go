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

// ResourceIdentifier represents a JSON:API resource identifier.
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

// LinkObject represents a JSON:API object link.
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

// LinkObject represents a JSON:API URI-reference link.
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

// A ResourceLinker should return all links for
// resource r with JSON:API formatted identifier id.
type ResourceLinker func(r any, id ResourceIdentifier) (map[string]Link, error)

// A RelationshipLinker should return all links for the given
// relationship, where r is the parent resource with JSON:API
// formatted identifier id, rel is the name of the relationship
// and data is the JSON:API formatted identifiers for all related data.
type RelationshipLinker func(r any, id ResourceIdentifier, rel string, toOne bool, data ...ResourceIdentifier) (map[string]Link, error)

// WithResourceLinks returns a resource marshaling option that will
// retrieve resource links with the supplied function.
func WithResourceLinker(links ResourceLinker) marshalResourceOpt {
	return func(opts marshalResourceOpts) marshalResourceOpts {
		opts.resourceLinks = links
		return opts
	}
}

// WithRelationshipLinker returns a resource marshaling option that will
// retrieve relationship links with the supplied function.
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
