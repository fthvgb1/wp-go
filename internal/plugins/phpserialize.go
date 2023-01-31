package plugins

import (
	"github.com/elliotchance/phpserialize"
	"github.com/fthvgb1/wp-go/helper/maps"
)

func UnPHPSerialize[T any](s string) (r T, err error) {
	var rr map[any]any
	err = phpserialize.Unmarshal([]byte(s), &rr)
	if err == nil {
		rx := maps.AnyAnyToStrAny(rr)
		r, err = maps.StrAnyMapToStruct[T](rx)
	}
	return
}
