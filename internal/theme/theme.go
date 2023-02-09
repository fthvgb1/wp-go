package theme

import (
	"github.com/fthvgb1/wp-go/internal/pkg/config"
	"github.com/fthvgb1/wp-go/internal/theme/twentyfifteen"
	"github.com/fthvgb1/wp-go/internal/theme/twentyseventeen"
	"github.com/fthvgb1/wp-go/internal/wpconfig"
)

func InitThemeAndTemplateFuncMap() {
	addThemeHookFunc(twentyfifteen.ThemeName, twentyfifteen.Hook)
	addThemeHookFunc(twentyseventeen.ThemeName, twentyseventeen.Hook)
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
