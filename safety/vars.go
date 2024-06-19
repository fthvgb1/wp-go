package safety

import (
	"sync/atomic"
	"unsafe"
)

type Var[T any] struct {
	val T
	p   unsafe.Pointer
}

func NewVar[T any](vals ...T) *Var[T] {
	var v T
	if len(vals) > 0 {
		v = vals[0]
	}
	return &Var[T]{val: v, p: unsafe.Pointer(&v)}
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

func (r *Var[T]) Flush() {
	for {
		px := atomic.LoadPointer(&r.p)
		var v T
		if atomic.CompareAndSwapPointer(&r.p, px, unsafe.Pointer(&v)) {
			return
		}
	}
}
