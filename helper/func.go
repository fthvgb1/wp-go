package helper

import (
	"context"
	"fmt"
	str "github.com/fthvgb1/wp-go/helper/strings"
	"net/url"
	"reflect"
	"strconv"
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

func IsZero[T comparable](t T) bool {
	var vv T
	return vv != t
}
func IsZeros(v any) bool {
	switch v.(type) {
	case int64, int, int8, int16, int32, uint64, uint, uint8, uint16, uint32:
		i := fmt.Sprintf("%d", v)
		return str.ToInt[int64](i) == 0
	case float32, float64:
		f := fmt.Sprintf("%v", v)
		ff, _ := strconv.ParseFloat(f, 64)
		return ff == float64(0)
	case bool:
		return v.(bool) == false
	case string:
		s := v.(string)
		return s == ""
	}
	return reflect.DeepEqual(v, reflect.Zero(reflect.TypeOf(v)).Interface())
}

func ToBool[T comparable](t T) bool {
	v := any(t)
	switch v.(type) {
	case string:
		s := v.(string)
		return s != "" && s != "0"
	}
	var vv T
	return vv != t
}

func ToBoolInt(t any) int8 {
	if IsZeros(t) {
		return 0
	}
	return 1
}

func GetContextVal[V, K any](ctx context.Context, k K, defaults V) V {
	v := ctx.Value(k)
	if v == nil {
		return defaults
	}
	vv, ok := v.(V)
	if !ok {
		return defaults
	}
	return vv
}

func IsImplements[T, A any](i A) (T, bool) {
	var a any = i
	t, ok := a.(T)
	return t, ok
}
