package common

import (
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
	Theme      string
	Session    sessions.Session
	GinH       gin.H
	Password   string
	Scene      int
	Code       int
	Stats      int
	Templ      string
	Class      []string
	Components map[string][]Components
	ThemeMods  wpconfig.ThemeMods
	HandleFns  map[int][]HandleCall
	Error      error
}

func NewHandle(c *gin.Context, scene int, theme string) *Handle {
	mods, err := wpconfig.GetThemeMods(theme)
	logs.ErrPrintln(err, "获取mods失败")
	return &Handle{
		C:          c,
		Theme:      theme,
		Session:    sessions.Default(c),
		GinH:       gin.H{},
		Scene:      scene,
		Stats:      constraints.Ok,
		ThemeMods:  mods,
		Components: make(map[string][]Components),
		HandleFns:  make(map[int][]HandleCall),
	}
}

// Components Order 为执行顺序，降序执行
type Components struct {
	Fn    func(*Handle) string
	Order int
}

func NewComponents(fn func(*Handle) string, order int) Components {
	return Components{Fn: fn, Order: order}
}

func (h *Handle) PushHandleFn(stats int, fns ...HandleCall) {
	h.HandleFns[stats] = append(h.HandleFns[stats], fns...)
}

func (h *Handle) AddComponent(name string, fn func(*Handle) string) {
	v, ok := reload.GetStr(name)
	if !ok {
		v = fn(h)
		reload.SetStr(name, v)
	}
	h.GinH[name] = v
}

func (h *Handle) PushHeadScript(fn ...Components) {
	h.Components[constraints.HeadScript] = append(h.Components[constraints.HeadScript], fn...)
}
func (h *Handle) PushFooterScript(fn ...Components) {
	h.Components[constraints.FooterScript] = append(h.Components[constraints.FooterScript], fn...)
}

func (h *Handle) GetPassword() {
	pw := h.Session.Get("post_password")
	if pw != nil {
		h.Password = pw.(string)
	}
}

func (h *Handle) ExecHandleFns() {
	calls, ok := h.HandleFns[h.Stats]
	if ok {
		slice.SortSelf(calls, func(i, j HandleCall) bool {
			return i.Order > j.Order
		})
		for _, call := range calls {
			call.Fn(h)
		}
	}
	fns, ok := h.HandleFns[constraints.AllStats]
	if ok {
		for _, fn := range fns {
			fn.Fn(h)
		}
	}

}

func (h *Handle) PreTemplate() {
	if h.Templ == "" {
		h.Templ = str.Join(h.Theme, "/posts/index.gohtml")
		if h.Scene == constraints.Detail {
			h.Templ = str.Join(h.Theme, "/posts/detail.gohtml")
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
	h.ExecHandleFns()
	h.PushHeadScript(Components{CalSiteIcon, 10}, Components{CalCustomCss, -1})
	h.AddComponent("customLogo", CalCustomLogo)
	h.CalMultipleComponents()
	h.CalBodyClass()
	h.C.HTML(h.Code, h.Templ, h.GinH)
}

func (h *Handle) PushComponents(name string, components ...Components) {
	h.Components[name] = append(h.Components[name], components...)
}

func (h *Handle) CalMultipleComponents() {
	for k, ss := range h.Components {
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
		h.GinH[k] = v
	}
}

type HandleFn[T any] func(T)

type HandlePipeFn[T any] func(HandleFn[T], T)

type HandleCall struct {
	Fn    HandleFn[*Handle]
	Order int
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
	switch h.Scene {
	case constraints.Detail:
		h.Detail.Render()
	default:
		h.Index.Render()
	}
}
