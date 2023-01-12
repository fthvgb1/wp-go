package wpconfig

import (
	"context"
	wp2 "github/fthvgb1/wp-go/internal/models"
	"github/fthvgb1/wp-go/model"
	"github/fthvgb1/wp-go/safety"
)

var Terms safety.Map[uint64, wp2.Terms]
var TermTaxonomies safety.Map[uint64, wp2.TermTaxonomy]

func InitTerms() (err error) {
	ctx := context.Background()
	terms, err := model.SimpleFind[wp2.Terms](ctx, nil, "*")
	if err != nil {
		return err
	}
	for _, wpTerms := range terms {
		Terms.Store(wpTerms.TermId, wpTerms)
	}
	termTax, err := model.SimpleFind[wp2.TermTaxonomy](ctx, nil, "*")
	if err != nil {
		return err
	}
	for _, taxonomy := range termTax {
		TermTaxonomies.Store(taxonomy.TermTaxonomyId, taxonomy)
	}
	return
}
