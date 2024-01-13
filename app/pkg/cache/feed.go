package cache

import (
	"context"
	"errors"
	"fmt"
	"github.com/fthvgb1/wp-go/app/pkg/logs"
	"github.com/fthvgb1/wp-go/app/pkg/models"
	"github.com/fthvgb1/wp-go/app/plugins"
	"github.com/fthvgb1/wp-go/app/plugins/wpposts"
	"github.com/fthvgb1/wp-go/app/wpconfig"
	"github.com/fthvgb1/wp-go/cache"
	"github.com/fthvgb1/wp-go/cache/cachemanager"
	"github.com/fthvgb1/wp-go/helper/slice"
	str "github.com/fthvgb1/wp-go/helper/strings"
	"github.com/fthvgb1/wp-go/plugin/digest"
	"github.com/fthvgb1/wp-go/rss2"
	"strings"
	"time"
)

var timeFormat = "Mon, 02 Jan 2006 15:04:05 +0000"
var templateRss rss2.Rss2

func InitFeed() {
	templateRss = rss2.Rss2{
		Title:           wpconfig.GetOption("blogname"),
		AtomLink:        fmt.Sprintf("%s/feed", wpconfig.GetOption("home")),
		Link:            wpconfig.GetOption("siteurl"),
		Description:     wpconfig.GetOption("blogdescription"),
		Language:        wpconfig.GetLang(),
		UpdatePeriod:    "hourly",
		UpdateFrequency: 1,
		Generator:       wpconfig.GetOption("home"),
	}
}

func CommentsFeedCache() *cache.VarCache[[]string] {
	r, _ := cachemanager.GetVarCache[[]string]("commentsFeed")
	return r
}

func FeedCache() *cache.VarCache[[]string] {
	r, _ := cachemanager.GetVarCache[[]string]("feed")
	return r
}

// PostFeedCache query func see PostFeed
func PostFeedCache() *cache.MapCache[string, string] {
	r, _ := cachemanager.GetMapCache[string, string]("postFeed")
	return r
}

func feed(c context.Context, _ ...any) (xml []string, err error) {
	r := RecentPosts(c, 10)
	ids := slice.Map(r, func(t models.Posts) uint64 {
		return t.Id
	})
	posts, err := GetPostsByIds(c, ids)
	if err != nil {
		return
	}
	site := wpconfig.GetOption("siteurl")
	rs := templateRss
	rs.LastBuildDate = time.Now().Format(timeFormat)
	rs.Items = slice.Map(posts, func(t models.Posts) rss2.Item {
		desc := "无法提供摘要。这是一篇受保护的文章。"
		if t.PostPassword != "" {
			wpposts.PasswordProjectTitle(&t)
			wpposts.PasswdProjectContent(&t)
		} else {
			desc = plugins.Digests(t.PostContent, t.Id, 55, nil)
		}
		l := ""
		if t.CommentStatus == "open" && t.CommentCount > 0 {
			l = fmt.Sprintf("%s/p/%d#comments", site, t.Id)
		} else if t.CommentStatus == "open" && t.CommentCount == 0 {
			l = fmt.Sprintf("%s/p/%d#respond", site, t.Id)
		}
		user := GetUserById(c, t.PostAuthor)

		return rss2.Item{
			Title:         t.PostTitle,
			Creator:       user.DisplayName,
			Guid:          t.Guid,
			SlashComments: int(t.CommentCount),
			Content:       t.PostContent,
			Category:      strings.Join(t.Categories, "、"),
			CommentLink:   l,
			CommentRss:    fmt.Sprintf("%s/p/%d/feed", site, t.Id),
			Link:          fmt.Sprintf("%s/p/%d", site, t.Id),
			Description:   desc,
			PubDate:       t.PostDateGmt.Format(timeFormat),
		}
	})
	xml = []string{rs.GetXML()}
	return
}

