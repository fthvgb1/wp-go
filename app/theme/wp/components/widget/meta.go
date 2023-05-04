package widget

import (
	"fmt"
	"github.com/fthvgb1/wp-go/app/cmd/reload"
	"github.com/fthvgb1/wp-go/app/pkg/constraints/widgets"
	"github.com/fthvgb1/wp-go/app/theme/wp"
	"github.com/fthvgb1/wp-go/app/wpconfig"
	"github.com/fthvgb1/wp-go/helper/maps"
	"github.com/fthvgb1/wp-go/helper/slice"
	str "github.com/fthvgb1/wp-go/helper/strings"
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
		"{$aria_label}": "",
		"{$title}":      "",
	}
}

func Meta(h *wp.Handle, id string) string {
	args := reload.GetAnyValBys("widget-meta-args", h, func(h *wp.Handle) map[string]string {
		commonArgs := wp.GetComponentsArgs(h, widgets.Widget, map[string]string{})
		metaArgs := metaArgs()
		args := wp.GetComponentsArgs(h, widgets.Meta, metaArgs)
		args = maps.FilterZeroMerge(metaArgs, CommonArgs(), commonArgs, args)
		args["{$before_widget}"] = fmt.Sprintf(args["{$before_widget}"], str.Join("meta-", id), str.Join("widget widget_", "meta"))
		args["{$title}"] = wpconfig.GetPHPArrayVal("widget_meta", "其它操作", int64(2), "title")
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
		return args
	})

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
