package jsonapi

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type simpleStruct struct {
	Int int `json:"int,omitempty"`
}

// resource id string type
type rscIdString struct {
	Id string `json:"id" jsonapi:"id,type"`
}

var rscIdStringValue = rscIdString{
	Id: "id",
}

const rscIdStringJson = `
{
	"type": "type",
	"id": "id"
}`

// resource id string ptr type
type rscIdStringPtr struct {
	Id *string `jsonapi:"id,type"`
}

var rscIdStringPtrValue = rscIdStringPtr{
	Id: addrOf("id"),
}

// resource id int type
type rscIdInt struct {
	Id int `jsonapi:"id,type"`
}

var rscIdIntValue = rscIdInt{
	Id: -1,
}

// resource id int ptr type
type rscIdIntPtr struct {
	Id *int `jsonapi:"id,type"`
}

var rscIdIntPtrValue = rscIdIntPtr{
	Id: addrOf(-1),
}

const rscIdIntJson = `
{
	"type": "type",
	"id": -1
}`

// resource id struct type
type rscIdStruct struct {
	Id simpleStruct `jsonapi:"id,type"`
}

var rscIdStructValue = rscIdStruct{
	Id: simpleStruct{
		Int: -2,
	},
}

// resource id struct ptr type
type rscIdStructPtr struct {
	Id *simpleStruct `jsonapi:"id,type"`
}

var rscIdStructPtrValue = rscIdStructPtr{
	Id: &simpleStruct{
		Int: -2,
	},
}

const rscIdStructJson = `
{
	"type": "type",
	"id": {
		"int": -2
	}
}`

func TestMarshalResource_RscId(t *testing.T) {
	type testCase struct {
		In       any
		Expected string
	}

	testCases := []testCase{
		{rscIdStringValue, rscIdStringJson},
		{rscIdStringPtrValue, rscIdStringJson},
		{rscIdIntValue, rscIdIntJson},
		{rscIdIntPtrValue, rscIdIntJson},
		{rscIdStructValue, rscIdStructJson},
		{rscIdStructPtrValue, rscIdStructJson},
		// nil pointers
		{rscIdStructPtr{}, `{"id": null, "type": "type"}`},
		{rscIdStringPtr{}, `{"id": null, "type": "type"}`},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%T", tc.In), func(t *testing.T) {
			got, err := MarshalResource(tc.In)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, fmtJson(t, []byte(tc.Expected)), fmtJson(t, got))
		})
	}
}

func TestUnmarshalResource_RscId(t *testing.T) {
	type testCase struct {
		In       any
		Data     string
		Expected any
	}

	testCases := []testCase{
		{&rscIdString{}, rscIdStringJson, &rscIdStringValue},
		{&rscIdStringPtr{}, rscIdStringJson, &rscIdStringPtrValue},
		{&rscIdInt{}, rscIdIntJson, &rscIdIntValue},
		{&rscIdIntPtr{}, rscIdIntJson, &rscIdIntPtrValue},
		{&rscIdStruct{}, rscIdStructJson, &rscIdStructValue},
		{&rscIdStructPtr{}, rscIdStructJson, &rscIdStructPtrValue},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%T", tc.In), func(t *testing.T) {
			if err := UnmarshalResource([]byte(tc.Data), &tc.In); err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tc.Expected, tc.In)
		})
	}
}

func TestMarshalResource_RscId_OmitEmpty(t *testing.T) {
	type rscIdString struct {
		String string `jsonapi:"id,type,omitempty"`
	}

	type rscIdStringPtr struct {
		String string `jsonapi:"id,type,omitempty"`
	}

	type rscIdStruct struct {
		Struct simpleStruct `jsonapi:"id,type,omitempty"`
	}

	type rscIdStructPtr struct {
		Struct *simpleStruct `jsonapi:"id,type,omitempty"`
	}

	want := `{"type": "type"}`
	testCases := []any{&rscIdString{}, rscIdStringPtr{}, rscIdStruct{}, rscIdStructPtr{}}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%T", tc), func(t *testing.T) {
			got, err := MarshalResource(tc)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, fmtJson(t, []byte(want)), fmtJson(t, got))
		})
	}
}

func TestUnmarshalResource_RscId_EmptyJson(t *testing.T) {
	type testCase struct {
		In       any
		Expected any
	}

	data := "{}"

	testCases := []testCase{
		{rscIdString{}, rscIdString{}},
		{rscIdStringPtr{}, rscIdStringPtr{}},
		{rscIdInt{}, rscIdInt{}},
		{rscIdIntPtr{}, rscIdIntPtr{}},
		{rscIdStruct{}, rscIdStruct{}},
		{rscIdStructPtr{}, rscIdStructPtr{}},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%T", tc.In), func(t *testing.T) {
			if err := UnmarshalResource([]byte(data), &tc.In); err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tc.Expected, tc.In)
		})
	}
}

// attributes of all primitive types
type attrsPrimitive struct {
	Bool      bool    `jsonapi:"attr,bool"`
	Int       int     `jsonapi:"attr,int"`
	Int8      int8    `jsonapi:"attr,int8"`
	Int16     int16   `jsonapi:"attr,int16"`
	Int32     int32   `jsonapi:"attr,int32"`
	Int64     int64   `jsonapi:"attr,int64"`
	Uint      uint    `jsonapi:"attr,uint"`
	Uint8     uint8   `jsonapi:"attr,uint8"`
	Uint16    uint16  `jsonapi:"attr,uint16"`
	Uint32    uint32  `jsonapi:"attr,uint32"`
	Uint64    uint64  `jsonapi:"attr,uint64"`
	Float32   float32 `jsonapi:"attr,float32"`
	Float64   float64 `jsonapi:"attr,float64"`
	String    string  `jsonapi:"attr,string"`
	Rune      rune    `jsonapi:"attr,rune"`
	Byte      byte    `jsonapi:"attr,byte"`
	SliceByte []byte  `jsonapi:"attr,[]byte"`
}

var attrsPrimitiveValue = attrsPrimitive{
	Bool: true,
	Int:  -1, Int8: -2, Int16: -3, Int32: -4, Int64: -5,
	Uint: 6, Uint8: 7, Uint16: 8, Uint32: 9, Uint64: 10,
	Float32: 11.32, Float64: 12.64,
	String: "str-13", Rune: -14, Byte: 15, SliceByte: []byte("str-16"),
}

const attrsPrimitiveJson = `
{
	"attributes": {
		"bool": true,
		"int": -1,
		"int8": -2,
		"int16": -3,
		"int32": -4,
		"int64": -5,
		"uint": 6,
		"uint8": 7,
		"uint16": 8,
		"uint32": 9,
		"uint64": 10,
		"float32": 11.32,
		"float64": 12.64,
		"string": "str-13",
		"rune": -14,
		"byte": 15,
		"[]byte": "c3RyLTE2"
	}
}`

func TestMarshalResource_Attrs_Primitive(t *testing.T) {
	got, err := MarshalResource(attrsPrimitiveValue)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, fmtJson(t, []byte(attrsPrimitiveJson)), fmtJson(t, got))
}

func TestUnmarshalResource_Attrs_Primitive(t *testing.T) {
	got := &attrsPrimitive{}
	if err := UnmarshalResource([]byte(attrsPrimitiveJson), got); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, &attrsPrimitiveValue, got)
}

// attributes of all primitive ptr types
type attrsPrimitivePtr struct {
	Bool      *bool    `jsonapi:"attr,bool"`
	Int       *int     `jsonapi:"attr,int"`
	Int8      *int8    `jsonapi:"attr,int8"`
	Int16     *int16   `jsonapi:"attr,int16"`
	Int32     *int32   `jsonapi:"attr,int32"`
	Int64     *int64   `jsonapi:"attr,int64"`
	Uint      *uint    `jsonapi:"attr,uint"`
	Uint8     *uint8   `jsonapi:"attr,uint8"`
	Uint16    *uint16  `jsonapi:"attr,uint16"`
	Uint32    *uint32  `jsonapi:"attr,uint32"`
	Uint64    *uint64  `jsonapi:"attr,uint64"`
	Float32   *float32 `jsonapi:"attr,float32"`
	Float64   *float64 `jsonapi:"attr,float64"`
	String    *string  `jsonapi:"attr,string"`
	Rune      *rune    `jsonapi:"attr,rune"`
	Byte      *byte    `jsonapi:"attr,byte"`
	SliceByte *[]byte  `jsonapi:"attr,[]byte"`
}

var attrsPrimitivePtrValue = attrsPrimitivePtr{
	Bool: addrOf(true),
	Int:  addrOf(-1), Int8: addrOf(int8(-2)), Int16: addrOf(int16(-3)), Int32: addrOf(int32(-4)), Int64: addrOf(int64(-5)),
	Uint: addrOf(uint(6)), Uint8: addrOf(uint8(7)), Uint16: addrOf(uint16(8)), Uint32: addrOf(uint32(9)), Uint64: addrOf(uint64(10)),
	Float32: addrOf(float32(11.32)), Float64: addrOf(12.64),
	String: addrOf("str-13"), Rune: addrOf(rune(-14)), Byte: addrOf(byte(15)), SliceByte: addrOf([]byte("str-16")),
}

func TestMarshalResource_Attrs_PrimitivePtr(t *testing.T) {
	got, err := MarshalResource(attrsPrimitivePtrValue)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, fmtJson(t, []byte(attrsPrimitiveJson)), fmtJson(t, got))
}

func TestMarshalResource_Attrs_PrimitiveNilPtr(t *testing.T) {
	got, err := MarshalResource(&attrsPrimitivePtr{})
	if err != nil {
		t.Fatal(err)
	}

	want := `
	{
		"attributes": {
			"bool": null,
			"int": null, "int8": null, "int16": null, "int32": null, "int64": null,
			"uint": null, "uint8": null, "uint16": null, "uint32": null, "uint64": null,
			"float32": null, "float64": null,
			"string": null, "rune": null, "byte": null, "[]byte": null
		}
	}`

	assert.Equal(t, fmtJson(t, []byte(want)), fmtJson(t, got))
}

