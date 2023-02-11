package twentyfifteen

import (
	"github.com/fthvgb1/wp-go/internal/pkg/constraints"
	"github.com/fthvgb1/wp-go/internal/plugins"
	"github.com/fthvgb1/wp-go/internal/theme/common"
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

	err := i.BuildIndexData(common.NewIndexParams(i.C))
	if err != nil {
		i.C.HTML(i.Code, i.Templ, i.GinH)
		return
	}
	i.ExecPostsPlugin()
	i.PageEle = plugins.TwentyFifteenPagination()
	i.Pagination()
	i.C.HTML(i.Code, i.Templ, i.GinH)
}

func (d *detailHandle) Detail() {
	d.Templ = "twentyfifteen/posts/detail.gohtml"

	err := d.BuildDetailData()
	if err != nil {
		d.Stats = constraints.Error404
		d.C.HTML(d.Code, d.Templ, d.GinH)
		return
	}
	d.PasswordProject()
	d.CommentRender = plugins.CommentRender()
	d.RenderComment()
	d.C.HTML(d.Code, d.Templ, d.GinH)
}
