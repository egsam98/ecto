package slices

import (
	"github.com/samber/lo"

	"github.com/egsam98/ecto"
)

// Min restricts slice length with a lower inclusive bound
func Min[T any](length uint) ecto.Test[[]T] {
	return ecto.Test[[]T]{
		Error: ecto.Errorf("must contain at least %d items", length),
		Func:  func(v *[]T) bool { return len(*v) >= int(length) },
	}
}

// Max restricts slice length with an upper inclusive bound
func Max[T any](length uint) ecto.Test[[]T] {
	return ecto.Test[[]T]{
		Error: ecto.Errorf("must contain at most %d items", length),
		Func:  func(v *[]T) bool { return len(*v) <= int(length) },
	}
}

// Unique makes sure to have all slice elements unique
func Unique[T comparable]() ecto.Test[[]T] { return UniqueBy(func(t T) T { return t }) }

// UniqueBy makes sure to have all slice elements unique by key function applied for every element
func UniqueBy[T any, K comparable](key func(T) K) ecto.Test[[]T] {
	return ecto.Test[[]T]{
		Error: "items must be unique",
		Func: func(v *[]T) bool {
			uniques := lo.Associate(*v, func(elem T) (K, struct{}) { return key(elem), struct{}{} })
			return len(*v) == len(uniques)
		},
	}
}
