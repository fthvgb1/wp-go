package widget

import (
	"fmt"
	"github.com/fthvgb1/wp-go/helper/maps"
	"github.com/fthvgb1/wp-go/helper/slice"
	str "github.com/fthvgb1/wp-go/helper/strings"
	"github.com/fthvgb1/wp-go/internal/pkg/cache"
	"github.com/fthvgb1/wp-go/internal/pkg/constraints"
	"github.com/fthvgb1/wp-go/internal/pkg/constraints/widgets"
	"github.com/fthvgb1/wp-go/internal/pkg/models"
	"github.com/fthvgb1/wp-go/internal/theme/wp"
	"github.com/fthvgb1/wp-go/internal/wpconfig"
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
		"{$before_widget}":  `<aside id="recent-posts-2" class="widget widget_recent_entries">`,
		"{$after_widget}":   "</aside>",
		"{$before_title}":   `<h2 class="widget-title">`,
		"{$after_title}":    "</h2>",
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

func RecentPosts(h *wp.Handle) string {
	recentPostsArgs := recentPostsArgs()
	recentConf := recentConf()
	args := wp.GetComponentsArgs(h, widgets.RecentPosts, recentPostsArgs)
	args = maps.FilterZeroMerge(recentPostsArgs, args)
	conf := wpconfig.GetPHPArrayVal[map[any]any]("widget_recent-posts", recentConf, int64(2))
	conf = maps.FilterZeroMerge(recentConf, conf)
	if id, ok := args["{$id}"]; ok && id != "" {
		args["{$before_widget}"] = strings.ReplaceAll(args["{$before_widget}"], "2", args["{$id}"])
	}
	args["{$title}"] = str.Join(args["{$before_title}"], conf["title"].(string), args["{$after_title}"])
	if slice.IsContained(h.CommonThemeMods().ThemeSupport.HTML5, "navigation-widgets") {
		args["{$nav}"] = fmt.Sprintf(`<nav aria-label="%s">`, conf["title"])
		args["{$navCloser}"] = "</nav>"
	}
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
