package wp

import (
	"github.com/fthvgb1/wp-go/helper/slice"
	"github.com/fthvgb1/wp-go/internal/cmd/reload"
	"github.com/fthvgb1/wp-go/internal/pkg/cache"
	"github.com/fthvgb1/wp-go/internal/pkg/constraints"
	"github.com/fthvgb1/wp-go/internal/wpconfig"
)

func (h *Handle) WidgetAreaData() {
	h.Archives()
	h.RecentPosts()
	h.ginH["categories"] = cache.CategoriesTags(h.C, constraints.Category)
	h.ginH["recentComments"] = cache.RecentComments(h.C, 5)
}

var recentConf = map[any]any{
	"number":    int64(5),
	"show_date": 0,
	"title":     "近期文章",
}

func (h *Handle) RecentPosts() {
	set := reload.GetAnyValBy("recentPostsConfig", func() map[any]any {
		return wpconfig.GetPHPArrayVal[map[any]any]("widget_recent-posts", recentConf, int64(2))
	})
	h.ginH["recentPostsConfig"] = set
	h.ginH["recentPosts"] = slice.Map(cache.RecentPosts(h.C, int(set["number"].(int64))), ProjectTitle)
}

var archivesConfig = map[any]any{
	"count":    0,
	"dropdown": 0,
	"title":    "归档",
}

func (h *Handle) Archives() {
	h.ginH["archivesConfig"] = reload.GetAnyValBy("archivesConfig", func() map[any]any {
		return wpconfig.GetPHPArrayVal[map[any]any]("widget_archives", archivesConfig, int64(2))
	})
	h.ginH["archives"] = cache.Archives(h.C)
}
