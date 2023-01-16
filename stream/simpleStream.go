package stream

import (
	"github/fthvgb1/wp-go/helper"
	"github/fthvgb1/wp-go/safety"
	"github/fthvgb1/wp-go/taskPools"
)

func Reduce[T any, S any](s SimpleSliceStream[S], fn func(S, T) T, init T) (r T) {
	return helper.SliceReduce(s.arr, fn, init)
}

type SimpleSliceStream[T any] struct {
	arr []T
}

func NewSimpleSliceStream[T any](arr []T) SimpleSliceStream[T] {
	return SimpleSliceStream[T]{arr: arr}
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
	p := taskPools.NewPools(c)
	var x []T
	rr := safety.NewSlice(x)
	for _, t := range r.arr {
		t := t
		p.Execute(func() {
			if fn(t) {
				rr.Append(t)
			}
		})
	}
	p.Wait()
	return SimpleSliceStream[T]{rr.Load()}
}
func (r SimpleSliceStream[T]) Filter(fn func(T) bool) SimpleSliceStream[T] {
	r.arr = helper.SliceFilter(r.arr, fn)
	return r
}

func (r SimpleSliceStream[T]) ParallelMap(fn func(T) T, c int) SimpleSliceStream[T] {
	p := taskPools.NewPools(c)
	var x []T
	rr := safety.NewSlice(x)
	for _, t := range r.arr {
		t := t
		p.Execute(func() {
			rr.Append(fn(t))
		})
	}
	p.Wait()
	return SimpleSliceStream[T]{rr.Load()}
}
func SimpleParallelFilterAndMap[R, T any](a SimpleSliceStream[T], fn func(T) (R, bool), c int) SimpleSliceStream[R] {
	p := taskPools.NewPools(c)
	var x []R
	rr := safety.NewSlice(x)
	for _, t := range a.arr {
		t := t
		p.Execute(func() {
			y, ok := fn(t)
			if ok {
				rr.Append(y)
			}
		})
	}
	p.Wait()
	return SimpleSliceStream[R]{rr.Load()}
}

func SimpleStreamFilterAndMap[R, T any](a SimpleSliceStream[T], fn func(T) (R, bool)) SimpleSliceStream[R] {
	return NewSimpleSliceStream(helper.SliceFilterAndMap(a.arr, fn))
}

func SimpleParallelMap[R, T any](a SimpleSliceStream[T], fn func(T) R, c int) SimpleSliceStream[R] {
	p := taskPools.NewPools(c)
	var x []R
	rr := safety.NewSlice(x)
	for _, t := range a.arr {
		t := t
		p.Execute(func() {
			rr.Append(fn(t))
		})
	}
	p.Wait()
	return SimpleSliceStream[R]{rr.Load()}
}
func SimpleStreamMap[R, T any](a SimpleSliceStream[T], fn func(T) R) SimpleSliceStream[R] {
	return NewSimpleSliceStream(helper.SliceMap(a.arr, fn))
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
	return SimpleSliceStream[T]{r.arr[offset : offset+limit]}
}

func (r SimpleSliceStream[T]) Reverse() SimpleSliceStream[T] {
	helper.SliceSelfReverse(r.arr)
	return r
}

func (r SimpleSliceStream[T]) Result() []T {
	return r.arr
}
