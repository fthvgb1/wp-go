package helper

import (
	"reflect"
)

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

func RangeSlice[T ~int | ~uint | ~int64 | ~int8 | ~int16 | ~int32 | ~uint64](start, end, step T) []T {
	r := make([]T, 0, int(end/step+1))
	for i := start; i <= end; {
		r = append(r, i)
		i = i + step
	}
	return r
}
