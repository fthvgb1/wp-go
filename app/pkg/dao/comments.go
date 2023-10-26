package dao

import (
	"context"
	"github.com/fthvgb1/wp-go/app/pkg/models"
	"github.com/fthvgb1/wp-go/helper/number"
	"github.com/fthvgb1/wp-go/helper/slice"
	"github.com/fthvgb1/wp-go/model"
)

// RecentComments
// param context.Context
func RecentComments(a ...any) (r []models.Comments, err error) {
	ctx := a[0].(context.Context)
	n := a[1].(int)
	return model.Finds[models.Comments](ctx, model.Conditions(
		model.Where(model.SqlBuilder{
			{"comment_approved", "1"},
			{"post_status", "publish"},
		}),
		model.Fields("comment_ID,comment_author,comment_post_ID,post_title"),
		model.Order(model.SqlBuilder{{"comment_date_gmt", "desc"}}),
		model.Join(model.SqlBuilder{{"a", "left join", "wp_posts b", "a.comment_post_ID=b.ID"}}),
		model.Limit(n),
	))
}

// PostComments
// param1 context.Context
// param2 postId
func PostComments(ctx context.Context, postId uint64, _ ...any) ([]uint64, error) {
	r, err := model.Finds[models.Comments](ctx, model.Conditions(
		model.Where(model.SqlBuilder{
			{"comment_approved", "1"},
			{"comment_post_ID", "=", number.IntToString(postId), "int"},
		}),
		model.Fields("comment_ID"),
		model.Order(model.SqlBuilder{
			{"comment_date_gmt", "asc"},
			{"comment_ID", "asc"},
		})),
	)
	if err != nil {
		return nil, err
	}
	return slice.Map(r, func(t models.Comments) uint64 {
		return t.CommentId
	}), err
}

func GetCommentByIds(ctx context.Context, ids []uint64, _ ...any) (map[uint64]models.Comments, error) {
	m := make(map[uint64]models.Comments)
	r, err := model.SimpleFind[models.Comments](ctx, model.SqlBuilder{
		{"comment_ID", "in", ""}, {"comment_approved", "1"},
	}, "*", slice.ToAnySlice(ids))
	if err != nil {
		return m, err
	}
	return slice.SimpleToMap(r, func(t models.Comments) uint64 {
		return t.CommentId
	}), err
}
