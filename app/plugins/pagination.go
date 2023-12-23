package plugins

import (
	"fmt"
	"github.com/fthvgb1/wp-go/app/wpconfig"
	"github.com/fthvgb1/wp-go/helper"
	str "github.com/fthvgb1/wp-go/helper/strings"
	"net/url"
	"regexp"
	"strings"
)

type PageEle struct {
	PrevEle    string
	NextEle    string
	DotsEle    string
	MiddleEle  string
	CurrentEle string
}

func TwentyFifteenPagination() PageEle {
	return twentyFifteen
}
func TwentyFifteenCommentPagination() CommentPageEle {
	return twentyFifteenComment
}

type CommentPageEle struct {
	PageEle
}

var twentyFifteen = PageEle{
	PrevEle: `<a class="prev page-numbers" href="%s">上一页</a>`,
	NextEle: `<a class="next page-numbers" href="%s">下一页</a>`,
	DotsEle: `<span class="page-numbers dots">…</span>`,
	MiddleEle: `<a class="page-numbers" href="%s"><span class="meta-nav screen-reader-text">页 </span>%d</a>
`,
	CurrentEle: `<span aria-current="page" class="page-numbers current">
            <span class="meta-nav screen-reader-text">页 </span>%d</span>`,
}
var twentyFifteenComment = CommentPageEle{
	PageEle{
		PrevEle: `<div class="nav-previous"><a href="%s">%s</a></div>`,
		NextEle: `<div class="nav-next"><a href="%s">%s</a></div>`,
	},
}

func (p PageEle) Current(page, totalPage, totalRow int) string {
	return fmt.Sprintf(p.CurrentEle, page)
}

func (p PageEle) Prev(url string) string {
	return fmt.Sprintf(p.PrevEle, url)
}

func (p PageEle) Next(url string) string {
	return fmt.Sprintf(p.NextEle, url)
}

func (p PageEle) Dots() string {
	return p.DotsEle
}

func (p PageEle) Step() int {
	return 0
}

func (p PageEle) Middle(page int, url string) string {
	return fmt.Sprintf(p.MiddleEle, url, page)
}

var reg = regexp.MustCompile(`(/page)/(\d+)`)
var commentReg = regexp.MustCompile(`/comment-page-(\d+)`)

func (p PageEle) Urls(u url.URL, page int, isTLS bool) string {
	var path, query = u.Path, u.RawQuery
	if !strings.Contains(path, "/page/") {
		path = fmt.Sprintf("%s%s", path, "/page/1")
	}
	if page == 1 {
		path = reg.ReplaceAllString(path, "")
	} else {
		s := fmt.Sprintf("$1/%d", page)
		path = reg.ReplaceAllString(path, s)
	}
	path = strings.Replace(path, "//", "/", -1)
	if path == "" {
		path = "/"
	}
	return str.Join(path, query)
}

func (p CommentPageEle) Urls(u url.URL, page int, isTLS bool) string {
	var path, query = u.Path, u.RawQuery
	if !strings.Contains(path, "/comment-page-") {
		path = fmt.Sprintf("%s%s", path, "/comment-page-1#comments")
	}
	path = commentReg.ReplaceAllString(path, fmt.Sprintf("/comment-page-%d#comments", page))
	path = strings.Replace(path, "//", "/", -1)
	if path == "" {
		path = "/"
	}
	return str.Join(path, query)
}

func (p CommentPageEle) Middle(page int, url string) string {
	return ""
}
func (p CommentPageEle) Dots() string {
	return ""
}
func (p CommentPageEle) Current(page, totalPage, totalRow int) string {
	return ""
}
func (p CommentPageEle) Prev(url string) string {
	return fmt.Sprintf(p.PrevEle, url, helper.Or(wpconfig.GetOption("comment_order") == "asc", "较早评论", "较新评论"))
}

func (p CommentPageEle) Step() int {
	return 0
}

func (p CommentPageEle) Next(url string) string {
	return fmt.Sprintf(p.NextEle, url, helper.Or(wpconfig.GetOption("comment_order") == "asc", "较新评论", "较早评论"))
}

type PaginationNav struct {
	Currents func(page, totalPage, totalRows int) string
	Prevs    func(url string) string
	Nexts    func(url string) string
	Dotss    func() string
	Middles  func(page int, url string) string
	Urlss    func(u url.URL, page int, isTLS bool) string
	Steps    func() int
}

func (p PaginationNav) Current(page, totalPage, totalRows int) string {
	return p.Currents(page, totalPage, totalRows)
}

func (p PaginationNav) Prev(url string) string {
	return p.Prevs(url)
}

func (p PaginationNav) Next(url string) string {
	return p.Nexts(url)
}

func (p PaginationNav) Dots() string {
	return p.Dotss()
}

func (p PaginationNav) Middle(page int, url string) string {
	return p.Middles(page, url)
}

func (p PaginationNav) Urls(u url.URL, page int, isTLS bool) string {
	return p.Urlss(u, page, isTLS)
}

func (p PaginationNav) Step() int {
	return p.Steps()
}
