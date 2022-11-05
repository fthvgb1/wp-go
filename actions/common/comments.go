package common

import (
	"context"
	"github/fthvgb1/wp-go/helper"
	"github/fthvgb1/wp-go/logs"
	"github/fthvgb1/wp-go/models"
	"github/fthvgb1/wp-go/models/wp"
	"strconv"
	"time"
)

func RecentComments(ctx context.Context, n int) (r []wp.WpComments) {
	r, err := recentCommentsCaches.GetCache(ctx, time.Second)
	if len(r) > n {
		r = r[0:n]
	}
	logs.ErrPrintln(err, "get recent comment")
	return
}
func recentComments(...any) (r []wp.WpComments, err error) {
	return models.Find[wp.WpComments](models.SqlBuilder{
		{"comment_approved", "1"},
		{"post_status", "publish"},
	}, "comment_ID,comment_author,comment_post_ID,post_title", "", models.SqlBuilder{{"comment_date_gmt", "desc"}}, models.SqlBuilder{
		{"a", "left join", "wp_posts b", "a.comment_post_ID=b.ID"},
	}, nil, 10)
}

func PostComments(ctx context.Context, Id uint64) ([]wp.WpComments, error) {
	ids, err := postCommentCaches.GetCache(ctx, Id, time.Second, Id)
	if err != nil {
		return nil, err
	}
	return GetCommentByIds(ctx, ids)
}

func postComments(args ...any) ([]uint64, error) {
	postId := args[0].(uint64)
	r, err := models.Find[wp.WpComments](models.SqlBuilder{
		{"comment_approved", "1"},
		{"comment_post_ID", "=", strconv.FormatUint(postId, 10), "int"},
	}, "comment_ID", "", models.SqlBuilder{
		{"comment_date_gmt", "asc"},
		{"comment_ID", "asc"},
	}, nil, nil, 0)
	if err != nil {
		return nil, err
	}
	return helper.SliceMap(r, func(t wp.WpComments) uint64 {
		return t.CommentId
	}), err
}

func GetCommentById(ctx context.Context, id uint64) (wp.WpComments, error) {
	return commentsCache.GetCache(ctx, id, time.Second, id)
}

func GetCommentByIds(ctx context.Context, ids []uint64) ([]wp.WpComments, error) {
	return commentsCache.GetCacheBatch(ctx, ids, time.Second, ids)
}

func getCommentByIds(args ...any) (map[uint64]wp.WpComments, error) {
	ids := args[0].([]uint64)
	m := make(map[uint64]wp.WpComments)
	r, err := models.SimpleFind[wp.WpComments](models.SqlBuilder{
		{"comment_ID", "in", ""}, {"comment_approved", "1"},
	}, "*", helper.SliceMap(ids, helper.ToAny[uint64]))
	if err != nil {
		return m, err
	}
	return helper.SimpleSliceToMap(r, func(t wp.WpComments) uint64 {
		return t.CommentId
	}), err
}