func TestUnmarshalResource_Attrs_PrimitivePtr(t *testing.T) {
	got := &attrsPrimitivePtr{}
	if err := UnmarshalResource([]byte(attrsPrimitiveJson), got); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, &attrsPrimitivePtrValue, got)
}

// attributes of various composite types
type attrsComposite struct {
	ArrBool                 [2]bool                 `jsonapi:"attr,[2]bool"`
	SliceInt                []int                   `jsonapi:"attr,[]int"`
	MapStringUint           map[string]uint         `jsonapi:"attr,map[string]uint"`
	MapStringSliceFloat     map[string][]float64    `jsonapi:"attr,map[string][]float64"`
	SliceMapStringString    []map[string]string     `jsonapi:"attr,[]map[string]string"`
	SliceMapStringSliceByte []map[string][]byte     `jsonapi:"attr,[]map[string][]byte"`
	Struct                  simpleStruct            `jsonapi:"attr,struct"`
	SliceStruct             []simpleStruct          `jsonapi:"attr,[]simpleStruct"`
	MapStringStruct         map[string]simpleStruct `jsonapi:"attr,map[string]simpleStruct"`
}

var attrsCompositeValue = attrsComposite{
	ArrBool:       [2]bool{true, false},
	SliceInt:      []int{-1, -2},
	MapStringUint: map[string]uint{"key3": 4},
	MapStringSliceFloat: map[string][]float64{
		"key5": {6.1, 7.2},
	},
	SliceMapStringString: []map[string]string{
		{"key8": "elem9"},
	},
	SliceMapStringSliceByte: []map[string][]byte{
		{"key10": []byte("elem11")},
	},
	Struct: simpleStruct{
		Int: 12,
	},
	SliceStruct:     []simpleStruct{{Int: 13}},
	MapStringStruct: map[string]simpleStruct{"key14": {Int: 15}},
}

const attrsCompositeJson = `
{
	"attributes": {
		"[2]bool": [ true, false ],
  		"[]int": [ -1, -2 ],
		"map[string]uint": {
			"key3": 4
		},
		"map[string][]float64": {
			"key5": [ 6.1, 7.2 ]
		},
		"[]map[string]string": [
			{
				"key8": "elem9"
			}
		],
		"[]map[string][]byte": [
			{
				"key10": "ZWxlbTEx"
			}
		],
		"struct": {
			"int": 12
		},
		"[]simpleStruct": [
			{
				"int": 13
			}
		],
		"map[string]simpleStruct": {
			"key14": {
				"int": 15
			}
		}
	}
}`

func TestMarshalResource_Attrs_Composite(t *testing.T) {
	got, err := MarshalResource(attrsCompositeValue)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, fmtJson(t, []byte(attrsCompositeJson)), fmtJson(t, got))
}

func TestUnmarshalResource_Attrs_Composite(t *testing.T) {
	got := &attrsComposite{}
	if err := UnmarshalResource([]byte(attrsCompositeJson), got); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, &attrsCompositeValue, got)
}

// attributes of various composite ptr types
type attrsCompositePtr struct {
	ArrBool                 *[2]*bool                 `jsonapi:"attr,[2]bool"`
	SliceInt                *[]*int                   `jsonapi:"attr,[]int"`
	MapStringUint           *map[string]*uint         `jsonapi:"attr,map[string]uint"`
	MapStringSliceFloat     *map[string][]*float64    `jsonapi:"attr,map[string][]float64"`
	SliceMapStringString    *[]map[string]*string     `jsonapi:"attr,[]map[string]string"`
	SliceMapStringSliceByte *[]map[string][]byte      `jsonapi:"attr,[]map[string][]byte"`
	Struct                  *simpleStruct             `jsonapi:"attr,struct"`
	SliceStruct             *[]*simpleStruct          `jsonapi:"attr,[]simpleStruct"`
	MapStringStruct         *map[string]*simpleStruct `jsonapi:"attr,map[string]simpleStruct"`
}

var attrsCompositePtrValue = attrsCompositePtr{
	ArrBool:       addrOf([2]*bool{addrOf(true), addrOf(false)}),
	SliceInt:      addrOf([]*int{addrOf(-1), addrOf(-2)}),
	MapStringUint: addrOf(map[string]*uint{"key3": addrOf(uint(4))}),
	MapStringSliceFloat: addrOf(map[string][]*float64{
		"key5": {addrOf(6.1), addrOf(7.2)},
	}),
	SliceMapStringString: addrOf([]map[string]*string{
		{"key8": addrOf("elem9")},
	}),
	SliceMapStringSliceByte: addrOf([]map[string][]byte{
		{"key10": []byte("elem11")},
	}),
	Struct: addrOf(simpleStruct{
		Int: 12,
	}),
	SliceStruct:     addrOf([]*simpleStruct{{Int: 13}}),
	MapStringStruct: addrOf(map[string]*simpleStruct{"key14": {Int: 15}}),
}

func TestMarshalResource_Attrs_CompositePtr(t *testing.T) {
	got, err := MarshalResource(attrsCompositeValue)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, fmtJson(t, []byte(attrsCompositeJson)), fmtJson(t, got))
}

func TestMarshalResource_Attrs_CompositeNilPtr(t *testing.T) {
	got, err := MarshalResource(attrsCompositePtr{})
	if err != nil {
		t.Fatal(err)
	}

	want := `
	{
		"attributes": {
			"[2]bool": null,
			"[]int": null,
			"map[string]uint": null,
			"map[string][]float64": null,
			"[]map[string]string": null,
			"[]map[string][]byte": null,
			"struct": null,
			"[]simpleStruct": null,
			"map[string]simpleStruct": null
		}
	}`

	assert.Equal(t, fmtJson(t, []byte(want)), fmtJson(t, got))
}

func TestMarshalResource_Attrs_OmitEmpty(t *testing.T) {
	type tp struct {
		String    string        `jsonapi:"attr,string,omitempty"`
		IntPtr    *string       `jsonapi:"attr,int,omitempty"`
		Struct    simpleStruct  `jsonapi:"attr,struct,omitempty"`
		StructPtr *simpleStruct `jsonapi:"attr,structPtr,omitempty"`
	}

	in := &tp{}

	got, err := MarshalResource(in)
	if err != nil {
		t.Fatal(err)
	}

	want := "{}"

	assert.Equal(t, fmtJson(t, []byte(want)), fmtJson(t, got))
}

func TestUnmarshalResource_Attrs_CompositePtr(t *testing.T) {
	got := &attrsCompositePtr{}
	if err := UnmarshalResource([]byte(attrsCompositeJson), got); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, &attrsCompositePtrValue, got)
}

func TestUnmarshalResource_Attrs_EmptyJson(t *testing.T) {
	type testCase struct {
		In       any
		Expected any
	}

	data := "{}"

	testCases := []testCase{
		{attrsPrimitive{}, attrsPrimitive{}},
		{attrsPrimitivePtr{}, attrsPrimitivePtr{}},
		{attrsComposite{}, attrsComposite{}},
		{attrsCompositePtr{}, attrsCompositePtr{}},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%T", tc.In), func(t *testing.T) {
			if err := UnmarshalResource([]byte(data), &tc.In); err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tc.Expected, tc.In)
		})
	}
}

// to-one relations of all primitive types
type relsPrimitive struct {
	Bool      bool    `jsonapi:"rel,bool,rel-bool"`
	Int       int     `jsonapi:"rel,int,rel-int"`
	Int8      int8    `jsonapi:"rel,int8,rel-int8"`
	Int16     int16   `jsonapi:"rel,int16,rel-int16"`
	Int32     int32   `jsonapi:"rel,int32,rel-int32"`
	Int64     int64   `jsonapi:"rel,int64,rel-int64"`
	Uint      uint    `jsonapi:"rel,uint,rel-uint"`
	Uint8     uint    `jsonapi:"rel,uint8,rel-uint8"`
	Uint16    uint16  `jsonapi:"rel,uint16,rel-uint16"`
	Uint32    uint32  `jsonapi:"rel,uint32,rel-uint32"`
	Uint64    uint64  `jsonapi:"rel,uint64,rel-uint64"`
	Float32   float32 `jsonapi:"rel,float32,rel-float32"`
	Float64   float64 `jsonapi:"rel,float64,rel-float64"`
	String    string  `jsonapi:"rel,string,rel-string"`
	Rune      rune    `jsonapi:"rel,rune,rel-rune"`
	Byte      byte    `jsonapi:"rel,byte,rel-byte"`
	SliceByte []byte  `jsonapi:"rel,[]byte,rel-[]byte"`
}

var relsPrimitiveValue = relsPrimitive{
	Bool: true,
	Int:  -1, Int8: -2, Int16: -3, Int32: -4, Int64: -5,
	Uint: 6, Uint8: 7, Uint16: 8, Uint32: 9, Uint64: 10,
	Float32: 11.32, Float64: 12.64,
	String: "str-13", Rune: -14, Byte: 15, SliceByte: []byte("bts-16"),
}

const relsPrimitiveJson = `
{
	"relationships": {
		"bool": {
			"data": { "type": "rel-bool", "id": true }
		},
		"int": {
			"data": { "type": "rel-int", "id": -1 }
		},
		"int8": {
			"data": { "type": "rel-int8", "id": -2 }
		},
		"int16": {
			"data": { "type": "rel-int16", "id": -3 }
		},
		"int32": {
			"data": { "type": "rel-int32", "id": -4 }
		},
		"int64": {
			"data": { "type": "rel-int64", "id": -5 }
		},
		"uint": {
			"data": { "type": "rel-uint", "id": 6 }
		},
		"uint8": {
			"data": { "type": "rel-uint8", "id": 7 }
		},
		"uint16": {
			"data": { "type": "rel-uint16", "id": 8 }
		},
		"uint32": {
			"data": { "type": "rel-uint32", "id": 9 }
		},
		"uint64": {
			"data": { "type": "rel-uint64", "id": 10 }
		},
		"float32": {
			"data": { "type": "rel-float32", "id": 11.32 }
		},
		"float64": {
			"data": { "type": "rel-float64", "id": 12.64 }
		},
		"string": {
			"data": { "type": "rel-string", "id": "str-13" }
		},
		"rune": {
			"data": { "type": "rel-rune", "id": -14 }
		},
		"byte": {
			"data": { "type": "rel-byte", "id": 15 }
		},
		"[]byte": {
			"data": { "type": "rel-[]byte", "id": "YnRzLTE2" }
		}
	}
}`

