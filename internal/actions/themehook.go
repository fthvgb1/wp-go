package actions

import (
	"github.com/fthvgb1/wp-go/internal/theme"
	"github.com/fthvgb1/wp-go/internal/theme/common"
	"github.com/gin-gonic/gin"
)

func ThemeHook(scene int) func(*gin.Context) {
	return func(ctx *gin.Context) {
		t := theme.GetTemplateName()
		theme.Hook(t, common.NewHandle(ctx, scene, t))
	}
}
