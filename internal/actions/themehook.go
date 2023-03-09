package actions

import (
	"github.com/fthvgb1/wp-go/internal/theme"
	"github.com/fthvgb1/wp-go/internal/theme/wp"
	"github.com/gin-gonic/gin"
)

func ThemeHook(scene int) func(*gin.Context) {
	return func(ctx *gin.Context) {
		t := theme.GetCurrentTemplateName()
		h := wp.NewHandle(ctx, scene, t)
		h.Index = wp.NewIndexHandle(h)
		h.Detail = wp.NewDetailHandle(h)
		theme.Hook(t, h)
	}
}
