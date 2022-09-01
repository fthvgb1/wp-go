package models

var Options = make(map[string]string)

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

func InitTerms() {
	//terms,err := WpTermsM.SimplePagination()
}