func TestMarshalResource_ToOneRel_Primitive(t *testing.T) {
	got, err := MarshalResource(&relsPrimitiveValue)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, fmtJson(t, []byte(relsPrimitiveJson)), fmtJson(t, got))
}

func TestUnmarshalResource_ToOneRel_Primitive(t *testing.T) {
	got := &relsPrimitive{}
	if err := UnmarshalResource([]byte(relsPrimitiveJson), got); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, &relsPrimitiveValue, got)
}

// to-one relations of all primitive ptr types
type relsPrimitivePtr struct {
	Bool      *bool    `jsonapi:"rel,bool,rel-bool"`
	Int       *int     `jsonapi:"rel,int,rel-int"`
	Int8      *int8    `jsonapi:"rel,int8,rel-int8"`
	Int16     *int16   `jsonapi:"rel,int16,rel-int16"`
	Int32     *int32   `jsonapi:"rel,int32,rel-int32"`
	Int64     *int64   `jsonapi:"rel,int64,rel-int64"`
	Uint      *uint    `jsonapi:"rel,uint,rel-uint"`
	Uint8     *uint8   `jsonapi:"rel,uint8,rel-uint8"`
	Uint16    *uint16  `jsonapi:"rel,uint16,rel-uint16"`
	Uint32    *uint32  `jsonapi:"rel,uint32,rel-uint32"`
	Uint64    *uint64  `jsonapi:"rel,uint64,rel-uint64"`
	Float32   *float32 `jsonapi:"rel,float32,rel-float32"`
	Float64   *float64 `jsonapi:"rel,float64,rel-float64"`
	String    *string  `jsonapi:"rel,string,rel-string"`
	Rune      *rune    `jsonapi:"rel,rune,rel-rune"`
	Byte      *byte    `jsonapi:"rel,byte,rel-byte"`
	SliceByte *[]byte  `jsonapi:"rel,[]byte,rel-[]byte"`
}

var relsPrimitivePtrValue = relsPrimitivePtr{
	Bool: addrOf(true),
	Int:  addrOf(-1), Int8: addrOf(int8(-2)), Int16: addrOf(int16(-3)), Int32: addrOf(int32(-4)), Int64: addrOf(int64(-5)),
	Uint: addrOf(uint(6)), Uint8: addrOf(uint8(7)), Uint16: addrOf(uint16(8)), Uint32: addrOf(uint32(9)), Uint64: addrOf(uint64(10)),
	Float32: addrOf(float32(11.32)), Float64: addrOf(12.64),
	String: addrOf("str-13"), Rune: addrOf(rune(-14)), Byte: addrOf(byte(15)), SliceByte: addrOf([]byte("bts-16")),
}

func TestMarshalResource_ToOneRel_PrimitivePtr(t *testing.T) {
	got, err := MarshalResource(&relsPrimitivePtrValue)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, fmtJson(t, []byte(relsPrimitiveJson)), fmtJson(t, got))
}

func TestUnmarshalResource_ToOneRels_PrimitivePtr(t *testing.T) {
	got := &relsPrimitivePtr{}
	if err := UnmarshalResource([]byte(relsPrimitiveJson), got); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, &relsPrimitivePtrValue, got)
}

func TestUnmarshalResource_ToOneRels_EmptyJson(t *testing.T) {
	type testCase struct {
		In       any
		Expected any
	}

	data := "{}"

	testCases := []testCase{
		{relsPrimitive{}, relsPrimitive{}},
		{relsPrimitivePtr{}, relsPrimitivePtr{}},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%T", tc.In), func(t *testing.T) {
			if err := UnmarshalResource([]byte(data), &tc.In); err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tc.Expected, tc.In)
		})
	}
}

// to-one relations of all composite types. NB slices,
// arrays and maps are treated as to-many relationships
type relsComposite struct {
	Struct simpleStruct `jsonapi:"rel,rel-struct,struct"`
}

var relsCompositeValue = relsComposite{
	Struct: simpleStruct{
		Int: -1,
	},
}

const relsCompositeJson = `
{
	"relationships": {
		"rel-struct": {
			"data": { "type": "struct", "id": { "int": -1 } }
		}
	}
}`

func TestMarshalResource_ToOneRel_Composite(t *testing.T) {
	got, err := MarshalResource(&relsCompositeValue)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, fmtJson(t, []byte(relsCompositeJson)), fmtJson(t, got))
}

func TestUnmarshalResource_ToOneRel_Composite(t *testing.T) {
	got := relsComposite{}
	if err := UnmarshalResource([]byte(relsCompositeJson), &got); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, relsCompositeValue, got)
}

// to-one relations of all composite ptr types
type relsCompositePtr struct {
	Struct *simpleStruct `jsonapi:"rel,rel-struct,struct"`
}

var relsCompositePtrValue = relsCompositePtr{
	Struct: &simpleStruct{
		Int: -1,
	},
}

func TestMarshalResource_ToOneRel_CompositePtr(t *testing.T) {
	got, err := MarshalResource(&relsCompositePtrValue)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, fmtJson(t, []byte(relsCompositeJson)), fmtJson(t, got))
}

func TestMarshalResource_ToOneRel_OmitEmpty(t *testing.T) {
	type tp struct {
		String    string        `jsonapi:"rel,string,rel-string,omitempty"`
		IntPtr    *string       `jsonapi:"rel,int-ptr,rel-int-ptr,omitempty"`
		Struct    simpleStruct  `jsonapi:"rel,struct,rel-struct,omitempty"`
		StructPtr *simpleStruct `jsonapi:"rel,struct-ptr,rel-struct-ptr,omitempty"`
	}

	in := &tp{}

	got, err := MarshalResource(in)
	if err != nil {
		t.Fatal(err)
	}

	want := "{}"
	assert.Equal(t, fmtJson(t, []byte(want)), fmtJson(t, got))
}

func TestUnmarshalResource_ToOneRel_CompositePtr(t *testing.T) {
	got := relsCompositePtr{}
	if err := UnmarshalResource([]byte(relsCompositeJson), &got); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, relsCompositePtrValue, got)
}

// to-many relations of all primitive types
type relsToManyPrimitive struct {
	Bool    []bool    `jsonapi:"rel,bool,rel-bool"`
	Int     []int     `jsonapi:"rel,int,rel-int"`
	Int16   []int16   `jsonapi:"rel,int16,rel-int16"`
	Int32   []int32   `jsonapi:"rel,int32,rel-int32"`
	Int64   []int64   `jsonapi:"rel,int64,rel-int64"`
	Uint    []uint    `jsonapi:"rel,uint,rel-uint"`
	Uint16  []uint16  `jsonapi:"rel,uint16,rel-uint16"`
	Uint32  []uint32  `jsonapi:"rel,uint32,rel-uint32"`
	Uint64  []uint64  `jsonapi:"rel,uint64,rel-uint64"`
	Float32 []float32 `jsonapi:"rel,float32,rel-float32"`
	Float64 []float64 `jsonapi:"rel,float64,rel-float64"`
	String  []string  `jsonapi:"rel,string,rel-string"`
	Rune    []rune    `jsonapi:"rel,rune,rel-rune"`
	Byte    [][]byte  `jsonapi:"rel,byte,rel-byte"`
}

var relsToManyPrimitiveValue = relsToManyPrimitive{
	Bool: []bool{true, false},
	Int:  []int{-1, -2}, Int16: []int16{-3, -4}, Int32: []int32{-5, -6}, Int64: []int64{-7, -8},
	Uint: []uint{9, 10}, Uint16: []uint16{11, 12}, Uint32: []uint32{13, 14}, Uint64: []uint64{15, 16},
	Float32: []float32{17.32, 18.01}, Float64: []float64{19.64, 20.65},
	String: []string{"str-21", "str-22"}, Rune: []rune{-23, -24}, Byte: [][]byte{[]byte("bts-25"), []byte("bts-26")},
}

const relsToManyPrimitiveJson = `
{
	"relationships": {
		"bool": {
			"data": [ { "type": "rel-bool", "id": true }, { "type": "rel-bool", "id": false } ]
		},
		"int": {
			"data": [ { "type": "rel-int", "id": -1 }, { "type": "rel-int", "id": -2 } ]
		},
		"int16": {
			"data": [ { "type": "rel-int16", "id": -3 }, { "type": "rel-int16", "id": -4 } ]
		},
		"int32": {
			"data": [ { "type": "rel-int32", "id": -5 }, { "type": "rel-int32", "id": -6 } ]
		},
		"int64": {
			"data": [ { "type": "rel-int64", "id": -7 }, { "type": "rel-int64", "id": -8 } ]
		},
		"uint": {
			"data": [ { "type": "rel-uint", "id": 9 }, { "type": "rel-uint", "id": 10 } ]
		},
		"uint": {
			"data": [ { "type": "rel-uint", "id": 9 }, { "type": "rel-uint", "id": 10 } ]
		},
		"uint16": {
			"data": [ { "type": "rel-uint16", "id": 11 }, { "type": "rel-uint16", "id": 12 } ]
		},
		"uint32": {
			"data": [ { "type": "rel-uint32", "id": 13 }, { "type": "rel-uint32", "id": 14 } ]
		},
		"uint64": {
			"data": [ { "type": "rel-uint64", "id": 15 }, { "type": "rel-uint64", "id": 16 } ]
		},
		"float32": {
			"data": [ { "type": "rel-float32", "id": 17.32 }, { "type": "rel-float32", "id": 18.01 } ]
		},
		"float64": {
			"data": [ { "type": "rel-float64", "id": 19.64 }, { "type": "rel-float64", "id": 20.65 } ]
		},
		"string": {
			"data": [ { "type": "rel-string", "id": "str-21" }, { "type": "rel-string", "id": "str-22" } ]
		},
		"rune": {
			"data": [ { "type": "rel-rune", "id": -23 }, { "type": "rel-rune", "id": -24 } ]
		},
		"byte": {
			"data": [ { "type": "rel-byte", "id": "YnRzLTI1" }, { "type": "rel-byte", "id": "YnRzLTI2" } ]
		}
	}
}`

