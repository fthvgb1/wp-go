package cache

import (
	"context"
	"fmt"
	"github.com/fthvgb1/wp-go/helper/slice"
	"reflect"
	"strings"
	"testing"
	"time"
)

var ca MapCache[string, string]
var fn func(a ...any) (string, error)
var batchFn func(a ...any) (map[string]string, error)
var ct context.Context

func init() {
	fn = func(a ...any) (string, error) {
		aa := a[1].(string)
		return strings.Repeat(aa, 2), nil
	}
	ct = context.Background()
	batchFn = func(a ...any) (map[string]string, error) {
		arr := a[1].([]string)
		return slice.SimpleToMap(arr, func(t string) string {
			return strings.Repeat(t, 2)
		}), nil
	}
	ca = *NewMemoryMapCacheByFn[string, string](fn, time.Second*2)
	ca.SetCacheBatchFn(batchFn)
	_, _ = ca.GetCache(ct, "aa", time.Second, ct, "aa")
	_, _ = ca.GetCache(ct, "bb", time.Second, ct, "bb")
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
	ca := *NewMemoryMapCacheByFn[string, string](fn, time.Second)
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
				params:  []any{ct, []string{"xx", "oo"}},
			},
			want:    []string{"xxxx", "oooo"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
		fn func(...any) (map[K]V, error)
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
	type args[V any] struct {
		fn func(...any) (V, error)
	}
	type testCase[K comparable, V any] struct {
		name string
		m    MapCache[K, V]
		args args[V]
	}
	tests := []testCase[string, string]{
		{
			name: "t1",
			m:    ca,
			args: args[string]{fn: fn},
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
			want: ca.expireTime - tx.Sub(txx),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Printf("过期时间=%v \nttl=%v \n当前时间   =%v\n最后设置时间=%v\n当时时间-最后设置时间=%v ", ca.expireTime, ca.Ttl(ct, "aa"), tx, txx, tx.Sub(txx))
			if got := tt.m.Ttl(tt.args.ct, tt.args.k); got != tt.want {
				t.Errorf("Ttl() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMapCache_setCacheFn(t *testing.T) {
	type args[K comparable, V any] struct {
		fn func(...any) (map[K]V, error)
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
			tt.m.setCacheFn(tt.args.fn)
			fmt.Println(ca.GetCache(ct, "xx", time.Second, ct, "xx"))
		})
	}
}
