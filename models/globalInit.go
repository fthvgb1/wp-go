package models

var Options = make(map[string]string)
var TermsIds []uint64

func InitOptions() error {
	ops, err := SimpleFind[WpOptions](SqlBuilder{{"autoload", "yes"}}, "option_name, option_value")
	if err != nil {
		return err
	}
	if len(ops) == 0 {
		ops, err = SimpleFind[WpOptions](nil, "option_name, option_value")
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
	var themes []interface{}
	themes = append(themes, "wp_theme")
	var name []interface{}
	name = append(name, "twentyfifteen")
	terms, err := Find[WpTerms](SqlBuilder{{
		"tt.taxonomy", "in", "",
	}, {"t.name", "in", ""}}, "t.term_id", nil, SqlBuilder{{
		"t", "inner join", "wp_term_taxonomy tt", "t.term_id = tt.term_id",
	}}, 1, themes, name)
	for _, wpTerms := range terms {
		TermsIds = append(TermsIds, wpTerms.TermId)
	}
	return
}
