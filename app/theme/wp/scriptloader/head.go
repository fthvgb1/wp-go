package scriptloader

import (
	"encoding/json"
	"github.com/fthvgb1/wp-go/app/cmd/reload"
	"github.com/fthvgb1/wp-go/app/pkg/config"
	"github.com/fthvgb1/wp-go/app/pkg/logs"
	"github.com/fthvgb1/wp-go/app/theme/wp"
	"github.com/fthvgb1/wp-go/helper/slice"
	str "github.com/fthvgb1/wp-go/helper/strings"
	"os"
	"path/filepath"
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
		path := filepath.Join(i.path)
		css := reload.GetAnyValMapBy("script-loader-MaybeInlineStyles", i.handle, path, func(a string) string {
			css, err := os.ReadFile(i.path)
			if err != nil {
				logs.Error(err, "read file ", i.path)
				return ""
			}
			return string(css)
		})

		s, _ := __styles.Load(i.handle)
		s.Src = ""
		a := s.Extra["after"]
		if a == nil {
			a = []string{}
		}
		slice.Unshift(&a, css)
		s.Extra["after"] = a
	}
}

func emojiDetectionScript(h *wp.Handle) {
	settings := map[string]any{
		"baseUrl": "https://s.w.org/images/core/emoji/14.0.0/72x72/",
		"ext":     ".png",
		"svgUrl":  "https://s.w.org/images/core/emoji/14.0.0/svg/", "svgExt": ".svg",
		"source": map[string]any{
			"concatemoji": "/wp-includes/js/wp-emoji-release.min.js?ver=6.2.2",
		},
	}
	setting, _ := json.Marshal(settings)
	dir := config.GetConfig().WpDir
	emotion := reload.GetAnyValBys("script-loader-emoji", struct{}{}, func(_ struct{}) string {
		f, err := os.ReadFile(dir)
		if err != nil {
			logs.Error(err, "load emoji css fail", dir)
			return ""
		}
		return string(f)
	})
	s := str.Join("window._wpemojiSettings = ", string(setting), "\n", emotion)
	PrintInlineScriptTag(h, s, nil)
}

func PrintInlineScriptTag(h *wp.Handle, script string, attr map[string]string) {
	ss := wp.GetComponentsArgs(h, "inlineScript", "")
	s := str.NewBuilder()
	s.WriteString(ss)
	s.WriteString("<script")
	for k, v := range attr {
		s.Sprintf(` %s="%s"`, k, v)
	}
	s.Sprintf(">%s</script>\n", script)
	wp.SetComponentsArgs(h, "inlineScript", s.String())
}

func PrintStyles(h *wp.Handle) {

}
