package theme

import (
	"embed"
	"github.com/fthvgb1/wp-go/multipTemplate"
	"github.com/fthvgb1/wp-go/safety"
	"html/template"
	"io/fs"
	"path/filepath"
	"strings"
)

//go:embed *[^.go]
var TemplateFs embed.FS

var templates = safety.NewMap[string, *template.Template]() //方便外部获取模板render后的字符串，不然在gin中获取不了

var multiple *multipTemplate.MultipleFsTemplate

func BuildTemplate() *multipTemplate.MultipleFsTemplate {
	if multiple != nil {
		tt := multipTemplate.NewFsTemplate(TemplateFs)
		commonTemplate(tt)
		for k, v := range map[string]*template.Template(any(tt.Template).(multipTemplate.TemplateMaps)) {
			multiple.Template.Store(k, v)
		}
	} else {
		multiple = multipTemplate.NewFsTemplates(TemplateFs, templates)
		multiple.FuncMap = FuncMap()
		commonTemplate(multiple)
	}

	/*t.AddTemplate("twentyfifteen/*[^layout]/*.gohtml", FuncMap(), "twentyfifteen/layout/*.gohtml","wp/template.gohtml"). //单个主题设置
	AddTemplate("twentyseventeen/*[^layout]/*.gohtml", FuncMap(), "twentyseventeen/layout/*.gohtml","wp/template.gohtml")*/
	return multiple
}

func GetMultipleTemplate() *multipTemplate.MultipleFsTemplate {
	if multiple == nil {
		BuildTemplate()
	}
	return multiple
}

func GetTemplate(name string) (*template.Template, bool) {
	t, ok := templates.Load(name)
	return t, ok
}

// 所有主题模板通用设置
func commonTemplate(t *multipTemplate.MultipleFsTemplate) {
	m, err := fs.Glob(t.Fs, "*/posts/*.gohtml")
	if err != nil {
		panic(err)
	}
	for _, main := range m {
		file := filepath.Base(main)
		dir := strings.Split(main, "/")[0]
		templ := template.Must(template.New(file).Funcs(t.FuncMap).ParseFS(t.Fs, main, filepath.Join(dir, "layout/*.gohtml"), "wp/template.gohtml"))
		t.SetTemplate(main, templ)
	}
}

func IsTemplateDirExists(tml string) bool {
	arr, err := TemplateFs.ReadDir(tml)
	if err != nil {
		return false
	}
	if len(arr) > 0 {
		return true
	}
	return false
}

func IsTemplateExists(tml string) bool {
	t, ok := templates.Load(tml)
	return ok && t != nil
}

func SetTemplate(name string, val *template.Template) {
	templates.Store(name, val)
}
