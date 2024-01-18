package wp

import (
	"fmt"
	"github.com/fthvgb1/wp-go/app/pkg/cache"
	"github.com/fthvgb1/wp-go/app/pkg/constraints"
	"github.com/fthvgb1/wp-go/app/pkg/models"
	"github.com/fthvgb1/wp-go/app/wpconfig"
	"github.com/fthvgb1/wp-go/helper/number"
	"github.com/fthvgb1/wp-go/helper/slice"
	str "github.com/fthvgb1/wp-go/helper/strings"
	"strconv"
	"strings"
)

func bodyClass(h *Handle) func() string {
	return func() string {
		return h.BodyClass()
	}
}

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
		if len(h.GetIndexHandle().Posts) > 0 {
			s = "search-results"
		}
		class = append(class, "search", s)

	case constraints.Category, constraints.Tag:
		class = append(class, "archive", "category")
		cat := h.GetIndexHandle().Param.Category
		if cat == "" {
			break
		}
		_, cate := slice.SearchFirst(cache.CategoriesTags(h.C, h.scene), func(my models.TermsMy) bool {
			return my.Name == cat
		})
		if cate.Slug[0] != '%' {
			class = append(class, str.Join("category-", cate.Slug))
		}
		class = append(class, str.Join("category-", number.IntToString(cate.Terms.TermId)))

	case constraints.Author:
		class = append(class, "archive", "author")
		author := h.GetIndexHandle().Param.Author
		user, _ := cache.GetUserByName(h.C, author)
		class = append(class, str.Join("author-", number.IntToString(user.Id)))
		if user.DisplayName[0] != '%' {
			class = append(class, str.Join("author-", user.DisplayName))
		}

	case constraints.Detail:
		class = append(class, "post-template-default", "single", "single-post")
		class = append(class, str.Join("postid-", number.IntToString(h.GetDetailHandle().Post.Id)))
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
	return h.DoActionFilter("bodyClass", strings.Join(class, " "))
}

func postClass(h *Handle) func(posts models.Posts) string {
	return func(posts models.Posts) string {
		return h.PostClass(posts)
	}
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
		if h.GetPassword() != posts.PostPassword {
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

	return h.DoActionFilter("postClass", strings.Join(class, " "))
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
