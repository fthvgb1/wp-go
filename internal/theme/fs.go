package theme

import (
	"embed"
	"github.com/fthvgb1/wp-go/multipTemplate"
	"html/template"
	"io/fs"
	"path/filepath"
	"strings"
)

//go:embed *[^.go]
var TemplateFs embed.FS

var templates map[string]*template.Template //方便外部获取模板render后的字符串，不然在gin中获取不了

func GetTemplate() *multipTemplate.MultipleFsTemplate {
	t := multipTemplate.NewFsTemplate(TemplateFs)
	templates = t.Template
	t.FuncMap = FuncMap()
	commonTemplate(t)
	/*t.AddTemplate("twentyfifteen/*[^layout]/*.gohtml", FuncMap(), "twentyfifteen/layout/*.gohtml"). //单个主题设置
	AddTemplate("twentyseventeen/*[^layout]/*.gohtml", FuncMap(), "twentyseventeen/layout/*.gohtml")*/
	return t
}

// 所有主题模板通用设置
func commonTemplate(t *multipTemplate.MultipleFsTemplate) {
	m, err := fs.Glob(t.Fs, "*/*[^layout]/*.gohtml")
	if err != nil {
		panic(err)
	}
	for _, main := range m {
		file := filepath.Base(main)
		dir := strings.Split(main, "/")[0]
		templ := template.Must(template.New(file).Funcs(t.FuncMap).ParseFS(t.Fs, main, filepath.Join(dir, "layout/*.gohtml")))
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
	t, ok := templates[tml]
	return ok && t != nil
}
