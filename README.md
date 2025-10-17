# jsonapi #

The `jsonapi` package is a utility for marshaling and unmarshaling Go structs to and from [JSON:API v1.1](https://jsonapi.org/) formatted JSON.

Features:
- Struct tags define the mapping between struct fields and the JSON:API id, attributes, relationships and metadata.
- Marshaling and unmarshaling behaviour can be customised by implementing the `ResourceMarshaler` and `ResourceUnmarshaler` interfaces, respectively.
- Exposes an API similar to the standard `encoding/json` package. 
- Supports anonymous/embedded struct fields.

Planned feaures:
- Strict mode that enforces JSON:API compliant output.
- Marshaling and unmarshaling arrays of resources.
- Marshaling and unmarshaling [top-level](https://jsonapi.org/format/1.0/#document-top-level) JSON:API documents

## Usage ##

Import with 

```
import (
	"github.com/max-waters/jsonapi"
)
```

Two functions are exposed:

```Go
MarshalResource(a any) ([]byte, error)
UnmarshalResource(data []byte, a any) error
```

`MarshalResource` returns the JSON:API encoding of `a`, and `UnmarshalResource` parses the JSON:API-encoded bytes `data` and stores the result in the value pointed to by `a`.

### Example ###

Go code:

```Go
type Article struct {
    ID       int    `jsonapi:"id,articles,string"`
    Title    string `jsonapi:"attr,title"`
    Author   int    `jsonapi:"rel,author,people,string"`
    Comments []int  `jsonapi:"rel,comments,comments,string"`
    Deleted  bool   `jsonapi:"meta,deleted"`
}

a := Article{
    ID:       1,
    Title:    "Hello World",
    Author:   2,
    Comments: []int{3, 4},
    Deleted:  false,
}

b, err := jsonapi.MarshalResource(&a)
if err != nil {
  // handle error
}
fmt.Println(string(b))
```

The resulting JSON:API:

```json
{
  "type": "articles",
  "id": "1",
  "meta": {
    "deleted": false
  },
  "attributes": {
    "title": "Hello World"
  },
  "relationships": {
    "author": {
      "data": {
        "type": "people",
        "id": "2"
      }
    },
    "comments": {
      "data": [
        {
          "type": "comments",
          "id": "3"
        },
        {
          "type": "comments",
          "id": "4"
        }
      ]
    }
  }
}
```

## Mapping structs to JSON:API ##

The mapping between struct fields and the JSON:API id, attributes, relationships and metadata is defined with struct tags. The marshal and unmarshal functions will look for these tags in the top-level fields of in the input struct, and those in any anonymous struct fields. The values of these fields are marshaled and unmarshaled into the appropriate location in the resulting JSON:API using the `encoding/json` package.

Note that if a struct field has no `jsonapi` tag, then it is assumed to be an attribute (see below) with the `encoding/json` default name. A struct field can be exlcuded from the mapping with the "ignore" tag, `jsonapi:"-"`.

### IDs ###

The `id` tag defines the resource's primary id:

```Go
`jsonapi:"id,{type},[options]"`
```

The tagged field's value is mapped to the resource's `"id"` field, and the `{type}` argument defines the content of the `"type"` field. The field value is marshaled and unmarshaled with the `encoding/json` package.

Note that the `jsonapi` package does not (currently) enforce the JSON:API requirement that the `"id"` field be a string. However, the `string` option will encode floating point or integer values as JSON strings, allowing them to be used as valid JSON:API identifiers.

The `omitempty` option will exclude zero-valued values from the resulting JSON, allowing for empty IDs (eg for server-side ID generation).

#### Example ID with `string` option ####

Struct tags:

```Go
type Article struct {
    ID int    `jsonapi:"id,articles,string"`
}

a := Article{
    ID: 1,
}
```

JSON:API:

```json
{
  "type": "articles",
  "id": "1",
}
```

#### Example ID with `omitempty` option ####

Struct tags:

```Go
type Article struct {
    ID int `jsonapi:"id,articles,string,omitempty"`
}

a := Article{}
```

JSON:API:

```json
{
  "type": "articles"
}
```


### Attributes ###

An attribute is defined by providing either an `attr` tag, or no `jsonapi` tag at all:

```Go
`jsonapi:"attr,{name},[options]"`
```

The field's value will be mapped to an attribute with the key specified by `{name}`. If no `jsonapi` tag is defined, or the `{name}` argument is empty, then the `encoding/json` default is used instead, ie either the name defined in the `json` tag, or the declared field name if none is found. The field value is marshaled and unmarshaled with the `encoding/json` package.

The `attr` tag supports the `string` and `omitempty` options, which encode numeric values as JSON strings, and omit zero-valued fields, respectively.

#### Example Attributes ####

Struct tags:

```Go
type Copyright struct {
    Owner string    `json:"owner"`
    Date  time.Time `json:"date"`
}

type Article struct {
    Title    string     `json:"title"`
    Content  string     `jsonapi:"content,omitempty"`
    Copyright Copyright `jsonapi:"attr,copyright"`
}

a := Article{
    Title: "Hello World",
    Copyright: Copyright {
        Owner: "Publishing Ltd",
        Date:  time.Now()
    }
}
```

JSON:API:

```json
{
  "attributes": {
    "title": "Hello World",
    "copyright": {
        "owner": "Publishing Ltd",
        "date": "2024-12-12T21:46:43.552855+11:00"
    }
  },
}
```

### Relationships ###

The `rel` tag defines a relationship:

```Go
`jsonapi:"rel,{name},{type},[options]"`
```

Any field annotated with a `rel` tag will be mapped to a relationship with the key specified by `{name}`. If the `{name}` argument is empty, then the `encoding/json` default is used instead, ie either the name defined in the `json` tag, or the declared field name if none is found. 

The field's declared type determines whether it maps to a to-one or a to-many relationship. Array and slices of anything except bytes (or pointers to these), will be mapped to a to-many relationship, and all other types are mapped to a to-one relationship. For to-one relationships, the field's value maps to the relationship's `"id"` field, and the `{type}` argument defines the `"type"` field. For to-many relationships, each element in the array or slice defines the `"id"` of a related resource. The IDs are marshaled and unmarshaled with the `encoding/json` package.

While the `jsonapi` package does not (currently) enforce the JSON:API requirement that the `"id"` field be a string, the `string` option will encode floating point or integer IDs as JSON strings, allowing them to be used as valid JSON:API identifiers. And the `omitempty` option will exclude relationships with zero-valued valued IDs from the resulting JSON.

#### Example To-One and To-Many Relationships with `string` option ####

Struct tags:

```Go
type Article struct {
    Author   int    `jsonapi:"rel,author,people,string"`
    Comments []int  `jsonapi:"rel,comments,comments,string"`
}

a := Article{
    Author:   2,
    Comments: []int{3, 4},
}
```

JSON:API:

```json
{
  "relationships": {
    "author": {
      "data": {
        "type": "people",
        "id": "2"
      }
    },
    "comments": {
      "data": [
        {
          "type": "comments",
          "id": "3"
        },
        {
          "type": "comments",
          "id": "4"
        }
      ]
    }
  }
}
```

#### Example Relationship with `omitempty` option ####

Struct tags:

```Go
type Article struct {
    Author   int    `jsonapi:"rel,author,people,string"`
    Comments []int  `jsonapi:"rel,comments,comments,string"`
}

a := Article{
    Comments: []int{3, 4},
}
```

JSON:API:

```json
{
  "relationships": {
    "comments": {
      "data": [
        {
          "type": "comments",
          "id": "3"
        },
        {
          "type": "comments",
          "id": "4"
        }
      ]
    }
  }
}
```

### Metadata ###

The `meta` tag defines a metadata item:

```Go
`jsonapi:"meta,{name},[options]"`
```

The field's value will be mapped to a metadata item with the key specified by `{name}`. If the `{name}` argument is empty, then the `encoding/json` default is used instead, ie either the name defined in the `json` tag, or the declared field name if none is found.  The field value is marshaled and unmarshaled with the `encoding/json` package.

The `meta` tag supports the `string` and `omitempty` options, which encode numeric values as JSON strings, and omit zero-valued fields, respectively.

### Links ###

The `link` tag defines a link:

```Go
`jsonapi:"link,{name},[options]"`
```

The field's value will be mapped to a link with the key specified by `{name}`. If the `{name}` argument is empty, then the `encoding/json` default is used instead, ie either the name defined in the `json` tag, or the declared field name if none is found.  The field value is marshaled and unmarshaled with the `encoding/json` package.

Note that the `jsonapi` package does not (currently) enforce the requirement that all links be either valid URIs or link objects.

The `link` tag supports the `string` and `omitempty` options, which encode numeric values as JSON strings, and omit zero-valued fields, respectively.

## Anonymous Struct Fields ##

Anonymous (ie, embedded) struct fields are "promoted" and treated as though their members are declared in their parent type:

```Go
type SessionAuthz struct {
    Editable  bool `jsonapi:"attr,editable"`
    Deletable bool `jsonapi:"attr,deletable"`
}

type Article struct {
    SessionAuthz
    Title string `jsonapi:"attr,title"`
}

a := Article{
    SessionAuthz: SessionAuthz{
        Editable: true,
        Deletable: false,
    },
    Title: "Hello World",
}
```

JSON:API:

```JSON
{
  "attributes": {
    "deletable": false,
    "editable": true,
    "title": "Hello World"
  }
}
```

Names clashes are resolved with standard Go promotion rules, as used by the `encoding/json` package. If two or more `attr`, `rel` or `meta` fields have the same name, then a selection is made based on the fields' nesting depth, then the presence of a `jsonapi` tag, then the presence of a `json` tag. If no single preferred field is found, then all clashing fields are excluded from the marshaling and unmarshaling.

## Customising Resource Marshaling and Unmarshaling ##

The `jsonapi` package provides two interfaces and a number of functional options to help with custom marshaling and unmarshaling.

### `ResourceMarshaler` and `ResourceUnmarshaler` ###

These interfaces allow a type to marshal or unmarshal itself to or from JSON:API:

```Go
type ResourceMarshaler interface {
    MarshalJsonApiResource() ([]byte, error)
}

type ResourceUnmarshaler interface {
    UnmarshalJsonApiResource([]byte) error
}
```

#### Example `ResourceMarshaler` and `ResourceUnmarshaler` ####

In this example, the `Article` type formats the `created` attribute as an `RFC3339` timestamp by creating an alias type, and then calling the `jsonapi` marshaling and unmarshaling functions:

```Go
type Article struct {
    ID      int
    Created time.Time
}

func (a *Article) MarshalJsonApiResource() ([]byte, error) {
    type alias struct {
        ID      int    `jsonapi:"id,articles,string"`
        Created string `jsonapi:"attr,created"`
    }

    b := alias{
        ID:      a.ID
        Created: a.Created.Format(time.RFC3339),
    }

    return jsonapi.MarshalResource(&b)
}

func (a *Article) UnmarshalJsonApiResource(data []byte) error {
    type alias struct {
        ID      int    `jsonapi:"id,articles,string"`
        Created string `jsonapi:"attr,created"`
    }

    b := alias{}

    if err := jsonapi.UnmarshalResource(data, &b); err != nil {
        return err
    }

    created, err := time.Parse(time.RFC3339, b.Created)
    if err != nil {
        return err
    }

    a.ID = b.ID
    a.Created = created
    return nil
}
```

### Custom Link Marshaling with Functional Options ###

The functional options `WithResourceLinker` and `WithRelationshipLinker` allow for resource and relationship links to be generated from struct tags as well as information not contained in the resource.

The `WithResourceLinker` option accepts a function with the signature `func(a any, r jsonapi.ResourceIdentifier) (map[string]jsonapi.Link, error)` that should return all links for the supplied resource and JSON:API id. The function will be called for every resource.

The `WithRelationshipLinker` option accepts a function with the signature `func(r any, id jsonapi.ResourceIdentifier, rel string, toOne bool, data ...jsonapi.ResourceIdentifier) (map[string]jsonapi.Link, error)` that should return all links for the supplied relationship on the supplied resource. The function will be called for every relationship on every resource.

#### Example `WithResourceLinker` and `WithRelationshipLinker` Function ####

In this example, `self` and `related` links are generated from the type names declared in the resource's struct tags and a pre-configured URL base:

```Go
func InitResourceLinker(urlBase string) jsonapi.ResourceLinker {
	return func(r any, id jsonapi.ResourceIdentifier) (map[string]jsonapi.Link, error) {
		return map[string]jsonapi.Link{
			"self": jsonapi.LinkUri{Uri: fmt.Sprintf("%s/%s/%s", urlBase, id.Type, id.Id)},
		}, nil
	}
}

func InitRelationshipLinker(urlBase string) jsonapi.RelationshipLinker {
	return func(r any, id jsonapi.ResourceIdentifier, rel string, toOne bool, data ...jsonapi.ResourceIdentifier) (map[string]jsonapi.Link, error) {
		return map[string]Link{
			"self":    jsonapi.LinkUri{Uri: fmt.Sprintf("%s/%s/%s/relationships/%s", urlBase, id.Type, id.Id, rel)},
			"related": jsonapi.LinkUri{Uri: fmt.Sprintf("%s/%s/%s/%s", urlBase, id.Type, id.Id, rel)},
		}, nil
	}
}

type Article struct {
    ID       int
    Comments []int
}

a := Article{
    Id: 4,
    Comments: []int{5}
}

urlBase := "https://example.com"

jsonapi.MarshalResource(a, 
  jsonapi.WithResourceLinker(InitResourceLinker(urlBase)),
  jsonapi.WithRelationshipLinker(InitRelationshipLinker(urlBase)),
)
```

JSON:API:

```JSON
{
  "type": "articles",
  "id": "4",
  "relationships": {
    "comments": {
      "data": [ 
        { 
          "type": "comment", 
          "id": 5 
        } 
      ],
      "links": {
        "related": "https://example.com/articles/4/comments",
        "self": "https://example.com/articles/4/relationships/comments"
      }
    }
  },
  "links": {
    "self": "https://example.com/articles/4"
  }
}
```
