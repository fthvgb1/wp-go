package common

import (
	"context"
	"github/fthvgb1/wp-go/helper"
	"github/fthvgb1/wp-go/logs"
	"github/fthvgb1/wp-go/models"
	"strconv"
	"time"
)

func RecentComments(ctx context.Context, n int) (r []models.WpComments) {
	r, err := recentCommentsCaches.GetCache(ctx, time.Second)
	if len(r) > n {
		r = r[0:n]
	}
	logs.ErrPrintln(err, "get recent comment")
	return
}
func recentComments(...any) (r []models.WpComments, err error) {
	return models.Find[models.WpComments](models.SqlBuilder{
		{"comment_approved", "1"},
		{"post_status", "publish"},
	}, "comment_ID,comment_author,comment_post_ID,post_title", "", models.SqlBuilder{{"comment_date_gmt", "desc"}}, models.SqlBuilder{
		{"a", "left join", "wp_posts b", "a.comment_post_ID=b.ID"},
	}, 10)
}

func PostComments(ctx context.Context, Id uint64) ([]models.WpComments, error) {
	return postCommentCaches.GetCache(ctx, Id, time.Second, Id)
}

func postComments(args ...any) ([]models.WpComments, error) {
	postId := args[0].(uint64)
	return models.Find[models.WpComments](models.SqlBuilder{
		{"comment_approved", "1"},
		{"comment_post_ID", "=", strconv.FormatUint(postId, 10), "int"},
	}, "*", "", models.SqlBuilder{
		{"comment_date_gmt", "asc"},
		{"comment_ID", "asc"},
	}, nil, 0)
}

func GetCommentById(ctx context.Context, id uint64) (models.WpComments, error) {
	return commentsCache.GetCache(ctx, id, time.Second, id)
}

func GetCommentByIds(ctx context.Context, ids []uint64) ([]models.WpComments, error) {
	return commentsCache.GetCacheBatch(ctx, ids, time.Second, ids)
}

func getCommentByIds(args ...any) (map[uint64]models.WpComments, error) {
	ids := args[0].([]uint64)
	m := make(map[uint64]models.WpComments)
	r, err := models.Find[models.WpComments](models.SqlBuilder{
		{"comment_ID", "in", ""}, {"comment_approved", "1"},
	}, "*", "", nil, nil, 0, helper.SliceMap(ids, helper.ToAny[uint64]))
	if err != nil {
		return m, err
	}
	return helper.SimpleSliceToMap(r, func(t models.WpComments) uint64 {
		return t.CommentId
	}), err
}
