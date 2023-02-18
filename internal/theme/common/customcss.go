package common

import (
	"fmt"
	"github.com/fthvgb1/wp-go/helper/html"
	"github.com/fthvgb1/wp-go/internal/cmd/reload"
	"github.com/fthvgb1/wp-go/internal/pkg/cache"
	"github.com/fthvgb1/wp-go/internal/pkg/constraints"
)

var css = reload.Vars(constraints.Defaults)

func (h *Handle) CalCustomCss() (r string) {
	if h.ThemeMods.CustomCssPostId < 1 {
		return
	}
	post, err := cache.GetPostById(h.C, uint64(h.ThemeMods.CustomCssPostId))
	if err != nil || post.Id < 1 {
		return
	}
	r = fmt.Sprintf(`<style id="wp-custom-css">%s</style>`, html.StripTags(post.PostContent, ""))
	return
}

func (h *Handle) CustomCss() {
	cs := css.Load()
	if cs == constraints.Defaults {
		cs = h.CalCustomCss()
		css.Store(cs)
	}
	h.GinH["customCss"] = cs
}
