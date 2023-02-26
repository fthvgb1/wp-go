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

func (h *Handle) CalBodyClass() {
	h.GinH["bodyClass"] = h.BodyClass(h.Class...)
}

func (h *Handle) BodyClass(class ...string) string {
	if constraints.Ok != h.Stats {
		return "error404"
	}
	switch h.Scene {
	case constraints.Home:
		class = append(class, "home", "blog")

	case constraints.Archive:
		class = append(class, "archive", "date")

	case constraints.Search:
		s := "search-no-results"
		if len(h.Index.Posts) > 0 {
			s = "search-results"
		}
		class = append(class, "search", s)

	case constraints.Category, constraints.Tag:
		class = append(class, "archive", "category")
		cat := h.Index.Param.Category
		_, cate := slice.SearchFirst(cache.CategoriesTags(h.C, h.Scene), func(my models.TermsMy) bool {
			return my.Name == cat
		})
		if cate.Slug[0] != '%' {
			class = append(class, str.Join("category-", cate.Slug))
		}
		class = append(class, str.Join("category-", number.ToString(cate.Terms.TermId)))

	case constraints.Author:
		class = append(class, "archive", "author")
		author := h.Index.Param.Author
		user, _ := cache.GetUserByName(h.C, author)
		class = append(class, str.Join("author-", number.ToString(user.Id)))
		if user.UserLogin[0] != '%' {
			class = append(class, str.Join("author-", user.UserLogin))
		}

	case constraints.Detail:
		class = append(class, "post-template-default", "single", "single-post")
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
	return strings.Join(class, " ")
}
