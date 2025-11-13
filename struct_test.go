package ecto_test

import (
	"regexp"
	"testing"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/egsam98/ecto"
	ectosl "github.com/egsam98/ecto/slices"
	ectos "github.com/egsam98/ecto/strings"
)

var structSchema = ecto.Struct[Data](ecto.M{
	"A": ecto.StringFrom[Hello]().
		Test(ectos.Regex(regexp.MustCompile("a")), ectos.Min(1)),
	"B": ecto.String().Required(),
	"C": ecto.Struct[C](map[string]ecto.Schema{
		"C1": ecto.String().Required().Test(ectos.Min(1)),
	}),
	"D": ecto.Slice[*string](
		ecto.Ptr[string](ecto.String().Required().Test(ectos.URL())),
	).Test(ectosl.Min[*string](2)),
	"E": ecto.Atomic[uuid.UUID]().Default(uuid.New()),
	"F": ecto.Slice[F](
		ecto.Struct[F](ecto.M{
			"F1": ecto.String().Required(),
		}),
	).Test(ectosl.Min[F](1)),
	"G": ecto.Ptr[int](ecto.Int().Required()).Required(),
})

type Hello string

type Data struct {
	A Hello  `validate:"required"`
	B string `validate:"required"`
	C C
	D []*string `validate:"required,min=2,dive,omitnil,url"`
	E uuid.UUID `validate:"required"`
	F []F       `json:"f" validate:"min=1,dive"`
	G *int      `validate:"omitnil,required"`
}

type C struct {
	C1 string `validate:"required"`
}

type F struct {
	F1 string `json:"f1" validate:"required"`
}

var data = Data{
	A: "a",
	B: "b",
	C: struct {
		C1 string `validate:"required"`
	}{C1: "c1"},
	D: []*string{lo.ToPtr("http://wikipedia.org"), lo.ToPtr("http://wikipedia.org")},
	E: uuid.New(),
	F: []F{{F1: "f1"}},
	G: lo.ToPtr(-1),
}

func TestStruct(t *testing.T) {
	ecto.Struct[Data](nil)
	ecto.Struct[Data](ecto.M{"B": ecto.String()})
	assert.Panics(t, func() { ecto.Struct[int](ecto.M{"B": ecto.String()}) })
	assert.Panics(t, func() { ecto.Struct[Data](ecto.M{"unknown": ecto.Int()}) })
}

func TestStructSchema_Process(t *testing.T) {
	assert.NoError(t, structSchema.Process(&data))

	t.Run("invalid F", func(t *testing.T) {
		data := data
		data.F = []F{}
		assert.EqualError(t, structSchema.Process(&data), `{"f":["must contain at least 1 items"]}`)
	})
}

func BenchmarkEcto(b *testing.B) {
	for i := 0; i < b.N; i++ {
		require.NoError(b, structSchema.Process(&data))
	}
}
