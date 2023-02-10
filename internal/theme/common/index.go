package common

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/fthvgb1/wp-go/helper/maps"
	"github.com/fthvgb1/wp-go/helper/number"
	"github.com/fthvgb1/wp-go/helper/slice"
	str "github.com/fthvgb1/wp-go/helper/strings"
	"github.com/fthvgb1/wp-go/internal/pkg/cache"
	"github.com/fthvgb1/wp-go/internal/pkg/config"
	"github.com/fthvgb1/wp-go/internal/pkg/constraints"
	"github.com/fthvgb1/wp-go/internal/pkg/dao"
	"github.com/fthvgb1/wp-go/internal/pkg/models"
	"github.com/fthvgb1/wp-go/internal/wpconfig"
	"github.com/fthvgb1/wp-go/model"
	"github.com/fthvgb1/wp-go/plugin/pagination"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

type IndexParams struct {
	c              *gin.Context
	Page           int
	PageSize       int
	Title          string
	titleL         string
	titleR         string
	search         string
	author         string
	totalPage      int
	category       string
	categoryType   string
	Where          model.SqlBuilder
	OrderBy        string
	Order          string
	Month          string
	Year           string
	Join           model.SqlBuilder
	PostType       []any
	PostStatus     []any
	header         string
	PaginationStep int
	scene          int
	stats          int
	CacheKey       string
	blogName       string
}

var months = slice.SimpleToMap(number.Range(1, 12, 1), func(v int) int {
	return v
})

var orders = map[string]struct{}{"asc": {}, "desc": {}}

func (i *IndexParams) setTitleLR(l, r string) {
	i.titleL = l
	i.titleR = r
}

func (i *IndexParams) getTitle() string {
	i.Title = fmt.Sprintf("%s-%s", i.titleL, i.titleR)
	return i.Title
}

func (i *IndexParams) getSearchKey() string {
	return fmt.Sprintf("action:%s|%s|%s|%s|%s|%s|%d|%d", i.author, i.search, i.OrderBy, i.Order, i.category, i.categoryType, i.Page, i.PageSize)
}

func newIndexHandle(ctx *gin.Context) *IndexParams {
	blogName := wpconfig.Options.Value("blogname")
	size := str.ToInteger(wpconfig.Options.Value("posts_per_page"), 10)
	return &IndexParams{
		c:              ctx,
		Page:           1,
		PageSize:       size,
		PaginationStep: number.Max(1, config.GetConfig().PaginationStep),
		titleL:         blogName,
		titleR:         wpconfig.Options.Value("blogdescription"),
		Where: model.SqlBuilder{
			{"post_type", "in", ""},
			{"post_status", "in", ""},
		},
		OrderBy:    "post_date",
		Join:       model.SqlBuilder{},
		PostType:   []any{"post"},
		PostStatus: []any{"publish"},
		scene:      constraints.Home,
		stats:      constraints.Ok,
		blogName:   wpconfig.Options.Value("blogname"),
	}
}

func (i *IndexParams) ParseSearch() {
	s := i.c.Query("s")
	if s != "" && strings.Replace(s, " ", "", -1) != "" {
		q := str.Join("%", s, "%")
		i.Where = append(i.Where, []string{
			"and", "post_title", "like", q, "",
			"or", "post_content", "like", q, "",
			"or", "post_excerpt", "like", q, "",
		}, []string{"post_password", ""})
		i.PostType = append(i.PostType, "Page", "attachment")
		i.header = fmt.Sprintf("<span>%s</span>的搜索结果", s)
		i.setTitleLR(str.Join(`"`, s, `"`, "的搜索结果"), i.blogName)
		i.search = s
		i.scene = constraints.Search
	}
}
func (i *IndexParams) ParseArchive() error {
	year := i.c.Param("year")
	if year != "" {
		y := str.ToInteger(year, -1)
		if y > time.Now().Year() || y <= 1970 {
			return errors.New(str.Join("year err : ", year))
		}
		i.Where = append(i.Where, []string{
			"year(post_date)", year,
		})
		i.Year = year
	}
	month := i.c.Param("month")
	if month != "" {
		m := str.ToInteger(month, -1)
		if !maps.IsExists(months, m) {
			return errors.New(str.Join("months err ", month))
		}

		i.Where = append(i.Where, []string{
			"month(post_date)", month,
		})
		ss := fmt.Sprintf("%s年%s月", year, strings.TrimLeft(month, "0"))
		i.header = fmt.Sprintf("月度归档： <span>%s</span>", ss)
		i.setTitleLR(ss, i.blogName)
		i.scene = constraints.Archive
		i.Month = month
	}
	return nil
}

