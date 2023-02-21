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
	ops, err := model.FindToStringMap[models.Options](ctx, model.Conditions(
		model.Where(model.SqlBuilder{{"autoload", "yes"}}),
		model.Fields("option_name k, option_value v"),
	))
	if err != nil {
		return err
	}
	for _, option := range ops {
		options.Store(option["k"], option["v"])
	}
	return nil
}

func GetOption(k string) string {
	v, ok := options.Load(k)
	if ok {
		return v
	}
	vv, err := model.GetField[models.Options](ctx, "option_value", model.Conditions(model.Where(model.SqlBuilder{{"option_name", k}})))
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
