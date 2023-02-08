package slice

import "sort"

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

func SortSelf[T any](arr []T, fn func(i, j T) bool) {
	slice := anyArr[T]{
		data: arr,
		fn:   fn,
	}
	sort.Sort(slice)
	return
}

func Sort[T any](arr []T, fn func(i, j T) bool) (r []T) {
	r = make([]T, len(arr))
	copy(r, arr)
	slice := anyArr[T]{
		data: r,
		fn:   fn,
	}
	sort.Sort(slice)
	return
}
