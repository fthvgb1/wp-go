package safety

import "sync"

type RWMap[K comparable, V any] struct {
	m   map[K]V
	mux sync.RWMutex
}

func NewRWMap[K comparable, V any](val ...map[K]V) *RWMap[K, V] {
	var m map[K]V
	if len(val) < 1 {
		m = make(map[K]V)
	} else {
		m = val[0]
	}
	return &RWMap[K, V]{m: m, mux: sync.RWMutex{}}
}

func (v *RWMap[K, V]) Store(key K, val V) {
	v.mux.Lock()
	defer v.mux.Unlock()
	v.m[key] = val
}

func (v *RWMap[K, V]) Load(key K) (V, bool) {
	v.mux.RLock()
	defer v.mux.RUnlock()
	val, ok := v.m[key]
	return val, ok
}

func (v *RWMap[K, V]) Del(keys ...K) {
	v.mux.Lock()
	defer v.mux.Unlock()
	for _, key := range keys {
		delete(v.m, key)
	}
}

func (v *RWMap[K, V]) Copy() map[K]V {
	v.mux.RLock()
	defer v.mux.RUnlock()
	var m = make(map[K]V)
	for k, val := range v.m {
		m[k] = val
	}
	return m
}

func (v *RWMap[K, V]) Len() int {
	v.mux.RLock()
	defer v.mux.RUnlock()
	return len(v.m)
}

func (v *RWMap[K, V]) Range(fn func(K, V) bool) {
	v.mux.RLock()
	defer v.mux.RUnlock()
	for key, val := range v.m {
		if !fn(key, val) {
			break
		}
	}
}

func (v *RWMap[K, V]) Set(m map[K]V) {
	v.mux.Lock()
	defer v.mux.Unlock()
	v.m = m
}
