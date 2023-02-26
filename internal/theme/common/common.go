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
	Index     *IndexHandle
	Detail    *DetailHandle
	C         *gin.Context
	Theme     string
	Session   sessions.Session
	GinH      gin.H
	Password  string
	Scene     int
	Code      int
	Stats     int
	Templ     string
	Class     []string
	Scripts   map[string][]func(*Handle) string
	ThemeMods wpconfig.ThemeMods
	HandleFns []HandleFn[*Handle]
}

func NewHandle(c *gin.Context, scene int, theme string) *Handle {
	mods, err := wpconfig.GetThemeMods(theme)
	logs.ErrPrintln(err, "获取mods失败")
	return &Handle{
		C:         c,
		Theme:     theme,
		Session:   sessions.Default(c),
		GinH:      gin.H{},
		Scene:     scene,
		Code:      http.StatusOK,
		Stats:     constraints.Ok,
		ThemeMods: mods,
		Scripts:   make(map[string][]func(*Handle) string),
	}
}

func (h *Handle) PushHandleFn(fns ...HandleFn[*Handle]) {
	h.HandleFns = append(h.HandleFns, fns...)
}

func (h *Handle) PlushComponent(name string, fn func(*Handle) string) {
	v, ok := reload.GetStr(name)
	if !ok {
		v = fn(h)
		reload.SetStr(name, v)
	}
	h.GinH[name] = v
}

func (h *Handle) PushHeadScript(name string, fn ...func(*Handle) string) {
	h.Scripts[name] = append(h.Scripts[name], fn...)
}

func Default[T any](t T) T {
	return t
}

func (h *Handle) GetPassword() {
	pw := h.Session.Get("post_password")
	if pw != nil {
		h.Password = pw.(string)
	}
}

func (h *Handle) Render() {
	if h.Templ == "" {
		h.Templ = str.Join(h.Theme, "/posts/index.gohtml")
		if h.Scene == constraints.Detail {
			h.Templ = str.Join(h.Theme, "/posts/detail.gohtml")
		}
	}
	for _, fn := range h.HandleFns {
		fn(h)
	}
	h.PushHeadScript(constraints.HeadScript, CalSiteIcon, CalCustomCss)
	h.PlushComponent("customLogo", CalCustomLogo)
	h.CalMultipleScript()
	h.CalBodyClass()

	h.C.HTML(h.Code, h.Templ, h.GinH)
}

func (h *Handle) CalMultipleScript() {
	for k, ss := range h.Scripts {
		v, ok := reload.GetStr(k)
		if !ok {
			v = strings.Join(slice.FilterAndMap(ss, func(t func(*Handle) string) (string, bool) {
				s := t(h)
				return s, s != ""
			}), "\n")
			reload.SetStr(k, v)
		}
		h.GinH[k] = v
	}
}

type HandleFn[T any] func(T)

type HandlePipeFn[T any] func(HandleFn[T], T)

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
