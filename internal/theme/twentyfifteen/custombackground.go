package twentyfifteen

import (
	"fmt"
	"github.com/fthvgb1/wp-go/helper"
	"github.com/fthvgb1/wp-go/helper/maps"
	str "github.com/fthvgb1/wp-go/helper/strings"
	"github.com/fthvgb1/wp-go/internal/cmd/reload"
	"github.com/fthvgb1/wp-go/internal/pkg/constraints"
	"github.com/fthvgb1/wp-go/internal/wpconfig"
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

var backgroud = reload.Vars(constraints.Defaults)

func (h *handle) CustomBackGround() {
	b := backgroud.Load()
	if b == constraints.Defaults {
		b = h.CalCustomBackGround()
		backgroud.Store(b)
	}
	h.IndexHandle.GinH["customBackground"] = b
}

func (h *handle) CalCustomBackGround() (r string) {
	mods, err := wpconfig.GetThemeMods(h.IndexHandle.Theme)
	if err != nil {
		return
	}
	if mods.BackgroundImage == "" && mods.BackgroundColor == themesupport.CustomBackground.DefaultColor {
		return
	}
	s := str.NewBuilder()
	if mods.BackgroundImage != "" {
		s.Sprintf(` background-image: url("%s");`, helper.CutUrlHost(mods.BackgroundImage))
	}
	backgroundPositionX := helper.Defaults(mods.BackgroundPositionX, themesupport.CustomBackground.DefaultPositionX)
	backgroundPositionX = maps.WithDefaultVal(postx, backgroundPositionX, "left")

	backgroundPositionY := helper.Defaults(mods.BackgroundPositionY, themesupport.CustomBackground.DefaultPositionY)
	backgroundPositionY = maps.WithDefaultVal(posty, backgroundPositionY, "top")
	positon := fmt.Sprintf(" background-position: %s %s;", backgroundPositionX, backgroundPositionY)

	siz := helper.DefaultVal(mods.BackgroundSize, themesupport.CustomBackground.DefaultSize)
	siz = maps.WithDefaultVal(size, siz, "auto")
	siz = fmt.Sprintf("  background-size: %s;", siz)

	repeats := helper.Defaults(mods.BackgroundRepeat, themesupport.CustomBackground.DefaultRepeat)
	repeats = maps.WithDefaultVal(repeat, repeats, "repeat")
	repeats = fmt.Sprintf(" background-repeat: %s;", repeats)

	attachment := helper.Defaults(mods.BackgroundAttachment, themesupport.CustomBackground.DefaultAttachment)
	attachment = helper.Or(attachment == "fixed", "fixed", "scroll")
	attachment = fmt.Sprintf(" background-attachment: %s;", attachment)
	s.WriteString(positon, siz, repeats, attachment)
	r = fmt.Sprintf(`<style id="custom-background-css">
body.custom-background {%s}
</style>`, s.String())
	return
}
