package cache

import (
	"context"
	"github.com/fthvgb1/wp-go/app/pkg/logs"
	"github.com/fthvgb1/wp-go/app/pkg/models"
	"github.com/fthvgb1/wp-go/cache/cachemanager"
	"github.com/fthvgb1/wp-go/helper/slice"
	"time"
)

// CategoriesTags  get all categories or tags
//
// query func see dao.CategoriesAndTags
//
// t is constraints.Tag or constraints.Category
func CategoriesTags(ctx context.Context, t ...string) []models.TermsMy {
	tt := ""
	if len(t) > 0 {
		tt = t[0]
	}
	r, err := cachemanager.GetBy[[]models.TermsMy]("categoryAndTagsData", ctx, tt, time.Second)
	logs.IfError(err, "get category fail")
	return r
}
func AllCategoryTagsNames(ctx context.Context, t ...string) map[string]struct{} {
	tt := ""
	if len(t) > 0 {
		tt = t[0]
	}
	r, err := cachemanager.GetBy[[]models.TermsMy]("categoryAndTagsData", ctx, tt, time.Second)
	if err != nil {
		logs.Error(err, "get category fail")
		return nil
	}
	return slice.ToMap(r, func(t models.TermsMy) (string, struct{}) {
		return t.Name, struct{}{}
	}, true)
}
