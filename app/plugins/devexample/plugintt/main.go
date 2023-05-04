package main

import (
	"github.com/fthvgb1/wp-go/app/pkg/constraints"
	"github.com/fthvgb1/wp-go/app/theme"
	"github.com/fthvgb1/wp-go/app/theme/wp"
	"plugintt/xx"
)

func init() {
	// here can register theme
	theme.AddThemeHookFunc("xxx", Xo)
}

func Xo(h *wp.Handle) {
	// action or hook or config theme
	h.PushHandler(constraints.PipeRender, constraints.Home, wp.HandleCall{
		Fn: func(handle *wp.Handle) {
			xx.Xo()
		},
		Order: 100,
		Name:  "xxxx",
	})
}
