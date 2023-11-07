package cachemanager

import (
	"context"
	"fmt"
	"github.com/fthvgb1/wp-go/helper/number"
	"github.com/fthvgb1/wp-go/helper/slice"
	str "github.com/fthvgb1/wp-go/helper/strings"
	"github.com/fthvgb1/wp-go/taskPools"
	"reflect"
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
		ChangeExpireTime(3*time.Second, "xx")
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
		fmt.Println(GetMultipleToMap[string]("xx", ctx, []string{"fff", "ffx", "ww", "kkkk"}, time.Second))
		v := NewVarMemoryCache(func(ct context.Context, a ...any) (string, error) {
			return "ssss", nil
		}, 3*time.Second, "ff")
		vv, _ := GetVarCache[string]("ff")
		fmt.Println(reflect.DeepEqual(v, vv))
	})
}
