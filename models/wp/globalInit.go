package wp

import (
	"context"
	"github/fthvgb1/wp-go/models"
)

var Option = make(map[string]string)
var Terms = map[uint64]WpTerms{}
var TermTaxonomies = map[uint64]TermTaxonomy{}

func InitOptions() error {
	ctx := context.Background()
	ops, err := models.SimpleFind[Options](ctx, models.SqlBuilder{{"autoload", "yes"}}, "option_name, option_value")
	if err != nil {
		return err
	}
	if len(ops) == 0 {
		ops, err = models.SimpleFind[Options](ctx, nil, "option_name, option_value")
		if err != nil {
			return err
		}
	}
	for _, options := range ops {
		Option[options.OptionName] = options.OptionValue
	}
	return nil
}

func InitTerms() (err error) {
	ctx := context.Background()
	terms, err := models.SimpleFind[WpTerms](ctx, nil, "*")
	if err != nil {
		return err
	}
	for _, wpTerms := range terms {
		Terms[wpTerms.TermId] = wpTerms
	}
	termTax, err := models.SimpleFind[TermTaxonomy](ctx, nil, "*")
	if err != nil {
		return err
	}
	for _, taxonomy := range termTax {
		TermTaxonomies[taxonomy.TermTaxonomyId] = taxonomy
	}
	return
}
