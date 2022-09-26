package plugins

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github/fthvgb1/wp-go/cache"
	"github/fthvgb1/wp-go/helper"
	"github/fthvgb1/wp-go/models"
	"github/fthvgb1/wp-go/vars"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"
)

var removeWpBlock = regexp.MustCompile("<!-- /?wp:.*-->")
var more = regexp.MustCompile("<!--more(.*?)?-->")
var digestCache *cache.MapCache[uint64, string]
var quto = regexp.MustCompile(`&quot; *|&amp; *|&lt; *|&gt; ?|&nbsp; *`)

func InitDigest() {
	digestCache = cache.NewMapCache[uint64](digestRaw, time.Second)
}

func digestRaw(arg ...any) (string, error) {
	str := arg[0].(string)
	id := arg[1].(uint64)
	limit := vars.Conf.DigestWordCount
	if limit < 0 {
		return str, nil
	} else if limit == 0 {
		return "", nil
	}
	return DigestRaw(str, limit, id), nil
}

func DigestCache(ctx *gin.Context, id uint64, str string) string {
	content, _ := digestCache.GetCache(ctx, id, time.Second, str, id)
	return content
}

func DigestRaw(str string, limit int, id uint64) string {
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
		tmp := `%s......%s<p class="read-more"><a href="/p/%d">继续阅读</a></p>`
		if strings.Contains(closeTag, "pre") || strings.Contains(closeTag, "code") {
			tmp = `%s%s......<p class="read-more"><a href="/p/%d">继续阅读</a></p>`
		}
		content = fmt.Sprintf(tmp, content, closeTag, id)
	}

	return content
}

func Digest(p *Plugin[models.WpPosts], c *gin.Context, post *models.WpPosts, scene uint) {
	if scene == Detail {
		return
	}
	//post.PostContent = DigestCache(c, post.Id, post.PostContent)
	post.PostContent = DigestRaw(post.PostContent, vars.Conf.DigestWordCount, post.Id)
}
