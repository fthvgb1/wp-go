package actions

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github/fthvgb1/wp-go/actions/common"
	"github/fthvgb1/wp-go/cache"
	"github/fthvgb1/wp-go/helper"
	"github/fthvgb1/wp-go/logs"
	"github/fthvgb1/wp-go/models"
	"github/fthvgb1/wp-go/plugins"
	"github/fthvgb1/wp-go/rss2"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var feedCache = cache.NewSliceCache[string](feed, time.Hour)
var postFeedCache = cache.NewMapCacheByFn[string, string](postFeed, time.Hour)
var tmp = "Mon, 02 Jan 2006 15:04:05 GMT"
var templateRss rss2.Rss2

func InitFeed() {
	templateRss = rss2.Rss2{
		Title:           models.Options["blogname"],
		AtomLink:        fmt.Sprintf("%s/feed", models.Options["home"]),
		Link:            models.Options["siteurl"],
		Description:     models.Options["blogdescription"],
		Language:        "zh-CN",
		UpdatePeriod:    "hourly",
		UpdateFrequency: 1,
		Generator:       models.Options["home"],
	}
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
	if !isCacheExpired(c, feedCache.SetTime()) {
		c.Status(http.StatusNotModified)
	} else {
		setFeed(feedCache, c)
	}
}

func feed(arg ...any) (xml []string, err error) {
	c := arg[0].(*gin.Context)
	r := common.RecentPosts(c, 10)
	ids := helper.SliceMap(r, func(t models.WpPosts) uint64 {
		return t.Id
	})
	posts, err := common.GetPostsByIds(c, ids)
	if err != nil {
		return
	}
	rs := templateRss
	rs.LastBuildDate = time.Now().Format(time.RFC1123Z)
	rs.Items = helper.SliceMap(posts, func(t models.WpPosts) rss2.Item {
		desc := "无法提供摘要。这是一篇受保护的文章。"
		common.PasswordProjectTitle(&t)
		if t.PostPassword != "" {
			common.PasswdProjectContent(&t)
		} else {
			desc = plugins.DigestRaw(t.PostContent, 55, t.Id)
		}
		l := ""
		if t.CommentStatus == "open" && t.CommentCount > 0 {
			l = fmt.Sprintf("%s/p/%d#comments", models.Options["siteurl"], t.Id)
		} else if t.CommentStatus == "open" && t.CommentCount == 0 {
			l = fmt.Sprintf("%s/p/%d#respond", models.Options["siteurl"], t.Id)
		}
		user := common.GetUser(c, t.PostAuthor)

		return rss2.Item{
			Title:         t.PostTitle,
			Creator:       user.DisplayName,
			Guid:          t.Guid,
			SlashComments: int(t.CommentCount),
			Content:       t.PostContent,
			Category:      strings.Join(t.Categories, "、"),
			CommentLink:   l,
			CommentRss:    fmt.Sprintf("%s/p/%d/feed", models.Options["siteurl"], t.Id),
			Link:          fmt.Sprintf("%s/p/%d", models.Options["siteurl"], t.Id),
			Description:   desc,
			PubDate:       t.PostDateGmt.Format(time.RFC1123Z),
		}
	})
	xml = []string{rs.GetXML()}
	return
}

func setFeed(sliceCache *cache.SliceCache[string], c *gin.Context) {
	s, err := sliceCache.GetCache(c, time.Second, c)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		c.Abort()
		c.Error(err)
		return
	}
	lastTimeGMT := sliceCache.SetTime().Format(tmp)
	eTag := helper.StringMd5(lastTimeGMT)
	c.Header("Content-Type", "application/rss+xml; charset=UTF-8")
	c.Header("Last-Modified", lastTimeGMT)
	c.Header("ETag", eTag)
	c.String(http.StatusOK, s[0])
}

func PostFeed(c *gin.Context) {
	id := c.Param("id")
	if !isCacheExpired(c, postFeedCache.GetSetTime(id)) {
		c.Status(http.StatusNotModified)
	} else {
		s, err := postFeedCache.GetCache(c, id, time.Second, c, id)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			c.Abort()
			c.Error(err)
			return
		}
		lastTimeGMT := postFeedCache.GetSetTime(id).Format(tmp)
		eTag := helper.StringMd5(lastTimeGMT)
		c.Header("Content-Type", "application/rss+xml; charset=UTF-8")
		c.Header("Last-Modified", lastTimeGMT)
		c.Header("ETag", eTag)
		c.String(http.StatusOK, s)
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
	maxId, err := common.GetMaxPostId(c)
	logs.ErrPrintln(err, "get max post id")
	if ID > maxId || err != nil {
		return
	}
	post, err := common.GetPostAndCache(c, ID)
	if post.Id == 0 || err != nil {
		return
	}
	common.PasswordProjectTitle(&post)
	comments, err := common.PostComments(c, post.Id)
	if err != nil {
		return
	}
	rs := templateRss

	rs.Title = fmt.Sprintf("《%s》的评论", post.PostTitle)
	rs.AtomLink = fmt.Sprintf("%s/p/%d/feed", models.Options["siteurl"], post.Id)
	rs.Link = fmt.Sprintf("%s/p/%d", models.Options["siteurl"], post.Id)
	rs.LastBuildDate = time.Now().Format(time.RFC1123Z)
	if post.PostPassword != "" {
		if len(comments) > 0 {
			common.PasswdProjectContent(&post)
			t := comments[len(comments)-1]
			rs.Items = []rss2.Item{
				{
					Title:       fmt.Sprintf("评价者：%s", t.CommentAuthor),
					Link:        fmt.Sprintf("%s/p/%d#comment-%d", models.Options["siteurl"], post.Id, t.CommentId),
					Creator:     t.CommentAuthor,
					PubDate:     t.CommentDateGmt.Format(time.RFC1123Z),
					Guid:        fmt.Sprintf("%s#comment-%d", post.Guid, t.CommentId),
					Description: "评论受保护：要查看请输入密码。",
					Content:     post.PostContent,
				},
			}
		}
	} else {
		rs.Items = helper.SliceMap(comments, func(t models.WpComments) rss2.Item {
			return rss2.Item{
				Title:   fmt.Sprintf("评价者：%s", t.CommentAuthor),
				Link:    fmt.Sprintf("%s/p/%d#comment-%d", models.Options["siteurl"], post.Id, t.CommentId),
				Creator: t.CommentAuthor,
				PubDate: t.CommentDateGmt.Format(time.RFC1123Z),
				Guid:    fmt.Sprintf("%s#comment-%d", post.Guid, t.CommentId),
				Content: t.CommentContent,
			}
		})
	}

	x = rs.GetXML()
	return
}
