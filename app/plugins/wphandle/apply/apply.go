package apply

import "github.com/fthvgb1/wp-go/safety"

var fn safety.Var[any]

func SetFn(f any) {
	fn.Store(f)
}

func UsePlugins() any {
	return fn.Load()
}
