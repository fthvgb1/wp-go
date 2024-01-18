package wp

import (
	"github.com/fthvgb1/wp-go/app/pkg/constraints"
	"github.com/fthvgb1/wp-go/cache/reload"
	"github.com/fthvgb1/wp-go/helper/slice"
	str "github.com/fthvgb1/wp-go/helper/strings"
)

type HandlePipeFn[T any] func(HandleFn[T], T)

type Pipe struct {
	Name  string
	Order float64
	Fn    HandlePipeFn[*Handle]
}

func NewPipe(name string, order float64, fn HandlePipeFn[*Handle]) Pipe {
	return Pipe{Name: name, Order: order, Fn: fn}
}

// HandlePipe  方便把功能写在其它包里
func HandlePipe[T any](initial func(T), fns ...HandlePipeFn[T]) HandleFn[T] {
	return slice.ReverseReduce(fns, func(next HandlePipeFn[T], f HandleFn[T]) HandleFn[T] {
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
	v, ok := handlerss.Load(pipScene)
	if !ok {
		v = make(map[string][]HandleCall)
	}
	v[scene] = append(v[scene], fns...)
	handlerss.Store(pipScene, v)
}

func (h *Handle) PushRender(statsOrScene string, fns ...HandleCall) {
	h.PushHandler(constraints.PipeRender, statsOrScene, fns...)
}
func (h *Handle) PushDataHandler(scene string, fns ...HandleCall) {
	h.PushHandler(constraints.PipeData, scene, fns...)
}

func BuildHandlers(pipeScene string, keyFn func(*Handle, string) string,
	fn func(*Handle, map[string][]HandleCall, string) []HandleCall) func(HandleFn[*Handle], *Handle) {

	pipeHandlerFn := reload.BuildMapFn[string]("pipeHandlers", BuildHandler(pipeScene, keyFn, fn))

	return func(next HandleFn[*Handle], h *Handle) {
		key := keyFn(h, pipeScene)
		handlers := pipeHandlerFn(key, h)
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

func BuildHandler(pipeScene string, keyFn func(*Handle, string) string,
	fn func(*Handle, map[string][]HandleCall, string) []HandleCall) func(*Handle) []HandleCall {

	return func(h *Handle) []HandleCall {
		key := keyFn(h, pipeScene)
		mut := reload.GetGlobeMutex()
		mut.Lock()
		hookers, _ := handleHooks.Load(pipeScene)
		hh, _ := handlerss.Load(pipeScene)
		mut.Unlock()
		calls := fn(h, hh, key)
		calls = slice.FilterAndMap(calls, func(call HandleCall) (HandleCall, bool) {
			ok := true
			for _, hook := range hookers {
				call, ok = hook(call)
				if !ok {
					break
				}
			}
			return call, ok
		})
		slice.SimpleSort(calls, slice.DESC, func(t HandleCall) float64 {
			return t.Order
		})
		return calls
	}
}

func PipeKey(h *Handle, pipScene string) string {
	key := str.Join("pipekey", "-", pipScene, "-", h.scene, "-", h.Stats)
	return h.DoActionFilter("pipeKey", key, pipScene)
}

var pipeInitFn = reload.BuildMapFn[string]("pipeInit", BuildPipe)

func Run(h *Handle, conf func(*Handle)) {
	if !h.isInited {
		InitHandle(conf, h)
	}
	pipeInitFn(h.scene, h.scene)(h)
}

func BuildPipe(scene string) func(*Handle) {
	pipees := GetFn[Pipe]("pipe", constraints.AllScene)
	pipees = append(pipees, GetFn[Pipe]("pipe", scene)...)
	pipes := slice.FilterAndMap(pipees, func(pipe Pipe) (Pipe, bool) {
		var ok bool
		mut := reload.GetGlobeMutex()
		mut.Lock()
		hooks := GetFnHook[func(Pipe) (Pipe, bool)]("pipeHook", constraints.AllScene)
		hooks = append(hooks, GetFnHook[func(Pipe) (Pipe, bool)]("pipeHook", scene)...)
		mut.Unlock()
		for _, fn := range hooks {
			pipe, ok = fn(pipe)
			if !ok {
				return pipe, false
			}
		}
		return pipe, pipe.Fn != nil
	})
	slice.SimpleSort(pipes, slice.DESC, func(t Pipe) float64 {
		return t.Order
	})

	arr := slice.Map(pipes, func(t Pipe) HandlePipeFn[*Handle] {
		return t.Fn
	})
	return HandlePipe(NothingToDo, arr...)
}

func MiddlewareKey(h *Handle, pipScene string) string {
	return h.DoActionFilter("middleware", str.Join("pipe-middleware-", h.scene), pipScene)
}

func PipeMiddlewareHandle(h *Handle, middlewares map[string][]HandleCall, key string) (handlers []HandleCall) {
	handlers = append(handlers, middlewares[h.scene]...)
	handlers = append(handlers, middlewares[constraints.AllScene]...)
	handlers = h.PipeHandleHook("PipeMiddlewareHandle", handlers, middlewares, key)
	return
}

func PipeDataHandle(h *Handle, dataHandlers map[string][]HandleCall, key string) (handlers []HandleCall) {
	handlers = append(handlers, dataHandlers[h.scene]...)
	handlers = append(handlers, dataHandlers[constraints.AllScene]...)
	handlers = h.PipeHandleHook("PipeDataHandle", handlers, dataHandlers, key)
	return
}

func PipeRender(h *Handle, renders map[string][]HandleCall, key string) (handlers []HandleCall) {
	handlers = append(handlers, renders[h.Stats]...)
	handlers = append(handlers, renders[h.scene]...)
	handlers = append(handlers, renders[constraints.AllStats]...)
	handlers = append(handlers, renders[constraints.AllScene]...)
	handlers = h.PipeHandleHook("PipeRender", handlers, renders, key)
	return
}

// DeleteHandle 写插件的时候用
func (h *Handle) DeleteHandle(pipeScene string, name string) {
	v, _ := handleHooks.Load(pipeScene)
	v = append(v, func(call HandleCall) (HandleCall, bool) {
		return call, name != call.Name
	})
	handleHooks.Store(pipeScene, v)
}

// ReplaceHandle 写插件的时候用
func (h *Handle) ReplaceHandle(pipeScene, name string, fn HandleFn[*Handle]) {
	v, _ := handleHooks.Load(pipeScene)
	v = append(v, func(call HandleCall) (HandleCall, bool) {
		if name == call.Name {
			call.Fn = fn
		}
		return call, true
	})
	handleHooks.Store(pipeScene, v)
}

// HookHandle 写插件的时候用
func (h *Handle) HookHandle(pipeScene string, hook func(HandleCall) (HandleCall, bool)) {
	v, _ := handleHooks.Load(pipeScene)
	v = append(v, hook)
	handleHooks.Store(pipeScene, v)
}

func (h *Handle) PushPipeHandleHook(name string, fn ...func([]HandleCall) []HandleCall) error {
	return PushFnHook("pipeHandleHook", name, fn...)
}

func (h *Handle) PipeHandleHook(name string, calls []HandleCall, m map[string][]HandleCall, key string) []HandleCall {
	fn := GetFnHook[func(*Handle, []HandleCall, map[string][]HandleCall, string) []HandleCall]("pipeHandleHook", name)
	return slice.Reduce(fn, func(t func(*Handle, []HandleCall, map[string][]HandleCall, string) []HandleCall, r []HandleCall) []HandleCall {
		return t(h, r, m, key)
	}, calls)
}

func InitPipe(h *Handle) {
	h.PushPipe(constraints.AllScene, NewPipe(constraints.PipeMiddleware, 300,
		BuildHandlers(constraints.PipeMiddleware, MiddlewareKey, PipeMiddlewareHandle)))

	h.PushPipe(constraints.AllScene, NewPipe(constraints.PipeData, 200,
		BuildHandlers(constraints.PipeData, PipeKey, PipeDataHandle)))
	h.PushPipe(constraints.AllScene, NewPipe(constraints.PipeRender, 100,
		BuildHandlers(constraints.PipeRender, PipeKey, PipeRender)))
}
