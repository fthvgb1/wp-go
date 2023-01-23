package multipTemplate

import (
	"embed"
	"github.com/gin-gonic/gin/render"
	"html/template"
	"io/fs"
	"path/filepath"
)

type MultipleFileTemplate struct {
	Template map[string]*template.Template
	FuncMap  template.FuncMap
}
type MultipleFsTemplate struct {
	MultipleFileTemplate
	Fs embed.FS
}

func (t *MultipleFileTemplate) AppendTemplate(name string, templates ...string) *MultipleFileTemplate {
	tmpl, ok := t.Template[name]
	if ok {
		t.Template[name] = template.Must(tmpl.ParseFiles(templates...))
	}
	return t
}

func (t *MultipleFsTemplate) AppendTemplate(name string, templates ...string) *MultipleFsTemplate {
	tmpl, ok := t.Template[name]
	if ok {
		t.Template[name] = template.Must(tmpl.ParseFS(t.Fs, templates...))
	}
	return t
}

func NewFileTemplate() *MultipleFileTemplate {
	return &MultipleFileTemplate{
		Template: make(map[string]*template.Template),
		FuncMap:  make(template.FuncMap),
	}
}
func NewFsTemplate(f embed.FS) *MultipleFsTemplate {
	return &MultipleFsTemplate{
		MultipleFileTemplate: MultipleFileTemplate{
			Template: make(map[string]*template.Template),
			FuncMap:  make(template.FuncMap),
		},
		Fs: f,
	}
}

func (t *MultipleFileTemplate) SetTemplate(name string, templ *template.Template) *MultipleFileTemplate {
	if _, ok := t.Template[name]; ok {
		panic("exists same template " + name)
	}
	t.Template[name] = templ
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
		t.Template[mainTemplate] = template.Must(template.New(file).Funcs(fnMap).ParseFiles(pattern...))
	}
	return t
}

func (t *MultipleFileTemplate) Instance(name string, data any) render.Render {
	return render.HTML{
		Template: t.Template[name],
		Data:     data,
	}
}

func (t *MultipleFsTemplate) AddTemplate(mainTemplatePattern string, fnMap template.FuncMap, layoutTemplatePattern ...string) *MultipleFsTemplate {
	mainTemplates, err := fs.Glob(t.Fs, mainTemplatePattern)
	if err != nil {
		panic(err)
	}
	for _, mainTemplate := range mainTemplates {
		if _, ok := t.Template[mainTemplate]; ok {
			panic("exists same Template " + mainTemplate)
		}
		file := filepath.Base(mainTemplate)
		pattern := append([]string{mainTemplate}, layoutTemplatePattern...)
		t.Template[mainTemplate] = template.Must(template.New(file).Funcs(fnMap).ParseFS(t.Fs, pattern...))
	}
	return t
}
