package scriptloader

import (
	"encoding/json"
	"fmt"
	"github.com/dlclark/regexp2"
	"github.com/fthvgb1/wp-go/app/cmd/reload"
	"github.com/fthvgb1/wp-go/app/theme/wp"
	"github.com/fthvgb1/wp-go/app/wpconfig"
	"github.com/fthvgb1/wp-go/helper"
	"github.com/fthvgb1/wp-go/helper/maps"
	"github.com/fthvgb1/wp-go/helper/number"
	"github.com/fthvgb1/wp-go/helper/slice"
	str "github.com/fthvgb1/wp-go/helper/strings"
	"github.com/fthvgb1/wp-go/safety"
	"html"
	"math"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var __styles = reload.MapBy(func(m *safety.Map[string, *Script]) {
	defaultStyles(m, ".css")
})
var __scripts = reload.MapBy[string, *Script](func(m *safety.Map[string, *Script]) {
	suffix := ".min"
	defaultScripts(m, suffix)

})

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
	AddData(handle, "data", localize(handle, objectname, l10n))
}
func AddDynamicLocalize(h *wp.Handle, handle, objectname string, l10n map[string]any) {
	AddDynamicData(h, handle, "data", localize(handle, objectname, l10n))
}

func getData(handle, key string) string {
	h, ok := __scripts.Load(handle)
	if !ok {
		return ""
	}
	return strings.Join(h.Extra[key], "\n")
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

func AddData(handle, key, data string, t ...int) {
	var s *Script
	var ok bool
	if t != nil {
		s, ok = __styles.Load(handle)
	} else {
		s, ok = __scripts.Load(handle)
	}
	if !ok {
		s = NewScript(handle, "", nil, "", nil)
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
	AddData(handle, position, data)
}

func AddInlineStyle(handle, data string) {
	if handle == "" || data == "" {
		return
	}
	AddData(handle, "after", data, style)
}

func InlineScripts(handle, position string, display bool) string {
	ss := getData(handle, position)
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
		AddData(h[0], "group", "1")
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

type Script struct {
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
	return &Script{Handle: handle, Src: src, Deps: deps, Ver: ver, Args: args}
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

var __elements = map[string]string{
	"link":    "a:where(:not(.wp-element-button))", // The `where` is needed to lower the specificity.
	"heading": "h1, h2, h3, h4, h5, h6",
	"h1":      "h1",
	"h2":      "h2",
	"h3":      "h3",
	"h4":      "h4",
	"h5":      "h5",
	"h6":      "h6",
	// We have the .wp-block-button__link class so that this will target older buttons that have been serialized.
	"button": ".wp-element-button, .wp-block-button__link",
	// The block classes are necessary to target older content that won't use the new class names.
	"caption": ".wp-element-caption, .wp-block-audio figcaption, .wp-block-embed figcaption, .wp-block-gallery figcaption, .wp-block-image figcaption, .wp-block-table figcaption, .wp-block-video figcaption",
	"cite":    "cite",
}

const RootBlockSelector = "body"

type node struct {
	Path     []string
	Selector string
	Name     string
}

var __validElementPseudoSelectors = map[string][]string{
	"link":   {":link", ":any-link", ":visited", ":hover", ":focus", ":active"},
	"button": {":link", ":any-link", ":visited", ":hover", ":focus", ":active"},
}

var blockSupportFeatureLevelSelectors = map[string]string{
	"__experimentalBorder": "border",
	"color":                "color",
	"spacing":              "spacing",
	"typography":           "typography",
}

func appendToSelector(selector, toAppend, position string) string {
	s := strings.Split(selector, ",")
	if position == "" {
		position = "right"
	}
	return strings.Join(slice.Map(s, func(t string) string {
		var l, r string
		if position == "right" {
			l = t
			r = toAppend
		} else {
			l = toAppend
			r = t
		}
		return str.Join(l, r)
	}), ",")
}

func __removeComment(m map[string]any) {
	delete(m, "//")
	for _, v := range m {
		mm, ok := v.(map[string]any)
		if ok {
			__removeComment(mm)
		}
	}
}

func scopeSelector(scope, selector string) string {
	scopes := strings.Split(scope, ",")
	selectors := strings.Split(selector, ",")
	var a []string
	for _, outer := range scopes {
		outer = strings.TrimSpace(outer)
		for _, inner := range selectors {
			inner = strings.TrimSpace(inner)
			if outer != "" && inner != "" {
				a = append(a, str.Join(outer, " ", inner))
			} else if outer == "" {
				a = append(a, inner)
			} else if inner == "" {
				a = append(a, outer)
			}
		}
	}
	return strings.Join(a, ", ")
}

func getBlockNodes(m map[string]any) []map[string]any {
	selectors, _ := maps.GetStrAnyVal[map[string]any](m, "blocks_metadata")
	mm, _ := maps.GetStrAnyVal[map[string]any](m, "theme_json.styles.blocks")
	var nodes []map[string]any
	for k, v := range mm {
		vv, ok := v.(map[string]any)
		if !ok {
			continue
		}
		s, _ := maps.GetStrAnyVal[string](selectors, str.Join(k, ".supports.selector"))
		d, _ := maps.GetStrAnyVal[string](selectors, str.Join(k, ".duotone"))
		f, _ := maps.GetStrAnyVal[string](selectors, str.Join(k, ".features"))
		n, ok := maps.GetStrAnyVal[map[string]any](vv, "variations")
		var variationSelectors []node
		if ok {
			for variation := range n {
				ss, _ := maps.GetStrAnyVal[string](selectors, str.Join(k, ".styleVariations.", variation))
				variationSelectors = append(variationSelectors, node{
					Path:     []string{"styles", "blocks", k, "variations", variation},
					Selector: ss,
				})
			}
		}
		nodes = append(nodes, map[string]any{
			"name":       k,
			"path":       []string{"styles", "blocks", k},
			"selector":   s,
			"duotone":    d,
			"features":   f,
			"variations": variationSelectors,
		})
		e, ok := maps.GetStrAnyVal[map[string]any](vv, "elements")
		if !ok {
			continue
		}
		for element, vvv := range e {
			_, ok = vvv.(map[string]any)
			if !ok {
				continue
			}
			key := str.Join(k, ".elements.", element)
			selector, _ := maps.GetStrAnyVal[string](selectors, key)
			nodes = append(nodes, map[string]any{
				"path":     []string{"styles", "blocks", k, "elements", element},
				"selector": selector,
			})
			if val, ok := __validElementPseudoSelectors[element]; ok {
				for _, ss := range val {
					_, ok = maps.GetStrAnyVal[string](vv, str.Join("elements.", ss))
					if !ok {
						continue
					}
					nodes = append(nodes, map[string]any{
						"path":     []string{"styles", "blocks", k, "elements", element},
						"selector": appendToSelector(selector, ss, ""),
					})
				}
			}
		}
	}
	return nodes
}

func getSettingNodes(m, setting map[string]any) []node {
	var nodes []node
	nodes = append(nodes, node{
		Path:     []string{"settings"},
		Selector: RootBlockSelector,
	})
	selectors, _ := maps.GetStrAnyVal[map[string]any](m, "blocks_metadata")
	s, ok := maps.GetStrAnyVal[map[string]any](setting, "settings.blocks")
	if !ok {
		return nil
	}
	for name := range s {
		selector, ok := maps.GetStrAnyVal[string](selectors, str.Join(name, ".supports.selector"))
		if ok {
			nodes = append(nodes, node{
				Path:     []string{"settings", "blocks", name},
				Selector: selector,
			})
		}
	}
	return nodes
}

type ThemeJson struct {
	blocksMetaData map[string]any
	themeJson      map[string]any
}

var validOrigins = []string{"default", "blocks", "theme", "custom"}

var layoutSelectorReg = regexp.MustCompile(`^[a-zA-Z0-9\-. *+>:()]*$`)

var validDisplayModes = []string{"block", "flex", "grid"}

func (j ThemeJson) getLayoutStyles(nodes node) string {
	//todo current theme supports disable-layout-styles
	var blockType map[string]any
	var s strings.Builder
	if nodes.Name != "" {
		v, ok := maps.GetStrAnyVal[map[string]any](j.blocksMetaData, nodes.Name)
		if ok {
			vv, ok := maps.GetStrAnyVal[map[string]any](v, "supports.__experimentalLayout")
			if ok && vv == nil {
				return ""
			}
			blockType = vv
		}
	}
	gap, hasBlockGapSupport := maps.GetStrAnyVal[map[string]any](j.themeJson, "settings.spacing.blockGap")
	if gap == nil {
		hasBlockGapSupport = false
	}
	_, ok := maps.GetStrAnyVal[map[string]any](j.themeJson, strings.Join(nodes.Path, "."))
	if !ok {
		return ""
	}
	layoutDefinitions, ok := maps.GetStrAnyVal[map[string]any](j.themeJson, "settings.layout.definitions")
	if !ok {
		return ""
	}
	blockGapValue := ""
	if !hasBlockGapSupport {
		if RootBlockSelector == nodes.Selector {
			blockGapValue = "0.5em"
		}
		if blockType != nil {
			blockGapValue, _ = maps.GetStrAnyVal[string](blockType, "supports.spacing.blockGap.__experimentalDefault")
		}
	} else {
		//todo getPropertyValue()
	}
	if blockGapValue != "" {
		for key, v := range layoutDefinitions {
			definition, ok := v.(map[string]any)
			if !ok {
				continue
			}
			if !hasBlockGapSupport && "flex" != key {
				continue
			}
			className := maps.GetStrAnyValWithDefaults(definition, "className", "")
			spacingRules := maps.GetStrAnyValWithDefaults(definition, "spacingStyles", []any{})
			if className == "" || spacingRules == nil {
				continue
			}
			for _, rule := range spacingRules {
				var declarations []declaration
				spacingRule, ok := rule.(map[string]any)
				if !ok {
					continue
				}
				selector := maps.GetStrAnyValWithDefaults(spacingRule, "selector", "")
				rules := maps.GetStrAnyValWithDefaults(spacingRule, "rules", map[string]any(nil))
				if selector != "" && !layoutSelectorReg.MatchString(selector) || rules == nil {
					continue
				}
				for property, v := range rules {
					value, ok := v.(string)
					if !ok || value == "" {
						value = blockGapValue
					}
					if isSafeCssDeclaration(property, value) {
						declarations = append(declarations, declaration{property, value})
					}
				}
				format := ""
				if !hasBlockGapSupport {
					format = helper.Or(RootBlockSelector == nodes.Selector, ":where(.%[2]s%[3]s)", ":where(%[1]s.%[2]s%[3]s)")
				} else {
					format = helper.Or(RootBlockSelector == nodes.Selector, "%s .%s%s", "%s.%s%s")
				}
				layoutSelector := fmt.Sprintf(format, nodes.Selector, className, selector)
				s.WriteString(toRuleset(layoutSelector, declarations))
			}
		}
	}

	if RootBlockSelector == nodes.Selector {
		for _, v := range layoutDefinitions {
			definition, ok := v.(map[string]any)
			if !ok {
				continue
			}
			className := maps.GetStrAnyValWithDefaults(definition, "className", "")
			baseStyleRules := maps.GetStrAnyValWithDefaults(definition, "baseStyles", []any{})
			if className == "" || nil == baseStyleRules {
				continue
			}
			displayMode := maps.GetStrAnyValWithDefaults(definition, "displayMode", "")
			if displayMode != "" && slice.IsContained(validDisplayModes, displayMode) {
				layoutSelector := str.Join(nodes.Selector, " .", className)
				s.WriteString(toRuleset(layoutSelector, []declaration{{"display", displayMode}}))
			}
			for _, rule := range baseStyleRules {
				var declarations []declaration
				r, ok := rule.(map[string]any)
				if !ok {
					continue
				}
				selector := maps.GetStrAnyValWithDefaults(r, "selector", "")
				rules := maps.GetStrAnyValWithDefaults(r, "rules", map[string]any(nil))
				if selector != "" && !layoutSelectorReg.MatchString(selector) || rules == nil {
					continue
				}
				for property, value := range rules {
					val, ok := value.(string)
					if !ok || val == "" {
						continue
					}
					if isSafeCssDeclaration(property, val) {
						declarations = append(declarations, declaration{property, val})
					}
				}
				layoutSelector := str.Join(nodes.Selector, " .", className, selector)
				s.WriteString(toRuleset(layoutSelector, declarations))
			}
		}
	}

	return s.String()
}

func isSafeCssDeclaration(name, value string) bool {
	css := str.Join(name, ": ", value)
	s := safeCSSFilterAttr(css)
	return "" != html.EscapeString(s)
}

var kses = regexp.MustCompile(`[\x00-\x08\x0B\x0C\x0E-\x1F]`)
var ksesop = regexp.MustCompile(`\\\\+0+`)

func KsesNoNull(s string, op ...bool) string {
	s = kses.ReplaceAllString(s, "")
	ops := true
	if len(op) > 0 {
		ops = op[0]
	}
	if ops {
		s = ksesop.ReplaceAllString(s, "")
	}
	return s
}

var allowedProtocols = []string{
	"http", "https", "ftp", "ftps", "mailto", "news", "irc", "irc6", "ircs", "gopher", "nntp",
	"feed", "telnet", "mms", "rtsp", "sms", "svn", "tel", "fax", "xmpp", "webcal", "urn",
}

var allowCssAttr = []string{
	"background", "background-color", "background-image", "background-position", "background-size", "background-attachment", "background-blend-mode",
	"border", "border-radius", "border-width", "border-color", "border-style", "border-right", "border-right-color",
	"border-right-style", "border-right-width", "border-bottom", "border-bottom-color", "border-bottom-left-radius",
	"border-bottom-right-radius", "border-bottom-style", "border-bottom-width", "border-bottom-right-radius",
	"border-bottom-left-radius", "border-left", "border-left-color", "border-left-style", "border-left-width",
	"border-top", "border-top-color", "border-top-left-radius", "border-top-right-radius", "border-top-style",
	"border-top-width", "border-top-left-radius", "border-top-right-radius", "border-spacing", "border-collapse",
	"caption-side", "columns", "column-count", "column-fill", "column-gap", "column-rule", "column-span", "column-width",
	"color", "filter", "font", "font-family", "font-size", "font-style", "font-variant", "font-weight",
	"letter-spacing", "line-height", "text-align", "text-decoration", "text-indent", "text-transform", "height",
	"min-height", "max-height", "width", "min-width", "max-width", "margin", "margin-right", "margin-bottom",
	"margin-left", "margin-top", "margin-block-start", "margin-block-end", "margin-inline-start", "margin-inline-end",
	"padding", "padding-right", "padding-bottom", "padding-left", "padding-top", "padding-block-start",
	"padding-block-end", "padding-inline-start", "padding-inline-end", "flex", "flex-basis", "flex-direction",
	"flex-flow", "flex-grow", "flex-shrink", "flex-wrap", "gap", "column-gap", "row-gap", "grid-template-columns",
	"grid-auto-columns", "grid-column-start", "grid-column-end", "grid-column-gap", "grid-template-rows",
	"grid-auto-rows", "grid-row-start", "grid-row-end", "grid-row-gap", "grid-gap", "justify-content", "justify-items",
	"justify-self", "align-content", "align-items", "align-self", "clear", "cursor", "direction", "float",
	"list-style-type", "object-fit", "object-position", "overflow", "vertical-align", "position", "top", "right",
	"bottom", "left", "z-index", "aspect-ratio", "--*",
}
var allowCssAttrMap = slice.FilterAndToMap(allowCssAttr, func(t string) (string, struct{}, bool) {
	return t, struct{}{}, true
})

var __strReg = regexp.MustCompile(`^--[a-zA-Z0-9-_]+$`)
var cssUrlDataTypes = []string{"background", "background-image", "cursor", "list-style", "list-style-image"}
var cssGradientDataTypes = []string{"background", "background-image"}
var __allowCSSReg = regexp2.MustCompile(`\b(?:var|calc|min|max|minmax|clamp)(\((?:[^()]|(?!1))*\))`, regexp2.None)
var __disallowCSSReg = regexp.MustCompile(`[\\(&=}]|/\*`)

func safeCSSFilterAttr(css string) string {
	css = KsesNoNull(css)
	css = strings.TrimSpace(css)
	css = str.Replaces(css, []string{"\n", "\r", "\t", ""})
	cssArr := strings.Split(css, ";")
	var isCustomVar, found, urlAttr, gradientAttr bool
	var cssValue string
	var ss strings.Builder
	for _, s := range cssArr {
		if s == "" {
			continue
		}
		if !strings.Contains(s, ":") {
			found = true
		} else {
			parts := strings.SplitN(s, ":", 2)
			selector := strings.TrimSpace(parts[0])
			if maps.IsExists(allowCssAttrMap, "--*") {
				i := __strReg.FindStringIndex(selector)
				if len(i) > 0 && i[0] > 0 {
					isCustomVar = true
					allowCssAttr = append(allowCssAttr, selector)
				}
			}
			if maps.IsExists(allowCssAttrMap, selector) {
				found = true
				urlAttr = slice.IsContained(cssUrlDataTypes, selector)
				gradientAttr = slice.IsContained(cssGradientDataTypes, selector)
			}
			if isCustomVar {
				cssValue = strings.TrimSpace(parts[1])
				urlAttr = cssValue[0:4] == "url("
				gradientAttr = strings.Contains(cssValue, "-gradient(")
			}
		}

		if found {
			//todo wtf 🤮
			if urlAttr {

			}
			if gradientAttr {

			}
			cssTest, _ := __allowCSSReg.Replace(s, "", 0, -1)
			isAllow := !__disallowCSSReg.MatchString(cssTest)
			if isAllow {
				ss.WriteString(s)
				ss.WriteString(";")
			}
		}

	}
	return strings.TrimRight(ss.String(), ";")
}

func (j ThemeJson) getStyletSheet(types, origins []string, options map[string]string) string {
	if origins == nil {
		origins = append(validOrigins)
	}
	styleNodes := getStyleNodes(j)
	settingsNodes := getSettingNodes(j.blocksMetaData, j.themeJson)
	rootStyleKey, _ := slice.SearchFirst(styleNodes, func(n node) bool {
		return n.Selector == RootBlockSelector
	})
	rootSettingsKey, _ := slice.SearchFirst(settingsNodes, func(n node) bool {
		return n.Selector == RootBlockSelector
	})
	if os, ok := options["scope"]; ok {
		for i := range settingsNodes {
			settingsNodes[i].Selector = scopeSelector(os, settingsNodes[i].Selector)
		}
		for i := range styleNodes {
			styleNodes[i].Selector = scopeSelector(os, styleNodes[i].Selector)
		}
	}
	if or, ok := options["root_selector"]; ok && or != "" {
		if rootSettingsKey > -1 {
			settingsNodes[rootSettingsKey].Selector = or
		}
		if rootStyleKey > -1 && rootStyleKey < len(settingsNodes) {
			settingsNodes[rootStyleKey].Selector = or
		}
	}
	stylesSheet := ""
	if slice.IsContained(types, "variables") {
		stylesSheet = j.getCssVariables(settingsNodes, origins)
	}

	if slice.IsContained(types, "styles") {
		stylesSheet = j.getCssVariables(settingsNodes, origins)
		if rootStyleKey > -1 {
			//todo getRootLayoutRules
			stylesSheet = str.Join(stylesSheet)
		}
	} else if slice.IsContained(types, "base-layout-styles") {
		rootSelector := RootBlockSelector
		columnsSelector := ".wp-block-columns"
		scope, ok := options["scope"]
		if ok && scope != "" {
			rootSelector = scopeSelector(scope, rootSelector)
			columnsSelector = scopeSelector(scope, columnsSelector)
		}
		rs, ok := options["root_selector"]
		if ok && rs != "" {
			rootSelector = rs
		}
		baseStylesNodes := []node{
			{
				Path:     []string{"styles"},
				Selector: rootSelector,
			},
			{
				Path:     []string{"styles", "blocks", "core/columns"},
				Selector: columnsSelector,
				Name:     "core/columns",
			},
		}
		for _, stylesNode := range baseStylesNodes {
			stylesSheet = str.Join(stylesSheet, j.getLayoutStyles(stylesNode))
		}
	}

	if slice.IsContained(types, "presets") {
		stylesSheet = str.Join(stylesSheet, j.getPresetClasses(settingsNodes, origins))
	}

	return stylesSheet
}

func (j ThemeJson) getPresetClasses(nodes []node, origins []string) string {
	var presetRules strings.Builder
	for _, n := range nodes {
		if n.Selector == "" {
			continue
		}
		no, ok := maps.GetStrAnyVal[map[string]any](j.themeJson, strings.Join(n.Path, "."))
		if !ok {
			continue
		}
		presetRules.WriteString(computePresetClasses(no, n.Selector, origins))
	}
	return presetRules.String()
}

func computePresetClasses(m map[string]any, selector string, origins []string) string {
	if selector == RootBlockSelector {
		selector = ""
	}
	var s strings.Builder
	for _, meta := range presetsMetadata {
		slugs := getSettingsSlugs(m, meta, origins)
		for class, property := range meta.classes {
			for _, slug := range slugs {
				cssVar := strings.ReplaceAll(meta.cssVars, "$slug", slug)
				className := strings.ReplaceAll(class, "$slug", slug)
				s.WriteString(toRuleset(appendToSelector(selector, className, ""), []declaration{
					{property, str.Join("var(", cssVar, ") !important")},
				}))
			}
		}
	}
	return s.String()
}

func getSettingsSlugs(settings map[string]any, meta presetMeta, origins []string) map[string]string {
	if origins == nil {
		origins = validOrigins
	}

	presetPerOrigin, ok := maps.GetStrAnyVal[map[string]any](settings, strings.Join(meta.path, "."))
	if !ok {
		return nil
	}
	m := map[string]string{}
	for _, origin := range origins {
		o, ok := maps.GetStrAnyVal[[]map[string]string](presetPerOrigin, origin)
		if !ok {
			continue
		}
		for _, mm := range o {
			slug := toKebabCase(mm["slug"])
			m[slug] = slug
		}
	}
	return m
}

func toKebabCase(s string) string {
	s = strings.ReplaceAll(s, "'", "")
	r, err := __kebabCaseReg.FindStringMatch(s)
	if err != nil {
		return s
	}
	var ss []string
	for r != nil {
		if r.GroupCount() < 1 {
			break
		}

		ss = append(ss, r.Groups()[0].String())
		r, _ = __kebabCaseReg.FindNextMatch(r)
	}

	return strings.ToLower(strings.Join(ss, "-"))
}

var __kebabCaseReg = func() *regexp2.Regexp {
	rsLowerRange := "a-z\\xdf-\\xf6\\xf8-\\xff"
	rsNonCharRange := "\\x00-\\x2f\\x3a-\\x40\\x5b-\\x60\\x7b-\\xbf"
	rsPunctuationRange := "\\x{2000}-\\x{206f}"
	rsSpaceRange := " \\t\\x0b\\f\\xa0\\x{feff}\\n\\r\\x{2028}\\x{2029}\\x{1680}\\x{180e}\\x{2000}\\x{2001}\\x{2002}\\x{2003}\\x{2004}\\x{2005}\\x{2006}\\x{2007}\\x{2008}\\x{2009}\\x{200a}\\x{202f}\\x{205f}\\x{3000}"
	rsUpperRange := "A-Z\\xc0-\\xd6\\xd8-\\xde"
	rsBreakRange := rsNonCharRange + rsPunctuationRange + rsSpaceRange

	/** Used to compose unicode capture groups. */
	rsBreak := "[" + rsBreakRange + "]"
	rsDigits := "\\d+" // The last lodash version in GitHub uses a single digit here and expands it when in use.
	rsLower := "[" + rsLowerRange + "]"
	rsMisc := "[^" + rsBreakRange + rsDigits + rsLowerRange + rsUpperRange + "]"
	rsUpper := "[" + rsUpperRange + "]"

	/** Used to compose unicode regexes. */
	rsMiscLower := "(?:" + rsLower + "|" + rsMisc + ")"
	rsMiscUpper := "(?:" + rsUpper + "|" + rsMisc + ")"
	rsOrdLower := "\\d*(?:1st|2nd|3rd|(?![123])\\dth)(?=\\b|[A-Z_])"
	rsOrdUpper := "\\d*(?:1ST|2ND|3RD|(?![123])\\dTH)(?=\\b|[a-z_])"

	reg := strings.Join([]string{
		rsUpper + "?" + rsLower + "+(?=" + strings.Join([]string{rsBreak, rsUpper, "$"}, "|") + ")",
		rsMiscUpper + "+(?=" + strings.Join([]string{rsBreak, rsUpper + rsMiscLower, "$"}, "|") + ")",
		rsUpper + "?" + rsMiscLower + "+",
		rsUpper + "+",
		rsOrdUpper,
		rsOrdLower,
		rsDigits,
	}, "|")
	return regexp2.MustCompile(reg, regexp2.Unicode)
}()

var presetsMetadata = []presetMeta{
	{
		path:            []string{"color", "palette"},
		preventOverride: []string{"color", "defaultPalette"},
		useDefaultNames: false,
		valueKey:        "color",
		valueFunc:       nil,
		cssVars:         "--wp--preset--color--$slug",
		classes: map[string]string{
			".has-$slug-color":            "color",
			".has-$slug-background-color": "background-color",
			".has-$slug-border-color":     "border-color",
		},
		properties: []string{"color", "background-color", "border-color"},
	}, {
		path:            []string{"color", "gradients"},
		preventOverride: []string{"color", "defaultGradients"},
		useDefaultNames: false,
		valueKey:        "gradient",
		valueFunc:       nil,
		cssVars:         "--wp--preset--gradient--$slug",
		classes: map[string]string{
			".has-$slug-gradient-background": "background",
		},
		properties: []string{"background"},
	}, {
		path:            []string{"color", "duotone"},
		preventOverride: []string{"color", "defaultDuotone"},
		useDefaultNames: false,
		valueKey:        "",
		valueFunc:       wpGetDuotoneFilterProperty,
		cssVars:         "--wp--preset--duotone--$slug",
		classes:         map[string]string{},
		properties:      []string{"filter"},
	}, {
		path:            []string{"typography", "fontSizes"},
		preventOverride: []string{},
		useDefaultNames: true,
		valueKey:        "",
		valueFunc:       wpGetTypographyFontSizeValue,
		cssVars:         "--wp--preset--font-size--$slug",
		classes: map[string]string{
			".has-$slug-font-size": "font-size",
		},
		properties: []string{"font-size"},
	}, {
		path:            []string{"typography", "fontFamilies"},
		preventOverride: []string{},
		useDefaultNames: false,
		valueKey:        "fontFamily",
		valueFunc:       nil,
		cssVars:         "--wp--preset--font-family--$slug",
		classes: map[string]string{
			".has-$slug-font-family": "font-family",
		},
		properties: []string{"font-family"},
	}, {
		path:            []string{"spacing", "spacingSizes"},
		preventOverride: []string{},
		useDefaultNames: true,
		valueKey:        "size",
		valueFunc:       nil,
		cssVars:         "--wp--preset--spacing--$slug",
		classes:         map[string]string{},
		properties:      []string{"padding", "margin"},
	}, {
		path:            []string{"shadow", "presets"},
		preventOverride: []string{"shadow", "defaultPresets"},
		useDefaultNames: false,
		valueKey:        "shadow",
		valueFunc:       nil,
		cssVars:         "--wp--preset--shadow--$slug",
		classes:         map[string]string{},
		properties:      []string{"box-shadow"},
	},
}

type presetMeta struct {
	path            []string
	preventOverride []string
	useDefaultNames bool
	valueFunc       func(map[string]string, map[string]any) string
	valueKey        string
	cssVars         string
	classes         map[string]string
	properties      []string
}

type declaration struct {
	name  string
	value string
}

func wpGetDuotoneFilterProperty(preset map[string]string, _ map[string]any) string {
	v, ok := preset["colors"]
	if ok && "unset" == v {
		return "none"
	}
	id, ok := preset["slug"]
	if ok {
		id = str.Join("wp-duotone-", id)
	}
	return str.Join(`url('#`, id, "')")
}

func wpGetTypographyFontSizeValue(preset map[string]string, m map[string]any) string {
	size, ok := preset["size"]
	if !ok {
		return ""
	}
	if size == "" || size == "0" {
		return size
	}
	origin := "custom"
	if !wpconfig.HasThemeJson() {
		origin = "theme"
	}
	typographySettings, ok := maps.GetStrAnyVal[map[string]any](m, "typography")
	if !ok {
		return size
	}
	fluidSettings, ok := maps.GetStrAnyVal[map[string]any](typographySettings, "fluid")

	//todo  so complex dying 👻
	_ = fluidSettings
	_ = origin
	return size
}

var __themeJson = reload.Vars(ThemeJson{})

func GetThemeJson() ThemeJson {
	return __themeJson.Load()
}

func wpGetTypographyValueAndUnit(value string, options map[string]any) {
	/*options := maps.Merge(options, map[string]any{
		"coerce_to":        "",
		"root_size_value":  16,
		"acceptable_units": []string{"rem", "px", "em"},
	})
	u, _ := maps.GetStrAnyVal[[]string](options, "acceptable_units")
	acceptableUnitsGroup := strings.Join(u, "|")*/

}

func computeThemeVars(m map[string]any) []declaration {
	//todo ......
	return nil
}

func computePresetVars(m map[string]any, origins []string) []declaration {
	var declarations []declaration
	for _, metadatum := range presetsMetadata {
		slug := getSettingsValuesBySlug(m, metadatum, origins)
		for k, v := range slug {
			declarations = append(declarations, declaration{
				name:  strings.Replace(metadatum.cssVars, "$slug", k, -1),
				value: v,
			})
		}
	}
	return declarations
}

func getSettingsValuesBySlug(m map[string]any, meta presetMeta, origins []string) map[string]string {
	perOrigin := maps.GetStrAnyValWithDefaults[map[string]any](m, strings.Join(meta.path, "."), nil)
	r := map[string]string{}
	for _, origin := range origins {
		if vv, ok := maps.GetStrAnyVal[[]map[string]string](perOrigin, origin); ok {
			for _, preset := range vv {
				slug := preset["slug"]
				value := ""
				if vv, ok := preset[meta.valueKey]; ok && vv != "" {
					value = vv
				} else if meta.valueFunc != nil {
					value = meta.valueFunc(preset, m)
				}
				r[slug] = value
			}
		}
	}
	return r
}

func setSpacingSizes(t ThemeJson) {
	m, _ := maps.GetStrAnyVal[map[string]any](t.themeJson, "settings.spacing.spacingScale")
	unit, _ := maps.GetStrAnyVal[string](m, "unit")
	currentStep, _ := maps.GetStrAnyVal[float64](m, "mediumStep")
	mediumStep := currentStep
	steps, _ := maps.GetStrAnyVal[float64](m, "steps")
	operator, _ := maps.GetStrAnyVal[string](m, "operator")
	increment, _ := maps.GetStrAnyVal[float64](m, "increment")
	stepsMidPoint := math.Round(steps / 2)
	reminder := float64(0)
	xSmallCount := ""
	var sizes []map[string]string
	slug := 40
	for i := stepsMidPoint - 1; i > 0 && slug > 0 && steps > 1; i-- {
		if "+" == operator {
			currentStep -= increment
		} else if increment > 1 {
			currentStep /= increment
		} else {
			currentStep *= increment
		}
		if currentStep <= 0 {
			reminder = i
			break
		}
		name := "small"
		if i != stepsMidPoint-1 {
			name = str.Join(xSmallCount, "X-Small")
		}
		sizes = append(sizes, map[string]string{
			"name": name,
			"slug": number.IntToString(slug),
			"size": fmt.Sprintf("%v%s", number.Round(currentStep, 2), unit),
		})
		if i == stepsMidPoint-2 {
			xSmallCount = strconv.Itoa(2)
		}
		if i < stepsMidPoint-2 {
			n := str.ToInt[int](xSmallCount)
			n++
			xSmallCount = strconv.Itoa(n)
		}
		slug -= 10
	}
	slice.ReverseSelf(sizes)
	sizes = append(sizes, map[string]string{
		"name": "Medium",
		"slug": "50",
		"size": str.Join(number.ToString(mediumStep), unit),
	})
	currentStep = mediumStep
	slug = 60
	xLargeCount := ""
	stepsAbove := steps - stepsMidPoint + reminder
	for aboveMidpointCount := float64(0); aboveMidpointCount < stepsAbove; aboveMidpointCount++ {
		if "+" == operator {
			currentStep += increment
		} else if increment >= 1 {
			currentStep *= increment
		} else {
			currentStep /= increment
		}
		name := "Large"
		if 0 != aboveMidpointCount {
			name = str.Join(xLargeCount, "X-Large")
		}
		sizes = append(sizes, map[string]string{
			"name": name,
			"slug": strconv.Itoa(slug),
			"size": fmt.Sprintf("%v%s", number.Round(currentStep, 2), unit),
		})
		if aboveMidpointCount == 1 {
			xLargeCount = strconv.Itoa(2)
		}
		if aboveMidpointCount > 1 {
			x := str.ToInt[int](xLargeCount)
			x++
			xLargeCount = strconv.Itoa(x)
		}
		slug += 10
	}
	if steps <= 7 {
		for i := 0; i < len(sizes); i++ {
			sizes[i]["name"] = strconv.Itoa(i + 1)
		}
	}
	maps.SetStrAnyVal(t.themeJson, "settings.spacing.spacingSizes.default", sizes)
}

func (j ThemeJson) getCssVariables(settingNodes []node, origins []string) string {
	var s strings.Builder
	for _, settingNode := range settingNodes {
		if "" == settingNode.Selector {
			continue
		}
		n := maps.GetStrAnyValWithDefaults[map[string]any](j.themeJson, strings.Join(settingNode.Path, "."), nil)
		declarations := computePresetVars(n, origins)
		declarations = append(declarations, computeThemeVars(j.themeJson)...)
		s.WriteString(toRuleset(settingNode.Selector, declarations))
	}
	return s.String()
}

func toRuleset(selector string, declarations []declaration) string {
	if len(declarations) < 1 {
		return ""
	}
	s := slice.Reduce(declarations, func(t declaration, r string) string {
		return str.Join(r, t.name, ":", t.value, ";")
	}, "")
	return str.Join(selector, "{", s, "}")
}

func getStyleNodes(t ThemeJson) []node {
	var styleNodes = []node{
		{[]string{"styles"}, "body", ""},
	}
	m := maps.GetStrAnyValWithDefaults[map[string]any](t.themeJson, "styles.elements", nil)
	if len(m) < 1 {
		return nil
	}
	for e, s := range __elements {
		_, ok := m[e]
		if !ok {
			continue
		}
		styleNodes = append(styleNodes, node{[]string{"styles", "elements", e}, s, ""})
		ss, ok := __validElementPseudoSelectors[e]
		if ok {
			for _, sel := range ss {
				if maps.IsExists(__elements, sel) {
					styleNodes = append(styleNodes, node{
						Path:     []string{"styles", "elements", e},
						Selector: appendToSelector(s, sel, ""),
					})
				}
			}
		}
	}
	blocks, _ := maps.GetStrAnyVal[map[string]any](t.blocksMetaData, "theme_json.styles.blocks")
	maps.SetStrAnyVal(t.themeJson, "styles.blocks", blocks)
	blockNodes := getBlockNodes(t.blocksMetaData)
	for _, blockNode := range blockNodes {
		p, ok := maps.GetStrAnyVal[[]string](blockNode, "path")
		s, oo := maps.GetStrAnyVal[string](blockNode, "selector")
		if ok && oo {
			_ = append(styleNodes, node{Path: p, Selector: s})
			//styleNodes = append(styleNodes, node{Path: p, Selector: s})
		}
	}

	return styleNodes
}

func GetGlobalStyletSheet() string {
	t := __themeJson.Load()
	var types, origins []string
	if types == nil && !wpconfig.HasThemeJson() {
		types = []string{"variables", "presets", "base-layout-styles"}
	} else if types == nil {
		types = []string{"variables", "styles", "presets"}
	}
	styleSheet := ""
	if slice.IsContained(types, "variables") {
		origins = []string{"default", "theme", "custom"}
		styleSheet = t.getStyletSheet([]string{"variables"}, origins, nil)
		slice.Delete(&types, slice.IndexOf(types, "variables"))
	}

	if len(types) > 0 {
		origins = []string{"default", "theme", "custm"}
		if !wpconfig.HasThemeJson() {
			origins = []string{"default"}
		}
		styleSheet = str.Join(styleSheet, t.getStyletSheet(types, origins, nil))
	}

	return styleSheet
}

/*func (j ThemeJson) getStylesForBlock(blockMeta map[string]any) {
	path, _ := maps.GetStrAnyVal[[]string](blockMeta, "path")
	node, _ := maps.GetStrAnyVal[map[string]any](j.themeJson, strings.Join(path, "."))
	useRootPadding := maps.GetStrAnyValWithDefaults(j.themeJson, "settings.useRootPaddingAwareAlignments", false)
	settings, _ := maps.GetStrAnyVal(j.themeJson, "settings")
	is_processing_element := slice.IsContained(path, "elements")
	currentElement := ""
	if is_processing_element {
		currentElement = path[len(path)-1]
	}
	element_pseudo_allowed := __validElementPseudoSelectors[currentElement]

}

func computeStyleProperties(styles, settings, properties, themeJson map[string]any, selector string, useRootPadding bool) {
	if properties == nil {
		//properties =
	}
}
*/
