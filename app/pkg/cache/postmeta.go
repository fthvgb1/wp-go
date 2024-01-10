package cache

import (
	"context"
	"github.com/fthvgb1/wp-go/cache/cachemanager"
	"time"
)

func GetPostMetaByPostIds(ctx context.Context, ids []uint64) ([]map[string]any, error) {
	return cachemanager.GetBatchBy[map[string]any]("postMetaData", ctx, ids, time.Second)
}
func GetPostMetaByPostId(ctx context.Context, id uint64) (map[string]any, error) {
	return cachemanager.GetBy[map[string]any]("postMetaData", ctx, id, time.Second)
}
