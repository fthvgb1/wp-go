package widget

import (
	"fmt"
	"github.com/fthvgb1/wp-go/app/cmd/reload"
	"github.com/fthvgb1/wp-go/app/pkg/cache"
	"github.com/fthvgb1/wp-go/app/pkg/constraints"
	"github.com/fthvgb1/wp-go/app/pkg/constraints/widgets"
	"github.com/fthvgb1/wp-go/app/pkg/models"
	"github.com/fthvgb1/wp-go/app/theme/wp"
	"github.com/fthvgb1/wp-go/app/wpconfig"
	"github.com/fthvgb1/wp-go/helper/maps"
	"github.com/fthvgb1/wp-go/helper/slice"
	str "github.com/fthvgb1/wp-go/helper/strings"
	"strings"
)

var recentPostsTemplate = `{$before_widget}
{$nav}
{$title}
<ul>
	{$li}
</ul>
{$navCloser}
{$after_widget}
`

func recentPostsArgs() map[string]string {
	return map[string]string{
		"{$before_sidebar}": "",
		"{$after_sidebar}":  "",
		"{$nav}":            "",
		"{$navCloser}":      "",
		"{$title}":          "",
	}
}

func recentConf() map[any]any {
	return map[any]any{
		"number":    int64(5),
		"show_date": false,
		"title":     "近期文章",
	}
}

func RecentPosts(h *wp.Handle, id string) string {
	conf := reload.GetAnyValBys("widget-recent-posts-conf", h, func(h *wp.Handle) map[any]any {
		recent := recentConf()
		conf := wpconfig.GetPHPArrayVal[map[any]any]("widget_recent-posts", recent, int64(2))
		conf = maps.FilterZeroMerge(recent, conf)
		return conf
	})

	args := reload.GetAnyValBys("widget-recent-posts-args", h, func(h *wp.Handle) map[string]string {
		recent := recentPostsArgs()
		commonArgs := wp.GetComponentsArgs(h, widgets.Widget, map[string]string{})
		args := wp.GetComponentsArgs(h, widgets.RecentPosts, recent)
		args = maps.FilterZeroMerge(recent, CommonArgs(), commonArgs, args)
		args["{$before_widget}"] = fmt.Sprintf(args["{$before_widget}"], str.Join("recent-posts-", id), str.Join("widget widget_", "recent_entries"))
		args["{$title}"] = str.Join(args["{$before_title}"], conf["title"].(string), args["{$after_title}"])
		if slice.IsContained(h.CommonThemeMods().ThemeSupport.HTML5, "navigation-widgets") {
			args["{$nav}"] = fmt.Sprintf(`<nav aria-label="%s">`, conf["title"])
			args["{$navCloser}"] = "</nav>"
		}
		return args
	})

	currentPostId := uint64(0)
	if h.Scene() == constraints.Detail {
		currentPostId = str.ToInteger(h.C.Param("id"), uint64(0))
	}
	posts := slice.Map(cache.RecentPosts(h.C, int(conf["number"].(int64))), func(t models.Posts) string {
		t = wp.ProjectTitle(t)
		date := ""
		if v, ok := conf["show_date"].(bool); ok && v {
			date = fmt.Sprintf(`<span class="post-date">%s</span>`, t.PostDateGmt.Format("2006年01月02日"))
		}
		ariaCurrent := ""
		if currentPostId == t.Id {
			ariaCurrent = ` aria-current="page"`
		}
		return fmt.Sprintf(`	<li>
		<a href="/p/%v"%s>%s</a>
		%s
	</li>`, t.Id, ariaCurrent, t.PostTitle, date)
	})
	s := strings.ReplaceAll(recentPostsTemplate, "{$li}", strings.Join(posts, "\n"))
	return h.ComponentFilterFnHook(widgets.RecentPosts, str.Replace(s, args))
}
