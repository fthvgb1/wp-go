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
	templateName := config.GetConfig().Theme
	if templateName == "" {
		templateName = wpconfig.GetOption("template")
	}
	if !IsTemplateDirExists(templateName) {
		templateName = "twentyfifteen"
	}
	return templateName
}
