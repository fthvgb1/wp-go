package cache

import (
	"context"
	"fmt"
	"testing"
	"time"
)

var cc = *NewVarCache(NewVarMemoryCache[int](func() time.Duration {
	return time.Minute
}), func(ctx context.Context, a ...any) (int, error) {
	return 1, nil
}, nil, nil)

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
			tt.c.Flush(ctx)
			fmt.Println(tt.c.GetCache(c, time.Second))
		})
	}
}
