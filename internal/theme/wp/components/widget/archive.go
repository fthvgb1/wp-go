package widget

import (
	"fmt"
	"github.com/fthvgb1/wp-go/helper/maps"
	"github.com/fthvgb1/wp-go/helper/slice"
	str "github.com/fthvgb1/wp-go/helper/strings"
	"github.com/fthvgb1/wp-go/internal/pkg/cache"
	"github.com/fthvgb1/wp-go/internal/pkg/constraints/widgets"
	"github.com/fthvgb1/wp-go/internal/pkg/models"
	"github.com/fthvgb1/wp-go/internal/theme/wp"
	"github.com/fthvgb1/wp-go/internal/wpconfig"
	"github.com/fthvgb1/wp-go/safety"
	"strings"
)

var archivesConfig = func() safety.Var[map[any]any] {
	v := safety.Var[map[any]any]{}
	v.Store(map[any]any{
		"count":    int64(0),
		"dropdown": int64(0),
		"title":    "归档",
	})
	archiveArgs.Store(map[string]string{
		"{$before_widget}":  `<aside id="archives-2" class="widget widget_archive">`,
		"{$after_widget}":   "</aside>",
		"{$before_title}":   `<h2 class="widget-title">`,
		"{$after_title}":    "</h2>",
		"{$before_sidebar}": "",
		"{$after_sidebar}":  "",
		"{$nav}":            "",
		"{$navCloser}":      "",
		"{$title}":          "",
		"{$dropdown_id}":    "archives-dropdown-2",
		"{$dropdown_type}":  "monthly",
		"{$dropdown_label}": "选择月份",
	})
	return v
}()

var archiveArgs = safety.Var[map[string]string]{}

var archiveTemplate = `{$before_widget}
{$title}
{$nav}
{$html}
{$navCloser}
{$after_widget}
`

func Archive(h *wp.Handle) string {
	args := wp.GetComponentsArgs(h, widgets.ArchiveArgs, archiveArgs.Load())
	args = maps.FilterZeroMerge(archiveArgs.Load(), args)
	conf := wpconfig.GetPHPArrayVal("widget_archives", archivesConfig.Load(), int64(2))
	args["{$title}"] = str.Join(args["{$before_title}"], conf["title"].(string), args["{$after_title}"])
	s := archiveTemplate
	if int64(1) == conf["dropdown"].(int64) {
		s = strings.ReplaceAll(s, "{$html}", archiveDropDown(h, conf, args, cache.Archives(h.C)))
	} else {
		s = strings.ReplaceAll(s, "{$html}", archiveUl(h, conf, args, cache.Archives(h.C)))
	}
	return h.ComponentFilterFnHook(widgets.ArchiveArgs, str.Replace(s, args))
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
	month := strings.TrimLeft(h.Index.Param.Month, "0")
	showCount := conf["count"].(int64)
	for _, archive := range archives {
		sel := ""
		if h.Index.Param.Year == archive.Year && month == archive.Month {
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
	if slice.IsContained(h.CommonThemeMods().ThemeSupport.HTML5, "navigation-widgets") {
		args["{$nav}"] = fmt.Sprintf(`<nav aria-label="%s">`, conf["title"].(string))
		args["{$navCloser}"] = "</nav>"
	}
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
