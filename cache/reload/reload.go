package reload

import (
	"github.com/fthvgb1/wp-go/helper"
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

var mut = &sync.Mutex{}

func GetGlobeMutex() *sync.Mutex {
	return mut
}

var waitReloadCalls = safety.NewSlice(make([]queue, 0))
var callMap = safety.NewMap[string, func()]()

var setFnVal = safety.NewMap[string, any]()

type SafetyVar[T, A any] struct {
	Val   *safety.Var[Val[T]]
	Mutex sync.Mutex
}
type Val[T any] struct {
	V  T
	Ok bool
}
type SafetyMap[K comparable, V, A any] struct {
	Val   *safety.Map[K, V]
	Mutex sync.Mutex
}

var safetyMaps = safety.NewMap[string, any]()
var safetyMapLock = sync.Mutex{}

var deleteMapFn = safety.NewMap[string, func(any)]()

// GetValMap can get stored map value with namespace which called BuildSafetyMap, BuildMapFnWithConfirm, BuildMapFn, BuildMapFnWithAnyParams
func GetValMap[K comparable, V any](namespace string) (*safety.Map[K, V], bool) {
	m, ok := safetyMaps.Load(namespace)
	if !ok {
		return nil, false
	}
	v, ok := m.(*safety.Map[K, V])
	return v, ok
}

func DeleteMapVal[T any](namespace string, key ...T) {
	fn, ok := deleteMapFn.Load(namespace)
	if !ok || len(key) < 1 {
		return
	}
	fn(key)
}

func Reloads(namespaces ...string) {
	mut.Lock()
	defer mut.Unlock()
	for _, namespace := range namespaces {
		fn, ok := callMap.Load(namespace)
		if !ok {
			continue
		}
		fn()
	}
}

// BuildMapFnWithConfirm same as BuildMapFn
func BuildMapFnWithConfirm[K comparable, V, A any](namespace string, fn func(A) (V, bool), a ...any) func(key K, args A) V {
	m := BuildSafetyMap[K, V, A](namespace, a...)
	return func(key K, a A) V {
		v, ok := m.Val.Load(key)
		if ok {
			return v
		}
		m.Mutex.Lock()
		defer m.Mutex.Unlock()
		v, ok = m.Val.Load(key)
		if ok {
			return v
		}
		v, ok = fn(a)
		if ok {
			m.Val.Store(key, v)
		}
		return v
	}
}

// BuildMapFn build given fn with a new fn which returned value can be saved and flushed when called Reload or Reloads
// with namespace
//
// if give a float then can be reloaded early or lately, more bigger more earlier
//
// if give a bool false will not flushed when called Reload, then can called GetValMap to flush manually
func BuildMapFn[K comparable, V, A any](namespace string, fn func(A) V, a ...any) func(key K, args A) V {
	m := BuildSafetyMap[K, V, A](namespace, a...)
	return func(key K, a A) V {
		v, ok := m.Val.Load(key)
		if ok {
			return v
		}
		m.Mutex.Lock()
		defer m.Mutex.Unlock()
		v, ok = m.Val.Load(key)
		if ok {
			return v
		}
		v = fn(a)
		m.Val.Store(key, v)
		return v
	}
}

// BuildMapFnWithAnyParams same as BuildMapFn use multiple params
func BuildMapFnWithAnyParams[K comparable, V any](namespace string, fn func(...any) V, a ...any) func(key K, a ...any) V {
	m := BuildSafetyMap[K, V, any](namespace, a...)
	return func(key K, a ...any) V {
		v, ok := m.Val.Load(key)
		if ok {
			return v
		}
		m.Mutex.Lock()
		defer m.Mutex.Unlock()
		v, ok = m.Val.Load(key)
		if ok {
			return v
		}
		v = fn(a...)
		m.Val.Store(key, v)
		return v
	}
}

func BuildSafetyMap[K comparable, V, A any](namespace string, args ...any) *SafetyMap[K, V, A] {
	vv, ok := safetyMaps.Load(namespace)
	var m *SafetyMap[K, V, A]
	if ok {
		m = vv.(*SafetyMap[K, V, A])
		return m
	}
	safetyMapLock.Lock()
	defer safetyMapLock.Unlock()
	vv, ok = safetyMaps.Load(namespace)
	if ok {
		m = vv.(*SafetyMap[K, V, A])
		return m
	}
	m = &SafetyMap[K, V, A]{safety.NewMap[K, V](), sync.Mutex{}}
	args = append(args, namespace)
	deleteMapFn.Store(namespace, func(a any) {
		k, ok := a.([]K)
		if !ok && len(k) > 0 {
			return
		}
		for _, key := range k {
			m.Val.Delete(key)
		}
	})
	Push(m.Val.Flush, args...)
	safetyMaps.Store(namespace, m)
	return m
}

func GetAnyValMapBy[K comparable, V, A any](namespace string, key K, a A, fn func(A) (V, bool), args ...any) V {
	m := BuildSafetyMap[K, V, A](namespace, args...)
	v, ok := m.Val.Load(key)
	if ok {
		return v
	}
	m.Mutex.Lock()
	defer m.Mutex.Unlock()
	v, ok = m.Val.Load(key)
	if ok {
		return v
	}
	v, ok = fn(a)
	if ok {
		m.Val.Store(key, v)
	}
	return v
}

func BuildAnyVal[T, A any](namespace string, counter bool, args ...any) *SafetyVar[T, A] {
	var vv *SafetyVar[T, A]
	vvv, ok := safetyMaps.Load(namespace)
	if ok {
		vv = vvv.(*SafetyVar[T, A])
	}
	safetyMapLock.Lock()
	defer safetyMapLock.Unlock()
	vvv, ok = safetyMaps.Load(namespace)
	if ok {
		vv = vvv.(*SafetyVar[T, A])
		return vv
	}
	v := Val[T]{}
	vv = &SafetyVar[T, A]{safety.NewVar(v), sync.Mutex{}}
	args = append(args, namespace)
	Push(vv.Val.Flush, args...)
	safetyMaps.Store(namespace, vv)
	return vv
}

func GetAnyValBys[T, A any](namespace string, a A, fn func(A) (T, bool), args ...any) T {
	var vv = BuildAnyVal[T, A](namespace, false, args...)
	v := vv.Val.Load()
	if v.Ok {
		return v.V
	}
	vv.Mutex.Lock()
	defer vv.Mutex.Unlock()
	v = vv.Val.Load()
	if v.Ok {
		return v.V
	}
	v.V, v.Ok = fn(a)
	vv.Val.Store(v)
	return v.V
}

// BuildValFnWithConfirm same as BuildValFn
//
// if give a int and value bigger than 1 will be a times which built fn called return false
func BuildValFnWithConfirm[T, A any](namespace string, fn func(A) (T, bool), args ...any) func(A) T {
	var vv = BuildAnyVal[T, A](namespace, false, args...)
	tryTimes := helper.ParseArgs(1, args...)
	var counter func() int
	if tryTimes > 1 {
		counter = number.Counters[int]()
	}
	return func(a A) T {
		v := vv.Val.Load()
		if v.Ok {
			return v.V
		}
		vv.Mutex.Lock()
		defer vv.Mutex.Unlock()
		v = vv.Val.Load()
		if v.Ok {
			return v.V
		}
		v.V, v.Ok = fn(a)
		if v.Ok {
			vv.Val.Store(v)
			return v.V
		}
		if counter == nil {
			return v.V
		}
		times := counter()
		if times >= tryTimes {
			v.Ok = true
			vv.Val.Store(v)
		}
		return v.V
	}
}

// BuildValFn build given fn a new fn which return value can be saved and flushed when called Reload or Reloads
// with namespace.
//
// note:  namespace should be not same as BuildMapFn and related fn, they stored same safety.Map[string,any].
//
// if give a float then can be reloaded early or lately, more bigger more earlier
//
// if give a bool false will not flushed when called Reload, but can call GetValMap or Reloads to flush manually
func BuildValFn[T, A any](namespace string, fn func(A) T, args ...any) func(A) T {
	var vv = BuildAnyVal[T, A](namespace, false, args...)
	return func(a A) T {
		v := vv.Val.Load()
		if v.Ok {
			return v.V
		}
		vv.Mutex.Lock()
		defer vv.Mutex.Unlock()
		v = vv.Val.Load()
		if v.Ok {
			return v.V
		}
		v.V = fn(a)
		v.Ok = true
		vv.Val.Store(v)
		return v.V
	}
}

// BuildValFnWithAnyParams same as BuildValFn use multiple params
func BuildValFnWithAnyParams[T any](namespace string, fn func(...any) T, args ...any) func(...any) T {
	var vv = BuildAnyVal[T, any](namespace, false, args...)
	return func(a ...any) T {
		v := vv.Val.Load()
		if v.Ok {
			return v.V
		}
		vv.Mutex.Lock()
		defer vv.Mutex.Unlock()
		v = vv.Val.Load()
		if v.Ok {
			return v.V
		}
		v.V = fn(a...)
		v.Ok = true
		vv.Val.Store(v)
		return v.V
	}
}

// Vars get default value and whenever reloaded assign default value
//
// args same as Push
//
// if give a name, then can be flushed by calls Reloads
//
// if give a float then can be reloaded early or lately, more bigger more earlier
//
// if give a bool false will not flushed when called Reload, but can call GetValMap or Reloads to flush manually
func Vars[T any](defaults T, args ...any) *safety.Var[T] {
	ss := safety.NewVar(defaults)

	Push(func() {
		ss.Store(defaults)
	}, args...)
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

// VarsBy
//
// args same as Push
// if give a name, then can be flushed by calls Reloads
func VarsBy[T any](fn func() T, args ...any) *safety.Var[T] {
	ss := safety.NewVar(fn())
	Push(func() {
		ss.Store(fn())
	}, args...)
	return ss
}
func MapBy[K comparable, T any](fn func(*safety.Map[K, T]), args ...any) *safety.Map[K, T] {
	m := safety.NewMap[K, T]()
	if fn != nil {
		fn(m)
	}
	Push(func() {
		m.Flush()
		if fn != nil {
			fn(m)
		}
	}, args...)
	return m
}

func SafeMap[K comparable, T any](args ...any) *safety.Map[K, T] {
	m := safety.NewMap[K, T]()
	Push(m.Flush, args...)
	return m
}

// Push the func that will be called whenever Reload called
//
// if give a name, then can be called by called Reloads
//
// if give a float then can be called early or lately when called Reload, more bigger more earlier
//
// if give a bool false will not flushed when called Reload, then can called GetValMap to flush manually
func Push(fn func(), a ...any) {
	ord, name := parseArgs(a...)
	auto := helper.ParseArgs(true, a...)
	if name != "" && !auto {
		callMap.Store(name, fn)
		return
	}
	waitReloadCalls.Append(queue{fn, ord, name})
	if name != "" {
		callMap.Store(name, fn)
	}
}

func Reload() {
	mut.Lock()
	defer mut.Unlock()
	deleteMapFn.Flush()
	reloadCalls := waitReloadCalls.Load()
	slice.SimpleSort(reloadCalls, slice.DESC, func(t queue) float64 {
		return t.order
	})
	for _, call := range reloadCalls {
		call.fn()
	}
	return
}

type Any[T any] struct {
	fn       func() T
	v        *safety.Var[T]
	isManual *safety.Var[bool]
}

// BuildFnVal build a new fn which can be set value by SetFnVal with name or set default value by given fn when called Reload
func BuildFnVal[T any](name string, t T, fn func() T) func() T {
	if fn == nil {
		fn = func() T {
			return t
		}
	} else {
		t = fn()
	}
	p := safety.NewVar(t)
	e := Any[T]{
		fn:       fn,
		v:        p,
		isManual: safety.NewVar(false),
	}
	Push(func() {
		if !e.isManual.Load() {
			e.v.Store(fn())
		}
	})
	setFnVal.Store(name, e)
	return func() T {
		return e.v.Load()
	}
}

func SetFnVal[T any](name string, val T, onlyManual bool) {
	v, ok := setFnVal.Load(name)
	if !ok {
		return
	}
	vv, ok := v.(Any[T])
	if !ok {
		return
	}
	if onlyManual && !vv.isManual.Load() {
		vv.isManual.Store(true)
	}
	vv.v.Store(val)
}
