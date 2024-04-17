package cachemanager

import (
	"context"
	"github.com/fthvgb1/wp-go/cache"
	"github.com/fthvgb1/wp-go/cache/reload"
	"github.com/fthvgb1/wp-go/helper"
	"github.com/fthvgb1/wp-go/helper/slice/mockmap"
	str "github.com/fthvgb1/wp-go/helper/strings"
	"github.com/fthvgb1/wp-go/safety"
	"runtime"
	"sync"
	"time"
)

var mutex sync.Mutex

type Fn func(context.Context)

type clearExpired interface {
	ClearExpired(ctx context.Context)
}

var clears = safety.NewVar(mockmap.Map[string, Fn]{})

var flushes = safety.NewVar(mockmap.Map[string, Fn]{})

func Flush() {
	ctx := context.WithValue(context.Background(), "execFlushBy", "mangerFlushFn")
	for _, f := range flushes.Load() {
		f.Value(ctx)
	}
}

func Flushes(ctx context.Context, names ...string) {
	execute(ctx, flushes, names...)
}

func execute(ctx context.Context, q *safety.Var[mockmap.Map[string, Fn]], names ...string) {
	queues := q.Load()
	for _, name := range names {
		queue := queues.Get(name)
		if queue.Value != nil {
			queue.Value(ctx)
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
			PushOrSetFlush(mockmap.Item[string, Fn]{
				Name:  name,
				Value: lo.Flush,
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
func pushOrSet(q *safety.Var[mockmap.Map[string, Fn]], queues ...mockmap.Item[string, Fn]) {
	mutex.Lock()
	defer mutex.Unlock()
	qu := q.Load()
	for _, queue := range queues {
		v := qu.Get(queue.Name)
		if v.Value != nil {
			qu.Set(queue.Name, queue.Value)
		} else {
			qu = append(qu, queue)
		}
	}
	q.Store(qu)
}

// PushOrSetFlush will execute flush func when call Flush or Flushes
func PushOrSetFlush(queues ...mockmap.Item[string, Fn]) {
	pushOrSet(flushes, queues...)
}

// PushOrSetClearExpired will execute clearExpired func when call ClearExpired or ClearExpireds
func PushOrSetClearExpired(queues ...mockmap.Item[string, Fn]) {
	pushOrSet(clears, queues...)
}

func del(q *safety.Var[mockmap.Map[string, Fn]], names ...string) {
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
		queue.Value(ctx)
	}
}
