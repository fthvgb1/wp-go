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
