package cache

import (
	"context"
	"fmt"
	"github.com/fthvgb1/wp-go/helper/slice"
	"github.com/fthvgb1/wp-go/taskPools"
	"reflect"
	"strings"
	"testing"
	"time"
)

var ca MapCache[string, string]
var fn MapSingleFn[string, string]
var batchFn MapBatchFn[string, string]
var ct context.Context

func init() {
	fn = func(ctx context.Context, aa string, a ...any) (string, error) {
		return strings.Repeat(aa, 2), nil
	}
	ct = context.Background()
	batchFn = func(ctx context.Context, arr []string, a ...any) (map[string]string, error) {
		fmt.Println(a)
		return slice.FilterAndToMap(arr, func(t string) (string, string, bool) {
			return t, strings.Repeat(t, 2), true
		}), nil
	}

}
func TestMapCache_ClearExpired(t *testing.T) {
	type args struct {
		ct context.Context
	}
	type testCase[K comparable, V any] struct {
		name string
		m    MapCache[K, V]
		args args
	}
	tests := []testCase[string, string]{
		{
			name: "t1",
			m:    ca,
			args: args{
				ct: ct,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Println(ca.Get(ct, "aa"))
			fmt.Println(ca.Get(ct, "bb"))
			time.Sleep(time.Second * 3)
			tt.m.ClearExpired(tt.args.ct)
			fmt.Println(ca.Get(ct, "bb"))
		})
	}
}

func TestMapCache_Flush(t *testing.T) {
	type args struct {
		ct context.Context
	}
	type testCase[K comparable, V any] struct {
		name string
		m    MapCache[K, V]
		args args
	}
	ca := *NewMapCache[string, string](NewMemoryMapCache[string, string](func() time.Duration {
		return time.Second
	}), fn, nil)
	_, _ = ca.GetCache(ct, "aa", time.Second, ct, "aa")
	tests := []testCase[string, string]{
		{
			name: "t1",
			m:    ca,
			args: args{
				ct,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Println(ca.Get(ct, "aa"))
			tt.m.Flush(tt.args.ct)
			fmt.Println(ca.Get(ct, "aa"))
		})
	}
}

