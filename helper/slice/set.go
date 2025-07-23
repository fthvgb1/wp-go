package slice

func IsContained[T comparable](arr []T, a T) bool {
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
	if len(b) < 1 {
		r = a
		return
	}
	bb := b[0]
	for _, ts := range b[1:] {
		bb = append(bb, ts...)
	}
	r = Diffs(a, bb)
	return
}

// Diffs return elements which in a and not in b
func Diffs[T comparable](a, b []T) (r []T) {
	if len(b) < 1 {
		r = Copy(a)
		return
	}
	m := map[T]struct{}{}
	i := 0
	for _, v := range a {
		if _, ok := m[v]; ok {
			continue
		}
		if i == len(b)-1 {
			r = append(r, v)
			continue
		}
		for ; i < len(b); i++ {
			m[b[i]] = struct{}{}
			if v == b[i] {
				continue
			}
		}
		r = append(r, v)
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

func DiffNewByFn[T, V any](a []T, fn func(i, j T) bool, fnV func(T) V, b ...[]T) (r []V) {
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
		r = append(r, fnV(t))
	}
	return
}

func Intersect[T comparable](a []T, b ...[]T) (r []T) {
	for _, t := range a {
		f := false
		for _, ts := range b {
			if !IsContained(ts, t) {
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

func IntersectNewByFn[T, V any](a []T, fn func(i, j T) bool, fnV func(T) V, b ...[]T) (r []V) {
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
		r = append(r, fnV(t))
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
