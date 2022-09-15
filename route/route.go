package route

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github/fthvgb1/wp-go/middleware"
	"github/fthvgb1/wp-go/static"
	"github/fthvgb1/wp-go/templates"
	"html/template"
	"net/http"
	"time"
)

func SetupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()
	r.Use(middleware.SetStaticFileCache)
	r.Use(gzip.Gzip(gzip.DefaultCompression, gzip.WithExcludedPaths([]string{
		"/wp-includes/", "/wp-content/",
	})))
	r.SetFuncMap(template.FuncMap{
		"unescaped": func(s string) interface{} {
			return template.HTML(s)
		},
		"dateCh": func(t time.Time) interface{} {
			return t.Format("2006年01月02日")
		},
	})
	f := static.Fs{FS: static.FsEx, Path: "wp-includes"}
	r.StaticFS("/wp-includes", http.FS(f))
	r.StaticFS("/wp-content", http.FS(static.Fs{
		FS:   static.FsEx,
		Path: "wp-content",
	}))
	loadTemplates(r, "**/*")
	r.GET("/", index)
	r.GET("/page/:page", index)

	return r
}

func loadTemplates(engine *gin.Engine, pattern string) {
	templ := template.New("").Funcs(engine.FuncMap).Delims("{{", "}}")
	templ = template.Must(templ.ParseFS(templates.TemplateFs, pattern))
	engine.SetHTMLTemplate(templ)
}
