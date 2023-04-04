package wp

import (
	"fmt"
	"github.com/fthvgb1/wp-go/helper/slice"
	str "github.com/fthvgb1/wp-go/helper/strings"
	"github.com/fthvgb1/wp-go/internal/cmd/reload"
	"github.com/fthvgb1/wp-go/internal/pkg/cache"
	"github.com/fthvgb1/wp-go/internal/pkg/logs"
	"github.com/fthvgb1/wp-go/internal/pkg/models"
	"github.com/fthvgb1/wp-go/internal/wpconfig"
	"github.com/fthvgb1/wp-go/model"
)

func (h *Handle) DisplayHeaderText() bool {
	return h.themeMods.ThemeSupport.CustomHeader.HeaderText && "blank" != h.themeMods.HeaderTextcolor
}

func (h *Handle) GetCustomHeader() (r models.PostThumbnail, isRand bool) {
	var err error
	img := reload.GetAnyValBys("headerImages", h.theme, func(theme string) []models.PostThumbnail {
		hs, er := h.GetHeaderImages(h.theme)
		if er != nil {
			err = er
			return nil
		}
		return hs
	})
	if err != nil {
		logs.ErrPrintln(err, "获取页眉背景图失败")
		return
	}
	hs := slice.Copy(img)

	if len(hs) < 1 {
		return
	}
	if len(hs) > 1 {
		isRand = true
	}
	r, _ = slice.RandPop(&hs)
	return
}

func (h *Handle) GetHeaderImages(theme string) (r []models.PostThumbnail, err error) {
	meta, err := wpconfig.GetThemeMods(theme)
	if err != nil || meta.HeaderImage == "" {
		return
	}
	if "random-uploaded-image" != meta.HeaderImage {
		m, er := cache.GetPostById(h.C, uint64(meta.HeaderImagData.AttachmentId))
		if er != nil {
			err = er
			return
		}
		m.Thumbnail = thumb(m, theme)
		r = []models.PostThumbnail{m.Thumbnail}
		return
	}

	headers, er := model.Finds[models.Posts](h.C, model.Conditions(
		model.Where(model.SqlBuilder{
			{"post_type", "attachment"},
			{"post_status", "inherit"},
			{"meta_value", theme},
			{"meta_key", "_wp_attachment_is_custom_header"},
		}),
		model.Fields("a.ID"),
		model.Group("a.ID"),
		model.Join(model.SqlBuilder{
			{" a", "left join", "wp_postmeta b", "a.ID=b.post_id"},
		}),
	))

	if er != nil {
		err = er
		return
	}
	if len(headers) > 0 {
		posts, er := cache.GetPostsByIds(h.C, slice.Map(headers, func(t models.Posts) uint64 {
			return t.Id
		}))
		if er != nil {
			err = er
			return
		}
		r = slice.Map(posts, func(m models.Posts) models.PostThumbnail {
			return thumb(m, theme)
		})
	}
	return

}

func thumb(m models.Posts, theme string) models.PostThumbnail {
	m.Thumbnail = wpconfig.Thumbnail(m.AttachmentMetadata, "thumbnail", "", "thumbnail", "post-thumbnail", fmt.Sprintf("%s-thumbnail-avatar", theme))
	m.Thumbnail.Width = m.AttachmentMetadata.Width
	m.Thumbnail.Height = m.AttachmentMetadata.Height
	if m.Thumbnail.Path != "" {
		if len(m.AttachmentMetadata.Sizes) > 0 {
			m.Thumbnail.Srcset = str.Join(m.Thumbnail.Path, " 2000w, ", m.Thumbnail.Srcset)
		}
	}
	return m.Thumbnail
}