func TestMarshalResource_ToManyRel_Primitive(t *testing.T) {
	got, err := MarshalResource(&relsToManyPrimitiveValue)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, fmtJson(t, []byte(relsToManyPrimitiveJson)), fmtJson(t, got))
}

func TestUnmarshalResource_ToManyRels_Primitive(t *testing.T) {
	got := &relsToManyPrimitive{}
	if err := UnmarshalResource([]byte(relsToManyPrimitiveJson), got); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, &relsToManyPrimitiveValue, got)
}

// to-many relations of all primitive types
type relsToManyPrimitivePtr struct {
	Bool    []*bool    `jsonapi:"rel,bool,rel-bool"`
	Int     []*int     `jsonapi:"rel,int,rel-int"`
	Int16   []*int16   `jsonapi:"rel,int16,rel-int16"`
	Int32   []*int32   `jsonapi:"rel,int32,rel-int32"`
	Int64   []*int64   `jsonapi:"rel,int64,rel-int64"`
	Uint    []*uint    `jsonapi:"rel,uint,rel-uint"`
	Uint16  []*uint16  `jsonapi:"rel,uint16,rel-uint16"`
	Uint32  []*uint32  `jsonapi:"rel,uint32,rel-uint32"`
	Uint64  []*uint64  `jsonapi:"rel,uint64,rel-uint64"`
	Float32 []*float32 `jsonapi:"rel,float32,rel-float32"`
	Float64 []*float64 `jsonapi:"rel,float64,rel-float64"`
	String  []*string  `jsonapi:"rel,string,rel-string"`
	Rune    []*rune    `jsonapi:"rel,rune,rel-rune"`
	Byte    []*[]byte  `jsonapi:"rel,byte,rel-byte"`
}

var relsToManyPrimitivePtrValue = relsToManyPrimitivePtr{
	Bool: []*bool{addrOf(true), addrOf(false)},
	Int:  []*int{addrOf(-1), addrOf(-2)}, Int16: []*int16{addrOf(int16(-3)), addrOf(int16(-4))},
	Int32: []*int32{addrOf(int32(-5)), addrOf(int32(-6))}, Int64: []*int64{addrOf(int64(-7)), addrOf(int64(-8))},
	Uint: []*uint{addrOf(uint(9)), addrOf(uint(10))}, Uint16: []*uint16{addrOf(uint16(11)), addrOf(uint16(12))},
	Uint32: []*uint32{addrOf(uint32(13)), addrOf(uint32(14))}, Uint64: []*uint64{addrOf(uint64(15)), addrOf(uint64(16))},
	Float32: []*float32{addrOf(float32(17.32)), addrOf(float32(18.01))}, Float64: []*float64{addrOf(float64(19.64)), addrOf(float64(20.65))},
	String: []*string{addrOf("str-21"), addrOf("str-22")}, Rune: []*rune{addrOf(rune(-23)), addrOf(rune(-24))},
	Byte: []*[]byte{addrOf([]byte("bts-25")), addrOf([]byte("bts-26"))},
}

func TestMarshalResource_ToManyRel_PrimitivePtr(t *testing.T) {
	got, err := MarshalResource(&relsToManyPrimitivePtrValue)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, fmtJson(t, []byte(relsToManyPrimitiveJson)), fmtJson(t, got))
}

func TestUnmarshalResource_ToManyRels_PrimitivePtr(t *testing.T) {
	got := &relsToManyPrimitivePtr{}
	if err := UnmarshalResource([]byte(relsToManyPrimitiveJson), got); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, &relsToManyPrimitivePtrValue, got)
}

type relsToManyComposite struct {
	ArrStruct  []simpleStruct `jsonapi:"rel,arr-struct,struct"`
	ArrByteArr [][]byte       `jsonapi:"rel,arr-byte-arr,byte-arr"`
}

var relsToManyCompositeValue = relsToManyComposite{
	ArrStruct:  []simpleStruct{{Int: -1}, {Int: 2}},
	ArrByteArr: [][]byte{[]byte("3"), []byte("4")},
}

const relsToManyCompositeJson = `{
	"relationships": {
		"arr-struct": {
			"data": [ { "type": "struct", "id": { "int": -1  } }, { "type": "struct", "id": { "int": 2 } } ]
		},
		"arr-byte-arr": {
			"data": [ { "type": "byte-arr", "id": "Mw==" }, { "type": "byte-arr", "id": "NA==" } ]
		}
	}

}`

func TestMarshalResource_ToManyRel_Composite(t *testing.T) {
	got, err := MarshalResource(&relsToManyCompositeValue)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, fmtJson(t, []byte(relsToManyCompositeJson)), fmtJson(t, got))
}

func TestUnmarshalResource_ToManyRels_Composite(t *testing.T) {
	got := &relsToManyComposite{}
	if err := UnmarshalResource([]byte(relsToManyCompositeJson), got); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, &relsToManyCompositeValue, got)
}

type relsToManyCompositePtr struct {
	ArrStruct  []*simpleStruct `jsonapi:"rel,arr-struct,struct"`
	ArrByteArr []*[]byte       `jsonapi:"rel,arr-byte-arr,byte-arr"`
}

var relsToManyCompositePtrValue = relsToManyCompositePtr{
	ArrStruct:  []*simpleStruct{{Int: -1}, {Int: 2}},
	ArrByteArr: []*[]byte{addrOf([]byte("3")), addrOf([]byte("4"))},
}

func TestMarshalResource_ToManyRel_CompositePtr(t *testing.T) {
	got, err := MarshalResource(&relsToManyCompositePtrValue)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, fmtJson(t, []byte(relsToManyCompositeJson)), fmtJson(t, got))
}

func TestUnmarshalResource_ToManyRels_CompositePtr(t *testing.T) {
	got := &relsToManyCompositePtr{}
	if err := UnmarshalResource([]byte(relsToManyCompositeJson), got); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, &relsToManyCompositePtrValue, got)
}

func TestMarshalResource_ToManyRel_OmitEmpty(t *testing.T) {
	type tp struct {
		SliceString      []string        `jsonapi:"rel,[]string,rel-string,omitempty"`
		SliceIntPtrSlice []*string       `jsonapi:"rel,[]int-ptr,rel-int-ptr,omitempty"`
		SliceStruct      []simpleStruct  `jsonapi:"rel,[]struct,rel-struct,omitempty"`
		SliceStructPtr   []*simpleStruct `jsonapi:"rel,[]struct-ptr,rel-struct-ptr,omitempty"`
		PtrSliceString   *[]string       `jsonapi:"rel,ptr-[]string,rel-string,omitempty"`
	}

	in := &tp{
		PtrSliceString: addrOf([]string{}),
	}

	got, err := MarshalResource(in)
	if err != nil {
		t.Fatal(err)
	}

	want := "{}"
	assert.Equal(t, fmtJson(t, []byte(want)), fmtJson(t, got))
}

func TestUnmarshalResource_ToManyRels_EmptyJson(t *testing.T) {
	type testCase struct {
		In       any
		Expected any
	}

	data := "{}"

	testCases := []testCase{
		{relsToManyPrimitive{}, relsToManyPrimitive{}},
		{relsToManyPrimitivePtr{}, relsToManyPrimitivePtr{}},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%T", tc.In), func(t *testing.T) {
			if err := UnmarshalResource([]byte(data), &tc.In); err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tc.Expected, tc.In)
		})
	}
}

// meta of all primitive types
type metaPrimitive struct {
	Bool      bool    `jsonapi:"meta,bool"`
	Int       int     `jsonapi:"meta,int"`
	Int8      int8    `jsonapi:"meta,int8"`
	Int16     int16   `jsonapi:"meta,int16"`
	Int32     int32   `jsonapi:"meta,int32"`
	Int64     int64   `jsonapi:"meta,int64"`
	Uint      uint    `jsonapi:"meta,uint"`
	Uint8     uint8   `jsonapi:"meta,uint8"`
	Uint16    uint16  `jsonapi:"meta,uint16"`
	Uint32    uint32  `jsonapi:"meta,uint32"`
	Uint64    uint64  `jsonapi:"meta,uint64"`
	Float32   float32 `jsonapi:"meta,float32"`
	Float64   float64 `jsonapi:"meta,float64"`
	String    string  `jsonapi:"meta,string"`
	Rune      rune    `jsonapi:"meta,rune"`
	Byte      byte    `jsonapi:"meta,byte"`
	SliceByte []byte  `jsonapi:"meta,[]byte"`
}

var metaPrimitiveValue = metaPrimitive{
	Bool: true,
	Int:  -1, Int8: -2, Int16: -3, Int32: -4, Int64: -5,
	Uint: 6, Uint8: 7, Uint16: 8, Uint32: 9, Uint64: 10,
	Float32: 11.32, Float64: 12.64,
	String: "str-13", Rune: -14, Byte: 15, SliceByte: []byte("str-16"),
}

const metaPrimitiveJson = `
{
	"meta": {
		"bool": true,
		"int": -1,
		"int8": -2,
		"int16": -3,
		"int32": -4,
		"int64": -5,
		"uint": 6,
		"uint8": 7,
		"uint16": 8,
		"uint32": 9,
		"uint64": 10,
		"float32": 11.32,
		"float64": 12.64,
		"string": "str-13",
		"rune": -14,
		"byte": 15,
		"[]byte": "c3RyLTE2"
	}
}`

func TestMarshalResource_Meta_Primitive(t *testing.T) {
	got, err := MarshalResource(metaPrimitiveValue)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, fmtJson(t, []byte(metaPrimitiveJson)), fmtJson(t, got))
}