func TestMapCache_Get(t *testing.T) {
	type args[K comparable] struct {
		ct context.Context
		k  K
	}
	type testCase[K comparable, V any] struct {
		name  string
		m     MapCache[K, V]
		args  args[K]
		want  V
		want1 bool
	}
	tests := []testCase[string, string]{
		{
			name:  "t1",
			m:     ca,
			args:  args[string]{ct, "aa"},
			want:  "aaaa",
			want1: true,
		},
		{
			name:  "t2",
			m:     ca,
			args:  args[string]{ct, "cc"},
			want:  "",
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.m.Get(tt.args.ct, tt.args.k)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Get() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestMapCache_GetCache(t *testing.T) {
	type args[K comparable] struct {
		c       context.Context
		key     K
		timeout time.Duration
		params  []any
	}
	type testCase[K comparable, V any] struct {
		name    string
		m       MapCache[K, V]
		args    args[K]
		want    V
		wantErr bool
	}
	tests := []testCase[string, string]{
		{
			name:    "t1",
			m:       ca,
			args:    args[string]{c: ct, key: "xx", timeout: time.Second, params: []any{ct, "xx"}},
			want:    "xxxx",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.GetCache(tt.args.c, tt.args.key, tt.args.timeout, tt.args.params...)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCache() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCache() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMapCache_GetCacheBatch(t *testing.T) {
	type args[K comparable] struct {
		c       context.Context
		key     []K
		timeout time.Duration
		params  []any
	}
	type testCase[K comparable, V any] struct {
		name    string
		m       MapCache[K, V]
		args    args[K]
		want    []V
		wantErr bool
	}
	tests := []testCase[string, string]{
		{
			name: "t1",
			m:    ca,
			args: args[string]{
				c:       ct,
				key:     []string{"xx", "oo"},
				timeout: time.Second,
				params:  []any{ct, []string{"xx", "oo", "aa"}},
			},
			want:    []string{"xxxx", "oooo", "aaaa"},
			wantErr: false,
		},
	}
	time.Sleep(2 * time.Second)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := taskPools.NewPools(10)
			for i := 0; i < 800000; i++ {
				p.Execute(func() {
					c := context.Background()
					//time.Sleep(time.Millisecond * number.Rand[time.Duration](200, 400))
					a, err := ca.GetCacheBatch(c, []string{"xx", "oo", "aa"}, time.Hour, c, []string{"xx", "oo", "aa"})
					if err != nil {
						panic(err)
						return
					}

					if a[0] == "xxxx" && a[1] == "oooo" && a[2] == "aaaa" {

					} else {
						fmt.Println(a)
						panic("xxx")
					}
					//fmt.Println(x)
				})
			}
			p.Wait()
			got, err := tt.m.GetCacheBatch(tt.args.c, tt.args.key, tt.args.timeout, tt.args.params...)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCacheBatch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCacheBatch() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMapCache_GetLastSetTime(t *testing.T) {
	type args[K comparable] struct {
		ct context.Context
		k  K
	}
	type testCase[K comparable, V any] struct {
		name  string
		m     MapCache[K, V]
		args  args[K]
		wantT time.Time
	}
	tests := []testCase[string, string]{
		{
			name:  "t1",
			m:     ca,
			args:  args[string]{ct, "aa"},
			wantT: ttt,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotT := tt.m.GetLastSetTime(tt.args.ct, tt.args.k); !reflect.DeepEqual(gotT, tt.wantT) {
				t.Errorf("GetLastSetTime() = %v, want %v", gotT, tt.wantT)
			}
		})
	}
}

func TestMapCache_Set(t *testing.T) {
	type args[K comparable, V any] struct {
		ct context.Context
		k  K
		v  V
	}
	type testCase[K comparable, V any] struct {
		name string
		m    MapCache[K, V]
		args args[K, V]
	}
	tests := []testCase[string, string]{
		{
			name: "t1",
			m:    ca,
			args: args[string, string]{
				ct, "xx", "yy",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Println(ca.Get(ct, "xx"))
			tt.m.Set(tt.args.ct, tt.args.k, tt.args.v)
			fmt.Println(ca.Get(ct, "xx"))
		})
	}
}

func TestMapCache_SetCacheBatchFn(t *testing.T) {
	type args[K comparable, V any] struct {
		fn MapBatchFn[K, V]
	}
	type testCase[K comparable, V any] struct {
		name string
		m    MapCache[K, V]
		args args[K, V]
	}
	tests := []testCase[string, string]{
		{
			name: "t1",
			m:    ca,
			args: args[string, string]{batchFn},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.SetCacheBatchFn(tt.args.fn)
		})
	}
}

func TestMapCache_SetCacheFunc(t *testing.T) {
	type args[K comparable, V any] struct {
		fn MapSingleFn[K, V]
	}
	type testCase[K comparable, V any] struct {
		name string
		m    MapCache[K, V]
		args args[K, V]
	}
	tests := []testCase[string, string]{
		{
			name: "t1",
			m:    ca,
			args: args[string, string]{fn: fn},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.SetCacheFunc(tt.args.fn)
		})
	}
}

func TestMapCache_Ttl(t *testing.T) {
	type args[K comparable] struct {
		ct context.Context
		k  K
	}
	type testCase[K comparable, V any] struct {
		name string
		m    MapCache[K, V]
		args args[K]
		want time.Duration
	}
	tx := time.Now()
	txx := ca.GetLastSetTime(ct, "aa")
	tests := []testCase[string, string]{
		{
			name: "t1",
			m:    ca,
			args: args[string]{ct, "aa"},
			want: ca.GetExpireTime(ct) - tx.Sub(txx),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Printf("过期时间=%v \nttl=%v \n当前时间   =%v\n最后设置时间=%v\n当时时间-最后设置时间=%v ", ca.GetExpireTime(ct), ca.Ttl(ct, "aa"), tx, txx, tx.Sub(txx))
			if got := tt.m.Ttl(tt.args.ct, tt.args.k); got != tt.want {
				t.Errorf("Ttl() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMapCache_setCacheFn(t *testing.T) {
	type args[K comparable, V any] struct {
		fn MapBatchFn[K, V]
	}
	type testCase[K comparable, V any] struct {
		name string
		m    MapCache[K, V]
		args args[K, V]
	}
	tests := []testCase[string, string]{
		{
			name: "t1",
			m:    ca,
			args: args[string, string]{batchFn},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ca.cacheFunc = nil
			tt.m.setDefaultCacheFn(tt.args.fn)
			fmt.Println(ca.GetCache(ct, "xx", time.Second, ct, "xx"))
		})
	}
}
