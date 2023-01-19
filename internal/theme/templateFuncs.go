package theme

import (
	"github.com/fthvgb1/wp-go/internal/wpconfig"
	"html/template"
	"time"
)

var funcs = template.FuncMap{
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
	return funcs
}

func AddTemplateFunc(fnName string, fn any) {
	if _, ok := funcs[fnName]; ok {
		panic("exists same name func")
	}
	funcs[fnName] = fn
}
