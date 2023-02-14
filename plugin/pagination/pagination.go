package pagination

import (
	"github.com/fthvgb1/wp-go/helper/number"
	"strings"
)

type Elements interface {
	Current(page, totalPage, totalRows int) string
	Prev(url string) string
	Next(url string) string
	Dots() string
	Middle(page int, url string) string
	Url(path, query string, page int) string
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

func NewParsePagination(totalRaw int, pageSize int, currentPage, step int, query string, path string) ParsePagination {
	return ParsePagination{TotalPage: number.CalTotalPage(totalRaw, pageSize), TotalRaw: totalRaw, PageSize: pageSize, CurrentPage: currentPage, Query: query, Path: path, Step: step}
}

func Paginate(e Elements, p ParsePagination) string {
	p.Elements = e
	return p.ToHtml()
}

func (p ParsePagination) ToHtml() (html string) {
	if p.TotalPage < 2 {
		return
	}
	s := strings.Builder{}
	if p.CurrentPage > p.TotalPage {
		p.CurrentPage = p.TotalPage
	}
	start := p.CurrentPage - p.Step
	end := p.CurrentPage + p.Step
	if start < 1 {
		start = 1
	}
	if p.CurrentPage > 1 {
		s.WriteString(p.Prev(p.Url(p.Path, p.Query, p.CurrentPage-1)))
	}
	if p.CurrentPage >= p.Step+2 {
		d := false
		if p.CurrentPage > p.Step+2 {
			d = true
		}
		s.WriteString(p.Middle(1, p.Url(p.Path, p.Query, 1)))
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
			h = p.Current(page, p.TotalPage, p.TotalRaw)
		} else {
			h = p.Middle(page, p.Url(p.Path, p.Query, page))
		}
		s.WriteString(h)

	}
	if p.TotalPage >= p.CurrentPage+p.Step+1 {
		d := false
		if p.TotalPage > p.CurrentPage+p.Step+1 {
			d = true
		}
		if d {
			s.WriteString(p.Dots())
		}
		s.WriteString(p.Middle(p.TotalPage, p.Url(p.Path, p.Query, p.TotalPage)))
	}
	if p.CurrentPage < p.TotalPage {
		s.WriteString(p.Next(p.Url(p.Path, p.Query, p.CurrentPage+1)))
	}
	html = s.String()
	return
}
