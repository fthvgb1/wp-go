package common

import (
	"context"
	"errors"
	"github.com/fthvgb1/wp-go/helper/maps"
	"github.com/fthvgb1/wp-go/helper/slice"
	str "github.com/fthvgb1/wp-go/helper/strings"
	"github.com/fthvgb1/wp-go/internal/pkg/config"
	"github.com/fthvgb1/wp-go/internal/pkg/logs"
	"github.com/fthvgb1/wp-go/internal/pkg/models"
	"github.com/fthvgb1/wp-go/internal/plugins"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type Handle struct {
	C        *gin.Context
	Session  sessions.Session
	GinH     gin.H
	Password string
	Scene    int
	Code     int
	Stats    int
}

func (h Handle) Detail() {

}

func (h Handle) Index() {

}

func (h Handle) ExecListPagePlugin(m map[string]Plugin[models.Posts], calls ...func(*models.Posts)) {

	pluginConf := config.GetConfig().ListPagePlugins

	plugin := GetListPostPlugins(pluginConf, m)

	posts, ok := maps.GetStrMapAnyVal[[]models.Posts](h.GinH, "posts")

	if ok {
		h.GinH["posts"] = slice.Map(posts, PluginFn[models.Posts](plugin, h, Defaults(calls...)))
	}
}

/*func (h Handle) Pagination(paginate pagination)  {

}*/

type Fn[T any] func(T) T
type Plugin[T any] func(next Fn[T], h Handle, t T) T

func PluginFn[T any](a []Plugin[T], h Handle, fn Fn[T]) Fn[T] {
	return slice.ReverseReduce(a, func(next Plugin[T], forward Fn[T]) Fn[T] {
		return func(t T) T {
			return next(forward, h, t)
		}
	}, fn)
}

var pluginFns = map[string]Plugin[models.Posts]{
	"passwordProject": PasswordProject,
	"digest":          Digest,
}

func ListPostPlugins() map[string]Plugin[models.Posts] {
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

func GetListPostPlugins(name []string, m map[string]Plugin[models.Posts]) []Plugin[models.Posts] {
	return slice.FilterAndMap(name, func(t string) (Plugin[models.Posts], bool) {
		v, ok := m[t]
		if ok {
			return v, true
		}
		logs.ErrPrintln(errors.New(str.Join("插件", t, "不存在")), "")
		return nil, false
	})
}

// PasswordProject 标题和内容密码保护
func PasswordProject(next Fn[models.Posts], h Handle, post models.Posts) (r models.Posts) {
	r = post
	if post.PostPassword != "" {
		plugins.PasswordProjectTitle(&r)
		if h.Password != post.PostPassword {
			plugins.PasswdProjectContent(&r)
			return
		}
	}
	r = next(r)
	return
}

func ProjectTitle(t models.Posts) models.Posts {
	if t.PostPassword != "" {
		plugins.PasswordProjectTitle(&t)
	}
	return t
}

// Digest 生成摘要
func Digest(next Fn[models.Posts], h Handle, post models.Posts) models.Posts {
	if post.PostExcerpt != "" {
		plugins.PostExcerpt(&post)
	} else {
		plugins.Digest(h.C, &post, config.GetConfig().DigestWordCount)
	}
	return next(post)
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
