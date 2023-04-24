package wp

import (
	"github.com/fthvgb1/wp-go/helper/maps"
	"github.com/fthvgb1/wp-go/helper/slice"
	str "github.com/fthvgb1/wp-go/helper/strings"
	"github.com/fthvgb1/wp-go/internal/cmd/reload"
	"github.com/fthvgb1/wp-go/internal/pkg/constraints"
	"github.com/fthvgb1/wp-go/internal/pkg/logs"
	"github.com/fthvgb1/wp-go/internal/wpconfig"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"html/template"
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
	scene             string
	Code              int
	Stats             string
	templ             string
	components        map[string][]Components[string]
	themeMods         wpconfig.ThemeMods
	handlers          map[string]map[string][]HandleCall
	handleHook        map[string][]func(HandleCall) (HandleCall, bool)
	err               error
	abort             bool
	stopPipe          bool
	componentsArgs    map[string]any
	componentFilterFn map[string][]func(*Handle, string, ...any) string
	template          *template.Template
}

func (h *Handle) SetTemplate(template *template.Template) {
	h.template = template
}

func (h *Handle) Template() *template.Template {
	return h.template
}

type HandlePlugins map[string]HandleFn[*Handle]

// Components Order 为执行顺序，降序执行
type Components[T any] struct {
	Val      T
	Fn       func(*Handle) T
	Order    int
	CacheKey string
}

type HandleFn[T any] func(T)

type HandleCall struct {
	Fn    HandleFn[*Handle]
	Order int
	Name  string
}

func InitThemeArgAndConfig(fn func(*Handle), h *Handle) {
	var inited = false
	hh := reload.GetAnyValBys("themeArgAndConfig", h, func(h *Handle) Handle {
		h.components = make(map[string][]Components[string])
		h.componentsArgs = make(map[string]any)
		h.componentFilterFn = make(map[string][]func(*Handle, string, ...any) string)
		h.handlers = make(map[string]map[string][]HandleCall)
		h.handleHook = make(map[string][]func(HandleCall) (HandleCall, bool))
		h.ginH = gin.H{}
		fn(h)
		inited = true
		return *h
	})
	m := make(map[string][]Components[string])
	for k, v := range hh.components {
		vv := make([]Components[string], len(v))
		copy(vv, v)
		m[k] = vv
	}
	h.components = m
	h.ginH = maps.Copy(hh.ginH)
	h.ginH["calPostClass"] = postClass(h)
	h.ginH["calBodyClass"] = bodyClass(h)
	h.ginH["customLogo"] = customLogo(h)
	if inited {
		return
	}
	h.Index.postsPlugin = hh.Index.postsPlugin
	h.Index.pageEle = hh.Index.pageEle
	h.Detail.CommentRender = hh.Detail.CommentRender
	h.handlers = hh.handlers
	h.handleHook = hh.handleHook
	h.componentsArgs = hh.componentsArgs
	h.componentFilterFn = hh.componentFilterFn
}

func (h *Handle) ComponentFilterFn(name string) ([]func(*Handle, string, ...any) string, bool) {
	fn, ok := h.componentFilterFn[name]
	return fn, ok
}

func (h *Handle) PushComponentFilterFn(name string, fns ...func(*Handle, string, ...any) string) {
	h.componentFilterFn[name] = append(h.componentFilterFn[name], fns...)
}
func (h *Handle) ComponentFilterFnHook(name, s string, args ...any) string {
	calls, ok := h.componentFilterFn[name]
	if ok {
		return slice.Reduce(calls, func(fn func(*Handle, string, ...any) string, r string) string {
			return fn(h, r, args...)
		}, s)
	}
	return s
}

