package apply

import "github.com/fthvgb1/wp-go/safety"

var contains = safety.NewMap[string, any]()

func SetVal(key string, val any) {
	contains.Store(key, val)
}

func DelVal(key string) {
	contains.Delete(key)
}

func GetVal[V any](key string) (V, bool) {
	v, ok := contains.Load(key)
	if !ok {
		var vv V
		return vv, ok
	}
	return v.(V), ok
}
func GetRawVal(key string) (any, bool) {
	return contains.Load(key)
}

func GetPlugins() any {
	v, _ := contains.Load("wp-plugins")
	return v
}
