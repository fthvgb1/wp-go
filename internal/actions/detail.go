package actions

import (
	"fmt"
	str "github.com/fthvgb1/wp-go/helper/strings"
	"github.com/fthvgb1/wp-go/internal/pkg/cache"
	"github.com/fthvgb1/wp-go/internal/pkg/logs"
	"github.com/fthvgb1/wp-go/internal/pkg/models"
	"github.com/fthvgb1/wp-go/internal/plugins"
	"github.com/fthvgb1/wp-go/internal/theme"
	"github.com/fthvgb1/wp-go/internal/theme/common"
	"github.com/fthvgb1/wp-go/internal/wpconfig"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Detail(c *gin.Context) {
	var err error
	var post models.Posts
	recent := cache.RecentPosts(c, 5, true)
	archive := cache.Archives(c)
	categoryItems := cache.CategoriesTags(c, plugins.Category)
	recentComments := cache.RecentComments(c, 5)
	var ginH = gin.H{
		"title":          wpconfig.Options.Value("blogname"),
		"recentPosts":    recent,
		"archives":       archive,
		"categories":     categoryItems,
		"recentComments": recentComments,
		"post":           post,
	}
	isApproveComment := false
	status := plugins.Ok
	pw := sessions.Default(c).Get("post_password")

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
		theme.Hook(t, common.Handle{
			C:        c,
			GinH:     ginH,
			Password: "",
			Scene:    plugins.Detail,
			Code:     code,
			Stats:    status,
		})
	}()
	ID := str.ToInteger[uint64](c.Param("id"), 0)

	maxId, err := cache.GetMaxPostId(c)
	logs.ErrPrintln(err, "get max post id")
	if ID > maxId || ID <= 0 || err != nil {
		return
	}
	post, err = cache.GetPostById(c, ID)
	if post.Id == 0 || err != nil || post.PostStatus != "publish" {
		return
	}
	showComment := false
	if post.CommentCount > 0 || post.CommentStatus == "open" {
		showComment = true
	}
	user := cache.GetUserById(c, post.PostAuthor)

	if post.PostPassword != "" {
		plugins.PasswordProjectTitle(&post)
		if pw != post.PostPassword {
			plugins.PasswdProjectContent(&post)
			showComment = false
		}
	} else if s, ok := cache.NewCommentCache().Get(c, c.Request.URL.RawQuery); ok && s != "" && (post.PostPassword == "" || post.PostPassword != "" && pw == post.PostPassword) {
		c.Writer.WriteHeader(http.StatusOK)
		c.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, err = c.Writer.WriteString(s)
		isApproveComment = true
		return
	}
	comments, err := cache.PostComments(c, post.Id)
	logs.ErrPrintln(err, "get post comment", post.Id)
	ginH["comments"] = comments
	prev, next, err := cache.GetContextPost(c, post.Id, post.PostDate)
	logs.ErrPrintln(err, "get pre and next post", post.Id, post.PostDate)
	ginH["title"] = fmt.Sprintf("%s-%s", post.PostTitle, wpconfig.Options.Value("blogname"))
	ginH["post"] = post
	ginH["showComment"] = showComment
	ginH["prev"] = prev
	d := str.ToInteger(wpconfig.Options.Value("thread_comments_depth"), 5)
	ginH["maxDep"] = d
	ginH["next"] = next
	ginH["user"] = user
	ginH["scene"] = plugins.Detail
}
