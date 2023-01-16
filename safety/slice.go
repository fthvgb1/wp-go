package safety

import "sync"

type Slice[T any] struct {
	Var[[]T]
	mu sync.Mutex
}

func NewSlice[T any](a []T) *Slice[T] {
	return &Slice[T]{
		NewVar(a),
		sync.Mutex{},
	}
}

func (r *Slice[T]) Append(t ...T) {
	r.mu.Lock()
	ts := append(r.Load(), t...)
	r.Store(ts)
	r.mu.Unlock()
}
