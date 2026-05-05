package maps

import (
	"github.com/egsam98/ecto"
)

// Min restricts map length with a lower inclusive bound
func Min[M ~map[K]V, K comparable, V any](length uint) ecto.Test[M] {
	return ecto.Test[M]{
		Error: ecto.Errorf("must contain at least %d entries", length),
		Func:  func(v *M) bool { return len(*v) >= int(length) },
	}
}

// Max restricts map length with an upper inclusive bound
func Max[M ~map[K]V, K comparable, V any](length uint) ecto.Test[M] {
	return ecto.Test[M]{
		Error: ecto.Errorf("must contain at most %d entries", length),
		Func:  func(v *M) bool { return len(*v) <= int(length) },
	}
}
