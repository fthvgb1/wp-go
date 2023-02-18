package common

import (
	"context"
	"errors"
	"github.com/fthvgb1/wp-go/helper/maps"
	"github.com/fthvgb1/wp-go/helper/slice"
	str "github.com/fthvgb1/wp-go/helper/strings"
	"github.com/fthvgb1/wp-go/internal/pkg/config"
	"github.com/fthvgb1/wp-go/internal/pkg/constraints"
	"github.com/fthvgb1/wp-go/internal/pkg/logs"
	"github.com/fthvgb1/wp-go/internal/pkg/models"
	"github.com/fthvgb1/wp-go/internal/plugins"
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

func (h *Handle) GetPassword() {
	pw := h.Session.Get("post_password")
	if pw != nil {
		h.Password = pw.(string)
	}
}

func (i *IndexHandle) ExecListPagePlugin(m map[string]Plugin[models.Posts, *Handle], calls ...func(*models.Posts)) {

	pluginConf := config.GetConfig().ListPagePlugins

	plugin := GetListPostPlugins(pluginConf, m)

	i.GinH["posts"] = slice.Map(i.Posts, PluginFn[models.Posts, *Handle](plugin, i.Handle, Defaults(calls...)))

}

func ListPostPlugins() map[string]Plugin[models.Posts, *Handle] {
	return maps.Copy(pluginFns)
}

func Defaults(call ...func(*models.Posts)) Fn[models.Posts] {
	return func(posts models.Posts) models.Posts {
		for _, fn := range call {
			fn(&posts)
		}
		return posts
	}
}

func Default[T any](t T) T {
	return t
}

func ProjectTitle(t models.Posts) models.Posts {
	if t.PostPassword != "" {
		plugins.PasswordProjectTitle(&t)
	}
	return t
}

func GetListPostPlugins(name []string, m map[string]Plugin[models.Posts, *Handle]) []Plugin[models.Posts, *Handle] {
	return slice.FilterAndMap(name, func(t string) (Plugin[models.Posts, *Handle], bool) {
		v, ok := m[t]
		if ok {
			return v, true
		}
		logs.ErrPrintln(errors.New(str.Join("插件", t, "不存在")), "")
		return nil, false
	})
}

func DigestsAndOthers(ctx context.Context, calls ...func(*models.Posts)) Fn[models.Posts] {
	return func(post models.Posts) models.Posts {
		if post.PostExcerpt != "" {
			plugins.PostExcerpt(&post)
		} else {
			plugins.Digest(ctx, &post, config.GetConfig().DigestWordCount)
		}
		if len(calls) > 0 {
			for _, call := range calls {
				call(&post)
			}
		}
		return post
	}
}
