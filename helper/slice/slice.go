package slice

import (
	"github.com/fthvgb1/wp-go/helper"
	"github.com/fthvgb1/wp-go/helper/number"
)

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

func Walk[T any](arr []T, fn func(*T)) {
	for i := 0; i < len(arr); i++ {
		fn(&arr[i])
	}
}

func SearchFirst[T any](arr []T, fn func(T) bool) (int, T) {
	for i, t := range arr {
		if fn(t) {
			return i, t
		}
	}
	var r T
	return -1, r
}

func SearchLast[T any](arr []T, fn func(T) bool) (int, T) {
	for i := len(arr) - 1; i > 0; i-- {
		if fn(arr[i]) {
			return i, arr[i]
		}
	}
	var r T
	return -1, r
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

// Fill 用指定值填充一个切片
func Fill[T any](start, len int, v T) []T {
	r := make([]T, start+len)
	for i := 0; i < len; i++ {
		r[start+i] = v
	}
	return r
}

// Pad 以指定长度将一个值填充进切片 returns a copy of the array padded to size specified by length with value. If length is positive then the array is padded on the right, if it's negative then on the left. If the absolute value of length is less than or equal to the length of the array then no padding takes place.
func Pad[T any](a []T, length int, v T) []T {
	l := len(a)
	if length > l {
		return append(a, Fill(0, length-l, v)...)
	} else if length < 0 && -length > l-1 {
		return append(Fill(0, -length-2, v), a...)
	}
	return a
}

// Pop 弹出最后一个元素
func Pop[T any](a *[]T) T {
	arr := *a
	v := arr[len(arr)-1]

	*a = append(arr[:len(arr)-1])
	return v
}

// Rand 随机取一个元素
func Rand[T any](a []T) (int, T) {
	i := number.Rand(0, len(a)-1)
	return i, a[i]
}

// RandPop 随机弹出一个元素并返回那个剩余长度
func RandPop[T any](a *[]T) (T, int) {
	arr := *a
	if len(arr) == 0 {
		var r T
		return r, 0
	}
	i := number.Rand(0, len(arr)-1)
	v := arr[i]
	if len(arr)-1 == i {
		*a = append(arr[:i])
	} else {
		*a = append(arr[:i], arr[i+1:]...)
	}
	return v, len(arr) - 1
}

// Shift 将切片的第一个单元移出并作为结果返回
func Shift[T any](a *[]T) (T, int) {
	arr := *a
	l := len(arr)
	if l > 0 {
		v := arr[0]
		*a = arr[1:]
		return v, l - 1
	}
	var r T
	return r, 0
}
