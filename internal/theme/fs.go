package theme

import (
	"embed"
	"github.com/gin-gonic/gin/render"
	"html/template"
	"io/fs"
	"path/filepath"
	"strings"
)

//go:embed *[^.go]
var TemplateFs embed.FS

var templates map[string]*template.Template

type FsTemplate struct {
	Templates map[string]*template.Template
	FuncMap   template.FuncMap
}

func NewFsTemplate(funcMap template.FuncMap) *FsTemplate {
	templates = make(map[string]*template.Template)
	return &FsTemplate{FuncMap: funcMap, Templates: templates}
}

func (t FsTemplate) SetTemplate() *FsTemplate {
	mainTemplates, err := fs.Glob(TemplateFs, `*/*[^layout]/*.gohtml`)
	if err != nil {
		panic(err)
	}
	for _, include := range mainTemplates {
		name := filepath.Base(include)
		c := strings.Split(include, "/")
		base := c[0]
		t.Templates[include] = template.Must(template.New(name).Funcs(t.FuncMap).ParseFS(TemplateFs, include, filepath.Join(base, "layout/*.gohtml")))
	}
	return &t
}

func (t FsTemplate) Instance(name string, data any) render.Render {
	r := t.Templates[name]
	return render.HTML{
		Template: r,
		Data:     data,
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

func GetTemplate(tml string) *template.Template {
	return templates[tml]
}
