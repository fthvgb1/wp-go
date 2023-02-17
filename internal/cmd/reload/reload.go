package reload

import "github.com/fthvgb1/wp-go/safety"

var calls []func()

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
}