func (h *Handle) Abort() {
	h.abort = true
}
func (h *Handle) StopPipe() {
	h.stopPipe = true
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

func (h *Handle) Scene() string {
	return h.scene
}

func (h *Handle) SetDatas(GinH gin.H) {
	maps.Merge(h.ginH, GinH)
}
func (h *Handle) SetData(k string, v any) {
	h.ginH[k] = v
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

func NewHandle(c *gin.Context, scene string, theme string) *Handle {
	mods, err := wpconfig.GetThemeMods(theme)
	logs.IfError(err, "获取mods失败")
	return &Handle{
		C:         c,
		theme:     theme,
		Session:   sessions.Default(c),
		scene:     scene,
		Stats:     constraints.Ok,
		themeMods: mods,
	}
}

func (h *Handle) NewCacheComponent(name string, order int, fn func(handle *Handle) string) Components[string] {
	return Components[string]{Fn: fn, CacheKey: name, Order: order}
}

func (h *Handle) AddCacheComponent(name string, fn func(*Handle) string) {
	h.components[name] = append(h.components[name], h.NewCacheComponent(name, 10, fn))
}

func (h *Handle) PushHeadScript(fn ...Components[string]) {
	h.PushComponents(constraints.HeadScript, fn...)
}
func (h *Handle) PushGroupHeadScript(order int, str ...string) {
	h.PushGroupComponentStrs(constraints.HeadScript, order, str...)
}
func (h *Handle) PushCacheGroupHeadScript(key string, order int, fns ...func(*Handle) string) {
	h.PushGroupCacheComponentFn(constraints.HeadScript, key, order, fns...)
}

func (h *Handle) PushFooterScript(fn ...Components[string]) {
	h.PushComponents(constraints.FooterScript, fn...)
}

func (h *Handle) PushGroupFooterScript(order int, fns ...string) {
	h.PushGroupComponentStrs(constraints.FooterScript, order, fns...)
}

func (h *Handle) PushCacheGroupFooterScript(key string, order int, fns ...func(*Handle) string) {
	h.PushGroupCacheComponentFn(constraints.FooterScript, key, order, fns...)
}
func (h *Handle) PushGroupCacheComponentFn(name, key string, order int, fns ...func(*Handle) string) {
	h.PushComponents(name, h.NewCacheComponent(key, order, func(h *Handle) string {
		return strings.Join(slice.Map(fns, func(t func(*Handle) string) string {
			return t(h)
		}), "\n")
	}))
}

func (h *Handle) GetPassword() {
	pw := h.Session.Get("post_password")
	if pw != nil {
		h.password = pw.(string)
	}
}

func PreTemplate(h *Handle) {
	if h.templ == "" {
		h.templ = str.Join(h.theme, "/posts/index.gohtml")
		if h.scene == constraints.Detail {
			h.templ = str.Join(h.theme, "/posts/detail.gohtml")
		}
	}
}
func PreCodeAndStats(h *Handle) {
	if h.Stats != "" && h.Code != 0 {
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

func (h *Handle) CommonComponents() {
	h.PushCacheGroupHeadScript("siteIconAndCustomCss", 0, CalSiteIcon, CalCustomCss)
	h.PushRender(constraints.AllStats, NewHandleFn(CalComponents, 10, "wp.CalComponents"))
	h.PushRender(constraints.AllStats, NewHandleFn(RenderTemplate, 0, "wp.RenderTemplate"))
}

func RenderTemplate(h *Handle) {
	h.C.HTML(h.Code, h.templ, h.ginH)
	h.StopPipe()
}

func (h *Handle) PushComponents(name string, components ...Components[string]) {
	h.components[name] = append(h.components[name], components...)
}

func (h *Handle) PushGroupComponentStrs(name string, order int, str ...string) {
	var calls []Components[string]
	for _, fn := range str {
		calls = append(calls, Components[string]{
			Val:   fn,
			Order: order,
		})
	}
	h.components[name] = append(h.components[name], calls...)
}
func (h *Handle) PushGroupComponentFns(name string, order int, fns ...func(*Handle) string) {
	var calls []Components[string]
	for _, fn := range fns {
		calls = append(calls, Components[string]{
			Fn:    fn,
			Order: order,
		})
	}
	h.components[name] = append(h.components[name], calls...)
}

func CalComponents(h *Handle) {
	for k, ss := range h.components {
		slice.Sort(ss, func(i, j Components[string]) bool {
			return i.Order > j.Order
		})
		var s []string
		for _, component := range ss {
			if component.Val != "" {
				s = append(s, component.Val)
				continue
			}
			if component.Fn != nil {
				v := ""
				if component.CacheKey != "" {
					v = reload.SafetyMapBy("calComponent", component.CacheKey, h, component.Fn)
				} else {
					v = component.Fn(h)
				}
				if v != "" {
					s = append(s, v)
				}
			}
		}
		h.ginH[k] = strings.Join(s, "\n")
	}
}

func NewHandleFn(fn HandleFn[*Handle], order int, name string) HandleCall {
	return HandleCall{Fn: fn, Order: order, Name: name}
}

func NothingToDo(*Handle) {

}
