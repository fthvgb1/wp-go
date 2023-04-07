package wp

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
	"github.com/fthvgb1/wp-go/internal/plugins/wpposts"
)

type Fn[T any] func(T) T
type Plugin[T, H any] func(initialFn Fn[T], h H, t T) T

func PluginFn[T, H any](a []Plugin[T, H], h H, fn Fn[T]) Fn[T] {
	return slice.ReverseReduce(a, func(next Plugin[T, H], fn Fn[T]) Fn[T] {
		return func(t T) T {
			return next(fn, h, t)
		}
	}, fn)
}

var pluginFns = map[string]Plugin[models.Posts, *Handle]{
	"passwordProject": PasswordProject,
	"digest":          Digest,
}

// PasswordProject 标题和内容密码保护
func PasswordProject(next Fn[models.Posts], h *Handle, post models.Posts) (r models.Posts) {
	r = post
	if post.PostPassword != "" {
		wpposts.PasswordProjectTitle(&r)
		if h.password != post.PostPassword {
			wpposts.PasswdProjectContent(&r)
			return
		}
	}
	r = next(r)
	return
}

// Digest 生成摘要 注意放到最后，不继续往下执行
func Digest(next Fn[models.Posts], h *Handle, post models.Posts) models.Posts {
	if post.PostExcerpt != "" {
		plugins.PostExcerpt(&post)
	} else {
		plugins.Digest(h.C, &post, config.GetConfig().DigestWordCount)
	}
	return post
}

func ListPostPlugins() map[string]Plugin[models.Posts, *Handle] {
	return maps.Copy(pluginFns)
}

func ProjectTitle(t models.Posts) models.Posts {
	if t.PostPassword != "" {
		wpposts.PasswordProjectTitle(&t)
	}
	return t
}

func GetListPostPlugins(name []string, m map[string]Plugin[models.Posts, *Handle]) []Plugin[models.Posts, *Handle] {
	return slice.FilterAndMap(name, func(t string) (Plugin[models.Posts, *Handle], bool) {
		v, ok := m[t]
		if ok {
			return v, true
		}
		logs.IfError(errors.New(str.Join("插件", t, "不存在")), "")
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

func (i *IndexHandle) ExecListPagePlugin(m map[string]Plugin[models.Posts, *Handle], calls ...func(*models.Posts)) {

	pluginConf := config.GetConfig().ListPagePlugins

	plugin := GetListPostPlugins(pluginConf, m)

	i.ginH["posts"] = slice.Map(i.Posts, PluginFn[models.Posts, *Handle](plugin, i.Handle, Defaults(calls...)))

}

func Defaults(call ...func(*models.Posts)) Fn[models.Posts] {
	return func(posts models.Posts) models.Posts {
		for _, fn := range call {
			fn(&posts)
		}
		return posts
	}
}
