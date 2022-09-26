package plugins

import (
	"github.com/gin-gonic/gin"
	"github/fthvgb1/wp-go/models"
)

func NewPostPlugin(ctx *gin.Context, scene uint) *Plugin[models.WpPosts] {
	p := NewPlugin[models.WpPosts](nil, -1, nil, scene, ctx)
	p.Push(Digest)
	return p
}

func ApplyPlugin(p *Plugin[models.WpPosts], post *models.WpPosts) {
	p.post = post
	p.Next()
	p.index = -1
}
