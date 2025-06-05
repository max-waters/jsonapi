package jsonapi

import (
	"cmp"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"unsafe"

	"reflect"
	"slices"
	"strings"
)

const (
	// tag keys
	TagKeyJson    = "json"
	TagKeyJsonApi = "jsonapi"
	// tag values
	TagValueIgnore = "-"
	TagValueId     = "id"
	TagValueAttr   = "attr"
	TagValueRel    = "rel"
	TagValueMeta   = "meta"
	TagValueLink   = "link"
	// options
	TagValueOmitEmpty = "omitempty"
	TagValueString    = "string"
)

var NullJson = json.RawMessage([]byte("null"))

type TagErr struct {
	Field string
	Err   error
}

func (e *TagErr) Error() string {
	return "tag error on field '" + e.Field + "': " + e.Err.Error()
}

type UnmarshalErr struct {
	Field string
	Err   error
}

func (e *UnmarshalErr) Error() string {
	return "unmarshal error on field '" + e.Field + "': " + e.Err.Error()
}

type MarshalErr struct {
	Field string
	Err   error
}

func (e *MarshalErr) Error() string {
	return "marshal error on field '" + e.Field + "': " + e.Err.Error()
}

type UnsupportedTypeErr struct {
	Field string
	Kind  reflect.Kind
}

func (e *UnsupportedTypeErr) Error() string {
	return "unsupported type on field " + e.Field + "': " + e.Kind.String()
}

var (
	ErrNotStructPtr = fmt.Errorf("not a struct pointer")
	ErrNotStruct    = fmt.Errorf("not a struct")
	ErrSelfRefPtr   = fmt.Errorf("self-referential pointer")
)

type ResourceUnmarshaler interface {
	UnmarshalJsonApiResource([]byte) error
}

type ResourceMarshaler interface {
	MarshalJsonApiResource() ([]byte, error)
}

var (
	resourceMarshalerType   = reflect.TypeFor[ResourceMarshaler]()
	resourceUnmarshalerType = reflect.TypeFor[ResourceUnmarshaler]()
)

type resourceIdentifier struct {
	Type string          `json:"type,omitempty"`
	Id   json.RawMessage `json:"id,omitempty"`
}

type toOneRelationship struct {
	Data  resourceIdentifier         `json:"data,omitempty"`
	Links map[string]json.RawMessage `json:"links,omitempty"`
	Meta  map[string]json.RawMessage `json:"meta,omitempty"`
}

type toManyRelationship struct {
	Data  []resourceIdentifier       `json:"data,omitempty"`
	Links map[string]json.RawMessage `json:"links,omitempty"`
	Meta  map[string]json.RawMessage `json:"meta,omitempty"`
}

type resource struct {
	resourceIdentifier
	Attributes          map[string]json.RawMessage
	ToOneRelationships  map[string]toOneRelationship
	ToManyRelationships map[string]toManyRelationship
	Links               map[string]json.RawMessage
	Meta                map[string]json.RawMessage
}

func newResource() resource {
	return resource{
		resourceIdentifier:  resourceIdentifier{},
		Attributes:          map[string]json.RawMessage{},
		ToOneRelationships:  map[string]toOneRelationship{},
		ToManyRelationships: map[string]toManyRelationship{},
		Meta:                map[string]json.RawMessage{},
		Links:               map[string]json.RawMessage{},
	}
}

func (r *resource) MarshalJSON() ([]byte, error) {
	type alias struct {
		resourceIdentifier
		Attributes    map[string]json.RawMessage `json:"attributes,omitempty"`
		Relationships map[string]any             `json:"relationships,omitempty"`
		Links         map[string]json.RawMessage `json:"links,omitempty"`
		Meta          map[string]json.RawMessage `json:"meta,omitempty"`
	}
	a := alias{
		resourceIdentifier: r.resourceIdentifier,
		Attributes:         r.Attributes,
		Links:              r.Links,
		Meta:               r.Meta,
	}

	if len(r.ToOneRelationships)+len(r.ToManyRelationships) > 0 {
		a.Relationships = make(map[string]any, len(r.ToOneRelationships)+len(r.ToManyRelationships))

		for k, v := range r.ToOneRelationships {
			a.Relationships[k] = v
		}
		for k, v := range r.ToManyRelationships {
			a.Relationships[k] = v
		}
	}
	return json.Marshal(a)
}

