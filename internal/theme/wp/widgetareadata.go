package wp

import (
	"github.com/fthvgb1/wp-go/helper/slice"
	"github.com/fthvgb1/wp-go/internal/pkg/cache"
	"github.com/fthvgb1/wp-go/internal/pkg/constraints"
)

func (h *Handle) WidgetAreaData() {
	h.ginH["archives"] = cache.Archives(h.C)
	h.ginH["recentPosts"] = slice.Map(cache.RecentPosts(h.C, 5), ProjectTitle)
	h.ginH["categories"] = cache.CategoriesTags(h.C, constraints.Category)
	h.ginH["recentComments"] = cache.RecentComments(h.C, 5)
}
