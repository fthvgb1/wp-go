package safety

import (
	"fmt"
	"github/fthvgb1/wp-go/helper"
	"testing"
	"time"
)

func TestSlice_Append(t *testing.T) {
	type args[T any] struct {
		t []T
	}
	type testCase[T any] struct {
		name string
		r    Slice[T]
		args args[T]
	}
	tests := []testCase[int]{
		{
			name: "t1",
			r:    *NewSlice([]int{}),
			args: args[int]{helper.RangeSlice(1, 10, 1)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fn := func() {
				tt.r.Append(tt.args.t...)
			}
			go fn()
			go fn()
			go fn()
			go fn()
			time.Sleep(time.Second)
			fmt.Println(tt.r.Load())
		})
	}
}
