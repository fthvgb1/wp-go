package plugins

import (
	"context"
	"fmt"
	"github.com/fthvgb1/wp-go/app/pkg/config"
	"github.com/fthvgb1/wp-go/app/pkg/logs"
	"github.com/fthvgb1/wp-go/app/pkg/models"
	"github.com/fthvgb1/wp-go/cache/cachemanager"
	"github.com/fthvgb1/wp-go/cache/reload"
	"github.com/fthvgb1/wp-go/helper"
	str "github.com/fthvgb1/wp-go/helper/strings"
	"github.com/fthvgb1/wp-go/plugin/digest"
	"github.com/fthvgb1/wp-go/safety"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"
)

var more = regexp.MustCompile("<!--more(.*?)?-->")

var removeWpBlock = regexp.MustCompile("<!-- /?wp:.*-->")

type DigestConfig struct {
	DigestWordCount    int    `yaml:"digestWordCount"`
	DigestAllowTag     string `yaml:"digestAllowTag"`
	DigestRegex        string `yaml:"digestRegex"`
	DigestTagOccupyNum []struct {
		Tag             string `yaml:"tag"`
		Num             int    `yaml:"num"`
		ChuckOvered     bool   `yaml:"chuckOvered"`
		EscapeCharacter []struct {
			Tags        string   `yaml:"tags"`
			Character   []string `yaml:"character"`
			Num         int      `yaml:"num"`
			ChuckOvered bool     `yaml:"chuckOvered"`
		} `yaml:"escapeCharacter"`
	} `yaml:"digestTagOccupyNum"`
	specialSolve map[string]digest.SpecialSolveConf
}

var digestConfig *safety.Var[DigestConfig]

func InitDigestCache() {
	cachemanager.NewMemoryMapCache(nil, digestRaw, config.GetConfig().CacheTime.DigestCacheTime, "digestPlugin", func() time.Duration {
		return config.GetConfig().CacheTime.DigestCacheTime
	})

	digestConfig = reload.VarsBy(func() DigestConfig {
		c, err := config.GetCustomizedConfig[DigestConfig]()
		if err != nil {
			logs.Error(err, "get digest config fail")
			c.DigestWordCount = config.GetConfig().DigestWordCount
			c.DigestAllowTag = config.GetConfig().DigestAllowTag
			return c
		}
		if c.DigestRegex != "" {
			digest.SetQutos(c.DigestRegex)
		}
		if len(c.DigestTagOccupyNum) <= 1 {
			return c
		}
		c.specialSolve = map[string]digest.SpecialSolveConf{}
		for _, item := range c.DigestTagOccupyNum {
			tags := strings.Split(strings.ReplaceAll(item.Tag, " ", ""), "<")
			for _, tag := range tags {
				if tag == "" {
					continue
				}
				tag = str.Join("<", tag)
				var ec map[rune]digest.SpecialSolve
				var specialTags map[string]digest.SpecialSolve
				if len(item.EscapeCharacter) > 0 {
					ec = make(map[rune]digest.SpecialSolve)
					for _, esc := range item.EscapeCharacter {
						for _, i := range esc.Character {
							s := []rune(i)
							if len(s) == 1 {
								ec[s[0]] = digest.SpecialSolve{
									Num:         esc.Num,
									ChuckOvered: esc.ChuckOvered,
								}
							}
						}
						if esc.Tags == "" {
							continue
						}
						tagss := strings.Split(strings.ReplaceAll(esc.Tags, " ", ""), "<")
						specialTags = make(map[string]digest.SpecialSolve)
						for _, t := range tagss {
							if t == "" {
								continue
							}
							t = str.Join("<", t)
							specialTags[t] = digest.SpecialSolve{
								Num:         esc.Num,
								ChuckOvered: esc.ChuckOvered,
							}
						}
					}

				}
				c.specialSolve[tag] = digest.SpecialSolveConf{
					Num:             item.Num,
					ChuckOvered:     item.ChuckOvered,
					EscapeCharacter: ec,
					Tags:            specialTags,
				}
			}
		}
		return c
	}, "digestConfig")
}

func RemoveWpBlock(s string) string {
	return removeWpBlock.ReplaceAllString(s, "")
}

func digestRaw(ctx context.Context, id uint64, arg ...any) (string, error) {
	s := arg[1].(string)
	limit := arg[3].(int)
	if limit < 0 {
		return s, nil
	} else if limit == 0 {
		return "", nil
	}

	s = more.ReplaceAllString(s, "")
	fn := helper.GetContextVal(ctx, "postMoreFn", PostsMore)
	return Digests(s, id, limit, fn), nil
}

func Digests(content string, id uint64, limit int, fn func(id uint64, content, closeTag string) string) string {
	closeTag := ""
	content = RemoveWpBlock(content)
	c := digestConfig.Load()
	tag := c.DigestAllowTag
	if tag == "" {
		tag = "<a><b><blockquote><br><cite><code><dd><del><div><dl><dt><em><h1><h2><h3><h4><h5><h6><i><img><li><ol><p><pre><span><strong><ul>"
	}
	content = digest.StripTags(content, tag)
	length := utf8.RuneCountInString(content) + 1
	if length <= limit {
		return content
	}
	if len(c.specialSolve) > 0 {
		content, closeTag = digest.CustomizeHtml(content, limit, c.specialSolve)
	} else {
		content, closeTag = digest.Html(content, limit)
	}

	if fn == nil {
		return PostsMore(id, content, closeTag)
	}
	return fn(id, content, closeTag)
}

func PostsMore(id uint64, content, closeTag string) string {
	tmp := `%s......%s<p class="read-more"><a href="/p/%d">继续阅读</a></p>`
	if strings.Contains(closeTag, "pre") || strings.Contains(closeTag, "code") {
		tmp = `%s%s......<p class="read-more"><a href="/p/%d">继续阅读</a></p>`
	}
	content = fmt.Sprintf(tmp, content, closeTag, id)
	return content
}

func Digest(ctx context.Context, post *models.Posts, limit int) {
	content, _ := cachemanager.GetBy[string]("digestPlugin", ctx, post.Id, time.Second, ctx, post.PostContent, post.Id, limit)
	post.PostContent = content
}

func PostExcerpt(post *models.Posts) {
	post.PostContent = strings.Replace(post.PostExcerpt, "\n", "<br/>", -1)
}
