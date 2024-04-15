package cachemanager

import (
	"context"
	"github.com/fthvgb1/wp-go/cache"
	"github.com/fthvgb1/wp-go/cache/reload"
	"github.com/fthvgb1/wp-go/helper"
	"github.com/fthvgb1/wp-go/helper/slice"
	str "github.com/fthvgb1/wp-go/helper/strings"
	"github.com/fthvgb1/wp-go/safety"
	"runtime"
	"sync"
	"time"
)

var mutex sync.Mutex

type Queue struct {
	Name string
	Fn   func(context.Context)
}

type Queues []Queue

func (q *Queues) Get(name string) Queue {
	_, v := slice.SearchFirst(*q, func(t Queue) bool {
		return name == t.Name
	})
	return v
}

func (q *Queues) Set(name string, fn func(context.Context)) {
	i := slice.IndexOfBy(*q, func(t Queue) bool {
		return name == t.Name
	})
	if i > -1 {
		(*q)[i].Fn = fn
		return
	}
	*q = append(*q, Queue{name, fn})
}

func (q *Queues) Del(name string) {
	i := slice.IndexOfBy(*q, func(t Queue) bool {
		return name == t.Name
	})
	if i > -1 {
		slice.Delete((*[]Queue)(q), i)
	}
}

type clearExpired interface {
	ClearExpired(ctx context.Context)
}

var clears = safety.NewVar(Queues{})

var flushes = safety.NewVar(Queues{})

func Flush() {
	ctx := context.WithValue(context.Background(), "execFlushBy", "mangerFlushFn")
	for _, f := range flushes.Load() {
		f.Fn(ctx)
	}
}

func Flushes(ctx context.Context, names ...string) {
	execute(ctx, flushes, names...)
}

func execute(ctx context.Context, q *safety.Var[Queues], names ...string) {
	queues := q.Load()
	for _, name := range names {
		queue := queues.Get(name)
		if queue.Fn != nil {
			queue.Fn(ctx)
		}
	}
}

func parseArgs(args ...any) (string, func() time.Duration) {
	var name string
	var fn func() time.Duration
	for _, arg := range args {
		v, ok := arg.(string)
		if ok {
			name = v
			continue
		}
		vv, ok := arg.(func() time.Duration)
		if ok {
			fn = vv
		}

	}
	return name, fn
}

func buildLockFn[K comparable](args ...any) cache.LockFn[K] {
	lockFn := helper.ParseArgs(cache.LockFn[K](nil), args...)
	name := helper.ParseArgs("", args...)
	num := helper.ParseArgs(runtime.NumCPU(), args...)
	loFn := func() int {
		return num
	}
	loFn = helper.ParseArgs(loFn, args...)
	if name != "" {
		loFn = reload.BuildFnVal(str.Join("cachesLocksNum-", name), num, loFn)
	}
	if lockFn == nil {
		looo := helper.ParseArgs(cache.Lockss[K](nil), args...)
		if looo != nil {
			lockFn = looo.GetLock
			loo, ok := any(looo).(cache.LocksNum)
			if ok && loo != nil {
				loo.SetLockNum(num)
			}
		} else {
			lo := cache.NewLocks[K](loFn)
			lockFn = lo.GetLock
			PushOrSetFlush(Queue{
				Name: name,
				Fn:   lo.Flush,
			})
		}

	}
	return lockFn
}

func SetExpireTime(c cache.SetTime, name string, expireTime time.Duration, expireTimeFn func() time.Duration) {
	if name == "" {
		return
	}
	fn := reload.BuildFnVal(str.Join("cacheManger-", name, "-expiredTime"), expireTime, expireTimeFn)
	c.SetExpiredTime(fn)
}

func ChangeExpireTime(t time.Duration, coverConf bool, name ...string) {
	for _, s := range name {
		reload.SetFnVal(s, t, coverConf)
	}
}
func pushOrSet(q *safety.Var[Queues], queues ...Queue) {
	mutex.Lock()
	defer mutex.Unlock()
	qu := q.Load()
	for _, queue := range queues {
		v := qu.Get(queue.Name)
		if v.Fn != nil {
			qu.Set(queue.Name, queue.Fn)
		} else {
			qu = append(qu, queue)
		}
	}
	q.Store(qu)
}

// PushOrSetFlush will execute flush func when call Flush or Flushes
func PushOrSetFlush(queues ...Queue) {
	pushOrSet(flushes, queues...)
}

// PushOrSetClearExpired will execute clearExpired func when call ClearExpired or ClearExpireds
func PushOrSetClearExpired(queues ...Queue) {
	pushOrSet(clears, queues...)
}

func del(q *safety.Var[Queues], names ...string) {
	mutex.Lock()
	defer mutex.Unlock()
	queues := q.Load()
	for _, name := range names {
		queues.Del(name)
	}
	q.Store(queues)
}
func DelFlush(names ...string) {
	del(flushes, names...)
}

func DelClearExpired(names ...string) {
	del(clears, names...)
}

func ClearExpireds(ctx context.Context, names ...string) {
	execute(ctx, clears, names...)
}

func ClearExpired() {
	ctx := context.WithValue(context.Background(), "execClearExpired", "mangerClearExpiredFn")
	for _, queue := range clears.Load() {
		queue.Fn(ctx)
	}
}
