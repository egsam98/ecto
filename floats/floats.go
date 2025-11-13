package floats

import (
	"strconv"
	"strings"

	"github.com/egsam98/ecto"
)

// Min restricts value with lower inclusive bound
func Min(value float64) ecto.Test[float64] {
	return ecto.Test[float64]{
		Error: ecto.Errorf("must be %s minimum", strconv.FormatFloat(value, 'f', -1, 64)),
		Func:  func(v *float64) bool { return *v >= value },
	}
}

// Max restricts value with upper inclusive bound
func Max(value float64) ecto.Test[float64] {
	return ecto.Test[float64]{
		Error: ecto.Errorf("must be %s maximum", strconv.FormatFloat(value, 'f', -1, 64)),
		Func:  func(v *float64) bool { return *v <= value },
	}
}

// MaxPrecision restricts precision/scale/fractional number of digits
func MaxPrecision(value uint) ecto.Test[float64] {
	return ecto.Test[float64]{
		Error: ecto.Errorf("has more than %d precision digits", value),
		Func: func(v *float64) bool {
			_, prec, _ := strings.Cut(strconv.FormatFloat(*v, 'f', -1, 64), ".")
			return uint(len(prec)) <= value
		},
	}
}
