package common

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github/fthvgb1/wp-go/helper"
	"github/fthvgb1/wp-go/internal/wp"
	"github/fthvgb1/wp-go/logs"
	"github/fthvgb1/wp-go/models"
	"strings"
	"sync/atomic"
	"time"
)

func GetPostById(ctx context.Context, id uint64) (wp.Posts, error) {
	return postsCache.GetCache(ctx, id, time.Second, ctx, id)
}

func GetPostsByIds(ctx context.Context, ids []uint64) ([]wp.Posts, error) {
	return postsCache.GetCacheBatch(ctx, ids, time.Second, ctx, ids)
}

func SearchPost(ctx context.Context, key string, args ...any) (r []wp.Posts, total int, err error) {
	ids, err := searchPostIdsCache.GetCache(ctx, key, time.Second, args...)
	if err != nil {
		return
	}
	total = ids.Length
	r, err = GetPostsByIds(ctx, ids.Ids)
	return
}

func getPostsByIds(ids ...any) (m map[uint64]wp.Posts, err error) {
	ctx := ids[0].(context.Context)
	m = make(map[uint64]wp.Posts)
	id := ids[1].([]uint64)
	arg := helper.SliceMap(id, helper.ToAny[uint64])
	rawPosts, err := models.Find[wp.Posts](ctx, models.SqlBuilder{{
		"Id", "in", "",
	}}, "a.*,ifnull(d.name,'') category_name,ifnull(taxonomy,'') `taxonomy`", "", nil, models.SqlBuilder{{
		"a", "left join", "wp_term_relationships b", "a.Id=b.object_id",
	}, {
		"left join", "wp_term_taxonomy c", "b.term_taxonomy_id=c.term_taxonomy_id",
	}, {
		"left join", "wp_terms d", "c.term_id=d.term_id",
	}}, nil, 0, arg)
	if err != nil {
		return m, err
	}
	postsMap := make(map[uint64]wp.Posts)
	for i, post := range rawPosts {
		v, ok := postsMap[post.Id]
		if !ok {
			v = rawPosts[i]
		}
		if post.Taxonomy == "category" {
			v.Categories = append(v.Categories, post.CategoryName)
		} else if post.Taxonomy == "post_tag" {
			v.Tags = append(v.Tags, post.CategoryName)
		}
		postsMap[post.Id] = v
	}
	meta, _ := getPostMetaByPostIds(ctx, id)
	for k, pp := range postsMap {
		if len(pp.Categories) > 0 {
			t := make([]string, 0, len(pp.Categories))
			for _, cat := range pp.Categories {
				t = append(t, fmt.Sprintf(`<a href="/p/category/%s" rel="category tag">%s</a>`, cat, cat))
			}
			pp.CategoriesHtml = strings.Join(t, "、")
			_, ok := meta[pp.Id]
			if ok {
				thumb := ToPostThumbnail(ctx, pp.Id)
				if thumb.Path != "" {
					pp.Thumbnail = thumb
				}
			}
		}
		if len(pp.Tags) > 0 {
			t := make([]string, 0, len(pp.Tags))
			for _, cat := range pp.Tags {
				t = append(t, fmt.Sprintf(`<a href="/p/tag/%s" rel="tag">%s</a>`, cat, cat))
			}
			pp.TagsHtml = strings.Join(t, "、")
		}
		m[k] = pp
	}
	return
}

func PostLists(ctx context.Context, key string, args ...any) (r []wp.Posts, total int, err error) {
	ids, err := postListIdsCache.GetCache(ctx, key, time.Second, args...)
	if err != nil {
		return
	}
	total = ids.Length
	r, err = GetPostsByIds(ctx, ids.Ids)
	return
}

func searchPostIds(args ...any) (ids PostIds, err error) {
	ctx := args[0].(context.Context)
	where := args[1].(models.SqlBuilder)
	page := args[2].(int)
	limit := args[3].(int)
	order := args[4].(models.SqlBuilder)
	join := args[5].(models.SqlBuilder)
	postType := args[6].([]any)
	postStatus := args[7].([]any)
	res, total, err := models.SimplePagination[wp.Posts](ctx, where, "ID", "", page, limit, order, join, nil, postType, postStatus)
	for _, posts := range res {
		ids.Ids = append(ids.Ids, posts.Id)
	}
	ids.Length = total
	totalR := int(atomic.LoadInt64(&TotalRaw))
	if total > totalR {
		tt := int64(total)
		atomic.StoreInt64(&TotalRaw, tt)
	}
	return
}

