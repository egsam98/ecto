package ecto_test

import (
	"encoding/json"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"

	"github.com/egsam98/ecto"
)

func TestPtr(t *testing.T) {
	ecto.Ptr[json.Number](ecto.StringFrom[json.Number]())
	assert.Panics(t, func() { ecto.Ptr[string](ecto.StringFrom[json.Number]()) })
}

func TestPtr_Process(t *testing.T) {
	schema := ecto.Ptr[int](ecto.Int().Required())
	assert.NoError(t, schema.Process(nil))
	assert.NoError(t, schema.Process(lo.ToPtr(1)))
	assert.EqualError(t, schema.Process(lo.ToPtr(0)), `["required"]`)

	t.Run("required", func(t *testing.T) {
		schema := ecto.Ptr[int](ecto.Int()).Required()
		assert.EqualError(t, schema.Process(nil), `["required"]`)
	})
}
