package cache

import (
	"context"
	"github.com/fthvgb1/wp-go/app/pkg/dao"
	"github.com/fthvgb1/wp-go/app/pkg/logs"
	"github.com/fthvgb1/wp-go/app/pkg/models"
	"github.com/fthvgb1/wp-go/cache"
	"github.com/fthvgb1/wp-go/cache/cachemanager"
	"github.com/fthvgb1/wp-go/helper"
	"github.com/fthvgb1/wp-go/helper/number"
	"time"
)

func RecentComments(ctx context.Context, n int) (r []models.Comments) {
	nn := number.Max(n, 10)
	r, err := cachemanager.GetVarVal[[]models.Comments]("recentComments", ctx, time.Second, ctx, nn)
	if len(r) > n {
		r = r[0:n]
	}
	logs.IfError(err, "get recent comment fail")
	return
}

func PostComments(ctx context.Context, Id uint64) ([]models.Comments, error) {
	ids, err := cachemanager.Get[[]uint64]("PostCommentsIds", ctx, Id, time.Second)
	if err != nil {
		return nil, err
	}
	return GetCommentDataByIds(ctx, ids)
}

func GetCommentById(ctx context.Context, id uint64) (models.Comments, error) {
	return cachemanager.Get[models.Comments]("postCommentData", ctx, id, time.Second)
}

func GetCommentDataByIds(ctx context.Context, ids []uint64) ([]models.Comments, error) {
	return cachemanager.GetMultiple[models.Comments]("postCommentData", ctx, ids, time.Second)
}

func NewCommentCache() *cache.MapCache[string, string] {
	r, _ := cachemanager.GetMapCache[string, string]("NewComment")
	return r
}

func PostTopComments(ctx context.Context, _ string, a ...any) (helper.PaginationData[uint64], error) {
	postId := a[0].(uint64)
	page := a[1].(int)
	limit := a[2].(int)
	total := a[3].(int)
	v, total, err := dao.PostCommentsIds(ctx, postId, page, limit, total)
	if err != nil {
		return helper.PaginationData[uint64]{}, err
	}
	return helper.PaginationData[uint64]{
		Data:     v,
		TotalRaw: total,
	}, nil
}
