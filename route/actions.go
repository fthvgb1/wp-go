package route

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github/fthvgb1/wp-go/helper"
	"github/fthvgb1/wp-go/models"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

var PostsCache sync.Map

func index(c *gin.Context) {
	page := 1
	pageSize := 10
	order := c.Query("order")
	if !helper.IsContainInArr(order, []string{"asc", "desc"}) {
		order = "asc"
	}
	where := models.SqlBuilder{{
		"post_type", "post",
	}, {"post_status", "in", ""}}
	p := c.Query("paged")
	year := c.Param("year")
	if year != "" {
		where = append(where, []string{
			"year(post_date)", year,
		})
	}
	month := c.Param("month")
	if month != "" {
		where = append(where, []string{
			"month(post_date)", month,
		})
	}
	tt := ""
	category := c.Param("category")
	if category == "" {
		category = c.Param("tag")
		if category != "" {
			tt = "post_tag"
		}
	} else {
		tt = "category"
	}
	var join models.SqlBuilder
	if category != "" {
		where = append(where, []string{
			"d.name", category,
		}, []string{"taxonomy", tt})
		join = append(join, []string{
			"a", "left join", "wp_term_relationships b", "a.Id=b.object_id",
		}, []string{
			"left join", "wp_term_taxonomy c", "b.term_taxonomy_id=c.term_taxonomy_id",
		}, []string{
			"left join", "wp_terms d", "c.term_id=d.term_id",
		})
	}
	if p == "" {
		p = c.Param("page")
	}
	if p != "" {
		if pa, err := strconv.Atoi(p); err == nil {
			page = pa
		}
	}

	status := []interface{}{"publish", "private"}
	postIds, totalRaw, err := models.SimplePagination[models.WpPosts](where, "ID", "", page, pageSize, models.SqlBuilder{{"post_date", order}}, join, status)
	defer func() {
		if err != nil {
			c.Error(err)
		}
	}()
	if err != nil {
		return
	}
	var all []uint64
	var allPosts []models.WpPosts
	var needQuery []interface{}
	for _, wpPosts := range postIds {
		all = append(all, wpPosts.Id)
		if _, ok := PostsCache.Load(wpPosts.Id); !ok {
			needQuery = append(needQuery, wpPosts.Id)
		}
	}
	if len(needQuery) > 0 {
		rawPosts, err := models.Find[models.WpPosts](models.SqlBuilder{{
			"Id", "in", "",
		}}, "a.*,d.name category_name,taxonomy", "", nil, models.SqlBuilder{{
			"a", "left join", "wp_term_relationships b", "a.Id=b.object_id",
		}, {
			"left join", "wp_term_taxonomy c", "b.term_taxonomy_id=c.term_taxonomy_id",
		}, {
			"left join", "wp_terms d", "c.term_id=d.term_id",
		}}, 0, needQuery)
		if err != nil {
			return
		}
		postsMap := make(map[uint64]*models.WpPosts)
		for i, post := range rawPosts {
			v, ok := postsMap[post.Id]
			if !ok {
				v = &rawPosts[i]
			}
			if post.Taxonomy == "category" {
				v.Categories = append(v.Categories, post.CategoryName)
			} else if post.Taxonomy == "post_tag" {
				v.Tags = append(v.Tags, post.CategoryName)
			}
			postsMap[post.Id] = v
		}
		for _, pp := range postsMap {
			if len(pp.Categories) > 0 {
				t := make([]string, 0, len(pp.Categories))
				for _, cat := range pp.Categories {
					t = append(t, fmt.Sprintf(`<a href="/p/category/%s" rel="category tag">%s</a>`, cat, cat))
				}
				pp.CategoriesHtml = strings.Join(t, "、")
			}
			if len(pp.Tags) > 0 {
				t := make([]string, 0, len(pp.Tags))
				for _, cat := range pp.Tags {
					t = append(t, fmt.Sprintf(`<a href="/p/tag/%s" rel="tag">%s</a>`, cat, cat))
				}
				pp.TagsHtml = strings.Join(t, "、")
			}
			PostsCache.Store(pp.Id, pp)
		}
	}

	for _, id := range all {
		post, _ := PostsCache.Load(id)
		pp := post.(*models.WpPosts)
		allPosts = append(allPosts, *pp)
	}
	recent, err := recentPosts()
	archive, err := archives()
	categoryItems, err := categories()
	totalPage := int(math.Ceil(float64(totalRaw) / float64(pageSize)))
	q := c.Request.URL.Query().Encode()
	if q != "" {
		q = fmt.Sprintf("?%s", q)
	}
	c.HTML(http.StatusOK, "index.html", gin.H{
		"posts":       allPosts,
		"options":     models.Options,
		"recentPosts": recent,
		"archives":    archive,
		"categories":  categoryItems,
		"totalPage":   totalPage,
		"queryRaw":    q,
		"pagination":  pagination(page, totalPage, 2, c.Request.URL.Path, q),
	})
}

