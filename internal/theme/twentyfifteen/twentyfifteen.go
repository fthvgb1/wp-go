package twentyfifteen

import (
	"embed"
	"encoding/json"
	"github.com/fthvgb1/wp-go/internal/pkg/constraints"
	"github.com/fthvgb1/wp-go/internal/pkg/logs"
	"github.com/fthvgb1/wp-go/internal/theme/wp"
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

var pipe = wp.HandlePipe(wp.Render, dispatch)

func Hook(h *wp.Handle) {
	pipe(h)
}

func dispatch(next wp.HandleFn[*wp.Handle], h *wp.Handle) {
	h.WidgetAreaData()
	h.GetPassword()
	h.PushHeadScript(
		wp.NewComponents(CalCustomBackGround, 10),
		wp.NewComponents(colorSchemeCss, 10),
	)
	h.PushHandleFn(constraints.AllStats, wp.NewHandleFn(customHeader, 10))
	switch h.Scene() {
	case constraints.Detail:
		detail(next, h.Detail)
	default:
		index(next, h.Index)
	}
}

func index(next wp.HandleFn[*wp.Handle], i *wp.IndexHandle) {
	i.Indexs()
}

func detail(fn wp.HandleFn[*wp.Handle], d *wp.DetailHandle) {
	d.Details()
}
