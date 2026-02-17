package ecto

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/egsam98/errors"
	"github.com/samber/lo"
	"golang.org/x/exp/constraints"
)

var _ Schema = (*AtomicSchema[any, any])(nil)
var _ IAtomicOrPtrSchema = (*AtomicSchema[any, any])(nil)
var typeStringer = reflect.TypeFor[fmt.Stringer]()
var typeJsonNumber = reflect.TypeFor[json.Number]()

// AtomicSchema is a schema designed to describe scalar Go types or those that do not need to be recursively
// processed internally (ex. decimal.Decimal).
// Features:
// - Convert input data into another (from T to R) for further validations
// - Mark as required (check if value is Go zero-value)
// - Set default if value is zero
// - Set of Test predicates for validation
type AtomicSchema[T comparable, R any] struct {
	required, omitZero bool
	defaultValue       *T
	convert            func(*T) (*R, error)
	tests              []Test[R]
}

func Atomic[T comparable]() AtomicSchema[T, T] {
	return AtomicFrom(func(t *T) (*T, error) { return t, nil })
}

func AtomicFrom[T comparable, R any](convert func(*T) (*R, error)) AtomicSchema[T, R] {
	return AtomicSchema[T, R]{convert: convert}
}

// Int shorthand
func Int() AtomicSchema[int, int] { return Atomic[int]() }

// IntFrom shorthand with conversion
func IntFrom[T constraints.Integer]() AtomicSchema[T, int] {
	return AtomicFrom[T, int](func(t *T) (*int, error) { return lo.ToPtr(int(*t)), nil })
}

// String shorthand
func String() AtomicSchema[string, string] { return Atomic[string]() }

// StringFrom shorthand with conversion. Supported types in order:
// - fmt.Stringer
// - string and its type definitions
func StringFrom[T comparable]() AtomicSchema[T, string] {
	var convert func(*T) (*string, error)
	switch rt := reflect.TypeFor[T](); {
	case rt.Implements(typeStringer):
		convert = func(v *T) (*string, error) {
			return lo.ToPtr(any(*v).(fmt.Stringer).String()), nil
		}
	case rt.Kind() == reflect.String:
		convert = func(v *T) (*string, error) {
			return lo.ToPtr(reflect.ValueOf(*v).String()), nil
		}
	default:
		panic(errors.Errorf("%s is neither string nor %s", rt, typeStringer))
	}

	return AtomicFrom[T, string](convert)
}

// Float shorthand
func Float() AtomicSchema[float64, float64] { return Atomic[float64]() }

// FloatFrom shorthand with conversion. Supported types in order:
// - float32/64 and its type definitions
// - json.Number
func FloatFrom[T comparable]() AtomicSchema[T, float64] {
	rt := reflect.TypeFor[T]()

	var convert func(*T) (*float64, error)
	switch kind := rt.Kind(); {
	case kind == reflect.Float32 || kind == reflect.Float64:
		convert = func(v *T) (*float64, error) {
			return lo.ToPtr(reflect.ValueOf(*v).Float()), nil
		}
	case rt == typeJsonNumber:
		convert = func(t *T) (*float64, error) {
			f, err := any(*t).(json.Number).Float64()
			if err != nil {
				return nil, errors.New("invalid number")
			}
			return &f, nil
		}
	default:
		panic(errors.Errorf("%s is neither float32/64 nor %s", rt, typeJsonNumber))
	}

	return AtomicFrom(convert)
}

func (s AtomicSchema[T, R]) Required() AtomicSchema[T, R] {
	s.required = true
	return s
}

// OmitZero
// Deprecated: for optional parameters use PtrSchema, this method is for backward compatibility
func (s AtomicSchema[T, R]) OmitZero() AtomicSchema[T, R] {
	s.omitZero = true
	return s
}

func (s AtomicSchema[T, R]) Default(value T) AtomicSchema[T, R] {
	s.defaultValue = &value
	return s
}

func (s AtomicSchema[T, R]) Test(tests ...Test[R]) AtomicSchema[T, R] {
	s.tests = tests
	return s
}

// Process may return ListError
func (s AtomicSchema[T, R]) Process(data *T) error { return s.process(data) }

func (AtomicSchema[T, R]) ForType() reflect.Type { return reflect.TypeFor[T]() }

func (s AtomicSchema[T, R]) process(ptrAny any) error {
	ptr := ptrAny.(*T)
	if lo.IsEmpty(*ptr) {
		if s.required {
			return ListError{errRequired}
		}
		if s.omitZero {
			return nil
		}
		if s.defaultValue != nil {
			*ptr = *s.defaultValue
		}
	}

	ptrConv, err := s.convert(ptr)
	if err != nil {
		return ListError{Error(err.Error())}
	}

	var errs ListError
	for _, test := range s.tests {
		if err := test.Run(ptrConv); err != nil {
			errs = append(errs, *err)
		}
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}

func (s AtomicSchema[T, R]) IsRequired() bool { return s.required || (!s.omitZero && len(s.tests) > 0) }

func (s AtomicSchema[T, R]) WithRequired(value bool) IAtomicOrPtrSchema {
	s.required = value
	s.omitZero = !value
	return s
}
