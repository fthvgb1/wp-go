package actions

import (
	"fmt"
	"github.com/fthvgb1/wp-go/internal/pkg/cache"
	"github.com/fthvgb1/wp-go/internal/pkg/logs"
	"github.com/fthvgb1/wp-go/internal/pkg/models"
	"github.com/fthvgb1/wp-go/internal/plugins"
	"github.com/fthvgb1/wp-go/internal/theme"
	"github.com/fthvgb1/wp-go/internal/wpconfig"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

type detailHandler struct {
	*gin.Context
}

func Detail(c *gin.Context) {
	var err error
	hh := detailHandler{
		c,
	}
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
	} else if s, ok := cache.NewCommentCache().Get(c.Request.URL.RawQuery); ok && s != "" && (post.PostPassword == "" || post.PostPassword != "" && pw == post.PostPassword) {
		c.Writer.WriteHeader(http.StatusOK)
		c.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, err = c.Writer.WriteString(s)
		isApproveComment = true
		return
	}
	plugins.ApplyPlugin(plugins.NewPostPlugin(c, plugins.Detail), &post)
	comments, err := cache.PostComments(c, post.Id)
	logs.ErrPrintln(err, "get post comment", post.Id)
	commentss := treeComments(comments)
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
	ginH["comments"] = hh.formatComment(commentss, 1, d)
	ginH["next"] = next
	ginH["user"] = user
}

type Comment struct {
	models.Comments
	Children []*Comment
}

type Comments []*Comment

func (c Comments) Len() int {
	return len(c)
}

func (c Comments) Less(i, j int) bool {
	return c[i].CommentId < c[j].CommentId
}

func (c Comments) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (d detailHandler) formatComment(comments Comments, depth, maxDepth int) (html string) {
	s := strings.Builder{}
	if depth > maxDepth {
		comments = findComments(comments)
	}
	sort.Sort(comments)
	for i, comment := range comments {
		eo := "even"
		if (i+1)%2 == 0 {
			eo = "odd"
		}
		p := ""
		fl := false
		if len(comment.Children) > 0 && depth < maxDepth+1 {
			p = "parent"
			fl = true
		}
		s.WriteString(d.formatLi(comment.Comments, depth, eo, p))
		if fl {
			depth++
			s.WriteString(`<ol class="children">`)
			s.WriteString(d.formatComment(comment.Children, depth, maxDepth))
			s.WriteString(`</ol>`)
		}
		s.WriteString("</li><!-- #comment-## -->")
		i++
		if i >= len(comments) {
			break
		}
	}

	html = s.String()
	return
}

func findComments(comments Comments) Comments {
	var r Comments
	for _, comment := range comments {
		tmp := *comment
		tmp.Children = nil
		r = append(r, &tmp)
		if len(comment.Children) > 0 {
			t := findComments(comment.Children)
			r = append(r, t...)
		}
	}
	return r
}

func treeComments(comments []models.Comments) Comments {
	var r = map[uint64]*Comment{
		0: {
			Comments: models.Comments{},
		},
	}
	var top []*Comment
	for _, comment := range comments {
		c := Comment{
			Comments: comment,
			Children: []*Comment{},
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

func (d detailHandler) formatLi(comments models.Comments, depth int, eo, parent string) string {
	li := `
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
	for k, v := range map[string]string{
		"{{CommentId}}":        strconv.FormatUint(comments.CommentId, 10),
		"{{Depth}}":            strconv.Itoa(depth),
		"{{Gravatar}}":         plugins.Gravatar(comments.CommentAuthorEmail, d.Context.Request.TLS != nil),
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
