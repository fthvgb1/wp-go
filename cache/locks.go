package cache

import (
	"context"
	"github.com/fthvgb1/wp-go/safety"
	"sync"
	"sync/atomic"
)

type Locks[K comparable] struct {
	numFn   func() int
	locks   []*sync.Mutex
	m       *safety.Map[K, *sync.Mutex]
	counter *int64
}

func (l *Locks[K]) Flush(_ context.Context) {
	l.m.Flush()
	atomic.StoreInt64(l.counter, 0)
}

func (l *Locks[K]) SetNumFn(numFn func() int) {
	l.numFn = numFn
}

func NewLocks[K comparable](num func() int) *Locks[K] {
	var i int64
	return &Locks[K]{numFn: num, m: safety.NewMap[K, *sync.Mutex](), counter: &i}
}

func (l *Locks[K]) SetLockNum(num int) {
	if num > 0 {
		l.locks = make([]*sync.Mutex, num)
		for i := 0; i < num; i++ {
			l.locks[i] = &sync.Mutex{}
		}
	}
}

func (l *Locks[K]) GetLock(ctx context.Context, gMut *sync.Mutex, keys ...K) *sync.Mutex {
	k := keys[0]
	lo, ok := l.m.Load(k)
	if ok {
		return lo
	}
	num := l.numFn()
	if num == 1 {
		return gMut
	}
	gMut.Lock()
	defer gMut.Unlock()
	lo, ok = l.m.Load(k)
	if ok {
		return lo
	}
	if num <= 0 {
		lo = &sync.Mutex{}
		l.m.Store(k, lo)
		return lo
	}
	if len(l.locks) == 0 {
		l.SetLockNum(num)
	}
	counter := int(atomic.LoadInt64(l.counter))
	if counter > len(l.locks)-1 {
		atomic.StoreInt64(l.counter, 0)
		counter = 0
	}
	lo = l.locks[counter]
	l.m.Store(k, lo)
	atomic.AddInt64(l.counter, 1)
	if len(l.locks) < num {
		for i := 0; i < num-len(l.locks); i++ {
			l.locks = append(l.locks, &sync.Mutex{})
		}
	}
	return lo
}