func (r *resource) UnmarshalJSON(data []byte) error {
	type relAlias struct {
		Meta  map[string]json.RawMessage `json:"meta"`
		Data  json.RawMessage            `json:"data"`
		Links map[string]json.RawMessage `json:"links"`
	}

	type alias struct {
		resourceIdentifier
		Attributes    map[string]json.RawMessage `json:"attributes"`
		Relationships map[string]relAlias        `json:"relationships"`
		Links         map[string]json.RawMessage `json:"links"`
		Meta          map[string]json.RawMessage `json:"meta"`
	}

	a := alias{}

	if err := json.Unmarshal(data, &a); err != nil {
		return err
	}

	r.resourceIdentifier = a.resourceIdentifier
	r.Attributes = a.Attributes
	r.Links = a.Links
	r.Meta = a.Meta

	nToOne := 0
	nToMany := 0
	for _, rel := range a.Relationships {
		switch rel.Data[0] {
		case '[':
			nToMany++
		case '{':
			nToOne++
		default:
			return fmt.Errorf("cannot unmarshal into relationship data")
		}
	}

	if nToMany > 0 {
		r.ToManyRelationships = make(map[string]toManyRelationship, nToMany)
	}
	if nToOne > 0 {
		r.ToOneRelationships = make(map[string]toOneRelationship, nToOne)
	}

	for name, rel := range a.Relationships {
		switch rel.Data[0] {
		case '[':
			if r.ToManyRelationships == nil {

			}
			ids := []resourceIdentifier{}
			if err := json.Unmarshal(rel.Data, &ids); err != nil {
				return err
			}
			r.ToManyRelationships[name] = toManyRelationship{
				Meta:  rel.Meta,
				Data:  ids,
				Links: rel.Links,
			}
		case '{':
			if r.ToOneRelationships == nil {

			}
			id := resourceIdentifier{}
			if err := json.Unmarshal(rel.Data, &id); err != nil {
				return err
			}
			r.ToOneRelationships[name] = toOneRelationship{
				Meta:  rel.Meta,
				Data:  id,
				Links: rel.Links,
			}
		default:
			return fmt.Errorf("cannot unmarshal into relationship data")
		}
	}

	return nil
}

func MarshalResource(a any, opts ...marshalResourceOpt) ([]byte, error) {
	marshalOpts := marshalResourceOpts{}
	for _, opt := range opts {
		if opt == nil {
			return nil, fmt.Errorf("nil option")
		}
		marshalOpts = opt(marshalOpts)
	}

	v := reflect.ValueOf(a)

	v, err := derefInput(v, resourceMarshalerType)
	if err != nil {
		return nil, fmt.Errorf("jsonapi: dereferencing input: %w", err)
	}

	if v.Type().Implements(resourceMarshalerType) {
		return v.Interface().(ResourceMarshaler).MarshalJsonApiResource()
	}

	if v.Type().Kind() != reflect.Struct {
		return nil, ErrNotStruct
	}

	r, err := format(a, v, marshalOpts)
	if err != nil {
		return nil, fmt.Errorf("jsonapi: %w", err)
	}

	data, err := json.Marshal(&r)
	if err != nil {
		return nil, fmt.Errorf("jsonapi: marshaling resource: %w", err)
	}

	return data, nil
}

func format(a any, v reflect.Value, opts marshalResourceOpts) (resource, error) {
	fields, err := parseTags(v)
	if err != nil {
		return resource{}, fmt.Errorf("parsing tags: %w", err)
	}

	r := newResource()

	for _, f := range fields {
		if err := marshalField(v, &r, f); err != nil {
			return resource{}, fmt.Errorf("marshaling field "+f.tag.name+": %w", err)
		}
	}

	if err := applyMarshalOpts(a, r, opts); err != nil {
		return resource{}, err
	}

	return r, nil
}

