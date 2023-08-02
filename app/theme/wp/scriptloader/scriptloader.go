package scriptloader

import (
	"encoding/json"
	"fmt"
	"github.com/fthvgb1/wp-go/app/cmd/reload"
	"github.com/fthvgb1/wp-go/app/theme/wp"
	"github.com/fthvgb1/wp-go/app/wpconfig"
	"github.com/fthvgb1/wp-go/helper"
	"github.com/fthvgb1/wp-go/helper/maps"
	"github.com/fthvgb1/wp-go/helper/slice"
	str "github.com/fthvgb1/wp-go/helper/strings"
	"github.com/fthvgb1/wp-go/safety"
	"net/url"
	"path/filepath"
	"strings"
)

var __styles = reload.MapBy(func(m *safety.Map[string, *Style]) {
	defaultStyles(m, ".css")
})
var __scripts = reload.MapBy[string, *Script](func(m *safety.Map[string, *Script]) {
	suffix := ".min"
	defaultScripts(m, suffix)

})

type Style struct {
	Dependencies
}

type Script struct {
	Dependencies
}

func addScript(handle string, src string, deps []string, ver string, args any) {
	script := NewScript(handle, src, deps, ver, args)
	__scripts.Store(handle, script)
}

func localize(handle, objectname string, l10n map[string]any) string {
	if "jquery" == handle {
		handle = "jquery-core"
	}
	after, ok := maps.GetStrAnyVal[string](l10n, "l10n_print_after")
	if ok {
		delete(l10n, "l10n_print_after")
	}
	v, _ := json.Marshal(l10n)
	script := fmt.Sprintf("var %s = %s;", objectname, string(v))
	if after != "" {
		script = str.Join(script, "\n", after, ";")
	}
	return script
}

func AddStaticLocalize(handle, objectname string, l10n map[string]any) {
	AddScriptData(handle, "data", localize(handle, objectname, l10n))
}
func AddDynamicLocalize(h *wp.Handle, handle, objectname string, l10n map[string]any) {
	AddDynamicData(h, handle, "data", localize(handle, objectname, l10n))
}

func (d *Dependencies) getData(key string) string {
	return strings.Join(d.Extra[key], "\n")
}
func GetData(h *wp.Handle, handle, key string) string {
	hh, ok := __scripts.Load(handle)
	if !ok {
		return ""
	}
	d := hh.Extra[key]
	d = append(d, GetDynamicData(h, handle, key))
	return strings.Join(d, "\n")
}

func AddScriptData(handle, key, data string) {
	var s *Script
	var ok bool
	s, ok = __scripts.Load(handle)
	if !ok {
		s = NewScript(handle, "", nil, "", nil)
	}
	if s.Extra == nil {
		s.Extra = make(map[string][]string)
	}
	s.Extra[key] = append(s.Extra[key], data)
}

func AddStyleData(handle, key, data string) {
	var s *Style
	var ok bool
	s, ok = __styles.Load(handle)
	if !ok {
		s = NewStyle(handle, "", nil, "", nil)
	}
	if s.Extra == nil {
		s.Extra = make(map[string][]string)
	}
	s.Extra[key] = append(s.Extra[key], data)
}

func AddInlineScript(handle, data, position string) {
	if handle == "" || data == "" {
		return
	}
	if position != "after" {
		position = "before"
	}
	AddScriptData(handle, position, data)
}

func AddInlineStyle(handle, data string) {
	if handle == "" || data == "" {
		return
	}
	AddStyleData(handle, "after", data)
}

func InlineScripts(handle, position string, display bool) string {
	v, _ := __scripts.Load(handle)
	ss := v.getData(position)
	if ss == "" {
		return ""
	}
	scp := strings.Trim(ss, "\n")
	if display {
		return fmt.Sprintf("<script id='%s-js-%s'>\n%s\n</script>\n", handle, position, scp)
	}
	return scp
}

func AddScript(handle string, src string, deps []string, ver string, args any) {
	script := NewScript(handle, src, deps, ver, args)
	__scripts.Store(handle, script)
}

const (
	style = iota
	script
)

