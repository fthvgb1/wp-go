package wpconfig

import (
	"github.com/fthvgb1/wp-go/app/pkg/models"
	"github.com/fthvgb1/wp-go/model"
	"github.com/fthvgb1/wp-go/safety"
)

var terms = safety.NewMap[uint64, models.Terms]()
var termTaxonomies = safety.NewMap[uint64, models.TermTaxonomy]()

var my = safety.NewMap[uint64, models.TermsMy]()

func GetTerm(termId uint64) (models.Terms, bool) {
	return terms.Load(termId)
}

func GetTermTaxonomy(termId uint64) (models.TermTaxonomy, bool) {
	return termTaxonomies.Load(termId)
}
func GetTermMy(termId uint64) (models.TermsMy, bool) {
	return my.Load(termId)
}

func InitTerms() (err error) {
	terms.Flush()
	termTaxonomies.Flush()
	term, err := model.SimpleFind[models.Terms](ctx, nil, "*")
	if err != nil {
		return err
	}
	for _, wpTerms := range term {
		terms.Store(wpTerms.TermId, wpTerms)
	}
	termTax, err := model.SimpleFind[models.TermTaxonomy](ctx, nil, "*")
	if err != nil {
		return err
	}
	for _, taxonomy := range termTax {
		termTaxonomies.Store(taxonomy.TermTaxonomyId, taxonomy)
		if term, ok := terms.Load(taxonomy.TermId); ok {
			my.Store(taxonomy.TermId, models.TermsMy{
				Terms:        term,
				TermTaxonomy: taxonomy,
			})
		}
	}
	return
}
