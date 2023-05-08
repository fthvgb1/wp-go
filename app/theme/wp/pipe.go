package wp

import (
	"github.com/fthvgb1/wp-go/app/cmd/reload"
	"github.com/fthvgb1/wp-go/app/pkg/constraints"
	"github.com/fthvgb1/wp-go/helper"
	"github.com/fthvgb1/wp-go/helper/slice"
	str "github.com/fthvgb1/wp-go/helper/strings"
)

type HandlePipeFn[T any] func(HandleFn[T], T)

type Pipe struct {
	Name  string
	Order int
	Fn    HandlePipeFn[*Handle]
}

func NewPipe(name string, order int, fn HandlePipeFn[*Handle]) Pipe {
	return Pipe{Name: name, Order: order, Fn: fn}
}

// HandlePipe  方便把功能写在其它包里
func HandlePipe[T any](initial func(T), fns ...HandlePipeFn[T]) HandleFn[T] {
	return slice.ReverseReduce(fns, func(next HandlePipeFn[T], f func(t T)) func(t T) {
		return func(t T) {
			next(f, t)
		}
	}, initial)
}

func (h *Handle) PushPipe(scene string, pipes ...Pipe) error {
	return PushFn("pipe", scene, pipes...)
}
func (h *Handle) PushPipeHook(scene string, pipes ...func(Pipe) (Pipe, bool)) error {
	return PushFnHook("pipeHook", scene, pipes...)
}

func (h *Handle) DeletePipe(scene, pipeName string) error {
	return PushFnHook("pipeHook", scene, func(pipe Pipe) (Pipe, bool) {
		return pipe, pipeName != pipe.Name
	})
}

func (h *Handle) ReplacePipe(scene, pipeName string, pipe Pipe) error {
	return PushFnHook("pipeHook", scene, func(p Pipe) (Pipe, bool) {
		if pipeName == p.Name {
			p = pipe
		}
		return p, true
	})
}

func (h *Handle) PushHandler(pipScene string, scene string, fns ...HandleCall) {
	if _, ok := h.handlers[pipScene]; !ok {
		h.handlers[pipScene] = make(map[string][]HandleCall)
	}
	h.handlers[pipScene][scene] = append(h.handlers[pipScene][scene], fns...)
}

func (h *Handle) PushRender(statsOrScene string, fns ...HandleCall) {
	h.PushHandler(constraints.PipeRender, statsOrScene, fns...)
}
func (h *Handle) PushDataHandler(scene string, fns ...HandleCall) {
	h.PushHandler(constraints.PipeData, scene, fns...)
}

