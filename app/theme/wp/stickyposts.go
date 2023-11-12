package wp

import (
	"fmt"
	"github.com/elliotchance/phpserialize"
	"github.com/fthvgb1/wp-go/app/pkg/cache"
	"github.com/fthvgb1/wp-go/app/pkg/logs"
	"github.com/fthvgb1/wp-go/app/pkg/models"
	"github.com/fthvgb1/wp-go/app/wpconfig"
	"github.com/fthvgb1/wp-go/cache/reload"
	"github.com/fthvgb1/wp-go/helper/maps"
	"github.com/fthvgb1/wp-go/helper/slice"
	str "github.com/fthvgb1/wp-go/helper/strings"
)

func (h *Handle) StickPosts() []models.Posts {
	return reload.GetAnyValBys("stickPostsSlice", h, func(h *Handle) (r []models.Posts, ok bool) {
		v := wpconfig.GetOption("sticky_posts")
		if v == "" {
			return
		}
		array, err := phpserialize.UnmarshalIndexedArray([]byte(v))
		if err != nil {
			logs.Error(err, "解析option sticky_posts错误", v)
			return
		}
		r = slice.FilterAndMap(array, func(t any) (models.Posts, bool) {
			id := str.ToInt[uint64](fmt.Sprintf("%v", t))
			post, err := cache.GetPostById(h.C, id)
			post.IsSticky = true
			return post, err == nil
		})
		ok = true
		return
	})
}

func (h *Handle) StickMapPosts() map[uint64]models.Posts {
	return reload.GetAnyValBys("stickPostsMap", h, func(h *Handle) (map[uint64]models.Posts, bool) {
		return slice.SimpleToMap(h.StickPosts(), func(v models.Posts) uint64 {
			return v.Id
		}), true
	})
}

func (h *Handle) IsStick(id uint64) bool {
	return maps.IsExists(h.StickMapPosts(), id)
}
