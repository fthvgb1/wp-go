package actions

import (
	"context"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github/fthvgb1/wp-go/actions/common"
	"github/fthvgb1/wp-go/helper"
	"github/fthvgb1/wp-go/models"
	"github/fthvgb1/wp-go/plugins"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type indexHandle struct {
	c              *gin.Context
	session        sessions.Session
	page           int
	pageSize       int
	title          string
	titleL         string
	titleR         string
	search         string
	totalPage      int
	category       string
	categoryType   string
	where          models.SqlBuilder
	orderBy        models.SqlBuilder
	order          string
	join           models.SqlBuilder
	postType       []any
	status         []any
	header         string
	paginationStep int
	scene          uint
}

func newIndexHandle(ctx *gin.Context) *indexHandle {
	return &indexHandle{
		c:              ctx,
		session:        sessions.Default(ctx),
		page:           1,
		pageSize:       10,
		paginationStep: 1,
		titleL:         models.Options["blogname"],
		titleR:         models.Options["blogdescription"],
		where: models.SqlBuilder{
			{"post_type", "in", ""},
			{"post_status", "in", ""},
		},
		orderBy:  models.SqlBuilder{},
		join:     models.SqlBuilder{},
		postType: []any{"post"},
		status:   []any{"publish"},
		scene:    plugins.Home,
	}
}
func (h *indexHandle) setTitleLR(l, r string) {
	h.titleL = l
	h.titleR = r
}

func (h *indexHandle) getTitle() string {
	h.title = fmt.Sprintf("%s-%s", h.titleL, h.titleR)
	return h.title
}

func (h *indexHandle) parseParams() {
	h.order = h.c.Query("order")
	if !helper.IsContainInArr(h.order, []string{"asc", "desc"}) {
		h.order = "asc"
	}
	year := h.c.Param("year")
	if year != "" {
		h.where = append(h.where, []string{
			"year(post_date)", year,
		})
	}
	month := h.c.Param("month")
	if month != "" {
		h.where = append(h.where, []string{
			"month(post_date)", month,
		})
		ss := fmt.Sprintf("%s年%s月", year, strings.TrimLeft(month, "0"))
		h.header = fmt.Sprintf("月度归档： <span>%s</span>", ss)
		h.setTitleLR(ss, models.Options["blogname"])
		h.scene = plugins.Archive
	}
	category := h.c.Param("category")
	if category == "" {
		category = h.c.Param("tag")
		if category != "" {
			h.categoryType = "post_tag"
			h.header = fmt.Sprintf("标签： <span>%s</span>", category)
		}
	} else {
		h.categoryType = "category"
		h.header = fmt.Sprintf("分类： <span>%s</span>", category)
	}
	h.category = category

	if category != "" {
		h.where = append(h.where, []string{
			"d.name", category,
		}, []string{"taxonomy", h.categoryType})
		h.join = append(h.join, []string{
			"a", "left join", "wp_term_relationships b", "a.Id=b.object_id",
		}, []string{
			"left join", "wp_term_taxonomy c", "b.term_taxonomy_id=c.term_taxonomy_id",
		}, []string{
			"left join", "wp_terms d", "c.term_id=d.term_id",
		})
		h.setTitleLR(category, models.Options["blogname"])
		h.scene = plugins.Category
	}
	s := h.c.Query("s")
	if s != "" && strings.Replace(s, " ", "", -1) != "" {
		q := helper.StrJoin("%", s, "%")
		h.where = append(h.where, []string{
			"and", "post_title", "like", q, "",
			"or", "post_content", "like", q, "",
			"or", "post_excerpt", "like", q, "",
		}, []string{"post_password", ""})
		h.postType = append(h.postType, "page", "attachment")
		h.header = fmt.Sprintf("%s的搜索结果", s)
		h.setTitleLR(helper.StrJoin(`"`, s, `"`, "的搜索结果"), models.Options["blogname"])
		h.search = s
		h.scene = plugins.Search
	} else {
		h.status = append(h.status, "private")
	}
	p := h.c.Query("paged")
	if p == "" {
		p = h.c.Param("page")
	}
	if p != "" {
		if pa, err := strconv.Atoi(p); err == nil {
			h.page = pa
		}
	}
	if h.page > 1 && (h.category != "" || h.search != "" || month != "") {
		h.setTitleLR(fmt.Sprintf("%s-第%d页", h.titleL, h.page), models.Options["blogname"])
	}
}

func (h *indexHandle) getTotalPage(totalRaws int) int {
	h.totalPage = int(math.Ceil(float64(totalRaws) / float64(h.pageSize)))
	return h.totalPage
}

func Index(c *gin.Context) {
	h := newIndexHandle(c)
	h.parseParams()
	ctx := context.TODO()
	archive := common.Archives()
	recent := common.RecentPosts(ctx)
	categoryItems := common.Categories(ctx)
	recentComments := common.RecentComments(ctx)
	ginH := gin.H{
		"options":        models.Options,
		"recentPosts":    recent,
		"archives":       archive,
		"categories":     categoryItems,
		"search":         h.search,
		"header":         h.header,
		"title":          h.getTitle(),
		"recentComments": recentComments,
	}
	postIds, totalRaw, err := models.SimplePagination[models.WpPosts](h.where, "ID", "", h.page, h.pageSize, h.orderBy, h.join, h.postType, h.status)
	defer func() {
		c.HTML(http.StatusOK, "posts/index.gohtml", ginH)
		if err != nil {
			c.Error(err)
		}
	}()
	if err != nil {
		return
	}
	if len(postIds) < 1 && h.category != "" {
		h.titleL = "未找到页面"
	}
	err = common.QueryAndSetPostCache(postIds)
	pw := h.session.Get("post_password")
	c.Set("post_password", pw)
	plug := plugins.NewPostPlugin(c, h.scene)
	for i, v := range postIds {
		post, _ := common.PostsCache.Load(v.Id)
		pp := post.(*models.WpPosts)
		px := *pp
		common.PasswordProjectTitle(&px)
		if px.PostPassword != "" && pw != px.PostPassword {
			common.PasswdProjectContent(&px)
		}
		postIds[i] = px
		plugins.ApplyPlugin(plug, &postIds[i])

	}
	for i, post := range recent {

		if post.PostPassword != "" && pw != post.PostPassword {
			common.PasswdProjectContent(&recent[i])
		}
	}
	q := c.Request.URL.Query().Encode()
	ginH["posts"] = postIds
	ginH["totalPage"] = h.getTotalPage(totalRaw)
	ginH["pagination"] = pagination(h.page, h.totalPage, h.paginationStep, c.Request.URL.Path, q)
	ginH["title"] = h.getTitle()
}

var complie = regexp.MustCompile(`(/page)/(\d+)`)

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
	r := complie
	start := currentPage - step
	end := currentPage + step
	if start < 1 {
		start = 1
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
