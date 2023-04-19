package wp

import (
	"database/sql"
	"fmt"
	"github.com/fthvgb1/wp-go/helper/number"
	"github.com/fthvgb1/wp-go/helper/slice"
	"github.com/fthvgb1/wp-go/internal/pkg/cache"
	"github.com/fthvgb1/wp-go/internal/pkg/constraints"
	"github.com/fthvgb1/wp-go/internal/pkg/models"
	"github.com/fthvgb1/wp-go/internal/plugins"
	"github.com/fthvgb1/wp-go/model"
	"github.com/fthvgb1/wp-go/plugin/pagination"
	"strings"
)

type IndexHandle struct {
	*Handle
	Param       *IndexParams
	Posts       []models.Posts
	pageEle     pagination.Elements
	TotalRows   int
	postsPlugin PostsPlugin
}

func (i *IndexHandle) ListPlugin() func(*Handle, *models.Posts) {
	return i.postsPlugin
}

func (i *IndexHandle) SetListPlugin(listPlugin func(*Handle, *models.Posts)) {
	i.postsPlugin = listPlugin
}

func (i *IndexHandle) PageEle() pagination.Elements {
	return i.pageEle
}

func (i *IndexHandle) SetPageEle(pageEle pagination.Elements) {
	i.pageEle = pageEle
}

func NewIndexHandle(handle *Handle) *IndexHandle {
	return &IndexHandle{Handle: handle}
}

func (i *IndexHandle) ParseIndex(parm *IndexParams) (err error) {
	i.Param = parm
	switch i.scene {
	case constraints.Search:
		i.Param.ParseSearch()
	case constraints.Category:
		err = i.Param.ParseCategory()
	case constraints.Tag:
		err = i.Param.ParseTag()
	case constraints.Archive:
		err = i.Param.ParseArchive()
	case constraints.Author:
		err = i.Param.ParseAuthor()
	}
	if err != nil {
		i.Stats = constraints.ParamError
		return
	}
	i.Param.ParseParams()
	i.Param.CacheKey = i.Param.getSearchKey()
	i.ginH["title"] = i.Param.getTitle()
	i.ginH["search"] = i.Param.Search
	i.ginH["header"] = i.Param.Header
	return
}

func (i *IndexHandle) GetIndexData() (posts []models.Posts, totalRaw int, err error) {

	q := model.QueryCondition{
		Where: i.Param.Where,
		Order: model.SqlBuilder{{i.Param.OrderBy, i.Param.Order}},
		Join:  i.Param.Join,
		In:    [][]any{i.Param.PostType, i.Param.PostStatus},
	}
	switch i.scene {
	case constraints.Home, constraints.Category, constraints.Tag, constraints.Author:

		posts, totalRaw, err = cache.PostLists(i.C, i.Param.CacheKey, i.C, q, i.Param.Page, i.Param.PageSize)
		if i.scene == constraints.Home && i.Param.Page == 1 {
			i.MarkSticky(&posts)
		}

	case constraints.Search:

		posts, totalRaw, err = cache.SearchPost(i.C, i.Param.CacheKey, i.C, q, i.Param.Page, i.Param.PageSize)

	case constraints.Archive:
		i.ginH["archiveYear"] = i.Param.Year
		i.ginH["archiveMonth"] = strings.TrimLeft(i.Param.Month, "0")
		posts, totalRaw, err = cache.GetMonthPostIds(i.C, i.Param.Year, i.Param.Month, i.Param.Page, i.Param.PageSize, i.Param.Order)

	}
	return

}

func (i *IndexHandle) Pagination() {
	if i.pageEle == nil {
		i.pageEle = plugins.TwentyFifteenPagination()
	}
	q := i.C.Request.URL.Query().Encode()
	if q != "" {
		q = fmt.Sprintf("?%s", q)
	}
	paginations := pagination.NewParsePagination(i.TotalRows, i.Param.PageSize, i.Param.Page, i.Param.PaginationStep, q, i.C.Request.URL.Path)
	i.ginH["pagination"] = pagination.Paginate(i.pageEle, paginations)

}

func (i *IndexHandle) BuildIndexData(parm *IndexParams) (err error) {
	err = i.ParseIndex(parm)
	if err != nil {
		i.Stats = constraints.ParamError
		return
	}
	posts, totalRows, err := i.GetIndexData()
	if err != nil && err != sql.ErrNoRows {
		i.Stats = constraints.Error404
		return
	}
	i.Posts = posts
	i.TotalRows = totalRows
	i.ginH["totalPage"] = number.CalTotalPage(totalRows, i.Param.PageSize)
	return
}

func (i *IndexHandle) ExecPostsPlugin() {
	if i.postsPlugin != nil {
		for j := range i.Posts {
			i.postsPlugin(i.Handle, &i.Posts[j])
		}
	}
}

func IndexRender(h *Handle) {
	if h.scene == constraints.Detail || h.Stats != constraints.Ok {
		return
	}
	i := h.Index
	i.ExecPostsPlugin()
	i.Pagination()
	i.ginH["posts"] = i.Posts
}

func Indexs(h *Handle) {
	if h.Scene() == constraints.Detail {
		return
	}
	i := h.Index
	_ = i.BuildIndexData(NewIndexParams(i.C))
}

func (i *IndexHandle) MarkSticky(posts *[]models.Posts) {
	a := i.StickPosts()
	if len(a) < 1 {
		return
	}
	m := i.StickMapPosts()
	*posts = append(a, slice.Filter(*posts, func(post models.Posts, _ int) bool {
		_, ok := m[post.Id]
		return !ok
	})...)
}
