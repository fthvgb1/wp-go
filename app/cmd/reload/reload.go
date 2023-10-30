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
	name  string
}

var calls = safety.NewSlice(make([]queue, 0))
var callsM = safety.NewMap[string, func()]()

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

var flushMapFn = safety.NewMap[string, func(any)]()

func FlushMapVal[T any](namespace string, key ...T) {
	fn, ok := flushMapFn.Load(namespace)
	if !ok || len(key) < 1 {
		return
	}
	fn(key)
}

func FlushAnyVal(namespaces ...string) {
	for _, namespace := range namespaces {
		fn, ok := callsM.Load(namespace)
		if !ok {
			continue
		}
		fn()
	}
}

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

func safetyMapFn[K comparable, V, A any](namespace string, args ...any) *safetyMap[K, V, A] {
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
			ord, _ := parseArgs(args...)
			flushMapFn.Store(namespace, func(a any) {
				k, ok := a.([]K)
				if !ok && len(k) > 0 {
					return
				}
				for _, key := range k {
					m.val.Delete(key)
				}
			})
			Push(func() {
				m.val.Flush()
			}, ord, namespace)
			safetyMaps.Store(namespace, m)
		}
		safetyMapLock.Unlock()
	}
	return m
}

func GetAnyValMapBy[K comparable, V, A any](namespace string, key K, a A, fn func(A) V, args ...any) V {
	m := safetyMapFn[K, V, A](namespace, args...)
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

func anyVal[T, A any](namespace string, counter bool, args ...any) *safetyVar[T, A] {
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
			ord, _ := parseArgs(args...)
			Push(func() {
				vv.Val.Flush()
			}, ord, namespace)
			safetyMaps.Store(namespace, vv)
		}
		safetyMapLock.Unlock()
	}
	return vv
}

func GetAnyValBy[T, A any](namespace string, tryTimes int, a A, fn func(A) (T, bool), args ...any) T {
	var vv = anyVal[T, A](namespace, true, args...)
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

func GetAnyValBys[T, A any](namespace string, a A, fn func(A) T, args ...any) T {
	var vv = anyVal[T, A](namespace, false, args...)
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

func Vars[T any](defaults T, args ...any) *safety.Var[T] {
	ss := safety.NewVar(defaults)
	ord, name := parseArgs(args...)
	Push(func() {
		ss.Store(defaults)
	}, ord, name)
	return ss
}

func parseArgs(a ...any) (ord float64, name string) {
	if len(a) > 0 {
		for _, arg := range a {
			v, ok := arg.(float64)
			if ok {
				ord = v
			}
			vv, ok := arg.(string)
			if ok {
				name = vv
			}
		}
	}
	return ord, name
}
func VarsBy[T any](fn func() T, args ...any) *safety.Var[T] {
	ss := safety.NewVar(fn())
	ord, name := parseArgs(args...)
	Push(func() {
		ss.Store(fn())
	}, ord, name)
	return ss
}
func MapBy[K comparable, T any](fn func(*safety.Map[K, T]), args ...any) *safety.Map[K, T] {
	m := safety.NewMap[K, T]()
	if fn != nil {
		fn(m)
	}
	ord, name := parseArgs(args...)
	Push(func() {
		m.Flush()
		if fn != nil {
			fn(m)
		}
	}, ord, name)
	return m
}

func SafeMap[K comparable, T any](args ...any) *safety.Map[K, T] {
	m := safety.NewMap[K, T]()
	ord, name := parseArgs(args...)
	Push(func() {
		m.Flush()
	}, ord, name)
	return m
}

func Push(fn func(), a ...any) {
	ord, name := parseArgs(a...)
	calls.Append(queue{fn, ord, name})
	//calls = append(calls, queue{fn, ord, name})
	if name != "" {
		callsM.Store(name, fn)
	}
}

func Reload() {
	anyMap.Flush()
	safetyMaps.Flush()
	callsM.Flush()
	flushMapFn.Flush()
	slice.Sort(calls.Load(), func(i, j queue) bool {
		return i.order > j.order
	})
	for _, call := range calls.Load() {
		call.fn()
	}
	return
}
