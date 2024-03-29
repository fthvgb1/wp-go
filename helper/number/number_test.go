package number

import (
	"fmt"
	"github.com/fthvgb1/wp-go/taskPools"
	"golang.org/x/exp/constraints"
	"reflect"
	"testing"
)

func TestRange(t *testing.T) {
	type args struct {
		start int
		end   int
		step  int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "t1",
			args: args{
				start: 1,
				end:   5,
				step:  1,
			},
			want: []int{1, 2, 3, 4, 5},
		},
		{
			name: "t2",
			args: args{
				start: 0,
				end:   5,
				step:  2,
			},
			want: []int{0, 2, 4},
		},
		{
			name: "t3",
			args: args{
				start: 1,
				end:   11,
				step:  3,
			},
			want: []int{1, 4, 7, 10},
		},
		{
			name: "t4",
			args: args{
				start: 0,
				end:   -5,
				step:  -1,
			},
			want: []int{0, -1, -2, -3, -4, -5},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Range(tt.args.start, tt.args.end, tt.args.step); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Range() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMin(t *testing.T) {
	type args struct {
		a []int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "t1",
			args: args{a: []int{1, 2, 3}},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Min(tt.args.a...); got != tt.want {
				t.Errorf("Min() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMax(t *testing.T) {
	type args struct {
		a []int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "t1",
			args: args{a: []int{1, 2, 3}},
			want: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Max(tt.args.a...); got != tt.want {
				t.Errorf("Max() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSum(t *testing.T) {
	type args struct {
		a []int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "t1",
			args: args{a: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
			want: 55,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Sum(tt.args.a...); got != tt.want {
				t.Errorf("Sum() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToString(t *testing.T) {
	type args struct {
		n float64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "t1",
			args: args{n: 111},
			want: "111",
		},
		{
			name: "t2",
			args: args{n: 111.222222},
			want: "111.222222",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToString(tt.args.n); got != tt.want {
				t.Errorf("NumberToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRand(t *testing.T) {
	type args struct {
		start int
		end   int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "t1",
			args: args{
				start: 1,
				end:   2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i := 0; i < 100; i++ {
				got := Rand(tt.args.start, tt.args.end)
				if got > tt.args.end || got < tt.args.start {
					t.Errorf("RandNum() = %v, range error", got)
				}
				fmt.Println(got)
			}
		})
	}
}

func TestAbs(t *testing.T) {
	type args[T constraints.Integer | constraints.Float] struct {
		n T
	}
	type testCase[T constraints.Integer | constraints.Float] struct {
		name string
		args args[T]
		want T
	}
	tests := []testCase[int]{
		{
			name: "t1",
			args: args[int]{-1},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Abs(tt.args.n); got != tt.want {
				t.Errorf("Abs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalTotalPage(t *testing.T) {
	type args[T constraints.Integer] struct {
		totalRows T
		size      T
	}
	type testCase[T constraints.Integer] struct {
		name string
		args args[T]
		want T
	}
	tests := []testCase[int]{
		{
			name: "t1",
			args: args[int]{5, 2},
			want: 3,
		},
		{
			name: "t1",
			args: args[int]{4, 2},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DivideCeil(tt.args.totalRows, tt.args.size); got != tt.want {
				t.Errorf("DivideCeil() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCounters(t *testing.T) {
	type testCase[T constraints.Integer] struct {
		name string
		want func() T
	}
	var c = 0
	tests := []testCase[int]{
		{
			name: "t1",
			want: func() int {
				c++
				return c
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Counters[int]()
			if !reflect.DeepEqual(got(), tt.want()) {
				t.Errorf("Counters() = %v, want %v", got(), tt.want())
			}
			got()
			got()
			got()
			p := taskPools.NewPools(6)
			for i := 0; i < 50; i++ {
				p.Execute(func() {
					got()
				})
			}
			p.Wait()
			fmt.Println("got ", got())
		})

	}
}
