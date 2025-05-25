package jsonapi

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewResourceIdentifier(t *testing.T) {
	type testCase struct {
		In   resourceIdentifier
		Want ResourceIdentifier
	}

	testCases := []testCase{
		{
			In: resourceIdentifier{
				Type: "type",
				Id:   json.RawMessage([]byte(`"id"`)),
			},
			Want: ResourceIdentifier{
				Type: "type",
				Id:   "id",
			},
		},
		{
			In: resourceIdentifier{
				Type: "type",
				Id:   json.RawMessage([]byte(`{"k":"v"}`)),
			},
			Want: ResourceIdentifier{
				Type: "type",
				Id:   `{"k":"v"}`,
			},
		},
		{
			In: resourceIdentifier{
				Type: "type",
				Id:   nil,
			},
			Want: ResourceIdentifier{
				Type: "type",
				Id:   "",
			},
		},
	}

	for _, tc := range testCases {
		assert.Equal(t, tc.Want, newResourceIdentifier(tc.In))
	}
}

func TestLinkObject(t *testing.T) {
	want := link{
		LinkObject: &linkObject{
			Href:        "href",
			Title:       "title",
			Type:        "type",
			DescribedBy: &link{Uri: "uri"},
			HrefLang:    []string{"lang"},
			Meta:        map[string]any{"key": "value"},
		},
	}

	got := LinkObject("href", "title", "type", &link{Uri: "uri"}, []string{"lang"}, map[string]any{"key": "value"})
	assert.Equal(t, want, got)
}

func TestLinkUri(t *testing.T) {
	want := link{Uri: "uri"}
	got := UriLink("uri")
	assert.Equal(t, want, got)
}

func TestWithResourceLinker(t *testing.T) {
	type testType struct {
		Id   string `jsonapi:"id,type"`
		Link string `jsonapi:"link,annotated_link"`
		Self string `jsonapi:"link,self"`
	}

	in := testType{Id: "id", Link: "link", Self: "overridden"}

	var gotArg ResourceIdentifier
	gotBytes, err := MarshalResource(in, WithResourceLinker(func(r ResourceIdentifier) (map[string]Link, error) {
		gotArg = r
		return map[string]Link{
			"self": UriLink(fmt.Sprintf("something.com/%s/%s", r.Type, r.Id)),
			"null": nil,
		}, nil
	}))

	if err != nil {
		t.Fatal(err)
	}

	wantBytes := []byte(`
	{
		"id": "id",
		"type": "type",
		"links": {
			"self": "something.com/type/id",
			"annotated_link": "link",
			"null": null
		}
	}`)
	wantArg := ResourceIdentifier{
		Type: "type",
		Id:   "id",
	}

	assert.Equal(t, fmtJson(t, wantBytes), fmtJson(t, gotBytes))
	assert.Equal(t, wantArg, gotArg)
}

func TestWithRelationshipLinker(t *testing.T) {
	type testType struct {
		Id        string   `jsonapi:"id,type"`
		ToOne     string   `jsonapi:"rel,rel1,relTyp1"`
		ToMany    []string `jsonapi:"rel,rel2,relTyp2"`
		Empty     string   `jsonapi:"rel,rel3,relTyp3"`
		OmitEmpty string   `jsonapi:"rel,rel4,relTyp3,omitempty"`
	}

	in := testType{Id: "id1", ToOne: "id2", ToMany: []string{"id3", "id4"}, Empty: ""}

	gotBytes, err := MarshalResource(in, WithRelationshipLinker(func(rsc ResourceIdentifier, rel string, data ...ResourceIdentifier) (map[string]Link, error) {
		if rel == "rel1" {
			return map[string]Link{
				"self": UriLink(fmt.Sprintf("something.com/%s/%s/%s/%s", rsc.Type, rsc.Id, rel, data[0].Id)),
			}, nil
		}

		if rel == "rel2" {
			return map[string]Link{
				"self": UriLink(fmt.Sprintf("something.com/%s/%s/%s", rsc.Type, rsc.Id, rel)),
			}, nil
		}

		if rel == "rel3" {
			return map[string]Link{
				"self": nil,
			}, nil
		}

		t.Fatalf("unexpected relationship: %s", rel)

		return nil, nil

	}))

	if err != nil {
		t.Fatal(err)
	}

	wantBytes := []byte(`
	{
		"id": "id1",
		"type": "type",
		"relationships": {
			"rel1": {
				"data": { "id": "id2", "type": "relTyp1" },
				"links": {
					"self": "something.com/type/id1/rel1/id2"
				}
				
			},
			"rel2": {
				"data": [ { "id": "id3", "type": "relTyp2" }, { "id": "id4", "type": "relTyp2" }  ],
				"links": {
					"self": "something.com/type/id1/rel2"
				}	
			},
			"rel3": {
				"data": { "id": "", "type": "relTyp3" },
				"links": {
					"self": null
				}	
			}
		}
	}`)

	assert.Equal(t, fmtJson(t, wantBytes), fmtJson(t, gotBytes))
}
