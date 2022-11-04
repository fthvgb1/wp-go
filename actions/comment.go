package actions

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github/fthvgb1/wp-go/actions/common"
	"github/fthvgb1/wp-go/config"
	"github/fthvgb1/wp-go/logs"
	"github/fthvgb1/wp-go/mail"
	"github/fthvgb1/wp-go/models"
	"io"
	"net/http"
	"net/http/cookiejar"
	"strconv"
	"strings"
	"time"
)

func PostComment(c *gin.Context) {
	jar, _ := cookiejar.New(nil)
	cli := &http.Client{
		Jar:     jar,
		Timeout: time.Second * 3,
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
	for k, v := range c.Request.Header {
		req.Header.Set(k, v[0])
	}
	res, err := cli.Do(req)
	if err != nil {
		return
	}
	if res.Request.Response != nil && res.Request.Response.StatusCode == http.StatusFound {
		for _, cookie := range res.Request.Response.Cookies() {
			c.SetCookie(cookie.Name, cookie.Value, cookie.MaxAge, cookie.Path, cookie.Domain, cookie.Secure, cookie.HttpOnly)
		}
		c.Redirect(http.StatusFound, res.Request.Response.Header.Get("Location"))
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
			su := fmt.Sprintf("%s: %s[%s]发表了评论对文档[%v]的评论", models.Options["siteurl"], author, m, post.PostTitle)
			err = mail.SendMail([]string{config.Conf.Mail.User}, su, comment)
			logs.ErrPrintln(err, "发送邮件")
		}()
		return
	}
	s, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}
	err = errors.New(string(s))
}
