package actions

import (
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"github.com/fthvgb1/wp-go/app/mail"
	"github.com/fthvgb1/wp-go/app/pkg/cache"
	"github.com/fthvgb1/wp-go/app/pkg/config"
	"github.com/fthvgb1/wp-go/app/pkg/logs"
	"github.com/fthvgb1/wp-go/app/wpconfig"
	"github.com/fthvgb1/wp-go/helper/slice"
	str "github.com/fthvgb1/wp-go/helper/strings"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type CommentForm struct {
	CommentPostId uint64 `form:"comment_post_ID" binding:"required" json:"comment_post_ID"`
	Author        string `form:"author" binding:"required" label:"显示名称" json:"author"`
	Email         string `form:"email" binding:"required,email"`
	Comment       string `form:"comment" binding:"required" label:"评论" json:"comment"`
}

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
			var v validator.ValidationErrors
			if errors.As(err, &v) {
				e := v.Translate(config.GetZh())
				for _, v := range e {
					fmt.Fprintf(c.Writer, fail, v)
					return
				}
			} else {
				c.Writer.WriteString("评论出错，请联系管理员或稍后再度")
			}

		}
	}()
	conf := config.GetConfig()
	if err != nil {
		logs.Error(err, "获取评论数据错误")
		return
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(data))
	var comment CommentForm
	if err = c.ShouldBind(&comment); err != nil {
		return
	}
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
	if err != nil && !errors.Is(err, http.ErrUseLastResponse) {
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
			if gin.Mode() != gin.ReleaseMode {
				return
			}
			id := comment.CommentPostId
			if id <= 0 {
				logs.Error(errors.New("获取文档id错误"), "", comment.CommentPostId)
				return
			}
			post, err := cache.GetPostById(cc, id)
			if err != nil {
				logs.Error(err, "获取文档错误", id)
				return
			}
			su := fmt.Sprintf("%s: %s[%s]发表了评论对文档[%v]的评论", wpconfig.GetOption("siteurl"), comment.Author, comment.Email, post.PostTitle)
			err = mail.SendMail([]string{conf.Mail.User}, su, comment.Comment)
			logs.IfError(err, "发送邮件", conf.Mail.User, su, comment)
		}()

		s, er := io.ReadAll(ress.Body)
		if er != nil {
			err = er
			return
		}
		cache.NewCommentCache().Set(c, up.RawQuery, string(s))
		uu, _ := url.Parse(res.Header.Get("Location"))
		uuu := str.Join(uu.Path, "?", uu.RawQuery)
		c.Redirect(http.StatusFound, uuu)
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

var fail = `

<!DOCTYPE html>
<html lang="zh-CN">
<head>
	<meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
	<meta name="viewport" content="width=device-width">
		<meta name='robots' content='max-image-preview:large, noindex, follow' />
	<title>评论提交失败</title>
	<style type="text/css">
		html {
			background: #f1f1f1;
		}
		body {
			background: #fff;
			border: 1px solid #ccd0d4;
			color: #444;
			font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Oxygen-Sans, Ubuntu, Cantarell, "Helvetica Neue", sans-serif;
			margin: 2em auto;
			padding: 1em 2em;
			max-width: 700px;
			-webkit-box-shadow: 0 1px 1px rgba(0, 0, 0, .04);
			box-shadow: 0 1px 1px rgba(0, 0, 0, .04);
		}
		h1 {
			border-bottom: 1px solid #dadada;
			clear: both;
			color: #666;
			font-size: 24px;
			margin: 30px 0 0 0;
			padding: 0;
			padding-bottom: 7px;
		}
		#error-page {
			margin-top: 50px;
		}
		#error-page p,
		#error-page .wp-die-message {
			font-size: 14px;
			line-height: 1.5;
			margin: 25px 0 20px;
		}
		#error-page code {
			font-family: Consolas, Monaco, monospace;
		}
		ul li {
			margin-bottom: 10px;
			font-size: 14px ;
		}
		a {
			color: #0073aa;
		}
		a:hover,
		a:active {
			color: #006799;
		}
		a:focus {
			color: #124964;
			-webkit-box-shadow:
				0 0 0 1px #5b9dd9,
				0 0 2px 1px rgba(30, 140, 190, 0.8);
			box-shadow:
				0 0 0 1px #5b9dd9,
				0 0 2px 1px rgba(30, 140, 190, 0.8);
			outline: none;
		}
		.button {
			background: #f3f5f6;
			border: 1px solid #016087;
			color: #016087;
			display: inline-block;
			text-decoration: none;
			font-size: 13px;
			line-height: 2;
			height: 28px;
			margin: 0;
			padding: 0 10px 1px;
			cursor: pointer;
			-webkit-border-radius: 3px;
			-webkit-appearance: none;
			border-radius: 3px;
			white-space: nowrap;
			-webkit-box-sizing: border-box;
			-moz-box-sizing:    border-box;
			box-sizing:         border-box;

			vertical-align: top;
		}

		.button.button-large {
			line-height: 2.30769231;
			min-height: 32px;
			padding: 0 12px;
		}

		.button:hover,
		.button:focus {
			background: #f1f1f1;
		}

		.button:focus {
			background: #f3f5f6;
			border-color: #007cba;
			-webkit-box-shadow: 0 0 0 1px #007cba;
			box-shadow: 0 0 0 1px #007cba;
			color: #016087;
			outline: 2px solid transparent;
			outline-offset: 0;
		}

		.button:active {
			background: #f3f5f6;
			border-color: #7e8993;
			-webkit-box-shadow: none;
			box-shadow: none;
		}

			</style>
</head>
<body id="error-page">
	<div class="wp-die-message"><p><strong>错误：</strong>%s</p></div>
<p><a href='javascript:history.back()'>&laquo; 返回</a></p></body>
</html>
	
`
