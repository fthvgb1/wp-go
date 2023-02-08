package slice

func Splice[T any](a *[]T, offset, length int, replacement []T) []T {
	arr := *a
	l := len(arr)
	if length < 0 {
		panic("length must > 0")
	}
	if offset >= 0 {
		if offset+length > l {
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

func Copy[T any](a []T) []T {
	dst := make([]T, len(a))
	copy(dst, a)
	return dst
}
