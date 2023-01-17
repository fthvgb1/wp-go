package helper

import "encoding/json"

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

func MapToSlice[T any, K comparable, V any](m map[K]V, fn func(K, V) (T, bool)) (r []T) {
	for k, v := range m {
		vv, ok := fn(k, v)
		if ok {
			r = append(r, vv)
		}
	}
	return
}

// MapAnyAnyToStrAny map[any]any => map[string]any 方便json转换
func MapAnyAnyToStrAny(m map[any]any) (r map[string]any) {
	r = make(map[string]any)
	for k, v := range m {
		kk, ok := k.(string)
		if ok {
			vv, ok := v.(map[any]any)
			if ok {
				r[kk] = MapAnyAnyToStrAny(vv)
			} else {
				r[kk] = v
			}
		}

	}
	return
}
