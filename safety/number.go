package safety

import (
	"golang.org/x/exp/constraints"
	"sync/atomic"
)

func Counter[T constraints.Integer]() func() T {
	var counter int64
	return func() T {
		return T(atomic.AddInt64(&counter, 1))
	}
}