func TestUnmarshalResource_Meta_Primitive(t *testing.T) {
	got := &metaPrimitive{}
	if err := UnmarshalResource([]byte(metaPrimitiveJson), got); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, &metaPrimitiveValue, got)
}

// meta of all primitive ptr types
type metaPrimitivePtr struct {
	Bool      *bool    `jsonapi:"meta,bool"`
	Int       *int     `jsonapi:"meta,int"`
	Int8      *int8    `jsonapi:"meta,int8"`
	Int16     *int16   `jsonapi:"meta,int16"`
	Int32     *int32   `jsonapi:"meta,int32"`
	Int64     *int64   `jsonapi:"meta,int64"`
	Uint      *uint    `jsonapi:"meta,uint"`
	Uint8     *uint8   `jsonapi:"meta,uint8"`
	Uint16    *uint16  `jsonapi:"meta,uint16"`
	Uint32    *uint32  `jsonapi:"meta,uint32"`
	Uint64    *uint64  `jsonapi:"meta,uint64"`
	Float32   *float32 `jsonapi:"meta,float32"`
	Float64   *float64 `jsonapi:"meta,float64"`
	String    *string  `jsonapi:"meta,string"`
	Rune      *rune    `jsonapi:"meta,rune"`
	Byte      *byte    `jsonapi:"meta,byte"`
	SliceByte *[]byte  `jsonapi:"meta,[]byte"`
}

var metaPrimitivePtrValue = metaPrimitivePtr{
	Bool: addrOf(true),
	Int:  addrOf(-1), Int8: addrOf(int8(-2)), Int16: addrOf(int16(-3)), Int32: addrOf(int32(-4)), Int64: addrOf(int64(-5)),
	Uint: addrOf(uint(6)), Uint8: addrOf(uint8(7)), Uint16: addrOf(uint16(8)), Uint32: addrOf(uint32(9)), Uint64: addrOf(uint64(10)),
	Float32: addrOf(float32(11.32)), Float64: addrOf(12.64),
	String: addrOf("str-13"), Rune: addrOf(rune(-14)), Byte: addrOf(byte(15)), SliceByte: addrOf([]byte("str-16")),
}

func TestMarshalResource_Meta_PrimitivePtr(t *testing.T) {
	got, err := MarshalResource(metaPrimitivePtrValue)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, fmtJson(t, []byte(metaPrimitiveJson)), fmtJson(t, got))
}

func TestUnmarshalResource_Meta_PrimitivePtr(t *testing.T) {
	got := &metaPrimitivePtr{}
	if err := UnmarshalResource([]byte(metaPrimitiveJson), got); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, &metaPrimitivePtrValue, got)
}

func TestMarshalResource_Meta_PrimitiveNilPtr(t *testing.T) {
	got, err := MarshalResource(&metaPrimitivePtr{})
	if err != nil {
		t.Fatal(err)
	}

	want := `
	{
		"meta": {
			"bool": null,
			"int": null, "int8": null, "int16": null, "int32": null, "int64": null,
			"uint": null, "uint8": null, "uint16": null, "uint32": null, "uint64": null,
			"float32": null, "float64": null,
			"string": null, "rune": null, "byte": null, "[]byte": null
		}
	}`

	assert.Equal(t, fmtJson(t, []byte(want)), fmtJson(t, got))
}

// meta of various composite types
type metaComposite struct {
	ArrBool                 [2]bool                 `jsonapi:"meta,[2]bool"`
	SliceInt                []int                   `jsonapi:"meta,[]int"`
	MapStringUint           map[string]uint         `jsonapi:"meta,map[string]uint"`
	MapStringSliceFloat     map[string][]float64    `jsonapi:"meta,map[string][]float64"`
	SliceMapStringString    []map[string]string     `jsonapi:"meta,[]map[string]string"`
	SliceMapStringSliceByte []map[string][]byte     `jsonapi:"meta,[]map[string][]byte"`
	Struct                  simpleStruct            `jsonapi:"meta,struct"`
	SliceStruct             []simpleStruct          `jsonapi:"meta,[]simpleStruct"`
	MapStringStruct         map[string]simpleStruct `jsonapi:"meta,map[string]simpleStruct"`
}

var metaCompositeValue = metaComposite{
	ArrBool:       [2]bool{true, false},
	SliceInt:      []int{-1, -2},
	MapStringUint: map[string]uint{"key3": 4},
	MapStringSliceFloat: map[string][]float64{
		"key5": {6.1, 7.2},
	},
	SliceMapStringString: []map[string]string{
		{"key8": "elem9"},
	},
	SliceMapStringSliceByte: []map[string][]byte{
		{"key10": []byte("elem11")},
	},
	Struct: simpleStruct{
		Int: 12,
	},
	SliceStruct:     []simpleStruct{{Int: 13}},
	MapStringStruct: map[string]simpleStruct{"key14": {Int: 15}},
}

const metaCompositeJson = `
{
	"meta": {
		"[2]bool": [ true, false ],
  		"[]int": [ -1, -2 ],
		"map[string]uint": {
			"key3": 4
		},
		"map[string][]float64": {
			"key5": [ 6.1, 7.2 ]
		},
		"[]map[string]string": [
			{
				"key8": "elem9"
			}
		],
		"[]map[string][]byte": [
			{
				"key10": "ZWxlbTEx"
			}
		],
		"struct": {
			"int": 12
		},
		"[]simpleStruct": [
			{
				"int": 13
			}
		],
		"map[string]simpleStruct": {
			"key14": {
				"int": 15
			}
		}
	}
}`

func TestMarshalResource_Meta_Composite(t *testing.T) {
	got, err := MarshalResource(metaCompositeValue)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, fmtJson(t, []byte(metaCompositeJson)), fmtJson(t, got))
}

func TestUnmarshalResource_Meta_Composite(t *testing.T) {
	got := &metaComposite{}
	if err := UnmarshalResource([]byte(metaCompositeJson), got); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, &metaCompositeValue, got)
}

// meta of various composite ptr types
type metaCompositePtr struct {
	ArrBool                 *[2]*bool                 `jsonapi:"meta,[2]bool"`
	SliceInt                *[]*int                   `jsonapi:"meta,[]int"`
	MapStringUint           *map[string]*uint         `jsonapi:"meta,map[string]uint"`
	MapStringSliceFloat     *map[string][]*float64    `jsonapi:"meta,map[string][]float64"`
	SliceMapStringString    *[]map[string]*string     `jsonapi:"meta,[]map[string]string"`
	SliceMapStringSliceByte *[]map[string][]byte      `jsonapi:"meta,[]map[string][]byte"`
	Struct                  *simpleStruct             `jsonapi:"meta,struct"`
	SliceStruct             *[]*simpleStruct          `jsonapi:"meta,[]simpleStruct"`
	MapStringStruct         *map[string]*simpleStruct `jsonapi:"meta,map[string]simpleStruct"`
}

var metaCompositePtrValue = metaCompositePtr{
	ArrBool:       addrOf([2]*bool{addrOf(true), addrOf(false)}),
	SliceInt:      addrOf([]*int{addrOf(-1), addrOf(-2)}),
	MapStringUint: addrOf(map[string]*uint{"key3": addrOf(uint(4))}),
	MapStringSliceFloat: addrOf(map[string][]*float64{
		"key5": {addrOf(6.1), addrOf(7.2)},
	}),
	SliceMapStringString: addrOf([]map[string]*string{
		{"key8": addrOf("elem9")},
	}),
	SliceMapStringSliceByte: addrOf([]map[string][]byte{
		{"key10": []byte("elem11")},
	}),
	Struct: addrOf(simpleStruct{
		Int: 12,
	}),
	SliceStruct:     addrOf([]*simpleStruct{{Int: 13}}),
	MapStringStruct: addrOf(map[string]*simpleStruct{"key14": {Int: 15}}),
}

func TestMarshalResource_Meta_CompositePtr(t *testing.T) {
	got, err := MarshalResource(metaCompositePtrValue)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, fmtJson(t, []byte(metaCompositeJson)), fmtJson(t, got))
}

func TestMarshalResource_Meta_CompositeNilPtr(t *testing.T) {
	got, err := MarshalResource(metaCompositePtr{})
	if err != nil {
		t.Fatal(err)
	}

	want := `
	{
		"meta": {
			"[2]bool": null,
			"[]int": null,
			"map[string]uint": null,
			"map[string][]float64": null,
			"[]map[string]string": null,
			"[]map[string][]byte": null,
			"struct": null,
			"[]simpleStruct": null,
			"map[string]simpleStruct": null
		}
	}`

	assert.Equal(t, fmtJson(t, []byte(want)), fmtJson(t, got))
}

func TestUnmarshalResource_Meta_CompositePtr(t *testing.T) {
	got := &metaCompositePtr{}
	if err := UnmarshalResource([]byte(metaCompositeJson), got); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, &metaCompositePtrValue, got)
}

func TestUnmarshalResource_Meta_EmptyJson(t *testing.T) {
	type testCase struct {
		In       any
		Expected any
	}

	data := "{}"

	testCases := []testCase{
		{metaPrimitive{}, metaPrimitive{}},
		{metaPrimitivePtr{}, metaPrimitivePtr{}},
		{metaComposite{}, metaComposite{}},
		{metaCompositePtr{}, metaCompositePtr{}},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%T", tc.In), func(t *testing.T) {
			if err := UnmarshalResource([]byte(data), &tc.In); err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tc.Expected, tc.In)
		})
	}
}

func TestMarshalResource_Meta_OmitEmpty(t *testing.T) {
	type tp struct {
		String    string        `jsonapi:"meta,string,omitempty"`
		IntPtr    *string       `jsonapi:"meta,int,omitempty"`
		Struct    simpleStruct  `jsonapi:"meta,struct,omitempty"`
		StructPtr *simpleStruct `jsonapi:"meta,structPtr,omitempty"`
	}

	in := &tp{}

	got, err := MarshalResource(in)
	if err != nil {
		t.Fatal(err)
	}

	want := "{}"

	assert.Equal(t, fmtJson(t, []byte(want)), fmtJson(t, got))
}