func PipeHandle(pipeScene string, keyFn func(*Handle, string) string, fn func(*Handle, map[string][]HandleCall, string) []HandleCall) func(HandleFn[*Handle], *Handle) {
	return func(next HandleFn[*Handle], h *Handle) {
		key := keyFn(h, pipeScene)
		handlers := reload.GetAnyValMapBy("pipeHandlers", key, h, func(h *Handle) []HandleCall {
			conf := h.handleHook[pipeScene]
			calls := fn(h, h.handlers[pipeScene], key)
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

func PipeKey(h *Handle, pipScene string) string {
	key := str.Join("pipekey", "-", pipScene, "-", h.scene, "-", h.Stats)
	return h.ComponentFilterFnHook("pipeKey", key, pipScene)
}

func Run(h *Handle, conf func(*Handle)) {
	if !helper.GetContextVal(h.C, "inited", false) {
		InitHandle(conf, h)
	}
	reload.GetAnyValBys(str.Join("pipeInit-", h.scene), h, func(h *Handle) func(*Handle) {
		p := GetFn[Pipe]("pipe", constraints.AllScene)
		p = append(p, GetFn[Pipe]("pipe", h.scene)...)
		pipes := slice.FilterAndMap(p, func(pipe Pipe) (Pipe, bool) {
			var ok bool
			hooks := GetFnHook[func(Pipe) (Pipe, bool)]("pipeHook", constraints.AllScene)
			hooks = append(hooks, GetFnHook[func(Pipe) (Pipe, bool)]("pipeHook", h.scene)...)
			for _, fn := range hooks {
				pipe, ok = fn(pipe)
				if !ok {
					return pipe, false
				}
			}
			return pipe, pipe.Fn != nil
		})
		slice.Sort(pipes, func(i, j Pipe) bool {
			return i.Order > j.Order
		})
		arr := slice.Map(pipes, func(t Pipe) HandlePipeFn[*Handle] {
			return t.Fn
		})
		return HandlePipe(NothingToDo, arr...)
	})(h)
}

func MiddlewareKey(h *Handle, pipScene string) string {
	return h.ComponentFilterFnHook("middleware", "middleware", pipScene)
}

func PipeMiddlewareHandle(h *Handle, middlewares map[string][]HandleCall, key string) (handlers []HandleCall) {
	handlers = append(handlers, middlewares[h.scene]...)
	handlers = append(handlers, middlewares[constraints.AllScene]...)
	handlers = *h.PipeHandleHook("PipeMiddlewareHandle", &handlers, key)
	return
}

func PipeDataHandle(h *Handle, dataHandlers map[string][]HandleCall, key string) (handlers []HandleCall) {
	handlers = append(handlers, dataHandlers[h.scene]...)
	handlers = append(handlers, dataHandlers[constraints.AllScene]...)
	handlers = *h.PipeHandleHook("PipeDataHandle", &handlers, key)
	return
}

func PipeRender(h *Handle, renders map[string][]HandleCall, key string) (handlers []HandleCall) {
	handlers = append(handlers, renders[h.Stats]...)
	handlers = append(handlers, renders[h.scene]...)
	handlers = append(handlers, renders[constraints.AllStats]...)
	handlers = append(handlers, renders[constraints.AllScene]...)
	handlers = *h.PipeHandleHook("PipeRender", &handlers, key)
	return
}

// DeleteHandle 写插件的时候用
func (h *Handle) DeleteHandle(pipeScene string, name string) {
	h.handleHook[pipeScene] = append(h.handleHook[pipeScene], func(call HandleCall) (HandleCall, bool) {
		return call, name != call.Name
	})
}

// ReplaceHandle 写插件的时候用
func (h *Handle) ReplaceHandle(pipeScene, name string, fn HandleFn[*Handle]) {
	h.handleHook[pipeScene] = append(h.handleHook[pipeScene], func(call HandleCall) (HandleCall, bool) {
		if name == call.Name {
			call.Fn = fn
		}
		return call, true
	})
}

// HookHandle 写插件的时候用
func (h *Handle) HookHandle(pipeScene string, hook func(HandleCall) (HandleCall, bool)) {
	h.handleHook[pipeScene] = append(h.handleHook[pipeScene], hook)
}

func (h *Handle) PushPipeHandleHook(name string, fn ...func([]HandleCall) []HandleCall) error {
	return PushFnHook("pipeHandleHook", name, fn...)
}

func (h *Handle) PipeHandleHook(name string, calls *[]HandleCall, key string) *[]HandleCall {
	fn := GetFnHook[func(*Handle, *[]HandleCall, string) *[]HandleCall]("pipeHandleHook", name)
	return slice.Reduce(fn, func(t func(*Handle, *[]HandleCall, string) *[]HandleCall, r *[]HandleCall) *[]HandleCall {
		return t(h, r, key)
	}, calls)
}

func InitPipe(h *Handle) {
	h.PushPipe(constraints.Home, NewPipe(constraints.PipeMiddleware, 300,
		PipeHandle(constraints.PipeMiddleware, MiddlewareKey, PipeMiddlewareHandle)))

	h.PushPipe(constraints.AllScene, NewPipe(constraints.PipeData, 200,
		PipeHandle(constraints.PipeData, PipeKey, PipeDataHandle)))
	h.PushPipe(constraints.AllScene, NewPipe(constraints.PipeRender, 100,
		PipeHandle(constraints.PipeRender, PipeKey, PipeRender)))
}
