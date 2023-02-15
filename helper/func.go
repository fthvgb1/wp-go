package helper

import (
	"net/url"
	"reflect"
	"strings"
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

func UrlScheme(u string, isHttps bool) string {
	return Or(isHttps,
		strings.Replace(u, "http://", "https://", 1),
		strings.Replace(u, "https://", "http://", 1),
	)
}

func CutUrlHost(u string) string {
	ur, err := url.Parse(u)
	if err != nil {
		return u
	}
	ur.Scheme = ""
	ur.Host = ""
	return ur.String()
}

func Defaults[T comparable](v, defaults T) T {
	var zero T
	if v == zero {
		return defaults
	}
	return v
}
func DefaultVal[T any](v, defaults T) T {
	var zero T
	if reflect.DeepEqual(zero, v) {
		return defaults
	}
	return v
}
