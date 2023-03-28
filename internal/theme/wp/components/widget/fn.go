package widget

import "github.com/fthvgb1/wp-go/internal/theme/wp"

func Fn(id string, fn func(*wp.Handle, string) string) func(h *wp.Handle) string {
	return func(h *wp.Handle) string {
		return fn(h, id)
	}
}
