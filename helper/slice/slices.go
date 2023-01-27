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
