package widget

import (
	"github.com/fthvgb1/wp-go/helper/maps"
	str "github.com/fthvgb1/wp-go/helper/strings"
	"github.com/fthvgb1/wp-go/internal/cmd/reload"
	"github.com/fthvgb1/wp-go/internal/theme/wp"
	"github.com/fthvgb1/wp-go/internal/wpconfig"
)

func Fn(id string, fn func(*wp.Handle, string) string) func(h *wp.Handle) string {
	return func(h *wp.Handle) string {
		return fn(h, id)
	}
}

func configs[M ~map[K]V, K comparable, V any](m M, key string, a ...any) M {
	return reload.GetAnyValBys(str.Join("widget-config-", key), key, func(_ string) M {
		c := wpconfig.GetPHPArrayVal[M](key, nil, a...)
		return maps.FilterZeroMerge(maps.Copy(m), c)
	})
}
