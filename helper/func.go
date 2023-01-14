package helper

import (
	"fmt"
	"math/rand"
	"reflect"
	"sort"
)

type IntNumber interface {
	~int | ~int64 | ~int32 | ~int8 | ~int16 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

func ToAny[T any](v T) any {
	return v
}

func IsContainInArr[T comparable](a T, arr []T) bool {
	for _, v := range arr {
		if a == v {
			return true
		}
	}
	return false
}

func StructColumn[T any, M any](arr []M, field string) (r []T) {
	for i := 0; i < len(arr); i++ {
		v := reflect.ValueOf(arr[i]).FieldByName(field).Interface()
		if val, ok := v.(T); ok {
			r = append(r, val)
		}
	}
	return
}

func RandNum[T IntNumber](start, end T) T {
	end++
	return T(rand.Int63n(int64(end-start))) + start
}

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

func SimpleSort[T any](arr []T, fn func(i, j T) bool) {
	slice := anyArr[T]{
		data: arr,
		fn:   fn,
	}
	sort.Sort(slice)
	return
}

func SimpleSortR[T any](arr []T, fn func(i, j T) bool) (r []T) {
	r = make([]T, 0, len(arr))
	for _, t := range arr {
		r = append(r, t)
	}
	slice := anyArr[T]{
		data: r,
		fn:   fn,
	}
	sort.Sort(slice)
	return
}

func Min[T IntNumber | ~float64 | ~float32](a ...T) T {
	min := a[0]
	for _, t := range a {
		if min > t {
			min = t
		}
	}
	return min
}

func Max[T IntNumber | ~float64 | ~float32](a ...T) T {
	max := a[0]
	for _, t := range a {
		if max < t {
			max = t
		}
	}
	return max
}

func Sum[T IntNumber | ~float64 | ~float32](a ...T) T {
	s := T(0)
	for _, t := range a {
		s += t
	}
	return s
}

func NumberToString[T IntNumber | ~float64 | ~float32](n T) string {
	return fmt.Sprintf("%v", n)
}
