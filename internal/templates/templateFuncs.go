package templates

import (
	"errors"
	"github/fthvgb1/wp-go/internal/wpconfig"
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
}

func FuncMap() template.FuncMap {
	return funcs
}

func InitTemplateFunc() {

}

func AddTemplateFunc(fnName string, fn any) error {
	if _, ok := funcs[fnName]; ok {
		return errors.New("a same name func exists")
	}
	funcs[fnName] = fn
	return nil
}
