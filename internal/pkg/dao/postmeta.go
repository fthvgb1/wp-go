package dao

import (
	"context"
	"fmt"
	"github.com/fthvgb1/wp-go/helper"
	"github.com/fthvgb1/wp-go/helper/maps"
	"github.com/fthvgb1/wp-go/helper/slice"
	"github.com/fthvgb1/wp-go/internal/pkg/logs"
	"github.com/fthvgb1/wp-go/internal/pkg/models"
	"github.com/fthvgb1/wp-go/model"
	"strconv"
	"strings"
)

func GetPostMetaByPostIds(args ...any) (r map[uint64]map[string]any, err error) {
	r = make(map[uint64]map[string]any)
	ctx := args[0].(context.Context)
	ids := args[1].([]uint64)
	rr, err := model.Find[models.Postmeta](ctx, model.SqlBuilder{
		{"post_id", "in", ""},
	}, "*", "", nil, nil, nil, 0, slice.Map(ids, helper.ToAny[uint64]))
	if err != nil {
		return
	}
	for _, postmeta := range rr {
		if _, ok := r[postmeta.PostId]; !ok {
			r[postmeta.PostId] = make(map[string]any)
		}
		if postmeta.MetaKey == "_wp_attachment_metadata" {
			metadata, err := models.AttachmentMetadata(postmeta.MetaValue)
			if err != nil {
				logs.ErrPrintln(err, "解析postmeta失败", postmeta.MetaId, postmeta.MetaValue)
				continue
			}
			r[postmeta.PostId][postmeta.MetaKey] = metadata

		} else {
			r[postmeta.PostId][postmeta.MetaKey] = postmeta.MetaValue
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
				r = thumbnail(metadata, "post-thumbnail", host)
			}
		}
	}
	return
}

func thumbnail(metadata models.WpAttachmentMetadata, thumbType, host string) (r models.PostThumbnail) {
	if _, ok := metadata.Sizes[thumbType]; ok {
		r.Path = fmt.Sprintf("%s/wp-content/uploads/%s", host, metadata.File)
		r.Width = metadata.Sizes[thumbType].Width
		r.Height = metadata.Sizes[thumbType].Height
		up := strings.Split(metadata.File, "/")
		r.Srcset = strings.Join(maps.FilterToSlice[string](metadata.Sizes, func(s string, size models.MetaDataFileSize) (r string, ok bool) {
			up[2] = size.File
			if s == "post-thumbnail" {
				return
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
	}
	return
}
