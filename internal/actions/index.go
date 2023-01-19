package actions

import (
	"fmt"
	"github.com/fthvgb1/wp-go/helper"
	"github.com/fthvgb1/wp-go/internal/pkg/cache"
	dao "github.com/fthvgb1/wp-go/internal/pkg/dao"
	"github.com/fthvgb1/wp-go/internal/pkg/logs"
	"github.com/fthvgb1/wp-go/internal/pkg/models"
	"github.com/fthvgb1/wp-go/internal/plugins"
	"github.com/fthvgb1/wp-go/internal/theme"
	"github.com/fthvgb1/wp-go/internal/wpconfig"
	"github.com/fthvgb1/wp-go/model"
	"github.com/fthvgb1/wp-go/plugin/pagination"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"
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
	author         string
	totalPage      int
	category       string
	categoryType   string
	where          model.SqlBuilder
	orderBy        string
	order          string
	join           model.SqlBuilder
	postType       []any
	status         []any
	header         string
	paginationStep int
	scene          uint
}

func newIndexHandle(ctx *gin.Context) *indexHandle {
	size := wpconfig.Options.Value("posts_per_page")
	si, _ := strconv.Atoi(size)
	return &indexHandle{
		c:              ctx,
		session:        sessions.Default(ctx),
		page:           1,
		pageSize:       si,
		paginationStep: 1,
		titleL:         wpconfig.Options.Value("blogname"),
		titleR:         wpconfig.Options.Value("blogdescription"),
		where: model.SqlBuilder{
			{"post_type", "in", ""},
			{"post_status", "in", ""},
		},
		orderBy:  "post_date",
		join:     model.SqlBuilder{},
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

func (h *indexHandle) getSearchKey() string {
	return fmt.Sprintf("action:%s|%s|%s|%s|%s|%s|%d|%d", h.author, h.search, h.orderBy, h.order, h.category, h.categoryType, h.page, h.pageSize)
}

func (h *indexHandle) parseParams() (err error) {
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
		h.setTitleLR(ss, wpconfig.Options.Value("blogname"))
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
	username := h.c.Param("author")
	if username != "" {
		user, er := cache.GetUserByName(h.c, username)
		if er != nil {
			err = er
			return
		}
		h.author = username
		h.where = append(h.where, []string{
			"post_author", "=", strconv.FormatUint(user.Id, 10), "int",
		})
	}
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
		h.setTitleLR(category, wpconfig.Options.Value("blogname"))
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
		h.setTitleLR(helper.StrJoin(`"`, s, `"`, "的搜索结果"), wpconfig.Options.Value("blogname"))
		h.search = s
		h.scene = plugins.Search
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
	total := int(atomic.LoadInt64(&dao.TotalRaw))
	if total > 0 && total < (h.page-1)*h.pageSize {
		h.page = 1
	}
	if h.page > 1 && (h.category != "" || h.search != "" || month != "") {
		h.setTitleLR(fmt.Sprintf("%s-第%d页", h.titleL, h.page), wpconfig.Options.Value("blogname"))
	}
	return
}

func (h *indexHandle) getTotalPage(totalRaws int) int {
	h.totalPage = int(math.Ceil(float64(totalRaws) / float64(h.pageSize)))
	return h.totalPage
}

func Index(c *gin.Context) {
	h := newIndexHandle(c)
	var postIds []models.Posts
	var totalRaw int
	var err error
	archive := cache.Archives(c)
	recent := cache.RecentPosts(c, 5)
	categoryItems := cache.Categories(c)
	recentComments := cache.RecentComments(c, 5)
	ginH := gin.H{
		"options":        wpconfig.Options,
		"recentPosts":    recent,
		"archives":       archive,
		"categories":     categoryItems,
		"search":         h.search,
		"header":         h.header,
		"recentComments": recentComments,
	}
	defer func() {
		stat := http.StatusOK
		if err != nil {
			c.Error(err)
			stat = http.StatusInternalServerError
			return
		}
		t := getTemplateName()
		theme.Hook(t, stat, c, ginH, int(h.scene))
	}()
	err = h.parseParams()
	if err != nil {
		return
	}
	ginH["title"] = h.getTitle()
	if c.Param("month") != "" {
		postIds, totalRaw, err = cache.GetMonthPostIds(c, c.Param("year"), c.Param("month"), h.page, h.pageSize, h.order)
		if err != nil {
			return
		}
	} else if h.search != "" {
		postIds, totalRaw, err = cache.SearchPost(c, h.getSearchKey(), c, h.where, h.page, h.pageSize, model.SqlBuilder{{h.orderBy, h.order}}, h.join, h.postType, h.status)
	} else {
		postIds, totalRaw, err = cache.PostLists(c, h.getSearchKey(), c, h.where, h.page, h.pageSize, model.SqlBuilder{{h.orderBy, h.order}}, h.join, h.postType, h.status)
	}
	if err != nil {
		logs.ErrPrintln(err, "获取数据错误")
		return
	}

	if len(postIds) < 1 && h.category != "" {
		h.titleL = "未找到页面"
		h.scene = plugins.Empty404
	}

	pw := h.session.Get("post_password")
	plug := plugins.NewPostPlugin(c, h.scene)
	for i, post := range postIds {
		plugins.PasswordProjectTitle(&postIds[i])
		if post.PostPassword != "" && pw != post.PostPassword {
			plugins.PasswdProjectContent(&postIds[i])
		} else {
			plugins.ApplyPlugin(plug, &postIds[i])
		}
	}
	for i, post := range recent {
		if post.PostPassword != "" && pw != post.PostPassword {
			plugins.PasswdProjectContent(&recent[i])
		}
	}
	q := c.Request.URL.Query().Encode()
	if q != "" {
		q = fmt.Sprintf("?%s", q)
	}
	ginH["posts"] = postIds
	ginH["totalPage"] = h.getTotalPage(totalRaw)
	ginH["currentPage"] = h.page
	ginH["title"] = h.getTitle()
	ginH["pagination"] = pagination.NewParsePagination(totalRaw, h.pageSize, h.page, q, c.Request.URL.Path, h.paginationStep)
}

func getTemplateName() string {
	tmlp := wpconfig.Options.Value("template")
	if i, err := theme.IsTemplateIsExist(tmlp); err != nil || !i {
		tmlp = "twentyfifteen"
	}
	return tmlp
}
