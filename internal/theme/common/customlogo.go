package common

import (
	"fmt"
	"github.com/fthvgb1/wp-go/helper/maps"
	str "github.com/fthvgb1/wp-go/helper/strings"
	"github.com/fthvgb1/wp-go/internal/cmd/reload"
	"github.com/fthvgb1/wp-go/internal/pkg/cache"
	"github.com/fthvgb1/wp-go/internal/pkg/constraints"
	"github.com/fthvgb1/wp-go/internal/wpconfig"
)

var logo = reload.Vars(constraints.Defaults)

func (h *Handle) CalCustomLogo() (r string) {
	mods, err := wpconfig.GetThemeMods(h.Theme)
	if err != nil {
		return
	}
	id := uint64(mods.CustomLogo)
	if id < 1 {
		id = str.ToInteger[uint64](wpconfig.Options.Value("site_logo"), 0)
		if id < 1 {
			return
		}
	}
	logo, err := cache.GetPostById(h.C, id)
	if err != nil || logo.AttachmentMetadata.File == "" {
		return
	}
	siz := "full"
	meta, _ := cache.GetPostMetaByPostId(h.C, id)
	alt := maps.WithDefaultVal(meta, "_wp_attachment_image_alt", any(wpconfig.Options.Value("blogname")))
	desc := alt.(string)
	imgx := map[string]string{
		"class":    "custom-logo",
		"alt":      desc,
		"decoding": "async",
		//"loading":"lazy",
	}
	img := wpconfig.Thumbnail(logo.AttachmentMetadata, siz, "", "")
	imgx["srcset"] = img.Srcset
	imgx["sizes"] = img.Sizes
	imgx["src"] = img.Path
	r = fmt.Sprintf("%s />", maps.Reduce(imgx, func(k string, v string, t string) string {
		return fmt.Sprintf(`%s %s="%s"`, t, k, v)
	}, fmt.Sprintf(`<img wight="%v" height="%v"`, img.Width, img.Height)))
	r = fmt.Sprintf(`<a href="%s" class="custom-logo-link" rel="home"%s>%s</a>`, "/", ` aria-current="page"`, r)
	return
}

func (h *Handle) CustomLogo() {
	s := logo.Load()
	if s == constraints.Defaults {
		s = h.CalCustomLogo()
		logo.Store(s)
	}
	h.GinH["customLogo"] = s
}
