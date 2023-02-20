package wpconfig

import (
	"context"
	"github.com/fthvgb1/wp-go/internal/pkg/models"
	"github.com/fthvgb1/wp-go/model"
	"github.com/fthvgb1/wp-go/safety"
	"strings"
)

var options safety.Map[string, string]

var ctx context.Context

func InitOptions() error {
	options.Flush()
	if ctx == nil {
		ctx = context.Background()
	}
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
	for _, option := range ops {
		options.Store(option.OptionName, option.OptionValue)
	}
	return nil
}

func GetOption(k string) string {
	v, ok := options.Load(k)
	if ok {
		return v
	}
	vv, err := model.GetField[models.Options, string](ctx, "option_value", model.Conditions(model.Where(model.SqlBuilder{{"option_name", k}})))
	options.Store(k, vv)
	if err != nil {
		return ""
	}
	return vv
}

func GetLang() string {
	s, ok := options.Load("WPLANG")
	if !ok {
		s = "zh-CN"
	}
	return strings.Replace(s, "_", "-", 1)
}
