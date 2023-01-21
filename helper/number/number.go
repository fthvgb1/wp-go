package number

import (
	"fmt"
	"math/rand"
)

type IntNumber interface {
	~int | ~int64 | ~int32 | ~int8 | ~int16 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

func Range[T IntNumber](start, end, step T) []T {
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

func Rand[T IntNumber](start, end T) T {
	end++
	return T(rand.Int63n(int64(end-start))) + start
}

func Min[T IntNumber | ~float64 | ~float32](a ...T) T {
	min := a[0]
	for _, t := range a {
		if min > t {
			min = t
		}
	}
	return min
}

func Max[T IntNumber | ~float64 | ~float32](a ...T) T {
	max := a[0]
	for _, t := range a {
		if max < t {
			max = t
		}
	}
	return max
}

func Sum[T IntNumber | ~float64 | ~float32](a ...T) T {
	s := T(0)
	for _, t := range a {
		s += t
	}
	return s
}

func ToString[T IntNumber | ~float64 | ~float32](n T) string {
	return fmt.Sprintf("%v", n)
}