func PostFeed(c context.Context, id string, _ ...any) (x string, err error) {
	ID := str.ToInteger[uint64](id, 0)
	maxId, err := GetMaxPostId(c)
	logs.IfError(err, "get max post id")
	if ID < 1 || ID > maxId || err != nil {
		return
	}
	post, err := GetPostById(c, ID)
	if post.Id == 0 || err != nil {
		return
	}
	limit := str.ToInteger(wpconfig.GetOption("comments_per_page"), 10)
	ids, err := PostTopLevelCommentIds(c, ID, 1, limit, 0, "desc", "latest-comment")
	if err != nil {
		return
	}
	comments, err := GetCommentDataByIds(c, ids)
	if err != nil {
		return
	}
	rs := templateRss
	site := wpconfig.GetOption("siteurl")

	rs.Title = fmt.Sprintf("《%s》的评论", post.PostTitle)
	rs.AtomLink = fmt.Sprintf("%s/p/%d/feed", site, post.Id)
	rs.Link = fmt.Sprintf("%s/p/%d", site, post.Id)
	rs.LastBuildDate = time.Now().Format(timeFormat)
	if post.PostPassword != "" {
		wpposts.PasswordProjectTitle(&post)
		wpposts.PasswdProjectContent(&post)
		if len(comments) > 0 {
			t := comments[len(comments)-1]
			u, err := GetCommentUrl(c, t.CommentId, t.CommentPostId)
			if err != nil {
				return "", err
			}
			rs.Items = []rss2.Item{
				{
					Title:       fmt.Sprintf("评价者：%s", t.CommentAuthor),
					Link:        fmt.Sprintf("%s%s", site, u),
					Creator:     t.CommentAuthor,
					PubDate:     t.CommentDateGmt.Format(timeFormat),
					Guid:        fmt.Sprintf("%s#comment-%d", post.Guid, t.CommentId),
					Description: "评论受保护：要查看请输入密码。",
					Content:     post.PostContent,
				},
			}
		}
	} else {
		rs.Items = slice.Map(comments, func(t models.Comments) rss2.Item {
			u, er := GetCommentUrl(c, t.CommentId, t.CommentPostId)
			if er != nil {
				err = errors.Join(err, er)
				return rss2.Item{}
			}
			return rss2.Item{
				Title:   fmt.Sprintf("评价者：%s", t.CommentAuthor),
				Link:    fmt.Sprintf("%s%s", site, u),
				Creator: t.CommentAuthor,
				PubDate: t.CommentDateGmt.Format(timeFormat),
				Guid:    fmt.Sprintf("%s#comment-%d", post.Guid, t.CommentId),
				Content: t.CommentContent,
			}
		})
	}

	x = rs.GetXML()
	return
}

func commentsFeed(c context.Context, _ ...any) (r []string, err error) {
	commens := RecentComments(c, 10)
	rs := templateRss
	rs.Title = fmt.Sprintf("\"%s\"的评论", wpconfig.GetOption("blogname"))
	rs.LastBuildDate = time.Now().Format(timeFormat)
	site := wpconfig.GetOption("siteurl")
	rs.AtomLink = fmt.Sprintf("%s/comments/feed", site)
	com, err := GetCommentDataByIds(c, slice.Map(commens, func(t models.Comments) uint64 {
		return t.CommentId
	}))
	if nil != err {
		return []string{}, err
	}
	rs.Items = slice.Map(com, func(t models.Comments) rss2.Item {
		post, _ := GetPostById(c, t.CommentPostId)
		desc := "评论受保护：要查看请输入密码。"
		content := t.CommentContent
		if post.PostPassword != "" {
			wpposts.PasswordProjectTitle(&post)
			wpposts.PasswdProjectContent(&post)
			content = post.PostContent
		} else {
			content = digest.StripTags(t.CommentContent, "")
		}
		u, er := GetCommentUrl(c, t.CommentId, t.CommentPostId)
		if er != nil {
			errors.Join(err, er)
		}
		u = str.Join(site, u)
		return rss2.Item{
			Title:       fmt.Sprintf("%s对《%s》的评论", t.CommentAuthor, post.PostTitle),
			Link:        u,
			Creator:     t.CommentAuthor,
			Description: desc,
			PubDate:     t.CommentDateGmt.Format(timeFormat),
			Guid:        fmt.Sprintf("%s#commment-%d", post.Guid, t.CommentId),
			Content:     content,
		}
	})
	r = []string{rs.GetXML()}
	return
}
