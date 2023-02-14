package theme

import (
	"github.com/fthvgb1/wp-go/internal/pkg/config"
	"github.com/fthvgb1/wp-go/internal/theme/common"
	"github.com/fthvgb1/wp-go/internal/theme/twentyfifteen"
	"github.com/fthvgb1/wp-go/internal/theme/twentyseventeen"
	"github.com/fthvgb1/wp-go/internal/wpconfig"
)

func InitTheme() {
	addThemeHookFunc(twentyfifteen.ThemeName, twentyfifteen.Hook)
	addThemeHookFunc(twentyseventeen.ThemeName, twentyseventeen.Hook)
	common.AddThemeSupport(twentyfifteen.ThemeName, twentyfifteen.ThemeSupport())
}

func GetTemplateName() string {
	tmlp := config.GetConfig().Theme
	if tmlp == "" {
		tmlp = wpconfig.Options.Value("template")
	}
	if !IsTemplateDirExists(tmlp) {
		tmlp = "twentyfifteen"
	}
	return tmlp
}
