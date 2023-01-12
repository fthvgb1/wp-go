package cache

import (
	"context"
	"github/fthvgb1/wp-go/helper"
	wp2 "github/fthvgb1/wp-go/internal/models"
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
