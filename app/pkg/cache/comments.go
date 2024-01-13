package cache

import (
	"context"
	"fmt"
	"github.com/fthvgb1/wp-go/app/pkg/dao"
	"github.com/fthvgb1/wp-go/app/pkg/logs"
	"github.com/fthvgb1/wp-go/app/pkg/models"
	"github.com/fthvgb1/wp-go/app/wpconfig"
	"github.com/fthvgb1/wp-go/cache"
	"github.com/fthvgb1/wp-go/cache/cachemanager"
	"github.com/fthvgb1/wp-go/helper"
	"github.com/fthvgb1/wp-go/helper/number"
	str "github.com/fthvgb1/wp-go/helper/strings"
	"time"
)

// RecentComments query func see RecentComment
func RecentComments(ctx context.Context, n int) (r []models.Comments) {
	nn := number.Max(n, 10)
	r, err := cachemanager.GetVarVal[[]models.Comments]("recentComments", ctx, time.Second, ctx, nn)
	if len(r) > n {
		r = r[0:n]
	}
	logs.IfError(err, "get recent comment fail")
	return
}

// PostTopLevelCommentIds query func see PostTopComments
func PostTopLevelCommentIds(ctx context.Context, postId uint64, page, limit, total int, order string, a ...any) ([]uint64, error) {
	var key string
	if len(a) > 0 {
		key = helper.ParseArgs("", a...)
	}
	if key == "" {
		key = fmt.Sprintf("%d-%d-%d-%d-%s", postId, page, limit, total, order)
	}
	return cachemanager.GetBy[[]uint64]("PostCommentsIds", ctx,
		key, time.Second, postId, page, limit, 0, order)
}

// GetCommentById query func see dao.GetCommentByIds
func GetCommentById(ctx context.Context, id uint64) (models.Comments, error) {
	return cachemanager.GetBy[models.Comments]("postCommentData", ctx, id, time.Second)
}

// GetCommentDataByIds query func see dao.GetCommentByIds
func GetCommentDataByIds(ctx context.Context, ids []uint64) ([]models.Comments, error) {
	return cachemanager.GetBatchBy[models.Comments]("postCommentData", ctx, ids, time.Second)
}

func NewCommentCache() *cache.MapCache[string, string] {
	r, _ := cachemanager.GetMapCache[string, string]("NewComment")
	return r
}

func PostTopComments(ctx context.Context, _ string, a ...any) ([]uint64, error) {
	postId := a[0].(uint64)
	page := a[1].(int)
	limit := a[2].(int)
	total := a[3].(int)
	order := helper.ParseArgs("", a...)
	if order == "" {
		order = wpconfig.GetOption("comment_order")
	}
	v, _, err := dao.PostCommentsIds(ctx, postId, page, limit, total, order)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func RecentComment(ctx context.Context, a ...any) (r []models.Comments, err error) {
	r, err = dao.RecentComments(ctx, a...)
	if err != nil {
		return r, err
	}
	for i, comment := range r {
		r[i].CommentAuthorUrl, err = GetCommentUrl(ctx, comment.CommentId, comment.CommentPostId)
		if err != nil {
			return nil, err
		}
	}
	return
}

func GetCommentUrl(ctx context.Context, commentId, postId uint64) (string, error) {
	if wpconfig.GetOption("page_comments") != "1" {
		return fmt.Sprintf("/p/%d#comment-%d", postId, commentId), nil
	}
	commentsPerPage := str.ToInteger(wpconfig.GetOption("comments_per_page"), 5)
	topCommentId, err := AncestorCommentId(ctx, commentId)
	if err != nil {
		return "", err
	}
	totalNum, err := cachemanager.GetBy[int]("postTopCommentsNum", ctx, postId, time.Second)
	if err != nil {
		return "", err
	}
	if totalNum <= commentsPerPage {
		return fmt.Sprintf("/p/%d#comment-%d", postId, commentId), nil
	}
	num, err := dao.PreviousCommentNum(ctx, topCommentId, postId)
	if err != nil {
		return "", err
	}
	order := wpconfig.GetOption("comment_order")
	page := number.DivideCeil(num+1, commentsPerPage)
	if order == "desc" {
		page = number.DivideCeil(totalNum-num, commentsPerPage)
	}
	return fmt.Sprintf("/p/%d/comment-page-%d/#comment-%d", postId, page, commentId), nil
}

func AncestorCommentId(ctx context.Context, commentId uint64) (uint64, error) {
	comment, err := GetCommentById(ctx, commentId)
	if err != nil {
		return 0, err
	}
	if comment.CommentParent == 0 {
		return comment.CommentId, nil
	}
	return AncestorCommentId(ctx, comment.CommentParent)
}
