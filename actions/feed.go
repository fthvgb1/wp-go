package actions

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github/fthvgb1/wp-go/actions/common"
	"github/fthvgb1/wp-go/cache"
	"github/fthvgb1/wp-go/helper"
	"github/fthvgb1/wp-go/models"
	"github/fthvgb1/wp-go/plugins"
	"github/fthvgb1/wp-go/rss2"
	"net/http"
	"strings"
	"time"
)

var feedCache = cache.NewSliceCache[string](feed, time.Hour)
var tmp = "Mon, 02 Jan 2006 15:04:05 GMT"

func FeedCached(c *gin.Context) {
	if !isCacheExpired(c, feedCache.SetTime()) {
		c.Status(http.StatusNotModified)
		c.Abort()
		return
	}
	c.Next()
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
	s, err := feedCache.GetCache(c, time.Second, c)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		c.Abort()
		c.Error(err)
		return
	}
	lastTimeGMT := feedCache.SetTime().Format(tmp)
	eTag := helper.StringMd5(lastTimeGMT)
	c.Header("Content-Type", "application/rss+xml; charset=UTF-8")
	c.Header("Last-Modified", lastTimeGMT)
	c.Header("ETag", eTag)
	c.String(http.StatusOK, s[0])
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
	rs := rss2.Rss2{
		Title:           models.Options["blogname"],
		AtomLink:        fmt.Sprintf("%s/feed", models.Options["home"]),
		Link:            models.Options["siteurl"],
		Description:     models.Options["blogdescription"],
		LastBuildDate:   time.Now().Format(time.RFC1123Z),
		Language:        "zh-CN",
		UpdatePeriod:    "hourly",
		UpdateFrequency: 1,
		Generator:       models.Options["home"],
		Items:           nil,
	}

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
