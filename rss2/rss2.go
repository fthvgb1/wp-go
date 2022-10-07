package rss2

import (
	"fmt"
	"github/fthvgb1/wp-go/helper"
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
                <link>{$link}</link>
				{$comments}
                <dc:creator><![CDATA[{$author}]]></dc:creator>
                <pubDate>{$pubDate}</pubDate>
                {$category}
                <guid isPermaLink="false">{$guid}</guid>
                {$description}
                <content:encoded><![CDATA[{$content}]]></content:encoded>
                {$commentRss}
				{$commentNumber}

            </item>
`

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
		"{$items}": strings.Join(helper.SliceMap(r.Items, func(t Item) string {
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
		"{$title}":         i.Title,
		"{$link}":          i.Link,
		"{$comments}":      i.getComments(),
		"{$author}":        i.Creator,
		"{$pubDate}":       i.PubDate,
		"{$category}":      i.getCategory(),
		"{$guid}":          i.Guid,
		"{$description}":   i.getDescription(),
		"{$content}":       i.Content,
		"{$commentRss}":    i.getCommentRss(),
		"{$commentNumber}": i.getSlashComments(),
	} {
		xml = strings.Replace(xml, k, v, -1)
	}
	return
}

func (i Item) getCategory() string {
	r := ""
	if i.Category != "" {
		r = fmt.Sprintf("<category><![CDATA[%s]]></category>", i.CommentLink)
	}
	return r
}
func (i Item) getDescription() string {
	r := ""
	if i.Description != "" {
		r = fmt.Sprintf("<description><![CDATA[%s]]></description>", i.Description)
	}
	return r
}
func (i Item) getComments() string {
	r := ""
	if i.CommentLink != "" {
		r = fmt.Sprintf("<comments>%s</comments>", i.CommentLink)
	}
	return r
}

func (i Item) getCommentRss() (r string) {
	if i.CommentLink != "" && i.SlashComments > 0 {
		r = fmt.Sprintf("<wfw:commentRss>%s</wfw:commentRss>", i.CommentRss)
	}
	return
}
func (i Item) getSlashComments() (r string) {
	if i.SlashComments > 0 {
		r = fmt.Sprintf("<slash:comments>%d</slash:comments>", i.SlashComments)
	}
	return
}
