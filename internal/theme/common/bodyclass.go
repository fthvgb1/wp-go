package common

import (
	"fmt"
	"github.com/fthvgb1/wp-go/helper/slice"
	str "github.com/fthvgb1/wp-go/helper/strings"
	"github.com/fthvgb1/wp-go/internal/pkg/cache"
	"github.com/fthvgb1/wp-go/internal/pkg/constraints"
	"github.com/fthvgb1/wp-go/internal/pkg/models"
	"github.com/fthvgb1/wp-go/internal/wpconfig"
	"strings"
)

var commonClass = map[int]string{
	constraints.Home:     "home blog ",
	constraints.Archive:  "archive date ",
	constraints.Category: "archive category ",
	constraints.Tag:      "archive category ",
	constraints.Search:   "search ",
	constraints.Detail:   "post-template-default single single-post ",
}

func (h *Handle) CalBodyClass() {
	h.GinH["bodyClass"] = h.BodyClass(h.Class...)
}

func (h *Handle) BodyClass(class ...string) string {
	s := ""
	if constraints.Ok != h.Stats {
		return "error404"
	}
	switch h.Scene {
	case constraints.Search:
		s = "search-no-results"
		if len(h.GinH["posts"].([]models.Posts)) > 0 {
			s = "search-results"
		}
	case constraints.Category, constraints.Tag:
		cat := h.C.Param("category")
		if cat == "" {
			cat = h.C.Param("tag")
		}
		_, cate := slice.SearchFirst(cache.CategoriesTags(h.C, h.Scene), func(my models.TermsMy) bool {
			return my.Name == cat
		})
		if cate.Slug[0] != '%' {
			s = cate.Slug
		}
		s = fmt.Sprintf("category-%v category-%v", s, cate.Terms.TermId)
	case constraints.Detail:
		s = fmt.Sprintf("postid-%d", h.GinH["post"].(models.Posts).Id)
		if len(h.ThemeMods.ThemeSupport.PostFormats) > 0 {
			s = str.Join(s, " single-format-standard")
		}
	}
	class = append(class, s)

	if wpconfig.IsCustomBackground(h.Theme) {
		class = append(class, "custom-background")
	}
	if h.ThemeMods.CustomLogo > 0 {
		class = append(class, "wp-custom-logo")
	}
	if h.ThemeMods.ThemeSupport.ResponsiveEmbeds {
		class = append(class, "wp-embed-responsive")
	}

	return str.Join(commonClass[h.Scene], strings.Join(class, " "))
}
