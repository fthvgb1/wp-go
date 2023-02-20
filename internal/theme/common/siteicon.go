package common

import (
	"fmt"
	"github.com/fthvgb1/wp-go/helper/slice"
	str "github.com/fthvgb1/wp-go/helper/strings"
	"github.com/fthvgb1/wp-go/internal/pkg/cache"
	"github.com/fthvgb1/wp-go/internal/wpconfig"
	"strings"
)

var sizes = []string{"site_icon-270", "site_icon-32", "site_icon-192", "site_icon-180"}

func (h *Handle) CalSiteIcon() (r string) {
	id := str.ToInteger[uint64](wpconfig.GetOption("site_icon"), 0)
	if id < 1 {
		return
	}
	icon, err := cache.GetPostById(h.C, id)
	if err != nil || icon.AttachmentMetadata.File == "" {
		return
	}
	m := strings.Join(strings.Split(icon.AttachmentMetadata.File, "/")[:2], "/")
	size := slice.FilterAndMap(sizes, func(t string) (string, bool) {
		s, ok := icon.AttachmentMetadata.Sizes[t]
		if !ok {
			return "", false
		}
		switch t {
		case "site_icon-270":
			return fmt.Sprintf(`<meta name="msapplication-TileImage" content="/wp-content/uploads/%s/%s" />`, m, s.File), true
		case "site_icon-180":
			return fmt.Sprintf(`<link rel="apple-touch-icon" href="/wp-content/uploads/%s/%s" />`, m, s.File), true
		default:
			ss := strings.Replace(t, "site_icon-", "", 1)
			return fmt.Sprintf(`<link rel="icon" href="/wp-content/uploads/%s/%s" sizes="%sx%s" />`, m, s.File, ss, ss), true
		}
	})
	r = strings.Join(size, "\n")
	return
}
