package wp

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/fthvgb1/wp-go/app/pkg/cache"
	"github.com/fthvgb1/wp-go/app/pkg/constraints"
	"github.com/fthvgb1/wp-go/app/pkg/models"
	"github.com/fthvgb1/wp-go/app/plugins"
	"github.com/fthvgb1/wp-go/cache/reload"
	"github.com/fthvgb1/wp-go/helper/number"
	"github.com/fthvgb1/wp-go/helper/slice"
	"github.com/fthvgb1/wp-go/model"
	"github.com/fthvgb1/wp-go/plugin/pagination"
	"strings"
)

type IndexHandle struct {
	*Handle
	Param       *IndexParams
	Posts       []models.Posts
	pageEle     pagination.Render
	TotalRows   int
	postsPlugin PostsPlugin
}

func (h *Handle) GetIndexHandle() *IndexHandle {
	v, ok := h.C.Get("indexHandle")
	if !ok {
		vv := NewIndexHandle(h)
		h.C.Set("indexHandle", vv)
		return vv
	}
	return v.(*IndexHandle)
}

func (i *IndexHandle) ListPlugin() func(*Handle, *models.Posts) {
	return i.postsPlugin
}

func (i *IndexHandle) SetListPlugin(listPlugin func(*Handle, *models.Posts)) {
	i.postsPlugin = listPlugin
}

func (i *IndexHandle) PageEle() pagination.Render {
	return i.pageEle
}

func (i *IndexHandle) SetPageEle(pageEle pagination.Render) {
	i.pageEle = pageEle
}

func NewIndexHandle(handle *Handle) *IndexHandle {
	return &IndexHandle{Handle: handle}
}

func PushIndexHandler(pipeScene string, h *Handle, call HandleCall) {
	h.PushHandlers(pipeScene, call, constraints.Home,
		constraints.Category, constraints.Search, constraints.Tag,
		constraints.Archive, constraints.Author,
	)
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

	q := &model.QueryCondition{
		Where: i.Param.Where,
		Order: model.SqlBuilder{{i.Param.OrderBy, i.Param.Order}},
		Join:  i.Param.Join,
		In:    [][]any{i.Param.PostType, i.Param.PostStatus},
	}
	switch i.scene {
	case constraints.Home, constraints.Category, constraints.Tag, constraints.Author:

		posts, totalRaw, err = cache.PostLists(i.C, i.Param.CacheKey, q, i.Param.Page, i.Param.PageSize)
		if i.scene == constraints.Home && i.Param.Page == 1 {
			i.MarkSticky(&posts)
		}

	case constraints.Search:

		posts, totalRaw, err = cache.SearchPost(i.C, i.Param.CacheKey, q, i.Param.Page, i.Param.PageSize)

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
	i.ginH["pagination"] = pagination.Paginate(i.pageEle, i.TotalRows, i.Param.PageSize, i.Param.Page, i.Param.PaginationStep, *i.C.Request.URL, i.IsHttps())

}

func (i *IndexHandle) BuildIndexData() (err error) {
	if i.Param == nil {
		i.Param = NewIndexParams(i.C)
	}
	err = i.ParseIndex(i.Param)
	if err != nil {
		i.Stats = constraints.ParamError
		return
	}
	posts, totalRows, err := i.GetIndexData()
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		i.Stats = constraints.Error404
		return
	}
	i.Posts = posts
	i.TotalRows = totalRows
	i.ginH["totalPage"] = number.DivideCeil(totalRows, i.Param.PageSize)
	return
}

var GetPostsPlugin = reload.BuildValFnWithAnyParams("postPlugins", UsePostsPlugins)

func (i *IndexHandle) ExecPostsPlugin() {
	fn := i.postsPlugin
	if fn == nil {
		fn = GetPostsPlugin()
	}
	for j := range i.Posts {
		fn(i.Handle, &i.Posts[j])
	}
}

func IndexRender(h *Handle) {
	i := h.GetIndexHandle()
	i.ExecPostsPlugin()
	i.Pagination()
	i.ginH["posts"] = i.Posts
}

func Index(h *Handle) {
	i := h.GetIndexHandle()
	err := i.BuildIndexData()
	if err != nil {
		i.SetErr(err, High)
	}
	h.SetData("scene", h.Scene())
}

func (i *IndexHandle) MarkSticky(posts *[]models.Posts) {
	a := GetStickPosts(i.Handle)
	if len(a) < 1 {
		return
	}
	m := GetStickMapPosts(i.Handle)
	*posts = append(a, slice.Filter(*posts, func(post models.Posts, _ int) bool {
		_, ok := m[post.Id]
		return !ok
	})...)
}
