package templates

import (
	"embed"
	"github.com/gin-gonic/gin/render"
	"html/template"
	"io/fs"
	"path/filepath"
	"strings"
)

//go:embed *[^templatefs.go]
var TemplateFs embed.FS

type FsTemplate struct {
	Templates map[string]*template.Template
	FuncMap   template.FuncMap
}

func NewFsTemplate(funcMap template.FuncMap) *FsTemplate {
	return &FsTemplate{FuncMap: funcMap, Templates: make(map[string]*template.Template)}
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
