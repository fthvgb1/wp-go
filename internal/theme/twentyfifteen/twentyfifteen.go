package twentyfifteen

import (
	"embed"
	"encoding/json"
	"github.com/fthvgb1/wp-go/internal/pkg/constraints"
	"github.com/fthvgb1/wp-go/internal/pkg/logs"
	"github.com/fthvgb1/wp-go/internal/theme/common"
)

const ThemeName = "twentyfifteen"

func Init(fs embed.FS) {
	b, err := fs.ReadFile("twentyfifteen/themesupport.json")
	if err != nil {
		return
	}
	err = json.Unmarshal(b, &themesupport)
	logs.ErrPrintln(err, "解析themesupport失败")
	bytes, err := fs.ReadFile("twentyfifteen/colorscheme.json")
	if err != nil {
		return
	}
	err = json.Unmarshal(bytes, &colorscheme)
	if err != nil {
		return
	}
	logs.ErrPrintln(err, "解析colorscheme失败")
}

type handle struct {
	*common.IndexHandle
	*common.DetailHandle
}

func newHandle(iHandle *common.IndexHandle, dHandle *common.DetailHandle) *handle {
	return &handle{iHandle, dHandle}
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
	h.IndexHandle.AutoCal("colorScheme", h.colorSchemeCss)
	h.IndexHandle.AutoCal("customBackground", h.CalCustomBackGround)
	h.Indexs()
}

func (h *handle) Detail() {
	h.CustomHeader()
	h.IndexHandle.AutoCal("colorScheme", h.colorSchemeCss)
	h.IndexHandle.AutoCal("customBackground", h.CalCustomBackGround)
	h.Details()
}
