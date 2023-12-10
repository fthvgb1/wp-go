package safety

import (
	"fmt"
	"github.com/fthvgb1/wp-go/helper/number"
	"github.com/fthvgb1/wp-go/taskPools"
	"testing"
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
			args: args[int]{number.Range(1, 10, 1)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fn := func() {
				switch number.Rand(1, 3) {
				case 1:
					f := tt.r.Load()
					fmt.Println(f)
				case 2:
					tt.r.Append(tt.args.t...)
				case 3:
					/*s := tt.r.Load()
					if len(s) < 1 {
						break
					}
					ii, v := slice.Rand(number.Range(0, len(s)))
					s[ii] = v*/
				}

			}
			p := taskPools.NewPools(20)
			for i := 0; i < 50; i++ {
				p.Execute(fn)
			}
			p.Wait()
			fmt.Println(tt.r.Load())
		})
	}
}
