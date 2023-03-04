package reload

import (
	"github.com/fthvgb1/wp-go/helper/slice"
	"github.com/fthvgb1/wp-go/safety"
	"strings"
)

var calls []func()

var str = safety.NewMap[string, string]()

func GetStr(name string) (string, bool) {
	return str.Load(name)
}
func SetStr(name, val string) {
	str.Store(name, val)
}

func GetStrBy[T any](key, delimiter string, t T, fn ...func(T) string) string {
	v, ok := str.Load(key)
	if ok {
		return v
	}
	v = strings.Join(slice.Map(fn, func(vv func(T) string) string {
		return vv(t)
	}), delimiter)
	str.Store(key, v)
	return v
}

func Vars[T any](defaults T) *safety.Var[T] {
	ss := safety.NewVar(defaults)
	calls = append(calls, func() {
		ss.Store(defaults)
	})
	return ss
}
func VarsBy[T any](fn func() T) *safety.Var[T] {
	ss := safety.NewVar(fn())
	calls = append(calls, func() {
		ss.Store(fn())
	})
	return ss
}

func Push(fn ...func()) {
	calls = append(calls, fn...)
}

func Reload() {
	for _, call := range calls {
		call()
	}
	str.Flush()
}
