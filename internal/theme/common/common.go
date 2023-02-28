package common

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
}

func (h *Handle) CommonThemeMods() wpconfig.ThemeMods {
	return h.themeMods
}

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

func (h *Handle) PushHandleFn(stats int, fns ...HandleCall) {
	h.handleFns[stats] = append(h.handleFns[stats], fns...)
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
	h.components[constraints.HeadScript] = append(h.components[constraints.HeadScript], fn...)
}
func (h *Handle) PushFooterScript(fn ...Components) {
	h.components[constraints.FooterScript] = append(h.components[constraints.FooterScript], fn...)
}

func (h *Handle) GetPassword() {
	pw := h.Session.Get("post_password")
	if pw != nil {
		h.password = pw.(string)
	}
}

func (h *Handle) ExecHandleFns() {
	calls, ok := h.handleFns[h.Stats]
	if ok {
		slice.SortSelf(calls, func(i, j HandleCall) bool {
			return i.Order > j.Order
		})
		for _, call := range calls {
			call.Fn(h)
		}
	}
	fns, ok := h.handleFns[constraints.AllStats]
	if ok {
		for _, fn := range fns {
			fn.Fn(h)
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

	h.PushHeadScript(Components{CalSiteIcon, 100}, Components{CalCustomCss, 0})

	h.PushHandleFn(constraints.AllStats, NewHandleFn(func(h *Handle) {
		h.CalMultipleComponents()
		h.CalBodyClass()
	}, 5))

	h.ExecHandleFns()

	h.C.HTML(h.Code, h.templ, h.ginH)
}

func (h *Handle) PushComponents(name string, components ...Components) {
	h.components[name] = append(h.components[name], components...)
}

func (h *Handle) CalMultipleComponents() {
	for k, ss := range h.components {
		v, ok := reload.GetStr(k)
		if !ok {
			slice.SortSelf(ss, func(i, j Components) bool {
				return i.Order > j.Order
			})
			v = strings.Join(slice.FilterAndMap(ss, func(t Components) (string, bool) {
				s := t.Fn(h)
				return s, s != ""
			}), "\n")
			reload.SetStr(k, v)
		}
		h.ginH[k] = v
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
