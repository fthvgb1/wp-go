package common

import (
	"context"
	"fmt"
	"github/fthvgb1/wp-go/helper"
	"github/fthvgb1/wp-go/models"
	"strings"
	"time"
)

func GetPostAndCache(ctx context.Context, id uint64) (models.WpPosts, error) {

	return postsCache.GetCache(ctx, id, time.Second, id)
}

func GetPost(id uint64) models.WpPosts {
	return postsCache.Get(id)
}

func GetPosts(ctx context.Context, ids []uint64) ([]models.WpPosts, error) {
	return postsCache.GetCacheBatch(ctx, ids, time.Second, ids)
}

func SetPostCache(ids []models.WpPosts) error {
	var arg []uint64
	for _, posts := range ids {
		arg = append(arg, posts.Id)
	}
	return postsCache.SetByBatchFn(arg)
}

func getPost(id ...any) (post models.WpPosts, err error) {
	Id := id[0].(uint64)
	posts, err := getPosts([]uint64{Id})
	if err != nil {
		return models.WpPosts{}, err
	}
	post = posts[Id]
	return
}

func getPosts(ids ...any) (m map[uint64]models.WpPosts, err error) {
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
