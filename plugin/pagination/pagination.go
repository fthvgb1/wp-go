package pagination

import (
	"github.com/fthvgb1/wp-go/helper/number"
	"net/url"
	"strings"
)

type Render interface {
	Current(page, totalPage, totalRows int) string
	Prev(url string) string
	Next(url string) string
	Dots() string
	Middle(page int, url string) string
	Urls(u url.URL, page int, isTLS bool) string
	Step() int
}

type parser struct {
	Render
	TotalPage   int
	TotalRaw    int
	PageSize    int
	CurrentPage int
	Url         url.URL
	Step        int
	IsTLS       bool
}

func Paginate(e Render, totalRaw int, pageSize int, currentPage, step int, u url.URL, isTLS bool) string {
	st := e.Step()
	if st > 0 {
		step = st
	}
	return parser{
		Render:      e,
		TotalPage:   number.DivideCeil(totalRaw, pageSize),
		TotalRaw:    totalRaw,
		PageSize:    pageSize,
		CurrentPage: currentPage,
		Url:         u,
		Step:        step,
		IsTLS:       isTLS,
	}.ToHtml()
}

func (p parser) ToHtml() (html string) {
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
		s.WriteString(p.Prev(p.Urls(p.Url, p.CurrentPage-1, p.IsTLS)))
	}
	if p.CurrentPage >= p.Step+2 {
		d := false
		if p.CurrentPage > p.Step+2 {
			d = true
		}
		s.WriteString(p.Middle(1, p.Urls(p.Url, 1, p.IsTLS)))
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
			h = p.Middle(page, p.Urls(p.Url, page, p.IsTLS))
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
		s.WriteString(p.Middle(p.TotalPage, p.Urls(p.Url, p.TotalPage, p.IsTLS)))
	}
	if p.CurrentPage < p.TotalPage {
		s.WriteString(p.Next(p.Urls(p.Url, p.CurrentPage+1, p.IsTLS)))
	}
	html = s.String()
	return
}
