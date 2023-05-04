package dao

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/fthvgb1/wp-go/app/pkg/models"
	"github.com/fthvgb1/wp-go/app/wpconfig"
	"github.com/fthvgb1/wp-go/helper/slice"
	"github.com/fthvgb1/wp-go/model"
	"strings"
	"sync/atomic"
	"time"
)

func GetPostsByIds(a ...any) (m map[uint64]models.Posts, err error) {
	ctx := a[0].(context.Context)
	m = make(map[uint64]models.Posts)
	ids := a[1].([]uint64)
	rawPosts, err := model.Finds[models.Posts](ctx, model.Conditions(
		model.Where(model.SqlBuilder{{"Id", "in", ""}}),
		model.Join(model.SqlBuilder{
			{"a", "left join", "wp_term_relationships b", "a.Id=b.object_id"},
			{"left join", "wp_term_taxonomy c", "b.term_taxonomy_id=c.term_taxonomy_id"},
			{"left join", "wp_terms d", "c.term_id=d.term_id"},
		}),
		model.Fields("a.*,ifnull(d.name,'') category_name,ifnull(c.term_id,0) terms_id,ifnull(taxonomy,'') `taxonomy`"),
		model.In(slice.ToAnySlice(ids)),
	))

	if err != nil {
		return m, err
	}
	postsMap := make(map[uint64]models.Posts)
	for i, post := range rawPosts {
		v, ok := postsMap[post.Id]
		if !ok {
			v = rawPosts[i]
		}
		if post.Taxonomy == "category" {
			v.Categories = append(v.Categories, post.CategoryName)
		} else if post.Taxonomy == "post_tag" {
			v.Tags = append(v.Tags, post.CategoryName)
		}
		if post.TermsId > 0 {
			v.TermIds = append(v.TermIds, post.TermsId)
		}
		postsMap[post.Id] = v
	}
	//host, _ := wpconfig.Options.Load("siteurl")
	host := ""
	meta, _ := GetPostMetaByPostIds(ctx, ids)
	for k, pp := range postsMap {
		if len(pp.Categories) > 0 {
			t := make([]string, 0, len(pp.Categories))
			for _, cat := range pp.Categories {
				t = append(t, fmt.Sprintf(`<a href="/p/category/%s" rel="category tag">%s</a>`, cat, cat))
			}
			pp.CategoriesHtml = strings.Join(t, "、")
		}
		mm, ok := meta[pp.Id]
		if ok {
			attMeta, ok := mm["_wp_attachment_metadata"]
			if ok {
				att, ok := attMeta.(models.WpAttachmentMetadata)
				if ok {
					pp.AttachmentMetadata = att
				}
			}
			if pp.PostType != "attachment" {
				thumb := ToPostThumb(ctx, mm, host)
				if thumb.Path != "" {
					pp.Thumbnail = thumb
				}
			} else if pp.PostType == "attachment" && pp.AttachmentMetadata.File != "" {
				thumb := wpconfig.Thumbnail(pp.AttachmentMetadata, "thumbnail", host, "thumbnail", "post-thumbnail")
				if thumb.Path != "" {
					pp.Thumbnail = thumb
				}
			}

		}
		if len(pp.Tags) > 0 {
			t := make([]string, 0, len(pp.Tags))
			for _, cat := range pp.Tags {
				t = append(t, fmt.Sprintf(`<a href="/p/tag/%s" rel="tag">%s</a>`, cat, cat))
			}
			pp.TagsHtml = strings.Join(t, "、")
		}
		m[k] = pp
	}
	return
}

func SearchPostIds(args ...any) (ids PostIds, err error) {
	ctx := args[0].(context.Context)
	q := args[1].(model.QueryCondition)
	page := args[2].(int)
	pageSize := args[3].(int)
	q.Fields = "ID"
	res, total, err := model.Pagination[models.Posts](ctx, q, page, pageSize)
	for _, posts := range res {
		ids.Ids = append(ids.Ids, posts.Id)
	}
	ids.Length = total
	totalR := int(atomic.LoadInt64(&TotalRaw))
	if total > totalR {
		tt := int64(total)
		atomic.StoreInt64(&TotalRaw, tt)
	}
	return
}

func GetMaxPostId(a ...any) (uint64, error) {
	ctx := a[0].(context.Context)
	r, err := model.SimpleFind[models.Posts](ctx,
		model.SqlBuilder{{"post_type", "post"}, {"post_status", "publish"}},
		"max(ID) ID",
	)
	var id uint64
	if len(r) > 0 {
		id = r[0].Id
	}
	return id, err
}

func RecentPosts(a ...any) (r []models.Posts, err error) {
	ctx := a[0].(context.Context)
	num := a[1].(int)
	r, err = model.Finds[models.Posts](ctx, model.Conditions(
		model.Where(model.SqlBuilder{
			{"post_type", "post"},
			{"post_status", "publish"},
		}),
		model.Fields("ID,post_title,post_password,post_date_gmt"),
		model.Order([][]string{{"post_date", "desc"}}),
		model.Limit(num),
	))
	return
}

func GetPostContext(arg ...any) (r PostContext, err error) {
	ctx := arg[0].(context.Context)
	t := arg[1].(time.Time)
	next, err := model.FirstOne[models.Posts](ctx, model.SqlBuilder{
		{"post_date", ">", t.Format("2006-01-02 15:04:05")},
		{"post_status", "in", ""},
		{"post_type", "post"},
	}, "ID,post_title,post_password", nil, []any{"publish"})
	if err == sql.ErrNoRows {
		err = nil
	}
	if err != nil {
		return
	}
	prev, err := model.FirstOne[models.Posts](ctx, model.SqlBuilder{
		{"post_date", "<", t.Format("2006-01-02 15:04:05")},
		{"post_status", "in", ""},
		{"post_type", "post"},
	}, "ID,post_title", model.SqlBuilder{{"post_date", "desc"}}, []any{"publish"})
	if err == sql.ErrNoRows {
		err = nil
	}
	if err != nil {
		return
	}
	r = PostContext{
		Prev: prev,
		Next: next,
	}
	return
}

func MonthPost(args ...any) (r []uint64, err error) {
	ctx := args[0].(context.Context)
	year, month := args[1].(string), args[2].(string)
	where := model.SqlBuilder{
		{"post_type", "post"},
		{"post_status", "publish"},
		{"year(post_date)", year},
		{"month(post_date)", month},
	}
	r, err = model.Column[models.Posts, uint64](ctx, func(v models.Posts) (uint64, bool) {
		return v.Id, true
	}, model.Conditions(
		model.Fields("ID"),
		model.Where(where),
	))
	l := int64(len(r))
	if l > atomic.LoadInt64(&TotalRaw) {
		atomic.StoreInt64(&TotalRaw, l)
	}
	return
}
