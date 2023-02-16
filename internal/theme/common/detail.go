package common

import (
	"fmt"
	str "github.com/fthvgb1/wp-go/helper/strings"
	"github.com/fthvgb1/wp-go/internal/pkg/cache"
	"github.com/fthvgb1/wp-go/internal/pkg/constraints"
	"github.com/fthvgb1/wp-go/internal/pkg/logs"
	"github.com/fthvgb1/wp-go/internal/pkg/models"
	"github.com/fthvgb1/wp-go/internal/plugins"
	"github.com/fthvgb1/wp-go/internal/wpconfig"
	"net/http"
)

type DetailHandle struct {
	*Handle
	CommentRender plugins.CommentHtml
	Comments      []models.Comments
	Post          models.Posts
}

func NewDetailHandle(handle *Handle) *DetailHandle {
	return &DetailHandle{Handle: handle}
}

func (d *DetailHandle) BuildDetailData() (err error) {
	d.GinH["title"] = wpconfig.Options.Value("blogname")
	err = d.CheckAndGetPost()
	if err != nil {
		return
	}
	d.WidgetAreaData()
	d.GetPassword()
	d.Comment()
	d.ContextPost()
	return
}

func (d *DetailHandle) CheckAndGetPost() (err error) {
	id := str.ToInteger[uint64](d.C.Param("id"), 0)
	maxId, err := cache.GetMaxPostId(d.C)
	logs.ErrPrintln(err, "get max post id")
	if id > maxId || id <= 0 || err != nil {
		return
	}
	post, err := cache.GetPostById(d.C, id)
	if post.Id == 0 || err != nil || post.PostStatus != "publish" {
		return
	}

	d.GinH["post"] = post
	d.Post = post
	d.GinH["user"] = cache.GetUserById(d.C, post.PostAuthor)
	d.GinH["title"] = fmt.Sprintf("%s-%s", post.PostTitle, wpconfig.Options.Value("blogname"))
	return
}

func (d *DetailHandle) PasswordProject() {
	if d.Post.PostPassword != "" {
		plugins.PasswordProjectTitle(&d.Post)
		if d.Password != d.Post.PostPassword {
			plugins.PasswdProjectContent(&d.Post)
		}
		d.GinH["post"] = d.Post
	}
}
func (d *DetailHandle) Comment() {
	comments, err := cache.PostComments(d.C, d.Post.Id)
	logs.ErrPrintln(err, "get d.Post comment", d.Post.Id)
	d.GinH["comments"] = comments
	d.Comments = comments

}

func (d *DetailHandle) RenderComment() {
	ableComment := true
	if d.Post.CommentStatus != "open" ||
		(d.Post.PostPassword != "" && d.Password != d.Post.PostPassword) {
		ableComment = false
	}
	d.GinH["showComment"] = ableComment
	if len(d.Comments) > 0 && ableComment {
		dep := str.ToInteger(wpconfig.Options.Value("thread_comments_depth"), 5)
		d.GinH["comments"] = plugins.FormatComments(d.C, d.CommentRender, d.Comments, dep)
	}
}

func (d *DetailHandle) ContextPost() {
	prev, next, err := cache.GetContextPost(d.C, d.Post.Id, d.Post.PostDate)
	logs.ErrPrintln(err, "get pre and next post", d.Post.Id, d.Post.PostDate)
	d.GinH["next"] = next
	d.GinH["prev"] = prev
}

func (d *DetailHandle) Render() {
	d.PasswordProject()
	if d.CommentRender == nil {
		d.CommentRender = plugins.CommentRender()
	}
	d.SiteIcon()
	d.CustomLogo()
	d.CustomCss()
	d.RenderComment()
	d.CalBodyClass()
	if d.Templ == "" {
		d.Templ = fmt.Sprintf("%s/posts/detail.gohtml", d.Theme)
	}
	d.CustomBackGround()
	d.C.HTML(d.Code, d.Templ, d.GinH)
}

func (d *DetailHandle) Details() {
	err := d.BuildDetailData()
	if err != nil {
		d.Stats = constraints.Error404
		d.Code = http.StatusNotFound
		d.C.HTML(d.Code, d.Templ, d.GinH)
		return
	}
	d.Render()
}
