package wp

import (
	"errors"
	"fmt"
	"github.com/fthvgb1/wp-go/app/pkg/cache"
	"github.com/fthvgb1/wp-go/app/pkg/constraints"
	"github.com/fthvgb1/wp-go/app/pkg/logs"
	"github.com/fthvgb1/wp-go/app/pkg/models"
	"github.com/fthvgb1/wp-go/app/plugins"
	"github.com/fthvgb1/wp-go/app/plugins/wpposts"
	"github.com/fthvgb1/wp-go/app/wpconfig"
	"github.com/fthvgb1/wp-go/cache/cachemanager"
	"github.com/fthvgb1/wp-go/helper"
	"github.com/fthvgb1/wp-go/helper/number"
	str "github.com/fthvgb1/wp-go/helper/strings"
	"github.com/fthvgb1/wp-go/plugin/pagination"
	"net/http"
	"time"
)

type DetailHandle struct {
	*Handle
	CommentRender  plugins.CommentHtml
	Comments       []uint64
	Page           int
	Limit          int
	Post           models.Posts
	CommentPageEle pagination.Render
	TotalRaw       int
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

func ShowPreComment(h *Handle) {
	v, ok := cache.NewCommentCache().Get(h.C, h.C.Request.URL.RawQuery)
	if ok {
		h.C.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
		h.C.Writer.WriteHeader(http.StatusOK)
		_, _ = h.C.Writer.Write([]byte(v))
		h.Abort()
	}
}

func (d *DetailHandle) CheckAndGetPost() (err error) {
	id := str.ToInteger[uint64](d.C.Param("id"), 0)
	maxId, err := cache.GetMaxPostId(d.C)
	logs.IfError(err, "get max post id")
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
		logs.IfError(err, "获取id失败")
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
		if d.GetPassword() != d.Post.PostPassword {
			wpposts.PasswdProjectContent(&d.Post)
		}
	}
}
func (d *DetailHandle) Comment() {
	d.ginH["totalCommentNum"] = 0
	d.ginH["totalCommentPage"] = 1
	d.ginH["commentPageNav"] = ""
	d.ginH["commentOrder"] = wpconfig.GetOption("comment_order")
	d.Page = str.ToInteger(d.C.Param("page"), 1)
	d.ginH["currentPage"] = d.Page
	d.Limit = str.ToInteger(wpconfig.GetOption("comments_per_page"), 5)
	key := fmt.Sprintf("%d-%d-%d", d.Post.Id, d.Page, d.Limit)
	data, err := cachemanager.Get[helper.PaginationData[uint64]]("PostCommentsIds", d.C, key, time.Second, d.Post.Id, d.Page, d.Limit, 0)
	if err != nil {
		d.SetErr(err)
		return
	}
	ids := data.Data
	totalCommentNum := data.TotalRaw
	d.TotalRaw = totalCommentNum
	num, err := cachemanager.Get[int]("commentNumber", d.C, d.Post.Id, time.Second)
	if err != nil {
		d.SetErr(err)
		return
	}
	d.ginH["totalCommentNum"] = num
	d.ginH["totalCommentPage"] = number.DivideCeil(totalCommentNum, d.Limit)
	d.Comments = ids
}

func (d *DetailHandle) RenderComment() {
	if d.CommentRender == nil {
		d.CommentRender = plugins.CommentRender()
	}
	ableComment := true
	if d.Post.CommentStatus != "open" ||
		(d.Post.PostPassword != "" && d.GetPassword() != d.Post.PostPassword) {
		ableComment = false
	}
	d.ginH["showComment"] = ableComment
	d.ginH["comments"] = ""
	if len(d.Comments) < 0 || !ableComment {
		return
	}
	var err error
	d.ginH["comments"], err = RenderComment(d.C, d.Page, d.CommentRender, d.Comments, 2*time.Second, d.IsHttps())
	if err != nil {
		d.SetErr(err)
		return
	}
	if d.CommentPageEle == nil {
		d.CommentPageEle = plugins.TwentyFifteenCommentPagination()
	}
	d.ginH["commentPageNav"] = pagination.Paginate(d.CommentPageEle, d.TotalRaw, d.Limit, d.Page, 1, *d.C.Request.URL, d.IsHttps())

}

func (d *DetailHandle) ContextPost() {
	prev, next, err := cache.GetContextPost(d.C, d.Post.Id, d.Post.PostDate)
	logs.IfError(err, "get pre and next post", d.Post.Id, d.Post.PostDate)
	d.ginH["next"] = next
	d.ginH["prev"] = prev
}

func DetailRender(h *Handle) {
	if h.Stats != constraints.Ok {
		return
	}
	d := h.Detail
	d.PasswordProject()
	d.RenderComment()
	d.ginH["post"] = d.Post
}

func Detail(h *Handle) {
	err := h.Detail.BuildDetailData()
	if err != nil {
		h.Detail.SetErr(err)
	}
	h.SetData("scene", h.Scene())
}

func ReplyCommentJs(h *Handle) {
	h.PushFooterScript(constraints.Detail, NewComponent("comment-reply.js", "", false, 10, func(h *Handle) string {
		reply := ""
		if h.Detail.Post.CommentStatus == "open" && wpconfig.GetOption("thread_comments") == "1" {
			reply = `<script src='/wp-includes/js/comment-reply.min.js' id='comment-reply-js'></script>`
		}
		return reply
	}))
}
