package reload

import (
	"github.com/fthvgb1/wp-go/safety"
	"sync"
)

var calls []func()

var str = safety.NewMap[string, string]()

var anyMap = safety.NewMap[string, any]()

type safetyVar[T, A any] struct {
	Val   *safety.Var[val[T]]
	mutex sync.Mutex
}
type val[T any] struct {
	v  T
	ok bool
}
type safetyMap[K comparable, V, A any] struct {
	val   *safety.Map[K, V]
	mutex sync.Mutex
}

var safetyMaps = safety.NewMap[string, any]()
var safetyMapLock = sync.Mutex{}

func SafetyMapByFn[K comparable, V, A any](namespace string, fn func(A) V) func(key K, args A) V {
	m := safetyMapFn[K, V, A](namespace)
	return func(key K, a A) V {
		v, ok := m.val.Load(key)
		if ok {
			return v
		}
		m.mutex.Lock()
		v, ok = m.val.Load(key)
		if ok {
			m.mutex.Unlock()
			return v
		}
		v = fn(a)
		m.val.Store(key, v)
		m.mutex.Unlock()
		return v
	}
}

func safetyMapFn[K comparable, V, A any](namespace string) *safetyMap[K, V, A] {
	vv, ok := safetyMaps.Load(namespace)
	var m *safetyMap[K, V, A]
	if ok {
		m = vv.(*safetyMap[K, V, A])
	} else {
		safetyMapLock.Lock()
		vv, ok = safetyMaps.Load(namespace)
		if ok {
			m = vv.(*safetyMap[K, V, A])
		} else {
			m = &safetyMap[K, V, A]{safety.NewMap[K, V](), sync.Mutex{}}
			Push(func() {
				m.val.Flush()
			})
			safetyMaps.Store(namespace, m)
		}
		safetyMapLock.Unlock()
	}
	return m
}

func SafetyMapBy[K comparable, V, A any](namespace string, key K, a A, fn func(A) V) V {
	m := safetyMapFn[K, V, A](namespace)
	v, ok := m.val.Load(key)
	if ok {
		return v
	}
	m.mutex.Lock()
	v, ok = m.val.Load(key)
	if ok {
		m.mutex.Unlock()
		return v
	}
	v = fn(a)
	m.val.Store(key, v)
	m.mutex.Unlock()
	return v
}

func GetAnyValBys[T, A any](namespace string, a A, fn func(A) T) T {
	var vv *safetyVar[T, A]
	vvv, ok := safetyMaps.Load(namespace)
	if ok {
		vv = vvv.(*safetyVar[T, A])
	} else {
		safetyMapLock.Lock()
		vvv, ok = safetyMaps.Load(namespace)
		if ok {
			vv = vvv.(*safetyVar[T, A])
		} else {
			vv = &safetyVar[T, A]{safety.NewVar(val[T]{}), sync.Mutex{}}
			Push(func() {
				vv.Val.Flush()
			})
			safetyMaps.Store(namespace, vv)
		}
		safetyMapLock.Unlock()
	}
	v := vv.Val.Load()
	if v.ok {
		return v.v
	}
	vv.mutex.Lock()
	v = vv.Val.Load()
	if v.ok {
		vv.mutex.Unlock()
		return v.v
	}
	v.v = fn(a)
	v.ok = true
	vv.Val.Store(v)
	vv.mutex.Unlock()
	return v.v
}

func Vars[T any](defaults T) *safety.Var[T] {
	ss := safety.NewVar(defaults)
	calls = append(calls, func() {
		ss.Store(defaults)
	})
	return ss
}
func VarsBy[T any](fn func() T) *safety.Var[T] {
	ss := safety.NewVar(fn())
	calls = append(calls, func() {
		ss.Store(fn())
	})
	return ss
}

func Push(fn ...func()) {
	calls = append(calls, fn...)
}

func Reload() {
	for _, call := range calls {
		call()
	}
	anyMap.Flush()
	str.Flush()
	safetyMaps.Flush()
}
