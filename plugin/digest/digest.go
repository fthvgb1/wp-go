package digest

import (
	"github.com/fthvgb1/wp-go/helper/html"
	"github.com/fthvgb1/wp-go/helper/slice"
	str "github.com/fthvgb1/wp-go/helper/strings"
	"regexp"
	"strings"
	"unicode/utf8"
)

var quto = regexp.MustCompile(`&quot;*|&amp;*|&lt;*|&gt;*|&nbsp;*|&#91;*|&#93;*|&emsp;*`)

func SetQutos(reg string) {
	quto = regexp.MustCompile(reg)
}

type SpecialSolveConf struct {
	Num             int
	ChuckOvered     bool
	EscapeCharacter map[rune]SpecialSolve
	Tags            map[string]SpecialSolve
}
type SpecialSolve struct {
	Num         int
	ChuckOvered bool
}

var selfCloseTags = html.GetSelfCloseTags()

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
		if end >= limit || i >= total {
			break
		}
		for len(runeIndex) > 0 && i == runeIndex[0][0] {
			i = runeIndex[0][1]
			runeIndex = runeIndex[1:]
			end++
			continue
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
		if tagIn == false && ru[i] != '\n' {
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

func CustomizeHtml(content string, limit int, m map[string]SpecialSolveConf) (string, string) {
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
	count := 0
	runeContent := []rune(content)
	tagIn := false
	runeTotal := len(runeContent)
	l, r := '<', '>'
	i := -1
	var currentTag, parentTag string
	var allTags = []string{"<top>"}
	var tag []rune
	var tagLocal = 0
	for {
		i++
		if count >= limit || i >= runeTotal {
			break
		}
		for len(runeIndex) > 0 && i == runeIndex[0][0] {
			i = runeIndex[0][1]
			runeIndex = runeIndex[1:]
			count++
			continue
		}

		if count >= limit || i >= runeTotal {
			break
		}

		if runeContent[i] == l {
			tagLocal = i
			tagIn = true
			continue
		}
		if tagIn && runeContent[i] == r {
			tagIn = false
			tags := str.Join("<", string(tag), ">")
			if strings.Contains(tags, " ") {
				tags = str.Join("<", strings.Split(string(tag), " ")[0], ">")
			}
			currentTag = tags
			rawTag := strings.ReplaceAll(strings.Trim(tags, "<>"), "/", "")
			_, ok := selfCloseTags[rawTag]
			if !ok {
				if '/' == tags[1] {
					parentTag = allTags[len(allTags)-2]
					allTags = allTags[:len(allTags)-1]
				} else {
					parentTag = allTags[len(allTags)-1]
					allTags = append(allTags, currentTag)
				}
			} else {
				parentTag = allTags[len(allTags)-1]
			}
			tag = tag[:0]
			if len(m) > 0 {
				nn, ok := m[parentTag]
				if ok {
					if n, ok := nn.Tags[tags]; ok {
						if (count+n.Num) > limit && n.ChuckOvered {
							i = tagLocal
							break
						}
						count += n.Num
						continue
					}
				}
				if n, ok := m[tags]; ok {
					if (count+n.Num) > limit && n.ChuckOvered {
						i = tagLocal
						break
					}
					count += n.Num
				}
			}
			continue
		}
		if tagIn {
			tag = append(tag, runeContent[i])
			continue
		}
		currentTags := allTags[len(allTags)-1]
		mm, ok := m[currentTags]
		if !ok {
			count++
		} else if len(mm.EscapeCharacter) > 0 {
			if n, ok := mm.EscapeCharacter[runeContent[i]]; ok {
				if (count+n.Num) > limit && n.ChuckOvered {
					break
				}
				count += n.Num
			} else {
				count++
			}
		}
	}
	if i > runeTotal {
		i = runeTotal
	}
	content = string(runeContent[:i])
	closeTag = html.CloseTag(content)
	return content, closeTag
}
