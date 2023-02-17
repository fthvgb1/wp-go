package common

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
	"net/http"
)

func (i *IndexHandle) ParseIndex(parm *IndexParams) (err error) {
	i.Param = parm
	switch i.Scene {
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
	i.GinH["title"] = i.Param.getTitle()
	i.GinH["search"] = i.Param.Search
	i.GinH["header"] = i.Param.Header
	return
}

func (i *IndexHandle) GetIndexData() (posts []models.Posts, totalRaw int, err error) {

	switch i.Scene {
	case constraints.Home, constraints.Category, constraints.Tag, constraints.Author:

		posts, totalRaw, err = cache.PostLists(i.C, i.Param.CacheKey, i.C, i.Param.Where, i.Param.Page, i.Param.PageSize,
			model.SqlBuilder{{i.Param.OrderBy, i.Param.Order}}, i.Param.Join, i.Param.PostType, i.Param.PostStatus)

	case constraints.Search:

		posts, totalRaw, err = cache.SearchPost(i.C, i.Param.CacheKey, i.C, i.Param.Where, i.Param.Page, i.Param.PageSize,
			model.SqlBuilder{{i.Param.OrderBy, i.Param.Order}}, i.Param.Join, i.Param.PostType, i.Param.PostStatus)

	case constraints.Archive:

		posts, totalRaw, err = cache.GetMonthPostIds(i.C, i.Param.Year, i.Param.Month, i.Param.Page, i.Param.PageSize, i.Param.Order)

	}
	return

}

func (i *IndexHandle) Pagination() {
	q := i.C.Request.URL.Query().Encode()
	if q != "" {
		q = fmt.Sprintf("?%s", q)
	}
	paginations := pagination.NewParsePagination(i.TotalRows, i.Param.PageSize, i.Param.Page, i.Param.PaginationStep, q, i.C.Request.URL.Path)

	i.GinH["pagination"] = pagination.Paginate(i.PageEle, paginations)

}

func (i *IndexHandle) BuildIndexData(parm *IndexParams) (err error) {
	err = i.ParseIndex(parm)
	if err != nil {
		return
	}
	posts, totalRows, err := i.GetIndexData()
	if err != nil && err != sql.ErrNoRows {
		return
	}
	i.GinH["posts"] = posts
	i.Posts = posts
	i.TotalRows = totalRows

	i.GinH["totalPage"] = number.CalTotalPage(totalRows, i.Param.PageSize)

	return
}

func (i *IndexHandle) ExecPostsPlugin(calls ...func(*models.Posts)) {

	pluginConf := config.GetConfig().ListPagePlugins

	postsPlugins := i.PostsPlugins
	if postsPlugins == nil {
		postsPlugins = pluginFns
	}
	plugin := GetListPostPlugins(pluginConf, postsPlugins)

	i.GinH["posts"] = slice.Map(i.Posts, PluginFn[models.Posts](plugin, i.Handle, Defaults(calls...)))

}

func (i *IndexHandle) Render() {
	i.ExecPostsPlugin()
	if i.PageEle == nil {
		i.PageEle = plugins.TwentyFifteenPagination()
	}
	i.Pagination()
	i.SiteIcon()
	i.CustomLogo()
	i.CustomCss()
	i.CalBodyClass()
	if i.Templ == "" {
		i.Templ = fmt.Sprintf("%s/posts/index.gohtml", i.Theme)
	}
	i.C.HTML(i.Code, i.Templ, i.GinH)
}

func (i *IndexHandle) Indexs() {
	err := i.BuildIndexData(NewIndexParams(i.C))
	if err != nil {
		i.Stats = constraints.Error404
		i.Code = http.StatusNotFound
		i.C.HTML(i.Code, i.Templ, i.GinH)
		return
	}
	i.Render()
}
