package cache

import (
	"context"
	"errors"
	"fmt"
	"github.com/fthvgb1/wp-go/helper"
	"github.com/fthvgb1/wp-go/helper/number"
	str "github.com/fthvgb1/wp-go/helper/strings"
	"github.com/fthvgb1/wp-go/safety"
	"strings"
	"time"
)

type Pagination[K comparable, V any] struct {
	*MapCache[string, helper.PaginationData[V]]
	maxNum        func() int
	isSwitch      *safety.Map[K, bool]
	dbFn          func(ctx context.Context, k K, page, limit, totalRaw int, a ...any) ([]V, int, error)
	localFn       func(ctx context.Context, data []V, k K, page, limit int, a ...any) ([]V, int, error)
	batchFetchNum func() int
	localKeyFn    func(K K, a ...any) string
	dbKeyFn       func(K K, a ...any) string
	name          string
}

var switchDb = errors.New("switch Db")

type DbFn[K comparable, V any] func(ctx context.Context, k K, page, limit, totalRaw int, a ...any) ([]V, int, error)

type LocalFn[K comparable, V any] func(ctx context.Context, data []V, k K, page, limit int, a ...any) ([]V, int, error)

func (p *Pagination[K, V]) IsSwitchDB(k K) bool {
	v, _ := p.isSwitch.Load(k)
	return v == true
}

func NewPagination[K comparable, V any](m *MapCache[string, helper.PaginationData[V]], maxNum func() int,
	dbFn DbFn[K, V], localFn LocalFn[K, V], dbKeyFn, localKeyFn func(K, ...any) string,
	batchFetchNum func() int, name string) *Pagination[K, V] {
	if dbKeyFn == nil {
		dbKeyFn = func(k K, a ...any) string {
			s := str.NewBuilder()
			for _, v := range append([]any{k}, a...) {
				s.Sprintf("%v|", v)
			}

			return strings.TrimRight(s.String(), "|")
		}
	}
	if localKeyFn == nil {
		localKeyFn = func(k K, a ...any) string {
			return fmt.Sprintf("%v", k)
		}
	}
	return &Pagination[K, V]{
		MapCache:      m,
		maxNum:        maxNum,
		isSwitch:      safety.NewMap[K, bool](),
		dbFn:          dbFn,
		localFn:       localFn,
		batchFetchNum: batchFetchNum,
		name:          name,
		dbKeyFn:       dbKeyFn,
		localKeyFn:    localKeyFn,
	}
}

func (p *Pagination[K, V]) Pagination(ctx context.Context, timeout time.Duration, k K, page, limit int, a ...any) ([]V, int, error) {
	if is, _ := p.isSwitch.Load(k); is {
		return p.paginationByDB(ctx, timeout, k, page, limit, 0, a...)
	}
	data, total, err := p.paginationByLocal(ctx, timeout, k, page, limit, a...)
	if err != nil {
		if errors.Is(err, switchDb) {
			p.isSwitch.Store(k, true)
			err = nil
			return p.paginationByDB(ctx, timeout, k, page, limit, total, a...)
		}
		return nil, 0, err
	}
	return data, total, err
}

func (p *Pagination[K, V]) paginationByLocal(ctx context.Context, timeout time.Duration, k K, page, limit int, a ...any) ([]V, int, error) {
	key := p.localKeyFn(k)
	data, ok := p.Get(ctx, key)
	if ok {
		if p.increaseUpdate != nil && p.refresh != nil {
			dat, err := p.increaseUpdates(ctx, timeout, data, key, a...)
			if err != nil {
				return nil, 0, err
			}
			if dat.TotalRaw >= p.maxNum() {
				return nil, 0, switchDb
			}
			data = dat
		}
		return p.localFn(ctx, data.Data, k, page, limit, a...)
	}
	p.mux.Lock()
	defer p.mux.Unlock()
	data, ok = p.Get(ctx, key)
	if ok {
		return data.Data, data.TotalRaw, nil
	}
	batchNum := p.batchFetchNum()
	da, totalRaw, err := p.fetchDb(ctx, timeout, k, 1, 0, 0, a...)
	if err != nil {
		return nil, 0, err
	}
	if totalRaw < 1 {
		data.Data = nil
		data.TotalRaw = 0
		p.Set(ctx, key, data)
		return nil, 0, nil
	}
	if totalRaw >= p.maxNum() {
		return nil, totalRaw, switchDb
	}
	totalPage := number.DivideCeil(totalRaw, batchNum)
	for i := 1; i <= totalPage; i++ {
		daa, _, err := p.fetchDb(ctx, timeout, k, i, batchNum, totalRaw, a...)
		if err != nil {
			return nil, 0, err
		}
		da = append(da, daa...)
	}
	data.Data = da
	data.TotalRaw = totalRaw
	p.Set(ctx, key, data)

	return p.localFn(ctx, data.Data, k, page, limit, a...)
}

func (p *Pagination[K, V]) dbGet(ctx context.Context, key string) (helper.PaginationData[V], bool) {
	data, ok := p.Get(ctx, key)
	if ok && p.increaseUpdate != nil && p.increaseUpdate.CycleTime() > p.GetExpireTime(ctx)-p.Ttl(ctx, key) {
		return data, true
	}
	return data, false
}

func (p *Pagination[K, V]) paginationByDB(ctx context.Context, timeout time.Duration, k K, page, limit, totalRaw int, a ...any) ([]V, int, error) {
	key := p.dbKeyFn(k, append([]any{page, limit}, a...)...)
	data, ok := p.dbGet(ctx, key)
	if ok {
		return data.Data, data.TotalRaw, nil
	}
	p.mux.Lock()
	defer p.mux.Unlock()
	data, ok = p.dbGet(ctx, key)
	if ok {
		return data.Data, data.TotalRaw, nil
	}
	dat, total, err := p.fetchDb(ctx, timeout, k, page, limit, totalRaw, a...)
	if err != nil {
		return nil, 0, err
	}
	data.Data, data.TotalRaw = dat, total
	p.Set(ctx, key, data)
	return data.Data, data.TotalRaw, err
}

func (p *Pagination[K, V]) fetchDb(ctx context.Context, timeout time.Duration, k K, page, limit, totalRaw int, a ...any) ([]V, int, error) {
	var data helper.PaginationData[V]
	var err error
	fn := func() {
		da, total, er := p.dbFn(ctx, k, page, limit, totalRaw, a...)
		if er != nil {
			err = er
			return
		}
		data.Data = da
		data.TotalRaw = total
	}
	if timeout > 0 {
		er := helper.RunFnWithTimeout(ctx, timeout, fn, fmt.Sprintf("fetch %s-[%v]-page[%d]-limit[%d] from db fail", p.name, k, page, limit))
		if err == nil && er != nil {
			err = er
		}
	} else {
		fn()
	}
	return data.Data, data.TotalRaw, err
}
