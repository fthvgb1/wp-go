package actions

import (
	"github.com/fthvgb1/wp-go/internal/theme"
	"github.com/fthvgb1/wp-go/internal/theme/common"
	"github.com/gin-gonic/gin"
)

func ThemeHook(scene int) func(*gin.Context) {
	return func(ctx *gin.Context) {
		theme.Hook(theme.GetTemplateName(), common.NewHandle(ctx, scene))
	}
}
