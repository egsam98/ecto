package ecto

import (
	"cmp"
	"reflect"
	"strings"

	"github.com/egsam98/errors"
	"github.com/samber/lo"
)

var _ Schema = (*StructSchema[any])(nil)

// StructSchema represents schema for struct types via hashmap as a struct field to its schema
type StructSchema[T any] struct {
	Fields M
	typ    reflect.Type
	meta   map[string]fieldMeta
}

type fieldMeta struct {
	index int
	tag   string
}

type M = map[string]Schema

func Struct[T any](fields M) StructSchema[T] {
	typ := reflect.TypeFor[T]()
	if typ.Kind() != reflect.Struct {
		panic(errors.Errorf("%s: not a struct", typ))
	}

	meta := make(map[string]fieldMeta)
	for i := range typ.NumField() {
		field := typ.Field(i)
		tag, _, _ := strings.Cut(field.Tag.Get("json"), ",")
		meta[field.Name] = fieldMeta{
			index: i,
			tag:   cmp.Or(tag, field.Name),
		}
	}
	self := StructSchema[T]{Fields: fields, typ: typ, meta: meta}

	for key, schema := range fields {
		field, ok := typ.FieldByName(key)
		if !ok {
			self.panicMissingKey(key)
		}
		if err := validateSchema(field.Type, schema); err != nil {
			panic(errors.Wrapf(err, "%T: %s", self, key))
		}
	}
	return self
}

// Process may return MapError
func (s StructSchema[T]) Process(ptr *T) error { return s.process(ptr) }

// Cast deserializes bytes into type and runs Process.
func (s StructSchema[T]) Cast(src []byte, deserialize func([]byte, any) error, opts ...CastOpt) (T, error) {
	var data T
	if err := deserialize(src, &data); err != nil {
		return data, errors.Wrapf(err, "deserialize %s into %T", src, data)
	}

	var cfg castConfig
	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}
	if cfg.scrub {
		scrub(reflect.ValueOf(&data))
	}

	return data, s.Process(&data)
}

// Extend existing schema
func (s StructSchema[T]) Extend(fields M) StructSchema[T] {
	return Struct[T](lo.Assign(s.Fields, fields))
}

func (s StructSchema[T]) process(ptrStruct any) error {
	rv := reflect.ValueOf(ptrStruct).Elem()
	if rv.Type() != s.typ {
		panic(errors.Errorf("expected type %s to process, got: %s", s.typ, rv.Type()))
	}

	var errs MapError
	for key, schema := range s.Fields {
		keyMeta, ok := s.meta[key]
		if !ok {
			s.panicMissingKey(key)
		}

		ptr := rv.Field(keyMeta.index).Addr().Interface()
		if err := schema.process(ptr); err != nil {
			errs.Add(keyMeta.tag, err)
		}
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}

func (s StructSchema[T]) forType() reflect.Type { return s.typ }

func (s StructSchema[T]) panicMissingKey(key string) {
	panic(errors.Errorf("%T: Missing struct schema key: %s", s, key))
}

type CastOpt func(*castConfig)

// Scrub replaces empty strings with nil for every `*string` type
func Scrub() CastOpt {
	return func(cfg *castConfig) { cfg.scrub = true }
}

type castConfig struct {
	scrub bool
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
