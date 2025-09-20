package jsonapi

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func BenchmarkResourceMarshalJSON_Empty(b *testing.B) {
	b.ReportAllocs()

	// currently: 2 allocations
	r := resource{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := r.MarshalJSON(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkResourceMarshalJSON(b *testing.B) {
	b.ReportAllocs()

	// currently: 18 allocations
	r := resource{
		resourceIdentifier: resourceIdentifier{
			Id: json.RawMessage([]byte(`"1"`)), Type: "type",
		},
		Attributes: map[string]json.RawMessage{
			"a": json.RawMessage([]byte(`"val"`)),
		},
		ToOneRelationships: map[string]*toOneRelationship{
			"r1": {
				Data: resourceIdentifier{
					Id: json.RawMessage([]byte(`"1"`)), Type: "type",
				},
			},
		},
		ToManyRelationships: map[string]*toManyRelationship{
			"r2": {
				Data: []resourceIdentifier{
					{Id: json.RawMessage([]byte(`"1"`)), Type: "type"},
				},
			},
		},
		Links: map[string]json.RawMessage{
			"m": json.RawMessage([]byte(`"val"`)),
		},
		Meta: map[string]json.RawMessage{
			"l": json.RawMessage([]byte(`"val"`)),
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := r.MarshalJSON(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkResourceUnmarshalJSON_Empty(b *testing.B) {
	b.ReportAllocs()

	// currently: 5 allocations
	r := resource{}
	data := []byte(`{}`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := r.UnmarshalJSON(data); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkResourceUnmarshalJSON(b *testing.B) {
	b.ReportAllocs()

	// currently: 57 allocations
	r := resource{}
	data := []byte(`{"type":"type","id":"1","attributes":{"a":"val"},"relationships":{"r1":{"data":{"type":"type","id":"1"}},"r2":{"data":[{"type":"type","id":"1"}]}},"links":{"m":"val"},"meta":{"l":"val"}}`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := r.UnmarshalJSON(data); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkFormat_Empty(b *testing.B) {
	b.ReportAllocs()

	// currently: 5 allocations
	type T struct{}
	in := T{}
	v := reflect.ValueOf(in)
	o := marshalResourceOpts{}

	b.ResetTimer()
	for _, cache := range []bool{false, true} {
		b.Run(fmt.Sprintf("cache=%v", cache), func(b *testing.B) {
			fieldCache.Clear()
			enableFieldCache = cache
			for i := 0; i < b.N; i++ {
				if _, err := format(in, v, o); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkFormat_Id(b *testing.B) {
	b.ReportAllocs()

	// currently: 10 allocations
	type T struct {
		Id int `jsonapi:"id,type"`
	}
	in := T{
		Id: 1,
	}
	v := reflect.ValueOf(in)
	o := marshalResourceOpts{}

	b.ResetTimer()
	for _, cache := range []bool{false, true} {
		b.Run(fmt.Sprintf("cache=%v", cache), func(b *testing.B) {
			fieldCache.Clear()
			enableFieldCache = cache
			for i := 0; i < b.N; i++ {
				if _, err := format(in, v, o); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkFormat_Attributes(b *testing.B) {
	b.ReportAllocs()

	// currently: 18 allocations
	type T struct {
		A, B, C int
	}
	in := T{
		A: 1, B: 2, C: 3,
	}
	v := reflect.ValueOf(in)
	o := marshalResourceOpts{}

	b.ResetTimer()
	for _, cache := range []bool{false, true} {
		b.Run(fmt.Sprintf("cache=%v", cache), func(b *testing.B) {
			fieldCache.Clear()
			enableFieldCache = cache
			for i := 0; i < b.N; i++ {
				if _, err := format(in, v, o); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkFormat_ToOneRelationships(b *testing.B) {
	b.ReportAllocs()

	// currently: 17 allocations
	type T struct {
		A int `jsonapi:"rel,type-a,name-a"`
		B int `jsonapi:"rel,type-b,name-b"`
	}
	in := T{
		A: 1, B: 2,
	}
	v := reflect.ValueOf(in)
	o := marshalResourceOpts{}

	b.ResetTimer()
	for _, cache := range []bool{false, true} {
		b.Run(fmt.Sprintf("cache=%v", cache), func(b *testing.B) {
			fieldCache.Clear()
			enableFieldCache = cache
			for i := 0; i < b.N; i++ {
				if _, err := format(in, v, o); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkFormat_ToManyRelationships(b *testing.B) {
	b.ReportAllocs()

	// currently: 36 allocations
	type T struct {
		A []int `jsonapi:"rel,type-a,name-a"`
		B []int `jsonapi:"rel,type-b,name-b"`
	}
	in := T{
		A: []int{1, 2}, B: []int{3, 4},
	}
	v := reflect.ValueOf(in)
	o := marshalResourceOpts{}

	b.ResetTimer()
	for _, cache := range []bool{false, true} {
		b.Run(fmt.Sprintf("cache=%v", cache), func(b *testing.B) {
			fieldCache.Clear()
			enableFieldCache = cache
			for i := 0; i < b.N; i++ {
				if _, err := format(in, v, o); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkDeformat_Empty(b *testing.B) {
	b.ReportAllocs()

	// currently: 0 allocations
	type T struct{}
	v, err := derefInput(reflect.ValueOf(&T{}), resourceUnmarshalerType)
	if err != nil {
		b.Fatal(err)
	}
	r := resource{}

	b.ResetTimer()
	for _, cache := range []bool{false, true} {
		b.Run(fmt.Sprintf("cache=%v", cache), func(b *testing.B) {
			fieldCache.Clear()
			enableFieldCache = cache
			for i := 0; i < b.N; i++ {
				if err := deformat(v, r); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkDeformat_Id(b *testing.B) {
	b.ReportAllocs()

	// currently: 6 allocations
	type T struct {
		Id int `jsonapi:"id,type,string"`
	}

	v, err := derefInput(reflect.ValueOf(&T{}), resourceUnmarshalerType)
	if err != nil {
		b.Fatal(err)
	}
	r := resource{
		resourceIdentifier: resourceIdentifier{
			Id:   json.RawMessage([]byte(`"1"`)),
			Type: "type",
		},
	}

	b.ResetTimer()
	for _, cache := range []bool{false, true} {
		b.Run(fmt.Sprintf("cache=%v", cache), func(b *testing.B) {
			fieldCache.Clear()
			enableFieldCache = cache
			for i := 0; i < b.N; i++ {
				if err := deformat(v, r); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkDeformat_Attributes(b *testing.B) {
	b.ReportAllocs()

	// currently: 14 allocations
	type T struct {
		A, B, C int
	}

	v, err := derefInput(reflect.ValueOf(&T{}), resourceUnmarshalerType)
	if err != nil {
		b.Fatal(err)
	}
	r := resource{
		Attributes: map[string]json.RawMessage{
			"A": json.RawMessage([]byte(`1`)),
			"B": json.RawMessage([]byte(`1`)),
			"C": json.RawMessage([]byte(`1`)),
		},
	}

	b.ResetTimer()
	for _, cache := range []bool{false, true} {
		b.Run(fmt.Sprintf("cache=%v", cache), func(b *testing.B) {
			fieldCache.Clear()
			enableFieldCache = cache
			for i := 0; i < b.N; i++ {
				if err := deformat(v, r); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkDeformat_ToOneRelationships(b *testing.B) {
	b.ReportAllocs()

	// currently: 6 allocations
	type T struct {
		A int `jsonapi:"rel,type-a,name-a"`
		B int `jsonapi:"rel,type-b,name-b"`
	}
	v, err := derefInput(reflect.ValueOf(&T{}), resourceUnmarshalerType)
	if err != nil {
		b.Fatal(err)
	}
	r := resource{
		ToOneRelationships: map[string]*toOneRelationship{
			"name-a": {
				Data: resourceIdentifier{
					Type: "type-a",
					Id:   json.RawMessage([]byte(`1`)),
				},
			},
			"name-b": {
				Data: resourceIdentifier{
					Type: "type-b",
					Id:   json.RawMessage([]byte(`2`)),
				},
			},
		},
	}

	b.ResetTimer()
	for _, cache := range []bool{false, true} {
		b.Run(fmt.Sprintf("cache=%v", cache), func(b *testing.B) {
			fieldCache.Clear()
			enableFieldCache = cache
			for i := 0; i < b.N; i++ {
				if err := deformat(v, r); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkDeformat_ToManyRelationships(b *testing.B) {
	b.ReportAllocs()

	// currently: 6 allocations
	type T struct {
		A []int `jsonapi:"rel,type-a,name-a"`
		B []int `jsonapi:"rel,type-b,name-b"`
	}
	v, err := derefInput(reflect.ValueOf(&T{}), resourceUnmarshalerType)
	if err != nil {
		b.Fatal(err)
	}
	r := resource{
		ToManyRelationships: map[string]*toManyRelationship{
			"name-a": {
				Data: []resourceIdentifier{
					{Type: "type-a", Id: json.RawMessage([]byte(`1`))},
				},
			},
			"name-b": {
				Data: []resourceIdentifier{
					{Type: "type-b", Id: json.RawMessage([]byte(`2`))},
				},
			},
		},
	}

	b.ResetTimer()
	for _, cache := range []bool{false, true} {
		b.Run(fmt.Sprintf("cache=%v", cache), func(b *testing.B) {
			fieldCache.Clear()
			enableFieldCache = cache
			for i := 0; i < b.N; i++ {
				if err := deformat(v, r); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkParseTags_Empty(b *testing.B) {
	b.ReportAllocs()

	// currently: 0 allocations
	type T struct{}
	v := reflect.ValueOf(T{})

	b.ResetTimer()
	for _, cache := range []bool{false, true} {
		b.Run(fmt.Sprintf("cache=%v", cache), func(b *testing.B) {
			fieldCache.Clear()
			enableFieldCache = cache
			for i := 0; i < b.N; i++ {
				if _, err := parseTags(v); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkParseTags_Id(b *testing.B) {
	b.ReportAllocs()

	// currently: num allocations = (2 x num fields) + 2 = 4
	type T struct {
		Id int `jsonapi:"id,type"`
	}
	v := reflect.ValueOf(T{})

	b.ResetTimer()
	for _, cache := range []bool{false, true} {
		b.Run(fmt.Sprintf("cache=%v", cache), func(b *testing.B) {
			fieldCache.Clear()
			enableFieldCache = cache
			for i := 0; i < b.N; i++ {
				if _, err := parseTags(v); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkParseTags_Attributes(b *testing.B) {
	b.ReportAllocs()

	// currently: num allocations = (2 x num fields) + 2 = 8
	type T struct {
		A, B, C int
	}
	v := reflect.ValueOf(T{})

	b.ResetTimer()
	for _, cache := range []bool{false, true} {
		b.Run(fmt.Sprintf("cache=%v", cache), func(b *testing.B) {
			fieldCache.Clear()
			enableFieldCache = cache
			for i := 0; i < b.N; i++ {
				if _, err := parseTags(v); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkParseTags_ToOneRelationships(b *testing.B) {
	b.ReportAllocs()

	// currently: num allocations = (2 x num fields) + 2 = 6
	type T struct {
		A int `jsonapi:"rel,type-a,name-a"`
		B int `jsonapi:"rel,type-b,name-b"`
	}
	v := reflect.ValueOf(T{})

	b.ResetTimer()
	for _, cache := range []bool{false, true} {
		b.Run(fmt.Sprintf("cache=%v", cache), func(b *testing.B) {
			fieldCache.Clear()
			enableFieldCache = cache
			for i := 0; i < b.N; i++ {
				if _, err := parseTags(v); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkParseTags_ToManyRelationships(b *testing.B) {
	b.ReportAllocs()

	// currently: num allocations = (2 x num fields) + 2 = 6
	type T struct {
		A []int `jsonapi:"rel,type-a,name-a"`
		B []int `jsonapi:"rel,type-b,name-b"`
	}
	v := reflect.ValueOf(T{})

	b.ResetTimer()
	for _, cache := range []bool{false, true} {
		b.Run(fmt.Sprintf("cache=%v", cache), func(b *testing.B) {
			fieldCache.Clear()
			enableFieldCache = cache
			for i := 0; i < b.N; i++ {
				if _, err := parseTags(v); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkParseTags_Anons(b *testing.B) {
	b.ReportAllocs()

	// currently: num allocations = (2 x num anon fields)) + 2 = 6
	type A2 struct{}

	type A1 struct {
		A2
	}

	type T struct {
		A1
	}

	v := reflect.ValueOf(T{})

	b.ResetTimer()
	for _, cache := range []bool{false, true} {
		b.Run(fmt.Sprintf("cache=%v", cache), func(b *testing.B) {
			fieldCache.Clear()
			enableFieldCache = cache
			for i := 0; i < b.N; i++ {
				if _, err := parseTags(v); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkParseTags_AnonInterfaces(b *testing.B) {
	b.ReportAllocs()

	// currently: num allocations = (2 x num anon fields)) + 2 = 6
	type I1 any
	type I2 any

	type T2 struct {
		I2
	}

	type T struct {
		I1
	}

	v := reflect.ValueOf(T{
		I1: T2{
			I2: T2{},
		},
	})

	b.ResetTimer()
	for _, cache := range []bool{false, true} {
		b.Run(fmt.Sprintf("cache=%v", cache), func(b *testing.B) {
			fieldCache.Clear()
			enableFieldCache = cache
			for i := 0; i < b.N; i++ {
				if _, err := parseTags(v); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkParseTag_Id(b *testing.B) {
	b.ReportAllocs()

	type T struct {
		Id int `jsonapi:"id,type"`
	}

	t := reflect.TypeOf(T{})
	f := t.Field(0)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := parseTag(f, "id", "type"); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkParseTag_Attribute(b *testing.B) {
	b.ReportAllocs()

	type T struct {
		A int `jsonapi:"attr,a"`
	}

	t := reflect.TypeOf(T{})
	f := t.Field(0)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := parseTag(f, "attr", "a"); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkParseTag_ToOneRelationship(b *testing.B) {
	b.ReportAllocs()

	type T struct {
		A int `jsonapi:"rel,type,name"`
	}

	t := reflect.TypeOf(T{})
	f := t.Field(0)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := parseTag(f, "rel", "type,name"); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkParseTag_ToManyRelationship(b *testing.B) {
	b.ReportAllocs()

	type T struct {
		A []int `jsonapi:"rel,type,name"`
	}

	t := reflect.TypeOf(T{})
	f := t.Field(0)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := parseTag(f, "rel", "type,name"); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkParseAttrTag(b *testing.B) {
	b.ReportAllocs()

	type T struct {
		A int `jsonapi:"attr,a"`
	}

	t := reflect.TypeOf(T{})
	f := t.Field(0)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := parseAttrTag(f, "a"); err != nil {
			b.Fatal(err)
		}
	}
}
