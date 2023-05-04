package plugins

import (
	"context"
	"fmt"
	"github.com/fthvgb1/wp-go/app/cmd/cachemanager"
	"github.com/fthvgb1/wp-go/app/pkg/config"
	"github.com/fthvgb1/wp-go/app/pkg/models"
	"github.com/fthvgb1/wp-go/cache"
	"github.com/fthvgb1/wp-go/helper"
	"github.com/fthvgb1/wp-go/plugin/digest"
	"regexp"
	"strings"
	"time"
)

var digestCache *cache.MapCache[uint64, string]

var more = regexp.MustCompile("<!--more(.*?)?-->")

var removeWpBlock = regexp.MustCompile("<!-- /?wp:.*-->")

func InitDigestCache() {
	digestCache = cachemanager.MapCacheBy[uint64](digestRaw, config.GetConfig().CacheTime.DigestCacheTime)
}

func RemoveWpBlock(s string) string {
	return removeWpBlock.ReplaceAllString(s, "")
}

func digestRaw(arg ...any) (string, error) {
	ctx := arg[0].(context.Context)
	s := arg[1].(string)
	id := arg[2].(uint64)
	limit := arg[3].(int)
	if limit < 0 {
		return s, nil
	} else if limit == 0 {
		return "", nil
	}

	s = more.ReplaceAllString(s, "")
	fn := helper.GetContextVal(ctx, "postMoreFn", PostsMore)
	return Digests(s, id, limit, fn), nil
}

func Digests(content string, id uint64, limit int, fn func(id uint64, content, closeTag string) string) string {
	closeTag := ""
	content = RemoveWpBlock(content)
	tag := config.GetConfig().DigestAllowTag
	if tag == "" {
		tag = "<a><b><blockquote><br><cite><code><dd><del><div><dl><dt><em><h1><h2><h3><h4><h5><h6><i><img><li><ol><p><pre><span><strong><ul>"
	}
	content = digest.StripTags(content, tag)
	content, closeTag = digest.Html(content, limit)
	if fn == nil {
		return PostsMore(id, content, closeTag)
	}
	return fn(id, content, closeTag)
}

func PostsMore(id uint64, content, closeTag string) string {
	tmp := `%s......%s<p class="read-more"><a href="/p/%d">继续阅读</a></p>`
	if strings.Contains(closeTag, "pre") || strings.Contains(closeTag, "code") {
		tmp = `%s%s......<p class="read-more"><a href="/p/%d">继续阅读</a></p>`
	}
	content = fmt.Sprintf(tmp, content, closeTag, id)
	return content
}

func Digest(ctx context.Context, post *models.Posts, limit int) {
	content, _ := digestCache.GetCache(ctx, post.Id, time.Second, ctx, post.PostContent, post.Id, limit)
	post.PostContent = content
}

func PostExcerpt(post *models.Posts) {
	post.PostContent = strings.Replace(post.PostExcerpt, "\n", "<br/>", -1)
}