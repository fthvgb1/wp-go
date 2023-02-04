package stream

import (
	"fmt"
	"github.com/fthvgb1/wp-go/helper/number"
	"github.com/fthvgb1/wp-go/helper/slice"
	"reflect"
	"strconv"
	"testing"
)

var s = NewStream(number.Range(1, 10, 1))

func TestSimpleSliceStream_ForEach(t *testing.T) {
	type args[T int] struct {
		fn func(T)
	}
	type testCase[T int] struct {
		name string
		r    Stream[T]
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
		r    Stream[T]
		args args
		want Stream[T]
	}
	tests := []testCase[int]{
		{
			name: "t1",
			r:    s,
			args: args{
				limit:  3,
				offset: 5,
			},
			want: Stream[int]{number.Range(6, 8, 1)},
		},
		{
			name: "t2",
			r:    s,
			args: args{
				limit:  3,
				offset: 9,
			},
			want: Stream[int]{number.Range(10, 10, 1)},
		},
		{
			name: "t3",
			r:    s,
			args: args{
				limit:  3,
				offset: 11,
			},
			want: Stream[int]{},
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

func TestSimpleSliceStream_Result(t *testing.T) {
	type testCase[T int] struct {
		name string
		r    Stream[T]
		want []T
	}
	tests := []testCase[int]{
		{
			name: "t1",
			r:    s,
			want: number.Range(1, 10, 1),
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
		r    Stream[T]
		args args[T]
		want Stream[T]
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
			want: Stream[int]{number.Range(10, 1, -1)},
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
		r    Stream[T]
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

func TestReduce(t *testing.T) {
	type args[S, T int] struct {
		s    Stream[S]
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
			wantR: number.Sum(s.Result()...),
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
		r    Stream[T]
		want Stream[T]
	}
	tests := []testCase[int]{
		{
			name: "t1",
			r:    NewStream(number.Range(1, 10, 1)),
			want: Stream[int]{number.Range(10, 1, -1)},
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

var x = number.Range(1, 100000, 1)

func TestSimpleStreamMap(t *testing.T) {
	type args[T int, R string] struct {
		a  Stream[T]
		fn func(T) R
	}
	type testCase[T int, R string] struct {
		name string
		args args[T, R]
		want Stream[R]
	}
	tests := []testCase[int, string]{
		{
			name: "t1",
			args: args[int, string]{
				a:  NewStream(x),
				fn: strconv.Itoa,
			},
			want: Stream[string]{
				slice.Map(x, strconv.Itoa),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MapNewStream(tt.args.a, tt.args.fn); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MapNewStream() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSimpleParallelFilterAndMap(t *testing.T) {
	type args[T string, R int] struct {
		a  Stream[string]
		fn func(T) (R, bool)
		c  int
	}
	type testCase[T string, R int] struct {
		name string
		args args[T, R]
		want Stream[R]
	}
	tests := []testCase[string, int]{
		{
			name: "t1",
			args: args[string, int]{
				a: NewStream(slice.Map(x, strconv.Itoa)),
				fn: func(s string) (int, bool) {
					i, _ := strconv.Atoi(s)
					if i > 50000 {
						return i, true
					}
					return 0, false
				},
				c: 6,
			},
			want: NewStream(x[50000:]),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParallelFilterAndMap(tt.args.a, tt.args.fn, tt.args.c).Sort(func(i, j int) bool {
				return i < j
			}); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParallelFilterAndMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSimpleStreamFilterAndMap(t *testing.T) {
	type args[T string, R int] struct {
		a  Stream[T]
		fn func(T) (R, bool)
	}
	type testCase[T any, R any] struct {
		name string
		args args[string, int]
		want Stream[R]
	}
	tests := []testCase[string, int]{
		{
			name: "t1",
			args: args[string, int]{
				a: NewStream(slice.Map(x, strconv.Itoa)),
				fn: func(s string) (int, bool) {
					i, _ := strconv.Atoi(s)
					if i > 50000 {
						return i, true
					}
					return 0, false
				},
			},
			want: NewStream(x[50000:]),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FilterAndMapNewStream(tt.args.a, tt.args.fn); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FilterAndMapNewStream() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSimpleSliceStream_Len(t *testing.T) {
	type testCase[T int] struct {
		name string
		r    Stream[T]
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

func TestSimpleParallelFilterAndMapToMap(t *testing.T) {
	type args[T int, K int, V int] struct {
		a  Stream[V]
		fn func(t T) (K, V, bool)
		c  int
	}
	type testCase[T int, K int, V int] struct {
		name  string
		args  args[T, K, V]
		wantR MapStream[K, V]
	}
	tests := []testCase[int, int, int]{
		{
			name: "t1",
			args: args[int, int, int]{
				a: NewStream(x),
				fn: func(v int) (int, int, bool) {
					if v >= 50000 {
						return v, v, true
					}
					return 0, 0, false
				},
				c: 6,
			},
			wantR: NewSimpleMapStream(slice.ToMap(x[50000:], func(t int) (int, int) {
				return t, t
			}, true)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotR := ParallelFilterAndMapToMapStream(tt.args.a, tt.args.fn, tt.args.c); !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("ParallelFilterAndMapToMapStream() = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestSimpleSliceFilterAndMapToMap(t *testing.T) {
	type args[T int, K int, V int] struct {
		a           Stream[T]
		fn          func(t T) (K, V, bool)
		isCoverPrev bool
	}
	type testCase[T int, K int, V int] struct {
		name  string
		args  args[T, K, V]
		wantR MapStream[K, V]
	}
	tests := []testCase[int, int, int]{
		{
			name: "t1",
			args: args[int, int, int]{
				a: NewStream(number.Range(1, 10, 1)),
				fn: func(i int) (int, int, bool) {
					if i > 6 {
						return i, i, true
					}
					return 0, 0, false
				},
			},
			wantR: NewSimpleMapStream(slice.ToMap(number.Range(7, 10, 1), func(t int) (int, int) {
				return t, t
			}, true)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotR := SliceFilterAndMapToMapStream(tt.args.a, tt.args.fn, tt.args.isCoverPrev); !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("SliceFilterAndMapToMapStream() = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestStream_ParallelFilterAndMap(t *testing.T) {
	type xy struct {
		args int
		res  string
	}
	type args[T any] struct {
		fn func(T) (T, bool)
		c  int
	}
	type testCase[T xy] struct {
		name string
		r    Stream[T]
		args args[T]
		want Stream[T]
	}
	tests := []testCase[xy]{
		{
			name: "t1",
			r: NewStream(slice.Map(number.Range(1, 10, 1), func(t int) xy {
				return xy{args: t}
			})),
			args: args[xy]{func(v xy) (xy, bool) {
				v.res = strconv.Itoa(v.args)
				return v, true
			}, 6},
			want: NewStream(slice.Map(number.Range(1, 10, 1), func(t int) xy {
				return xy{args: t, res: strconv.Itoa(t)}
			})),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.ParallelFilterAndMap(tt.args.fn, tt.args.c).Sort(func(i, j xy) bool {
				return i.args < j.args
			}); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParallelFilterAndMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStream_Reduce(t *testing.T) {
	type aa struct {
		args int
		res  int
	}
	type args[T any] struct {
		fn   func(v, r T) T
		init T
	}
	type testCase[T any] struct {
		name string
		r    Stream[T]
		args args[T]
		want T
	}
	tests := []testCase[aa]{
		{
			name: "t1",
			r: NewStream(slice.Map(number.Range(1, 10, 1), func(t int) aa {
				return aa{args: t}
			})),
			args: args[aa]{func(v, r aa) aa {
				return aa{res: v.args + r.res}
			}, aa{}},
			want: aa{res: 55},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.Reduce(tt.args.fn, tt.args.init); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Reduce() = %v, want %v", got, tt.want)
			}
		})
	}
}
