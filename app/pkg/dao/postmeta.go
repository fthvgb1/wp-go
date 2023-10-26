package dao

import (
	"context"
	"github.com/fthvgb1/wp-go/app/phphelper"
	"github.com/fthvgb1/wp-go/app/pkg/logs"
	"github.com/fthvgb1/wp-go/app/pkg/models"
	"github.com/fthvgb1/wp-go/app/wpconfig"
	"github.com/fthvgb1/wp-go/helper/slice"
	"github.com/fthvgb1/wp-go/model"
	"strconv"
)

func GetPostMetaByPostIds(ctx context.Context, ids []uint64, _ ...any) (r map[uint64]map[string]any, err error) {
	r = make(map[uint64]map[string]any)
	rr, err := model.Finds[models.PostMeta](ctx, model.Conditions(
		model.Where(model.SqlBuilder{{"post_id", "in", ""}}),
		model.In(slice.ToAnySlice(ids)),
	))
	if err != nil {
		return
	}
	for _, postmeta := range rr {
		if _, ok := r[postmeta.PostId]; !ok {
			r[postmeta.PostId] = make(map[string]any)
		}
		r[postmeta.PostId][postmeta.MetaKey] = postmeta.MetaValue
		if postmeta.MetaKey == "_wp_attachment_metadata" {
			metadata, err := phphelper.UnPHPSerializeToStruct[models.WpAttachmentMetadata](postmeta.MetaValue)
			if err != nil {
				logs.Error(err, "解析postmeta失败", postmeta.MetaId, postmeta.MetaValue)
				continue
			}
			r[postmeta.PostId][postmeta.MetaKey] = metadata
		}
	}
	return
}

func ToPostThumb(c context.Context, meta map[string]any, host string) (r models.PostThumbnail) {
	if meta != nil {
		m, ok := meta["_thumbnail_id"]
		if !ok {
			return
		}
		id, err := strconv.ParseUint(m.(string), 10, 64)
		if err != nil {
			return
		}
		mx, err := GetPostMetaByPostIds(c, []uint64{id})
		if err != nil || mx == nil {
			return
		}
		mm, ok := mx[id]
		if !ok || mm == nil {
			return
		}
		x, ok := mm["_wp_attachment_metadata"]
		if ok {
			metadata, ok := x.(models.WpAttachmentMetadata)
			if ok {
				r = wpconfig.Thumbnail(metadata, "post-thumbnail", host, "thumbnail")
			}
		}
	}
	return
}
