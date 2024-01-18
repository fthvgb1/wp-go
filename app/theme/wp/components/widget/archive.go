package widget

import (
	"fmt"
	"github.com/fthvgb1/wp-go/app/pkg/cache"
	"github.com/fthvgb1/wp-go/app/pkg/constraints/widgets"
	"github.com/fthvgb1/wp-go/app/pkg/models"
	"github.com/fthvgb1/wp-go/app/theme/wp"
	"github.com/fthvgb1/wp-go/cache/reload"
	"github.com/fthvgb1/wp-go/helper/maps"
	"github.com/fthvgb1/wp-go/helper/slice"
	str "github.com/fthvgb1/wp-go/helper/strings"
	"strings"
)

var archiveTemplate = `{$before_widget}
{$title}
{$nav}
{$html}
{$navCloser}
{$after_widget}
`

func archiveArgs() map[string]string {
	return map[string]string{
		"{$before_sidebar}": "",
		"{$after_sidebar}":  "",
		"{$nav}":            "",
		"{$navCloser}":      "",
		"{$title}":          "",
		"{$dropdown_id}":    "archives-dropdown-2",
		"{$dropdown_type}":  "monthly",
		"{$dropdown_label}": "选择月份",
	}
}

var archivesConfig = map[any]any{
	"count":    int64(0),
	"dropdown": int64(0),
	"title":    "归档",
}

var GetArchiveConf = BuildconfigFn(archivesConfig, "widget_archives", int64(2))
var GetArchiveArgs = reload.BuildValFnWithAnyParams("", archiveArgsFn)

func archiveArgsFn(a ...any) map[string]string {
	h := a[0].(*wp.Handle)
	conf := a[1].(map[any]any)
	id := a[2].(string)
	archiveArgs := archiveArgs()
	commonArgs := wp.GetComponentsArgs(widgets.Widget, CommonArgs())
	args := wp.GetComponentsArgs(widgets.Archive, archiveArgs)
	args = maps.FilterZeroMerge(archiveArgs, CommonArgs(), commonArgs, args)
	args["{$before_widget}"] = fmt.Sprintf(args["{$before_widget}"], str.Join("archives-", id), str.Join("widget widget_", "archive"))
	args["{$title}"] = str.Join(args["{$before_title}"], conf["title"].(string), args["{$after_title}"])
	if conf["dropdown"].(int64) == 0 && slice.IsContained(h.CommonThemeMods().ThemeSupport.HTML5, "navigation-widgets") {
		args["{$nav}"] = fmt.Sprintf(`<nav aria-label="%s">`, conf["title"].(string))
		args["{$navCloser}"] = "</nav>"
	}
	return args
}

func Archive(h *wp.Handle, id string) string {
	conf := GetArchiveConf()
	args := GetArchiveArgs(h, conf, id)

	s := archiveTemplate
	if int64(1) == conf["dropdown"].(int64) {
		s = strings.ReplaceAll(s, "{$html}", archiveDropDown(h, conf, args, cache.Archives(h.C)))
	} else {
		s = strings.ReplaceAll(s, "{$html}", archiveUl(h, conf, args, cache.Archives(h.C)))
	}
	return h.DoActionFilter(widgets.Archive, str.Replace(s, args))
}

var dropdownScript = `
<script>
            /* <![CDATA[ */
            (function() {
                const dropdown = document.getElementById("archives-dropdown-2");
                function onSelectChange() {
                    if ( dropdown.options[ dropdown.selectedIndex ].value !== '' ) {
                        document.location.href = this.options[ this.selectedIndex ].value;
                    }
                }
                dropdown.onchange = onSelectChange;
            })();
            /* ]]> */
        </script>`

func archiveDropDown(h *wp.Handle, conf map[any]any, args map[string]string, archives []models.PostArchive) string {
	option := str.NewBuilder()
	option.Sprintf(`<option value="">%s</option>`, args["{$dropdown_label}"])
	i := h.GetIndexHandle()
	month := strings.TrimLeft(i.Param.Month, "0")
	showCount := conf["count"].(int64)
	for _, archive := range archives {
		sel := ""
		if i.Param.Year == archive.Year && month == archive.Month {
			sel = "selected"
		}
		count := ""
		if showCount == int64(1) {
			count = fmt.Sprintf("(%v)", archive.Posts)
		}
		option.Sprintf(`<option %s value="/p/date/%s/%02s" >%s年%s月 %s</option>
`, sel, archive.Year, archive.Month, archive.Year, archive.Month, count)
	}
	s := str.NewBuilder()
	s.Sprintf(`<label class="screen-reader-text" for="%s">%s</label>
<select id="%s" name="archive-dropdown">%s</select>%s
`, args["{$dropdown_id}"], args["{$title}"], args["{$dropdown_id}"], option.String(), dropdownScript)
	return s.String()
}

func archiveUl(h *wp.Handle, conf map[any]any, args map[string]string, archives []models.PostArchive) string {
	s := str.NewBuilder()
	s.WriteString(`<ul>`)
	showCount := conf["count"].(int64)
	for _, archive := range archives {
		count := ""
		if showCount == 1 {
			count = fmt.Sprintf("(%v)", archive.Posts)
		}
		s.Sprintf(`<li><a href="/p/date/%[1]s/%02[2]s">%[1]s年%[2]s月%[3]s</a></li>`, archive.Year, archive.Month, count)
	}
	return s.String()
}