func marshalField(v reflect.Value, r *resource, f field) error {
	if f.tag.typ == TagValueRel {
		return marshalRel(v, r, f)
	}

	if f.tag.typ == TagValueId {
		r.Type = f.tag.rscType
	}

	v, err := fieldByIndex(v, f.idxs)
	if err != nil {
		return err
	}

	v, err = derefValue(v)
	if err != nil {
		return err
	}

	if f.tag.omitempty && isEmpty(v) {
		return nil
	}

	j, err := marshalJson(v, f.tag.quote)
	if err != nil {
		return &MarshalErr{f.tag.name, err}
	}

	switch f.tag.typ {
	case TagValueId:
		r.resourceIdentifier.Id = j
	case TagValueAttr:
		r.Attributes[f.tag.name] = j
	case TagValueMeta:
		r.Meta[f.tag.name] = j
	case TagValueLink:
		r.Links[f.tag.name] = j
	default:
		return errors.New("unknown tag type " + f.tag.typ)
	}

	return nil
}

func UnmarshalResource(data []byte, a any) error {
	v := reflect.ValueOf(a)

	if v.Kind() != reflect.Pointer {
		return ErrNotStructPtr
	}

	v, err := derefInput(v, resourceUnmarshalerType)
	if err != nil {
		return fmt.Errorf("jsonapi: dereferencing input: %w", err)
	}

	if v.Type().Implements(resourceUnmarshalerType) {
		return v.Interface().(ResourceUnmarshaler).UnmarshalJsonApiResource(data)
	}

	if v.Type().Kind() != reflect.Struct {
		return ErrNotStructPtr
	}

	r := newResource()
	if err := json.Unmarshal(data, &r); err != nil {
		return fmt.Errorf("jsonapi: unmarshaling resource: %w", err)
	}

	return deformat(v, r)
}

func deformat(v reflect.Value, r resource) error {
	fields, err := parseTags(v)
	if err != nil {
		return fmt.Errorf("jsonapi: parsing tags: %w", err)
	}

	for _, f := range fields {
		if err := unmarshalField(v, &r, f); err != nil {
			return fmt.Errorf("jsonapi: unmarshaling field '"+f.tag.name+"': %w", err)
		}
	}
	return nil
}

func unmarshalField(v reflect.Value, r *resource, f field) error {
	if f.tag.typ == TagValueRel {
		return unmarshalRel(v, r, f)
	}

	var j json.RawMessage
	switch f.tag.typ {
	case TagValueId:
		j = r.resourceIdentifier.Id
	case TagValueAttr:
		j = r.Attributes[f.tag.name]
	case TagValueMeta:
		j = r.Meta[f.tag.name]
	case TagValueLink:
		j = r.Links[f.tag.name]
	default:
		return errors.New("unknown tag type " + f.tag.typ)
	}

	if len(j) == 0 {
		return nil
	}

	v, err := initFieldByIndex(v, f.idxs)
	if err != nil {
		return err
	}

	if err := unmarshalJson(j, v, f.tag.quote); err != nil {
		return &UnmarshalErr{f.tag.name, err}
	}
	return nil
}

var enableFieldCache = true
var fieldCache sync.Map // map[reflect.Type][]field

