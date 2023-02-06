package helper

import (
	"reflect"
)

func ToAny[T any](v T) any {
	return v
}

func Or[T any](is bool, left, right T) T {
	if is {
		return left
	}
	return right
}

func StructColumnToSlice[T any, M any](arr []M, field string) (r []T) {
	for i := 0; i < len(arr); i++ {
		v := reflect.ValueOf(arr[i]).FieldByName(field).Interface()
		if val, ok := v.(T); ok {
			r = append(r, val)
		}
	}
	return
}
