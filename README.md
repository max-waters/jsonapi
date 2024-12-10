# jsonapi #

Package `jsonapi` is a utility for marshaling and unmarshaling structs to and from [JSON:API v1.1](https://jsonapi.org/) formatted JSON.

```
type Article struct {
    ID       int    `jsonapi:"id,articles,string"`
    Title    string `jsonapi:"attr,title"`
    Author   int    `jsonapi:"rel,author,people,string"`
    Comments []int  `jsonapi:"rel,comments,comments"`
    Deleted  bool   `jsonapi:"meta,deleted"`
}

a := Article{
    ID:       1,
    Title:    "Hello World",
    Author:   2,
    Comments: []int{3, 4},
    Deleted:  false,
}

b, err := MarshalResource(&a)
if err != nil {
    panic(err)
}
fmt.Println(string(b))

// {
// 	"type": "articles",
// 	"id": "1",
// 	"meta": {
// 	  "deleted": false
// 	},
// 	"attributes": {
// 	  "title": "Hello World"
// 	},
// 	"relationships": {
// 	  "author": {
// 		"data": {
// 		  "type": "people",
// 		  "id": "2"
// 		}
// 	  },
// 	  "comments": {
// 		"data": [
// 		  {
// 			"type": "comments",
// 			"id": 3
// 		  },
// 		  {
// 			"type": "comments",
// 			"id": 4
// 		  }
// 		]
// 	  }
// 	}
// }
```

## Struct Tags ##

The `string` option signals that a field is stored as JSON inside a JSON-encoded string. 
It applies only to fields of floating point or integer types.

### IDs ###

The `id` tag defines the resource's primary id:

`jsonapi:"id,{type},[options]"`

The field's value is stored in the resource's `"id"` field, and the `{type}` argument is stored in the `"type"` field.

Marshaling and unmarshaling of both `"id"` and `"type"` use the `encoding/json` package.

Note that fields of floating point or integer types can be used as valid JSON:API identifiers with the `string` option.

### attributes ###

An `attr` tag defines an attribute:

`jsonapi:"attr,{name},[options]"`

The field's value is stored in the resource's attributes, under the key defined by `{name}`.

If the `{name}` argument is omitted, then the `encoding/json` default is used instead, ie either the name defined in the `json` tag, or the declared field name if none is found.

### relationships ###

A `rel` tag defines a relationship:

`jsonapi:"rel,{name},{type},[options]"`

### metadata ###

A `meta` tag defines a metadata item:

`jsonapi:"meta,{name},[options]"`

The field's value is stored in the resource's metadata, under the key defined by `{name}`.