// parseTags retrieves all attributes, relationships,
// etc from the input value.
//   - performs a breadth-first search over the value
//     rooted at v
//   - if a struct field is found with no value, the
//     search continues over the type tree rooted at
//     f's type
//   - modelled on the equivalent function in the
//     encoding/json package to reduce heap allocs
//     (see issue #1)
func parseTags(v reflect.Value) ([]field, error) {
	if fields, ok := fieldCache.Load(v.Type()); ok {
		return fields.([]field), nil
	}

	// every element in the queue represents a
	// struct, either a type or a value
	type structElem struct {
		t    reflect.Type
		v    reflect.Value
		ok   bool  // true if the value is present
		idxs []int // path to this structElem
	}

	var fields []field

	types := map[reflect.Type]bool{}

	next := []structElem{{t: v.Type(), v: v, ok: true}}
	var current []structElem

	// nb no allocations happen until needed
	nextCount := map[reflect.Type]int{}

	cache := true
	for len(next) > 0 {
		current, next = next, current[:0]
		clear(nextCount)

		// count struct fields
		nfs := 0
		for _, c := range current {
			nfs += c.t.NumField()
		}

		// pre-allocate space in one go
		fields = slices.Grow(fields, nfs) // alloc
		next = slices.Grow(next, nfs)     // alloc

		for _, c := range current {
			types[c.t] = true

			for i := range c.t.NumField() {
				f := c.t.Field(i) // alloc (!)

				if !f.IsExported() && !f.Anonymous {
					continue
				}

				typ, opts, ok := splitTypeAndOpts(f)

				if typ == TagValueIgnore {
					continue
				}

				fIdxs := make([]int, len(c.idxs)+1) // alloc
				copy(fIdxs, c.idxs)
				fIdxs[len(fIdxs)-1] = i

				if !ok {
					if f.Anonymous {
						ft := derefType(f.Type)
						if ft.Kind() == reflect.Struct {
							if !types[ft] && nextCount[ft] < 2 {
								nextCount[ft] = nextCount[ft] + 1
								if c.ok {
									fv, err := derefValue(c.v.Field(i))
									if err != nil {
										return nil, err
									}
									if fv.Kind() == reflect.Struct {
										next = append(next, structElem{ft, fv, true, fIdxs}) // alloc
									}
								} else {
									next = append(next, structElem{ft, reflect.Value{}, false, fIdxs}) // alloc
								}
							}
						}

						if ft.Kind() == reflect.Interface {
							cache = false
							if c.ok {
								fv, err := derefValue(c.v.Field(i))
								if err != nil {
									return nil, err
								}

								if fv.Kind() == reflect.Struct {
									fvt := fv.Type()
									if !types[fvt] && nextCount[fvt] < 2 {
										next = append(next, structElem{fvt, fv, true, fIdxs}) // alloc
										nextCount[fvt] = nextCount[fvt] + 1
									}
								}
							}
						}

						continue
					}

					typ = TagValueAttr
				}

				tag, err := parseTag(f, typ, opts)
				if err != nil {
					return nil, err
				}

				fld := field{
					tag:  tag,
					idxs: fIdxs,
				}

				fields = append(fields, fld)
			}
		}
	}

	// sort by type, then name, then depth, then name precedence
	slices.SortFunc(fields, func(a, b field) int {
		if c := cmp.Compare(a.tag.typ, b.tag.typ); c != 0 {
			return c
		}

		if c := cmp.Compare(a.tag.name, b.tag.name); c != 0 {
			return c
		}

		if c := cmp.Compare(len(a.idxs), len(b.idxs)); c != 0 {
			return c
		}

		return -cmp.Compare(a.tag.namePrec, b.tag.namePrec)
	})

	// now filter all fields that are overridden by those
	// with a higher precedence
	nFiltered := 0
	for nType, i := 0, 0; i < len(fields); i += nType {
		// find subslice of all fields of the same type
		typ := fields[i].tag.typ
		for nType = 1; i+nType < len(fields); nType++ {
			if fields[i+nType].tag.typ != typ {
				break
			}
		}

		for nName, j := 0, i; j < i+nType; j += nName {
			// find subslice of all fields with the same name (and type)
			name := fields[j].tag.name
			for nName = 1; j+nName < i+nType; nName++ {
				if fields[j+nName].tag.name != name {
					break
				}
			}

			// if there are multiple with the same name and type,
			// get the dominant field
			field, ok := getDominantField(fields[j : j+nName])
			if ok {
				// copy back into original slice to save allocs
				fields[nFiltered] = field
				nFiltered++
			}

		}
	}

	fields = fields[:nFiltered]
	if enableFieldCache && cache {
		fieldCache.Store(v.Type(), fields) // 4 allocs
	}

	return fields, nil
}

