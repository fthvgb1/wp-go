package slice

import (
	"fmt"
	"github.com/fthvgb1/wp-go/helper/number"
	"reflect"
	"testing"
)

func TestSplice(t *testing.T) {
	type args[T int] struct {
		a           *[]T
		offset      int
		length      int
		replacement []T
	}
	type testCase[T int] struct {
		name string
		args args[T]
		want []T
	}
	a := number.Range(1, 10, 1)
	b := number.Range(1, 10, 1)
	c := number.Range(1, 10, 1)
	d := number.Range(1, 10, 1)
	tests := []testCase[int]{
		{
			name: "t1",
			args: args[int]{
				a:           &a,
				offset:      3,
				length:      2,
				replacement: nil,
			},
			want: []int{4, 5},
		},
		{
			name: "t2",
			args: args[int]{
				a:           &b,
				offset:      3,
				length:      2,
				replacement: []int{11, 12, 15},
			},
			want: []int{4, 5},
		},
		{
			name: "t3",
			args: args[int]{
				a:           &c,
				offset:      -1,
				length:      2,
				replacement: nil, //[]int{11, 12, 15},
			},
			want: []int{10},
		},
		{
			name: "t4",
			args: args[int]{
				a:           &d,
				offset:      -3,
				length:      5,
				replacement: []int{11, 12, 15},
			},
			want: []int{8, 9, 10},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Splice(tt.args.a, tt.args.offset, tt.args.length, tt.args.replacement); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Splice() = %v, want %v", got, tt.want)
			}
		})
	}
	fmt.Println(a)
	fmt.Println(b)
	fmt.Println(c)
	fmt.Println(d)
}