func getMaxPostId(a ...any) ([]uint64, error) {
	ctx := a[0].(context.Context)
	r, err := models.SimpleFind[wp.Posts](ctx, models.SqlBuilder{{"post_type", "post"}, {"post_status", "publish"}}, "max(ID) ID")
	var id uint64
	if len(r) > 0 {
		id = r[0].Id
	}
	return []uint64{id}, err
}

func GetMaxPostId(ctx *gin.Context) (uint64, error) {
	Id, err := maxPostIdCache.GetCache(ctx, time.Second, ctx)
	return Id[0], err
}

func RecentPosts(ctx context.Context, n int) (r []wp.Posts) {
	r, err := recentPostsCaches.GetCache(ctx, time.Second, ctx)
	if n < len(r) {
		r = r[:n]
	}
	logs.ErrPrintln(err, "get recent post")
	return
}
func recentPosts(a ...any) (r []wp.Posts, err error) {
	ctx := a[0].(context.Context)
	r, err = models.Find[wp.Posts](ctx, models.SqlBuilder{{
		"post_type", "post",
	}, {"post_status", "publish"}}, "ID,post_title,post_password", "", models.SqlBuilder{{"post_date", "desc"}}, nil, nil, 10)
	for i, post := range r {
		if post.PostPassword != "" {
			PasswordProjectTitle(&r[i])
		}
	}
	return
}

func GetContextPost(ctx context.Context, id uint64, date time.Time) (prev, next wp.Posts, err error) {
	postCtx, err := postContextCache.GetCache(ctx, id, time.Second, ctx, date)
	if err != nil {
		return wp.Posts{}, wp.Posts{}, err
	}
	prev = postCtx.prev
	next = postCtx.next
	return
}

func getPostContext(arg ...any) (r PostContext, err error) {
	ctx := arg[0].(context.Context)
	t := arg[1].(time.Time)
	next, err := models.FirstOne[wp.Posts](ctx, models.SqlBuilder{
		{"post_date", ">", t.Format("2006-01-02 15:04:05")},
		{"post_status", "in", ""},
		{"post_type", "post"},
	}, "ID,post_title,post_password", nil, []any{"publish"})
	if err == sql.ErrNoRows {
		err = nil
	}
	if err != nil {
		return
	}
	prev, err := models.FirstOne[wp.Posts](ctx, models.SqlBuilder{
		{"post_date", "<", t.Format("2006-01-02 15:04:05")},
		{"post_status", "in", ""},
		{"post_type", "post"},
	}, "ID,post_title", models.SqlBuilder{{"post_date", "desc"}}, []any{"publish"})
	if err == sql.ErrNoRows {
		err = nil
	}
	if err != nil {
		return
	}
	r = PostContext{
		prev: prev,
		next: next,
	}
	return
}

func GetMonthPostIds(ctx context.Context, year, month string, page, limit int, order string) (r []wp.Posts, total int, err error) {
	res, err := monthPostsCache.GetCache(ctx, fmt.Sprintf("%s%s", year, month), time.Second, ctx, year, month)
	if err != nil {
		return
	}
	if order == "desc" {
		res = helper.SliceReverse(res)
	}
	total = len(res)
	rr := helper.SlicePagination(res, page, limit)
	r, err = GetPostsByIds(ctx, rr)
	return
}

func monthPost(args ...any) (r []uint64, err error) {
	ctx := args[0].(context.Context)
	year, month := args[1].(string), args[2].(string)
	where := models.SqlBuilder{
		{"post_type", "in", ""},
		{"post_status", "in", ""},
		{"year(post_date)", year},
		{"month(post_date)", month},
	}
	postType := []any{"post"}
	status := []any{"publish"}
	ids, err := models.Find[wp.Posts](ctx, where, "ID", "", models.SqlBuilder{{"Id", "asc"}}, nil, nil, 0, postType, status)
	if err != nil {
		return
	}
	for _, post := range ids {
		r = append(r, post.Id)
	}
	return
}
