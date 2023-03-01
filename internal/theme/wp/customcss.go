package wp

import (
	"fmt"
	"github.com/fthvgb1/wp-go/helper/html"
	"github.com/fthvgb1/wp-go/internal/pkg/cache"
)

func CalCustomCss(h *Handle) (r string) {
	if h.themeMods.CustomCssPostId < 1 {
		return
	}
	post, err := cache.GetPostById(h.C, uint64(h.themeMods.CustomCssPostId))
	if err != nil || post.Id < 1 {
		return
	}
	r = fmt.Sprintf(`<style id="wp-custom-css">%s</style>`, html.StripTags(post.PostContent, ""))
	return
}
