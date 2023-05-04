package cache

import (
	"context"
	"time"
)

func GetPostMetaByPostIds(ctx context.Context, ids []uint64) (r []map[string]any, err error) {
	r, err = postMetaCache.GetCacheBatch(ctx, ids, time.Second, ctx, ids)
	return
}
func GetPostMetaByPostId(ctx context.Context, id uint64) (r map[string]any, err error) {
	r, err = postMetaCache.GetCache(ctx, id, time.Second, ctx, id)
	return
}
