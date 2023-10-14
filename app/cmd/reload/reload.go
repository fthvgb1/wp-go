package reload

import (
	"github.com/fthvgb1/wp-go/helper/number"
	"github.com/fthvgb1/wp-go/helper/slice"
	"github.com/fthvgb1/wp-go/safety"
	"sync"
)

type queue struct {
	fn    func()
	order float64
}

var calls []queue

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

func safetyMapFn[K comparable, V, A any](namespace string, order ...float64) *safetyMap[K, V, A] {
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
			}, order...)
			safetyMaps.Store(namespace, m)
		}
		safetyMapLock.Unlock()
	}
	return m
}

func GetAnyValMapBy[K comparable, V, A any](namespace string, key K, a A, fn func(A) V, order ...float64) V {
	m := safetyMapFn[K, V, A](namespace, order...)
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

func anyVal[T, A any](namespace string, counter bool, order ...float64) *safetyVar[T, A] {
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
			}, getOrder(order...))
			safetyMaps.Store(namespace, vv)
		}
		safetyMapLock.Unlock()
	}
	return vv
}

func GetAnyValBy[T, A any](namespace string, tryTimes int, a A, fn func(A) (T, bool), order ...float64) T {
	var vv = anyVal[T, A](namespace, true, order...)
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

func GetAnyValBys[T, A any](namespace string, a A, fn func(A) T, order ...float64) T {
	var vv = anyVal[T, A](namespace, false, order...)
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

func Vars[T any](defaults T, order ...float64) *safety.Var[T] {
	ss := safety.NewVar(defaults)
	ord := getOrder(order...)
	calls = append(calls, queue{func() {
		ss.Store(defaults)
	}, ord})
	return ss
}

func getOrder(order ...float64) float64 {
	var ord float64
	if len(order) > 0 {
		ord = order[0]
	}
	return ord
}
func VarsBy[T any](fn func() T, order ...float64) *safety.Var[T] {
	ss := safety.NewVar(fn())
	ord := getOrder(order...)
	calls = append(calls, queue{
		func() {
			ss.Store(fn())
		}, ord,
	})
	return ss
}
func MapBy[K comparable, T any](fn func(*safety.Map[K, T]), order ...float64) *safety.Map[K, T] {
	m := safety.NewMap[K, T]()
	if fn != nil {
		fn(m)
	}
	ord := getOrder(order...)
	calls = append(calls, queue{
		func() {
			m.Flush()
			if fn != nil {
				fn(m)
			}
		}, ord,
	})
	return m
}

func SafeMap[K comparable, T any](order ...float64) *safety.Map[K, T] {
	m := safety.NewMap[K, T]()
	ord := getOrder(order...)
	calls = append(calls, queue{func() {
		m.Flush()
	}, ord})
	return m
}

func Push(fn func(), order ...float64) {
	ord := getOrder(order...)
	calls = append(calls, queue{fn, ord})
}

func Reload() {
	slice.Sort(calls, func(i, j queue) bool {
		return i.order > j.order
	})
	for _, call := range calls {
		call.fn()
	}
	anyMap.Flush()
	safetyMaps.Flush()
}