var scriptQueues = reload.Vars(scriptQueue{})
var styleQueues = reload.Vars(scriptQueue{})

type scriptQueue struct {
	Register             map[string]struct{}
	Queue                []string
	Args                 map[string]string
	queuedBeforeRegister map[string]string
}

func EnqueueStyle(handle, src string, deps []string, ver, media string) {
	if media == "" {
		media = "all"
	}

	h := strings.Split(handle, "?")
	if src != "" {
		AddScript(h[0], src, deps, ver, media)
	}
	enqueue(handle, style)
}
func EnqueueStyles(handle, src string, deps []string, ver, media string) {
	if src != "" {
		src = GetThemeFileUri(src)
	}
	EnqueueStyle(handle, src, deps, ver, media)
}
func EnqueueScript(handle, src string, deps []string, ver string, inFooter bool) {
	h := strings.Split(handle, "?")
	if src != "" {
		AddScript(h[0], src, deps, ver, nil)
	}
	if inFooter {
		AddScriptData(h[0], "group", "1")
	}
	enqueue(handle, script)
}
func EnqueueScripts(handle, src string, deps []string, ver string, inFooter bool) {
	if src != "" {
		src = GetThemeFileUri(src)
	}
	EnqueueScript(handle, src, deps, ver, inFooter)
}

func enqueue(handle string, t int) {
	h := strings.Split(handle, "?")
	ss := styleQueues.Load()
	if t == 1 {
		ss = scriptQueues.Load()
	}
	if slice.IsContained(ss.Queue, h[0]) && maps.IsExists(ss.Register, h[0]) {
		ss.Queue = append(ss.Queue, h[0])
	} else if maps.IsExists(ss.Register, h[0]) {
		ss.queuedBeforeRegister[h[0]] = ""
		if len(h) > 1 {
			ss.queuedBeforeRegister[h[0]] = h[1]
		}
	}
}

func GetStylesheetUri() string {
	return GetThemeFileUri("/styles.css")
}

func GetThemeFileUri(file string) string {
	return filepath.Join("/wp-content/themes", wpconfig.GetOption("template"), file)
}

type Dependencies struct {
	Handle           string              `json:"handle,omitempty"`
	Src              string              `json:"src,omitempty"`
	Deps             []string            `json:"deps,omitempty"`
	Ver              string              `json:"ver,omitempty"`
	Args             any                 `json:"args,omitempty"`
	Extra            map[string][]string `json:"extra,omitempty"`
	Textdomain       string              `json:"textdomain,omitempty"`
	TranslationsPath string              `json:"translations_path,omitempty"`
}

func NewScript(handle string, src string, deps []string, ver string, args any) *Script {
	return &Script{Dependencies{Handle: handle, Src: src, Deps: deps, Ver: ver, Args: args}}
}

func NewStyle(handle string, src string, deps []string, ver string, args any) *Style {
	return &Style{Dependencies{Handle: handle, Src: src, Deps: deps, Ver: ver, Args: args}}
}

func AddDynamicData(h *wp.Handle, handle, key, data string) {
	da := helper.GetContextVal(h.C, "__scriptDynamicData__", map[string]map[string][]string{})
	m, ok := da[handle]
	if !ok {
		m = map[string][]string{}
	}
	m[key] = append(m[key], data)
	da[handle] = m
}

func GetDynamicData(h *wp.Handle, handle, key string) string {
	da := helper.GetContextVal(h.C, "__scriptDynamicData__", map[string]map[string][]string{})
	if len(da) < 1 {
		return ""
	}
	m, ok := da[handle]
	if !ok {
		return ""
	}
	mm, ok := m[key]
	if !ok {
		return ""
	}
	return strings.Join(mm, "\n")
}

func SetTranslation(handle, domain, path string) {
	hh, ok := __scripts.Load(handle)
	if !ok {
		return
	}
	if !slice.IsContained(hh.Deps, handle) {
		hh.Deps = append(hh.Deps, "wp-i18n")
	}
	if domain == "" {
		domain = "default"
	}
	hh.Textdomain = domain
	hh.TranslationsPath = path
}

