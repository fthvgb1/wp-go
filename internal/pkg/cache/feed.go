package cache

import (
	"fmt"
	"github.com/fthvgb1/wp-go/cache"
	"github.com/fthvgb1/wp-go/helper"
	"github.com/fthvgb1/wp-go/internal/pkg/logs"
	"github.com/fthvgb1/wp-go/internal/pkg/models"
	"github.com/fthvgb1/wp-go/internal/plugins"
	"github.com/fthvgb1/wp-go/internal/wpconfig"
	"github.com/fthvgb1/wp-go/plugin/digest"
	"github.com/fthvgb1/wp-go/rss2"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
	"time"
)

var timeFormat = "Mon, 02 Jan 2006 15:04:05 +0000"
var templateRss rss2.Rss2

func InitFeed() {
	templateRss = rss2.Rss2{
		Title:           wpconfig.Options.Value("blogname"),
		AtomLink:        fmt.Sprintf("%s/feed", wpconfig.Options.Value("home")),
		Link:            wpconfig.Options.Value("siteurl"),
		Description:     wpconfig.Options.Value("blogdescription"),
		Language:        "zh-CN",
		UpdatePeriod:    "hourly",
		UpdateFrequency: 1,
		Generator:       wpconfig.Options.Value("home"),
	}
}

func CommentsFeedCache() *cache.SliceCache[string] {
	return commentsFeedCache
}

func FeedCache() *cache.SliceCache[string] {
	return feedCache
}

func PostFeedCache() *cache.MapCache[string, string] {
	return postFeedCache
}

func feed(arg ...any) (xml []string, err error) {
	c := arg[0].(*gin.Context)
	r := RecentPosts(c, 10)
	ids := helper.SliceMap(r, func(t models.Posts) uint64 {
		return t.Id
	})
	posts, err := GetPostsByIds(c, ids)
	if err != nil {
		return
	}
	rs := templateRss
	rs.LastBuildDate = time.Now().Format(timeFormat)
	rs.Items = helper.SliceMap(posts, func(t models.Posts) rss2.Item {
		desc := "无法提供摘要。这是一篇受保护的文章。"
		plugins.PasswordProjectTitle(&t)
		if t.PostPassword != "" {
			plugins.PasswdProjectContent(&t)
		} else {
			desc = digest.Raw(t.PostContent, 55, fmt.Sprintf("/p/%d", t.Id))
		}
		l := ""
		if t.CommentStatus == "open" && t.CommentCount > 0 {
			l = fmt.Sprintf("%s/p/%d#comments", wpconfig.Options.Value("siteurl"), t.Id)
		} else if t.CommentStatus == "open" && t.CommentCount == 0 {
			l = fmt.Sprintf("%s/p/%d#respond", wpconfig.Options.Value("siteurl"), t.Id)
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
			CommentRss:    fmt.Sprintf("%s/p/%d/feed", wpconfig.Options.Value("siteurl"), t.Id),
			Link:          fmt.Sprintf("%s/p/%d", wpconfig.Options.Value("siteurl"), t.Id),
			Description:   desc,
			PubDate:       t.PostDateGmt.Format(timeFormat),
		}
	})
	xml = []string{rs.GetXML()}
	return
}

func postFeed(arg ...any) (x string, err error) {
	c := arg[0].(*gin.Context)
	id := arg[1].(string)
	Id := 0
	if id != "" {
		Id, err = strconv.Atoi(id)
		if err != nil {
			return
		}
	}
	ID := uint64(Id)
	maxId, err := GetMaxPostId(c)
	logs.ErrPrintln(err, "get max post id")
	if ID > maxId || err != nil {
		return
	}
	post, err := GetPostById(c, ID)
	if post.Id == 0 || err != nil {
		return
	}
	plugins.PasswordProjectTitle(&post)
	comments, err := PostComments(c, post.Id)
	if err != nil {
		return
	}
	rs := templateRss

	rs.Title = fmt.Sprintf("《%s》的评论", post.PostTitle)
	rs.AtomLink = fmt.Sprintf("%s/p/%d/feed", wpconfig.Options.Value("siteurl"), post.Id)
	rs.Link = fmt.Sprintf("%s/p/%d", wpconfig.Options.Value("siteurl"), post.Id)
	rs.LastBuildDate = time.Now().Format(timeFormat)
	if post.PostPassword != "" {
		if len(comments) > 0 {
			plugins.PasswdProjectContent(&post)
			t := comments[len(comments)-1]
			rs.Items = []rss2.Item{
				{
					Title:       fmt.Sprintf("评价者：%s", t.CommentAuthor),
					Link:        fmt.Sprintf("%s/p/%d#comment-%d", wpconfig.Options.Value("siteurl"), post.Id, t.CommentId),
					Creator:     t.CommentAuthor,
					PubDate:     t.CommentDateGmt.Format(timeFormat),
					Guid:        fmt.Sprintf("%s#comment-%d", post.Guid, t.CommentId),
					Description: "评论受保护：要查看请输入密码。",
					Content:     post.PostContent,
				},
			}
		}
	} else {
		rs.Items = helper.SliceMap(comments, func(t models.Comments) rss2.Item {
			return rss2.Item{
				Title:   fmt.Sprintf("评价者：%s", t.CommentAuthor),
				Link:    fmt.Sprintf("%s/p/%d#comment-%d", wpconfig.Options.Value("siteurl"), post.Id, t.CommentId),
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

func commentsFeed(args ...any) (r []string, err error) {
	c := args[0].(*gin.Context)
	commens := RecentComments(c, 10)
	rs := templateRss
	rs.Title = fmt.Sprintf("\"%s\"的评论", wpconfig.Options.Value("blogname"))
	rs.LastBuildDate = time.Now().Format(timeFormat)
	rs.AtomLink = fmt.Sprintf("%s/comments/feed", wpconfig.Options.Value("siteurl"))
	com, err := GetCommentByIds(c, helper.SliceMap(commens, func(t models.Comments) uint64 {
		return t.CommentId
	}))
	if nil != err {
		return []string{}, err
	}
	rs.Items = helper.SliceMap(com, func(t models.Comments) rss2.Item {
		post, _ := GetPostById(c, t.CommentPostId)
		plugins.PasswordProjectTitle(&post)
		desc := "评论受保护：要查看请输入密码。"
		content := t.CommentContent
		if post.PostPassword != "" {
			plugins.PasswdProjectContent(&post)
			content = post.PostContent
		} else {
			desc = digest.ClearHtml(t.CommentContent)
			content = desc
		}
		return rss2.Item{
			Title:       fmt.Sprintf("%s对《%s》的评论", t.CommentAuthor, post.PostTitle),
			Link:        fmt.Sprintf("%s/p/%d#comment-%d", wpconfig.Options.Value("siteurl"), post.Id, t.CommentId),
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