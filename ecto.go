package ecto

import (
	"reflect"

	"github.com/egsam98/errors"
	"github.com/samber/lo"
)

// Schema common interface to process input data (conversions, validations etc.)
type Schema interface {
	process(ptr any) error
	forType() reflect.Type
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

func validateSchema(typ reflect.Type, schema Schema) error {
	if innerType := schema.forType(); typ != innerType {
		return errors.Errorf("schema must have inner type %s, got %s. Probably you'd like to use `AtomicFrom()`"+
			" schema (ex. IntFrom[T])", typ, innerType)
	}
	return nil
}
