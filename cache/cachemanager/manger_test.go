package cachemanager

import (
	"context"
	"fmt"
	"github.com/fthvgb1/wp-go/helper/number"
	"github.com/fthvgb1/wp-go/helper/slice"
	str "github.com/fthvgb1/wp-go/helper/strings"
	"github.com/fthvgb1/wp-go/taskPools"
	"reflect"
	"strings"
	"testing"
	"time"
)

var ctx = context.Background()

func TestFlushMapVal(t *testing.T) {
	_ = number.Range(1, 5, 0)
	t.Run("t1", func(t *testing.T) {
		count := 0
		vv := NewMemoryMapCache(func(ctx2 context.Context, ks []int, a ...any) (map[int]int, error) {
			r := make(map[int]int)
			for _, k := range ks {
				r[k] = k * k
			}
			count++
			return r, nil
		}, nil, time.Second, "test")

		gets, err := GetBatchBy[int]("test", ctx, number.Range(1, 10), time.Second)
		if err != nil {
			t.Fatal(t, "err:", err)
		}
		p := taskPools.NewPools(10)
		for i := 0; i < 20; i++ {
			i := i
			p.Execute(func() {
				if i%2 == 0 {
					vv.Get(ctx, 5)
				} else {
					vv.Set(ctx, i, i)
				}
			})
		}
		p.Wait()
		fmt.Println(gets, count)
		DelMapCacheVal("test", 3, 4)
		fmt.Println(vv.Get(ctx, 3))
		fmt.Println(vv.Get(ctx, 4))
		get, err := GetBy[int]("test", ctx, 3, time.Second)
		if err != nil {
			t.Fatal(t, "err", err)
		}
		fmt.Println(get, count)
		fmt.Println(vv.Get(ctx, 5))
		Flushes(ctx, "test")
		fmt.Println(vv.Get(ctx, 5))
		fmt.Println(vv.Get(ctx, 6))
		//fmt.Println(GetVarCache("test"))
	})
}

func TestSetExpireTime(t *testing.T) {
	t.Run("t1", func(t *testing.T) {
		c := NewMemoryMapCache[string, string](func(ctx2 context.Context, strings []string, a ...any) (map[string]string, error) {
			return slice.ToMap(strings, func(v string) (string, string) {
				return v, str.Join(v, "__", v)
			}, false), nil
		}, nil, time.Second, "xx")
		c.Set(ctx, "xx", "yy")
		fmt.Println(c.Get(ctx, "xx"))
		time.Sleep(time.Second)
		fmt.Println(c.Get(ctx, "xx"))
		ChangeExpireTime(3*time.Second, true, "xx")
		c.Set(ctx, "xx", "yyy")
		time.Sleep(time.Second)
		fmt.Println(c.Get(ctx, "xx"))
		time.Sleep(3 * time.Second)
		fmt.Println(c.Get(ctx, "xx"))
		cc, _ := GetMapCache[string, string]("xx")
		fmt.Println(reflect.DeepEqual(c, cc))
		cc.Set(ctx, "fff", "xxxx")
		cc.Set(ctx, "ffx", "eex")
		cc.Set(ctx, "ww", "vv")
		m, err := cc.GetBatchToMap(ctx, []string{"fff", "ffx", "ww", "kkkk"}, time.Second)
		fmt.Println(m, err)
		fmt.Println(GetBatchByToMap[string]("xx", ctx, []string{"fff", "ffx", "ww", "kkkk"}, time.Second))
		v := NewVarMemoryCache(func(ct context.Context, a ...any) (string, error) {
			return "ssss", nil
		}, 3*time.Second, "ff")
		vv, _ := GetVarCache[string]("ff")
		fmt.Println(reflect.DeepEqual(v, vv))
	})
}

func TestSetMapCache(t *testing.T) {
	t.Run("t1", func(t *testing.T) {
		x := NewMemoryMapCache(nil, func(ctx2 context.Context, k string, a ...any) (string, error) {
			fmt.Println("memory cache")
			return strings.Repeat(k, 2), nil
		}, time.Hour, "test")
		fmt.Println(GetBy[string]("test", ctx, "test", time.Second))

		NewMapCache[string, string](xx[string, string]{m: map[string]string{}}, nil, func(ctx2 context.Context, k string, a ...any) (string, error) {
			fmt.Println("other cache drives. eg: redis,file.....")
			return strings.Repeat(k, 2), nil
		}, "test", time.Hour)

		if err := SetMapCache("kkk", x); err != nil {
			t.Errorf("SetMapCache() error = %v, wantErr %v", err, nil)
		}
		fmt.Println(GetBy[string]("test", ctx, "test", time.Second))
	})
}

type xx[K comparable, V any] struct {
	m map[K]V
}

func (x xx[K, V]) Get(ctx context.Context, key K) (V, bool) {
	v, ok := x.m[key]
	return v, ok
}

func (x xx[K, V]) Set(ctx context.Context, key K, val V) {
	x.m[key] = val
}

func (x xx[K, V]) GetExpireTime(ctx context.Context) time.Duration {
	//TODO implement me
	panic("implement me")
}

func (x xx[K, V]) Ttl(ctx context.Context, key K) time.Duration {
	//TODO implement me
	panic("implement me")
}

func (x xx[K, V]) Flush(ctx context.Context) {
	//TODO implement me
	panic("implement me")
}

func (x xx[K, V]) Del(ctx context.Context, key ...K) {
	//TODO implement me
	panic("implement me")
}

func (x xx[K, V]) ClearExpired(ctx context.Context) {
	//TODO implement me
	panic("implement me")
}

func TestSetVarCache(t *testing.T) {
	t.Run("t1", func(t *testing.T) {
		bak := NewVarMemoryCache(func(ctx2 context.Context, a ...any) (string, error) {
			fmt.Println("memory cache")
			return "xxx", nil
		}, time.Hour, "test")
		fmt.Println(GetVarVal[string]("test", ctx, time.Second))
		NewVarCache[string](oo[string]{}, func(ctx2 context.Context, a ...any) (string, error) {
			fmt.Println("other cache drives. eg: redis,file.....")
			return "ooo", nil
		}, "test")
		if err := SetVarCache("xx", bak); err != nil {
			t.Errorf("SetVarCache() error = %v, wantErr %v", err, nil)
		}
		fmt.Println(GetVarVal[string]("test", ctx, time.Second))
	})
}

type oo[T any] struct {
	val T
}

func (o oo[T]) Get(ctx context.Context) (T, bool) {
	return o.val, false
}

func (o oo[T]) Set(ctx context.Context, v T) {
	o.val = v
}

func (o oo[T]) Flush(ctx context.Context) {
	//TODO implement me
	panic("implement me")
}

func (o oo[T]) GetLastSetTime(ctx context.Context) time.Time {
	//TODO implement me
	panic("implement me")
}
