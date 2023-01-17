package stream

import (
	"github/fthvgb1/wp-go/helper"
	"github/fthvgb1/wp-go/safety"
	"github/fthvgb1/wp-go/taskPools"
)

func SimpleParallelFilterAndMap[R, T any](a SimpleSliceStream[T], fn func(T) (R, bool), c int) SimpleSliceStream[R] {
	var x []R
	rr := safety.NewSlice(x)
	a.ParallelForEach(func(t T) {
		y, ok := fn(t)
		if ok {
			rr.Append(y)
		}
	}, c)
	return SimpleSliceStream[R]{rr.Load()}
}

func SimpleParallelFilterAndMapToMap[K comparable, V any, T any](a SimpleSliceStream[T], fn func(t T) (K, V, bool), c int) (r SimpleMapStream[K, V]) {
	m := newMapX[K, V]()
	a.ParallelForEach(func(t T) {
		k, v, ok := fn(t)
		if ok {
			m.set(k, v)
		}
	}, c)
	var mm = map[K]V{}
	r = NewSimpleMapStream(mm)
	return
}

func SimpleSliceFilterAndMapToMap[K comparable, V any, T any](a SimpleSliceStream[T], fn func(t T) (K, V, bool), isCoverPrev bool) (r SimpleMapStream[K, V]) {
	m := make(map[K]V)
	a.ForEach(func(t T) {
		k, v, ok := fn(t)
		if ok {
			_, ok = m[k]
			if isCoverPrev || !ok {
				m[k] = v
			}
		}
	})
	r.m = m
	return
}

func SimpleStreamFilterAndMap[R, T any](a SimpleSliceStream[T], fn func(T) (R, bool)) SimpleSliceStream[R] {
	return NewSimpleSliceStream(helper.SliceFilterAndMap(a.arr, fn))
}

func SimpleParallelMap[R, T any](a SimpleSliceStream[T], fn func(T) R, c int) SimpleSliceStream[R] {
	var x []R
	rr := safety.NewSlice(x)
	a.ParallelForEach(func(t T) {
		rr.Append(fn(t))
	}, c)
	return SimpleSliceStream[R]{rr.Load()}
}
func SimpleStreamMap[R, T any](a SimpleSliceStream[T], fn func(T) R) SimpleSliceStream[R] {
	return NewSimpleSliceStream(helper.SliceMap(a.arr, fn))
}

func Reduce[T any, S any](s SimpleSliceStream[S], fn func(S, T) T, init T) (r T) {
	return helper.SliceReduce(s.arr, fn, init)
}

func NewSimpleSliceStream[T any](arr []T) SimpleSliceStream[T] {
	return SimpleSliceStream[T]{arr: arr}
}

type SimpleSliceStream[T any] struct {
	arr []T
}

func (r SimpleSliceStream[T]) ForEach(fn func(T)) {
	for _, t := range r.arr {
		fn(t)
	}
}

func (r SimpleSliceStream[T]) ParallelForEach(fn func(T), c int) {
	p := taskPools.NewPools(c)
	for _, t := range r.arr {
		t := t
		p.Execute(func() {
			fn(t)
		})
	}
	p.Wait()
}

func (r SimpleSliceStream[T]) ParallelFilter(fn func(T) bool, c int) SimpleSliceStream[T] {
	rr := safety.NewSlice([]T{})
	r.ParallelForEach(func(t T) {
		if fn(t) {
			rr.Append(t)
		}
	}, c)
	return SimpleSliceStream[T]{rr.Load()}
}
func (r SimpleSliceStream[T]) Filter(fn func(T) bool) SimpleSliceStream[T] {
	r.arr = helper.SliceFilter(r.arr, fn)
	return r
}

func (r SimpleSliceStream[T]) ParallelMap(fn func(T) T, c int) SimpleSliceStream[T] {
	rr := safety.NewSlice([]T{})
	r.ParallelForEach(func(t T) {
		rr.Append(fn(t))
	}, c)
	return SimpleSliceStream[T]{rr.Load()}
}

func (r SimpleSliceStream[T]) Map(fn func(T) T) SimpleSliceStream[T] {
	r.arr = helper.SliceMap(r.arr, fn)
	return r
}

func (r SimpleSliceStream[T]) Sort(fn func(i, j T) bool) SimpleSliceStream[T] {
	helper.SimpleSort(r.arr, fn)
	return r
}

func (r SimpleSliceStream[T]) Len() int {
	return len(r.arr)
}

func (r SimpleSliceStream[T]) Limit(limit, offset int) SimpleSliceStream[T] {
	l := len(r.arr)
	if offset >= l {
		return SimpleSliceStream[T]{}
	}
	ll := offset + limit
	if ll > l {
		ll = l
	}
	return SimpleSliceStream[T]{r.arr[offset:ll]}
}

func (r SimpleSliceStream[T]) Reverse() SimpleSliceStream[T] {
	helper.SliceSelfReverse(r.arr)
	return r
}

func (r SimpleSliceStream[T]) Result() []T {
	return r.arr
}
