package slice

func Splice[T any](a *[]T, offset, length int, replacement []T) []T {
	arr := *a
	l := len(arr)
	if length < 0 {
		panic("length must > 0")
	}
	if offset >= 0 {
		if offset+length > l {
			if offset == 0 {
				*a = []T{}
				return arr[:l]
			}
			return nil
		} else if l > offset && l < offset+length {
			v := arr[offset:l]
			*a = append(arr[:offset], replacement...)
			return v
		} else if offset+length <= l {
			v := append([]T{}, arr[offset:offset+length]...)
			*a = append(arr[:offset], append(replacement, arr[offset+length:]...)...)
			return v
		}
	} else {
		if -offset > l {
			return nil
		} else if -offset <= l && l+offset+length < l {
			v := append([]T{}, arr[l+offset:l+offset+length]...)
			*a = append(arr[:l+offset], append(replacement, arr[l+offset+length:]...)...)
			return v
		} else if -offset <= l && l+offset+length >= l {
			v := append([]T{}, arr[l+offset:]...)
			*a = append(arr[:l+offset], replacement...)
			return v
		}
	}
	return nil
}

func Shuffle[T any](a *[]T) {
	if len(*a) < 1 {
		return
	}
	b := make([]T, 0, len(*a))
	for {
		v, l := RandPop(a)
		b = append(b, v)
		if l < 1 {
			break
		}

	}
	*a = b
}

func Delete[T any](a *[]T, index int) {
	if index >= len(*a) || index < 0 {
		return
	}
	arr := *a
	*a = append(arr[:index], arr[index+1:]...)
}

func Copy[T any](a []T, l ...int) []T {
	length := len(a)
	if len(l) > 0 {
		length = l[0]
	}
	var dst []T
	if len(a) < length {
		dst = make([]T, len(a), length)
	} else {
		dst = make([]T, length)
	}
	copy(dst, a)
	return dst
}

func Unshift[T any](a *[]T, e ...T) {
	*a = append(e, *a...)
}

func Push[T any](a *[]T, e ...T) {
	*a = append(*a, e...)
}

func Decompress[T any](a [][]T) (r []T) {
	for _, ts := range a {
		r = append(r, ts...)
	}
	return
}
func DecompressBy[T, R any](a [][]T, fn func(T) (R, bool)) (r []R) {
	for _, ts := range a {
		for _, t := range ts {
			v, ok := fn(t)
			if ok {
				r = append(r, v)
			}
		}
	}
	return
}

func Replace[T any](a []T, offset int, replacement []T) {
	aa := a[offset:]
	aa = aa[:0]
	aa = append(aa, replacement...)
}
