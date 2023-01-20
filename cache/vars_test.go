package cache

import (
	"context"
	"fmt"
	"testing"
	"time"
)

var cc = *NewVarCache(func(a ...any) (int, error) {
	return 1, nil
}, time.Minute)

func TestVarCache_Flush(t *testing.T) {
	type testCase[T any] struct {
		name string
		c    VarCache[T]
	}
	tests := []testCase[int]{
		{
			name: "1",
			c:    cc,
		},
	}
	c := context.Background()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Println(tt.c.GetCache(c, time.Second))
			tt.c.Flush()
			fmt.Println(tt.c.GetCache(c, time.Second))
		})
	}
}

func TestVarCache_IsExpired(t *testing.T) {
	type testCase[T any] struct {
		name string
		c    VarCache[T]
		want bool
	}
	tests := []testCase[int]{
		{
			name: "expired",
			c:    cc,
			want: true,
		},
		{
			name: "not expired",
			c: func() VarCache[int] {
				v := *NewVarCache(func(a ...any) (int, error) {
					return 1, nil
				}, time.Minute)
				_, _ = v.GetCache(context.Background(), time.Second)
				return v
			}(),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.IsExpired(); got != tt.want {
				t.Errorf("IsExpired() = %v, want %v", got, tt.want)
			}
		})
	}
}
