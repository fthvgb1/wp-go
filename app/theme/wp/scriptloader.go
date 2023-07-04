package wp

import (
	"encoding/json"
	"fmt"
	"github.com/fthvgb1/wp-go/app/cmd/reload"
	"github.com/fthvgb1/wp-go/app/wpconfig"
	"github.com/fthvgb1/wp-go/helper"
	"github.com/fthvgb1/wp-go/helper/maps"
	"github.com/fthvgb1/wp-go/helper/number"
	"github.com/fthvgb1/wp-go/helper/slice"
	str "github.com/fthvgb1/wp-go/helper/strings"
	"github.com/fthvgb1/wp-go/safety"
	"math"
	"strconv"
	"strings"
)

var scripts = reload.MapBy[string, *Script](func(m *safety.Map[string, *Script]) {
	suffix := ".min"
	defaultScripts(m, suffix)
})

func addScript(handle string, src string, deps []string, ver string, args any) {
	script := NewScript(handle, src, deps, ver, args)
	scripts.Store(handle, script)
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
	addData(handle, "data", localize(handle, objectname, l10n))
}
func AddDynamicLocalize(h *Handle, handle, objectname string, l10n map[string]any) {
	AddDynamicData(h, handle, "data", localize(handle, objectname, l10n))
}

func getData(handle, key string) string {
	h, ok := scripts.Load(handle)
	if !ok {
		return ""
	}
	return strings.Join(h.Extra[key], "\n")
}
func GetData(h *Handle, handle, key string) string {
	hh, ok := scripts.Load(handle)
	if !ok {
		return ""
	}
	d := hh.Extra[key]
	d = append(d, GetDynamicData(h, handle, key))
	return strings.Join(d, "\n")
}

func addData(handle, key, data string) {
	s, ok := scripts.Load(handle)
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
	addData(handle, position, data)
}

func AddInlineStyle(handle, data string) {
	if handle == "" || data == "" {
		return
	}
	addData(handle, "after", data)
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
	scripts.Store(handle, script)
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

func AddDynamicData(h *Handle, handle, key, data string) {
	da := helper.GetContextVal(h.C, "__scriptDynamicData__", map[string]map[string][]string{})
	m, ok := da[handle]
	if !ok {
		m = map[string][]string{}
	}
	m[key] = append(m[key], data)
	da[handle] = m
}

func GetDynamicData(h *Handle, handle, key string) string {
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
	hh, ok := scripts.Load(handle)
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
		outer = strings.Trim(outer, " \t\n\r\x00\x0B")
		for _, inner := range selectors {
			inner = strings.Trim(inner, " \t\n\r\x00\x0B")
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

var __themeJson *safety.Var[ThemeJson]

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
		{[]string{"styles"}, "body"},
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
		styleNodes = append(styleNodes, node{[]string{"styles", "elements", e}, s})
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
			styleNodes = append(styleNodes, node{Path: p, Selector: s})
		}
	}

	return styleNodes
}

var validOrigins = []string{"default", "theme", "custom"}

func GetStyletSheet(t ThemeJson, types, origins []string, options map[string]string) string {
	styleNodes := getStyleNodes(t)
	settingsNodes := getSettingNodes(t.blocksMetaData, t.themeJson)
	if types == nil && !wpconfig.HasThemeJson() {
		types = []string{"variables", "presets", "base-layout-styles"}
	} else if types == nil {
		types = []string{"variables", "styles", "presets"}
	}
	if origins == nil {
		origins = validOrigins
	}
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
	var s string
	if slice.IsContained(types, "variables") {
		s = t.getCssVariables(settingsNodes, origins)
		slice.Delete(&types, slice.IndexOf(types, "variables"))
	}
	if slice.IsContained(types, "styles") {
		s = t.getCssVariables(settingsNodes, origins)
	}

	return s
}
