package ecto

import (
	"reflect"
	"strconv"

	"github.com/egsam98/errors"
)

var _ Schema = (*SliceSchema[any])(nil)

// SliceSchema wraps inner Schema assuming input data as slice. Features:
// - Run slice-specific tests (see ecto/slices subpackage)
// - Process internal schema
type SliceSchema[T any] struct {
	inner Schema
	tests []Test[[]T]
}

func Slice[T any](inner Schema) SliceSchema[T] {
	self := SliceSchema[T]{inner: inner}

	if err := validateSchema(reflect.TypeFor[T](), inner); err != nil {
		panic(errors.Wrapf(err, "%T", self))
	}
	return self
}

func (s SliceSchema[T]) Test(tests ...Test[[]T]) SliceSchema[T] {
	s.tests = tests
	return s
}

// Process may return ListError (for list tests) or MapError for individual element errors.
// Map key is a stringified slice index
func (s SliceSchema[T]) Process(data []T) error { return s.process(&data) }

func (s SliceSchema[T]) process(ptrAny any) error {
	ptr := ptrAny.(*[]T)

	var errs ListError
	for _, test := range s.tests {
		if err := test.Run(ptr); err != nil {
			errs = append(errs, *err)
		}
	}
	if len(errs) > 0 {
		return errs
	}

	var innerErrs MapError
	for i, elem := range *ptr {
		if err := s.inner.process(&elem); err != nil {
			innerErrs.Add(strconv.Itoa(i), err)
		}
	}
	if len(innerErrs) > 0 {
		return innerErrs
	}
	return nil
}

func (SliceSchema[T]) forType() reflect.Type { return reflect.TypeFor[[]T]() }
