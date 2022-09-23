package plugins

import (
	"github.com/gin-gonic/gin"
)

const (
	Home = iota + 1
	Archive
	Category
	Search
	Detail
)

type PluginFunc[T any] func(*Plugin[T], *gin.Context, *T, uint)

type Plugin[T any] struct {
	calls []PluginFunc[T]
	index int
	post  *T
	scene uint
	c     *gin.Context
}

func NewPlugin[T any](calls []PluginFunc[T], index int, post *T, scene uint, c *gin.Context) *Plugin[T] {
	return &Plugin[T]{calls: calls, index: index, post: post, scene: scene, c: c}
}

func (p *Plugin[T]) Push(call ...PluginFunc[T]) {
	p.calls = append(p.calls, call...)
}

func (p *Plugin[T]) Next() {
	p.index++
	for ; p.index < len(p.calls); p.index++ {
		p.calls[p.index](p, p.c, p.post, p.scene)
	}
}
