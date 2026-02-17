package ecto

import (
	"reflect"

	"github.com/egsam98/errors"
	"github.com/samber/lo"
)

// Schema common interface to process input data (conversions, validations etc.)
type Schema interface {
	ForType() reflect.Type
	process(ptr any) error
}

type IAtomicOrPtrSchema interface {
	Schema
	IsRequired() bool
	WithRequired(value bool) IAtomicOrPtrSchema
}

type ISliceSchema interface {
	Schema
	Inner() Schema
	WithInner(inner Schema) ISliceSchema
}

type IStructSchema interface {
	Schema
	Fields() M
	WithFields(M) IStructSchema
	CastToAny(src []byte, unmarshal func([]byte, any) error, opts ...CastOpt) (any, error)
	Meta() map[string]FieldMeta
}

// Test holds predicate function to apply on validated data and returns Error in case of failure
type Test[T any] struct {
	Error Error
	Func  func(v *T) bool
}

// Run applies predicate
func (t *Test[T]) Run(ptr *T) *Error {
	if t.Func(ptr) {
		return nil
	}
	return &t.Error
}

// OneOf restricts value to limited variants
func OneOf[T comparable](variants ...T) Test[T] {
	set := lo.Keyify(variants)
	return Test[T]{
		Error: Errorf("must be one of %v", variants),
		Func:  func(v *T) bool { return lo.HasKey(set, *v) },
	}
}

// ScrubAny replaces "" with nil in *string recursively of any Go type
func ScrubAny(ptr any) {
	scrub(reflect.ValueOf(ptr))
}

// ScrubString replaces "" with nil of *string. Simplified non-reflect version of ScrubAny
func ScrubString[S ~string](ptr **S) {
	if *ptr != nil && **ptr == "" {
		*ptr = nil
	}
}

func scrub(rv reflect.Value) {
	switch rv.Kind() {
	case reflect.Ptr:
		rvElem := rv.Elem()
		if rvElem.Kind() == reflect.String && rvElem.IsZero() {
			if rv.CanSet() {
				rv.SetZero()
			}
			return
		}
		scrub(rvElem)
	case reflect.Struct:
		for i := range rv.NumField() {
			scrub(rv.Field(i))
		}
	case reflect.Array, reflect.Slice:
		for i := range rv.Len() {
			scrub(rv.Index(i))
		}
	case reflect.Map:
		it := rv.MapRange()
		for it.Next() {
			v := it.Value()

			var vElem reflect.Value
			if v.Kind() == reflect.Ptr {
				vElem = v.Elem()
			}

			// Cannot continue via recursion 'cause map values aren't addressable
			if vElem.Kind() == reflect.String && vElem.IsZero() {
				rv.SetMapIndex(it.Key(), reflect.Zero(v.Type()))
				continue
			}
			scrub(vElem)
		}
	default:
	}
}

func validateSchema(typ reflect.Type, schema Schema) error {
	if innerType := schema.ForType(); typ != innerType {
		return errors.Errorf("schema must have inner type %s, got %s. Probably you'd like to use `AtomicFrom()`"+
			" schema (ex. IntFrom[T])", typ, innerType)
	}
	return nil
}
