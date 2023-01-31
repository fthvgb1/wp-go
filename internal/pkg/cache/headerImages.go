package cache

import (
	"context"
	"fmt"
	"github.com/fthvgb1/wp-go/helper/slice"
	str "github.com/fthvgb1/wp-go/helper/strings"
	"github.com/fthvgb1/wp-go/internal/pkg/models"
	"github.com/fthvgb1/wp-go/internal/plugins"
	"github.com/fthvgb1/wp-go/internal/wpconfig"
	"github.com/fthvgb1/wp-go/model"
	"time"
)

func GetHeaderImages(ctx context.Context, theme string) (r []models.PostThumbnail, err error) {
	r, err = headerImagesCache.GetCache(ctx, theme, time.Second, ctx, theme)
	return
}

func getHeaderImages(a ...any) (r []models.PostThumbnail, err error) {
	ctx := a[0].(context.Context)
	theme := a[1].(string)
	mods, ok := wpconfig.Options.Load(fmt.Sprintf("theme_mods_%s", theme))
	if ok && mods != "" {
		meta, er := plugins.UnPHPSerialize[plugins.HeaderImageMeta](mods)
		if er != nil {
			err = er
			return
		}
		if "random-uploaded-image" == meta.HeaderImage {
			headers, er := model.Find[models.Posts](ctx, model.SqlBuilder{
				{"post_type", "attachment"},
				{"post_status", "inherit"},
				{"meta_value", theme},
				{"meta_key", "_wp_attachment_is_custom_header"},
			}, "a.ID", "a.ID", nil, model.SqlBuilder{
				{" a", "left join", "wp_postmeta b", "a.ID=b.post_id"},
			}, nil, 0)
			if er != nil {
				err = er
				return
			}
			if len(headers) > 0 {
				posts, er := GetPostsByIds(ctx, slice.Map(headers, func(t models.Posts) uint64 {
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
		m, er := GetPostById(ctx, uint64(meta.HeaderImagData.AttachmentId))
		if er != nil {
			err = er
			return
		}
		m.Thumbnail = thumb(m, theme)
		r = []models.PostThumbnail{m.Thumbnail}
	}
	return
}

func thumb(m models.Posts, theme string) models.PostThumbnail {
	m.Thumbnail = plugins.Thumbnail(m.AttachmentMetadata, "thumbnail", "", "thumbnail", "post-thumbnail", fmt.Sprintf("%s-thumbnail-avatar", theme))
	m.Thumbnail.Width = m.AttachmentMetadata.Width
	m.Thumbnail.Height = m.AttachmentMetadata.Height
	if m.Thumbnail.Path != "" {
		if len(m.AttachmentMetadata.Sizes) > 0 {
			m.Thumbnail.Srcset = str.Join(m.Thumbnail.Path, " 2000w, ", m.Thumbnail.Srcset)
		}
	}
	return m.Thumbnail
}
