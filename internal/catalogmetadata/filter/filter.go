package filter

import (
	"github.com/operator-framework/operator-controller/internal/catalogmetadata"
)

// Predicate returns true if the object should be kept when filtering
type Predicate[T catalogmetadata.Schemas] func(entity *T) bool

// Filter filters a slice accordingly to
func Filter[T catalogmetadata.Schemas](in []*T, test Predicate[T]) []*T {
	out := []*T{}
	for i := range in {
		if test(in[i]) {
			out = append(out, in[i])
		}
	}
	return out
}

func And[T catalogmetadata.Schemas](predicates ...Predicate[T]) Predicate[T] {
	return func(obj *T) bool {
		for _, predicate := range predicates {
			if !predicate(obj) {
				return false
			}
		}
		return true
	}
}

func Or[T catalogmetadata.Schemas](predicates ...Predicate[T]) Predicate[T] {
	return func(obj *T) bool {
		for _, predicate := range predicates {
			if predicate(obj) {
				return true
			}
		}
		return false
	}
}

func Not[T catalogmetadata.Schemas](predicate Predicate[T]) Predicate[T] {
	return func(obj *T) bool {
		return !predicate(obj)
	}
}
