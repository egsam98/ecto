package ecto_test

import (
	"strconv"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"

	"github.com/egsam98/ecto"
)

func TestAtomic(t *testing.T) {
	schema := ecto.Atomic[int]()

	assert.NoError(t, schema.Process(lo.ToPtr(0)))
	assert.Panics(t, func() { _ = schema.Process(nil) })

	t.Run("required", func(t *testing.T) {
		schema := ecto.Atomic[int]().Required()
		assert.EqualError(t, schema.Process(lo.ToPtr(0)), `["required"]`)
	})

	t.Run("default", func(t *testing.T) {
		schema := ecto.Atomic[int]().Default(5)
		v := 0
		assert.NoError(t, schema.Process(&v))
		assert.Equal(t, 5, v)
		v = 1
		assert.NoError(t, schema.Process(&v))
		assert.Equal(t, 1, v)
	})
}

func TestAtomicFrom(t *testing.T) {
	schema := ecto.AtomicFrom[string, int](func(s *string) (*int, error) {
		i, err := strconv.Atoi(*s)
		return &i, err
	})

	assert.NoError(t, schema.Process(lo.ToPtr("1")))
	assert.EqualError(t, schema.Process(lo.ToPtr("a")), `["strconv.Atoi: parsing \"a\": invalid syntax"]`)
}
