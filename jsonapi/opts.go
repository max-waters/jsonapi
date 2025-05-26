package jsonapi

import (
	"encoding/json"
	"fmt"
)

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

type link struct {
	Uri        string
	LinkObject *linkObject
}

func (link) linkTag() {}

type linkObject struct {
	Href        string         `json:"href"`
	DescribedBy *link          `json:"described_by,omitempty"`
	Title       string         `json:"title,omitempty"`
	Type        string         `json:"type,omitempty"`
	HrefLang    []string       `json:"hreflang,omitempty"`
	Meta        map[string]any `json:"meta,omitempty"`
}

// UriLink returns a Link that represents a uri reference.
func UriLink(uri string) Link {
	return link{Uri: uri}
}

// LinkObject returns Link represnting a link object.
func LinkObject(href, title, linkType string, describedBy Link, hrefLang []string, meta map[string]any) Link {
	return link{
		LinkObject: &linkObject{
			Href:        href,
			DescribedBy: describedBy.(*link),
			Title:       title,
			Type:        linkType,
			HrefLang:    hrefLang,
			Meta:        meta,
		},
	}
}

func (l link) MarshalJSON() ([]byte, error) {
	if l.Uri != "" && l.LinkObject != nil {
		return nil, fmt.Errorf("uri and object defined")
	}
	if l.Uri != "" {
		return json.Marshal(l.Uri)
	}
	return json.Marshal(l.LinkObject)
}

func (l link) UnmarshalJSON(data []byte) error {
	switch data[0] {
	case '"':
		return json.Unmarshal(data, &l.Uri)
	case '{':
		return json.Unmarshal(data, &l.LinkObject)
	default:
		return fmt.Errorf("cannot unmarshal into link data")
	}
}

type ResourceLinks func(a any, r ResourceIdentifier) (map[string]Link, error)
type RelationshipLinks func(a any, r ResourceIdentifier, relationship string, data ...ResourceIdentifier) (map[string]Link, error)

type marshalResourceOpts struct {
	resourceLinks     ResourceLinks
	relationshipLinks RelationshipLinks
	resourceMeta      ResourceMeta
	relationshipMeta  RelationshipMeta
}

type marshalResourceOpt func(opts marshalResourceOpts) marshalResourceOpts

func WithResourceLinks(links ResourceLinks) marshalResourceOpt {
	return func(opts marshalResourceOpts) marshalResourceOpts {
		opts.resourceLinks = links
		return opts
	}
}

func WithRelationshipLinks(links RelationshipLinks) marshalResourceOpt {
	return func(opts marshalResourceOpts) marshalResourceOpts {
		opts.relationshipLinks = links
		return opts
	}
}

type ResourceMeta func(a any, r ResourceIdentifier) (map[string]any, error)
type RelationshipMeta func(a any, r ResourceIdentifier, relationship string, data ...ResourceIdentifier) (map[string]any, error)

func WithResourceMeta(meta ResourceMeta) marshalResourceOpt {
	return func(opts marshalResourceOpts) marshalResourceOpts {
		opts.resourceMeta = meta
		return opts
	}
}

func WithRelationshipMeta(meta RelationshipMeta) marshalResourceOpt {
	return func(opts marshalResourceOpts) marshalResourceOpts {
		opts.relationshipMeta = meta
		return opts
	}
}

func applyMarshalOpts(a any, r resource, opts marshalResourceOpts) error {
	if opts.resourceLinks == nil && opts.relationshipLinks == nil &&
		opts.resourceMeta == nil && opts.relationshipMeta == nil {
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
				return &MarshalErr{Field: name, Err: err}
			}
			r.Links[name] = l
		}
	}

	if opts.relationshipLinks != nil {
		for relName, rel := range r.ToOneRelationships {
			links, err := opts.relationshipLinks(a, rscId, relName, newResourceIdentifier(rel.Data))
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

			links, err := opts.relationshipLinks(a, rscId, relName, data...)
			if err != nil {
				return fmt.Errorf("adding links for relationship %s: %w", relName, err)
			}

			rel.Links, err = formatMap(links)
			if err != nil {
				return err
			}
		}
	}

	if opts.resourceMeta != nil {
		meta, err := opts.resourceMeta(a, rscId)
		if err != nil {
			return fmt.Errorf("adding resource meta: %w", err)
		}

		for name, m := range meta {
			l, err := json.Marshal(m)
			if err != nil {
				return &MarshalErr{Field: name, Err: err}
			}
			r.Meta[name] = l
		}
	}

	if opts.relationshipMeta != nil {
		for relName, rel := range r.ToOneRelationships {
			meta, err := opts.relationshipMeta(a, rscId, relName, newResourceIdentifier(rel.Data))
			if err != nil {
				return fmt.Errorf("adding meta for relationship %s: %w", relName, err)
			}

			rel.Meta, err = formatMap(meta)
			if err != nil {
				return err
			}
		}

		for relName, rel := range r.ToManyRelationships {
			data := make([]ResourceIdentifier, len(rel.Data))
			for i, d := range rel.Data {
				data[i] = newResourceIdentifier(d)
			}

			meta, err := opts.relationshipMeta(a, rscId, relName, data...)
			if err != nil {
				return fmt.Errorf("adding meta for relationship %s: %w", relName, err)
			}

			rel.Meta, err = formatMap(meta)
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
			return nil, &MarshalErr{Field: k, Err: err}
		}
		formatted[k] = f
	}
	return formatted, nil
}
