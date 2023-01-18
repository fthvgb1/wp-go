package common

import (
	"context"
	"fmt"
	models2 "github.com/fthvgb1/wp-go/internal/pkg/models"
	"github.com/fthvgb1/wp-go/internal/wpconfig"
	"github.com/fthvgb1/wp-go/model"
)

var TotalRaw int64

type PostIds struct {
	Ids    []uint64
	Length int
}

type PostContext struct {
	Prev models2.Posts
	Next models2.Posts
}

func PasswordProjectTitle(post *models2.Posts) {
	if post.PostPassword != "" {
		post.PostTitle = fmt.Sprintf("密码保护：%s", post.PostTitle)
	}
}

func Categories(a ...any) (terms []models2.TermsMy, err error) {
	ctx := a[0].(context.Context)
	var in = []any{"category"}
	terms, err = model.Find[models2.TermsMy](ctx, model.SqlBuilder{
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

func Archives(ctx context.Context) ([]models2.PostArchive, error) {
	return model.Find[models2.PostArchive](ctx, model.SqlBuilder{
		{"post_type", "post"}, {"post_status", "publish"},
	}, "YEAR(post_date) AS `year`, MONTH(post_date) AS `month`, count(ID) as posts", "year,month", model.SqlBuilder{{"year", "desc"}, {"month", "desc"}}, nil, nil, 0)
}
