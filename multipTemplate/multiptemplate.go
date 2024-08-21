package multipTemplate

import (
	"embed"
	"github.com/gin-gonic/gin/render"
	"html/template"
	"io/fs"
	"path/filepath"
)

type MultipleFileTemplate struct {
	Template maps
}
type MultipleFsTemplate struct {
	MultipleFileTemplate
	Fs embed.FS
}

type TemplateMaps map[string]*template.Template

func (m TemplateMaps) Load(name string) (*template.Template, bool) {
	v, ok := m[name]
	return v, ok
}

func (m TemplateMaps) Store(name string, v *template.Template) {
	m[name] = v
}

type maps interface {
	Load(name string) (*template.Template, bool)
	Store(name string, v *template.Template)
}

func (t *MultipleFileTemplate) AppendTemplate(name string, templates ...string) *MultipleFileTemplate {
	tmpl, ok := t.Template.Load(name)
	if ok {
		t.Template.Store(name, template.Must(tmpl.ParseFiles(templates...)))
	}
	return t
}

func (t *MultipleFsTemplate) AppendTemplate(name string, templates ...string) *MultipleFsTemplate {
	tmpl, ok := t.Template.Load(name)
	if ok {
		t.Template.Store(name, template.Must(tmpl.ParseFS(t.Fs, templates...)))
	}
	return t
}

func NewFileTemplates(m maps) *MultipleFileTemplate {
	return &MultipleFileTemplate{
		Template: m,
	}
}
func NewFsTemplate(f embed.FS) *MultipleFsTemplate {
	return &MultipleFsTemplate{
		MultipleFileTemplate: MultipleFileTemplate{
			Template: TemplateMaps(make(map[string]*template.Template)),
		},
		Fs: f,
	}
}
func NewFsTemplates(f embed.FS, m maps) *MultipleFsTemplate {
	return &MultipleFsTemplate{
		MultipleFileTemplate: MultipleFileTemplate{
			Template: m,
		},
		Fs: f,
	}
}

func (t *MultipleFileTemplate) SetTemplate(name string, templ *template.Template) *MultipleFileTemplate {
	t.Template.Store(name, templ)
	return t
}

func (t *MultipleFileTemplate) AddTemplate(mainTemplatePattern string, fnMap template.FuncMap, layoutTemplatePattern ...string) *MultipleFileTemplate {
	mainTemplates, err := filepath.Glob(mainTemplatePattern)
	if err != nil {
		panic(err)
	}
	for _, mainTemplate := range mainTemplates {
		file := filepath.Base(mainTemplate)
		pattern := append([]string{mainTemplate}, layoutTemplatePattern...)
		t.Template.Store(mainTemplate, template.Must(template.New(file).Funcs(fnMap).ParseFiles(pattern...)))
	}
	return t
}

func (t *MultipleFileTemplate) Instance(name string, data any) render.Render {
	v, _ := t.Template.Load(name)
	return render.HTML{
		Template: v,
		Data:     data,
	}
}

func (t *MultipleFsTemplate) AddTemplate(mainTemplatePattern string, fnMap template.FuncMap, layoutTemplatePattern ...string) *MultipleFsTemplate {
	mainTemplates, err := fs.Glob(t.Fs, mainTemplatePattern)
	if err != nil {
		panic(err)
	}
	for _, mainTemplate := range mainTemplates {
		file := filepath.Base(mainTemplate)
		pattern := append([]string{mainTemplate}, layoutTemplatePattern...)
		t.Template.Store(mainTemplate, template.Must(template.New(file).Funcs(fnMap).ParseFS(t.Fs, pattern...)))
	}
	return t
}
