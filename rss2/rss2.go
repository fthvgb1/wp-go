package rss2

import (
	"fmt"
	"github.com/fthvgb1/wp-go/helper/slice"
	"strconv"
	"strings"
)

var template = `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0"
     xmlns:content="http://purl.org/rss/1.0/modules/content/"
     xmlns:wfw="http://wellformedweb.org/CommentAPI/"
     xmlns:dc="http://purl.org/dc/elements/1.1/"
     xmlns:atom="http://www.w3.org/2005/Atom"
     xmlns:sy="http://purl.org/rss/1.0/modules/syndication/"
     xmlns:slash="http://purl.org/rss/1.0/modules/slash/"
>

    <channel>
        <title>{$title}</title>
        <atom:link href="{$feedLink}" rel="self" type="application/rss+xml"/>
        <link>{$link}</link>
        <description>{$description}</description>
        <lastBuildDate>{$lastBuildDate}</lastBuildDate>
        <language>{$lang}</language>
        <sy:updatePeriod>
            {$updatePeriod}
        </sy:updatePeriod>
        <sy:updateFrequency>
            {$updateFrequency}
        </sy:updateFrequency>
        <generator>{$generator}</generator>
		{$items}

    </channel>
</rss>
`
var templateItems = `
			<item>
                <title>{$title}</title>
                {$link}
				{$comments}
                {$creator}
                <pubDate>{$pubDate}</pubDate>
                {$category}
                {$guid}
                {$description}
                <content:encoded><![CDATA[{$content}]]></content:encoded>
                {$commentRss}
				{$commentNumber}
            </item>
`
var templateReplace = map[string]string{
	"{$category}":      "<category><![CDATA[%s]]></category>",
	"{$link}":          "<link>%s</link>",
	"{$creator}":       "<dc:creator><![CDATA[%s]]></dc:creator>",
	"{$description}":   "<description><![CDATA[%s]]></description>",
	"{$comments}":      "<comments>%s</comments>",
	"{$commentRss}":    "<wfw:commentRss>%s</wfw:commentRss>",
	"{$commentNumber}": "<slash:comments>%s</slash:comments>",
	"{$guid}":          "<guid isPermaLink=\"false\">%s</guid>",
}

type Rss2 struct {
	Title           string
	AtomLink        string
	Link            string
	Description     string
	LastBuildDate   string
	Language        string
	UpdatePeriod    string
	UpdateFrequency int
	Generator       string
	Items           []Item
}

type Item struct {
	Title         string
	Link          string
	CommentLink   string
	Creator       string
	PubDate       string
	Category      string
	Guid          string
	Description   string
	Content       string
	CommentRss    string
	SlashComments int
}

func (r Rss2) GetXML() (xml string) {
	xml = template
	for k, v := range map[string]string{
		"{$title}":           r.Title,
		"{$link}":            r.Link,
		"{$feedLink}":        r.AtomLink,
		"{$description}":     r.Description,
		"{$lastBuildDate}":   r.LastBuildDate,
		"{$lang}":            r.Language,
		"{$updatePeriod}":    r.UpdatePeriod,
		"{$updateFrequency}": fmt.Sprintf("%d", r.UpdateFrequency),
		"{$generator}":       r.Generator,
		"{$items}": strings.Join(slice.Map(r.Items, func(t Item) string {
			return t.GetXml()
		}), ""),
	} {
		xml = strings.Replace(xml, k, v, -1)
	}
	return
}

func (i Item) GetXml() (xml string) {
	xml = templateItems
	for k, v := range map[string]string{
		"{$title}":   i.Title,
		"{$pubDate}": i.PubDate,
		"{$content}": i.Content,
	} {
		xml = strings.Replace(xml, k, v, -1)
	}

	m := map[string]string{
		"{$category}":    i.Category,
		"{$link}":        i.Link,
		"{$creator}":     i.Creator,
		"{$description}": i.Description,
		"{$comments}":    i.CommentLink,
		"{$guid}":        i.Guid,
		"{$commentRss}":  i.CommentRss,
	}

	if i.CommentRss != "" && i.SlashComments > 0 {
		m["{$commentRss}"] = i.CommentRss
	} else {
		m["{$commentRss}"] = ""
	}
	if i.SlashComments > 0 {
		m["{$commentNumber}"] = strconv.Itoa(i.SlashComments)
	} else {
		m["{$commentNumber}"] = ""
	}

	for k, v := range m {
		t := ""
		if v != "" {
			t = fmt.Sprintf(templateReplace[k], v)
		}
		xml = strings.Replace(xml, k, t, -1)
	}

	return
}
