package stream

import (
	"github.com/fthvgb1/wp-go/helper/slice"
	"github.com/fthvgb1/wp-go/safety"
	"github.com/fthvgb1/wp-go/taskPools"
)

func ParallelFilterAndMap[R, T any](a Stream[T], fn func(T) (R, bool), c int) Stream[R] {
	var x []R
	rr := safety.NewSlice(x)
	a.ParallelForEach(func(t T) {
		y, ok := fn(t)
		if ok {
			rr.Append(y)
		}
	}, c)
	return Stream[R]{rr.Load()}
}

func ParallelFilterAndMapToMapStream[K comparable, V any, T any](a Stream[T], fn func(t T) (K, V, bool), c int) (r MapStream[K, V]) {
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

func SliceFilterAndMapToMapStream[K comparable, V any, T any](a Stream[T], fn func(t T) (K, V, bool), isCoverPrev bool) (r MapStream[K, V]) {
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

func FilterAndMapNewStream[R, T any](a Stream[T], fn func(T) (R, bool)) Stream[R] {
	return NewStream(slice.FilterAndMap(a.arr, fn))
}

func MapNewStream[R, T any](a Stream[T], fn func(T) R) Stream[R] {
	return NewStream(slice.Map(a.arr, fn))
}

func Reduce[T any, S any](s Stream[S], fn func(S, T) T, init T) (r T) {
	return slice.Reduce(s.arr, fn, init)
}

func NewStream[T any](arr []T) Stream[T] {
	return Stream[T]{arr: arr}
}

type Stream[T any] struct {
	arr []T
}

func (r Stream[T]) ForEach(fn func(T)) {
	for _, t := range r.arr {
		fn(t)
	}
}

func (r Stream[T]) ParallelForEach(fn func(T), c int) {
	p := taskPools.NewPools(c)
	for _, t := range r.arr {
		t := t
		p.Execute(func() {
			fn(t)
		})
	}
	p.Wait()
}

func (r Stream[T]) ParallelFilterAndMap(fn func(T) (T, bool), c int) Stream[T] {
	rr := safety.NewSlice([]T{})
	r.ParallelForEach(func(t T) {
		v, ok := fn(t)
		if ok {
			rr.Append(v)
		}
	}, c)
	return Stream[T]{rr.Load()}
}

func (r Stream[T]) FilterAndMap(fn func(T) (T, bool)) Stream[T] {
	r.arr = slice.FilterAndMap(r.arr, fn)
	return r
}

func (r Stream[T]) Reduce(fn func(v, r T) T, init T) T {
	return slice.Reduce[T, T](r.arr, fn, init)
}

func (r Stream[T]) Sort(fn func(i, j T) bool) Stream[T] {
	slice.SortSelf(r.arr, fn)
	return r
}

func (r Stream[T]) Len() int {
	return len(r.arr)
}

func (r Stream[T]) Limit(limit, offset int) Stream[T] {
	l := len(r.arr)
	if offset >= l {
		return Stream[T]{}
	}
	ll := offset + limit
	if ll > l {
		ll = l
	}
	return Stream[T]{r.arr[offset:ll]}
}

func (r Stream[T]) Reverse() Stream[T] {
	slice.ReverseSelf(r.arr)
	return r
}

func (r Stream[T]) Result() []T {
	return r.arr
}
