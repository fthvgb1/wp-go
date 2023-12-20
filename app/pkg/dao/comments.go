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

func GetIncreaseComment(ctx context.Context, currentData []uint64, k uint64, _ time.Time, _ ...any) (data []uint64, save bool, refresh bool, err error) {
	r, err := model.ChunkFind[models.Comments](ctx, 1000, model.Conditions(
		model.Where(model.SqlBuilder{
			{"comment_approved", "1"},
			{"comment_post_ID", "=", number.IntToString(k), "int"},
			//{"comment_date", ">=", t.Format(time.DateTime)},
			{"comment_ID", ">", number.IntToString(currentData[len(currentData)-1])},
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

func CommentDates(ctx context.Context, CommentIds []uint64, _ ...any) (map[uint64]models.PostComments, error) {
	if len(CommentIds) < 1 {
		return nil, nil
	}
	m := make(map[uint64]models.PostComments)
	off := 0
	threadComments := wpconfig.GetOption("thread_comments")
	where := model.SqlBuilder{
		{"comment_approved", "1"},
	}
	var in [][]any
	if threadComments == "1" {
		where = append(where, []string{"and", "comment_ID", "in", "", "", "or", "comment_parent", "in", "", ""})
	} else {
		where = append(where, []string{"comment_ID", "in", ""})
	}
	for {
		id := slice.Slice(CommentIds, off, 200)
		if len(id) < 1 {
			break
		}
		if threadComments == "1" {
			in = [][]any{slice.ToAnySlice(id), slice.ToAnySlice(id)}
		} else {
			in = [][]any{slice.ToAnySlice(id)}
		}
		r, err := model.Finds[models.Comments](ctx, model.Conditions(
			model.Where(where),
			model.Fields("*"),
			model.In(in...),
		))
		if err != nil {
			return m, err
		}
		rr := slice.GroupBy(r, func(t models.Comments) (uint64, models.Comments) {
			return t.CommentParent, t
		})
		mm := map[uint64][]uint64{}
		for u, comments := range rr {
			slice.SimpleSort(comments, slice.ASC, func(t models.Comments) uint64 {
				return t.CommentId
			})
			mm[u] = slice.Map(comments, func(t models.Comments) uint64 {
				return t.CommentId
			})
		}
		for _, comments := range r {
			var children []uint64
			if threadComments == "1" {
				children = mm[comments.CommentId]
			}
			v := models.PostComments{
				Comments: comments,
				Children: children,
			}
			m[comments.CommentId] = v
		}
		off += 200
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
func PostCommentLocal(_ context.Context, data []uint64, _ uint64, page, limit int, _ ...any) ([]uint64, int, error) {
	/*order := wpconfig.GetOption("comment_order")
	if order == "desc" {
		r := slice.ReversePagination(data, page, limit)
		if len(r) < 1 {
			return nil, 0, nil
		}
		r = slice.SortsNew(r, slice.DESC)
		return r, len(data), nil
	}*/
	return slice.Pagination(data, page, limit), len(data), nil
}

func PostCommentsIds(ctx context.Context, postId uint64, page, limit, totalRaw int, _ ...any) ([]uint64, int, error) {
	order := wpconfig.GetOption("comment_order")
	threadComments := wpconfig.GetOption("thread_comments")
	where := model.SqlBuilder{
		{"comment_approved", "1"},
		{"comment_post_ID", "=", number.IntToString(postId), "int"},
	}
	if threadComments == "1" || "1" == wpconfig.GetOption("thread_comments_depth") {
		where = append(where, []string{"comment_parent", "0"})
	}
	r, total, err := model.Pagination[models.Comments](ctx, model.Conditions(
		model.Where(where),
		model.TotalRaw(totalRaw),
		model.Fields("comment_ID"),
		model.Order(model.SqlBuilder{
			{"comment_date_gmt", order},
			{"comment_ID", "asc"},
		}),
	), page, limit)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		err = nil
	}
	return slice.Map(r, func(t models.Comments) uint64 {
		return t.CommentId
	}), total, err
}

func IncreaseCommentData(ctx context.Context, postId, maxCommentId uint64, _ time.Time) ([]uint64, error) {
	r, err := model.ChunkFind[models.Comments](ctx, 1000, model.Conditions(
		model.Where(model.SqlBuilder{
			{"comment_approved", "1"},
			{"comment_post_ID", "=", number.IntToString(postId), "int"},
			{"comment_ID", ">", number.IntToString(maxCommentId), "int"},
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
			return nil, nil
		}
		return nil, err
	}
	return slice.Map(r, func(t models.Comments) uint64 {
		return t.CommentId
	}), err
}

func LatestCommentId(ctx context.Context, postId uint64) (uint64, error) {
	v, err := model.GetField[models.Comments](ctx, "comment_ID", model.Conditions(
		model.Where(model.SqlBuilder{
			{"comment_approved", "1"},
			{"comment_post_ID", "=", number.IntToString(postId), "int"},
		}),
		model.Order(model.SqlBuilder{{"comment_ID", "desc"}}),
		model.Limit(1),
	))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = nil
		}
		return 0, err
	}
	return str.ToInteger[uint64](v, 0), err
}
