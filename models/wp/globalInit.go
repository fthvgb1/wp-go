package wp

import "github/fthvgb1/wp-go/models"

var Options = make(map[string]string)
var Terms = map[uint64]WpTerms{}
var TermTaxonomy = map[uint64]WpTermTaxonomy{}

func InitOptions() error {
	ops, err := models.SimpleFind[WpOptions](models.SqlBuilder{{"autoload", "yes"}}, "option_name, option_value")
	if err != nil {
		return err
	}
	if len(ops) == 0 {
		ops, err = models.SimpleFind[WpOptions](nil, "option_name, option_value")
		if err != nil {
			return err
		}
	}
	for _, options := range ops {
		Options[options.OptionName] = options.OptionValue
	}
	return nil
}

func InitTerms() (err error) {
	terms, err := models.SimpleFind[WpTerms](nil, "*")
	if err != nil {
		return err
	}
	for _, wpTerms := range terms {
		Terms[wpTerms.TermId] = wpTerms
	}
	termTax, err := models.SimpleFind[WpTermTaxonomy](nil, "*")
	if err != nil {
		return err
	}
	for _, taxonomy := range termTax {
		TermTaxonomy[taxonomy.TermTaxonomyId] = taxonomy
	}
	return
}
