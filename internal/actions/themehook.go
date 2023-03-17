package actions

import (
	"github.com/fthvgb1/wp-go/internal/pkg/constraints"
	"github.com/fthvgb1/wp-go/internal/theme"
	"github.com/fthvgb1/wp-go/internal/theme/wp"
	"github.com/gin-gonic/gin"
)

func ThemeHook(scene int) func(*gin.Context) {
	return func(c *gin.Context) {
		s := scene
		if scene == constraints.Home {
			if _, ok := c.GetQuery("s"); ok {
				s = constraints.Search
			}
		}
		t := theme.GetCurrentTemplateName()
		h := wp.NewHandle(c, s, t)
		h.Index = wp.NewIndexHandle(h)
		h.Detail = wp.NewDetailHandle(h)
		theme.Hook(t, h)
	}
}
