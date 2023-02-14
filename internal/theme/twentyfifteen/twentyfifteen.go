package twentyfifteen

import (
	"fmt"
	"github.com/fthvgb1/wp-go/helper/slice"
	"github.com/fthvgb1/wp-go/internal/pkg/cache"
	"github.com/fthvgb1/wp-go/internal/pkg/constraints"
	"github.com/fthvgb1/wp-go/internal/pkg/logs"
	"github.com/fthvgb1/wp-go/internal/pkg/models"
	"github.com/fthvgb1/wp-go/internal/plugins"
	"github.com/fthvgb1/wp-go/internal/theme/common"
	"github.com/gin-gonic/gin"
	"net/http"
)

const ThemeName = "twentyfifteen"

type indexHandle struct {
	*common.IndexHandle
}

func newIndexHandle(iHandle *common.IndexHandle) *indexHandle {
	return &indexHandle{iHandle}
}

type detailHandle struct {
	*common.DetailHandle
}

func newDetailHandle(dHandle *common.DetailHandle) *detailHandle {
	return &detailHandle{DetailHandle: dHandle}
}

func Hook(h *common.Handle) {
	h.WidgetAreaData()
	h.GetPassword()
	if h.Scene == constraints.Detail {
		newDetailHandle(common.NewDetailHandle(h)).Detail()
		return
	}
	newIndexHandle(common.NewIndexHandle(h)).Index()
}

func (i *indexHandle) Index() {
	i.Templ = "twentyfifteen/posts/index.gohtml"
	img := getHeaderImage(i.C)
	fmt.Println(img)
	err := i.BuildIndexData(common.NewIndexParams(i.C))
	if err != nil {
		i.Stats = constraints.Error404
		i.Code = http.StatusNotFound
		i.C.HTML(i.Code, i.Templ, i.GinH)
		return
	}
	i.ExecPostsPlugin()
	i.PageEle = plugins.TwentyFifteenPagination()
	i.Pagination()
	i.CalBodyClass()
	i.C.HTML(i.Code, i.Templ, i.GinH)
}

func (d *detailHandle) Detail() {
	d.Templ = "twentyfifteen/posts/detail.gohtml"

	err := d.BuildDetailData()
	if err != nil {
		d.Stats = constraints.Error404
		d.Code = http.StatusNotFound
		d.C.HTML(d.Code, d.Templ, d.GinH)
		return
	}
	d.PasswordProject()
	d.CommentRender = plugins.CommentRender()
	d.RenderComment()
	d.CalBodyClass()
	d.C.HTML(d.Code, d.Templ, d.GinH)
}

func getHeaderImage(c *gin.Context) (r models.PostThumbnail) {
	r.Path = "/wp-content/themes/twentyseventeen/assets/images/header.jpg"
	r.Width = 2000
	r.Height = 1200
	hs, err := cache.GetHeaderImages(c, ThemeName)
	if err != nil {
		logs.ErrPrintln(err, "获取页眉背景图失败")
	} else if len(hs) > 0 && err == nil {
		_, r = slice.Rand(hs)

	}
	r.Sizes = "100vw"
	return
}

func ThemeSupport() map[string]struct{} {
	return map[string]struct{}{
		"custom-background": {},
		"wp-custom-logo":    {},
		"responsive-embeds": {},
		"post-formats":      {},
	}
}
