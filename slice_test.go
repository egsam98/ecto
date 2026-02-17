package ecto_test

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	"github.com/egsam98/ecto"
	ectos "github.com/egsam98/ecto/strings"
)

func TestSlice(t *testing.T) {
	ecto.Slice[[]string](ecto.String())
	ecto.Slice[[]decimal.Decimal](ecto.StringFrom[decimal.Decimal]())
	assert.Panics(t, func() { ecto.Slice[[]string](ecto.Int()) })
}

func TestSlice_Process(t *testing.T) {
	schema := ecto.Slice[[]string](ecto.String().Test(ectos.URL()))
	assert.NoError(t, schema.Process([]string{"http://wikipedia.org"}))
	assert.EqualError(t, schema.Process([]string{"test", "http://wikipedia.org", ""}),
		`{"0":["invalid URL"],"2":["invalid URL"]}`)
}
