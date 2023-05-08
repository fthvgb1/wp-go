package apply

var fn any

func SetFn(f any) {
	fn = f
}

func UsePlugins() any {
	return fn
}
