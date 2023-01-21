package theme

import (
	"github.com/fthvgb1/wp-go/internal/wpconfig"
	"html/template"
	"time"
)

var comFn = template.FuncMap{
	"unescaped": func(s string) any {
		return template.HTML(s)
	},
	"dateCh": func(t time.Time) any {
		return t.Format("2006年 01月 02日")
	},
	"getOption": func(k string) string {
		return wpconfig.Options.Value(k)
	},
	"getLang": wpconfig.GetLang,
}

func FuncMap() template.FuncMap {
	return comFn
}

func AddTemplateFunc(fnName string, fn any) {
	if _, ok := comFn[fnName]; ok {
		panic("exists same name func")
	}
	comFn[fnName] = fn
}
