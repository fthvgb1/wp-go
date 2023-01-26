package actions

import (
	"errors"
	"fmt"
	"github.com/fthvgb1/wp-go/helper/number"
	"github.com/fthvgb1/wp-go/helper/slice"
	str "github.com/fthvgb1/wp-go/helper/strings"
	"github.com/fthvgb1/wp-go/internal/pkg/cache"
	"github.com/fthvgb1/wp-go/internal/pkg/config"
	"github.com/fthvgb1/wp-go/internal/pkg/dao"
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
	"time"
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
	postStatus     []any
	header         string
	paginationStep int
	scene          int
	status         int
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
		orderBy:    "post_date",
		join:       model.SqlBuilder{},
		postType:   []any{"post"},
		postStatus: []any{"publish"},
		scene:      plugins.Home,
		status:     plugins.Ok,
	}
}

var months = slice.SimpleToMap(number.Range(1, 12, 1), func(v int) int {
	return v
})

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

var orders = []string{"asc", "desc"}

func (h *indexHandle) parseParams() (err error) {
	h.order = h.c.Query("order")

	if !slice.IsContained(h.order, orders) {
		order := config.Conf.Load().PostOrder
		h.order = "asc"
		if order != "" && slice.IsContained(order, orders) {
			h.order = order
		}
	}
	year := h.c.Param("year")
	if year != "" {
		y, er := strconv.Atoi(year)
		if er != nil {
			return err
		}
		if y > time.Now().Year() || y <= 1970 {
			return errors.New(str.Join("year err : ", year))
		}
		h.where = append(h.where, []string{
			"year(post_date)", year,
		})
	}
	month := h.c.Param("month")
	if month != "" {
		m, err := strconv.Atoi(month)
		if err != nil {
			return err
		}
		if _, ok := months[m]; !ok {
			return errors.New(str.Join("months err ", month))
		}

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
			h.scene = plugins.Tag
			allNames := cache.AllTagsNames(h.c)
			if _, ok := allNames[category]; !ok {
				return errors.New(str.Join("not exists tag ", category))
			}
			h.categoryType = "post_tag"
			h.header = fmt.Sprintf("标签： <span>%s</span>", category)
		}
	} else {
		h.scene = plugins.Category
		allNames := cache.AllCategoryNames(h.c)
		if _, ok := allNames[category]; !ok {
			return errors.New(str.Join("not exists category ", category))
		}
		h.categoryType = "category"
		h.header = fmt.Sprintf("分类： <span>%s</span>", category)
	}
	h.category = category
	username := h.c.Param("author")
	if username != "" {
		allUsername, er := cache.GetAllUsername(h.c)
		if er != nil {
			err = er
			return
		}
		if _, ok := allUsername[username]; !ok {
			err = errors.New(str.Join("user ", username, " is not exists"))
			return
		}
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
	}
	s := h.c.Query("s")
	if s != "" && strings.Replace(s, " ", "", -1) != "" {
		q := str.Join("%", s, "%")
		h.where = append(h.where, []string{
			"and", "post_title", "like", q, "",
			"or", "post_content", "like", q, "",
			"or", "post_excerpt", "like", q, "",
		}, []string{"post_password", ""})
		h.postType = append(h.postType, "page", "attachment")
		h.header = fmt.Sprintf("%s的搜索结果", s)
		h.setTitleLR(str.Join(`"`, s, `"`, "的搜索结果"), wpconfig.Options.Value("blogname"))
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
	var posts []models.Posts
	var totalRaw int
	var err error
	archive := cache.Archives(c)
	recent := cache.RecentPosts(c, 5)
	categoryItems := cache.Categories(c)
	recentComments := cache.RecentComments(c, 5)
	ginH := gin.H{
		"err":            err,
		"options":        wpconfig.Options,
		"recentPosts":    recent,
		"archives":       archive,
		"categories":     categoryItems,
		"search":         h.search,
		"header":         h.header,
		"recentComments": recentComments,
	}
	defer func() {
		code := http.StatusOK
		if err != nil {
			code = http.StatusNotFound
			if h.status == plugins.InternalErr {
				code = http.StatusInternalServerError
				c.Error(err)
				return
			}
			c.Error(err)
			h.status = plugins.Error
		}
		t := theme.GetTemplateName()
		theme.Hook(t, code, c, ginH, h.scene, h.status)
	}()
	err = h.parseParams()
	if err != nil {
		return
	}
	ginH["title"] = h.getTitle()
	if c.Param("month") != "" {
		posts, totalRaw, err = cache.GetMonthPostIds(c, c.Param("year"), c.Param("month"), h.page, h.pageSize, h.order)
		if err != nil {
			return
		}
	} else if h.search != "" {
		posts, totalRaw, err = cache.SearchPost(c, h.getSearchKey(), c, h.where, h.page, h.pageSize, model.SqlBuilder{{h.orderBy, h.order}}, h.join, h.postType, h.postStatus)
	} else {
		posts, totalRaw, err = cache.PostLists(c, h.getSearchKey(), c, h.where, h.page, h.pageSize, model.SqlBuilder{{h.orderBy, h.order}}, h.join, h.postType, h.postStatus)
	}
	if err != nil {
		h.status = plugins.Error
		logs.ErrPrintln(err, "获取数据错误")
		return
	}

	if len(posts) < 1 && h.category != "" {
		h.titleL = "未找到页面"
		h.status = plugins.Empty404
	}

	pw := h.session.Get("post_password")
	plug := plugins.NewPostPlugin(c, h.scene)
	for i, post := range posts {
		plugins.PasswordProjectTitle(&posts[i])
		if post.PostPassword != "" && pw != post.PostPassword {
			plugins.PasswdProjectContent(&posts[i])
		} else {
			plugins.ApplyPlugin(plug, &posts[i])
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
	ginH["posts"] = posts
	ginH["totalPage"] = h.getTotalPage(totalRaw)
	ginH["currentPage"] = h.page
	ginH["title"] = h.getTitle()
	ginH["scene"] = h.scene
	ginH["pagination"] = pagination.NewParsePagination(totalRaw, h.pageSize, h.page, h.paginationStep, q, c.Request.URL.Path)
}
