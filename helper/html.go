package helper

import (
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

func htmlSpecialChars(text string, flags int) string {
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
func htmlSpecialCharsDecode(text string, flags int) string {
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
