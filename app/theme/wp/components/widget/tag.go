package widget

import (
	"github.com/fthvgb1/wp-go/app/pkg/cache"
	"github.com/fthvgb1/wp-go/app/pkg/constraints"
	"github.com/fthvgb1/wp-go/app/pkg/models"
	"github.com/fthvgb1/wp-go/app/theme/wp"
	"github.com/fthvgb1/wp-go/helper/slice"
	str "github.com/fthvgb1/wp-go/helper/strings"
)

func IsTag(h *wp.Handle) (models.TermsMy, bool) {
	if h.Scene() == constraints.Tag {
		id := str.ToInt[uint64](h.C.Query("tag"))
		i, t := slice.SearchFirst(cache.CategoriesTags(h.C, constraints.Tag), func(my models.TermsMy) bool {
			return id == my.Terms.TermId
		})
		if i > 0 {
			return t, true
		}
	}
	return models.TermsMy{}, false
}
