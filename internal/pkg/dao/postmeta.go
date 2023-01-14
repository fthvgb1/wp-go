package common

import (
	"context"
	"fmt"
	"github/fthvgb1/wp-go/helper"
	"github/fthvgb1/wp-go/internal/pkg/logs"
	"github/fthvgb1/wp-go/internal/pkg/models"
	"github/fthvgb1/wp-go/model"
	"strconv"
	"strings"
)

func GetPostMetaByPostIds(args ...any) (r map[uint64]map[string]any, err error) {
	r = make(map[uint64]map[string]any)
	ctx := args[0].(context.Context)
	ids := args[1].([]uint64)
	rr, err := model.Find[models.Postmeta](ctx, model.SqlBuilder{
		{"post_id", "in", ""},
	}, "*", "", nil, nil, nil, 0, helper.SliceMap(ids, helper.ToAny[uint64]))
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
		if ok {
			id, err := strconv.ParseUint(m.(string), 10, 64)
			if err == nil {
				mx, err := GetPostMetaByPostIds(c, []uint64{id})
				if err == nil && mx != nil {
					mm, ok := mx[id]
					if ok && mm != nil {
						f, ok := mm["_wp_attached_file"]
						if ok {
							ff, ok := f.(string)
							if ok && ff != "" {
								r.Path = ff
							}
						}
						x, ok := mm["_wp_attachment_metadata"]
						if ok {
							metadata, ok := x.(models.WpAttachmentMetadata)
							if ok {
								if _, ok := metadata.Sizes["post-thumbnail"]; ok {
									r.Width = metadata.Sizes["post-thumbnail"].Width
									r.Height = metadata.Sizes["post-thumbnail"].Height
									up := strings.Split(metadata.File, "/")
									r.Srcset = strings.Join(helper.MapToSlice[string](metadata.Sizes, func(s string, size models.MetaDataFileSize) (r string, ok bool) {
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
							}
						}
					}
				}
			}
		}
	}
	return
}
