package common

import (
	"fmt"
	"github.com/fthvgb1/wp-go/helper"
	"github.com/fthvgb1/wp-go/helper/maps"
	"github.com/fthvgb1/wp-go/internal/wpconfig"
	"github.com/fthvgb1/wp-go/safety"
	"strings"
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

var backgroud = safety.NewVar("")

func (h *Handle) CustomBackGround() {
	b := backgroud.Load()
	if b == "" {
		b = h.CalCustomBackGround()
		backgroud.Store(b)
	}
	h.GinH["customBackground"] = b
}

func (h *Handle) CalCustomBackGround() (r string) {
	mods, err := wpconfig.GetThemeMods(h.Theme)
	if err != nil {
		return
	}
	if mods.BackgroundImage == "" && mods.BackgroundColor == mods.ThemeSupport.CustomBackground.DefaultColor {
		return
	}
	s := strings.Builder{}
	if mods.BackgroundImage != "" {
		s.WriteString(` background-image: url("`)
		s.WriteString(helper.CutUrlHost(mods.BackgroundImage))
		s.WriteString(`");`)
	}
	backgroundPositionX := helper.Defaults(mods.BackgroundPositionX, mods.ThemeSupport.CustomBackground.DefaultPositionX)
	backgroundPositionX = maps.WithDefaultVal(postx, backgroundPositionX, "left")

	backgroundPositionY := helper.Defaults(mods.BackgroundPositionY, mods.ThemeSupport.CustomBackground.DefaultPositionY)
	backgroundPositionY = maps.WithDefaultVal(posty, backgroundPositionY, "top")
	positon := fmt.Sprintf(" background-position: %s %s;", backgroundPositionX, backgroundPositionY)

	siz := helper.DefaultVal(mods.BackgroundSize, mods.ThemeSupport.CustomBackground.DefaultSize)
	siz = maps.WithDefaultVal(size, siz, "auto")
	siz = fmt.Sprintf("  background-size: %s;", siz)

	repeats := helper.Defaults(mods.BackgroundRepeat, mods.ThemeSupport.CustomBackground.DefaultRepeat)
	repeats = maps.WithDefaultVal(repeat, repeats, "repeat")
	repeats = fmt.Sprintf(" background-repeat: %s;", repeats)

	attachment := helper.Defaults(mods.BackgroundAttachment, mods.ThemeSupport.CustomBackground.DefaultAttachment)
	attachment = helper.Or(attachment == "fixed", "fixed", "scroll")
	attachment = fmt.Sprintf(" background-attachment: %s;", attachment)
	s.WriteString(positon)
	s.WriteString(siz)
	s.WriteString(repeats)
	s.WriteString(attachment)
	r = fmt.Sprintf(`<style id="custom-background-css">
body.custom-background {%s}
</style>`, s.String())
	return
}
