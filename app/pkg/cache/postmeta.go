package cache

import (
	"context"
	"github.com/fthvgb1/wp-go/cache/cachemanager"
	"time"
)

// GetPostMetaByPostIds query func see dao.GetPostMetaByPostIds
func GetPostMetaByPostIds(ctx context.Context, ids []uint64) ([]map[string]any, error) {
	return cachemanager.GetBatchBy[map[string]any]("postMetaData", ctx, ids, time.Second)
}

// GetPostMetaByPostId query func see dao.GetPostMetaByPostIds
func GetPostMetaByPostId(ctx context.Context, id uint64) (map[string]any, error) {
	return cachemanager.GetBy[map[string]any]("postMetaData", ctx, id, time.Second)
}
