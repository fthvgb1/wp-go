package scriptloader

import (
	"encoding/json"
	"fmt"
	"github.com/fthvgb1/wp-go/app/cmd/reload"
	"github.com/fthvgb1/wp-go/app/pkg/config"
	"github.com/fthvgb1/wp-go/app/pkg/constraints"
	"github.com/fthvgb1/wp-go/app/pkg/logs"
	"github.com/fthvgb1/wp-go/app/theme/wp"
	"github.com/fthvgb1/wp-go/app/theme/wp/components/widget"
	"github.com/fthvgb1/wp-go/app/wpconfig"
	"github.com/fthvgb1/wp-go/helper"
	"github.com/fthvgb1/wp-go/helper/number"
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

func PrintInlineStyles(handle string) string {
	o, _ := __styles.Load(handle)
	out := o.getData("after")
	if out == "" {
		return ""
	}
	return fmt.Sprintf("<style id='%s-inline-css'%s>\n%s\n</style>\n", handle, "", out)
}

func PrintStyle(h *wp.Handle, s ...string) {
	out := wp.GetComponentsArgs(h, "wp_style_out", str.NewBuilder())
	out.WriteString(s...)
}
func PrintHead(h *wp.Handle, s ...string) {
	out := wp.GetComponentsArgs(h, "wp_head", str.NewBuilder())
	out.WriteString(s...)
}

func LinkHead(h *wp.Handle) {
	PrintHead(h, "<link rel=\"https://api.w.org/\" href=\"/wp-json\" />")
	if s := restGetQueriedResourceRoute(h); s != "" {
		PrintHead(h, "<link rel=\"alternate\" type=\"application/json\" href=", s, " />")
	}
}

func restGetQueriedResourceRoute(h *wp.Handle) string {
	if cate, ok := widget.IsCategory(h); ok {
		return fmt.Sprintf("/wp/v2/categories/%d", cate.Terms.TermId)
	}
	if tag, ok := widget.IsTag(h); ok {
		return fmt.Sprintf("/wp/v2/tags/%d", tag.Terms.TermId)
	}
	return ""
}

func RsdLink(h *wp.Handle) {
	PrintHead(h, fmt.Sprintf("<link rel=\"EditURI\" type=\"application/rsd+xml\" title=\"RSD\" href=\"%s\" />\n", "xmlrpc.php?rsd"))
}

func WlwmanifestLink(h *wp.Handle) {
	PrintHead(h, fmt.Sprintf("<link rel=\"wlwmanifest\" type=\"application/wlwmanifest+xml\" href=\"%s\" />\n", "/wp-includes/wlwmanifest.xml"))
}

func LocaleStylesheet(h *wp.Handle) {
	uri := reload.GetAnyValBys("printHead-localStylesheet", h, func(a *wp.Handle) string {
		ur := str.Join("wp-content/themes", h.Theme(), str.Join(wpconfig.GetLang(), ".css"))
		path := filepath.Join(config.GetConfig().WpDir, ur)
		if helper.FileExist(path) {
			return str.Join("/", ur)
		}
		return ""
	})
	if uri != "" {
		PrintHead(h, fmt.Sprintf("<link rel=\"stylesheet\" href=\"%s\"%s media=\"screen\" />", uri, ""))
	}
}

func TheGenerator(h *wp.Handle) {
	PrintHead(h, fmt.Sprintf(`<meta name="generator" content="WordPress %s"/>`, "6.2.2"))
}

func ShortLinkWpHead(h *wp.Handle) {
	if h.Scene() != constraints.Detail || h.Detail.Post.Id < 1 {
		return
	}
	shortlink := ""
	post := h.Detail.Post
	if post.PostType == "page" && wpconfig.GetOption("page_on_front") == number.IntToString(post.Id) &&
		wpconfig.GetOption("show_on_front") == "page" {
		shortlink = "/"
	} else {
		shortlink = str.Join("/p/", number.IntToString(post.Id))
	}
	if shortlink != "" {
		PrintHead(h, fmt.Sprintf(`<link rel='shortlink' href="%s" />`, shortlink))
	}
}

func customLogoHeaderStyles(h *wp.Handle) {
	mod := h.CommonThemeMods()
	if !mod.ThemeSupport.CustomHeader.HeaderText && mod.ThemeSupport.CustomLogo.HeaderText != "" {
		class := mod.ThemeSupport.CustomLogo.HeaderText
		attr := ""
		if !slice.IsContained(mod.ThemeSupport.HTML5, "style") {
			attr = ` type="text/css"`
		}
		PrintHead(h, fmt.Sprintf(`<style id="custom-logo-css"%s>
	.%s {
			position: absolute;
			clip: rect(1px, 1px, 1px, 1px);
		}
</style>`, attr, class))
	}
}

func PrintHeadToStr(h *wp.Handle) string {
	h.DoActionFilter("wp_head", "", h)
	return wp.GetComponentsArgs(h, "wp_head", str.NewBuilder()).String()
}
