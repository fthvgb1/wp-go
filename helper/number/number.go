// Package number
// 使用随机数时需要先 调用 rand.seed()函数
package number

import (
	"fmt"
	"golang.org/x/exp/constraints"
	"math"
	"math/rand"
)

func Range[T constraints.Integer](start, end, step T) []T {
	if step == 0 {
		panic("step can't be 0")
	}
	l := int((end-start+1)/step + 1)
	if l < 0 {
		l = 0 - l
	}
	r := make([]T, 0, l)
	for i := start; ; {
		r = append(r, i)
		i = i + step
		if (step > 0 && i > end) || (step < 0 && i < end) {
			break
		}
	}
	return r
}

// Rand 都为闭区间 [start,end]
func Rand[T constraints.Integer](start, end T) T {
	end++
	return T(rand.Int63n(int64(end-start))) + start
}

func Min[T constraints.Integer | constraints.Float](a ...T) T {
	min := a[0]
	for _, t := range a {
		if min > t {
			min = t
		}
	}
	return min
}

func Max[T constraints.Integer | constraints.Float](a ...T) T {
	max := a[0]
	for _, t := range a {
		if max < t {
			max = t
		}
	}
	return max
}

func Sum[T constraints.Integer | constraints.Float](a ...T) T {
	s := T(0)
	for _, t := range a {
		s += t
	}
	return s
}

func Add[T constraints.Integer | constraints.Float](i, j T) T {
	return i + j
}
func Sub[T constraints.Integer | constraints.Float](i, j T) T {
	return i - j
}

func ToString[T constraints.Integer | constraints.Float](n T) string {
	return fmt.Sprintf("%v", n)
}

func Abs[T constraints.Integer | constraints.Float](n T) T {
	if n >= 0 {
		return n
	}
	return -n
}

func Mul[T constraints.Integer | constraints.Float](i, j T) T {
	return i * j
}

func Divide[T constraints.Integer | constraints.Float](i, j T) T {
	return i / j
}

func CalTotalPage[T constraints.Integer](totalRows, size T) T {
	return T(math.Ceil(float64(totalRows) / float64(size)))
}
