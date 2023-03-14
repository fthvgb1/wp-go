package widget

import "github.com/fthvgb1/wp-go/internal/theme/wp"

func MiddleWare(call ...wp.HandlePipeFn[*wp.Handle]) []wp.HandlePipeFn[*wp.Handle] {
	return append([]wp.HandlePipeFn[*wp.Handle]{
		IsCategory,
	}, call...)
}
