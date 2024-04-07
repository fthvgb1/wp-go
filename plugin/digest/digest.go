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
	var runeIndex [][]int
	if len(index) > 0 {
		runeIndex = slice.Map(index, func(t []int) []int {
			return slice.Map(t, func(i int) int {
				return utf8.RuneCountInString(content[:i])
			})
		})
	}
	end := 0
	ru := []rune(content)
	tagIn := false
	total := len(ru)
	l, r := '<', '>'
	i := -1
	for {
		i++
		for len(runeIndex) > 0 && i >= runeIndex[0][0] {
			ints := runeIndex[0]
			if ints[0] <= i {
				i = ints[1]
				runeIndex = runeIndex[1:]
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
