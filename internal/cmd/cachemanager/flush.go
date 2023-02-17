package cachemanager

import (
	"context"
	"github.com/fthvgb1/wp-go/cache"
	"time"
)

var ctx = context.Background()

type flush interface {
	Flush(ctx context.Context)
}

type clear interface {
	ClearExpired(ctx context.Context)
}

var clears []clear

var flushes []flush

func Flush() {
	for _, f := range flushes {
		f.Flush(ctx)
	}
}

func MapCacheBy[K comparable, V any](fn func(...any) (V, error), expireTime time.Duration) *cache.MapCache[K, V] {
	m := cache.NewMemoryMapCacheByFn[K, V](fn, expireTime)
	FlushPush(m)
	ClearPush(m)
	return m
}
func MapBatchCacheBy[K comparable, V any](fn func(...any) (map[K]V, error), expireTime time.Duration) *cache.MapCache[K, V] {
	m := cache.NewMemoryMapCacheByBatchFn[K, V](fn, expireTime)
	FlushPush(m)
	ClearPush(m)
	return m
}

func FlushPush(f ...flush) {
	flushes = append(flushes, f...)
}
func ClearPush(c ...clear) {
	clears = append(clears, c...)
}

func ClearExpired() {
	for _, c := range clears {
		c.ClearExpired(ctx)
	}
}
