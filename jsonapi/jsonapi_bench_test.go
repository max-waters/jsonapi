package jsonapi

import (
	"encoding/json"
	"reflect"
	"testing"
)

func BenchmarkParseTags_Empty(b *testing.B) {
	// currently: 0 allocations
	type T struct{}

	v := reflect.ValueOf(T{})
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if _, err := parseTags(v); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkParseTags_Attrs(b *testing.B) {
	// currently: num allocations = (2 x num attrs) + 2
	type T struct {
		A int `jsonapi:"attr,a"`
		B int `jsonapi:"attr,b"`
		C int `jsonapi:"attr,c"`
		D int `jsonapi:"attr,d"`
		E int `jsonapi:"attr,e"`
	}
	v := reflect.ValueOf(T{})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := parseTags(v); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkParseTags_Anons(b *testing.B) {
	// currently: num allocations = (2 x num anon fields)) + 2
	type A4 struct{}

	type A3 struct {
		A4
	}

	type A2 struct {
		A3
	}

	type A1 struct {
		A2
	}

	type T struct {
		A1
	}

	v := reflect.ValueOf(T{})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := parseTags(v); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkParseTags_Attr(b *testing.B) {
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

func BenchmarkParseAttrTag(b *testing.B) {
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

func BenchmarkMarshalJson_Attr_Primitive(b *testing.B) {
	// 1 alloc
	type T struct {
		A, B, C, D, E int
	}
	t := &T{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := json.Marshal(t); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkUnmarshalJson_Attr_Primitive(b *testing.B) {
	// 4 allocs
	type T struct {
		A, B, C, D, E int
	}
	t := &T{}
	data := []byte(`{"A":1}`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := json.Unmarshal(data, t); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMarshalJson_Anons(b *testing.B) {
	// 1 alloc
	type A4 struct {
		A int
	}

	type A3 struct {
		A4
	}

	type A2 struct {
		A3
	}

	type A1 struct {
		A2
	}

	type T struct {
		A1
	}
	t := &T{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := json.Marshal(t); err != nil {
			b.Fatal(err)
		}
	}
}

func TestMarshalJson_Anons(t *testing.T) {
	// 1 alloc
	type A4 struct {
		A int
	}

	type A3 struct {
		A4
	}

	type A2 struct {
		A3
	}

	type A1 struct {
		A2
	}

	type T struct {
		A1
	}
	in := &T{}

	if _, err := json.Marshal(in); err != nil {
		t.Fatal(err)
	}

}

func BenchmarkUnmarshalJson_Anons(b *testing.B) {
	// 4 alloc no matter the depth
	type A4 struct {
		A int
	}

	type A3 struct {
		A4
	}

	type A2 struct {
		A3
	}

	type A1 struct {
		A2
	}

	type T struct {
		A1
	}
	t := &T{}
	data := []byte(`{"A":1}`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := json.Unmarshal(data, t); err != nil {
			b.Fatal(err)
		}
	}
}
