package common

import (
	"context"
	"github/fthvgb1/wp-go/helper"
	"github/fthvgb1/wp-go/models"
	"time"
)

func GetPostById(ctx context.Context, id uint64, ids ...uint64) (models.WpPosts, error) {

	return postsCache.GetCacheBatch(ctx, id, time.Second, ids)
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
	if err == nil {
		for _, v := range rawPosts {
			m[v.Id] = v
		}
	}
	return
}
