package widget

import (
	"github.com/fthvgb1/wp-go/app/theme/wp"
	"github.com/fthvgb1/wp-go/app/wpconfig"
	"github.com/fthvgb1/wp-go/cache/reload"
	"github.com/fthvgb1/wp-go/helper/maps"
	str "github.com/fthvgb1/wp-go/helper/strings"
)

func Fn(id string, fn func(*wp.Handle, string) string) func(h *wp.Handle) string {
	return func(h *wp.Handle) string {
		return fn(h, id)
	}
}

func configFns[K comparable, V any](m map[K]V, key string, a ...any) func(_ ...any) map[K]V {
	return func(_ ...any) map[K]V {
		c := wpconfig.GetPHPArrayVal[map[K]V](key, nil, a...)
		return maps.FilterZeroMerge(maps.Copy(m), c)
	}
}

func BuildconfigFn[K comparable, V any](m map[K]V, key string, a ...any) func(_ ...any) map[K]V {
	return reload.BuildValFnWithAnyParams(str.Join("widget-config-", key), configFns(m, key, a...))
}
