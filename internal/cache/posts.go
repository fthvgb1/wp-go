package cache

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github/fthvgb1/wp-go/helper"
	"github/fthvgb1/wp-go/internal/models"
	"github/fthvgb1/wp-go/logs"
	"time"
)

func GetPostById(ctx context.Context, id uint64) (models.Posts, error) {
	return postsCache.GetCache(ctx, id, time.Second, ctx, id)
}

func GetPostsByIds(ctx context.Context, ids []uint64) ([]models.Posts, error) {
	return postsCache.GetCacheBatch(ctx, ids, time.Second, ctx, ids)
}

func SearchPost(ctx context.Context, key string, args ...any) (r []models.Posts, total int, err error) {
	ids, err := searchPostIdsCache.GetCache(ctx, key, time.Second, args...)
	if err != nil {
		return
	}
	total = ids.Length
	r, err = GetPostsByIds(ctx, ids.Ids)
	return
}

func PostLists(ctx context.Context, key string, args ...any) (r []models.Posts, total int, err error) {
	ids, err := postListIdsCache.GetCache(ctx, key, time.Second, args...)
	if err != nil {
		return
	}
	total = ids.Length
	r, err = GetPostsByIds(ctx, ids.Ids)
	return
}

func GetMaxPostId(ctx *gin.Context) (uint64, error) {
	Id, err := maxPostIdCache.GetCache(ctx, time.Second, ctx)
	return Id[0], err
}

func RecentPosts(ctx context.Context, n int) (r []models.Posts) {
	r, err := recentPostsCaches.GetCache(ctx, time.Second, ctx)
	if n < len(r) {
		r = r[:n]
	}
	logs.ErrPrintln(err, "get recent post")
	return
}

func GetContextPost(ctx context.Context, id uint64, date time.Time) (prev, next models.Posts, err error) {
	postCtx, err := postContextCache.GetCache(ctx, id, time.Second, ctx, date)
	if err != nil {
		return models.Posts{}, models.Posts{}, err
	}
	prev = postCtx.Prev
	next = postCtx.Next
	return
}

func GetMonthPostIds(ctx context.Context, year, month string, page, limit int, order string) (r []models.Posts, total int, err error) {
	res, err := monthPostsCache.GetCache(ctx, fmt.Sprintf("%s%s", year, month), time.Second, ctx, year, month)
	if err != nil {
		return
	}
	if order == "desc" {
		res = helper.SliceReverse(res)
	}
	total = len(res)
	rr := helper.SlicePagination(res, page, limit)
	r, err = GetPostsByIds(ctx, rr)
	return
}
