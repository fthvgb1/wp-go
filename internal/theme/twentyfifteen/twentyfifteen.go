package twentyfifteen

import (
	"github.com/fthvgb1/wp-go/helper/slice"
	"github.com/fthvgb1/wp-go/internal/pkg/cache"
	"github.com/fthvgb1/wp-go/internal/pkg/constraints"
	"github.com/fthvgb1/wp-go/internal/pkg/logs"
	"github.com/fthvgb1/wp-go/internal/pkg/models"
	"github.com/fthvgb1/wp-go/internal/theme/common"
	"github.com/gin-gonic/gin"
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
	i.Indexs()
}

func (d *detailHandle) Detail() {
	d.Details()
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
