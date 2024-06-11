package safety

import "sync"

type Slice[T any] struct {
	*Var[[]T]
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
	defer r.mu.Unlock()
	ts := append(r.Var.Load(), t...)
	r.Store(ts)
}

func (r *Slice[T]) Set(index int, val T) {
	v := r.Var.Load()
	if index >= len(v) {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	v[index] = val
}

func (r *Slice[T]) Load() (a []T) {
	v := r.Var.Load()
	a = make([]T, len(v))
	copy(a, v)
	return
}
