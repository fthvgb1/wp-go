package wpconfig

import (
	"context"
	"github/fthvgb1/wp-go/models"
	"github/fthvgb1/wp-go/models/wp"
	"github/fthvgb1/wp-go/safety"
)

var Terms safety.Map[uint64, wp.Terms]
var TermTaxonomies safety.Map[uint64, wp.TermTaxonomy]

func InitTerms() (err error) {
	ctx := context.Background()
	terms, err := models.SimpleFind[wp.Terms](ctx, nil, "*")
	if err != nil {
		return err
	}
	for _, wpTerms := range terms {
		Terms.Store(wpTerms.TermId, wpTerms)
	}
	termTax, err := models.SimpleFind[wp.TermTaxonomy](ctx, nil, "*")
	if err != nil {
		return err
	}
	for _, taxonomy := range termTax {
		TermTaxonomies.Store(taxonomy.TermTaxonomyId, taxonomy)
	}
	return
}
