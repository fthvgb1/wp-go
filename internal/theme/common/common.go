package common

import (
	"github.com/fthvgb1/wp-go/helper/slice"
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
	Plugins   []HandlePluginFn[*Handle]
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

func (h *Handle) ExecHandlePlugin() {
	if len(h.Plugins) > 0 {
		HandlePlugin(h.Plugins, h)
	}
}

type HandleFn[T any] func(T)

type HandlePluginFn[T any] func(HandleFn[T], T) HandleFn[T]

// HandlePlugin 方便把功能写在其它包里
func HandlePlugin[T any](fns []HandlePluginFn[T], h T) HandleFn[T] {
	return slice.ReverseReduce(fns, func(t HandlePluginFn[T], r HandleFn[T]) HandleFn[T] {
		return t(r, h)
	}, func(t T) {})
}
