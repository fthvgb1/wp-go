package theme

import (
	"github.com/fthvgb1/wp-go/internal/pkg/config"
	"github.com/fthvgb1/wp-go/internal/theme/twentyseventeen"
	"github.com/fthvgb1/wp-go/internal/wpconfig"
)

func InitThemeAndTemplateFuncMap() {
	AddThemeHookFunc(twentyseventeen.ThemeName, twentyseventeen.Hook)
}

func GetTemplateName() string {
	tmlp := config.Conf.Load().Theme
	if tmlp == "" {
		tmlp = wpconfig.Options.Value("template")
	}
	if !IsTemplateDirExists(tmlp) {
		tmlp = "twentyfifteen"
	}
	return tmlp
}
