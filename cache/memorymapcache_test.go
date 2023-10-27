package cache

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"
)

var mm MemoryMapCache[string, string]
var ctx context.Context

var ttt time.Time

func init() {
	ctx = context.Background()
	mm = *NewMemoryMapCache[string, string](3 * time.Second)
	ttt = time.Now()
	mm.Store("aa", mapVal[string]{
		setTime: ttt,
		ver:     1,
		data:    "bb",
	})
	time.Sleep(60 * time.Millisecond)
	mm.Store("cc", mapVal[string]{
		setTime: time.Now(),
		ver:     1,
		data:    "dd",
	})
}

func TestMemoryMapCache_ClearExpired(t *testing.T) {
	type args struct {
		in0    context.Context
		expire time.Duration
	}
	type testCase[K string, V string] struct {
		name string
		m    MemoryMapCache[K, V]
		args args
	}
	tests := []testCase[string, string]{
		{
			name: "t1",
			m:    mm,
			args: args{
				in0:    ctx,
				expire: time.Second,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Println(tt.m)
			tt.m.ClearExpired(tt.args.in0)
			time.Sleep(time.Second)
			fmt.Println(tt.m)
		})
	}
}

func TestMemoryMapCache_Delete(t *testing.T) {
	type args[K comparable] struct {
		in0 context.Context
		key K
	}
	type testCase[K comparable, V any] struct {
		name string
		m    MemoryMapCache[K, V]
		args args[K]
	}
	tests := []testCase[string, string]{
		{
			name: "t1",
			m:    mm,
			args: args[string]{
				in0: ctx,
				key: "aa",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Println(mm.Get(ctx, "aa"))
			tt.m.Del(tt.args.in0, tt.args.key)
			fmt.Println(mm.Get(ctx, "aa"))

		})
	}
}

func TestMemoryMapCache_Flush(t *testing.T) {
	type args struct {
		in0 context.Context
	}
	type testCase[K comparable, V any] struct {
		name string
		m    MemoryMapCache[K, V]
		args args
	}
	tests := []testCase[string, string]{
		{
			name: "t1",
			m:    mm,
			args: args{
				in0: ctx,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.Flush(tt.args.in0)
			mm.Set(ctx, "aa", "xx")
			fmt.Println(mm.Get(ctx, "aa"))
		})
	}
}

func TestMemoryMapCache_Get(t *testing.T) {
	type args[K comparable] struct {
		in0 context.Context
		key K
	}
	type testCase[K comparable, V any] struct {
		name   string
		m      MemoryMapCache[K, V]
		args   args[K]
		wantR  V
		wantOk bool
	}
	tests := []testCase[string, string]{
		{
			name: "t1",
			m:    mm,
			args: args[string]{
				in0: ctx,
				key: "aa",
			},
			wantOk: true,
			wantR:  "bb",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotR, gotOk := tt.m.Get(tt.args.in0, tt.args.key)
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("Get() gotR = %v, want %v", gotR, tt.wantR)
			}
			if gotOk != tt.wantOk {
				t.Errorf("Get() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func TestMemoryMapCache_Set(t *testing.T) {
	type args[K comparable, V any] struct {
		in0 context.Context
		key K
		val V
		in3 time.Duration
	}
	type testCase[K comparable, V any] struct {
		name string
		m    MemoryMapCache[K, V]
		args args[K, V]
	}
	tests := []testCase[string, string]{
		{
			name: "t1",
			m:    mm,
			args: args[string, string]{
				in0: ctx,
				key: "ee",
				val: "ff",
				in3: time.Second,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.Set(tt.args.in0, tt.args.key, tt.args.val)
			fmt.Println(tt.m.Get(ctx, tt.args.key))
		})
	}
}

func TestMemoryMapCache_Ttl(t *testing.T) {
	type args[K comparable] struct {
		in0    context.Context
		key    K
		expire time.Duration
	}
	type testCase[K comparable, V any] struct {
		name string
		m    MemoryMapCache[K, V]
		args args[K]
		want time.Duration
	}
	tt := time.Now()
	tests := []testCase[string, string]{
		{
			name: "t1",
			m:    mm,
			args: args[string]{key: "aa", in0: ctx, expire: time.Second * 4},
			want: 4*time.Second - tt.Sub(ttt),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.Ttl(tt.args.in0, tt.args.key); got != tt.want {
				t.Errorf("Ttl() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemoryMapCache_Ver(t *testing.T) {
	type args[K comparable] struct {
		in0 context.Context
		key K
	}
	type testCase[K comparable, V any] struct {
		name string
		m    MemoryMapCache[K, V]
		args args[K]
		want int
	}
	mm.Set(ctx, "aa", "ff")
	tests := []testCase[string, string]{
		{
			name: "t1",
			m:    mm,
			args: args[string]{ctx, "aa"},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.Ver(tt.args.in0, tt.args.key); got != tt.want {
				t.Errorf("Ver() = %v, want %v", got, tt.want)
			}
		})
	}
}
