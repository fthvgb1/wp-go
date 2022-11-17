package actions

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github/fthvgb1/wp-go/actions/common"
	"github/fthvgb1/wp-go/cache"
	"github/fthvgb1/wp-go/config"
	"github/fthvgb1/wp-go/config/wpconfig"
	"github/fthvgb1/wp-go/helper"
	"github/fthvgb1/wp-go/logs"
	"github/fthvgb1/wp-go/mail"
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
			c.Writer.WriteHeader(http.StatusConflict)
			c.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
			c.Writer.WriteString(err.Error())
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
	req, err := http.NewRequest("POST", config.Conf.Load().PostCommentUrl, strings.NewReader(c.Request.PostForm.Encode()))
	if err != nil {
		return
	}
	defer req.Body.Close()
	req.Header = c.Request.Header.Clone()
	home, err := url.Parse(wpconfig.Options.Value("siteurl"))
	if err != nil {
		return
	}
	req.Host = home.Host
	res, err := cli.Do(req)
	if err != nil && err != http.ErrUseLastResponse {
		return
	}
	if res.StatusCode == http.StatusFound {
		for _, cookie := range res.Cookies() {
			c.SetCookie(cookie.Name, cookie.Value, cookie.MaxAge, cookie.Path, cookie.Domain, cookie.Secure, cookie.HttpOnly)
		}
		u := res.Header.Get("Location")
		up, er := url.Parse(u)
		if er != nil {
			err = er
			return
		}
		cu, er := url.Parse(config.Conf.Load().PostCommentUrl)
		if er != nil {
			err = er
			return
		}
		up.Host = cu.Host
		up.Scheme = "http"
		newReq, er := http.NewRequest("GET", up.String(), nil)
		if er != nil {
			err = er
			return
		}
		newReq.Host = home.Host
		newReq.Header.Set("Cookie", strings.Join(helper.SliceMap(c.Request.Cookies(), func(t *http.Cookie) string {
			return fmt.Sprintf("%s=%s", t.Name, t.Value)
		}), "; "))
		ress, er := http.DefaultClient.Do(newReq)
		if er != nil {
			err = er
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
			su := fmt.Sprintf("%s: %s[%s]发表了评论对文档[%v]的评论", wpconfig.Options.Value("siteurl"), author, m, post.PostTitle)
			err = mail.SendMail([]string{config.Conf.Load().Mail.User}, su, comment)
			logs.ErrPrintln(err, "发送邮件", config.Conf.Load().Mail.User, su, comment)
		}()

		s, er := io.ReadAll(ress.Body)
		if er != nil {
			err = er
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