func recentPosts() (r []models.WpPosts, err error) {
	r, err = models.Find[models.WpPosts](models.SqlBuilder{{
		"post_type", "post",
	}, {"post_status", "publish"}}, "ID,post_title", "", models.SqlBuilder{{"post_date", "desc"}}, nil, 5)
	return
}

func categories() (terms []models.WpTermsMy, err error) {
	var in = []interface{}{"category"}
	terms, err = models.Find[models.WpTermsMy](models.SqlBuilder{
		{"tt.count", ">", "0", "int"},
		{"tt.taxonomy", "in", ""},
	}, "t.term_id", "", models.SqlBuilder{
		{"t.name", "asc"},
	}, models.SqlBuilder{
		{"t", "inner join", "wp_term_taxonomy tt", "t.term_id = tt.term_id"},
	}, 0, in)
	for i := 0; i < len(terms); i++ {
		if v, ok := models.Terms[terms[i].WpTerms.TermId]; ok {
			terms[i].WpTerms = v
		}
		if v, ok := models.TermTaxonomy[terms[i].WpTerms.TermId]; ok {
			terms[i].WpTermTaxonomy = v
		}
	}

	return
}

func archives() (r []models.PostArchive, err error) {
	r, err = models.Find[models.PostArchive](models.SqlBuilder{
		{"post_type", "post"}, {"post_status", "publish"},
	}, "YEAR(post_date) AS `year`, MONTH(post_date) AS `month`, count(ID) as posts", "year,month", models.SqlBuilder{{"year", "desc"}, {"month", "desc"}}, nil, 0)
	return
}

func pagination(currentPage, totalPage, step int, path, query string) (html string) {
	if totalPage < 2 {
		return
	}
	pathx := path
	if !strings.Contains(path, "/page/") {
		pathx = fmt.Sprintf("%s%s", path, "/page/1")
	}
	s := strings.Builder{}
	if currentPage > totalPage {
		currentPage = totalPage
	}
	r := regexp.MustCompile(`(/page)/(\d+)`)
	start := currentPage - step
	end := currentPage + step
	if start < 1 {
		start = currentPage
	}
	if currentPage > 1 {
		pp := ""
		if currentPage >= 2 {
			pp = replacePage(r, pathx, currentPage-1)
		}
		s.WriteString(fmt.Sprintf(`<a class="prev page-numbers" href="%s%s">上一页</a>`, pp, query))
	}
	if currentPage >= step+2 {
		d := ""
		if currentPage > step+2 {
			d = `<span class="page-numbers dots">…</span>`
		}
		e := replacePage(r, path, 1)
		s.WriteString(fmt.Sprintf(`
<a class="page-numbers" href="%s%s"><span class="meta-nav screen-reader-text">页 </span>1</a>
%s
`, e, query, d))
	}
	if totalPage < end {
		end = totalPage
	}

	for page := start; page <= end; page++ {
		h := ""
		if currentPage == page {
			h = fmt.Sprintf(`
<span aria-current="page" class="page-numbers current">
            <span class="meta-nav screen-reader-text">页 </span>%d</span>
`, page)

		} else {
			d := replacePage(r, pathx, page)
			h = fmt.Sprintf(`
<a class="page-numbers" href="%s%s">
<span class="meta-nav screen-reader-text">页 </span>%d</a>
`, d, query, page)
		}
		s.WriteString(h)

	}
	if totalPage >= currentPage+step+1 {
		if totalPage > currentPage+step+1 {
			s.WriteString(`<span class="page-numbers dots">…</span>`)
		}
		dd := replacePage(r, pathx, totalPage)
		s.WriteString(fmt.Sprintf(`
<a class="page-numbers" href="%s%s"><span class="meta-nav screen-reader-text">页 </span>%d</a>`, dd, query, totalPage))
	}
	if currentPage < totalPage {
		dd := replacePage(r, pathx, currentPage+1)
		s.WriteString(fmt.Sprintf(`<a class="next page-numbers" href="%s%s">下一页</a>`, dd, query))
	}
	html = s.String()
	return
}

func replacePage(r *regexp.Regexp, path string, page int) (src string) {
	if page == 1 {
		src = r.ReplaceAllString(path, "")
	} else {
		s := fmt.Sprintf("$1/%d", page)
		src = r.ReplaceAllString(path, s)
	}
	src = strings.Replace(src, "//", "/", -1)
	if src == "" {
		src = "/"
	}
	return
}
