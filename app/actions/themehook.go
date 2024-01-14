package actions

import (
	"github.com/fthvgb1/wp-go/app/theme"
	"github.com/fthvgb1/wp-go/app/theme/wp"
	"github.com/gin-gonic/gin"
)

func ThemeHook(scene string) func(*gin.Context) {
	return func(c *gin.Context) {
		t := theme.GetCurrentTemplateName()
		h := wp.NewHandle(c, scene, t)
		h.Index = wp.NewIndexHandle(h)
		h.Detail = wp.NewDetailHandle(h)
		templ, _ := theme.GetTemplate(t)
		h.SetTemplate(templ)
		theme.Hook(t, h)
	}
}
