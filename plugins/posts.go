package plugins

import (
	"github.com/gin-gonic/gin"
	"github/fthvgb1/wp-go/models/wp"
)

func NewPostPlugin(ctx *gin.Context, scene uint) *Plugin[wp.Posts] {
	p := NewPlugin[wp.Posts](nil, -1, nil, scene, ctx)
	p.Push(Digest)
	return p
}

func ApplyPlugin(p *Plugin[wp.Posts], post *wp.Posts) {
	p.post = post
	p.Next()
	p.index = -1
}