type noJsonKey struct {
	A1 int `jsonapi:"attr"`
	A2 int `jsonapi:"attr,,omitempty"`
	A3 int
	R1 int   `jsonapi:"rel,,rel-type"`
	R2 []int `jsonapi:"rel,,rel-type"`
	M1 int   `jsonapi:"meta"`
	M2 int   `jsonapi:"meta,,omitempty"`
}

var noJsonKeyValue = noJsonKey{
	A1: 1, A2: 2, A3: 3, R1: 4, R2: []int{5, 6}, M1: 7, M2: 8,
}

const noJsonKeyJson = `
{
	"attributes": {
		"A1": 1,
		"A2": 2,
		"A3": 3
	},
	"relationships": {
		"R1": {
			"data": {
				"type": "rel-type",
				"id": 4
			}
		},
		"R2": {	
			"data": [ { "type": "rel-type", "id": 5 }, { "type": "rel-type", "id": 6 } ]
		}
	},
	"meta": {
		"M1": 7,
		"M2": 8
	}
}
`

func TestMarshalResource_NoJsonKey(t *testing.T) {
	got, err := MarshalResource(noJsonKeyValue)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, fmtJson(t, []byte(noJsonKeyJson)), fmtJson(t, got))
}

func TestUnmarshalResource_NoJsonKey(t *testing.T) {
	got := noJsonKey{}
	err := UnmarshalResource([]byte(noJsonKeyJson), &got)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, noJsonKeyValue, got)
}

type Anonymous2 struct {
	Id  string `json:"id" jsonapi:"id,embed"`
	Int int    `json:"int" jsonapi:"attr,int"`
}

type Anonymous1 struct {
	Anonymous2
	String string `json:"string" jsonapi:"attr,string"`
}

type anonymous struct {
	Anonymous1
	Float64 float64 `json:"float64" jsonapi:"attr,float64"`
}

var anonymousValue = anonymous{
	Anonymous1: Anonymous1{
		Anonymous2: Anonymous2{
			Id:  "1",
			Int: 2,
		},
		String: "3",
	},
	Float64: 4.1,
}

const anonymousJson = `
{
	"type": "embed",
	"id": "1",
	"attributes": {
		"int": 2,
		"string": "3",
		"float64": 4.1
	}
}`

func TestMarshalResource_Anonymous(t *testing.T) {
	got, err := MarshalResource(anonymousValue)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, fmtJson(t, []byte(anonymousJson)), fmtJson(t, got))
}

func TestUnmarshalResource_Anonymous(t *testing.T) {
	got := anonymous{}
	if err := UnmarshalResource([]byte(anonymousJson), &got); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, anonymousValue, got)
}

type Anonymous1Ptr struct {
	*Anonymous2
	String string `json:"string" jsonapi:"attr,string"`
}

type anonymousPtr struct {
	*Anonymous1Ptr
	Float64 float64 `json:"float64" jsonapi:"attr,float64"`
}

var anonymousPtrValue = anonymousPtr{
	Anonymous1Ptr: &Anonymous1Ptr{
		Anonymous2: &Anonymous2{
			Id:  "1",
			Int: 2,
		},
		String: "3",
	},
	Float64: 4.1,
}

func TestMarshalResource_AnonymousPtr(t *testing.T) {
	got, err := MarshalResource(anonymousPtrValue)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, fmtJson(t, []byte(anonymousJson)), fmtJson(t, got))
}

func TestUnmarshalResource_AnonymousPtr(t *testing.T) {
	got := anonymousPtr{}
	if err := UnmarshalResource([]byte(anonymousJson), &got); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, anonymousPtrValue, got)
}

type AnonymousOverride1 struct {
	// first two fields are overridden
	String  string  `json:"string" jsonapi:"attr,a"`
	Int     int     `json:"int" jsonapi:"id,anon"`
	Float64 float64 `json:"float64" jsonapi:"attr,float64"`
}

type anonymousOverride struct {
	AnonymousOverride1
	String string `json:"string" jsonapi:"id,anon"`
	Int    int    `json:"int" jsonapi:"attr,a"`
}

var anonymousOverrideMarshalValue = anonymousOverride{
	AnonymousOverride1: AnonymousOverride1{
		String:  "1",
		Int:     2,
		Float64: 3.1,
	},
	String: "4",
	Int:    5,
}

var anonymousOverrideUnmarshalValue = anonymousOverride{
	AnonymousOverride1: AnonymousOverride1{
		String:  "",
		Int:     0,
		Float64: 3.1,
	},
	String: "4",
	Int:    5,
}

const anonymousOverrideJson = `
{
	"type": "anon",
	"id": "4",
	"attributes": {
		"float64": 3.1,
		"a": 5
	}
}`

func TestMarshalResource_AnonymousOverride(t *testing.T) {
	got, err := MarshalResource(anonymousOverrideMarshalValue)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, fmtJson(t, []byte(anonymousOverrideJson)), fmtJson(t, got))
}

func TestUnmarshalResource_AnonymousOverride(t *testing.T) {
	got := anonymousOverride{}
	if err := UnmarshalResource([]byte(anonymousOverrideJson), &got); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, anonymousOverrideUnmarshalValue, got)
}

type anonWithExp struct {
	Int int `jsonapi:"attr,int"`
}

type anonUnexpExp struct {
	anonWithExp
}

var anonUnexpExpValue = anonUnexpExp{
	anonWithExp: anonWithExp{
		Int: 1,
	},
}

const anonWithExpJson = `
{
	"attributes": {
		"int": 1
	}
}`

func TestMarshalResource_AnonymousUnexportedExportedField(t *testing.T) {
	got, err := MarshalResource(anonUnexpExpValue)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, fmtJson(t, []byte(anonWithExpJson)), fmtJson(t, got))
}

func TestUnmarshalResource_AnonymousUnexportedExportedField(t *testing.T) {
	got := anonUnexpExp{}
	if err := UnmarshalResource([]byte(anonWithExpJson), &got); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, anonUnexpExpValue, got)
}

type PromotionPrecedence1 struct {
	Id      string  `jsonapi:"id,type"`
	Int     int     `jsonapi:"attr,int"`
	String  string  `jsonapi:"attr"`
	Float32 float32 `json:"flt32"`
	Float64 float64
	Uint    uint   `jsonapi:"attr,uint"`
	Uint16  uint16 `json:"Uint16"`
}

type PromotionPrecedence2 struct {
	Id      int     `jsonapi:"id,type"`
	Int     int     `jsonapi:"attr,int"`
	String  string  `jsonapi:"attr"`
	Float32 float32 `json:"flt32"`
	Float64 float64
	Uint    uint `json:"uint"`
	Uint16  uint16
}

type promotionPrecedence struct {
	PromotionPrecedence1
	PromotionPrecedence2 //nolint:govet // ignore the lint error, as we're testing this
}

var promotionPrecedenceMarshalValue = promotionPrecedence{
	PromotionPrecedence1: PromotionPrecedence1{
		Id:      "id",
		Int:     1,
		String:  "2",
		Float32: 3.1,
		Float64: 4.2,
		Uint:    5,
		Uint16:  6,
	},
	PromotionPrecedence2: PromotionPrecedence2{
		Id:      7,
		Int:     8,
		String:  "9",
		Float32: 10.1,
		Float64: 11.2,
		Uint:    12,
		Uint16:  13,
	},
}

var promotionPrecedenceUnmarshalValue = promotionPrecedence{
	PromotionPrecedence1: PromotionPrecedence1{
		Uint:   5,
		Uint16: 6,
	},
	PromotionPrecedence2: PromotionPrecedence2{},
}

const promotionPrecedenceJson = `{
	"attributes": {
		"uint": 5,
		"Uint16": 6
	}
}`

func TestMarshalResource_PromotionPrecedence(t *testing.T) {
	got, err := MarshalResource(promotionPrecedenceMarshalValue)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, fmtJson(t, []byte(promotionPrecedenceJson)), fmtJson(t, got))
}

func TestUnmarshalResource_PromotionPrecedence(t *testing.T) {
	got := promotionPrecedence{}
	if err := UnmarshalResource([]byte(promotionPrecedenceJson), &got); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, promotionPrecedenceUnmarshalValue, got)
}

type AnonymousEliminationBase struct {
	SimpleIface
	Flt float64 `jsonapi:"attr,flt"`
}

type AnonymousElimination2 struct {
	AnonymousEliminationBase
}

type AnonymousElimination1 struct {
	AnonymousEliminationBase
}

type anonymousElimination struct {
	AnonymousElimination1
	AnonymousElimination2
}

func TestMarshalResource_AnonymousElimination(t *testing.T) {
	in := &anonymousElimination{
		AnonymousElimination1: AnonymousElimination1{
			AnonymousEliminationBase: AnonymousEliminationBase{
				Flt: 1,
			},
		},
		AnonymousElimination2: AnonymousElimination2{
			AnonymousEliminationBase: AnonymousEliminationBase{
				Flt: 2,
			},
		},
	}
	got, err := MarshalResource(in)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, fmtJson(t, []byte("{}")), fmtJson(t, got))
}

func TestUnmarshalResource_AnonymousElimination(t *testing.T) {
	got := anonymousElimination{}
	in := `
	{
		"attributes": {
			"flt": 2
		}
	}
	`
	if err := UnmarshalResource([]byte(in), &got); err != nil {
		t.Fatal(err)
	}

	want := anonymousElimination{}

	assert.Equal(t, want, got)
}

func TestMarshalResource_AnonymousElimination_InterfaceValue(t *testing.T) {
	in := &anonymousElimination{
		AnonymousElimination1: AnonymousElimination1{
			AnonymousEliminationBase: AnonymousEliminationBase{
				SimpleIface: SimpleIfaceImpl{
					Int: 1,
				},
				Flt: 2,
			},
		},
		AnonymousElimination2: AnonymousElimination2{
			AnonymousEliminationBase: AnonymousEliminationBase{
				Flt: 3,
			},
		},
	}
	got, err := MarshalResource(in)
	if err != nil {
		t.Fatal(err)
	}

	want := `{
		"attributes": {
			"int": 1
		}
	}`

	assert.Equal(t, fmtJson(t, []byte(want)), fmtJson(t, got))
}

