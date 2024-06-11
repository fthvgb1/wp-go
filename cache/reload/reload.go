package reload

import (
	"github.com/fthvgb1/wp-go/helper"
	"github.com/fthvgb1/wp-go/helper/slice"
	str "github.com/fthvgb1/wp-go/helper/strings"
	"github.com/fthvgb1/wp-go/safety"
	"sync"
	"sync/atomic"
)

type Queue struct {
	Fn       func()
	Order    float64
	Name     string
	AutoExec bool
	Once     bool
}

var mut = &sync.Mutex{}

func GetGlobeMutex() *sync.Mutex {
	return mut
}

var reloadQueues = safety.NewSlice[Queue]()

var reloadQueueHookFns = safety.NewVar[[]func(queue Queue) (Queue, bool)](nil)

var setFnVal = safety.NewMap[string, any]()

func DeleteReloadQueue(names ...string) {
	hooks := reloadQueueHookFns.Load()
	for _, name := range names {
		hooks = append(hooks, func(queue Queue) (Queue, bool) {
			if name != queue.Name {
				return queue, true
			}
			return queue, false
		})
	}
	reloadQueueHookFns.Store(hooks)
}

func HookReloadQueue(fn func(queue Queue) (Queue, bool)) {
	a := reloadQueueHookFns.Load()
	a = append(a, fn)
	reloadQueueHookFns.Store(a)
}

func GetReloadFn(name string) func() {
	hookQueue()
	i, queue := slice.SearchFirst(reloadQueues.Load(), func(queue Queue) bool {
		return queue.Name == name
	})
	if i > -1 && queue.Fn != nil {
		return queue.Fn
	}
	return nil
}

func hookQueue() {
	hooks := reloadQueueHookFns.Load()
	queues := reloadQueues.Load()
	length := len(queues)
	for _, hook := range hooks {
		queues = slice.FilterAndMap(queues, hook)
	}
	if len(queues) != length {
		reloadQueues.Store(queues)
	}
	reloadQueueHookFns.Flush()
}

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
	hookQueue()
	queues := reloadQueues.Load()
	for _, name := range namespaces {
		i, queue := slice.SearchFirst(queues, func(queue Queue) bool {
			return name == queue.Name
		})
		if i < 0 {
			continue
		}
		queue.Fn()
		if queue.Once {
			slice.Delete(&queues, i)
			reloadQueues.Store(queues)
		}
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
	Append(m.Val.Flush, args...)
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

func BuildAnyVal[T, A any](namespace string, args ...any) *SafetyVar[T, A] {
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
	Append(vv.Val.Flush, args...)
	safetyMaps.Store(namespace, vv)
	return vv
}

func GetAnyValBys[T, A any](namespace string, a A, fn func(A) (T, bool), args ...any) T {
	var vv = BuildAnyVal[T, A](namespace, args...)
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
	var vv = BuildAnyVal[T, A](namespace, args...)
	tryTimes := helper.ParseArgs(1, args...)
	var counter int64
	if tryTimes > 1 {
		Append(func() {
			atomic.StoreInt64(&counter, 0)
		}, str.Join("reload-valFn-counter-", namespace))
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
		if atomic.LoadInt64(&counter) <= 1 {
			return v.V
		}
		atomic.AddInt64(&counter, 1)
		if atomic.LoadInt64(&counter) >= int64(tryTimes) {
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
	var vv = BuildAnyVal[T, A](namespace, args...)
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
	var vv = BuildAnyVal[T, any](namespace, args...)
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
// args same as Append
//
// if give a name, then can be flushed by calls Reloads
//
// if give a float then can be reloaded early or lately, more bigger more earlier
//
// if give a bool false will not flushed when called Reload, but can call GetValMap or Reloads to flush manually
//
// if give a int 1 will only execute once when called Reload or Reloads and then delete the flush fn
func Vars[T any](defaults T, args ...any) *safety.Var[T] {
	ss := safety.NewVar(defaults)

	Append(func() {
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
// args same as Append
// if give a name, then can be flushed by calls Reloads
func VarsBy[T any](fn func() T, args ...any) *safety.Var[T] {
	ss := safety.NewVar(fn())
	Append(func() {
		ss.Store(fn())
	}, args...)
	return ss
}
func MapBy[K comparable, T any](fn func(*safety.Map[K, T]), args ...any) *safety.Map[K, T] {
	m := safety.NewMap[K, T]()
	if fn != nil {
		fn(m)
	}
	Append(func() {
		m.Flush()
		if fn != nil {
			fn(m)
		}
	}, args...)
	return m
}

func SafeMap[K comparable, T any](args ...any) *safety.Map[K, T] {
	m := safety.NewMap[K, T]()
	Append(m.Flush, args...)
	return m
}

// Append the func that will be called whenever Reload called
//
// if give a name, then can be called by called Reloads
//
// if give a float then can be called early or lately when called Reload, more bigger more earlier
//
// if give a bool false will not execute when called Reload, then can called Reloads to execute manually
//
// if give a int 1 will only execute once when called Reload or Reloads and then delete the Queue
func Append(fn func(), a ...any) {
	ord, name := parseArgs(a...)
	autoExec := helper.ParseArgs(true, a...)
	once := helper.ParseArgs(0, a...)
	queues := reloadQueues.Load()
	queue := Queue{fn, ord, name, autoExec, once == 1}
	if name != "" {
		i, _ := slice.SearchFirst(queues, func(queue Queue) bool {
			return queue.Name == name
		})
		if i > -1 {
			queues[i] = queue
			reloadQueues.Store(queues)
			return
		}
	}
	reloadQueues.Append(queue)
}

// AppendOnceFn function and args same as Append, but func will execute only once when called Reload or Reloads and then will be deleted. Especially suitable for using to develop plugins, when uninstall plugin can clean or recover some progress's self relative data or behavior which was changed by plugin.
func AppendOnceFn(fn func(), a ...any) {
	a = append([]any{1}, a...)
	Append(fn, a...)
}

func Reload() {
	mut.Lock()
	defer mut.Unlock()
	deleteMapFn.Flush()
	hookQueue()
	queues := reloadQueues.Load()
	length := len(queues)
	slice.SimpleSort(queues, slice.DESC, func(t Queue) float64 {
		return t.Order
	})
	for i, queue := range queues {
		if !queue.AutoExec {
			continue
		}
		queue.Fn()
		if queue.Once {
			slice.Delete(&queues, i)
		}
	}
	if length != len(queues) {
		reloadQueues.Store(queues)
	}
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
	Append(func() {
		if !e.isManual.Load() {
			e.v.Store(fn())
		}
	}, str.Join("fnval-", name))
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
