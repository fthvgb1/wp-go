package html

import (
	"fmt"
	"github.com/dlclark/regexp2"
	"github.com/fthvgb1/wp-go/helper/slice"
	"regexp"
	"strings"
)

var entitlesMap = map[int][]string{
	EntCompat:   {"&amp;", "&quot;", "&lt;", "&gt;"},
	EntQuotes:   {"&amp;", "&quot;", "&#039;", "&lt;", "&gt;"},
	EntNoQuotes: {"&amp;", "&lt;", "&gt;"},
	EntSpace:    {"&nbsp;"},
}
var unEntitlesMap = map[int][]string{
	EntCompat:   {"&", "\"", "<", ">"},
	EntQuotes:   {"&", "\"", "'", "<", ">"},
	EntNoQuotes: {"&", "<", ">"},
	EntSpace:    {" "},
}

const (
	EntCompat   = 1
	EntQuotes   = 2
	EntNoQuotes = 4
	EntSpace    = 8
)

func htmlSpecialChars(text string, flags int) string {
	r, ok := unEntitlesMap[flags]
	e := entitlesMap[flags]
	if !ok {
		r = unEntitlesMap[EntCompat]
		e = entitlesMap[EntCompat]
	}
	if flags&EntSpace == EntSpace {
		r = append(r, unEntitlesMap[EntSpace]...)
		e = append(e, entitlesMap[EntSpace]...)
	}

	for i, entitle := range r {
		text = strings.Replace(text, entitle, e[i], -1)
	}
	return text
}
func htmlSpecialCharsDecode(text string, flags int) string {
	r, ok := entitlesMap[flags]
	u := unEntitlesMap[flags]
	if !ok {
		r = entitlesMap[EntCompat]
		u = unEntitlesMap[EntCompat]
	}
	if flags&EntSpace == EntSpace {
		r = append(r, entitlesMap[EntSpace]...)
		u = append(u, unEntitlesMap[EntSpace]...)
	}

	for i, entitle := range r {
		text = strings.Replace(text, entitle, u[i], -1)
	}
	return text
}

var allHtmlTag = regexp.MustCompile("</?.*>")

func StripTags(str, allowable string) string {
	html := ""
	if allowable == "" {
		return allHtmlTag.ReplaceAllString(str, "")
	}
	r := strings.Split(allowable, ">")
	re := ""
	for _, reg := range r {
		if reg == "" {
			continue
		}
		tag := strings.TrimLeft(reg, "<")
		ree := fmt.Sprintf(`%s|\/%s`, tag, tag)
		re = fmt.Sprintf("%s|%s", re, ree)
	}
	ree := strings.Trim(re, "|")
	reg := fmt.Sprintf("<(?!%s).*?>", ree)
	compile, err := regexp2.Compile(reg, regexp2.IgnoreCase)
	if err != nil {
		return str
	}
	html, err = compile.Replace(str, "", 0, -1)
	if err != nil {
		return str
	}
	return html
}

var tag = regexp.MustCompile(`<(.*?)>`)

func StripTagsX(str, allowable string) string {
	if allowable == "" {
		return allHtmlTag.ReplaceAllString(str, "")
	}
	tags := tag.ReplaceAllString(allowable, "$1|")
	or := strings.TrimRight(tags, "|")
	reg := fmt.Sprintf(`<(/?(%s).*?)>`, or)
	regx := fmt.Sprintf(`\{\[(/?(%s).*?)\]\}`, or)
	cp, err := regexp.Compile(reg)
	if err != nil {
		return str
	}
	rep := cp.ReplaceAllString(str, "{[$1]}")
	tmp := tag.ReplaceAllString(rep, "")
	rex, err := regexp.Compile(regx)
	if err != nil {
		return str
	}
	html := rex.ReplaceAllString(tmp, "<$1>")
	return html
}

var tagx = regexp.MustCompile(`(</?[a-z0-9]+?)( |>)`)
var selfCloseTags = map[string]string{"area": "", "base": "", "basefont": "", "br": "", "col": "", "command": "", "embed": "", "frame": "", "hr": "", "img": "", "input": "", "isindex": "", "link": "", "meta": "", "param": "", "source": "", "track": "", "wbr": ""}

func CloseTag(str string) string {
	tags := tag.FindAllString(str, -1)
	if len(tags) < 1 {
		return str
	}
	var tagss = make([]string, 0, len(tags))
	for _, s := range tags {
		ss := strings.TrimSpace(tagx.FindString(s))
		if ss[len(ss)-1] != '>' {
			ss = fmt.Sprintf("%s>", ss)
			if _, ok := selfCloseTags[ss]; ok {
				continue
			}
		}
		tagss = append(tagss, ss)
	}
	r := slice.Map(slice.Reverse(UnClosedTag(tagss)), func(s string) string {
		return fmt.Sprintf("</%s>", strings.Trim(s, "<>"))
	})
	return strings.Join(r, "")
}

func UnClosedTag(s []string) []string {
	i := 0
	for {
		if len(s[i:]) < 2 {
			return s
		}
		l := s[i]
		r := fmt.Sprintf(`</%s>`, strings.Trim(l, "<>"))
		if s[i+1] == r {
			if len(s[i+1:]) > 1 {
				ss := s[:i]
				s = append(ss, s[i+2:]...)

			} else {
				s = s[:i]
			}
			i = 0
			continue
		}
		i++
	}
}
