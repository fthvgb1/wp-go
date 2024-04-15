package cachemanager

import (
	"context"
	"errors"
	"github.com/fthvgb1/wp-go/cache"
	"github.com/fthvgb1/wp-go/cache/reload"
	"github.com/fthvgb1/wp-go/helper"
	str "github.com/fthvgb1/wp-go/helper/strings"
	"time"
)

func NewPaginationCache[K comparable, V any](m *cache.MapCache[string, helper.PaginationData[V]], maxNum int,
	dbFn cache.DbFn[K, V], localFn cache.LocalFn[K, V], dbKeyFn, localKeyFn func(K, ...any) string, fetchNum int, name string, a ...any) *cache.Pagination[K, V] {
	fn := helper.ParseArgs([]func() int(nil), a...)
	var ma, fet func() int
	if len(fn) > 0 {
		ma = fn[0]
		if len(fn) > 1 {
			fet = fn[1]
		}
	}
	if ma == nil {
		ma = reload.BuildFnVal(str.Join("paginationCache-", name, "-maxNum"), maxNum, nil)
	}
	if fet == nil {
		fet = reload.BuildFnVal(str.Join("paginationCache-", name, "-fetchNum"), fetchNum, nil)
	}
	p := cache.NewPagination(m, ma, dbFn, localFn, dbKeyFn, localKeyFn, fet, name)
	mapCache.Store(name, p)
	return p
}

func GetPaginationCache[K comparable, V any](name string) (*cache.Pagination[K, V], bool) {
	v, err := getPagination[K, V](name)
	return v, err == nil
}

func Pagination[V any, K comparable](name string, ctx context.Context, timeout time.Duration, k K, page, limit int, a ...any) ([]V, int, error) {
	v, err := getPagination[K, V](name)
	if err != nil {
		return nil, 0, err
	}
	return v.Pagination(ctx, timeout, k, page, limit, a...)
}

func getPagination[K comparable, T any](name string) (*cache.Pagination[K, T], error) {
	m, ok := mapCache.Load(name)
	if !ok {
		return nil, errors.New(str.Join("cache ", name, " doesn't exist"))
	}
	vk, ok := m.(*cache.Pagination[K, T])
	if !ok {
		return nil, errors.New(str.Join("cache ", name, " type error"))
	}
	return vk, nil
}
