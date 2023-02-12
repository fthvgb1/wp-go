package common

import (
	"github.com/fthvgb1/wp-go/helper/slice"
	"github.com/fthvgb1/wp-go/internal/pkg/config"
	"github.com/fthvgb1/wp-go/internal/pkg/models"
	"github.com/fthvgb1/wp-go/internal/plugins"
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
		plugins.PasswordProjectTitle(&r)
		if h.Password != post.PostPassword {
			plugins.PasswdProjectContent(&r)
			return
		}
	}
	r = next(r)
	return
}

// Digest 生成摘要
func Digest(next Fn[models.Posts], h *Handle, post models.Posts) models.Posts {
	if post.PostExcerpt != "" {
		plugins.PostExcerpt(&post)
	} else {
		plugins.Digest(h.C, &post, config.GetConfig().DigestWordCount)
	}
	return next(post)
}
