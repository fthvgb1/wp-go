package twentyfifteen

import (
	"fmt"
	"github.com/fthvgb1/wp-go/helper"
	"github.com/fthvgb1/wp-go/helper/maps"
	str "github.com/fthvgb1/wp-go/helper/strings"
	"github.com/fthvgb1/wp-go/internal/theme/common"
)

var postx = map[string]string{
	"left":   "left",
	"right":  "right",
	"center": "center",
}
var posty = map[string]string{
	"top":    "top",
	"bottom": "bottom",
	"center": "center",
}
var size = map[string]string{
	"auto":    "auto",
	"contain": "contain",
	"cover":   "cover",
}
var repeat = map[string]string{
	"repeat-x":  "repeat-x",
	"repeat-y":  "repeat-y",
	"repeat":    "repeat",
	"no-repeat": "no-repeat",
}

func CalCustomBackGround(h *common.Handle) (r string) {
	if h.ThemeMods.BackgroundImage == "" && h.ThemeMods.BackgroundColor == themesupport.CustomBackground.DefaultColor {
		return
	}
	s := str.NewBuilder()
	if h.ThemeMods.BackgroundImage != "" {
		s.Sprintf(` background-image: url("%s");`, helper.CutUrlHost(h.ThemeMods.BackgroundImage))
	}
	backgroundPositionX := helper.Defaults(h.ThemeMods.BackgroundPositionX, themesupport.CustomBackground.DefaultPositionX)
	backgroundPositionX = maps.WithDefaultVal(postx, backgroundPositionX, "left")

	backgroundPositionY := helper.Defaults(h.ThemeMods.BackgroundPositionY, themesupport.CustomBackground.DefaultPositionY)
	backgroundPositionY = maps.WithDefaultVal(posty, backgroundPositionY, "top")
	positon := fmt.Sprintf(" background-position: %s %s;", backgroundPositionX, backgroundPositionY)

	siz := helper.DefaultVal(h.ThemeMods.BackgroundSize, themesupport.CustomBackground.DefaultSize)
	siz = maps.WithDefaultVal(size, siz, "auto")
	siz = fmt.Sprintf("  background-size: %s;", siz)

	repeats := helper.Defaults(h.ThemeMods.BackgroundRepeat, themesupport.CustomBackground.DefaultRepeat)
	repeats = maps.WithDefaultVal(repeat, repeats, "repeat")
	repeats = fmt.Sprintf(" background-repeat: %s;", repeats)

	attachment := helper.Defaults(h.ThemeMods.BackgroundAttachment, themesupport.CustomBackground.DefaultAttachment)
	attachment = helper.Or(attachment == "fixed", "fixed", "scroll")
	attachment = fmt.Sprintf(" background-attachment: %s;", attachment)
	s.WriteString(positon, siz, repeats, attachment)
	r = fmt.Sprintf(`<style id="custom-background-css">
body.custom-background {%s}
</style>`, s.String())
	return
}
