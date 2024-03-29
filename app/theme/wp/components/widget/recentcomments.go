package widget

import (
	"fmt"
	"github.com/fthvgb1/wp-go/app/pkg/cache"
	"github.com/fthvgb1/wp-go/app/pkg/constraints/widgets"
	"github.com/fthvgb1/wp-go/app/pkg/models"
	"github.com/fthvgb1/wp-go/app/theme/wp"
	"github.com/fthvgb1/wp-go/cache/reload"
	"github.com/fthvgb1/wp-go/helper/maps"
	"github.com/fthvgb1/wp-go/helper/slice"
	str "github.com/fthvgb1/wp-go/helper/strings"
	"strings"
)

func recentCommentsArgs() map[string]string {
	return map[string]string{
		"{$before_sidebar}":     "",
		"{$after_sidebar}":      "",
		"{$nav}":                "",
		"{$navCloser}":          "",
		"{$title}":              "",
		"{$recent_comments_id}": "recentcomments",
	}
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

var GetRecentCommentConf = BuildconfigFn(recentCommentConf, "widget_recent-comments", int64(2))

var GetRecentCommentArgs = reload.BuildValFnWithAnyParams("widget-recent-comment-args", RecentCommentArgs)

func RecentCommentArgs(a ...any) map[string]string {
	h := a[0].(*wp.Handle)
	conf := a[1].(map[any]any)
	id := a[2].(string)
	commentsArgs := recentCommentsArgs()
	commonArgs := wp.GetComponentsArgs(widgets.Widget, map[string]string{})
	args := wp.GetComponentsArgs(widgets.RecentComments, commentsArgs)
	args = maps.FilterZeroMerge(commentsArgs, CommonArgs(), commonArgs, args)
	args["{$before_widget}"] = fmt.Sprintf(args["{$before_widget}"], str.Join("recent-comments-", id), str.Join("widget widget_", "recent_comments"))
	args["{$title}"] = str.Join(args["{$before_title}"], conf["title"].(string), args["{$after_title}"])
	if slice.IsContained(h.CommonThemeMods().ThemeSupport.HTML5, "navigation-widgets") {
		args["{$nav}"] = fmt.Sprintf(`<nav aria-label="%s">`, conf["title"])
		args["{$navCloser}"] = "</nav>"
	}
	return args
}
func RecentComments(h *wp.Handle, id string) string {
	conf := GetRecentCommentConf()
	args := GetRecentCommentArgs(h, conf, id)

	comments := slice.Map(cache.RecentComments(h.C, int(conf["number"].(int64))), func(t models.Comments) string {
		return fmt.Sprintf(`	<li>
<span class="comment-author-link">%s</span>发表在《
		<a href="%s">%s</a>
		》
	</li>`, t.CommentAuthor, t.CommentAuthorUrl, t.PostTitle)
	})
	s := strings.ReplaceAll(recentCommentsTemplate, "{$li}", strings.Join(comments, "\n"))
	return h.DoActionFilter(widgets.RecentComments, str.Replace(s, args))
}
