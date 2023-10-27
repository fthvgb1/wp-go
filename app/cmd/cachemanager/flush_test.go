package cachemanager

import (
	"context"
	"fmt"
	"github.com/fthvgb1/wp-go/helper/number"
	"github.com/fthvgb1/wp-go/taskPools"
	"testing"
	"time"
)

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

		gets, err := GetMultiple[int]("test", ctx, number.Range(1, 10), time.Second)
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
		FlushMapVal("test", 3, 4)
		fmt.Println(vv.Get(ctx, 3))
		fmt.Println(vv.Get(ctx, 4))
		get, err := Get[int]("test", ctx, 3, time.Second)
		if err != nil {
			t.Fatal(t, "err", err)
		}
		fmt.Println(get, count)
		fmt.Println(vv.Get(ctx, 5))
		FlushAnyVal("test")
		fmt.Println(vv.Get(ctx, 5))
		fmt.Println(vv.Get(ctx, 6))
	})
}
