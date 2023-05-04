package components

import (
	"github.com/fthvgb1/wp-go/app/theme/wp"
	"github.com/fthvgb1/wp-go/app/theme/wp/components/block"
	"github.com/fthvgb1/wp-go/app/wpconfig"
	str "github.com/fthvgb1/wp-go/helper/strings"
	"strings"
)

var blockFn = map[string]func(*wp.Handle, string, block.ParserBlock) (func() string, error){
	"core/categories": block.Category,
}

func Block(id string) (func(*wp.Handle) string, string) {
	content := wpconfig.GetPHPArrayVal("widget_block", "", str.ToInteger[int64](id, 0), "content")
	if content == "" {
		return nil, ""
	}
	var name string
	v := block.ParseBlock(content)
	if len(v.Output) > 0 {
		name = v.Output[0].Name
	}
	return func(h *wp.Handle) string {
		var out []string
		for _, parserBlock := range v.Output {
			fn, ok := blockFn[parserBlock.Name]
			if ok {
				s, err := fn(h, id, parserBlock)
				if err != nil {
					continue
				}
				out = append(out, s())

			}
		}
		return strings.Join(out, "\n")
	}, name
}
