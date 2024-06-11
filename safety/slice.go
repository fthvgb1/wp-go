package safety

import "sync"

type Slice[T any] struct {
	Var []T
	mu  sync.RWMutex
}

func NewSlice[T any](a ...[]T) *Slice[T] {
	var s []T
	if len(a) > 0 {
		s = a[0]
	}
	return &Slice[T]{
		s,
		sync.RWMutex{},
	}
}

func (r *Slice[T]) Append(t ...T) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Var = append(r.Var, t...)
}

func (r *Slice[T]) Store(a []T) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Var = a
}

func (r *Slice[T]) Load() (a []T) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	a = make([]T, len(r.Var))
	copy(a, r.Var)
	return
}
