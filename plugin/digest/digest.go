package digest

import (
	"fmt"
	"github.com/fthvgb1/wp-go/helper"
	"regexp"
	"strings"
	"unicode/utf8"
)

var removeWpBlock = regexp.MustCompile("<!-- /?wp:.*-->")
var more = regexp.MustCompile("<!--more(.*?)?-->")

var quto = regexp.MustCompile(`&quot; *|&amp; *|&lt; *|&gt; ?|&nbsp; *`)

func ClearHtml(str string) string {
	content := removeWpBlock.ReplaceAllString(str, "")
	content = strings.Trim(content, " \t\n\r\000\x0B")
	content = strings.Replace(content, "]]>", "]]&gt;", -1)
	content = helper.StripTagsX(content, "<a><b><blockquote><br><cite><code><dd><del><div><dl><dt><em><h1><h2><h3><h4><h5><h6><i><img><li><ol><p><pre><span><strong><ul>")
	return str
}

func Raw(str string, limit int, u string) string {
	if r := more.FindString(str); r != "" {
		m := strings.Split(str, r)
		str = m[0]
		return ""
	}
	content := removeWpBlock.ReplaceAllString(str, "")
	content = strings.Trim(content, " \t\n\r\000\x0B")
	content = strings.Replace(content, "]]>", "]]&gt;", -1)
	content = helper.StripTagsX(content, "<a><b><blockquote><br><cite><code><dd><del><div><dl><dt><em><h1><h2><h3><h4><h5><h6><i><img><li><ol><p><pre><span><strong><ul>")
	length := utf8.RuneCountInString(content) + 1
	if length > limit {
		index := quto.FindAllStringIndex(content, -1)
		end := 0
		ru := []rune(content)
		tagIn := false
		total := len(ru)
		l, r := '<', '>'
		i := -1
		for {
			i++
			for len(index) > 0 {
				ints := helper.SliceMap(index[0], func(t int) int {
					return utf8.RuneCountInString(content[:t])
				})
				if ints[0] <= i {
					i = i + i - ints[0] + ints[1] - ints[0]
					index = index[1:]
					end++
					continue
				} else {
					break
				}
			}

			if end >= limit || i >= total-1 {
				break
			}
			if ru[i] == l {
				tagIn = true
				continue
			} else if ru[i] == r {
				tagIn = false
				continue
			}
			if tagIn == false {
				end++
			}
		}
		if i > total-1 {
			i = total - 1
		}

		content = string(ru[:i])
		closeTag := helper.CloseHtmlTag(content)
		tmp := `%s......%s<p class="read-more"><a href="%s">继续阅读</a></p>`
		if strings.Contains(closeTag, "pre") || strings.Contains(closeTag, "code") {
			tmp = `%s%s......<p class="read-more"><a href="%s">继续阅读</a></p>`
		}
		content = fmt.Sprintf(tmp, content, closeTag, u)
	}
	return content
}
