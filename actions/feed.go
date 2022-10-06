package actions

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github/fthvgb1/wp-go/actions/common"
	"github/fthvgb1/wp-go/helper"
	"github/fthvgb1/wp-go/logs"
	"github/fthvgb1/wp-go/models"
	"github/fthvgb1/wp-go/plugins"
	"github/fthvgb1/wp-go/templates"
	"html"
	"html/template"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"
)

func Feed() func(ctx *gin.Context) {
	fs, err := template.ParseFS(templates.TemplateFs, "feed/feed.gohtml")
	if err != nil {
		panic(err)
	}
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/rss+xml; charset=UTF-8")
		c.Header("Cache-Control", "no-cache, must-revalidate, max-age=0")
		c.Header("Expires", "Wed, 11 Jan 1984 05:00:00 GMT")
		//c.Header("Last-Modified", "false")
		c.Header("ETag", helper.StringMd5("gmt"))
		r := common.RecentPosts(c, 10)
		ids := helper.SliceMap(r, func(t models.WpPosts) uint64 {
			return t.Id
		})
		posts, err := common.GetPostsByIds(c, ids)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			c.Abort()
			return
		}
		type p struct {
			models.WpPosts
			Cates       string
			CommentLink string
			Username    string
			Category    string
			Link        string
			Description string
			Date        string
		}
		rr := helper.SliceMap(posts, func(t models.WpPosts) p {
			common.PasswordProjectTitle(&t)
			if t.PostPassword != "" {
				common.PasswdProjectContent(&t)
			}
			l := ""
			if t.CommentStatus == "open" || t.CommentCount > 0 {
				l = fmt.Sprintf("%s/p/%d#comments", models.Options["siteurl"], t.Id)
			}
			user := common.GetUser(c, t.PostAuthor)
			content := plugins.DigestRaw(t.PostContent, utf8.RuneCountInString(t.PostContent), t.Id)
			t.PostContent = content
			return p{
				WpPosts:     t,
				Cates:       strings.Join(t.Categories, "、"),
				CommentLink: l,
				Username:    user.DisplayName,
				Link:        fmt.Sprintf("%s/p/%d", models.Options["siteurl"], t.Id),
				Description: plugins.DigestRaw(content, 55, t.Id),
				Date:        t.PostDateGmt.Format(time.RFC1123Z),
			}
		})
		h := gin.H{
			"posts":   rr,
			"options": models.Options,
			"now":     time.Now().Format(time.RFC1123Z),
		}

		var buf bytes.Buffer
		err = fs.Execute(&buf, h)
		if err != nil {
			logs.ErrPrintln(err, "parse template")
			return
		}

		c.String(http.StatusOK, html.UnescapeString(buf.String()))
	}

}
