package integer

import (
	"github.com/egsam98/ecto"
)

// Eq forces a value to be equal to another
func Eq(value int) ecto.Test[int] {
	return ecto.Test[int]{
		Error: ecto.Errorf("must be equal to %d", value),
		Func:  func(v *int) bool { return *v == value },
	}
}

// Min restricts value with lower inclusive bound
func Min(value int) ecto.Test[int] {
	return ecto.Test[int]{
		Error: ecto.Errorf("must be %d minimum", value),
		Func:  func(v *int) bool { return *v >= value },
	}
}

// Max restricts value with upper inclusive bound
func Max(value int) ecto.Test[int] {
	return ecto.Test[int]{
		Error: ecto.Errorf("must be %d maximum", value),
		Func:  func(v *int) bool { return *v <= value },
	}
}
