package wp

import (
	"github.com/fthvgb1/wp-go/app/pkg/config"
	"github.com/fthvgb1/wp-go/app/pkg/models"
	"github.com/fthvgb1/wp-go/app/plugins"
	"github.com/fthvgb1/wp-go/app/plugins/wpposts"
	"github.com/fthvgb1/wp-go/helper/maps"
	"github.com/fthvgb1/wp-go/helper/slice"
)

type PostsPlugin func(*Handle, *models.Posts)

func PostsPlugins(initial PostsPlugin, calls ...func(PostsPlugin, *Handle, *models.Posts)) PostsPlugin {
	return slice.ReverseReduce(calls, func(t func(PostsPlugin, *Handle, *models.Posts), r PostsPlugin) PostsPlugin {
		return func(handle *Handle, posts *models.Posts) {
			t(r, handle, posts)
		}
	}, initial)
}

var pluginFns = map[string]func(PostsPlugin, *Handle, *models.Posts){
	"passwordProject": PasswordProject,
	"digest":          Digest,
}

// PasswordProject 标题和内容密码保护
func PasswordProject(next PostsPlugin, h *Handle, post *models.Posts) {
	r := post
	if post.PostPassword != "" {
		wpposts.PasswordProjectTitle(r)
		if h.GetPassword() != post.PostPassword {
			wpposts.PasswdProjectContent(r)
			return
		}
	}
	next(h, r)
}

// Digest 生成摘要
func Digest(next PostsPlugin, h *Handle, post *models.Posts) {
	if post.PostExcerpt != "" {
		plugins.PostExcerpt(post)
	} else {
		plugins.Digest(h.C, post, config.GetConfig().DigestWordCount)
	}
	next(h, post)
}

func PostPlugin(calls ...func(h *Handle, posts *models.Posts)) PostsPlugin {
	return func(h *Handle, posts *models.Posts) {
		for _, call := range calls {
			call(h, posts)
		}
	}
}

func ListPostPlugins() map[string]func(PostsPlugin, *Handle, *models.Posts) {
	return maps.Copy(pluginFns)
}

func ProjectTitle(t models.Posts) models.Posts {
	if t.PostPassword != "" {
		wpposts.PasswordProjectTitle(&t)
	}
	return t
}

func GetListPostPlugins(name []string, m map[string]func(PostsPlugin, *Handle, *models.Posts)) []func(PostsPlugin, *Handle, *models.Posts) {
	return slice.FilterAndMap(name, func(t string) (func(PostsPlugin, *Handle, *models.Posts), bool) {
		v, ok := m[t]
		if ok {
			return v, true
		}
		return nil, false
	})
}
