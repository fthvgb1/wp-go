package actions

import (
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"github.com/fthvgb1/wp-go/helper/slice"
	str "github.com/fthvgb1/wp-go/helper/strings"
	"github.com/fthvgb1/wp-go/internal/mail"
	"github.com/fthvgb1/wp-go/internal/pkg/cache"
	"github.com/fthvgb1/wp-go/internal/pkg/config"
	"github.com/fthvgb1/wp-go/internal/pkg/logs"
	"github.com/fthvgb1/wp-go/internal/wpconfig"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

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
			c.Writer.WriteString("评论出错，请联系管理员或稍后再度")
		}
	}()
	conf := config.GetConfig()
	if err != nil {
		logs.Error(err, "获取评论数据错误")
		return
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(data))
	i := c.PostForm("comment_post_ID")
	author := c.PostForm("author")
	m := c.PostForm("email")
	comment := c.PostForm("comment")
	c.Request.Body = io.NopCloser(bytes.NewBuffer(data))
	req, err := http.NewRequest("POST", conf.PostCommentUrl, strings.NewReader(c.Request.PostForm.Encode()))
	if err != nil {
		logs.Error(err, "创建评论请求错误")
		return
	}
	defer req.Body.Close()
	req.Header = c.Request.Header.Clone()
	home, err := url.Parse(wpconfig.GetOption("siteurl"))
	if err != nil {
		logs.Error(err, "解析评论接口错误")
		return
	}
	req.Host = home.Host
	res, err := cli.Do(req)
	if err != nil && err != http.ErrUseLastResponse {
		logs.Error(err, "请求评论接口错误")
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
		cu, er := url.Parse(conf.PostCommentUrl)
		if er != nil {
			err = er
			return
		}
		up.Host = cu.Host
		up.Scheme = "http"
		newReq, _ := http.NewRequest("GET", up.String(), nil)
		newReq.Host = home.Host
		newReq.Header.Set("Cookie", strings.Join(slice.Map(c.Request.Cookies(), func(t *http.Cookie) string {
			return fmt.Sprintf("%s=%s", t.Name, t.Value)
		}), "; "))
		ress, er := http.DefaultClient.Do(newReq)
		if er != nil {
			err = er
			return
		}
		cc := c.Copy()
		go func() {
			id := str.ToInteger[uint64](i, 0)
			if id <= 0 {
				logs.Error(errors.New("获取文档id错误"), "", i)
				return
			}
			post, err := cache.GetPostById(cc, id)
			if err != nil {
				logs.Error(err, "获取文档错误", id)
				return
			}
			su := fmt.Sprintf("%s: %s[%s]发表了评论对文档[%v]的评论", wpconfig.GetOption("siteurl"), author, m, post.PostTitle)
			err = mail.SendMail([]string{conf.Mail.User}, su, comment)
			logs.IfError(err, "发送邮件", conf.Mail.User, su, comment)
		}()

		s, er := io.ReadAll(ress.Body)
		if er != nil {
			err = er
			return
		}
		cache.NewCommentCache().Set(c, up.RawQuery, string(s))
		c.Redirect(http.StatusFound, res.Header.Get("Location"))
		return
	}
	var r io.Reader
	if res.Header.Get("Content-Encoding") == "gzip" {
		r, err = gzip.NewReader(res.Body)
		if err != nil {
			logs.Error(err, "gzip解压错误")
			return
		}
	} else {
		r = res.Body
	}
	s, err := io.ReadAll(r)
	if err != nil {
		logs.Error(err, "读取结果错误")
		return
	}
	c.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	c.Writer.WriteHeader(res.StatusCode)
	_, _ = c.Writer.Write(s)

}
