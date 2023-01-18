package stream

import (
	"fmt"
	"github.com/fthvgb1/wp-go/helper"
	"reflect"
	"strconv"
	"testing"
)

func TestNewSimpleMapStream(t *testing.T) {
	type args[K int, V int] struct {
		m map[K]V
	}
	type testCase[K int, V int] struct {
		name string
		args args[K, V]
		want SimpleMapStream[K, V]
	}
	tests := []testCase[int, int]{
		{
			name: "t1",
			args: args[int, int]{make(map[int]int)},
			want: SimpleMapStream[int, int]{make(map[int]int)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSimpleMapStream(tt.args.m); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSimpleMapStream() = %v, want %v", got, tt.want)
			}
		})
	}
}

var y = helper.RangeSlice(1, 1000, 1)
var w = helper.SliceToMap(y, func(v int) (int, int) {
	return v, v
}, true)

func TestSimpleMapFilterAndMapToSlice(t *testing.T) {
	type args[K int, V int, R int] struct {
		mm SimpleMapStream[K, V]
		fn func(K, V) (R, bool)
		c  int
	}
	type testCase[K int, V int, R int] struct {
		name string
		args args[K, V, R]
		want SimpleSliceStream[R]
	}
	tests := []testCase[int, int, int]{
		{
			name: "t1",
			args: args[int, int, int]{
				mm: NewSimpleMapStream(w),
				fn: func(k, v int) (int, bool) {
					if v > 500 {
						return v, true
					}
					return 0, false
				},
				c: 6,
			},
			want: NewSimpleSliceStream(y[500:]),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SimpleMapFilterAndMapToSlice(tt.args.mm, tt.args.fn, tt.args.c).Sort(func(i, j int) bool {
				return i < j
			}); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SimpleMapFilterAndMapToSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSimpleMapParallelFilterAndMapToMap(t *testing.T) {
	type args[KK string, VV string, K int, V int] struct {
		mm SimpleMapStream[K, V]
		fn func(K, V) (KK, VV, bool)
		c  int
	}
	type testCase[KK string, VV string, K int, V int] struct {
		name string
		args args[KK, VV, K, V]
		want SimpleMapStream[KK, VV]
	}
	tests := []testCase[string, string, int, int]{
		{
			name: "t1",
			args: args[string, string, int, int]{
				mm: NewSimpleMapStream(w),
				fn: func(k, v int) (string, string, bool) {
					if v > 500 {
						t := strconv.Itoa(v)
						return t, t, true
					}
					return "", "", false
				},
				c: 6,
			},
			want: NewSimpleMapStream(helper.SliceToMap(y[500:], func(v int) (K, T string) {
				t := strconv.Itoa(v)
				return t, t
			}, true)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SimpleMapParallelFilterAndMapToMap(tt.args.mm, tt.args.fn, tt.args.c); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SimpleMapParallelFilterAndMapToMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSimpleMapStreamFilterAndMapToMap(t *testing.T) {
	type args[KK string, VV string, K int, V int] struct {
		a  SimpleMapStream[K, V]
		fn func(K, V) (KK, VV, bool)
	}
	type testCase[KK string, VV string, K int, V int] struct {
		name  string
		args  args[KK, VV, K, V]
		wantR SimpleMapStream[KK, VV]
	}
	tests := []testCase[string, string, int, int]{
		{
			name: "t1",
			args: args[string, string, int, int]{
				a: NewSimpleMapStream(w),
				fn: func(k, v int) (string, string, bool) {
					if v > 500 {
						t := strconv.Itoa(v)
						return t, t, true
					}
					return "", "", false
				},
			},
			wantR: NewSimpleMapStream(helper.SliceToMap(y[500:], func(v int) (K, T string) {
				t := strconv.Itoa(v)
				return t, t
			}, true)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotR := SimpleMapStreamFilterAndMapToMap(tt.args.a, tt.args.fn); !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("SimpleMapStreamFilterAndMapToMap() = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestSimpleMapStream_ForEach(t *testing.T) {
	type args[K int, V int] struct {
		fn func(K, V)
	}
	type testCase[K int, V int] struct {
		name string
		r    SimpleMapStream[K, V]
		args args[K, V]
	}
	tests := []testCase[int, int]{
		{
			name: "t1",
			r: NewSimpleMapStream(helper.SliceToMap(y[0:10], func(v int) (int, int) {
				return v, v
			}, true)),
			args: args[int, int]{
				fn: func(k, v int) {
					fmt.Println(k, v)
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

func TestSimpleMapStream_Len(t *testing.T) {
	type testCase[K int, V int] struct {
		name string
		r    SimpleMapStream[K, V]
		want int
	}
	tests := []testCase[int, int]{
		{
			name: "t1",
			r:    NewSimpleMapStream(w),
			want: len(w),
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

func TestSimpleMapStream_ParallelForEach(t *testing.T) {
	type args[K int, V int] struct {
		fn func(K, V)
		c  int
	}
	type testCase[K int, V int] struct {
		name string
		r    SimpleMapStream[K, V]
		args args[K, V]
	}
	tests := []testCase[int, int]{
		{
			name: "t1",
			r:    NewSimpleMapStream(w),
			args: args[int, int]{
				func(k, v int) {
					fmt.Println(k, v)
				},
				6,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.r.ParallelForEach(tt.args.fn, tt.args.c)
		})
	}
}

func TestSimpleMapStream_Result(t *testing.T) {
	type testCase[K int, V int] struct {
		name string
		r    SimpleMapStream[K, V]
		want map[K]V
	}
	tests := []testCase[int, int]{
		{
			name: "t1",
			r: NewSimpleMapStream(helper.SliceToMap(y, func(v int) (int, int) {
				return v, v
			}, true)),
			want: helper.SliceToMap(y, func(v int) (int, int) {
				return v, v
			}, true),
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
