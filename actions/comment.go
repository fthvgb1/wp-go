package actions

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github/fthvgb1/wp-go/actions/common"
	"github/fthvgb1/wp-go/cache"
	"github/fthvgb1/wp-go/config"
	"github/fthvgb1/wp-go/logs"
	"github/fthvgb1/wp-go/mail"
	"github/fthvgb1/wp-go/models/wp"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var commentCache = cache.NewMapCacheByFn[string, string](nil, 15*time.Minute)

func PostComment(c *gin.Context) {
	cli := &http.Client{
		Timeout: time.Second * 3,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	data, err := c.GetRawData()
	defer func() {
		if err != nil {
			c.String(http.StatusConflict, err.Error())
		}
	}()
	if err != nil {
		return
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(data))
	i := c.PostForm("comment_post_ID")
	author := c.PostForm("author")
	m := c.PostForm("email")
	comment := c.PostForm("comment")
	c.Request.Body = io.NopCloser(bytes.NewBuffer(data))
	req, err := http.NewRequest("POST", config.Conf.PostCommentUrl, strings.NewReader(c.Request.PostForm.Encode()))
	if err != nil {
		return
	}
	defer req.Body.Close()
	req.Header = c.Request.Header.Clone()
	res, err := cli.Do(req)
	if err != nil && err != http.ErrUseLastResponse {
		return
	}
	if res.StatusCode == http.StatusFound {
		for _, cookie := range res.Cookies() {
			c.SetCookie(cookie.Name, cookie.Value, cookie.MaxAge, cookie.Path, cookie.Domain, cookie.Secure, cookie.HttpOnly)
		}
		u := res.Header.Get("Location")
		up, err := url.Parse(u)
		if err != nil {
			return
		}
		cu, err := url.Parse(config.Conf.PostCommentUrl)
		if err != nil {
			return
		}
		up.Host = cu.Host

		ress, err := http.DefaultClient.Get(up.String())

		if err != nil {
			return
		}
		cc := c.Copy()
		go func() {
			id, err := strconv.ParseUint(i, 10, 64)
			if err != nil {
				logs.ErrPrintln(err, "获取文档id", i)
				return
			}
			post, err := common.GetPostById(cc, id)
			if err != nil {
				logs.ErrPrintln(err, "获取文档", id)
				return
			}
			su := fmt.Sprintf("%s: %s[%s]发表了评论对文档[%v]的评论", wp.Option["siteurl"], author, m, post.PostTitle)
			err = mail.SendMail([]string{config.Conf.Mail.User}, su, comment)
			logs.ErrPrintln(err, "发送邮件", config.Conf.Mail.User, su, comment)
		}()

		s, err := io.ReadAll(ress.Body)
		if err != nil {
			return
		}
		commentCache.Set(up.RawQuery, string(s))
		c.Redirect(http.StatusFound, res.Header.Get("Location"))
		return
	}
	s, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}
	err = errors.New(string(s))
}
