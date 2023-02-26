package common

import (
	"github.com/fthvgb1/wp-go/helper/number"
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
	constraints.Author:   "archive author ",
	constraints.Detail:   "post-template-default single single-post ",
}

func (h *Handle) CalBodyClass() {
	h.GinH["bodyClass"] = h.BodyClass(h.Class...)
}

func (h *Handle) BodyClass(class ...string) string {
	if constraints.Ok != h.Stats {
		return "error404"
	}
	switch h.Scene {
	case constraints.Search:
		s := "search-no-results"
		if len(h.Index.Posts) > 0 {
			s = "search-results"
		}
		class = append(class, s)
	case constraints.Category, constraints.Tag:
		cat := h.Index.Param.Category
		_, cate := slice.SearchFirst(cache.CategoriesTags(h.C, h.Scene), func(my models.TermsMy) bool {
			return my.Name == cat
		})
		if cate.Slug[0] != '%' {
			class = append(class, str.Join("category-", cate.Slug))
		}
		class = append(class, str.Join("category-", number.ToString(cate.Terms.TermId)))

	case constraints.Author:
		author := h.Index.Param.Author
		user, _ := cache.GetUserByName(h.C, author)
		class = append(class, str.Join("author-", number.ToString(user.Id)))
		if user.UserLogin[0] != '%' {
			class = append(class, str.Join("author-", user.UserLogin))
		}

	case constraints.Detail:
		class = append(class, str.Join("postid-", number.ToString(h.Detail.Post.Id)))
		if len(h.ThemeMods.ThemeSupport.PostFormats) > 0 {
			class = append(class, "single-format-standard")
		}
	}
	if wpconfig.IsCustomBackground(h.Theme) {
		class = append(class, "custom-background")
	}
	if h.ThemeMods.CustomLogo > 0 || str.ToInteger(wpconfig.GetOption("site_logo"), 0) > 0 {
		class = append(class, "wp-custom-logo")
	}
	if h.ThemeMods.ThemeSupport.ResponsiveEmbeds {
		class = append(class, "wp-embed-responsive")
	}
	class = append(class, strings.Fields(commonClass[h.Scene])...)
	return strings.Join(slice.Reverse(class), " ")
}
