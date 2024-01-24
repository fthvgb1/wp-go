package theme

import (
	"github.com/fthvgb1/wp-go/app/pkg/models"
	"github.com/fthvgb1/wp-go/app/wpconfig"
	"html/template"
	"time"
)

func postsFn(fn func(models.Posts) string, a models.Posts) string {
	return fn(a)
}

func FuncMap() template.FuncMap {
	return template.FuncMap{
		"unescaped": func(s string) any {
			return template.HTML(s)
		},
		"dateCh": func(t time.Time) any {
			return t.Format("2006年 01月 02日")
		},
		"timeFormat": func(t time.Time, format string) any {
			return t.Format(format)
		},
		"getOption": func(k string) string {
			return wpconfig.GetOption(k)
		},
		"getLang": wpconfig.GetLang,
		"postsFn": postsFn,
		"exec": func(fn func() string) template.HTML {
			return template.HTML(fn())
		},
		"callFuncString": func(fn func(string) string, s string) template.HTML {
			return template.HTML(fn(s))
		},
	}
}
