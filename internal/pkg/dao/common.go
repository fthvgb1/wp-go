package dao

import (
	"context"
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
	var in = []any{"category", "post_tag"}
	terms, err = model.Finds[models.TermsMy](ctx, model.Conditions(
		model.Where(model.SqlBuilder{
			{"tt.count", ">", "0", "int"},
			{"tt.taxonomy", "in", ""},
		}),
		model.Fields("t.term_id"),
		model.Order(model.SqlBuilder{{"t.name", "asc"}}),
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
