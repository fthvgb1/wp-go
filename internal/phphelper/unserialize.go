package phphelper

import (
	"github.com/elliotchance/phpserialize"
	"github.com/fthvgb1/wp-go/helper/maps"
)

// UnPHPSerializeToStruct 使用 json tag
func UnPHPSerializeToStruct[T any](s string) (r T, err error) {
	var rr map[any]any
	err = phpserialize.Unmarshal([]byte(s), &rr)
	if err == nil {
		rx := maps.AnyAnyToStrAny(rr)
		r, err = maps.StrAnyMapToStruct[T](rx)
	}
	return
}

func UnPHPSerializeToAnyMap(s string) (map[string]any, error) {
	m := map[string]any{}
	var r map[any]any
	err := phpserialize.Unmarshal([]byte(s), &r)
	if err != nil {
		return nil, err
	}

	m = maps.AnyAnyToStrAny(r)
	return m, err
}
