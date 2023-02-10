package maps

import "strings"

func GetStrMapAnyVal[T any](m map[string]any, key string) (r T, o bool) {
	k := strings.Split(key, ".")
	if len(k) > 1 {
		val, ok := m[k[0]]
		if ok {
			vx, ok := val.(map[string]any)
			if ok {
				r, o = GetStrMapAnyVal[T](vx, strings.Join(k[1:], "."))
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
