package helper

import "encoding/json"

func MapToStruct[T any](m map[string]any) (r T, err error) {
	str, err := json.Marshal(m)
	if err != nil {
		return
	}
	err = json.Unmarshal(str, &r)
	return
}

func StructToMap[T any](s T) (r map[string]any, err error) {
	marshal, err := json.Marshal(s)
	if err != nil {
		return
	}
	r = make(map[string]any)
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
