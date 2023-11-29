package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/fthvgb1/wp-go/app/pkg/models"
	"github.com/fthvgb1/wp-go/helper"
	"github.com/fthvgb1/wp-go/helper/number"
	"github.com/fthvgb1/wp-go/helper/slice"
	"github.com/fthvgb1/wp-go/model"
	"time"
)

// RecentComments
// param context.Context
func RecentComments(ctx context.Context, a ...any) (r []models.Comments, err error) {
	n := helper.ParseArgs(10, a...)
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
	r, err := model.ChunkFind[models.Comments](ctx, 300, model.Conditions(
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
	if len(ids) < 1 {
		return nil, nil
	}
	m := make(map[uint64]models.Comments)
	off := 0
	for {
		id := slice.Slice(ids, off, 500)
		if len(id) < 1 {
			break
		}
		r, err := model.Finds[models.Comments](ctx, model.Conditions(
			model.Where(model.SqlBuilder{
				{"comment_ID", "in", ""}, {"comment_approved", "1"},
			}),
			model.Fields("*"),
			model.In(slice.ToAnySlice(id)),
		))
		if err != nil {
			return m, err
		}
		for _, comments := range r {
			m[comments.CommentId] = comments
		}
		off += 500
	}

	return m, nil
}

func GetIncreaseComment(ctx context.Context, currentData []uint64, k uint64, t time.Time, _ ...any) (data []uint64, save bool, refresh bool, err error) {
	r, err := model.ChunkFind[models.Comments](ctx, 1000, model.Conditions(
		model.Where(model.SqlBuilder{
			{"comment_approved", "1"},
			{"comment_post_ID", "=", number.IntToString(k), "int"},
			{"comment_date", ">=", t.Format(time.DateTime)},
		}),
		model.Fields("comment_ID"),
		model.Order(model.SqlBuilder{
			{"comment_date_gmt", "asc"},
			{"comment_ID", "asc"},
		})),
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = nil
			refresh = true
		}
		return
	}
	if len(r) < 1 {
		refresh = true
		return
	}
	rr := slice.Map(r, func(t models.Comments) uint64 {
		return t.CommentId
	})
	data = append(currentData, rr...)
	save = true
	return
}
