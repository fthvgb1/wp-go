package taskPools

import (
	"context"
	"log"
	"sync"
	"time"
)

type Pools struct {
	ch chan struct{}
	wg *sync.WaitGroup
}

func NewPools(n int) *Pools {
	if n <= 0 {
		panic("n must >= 1")
	}
	c := make(chan struct{}, n)
	for i := 0; i < n; i++ {
		c <- struct{}{}
	}
	return &Pools{
		ch: c,
		wg: &sync.WaitGroup{},
	}
}

func (p *Pools) ExecuteWithTimeOut(timeout time.Duration, fn func(), args ...string) {
	if timeout <= 0 {
		p.Execute(fn)
		return
	}
	p.wg.Add(1)
	q := <-p.ch
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		done := make(chan struct{}, 1)
		defer func() {
			cancel()
			p.wg.Done()
			p.ch <- q
		}()
		go func() {
			fn()
			done <- struct{}{}
		}()
		select {
		case <-ctx.Done():
			if len(args) > 0 && args[0] != "" {
				log.Printf("执行%s超时", args[0])
			}
		case <-done:
		}
	}()
}

func (p *Pools) Execute(fn func()) {
	if cap(p.ch) == 1 {
		fn()
		return
	}
	p.wg.Add(1)
	q := <-p.ch
	go func() {
		defer func() {
			p.wg.Done()
			p.ch <- q
		}()
		fn()
	}()
}

func (p *Pools) Wait() {
	p.wg.Wait()
}
