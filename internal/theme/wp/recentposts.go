package wp

import (
	"fmt"
	"github.com/fthvgb1/wp-go/helper/maps"
	"github.com/fthvgb1/wp-go/helper/slice"
	str "github.com/fthvgb1/wp-go/helper/strings"
	"github.com/fthvgb1/wp-go/internal/pkg/cache"
	"github.com/fthvgb1/wp-go/internal/pkg/constraints"
	"github.com/fthvgb1/wp-go/internal/pkg/models"
	"github.com/fthvgb1/wp-go/internal/theme/wp/components"
	"github.com/fthvgb1/wp-go/internal/wpconfig"
	"strings"
)

var recentPostsArgs = map[string]string{
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

var recentPostsTemplate = `{$before_widget}
{$nav}
{$title}
<ul>
	{$li}
</ul>
{$navCloser}
{$after_widget}
`

var recentConf = map[any]any{
	"number":    int64(5),
	"show_date": false,
	"title":     "近期文章",
}

func RecentPosts(h *Handle) string {
	args := GetComponentsArgs(h, components.RecentPostsArgs, recentPostsArgs)
	args = maps.Merge(recentPostsArgs, args)
	conf := wpconfig.GetPHPArrayVal[map[any]any]("widget_recent-posts", recentConf, int64(2))
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
		t = ProjectTitle(t)
		date := ""
		if conf["show_date"].(bool) {
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
	return str.Replace(s, args)
}
