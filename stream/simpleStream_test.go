package stream

import (
	"fmt"
	"github/fthvgb1/wp-go/helper"
	"reflect"
	"strconv"
	"testing"
)

var s = NewSimpleSliceStream(helper.RangeSlice(1, 10, 1))

func TestSimpleSliceStream_Filter(t *testing.T) {
	type args[T int] struct {
		fn func(T) bool
	}
	type testCase[T int] struct {
		name string
		r    SimpleSliceStream[T]
		args args[T]
		want SimpleSliceStream[T]
	}
	tests := []testCase[int]{
		{
			name: "t1",
			r:    s,
			args: args[int]{
				func(t int) (r bool) {
					if t > 5 {
						r = true
					}
					return
				},
			},
			want: SimpleSliceStream[int]{helper.RangeSlice(6, 10, 1)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.Filter(tt.args.fn); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Filter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSimpleSliceStream_ForEach(t *testing.T) {
	type args[T int] struct {
		fn func(T)
	}
	type testCase[T int] struct {
		name string
		r    SimpleSliceStream[T]
		args args[T]
	}
	tests := []testCase[int]{
		{
			name: "t1",
			r:    s,
			args: args[int]{
				func(t int) {
					fmt.Println(t)
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.r.ForEach(tt.args.fn)
		})
	}
}

func TestSimpleSliceStream_Limit(t *testing.T) {
	type args struct {
		limit  int
		offset int
	}
	type testCase[T int] struct {
		name string
		r    SimpleSliceStream[T]
		args args
		want SimpleSliceStream[T]
	}
	tests := []testCase[int]{
		{
			name: "t1",
			r:    s,
			args: args{
				limit:  3,
				offset: 5,
			},
			want: SimpleSliceStream[int]{helper.RangeSlice(6, 8, 1)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.Limit(tt.args.limit, tt.args.offset); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Limit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSimpleSliceStream_Map(t *testing.T) {
	type args[T int] struct {
		fn func(T) T
	}
	type testCase[T int] struct {
		name string
		r    SimpleSliceStream[T]
		args args[T]
		want SimpleSliceStream[T]
	}
	tests := []testCase[int]{
		{
			name: "t1",
			r:    s,
			args: args[int]{
				func(t int) (r int) {
					return t * 2
				},
			},
			want: SimpleSliceStream[int]{helper.RangeSlice(2, 20, 2)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.Map(tt.args.fn); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Map() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSimpleSliceStream_Result(t *testing.T) {
	type testCase[T int] struct {
		name string
		r    SimpleSliceStream[T]
		want []T
	}
	tests := []testCase[int]{
		{
			name: "t1",
			r:    s,
			want: helper.RangeSlice(1, 10, 1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.Result(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Result() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSimpleSliceStream_Sort(t *testing.T) {
	type args[T int] struct {
		fn func(i, j T) bool
	}
	type testCase[T int] struct {
		name string
		r    SimpleSliceStream[T]
		args args[T]
		want SimpleSliceStream[T]
	}
	tests := []testCase[int]{
		{
			name: "t1",
			r:    s,
			args: args[int]{
				fn: func(i, j int) bool {
					return i > j
				},
			},
			want: SimpleSliceStream[int]{helper.RangeSlice(10, 1, -1)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.Sort(tt.args.fn); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Sort() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSimpleSliceStream_parallelForEach(t *testing.T) {
	type args[T int] struct {
		fn func(T)
		c  int
	}
	type testCase[T int] struct {
		name string
		r    SimpleSliceStream[T]
		args args[T]
	}
	tests := []testCase[int]{
		{
			name: "t1",
			r:    s,
			args: args[int]{
				fn: func(t int) {
					fmt.Println(t)
				},
				c: 3,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.r.ParallelForEach(tt.args.fn, tt.args.c)
		})
	}
}

func TestSimpleSliceStream_ParallelFilter(t *testing.T) {
	type args[T int] struct {
		fn func(T) bool
		c  int
	}
	type testCase[T int] struct {
		name string
		r    SimpleSliceStream[T]
		args args[T]
		want SimpleSliceStream[T]
	}
	tests := []testCase[int]{
		{
			name: "t1",
			r:    s,
			args: args[int]{
				fn: func(t int) bool {
					return t > 3
				},
				c: 6,
			},
			want: SimpleSliceStream[int]{helper.RangeSlice(4, 10, 1)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.ParallelFilter(tt.args.fn, tt.args.c).Sort(func(i, j int) bool {
				return i < j
			}); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParallelFilter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSimpleSliceStream_ParallelMap(t *testing.T) {
	type args[T int] struct {
		fn func(T) T
		c  int
	}
	type testCase[T int] struct {
		name string
		r    SimpleSliceStream[T]
		args args[T]
		want SimpleSliceStream[T]
	}
	tests := []testCase[int]{
		{
			name: "t1",
			r:    s,
			args: args[int]{
				fn: func(t int) int {
					return t * 2
				},
				c: 6,
			},
			want: SimpleSliceStream[int]{helper.RangeSlice(2, 20, 2)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.ParallelMap(tt.args.fn, tt.args.c).Sort(func(i, j int) bool {
				return i < j
			}); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SimpleParallelMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReduce(t *testing.T) {
	type args[S, T int] struct {
		s    SimpleSliceStream[S]
		fn   func(S, T) T
		init T
	}
	type testCase[S, T int] struct {
		name  string
		args  args[S, T]
		wantR T
	}
	tests := []testCase[int, int]{
		{
			name: "t1",
			args: args[int, int]{
				s, func(i, r int) int {
					return i + r
				},
				0,
			},
			wantR: helper.Sum(s.Result()...),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotR := Reduce(tt.args.s, tt.args.fn, tt.args.init); !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("Reduce() = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestSimpleSliceStream_Reverse(t *testing.T) {
	type testCase[T int] struct {
		name string
		r    SimpleSliceStream[T]
		want SimpleSliceStream[T]
	}
	tests := []testCase[int]{
		{
			name: "t1",
			r:    NewSimpleSliceStream(helper.RangeSlice(1, 10, 1)),
			want: SimpleSliceStream[int]{helper.RangeSlice(10, 1, -1)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.Reverse(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Reverse() = %v, want %v", got, tt.want)
			}
		})
	}
}

var x = helper.RangeSlice(1, 100000, 1)

func TestSimpleStreamMap(t *testing.T) {
	type args[T int, R string] struct {
		a  SimpleSliceStream[T]
		fn func(T) R
	}
	type testCase[T int, R string] struct {
		name string
		args args[T, R]
		want SimpleSliceStream[R]
	}
	tests := []testCase[int, string]{
		{
			name: "t1",
			args: args[int, string]{
				a:  NewSimpleSliceStream(x),
				fn: strconv.Itoa,
			},
			want: SimpleSliceStream[string]{
				helper.SliceMap(x, strconv.Itoa),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SimpleStreamMap(tt.args.a, tt.args.fn); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SimpleStreamMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSimpleParallelMap(t *testing.T) {
	type args[T string, R int] struct {
		a  SimpleSliceStream[string]
		fn func(T) R
		c  int
	}
	type testCase[T string, R int] struct {
		name string
		args args[T, R]
		want SimpleSliceStream[R]
	}

	tests := []testCase[string, int]{
		{
			name: "t1",
			args: args[string, int]{
				a: NewSimpleSliceStream(helper.SliceMap(x, strconv.Itoa)),
				fn: func(s string) int {
					i, _ := strconv.Atoi(s)
					return i
				},
				c: 6,
			},
			want: NewSimpleSliceStream(x),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SimpleParallelMap(tt.args.a, tt.args.fn, tt.args.c).Sort(func(i, j int) bool {
				return i < j
			}); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SimpleParallelMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSimpleParallelFilterAndMap(t *testing.T) {
	type args[T string, R int] struct {
		a  SimpleSliceStream[string]
		fn func(T) (R, bool)
		c  int
	}
	type testCase[T string, R int] struct {
		name string
		args args[T, R]
		want SimpleSliceStream[R]
	}
	tests := []testCase[string, int]{
		{
			name: "t1",
			args: args[string, int]{
				a: NewSimpleSliceStream(helper.SliceMap(x, strconv.Itoa)),
				fn: func(s string) (int, bool) {
					i, _ := strconv.Atoi(s)
					if i > 50000 {
						return i, true
					}
					return 0, false
				},
				c: 6,
			},
			want: NewSimpleSliceStream(x[50000:]),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SimpleParallelFilterAndMap(tt.args.a, tt.args.fn, tt.args.c).Sort(func(i, j int) bool {
				return i < j
			}); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SimpleParallelFilterAndMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSimpleStreamFilterAndMap(t *testing.T) {
	type args[T string, R int] struct {
		a  SimpleSliceStream[T]
		fn func(T) (R, bool)
	}
	type testCase[T any, R any] struct {
		name string
		args args[string, int]
		want SimpleSliceStream[R]
	}
	tests := []testCase[string, int]{
		{
			name: "t1",
			args: args[string, int]{
				a: NewSimpleSliceStream(helper.SliceMap(x, strconv.Itoa)),
				fn: func(s string) (int, bool) {
					i, _ := strconv.Atoi(s)
					if i > 50000 {
						return i, true
					}
					return 0, false
				},
			},
			want: NewSimpleSliceStream(x[50000:]),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SimpleStreamFilterAndMap(tt.args.a, tt.args.fn); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SimpleStreamFilterAndMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSimpleSliceStream_Len(t *testing.T) {
	type testCase[T int] struct {
		name string
		r    SimpleSliceStream[T]
		want int
	}
	tests := []testCase[int]{
		{
			name: "t1",
			r:    s,
			want: 10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.Len(); got != tt.want {
				t.Errorf("Len() = %v, want %v", got, tt.want)
			}
		})
	}
}
