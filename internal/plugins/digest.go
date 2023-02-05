package plugins

import (
	"context"
	"fmt"
	"github.com/fthvgb1/wp-go/cache"
	"github.com/fthvgb1/wp-go/internal/pkg/config"
	"github.com/fthvgb1/wp-go/internal/pkg/models"
	"github.com/fthvgb1/wp-go/plugin/digest"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

var digestCache *cache.MapCache[uint64, string]
var ctx context.Context

func InitDigestCache() {
	ctx = context.Background()
	digestCache = cache.NewMemoryMapCacheByFn[uint64](digestRaw, config.GetConfig().CacheTime.DigestCacheTime)
}

func ClearDigestCache() {
	digestCache.ClearExpired(ctx)
}
func FlushCache() {
	digestCache.Flush(ctx)
}

func digestRaw(arg ...any) (string, error) {
	str := arg[0].(string)
	id := arg[1].(uint64)
	limit := config.GetConfig().DigestWordCount
	if limit < 0 {
		return str, nil
	} else if limit == 0 {
		return "", nil
	}
	return digest.Raw(str, limit, fmt.Sprintf("/p/%d", id)), nil
}

func Digest(p *Plugin[models.Posts], c *gin.Context, post *models.Posts, scene int) {
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
