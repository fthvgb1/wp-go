package helper

import "strings"

func GetStrMapAnyVal[T any](key string, v map[string]any) (r T, o bool) {
	k := strings.Split(key, ".")
	if len(k) > 1 {
		val, ok := v[k[0]]
		if ok {
			vx, ok := val.(map[string]any)
			if ok {
				r, o = GetStrMapAnyVal[T](strings.Join(k[1:], "."), vx)
			}
		}
	} else {
		x, ok := v[k[0]]
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

func GetStrMapAnyValToAny(key string, v map[string]any) (r any, o bool) {
	k := strings.Split(key, ".")
	if len(k) > 1 {
		val, ok := v[k[0]]
		if ok {
			vx, ok := val.(map[string]any)
			if ok {
				r, o = GetStrMapAnyValToAny(strings.Join(k[1:], "."), vx)
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
