package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/fthvgb1/wp-go/app/pkg/models"
	"github.com/fthvgb1/wp-go/app/wpconfig"
	"github.com/fthvgb1/wp-go/helper"
	"github.com/fthvgb1/wp-go/helper/number"
	"github.com/fthvgb1/wp-go/helper/slice"
	str "github.com/fthvgb1/wp-go/helper/strings"
	"github.com/fthvgb1/wp-go/model"
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

func CommentNum(ctx context.Context, postId uint64, _ ...any) (int, error) {
	n, err := model.GetField[models.Posts](ctx, "comment_count", model.Conditions(
		model.Where(model.SqlBuilder{{"ID", "=", number.IntToString(postId), "int"}})))
	if err != nil {
		return 0, err
	}
	return str.ToInteger(n, 0), err
}

func PostTopCommentNum(ctx context.Context, postId uint64, _ ...any) (int, error) {
	v, err := model.GetField[models.Comments](ctx, "count(*) num", model.Conditions(
		model.Where(postTopCommentNumWhere(postId)),
	))
	if err != nil {
		return 0, err
	}
	return str.ToInteger(v, 0), nil
}

func postTopCommentNumWhere(postId uint64) model.SqlBuilder {
	threadComments := wpconfig.GetOption("thread_comments")
	pageComments := wpconfig.GetOption("page_comments")
	where := model.SqlBuilder{
		{"comment_approved", "1"},
		{"comment_post_ID", "=", number.IntToString(postId), "int"},
	}
	if pageComments != "1" || threadComments == "1" || "1" == wpconfig.GetOption("thread_comments_depth") {
		where = append(where, []string{"comment_parent", "0"})
	}
	return where
}

func PostCommentsIds(ctx context.Context, postId uint64, page, limit, totalRaw int, _ ...any) ([]uint64, int, error) {
	order := wpconfig.GetOption("comment_order")
	pageComments := wpconfig.GetOption("page_comments")
	condition := model.Conditions(
		model.Where(postTopCommentNumWhere(postId)),
		model.TotalRaw(totalRaw),
		model.Fields("comment_ID"),
		model.Order(model.SqlBuilder{
			{"comment_date_gmt", order},
			{"comment_ID", "asc"},
		}),
	)
	var r []models.Comments
	var total int
	var err error
	if pageComments != "1" {
		r, err = model.ChunkFind[models.Comments](ctx, 300, condition)
		total = len(r)
	} else {
		r, total, err = model.Pagination[models.Comments](ctx, condition, page, limit)
	}

	if err != nil && errors.Is(err, sql.ErrNoRows) {
		err = nil
	}
	return slice.Map(r, func(t models.Comments) uint64 {
		return t.CommentId
	}), total, err
}

func CommentChildren(ctx context.Context, commentIds []uint64, _ ...any) (r map[uint64][]uint64, err error) {
	rr, err := model.Finds[models.Comments](ctx, model.Conditions(
		model.Where(model.SqlBuilder{
			{"comment_parent", "in", ""},
			{"comment_approved", "1"},
		}),
		model.In(slice.ToAnySlice(commentIds)),
		model.Fields("comment_ID,comment_parent"),
	))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = nil
		}
		return
	}
	rrr := slice.GroupBy(rr, func(v models.Comments) (uint64, uint64) {
		return v.CommentParent, v.CommentId
	})
	r = make(map[uint64][]uint64)
	for _, id := range commentIds {
		r[id] = rrr[id]
	}
	return
}

func PreviousCommentNum(ctx context.Context, commentId, postId uint64) (int, error) {
	v, err := model.GetField[models.Comments](ctx, "count(*)", model.Conditions(
		model.Where(model.SqlBuilder{
			{"comment_approved", "1"},
			{"comment_post_ID", "=", number.IntToString(postId), "int"},
			{"comment_ID", "<", number.IntToString(commentId), "int"},
			{"comment_parent", "=", "0", "int"},
		}),
	))
	if err != nil {
		return 0, err
	}
	return str.ToInteger(v, 0), nil
}
