package html

import (
	"bytes"
	"fmt"
	"github.com/dlclark/regexp2"
	"github.com/fthvgb1/wp-go/helper/slice"
	strings2 "github.com/fthvgb1/wp-go/helper/strings"
	"html/template"
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

func SpecialChars(text string, flag ...int) string {
	if len(text) < 1 {
		return ""
	}
	flags := EntQuotes
	if len(flag) > 0 {
		flags = flag[0]
	}
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
func SpecialCharsDecode(text string, flag ...int) string {
	flags := EntQuotes
	if len(flag) > 0 {
		flags = flag[0]
	}
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

var tag = regexp.MustCompile(`(?is:<(.*?)>)`) //使用忽略大小写和包含换行符模式

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

var selfCloseTags = map[string]string{"area": "", "base": "", "basefont": "", "br": "", "col": "", "command": "", "fecolormatrix": "", "embed": "", "frame": "", "hr": "", "img": "", "input": "", "isindex": "", "link": "", "fecomposite": "", "fefuncr": "", "fefuncg": "", "fefuncb": "", "fefunca": "", "meta": "", "param": "", "!doctype": "", "source": "", "track": "", "wbr": ""}

func GetSelfCloseTags() map[string]string {
	return selfCloseTags
}
func SetSelfCloseTags(m map[string]string) {
	selfCloseTags = m
}

func CloseTag(str string) string {
	tags := tag.FindAllString(str, -1)
	if len(tags) < 1 {
		return str
	}
	var tagss = make([]string, 0, len(tags))
	for _, s := range tags {
		ss := strings.Split(s, " ")
		sss := strings2.Replace(ss[0], map[string]string{
			"\\":  "",
			"\n":  "",
			"\\/": "/",
		})
		if strings.Contains(sss, "<!") {
			continue
		}
		if sss[len(sss)-1] != '>' {
			sss = fmt.Sprintf("%s>", sss)
		}
		if _, ok := selfCloseTags[strings.Trim(strings.ToLower(sss), "\\/<>")]; ok {
			continue
		}
		tagss = append(tagss, sss)
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

func RenderedHtml(t *template.Template, data map[string]any) (r string, err error) {
	var buf bytes.Buffer
	err = t.Execute(&buf, data)
	if err != nil {
		return
	}
	r = buf.String()
	return
}

func BuildOptions[T any, K comparable](a []T, selected K, fn func(T) (K, any, string)) string {
	s := strings2.NewBuilder()
	for _, t := range a {
		k, v, attr := fn(t)
		ss := ""
		if k == selected {
			ss = "selected"
		}
		s.Sprintf(`<option %s %s value="%v">%v</option>`, ss, attr, v, k)
	}
	return s.String()
}
