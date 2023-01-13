package plugins

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github/fthvgb1/wp-go/cache"
	"github/fthvgb1/wp-go/internal/pkg/config"
	"github/fthvgb1/wp-go/internal/pkg/models"
	"github/fthvgb1/wp-go/plugin/digest"
	"strings"
	"time"
)

var digestCache *cache.MapCache[uint64, string]

func InitDigestCache() {
	digestCache = cache.NewMapCacheByFn[uint64](digestRaw, config.Conf.Load().DigestCacheTime)
}

func ClearDigestCache() {
	digestCache.ClearExpired()
}
func FlushCache() {
	digestCache.Flush()
}

func digestRaw(arg ...any) (string, error) {
	str := arg[0].(string)
	id := arg[1].(uint64)
	limit := config.Conf.Load().DigestWordCount
	if limit < 0 {
		return str, nil
	} else if limit == 0 {
		return "", nil
	}
	return digest.Raw(str, limit, fmt.Sprintf("/p/%d", id)), nil
}

func Digest(p *Plugin[models.Posts], c *gin.Context, post *models.Posts, scene uint) {
	if scene == Detail {
		return
	}
	if post.PostExcerpt != "" {
		post.PostContent = strings.Replace(post.PostExcerpt, "\n", "<br/>", -1)
	} else {
		post.PostContent = DigestCache(c, post.Id, post.PostContent)

	}
	p.Next()
}

func DigestCache(ctx *gin.Context, id uint64, str string) string {
	content, _ := digestCache.GetCache(ctx, id, time.Second, str, id)
	return content
}
