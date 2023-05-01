package maps

import (
	"encoding/json"
	"github.com/fthvgb1/wp-go/helper"
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
func StrAnyToAnyAny(m map[string]any) (r map[any]any) {
	r = make(map[any]any)
	for kk, v := range m {
		vv, ok := v.(map[string]any)
		if ok {
			r[kk] = StrAnyToAnyAny(vv)
		} else {
			r[kk] = v
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

func Copy[K comparable, V any](m map[K]V) map[K]V {
	r := make(map[K]V)
	for k, v := range m {
		r[k] = v
	}
	return r
}

func Merge[K comparable, V any](m ...map[K]V) map[K]V {
	if len(m) < 1 {
		return nil
	} else if len(m) < 2 {
		return m[0]
	}
	mm := m[0]
	if mm == nil {
		mm = make(map[K]V)
	}
	for _, m2 := range m[1:] {
		for k, v := range m2 {
			mm[k] = v
		}
	}
	return mm
}

func MergeBy[K comparable, V any](fn func(k K, v1, v2 V) (V, bool), m ...map[K]V) map[K]V {
	if len(m) < 1 {
		return nil
	} else if len(m) < 2 {
		return m[0]
	}
	mm := m[0]
	if mm == nil {
		mm = make(map[K]V)
	}
	for _, m2 := range m[1:] {
		for k, v := range m2 {
			vv, ok := mm[k]
			if ok {
				vvv, ok := fn(k, vv, v)
				if ok {
					v = vvv
				}
			}
			mm[k] = v
		}
	}
	return mm
}

func FilterZeroMerge[K comparable, V any](m ...map[K]V) map[K]V {
	if len(m) < 1 {
		return nil
	} else if len(m) < 2 {
		return m[0]
	}
	mm := m[0]
	if mm == nil {
		mm = make(map[K]V)
	}
	for _, m2 := range m[1:] {
		for k, v := range m2 {
			if helper.IsZeros(v) {
				continue
			}
			mm[k] = v
		}
	}
	return mm
}

func WithDefaultVal[K comparable, V any](m map[K]V, k K, defaults V) V {
	vv, ok := m[k]
	if ok {
		return vv
	}
	return defaults
}

func AnyAnyMap[K comparable, V any](m map[any]any, fn func(k, v any) (K, V, bool)) map[K]V {
	mm := make(map[K]V, 0)
	for k, v := range m {
		key, val, ok := fn(k, v)
		if ok {
			mm[key] = val
		}
	}
	return mm
}
