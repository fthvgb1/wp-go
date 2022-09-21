package actions

import (
	"context"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github/fthvgb1/wp-go/actions/common"
	"github/fthvgb1/wp-go/models"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

func Detail(c *gin.Context) {
	var err error
	ctx := context.TODO()
	recent := common.RecentPosts(ctx)
	archive := common.Archives()
	categoryItems := common.Categories(ctx)
	recentComments := common.RecentComments(ctx)
	var h = gin.H{
		"title":          models.Options["blogname"],
		"options":        models.Options,
		"recentPosts":    recent,
		"archives":       archive,
		"categories":     categoryItems,
		"recentComments": recentComments,
	}
	defer func() {
		c.HTML(http.StatusOK, "posts/detail.gohtml", h)
		if err != nil {
			c.Error(err)
		}
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
	post := common.GetPostFromCache(ID)
	if post.Id == 0 {
		er := common.QueryAndSetPostCache([]models.WpPosts{{Id: ID}})
		if er != nil {
			err = er
			return
		}
		post = common.GetPostFromCache(ID)
		if post.Id == 0 {
			return
		}
	}
	pw := sessions.Default(c).Get("post_password")
	showComment := false
	if post.CommentCount > 0 || post.CommentStatus == "open" {
		showComment = true
	}
	common.PasswordProjectTitle(&post)
	if post.PostPassword != "" && pw != post.PostPassword {
		common.PasswdProjectContent(&post)
		showComment = false
	}
	comments, err := common.PostComments(ctx, post.Id)
	commentss := treeComments(comments)
	prev, next, err := common.GetContextPost(post.Id, post.PostDate)
	h["title"] = fmt.Sprintf("%s-%s", post.PostTitle, models.Options["blogname"])
	h["post"] = post
	h["showComment"] = showComment
	h["prev"] = prev
	h["comments"] = formatComment(commentss, 1, 5)
	h["next"] = next
}

type Comment struct {
	models.WpComments
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

func formatComment(comments Comments, depth, maxDepth int) (html string) {
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
		s.WriteString(formatLi(comment.WpComments, depth, eo, p))
		if fl {
			depth++
			s.WriteString(`<ol class="children">`)
			s.WriteString(formatComment(comment.Children, depth, maxDepth))
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

func treeComments(comments []models.WpComments) Comments {
	var r = map[uint64]*Comment{
		0: {
			WpComments: models.WpComments{},
		},
	}
	var top []*Comment
	for _, comment := range comments {
		c := Comment{
			WpComments: comment,
			Children:   []*Comment{},
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

func formatLi(comments models.WpComments, d int, eo, parent string) string {
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
		"{{Depth}}":            strconv.Itoa(d),
		"{{Gravatar}}":         "http://1.gravatar.com/avatar/d7a973c7dab26985da5f961be7b74480?s=56&d=mm&r=g",
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
