package wp

import (
	"github.com/fthvgb1/wp-go/helper/slice"
	"github.com/fthvgb1/wp-go/internal/pkg/cache"
	"github.com/fthvgb1/wp-go/internal/pkg/constraints"
	"github.com/fthvgb1/wp-go/internal/wpconfig"
)

func (h *Handle) WidgetAreaData() {
	h.Archives()
	h.RecentPosts()
	h.RecentComments()
	h.ginH["searchConf"] = wpconfig.GetPHPArrayVal("widget_search", "", int64(2), "title")
	h.ginH["categories"] = cache.CategoriesTags(h.C, constraints.Category)
}

var recentConf = map[any]any{
	"number":    int64(5),
	"show_date": 0,
	"title":     "近期文章",
}

func (h *Handle) RecentPosts() {
	set := wpconfig.GetPHPArrayVal[map[any]any]("widget_recent-posts", recentConf, int64(2))
	h.ginH["recentPostsConfig"] = set
	h.ginH["recentPosts"] = slice.Map(cache.RecentPosts(h.C, int(set["number"].(int64))), ProjectTitle)
}

var recentCommentConf = map[any]any{
	"number": int64(5),
	"title":  "近期文章",
}

func (h *Handle) RecentComments() {
	set := wpconfig.GetPHPArrayVal[map[any]any]("widget_recent-comments", recentCommentConf, int64(2))
	h.ginH["recentCommentsConfig"] = set
	h.ginH["recentComments"] = cache.RecentComments(h.C, int(set["number"].(int64)))
}

var archivesConfig = map[any]any{
	"count":    0,
	"dropdown": 0,
	"title":    "归档",
}

func (h *Handle) Archives() {
	h.ginH["archivesConfig"] = wpconfig.GetPHPArrayVal[map[any]any]("widget_archives", archivesConfig, int64(2))
	h.ginH["archives"] = cache.Archives(h.C)
}
