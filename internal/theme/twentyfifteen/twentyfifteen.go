package twentyfifteen

import (
	"github.com/fthvgb1/wp-go/internal/pkg/constraints"
	"github.com/fthvgb1/wp-go/internal/theme/common"
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
	h.CustomBackGround()
	h.Indexs()
}

func (h *handle) Detail() {
	h.CustomHeader()
	h.CustomBackGround()
	h.Details()
}
