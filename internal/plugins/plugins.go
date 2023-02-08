package plugins

import (
	"github.com/gin-gonic/gin"
)

const (
	Home = iota + 1
	Archive
	Category
	Tag
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
	Tag:      {},
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
