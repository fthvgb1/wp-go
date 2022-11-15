package config

import (
	"context"
	"github/fthvgb1/wp-go/models"
	"github/fthvgb1/wp-go/models/wp"
	"github/fthvgb1/wp-go/safety"
)

var Options safety.Map[string, string]

func InitOptions() error {
	ctx := context.Background()
	ops, err := models.SimpleFind[wp.Options](ctx, models.SqlBuilder{{"autoload", "yes"}}, "option_name, option_value")
	if err != nil {
		return err
	}
	if len(ops) == 0 {
		ops, err = models.SimpleFind[wp.Options](ctx, nil, "option_name, option_value")
		if err != nil {
			return err
		}
	}
	for _, options := range ops {
		Options.Store(options.OptionName, options.OptionValue)
	}
	return nil
}
