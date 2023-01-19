package plugins

import (
	"fmt"
	"github.com/fthvgb1/wp-go/helper"
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

var twentyFifteen = PageEle{
	PrevEle: `<a class="prev page-numbers" href="%s">上一页</a>`,
	NextEle: `<a class="next page-numbers" href="%s">下一页</a>`,
	DotsEle: `<span class="page-numbers dots">…</span>`,
	MiddleEle: `<a class="page-numbers" href="%s"><span class="meta-nav screen-reader-text">页 </span>%d</a>
`,
	CurrentEle: `<span aria-current="page" class="page-numbers current">
            <span class="meta-nav screen-reader-text">页 </span>%d</span>`,
}

func (p PageEle) Current(page int) string {
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

func (p PageEle) Middle(page int, url string) string {
	return fmt.Sprintf(p.MiddleEle, url, page)
}

var reg = regexp.MustCompile(`(/page)/(\d+)`)

func (p PageEle) Url(path, query string, page int) string {
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
	return helper.StrJoin(path, query)
}
