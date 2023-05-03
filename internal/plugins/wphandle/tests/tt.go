package tests

import (
	"fmt"
	"github.com/fthvgb1/wp-go/internal/pkg/constraints"
	"github.com/fthvgb1/wp-go/internal/theme/wp"
)

func Tt(h *wp.Handle) {
	h.HookHandle(constraints.PipeMiddleware, func(call wp.HandleCall) (wp.HandleCall, bool) {
		return call, false
	})
	/*h.PushPipeHook(constraints.Home, func(pipe wp.Pipe) (wp.Pipe, bool) {
		return wp.Pipe{}, false
	})*/
	//h.DeletePipe(constraints.Home, constraints.PipeMiddleware)
	h.ReplacePipe(constraints.Home, constraints.PipeMiddleware, wp.NewPipe("log", 500, func(next wp.HandleFn[*wp.Handle], h *wp.Handle) {
		fmt.Println("ffff")
		next(h)
		fmt.Println("iiiii")
	}))
}