func (i *IndexParams) ParseCategory() error {
	category := i.c.Param("category")
	if category != "" {
		i.scene = constraints.Category
		if !maps.IsExists(cache.AllCategoryTagsNames(i.c, constraints.Category), category) {
			return errors.New(str.Join("not exists category ", category))
		}
		i.categoryType = "category"
		i.header = fmt.Sprintf("分类： <span>%s</span>", category)
		i.category = category
		i.CategoryCondition()
	}
	return nil
}
func (i *IndexParams) ParseTag() error {
	tag := i.c.Param("tag")
	if tag != "" {
		i.scene = constraints.Tag
		if !maps.IsExists(cache.AllCategoryTagsNames(i.c, constraints.Tag), tag) {
			return errors.New(str.Join("not exists tag ", tag))
		}
		i.categoryType = "post_tag"
		i.header = fmt.Sprintf("标签： <span>%s</span>", tag)
		i.category = tag
		i.CategoryCondition()
	}
	return nil
}

func (i *IndexParams) CategoryCondition() {
	if i.category != "" {
		i.Where = append(i.Where, []string{
			"d.name", i.category,
		}, []string{"taxonomy", i.categoryType})
		i.Join = append(i.Join, []string{
			"a", "left Join", "wp_term_relationships b", "a.Id=b.object_id",
		}, []string{
			"left Join", "wp_term_taxonomy c", "b.term_taxonomy_id=c.term_taxonomy_id",
		}, []string{
			"left Join", "wp_terms d", "c.term_id=d.term_id",
		})
		i.setTitleLR(i.category, i.blogName)
	}
}
func (i *IndexParams) ParseAuthor() (err error) {
	username := i.c.Param("author")
	if username != "" {
		allUsername, er := cache.GetAllUsername(i.c)
		if err != nil {
			err = er
			return
		}
		if !maps.IsExists(allUsername, username) {
			err = errors.New(str.Join("user ", username, " is not exists"))
			return
		}
		user, er := cache.GetUserByName(i.c, username)
		if er != nil {
			return
		}
		i.author = username
		i.Where = append(i.Where, []string{
			"post_author", "=", strconv.FormatUint(user.Id, 10), "int",
		})
	}
	return
}

func (i *IndexParams) parseParams() {
	i.Order = i.c.Query("Order")
	if !maps.IsExists(orders, i.Order) {
		order := config.GetConfig().PostOrder
		i.Order = "asc"
		if order != "" && maps.IsExists(orders, order) {
			i.Order = order
		}
	}

	i.Page = str.ToInteger(i.c.Param("page"), 1)
	total := int(atomic.LoadInt64(&dao.TotalRaw))
	if total > 0 && total < (i.Page-1)*i.PageSize {
		i.Page = 1
	}
	if i.Page > 1 && (i.category != "" || i.search != "" || i.Month != "") {
		i.setTitleLR(fmt.Sprintf("%s-第%d页", i.titleL, i.Page), i.blogName)
	}
	return
}

func (h *Handle) ParseIndex() (i *IndexParams, err error) {
	i = newIndexHandle(h.C)
	switch h.Scene {
	case constraints.Home, constraints.Search:
		i.ParseSearch()
	case constraints.Category:
		err = i.ParseCategory()
	case constraints.Tag:
		err = i.ParseTag()
	case constraints.Archive:
		err = i.ParseArchive()
	case constraints.Author:
		err = i.ParseAuthor()
	}
	if err != nil {
		h.Stats = constraints.ParamError
		return
	}
	i.CacheKey = i.getSearchKey()
	i.parseParams()
	h.GinH["title"] = i.getTitle()
	h.GinH["search"] = i.search
	h.GinH["header"] = i.header
	return
}

func (h *Handle) GetIndexData(i *IndexParams) (posts []models.Posts, totalRaw int, err error) {

	switch h.Scene {
	case constraints.Home, constraints.Category, constraints.Tag, constraints.Author:

		posts, totalRaw, err = cache.PostLists(h.C, i.CacheKey, h.C, i.Where, i.Page, i.PageSize,
			model.SqlBuilder{{i.OrderBy, i.Order}}, i.Join, i.PostType, i.PostStatus)

	case constraints.Search:

		posts, totalRaw, err = cache.SearchPost(h.C, i.CacheKey, h.C, i.Where, i.Page, i.PageSize,
			model.SqlBuilder{{i.OrderBy, i.Order}}, i.Join, i.PostType, i.PostStatus)

	case constraints.Archive:

		posts, totalRaw, err = cache.GetMonthPostIds(h.C, i.Year, i.Month, i.Page, i.PageSize, i.Order)

	}
	return

}

func (h *Handle) Indexs() (err error) {
	i, err := h.ParseIndex()
	if err != nil {
		h.Stats = constraints.ParamError
		h.Code = http.StatusNotFound
		return
	}
	posts, totalRows, err := h.GetIndexData(i)
	if err != nil && err != sql.ErrNoRows {
		h.Scene = constraints.InternalErr
		h.Code = http.StatusInternalServerError
		return
	}
	pw := h.Session.Get("post_password")
	if pw != nil {
		h.Password = pw.(string)
	}
	h.GinH["posts"] = posts
	h.GinH["totalPage"] = number.CalTotalPage(totalRows, i.PageSize)
	q := h.C.Request.URL.Query().Encode()
	if q != "" {
		q = fmt.Sprintf("?%s", q)
	}
	h.GinH["pagination"] = pagination.NewParsePagination(totalRows, i.PageSize, i.Page, i.PaginationStep, q, h.C.Request.URL.Path)
	return
}