// getDominantField returns the highest precedence
// field from the supplied list, with (zero, false)
// indicating a that no dominant tag can be determined.
// Assumes that the input list items all have the same name and
// type, and are sorted by depth then name precedence
func getDominantField(fs []field) (field, bool) {
	if len(fs) == 0 {
		return field{}, false
	}

	if len(fs) == 1 {
		return fs[0], true
	}

	// if the two first items have the same depth and name prec then
	// no dominant item can be determined
	if len(fs[0].idxs) == len(fs[1].idxs) && fs[0].tag.namePrec == fs[1].tag.namePrec {
		return field{}, false
	}

	// the first item must take precedence
	return fs[0], true
}

func parseTag(f reflect.StructField, typ, opts string) (tag, error) {
	k := derefType(f.Type).Kind()
	switch k {
	case reflect.Func, reflect.Chan, reflect.Complex64, reflect.Complex128:
		return tag{}, &UnsupportedTypeErr{Field: f.Name, Kind: k}
	}

	switch typ {
	case TagValueId:
		return parseIdTag(f, opts)
	case TagValueAttr:
		return parseAttrTag(f, opts)
	case TagValueMeta:
		return parseMetaTag(f, opts)
	case TagValueRel:
		return parseRelTag(f, opts)
	case TagValueLink:
		return parseLinkTag(f, opts)
	default:
		return tag{}, &TagErr{f.Name, errors.New("unknown tag type: " + typ)}
	}
}

// field represents a tagged struct field,
// with tag representing the annotated tag, and idxs
// uniquely identifying this field with its path
// from the top-level struct
type field struct {
	// the tag information annotated onto this struct field
	tag tag
	// idxs represents this and all ancestor fields' indexes
	// within their parent structs
	idxs []int
}

// tag represents a jsonapi struct tag
type tag struct {
	// The jsonapi tag type, eg attribute, relationship etc
	typ string
	// The name that will appear in the output JSON.
	name string
	// The precedence of the name, with a jsonapi tag
	// name being the highest, then a json tag, then
	// the declared field name
	namePrec int
	// If this type is relationship or id, this field
	// defines the resource type
	rscType string
	// whether the "string" flag was specified
	quote bool
	// whether the "omitempty" flag was specified
	omitempty bool
}

// parseIdTag parses an id tag, eg `jsonapi:"id,name,type,opt1,opt2..."`
func parseIdTag(f reflect.StructField, opts string) (tag, error) {
	rscType, opts := splitFirstAndOpts(opts)
	if rscType == "" {
		return tag{}, &TagErr{f.Name, fmt.Errorf("required: type")}
	}

	omitempty, quote := optFlags(opts)

	return tag{
		typ:       TagValueId,
		rscType:   rscType,
		omitempty: omitempty,
		quote:     quote,
	}, nil
}

// parseAttrTag parses an attribute tag, eg `jsonapi:"attr,name,opt1,opt2..."`
func parseAttrTag(f reflect.StructField, opts string) (tag, error) {
	name, namePrec, opts := splitNameAndOpts(f, opts)
	omitempty, quote := optFlags(opts)

	return tag{
		typ:       TagValueAttr,
		name:      name,
		namePrec:  namePrec,
		omitempty: omitempty,
		quote:     quote,
	}, nil
}

// parseRelTag parses a relationship tag, eg `jsonapi:"rel,name,type,opt1,opt2..."`
func parseRelTag(f reflect.StructField, opts string) (tag, error) {
	name, namePrec, opts := splitNameAndOpts(f, opts)
	rscType, opts := splitFirstAndOpts(opts)
	if rscType == "" {
		return tag{}, &TagErr{f.Name, fmt.Errorf("required: type")}
	}

	omitempty, quote := optFlags(opts)

	return tag{
		typ:       TagValueRel,
		name:      name,
		namePrec:  namePrec,
		rscType:   rscType,
		omitempty: omitempty,
		quote:     quote,
	}, nil
}

func marshalRel(v reflect.Value, r *resource, f field) error {
	v, err := fieldByIndex(v, f.idxs)
	if err != nil {
		return err
	}

	v, err = derefValue(v)
	if err != nil {
		return err
	}

	if f.tag.omitempty && isEmpty(v) {
		return nil
	}

	if isToOne(v) {
		return marshalToOneRel(v, r, f)
	}

	return marshalToManyRel(v, r, f)
}

