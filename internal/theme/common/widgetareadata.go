package common

import (
	"github.com/fthvgb1/wp-go/helper/slice"
	"github.com/fthvgb1/wp-go/internal/pkg/cache"
	"github.com/fthvgb1/wp-go/internal/pkg/constraints"
)

func (h *Handle) WidgetAreaData() {
	h.GinH["archives"] = cache.Archives(h.C)
	h.GinH["recentPosts"] = slice.Map(cache.RecentPosts(h.C, 5), ProjectTitle)
	h.GinH["categories"] = cache.CategoriesTags(h.C, constraints.Category)
	h.GinH["recentComments"] = cache.RecentComments(h.C, 5)
}
