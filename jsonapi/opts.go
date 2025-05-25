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

// The Link interface represents JSON:API link,
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

// UriLink returns a Link that represents a uri-reference.
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

type ResourceLinker func(r ResourceIdentifier) (map[string]Link, error)
type RelationshipLinker func(r ResourceIdentifier, relationship string, data ...ResourceIdentifier) (map[string]Link, error)

type marshalResourceOpts struct {
	resourceLinker     ResourceLinker
	relationshipLinker RelationshipLinker
}

type marshalResourceOpt func(opts marshalResourceOpts) marshalResourceOpts

func WithResourceLinker(linker ResourceLinker) marshalResourceOpt {
	return func(opts marshalResourceOpts) marshalResourceOpts {
		opts.resourceLinker = linker
		return opts
	}
}

func WithRelationshipLinker(linker RelationshipLinker) marshalResourceOpt {
	return func(opts marshalResourceOpts) marshalResourceOpts {
		opts.relationshipLinker = linker
		return opts
	}
}

func applyMarshalOpts(r resource, opts marshalResourceOpts) error {
	if opts.resourceLinker == nil && opts.relationshipLinker == nil {
		return nil
	}

	rscId := newResourceIdentifier(r.resourceIdentifier)

	if opts.resourceLinker != nil {
		links, err := opts.resourceLinker(rscId)
		if err != nil {
			return fmt.Errorf("formatting resource links: %w", err)
		}

		for name, link := range links {
			l, err := json.Marshal(link)
			if err != nil {
				return &MarshalErr{Field: name, Err: err}
			}
			r.Links[name] = l
		}
	}

	if opts.relationshipLinker != nil {
		for relName, rel := range r.ToOneRelationships {
			links, err := opts.relationshipLinker(rscId, relName, newResourceIdentifier(rel.Data))
			if err != nil {
				return fmt.Errorf("formatting links for relationship %s: %w", relName, err)
			}

			rel.Links, err = formatLinks(links)
			if err != nil {
				return err
			}
		}

		for relName, rel := range r.ToManyRelationships {
			data := make([]ResourceIdentifier, len(rel.Data))
			for i, d := range rel.Data {
				data[i] = newResourceIdentifier(d)
			}

			links, err := opts.relationshipLinker(rscId, relName, data...)
			if err != nil {
				return fmt.Errorf("formatting links for relationship %s: %w", relName, err)
			}

			rel.Links, err = formatLinks(links)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func formatLinks(links map[string]Link) (map[string]json.RawMessage, error) {
	jsonLinks := make(map[string]json.RawMessage, len(links))
	for name, link := range links {
		l, err := json.Marshal(link)
		if err != nil {
			return nil, &MarshalErr{Field: name, Err: err}
		}
		jsonLinks[name] = l
	}
	return jsonLinks, nil
}
