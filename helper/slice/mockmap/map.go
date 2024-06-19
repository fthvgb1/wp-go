package mockmap

import (
	"github.com/fthvgb1/wp-go/helper/slice"
)

type Item[K comparable, T any] struct {
	Name  K
	Value T
	Order float64
}

type Map[K comparable, T any] []Item[K, T]

func (q *Map[K, T]) Get(name K) Item[K, T] {
	_, v := slice.SearchFirst(*q, func(t Item[K, T]) bool {
		return name == t.Name
	})
	return v
}

func (q *Map[K, T]) Set(name K, value T, orders ...float64) {
	i := slice.IndexOfBy(*q, func(t Item[K, T]) bool {
		return name == t.Name
	})
	if i > -1 {
		(*q)[i].Value = value
		return
	}
	order := float64(0)
	if len(orders) > 0 {
		order = orders[0]
	}
	*q = append(*q, Item[K, T]{name, value, order})
}

func (q *Map[K, T]) Del(name K) {
	i := slice.IndexOfBy(*q, func(t Item[K, T]) bool {
		return name == t.Name
	})
	if i > -1 {
		slice.Delete((*[]Item[K, T])(q), i)
	}
}
func (q *Map[K, T]) DelByIndex(i int) {
	slice.Delete((*[]Item[K, T])(q), i)
}
