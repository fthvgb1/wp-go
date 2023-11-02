package actions

import (
	"github.com/fthvgb1/wp-go/app/pkg/cache"
	str "github.com/fthvgb1/wp-go/helper/strings"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

var tmp = "Mon, 02 Jan 2006 15:04:05 GMT"

func isCacheExpired(c *gin.Context, lastTime time.Time) bool {
	eTag := str.Md5(lastTime.Format(tmp))
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
	feed := cache.FeedCache()
	if !isCacheExpired(c, feed.GetLastSetTime(c)) {
		c.Status(http.StatusNotModified)
		return
	}

	r, err := feed.GetCache(c, time.Second, c)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		c.Abort()
		c.Error(err)
		return
	}
	setFeed(r[0], c, feed.GetLastSetTime(c))
}

func setFeed(s string, c *gin.Context, t time.Time) {
	lastTimeGMT := t.Format(tmp)
	eTag := str.Md5(lastTimeGMT)
	c.Header("Content-Type", "application/rss+xml; charset=UTF-8")
	c.Header("Last-Modified", lastTimeGMT)
	c.Header("ETag", eTag)
	c.String(http.StatusOK, s)
}

func PostFeed(c *gin.Context) {
	id := c.Param("id")
	postFeed := cache.PostFeedCache()
	if !isCacheExpired(c, postFeed.GetLastSetTime(c, id)) {
		c.Status(http.StatusNotModified)
		return
	}
	s, err := postFeed.GetCache(c, id, time.Second, c, id)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		c.Abort()
		c.Error(err)
		return
	}
	setFeed(s, c, postFeed.GetLastSetTime(c, id))
}

func CommentsFeed(c *gin.Context) {
	feed := cache.CommentsFeedCache()
	if !isCacheExpired(c, feed.GetLastSetTime(c)) {
		c.Status(http.StatusNotModified)
		return
	}
	r, err := feed.GetCache(c, time.Second, c)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		c.Abort()
		c.Error(err)
		return
	}
	setFeed(r[0], c, feed.GetLastSetTime(c))
}
