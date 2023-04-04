package wp

import (
	"fmt"
	"github.com/fthvgb1/wp-go/helper/maps"
	str "github.com/fthvgb1/wp-go/helper/strings"
	"github.com/fthvgb1/wp-go/internal/pkg/cache"
	"github.com/fthvgb1/wp-go/internal/wpconfig"
)

func CalCustomLogo(h *Handle) (r string) {
	id := uint64(h.themeMods.CustomLogo)
	if id < 1 {
		id = str.ToInteger[uint64](wpconfig.GetOption("site_logo"), 0)
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
	alt := maps.WithDefaultVal(meta, "_wp_attachment_image_alt", any(wpconfig.GetOption("blogname")))
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
