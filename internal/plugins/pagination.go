package plugins

import "fmt"

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
