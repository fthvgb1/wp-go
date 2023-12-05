package slice

import (
	"golang.org/x/exp/constraints"
	"sort"
)

const (
	ASC = iota
	DESC
)

type anyArr[T any] struct {
	data []T
	fn   func(i, j T) bool
}

func (r anyArr[T]) Len() int {
	return len(r.data)
}

func (r anyArr[T]) Swap(i, j int) {
	r.data[i], r.data[j] = r.data[j], r.data[i]
}

func (r anyArr[T]) Less(i, j int) bool {
	return r.fn(r.data[i], r.data[j])
}

// Sort fn i>j 为降序 desc，反之为升序 asc
func Sort[T any](arr []T, fn func(i, j T) bool) {
	slice := anyArr[T]{
		data: arr,
		fn:   fn,
	}
	sort.Sort(slice)
}

func StableSort[T any](arr []T, fn func(i, j T) bool) {
	slice := anyArr[T]{
		data: arr,
		fn:   fn,
	}
	sort.Stable(slice)
}

func Sorts[T constraints.Ordered](a []T, order int) {
	slice := anyArr[T]{
		data: a,
		fn: func(i, j T) bool {
			if order == DESC {
				return i > j
			}
			return i < j
		},
	}
	sort.Sort(slice)
}

func SimpleSort[T any, O constraints.Ordered](a []T, order int, fn func(t T) O) {
	slice := anyArr[T]{
		data: a,
		fn: func(i, j T) bool {
			if order == DESC {
				return fn(i) > fn(j)
			}
			return fn(i) < fn(j)
		},
	}
	sort.Sort(slice)
}
