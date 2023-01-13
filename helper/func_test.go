package helper

import (
	"fmt"
	"reflect"
	"testing"
)

type x struct {
	Id uint64
}

func c(x []*x) (r []uint64) {
	for i := 0; i < len(x); i++ {
		r = append(r, x[i].Id)
	}
	return
}

func getX() (r []*x) {
	for i := 0; i < 10; i++ {
		r = append(r, &x{
			uint64(i),
		})
	}
	return
}

func BenchmarkOr(b *testing.B) {
	y := getX()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c(y)
	}
}

func BenchmarkStructColumn(b *testing.B) {
	y := getX()
	fmt.Println(y)
	b.ResetTimer()
	//b.N = 2
	for i := 0; i < 1; i++ {
		StructColumn[int, *x](y, "Id")
	}
}

func TestStructColumn(t *testing.T) {
	type args struct {
		arr   []x
		field string
	}

	tests := []struct {
		name  string
		args  args
		wantR []uint64
	}{
		{name: "test1", args: args{
			arr: []x{
				{Id: 1},
				{2},
				{4},
				{6},
			},
			field: "Id",
		}, wantR: []uint64{1, 2, 4, 6}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotR := StructColumn[uint64, x](tt.args.arr, tt.args.field); !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("StructColumn() = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestToInterface(t *testing.T) {
	type args struct {
		v int
	}
	tests := []struct {
		name string
		args args
		want any
	}{
		{
			name: "t1",
			args: args{v: 1},
			want: any(1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToAny(tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToAny() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRandNum(t *testing.T) {
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
				got := RandNum(tt.args.start, tt.args.end)
				if got > tt.args.end || got < tt.args.start {
					t.Errorf("RandNum() = %v, range error", got)
				}
			}
		})
	}
}

func TestSimpleSort(t *testing.T) {
	type xy struct {
		x int
		y int
	}

	type args struct {
		arr []xy
		fn  func(i, j xy) bool
	}
	tests := []struct {
		name  string
		args  args
		wantR []xy
	}{
		{
			name: "t1",
			args: args{
				arr: []xy{
					{1, 2},
					{3, 4},
					{1, 3},
					{2, 1},
					{1, 6},
				},
				fn: func(i, j xy) bool {
					if i.x < j.x {
						return true
					}
					if i.x == j.x && i.y > i.y {
						return true
					}
					return false
				},
			},
			wantR: []xy{
				{1, 2},
				{1, 3},
				{1, 6},
				{2, 1},
				{3, 4},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if SimpleSort(tt.args.arr, tt.args.fn); !reflect.DeepEqual(tt.args.arr, tt.wantR) {
				t.Errorf("SimpleSort() = %v, want %v", tt.args.arr, tt.wantR)
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
			args: args{a: RangeSlice(1, 10, 1)},
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

func TestNumberToString(t *testing.T) {
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
			if got := NumberToString(tt.args.n); got != tt.want {
				t.Errorf("NumberToString() = %v, want %v", got, tt.want)
			}
		})
	}
}
