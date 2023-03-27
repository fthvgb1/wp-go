package block

import (
	"github.com/dlclark/regexp2"
	str "github.com/fthvgb1/wp-go/helper/strings"
)

type Block struct {
	Name        string
	Attrs       string
	Len         int
	StartOffset int
	Type        string
}

type BockParser struct {
	Document string
	Offset   int
	Output   []ParserBlock
}

type ParserBlock struct {
	Name         string
	Attrs        string
	InnerBlocks  string
	InnerHtml    string
	InnerContent string
}

var block = regexp2.MustCompile(`<!--\s+(?<closer>/)?wp:(?<namespace>[a-z][a-z0-9_-]*\/)?(?<name>[a-z][a-z0-9_-]*)\s+(?<attrs>{(?:(?:[^}]+|}+(?=})|(?!}\s+\/?-->).)*)?}\s+)?(?<void>\/)?-->`, regexp2.IgnoreCase|regexp2.Singleline)

func ParseBlock(content string) (r BockParser) {
	m, err := block.FindStringMatch(content)
	if err != nil {
		panic(err)
	}
	r.Document = content
	for m != nil {
		if m.GroupCount() < 1 {
			continue
		}

		b := token(m.Groups())
		bb := ParserBlock{}
		bb.Name = b.Name
		bb.Attrs = b.Attrs
		r.Output = append(r.Output, bb)
		m, _ = block.FindNextMatch(m)
	}
	return
}

func token(g []regexp2.Group) (b Block) {
	if len(g) < 1 {
		b.Type = "no-more-tokens"
		return
	}
	var closer, name, void, nameSpace = "", "", "", ""
	for i, group := range g {
		v := group.String()
		if v == "" {
			continue
		}
		switch i {
		case 0:
			b.Len = group.Length
			b.StartOffset = group.Index
		case 1:
			closer = v
		case 2:
			nameSpace = v
		case 3:
			name = v
		case 4:
			b.Attrs = v
		case 5:
			void = v
		default:
			continue
		}
	}
	if nameSpace == "" {
		nameSpace = "core/"
	}
	b.Name = str.Join(nameSpace, name)
	if void != "" {
		b.Type = "void-block"
		return
	}
	if closer != "" {
		b.Type = "block-closer"
		return
	}
	b.Type = "block-opener"
	return b
}
