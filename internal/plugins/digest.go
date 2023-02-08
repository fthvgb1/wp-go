package plugins

import (
	"context"
	"fmt"
	"github.com/fthvgb1/wp-go/cache"
	"github.com/fthvgb1/wp-go/internal/pkg/config"
	"github.com/fthvgb1/wp-go/internal/pkg/models"
	"github.com/fthvgb1/wp-go/plugin/digest"
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
	limit := arg[2].(int)
	if limit < 0 {
		return str, nil
	} else if limit == 0 {
		return "", nil
	}
	return digest.Raw(str, limit, fmt.Sprintf("/p/%d", id)), nil
}

func Digest(ctx context.Context, post *models.Posts, limit int) {
	content, _ := digestCache.GetCache(ctx, post.Id, time.Second, post.PostContent, post.Id, limit)
	post.PostContent = content
}

func PostExcerpt(post *models.Posts) {
	post.PostContent = strings.Replace(post.PostExcerpt, "\n", "<br/>", -1)
}
