package actions

import (
	"fmt"
	"github.com/fthvgb1/wp-go/internal/pkg/cache"
	"github.com/fthvgb1/wp-go/internal/pkg/logs"
	"github.com/fthvgb1/wp-go/internal/plugins"
	"github.com/fthvgb1/wp-go/internal/theme"
	"github.com/fthvgb1/wp-go/internal/wpconfig"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type detailHandler struct {
	*gin.Context
}

func Detail(c *gin.Context) {
	var err error
	recent := cache.RecentPosts(c, 5)
	archive := cache.Archives(c)
	categoryItems := cache.Categories(c)
	recentComments := cache.RecentComments(c, 5)
	var ginH = gin.H{
		"title":          wpconfig.Options.Value("blogname"),
		"options":        wpconfig.Options,
		"recentPosts":    recent,
		"archives":       archive,
		"categories":     categoryItems,
		"recentComments": recentComments,
	}
	isApproveComment := false
	status := plugins.Ok
	defer func() {
		code := http.StatusOK
		if err != nil {
			code = http.StatusNotFound
			c.Error(err)
			status = plugins.Error
			return
		}
		if isApproveComment == true {
			return
		}

		t := theme.GetTemplateName()
		theme.Hook(t, code, c, ginH, plugins.Detail, status)
	}()
	id := c.Param("id")
	Id := 0
	if id != "" {
		Id, err = strconv.Atoi(id)
		if err != nil {
			return
		}
	}
	ID := uint64(Id)
	maxId, err := cache.GetMaxPostId(c)
	logs.ErrPrintln(err, "get max post id")
	if ID > maxId || err != nil {
		return
	}
	post, err := cache.GetPostById(c, ID)
	if post.Id == 0 || err != nil || post.PostStatus != "publish" {
		return
	}
	pw := sessions.Default(c).Get("post_password")
	showComment := false
	if post.CommentCount > 0 || post.CommentStatus == "open" {
		showComment = true
	}
	user := cache.GetUserById(c, post.PostAuthor)
	plugins.PasswordProjectTitle(&post)
	if post.PostPassword != "" && pw != post.PostPassword {
		plugins.PasswdProjectContent(&post)
		showComment = false
	} else if s, ok := cache.NewCommentCache().Get(c, c.Request.URL.RawQuery); ok && s != "" && (post.PostPassword == "" || post.PostPassword != "" && pw == post.PostPassword) {
		c.Writer.WriteHeader(http.StatusOK)
		c.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, err = c.Writer.WriteString(s)
		isApproveComment = true
		return
	}
	plugins.ApplyPlugin(plugins.NewPostPlugin(c, plugins.Detail), &post)
	comments, err := cache.PostComments(c, post.Id)
	logs.ErrPrintln(err, "get post comment", post.Id)
	ginH["comments"] = comments
	prev, next, err := cache.GetContextPost(c, post.Id, post.PostDate)
	logs.ErrPrintln(err, "get pre and next post", post.Id, post.PostDate)
	ginH["title"] = fmt.Sprintf("%s-%s", post.PostTitle, wpconfig.Options.Value("blogname"))
	ginH["post"] = post
	ginH["showComment"] = showComment
	ginH["prev"] = prev
	depth := wpconfig.Options.Value("thread_comments_depth")
	d, err := strconv.Atoi(depth)
	if err != nil {
		logs.ErrPrintln(err, "get comment depth ", depth)
		d = 5
	}
	ginH["maxDep"] = d
	ginH["next"] = next
	ginH["user"] = user
	ginH["scene"] = plugins.Detail
}
