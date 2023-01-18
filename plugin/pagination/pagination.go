package pagination

import (
	"fmt"
	"github.com/fthvgb1/wp-go/helper"
	"math"
	"regexp"
	"strings"
)

type Elements interface {
	Current(page int) string
	Prev(url string) string
	Next(url string) string
	Dots() string
	Middle(page int, url string) string
}

type ParsePagination struct {
	Elements
	TotalPage   int
	TotalRaw    int
	PageSize    int
	CurrentPage int
	Query       string
	Path        string
	Step        int
}

func NewParsePagination(totalRaw int, pageSize int, currentPage int, query string, path string, step int) ParsePagination {
	allPage := int(math.Ceil(float64(totalRaw) / float64(pageSize)))
	return ParsePagination{TotalPage: allPage, TotalRaw: totalRaw, PageSize: pageSize, CurrentPage: currentPage, Query: query, Path: path, Step: step}
}

func Paginate(e Elements, p ParsePagination) string {
	p.Elements = e
	return p.ToHtml()
}

var complie = regexp.MustCompile(`(/page)/(\d+)`)

func (p ParsePagination) ToHtml() (html string) {
	if p.TotalRaw < 2 {
		return
	}
	pathx := p.Path
	if !strings.Contains(p.Path, "/page/") {
		pathx = fmt.Sprintf("%s%s", p.Path, "/page/1")
	}
	s := strings.Builder{}
	if p.CurrentPage > p.TotalPage {
		p.CurrentPage = p.TotalPage
	}
	r := complie
	start := p.CurrentPage - p.Step
	end := p.CurrentPage + p.Step
	if start < 1 {
		start = 1
	}
	if p.CurrentPage > 1 {
		pp := ""
		if p.CurrentPage >= 2 {
			pp = replacePage(r, pathx, p.CurrentPage-1)
		}
		s.WriteString(p.Prev(helper.StrJoin(pp, p.Query)))
	}
	if p.CurrentPage >= p.Step+2 {
		d := false
		if p.CurrentPage > p.Step+2 {
			d = true
		}
		e := replacePage(r, p.Path, 1)
		s.WriteString(p.Middle(1, helper.StrJoin(e, p.Query)))
		if d {
			s.WriteString(p.Dots())
		}
	}
	if p.TotalPage < end {
		end = p.TotalPage
	}

	for page := start; page <= end; page++ {
		h := ""
		if p.CurrentPage == page {
			h = p.Current(page)
		} else {
			d := replacePage(r, pathx, page)
			h = p.Middle(page, helper.StrJoin(d, p.Query))
		}
		s.WriteString(h)

	}
	if p.TotalPage >= p.CurrentPage+p.Step+1 {
		d := false
		if p.TotalPage > p.CurrentPage+p.Step+1 {
			d = true
		}
		dd := replacePage(r, pathx, p.TotalPage)
		if d {
			s.WriteString(p.Dots())
		}
		s.WriteString(p.Middle(p.TotalPage, helper.StrJoin(dd, p.Query)))
	}
	if p.CurrentPage < p.TotalPage {
		dd := replacePage(r, pathx, p.CurrentPage+1)
		s.WriteString(p.Next(helper.StrJoin(dd, p.Query)))
	}
	html = s.String()
	return
}

func replacePage(r *regexp.Regexp, path string, page int) (src string) {
	if page == 1 {
		src = r.ReplaceAllString(path, "")
	} else {
		s := fmt.Sprintf("$1/%d", page)
		src = r.ReplaceAllString(path, s)
	}
	src = strings.Replace(src, "//", "/", -1)
	if src == "" {
		src = "/"
	}
	return
}
