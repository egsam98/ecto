package ecto

import (
	"cmp"
	"encoding/json"
	"reflect"
	"strings"

	"github.com/egsam98/errors"
	"github.com/samber/lo"
)

var _ Schema = (*StructSchema[any])(nil)
var _ IStructSchema = (*StructSchema[any])(nil)

// StructSchema represents schema for struct types via hashmap as a struct field to its schema
type StructSchema[T any] struct {
	fields M
	meta   map[string]FieldMeta
}

type FieldMeta struct {
	Index int
	Tag   string
}

type M = map[string]Schema

func Struct[T any](fields M) StructSchema[T] {
	typ := reflect.TypeFor[T]()
	if typ.Kind() != reflect.Struct {
		panic(errors.Errorf("%s: not a struct", typ))
	}

	self := StructSchema[T]{fields: fields}
	self.makeMeta()
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
		ScrubAny(&data)
	}

	return data, s.Process(&data)
}

func (s StructSchema[T]) CastJSON(src []byte, opts ...CastOpt) (T, error) {
	return s.Cast(src, json.Unmarshal, opts...)
}

func (s StructSchema[T]) CastToAny(src []byte, deserialize func([]byte, any) error, opts ...CastOpt) (any, error) {
	return s.Cast(src, deserialize, opts...)
}

// Extend existing schema
func (s StructSchema[T]) Extend(fields M) StructSchema[T] {
	return Struct[T](lo.Assign(s.fields, fields))
}

func (s StructSchema[T]) ForType() reflect.Type { return reflect.TypeFor[T]() }

func (s StructSchema[T]) process(ptrStruct any) error {
	if len(s.fields) == 0 {
		return nil
	}

	rv := reflect.ValueOf(ptrStruct).Elem()
	var errs MapError
	for key, schema := range s.fields {
		keyMeta, ok := s.meta[key]
		if !ok {
			s.panicMissingKey(key)
		}

		ptr := rv.Field(keyMeta.Index).Addr().Interface()
		if err := schema.process(ptr); err != nil {
			errs.Add(keyMeta.Tag, err)
		}
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}

func (s StructSchema[T]) Fields() M { return s.fields }

func (s StructSchema[T]) WithFields(fields M) IStructSchema {
	s.fields = fields
	s.makeMeta()
	return s
}

func (s StructSchema[T]) Meta() map[string]FieldMeta { return s.meta }

func (s *StructSchema[T]) makeMeta() {
	s.meta = make(map[string]FieldMeta)
	if len(s.fields) == 0 {
		return
	}

	typ := reflect.TypeFor[T]()
	for i := range typ.NumField() {
		field := typ.Field(i)
		tag, _, _ := strings.Cut(field.Tag.Get("json"), ",")
		s.meta[field.Name] = FieldMeta{
			Index: i,
			Tag:   cmp.Or(tag, field.Name),
		}
	}
	for key, schema := range s.fields {
		field, ok := typ.FieldByName(key)
		if !ok {
			s.panicMissingKey(key)
		}
		if err := validateSchema(field.Type, schema); err != nil {
			panic(errors.Wrapf(err, "%T: %s", s, key))
		}
	}
}

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
