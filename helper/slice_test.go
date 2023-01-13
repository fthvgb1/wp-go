package helper

import (
	"reflect"
	"testing"
)

func TestSlicePagination(t *testing.T) {
	arr := RangeSlice[int](1, 10, 1)
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
			want: RangeSlice[int](1, 2, 1),
		}, {
			name: "t2",
			args: args{
				arr:      arr,
				page:     2,
				pageSize: 2,
			},
			want: RangeSlice[int](3, 4, 1),
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
			if got := SlicePagination(tt.args.arr, tt.args.page, tt.args.pageSize); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SlicePagination() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSliceReduce(t *testing.T) {
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
			args: args{arr: RangeSlice(1, 10, 1), fn: func(i int, i2 int) int {
				return i + i2
			}},
			want: 55,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SliceReduce(tt.args.arr, tt.args.fn, tt.args.r); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SliceReduce() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSliceFilter(t *testing.T) {
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
			args: args{arr: RangeSlice(1, 10, 1), fn: func(i int) bool {
				return i > 4
			}},
			want: RangeSlice(5, 10, 1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SliceFilter(tt.args.arr, tt.args.fn); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SliceFilter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSliceMap(t *testing.T) {
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
				arr: RangeSlice[int8](1, 10, 1),
				fn: func(i int8) int {
					return int(i)
				},
			},
			want: RangeSlice(1, 10, 1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SliceMap(tt.args.arr, tt.args.fn); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SliceMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSliceReverse(t *testing.T) {
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
			args: args{arr: RangeSlice(1, 10, 1)},
			want: RangeSlice(10, 1, -1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SliceReverse(tt.args.arr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SliceReverse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRangeSlice(t *testing.T) {
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
			if got := RangeSlice(tt.args.start, tt.args.end, tt.args.step); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RangeSlice() = %v, want %v", got, tt.want)
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
				arr:    RangeSlice(1, 10, 1),
				offset: 3,
				length: 2,
			},
			wantR: RangeSlice(4, 5, 1),
		},
		{
			name: "t2",
			args: args{
				arr:    RangeSlice(1, 10, 1),
				offset: 3,
				length: 0,
			},
			wantR: RangeSlice(4, 10, 1),
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
				arr: RangeSlice(1, 5, 1),
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
				arr: RangeSlice(1, 5, 1),
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
				arr: RangeSlice(1, 5, 1),
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

func TestSliceChunk(t *testing.T) {
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
				arr:  RangeSlice(1, 7, 1),
				size: 2,
			},
			want: [][]int{{1, 2}, {3, 4}, {5, 6}, {7}},
		},
		{
			name: "t2",
			args: args{
				arr:  RangeSlice(1, 8, 1),
				size: 2,
			},
			want: [][]int{{1, 2}, {3, 4}, {5, 6}, {7, 8}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SliceChunk(tt.args.arr, tt.args.size); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SliceChunk() = %v, want %v", got, tt.want)
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
			if got := SimpleSliceToMap(tt.args.arr, tt.args.fn); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SimpleSliceToMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSliceToMap(t *testing.T) {
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
			if got := SliceToMap(tt.args.arr, tt.args.fn, tt.args.isCoverPrev); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SliceToMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSliceSelfReverse(t *testing.T) {
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
				arr: RangeSlice(1, 10, 1),
			},
			want: RangeSlice(10, 1, -1),
		}, {
			name: "t2",
			args: args{
				arr: RangeSlice(1, 9, 1),
			},
			want: RangeSlice(9, 1, -1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SliceSelfReverse(tt.args.arr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SliceSelfReverse() = %v, want %v", got, tt.want)
			}
		})
	}
}
