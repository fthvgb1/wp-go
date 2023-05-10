package slice

import (
	"fmt"
	"golang.org/x/exp/constraints"
	"reflect"
	"testing"
)

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
		{
			name: "t1",
			args: args[xy]{
				arr: []xy{{1, 2}, {1, 3}, {1, 6}},
				fn: func(i, j xy) bool {
					return i.x > j.x
				},
			},
			wantR: []xy{{1, 2}, {1, 3}, {1, 6}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if Sort[xy](tt.args.arr, tt.args.fn); !reflect.DeepEqual(tt.args.arr, tt.wantR) {
				t.Errorf("SimpleSortR() = %v, want %v", tt.args.arr, tt.wantR)
			}
		})
	}
}

func TestSorts(t *testing.T) {
	type args[T constraints.Ordered] struct {
		a     []T
		order int
	}
	type testCase[T constraints.Ordered] struct {
		name string
		args args[T]
	}
	tests := []testCase[int]{
		{
			name: "asc",
			args: args[int]{
				a:     []int{1, -3, 6, 10, 3, 2, 8},
				order: ASC,
			}, //[-3 1 2 3 6 8 10]
		},
		{
			name: "desc",
			args: args[int]{
				a:     []int{1, -3, 6, 10, 3, 2, 8},
				order: DESC,
			}, //[10 8 6 3 2 1 -3]
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Sorts(tt.args.a, tt.args.order)
			fmt.Println(tt.args.a)
		})
	}
}

func TestStableSort(t *testing.T) {
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
		{
			name: "t1",
			args: args[xy]{
				arr: []xy{{1, 2}, {1, 3}, {1, 6}},
				fn: func(i, j xy) bool {
					return i.x > j.x
				},
			},
			wantR: []xy{{1, 2}, {1, 3}, {1, 6}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			StableSort(tt.args.arr, tt.args.fn)
		})
	}
}
