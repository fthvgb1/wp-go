package signs

import (
	"github.com/fthvgb1/wp-go/helper"
	"github.com/fthvgb1/wp-go/helper/slice"
	"github.com/fthvgb1/wp-go/helper/slice/mockmap"
	"os"
	"os/signal"
	"sync"
)

type Call func() bool

type HookFn func(mockmap.Item[string, Call]) (os.Signal, mockmap.Item[string, Call], bool)

var queues = map[os.Signal]mockmap.Map[string, Call]{}

var ch = make(chan os.Signal, 1)

var stopCh = make(chan struct{}, 1)

var hooks = map[os.Signal][]HookFn{}

var mux = sync.Mutex{}

func GetChannel() chan os.Signal {
	return ch
}

func Cancel(sings ...os.Signal) {
	if len(sings) < 1 {
		return
	}
	mux.Lock()
	defer mux.Unlock()
	for _, sing := range sings {
		delete(queues, sing)
	}
	signal.Reset(sings...)
}

func Hook(sign os.Signal, fn HookFn) {
	mux.Lock()
	defer mux.Unlock()
	hooks[sign] = append(hooks[sign], fn)
}

func hook(item []mockmap.Item[string, Call], sign os.Signal) []mockmap.Item[string, Call] {
	mux.Lock()
	defer mux.Unlock()
	if len(item) < 1 {
		delete(hooks, sign)
		return item
	}
	hooksFn, ok := hooks[sign]
	if !ok {
		return item
	}
	for _, fn := range hooksFn {
		item = slice.FilterAndMap(item, func(t mockmap.Item[string, Call]) (mockmap.Item[string, Call], bool) {
			s, c, ok := fn(t)
			if sign != s {
				install(s, t.Value, t.Name, t.Order)
				return t, false
			}
			return c, ok
		})

	}
	delete(hooks, sign)
	slice.SimpleSort(item, slice.DESC, func(t mockmap.Item[string, Call]) float64 {
		return t.Order
	})
	return item
}

func Install(sign os.Signal, fn Call, a ...any) {
	mux.Lock()
	defer mux.Unlock()
	arr := helper.ParseArgs([]os.Signal{}, a...)
	if len(arr) > 0 {
		for _, o := range arr {
			install(o, fn, a...)
		}
	}
	install(sign, fn, a...)
}
func install(sign os.Signal, fn Call, a ...any) {
	m, ok := queues[sign]
	if !ok {
		queues[sign] = make(mockmap.Map[string, Call], 0)
		signal.Notify(ch, sign)
	}
	m.Set(helper.ParseArgs("", a...), fn, helper.ParseArgs[float64](0, a...))
	queues[sign] = m
}

func del(queue mockmap.Map[string, Call], sign os.Signal, i int) {
	mux.Lock()
	defer mux.Unlock()
	queue.DelByIndex(i)
	queues[sign] = queue
}

func Stop() {
	stopCh <- struct{}{}
}

func Wait() {
	for {
		select {
		case <-stopCh:
			break
		case sign := <-ch:
			queue, ok := queues[sign]
			if !ok {
				break
			}
			queue = hook(queue, sign)
			queues[sign] = queue
			if len(queue) < 1 {
				signal.Reset(sign)
				continue
			}
			for i, item := range queue {
				if !item.Value() {
					del(queue, sign, i)
				}
			}
			if len(queues[sign]) < 1 {
				signal.Reset(sign)
			}
		}
	}
}
