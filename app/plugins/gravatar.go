package plugins

import (
	"fmt"
	"github.com/fthvgb1/wp-go/app/wpconfig"
	"github.com/fthvgb1/wp-go/helper/number"
	str "github.com/fthvgb1/wp-go/helper/strings"
	"net/url"
	"strings"
)

func Gravatar(email string, isTls bool) (u string) {
	email = strings.Trim(email, " \t\n\r\000\x0B")
	num := number.Rand(0, 2)
	h := ""
	if email != "" {
		h = str.Md5(strings.ToLower(email))
		num = int(h[0] % 3)
	}
	if isTls {
		u = fmt.Sprintf("%s%s", "https://secure.gravatar.com/avatar/", h)
	} else {
		u = fmt.Sprintf("http://%d.gravatar.com/avatar/%s", num, h)
	}
	q := url.Values{}
	q.Add("s", "112")
	q.Add("d", "mm")
	q.Add("r", strings.ToLower(wpconfig.GetOption("avatar_rating")))
	u = fmt.Sprintf("%s?%s", u, q.Encode())
	return
}
