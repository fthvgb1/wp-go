package tree

import "github.com/fthvgb1/wp-go/helper/slice"

type Node[T any, K comparable] struct {
	Data     T
	Children *[]Node[T, K]
	Parent   K
}

func (n *Node[T, K]) GetChildren() []T {
	return slice.Map(*n.Children, func(t Node[T, K]) T {
		return t.Data
	})
}

func (n *Node[T, K]) ChildrenByOrder(fn func(T, T) bool) []T {
	a := slice.Map(*n.Children, func(t Node[T, K]) T {
		return t.Data
	})
	slice.Sort(a, fn)
	return a
}

func (n *Node[T, K]) loop(fn func(T, int), deep int) {
	for _, nn := range *n.Children {
		fn(nn.Data, deep)
		if len(*nn.Children) > 0 {
			nn.loop(fn, deep+1)
		}
	}
}

func (n *Node[T, K]) Loop(fn func(T, int)) {
	n.loop(fn, 0)
}
func (n *Node[T, K]) orderByLoop(fn func(T, int), orderBy func(T, T) bool, deep int) {
	slice.Sort(*n.Children, func(i, j Node[T, K]) bool {
		return orderBy(i.Data, j.Data)
	})
	for _, nn := range *n.Children {
		fn(nn.Data, deep)
		if len(*nn.Children) > 0 {
			nn.orderByLoop(fn, orderBy, deep+1)
		}
	}
}

func (n *Node[T, K]) OrderByLoop(fn func(T, int), orderBy func(T, T) bool) {
	n.orderByLoop(fn, orderBy, 0)
}

func Root[T any, K comparable](a []T, top K, fn func(T) (child, parent K)) *Node[T, K] {
	m := make(map[K]*Node[T, K])
	m[top] = &Node[T, K]{Children: new([]Node[T, K])}
	for _, t := range a {
		c, p := fn(t)
		node := Node[T, K]{Parent: p, Data: t, Children: new([]Node[T, K])}
		m[c] = &node
		parent, ok := m[p]
		if !ok {
			m[p] = &Node[T, K]{Children: new([]Node[T, K])}
			parent = m[p]
		}
		*parent.Children = append(*parent.Children, node)
	}
	return m[top]
}
