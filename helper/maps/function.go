package maps

import (
	"github.com/fthvgb1/wp-go/helper/slice"
	"golang.org/x/exp/constraints"
	"strings"
)

func GetStrAnyVal[T any](m map[string]any, key string, delimiter ...string) (r T, o bool) {
	separator := "."
	if len(delimiter) > 0 && delimiter[0] != "" {
		separator = delimiter[0]
	}
	k := strings.Split(key, separator)
	if len(k) > 1 {
		val, ok := m[k[0]]
		if ok {
			vx, ok := val.(map[string]any)
			if ok {
				r, o = GetStrAnyVal[T](vx, strings.Join(k[1:], separator))
			}
		}
	} else {
		x, ok := m[k[0]]
		if ok {
			vv, ok := x.(T)
			if ok {
				o = true
				r = vv
			}
		}
	}
	return
}

func GetStrAnyValWithDefaults[T any](m map[string]any, key string, defaults T) (r T) {
	r = defaults
	v, ok := GetStrAnyVal[T](m, key)
	if !ok {
		return
	}
	r = v
	return
}

// GetStrMapAnyValWithAny 使用"." 分隔层级
func GetStrMapAnyValWithAny(v map[string]any, key string) (r any, o bool) {
	k := strings.Split(key, ".")
	if len(k) > 1 {
		val, ok := v[k[0]]
		if ok {
			vx, ok := val.(map[string]any)
			if ok {
				r, o = GetStrMapAnyValWithAny(vx, strings.Join(k[1:], "."))
			}
		}
	} else {
		x, ok := v[k[0]]
		if ok {
			o = true
			r = x
		}
	}
	return
}

func GetAnyAnyMapVal[T any](m map[any]any, k ...any) (r T, o bool) {
	if len(k) > 1 {
		val, ok := m[k[0]]
		if ok {
			vx, ok := val.(map[any]any)
			if ok {
				r, o = GetAnyAnyMapVal[T](vx, k[1:]...)
			}
		}
	} else {
		x, ok := m[k[0]]
		if ok {
			vv, ok := x.(T)
			if ok {
				o = true
				r = vv
			}
		}
	}
	return
}

func GetAnyAnyMapWithAny(v map[any]any, k ...any) (r any, o bool) {
	if len(k) > 1 {
		val, ok := v[k[0]]
		if ok {
			vx, ok := val.(map[any]any)
			if ok {
				r, o = GetAnyAnyMapWithAny(vx, k[1:]...)
			}
		}
	} else {
		x, ok := v[k[0]]
		if ok {
			o = true
			r = x
		}
	}
	return
}

func GetAnyAnyValWithDefaults[T any](m map[any]any, defaults T, key ...any) (r T) {
	r = defaults
	v, ok := GetAnyAnyMapVal[T](m, key...)
	if !ok {
		return
	}
	r = v
	return
}

func SetStrAnyVal[T any](m map[string]any, k string, v T, delimiter ...string) {
	del := "."
	if len(delimiter) > 0 && delimiter[0] != "" {
		del = delimiter[0]
	}
	kk := strings.Split(k, del)
	if len(kk) < 1 {
		return
	} else if len(kk) < 2 {
		m[k] = v
		return
	}
	mm, ok := GetStrAnyVal[map[string]any](m, strings.Join(kk[0:len(kk)-1], del))
	if ok {
		mm[kk[len(kk)-1]] = v
		return
	}
	mx, ok := GetStrAnyVal[map[string]any](m, kk[0])
	if !ok {
		m[kk[0]] = map[string]any{}
		mx = m[kk[0]].(map[string]any)
	}
	for i, _ := range kk[0 : len(kk)-2] {
		key := kk[i+1]
		mm, ok := mx[key]
		if !ok {
			mmm := map[string]any{}
			mx[key] = mmm
			mx = mmm

		} else {
			mx = mm.(map[string]any)
		}
	}
	mx[kk[len(kk)-1]] = v
}

func SetAnyAnyVal[T any](m map[any]any, v T, k ...any) {
	if len(k) < 1 {
		return
	} else if len(k) == 1 {
		m[k[0]] = v
		return
	}
	for i, _ := range k[0 : len(k)-1] {
		key := k[0 : i+1]
		mm, ok := GetAnyAnyMapVal[map[any]any](m, key...)
		if !ok {
			mm = map[any]any{}
			preKey := k[0:i]
			if len(preKey) == 0 {
				SetAnyAnyVal(m, mm, key...)
			} else {
				m, _ := GetAnyAnyMapVal[map[any]any](m, preKey...)
				SetAnyAnyVal(m, mm, k[i])
			}
		}
	}
	key := k[0 : len(k)-1]
	mm, _ := GetAnyAnyMapVal[map[any]any](m, key...)
	mm[k[len(k)-1]] = v
}

func AscEahByKey[K constraints.Ordered, V any](m map[K]V, fn func(K, V)) {
	orderedEahByKey(m, slice.ASC, fn)
}
func DescEahByKey[K constraints.Ordered, V any](m map[K]V, fn func(K, V)) {
	orderedEahByKey(m, slice.ASC, fn)
}

func orderedEahByKey[K constraints.Ordered, V any](m map[K]V, ordered int, fn func(K, V)) {
	keys := Keys(m)
	slice.Sorts(keys, ordered)
	for _, key := range keys {
		fn(key, m[key])
	}
}

func Flip[K, V comparable](m map[K]V) map[V]K {
	mm := make(map[V]K, len(m))
	for k, v := range m {
		mm[v] = k
	}
	return mm
}
