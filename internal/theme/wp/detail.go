package wp

import (
	"errors"
	"fmt"
	str "github.com/fthvgb1/wp-go/helper/strings"
	"github.com/fthvgb1/wp-go/internal/pkg/cache"
	"github.com/fthvgb1/wp-go/internal/pkg/constraints"
	"github.com/fthvgb1/wp-go/internal/pkg/logs"
	"github.com/fthvgb1/wp-go/internal/pkg/models"
	"github.com/fthvgb1/wp-go/internal/plugins"
	"github.com/fthvgb1/wp-go/internal/plugins/wpposts"
	"github.com/fthvgb1/wp-go/internal/wpconfig"
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
	d.ginH["title"] = wpconfig.GetOption("blogname")
	err = d.CheckAndGetPost()
	if err != nil {
		return
	}
	d.Comment()
	d.ContextPost()
	return
}

func (d *DetailHandle) CheckAndGetPost() (err error) {
	id := str.ToInteger[uint64](d.C.Param("id"), 0)
	maxId, err := cache.GetMaxPostId(d.C)
	logs.ErrPrintln(err, "get max post id")
	if id > maxId || id <= 0 {
		d.Stats = constraints.ParamError
		err = errors.New("无效的文档id")
	}
	if err != nil {
		return
	}
	post, err := cache.GetPostById(d.C, id)
	if post.Id == 0 || err != nil || post.PostStatus != "publish" {
		d.Stats = constraints.Error404
		logs.ErrPrintln(err, "获取id失败")
		err = errors.New(str.Join("无效的文档id "))
		return
	}

	d.Post = post
	d.ginH["user"] = cache.GetUserById(d.C, post.PostAuthor)
	d.ginH["title"] = fmt.Sprintf("%s-%s", post.PostTitle, wpconfig.GetOption("blogname"))
	return
}

func (d *DetailHandle) PasswordProject() {
	if d.Post.PostPassword != "" {
		wpposts.PasswordProjectTitle(&d.Post)
		if d.password != d.Post.PostPassword {
			wpposts.PasswdProjectContent(&d.Post)
		}
		d.ginH["post"] = d.Post
	}
}
func (d *DetailHandle) Comment() {
	comments, err := cache.PostComments(d.C, d.Post.Id)
	logs.ErrPrintln(err, "get d.Post comment", d.Post.Id)
	d.ginH["comments"] = comments
	d.Comments = comments

}

func (d *DetailHandle) RenderComment() {
	if d.CommentRender == nil {
		d.CommentRender = plugins.CommentRender()
	}
	ableComment := true
	if d.Post.CommentStatus != "open" ||
		(d.Post.PostPassword != "" && d.password != d.Post.PostPassword) {
		ableComment = false
	}
	d.ginH["showComment"] = ableComment
	if len(d.Comments) > 0 && ableComment {
		dep := str.ToInteger(wpconfig.GetOption("thread_comments_depth"), 5)
		d.ginH["comments"] = plugins.FormatComments(d.C, d.CommentRender, d.Comments, dep)
	}
}

func (d *DetailHandle) ContextPost() {
	prev, next, err := cache.GetContextPost(d.C, d.Post.Id, d.Post.PostDate)
	logs.ErrPrintln(err, "get pre and next post", d.Post.Id, d.Post.PostDate)
	d.ginH["next"] = next
	d.ginH["prev"] = prev
}

func (d *DetailHandle) Render() {
	reply := ""
	if d.Post.CommentStatus == "open" && wpconfig.GetOption("thread_comments") == "1" {
		reply = `<script src='/wp-includes/js/comment-reply.min.js' id='comment-reply-js'></script>`
	}
	d.PushGroupFooterScript(10, reply)

	d.PushHandleFn(constraints.Ok, NewHandleFn(func(h *Handle) {
		d.PasswordProject()
		d.RenderComment()
		d.ginH["post"] = d.Post
	}, 10))

	d.Handle.Render()
}

func Details(h *Handle) {
	_ = h.Detail.BuildDetailData()
	PreCodeAndStats(h)
	PreTemplate(h)
}
