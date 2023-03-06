package wp

import (
	"github.com/fthvgb1/wp-go/helper/slice"
	"github.com/fthvgb1/wp-go/internal/pkg/cache"
	"github.com/fthvgb1/wp-go/internal/pkg/constraints"
	"github.com/fthvgb1/wp-go/internal/wpconfig"
)

func (h *Handle) WidgetAreaData() {
	h.ginH["showArchiveCount"] = wpconfig.GetPHPArrayValWithDefaults[int64]("widget_archives", 0, int64(2), "count")
	h.ginH["archiveDropdown"] = wpconfig.GetPHPArrayValWithDefaults[int64]("widget_archives", 0, int64(2), "dropdown")
	h.ginH["archiveTitle"] = wpconfig.GetPHPArrayValWithDefaults[string]("widget_archives", "归档", int64(2), "title")

	h.ginH["archives"] = cache.Archives(h.C)
	h.ginH["recentPosts"] = slice.Map(cache.RecentPosts(h.C, 5), ProjectTitle)
	h.ginH["categories"] = cache.CategoriesTags(h.C, constraints.Category)
	h.ginH["recentComments"] = cache.RecentComments(h.C, 5)
}
