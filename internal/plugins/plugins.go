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

	Ok
	Empty404
	Error
	InternalErr
)

var IndexSceneMap = map[int]struct{}{
	Home:     {},
	Archive:  {},
	Category: {},
	Search:   {},
}

type Func[T any] func(*Plugin[T], *gin.Context, *T, int)

type Plugin[T any] struct {
	calls []Func[T]
	index int
	post  *T
	scene int
	c     *gin.Context
}

func NewPlugin[T any](calls []Func[T], index int, post *T, scene int, c *gin.Context) *Plugin[T] {
	return &Plugin[T]{calls: calls, index: index, post: post, scene: scene, c: c}
}

func (p *Plugin[T]) Push(call ...Func[T]) {
	p.calls = append(p.calls, call...)
}

func (p *Plugin[T]) Next() {
	p.index++
	for ; p.index < len(p.calls); p.index++ {
		p.calls[p.index](p, p.c, p.post, p.scene)
	}
}
