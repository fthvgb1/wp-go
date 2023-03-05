package plugins

import (
	"github.com/fthvgb1/wp-go/helper/slice"
	str "github.com/fthvgb1/wp-go/helper/strings"
	"github.com/fthvgb1/wp-go/internal/pkg/models"
	"github.com/fthvgb1/wp-go/internal/wpconfig"
	"github.com/gin-gonic/gin"
	"net/url"
	"strconv"
	"strings"
)

type CommentHandler struct {
	*gin.Context
	comments []*Comments
	maxDepth int
	depth    int
	isTls    bool
	i        CommentHtml
}

type Comments struct {
	models.Comments
	Children []*Comments
}

type CommentHtml interface {
	Sort(i, j *Comments) bool
	FormatLi(c *gin.Context, m models.Comments, depth int, isTls bool, eo, parent string) string
}

func FormatComments(c *gin.Context, i CommentHtml, comments []models.Comments, maxDepth int) string {
	tree := treeComments(comments)
	u := c.Request.Header.Get("Referer")
	var isTls bool
	if u != "" {
		uu, _ := url.Parse(u)
		if uu.Scheme == "https" {
			isTls = true
		}
	}
	h := CommentHandler{
		Context:  c,
		comments: tree,
		maxDepth: maxDepth,
		depth:    1,
		isTls:    isTls,
		i:        i,
	}
	return h.formatComment(h.comments, true)
}

func (d CommentHandler) formatComment(comments []*Comments, isTop bool) (html string) {
	s := str.NewBuilder()
	if d.depth > d.maxDepth {
		comments = d.findComments(comments)
	}
	slice.Sort(comments, d.i.Sort)
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
		s.WriteString(d.i.FormatLi(d.Context, comment.Comments, d.depth, d.isTls, eo, parent))
		if fl {
			d.depth++
			s.WriteString(`<ol class="children">`, d.formatComment(comment.Children, false), `</ol>`)
			if isTop {
				d.depth = 1
			}
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
		comment.Children = nil
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

func (c CommonCommentFormat) Sort(i, j *Comments) bool {
	order := wpconfig.GetOption("comment_order")
	if order == "asc" {
		return i.CommentDate.UnixNano() < j.CommentDate.UnixNano()
	}
	return i.CommentDate.UnixNano() > j.CommentDate.UnixNano()
}

func (c CommonCommentFormat) FormatLi(ctx *gin.Context, m models.Comments, depth int, isTls bool, eo, parent string) string {
	return FormatLi(CommonLi(), ctx, m, depth, isTls, eo, parent)
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
                <b class="fn">
                    <a href="{{CommentAuthorUrl}}" rel="external nofollow ugc"
                       class="url">{{CommentAuthor}}</a>
                </b>
                <span class="says">说道：</span></div><!-- .comment-author -->

            <div class="comment-metadata">
                <a href="/p/{{PostId}}#comment-{{CommentId}}">
                    <time datetime="{{CommentDateGmt}}">{{CommentDate}}</time>
                </a></div><!-- .comment-metadata -->

        </footer><!-- .comment-meta -->

        <div class="comment-content">
            <p>{{CommentContent}}</p>
        </div><!-- .comment-content -->

        <div class="reply">
            <a rel="nofollow" class="comment-reply-link"
               href="/p/{{PostId}}?replytocom={{CommentId}}#respond" data-commentid="{{CommentId}}" data-postid="{{PostId}}"
               data-belowelement="div-comment-{{CommentId}}" data-respondelement="respond"
               data-replyto="回复给{{CommentAuthor}}"
               aria-label="回复给{{CommentAuthor}}">回复</a>
        </div>
    </article><!-- .comment-body -->

`

func FormatLi(li string, c *gin.Context, comments models.Comments, depth int, isTls bool, eo, parent string) string {
	for k, v := range map[string]string{
		"{{CommentId}}":        strconv.FormatUint(comments.CommentId, 10),
		"{{Depth}}":            strconv.Itoa(depth),
		"{{Gravatar}}":         Gravatar(comments.CommentAuthorEmail, isTls),
		"{{CommentAuthorUrl}}": comments.CommentAuthorUrl,
		"{{CommentAuthor}}":    comments.CommentAuthor,
		"{{PostId}}":           strconv.FormatUint(comments.CommentPostId, 10),
		"{{CommentDateGmt}}":   comments.CommentDateGmt.String(),
		"{{CommentDate}}":      comments.CommentDate.Format("2006-01-02 15:04"),
		"{{CommentContent}}":   comments.CommentContent,
		"{{eo}}":               eo,
		"{{parent}}":           parent,
	} {
		li = strings.Replace(li, k, v, -1)
	}
	return li
}
