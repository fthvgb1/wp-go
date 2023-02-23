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

var detailPipe = common.HandlePipe(func(d *common.DetailHandle) {
	d.Render()
}, detail)
var indexPipe = common.HandlePipe(func(i *common.IndexHandle) {
	i.Render()
}, index)

func Hook(h *common.Handle) {
	h.WidgetAreaData()
	h.GetPassword()
	h.AutoCal("colorScheme", colorSchemeCss)
	h.AutoCal("customBackground", CalCustomBackGround)
	h.PushHandleFn(customHeader)
	switch h.Scene {
	case constraints.Detail:
		detailPipe(common.NewDetailHandle(h))
	default:
		indexPipe(common.NewIndexHandle(h))
	}
}

func index(next common.HandleFn[*common.IndexHandle], i *common.IndexHandle) {
	i.Indexs()
}

func detail(fn common.HandleFn[*common.DetailHandle], d *common.DetailHandle) {
	d.Details()
}
