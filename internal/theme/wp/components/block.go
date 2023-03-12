package components

import (
	"github.com/dlclark/regexp2"
)

type Block struct {
	Closer      string
	NameSpace   string
	Name        string
	Attrs       string
	Void        string
	Len         int
	StartOffset int
}

var block = regexp2.MustCompile(`<!--\s+(?<closer>/)?wp:(?<namespace>[a-z][a-z0-9_-]*\/)?(?<name>[a-z][a-z0-9_-]*)\s+(?<attrs>{(?:(?:[^}]+|}+(?=})|(?!}\s+\/?-->).)*)?}\s+)?(?<void>\/)?-->`, regexp2.IgnoreCase|regexp2.Singleline)

func ParseBlock(s string) []Block {
	m, err := block.FindStringMatch(s)
	if err != nil {
		panic(err)
	}
	var blocks []Block
	for m != nil {
		if m.GroupCount() < 1 {
			continue
		}

		b, _ := token(m.Groups())
		b.StartOffset = m.Group.Index
		b.Len = m.Length
		blocks = append(blocks, b)
		m, _ = block.FindNextMatch(m)
	}
	return blocks
}

func token(g []regexp2.Group) (Block, string) {
	if len(g) < 1 {
		return Block{}, "no-more-tokens"
	}
	b := Block{NameSpace: "core/"}
	for i, group := range g {
		v := group.String()
		if v == "" {
			continue
		}
		switch i {
		case 1:
			b.Closer = v
		case 2:
			b.NameSpace = v
		case 3:
			b.Name = v
		case 4:
			b.Attrs = v
		case 5:
			b.Void = v
		default:
			continue
		}
	}
	if b.Void != "" {
		return b, ""
	}
	if b.Closer != "" {
		return b, "block-closer"
	}
	return b, "block-opener"
}
