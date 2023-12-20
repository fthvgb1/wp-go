package cache

import (
	"context"
	"github.com/fthvgb1/wp-go/app/pkg/dao"
	"github.com/fthvgb1/wp-go/app/pkg/logs"
	"github.com/fthvgb1/wp-go/app/pkg/models"
	"github.com/fthvgb1/wp-go/app/wpconfig"
	"github.com/fthvgb1/wp-go/cache"
	"github.com/fthvgb1/wp-go/cache/cachemanager"
	"github.com/fthvgb1/wp-go/helper"
	"github.com/fthvgb1/wp-go/helper/number"
	"github.com/fthvgb1/wp-go/helper/slice"
	str "github.com/fthvgb1/wp-go/helper/strings"
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

func PostComments(ctx context.Context, Id uint64) ([]models.PostComments, error) {
	ids, err := cachemanager.Get[[]uint64]("PostCommentsIds", ctx, Id, time.Second)
	if err != nil {
		return nil, err
	}
	return GetCommentDataByIds(ctx, ids)
}

func GetCommentById(ctx context.Context, id uint64) (models.PostComments, error) {
	return cachemanager.Get[models.PostComments]("postCommentData", ctx, id, time.Second)
}

func GetCommentDataByIds(ctx context.Context, ids []uint64) ([]models.PostComments, error) {
	return cachemanager.GetMultiple[models.PostComments]("postCommentData", ctx, ids, time.Second)
}

func NewCommentCache() *cache.MapCache[string, string] {
	r, _ := cachemanager.GetMapCache[string, string]("NewComment")
	return r
}

func CommentDataIncreaseUpdates(_ context.Context, _ uint64, _ ...any) ([]models.Comments, error) {
	return nil, nil
}
func IncreaseUpdates(ctx context.Context, currentData []models.Comments, postId uint64, t time.Time, _ ...any) ([]models.Comments, bool, bool, error) {
	var maxId uint64
	if len(currentData) > 0 {
		maxId = currentData[len(currentData)-1].CommentId
	} else {
		maxId, err := dao.LatestCommentId(ctx, postId)
		return []models.Comments{{CommentId: maxId}}, true, false, err
	}
	v, err := dao.IncreaseCommentData(ctx, postId, maxId, t)
	if err != nil {
		return nil, false, false, err
	}
	if len(v) < 1 {
		return nil, false, true, nil
	}
	m, err := dao.CommentDates(ctx, v)
	if err != nil {
		return nil, false, false, err
	}
	CommentData, _ := cachemanager.GetMapCache[uint64, models.PostComments]("postCommentData")
	data := slice.Map(v, func(t uint64) models.Comments {
		comments := m[t].Comments
		if comments.CommentParent > 0 {
			vv, ok := CommentData.Get(ctx, comments.CommentParent)
			if ok && !slice.IsContained(vv.Children, comments.CommentId) {
				vv.Children = append(vv.Children, comments.CommentId)
				CommentData.Set(ctx, comments.CommentParent, vv)
			}
		}
		CommentData.Set(ctx, comments.CommentId, models.PostComments{Comments: comments})
		return comments
	})
	return data, true, false, nil
}

func CommentDataIncreaseUpdate(ctx context.Context, currentData helper.PaginationData[uint64], postId string, _ time.Time, _ ...any) (data helper.PaginationData[uint64], save bool, refresh bool, err error) {
	refresh = true
	increaseUpdateData, _ := cachemanager.GetMapCache[uint64, []models.Comments]("increaseComment30s")
	v, ok := increaseUpdateData.Get(ctx, str.ToInt[uint64](postId))
	if !ok {
		return
	}
	if len(v) < 1 {
		return
	}
	if len(currentData.Data) > 0 {
		if slice.IsContained(currentData.Data, v[0].CommentId) {
			return
		}
	}

	dat := slice.FilterAndMap(v, func(t models.Comments) (uint64, bool) {
		if wpconfig.GetOption("thread_comments") != "1" || "1" == wpconfig.GetOption("thread_comments_depth") {
			return t.CommentId, t.CommentId > 0
		}
		return t.CommentId, t.CommentId > 0 && t.CommentParent == 0
	})
	if len(dat) > 0 {
		save = true
		refresh = false
		var a []uint64
		a = append(currentData.Data, dat...)
		slice.Sorts(a, wpconfig.GetOption("comment_order"))
		data.Data = a
		data.TotalRaw = len(data.Data)
	}
	return data, save, refresh, err
}

func UpdateCommentCache(ctx context.Context, timeout time.Duration, postId uint64) (err error) {
	c, _ := cachemanager.GetPaginationCache[uint64, uint64]("PostCommentsIds")
	if c.IsSwitchDB(postId) {
		return
	}
	_, err = cachemanager.Get[[]models.Comments]("increaseComment30s", ctx, postId, timeout)
	return
}
