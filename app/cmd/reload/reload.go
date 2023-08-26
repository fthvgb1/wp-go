package reload

import (
	"github.com/fthvgb1/wp-go/helper/number"
	"github.com/fthvgb1/wp-go/safety"
	"sync"
)

var calls []func()

var anyMap = safety.NewMap[string, any]()

type safetyVar[T, A any] struct {
	Val   *safety.Var[val[T]]
	mutex sync.Mutex
}
type val[T any] struct {
	v       T
	ok      bool
	counter number.Counter[int]
}
type safetyMap[K comparable, V, A any] struct {
	val   *safety.Map[K, V]
	mutex sync.Mutex
}

var safetyMaps = safety.NewMap[string, any]()
var safetyMapLock = sync.Mutex{}

func GetAnyMapFnBys[K comparable, V, A any](namespace string, fn func(A) V) func(key K, args A) V {
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

func GetAnyValMapBy[K comparable, V, A any](namespace string, key K, a A, fn func(A) V) V {
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

func anyVal[T, A any](namespace string, counter bool) *safetyVar[T, A] {
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
			v := val[T]{}
			if counter {
				v.counter = number.Counters[int]()
			}
			vv = &safetyVar[T, A]{safety.NewVar(v), sync.Mutex{}}
			Push(func() {
				vv.Val.Flush()
			})
			safetyMaps.Store(namespace, vv)
		}
		safetyMapLock.Unlock()
	}
	return vv
}

func GetAnyValBy[T, A any](namespace string, tryTimes int, a A, fn func(A) (T, bool)) T {
	var vv = anyVal[T, A](namespace, true)
	var ok bool
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
	v.v, ok = fn(a)
	times := v.counter()
	if ok || times >= tryTimes {
		v.ok = true
		vv.Val.Store(v)
	}
	vv.mutex.Unlock()
	return v.v
}

func GetAnyValBys[T, A any](namespace string, a A, fn func(A) T) T {
	var vv = anyVal[T, A](namespace, false)
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
func MapBy[K comparable, T any](fn func(*safety.Map[K, T])) *safety.Map[K, T] {
	m := safety.NewMap[K, T]()
	if fn != nil {
		fn(m)
	}
	calls = append(calls, func() {
		m.Flush()
		if fn != nil {
			fn(m)
		}
	})
	return m
}

func SafeMap[K comparable, T any]() *safety.Map[K, T] {
	m := safety.NewMap[K, T]()
	calls = append(calls, func() {
		m.Flush()
	})
	return m
}

func Push(fn ...func()) {
	calls = append(calls, fn...)
}

func Reload() {
	for _, call := range calls {
		call()
	}
	anyMap.Flush()
	safetyMaps.Flush()
}
