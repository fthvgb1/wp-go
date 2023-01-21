package slice

import (
	"reflect"
	"testing"
)

func TestSortSelf(t *testing.T) {
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
			if SortSelf(tt.args.arr, tt.args.fn); !reflect.DeepEqual(tt.args.arr, tt.wantR) {
				t.Errorf("SimpleSort() = %v, want %v", tt.args.arr, tt.wantR)
			}
		})
	}
}

func TestSort(t *testing.T) {
	type xy struct {
		x int
		y int
	}
	type args[T any] struct {
		arr []T
		fn  func(i, j T) bool
	}
	type testCase[T any] struct {
		name  string
		args  args[T]
		wantR []T
	}
	tests := []testCase[xy]{
		{
			name: "t1",
			args: args[xy]{
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
			if gotR := Sort[xy](tt.args.arr, tt.args.fn); !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("SimpleSortR() = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}