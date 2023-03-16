package tree

import (
	"fmt"
	"github.com/fthvgb1/wp-go/helper/slice"
	"reflect"
	"testing"
)

type xx struct {
	id  int
	pid int
}

var ss []xx

func init() {
	ss = []xx{
		{1, 0},
		{2, 0},
		{3, 0},
		{4, 1},
		{5, 1},
		{6, 2},
		{7, 3},
		{8, 4},
		{9, 5},
		{10, 2},
	}
}

func ffn(x xx) (child, parent int) {
	return x.id, x.pid
}
func TestAncestor(t *testing.T) {
	type args[K comparable, T any] struct {
		root  map[K]*Node[T, K]
		top   K
		child *Node[T, K]
	}
	type testCase[K comparable, T any] struct {
		name string
		args args[K, T]
		want *Node[T, K]
	}
	r := Roots(ss, 0, ffn)
	tests := []testCase[int, xx]{
		{
			name: "t1",
			args: args[int, xx]{
				root: r,
				top:  0,
				child: &Node[xx, int]{
					Data:   xx{9, 5},
					Parent: 5,
				},
			},
			want: &Node[xx, int]{
				Data: xx{
					id:  1,
					pid: 0,
				},
				Children: &[]Node[xx, int]{
					{
						Data: xx{4, 1},
						Children: &[]Node[xx, int]{
							{
								Data:   xx{8, 4},
								Parent: 4,
							},
						},
					},
					{
						Data: xx{5, 1},
						Children: &[]Node[xx, int]{
							{
								Data:   xx{9, 5},
								Parent: 5,
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Ancestor(tt.args.root, tt.args.top, tt.args.child); !reflect.DeepEqual(*got, *tt.want) {
				t.Errorf("Ancestor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNode_ChildrenByOrder(t *testing.T) {
	type args[T any] struct {
		fn func(T, T) bool
	}
	type testCase[T any, K comparable] struct {
		name string
		n    Node[T, K]
		args args[T]
		want []T
	}
	tests := []testCase[xx, int]{
		{
			name: "t1",
			n:    *root(ss, 0, ffn)[2],
			args: args[xx]{
				fn: func(x xx, x2 xx) bool {
					return x.id < x2.id
				},
			},
			want: []xx{{6, 2}, {10, 2}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.n.ChildrenByOrder(tt.args.fn); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ChildrenByOrder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNode_GetChildren(t *testing.T) {
	type testCase[T any, K comparable] struct {
		name string
		n    Node[T, K]
		want []T
	}
	tests := []testCase[xx, int]{
		{
			name: "t1",
			n:    *root(ss, 0, ffn)[2],
			want: []xx{{6, 2}, {10, 2}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.n.GetChildren(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetChildren() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNode_Loop(t *testing.T) {
	type args[T any] struct {
		fn func(T, int)
	}
	type testCase[T any, K comparable] struct {
		name string
		n    Node[T, K]
		args args[T]
	}
	tests := []testCase[xx, int]{
		{
			name: "t1",
			n:    *Root(ss, 0, ffn),
			args: args[xx]{
				fn: func(x xx, i int) {
					fmt.Println(x, i)
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.n.Loop(tt.args.fn)
		})
	}
}

func TestNode_OrderByLoop(t *testing.T) {
	type args[T any] struct {
		fn      func(T, int)
		orderBy func(T, T) bool
	}
	type testCase[T any, K comparable] struct {
		name string
		n    Node[T, K]
		args args[T]
	}
	tests := []testCase[xx, int]{
		{
			name: "",
			n:    *Root(ss, 0, ffn),
			args: args[xx]{
				fn: func(x xx, i int) {
					fmt.Println(x)
				},
				orderBy: func(x xx, x2 xx) bool {
					return x.id > x2.id
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.n.OrderByLoop(tt.args.fn, tt.args.orderBy)
		})
	}
}

func TestNode_Posterity(t *testing.T) {
	type testCase[T any, K comparable] struct {
		name  string
		n     Node[T, K]
		wantR []T
	}
	tests := []testCase[xx, int]{
		{
			name:  "t1",
			n:     *Root(ss, 0, ffn),
			wantR: ss,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotR := tt.n.Posterity()
			fn := func(i, j xx) bool {
				if i.id < j.id {
					return true
				}
				if i.pid < j.pid {
					return true
				}
				return false
			}
			slice.Sort(gotR, fn)
			slice.Sort(tt.wantR, fn)
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("Posterity() = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestNode_loop(t *testing.T) {
	type args[T any] struct {
		fn   func(T, int)
		deep int
	}
	type testCase[T any, K comparable] struct {
		name string
		n    Node[T, K]
		args args[T]
	}
	tests := []testCase[xx, int]{
		{
			name: "t1",
			n:    *Root(ss, 0, ffn),
			args: args[xx]{fn: func(x xx, i int) {
				fmt.Println(x)
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.n.loop(tt.args.fn, tt.args.deep)
		})
	}
}

func TestNode_orderByLoop(t *testing.T) {
	type args[T any] struct {
		fn      func(T, int)
		orderBy func(T, T) bool
		deep    int
	}
	type testCase[T any, K comparable] struct {
		name string
		n    Node[T, K]
		args args[T]
	}
	tests := []testCase[xx, int]{
		{
			name: "t1",
			n:    *Root(ss, 0, ffn),
			args: args[xx]{fn: func(x xx, i int) {
				fmt.Println(x)
			}, orderBy: func(x xx, x2 xx) bool {
				return x.id > x2.id
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.n.orderByLoop(tt.args.fn, tt.args.orderBy, tt.args.deep)
		})
	}
}

func TestNode_posterity(t *testing.T) {
	type args[T any] struct {
		a *[]T
	}
	type testCase[T any, K comparable] struct {
		name string
		n    Node[T, K]
		args args[T]
	}
	tests := []testCase[xx, int]{
		{
			name: "t1",
			n:    *Root(ss, 0, ffn),
			args: args[xx]{a: new([]xx)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.n.posterity(tt.args.a)
		})
	}
}

func TestRoot(t *testing.T) {
	type args[T any, K comparable] struct {
		a   []T
		top K
		fn  func(T) (child, parent K)
	}
	type testCase[T any, K comparable] struct {
		name string
		args args[T, K]
		want *Node[T, K]
	}
	tests := []testCase[xx, int]{
		{
			name: "t1",
			args: args[xx, int]{
				a:   ss,
				top: 0,
				fn:  ffn,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Root(tt.args.a, tt.args.top, tt.args.fn); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Root() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRoots(t *testing.T) {
	type args[T any, K comparable] struct {
		a   []T
		top K
		fn  func(T) (child, parent K)
	}
	type testCase[T any, K comparable] struct {
		name string
		args args[T, K]
		want map[K]*Node[T, K]
	}
	tests := []testCase[xx, int]{
		{
			name: "t1",
			args: args[xx, int]{ss, 0, ffn},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Roots(tt.args.a, tt.args.top, tt.args.fn); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Roots() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_root(t *testing.T) {
	type args[T any, K comparable] struct {
		a   []T
		top K
		fn  func(T) (child, parent K)
	}
	type testCase[T any, K comparable] struct {
		name string
		args args[T, K]
		want map[K]*Node[T, K]
	}
	tests := []testCase[xx, int]{
		{
			name: "t1",
			args: args[xx, int]{ss, 0, ffn},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := root(tt.args.a, tt.args.top, tt.args.fn); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("root() = %v, want %v", got, tt.want)
			}
		})
	}
}
