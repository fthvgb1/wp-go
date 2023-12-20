package plugins

import (
	"context"
	"github.com/fthvgb1/wp-go/app/pkg/models"
	"github.com/fthvgb1/wp-go/app/wpconfig"
	"github.com/fthvgb1/wp-go/helper/number"
	"github.com/fthvgb1/wp-go/helper/slice"
	str "github.com/fthvgb1/wp-go/helper/strings"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

type CommentHandler struct {
	*gin.Context
	comments         []*Comments
	maxDepth         int
	depth            int
	isTls            bool
	i                CommentHtml
	isThreadComments bool
}

type Comments struct {
	models.Comments
	Children []*Comments
}

type CommentHtml interface {
	FormatLi(c context.Context, m models.Comments, depth, maxDepth, page int, isTls, isThreadComments bool, eo, parent string) string
	FloorOrder(wpOrder string, i, j models.PostComments) bool
}

func FormatComments(c *gin.Context, i CommentHtml, comments []models.Comments, maxDepth int) string {
	tree := treeComments(comments)

	var isTls bool
	if c.Request.TLS != nil {
		isTls = true
	} else {
		isTls = "https" == strings.ToLower(c.Request.Header.Get("X-Forwarded-Proto"))
	}
	h := CommentHandler{
		Context:  c,
		comments: tree,
		maxDepth: maxDepth,
		depth:    1,
		isTls:    isTls,
		i:        i,
	}
	return h.formatComment(h.comments)
}

func (d CommentHandler) formatComment(comments []*Comments) (html string) {
	s := str.NewBuilder()
	if d.depth >= d.maxDepth {
		comments = d.findComments(comments)
	}
	order := wpconfig.GetOption("comment_order")
	slice.Sort(comments, func(i, j *Comments) bool {
		if order == "asc" {
			return i.CommentDate.Sub(j.CommentDate) < 0
		}
		return i.CommentDate.Sub(j.CommentDate) > 0
	})
	for i, comment := range comments {
		eo := "even"
		if (i+1)%2 == 0 {
			eo = "odd"
		}
		parent := ""
		fl := false
		if len(comment.Children) > 0 && d.depth < d.maxDepth+1 {
			parent = "parent"
			fl = true
		}
		s.WriteString(d.i.FormatLi(d.Context, comment.Comments, d.depth, d.maxDepth, 1, d.isTls, d.isThreadComments, eo, parent))
		if fl {
			d.depth++
			s.WriteString(`<ol class="children">`, d.formatComment(comment.Children), `</ol>`)
			d.depth--
		}
		s.WriteString("</li><!-- #comment-## -->")
	}

	html = s.String()
	return
}

func (d CommentHandler) findComments(comments []*Comments) []*Comments {
	var r []*Comments
	for _, comment := range comments {
		tmp := *comment
		tmp.Children = nil
		r = append(r, &tmp)
		if len(comment.Children) > 0 {
			t := d.findComments(comment.Children)
			r = append(r, t...)
		}
	}
	return r
}

func treeComments(comments []models.Comments) []*Comments {
	var r = map[uint64]*Comments{
		0: {
			Comments: models.Comments{},
		},
	}
	var top []*Comments
	for _, comment := range comments {
		c := Comments{
			Comments: comment,
			Children: []*Comments{},
		}
		r[comment.CommentId] = &c
		if comment.CommentParent == 0 {
			top = append(top, &c)
		}
	}
	for id, son := range r {
		if id == son.CommentParent {
			continue
		}
		parent := r[son.CommentParent]
		parent.Children = append(parent.Children, son)
	}
	return top
}

func CommonLi() string {
	return li
}

var commonCommentFormat = CommonCommentFormat{}

func CommentRender() CommonCommentFormat {
	return commonCommentFormat
}

type CommonCommentFormat struct {
}

func (c CommonCommentFormat) FormatLi(_ context.Context, m models.Comments, currentDepth, maxDepth, page int, isTls, isThreadComments bool, eo, parent string) string {
	return FormatLi(CommonLi(), m, currentDepth, maxDepth, page, isTls, isThreadComments, eo, parent)
}

func (c CommonCommentFormat) FloorOrder(wpOrder string, i, j models.PostComments) bool {
	return i.CommentId > j.CommentId
}
func respond(m models.Comments, isShow bool) string {
	if !isShow {
		return ""
	}
	pId := number.IntToString(m.CommentPostId)
	cId := number.IntToString(m.CommentId)
	return str.Replace(respondTml, map[string]string{
		"{{PostId}}":        pId,
		"{{CommentId}}":     cId,
		"{{CommentAuthor}}": m.CommentAuthor,
	})
}

var li = `
<li id="comment-{{CommentId}}" class="comment {{eo}} thread-even depth-{{Depth}} {{parent}}">
    <article id="div-comment-{{CommentId}}" class="comment-body">
        <footer class="comment-meta">
            <div class="comment-author vcard">
                <img alt=""
                     src="{{Gravatar}}"
                     srcset="{{Gravatar}} 2x"
                     class="avatar avatar-56 photo" height="56" width="56" loading="lazy">
                <b class="fn">{{CommentAuthor}}</b>
                <span class="says">说道：</span></div><!-- .comment-author -->

            <div class="comment-metadata">
                <a href="/p/{{PostId}}/comment-page-{{page}}#comment-{{CommentId}}">
                    <time datetime="{{CommentDateGmt}}">{{CommentDate}}</time>
                </a></div><!-- .comment-metadata -->

        </footer><!-- .comment-meta -->

        <div class="comment-content">
            <p>{{CommentContent}}</p>
        </div><!-- .comment-content -->

        {{respond}}
    </article><!-- .comment-body -->

`

var respondTml = `<div class="reply">
            <a rel="nofollow" class="comment-reply-link"
               href="/p/{{PostId}}?replytocom={{CommentId}}#respond" data-commentid="{{CommentId}}" data-postid="{{PostId}}"
               data-belowelement="div-comment-{{CommentId}}" data-respondelement="respond"
               data-replyto="回复给{{CommentAuthor}}"
               aria-label="回复给{{CommentAuthor}}">回复</a>
        </div>`

func FormatLi(li string, comments models.Comments, currentDepth, maxDepth, page int, isTls, isThreadComments bool, eo, parent string) string {
	isShow := false
	if isThreadComments && currentDepth < maxDepth {
		isShow = true
	}
	for k, v := range map[string]string{
		"{{CommentId}}":        strconv.FormatUint(comments.CommentId, 10),
		"{{Depth}}":            strconv.Itoa(currentDepth),
		"{{Gravatar}}":         Gravatar(comments.CommentAuthorEmail, isTls),
		"{{CommentAuthorUrl}}": comments.CommentAuthorUrl,
		"{{CommentAuthor}}":    comments.CommentAuthor,
		"{{PostId}}":           strconv.FormatUint(comments.CommentPostId, 10),
		"{{page}}":             strconv.Itoa(page),
		"{{CommentDateGmt}}":   comments.CommentDateGmt.String(),
		"{{CommentDate}}":      comments.CommentDate.Format("2006-01-02 15:04"),
		"{{CommentContent}}":   comments.CommentContent,
		"{{eo}}":               eo,
		"{{parent}}":           parent,
		"{{respond}}":          respond(comments, isShow),
	} {
		li = strings.Replace(li, k, v, -1)
	}
	return li
}
