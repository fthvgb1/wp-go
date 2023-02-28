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

var pipe = common.HandlePipe(common.Render, dispatch)

func Hook(h *common.Handle) {
	pipe(h)
}

func dispatch(next common.HandleFn[*common.Handle], h *common.Handle) {
	h.WidgetAreaData()
	h.GetPassword()
	h.PushHeadScript(
		common.NewComponents(CalCustomBackGround, 10),
		common.NewComponents(colorSchemeCss, 10),
	)
	h.PushHandleFn(constraints.AllStats, customHeader)
	switch h.Scene {
	case constraints.Detail:
		detail(next, h.Detail)
	default:
		index(next, h.Index)
	}
}

func index(next common.HandleFn[*common.Handle], i *common.IndexHandle) {
	i.Indexs()
}

func detail(fn common.HandleFn[*common.Handle], d *common.DetailHandle) {
	d.Details()
}
