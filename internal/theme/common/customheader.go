package common

import (
	"github.com/fthvgb1/wp-go/helper/slice"
	"github.com/fthvgb1/wp-go/internal/pkg/cache"
	"github.com/fthvgb1/wp-go/internal/pkg/logs"
	"github.com/fthvgb1/wp-go/internal/pkg/models"
)

func (h *Handle) DisplayHeaderText() bool {
	return h.ThemeMods.ThemeSupport.CustomHeader.HeaderText && "blank" != h.ThemeMods.HeaderTextcolor
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
