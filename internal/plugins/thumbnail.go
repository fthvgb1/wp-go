package plugins

import (
	"fmt"
	"github.com/fthvgb1/wp-go/helper/maps"
	"github.com/fthvgb1/wp-go/internal/pkg/models"
	"strings"
)

func Thumbnail(metadata models.WpAttachmentMetadata, Type, host string, except ...string) (r models.PostThumbnail) {
	if _, ok := metadata.Sizes[Type]; ok {
		r.Path = fmt.Sprintf("%s/wp-content/uploads/%s", host, metadata.File)
		r.Width = metadata.Sizes[Type].Width
		r.Height = metadata.Sizes[Type].Height
		up := strings.Split(metadata.File, "/")
		r.Srcset = strings.Join(maps.FilterToSlice[string](metadata.Sizes, func(s string, size models.MetaDataFileSize) (r string, ok bool) {
			up[2] = size.File
			for _, s2 := range except {
				if s == s2 {
					return
				}
			}
			r = fmt.Sprintf("%s/wp-content/uploads/%s %dw", host, strings.Join(up, "/"), size.Width)
			ok = true
			return
		}), ", ")
		r.Sizes = fmt.Sprintf("(max-width: %dpx) 100vw, %dpx", r.Width, r.Width)
		if r.Width >= 740 && r.Width < 767 {
			r.Sizes = "(max-width: 706px) 89vw, (max-width: 767px) 82vw, 740px"
		} else if r.Width >= 767 {
			r.Sizes = "(max-width: 767px) 89vw, (max-width: 1000px) 54vw, (max-width: 1071px) 543px, 580px"
		}
		r.OriginAttachmentData = metadata
	}
	return
}
