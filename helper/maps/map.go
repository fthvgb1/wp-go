package maps

import (
	"encoding/json"
)

func StrAnyMapToStruct[T any, M any](m M) (r T, err error) {
	str, err := json.Marshal(m)
	if err != nil {
		return
	}
	err = json.Unmarshal(str, &r)
	return
}

func StructToAnyMap[K comparable, T any](s T) (r map[K]any, err error) {
	marshal, err := json.Marshal(s)
	if err != nil {
		return
	}
	r = make(map[K]any)
	err = json.Unmarshal(marshal, &r)
	return
}

func FilterToSlice[T any, K comparable, V any](m map[K]V, fn func(K, V) (T, bool)) (r []T) {
	for k, v := range m {
		vv, ok := fn(k, v)
		if ok {
			r = append(r, vv)
		}
	}
	return
}

// AnyAnyToStrAny map[any]any => map[string]any 方便json转换
func AnyAnyToStrAny(m map[any]any) (r map[string]any) {
	r = make(map[string]any)
	for k, v := range m {
		kk, ok := k.(string)
		if ok {
			vv, ok := v.(map[any]any)
			if ok {
				r[kk] = AnyAnyToStrAny(vv)
			} else {
				r[kk] = v
			}
		}

	}
	return
}

func IsExists[K comparable, V any](m map[K]V, k K) bool {
	_, ok := m[k]
	return ok
}

func Keys[K comparable, V any](m map[K]V) []K {
	return FilterToSlice(m, func(k K, v V) (K, bool) {
		return k, true
	})
}
func Values[K comparable, V any](m map[K]V) []V {
	return FilterToSlice(m, func(k K, v V) (V, bool) {
		return v, true
	})
}

func Reduce[T, V any, K comparable](m map[K]V, fn func(K, V, T) T, r T) T {
	for k, v := range m {
		r = fn(k, v, r)
	}
	return r
}

func Replace[K comparable, V any](m map[K]V, mm ...map[K]V) map[K]V {
	for _, n := range mm {
		for k, v := range n {
			_, ok := m[k]
			if ok {
				m[k] = v
			}
		}
	}
	return m
}
