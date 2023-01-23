package slice

import (
	"github.com/fthvgb1/wp-go/helper/number"
	"reflect"
	"testing"
)

type x struct {
	int
	y string
}

func y(i []int) []x {
	return Map(i, func(t int) x {
		return x{t, ""}
	})
}

func TestDiff(t *testing.T) {
	type args[T int] struct {
		a []T
		b [][]T
	}
	type testCase[T int] struct {
		name  string
		args  args[T]
		wantR []T
	}
	tests := []testCase[int]{
		{
			name: "t1",
			args: args[int]{
				a: number.Range(1, 10, 1),
				b: [][]int{number.Range(3, 7, 1), number.Range(6, 9, 1)},
			},
			wantR: []int{1, 2, 10},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotR := Diff(tt.args.a, tt.args.b...); !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("Diff() = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestDiffByFn(t *testing.T) {

	type args[T x] struct {
		a  []T
		fn func(i, j T) bool
		b  [][]T
	}
	type testCase[T x] struct {
		name  string
		args  args[T]
		wantR []T
	}
	tests := []testCase[x]{
		{
			name: "t1",
			args: args[x]{
				a: y(number.Range(1, 10, 1)),
				fn: func(i, j x) bool {
					return i.int == j.int
				},
				b: [][]x{
					y(number.Range(3, 7, 1)),
					y(number.Range(6, 9, 1)),
				},
			},
			wantR: []x{{1, ""}, {2, ""}, {10, ""}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotR := DiffByFn(tt.args.a, tt.args.fn, tt.args.b...); !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("DiffByFn() = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestIntersect(t *testing.T) {
	type args[T int] struct {
		a []T
		b [][]T
	}
	type testCase[T int] struct {
		name  string
		args  args[T]
		wantR []T
	}
	tests := []testCase[int]{
		{
			name: "t1",
			args: args[int]{
				a: number.Range(1, 10, 1),
				b: [][]int{number.Range(3, 7, 1), number.Range(6, 9, 1)},
			},
			wantR: []int{6, 7},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotR := Intersect(tt.args.a, tt.args.b...); !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("Intersect() = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestIntersectByFn(t *testing.T) {
	type args[T x] struct {
		a  []T
		fn func(i, j T) bool
		b  [][]T
	}
	type testCase[T x] struct {
		name  string
		args  args[T]
		wantR []T
	}
	tests := []testCase[x]{
		{
			name: "t1",
			args: args[x]{
				a: y(number.Range(1, 10, 1)),
				fn: func(i, j x) bool {
					return i.int == j.int
				},
				b: [][]x{
					y(number.Range(3, 7, 1)),
					y(number.Range(6, 9, 1)),
				},
			},
			wantR: []x{{6, ""}, {7, ""}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotR := IntersectByFn(tt.args.a, tt.args.fn, tt.args.b...); !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("IntersectByFn() = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestUnique(t *testing.T) {
	type args[T int] struct {
		a [][]T
	}
	type testCase[T int] struct {
		name  string
		args  args[T]
		wantR []T
	}
	tests := []testCase[int]{
		{
			name: "t1",
			args: args[int]{
				a: [][]int{
					number.Range(1, 5, 1),
					number.Range(3, 6, 1),
					number.Range(6, 15, 1),
				},
			},
			wantR: number.Range(1, 15, 1),
		},
		{
			name: "t2",
			args: args[int]{
				a: [][]int{
					{1, 1, 1, 2, 2, 2, 3, 3, 4, 5},
				},
			},
			wantR: number.Range(1, 5, 1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotR := Unique(tt.args.a...); !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("Unique() = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestUniqueByFn(t *testing.T) {
	type args[T x] struct {
		fn func(T, T) bool
		a  [][]T
	}
	type testCase[T x] struct {
		name  string
		args  args[T]
		wantR []T
	}
	tests := []testCase[x]{
		{
			name: "t1",
			args: args[x]{
				fn: func(i, j x) bool {
					return i.int == j.int
				},
				a: [][]x{y([]int{1, 1, 2, 2, 3, 3}), y([]int{2, 2, 4, 4})},
			},
			wantR: y([]int{1, 2, 3, 4}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotR := UniqueByFn(tt.args.fn, tt.args.a...); !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("UniqueByFn() = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestUniqueNewByFn(t *testing.T) {
	type args[T x, V int] struct {
		fn    func(T, T) bool
		fnVal func(T) V
		a     [][]T
	}
	type testCase[T x, V int] struct {
		name  string
		args  args[T, V]
		wantR []V
	}
	tests := []testCase[x, int]{
		{
			name: "t1",
			args: args[x, int]{
				fn: func(i, j x) bool {
					return i.int == j.int
				},
				fnVal: func(i x) int {
					return i.int
				},
				a: [][]x{y([]int{1, 1, 2, 2, 3, 3}), y([]int{2, 2, 4, 4})},
			},
			wantR: []int{1, 2, 3, 4},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotR := UniqueNewByFn(tt.args.fn, tt.args.fnVal, tt.args.a...); !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("UniqueNewByFn() = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}
