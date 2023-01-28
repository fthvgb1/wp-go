package slice

import (
	"fmt"
	"github.com/fthvgb1/wp-go/helper/number"
	"reflect"
	"testing"
)

func TestPagination(t *testing.T) {
	arr := number.Range[int](1, 10, 1)
	type args struct {
		arr      []int
		page     int
		pageSize int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "t1",
			args: args{
				arr:      arr,
				page:     1,
				pageSize: 2,
			},
			want: number.Range[int](1, 2, 1),
		}, {
			name: "t2",
			args: args{
				arr:      arr,
				page:     2,
				pageSize: 2,
			},
			want: number.Range[int](3, 4, 1),
		}, {
			name: "t3",
			args: args{
				arr:      arr,
				page:     4,
				pageSize: 3,
			},
			want: []int{10},
		}, {
			name: "t4",
			args: args{
				arr:      arr,
				page:     5,
				pageSize: 3,
			},
			want: []int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Pagination(tt.args.arr, tt.args.page, tt.args.pageSize); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Pagination() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReduce(t *testing.T) {
	type args struct {
		arr []int
		fn  func(int, int) int
		r   int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "t1",
			args: args{arr: number.Range(1, 10, 1), fn: func(i int, i2 int) int {
				return i + i2
			}},
			want: 55,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Reduce(tt.args.arr, tt.args.fn, tt.args.r); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Reduce() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilter(t *testing.T) {
	type args struct {
		arr []int
		fn  func(int) bool
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "t1",
			args: args{arr: number.Range(1, 10, 1), fn: func(i int) bool {
				return i > 4
			}},
			want: number.Range(5, 10, 1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Filter(tt.args.arr, tt.args.fn); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Filter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMap(t *testing.T) {
	type args struct {
		arr []int8
		fn  func(int8) int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "t1",
			args: args{
				arr: number.Range[int8](1, 10, 1),
				fn: func(i int8) int {
					return int(i)
				},
			},
			want: number.Range(1, 10, 1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Map(tt.args.arr, tt.args.fn); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Map() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReverse(t *testing.T) {
	type args struct {
		arr []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "t1",
			args: args{arr: number.Range(1, 10, 1)},
			want: number.Range(10, 1, -1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Reverse(tt.args.arr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Reverse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSlice(t *testing.T) {
	type args struct {
		arr    []int
		offset int
		length int
	}
	tests := []struct {
		name  string
		args  args
		wantR []int
	}{
		{
			name: "t1",
			args: args{
				arr:    number.Range(1, 10, 1),
				offset: 3,
				length: 2,
			},
			wantR: number.Range(4, 5, 1),
		},
		{
			name: "t2",
			args: args{
				arr:    number.Range(1, 10, 1),
				offset: 3,
				length: 0,
			},
			wantR: number.Range(4, 10, 1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotR := Slice(tt.args.arr, tt.args.offset, tt.args.length); !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("Slice() = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestComb(t *testing.T) {
	type args struct {
		arr []int
		m   int
	}
	tests := []struct {
		name  string
		args  args
		wantR [][]int
	}{
		{
			name: "t1",
			args: args{
				arr: number.Range(1, 5, 1),
				m:   2,
			},
			wantR: [][]int{
				{1, 2},
				{1, 3},
				{1, 4},
				{1, 5},
				{2, 3},
				{2, 4},
				{2, 5},
				{3, 4},
				{3, 5},
				{4, 5},
			},
		},
		{
			name: "t2",
			args: args{
				arr: number.Range(1, 5, 1),
				m:   3,
			},
			wantR: [][]int{
				{1, 2, 3},
				{1, 2, 4},
				{1, 2, 5},
				{1, 3, 4},
				{1, 3, 5},
				{1, 4, 5},
				{2, 3, 4},
				{2, 3, 5},
				{2, 4, 5},
				{3, 4, 5},
			},
		},
		{
			name: "t3",
			args: args{
				arr: number.Range(1, 5, 1),
				m:   4,
			},
			wantR: [][]int{
				{1, 2, 3, 4},
				{1, 2, 3, 5},
				{1, 2, 4, 5},
				{1, 3, 4, 5},
				{2, 3, 4, 5},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotR := Comb(tt.args.arr, tt.args.m); !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("Comb() = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestChunk(t *testing.T) {
	type args struct {
		arr  []int
		size int
	}
	tests := []struct {
		name string
		args args
		want [][]int
	}{
		{
			name: "t1",
			args: args{
				arr:  number.Range(1, 7, 1),
				size: 2,
			},
			want: [][]int{{1, 2}, {3, 4}, {5, 6}, {7}},
		},
		{
			name: "t2",
			args: args{
				arr:  number.Range(1, 8, 1),
				size: 2,
			},
			want: [][]int{{1, 2}, {3, 4}, {5, 6}, {7, 8}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Chunk(tt.args.arr, tt.args.size); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Chunk() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSimpleSliceToMap(t *testing.T) {
	type args struct {
		arr []int
		fn  func(int) int
	}
	tests := []struct {
		name string
		args args
		want map[int]int
	}{
		{
			name: "t1",
			args: args{arr: []int{1, 2, 3}, fn: func(i int) int {
				return i
			}},
			want: map[int]int{1: 1, 2: 2, 3: 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SimpleToMap(tt.args.arr, tt.args.fn); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SimpleToMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToMap(t *testing.T) {
	type ss struct {
		id int
		v  string
	}
	type args struct {
		arr         []ss
		fn          func(ss) (int, ss)
		isCoverPrev bool
	}
	tests := []struct {
		name string
		args args
		want map[int]ss
	}{
		{
			name: "t1",
			args: args{
				arr: []ss{{1, "k1"}, {2, "v2"}, {2, "v3"}},
				fn: func(s ss) (int, ss) {
					return s.id, s
				},
				isCoverPrev: true,
			},
			want: map[int]ss{1: {1, "k1"}, 2: {2, "v3"}},
		}, {
			name: "t2",
			args: args{
				arr: []ss{{1, "k1"}, {2, "v2"}, {2, "v3"}},
				fn: func(s ss) (int, ss) {
					return s.id, s
				},
				isCoverPrev: false,
			},
			want: map[int]ss{1: {1, "k1"}, 2: {2, "v2"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToMap(tt.args.arr, tt.args.fn, tt.args.isCoverPrev); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReverseSelf(t *testing.T) {
	type args struct {
		arr []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "t1",
			args: args{
				arr: number.Range(1, 10, 1),
			},
			want: number.Range(10, 1, -1),
		}, {
			name: "t2",
			args: args{
				arr: number.Range(1, 9, 1),
			},
			want: number.Range(9, 1, -1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ReverseSelf(tt.args.arr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReverseSelf() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilterAndMap(t *testing.T) {
	type a struct {
		x int
		y string
	}

	type args[T any, N any] struct {
		arr []T
		fn  func(T) (N, bool)
	}
	type testCase[T any, N any] struct {
		name  string
		args  args[T, N]
		wantR []N
	}
	tests := []testCase[a, string]{
		{
			name: "t1",
			args: args[a, string]{
				arr: []a{
					{1, "1"}, {2, "2"}, {3, "3"},
				},
				fn: func(t a) (r string, ok bool) {
					if t.x > 2 {
						r = t.y
						ok = true
					}
					return
				},
			},
			wantR: []string{"3"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotR := FilterAndMap[string](tt.args.arr, tt.args.fn); !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("FilterAndMap() = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestGroupBy(t *testing.T) {
	type args[T x, K int, V string] struct {
		a  []T
		fn func(T) (K, V)
	}
	type testCase[T x, K int, V string] struct {
		name string
		args args[T, K, V]
		want map[K][]V
	}
	tests := []testCase[x, int, string]{
		{
			name: "t1",
			args: args[x, int, string]{
				a: Map([]int{1, 1, 2, 2, 3, 3, 4, 4, 5, 5}, func(t int) x {
					return x{t, number.ToString(t)}
				}),
				fn: func(v x) (int, string) {
					return v.int, v.y
				},
			},
			want: map[int][]string{
				1: {"1", "1"},
				2: {"2", "2"},
				3: {"3", "3"},
				4: {"4", "4"},
				5: {"5", "5"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GroupBy(tt.args.a, tt.args.fn); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GroupBy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToAnySlice(t *testing.T) {
	type args[T int] struct {
		a []T
	}
	type testCase[T int] struct {
		name string
		args args[T]
		want []any
	}
	tests := []testCase[int]{
		{
			name: "t1",
			args: args[int]{number.Range(1, 5, 1)},
			want: []any{1, 2, 3, 4, 5},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToAnySlice(tt.args.a); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToAnySlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearchFirst(t *testing.T) {
	type args[T int] struct {
		arr []T
		fn  func(T) bool
	}
	type testCase[T int] struct {
		name  string
		args  args[T]
		want  int
		want1 T
	}
	tests := []testCase[int]{
		{
			name: "t1",
			args: args[int]{
				arr: number.Range(1, 10, 1),
				fn: func(t int) bool {
					return t == 5
				},
			},
			want:  4,
			want1: 5,
		}, {
			name: "t2",
			args: args[int]{
				arr: number.Range(1, 10, 1),
				fn: func(t int) bool {
					return t == 11
				},
			},
			want:  -1,
			want1: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := SearchFirst(tt.args.arr, tt.args.fn)
			if got != tt.want {
				t.Errorf("SearchFirst() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("SearchFirst() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestSearchLast(t *testing.T) {
	type args[T int] struct {
		arr []T
		fn  func(T) bool
	}
	type testCase[T int] struct {
		name  string
		args  args[T]
		want  int
		want1 T
	}
	tests := []testCase[int]{
		{
			name: "t1",
			args: args[int]{
				arr: []int{1, 55, 5, 5, 5, 5, 22},
				fn: func(t int) bool {
					return t == 5
				},
			},
			want:  5,
			want1: 5,
		}, {
			name: "t2",
			args: args[int]{
				arr: number.Range(1, 10, 1),
				fn: func(t int) bool {
					return t == 11
				},
			},
			want:  -1,
			want1: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := SearchLast(tt.args.arr, tt.args.fn)
			if got != tt.want {
				t.Errorf("SearchLast() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("SearchLast() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestWalk(t *testing.T) {
	type args[T int] struct {
		arr []T
		fn  func(*T)
	}
	type testCase[T int] struct {
		name string
		args args[T]
	}
	tests := []testCase[int]{
		{
			name: "t1",
			args: args[int]{
				arr: number.Range(1, 10, 1),
				fn: func(i *int) {
					*i = *i * 2
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Println(tt.args.arr)
			Walk(tt.args.arr, tt.args.fn)
			fmt.Println(tt.args.arr)
		})
	}
}

func TestFill(t *testing.T) {
	type args[T int] struct {
		start int
		len   int
		v     T
	}
	type testCase[T int] struct {
		name string
		args args[T]
		want []T
	}
	tests := []testCase[int]{
		{
			name: "t1",
			args: args[int]{
				start: 2,
				len:   3,
				v:     1,
			},
			want: []int{0, 0, 1, 1, 1},
		}, {
			name: "t2",
			args: args[int]{
				start: 0,
				len:   3,
				v:     2,
			},
			want: []int{2, 2, 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Fill(tt.args.start, tt.args.len, tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Fill() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPad(t *testing.T) {
	type args[T int] struct {
		a      []T
		length int
		v      T
	}
	type testCase[T int] struct {
		name string
		args args[T]
		want []T
	}
	tests := []testCase[int]{
		{
			name: "length >0",
			args: args[int]{
				a:      []int{1, 2},
				length: 5,
				v:      10,
			},
			want: []int{1, 2, 10, 10, 10},
		},
		{
			name: "length <0",
			args: args[int]{
				a:      []int{1, 2},
				length: -5,
				v:      10,
			},
			want: []int{10, 10, 10, 1, 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Pad(tt.args.a, tt.args.length, tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Pad() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPop(t *testing.T) {
	type args[T int] struct {
		a *[]T
	}
	type testCase[T int] struct {
		name string
		args args[T]
		want T
	}
	a := number.Range(1, 10, 1)
	tests := []testCase[int]{
		{
			name: "t1",
			args: args[int]{
				a: &a,
			},
			want: 10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Pop(tt.args.a); !reflect.DeepEqual(got, tt.want) && !reflect.DeepEqual(a, number.Range(1, 9, 1)) {
				t.Errorf("Pop() = %v, want %v", got, tt.want)
			}
			fmt.Println(a)
		})
	}
}

func TestRand(t *testing.T) {
	type args[T int] struct {
		a []T
	}
	type testCase[T int] struct {
		name string
		args args[T]
	}
	tests := []testCase[int]{
		{
			name: "t1",
			args: args[int]{
				number.Range(1, 5, 1),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i := 0; i < 50; i++ {
				got, got1 := Rand(tt.args.a)
				fmt.Println(got, got1)
			}
		})
	}
}

func TestRandPop(t *testing.T) {
	type args[T int] struct {
		a *[]T
	}
	type testCase[T int] struct {
		name string
		args args[T]
		want T
	}
	a := number.Range(1, 10, 1)
	tests := []testCase[int]{
		{
			name: "t1",
			args: args[int]{
				a: &a,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i := 0; i < 11; i++ {
				got, l := RandPop(tt.args.a)
				fmt.Println(got, l, a)
			}
		})
	}
}

func TestShift(t *testing.T) {
	type args[T int] struct {
		a *[]T
	}
	type testCase[T int] struct {
		name  string
		args  args[T]
		want  T
		want1 int
	}
	a := number.Range(1, 10, 1)
	tests := []testCase[int]{
		{
			name: "t1",
			args: args[int]{&a},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i := 0; i < 11; i++ {
				got, got1 := Shift(tt.args.a)
				fmt.Println(got, got1)
			}
		})
	}
}
