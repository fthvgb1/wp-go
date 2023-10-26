package cache

import (
	"context"
	"fmt"
	"github.com/fthvgb1/wp-go/app/pkg/logs"
	"github.com/fthvgb1/wp-go/app/pkg/models"
	"github.com/fthvgb1/wp-go/app/wpconfig"
	"github.com/fthvgb1/wp-go/helper/number"
	"github.com/fthvgb1/wp-go/helper/slice"
	str "github.com/fthvgb1/wp-go/helper/strings"
	"time"
)

func GetPostById(ctx context.Context, id uint64) (models.Posts, error) {
	//return cachemanager.Get[models.Posts]("postData", ctx, id, time.Second)
	return postsCache.GetCache(ctx, id, time.Second)
}

func GetPostsByIds(ctx context.Context, ids []uint64) ([]models.Posts, error) {
	//return cachemanager.GetMultiple[models.Posts]("postData", ctx, ids, time.Second)
	return postsCache.GetCacheBatch(ctx, ids, time.Second)
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

func GetMaxPostId(ctx context.Context) (uint64, error) {
	return maxPostIdCache.GetCache(ctx, time.Second, ctx)
}

func RecentPosts(ctx context.Context, n int) (r []models.Posts) {
	nn := n
	feedNum := str.ToInteger(wpconfig.GetOption("posts_per_rss"), 10)
	nn = number.Max(n, feedNum)
	r, err := recentPostsCaches.GetCache(ctx, time.Second, ctx, nn)
	if n < len(r) {
		r = r[:n]
	}
	logs.IfError(err, "get recent post")
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
		res = slice.Reverse(res)
	}
	total = len(res)
	rr := slice.Pagination(res, page, limit)
	r, err = GetPostsByIds(ctx, rr)
	return
}
