package wpconfig

import (
	"context"
	"github/fthvgb1/wp-go/internal/pkg/models"
	"github/fthvgb1/wp-go/model"
	"github/fthvgb1/wp-go/safety"
)

var Options safety.Map[string, string]

func InitOptions() error {
	ctx := context.Background()
	ops, err := model.SimpleFind[models.Options](ctx, model.SqlBuilder{{"autoload", "yes"}}, "option_name, option_value")
	if err != nil {
		return err
	}
	if len(ops) == 0 {
		ops, err = model.SimpleFind[models.Options](ctx, nil, "option_name, option_value")
		if err != nil {
			return err
		}
	}
	for _, options := range ops {
		Options.Store(options.OptionName, options.OptionValue)
	}
	return nil
}
