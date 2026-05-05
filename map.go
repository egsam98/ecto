package ecto

import (
	"reflect"
)

var _ Schema = (*MapSchema[map[any]any, any, any])(nil)
var _ IMapSchema = (*MapSchema[map[any]any, any, any])(nil)

// MapSchema represents a generic map type. There's no subschema support for key/value yet.
type MapSchema[M ~map[K]V, K comparable, V any] struct {
	tests []Test[M]
}

func Map[M ~map[K]V, K comparable, V any]() MapSchema[M, K, V] {
	return MapSchema[M, K, V]{}
}

func (m MapSchema[M, K, V]) Test(tests ...Test[M]) MapSchema[M, K, V] {
	m.tests = append(m.tests, tests...)
	return m
}

func (m MapSchema[M, K, V]) ForType() reflect.Type { return reflect.TypeFor[M]() }

// Process may return ListError
func (s MapSchema[M, K, V]) Process(data M) error { return s.process(&data) }

func (m MapSchema[M, K, V]) process(ptrAny any) error {
	ptr := ptrAny.(*M)
	var errs ListError
	for _, test := range m.tests {
		if err := test.Run(ptr); err != nil {
			errs = append(errs, *err)
		}
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}

func (MapSchema[M, K, V]) implIMapSchema() {}
