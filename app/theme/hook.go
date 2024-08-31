package theme

import (
	"github.com/fthvgb1/wp-go/app/theme/wp"
	str "github.com/fthvgb1/wp-go/helper/strings"
	"github.com/fthvgb1/wp-go/safety"
)

var themeMap = safety.NewMap[string, func(*wp.Handle)]()

func AddTheme(name string, fn func(handle *wp.Handle)) {
	themeMap.Store(name, fn)
}

func DelTheme(name string) {
	themeMap.Delete(name)
}

func GetTheme(name string) (func(*wp.Handle), bool) {
	return themeMap.Load(name)
}

func IsThemeHookFuncExist(name string) bool {
	_, ok := themeMap.Load(name)
	return ok
}

func Hook(themeName string, h *wp.Handle) {
	fn, ok := themeMap.Load(themeName)
	if ok && fn != nil {
		fn(h)
		return
	}
	panic(str.Join("theme ", themeName, " don't exist"))
}
