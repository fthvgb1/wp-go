package theme

import (
	"github.com/fthvgb1/wp-go/app/pkg/config"
	"github.com/fthvgb1/wp-go/app/theme/twentyfifteen"
	"github.com/fthvgb1/wp-go/app/theme/twentyseventeen"
	"github.com/fthvgb1/wp-go/app/wpconfig"
)

func InitTheme() {
	AddTheme(twentyfifteen.ThemeName, twentyfifteen.Hook)
	AddTheme(twentyseventeen.ThemeName, twentyseventeen.Hook)
}

func GetCurrentTheme() string {
	themeName := config.GetConfig().Theme
	if themeName == "" {
		themeName = wpconfig.GetOption("template")
	}
	if !IsTemplateDirExists(themeName) {
		themeName = "twentyfifteen"
	}
	return themeName
}
