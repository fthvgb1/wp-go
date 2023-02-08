package common

import (
	"github.com/fthvgb1/wp-go/helper/slice"
	"github.com/fthvgb1/wp-go/internal/pkg/config"
	"github.com/fthvgb1/wp-go/internal/pkg/models"
	"github.com/fthvgb1/wp-go/internal/plugins"
	"github.com/gin-gonic/gin"
)

type Handle struct {
	C        *gin.Context
	GinH     gin.H
	Password string
	Scene    int
	Code     int
	Stats    int
}

type Fn[T any] func(T) T
type Plugin[T any] func(next Fn[T], h Handle, t T) T

func PluginFn[T any](a []Plugin[T], h Handle, fn Fn[T]) Fn[T] {
	return slice.ReverseReduce(a, func(next Plugin[T], forward Fn[T]) Fn[T] {
		return func(t T) T {
			return next(forward, h, t)
		}
	}, fn)
}

func Default[T any](t T) T {
	return t
}

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

func Digest(next Fn[models.Posts], h Handle, post models.Posts) models.Posts {
	if post.PostExcerpt != "" {
		plugins.PostExcerpt(&post)
	} else {
		plugins.Digest(h.C, &post, config.GetConfig().DigestWordCount)
	}
	return next(post)
}
