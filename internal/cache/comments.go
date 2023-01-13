package cache

import (
	"context"
	"github/fthvgb1/wp-go/internal/logs"
	"github/fthvgb1/wp-go/internal/models"
	"time"
)

func RecentComments(ctx context.Context, n int) (r []models.Comments) {
	r, err := recentCommentsCaches.GetCache(ctx, time.Second, ctx)
	if len(r) > n {
		r = r[0:n]
	}
	logs.ErrPrintln(err, "get recent comment")
	return
}

func PostComments(ctx context.Context, Id uint64) ([]models.Comments, error) {
	ids, err := postCommentCaches.GetCache(ctx, Id, time.Second, ctx, Id)
	if err != nil {
		return nil, err
	}
	return GetCommentByIds(ctx, ids)
}

func GetCommentById(ctx context.Context, id uint64) (models.Comments, error) {
	return commentsCache.GetCache(ctx, id, time.Second, ctx, id)
}

func GetCommentByIds(ctx context.Context, ids []uint64) ([]models.Comments, error) {
	return commentsCache.GetCacheBatch(ctx, ids, time.Second, ctx, ids)
}