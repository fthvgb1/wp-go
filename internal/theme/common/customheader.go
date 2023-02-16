package common

import (
	"github.com/fthvgb1/wp-go/helper/slice"
	"github.com/fthvgb1/wp-go/internal/pkg/cache"
	"github.com/fthvgb1/wp-go/internal/pkg/logs"
	"github.com/fthvgb1/wp-go/internal/pkg/models"
	"github.com/fthvgb1/wp-go/internal/wpconfig"
)

func (h *Handle) DisplayHeaderText() bool {
	mods, err := wpconfig.GetThemeMods(h.Theme)
	if err != nil {
		return false
	}
	return mods.ThemeSupport.CustomHeader.HeaderText && "blank" != mods.HeaderTextcolor
}

func (h *Handle) GetCustomHeader() (r models.PostThumbnail, isRand bool) {
	hs, err := cache.GetHeaderImages(h.C, h.Theme)
	if err != nil {
		logs.ErrPrintln(err, "获取页眉背景图失败")
		return
	}
	if len(hs) < 1 {
		return
	}
	if len(hs) > 1 {
		isRand = true
	}
	r, _ = slice.RandPop(&hs)
	return
}
