package dao

import (
	"context"
	"github.com/fthvgb1/wp-go/helper"
	"github.com/fthvgb1/wp-go/internal/pkg/constraints"
	"github.com/fthvgb1/wp-go/internal/pkg/models"
	"github.com/fthvgb1/wp-go/internal/wpconfig"
	"github.com/fthvgb1/wp-go/model"
)

var TotalRaw int64

type PostIds struct {
	Ids    []uint64
	Length int
}

type PostContext struct {
	Prev models.Posts
	Next models.Posts
}

func CategoriesAndTags(a ...any) (terms []models.TermsMy, err error) {
	ctx := a[0].(context.Context)
	t, ok := a[1].(int)
	var in = []any{"category", "post_tag"}
	if ok {
		switch t {
		case constraints.Category:
			in = []any{"category"}
		case constraints.Tag:
			in = []any{"post_tag"}
		}
	}
	w := model.SqlBuilder{
		{"tt.taxonomy", "in", ""},
	}
	if helper.GetContextVal(ctx, "showOnlyTopLevel", false) {
		w = append(w, []string{"tt.parent", "=", "0", "int"})
	}
	if !helper.GetContextVal(ctx, "showEmpty", false) {
		w = append(w, []string{"tt.count", ">", "0", "int"})
	}
	order := []string{"name", "asc"}
	ord := helper.GetContextVal[[]string](ctx, "order", nil)
	if ord != nil {
		order = ord
	}
	terms, err = model.Finds[models.TermsMy](ctx, model.Conditions(
		model.Where(w),
		model.Fields("t.term_id"),
		model.Order(model.SqlBuilder{order}),
		model.Join(model.SqlBuilder{
			{"t", "inner join", "wp_term_taxonomy tt", "t.term_id = tt.term_id"},
		}),
		model.In(in),
	))
	for i := 0; i < len(terms); i++ {
		if v, ok := wpconfig.Terms.Load(terms[i].Terms.TermId); ok {
			terms[i].Terms = v
		}
		if v, ok := wpconfig.TermTaxonomies.Load(terms[i].Terms.TermId); ok {
			terms[i].TermTaxonomy = v
		}
	}
	return
}

func Archives(ctx context.Context) ([]models.PostArchive, error) {
	return model.Finds[models.PostArchive](ctx, model.Conditions(
		model.Where(model.SqlBuilder{
			{"post_type", "post"},
			{"post_status", "publish"},
		}),
		model.Fields("YEAR(post_date) AS `year`, MONTH(post_date) AS `month`, count(ID) as posts"),
		model.Group("year,month"),
		model.Order(model.SqlBuilder{{"year", "desc"}, {"month", "desc"}}),
	))
}
