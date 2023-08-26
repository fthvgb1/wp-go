package hiddenlogin

import (
	"github.com/fthvgb1/wp-go/app/pkg/constraints/widgets"
	"github.com/fthvgb1/wp-go/app/theme/wp"
	str "github.com/fthvgb1/wp-go/helper/strings"
)

func HiddenLogin(h *wp.Handle) {
	h.AddActionFilter(widgets.Meta, func(h *wp.Handle, s string, args ...any) string {
		return str.Replace(s, map[string]string{
			`<li><a href="/wp-login.php">登录</a></li>`:  "",
			`<li><a href="/feed">登录</a></li>`:          "",
			`<li><a href="/comments/feed">登录</a></li>`: "",
		})
	})
}
