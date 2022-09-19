package templates

import (
	"embed"
	"github.com/gin-gonic/gin/render"
	"html/template"
	"io/fs"
	"path/filepath"
)

//go:embed posts layout
var TemplateFs embed.FS

type FsTemplate struct {
	Templates map[string]*template.Template
	FuncMap   template.FuncMap
}

func NewFsTemplate(funcMap template.FuncMap) *FsTemplate {
	return &FsTemplate{FuncMap: funcMap, Templates: make(map[string]*template.Template)}
}

func (t *FsTemplate) AddTemplate() *FsTemplate {
	mainTemplates, err := fs.Glob(TemplateFs, "*[^layout]/*.gohtml")
	if err != nil {
		panic(err)
	}
	for _, include := range mainTemplates {
		name := filepath.Base(include)
		t.Templates[include] = template.Must(template.New(name).Funcs(t.FuncMap).ParseFS(TemplateFs, include, "layout/*.gohtml"))
	}
	return t
}

func (t FsTemplate) Instance(name string, data any) render.Render {
	r := t.Templates[name]
	return render.HTML{
		Template: r,
		Data:     data,
	}
}
