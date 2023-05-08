package theme

import (
	"github.com/fthvgb1/wp-go/app/pkg/config"
	"github.com/fthvgb1/wp-go/app/theme/twentyfifteen"
	"github.com/fthvgb1/wp-go/app/theme/twentyseventeen"
	"github.com/fthvgb1/wp-go/app/wpconfig"
)

func InitTheme() {
	AddThemeHookFunc(twentyfifteen.ThemeName, twentyfifteen.Hook)
	AddThemeHookFunc(twentyseventeen.ThemeName, twentyseventeen.Hook)
}

func GetCurrentTemplateName() string {
	tmlp := config.GetConfig().Theme
	if tmlp == "" {
		tmlp = wpconfig.GetOption("template")
	}
	if !IsTemplateDirExists(tmlp) {
		tmlp = "twentyfifteen"
	}
	return tmlp
}