func TestUnmarshalResource_AnonymousElimination_InterfaceValue(t *testing.T) {
	got := anonymousElimination{
		AnonymousElimination1: AnonymousElimination1{
			AnonymousEliminationBase: AnonymousEliminationBase{
				SimpleIface: &SimpleIfaceImpl{},
			},
		},
		AnonymousElimination2: AnonymousElimination2{
			AnonymousEliminationBase: AnonymousEliminationBase{},
		},
	}
	in := `
	{
		"attributes": {
			"int": 2
		}
	}
	`
	if err := UnmarshalResource([]byte(in), &got); err != nil {
		t.Fatal(err)
	}

	want := anonymousElimination{
		AnonymousElimination1: AnonymousElimination1{
			AnonymousEliminationBase: AnonymousEliminationBase{
				SimpleIface: &SimpleIfaceImpl{
					Int: 2,
				},
			},
		},
		AnonymousElimination2: AnonymousElimination2{
			AnonymousEliminationBase: AnonymousEliminationBase{},
		},
	}

	assert.Equal(t, want, got)
}

type SimpleIface interface {
	f()
}

type SimpleIfaceImpl struct {
	Int int `jsonapi:"attr,int"`
}

func (SimpleIfaceImpl) f() {}

type anonymousIface struct {
	SimpleIface
}

var anonymousIfaceValue = anonymousIface{
	SimpleIface: &SimpleIfaceImpl{
		Int: 1,
	},
}

const anonymousIfaceJson = `
{
	"attributes": {
		"int": 1
	}
}`

func TestMarshalResource_AnonymousIface(t *testing.T) {
	got, err := MarshalResource(anonymousIfaceValue)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, fmtJson(t, []byte(anonymousIfaceJson)), fmtJson(t, got))
}

func TestUnmarshalResource_AnonymousIface(t *testing.T) {
	got := anonymousIface{
		SimpleIface: &SimpleIfaceImpl{},
	}
	if err := UnmarshalResource([]byte(anonymousIfaceJson), &got); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, anonymousIfaceValue, got)
}

func TestMarshalResource_AnonymousIface_Value(t *testing.T) {
	in := anonymousIface{
		// NB not a pointer
		SimpleIface: SimpleIfaceImpl{
			Int: 1,
		},
	}
	got, err := MarshalResource(&in)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, fmtJson(t, []byte(anonymousIfaceJson)), fmtJson(t, got))
}

func TestUnmarshalResource_AnonymousIface_Value(t *testing.T) {
	got := anonymousIface{
		// NB not a pointer, not addressable
		SimpleIface: SimpleIfaceImpl{},
	}

	err := UnmarshalResource([]byte(anonymousIfaceJson), &got)
	assert.ErrorAs(t, err, addrOf(&UnmarshalErr{}))
}

var unsupportedTypes = []any{
	struct {
		Chan chan any `jsonapi:"id,type"`
	}{},
	struct {
		Chan chan any `jsonapi:"attr"`
	}{},
	struct {
		Chan chan any `jsonapi:"rel,type"`
	}{},
	struct {
		Chan chan any `jsonapi:"meta"`
	}{},
	struct {
		Func func() `jsonapi:"id,type"`
	}{},
	struct {
		Func func() `jsonapi:"attr"`
	}{},
	struct {
		Func func() `jsonapi:"rel,type"`
	}{},
	struct {
		Func func() `jsonapi:"meta"`
	}{},
	struct {
		Complex complex64 `jsonapi:"id,type"`
	}{},
	struct {
		Complex complex64 `jsonapi:"attr"`
	}{},
	struct {
		Complex complex64 `jsonapi:"rel,type"`
	}{},
	struct {
		Complex complex64 `jsonapi:"meta"`
	}{},
	struct {
		Complex complex128 `jsonapi:"id,type"`
	}{},
	struct {
		Complex complex128 `jsonapi:"attr"`
	}{},
	struct {
		Complex complex128 `jsonapi:"rel,type"`
	}{},
	struct {
		Complex complex128 `jsonapi:"meta"`
	}{},
}

func TestMarshalResource_UnsupportedTypes(t *testing.T) {
	for _, tc := range unsupportedTypes {
		t.Run("", func(t *testing.T) {
			bts, err := MarshalResource(tc)
			assert.Nil(t, bts)
			assert.ErrorAs(t, err, addrOf(&UnsupportedTypeErr{}))
		})
	}
}

func TestUnmarshalResource_UnsupportedTypes(t *testing.T) {
	data := []byte("{}")
	for _, tc := range unsupportedTypes {
		t.Run("", func(t *testing.T) {
			err := UnmarshalResource(data, &tc)
			assert.ErrorAs(t, err, addrOf(&UnsupportedTypeErr{}))
		})
	}
}

func TestMarshalResource_InputErr(t *testing.T) {
	data, err := MarshalResource(0)
	assert.Empty(t, data)
	assert.ErrorIs(t, err, ErrNotStruct)
}

func TestUnmarshalResource_InputTypeErr(t *testing.T) {
	type tp struct {
		Int int `jsonapi:"attr,int"`
	}

	type testCase struct {
		In       any
		Expected error
	}

	testCases := []testCase{
		{0, ErrNotStructPtr},
		{addrOf(0), ErrNotStructPtr},
		{tp{}, ErrNotStructPtr},
	}

	jsonData := []byte("{}")
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			err := UnmarshalResource(jsonData, tc.In)
			assert.ErrorIs(t, err, tc.Expected)
		})
	}
}

func TestMarshalResource_UnknownTagType(t *testing.T) {
	type tp struct {
		Int int `jsonapi:"xxx,int"`
	}

	data, err := MarshalResource(&tp{})
	assert.Empty(t, data)
	assert.ErrorAs(t, err, addrOf(&TagErr{}))
}

type ifaceFields struct {
	A any `jsonapi:"attr,a"`
	M any `jsonapi:"meta,m"`
	R any `jsonapi:"rel,name,type"`
}

var ifaceFieldsValue = ifaceFields{
	A: simpleStruct{Int: 1},
	M: addrOf(addrOf(&simpleStruct{Int: 2})),
	R: 3,
}

const ifaceFieldsJson = `
{
	"attributes": {
		"a": {
			"int": 1
		}
	},
	"relationships": {
		"name": {
			"data": { "type": "type", "id": 3 }
		}
	},
	"meta": {
		"m": {
			"int": 2
		}
	}
}`

func TestMarshalResource_InterfaceFields(t *testing.T) {
	got, err := MarshalResource(ifaceFieldsValue)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, fmtJson(t, []byte(ifaceFieldsJson)), fmtJson(t, got))
}

func TestUnmarshalResource_InterfaceFields(t *testing.T) {
	got := ifaceFields{
		A: simpleStruct{},
		M: addrOf(addrOf(&simpleStruct{})),
		R: 0,
	}
	if err := UnmarshalResource([]byte(ifaceFieldsJson), &got); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, ifaceFieldsValue, got)
}

func TestUnmarshalResource_UnitialisedInterfaceFields(t *testing.T) {
	got := ifaceFields{
		A: nil,
		M: nil,
		R: nil,
	}
	if err := UnmarshalResource([]byte(ifaceFieldsJson), &got); err != nil {
		t.Fatal(err)
	}

	want := ifaceFields{
		A: map[string]interface{}{
			"int": float64(1),
		},
		M: map[string]interface{}{
			"int": float64(2),
		},
		R: float64(3),
	}
	assert.Equal(t, want, got)
}

func TestMarshalResource_SelfRefPtr(t *testing.T) {
	// marshaling with a ptr cycle should return
	// a self-referential pointer err
	type T struct {
		A any `jsonapi:"attr"`
	}

	in := T{}
	p := &in.A
	in.A = &p

	_, err := MarshalResource(&in)
	assert.ErrorIs(t, err, ErrSelfRefPtr)
}

func TestUnmarshalResource_SelfRefPtr(t *testing.T) {
	// unmarshaling with a ptr cycle should work. this
	// is the same behaviour as the standard JSON lib
	type T struct {
		A any `jsonapi:"attr,a"`
	}

	got := T{}
	p := &got.A
	got.A = &p

	data := `{
		"attributes": {
			"a": 1
		}
	}`

	if err := UnmarshalResource([]byte(data), &got); err != nil {
		t.Fatal(err)
	}

	var i interface{} = float64(1)
	want := T{
		A: addrOf(addrOf(i)),
	}

	assert.Equal(t, want, got)
}

func TestMarshalResource_AnonymousSelfRefPtr(t *testing.T) {
	type I interface{}

	type T struct {
		I
	}

	in := &T{}
	p := &in.I
	in.I = &p

	_, err := MarshalResource(&in)
	assert.ErrorIs(t, err, ErrSelfRefPtr)
}

func TestUnmarshalResource_AnonymousSelfRefPtr(t *testing.T) {
	// TODO: does this even make sense?
}

type CycleType struct {
	*CycleType
	Int int `jsonapi:"attr,int"`
}

const cycleTypeJson = `
{
	"attributes": {
		"int": 3
	}
}
`

func TestMarshalResource_TypeCycle(t *testing.T) {
	in := CycleType{
		CycleType: &CycleType{
			CycleType: &CycleType{
				Int: 1,
			},
			Int: 2,
		},
		Int: 3,
	}

	got, err := MarshalResource(&in)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, fmtJson(t, []byte(cycleTypeJson)), fmtJson(t, got))
}

func TestUnmarshalResource_TypeCycle(t *testing.T) {
	got := CycleType{}
	if err := UnmarshalResource([]byte(cycleTypeJson), &got); err != nil {
		t.Fatal(err)
	}

	want := CycleType{
		Int: 3,
	}

	assert.Equal(t, want, got)
}

