package ecto

import (
	"reflect"
	"strconv"

	"github.com/egsam98/errors"
)

var _ Schema = (*SliceSchema[[]any, any])(nil)
var _ ISliceSchema = (*SliceSchema[[]any, any])(nil)

// SliceSchema wraps inner Schema assuming input data as slice. Features:
// - Run slice-specific tests (see ecto/slices subpackage)
// - Process internal schema
type SliceSchema[S ~[]T, T any] struct {
	inner Schema
	tests []Test[S]
}

func Slice[S ~[]T, T any](inner Schema) SliceSchema[S, T] {
	self := SliceSchema[S, T]{inner: inner}

	if err := validateSchema(reflect.TypeFor[T](), inner); err != nil {
		panic(errors.Wrapf(err, "%T", self))
	}
	return self
}

func (s SliceSchema[S, T]) Test(tests ...Test[S]) SliceSchema[S, T] {
	s.tests = tests
	return s
}

// Process may return ListError (for list tests) or MapError for individual element errors.
// Map key is a stringified slice index
func (s SliceSchema[S, T]) Process(data []T) error { return s.process(&data) }

func (s SliceSchema[S, T]) process(ptrAny any) error {
	ptr := ptrAny.(*S)

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

func (SliceSchema[S, T]) ForType() reflect.Type { return reflect.TypeFor[S]() }

func (s SliceSchema[S, T]) Inner() Schema { return s.inner }

func (s SliceSchema[S, T]) WithInner(inner Schema) ISliceSchema {
	s.inner = inner
	return s
}