func marshalToOneRel(v reflect.Value, r *resource, f field) error {
	j, err := marshalJson(v, f.tag.quote)
	if err != nil {
		return &MarshalErr{f.tag.name, err}
	}

	r.ToOneRelationships[f.tag.name] = toOneRelationship{
		Data: resourceIdentifier{
			Type: f.tag.rscType,
			Id:   j,
		},
	}
	return nil
}

func marshalToManyRel(v reflect.Value, r *resource, f field) error {
	r.ToManyRelationships[f.tag.name] = toManyRelationship{
		Data: make([]resourceIdentifier, v.Len()),
	}

	for i := range v.Len() {
		vi, err := derefValue(v.Index(i))
		if err != nil {
			return err
		}

		j, err := marshalJson(vi, f.tag.quote)
		if err != nil {
			return &MarshalErr{f.tag.name, err}
		}

		r.ToManyRelationships[f.tag.name].Data[i] = resourceIdentifier{
			Type: f.tag.rscType,
			Id:   j,
		}
	}

	return nil
}

func unmarshalRel(v reflect.Value, r *resource, f field) error {
	fv, err := fieldByIndex(v, f.idxs)
	if err != nil {
		return err
	}

	if isToOne(fv) {
		return unmarshalToOneRel(v, r, f)
	}
	return unmarshalToManyRel(v, r, f)
}

func unmarshalToOneRel(v reflect.Value, r *resource, f field) error {
	rel, ok := r.ToOneRelationships[f.tag.name]
	if !ok {
		return nil
	}

	if len(rel.Data.Id) == 0 {
		return nil
	}

	v, err := initFieldByIndex(v, f.idxs)
	if err != nil {
		return err
	}

	if err := unmarshalJson(rel.Data.Id, v, f.tag.quote); err != nil {
		return &UnmarshalErr{f.tag.name, err}
	}
	return nil
}

func unmarshalToManyRel(v reflect.Value, r *resource, f field) error {
	rels, ok := r.ToManyRelationships[f.tag.name]
	if !ok {
		return nil
	}

	if len(rels.Data) == 0 {
		return nil
	}

	v, err := initFieldByIndex(v, f.idxs)
	if err != nil {
		return err
	}

	v.Grow(len(rels.Data) - v.Cap())
	v.SetLen(len(rels.Data))
	for i, rel := range rels.Data {
		elem := v.Index(i)
		initValue(elem)
		if err := unmarshalJson(rel.Id, elem, f.tag.quote); err != nil {
			return &UnmarshalErr{f.tag.name, err}
		}
	}

	return nil
}

// isToOne returns whether the supplied value represents a to-one or
// to-many relationship. A to-many relationship must be an array, or a slice
// of anything that is not a byte.
func isToOne(fv reflect.Value) bool {
	return fv.Kind() != reflect.Array && (fv.Kind() != reflect.Slice || fv.Type().Elem().Kind() == reflect.Uint8)
}

// parseMetaTag parses a meta tag, eg `jsonapi:"meta,name,opt1,opt2..."`
func parseMetaTag(f reflect.StructField, opts string) (tag, error) {
	name, namePrec, opts := splitNameAndOpts(f, opts)
	omitempty, quote := optFlags(opts)

	return tag{
		typ:       TagValueMeta,
		name:      name,
		namePrec:  namePrec,
		omitempty: omitempty,
		quote:     quote,
	}, nil
}

// parseMetaTag parses a link tag, eg `jsonapi:"link,name,opt1,opt2..."`
func parseLinkTag(f reflect.StructField, opts string) (tag, error) {
	name, namePrec, opts := splitNameAndOpts(f, opts)
	omitempty, quote := optFlags(opts)

	return tag{
		typ:       TagValueLink,
		name:      name,
		namePrec:  namePrec,
		omitempty: omitempty,
		quote:     quote,
	}, nil
}

