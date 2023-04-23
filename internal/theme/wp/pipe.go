package wp

import (
	"fmt"
	"github.com/fthvgb1/wp-go/helper/slice"
	"github.com/fthvgb1/wp-go/internal/cmd/reload"
	"github.com/fthvgb1/wp-go/internal/pkg/constraints"
)

type HandlePipeFn[T any] func(HandleFn[T], T)

// HandlePipe  方便把功能写在其它包里
func HandlePipe[T any](initial func(T), fns ...HandlePipeFn[T]) HandleFn[T] {
	return slice.ReverseReduce(fns, func(next HandlePipeFn[T], f func(t T)) func(t T) {
		return func(t T) {
			next(f, t)
		}
	}, initial)
}

func (h *Handle) PushHandler(pipScene int, scene int, fns ...HandleCall) {
	if _, ok := h.handlers[pipScene]; !ok {
		h.handlers[pipScene] = make(map[int][]HandleCall)
	}
	h.handlers[pipScene][scene] = append(h.handlers[pipScene][scene], fns...)
}

func (h *Handle) PushRender(statsOrScene int, fns ...HandleCall) {
	h.PushHandler(constraints.PipRender, statsOrScene, fns...)
}
func (h *Handle) PushDataHandler(scene int, fns ...HandleCall) {
	h.PushHandler(constraints.PipData, scene, fns...)
}

func PipeHandle(pipeScene int, keyFn func(*Handle, int) string, fn func(*Handle, map[int][]HandleCall) []HandleCall) func(HandleFn[*Handle], *Handle) {
	return func(next HandleFn[*Handle], h *Handle) {
		handlers := reload.SafetyMapBy("pipHandlers", keyFn(h, pipeScene), h, func(h *Handle) []HandleCall {
			conf := h.handleHook[pipeScene]
			calls := fn(h, h.handlers[pipeScene])
			calls = slice.FilterAndMap(calls, func(call HandleCall) (HandleCall, bool) {
				ok := true
				for _, hook := range conf {
					call, ok = hook(call)
					if !ok {
						break
					}
				}
				return call, ok
			})
			slice.Sort(calls, func(i, j HandleCall) bool {
				return i.Order > j.Order
			})
			return calls
		})
		for _, handler := range handlers {
			handler.Fn(h)
			if h.abort {
				break
			}
		}
		if !h.stopPipe {
			next(h)
		}
	}
}

func PipeKey(h *Handle, pipScene int) string {
	return fmt.Sprintf("pipekey-%d-%d-%d", pipScene, h.scene, h.scene)
}

func PipeDataHandle(h *Handle, dataHandlers map[int][]HandleCall) (handlers []HandleCall) {
	handlers = append(handlers, dataHandlers[h.scene]...)
	handlers = append(handlers, dataHandlers[constraints.AllScene]...)
	return
}

func PipeRender(h *Handle, renders map[int][]HandleCall) (handlers []HandleCall) {
	handlers = append(handlers, renders[h.Stats]...)
	handlers = append(handlers, renders[h.scene]...)
	handlers = append(handlers, renders[constraints.AllStats]...)
	handlers = append(handlers, renders[constraints.AllScene]...)
	return
}

// DeleteHandle 写插件的时候用
func (h *Handle) DeleteHandle(pipeScene int, name string) {
	h.handleHook[pipeScene] = append(h.handleHook[pipeScene], func(call HandleCall) (HandleCall, bool) {
		return call, name != call.Name
	})
}

// ReplaceHandle 写插件的时候用
func (h *Handle) ReplaceHandle(pipeScene int, name string, fn HandleFn[*Handle]) {
	h.handleHook[pipeScene] = append(h.handleHook[pipeScene], func(call HandleCall) (HandleCall, bool) {
		if name == call.Name {
			call.Fn = fn
		}
		return call, true
	})
}

// HookHandle 写插件的时候用
func (h *Handle) HookHandle(pipeScene int, hook func(HandleCall) (HandleCall, bool)) {
	h.handleHook[pipeScene] = append(h.handleHook[pipeScene], hook)
}
