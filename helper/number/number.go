// Package number
// 使用随机数时需要先 调用 rand.seed()函数
package number

import (
	"fmt"
	"golang.org/x/exp/constraints"
	"math"
	"math/rand"
	"strconv"
)

func Range[T constraints.Integer](start, end T, steps ...T) []T {
	step := T(1)
	if len(steps) > 0 {
		step = steps[0]
	}
	var l int
	if step == 0 {
		l = int(end - start + 1)
	} else {
		l = int((end - start + 1) / step)
		if step*T(l) <= end && step != 1 {
			l++
		}
	}
	if l < 0 {
		l = -l
	}
	r := make([]T, 0, l)
	gap := start
	for i := 0; i < l; i++ {
		r = append(r, gap)
		gap += step
	}
	return r
}

// Rand 都为闭区间 [start,end]
func Rand[T constraints.Integer](start, end T) T {
	end++
	return T(rand.Int63n(int64(end-start))) + start
}

func Min[T constraints.Integer | constraints.Float](a ...T) T {
	mins := a[0]
	for _, t := range a {
		if mins > t {
			mins = t
		}
	}
	return mins
}

func Max[T constraints.Integer | constraints.Float](a ...T) T {
	maxs := a[0]
	for _, t := range a {
		if maxs < t {
			maxs = t
		}
	}
	return maxs
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

func IntToString[T constraints.Integer](i T) string {
	return strconv.FormatInt(int64(i), 10)
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

func Ceil[T constraints.Integer](num1, num2 T) int {
	return int((num1 + num2 - 1) / num2)
}

func DivideCeil[T constraints.Integer](num1, num2 T) T {
	return T(math.Ceil(float64(num1) / float64(num2)))
}

type Counter[T constraints.Integer] func() T

func Counters[T constraints.Integer]() func() T {
	var count T
	return func() T {
		count++
		return count
	}
}
