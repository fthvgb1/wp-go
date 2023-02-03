package actions

import (
	str "github.com/fthvgb1/wp-go/helper/strings"
	"github.com/fthvgb1/wp-go/internal/pkg/cache"
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
	if !isCacheExpired(c, cache.FeedCache().GetLastSetTime()) {
		c.Status(http.StatusNotModified)
	} else {
		r, err := cache.FeedCache().GetCache(c, time.Second, c)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			c.Abort()
			c.Error(err)
			return
		}
		setFeed(r[0], c, cache.FeedCache().GetLastSetTime())
	}
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
	if !isCacheExpired(c, cache.PostFeedCache().GetLastSetTime(c, id)) {
		c.Status(http.StatusNotModified)
	} else {
		s, err := cache.PostFeedCache().GetCache(c, id, time.Second, c, id)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			c.Abort()
			c.Error(err)
			return
		}
		setFeed(s, c, cache.PostFeedCache().GetLastSetTime(c, id))
	}
}

func CommentsFeed(c *gin.Context) {
	if !isCacheExpired(c, cache.CommentsFeedCache().GetLastSetTime()) {
		c.Status(http.StatusNotModified)
	} else {
		r, err := cache.CommentsFeedCache().GetCache(c, time.Second, c)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			c.Abort()
			c.Error(err)
			return
		}
		setFeed(r[0], c, cache.CommentsFeedCache().GetLastSetTime())
	}
}