func (item *__parseLoadItem) allDeps(handles []string, recursion bool, group int) bool {
	for _, handle := range handles {
		parts := strings.Split(handle, "?")
		queued := slice.IsContained(item.todo, parts[0])
		handle = parts[0]
		moved := item.setGroup(handle, group)
		if queued || !moved {
			continue
		}
		newGroup := item.groups[handle]
		keepGoing := true
		h, ok := __styles.Load(handle)
		if !ok {
			keepGoing = false
		}
		if len(h.Deps) > 0 && len(slice.Diff(h.Deps, __styles.Keys())) > 0 {
			keepGoing = false
		}
		if len(h.Deps) > 0 && item.allDeps(h.Deps, true, newGroup) {
			keepGoing = false
		}
		if !keepGoing {
			if recursion {
				return false
			} else {
				continue
			}
		}
		if len(parts) > 1 {
			item.args[handle] = parts[1]
		}
		item.todo = append(item.todo, handle)
	}
	return true
}

type __parseLoadItem struct {
	todo          []string
	done          []string
	groups        map[string]int
	args          map[string]string
	textDirection string
	concat        string
	doConcat      bool
}

func newParseLoadItem() *__parseLoadItem {
	return &__parseLoadItem{
		groups: map[string]int{},
	}
}

func (item *__parseLoadItem) setGroup(handle string, group int) bool {
	if v, ok := item.groups[handle]; ok && v <= group {
		return false
	}
	item.groups[handle] = group
	return true
}

func DoStyleItems(h *wp.Handle, handles []string, group int) []string {
	item := newParseLoadItem()
	item.allDeps(handles, false, 0)
	for i, handle := range item.todo {
		_, ok := __styles.Load(handle)
		if !slice.IsContained(item.done, handle) && ok {
			if DoStyleItem(h, item, handle, group) {
				item.done = append(item.done, handle)
			}
			slice.Delete(&item.todo, i)
		}
	}
	return item.done
}

func (s *Style) DoHeadItems() {

}
func (s *Style) DoItems(handle string) {

}

func DoStyleItem(h *wp.Handle, item *__parseLoadItem, handle string, group int) bool {
	obj, _ := __styles.Load(handle)
	ver := obj.Ver
	if item.args[handle] != "" {
		str.Join(ver, "&amp;", item.args[handle])
	}
	src := obj.Src
	var condBefore, condAfter, conditional, _ string
	if v, ok := obj.Extra["conditional"]; ok && v != nil {
		conditional = v[0]
	}
	if conditional != "" {
		condBefore = str.Join("<!==[if ", conditional, "]>\n")
		condAfter = "<![endif]-->\n"
	}
	inlineStyle := item.PrintInline(handle)
	if inlineStyle != "" {
		_ = fmt.Sprintf("<style id='%s-inline-css'%s>\n%s\n</style>\n", handle, "", inlineStyle)
	}
	href := item.CssHref(src, ver)
	ref := "stylesheet"
	if v, ok := obj.Extra["alt"]; ok && v != nil {
		ref = "alternate stylesheet"
	}
	title := ""
	if v, ok := obj.Extra["title"]; ok && v != nil {
		title = str.Join(" title='", v[len(v)-1], "'")
	}
	tag := fmt.Sprintf("<link rel='%s' id='%s-css'%s href='%s'%s media='%s' />\n", ref, handle, title, href, "", item.args)

	if !item.doConcat {
		PrintStyle(h, condBefore, tag, PrintInlineStyles(handle), condAfter)
	}
	return true
}

func (item *__parseLoadItem) CssHref(src, ver string) string {
	if ver != "" {
		u, _ := url.Parse(src)
		v := u.Query()
		v.Set("ver", ver)
		u.RawQuery = v.Encode()
		src = u.String()
	}
	return src
}

func (item *__parseLoadItem) PrintInline(handle string) string {
	sty, _ := __styles.Load(handle)
	out := sty.getData("after")
	if out == "" {
		return ""
	}
	s := str.NewBuilder()
	s.Sprintf("<style id='%s-inline-css'%s>\n%s\n</style>\n", handle)
	return s.String()
}
