package wp

import (
	"database/sql"
	"fmt"
	"github.com/fthvgb1/wp-go/helper/number"
	"github.com/fthvgb1/wp-go/helper/slice"
	"github.com/fthvgb1/wp-go/internal/pkg/cache"
	"github.com/fthvgb1/wp-go/internal/pkg/config"
	"github.com/fthvgb1/wp-go/internal/pkg/constraints"
	"github.com/fthvgb1/wp-go/internal/pkg/models"
	"github.com/fthvgb1/wp-go/internal/plugins"
	"github.com/fthvgb1/wp-go/model"
	"github.com/fthvgb1/wp-go/plugin/pagination"
)

type IndexHandle struct {
	*Handle
	Param        *IndexParams
	Posts        []models.Posts
	pageEle      pagination.Elements
	TotalRows    int
	postsPlugins map[string]Plugin[models.Posts, *Handle]
}

func (i *IndexHandle) PageEle() pagination.Elements {
	return i.pageEle
}

func (i *IndexHandle) SetPageEle(pageEle pagination.Elements) {
	i.pageEle = pageEle
}

func (i *IndexHandle) PostsPlugins() map[string]Plugin[models.Posts, *Handle] {
	return i.postsPlugins
}

func (i *IndexHandle) SetPostsPlugins(postsPlugins map[string]Plugin[models.Posts, *Handle]) {
	i.postsPlugins = postsPlugins
}

func NewIndexHandle(handle *Handle) *IndexHandle {
	return &IndexHandle{Handle: handle}
}

func (i *IndexHandle) ParseIndex(parm *IndexParams) (err error) {
	i.Param = parm
	switch i.scene {
	case constraints.Home, constraints.Search:
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

	case constraints.Search:

		posts, totalRaw, err = cache.SearchPost(i.C, i.Param.CacheKey, i.C, q, i.Param.Page, i.Param.PageSize)

	case constraints.Archive:

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

func (i *IndexHandle) ExecPostsPlugin(calls ...func(*models.Posts)) {

	pluginConf := config.GetConfig().ListPagePlugins

	postsPlugins := i.postsPlugins
	if postsPlugins == nil {
		postsPlugins = pluginFns
	}
	plugin := GetListPostPlugins(pluginConf, postsPlugins)

	i.Posts = slice.Map(i.Posts, PluginFn[models.Posts](plugin, i.Handle, Defaults(calls...)))

}

func (i *IndexHandle) Render() {
	i.PushHandleFn(constraints.Ok, NewHandleFn(func(h *Handle) {
		i.ExecPostsPlugin()
		i.Pagination()
	}, 10))
	i.ginH["posts"] = i.Posts
	i.Handle.Render()
}

func (i *IndexHandle) Indexs() {
	_ = i.BuildIndexData(NewIndexParams(i.C))
	i.Render()
}
