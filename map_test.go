package ecto_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/egsam98/ecto"
	ectom "github.com/egsam98/ecto/maps"
)

func TestMap_Process(t *testing.T) {
	type KV = map[string]int
	schema := ecto.Map[KV]().Test(ectom.Min[KV](2))

	assert.NoError(t, schema.Process(KV{"a": 1, "b": 2}))
	assert.EqualError(t, schema.Process(KV{"b": 1}), `["must contain at least 2 entries"]`)
}
