package plugins

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github/fthvgb1/wp-go/helper"
	"github/fthvgb1/wp-go/models"
	"regexp"
	"strings"
	"unicode/utf8"
)

var removeWpBlock = regexp.MustCompile("<!-- /?wp:.*-->")
var more = regexp.MustCompile("<!--more(.*?)?-->")
var tag = regexp.MustCompile(`<.*?>`)
var limit = 300

func ExceptRaw(str string, limit, id int) string {

	if r := more.FindString(str); r != "" {
		m := strings.Split(str, r)
		str = m[0]
		return ""
	}
	content := removeWpBlock.ReplaceAllString(str, "")
	content = strings.Trim(content, " \t\n\r\000\x0B")
	content = strings.Replace(content, "]]>", "]]&gt;", -1)
	content = helper.StripTags(content, "<a><b><blockquote><br><cite><code><dd><del><div><dl><dt><em><h1><h2><h3><h4><h5><h6><i><img><li><ol><p><pre><span><strong><ul>")
	length := utf8.RuneCountInString(content) + 1
	if length > limit {
		start, l := 0, limit+1
		end := l
		for {
			txt := string([]rune(content)[start:end])
			count := 0
			for _, ints := range tag.FindAllStringIndex(txt, -1) {
				t := []rune(content[ints[0]:ints[1]])
				count += len(t)
				l += len(t)
			}

			if count > 0 && length > l {
				start = end
				end += count
			} else if count > 0 && length < l {
				break
			} else {
				content = string([]rune(content)[:end])
				content = fmt.Sprintf(`%s...<p class="read-more"><a href="/p/%d">继续阅读</a></p>`, content, id)
				break
			}
		}
	}
	return content
}

func Except(p *Plugin[models.WpPosts], c *gin.Context, post *models.WpPosts, scene uint) {
	if scene == Detail {
		return
	}
	post.PostContent = ExceptRaw(post.PostContent, limit, int(post.Id))

}
