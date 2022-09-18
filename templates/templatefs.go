package templates

import (
	"embed"
	"github.com/gin-gonic/gin/render"
	"html/template"
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

func (t *FsTemplate) AddTemplate(name, main string, sub ...string) {
	tmp := []string{main}
	tmp = append(tmp, sub...)
	t.Templates[name] = template.Must(template.New(filepath.Base(main)).Funcs(t.FuncMap).ParseFS(TemplateFs, tmp...))
}

func (t FsTemplate) Instance(name string, data any) render.Render {
	r := t.Templates[name]
	return render.HTML{
		Template: r,
		Data:     data,
	}
}
