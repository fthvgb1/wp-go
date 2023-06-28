package maps

import "strings"

func GetStrAnyVal[T any](m map[string]any, key string) (r T, o bool) {
	k := strings.Split(key, ".")
	if len(k) > 1 {
		val, ok := m[k[0]]
		if ok {
			vx, ok := val.(map[string]any)
			if ok {
				r, o = GetStrAnyVal[T](vx, strings.Join(k[1:], "."))
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
func GetStrMapAnyValWithAny(key string, v map[string]any) (r any, o bool) {
	k := strings.Split(key, ".")
	if len(k) > 1 {
		val, ok := v[k[0]]
		if ok {
			vx, ok := val.(map[string]any)
			if ok {
				r, o = GetStrMapAnyValWithAny(strings.Join(k[1:], "."), vx)
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

func RecursiveSetStrVal[T any](m map[string]any, k string, v T) {
	kk := strings.Split(k, ".")
	if len(kk) < 1 {
		return
	} else if len(kk) < 2 {
		m[k] = v
		return
	}
	for i, _ := range kk[0 : len(kk)-1] {
		key := strings.Join(kk[0:i+1], ".")
		mm, ok := GetStrAnyVal[map[string]any](m, key)
		if !ok {
			mm = map[string]any{}
			preKey := strings.Join(kk[0:i], ".")
			if preKey == "" {
				RecursiveSetStrVal(m, key, mm)
			} else {
				m, _ := GetStrAnyVal[map[string]any](m, preKey)
				RecursiveSetStrVal(m, kk[i], mm)
			}
		}
	}
	key := strings.Join(kk[0:len(kk)-1], ".")
	mm, _ := GetStrAnyVal[map[string]any](m, key)
	mm[kk[len(kk)-1]] = v
}
