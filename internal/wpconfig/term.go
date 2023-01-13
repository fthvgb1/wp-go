package wpconfig

import (
	"context"
	"github/fthvgb1/wp-go/internal/pkg/models"
	"github/fthvgb1/wp-go/model"
	"github/fthvgb1/wp-go/safety"
)

var Terms safety.Map[uint64, models.Terms]
var TermTaxonomies safety.Map[uint64, models.TermTaxonomy]

func InitTerms() (err error) {
	ctx := context.Background()
	terms, err := model.SimpleFind[models.Terms](ctx, nil, "*")
	if err != nil {
		return err
	}
	for _, wpTerms := range terms {
		Terms.Store(wpTerms.TermId, wpTerms)
	}
	termTax, err := model.SimpleFind[models.TermTaxonomy](ctx, nil, "*")
	if err != nil {
		return err
	}
	for _, taxonomy := range termTax {
		TermTaxonomies.Store(taxonomy.TermTaxonomyId, taxonomy)
	}
	return
}
