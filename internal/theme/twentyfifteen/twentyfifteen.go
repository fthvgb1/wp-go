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

type handle struct {
	*common.IndexHandle
	*common.DetailHandle
}

func newHandle(iHandle *common.IndexHandle, dHandle *common.DetailHandle) *handle {
	return &handle{iHandle, dHandle}
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
	handle := newHandle(common.NewIndexHandle(h), common.NewDetailHandle(h))
	if h.Scene == constraints.Detail {
		handle.Detail()
		return
	}
	handle.Index()
}

func (h *handle) Index() {
	h.CustomHeader()
	h.Indexs()
}

func (h *handle) Detail() {
	h.CustomHeader()
	h.Details()
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