// splitTypeAndOpts extracts the jsonapi tag value from the supplied tag
// and returns the type string and all remaining options. The bool represents
// whether a tag was found.
// Eg `jsonapi:"attr,name,omitempty"` returns ("attribute", "name,omitempty", true )
func splitTypeAndOpts(f reflect.StructField) (string, string, bool) {
	value, ok := f.Tag.Lookup(TagKeyJsonApi)
	if !ok {
		return "", "", false
	}

	typ, opts, _ := strings.Cut(value, ",")
	return typ, opts, true
}

// splitNameAndOpts extracts the name and precedence from the supplied
// field and tag options. It returns the name, the name's precedence, and the
// remaining options.
// NB assumes that the opts string does not contain the type.
// If the opts string contains a declared name, then it is returned with
// precedence 3. If there is no declared name but there is a decalred json
// name, that is returned with precedence 2. Otherwise the field name is returned
// with precedence 1.
func splitNameAndOpts(f reflect.StructField, opts string) (string, int, string) {
	name, opts := splitFirstAndOpts(opts)
	if name != "" {
		return name, 3, opts
	}

	name, _, _ = strings.Cut(f.Tag.Get(TagKeyJson), ",")
	if name != "" {
		return name, 2, opts
	}

	return f.Name, 1, opts
}

// splitFirstAndOpts extracts the first opt from the opts list.
// Returns the extracted opt and all remaining opts.
func splitFirstAndOpts(opts string) (string, string) {
	fst, opts, _ := strings.Cut(opts, ",")
	return fst, opts
}

// optFlags gets the values of the omitempty and
// string flags from the supplied opts.
func optFlags(opts string) (bool, bool) {
	omitempty := false
	quote := false
	for opts != "" {
		opt, rest, _ := strings.Cut(opts, ",")
		switch opt {
		case TagValueOmitEmpty:
			omitempty = true
		case TagValueString:
			quote = true
		}
		opts = rest
	}
	return omitempty, quote
}

// marshalJson marshals the value represented by v to raw json.
func marshalJson(v reflect.Value, quote bool) (json.RawMessage, error) {
	if !v.IsValid() {
		return NullJson, nil
	}
	jsonBts, err := json.Marshal(v.Interface())
	if err != nil {
		return nil, err
	}
	if quote && quotable(v.Kind()) {
		jsonBts = []byte("\"" + string(jsonBts) + "\"")
	}
	return json.RawMessage(jsonBts), nil
}

// unmarshalJson unmarshals the raw json into a variable of the appropriate type
// and the sets this value in v.
func unmarshalJson(data json.RawMessage, v reflect.Value, quote bool) error {
	if len(data) == 0 {
		return nil
	}

	if quote && quotable(v.Kind()) {
		data = data[1 : len(data)-1]
	}

	for v.Kind() == reflect.Pointer {
		v = v.Elem()
	}

	if !v.CanAddr() {
		return fmt.Errorf("unaddressable value")
	}

	switch v.Type().Kind() {
	case reflect.Bool:
		var b bool
		if err := json.Unmarshal(data, &b); err != nil {
			return err
		}
		v.SetBool(b)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var i int64
		if err := json.Unmarshal(data, &i); err != nil {
			return err
		}
		v.SetInt(i)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		var u uint64
		if err := json.Unmarshal(data, &u); err != nil {
			return err
		}
		v.SetUint(u)
	case reflect.Float32, reflect.Float64:
		var f float64
		if err := json.Unmarshal(data, &f); err != nil {
			return err
		}
		v.SetFloat(f)
	case reflect.String:
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
		v.SetString(s)
	case reflect.Struct, reflect.Array, reflect.Slice, reflect.Map:
		var s = reflect.New(v.Type()).Interface()
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
		v.Set(reflect.ValueOf(s).Elem())
	case reflect.Interface:
		// if the interface has been initialised, unmarshal
		// into the supplied value
		e := v.Elem()
		var s any
		if e.IsValid() {
			s = reflect.New(e.Type()).Interface()
		} else {
			s = reflect.New(v.Type()).Interface()
		}
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
		v.Set(reflect.ValueOf(s).Elem())
	default:
		return &UnsupportedTypeErr{Kind: v.Type().Kind()}
	}

	return nil
}