type mapMarshalUnmarshaler struct {
	Id            string
	Attributes    map[string]interface{}
	Meta          map[string]interface{}
	Relationships map[string]interface{}
}

func (m *mapMarshalUnmarshaler) MarshalJsonApiResource() ([]byte, error) {
	j, err := json.Marshal(m.Id)
	if err != nil {
		return nil, err
	}
	r := &Resource{
		ResourceIdentifier: ResourceIdentifier{
			Type: "type",
			Id:   json.RawMessage(j),
			Meta: map[string]json.RawMessage{},
		},
		Attributes:         map[string]json.RawMessage{},
		ToOneRelationships: map[string]*ToOneResourceLinkage{},
	}

	for k, v := range m.Attributes {
		j, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		r.Attributes[k] = json.RawMessage(j)
	}

	for k, v := range m.Meta {
		j, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		r.Meta[k] = json.RawMessage(j)
	}

	for k, v := range m.Relationships {
		j, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		r.ToOneRelationships[k] = &ToOneResourceLinkage{
			Data: ResourceIdentifier{
				Type: "rel-type",
				Id:   json.RawMessage(j),
			},
		}
	}

	return json.Marshal(r)
}

func (m *mapMarshalUnmarshaler) UnmarshalJsonApiResource(data []byte) error {
	r := Resource{
		Attributes: map[string]json.RawMessage{},
	}
	if err := json.Unmarshal(data, &r); err != nil {
		return err
	}

	if err := json.Unmarshal(r.ResourceIdentifier.Id, &m.Id); err != nil {
		return err
	}

	m.Attributes = map[string]interface{}{}
	for k, v := range r.Attributes {
		var i interface{}
		if err := json.Unmarshal(v, &i); err != nil {
			return err
		}
		m.Attributes[k] = i
	}

	m.Meta = map[string]interface{}{}
	for k, v := range r.Meta {
		var i interface{}
		if err := json.Unmarshal(v, &i); err != nil {
			return err
		}
		m.Meta[k] = i
	}

	m.Relationships = map[string]interface{}{}
	for k, v := range r.ToOneRelationships {
		var i interface{}
		if err := json.Unmarshal(v.Data.Id, &i); err != nil {
			return err
		}
		m.Relationships[k] = i
	}

	return nil
}

var mapMarshalUnmarshalerValue = mapMarshalUnmarshaler{
	Id: "id",
	Attributes: map[string]interface{}{
		"int": float64(1),
	},
	Meta: map[string]interface{}{
		"float64": 2.1,
	},
	Relationships: map[string]interface{}{
		"name": "3",
	},
}

const mapMarshalUnmarshalerJson = `
{
	"id": "id",
	"type": "type",
	"attributes": {
		"int": 1
	},
	"meta": {
		"float64": 2.1
	},
	"relationships": {
		"name": { "data": { "type": "rel-type", "id": "3" } }
	}
}`

func TestMarshalResource_MapResourceMarshaler(t *testing.T) {
	got, err := MarshalResource(&mapMarshalUnmarshalerValue)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, fmtJson(t, []byte(mapMarshalUnmarshalerJson)), fmtJson(t, got))
}

func TestUnmarshalResource_MapResourceUnmarshaler(t *testing.T) {
	got := mapMarshalUnmarshaler{}
	if err := UnmarshalResource([]byte(mapMarshalUnmarshalerJson), &got); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, mapMarshalUnmarshalerValue, got)
}

type aliasMarshalUnmarshaler struct {
	S *simpleStruct
}

func (m *aliasMarshalUnmarshaler) MarshalJsonApiResource() ([]byte, error) {
	type alias struct {
		I int `jsonapi:"attr,i"`
	}

	a := &alias{
		I: m.S.Int,
	}

	return MarshalResource(&a)
}

func (m *aliasMarshalUnmarshaler) UnmarshalJsonApiResource(data []byte) error {
	type alias struct {
		I int `jsonapi:"attr,i"`
	}

	a := &alias{}

	if err := UnmarshalResource(data, &a); err != nil {
		return err
	}

	m.S = &simpleStruct{
		Int: a.I,
	}

	return nil
}

var aliasMarshalUnmarshalerValue = aliasMarshalUnmarshaler{
	S: &simpleStruct{
		Int: 1,
	},
}

const aliasMarshalUnmarshalerJson = `
{
	"attributes": {
		"i": 1
	}
}`

func TestMarshalResource_AliasResourceMarshaler(t *testing.T) {
	got, err := MarshalResource(&aliasMarshalUnmarshalerValue)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, fmtJson(t, []byte(aliasMarshalUnmarshalerJson)), fmtJson(t, got))
}

func TestUnmarshalResource_AliasResourceUnmarshaler(t *testing.T) {
	got := aliasMarshalUnmarshaler{}
	if err := UnmarshalResource([]byte(aliasMarshalUnmarshalerJson), &got); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, aliasMarshalUnmarshalerValue, got)
}

func TestMarshalResource_PtrToResourceMarshaler(t *testing.T) {
	in := &aliasMarshalUnmarshalerValue
	got, err := MarshalResource(addrOf(addrOf(in)))
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, fmtJson(t, []byte(aliasMarshalUnmarshalerJson)), fmtJson(t, got))
}

func TestUnmarshalResource_PtrToResourceUnmarshaler(t *testing.T) {
	got := &aliasMarshalUnmarshaler{}
	if err := UnmarshalResource([]byte(aliasMarshalUnmarshalerJson), addrOf(addrOf(got))); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, &aliasMarshalUnmarshalerValue, got)
}

type formatMarshalUnmarshaler struct {
	I     int    `jsonapi:"attr,i"`
	R     int    `jsonapi:"rel,r1,rel-type"`
	Link1 string `jsonapi:"-"`
	Link2 string `jsonapi:"-"`
	Link3 string `jsonapi:"-"`
	Link4 string `jsonapi:"-"`
}

func (m *formatMarshalUnmarshaler) MarshalJsonApiResource() ([]byte, error) {
	r, err := FormatResource(m)
	if err != nil {
		return nil, err
	}

	r.Links = map[string]*Link{
		"l1": {
			LinkString: m.Link1,
		},
		"l2": {
			LinkObject: LinkObject{
				Href: m.Link2,
			},
		},
	}

	r.ToOneRelationships["r1"].Links = map[string]*Link{
		"l3": {
			LinkString: m.Link3,
		},
		"l4": {
			LinkObject: LinkObject{
				Href: m.Link4,
			},
		},
	}

	return json.Marshal(r)
}

func (m *formatMarshalUnmarshaler) UnmarshalJsonApiResource(data []byte) error {
	r := &Resource{}
	if err := json.Unmarshal(data, r); err != nil {
		return err
	}

	if err := DeformatResource(r, m); err != nil {
		return err
	}

	m.Link1 = r.Links["l1"].LinkString
	m.Link2 = r.Links["l2"].LinkObject.Href
	m.Link3 = r.ToOneRelationships["r1"].Links["l3"].LinkString
	m.Link4 = r.ToOneRelationships["r1"].Links["l4"].LinkObject.Href
	return nil
}

var formatMarshalUnmarshalerValue = formatMarshalUnmarshaler{
	I:     1,
	R:     2,
	Link1: "http://test.com/1",
	Link2: "http://test.com/2",
	Link3: "http://test.com/3",
	Link4: "http://test.com/4",
}

const formatMarshalUnmarshalerJson = `
{
	"attributes": {
		"i": 1
	},
	"links": {
		"l1": "http://test.com/1",
		"l2": {
			"href": "http://test.com/2"
		}
	},
	"relationships": {
		"r1": {
			"data": {
				"type": "rel-type",
				"id": 2
			},
			"links": {
				"l3": "http://test.com/3",
				"l4": {
					"href": "http://test.com/4"
				}
			}
		}
	}
}
`

func TestMarshalResource_FormatResourceMarshaler(t *testing.T) {
	got, err := MarshalResource(&formatMarshalUnmarshalerValue)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, fmtJson(t, []byte(formatMarshalUnmarshalerJson)), fmtJson(t, got))
}

func TestUnmarshalResource_FormatResourceUnmarshaler(t *testing.T) {
	got := formatMarshalUnmarshaler{}
	if err := UnmarshalResource([]byte(formatMarshalUnmarshalerJson), &got); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, formatMarshalUnmarshalerValue, got)
}

type stringTag struct {
	Id     int     `jsonapi:"id,tp,string"`
	Attr   float32 `jsonapi:"attr,a,string"`
	Meta   string  `jsonapi:"meta,m,string"`
	ToOne  int     `jsonapi:"rel,r1,r1-type,string"`
	ToMany []int   `jsonapi:"rel,r2,r2-type,string"`
}

var stringTagValue stringTag = stringTag{
	Id:     1,
	Attr:   2.1,
	Meta:   "value",
	ToOne:  3,
	ToMany: []int{-4, 5},
}

const stringTagJson = `{
	"id": "1",
	"type": "tp",
	"attributes": {
		"a": "2.1"
	},
	"meta": {
		"m": "value"
	},
	"relationships": {
		"r1": { "data": { "type": "r1-type", "id": "3" } },
		"r2": {
			"data": [ 
				{ "type": "r2-type", "id": "-4" }, 
				{ "type": "r2-type", "id": "5" } 
			]
		}
	}
}`

func TestMarshalResource_StringTag(t *testing.T) {
	got, err := MarshalResource(stringTagValue)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, fmtJson(t, []byte(stringTagJson)), fmtJson(t, got))
}

func TestUnmarshalResource_StringTag(t *testing.T) {
	got := stringTag{}
	if err := UnmarshalResource([]byte(stringTagJson), &got); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, stringTagValue, got)
}

func fmtJson(t *testing.T, data []byte) string {
	m := map[string]interface{}{}
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatal(err)
	}

	data, err := json.MarshalIndent(m, "", "    ")
	if err != nil {
		t.Fatal(err)
	}

	return string(data)
}

func addrOf[A any](a A) *A {
	return &a
}
