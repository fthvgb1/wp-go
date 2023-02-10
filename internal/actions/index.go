package actions

import (
	"github.com/fthvgb1/wp-go/internal/pkg/constraints"
	"github.com/fthvgb1/wp-go/internal/theme"
	"github.com/fthvgb1/wp-go/internal/theme/common"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Index(scene int) func(*gin.Context) {
	return func(ctx *gin.Context) {
		theme.Hook(theme.GetTemplateName(), common.Handle{
			C:     ctx,
			GinH:  gin.H{},
			Scene: scene,
			Code:  http.StatusOK,
			Stats: constraints.Ok,
		})
	}
}
