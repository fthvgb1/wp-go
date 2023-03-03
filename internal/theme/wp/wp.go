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
	Index      *IndexHandle
	Detail     *DetailHandle
	C          *gin.Context
	theme      string
	Session    sessions.Session
	ginH       gin.H
	password   string
	scene      int
	Code       int
	Stats      int
	templ      string
	class      []string
	components map[string][]Components
	themeMods  wpconfig.ThemeMods
	handleFns  map[int][]HandleCall
	err        error
	abort      bool
}

type HandlePlugins map[string]HandleFn[*Handle]

// Components Order 为执行顺序，降序执行
type Components struct {
	Fn    func(*Handle) string
	Order int
}

type HandleFn[T any] func(T)

type HandlePipeFn[T any] func(HandleFn[T], T)

type HandleCall struct {
	Fn    HandleFn[*Handle]
	Order int
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

func NewHandle(c *gin.Context, scene int, theme string) *Handle {
	mods, err := wpconfig.GetThemeMods(theme)
	logs.ErrPrintln(err, "获取mods失败")
	return &Handle{
		C:          c,
		theme:      theme,
		Session:    sessions.Default(c),
		ginH:       gin.H{},
		scene:      scene,
		Stats:      constraints.Ok,
		themeMods:  mods,
		components: make(map[string][]Components),
		handleFns:  make(map[int][]HandleCall),
	}
}

func NewComponents(fn func(*Handle) string, order int) Components {
	return Components{Fn: fn, Order: order}
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

func (h *Handle) AddComponent(name string, fn func(*Handle) string) {
	v, ok := reload.GetStr(name)
	if !ok {
		v = fn(h)
		reload.SetStr(name, v)
	}
	h.ginH[name] = v
}

func (h *Handle) PushHeadScript(fn ...Components) {
	h.PushComponents(constraints.HeadScript, fn...)
}
func (h *Handle) PushGroupHeadScript(order int, fns ...func(*Handle) string) {
	h.PushGroupComponents(constraints.HeadScript, order, fns...)
}
func (h *Handle) PushFooterScript(fn ...Components) {
	h.PushComponents(constraints.FooterScript, fn...)
}

func (h *Handle) PushGroupFooterScript(order int, fns ...func(*Handle) string) {
	h.PushGroupComponents(constraints.FooterScript, order, fns...)
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
	h.PreCodeAndStats()
	h.PreTemplate()
	h.AddComponent("customLogo", CalCustomLogo)

	h.PushGroupHeadScript(0, CalSiteIcon, CalCustomCss)

	h.PushHandleFn(constraints.AllStats, NewHandleFn(func(h *Handle) {
		h.CalMultipleComponents()
		h.CalBodyClass()
	}, 10), NewHandleFn(func(h *Handle) {
		h.C.HTML(h.Code, h.templ, h.ginH)
	}, 0))

	h.ExecHandleFns()

}

func (h *Handle) PushComponents(name string, components ...Components) {
	k := h.componentKey(name)
	h.components[k] = append(h.components[k], components...)
}

func (h *Handle) PushGroupComponents(name string, order int, fns ...func(*Handle) string) {
	var calls []Components
	for _, fn := range fns {
		calls = append(calls, Components{fn, order})
	}
	k := h.componentKey(name)
	h.components[k] = append(h.components[k], calls...)
}

func (h *Handle) componentKey(name string) string {
	return fmt.Sprintf("%d_%s", h.scene, name)
}

func (h *Handle) CalMultipleComponents() {
	for k, ss := range h.components {
		v, ok := reload.GetStr(k)
		if !ok {
			slice.Sort(ss, func(i, j Components) bool {
				return i.Order > j.Order
			})
			v = strings.Join(slice.FilterAndMap(ss, func(t Components) (string, bool) {
				s := t.Fn(h)
				return s, s != ""
			}), "\n")
			reload.SetStr(k, v)
		}
		key := strings.Split(k, "_")[1]
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
