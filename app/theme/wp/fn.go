package wp

import (
	"errors"
	"github.com/fthvgb1/wp-go/safety"
)

var fnMap = safety.NewMap[string, map[string]any]()
var fnHook = safety.NewMap[string, map[string]any]()

func GetFn[T any](fnType string, name string) []T {
	v, ok := fnMap.Load(fnType)
	if !ok {
		return nil
	}
	vv, ok := v[name]
	if !ok {
		return nil
	}
	return vv.([]T)
}
func GetFnHook[T any](fnType string, name string) []T {
	v, ok := fnHook.Load(fnType)
	if !ok {
		return nil
	}
	vv, ok := v[name]
	if !ok {
		return nil
	}
	return vv.([]T)
}

func PushFn[T any](fnType string, name string, fns ...T) error {
	v, ok := fnMap.Load(fnType)
	if !ok {
		v = make(map[string]any)
		fnMap.Store(fnType, v)
		v[name] = fns
		return nil
	}
	vv, ok := v[name]
	if !ok || vv == nil {
		v[name] = fns
		return nil
	}
	s, ok := vv.([]T)
	if ok {
		s = append(s, fns...)
		v[name] = s
	}
	return errors.New("error fn type")
}

func PushFnHook[T any](fnType string, name string, fns ...T) error {
	v, ok := fnHook.Load(fnType)
	if !ok {
		v = make(map[string]any)
		fnHook.Store(fnType, v)
		v[name] = fns
		return nil
	}
	vv, ok := v[name]
	if !ok || vv == nil {
		v[name] = fns
		return nil
	}
	s, ok := vv.([]T)
	if ok {
		s = append(s, fns...)
		v[name] = s
	}
	return errors.New("error fn type")
}
