package scriptloader

import (
	"github.com/fthvgb1/wp-go/app/pkg/logs"
	"github.com/fthvgb1/wp-go/app/theme/wp"
	"github.com/fthvgb1/wp-go/helper/slice"
	"os"
)

type _style struct {
	handle string
	src    string
	path   string
	size   int64
}

func MaybeInlineStyles(h *wp.Handle) {
	totalInlineLimit := int64(0)
	var styles []_style
	ss := styleQueues.Load()
	for _, que := range ss.Queue {
		p, ok := __styles.Load(que)
		if !ok {
			continue
		}
		f, ok := p.Extra["path"]
		if !ok || f == nil {
			continue
		}
		ff := f[0]
		stat, err := os.Stat(ff)
		if err != nil {
			return
		}
		styles = append(styles, _style{
			handle: que,
			src:    p.Src,
			path:   ff,
			size:   stat.Size(),
		})
	}
	if len(styles) < 1 {
		return
	}
	slice.Sort(styles, func(i, j _style) bool {
		return i.size > j.size
	})
	totalInlineSize := int64(0)
	for _, i := range styles {
		if totalInlineSize+i.size > totalInlineLimit {
			break
		}
		css, err := os.ReadFile(i.path)
		if err != nil {
			logs.Error(err, "read file ", i.path)
			continue
		}
		s, _ := __styles.Load(i.handle)
		s.Src = ""
		s.Extra["after"] = append(s.Extra["after"], string(css))
	}
}
