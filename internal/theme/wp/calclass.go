package wp

import (
	"fmt"
	"github.com/fthvgb1/wp-go/helper/number"
	"github.com/fthvgb1/wp-go/helper/slice"
	str "github.com/fthvgb1/wp-go/helper/strings"
	"github.com/fthvgb1/wp-go/internal/pkg/cache"
	"github.com/fthvgb1/wp-go/internal/pkg/constraints"
	"github.com/fthvgb1/wp-go/internal/pkg/models"
	"github.com/fthvgb1/wp-go/internal/wpconfig"
	"strconv"
	"strings"
)

func (h *Handle) BodyClass() string {
	var class []string
	if constraints.Ok != h.Stats {
		class = append(class, "error404")
	}
	switch h.scene {
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
		if cat == "" {
			break
		}
		_, cate := slice.SearchFirst(cache.CategoriesTags(h.C, h.scene), func(my models.TermsMy) bool {
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
		if len(h.themeMods.ThemeSupport.PostFormats) > 0 {
			class = append(class, "single-format-standard")
		}
	}
	if wpconfig.IsCustomBackground(h.theme) {
		class = append(class, "custom-background")
	}
	if h.themeMods.CustomLogo > 0 || str.ToInteger(wpconfig.GetOption("site_logo"), 0) > 0 {
		class = append(class, "wp-custom-logo")
	}
	if h.themeMods.ThemeSupport.ResponsiveEmbeds {
		class = append(class, "wp-embed-responsive")
	}
	return h.ComponentFilterFnHook("bodyClass", strings.Join(class, " "))
}
func (h *Handle) PostClass(posts models.Posts) string {
	var class []string
	class = append(class, fmt.Sprintf("post-%d", posts.Id), posts.PostType,
		str.Join("type-", posts.PostType), str.Join("status-", posts.PostStatus),
		"hentry", "format-standard")
	if h.CommonThemeMods().ThemeSupport.PostThumbnails && posts.Thumbnail.Path != "" {
		class = append(class, "has-post-thumbnail")
	}

	if posts.PostPassword != "" {
		if h.password != posts.PostPassword {
			class = append(class, "post-password-required")
		} else {
			class = append(class, "post-password-projected")
		}
	}

	if h.scene == constraints.Home && h.IsStick(posts.Id) {
		class = append(class, "sticky")
	}
	for _, id := range posts.TermIds {
		term, ok := wpconfig.GetTermMy(id)
		if !ok || term.Slug == "" {
			continue
		}
		class = append(class, TermClass(term))
	}

	return h.ComponentFilterFnHook("postClass", strings.Join(class, " "))
}

func TermClass(term models.TermsMy) string {
	termClass := term.Slug
	if strings.Contains(term.Slug, "%") {
		termClass = strconv.FormatUint(term.TermTaxonomy.TermId, 10)
	}
	switch term.Taxonomy {
	case "post_tag":
		return str.Join("tag-", termClass)
	case "post_format":
		return fmt.Sprintf("format-%s", strings.ReplaceAll(term.Slug, "post-format-", ""))
	}
	return str.Join(term.Taxonomy, "-", termClass)
}
