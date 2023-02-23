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
)

type Handle struct {
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
	ThemeMods wpconfig.ThemeMods
	HandleFns []func(*Handle)
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
	}
}

func (h *Handle) AutoCal(name string, fn func() string) {
	v, ok := reload.GetStr(name)
	if !ok {
		v = fn()
		reload.SetStr(name, v)
	}
	h.GinH[name] = v
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
	h.AutoCal("siteIcon", h.CalSiteIcon)
	h.AutoCal("customLogo", h.CalCustomLogo)
	h.AutoCal("customCss", h.CalCustomCss)
	h.CalBodyClass()

	h.C.HTML(h.Code, h.Templ, h.GinH)
}

type HandleFn[T any] func(T)

type HandlePipeFn[T any] func(HandleFn[T], T)

// HandlePipe  方便把功能写在其它包里
func HandlePipe[T any](fns []HandlePipeFn[T], initial func(T)) HandleFn[T] {
	return slice.ReverseReduce(fns, func(next HandlePipeFn[T], f func(t T)) func(t T) {
		return func(t T) {
			next(f, t)
		}
	}, initial)
}
