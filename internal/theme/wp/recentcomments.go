package wp

import (
	"fmt"
	"github.com/fthvgb1/wp-go/helper/maps"
	"github.com/fthvgb1/wp-go/helper/slice"
	str "github.com/fthvgb1/wp-go/helper/strings"
	"github.com/fthvgb1/wp-go/internal/pkg/cache"
	"github.com/fthvgb1/wp-go/internal/pkg/models"
	"github.com/fthvgb1/wp-go/internal/theme/wp/components"
	"github.com/fthvgb1/wp-go/internal/wpconfig"
	"strings"
)

var recentCommentsArgs = map[string]string{
	"{$before_widget}":      `<aside id="recent-comments-2" class="widget widget_recent_comments">`,
	"{$after_widget}":       "</aside>",
	"{$before_title}":       `<h2 class="widget-title">`,
	"{$after_title}":        "</h2>",
	"{$before_sidebar}":     "",
	"{$after_sidebar}":      "",
	"{$nav}":                "",
	"{$navCloser}":          "",
	"{$title}":              "",
	"{$recent_comments_id}": "recentcomments",
}

var recentCommentConf = map[any]any{
	"number": int64(5),
	"title":  "近期评论",
}

var recentCommentsTemplate = `{$before_widget}
{$nav}
{$title}
<ul id="{$recent_comments_id}">
	{$li}
</ul>
{$navCloser}
{$after_widget}
`

func RecentComments(h *Handle) string {
	args := GetComponentsArgs(h, components.RecentCommentsArgs, recentCommentsArgs)
	args = maps.Merge(recentCommentsArgs, args)
	conf := wpconfig.GetPHPArrayVal[map[any]any]("widget_recent-comments", recentCommentConf, int64(2))
	args["{$title}"] = str.Join(args["{$before_title}"], conf["title"].(string), args["{$after_title}"])
	if slice.IsContained(h.CommonThemeMods().ThemeSupport.HTML5, "navigation-widgets") {
		args["{$nav}"] = fmt.Sprintf(`<nav aria-label="%s">`, conf["title"])
		args["{$navCloser}"] = "</nav>"
	}
	comments := slice.Map(cache.RecentComments(h.C, int(conf["number"].(int64))), func(t models.Comments) string {
		return fmt.Sprintf(`	<li>
<span class="comment-author-link">%s</span>发表在《
		<a href="/p/%v#comment-%v">%s</a>
		》
	</li>`, t.CommentAuthor, t.CommentId, t.CommentPostId, t.PostTitle)
	})
	s := strings.ReplaceAll(recentCommentsTemplate, "{$li}", strings.Join(comments, "\n"))
	return str.Replace(s, args)
}
