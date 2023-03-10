package wp

import (
	"fmt"
	"github.com/fthvgb1/wp-go/helper/maps"
	"github.com/fthvgb1/wp-go/helper/slice"
	str "github.com/fthvgb1/wp-go/helper/strings"
	"github.com/fthvgb1/wp-go/internal/cmd/reload"
	"github.com/fthvgb1/wp-go/internal/pkg/constraints"
	"github.com/fthvgb1/wp-go/internal/pkg/logs"
	"github.com/fthvgb1/wp-go/internal/wpconfig"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type Handle struct {
	Index             *IndexHandle
	Detail            *DetailHandle
	C                 *gin.Context
	theme             string
	Session           sessions.Session
	ginH              gin.H
	password          string
	scene             int
	Code              int
	Stats             int
	templ             string
	class             []string
	components        map[string][]Components
	themeMods         wpconfig.ThemeMods
	handleFns         map[int][]HandleCall
	err               error
	abort             bool
	componentsArgs    map[string]any
	componentFilterFn map[string][]func(*Handle, string) string
}

type HandlePlugins map[string]HandleFn[*Handle]

// Components Order 为执行顺序，降序执行
type Components struct {
	Str   string
	Fn    func(*Handle) string
	Order int
}

type HandleFn[T any] func(T)

type HandlePipeFn[T any] func(HandleFn[T], T)

type HandleCall struct {
	Fn    HandleFn[*Handle]
	Order int
}

func (h *Handle) ComponentFilterFn(name string) ([]func(*Handle, string) string, bool) {
	fn, ok := h.componentFilterFn[name]
	return fn, ok
}

func (h *Handle) PushComponentFilterFn(name string, fns ...func(*Handle, string) string) {
	h.componentFilterFn[name] = append(h.componentFilterFn[name], fns...)
}
func (h *Handle) ComponentFilterFnHook(name, s string) string {
	calls, ok := h.componentFilterFn[name]
	if ok {
		return slice.Reduce(calls, func(fn func(*Handle, string) string, r string) string {
			return fn(h, r)
		}, s)
	}
	return s
}

func (h *Handle) Abort() {
	h.abort = true
}

func (h *Handle) CommonThemeMods() wpconfig.ThemeMods {
	return h.themeMods
}

func (h *Handle) Err() error {
	return h.err
}

func (h *Handle) SetErr(err error) {
	h.err = err
}

func (h *Handle) Password() string {
	return h.password
}

func (h *Handle) SetTempl(templ string) {
	h.templ = templ
}

func (h *Handle) Scene() int {
	return h.scene
}

func (h *Handle) SetDatas(GinH gin.H) {
	maps.Merge(h.ginH, GinH)
}
func (h *Handle) SetData(k string, v any) {
	h.ginH[k] = v
}

func (h *Handle) PushClass(class ...string) {
	h.class = append(h.class, class...)
}

func GetComponentsArgs[T any](h *Handle, k string, defaults T) T {
	v, ok := h.componentsArgs[k]
	if ok {
		vv, ok := v.(T)
		if ok {
			return vv
		}
	}
	return defaults
}

func PushComponentsArgsForSlice[T any](h *Handle, name string, v ...T) {
	val, ok := h.componentsArgs[name]
	if !ok {
		var vv []T
		vv = append(vv, v...)
		h.componentsArgs[name] = vv
		return
	}
	vv, ok := val.([]T)
	if ok {
		vv = append(vv, v...)
		h.componentsArgs[name] = vv
	}
}
func SetComponentsArgsForMap[K comparable, V any](h *Handle, name string, key K, v V) {
	val, ok := h.componentsArgs[name]
	if !ok {
		vv := make(map[K]V)
		vv[key] = v
		h.componentsArgs[name] = vv
		return
	}
	vv, ok := val.(map[K]V)
	if ok {
		vv[key] = v
		h.componentsArgs[name] = vv
	}
}
func MergeComponentsArgsForMap[K comparable, V any](h *Handle, name string, m map[K]V) {
	val, ok := h.componentsArgs[name]
	if !ok {
		h.componentsArgs[name] = m
		return
	}
	vv, ok := val.(map[K]V)
	if ok {
		h.componentsArgs[name] = maps.Merge(vv, m)
	}
}

func (h *Handle) SetComponentsArgs(key string, value any) {
	h.componentsArgs[key] = value
}

func NewHandle(c *gin.Context, scene int, theme string) *Handle {
	mods, err := wpconfig.GetThemeMods(theme)
	logs.ErrPrintln(err, "获取mods失败")
	return &Handle{
		C:                 c,
		theme:             theme,
		Session:           sessions.Default(c),
		ginH:              gin.H{},
		scene:             scene,
		Stats:             constraints.Ok,
		themeMods:         mods,
		components:        make(map[string][]Components),
		handleFns:         make(map[int][]HandleCall),
		componentsArgs:    make(map[string]any),
		componentFilterFn: make(map[string][]func(*Handle, string) string),
	}
}

func (h *Handle) NewCacheComponent(name string, order int, fn func(handle *Handle) string) Components {
	return Components{Str: h.CacheStr(name, fn), Order: order}
}

func (h *Handle) PushHandleFn(statsOrScene int, fns ...HandleCall) {
	h.handleFns[statsOrScene] = append(h.handleFns[statsOrScene], fns...)
}

func (h *Handle) PushGroupHandleFn(statsOrScene, order int, fns ...HandleFn[*Handle]) {
	var calls []HandleCall
	for _, fn := range fns {
		calls = append(calls, HandleCall{fn, order})
	}
	h.handleFns[statsOrScene] = append(h.handleFns[statsOrScene], calls...)
}

func (h *Handle) AddCacheComponent(name string, fn func(*Handle) string) {
	h.ginH[name] = h.CacheStr(name, fn)
}

func (h *Handle) CacheStr(name string, fn func(*Handle) string) string {
	return reload.GetAnyValBy(name, func() string {
		return fn(h)
	})
}

func (h *Handle) PushHeadScript(fn ...Components) {
	h.PushComponents(constraints.HeadScript, fn...)
}
func (h *Handle) PushGroupHeadScript(order int, str ...string) {
	h.PushGroupComponentStrs(constraints.HeadScript, order, str...)
}
func (h *Handle) PushCacheGroupHeadScript(key string, order int, fns ...func(*Handle) string) {
	h.PushGroupCacheComponentFn(constraints.HeadScript, key, order, fns...)
}

func (h *Handle) PushFooterScript(fn ...Components) {
	h.PushComponents(constraints.FooterScript, fn...)
}

func (h *Handle) PushGroupFooterScript(order int, fns ...string) {
	h.PushGroupComponentStrs(constraints.FooterScript, order, fns...)
}

func (h *Handle) componentKey(name string) string {
	return fmt.Sprintf("theme_%d_%s", h.scene, name)
}

func (h *Handle) PushCacheGroupFooterScript(key string, order int, fns ...func(*Handle) string) {
	h.PushGroupCacheComponentFn(constraints.FooterScript, key, order, fns...)
}
func (h *Handle) PushGroupCacheComponentFn(name, key string, order int, fns ...func(*Handle) string) {
	v := reload.GetStrBy(key, "\n", h, fns...)
	h.PushGroupComponentStrs(name, order, v)
}

func (h *Handle) GetPassword() {
	pw := h.Session.Get("post_password")
	if pw != nil {
		h.password = pw.(string)
	}
}

func (h *Handle) ExecHandleFns() {
	calls, ok := h.handleFns[h.Stats]
	var fns []HandleCall
	if ok {
		fns = append(fns, calls...)
	}
	calls, ok = h.handleFns[h.scene]
	if ok {
		fns = append(fns, calls...)
	}
	calls, ok = h.handleFns[constraints.AllStats]
	if ok {
		fns = append(fns, calls...)
	}
	slice.Sort(fns, func(i, j HandleCall) bool {
		return i.Order > j.Order
	})
	for _, fn := range fns {
		fn.Fn(h)
		if h.abort {
			break
		}
	}
}

func (h *Handle) PreTemplate() {
	if h.templ == "" {
		h.templ = str.Join(h.theme, "/posts/index.gohtml")
		if h.scene == constraints.Detail {
			h.templ = str.Join(h.theme, "/posts/detail.gohtml")
		}
	}
}
func (h *Handle) PreCodeAndStats() {
	if h.Stats != 0 && h.Code != 0 {
		return
	}
	switch h.Stats {
	case constraints.Ok:
		h.Code = http.StatusOK
	case constraints.ParamError, constraints.Error404:
		h.Code = http.StatusNotFound
	case constraints.InternalErr:
		h.Code = http.StatusInternalServerError
	}
}

func (h *Handle) Render() {
	h.CommonComponents()
	h.ExecHandleFns()
}

func (h *Handle) CommonComponents() {
	h.AddCacheComponent("customLogo", CalCustomLogo)
	h.PushCacheGroupHeadScript("siteIconAndCustomCss", 0, CalSiteIcon, CalCustomCss)
	h.PushHandleFn(constraints.AllStats, NewHandleFn(func(h *Handle) {
		h.CalMultipleComponents()
		h.CalBodyClass()
	}, 10), NewHandleFn(func(h *Handle) {
		h.C.HTML(h.Code, h.templ, h.ginH)
	}, 0))
}

func (h *Handle) PushComponents(name string, components ...Components) {
	k := h.componentKey(name)
	h.components[k] = append(h.components[k], components...)
}

func (h *Handle) PushGroupComponentStrs(name string, order int, fns ...string) {
	var calls []Components
	for _, fn := range fns {
		calls = append(calls, Components{
			Str:   fn,
			Order: order,
		})
	}
	k := h.componentKey(name)
	h.components[k] = append(h.components[k], calls...)
}
func (h *Handle) PushGroupComponentFns(name string, order int, fns ...func(*Handle) string) {
	var calls []Components
	for _, fn := range fns {
		calls = append(calls, Components{
			Fn:    fn,
			Order: order,
		})
	}
	k := h.componentKey(name)
	h.components[k] = append(h.components[k], calls...)
}

func (h *Handle) CalMultipleComponents() {
	for k, ss := range h.components {
		slice.Sort(ss, func(i, j Components) bool {
			return i.Order > j.Order
		})
		v := strings.Join(slice.FilterAndMap(ss, func(t Components) (string, bool) {
			s := t.Str
			if s == "" && t.Fn != nil {
				s = t.Fn(h)
			}
			return s, s != ""
		}), "\n")
		kk := strings.Split(k, "_")
		key := kk[len(kk)-1]
		h.ginH[key] = v
	}
}

func NewHandleFn(fn HandleFn[*Handle], order int) HandleCall {
	return HandleCall{Fn: fn, Order: order}
}

// HandlePipe  方便把功能写在其它包里
func HandlePipe[T any](initial func(T), fns ...HandlePipeFn[T]) HandleFn[T] {
	return slice.ReverseReduce(fns, func(next HandlePipeFn[T], f func(t T)) func(t T) {
		return func(t T) {
			next(f, t)
		}
	}, initial)
}

func Render(h *Handle) {
	switch h.scene {
	case constraints.Detail:
		h.Detail.Render()
	default:
		h.Index.Render()
	}
}
