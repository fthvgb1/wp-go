package digest

import (
	"github.com/fthvgb1/wp-go/helper/html"
	"github.com/fthvgb1/wp-go/helper/slice"
	"regexp"
	"strings"
	"unicode/utf8"
)

var quto = regexp.MustCompile(`&quot; *|&amp; *|&lt; *|&gt; ?|&nbsp; *`)

func StripTags(content, allowTag string) string {
	content = strings.Trim(content, " \t\n\r\000\x0B")
	content = strings.Replace(content, "]]>", "]]&gt;", -1)
	content = html.StripTagsX(content, allowTag)
	return content
}

func Html(content string, limit int) (string, string) {
	closeTag := ""
	length := utf8.RuneCountInString(content) + 1
	if length <= limit {
		return content, ""
	}
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
			ints := slice.Map(index[0], func(t int) int {
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

		if end >= limit || i >= total {
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
	if i > total {
		i = total
	}
	content = string(ru[:i])
	closeTag = html.CloseTag(content)
	return content, closeTag
}
