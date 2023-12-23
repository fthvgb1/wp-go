package wp

import (
	"errors"
	"github.com/fthvgb1/wp-go/app/pkg/config"
	"github.com/fthvgb1/wp-go/app/pkg/constraints"
	"github.com/fthvgb1/wp-go/app/pkg/logs"
	"github.com/fthvgb1/wp-go/app/plugins/wphandle/apply"
	"github.com/fthvgb1/wp-go/app/wpconfig"
	"github.com/fthvgb1/wp-go/cache/reload"
	"github.com/fthvgb1/wp-go/helper/maps"
	str "github.com/fthvgb1/wp-go/helper/strings"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
	"os"
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
	components        map[string]map[string][]Components[string]
	componentHook     map[string][]func(Components[string]) (Components[string], bool)
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

func (h *Handle) Theme() string {
	return h.theme
}

func (h *Handle) GinH() gin.H {
	return h.ginH
}

func (h *Handle) SetScene(scene string) {
	h.scene = scene
}

func (h *Handle) Components() map[string]map[string][]Components[string] {
	return h.components
}

func (h *Handle) ComponentHook() map[string][]func(Components[string]) (Components[string], bool) {
	return h.componentHook
}

func (h *Handle) Handlers() map[string]map[string][]HandleCall {
	return h.handlers
}

func (h *Handle) HandleHook() map[string][]func(HandleCall) (HandleCall, bool) {
	return h.handleHook
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
	Name   string
	Val    T
	Fn     func(*Handle) T
	Order  float64
	Cached bool
}

type HandleFn[T any] func(T)

type HandleCall struct {
	Fn    HandleFn[*Handle]
	Order float64
	Name  string
}

func InitHandle(fn func(*Handle), h *Handle) {
	var inited = false
	hh := reload.GetAnyValBys("themeArgAndConfig", h, func(h *Handle) (Handle, bool) {
		h.components = make(map[string]map[string][]Components[string])
		h.componentsArgs = make(map[string]any)
		h.componentFilterFn = make(map[string][]func(*Handle, string, ...any) string)
		h.handlers = make(map[string]map[string][]HandleCall)
		h.handleHook = make(map[string][]func(HandleCall) (HandleCall, bool))
		h.ginH = gin.H{}
		fnMap.Flush()
		fnHook.Flush()
		fn(h)
		v := apply.UsePlugins()
		pluginFn, ok := v.(func(*Handle))
		if ok {
			pluginFn(h)
		}
		h.C.Set("inited", true)
		inited = true
		return *h, true
	})
	h.ginH = maps.Copy(hh.ginH)
	h.ginH["calPostClass"] = postClass(h)
	h.ginH["calBodyClass"] = bodyClass(h)
	h.ginH["customLogo"] = customLogo(h)
	if inited {
		return
	}
	h.components = hh.components
	h.Index.postsPlugin = hh.Index.postsPlugin
	h.Index.pageEle = hh.Index.pageEle
	h.Detail.CommentRender = hh.Detail.CommentRender
	h.Detail.CommentPageEle = hh.Detail.CommentPageEle
	h.handlers = hh.handlers
	h.handleHook = hh.handleHook
	h.componentHook = hh.componentHook
	h.componentsArgs = hh.componentsArgs
	h.componentFilterFn = hh.componentFilterFn
	h.C.Set("inited", true)
}

func (h *Handle) Abort() {
	h.abort = true
	h.stopPipe = true
}
func (h *Handle) StopPipe() {
	h.stopPipe = true
}
func (h *Handle) StopHandle() {
	h.abort = true
}

func (h *Handle) CommonThemeMods() wpconfig.ThemeMods {
	return h.themeMods
}

func (h *Handle) Err() error {
	return h.err
}

func (h *Handle) SetErr(err error) {
	h.err = errors.Join(err)
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

func (h *Handle) GetPassword() string {
	if h.password != "" {
		return h.password
	}
	pw := h.Session.Get("post_password")
	if pw != nil {
		h.password = pw.(string)
	}
	return h.password
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

var htmlContentType = []string{"text/html; charset=utf-8"}

func (h *Handle) RenderHtml(t *template.Template, statsCode int, name string) {
	header := h.C.Writer.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = htmlContentType
	}
	h.C.Status(statsCode)
	err := t.ExecuteTemplate(h.C.Writer, name, h.ginH)
	h.Abort()
	if err != nil {
		panic(err)
	}
}

func (h *Handle) PushHandlers(pipeScene string, call HandleCall, statsOrScene ...string) {
	for _, s := range statsOrScene {
		h.PushHandler(pipeScene, s, call)
	}
}

func (h *Handle) CommonComponents() {
	h.PushCacheGroupHeadScript(constraints.AllScene, "siteIconAndCustomCss", 0, CalSiteIcon, CalCustomCss)
	h.PushRender(constraints.AllStats, NewHandleFn(CalComponents, 10.001, "wp.CalComponents"))
	h.PushRender(constraints.AllStats, NewHandleFn(PreRenderTemplate, 0, "wp.PreRenderTemplate"))
	AdditionScript(h)
}

func AdditionScript(h *Handle) {
	s := config.GetConfig().ExternScript
	if len(s) < 1 {
		return
	}
	fn := func(f, name string) {
		if f == "" {
			return
		}
		ss, err := os.ReadFile(f)
		if err != nil {
			logs.Error(err, str.Join("解析", name, "失败"), f)
		} else {
			h.PushComponents(constraints.AllScene, constraints.HeadScript, NewComponent(name, string(ss), false, 0, nil))
		}
	}
	switch len(s) {
	case 1:
		fn(s[0], "externHead")
	case 2:
		fn(s[0], "externHead")
		fn(s[1], "externFooter")
	}
}

func PreRenderTemplate(h *Handle) {
	h.C.HTML(h.Code, h.templ, h.ginH)
	h.Abort()
}

func NewHandleFn(fn HandleFn[*Handle], order float64, name string) HandleCall {
	return HandleCall{Fn: fn, Order: order, Name: name}
}

func NothingToDo(h *Handle) {
	h.Abort()
}

func (h *Handle) IsHttps() bool {
	if h.C.Request.TLS != nil {
		return true
	}
	return "https" == strings.ToLower(h.C.Request.Header.Get("X-Forwarded-Proto"))
}
