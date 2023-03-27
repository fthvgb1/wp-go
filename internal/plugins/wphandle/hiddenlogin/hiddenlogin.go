package hiddenlogin

import (
	str "github.com/fthvgb1/wp-go/helper/strings"
	"github.com/fthvgb1/wp-go/internal/pkg/constraints/widgets"
	"github.com/fthvgb1/wp-go/internal/theme/wp"
)

func HiddenLogin(h *wp.Handle) {
	h.PushComponentFilterFn(widgets.Meta, func(h *wp.Handle, s string, args ...any) string {
		return str.Replace(s, map[string]string{
			`<li><a href="/wp-login.php">登录</a></li>`:  "",
			`<li><a href="/feed">登录</a></li>`:          "",
			`<li><a href="/comments/feed">登录</a></li>`: "",
		})
	})
}
