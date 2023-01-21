package dao

import (
	"context"
	"github.com/fthvgb1/wp-go/helper"
	"github.com/fthvgb1/wp-go/helper/slice"
	"github.com/fthvgb1/wp-go/internal/pkg/models"
	"github.com/fthvgb1/wp-go/model"
	"strconv"
)

// RecentComments
// param context.Context
func RecentComments(a ...any) (r []models.Comments, err error) {
	ctx := a[0].(context.Context)
	return model.Find[models.Comments](ctx, model.SqlBuilder{
		{"comment_approved", "1"},
		{"post_status", "publish"},
	}, "comment_ID,comment_author,comment_post_ID,post_title", "", model.SqlBuilder{{"comment_date_gmt", "desc"}}, model.SqlBuilder{
		{"a", "left join", "wp_posts b", "a.comment_post_ID=b.ID"},
	}, nil, 10)
}

// PostComments
// param1 context.Context
// param2 postId
func PostComments(args ...any) ([]uint64, error) {
	ctx := args[0].(context.Context)
	postId := args[1].(uint64)
	r, err := model.Find[models.Comments](ctx, model.SqlBuilder{
		{"comment_approved", "1"},
		{"comment_post_ID", "=", strconv.FormatUint(postId, 10), "int"},
	}, "comment_ID", "", model.SqlBuilder{
		{"comment_date_gmt", "asc"},
		{"comment_ID", "asc"},
	}, nil, nil, 0)
	if err != nil {
		return nil, err
	}
	return slice.Map(r, func(t models.Comments) uint64 {
		return t.CommentId
	}), err
}

func GetCommentByIds(args ...any) (map[uint64]models.Comments, error) {
	ctx := args[0].(context.Context)
	ids := args[1].([]uint64)
	m := make(map[uint64]models.Comments)
	r, err := model.SimpleFind[models.Comments](ctx, model.SqlBuilder{
		{"comment_ID", "in", ""}, {"comment_approved", "1"},
	}, "*", slice.Map(ids, helper.ToAny[uint64]))
	if err != nil {
		return m, err
	}
	return slice.SimpleToMap(r, func(t models.Comments) uint64 {
		return t.CommentId
	}), err
}
