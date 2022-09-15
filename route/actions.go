package route

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github/fthvgb1/wp-go/models"
	"math"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

var PostsCache sync.Map

func index(c *gin.Context) {
	page := 1
	pageSize := 10
	p := c.Query("paged")
	if p == "" {
		p = c.Param("page")
	}
	if p != "" {
		if pa, err := strconv.Atoi(p); err == nil {
			page = pa
		}
	}

	status := []interface{}{"publish", "private"}
	posts, totalRaw, err := models.SimplePagination[models.WpPosts](models.SqlBuilder{{
		"post_type", "post",
	}, {"post_status", "in", ""}}, "ID", page, pageSize, models.SqlBuilder{{"post_date", "desc"}}, nil, status)
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
	for _, wpPosts := range posts {
		all = append(all, wpPosts.Id)
		if _, ok := PostsCache.Load(wpPosts.Id); !ok {
			needQuery = append(needQuery, wpPosts.Id)
		}
	}
	if len(needQuery) > 0 {
		rawPosts, err := models.Find[models.WpPosts](models.SqlBuilder{{
			"Id", "in", "",
		}}, "a.*,d.name category_name", "", nil, models.SqlBuilder{{
			"a", "left join", "wp_term_relationships b", "a.Id=b.object_id",
		}, {
			"left join", "wp_term_taxonomy c", "b.term_taxonomy_id=c.term_taxonomy_id",
		}, {
			"left join", "wp_terms d", "c.term_id=d.term_id",
		}}, 0, needQuery)
		if err != nil {
			return
		}
		for _, post := range rawPosts {
			PostsCache.Store(post.Id, post)
		}
	}

	for _, id := range all {
		post, _ := PostsCache.Load(id)
		pp := post.(models.WpPosts)
		allPosts = append(allPosts, pp)
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
		"pagination":  pagination(page, totalPage, 1, q),
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

func pagination(currentPage, totalPage, step int, query string) (html string) {
	html = ""
	s := strings.Builder{}
	if currentPage > totalPage {
		currentPage = totalPage
	}

	start := currentPage - step
	end := currentPage + step
	if start < 1 {
		start = currentPage
	}
	if currentPage > 1 {
		pp := ""
		if currentPage > 2 {
			pp = fmt.Sprintf("page/%d", currentPage-1)
		}
		s.WriteString(fmt.Sprintf(`<a class="prev page-numbers" href="/%s%s">上一页</a>`, pp, query))
	}
	if currentPage >= step+2 {
		d := ""
		if currentPage > step+2 {
			d = `<span class="page-numbers dots">…</span>`
		}
		s.WriteString(fmt.Sprintf(`
<a class="page-numbers" href="/%s"><span class="meta-nav screen-reader-text">页 </span>1</a>
%s
`, query, d))
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
			d := fmt.Sprintf("/page/%d", page)
			if currentPage > page && page == 1 {
				d = "/"
			}
			h = fmt.Sprintf(`
<a class="page-numbers" href="%s%s">
<span class="meta-nav screen-reader-text">页 </span>%d</a>
`, d, query, page)
		}
		s.WriteString(h)

	}
	if totalPage > currentPage+step+2 {
		s.WriteString(fmt.Sprintf(`
<span class="page-numbers dots">…</span>
<a class="page-numbers" href="/page/%d%s"><span class="meta-nav screen-reader-text">页 </span>%d</a>`, totalPage, query, totalPage))
	}
	if currentPage < totalPage {
		s.WriteString(fmt.Sprintf(`<a class="next page-numbers" href="/page/%d%s">下一页</a>`, currentPage+1, query))
	}
	html = s.String()
	return
}
