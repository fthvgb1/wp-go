package widget

import (
	"fmt"
	"github.com/fthvgb1/wp-go/helper/maps"
	"github.com/fthvgb1/wp-go/helper/slice"
	str "github.com/fthvgb1/wp-go/helper/strings"
	"github.com/fthvgb1/wp-go/internal/cmd/reload"
	"github.com/fthvgb1/wp-go/internal/pkg/constraints/widgets"
	"github.com/fthvgb1/wp-go/internal/theme/wp"
	"github.com/fthvgb1/wp-go/internal/wpconfig"
	"strings"
)

var metaTemplate = `{$before_widget}
{$h2title}
{$nav}
<ul>
	{$li}
</ul>
{$navCloser}
{$after_widget}`

func metaArgs() map[string]string {
	return map[string]string{
		"{$before_widget}": `<aside id="meta-2" class="widget widget_meta">`,
		"{$after_widget}":  `</aside>`,
		"{$aria_label}":    "",
		"{$title}":         "",
		"":                 "",
		"{$before_title}":  `<h2 class="widget-title">`,
		"{$after_title}":   `</h2>`,
	}
}

func Meta(h *wp.Handle) string {
	args := reload.GetAnyValBys("widget-meta-args", h, func(h *wp.Handle) map[string]string {
		metaArgs := metaArgs()
		args := wp.GetComponentsArgs(h, widgets.Meta, metaArgs)
		args = maps.FilterZeroMerge(metaArgs, args)
		return args
	})
	args["{$title}"] = wpconfig.GetPHPArrayVal("widget_meta", "其它操作", int64(2), "title")
	if id, ok := args["{$id}"]; ok && id != "" {
		args["{$before_widget}"] = strings.ReplaceAll(args["{$before_widget}"], "2", args["{$id}"])
	}
	if args["{$title}"] == "" {
		args["{$title}"] = "其他操作"
	}
	if args["{$title}"] != "" {
		args["{$h2title}"] = str.Join(args["{$before_title}"], args["{$title}"], args["{$after_title}"])
	}
	if slice.IsContained(h.CommonThemeMods().ThemeSupport.HTML5, "navigation-widgets") {
		args["{$nav}"] = fmt.Sprintf(`<nav aria-label="%s">`, args["{$title}"])
		args["{$navCloser}"] = "</nav>"
	}
	ss := str.NewBuilder()
	if str.ToInteger(wpconfig.GetOption("users_can_register"), 0) > 0 {
		ss.Sprintf(`<li><a href="/wp-login.php?action=register">注册</li>`)
	}
	ss.Sprintf(`<li><a href="%s">登录</a></li>`, "/wp-login.php")
	ss.Sprintf(`<li><a href="%s">条目feed</a></li>`, "/feed")
	ss.Sprintf(`<li><a href="%s">评论feed</a></li>`, "/comments/feed")
	s := strings.ReplaceAll(metaTemplate, "{$li}", ss.String())
	return h.ComponentFilterFnHook(widgets.Meta, str.Replace(s, args))
}
