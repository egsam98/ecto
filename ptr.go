package ecto

import (
	"reflect"

	"github.com/pkg/errors"
)

var _ Schema = (*PtrSchema[any])(nil)
var _ IAtomicOrPtrSchema = (*PtrSchema[any])(nil)

// PtrSchema wraps inner Schema assuming input data as pointer. Features:
// - Mark a pointer as required (non-nil)
// - Process internal schema
type PtrSchema[To any] struct {
	inner    Schema
	required bool
}

func Ptr[To any](inner Schema) PtrSchema[To] {
	self := PtrSchema[To]{inner: inner}

	if err := validateSchema(reflect.TypeFor[To](), inner); err != nil {
		panic(errors.Wrapf(err, "%T", self))
	}
	return self
}

// Process may return ListError
func (s PtrSchema[T]) Process(data *T) error { return s.process(&data) }

func (s PtrSchema[T]) Required() PtrSchema[T] {
	s.required = true
	return s
}

func (s PtrSchema[T]) process(ptrAny any) error {
	ptr := *ptrAny.(**T)
	if ptr == nil {
		if s.required {
			return ListError{errRequired}
		}
		return nil
	}
	return s.inner.process(ptr)
}

func (s PtrSchema[T]) ForType() reflect.Type { return reflect.TypeFor[*T]() }

func (s PtrSchema[T]) IsRequired() bool { return s.required }

func (s PtrSchema[To]) WithRequired(value bool) IAtomicOrPtrSchema {
	s.required = value
	return s
}
