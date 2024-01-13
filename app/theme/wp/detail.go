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
	"github.com/fthvgb1/wp-go/helper/number"
	str "github.com/fthvgb1/wp-go/helper/strings"
	"github.com/fthvgb1/wp-go/plugin/pagination"
	"strings"
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
	TotalPage      int
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
	d.CommentData()
	d.ContextPost()
	return
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
func (d *DetailHandle) CommentData() {
	d.ginH["totalCommentNum"] = 0
	d.ginH["totalCommentPage"] = 1
	d.ginH["commentPageNav"] = ""
	order := wpconfig.GetOption("comment_order")
	d.ginH["commentOrder"] = order
	d.Limit = str.ToInteger(wpconfig.GetOption("comments_per_page"), 5)
	pageComments := wpconfig.GetOption("page_comments")
	num, err := cachemanager.GetBy[int]("commentNumber", d.C, d.Post.Id, time.Second)
	if err != nil {
		d.SetErr(err)
		return
	}
	if num < 1 {
		return
	}
	topNum, err := cachemanager.GetBy[int]("postTopCommentsNum", d.C, d.Post.Id, time.Second)
	if err != nil {
		d.SetErr(err)
		return
	}
	d.TotalPage = number.DivideCeil(topNum, d.Limit)
	if !strings.Contains(d.C.Request.URL.Path, "comment-page") {
		defaultCommentsPage := wpconfig.GetOption("default_comments_page")
		if order == "desc" && defaultCommentsPage == "oldest" || order == "asc" && defaultCommentsPage == "newest" {
			d.C.AddParam("page", number.IntToString(d.TotalPage))
		}
	}
	d.Page = str.ToInteger(d.C.Param("page"), 1)
	d.ginH["currentPage"] = d.Page
	var key string
	if pageComments != "1" {
		key = number.IntToString(d.Post.Id)
		d.Limit = 0
	} else {
		key = fmt.Sprintf("%d-%d-%d", d.Post.Id, d.Page, d.Limit)
	}
	d.ginH["page_comments"] = pageComments
	d.ginH["totalCommentPage"] = d.TotalPage
	if d.TotalPage < d.Page {
		d.SetErr(errors.New("curren page above total page"))
		return
	}
	data, err := cache.PostTopLevelCommentIds(d.C, d.Post.Id, d.Page, d.Limit, topNum, order, key)
	if err != nil {
		d.SetErr(err)
		return
	}
	d.TotalRaw = topNum
	d.ginH["totalCommentNum"] = num
	d.Comments = data
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
	if wpconfig.GetOption("page_comments") == "1" && d.TotalPage > 1 {
		d.ginH["commentPageNav"] = pagination.Paginate(d.CommentPageEle, d.TotalRaw, d.Limit, d.Page, 1, *d.C.Request.URL, d.IsHttps())
	}
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
