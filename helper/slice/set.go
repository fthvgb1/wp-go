package slice

func IsContained[T comparable](a T, arr []T) bool {
	for _, v := range arr {
		if a == v {
			return true
		}
	}
	return false
}

func IsContainedByFn[T any](a []T, e T, fn func(i, j T) bool) bool {
	for _, t := range a {
		if fn(e, t) {
			return true
		}
	}
	return false
}

// Diff return elements which in a and not in b,...
func Diff[T comparable](a []T, b ...[]T) (r []T) {
	for _, t := range a {
		f := false
		for _, ts := range b {
			if IsContained(t, ts) {
				f = true
				break
			}
		}
		if f {
			continue
		}
		r = append(r, t)
	}
	return
}

func DiffByFn[T any](a []T, fn func(i, j T) bool, b ...[]T) (r []T) {
	for _, t := range a {
		f := false
		for _, ts := range b {
			if IsContainedByFn(ts, t, fn) {
				f = true
				break
			}
		}
		if f {
			continue
		}
		r = append(r, t)
	}
	return
}

func Intersect[T comparable](a []T, b ...[]T) (r []T) {
	for _, t := range a {
		f := false
		for _, ts := range b {
			if !IsContained(t, ts) {
				f = true
				break
			}
		}
		if f {
			continue
		}
		r = append(r, t)
	}
	return
}

func IntersectByFn[T any](a []T, fn func(i, j T) bool, b ...[]T) (r []T) {
	for _, t := range a {
		f := false
		for _, ts := range b {
			if !IsContainedByFn(ts, t, fn) {
				f = true
				break
			}
		}
		if f {
			continue
		}
		r = append(r, t)
	}
	return
}

func Unique[T comparable](a ...[]T) (r []T) {
	m := map[T]struct{}{}
	for _, ts := range a {
		for _, t := range ts {
			if _, ok := m[t]; !ok {
				m[t] = struct{}{}
				r = append(r, t)
			} else {
				continue
			}
		}
	}
	return
}

func UniqueByFn[T any](fn func(T, T) bool, a ...[]T) (r []T) {
	for _, ts := range a {
		for _, t := range ts {
			if !IsContainedByFn(r, t, fn) {
				r = append(r, t)
			}
		}
	}
	return r
}
func UniqueNewByFn[T, V any](fn func(T, T) bool, fnVal func(T) V, a ...[]T) (r []V) {
	var rr []T
	for _, ts := range a {
		for _, t := range ts {
			if !IsContainedByFn(rr, t, fn) {
				rr = append(rr, t)
				r = append(r, fnVal(t))
			}
		}
	}
	return r
}
