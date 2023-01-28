package stream

import (
	"github.com/fthvgb1/wp-go/helper/maps"
	"github.com/fthvgb1/wp-go/taskPools"
	"sync"
)

type mapX[K comparable, V any] struct {
	m   map[K]V
	mut sync.Mutex
}

func (r *mapX[K, V]) set(k K, v V) {
	r.mut.Lock()
	defer r.mut.Unlock()
	r.m[k] = v
}

func newMapX[K comparable, V any]() mapX[K, V] {
	return mapX[K, V]{
		m:   map[K]V{},
		mut: sync.Mutex{},
	}
}

func SimpleMapFilterAndMapToSlice[R any, K comparable, V any](mm MapStream[K, V], fn func(K, V) (R, bool)) Stream[R] {
	return NewStream(maps.FilterToSlice(mm.m, fn))
}

func SimpleMapParallelFilterAndMapToMap[K comparable, V any, KK comparable, VV any](mm MapStream[KK, VV], fn func(KK, VV) (K, V, bool), c int) MapStream[K, V] {
	m := newMapX[K, V]()
	mm.ParallelForEach(func(kk KK, vv VV) {
		k, v, ok := fn(kk, vv)
		if ok {
			m.set(k, v)
		}
	}, c)
	return MapStream[K, V]{m.m}
}

func SimpleMapStreamFilterAndMapToMap[K comparable, V any, KK comparable, VV comparable](a MapStream[KK, VV], fn func(KK, VV) (K, V, bool)) (r MapStream[K, V]) {
	r = MapStream[K, V]{make(map[K]V)}
	for k, v := range a.m {
		kk, vv, ok := fn(k, v)
		if ok {
			r.m[kk] = vv
		}
	}
	return
}

func NewSimpleMapStream[K comparable, V any](m map[K]V) MapStream[K, V] {
	return MapStream[K, V]{m}
}

type MapStream[K comparable, V any] struct {
	m map[K]V
}

func (r MapStream[K, V]) ForEach(fn func(K, V)) {
	for k, v := range r.m {
		fn(k, v)
	}
}

func (r MapStream[K, V]) ParallelForEach(fn func(K, V), c int) {
	p := taskPools.NewPools(c)
	for k, v := range r.m {
		k := k
		v := v
		p.Execute(func() {
			fn(k, v)
		})
	}
	p.Wait()
}

func (r MapStream[K, V]) Len() int {
	return len(r.m)
}

func (r MapStream[K, V]) Result() map[K]V {
	return r.m
}
