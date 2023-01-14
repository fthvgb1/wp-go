package actions

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github/fthvgb1/wp-go/cache"
	"github/fthvgb1/wp-go/helper"
	cache3 "github/fthvgb1/wp-go/internal/pkg/cache"
	"github/fthvgb1/wp-go/internal/pkg/logs"
	models2 "github/fthvgb1/wp-go/internal/pkg/models"
	"github/fthvgb1/wp-go/internal/plugins"
	"github/fthvgb1/wp-go/internal/wpconfig"
	"github/fthvgb1/wp-go/plugin/digest"
	"github/fthvgb1/wp-go/rss2"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var feedCache = cache.NewSliceCache(feed, time.Hour)
var postFeedCache = cache.NewMapCacheByFn[string, string](postFeed, time.Hour)
var tmp = "Mon, 02 Jan 2006 15:04:05 GMT"
var timeFormat = "Mon, 02 Jan 2006 15:04:05 +0000"
var templateRss rss2.Rss2
var commentsFeedCache = cache.NewSliceCache(commentsFeed, time.Hour)

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

func ClearCache() {
	postFeedCache.ClearExpired()
	commentCache.ClearExpired()
}
func FlushCache() {
	postFeedCache.Flush()
	commentCache.Flush()
}

func isCacheExpired(c *gin.Context, lastTime time.Time) bool {
	eTag := helper.StringMd5(lastTime.Format(tmp))
	since := c.Request.Header.Get("If-Modified-Since")
	cTag := c.Request.Header.Get("If-None-Match")
	if since != "" && cTag != "" {
		cGMT, err := time.Parse(tmp, since)
		if err == nil && lastTime.Unix() <= cGMT.Unix() && eTag == cTag {
			c.Status(http.StatusNotModified)
			return false
		}
	}
	return true
}

func Feed(c *gin.Context) {
	if !isCacheExpired(c, feedCache.GetLastSetTime()) {
		c.Status(http.StatusNotModified)
	} else {
		r, err := feedCache.GetCache(c, time.Second, c)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			c.Abort()
			c.Error(err)
			return
		}
		setFeed(r[0], c, feedCache.GetLastSetTime())
	}
}

func feed(arg ...any) (xml []string, err error) {
	c := arg[0].(*gin.Context)
	r := cache3.RecentPosts(c, 10)
	ids := helper.SliceMap(r, func(t models2.Posts) uint64 {
		return t.Id
	})
	posts, err := cache3.GetPostsByIds(c, ids)
	if err != nil {
		return
	}
	rs := templateRss
	rs.LastBuildDate = time.Now().Format(timeFormat)
	rs.Items = helper.SliceMap(posts, func(t models2.Posts) rss2.Item {
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
		user := cache3.GetUserById(c, t.PostAuthor)

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

func setFeed(s string, c *gin.Context, t time.Time) {
	lastTimeGMT := t.Format(tmp)
	eTag := helper.StringMd5(lastTimeGMT)
	c.Header("Content-Type", "application/rss+xml; charset=UTF-8")
	c.Header("Last-Modified", lastTimeGMT)
	c.Header("ETag", eTag)
	c.String(http.StatusOK, s)
}

func PostFeed(c *gin.Context) {
	id := c.Param("id")
	if !isCacheExpired(c, postFeedCache.GetLastSetTime(id)) {
		c.Status(http.StatusNotModified)
	} else {
		s, err := postFeedCache.GetCache(c, id, time.Second, c, id)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			c.Abort()
			c.Error(err)
			return
		}
		setFeed(s, c, postFeedCache.GetLastSetTime(id))
	}
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
	maxId, err := cache3.GetMaxPostId(c)
	logs.ErrPrintln(err, "get max post id")
	if ID > maxId || err != nil {
		return
	}
	post, err := cache3.GetPostById(c, ID)
	if post.Id == 0 || err != nil {
		return
	}
	plugins.PasswordProjectTitle(&post)
	comments, err := cache3.PostComments(c, post.Id)
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
		rs.Items = helper.SliceMap(comments, func(t models2.Comments) rss2.Item {
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

func CommentsFeed(c *gin.Context) {
	if !isCacheExpired(c, commentsFeedCache.GetLastSetTime()) {
		c.Status(http.StatusNotModified)
	} else {
		r, err := commentsFeedCache.GetCache(c, time.Second, c)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			c.Abort()
			c.Error(err)
			return
		}
		setFeed(r[0], c, commentsFeedCache.GetLastSetTime())
	}
}

func commentsFeed(args ...any) (r []string, err error) {
	c := args[0].(*gin.Context)
	commens := cache3.RecentComments(c, 10)
	rs := templateRss
	rs.Title = fmt.Sprintf("\"%s\"的评论", wpconfig.Options.Value("blogname"))
	rs.LastBuildDate = time.Now().Format(timeFormat)
	rs.AtomLink = fmt.Sprintf("%s/comments/feed", wpconfig.Options.Value("siteurl"))
	com, err := cache3.GetCommentByIds(c, helper.SliceMap(commens, func(t models2.Comments) uint64 {
		return t.CommentId
	}))
	if nil != err {
		return []string{}, err
	}
	rs.Items = helper.SliceMap(com, func(t models2.Comments) rss2.Item {
		post, _ := cache3.GetPostById(c, t.CommentPostId)
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
