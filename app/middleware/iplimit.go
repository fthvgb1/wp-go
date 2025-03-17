package middleware

import (
	"github.com/fthvgb1/wp-go/helper/number"
	str "github.com/fthvgb1/wp-go/helper/strings"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type LimitMap[K comparable] struct {
	Mux      *sync.RWMutex
	Map      map[K]*int64
	LimitNum *int64
	ClearNum *int64
}

type FlowLimits[K comparable] struct {
	GetKeyFn     func(ctx *gin.Context) K
	LimitedFns   func(ctx *gin.Context)
	DeferClearFn func(c *gin.Context, m LimitMap[K], k K, v *int64)
	AddFn        func(c *gin.Context, m LimitMap[K], k K, v *int64)
}

func NewFlowLimits[K comparable](getKeyFn func(ctx *gin.Context) K,
	limitedFns func(ctx *gin.Context),
	deferClearFn func(c *gin.Context, m LimitMap[K], k K, v *int64),
	addFns ...func(c *gin.Context, m LimitMap[K], k K, v *int64),
) *FlowLimits[K] {

	f := &FlowLimits[K]{
		GetKeyFn:     getKeyFn,
		LimitedFns:   limitedFns,
		DeferClearFn: deferClearFn,
	}
	fn := f.Adds
	if len(addFns) > 0 {
		fn = addFns[0]
	}
	f.AddFn = fn
	return f
}

func (f FlowLimits[K]) GetKey(c *gin.Context) K {
	return f.GetKeyFn(c)
}

func (f FlowLimits[K]) Limit(c *gin.Context) {
	f.LimitedFns(c)
}
func (f FlowLimits[K]) DeferClear(c *gin.Context, m LimitMap[K], k K, v *int64) {
	f.DeferClearFn(c, m, k, v)
}

func (f FlowLimits[K]) Add(c *gin.Context, m LimitMap[K], k K, v *int64) {
	f.AddFn(c, m, k, v)
}

func (f FlowLimits[K]) Adds(_ *gin.Context, _ LimitMap[K], _ K, v *int64) {
	atomic.AddInt64(v, 1)
}

type MapFlowLimit[K comparable] interface {
	GetKey(c *gin.Context) K
	Limit(c *gin.Context)
	DeferClear(c *gin.Context, m LimitMap[K], k K, v *int64)
	Add(c *gin.Context, m LimitMap[K], k K, v *int64)
}

func CustomFlowLimit[K comparable](a MapFlowLimit[K], maxRequestNum int64, clearNum ...int64) (func(ctx *gin.Context), func(int64, ...int64)) {
	m := LimitMap[K]{
		Mux:      &sync.RWMutex{},
		Map:      make(map[K]*int64),
		LimitNum: new(int64),
		ClearNum: new(int64),
	}
	atomic.StoreInt64(m.LimitNum, maxRequestNum)
	if len(clearNum) > 0 {
		atomic.StoreInt64(m.ClearNum, clearNum[0])
	}

	fn := func(num int64, clearNum ...int64) {
		atomic.StoreInt64(m.LimitNum, num)
		if len(clearNum) > 0 {
			atomic.StoreInt64(m.ClearNum, clearNum[0])
		}
	}
	return func(c *gin.Context) {
		if atomic.LoadInt64(m.LimitNum) <= 0 {
			c.Next()
			return
		}
		key := a.GetKey(c)
		m.Mux.RLock()
		i, ok := m.Map[key]
		m.Mux.RUnlock()
		if !ok {
			m.Mux.Lock()
			i = new(int64)
			m.Map[key] = i
			m.Mux.Unlock()
		}
		a.Add(c, m, key, i)
		defer a.DeferClear(c, m, key, i)
		if atomic.LoadInt64(i) > atomic.LoadInt64(m.LimitNum) {
			a.Limit(c)
			return
		}
		c.Next()
	}, fn
}

func IpLimitClear[K comparable](_ *gin.Context, m LimitMap[K], key K, i *int64) {
	atomic.AddInt64(i, -1)
	if atomic.LoadInt64(i) <= 0 {
		cNum := int(atomic.LoadInt64(m.ClearNum))
		if cNum <= 0 {
			m.Mux.Lock()
			delete(m.Map, key)
			m.Mux.Unlock()
			return
		}

		m.Mux.RLock()
		l := len(m.Map)
		m.Mux.RUnlock()
		if l < cNum {
			m.Mux.Lock()
			for k, v := range m.Map {
				if atomic.LoadInt64(v) < 1 {
					delete(m.Map, k)
				}
			}
			m.Mux.Unlock()
		}
	}
}

func ToManyRequest(messages ...string) func(c *gin.Context) {
	message := "请求太多了，服务器君表示压力山大==!, 请稍后访问"
	if len(messages) > 0 {
		message = messages[0]
	}
	return func(c *gin.Context) {
		c.String(http.StatusForbidden, message)
		c.Abort()
	}
}

func IpLimit(num int64, clearNum ...int64) (func(ctx *gin.Context), func(int64, ...int64)) {
	a := NewFlowLimits(func(c *gin.Context) string {
		return c.ClientIP()
	}, ToManyRequest(), IpLimitClear)
	return CustomFlowLimit[string](a, num, clearNum...)
}

func IpMinuteLimit(num int64, clearNum ...int64) (func(ctx *gin.Context), func(int64, ...int64)) {
	total := new(int64)
	a := NewFlowLimits(func(c *gin.Context) string {
		return str.Join(c.ClientIP(), "|", time.Now().Format("2006-01-02 15:04"))
	},
		ToManyRequest(),
		IpMinuteLimitDeferFn(total),
		func(c *gin.Context, m LimitMap[string], k string, v *int64) {
			atomic.AddInt64(v, 1)
			atomic.AddInt64(total, 1)
		},
	)

	return CustomFlowLimit(a, num, clearNum...)
}

func IpMinuteLimitDeferFn(total *int64) func(_ *gin.Context, m LimitMap[string], k string, _ *int64) {
	return func(_ *gin.Context, m LimitMap[string], k string, _ *int64) {
		atomic.AddInt64(total, -1)
		cNum := number.Min(int(atomic.LoadInt64(m.ClearNum)), 1)
		minu := strings.Split(k, "|")[1]
		if int(atomic.LoadInt64(total)) < cNum {
			m.Mux.Lock()
			for key := range m.Map {
				t := strings.Split(key, "|")[1]
				if minu != t {
					delete(m.Map, key)
				}
			}
			m.Mux.Unlock()
		}
	}
}
