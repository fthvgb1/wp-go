package common

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github/fthvgb1/wp-go/helper"
	"github/fthvgb1/wp-go/models"
	"strings"
	"time"
)

func GetPostAndCache(ctx context.Context, id uint64) (models.WpPosts, error) {

	return postsCache.GetCache(ctx, id, time.Second, id)
}

func GetPostById(id uint64) models.WpPosts {
	return postsCache.Get(id)
}

func GetPostsByIds(ctx context.Context, ids []uint64) ([]models.WpPosts, error) {
	return postsCache.GetCacheBatch(ctx, ids, time.Second, ids)
}

func SearchPost(ctx context.Context, key string, args ...any) (r []models.WpPosts, total int, err error) {
	ids, err := searchPostIdsCache.GetCache(ctx, key, time.Second, args...)
	if err != nil {
		return
	}
	total = ids.Length
	r, err = GetPostsByIds(ctx, ids.Ids)
	return
}

func getPostById(id ...any) (post models.WpPosts, err error) {
	Id := id[0].(uint64)
	posts, err := getPostsByIds([]uint64{Id})
	if err != nil {
		return models.WpPosts{}, err
	}
	post = posts[Id]
	return
}

func getPostsByIds(ids ...any) (m map[uint64]models.WpPosts, err error) {
	m = make(map[uint64]models.WpPosts)
	id := ids[0].([]uint64)
	arg := helper.SliceMap(id, helper.ToAny[uint64])
	rawPosts, err := models.Find[models.WpPosts](models.SqlBuilder{{
		"Id", "in", "",
	}}, "a.*,ifnull(d.name,'') category_name,ifnull(taxonomy,'') `taxonomy`", "", nil, models.SqlBuilder{{
		"a", "left join", "wp_term_relationships b", "a.Id=b.object_id",
	}, {
		"left join", "wp_term_taxonomy c", "b.term_taxonomy_id=c.term_taxonomy_id",
	}, {
		"left join", "wp_terms d", "c.term_id=d.term_id",
	}}, 0, arg)
	if err != nil {
		return m, err
	}
	postsMap := make(map[uint64]models.WpPosts)
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
	for k, pp := range postsMap {
		if len(pp.Categories) > 0 {
			t := make([]string, 0, len(pp.Categories))
			for _, cat := range pp.Categories {
				t = append(t, fmt.Sprintf(`<a href="/p/category/%s" rel="category tag">%s</a>`, cat, cat))
			}
			pp.CategoriesHtml = strings.Join(t, "、")
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

func PostLists(ctx context.Context, key string, args ...any) (r []models.WpPosts, total int, err error) {
	ids, err := postListIdsCache.GetCache(ctx, key, time.Second, args...)
	if err != nil {
		return
	}
	total = ids.Length
	r, err = GetPostsByIds(ctx, ids.Ids)
	return
}

func searchPostIds(args ...any) (ids PostIds, err error) {
	where := args[0].(models.SqlBuilder)
	page := args[1].(int)
	limit := args[2].(int)
	order := args[3].(models.SqlBuilder)
	join := args[4].(models.SqlBuilder)
	postType := args[5].([]any)
	postStatus := args[6].([]any)
	res, total, err := models.SimplePagination[models.WpPosts](where, "ID", "", page, limit, order, join, postType, postStatus)
	for _, posts := range res {
		ids.Ids = append(ids.Ids, posts.Id)
	}
	ids.Length = total
	if total > TotalRaw {
		TotalRaw = total
	}
	return
}

func getMaxPostId(...any) ([]uint64, error) {
	r, err := models.Find[models.WpPosts](models.SqlBuilder{{"post_type", "post"}, {"post_status", "publish"}}, "max(ID) ID", "", nil, nil, 0)
	var id uint64
	if len(r) > 0 {
		id = r[0].Id
	}
	return []uint64{id}, err
}

func GetMaxPostId(ctx *gin.Context) (uint64, error) {
	Id, err := maxPostIdCache.GetCache(ctx, time.Second)
	return Id[0], err
}
