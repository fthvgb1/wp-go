package common

import (
	"context"
	"github.com/leeqvip/gophp"
	"github/fthvgb1/wp-go/helper"
	wp2 "github/fthvgb1/wp-go/internal/wp"
	"github/fthvgb1/wp-go/logs"
	"github/fthvgb1/wp-go/models"
	"strconv"
	"time"
)

func GetPostMetaByPostIds(ctx context.Context, ids []uint64) (r []map[string]any, err error) {
	r, err = postMetaCache.GetCacheBatch(ctx, ids, time.Second, ctx, ids)
	return
}
func GetPostMetaByPostId(ctx context.Context, id uint64) (r map[string]any, err error) {
	r, err = postMetaCache.GetCache(ctx, id, time.Second, ctx, id)
	return
}

func ToPostThumbnail(c context.Context, postId uint64) (r wp2.PostThumbnail) {
	meta, err := GetPostMetaByPostId(c, postId)
	if err == nil {
		m, ok := meta["_thumbnail_id"]
		if ok {
			id, err := strconv.ParseUint(m.(string), 10, 64)
			if err == nil {
				mm, err := GetPostMetaByPostId(c, id)
				if err == nil {
					f, ok := mm["_wp_attached_file"]
					if ok {
						ff, ok := f.(string)
						if ok && ff != "" {
							r.Path = ff
						}
					}
					tt, ok := helper.GetStrMapAnyVal[map[string]any]("_wp_attachment_metadata.sizes.post-thumbnail", mm)
					if ok && tt != nil {
						width, ok := tt["width"]
						if ok {
							w, ok := width.(int)
							if ok {
								r.Width = w
							}
						}
						height, ok := tt["height"]
						if ok {
							h, ok := height.(int)
							if ok {
								r.Height = h
							}
						}
					}
				}
			}
		}
	}
	return
}

func getPostMetaByPostIds(args ...any) (r map[uint64]map[string]any, err error) {
	r = make(map[uint64]map[string]any)
	ctx := args[0].(context.Context)
	ids := args[1].([]uint64)
	rr, err := models.Find[wp2.Postmeta](ctx, models.SqlBuilder{
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
			meta, err := gophp.Unserialize([]byte(postmeta.MetaValue))
			if err != nil {
				logs.ErrPrintln(err, "反序列化postmeta失败", postmeta.MetaValue)
				continue
			}
			metaVal, ok := meta.(map[string]any)
			if ok {
				r[postmeta.PostId][postmeta.MetaKey] = metaVal
			}
		} else {
			r[postmeta.PostId][postmeta.MetaKey] = postmeta.MetaValue
		}

	}
	return
}
