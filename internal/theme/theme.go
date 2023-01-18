package theme

import (
	"github.com/fthvgb1/wp-go/internal/theme/twentyseventeen"
)

func InitThemeAndTemplateFuncMap() {
	AddThemeHookFunc(twentyseventeen.ThemeName, twentyseventeen.Hook)
}
