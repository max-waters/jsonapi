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

func TestLinkObject_MarshalJSON(t *testing.T) {
	want := LinkObject{
		DescribedBy: LinkUri{Uri: "uri"},
	}

	got := LinkObject{}

	if err := json.Unmarshal([]byte(`{ "described_by": "uri" }`), &got); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, want, got)
}

func TestLinkObject_UnmarshalJSON(t *testing.T) {
	// uri string
	want := LinkObject{
		DescribedBy: LinkUri{Uri: "uri"},
	}

	got := LinkObject{}

	if err := json.Unmarshal([]byte(`{ "described_by": "uri" }`), &got); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, want, got)

	// link object
	want = LinkObject{
		DescribedBy: LinkObject{
			Title:       "test",
			DescribedBy: LinkUri{"uri"},
		},
	}

	got = LinkObject{}

	if err := json.Unmarshal([]byte(`{ "described_by": { "described_by": "uri", "title": "test" } }`), &got); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, want, got)

	// no value
	want = LinkObject{}
	got = LinkObject{}

	if err := json.Unmarshal([]byte(`{ }`), &got); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, want, got)
}

func TestLinkUri_MarshalJson(t *testing.T) {
	linkUri := LinkUri{Uri: "uri"}

	got, err := json.Marshal(linkUri)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, `"uri"`, string(got))
}

func TestLinkUri_UnmarshalJson(t *testing.T) {
	linkUri := LinkUri{}
	if err := json.Unmarshal([]byte(`"uri"`), &linkUri); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "uri", linkUri.Uri)
}

func TestWithResourceLinksFunc(t *testing.T) {
	type testType struct {
		Id   string `jsonapi:"id,type"`
		Link string `jsonapi:"link,annotated_link"`
		Self string `jsonapi:"link,self"`
	}

	in := testType{Id: "id", Link: "link", Self: "overridden"}

	var gotRsc ResourceIdentifier
	var gotA any
	gotBytes, err := MarshalResource(in, WithResourceLinker(func(a any, r ResourceIdentifier) (map[string]Link, error) {
		gotRsc = r
		gotA = a
		return map[string]Link{
			"self": LinkUri{Uri: fmt.Sprintf("something.com/%s/%s", r.Type, r.Id)},
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
	wantRsc := ResourceIdentifier{
		Type: "type",
		Id:   "id",
	}

	assert.Equal(t, fmtJson(t, wantBytes), fmtJson(t, gotBytes))
	assert.Equal(t, wantRsc, gotRsc)
	assert.Equal(t, in, gotA)
}

func TestWWithRelationshipLinksFunc(t *testing.T) {
	type testType struct {
		Id        string   `jsonapi:"id,type"`
		ToOne     string   `jsonapi:"rel,rel1,relTyp1"`
		ToMany    []string `jsonapi:"rel,rel2,relTyp2"`
		Empty     string   `jsonapi:"rel,rel3,relTyp3"`
		OmitEmpty string   `jsonapi:"rel,rel4,relTyp3,omitempty"`
	}

	in := testType{Id: "id1", ToOne: "id2", ToMany: []string{"id3", "id4"}, Empty: ""}

	var gotRsc ResourceIdentifier
	var gotA any

	gotBytes, err := MarshalResource(in, WithRelationshipLinker(func(a any, rsc ResourceIdentifier, rel string, toOne bool, data ...ResourceIdentifier) (map[string]Link, error) {
		gotRsc = rsc
		gotA = a

		if rel == "rel1" {
			return map[string]Link{
				"self": LinkUri{Uri: fmt.Sprintf("something.com/%s/%s/%s/%s", rsc.Type, rsc.Id, rel, data[0].Id)},
			}, nil
		}

		if rel == "rel2" {
			return map[string]Link{
				"self": LinkUri{fmt.Sprintf("something.com/%s/%s/%s", rsc.Type, rsc.Id, rel)},
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

	wantRsc := ResourceIdentifier{
		Type: "type",
		Id:   "id1",
	}

	assert.Equal(t, fmtJson(t, wantBytes), fmtJson(t, gotBytes))
	assert.Equal(t, wantRsc, gotRsc)
	assert.Equal(t, in, gotA)
}
