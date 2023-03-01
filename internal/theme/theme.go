package theme

import (
	"github.com/fthvgb1/wp-go/internal/pkg/config"
	"github.com/fthvgb1/wp-go/internal/theme/twentyfifteen"
	"github.com/fthvgb1/wp-go/internal/theme/twentyseventeen"
	"github.com/fthvgb1/wp-go/internal/wpconfig"
)

func InitTheme() {
	addThemeHookFunc(twentyfifteen.ThemeName, twentyfifteen.Hook)
	twentyfifteen.Init(TemplateFs)
	addThemeHookFunc(twentyseventeen.ThemeName, twentyseventeen.Hook)
	twentyseventeen.Init(TemplateFs)
}

func GetTemplateName() string {
	tmlp := config.GetConfig().Theme
	if tmlp == "" {
		tmlp = wpconfig.GetOption("template")
	}
	if !IsTemplateDirExists(tmlp) {
		tmlp = "twentyfifteen"
	}
	return tmlp
}
