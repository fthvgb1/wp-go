package safety

import (
	"sync/atomic"
	"unsafe"
)

type Var[T any] struct {
	val T
	p   unsafe.Pointer
}

func NewVar[T any](val T) Var[T] {
	return Var[T]{val: val, p: unsafe.Pointer(&val)}
}

func (r *Var[T]) Load() T {
	return *(*T)(atomic.LoadPointer(&r.p))
}

func (r *Var[T]) Delete() {
	for {
		px := atomic.LoadPointer(&r.p)
		if atomic.CompareAndSwapPointer(&r.p, px, nil) {
			return
		}
	}
}

func (r *Var[T]) Store(v T) {
	for {
		px := atomic.LoadPointer(&r.p)
		if atomic.CompareAndSwapPointer(&r.p, px, unsafe.Pointer(&v)) {
			return
		}
	}
}
