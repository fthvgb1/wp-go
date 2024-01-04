package wp

import (
	"context"
	"github.com/fthvgb1/wp-go/app/pkg/models"
	"github.com/fthvgb1/wp-go/app/plugins"
	"github.com/fthvgb1/wp-go/app/wpconfig"
	"github.com/fthvgb1/wp-go/cache"
	"github.com/fthvgb1/wp-go/cache/cachemanager"
	"github.com/fthvgb1/wp-go/helper/slice"
	str "github.com/fthvgb1/wp-go/helper/strings"
	"time"
)

func RenderComment(ctx context.Context, page int, render plugins.CommentHtml, ids []uint64, timeout time.Duration, isTLS bool) (string, error) {
	ca, _ := cachemanager.GetMapCache[uint64, models.Comments]("postCommentData")
	children, _ := cachemanager.GetMapCache[uint64, []uint64]("commentChildren")
	h := CommentHandle{
		maxDepth:       str.ToInteger(wpconfig.GetOption("thread_comments_depth"), 5),
		depth:          1,
		isTls:          isTLS,
		html:           render,
		order:          wpconfig.GetOption("comment_order"),
		ca:             ca,
		children:       children,
		threadComments: wpconfig.GetOption("thread_comments") == "1",
		page:           page,
	}
	return h.formatComments(ctx, ids, timeout)
}

type CommentHandle struct {
	maxDepth       int
	depth          int
	isTls          bool
	html           plugins.CommentHtml
	order          string
	page           int
	ca             *cache.MapCache[uint64, models.Comments]
	children       *cache.MapCache[uint64, []uint64]
	threadComments bool
}

func (c CommentHandle) findComments(ctx context.Context, timeout time.Duration, comments []models.Comments) ([]models.Comments, error) {
	parentIds := slice.Map(comments, func(t models.Comments) uint64 {
		return t.CommentId
	})
	children, err := c.childrenComment(ctx, parentIds, timeout)
	rr := slice.FilterAndMap(children, func(t []uint64) ([]uint64, bool) {
		return t, len(t) > 0
	})
	if len(rr) < 1 {
		slice.Sort(comments, func(i, j models.Comments) bool {
			return c.html.FloorOrder(c.order, i, j)
		})
		return comments, nil
	}
	ids := slice.Decompress(rr)
	r, err := c.ca.GetCacheBatch(ctx, ids, timeout)
	if err != nil {
		return nil, err
	}
	rrr, err := c.findComments(ctx, timeout, r)
	if err != nil {
		return nil, err
	}
	comments = append(comments, rrr...)
	return comments, nil
}

func (c CommentHandle) childrenComment(ctx context.Context, ids []uint64, timeout time.Duration) ([][]uint64, error) {
	v, err := c.children.GetCacheBatch(ctx, ids, timeout)
	if err != nil {
		return nil, err
	}

	return slice.Copy(v), nil
}

func (c CommentHandle) formatComments(ctx context.Context, ids []uint64, timeout time.Duration) (html string, err error) {
	comments, err := c.ca.GetCacheBatch(ctx, ids, timeout)
	if err != nil {
		return "", err
	}
	if c.depth > 1 && c.depth < c.maxDepth {
		comments = slice.Copy(comments)
		slice.Sort(comments, func(i, j models.Comments) bool {
			return c.html.FloorOrder(c.order, i, j)
		})
	}
	fixChildren := false
	if c.depth >= c.maxDepth {
		comments, err = c.findComments(ctx, timeout, comments)
		if err != nil {
			return "", err
		}
		fixChildren = true
	}
	s := str.NewBuilder()
	for i, comment := range comments {
		eo := "even"
		if (i+1)%2 == 0 {
			eo = "odd"
		}
		parent := ""
		fl := false
		var children []uint64
		if !fixChildren {
			children, err = c.children.GetCache(ctx, comment.CommentId, timeout)
		}

		if err != nil {
			return "", err
		}
		if c.threadComments && len(children) > 0 && c.depth < c.maxDepth+1 {
			parent = "parent"
			fl = true
		}
		s.WriteString(c.html.FormatLi(ctx, comment, c.depth, c.maxDepth, c.page, c.isTls, c.threadComments, eo, parent))
		if fl {
			c.depth++
			ss, err := c.formatComments(ctx, children, timeout)
			if err != nil {
				return "", err
			}
			s.WriteString(`<ol class="children">`, ss, `</ol>`)
			c.depth--
		}
		s.WriteString("</li><!-- #comment-## -->")
	}

	html = s.String()
	return
}
