package common

import (
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
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

type IndexParams struct {
	ParseSearch       func()
	ParseArchive      func() error
	ParseCategory     func() error
	ParseTag          func() error
	ParseAuthor       func() error
	CategoryCondition func()
	ParseParams       func()
	Ctx               *gin.Context
	Page              int
	PageSize          int
	Title             string
	TitleL            string
	TitleR            string
	Search            string
	Author            string
	TotalPage         int
	Category          string
	CategoryType      string
	Where             model.SqlBuilder
	OrderBy           string
	Order             string
	Month             string
	Year              string
	Join              model.SqlBuilder
	PostType          []any
	PostStatus        []any
	Header            string
	PaginationStep    int
	CacheKey          string
	BlogName          string
}

type IndexHandle struct {
	*Handle
	Param        *IndexParams
	Posts        []models.Posts
	PageEle      pagination.Elements
	TotalRows    int
	PostsPlugins map[string]Plugin[models.Posts]
}

func NewIndexHandle(handle *Handle) *IndexHandle {
	return &IndexHandle{Handle: handle}
}

var months = slice.SimpleToMap(number.Range(1, 12, 1), func(v int) int {
	return v
})

var orders = map[string]struct{}{"asc": {}, "desc": {}}

func (i *IndexParams) setTitleLR(l, r string) {
	i.TitleL = l
	i.TitleR = r
}

func (i *IndexParams) getTitle() string {
	i.Title = fmt.Sprintf("%s-%s", i.TitleL, i.TitleR)
	return i.Title
}

func (i *IndexParams) getSearchKey() string {
	return fmt.Sprintf("action:%s|%s|%s|%s|%s|%s|%d|%d", i.Author, i.Search, i.OrderBy, i.Order, i.Category, i.CategoryType, i.Page, i.PageSize)
}

func NewIndexParams(ctx *gin.Context) *IndexParams {
	blogName := wpconfig.Options.Value("blogname")
	size := str.ToInteger(wpconfig.Options.Value("posts_per_page"), 10)
	i := &IndexParams{
		Ctx:            ctx,
		Page:           1,
		PageSize:       size,
		PaginationStep: number.Max(1, config.GetConfig().PaginationStep),
		TitleL:         blogName,
		TitleR:         wpconfig.Options.Value("blogdescription"),
		Where: model.SqlBuilder{
			{"post_type", "in", ""},
			{"post_status", "in", ""},
		},
		OrderBy:    "post_date",
		Join:       model.SqlBuilder{},
		PostType:   []any{"post"},
		PostStatus: []any{"publish"},
		BlogName:   wpconfig.Options.Value("blogname"),
	}
	i.ParseSearch = i.parseSearch
	i.ParseArchive = i.parseArchive
	i.ParseCategory = i.parseCategory
	i.ParseTag = i.parseTag
	i.CategoryCondition = i.categoryCondition
	i.ParseAuthor = i.parseAuthor
	i.ParseParams = i.parseParams
	return i
}

func (i *IndexParams) parseSearch() {
	s := i.Ctx.Query("s")
	if s != "" && strings.Replace(s, " ", "", -1) != "" {
		q := str.Join("%", s, "%")
		i.Where = append(i.Where, []string{
			"and", "post_title", "like", q, "",
			"or", "post_content", "like", q, "",
			"or", "post_excerpt", "like", q, "",
		}, []string{"post_password", ""})
		i.PostType = append(i.PostType, "Page", "attachment")
		i.Header = fmt.Sprintf("<span>%s</span>的搜索结果", s)
		i.setTitleLR(str.Join(`"`, s, `"`, "的搜索结果"), i.BlogName)
		i.Search = s
	}
}
func (i *IndexParams) parseArchive() error {
	year := i.Ctx.Param("year")
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
	month := i.Ctx.Param("month")
	if month != "" {
		m := str.ToInteger(month, -1)
		if !maps.IsExists(months, m) {
			return errors.New(str.Join("months err ", month))
		}

		i.Where = append(i.Where, []string{
			"month(post_date)", month,
		})
		ss := fmt.Sprintf("%s年%s月", year, strings.TrimLeft(month, "0"))
		i.Header = fmt.Sprintf("月度归档： <span>%s</span>", ss)
		i.setTitleLR(ss, i.BlogName)
		i.Month = month
	}
	return nil
}

func (i *IndexParams) parseCategory() error {
	category := i.Ctx.Param("category")
	if category != "" {
		if !maps.IsExists(cache.AllCategoryTagsNames(i.Ctx, constraints.Category), category) {
			return errors.New(str.Join("not exists category ", category))
		}
		i.CategoryType = "category"
		i.Header = fmt.Sprintf("分类： <span>%s</span>", category)
		i.Category = category
		i.CategoryCondition()
	}
	return nil
}
func (i *IndexParams) parseTag() error {
	tag := i.Ctx.Param("tag")
	if tag != "" {
		if !maps.IsExists(cache.AllCategoryTagsNames(i.Ctx, constraints.Tag), tag) {
			return errors.New(str.Join("not exists tag ", tag))
		}
		i.CategoryType = "post_tag"
		i.Header = fmt.Sprintf("标签： <span>%s</span>", tag)
		i.Category = tag
		i.CategoryCondition()
	}
	return nil
}

func (i *IndexParams) categoryCondition() {
	if i.Category != "" {
		i.Where = append(i.Where, []string{
			"d.name", i.Category,
		}, []string{"taxonomy", i.CategoryType})
		i.Join = append(i.Join, []string{
			"a", "left Join", "wp_term_relationships b", "a.Id=b.object_id",
		}, []string{
			"left Join", "wp_term_taxonomy c", "b.term_taxonomy_id=c.term_taxonomy_id",
		}, []string{
			"left Join", "wp_terms d", "c.term_id=d.term_id",
		})
		i.setTitleLR(i.Category, i.BlogName)
	}
}
func (i *IndexParams) parseAuthor() (err error) {
	username := i.Ctx.Param("author")
	if username != "" {
		allUsername, er := cache.GetAllUsername(i.Ctx)
		if err != nil {
			err = er
			return
		}
		if !maps.IsExists(allUsername, username) {
			err = errors.New(str.Join("user ", username, " is not exists"))
			return
		}
		user, er := cache.GetUserByName(i.Ctx, username)
		if er != nil {
			return
		}
		i.Author = username
		i.Where = append(i.Where, []string{
			"post_author", "=", strconv.FormatUint(user.Id, 10), "int",
		})
	}
	return
}

func (i *IndexParams) parseParams() {
	i.Order = i.Ctx.Query("Order")
	if !maps.IsExists(orders, i.Order) {
		order := config.GetConfig().PostOrder
		i.Order = "asc"
		if order != "" && maps.IsExists(orders, order) {
			i.Order = order
		}
	}

	i.Page = str.ToInteger(i.Ctx.Param("page"), 1)
	total := int(atomic.LoadInt64(&dao.TotalRaw))
	if total > 0 && total < (i.Page-1)*i.PageSize {
		i.Page = 1
	}
	if i.Page > 1 && (i.Category != "" || i.Search != "" || i.Month != "") {
		i.setTitleLR(fmt.Sprintf("%s-第%d页", i.TitleL, i.Page), i.BlogName)
	}
	return
}