// quotable retuns true iff the kind can be converted to or
// from a string by wrapping or unwrapping in quotes. Currently
// only numeric kinds are supported.
func quotable(k reflect.Kind) bool {
	switch k {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return true
	default:
		return false
	}
}

// isEmpty returns true iff the value is should be omitted
// when the omitempty flag is set, ie if it is not valid,
// zero, or an empty array, slice or map.
// NB assumes that the input has been derefernced eg with
// derefValue.
func isEmpty(v reflect.Value) bool {
	if !v.IsValid() || v.IsZero() {
		return true
	}

	switch v.Kind() {
	case reflect.Array, reflect.Slice, reflect.Map:
		return v.Len() == 0
	default:
		return false
	}
}

// derefInput returns either:
// - the underlying value of v, found by following all pointers, or
// - an instance of type t, if one of the dereferenced values implements it.
// An error is returned if a loop of self-referential pointers is found.
func derefInput(v reflect.Value, t reflect.Type) (reflect.Value, error) {
	path := map[unsafe.Pointer]bool{}
	for {
		if v.Type().Implements(t) || (v.Kind() != reflect.Pointer && v.Kind() != reflect.Interface) {
			return v, nil
		}

		// check for a loop of self-referential pointers
		if v.Kind() == reflect.Pointer {
			ptr := v.UnsafePointer()
			if path[ptr] {
				return reflect.Value{}, ErrSelfRefPtr
			}
			path[ptr] = true
		}

		v = v.Elem()
	}
}

// fieldByIndex returns the value found by following the nested
// struct fields defined by the supplied indexes.
// It assumes that every value on the path is either a struct
// or a pointer to a struct.
func fieldByIndex(v reflect.Value, idxs []int) (reflect.Value, error) {
	var err error
	for _, idx := range idxs {
		v, err = derefValue(v)
		if err != nil {
			return reflect.Value{}, err
		}

		v = v.Field(idx)
	}
	return v, nil
}

// initFieldByIndex takes a value v and an array of indexes idxs, and
// initialises the struct field found in v at index idxs[0], then the
// struct field found at idxs[1] in the newly intialised struct, etc.
// All are initialised to their zero values.
// It assumes that all but the final value on the idx path is either a struct
// or a pointer to a struct.
func initFieldByIndex(v reflect.Value, idxs []int) (reflect.Value, error) {
	var err error
	for _, idx := range idxs {
		v, err = derefValue(v)
		if err != nil {
			return reflect.Value{}, err
		}

		v = v.Field(idx)
		initValue(v)
	}
	return v, nil
}

// initValue initialises v's underlying value, found by following
// all pointers, to its zero value.
// Any required pointers will also be initialised.
// NB if a self-referential pointer loop is discovered, the
// function returns immediately without error.
func initValue(v reflect.Value) {
	path := map[unsafe.Pointer]bool{}
	for {
		if v.Kind() == reflect.Interface {
			v = v.Elem()
			continue
		}
		if v.Kind() != reflect.Pointer {
			return
		}

		// Prevent infinite loop if v is an interface pointing to its own address:
		// type t struct {
		//     v any
		// }
		// s := t{}
		// s.v = &s.v
		ptr := v.UnsafePointer()
		if path[ptr] {
			return
		}
		path[ptr] = true

		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		v = v.Elem()
	}
}

// derefValue returns the value of v after following all pointers,
// or an error if a cycle of pointers is detected.
func derefValue(v reflect.Value) (reflect.Value, error) {
	path := map[unsafe.Pointer]bool{}
	for {
		if v.Kind() != reflect.Pointer && v.Kind() != reflect.Interface {
			return v, nil
		}

		// check for a loop of self-referential pointers
		if v.Kind() == reflect.Pointer {
			ptr := v.UnsafePointer()
			if path[ptr] {
				return reflect.Value{}, ErrSelfRefPtr
			}
			path[ptr] = true
		}

		v = v.Elem()
	}
}

// derefType returns the type of t after following all pointers.
func derefType(t reflect.Type) reflect.Type {
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	return t
}
