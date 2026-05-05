package ecto_test

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	"github.com/egsam98/ecto"
	ectosl "github.com/egsam98/ecto/slices"
	ectos "github.com/egsam98/ecto/strings"
)

func TestSlice(t *testing.T) {
	ecto.Slice[[]string](ecto.String())
	ecto.Slice[[]decimal.Decimal](ecto.StringFrom[decimal.Decimal]())
	assert.Panics(t, func() { ecto.Slice[[]string](ecto.Int()) })
}

func TestSlice_Process(t *testing.T) {
	schema := ecto.Slice[[]string](
		ecto.String().Test(ectos.URL()),
	).Test(ectosl.Min[[]string](2))

	assert.NoError(t, schema.Process([]string{"http://wikipedia.org", "http://example.com"}))
	assert.EqualError(t, schema.Process([]string{""}), `["must contain at least 2 items"]`)
	assert.EqualError(t, schema.Process([]string{"test", "http://example.com", ""}),
		`{"0":["invalid URL"],"2":["invalid URL"]}`)
}
