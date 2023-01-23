package slice

import "github.com/fthvgb1/wp-go/helper"

func Map[T, R any](arr []T, fn func(T) R) []R {
	r := make([]R, 0, len(arr))
	for _, t := range arr {
		r = append(r, fn(t))
	}
	return r
}

func FilterAndMap[N any, T any](arr []T, fn func(T) (N, bool)) (r []N) {
	for _, t := range arr {
		x, ok := fn(t)
		if ok {
			r = append(r, x)
		}
	}
	return
}

func Filter[T any](arr []T, fn func(T) bool) []T {
	var r []T
	for _, t := range arr {
		if fn(t) {
			r = append(r, t)
		}
	}
	return r
}

func Reduce[R, T any](arr []T, fn func(T, R) R, r R) R {
	for _, t := range arr {
		r = fn(t, r)
	}
	return r
}

func Reverse[T any](arr []T) []T {
	var r = make([]T, 0, len(arr))
	for i := len(arr); i > 0; i-- {
		r = append(r, arr[i-1])
	}
	return r
}

func ReverseSelf[T any](arr []T) []T {
	l := len(arr)
	half := l / 2
	for i := 0; i < half; i++ {
		arr[i], arr[l-i-1] = arr[l-i-1], arr[i]
	}
	return arr
}

func SimpleToMap[K comparable, V any](arr []V, fn func(V) K) map[K]V {
	return ToMap(arr, func(v V) (K, V) {
		return fn(v), v
	}, true)
}

func ToMap[K comparable, V, T any](arr []V, fn func(V) (K, T), isCoverPrev bool) map[K]T {
	m := make(map[K]T)
	for _, v := range arr {
		k, r := fn(v)
		if !isCoverPrev {
			if _, ok := m[k]; ok {
				continue
			}
		}
		m[k] = r
	}
	return m
}

func Pagination[T any](arr []T, page, pageSize int) []T {
	start := (page - 1) * pageSize
	l := len(arr)
	if start > l {
		start = l
	}
	end := page * pageSize
	if l < end {
		end = l
	}
	return arr[start:end]
}

func Chunk[T any](arr []T, size int) [][]T {
	var r [][]T
	i := 0
	for {
		if len(arr) <= size+i {
			r = append(r, arr[i:])
			break
		}
		r = append(r, arr[i:i+size])
		i += size
	}
	return r
}

func Slice[T any](arr []T, offset, length int) (r []T) {
	l := len(arr)
	if length == 0 {
		length = l - offset
	}
	if l > offset && l >= offset+length {
		r = append(make([]T, 0, length), arr[offset:offset+length]...)
		arr = append(arr[:offset], arr[offset+length:]...)
	} else if l <= offset {
		return
	} else if l > offset && l < offset+length {
		r = append(make([]T, 0, length), arr[offset:]...)
		arr = arr[:offset]
	}
	return
}

func FilterAndToMap[K comparable, V, T any](arr []T, fn func(T) (K, V, bool)) map[K]V {
	r := make(map[K]V)
	for _, t := range arr {
		k, v, ok := fn(t)
		if ok {
			r[k] = v
		}
	}
	return r
}

func Comb[T any](arr []T, m int) (r [][]T) {
	if m == 1 {
		for _, t := range arr {
			r = append(r, []T{t})
		}
		return
	}
	l := len(arr) - m
	for i := 0; i <= l; i++ {
		next := Slice(arr, i+1, 0)
		nexRes := Comb(next, m-1)
		for _, re := range nexRes {
			t := append([]T{arr[i]}, re...)
			r = append(r, t)
		}
	}
	return r
}

func GroupBy[K comparable, T, V any](a []T, fn func(T) (K, V)) map[K][]V {
	r := make(map[K][]V)
	for _, t := range a {
		k, v := fn(t)
		if _, ok := r[k]; !ok {
			r[k] = []V{v}
		} else {
			r[k] = append(r[k], v)
		}
	}
	return r
}

func ToAnySlice[T any](a []T) []any {
	return Map(a, helper.ToAny[T])
}
