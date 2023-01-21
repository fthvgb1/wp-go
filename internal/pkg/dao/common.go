package dao

import (
	"context"
	"fmt"
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

func PasswordProjectTitle(post *models.Posts) {
	if post.PostPassword != "" {
		post.PostTitle = fmt.Sprintf("密码保护：%s", post.PostTitle)
	}
}

func Categories(a ...any) (terms []models.TermsMy, err error) {
	ctx := a[0].(context.Context)
	var in = []any{"category"}
	terms, err = model.Find[models.TermsMy](ctx, model.SqlBuilder{
		{"tt.count", ">", "0", "int"},
		{"tt.taxonomy", "in", ""},
	}, "t.term_id", "", model.SqlBuilder{
		{"t.name", "asc"},
	}, model.SqlBuilder{
		{"t", "inner join", "wp_term_taxonomy tt", "t.term_id = tt.term_id"},
	}, nil, 0, in)
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
	return model.Find[models.PostArchive](ctx, model.SqlBuilder{
		{"post_type", "post"}, {"post_status", "publish"},
	}, "YEAR(post_date) AS `year`, MONTH(post_date) AS `month`, count(ID) as posts", "year,month", model.SqlBuilder{{"year", "desc"}, {"month", "desc"}}, nil, nil, 0)
}
