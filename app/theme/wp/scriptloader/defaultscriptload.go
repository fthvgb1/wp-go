package scriptloader

import (
	"encoding/json"
	"github.com/fthvgb1/wp-go/app/pkg/config"
	"github.com/fthvgb1/wp-go/app/pkg/logs"
	"github.com/fthvgb1/wp-go/helper/maps"
	str "github.com/fthvgb1/wp-go/helper/strings"
	"github.com/fthvgb1/wp-go/safety"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func defaultScripts(m *safety.Map[string, *Script], suffix string) {
	m.Store("utils", NewScript("utils", "/wp-includes/js/utils"+suffix+".js", nil, "", nil))
	m.Store("common", NewScript("common", "/wp-admin/js/common"+suffix+".js", []string{"jquery", "hoverIntent", "utils", "wp-i18n"}, "", 1))
	m.Store("wp-sanitize", NewScript("wp-sanitize", "/wp-includes/js/wp-sanitize"+suffix+".js", nil, "", 1))
	m.Store("sack", NewScript("sack", "/wp-includes/js/tw-sack"+suffix+".js", nil, "1.6.1", 1))
	m.Store("quicktags", NewScript("quicktags", "/wp-includes/js/quicktags"+suffix+".js", nil, "", 1))
	m.Store("colorpicker", NewScript("colorpicker", "/wp-includes/js/colorpicker"+suffix+".js", []string{"prototype"}, "3517m", nil))
	m.Store("editor", NewScript("editor", "/wp-admin/js/editor"+suffix+".js", []string{"utils", "jquery"}, "", 1))
	m.Store("clipboard", NewScript("clipboard", "/wp-includes/js/clipboard"+suffix+".js", nil, "2.0.11", 1))
	m.Store("wp-ajax-response", NewScript("wp-ajax-response", "/wp-includes/js/wp-ajax-response"+suffix+".js", []string{"jquery", "wp-a11y"}, "", 1))
	m.Store("wp-api-request", NewScript("wp-api-request", "/wp-includes/js/api-request"+suffix+".js", []string{"jquery"}, "", 1))
	m.Store("wp-pointer", NewScript("wp-pointer", "/wp-includes/js/wp-pointer"+suffix+".js", []string{"jquery-ui-core", "wp-i18n"}, "", 1))
	m.Store("autosave", NewScript("autosave", "/wp-includes/js/autosave"+suffix+".js", []string{"heartbeat"}, "", 1))
	m.Store("heartbeat", NewScript("heartbeat", "/wp-includes/js/heartbeat"+suffix+".js", []string{"jquery", "wp-hooks"}, "", 1))
	m.Store("wp-auth-check", NewScript("wp-auth-check", "/wp-includes/js/wp-auth-check"+suffix+".js", []string{"heartbeat", "wp-i18n"}, "", 1))
	m.Store("wp-lists", NewScript("wp-lists", "/wp-includes/js/wp-lists"+suffix+".js", []string{"wp-ajax-response", "jquery-color"}, "", 1))
	m.Store("prototype", NewScript("prototype", "https://ajax.googleapis.com/ajax/libs/prototype/1.7.1.0/prototype.js"+suffix+".js", nil, "1.7.1", nil))
	m.Store("scriptaculous-root", NewScript("scriptaculous-root", "https://ajax.googleapis.com/ajax/libs/scriptaculous/1.9.0/scriptaculous.js"+suffix+".js", []string{"prototype"}, "1.9.0", nil))
	m.Store("scriptaculous-builder", NewScript("scriptaculous-builder", "https://ajax.googleapis.com/ajax/libs/scriptaculous/1.9.0/builder.js"+suffix+".js", []string{"scriptaculous-root"}, "1.9.0", nil))
	m.Store("scriptaculous-dragdrop", NewScript("scriptaculous-dragdrop", "https://ajax.googleapis.com/ajax/libs/scriptaculous/1.9.0/dragdrop.js"+suffix+".js", []string{"scriptaculous-builder", "scriptaculous-effects"}, "1.9.0", nil))
	m.Store("scriptaculous-effects", NewScript("scriptaculous-effects", "https://ajax.googleapis.com/ajax/libs/scriptaculous/1.9.0/effects.js"+suffix+".js", []string{"scriptaculous-root"}, "1.9.0", nil))
	m.Store("scriptaculous-slider", NewScript("scriptaculous-slider", "https://ajax.googleapis.com/ajax/libs/scriptaculous/1.9.0/slider.js"+suffix+".js", []string{"scriptaculous-effects"}, "1.9.0", nil))
	m.Store("scriptaculous-sound", NewScript("scriptaculous-sound", "https://ajax.googleapis.com/ajax/libs/scriptaculous/1.9.0/sound.js"+suffix+".js", []string{"scriptaculous-root"}, "1.9.0", nil))
	m.Store("scriptaculous-controls", NewScript("scriptaculous-controls", "https://ajax.googleapis.com/ajax/libs/scriptaculous/1.9.0/controls.js"+suffix+".js", []string{"scriptaculous-root"}, "1.9.0", nil))
	m.Store("scriptaculous", NewScript("scriptaculous", ""+suffix+".js", []string{"scriptaculous-dragdrop", "scriptaculous-slider", "scriptaculous-controls"}, "", nil))
	m.Store("cropper", NewScript("cropper", "/wp-includes/js/crop/cropper.js"+suffix+".js", []string{"scriptaculous-dragdrop"}, "", nil))
	m.Store("jquery", NewScript("jquery", ""+suffix+".js", []string{"jquery-core", "jquery-migrate"}, "3.6.4", nil))
	m.Store("jquery-core", NewScript("jquery-core", "/wp-includes/js/jquery/jquery"+suffix+".js", nil, "3.6.4", nil))
	m.Store("jquery-migrate", NewScript("jquery-migrate", "/wp-includes/js/jquery/jquery-migrate"+suffix+".js", nil, "3.4.0", nil))
	m.Store("jquery-ui-core", NewScript("jquery-ui-core", "/wp-includes/js/jquery/ui/core"+suffix+".js", []string{"jquery"}, "1.13.2", 1))
	m.Store("jquery-effects-core", NewScript("jquery-effects-core", "/wp-includes/js/jquery/ui/effect"+suffix+".js", []string{"jquery"}, "1.13.2", 1))
	m.Store("jquery-effects-blind", NewScript("jquery-effects-blind", "/wp-includes/js/jquery/ui/effect-blind"+suffix+".js", []string{"jquery-effects-core"}, "1.13.2", 1))
	m.Store("jquery-effects-bounce", NewScript("jquery-effects-bounce", "/wp-includes/js/jquery/ui/effect-bounce"+suffix+".js", []string{"jquery-effects-core"}, "1.13.2", 1))
	m.Store("jquery-effects-clip", NewScript("jquery-effects-clip", "/wp-includes/js/jquery/ui/effect-clip"+suffix+".js", []string{"jquery-effects-core"}, "1.13.2", 1))
	m.Store("jquery-effects-drop", NewScript("jquery-effects-drop", "/wp-includes/js/jquery/ui/effect-drop"+suffix+".js", []string{"jquery-effects-core"}, "1.13.2", 1))
	m.Store("jquery-effects-explode", NewScript("jquery-effects-explode", "/wp-includes/js/jquery/ui/effect-explode"+suffix+".js", []string{"jquery-effects-core"}, "1.13.2", 1))
	m.Store("jquery-effects-fade", NewScript("jquery-effects-fade", "/wp-includes/js/jquery/ui/effect-fade"+suffix+".js", []string{"jquery-effects-core"}, "1.13.2", 1))
	m.Store("jquery-effects-fold", NewScript("jquery-effects-fold", "/wp-includes/js/jquery/ui/effect-fold"+suffix+".js", []string{"jquery-effects-core"}, "1.13.2", 1))
	m.Store("jquery-effects-highlight", NewScript("jquery-effects-highlight", "/wp-includes/js/jquery/ui/effect-highlight"+suffix+".js", []string{"jquery-effects-core"}, "1.13.2", 1))
	m.Store("jquery-effects-puff", NewScript("jquery-effects-puff", "/wp-includes/js/jquery/ui/effect-puff"+suffix+".js", []string{"jquery-effects-core", "jquery-effects-scale"}, "1.13.2", 1))
	m.Store("jquery-effects-pulsate", NewScript("jquery-effects-pulsate", "/wp-includes/js/jquery/ui/effect-pulsate"+suffix+".js", []string{"jquery-effects-core"}, "1.13.2", 1))
	m.Store("jquery-effects-scale", NewScript("jquery-effects-scale", "/wp-includes/js/jquery/ui/effect-scale"+suffix+".js", []string{"jquery-effects-core", "jquery-effects-size"}, "1.13.2", 1))
	m.Store("jquery-effects-shake", NewScript("jquery-effects-shake", "/wp-includes/js/jquery/ui/effect-shake"+suffix+".js", []string{"jquery-effects-core"}, "1.13.2", 1))
	m.Store("jquery-effects-size", NewScript("jquery-effects-size", "/wp-includes/js/jquery/ui/effect-size"+suffix+".js", []string{"jquery-effects-core"}, "1.13.2", 1))
	m.Store("jquery-effects-slide", NewScript("jquery-effects-slide", "/wp-includes/js/jquery/ui/effect-slide"+suffix+".js", []string{"jquery-effects-core"}, "1.13.2", 1))
	m.Store("jquery-effects-transfer", NewScript("jquery-effects-transfer", "/wp-includes/js/jquery/ui/effect-transfer"+suffix+".js", []string{"jquery-effects-core"}, "1.13.2", 1))
	m.Store("jquery-ui-accordion", NewScript("jquery-ui-accordion", "/wp-includes/js/jquery/ui/accordion"+suffix+".js", []string{"jquery-ui-core"}, "1.13.2", 1))
	m.Store("jquery-ui-autocomplete", NewScript("jquery-ui-autocomplete", "/wp-includes/js/jquery/ui/autocomplete"+suffix+".js", []string{"jquery-ui-menu", "wp-a11y"}, "1.13.2", 1))
	m.Store("jquery-ui-button", NewScript("jquery-ui-button", "/wp-includes/js/jquery/ui/button"+suffix+".js", []string{"jquery-ui-core", "jquery-ui-controlgroup", "jquery-ui-checkboxradio"}, "1.13.2", 1))
	m.Store("jquery-ui-datepicker", NewScript("jquery-ui-datepicker", "/wp-includes/js/jquery/ui/datepicker"+suffix+".js", []string{"jquery-ui-core"}, "1.13.2", 1))
	m.Store("jquery-ui-dialog", NewScript("jquery-ui-dialog", "/wp-includes/js/jquery/ui/dialog"+suffix+".js", []string{"jquery-ui-resizable", "jquery-ui-draggable", "jquery-ui-button"}, "1.13.2", 1))
	m.Store("jquery-ui-menu", NewScript("jquery-ui-menu", "/wp-includes/js/jquery/ui/menu"+suffix+".js", []string{"jquery-ui-core"}, "1.13.2", 1))
	m.Store("jquery-ui-mouse", NewScript("jquery-ui-mouse", "/wp-includes/js/jquery/ui/mouse"+suffix+".js", []string{"jquery-ui-core"}, "1.13.2", 1))
	m.Store("jquery-ui-progressbar", NewScript("jquery-ui-progressbar", "/wp-includes/js/jquery/ui/progressbar"+suffix+".js", []string{"jquery-ui-core"}, "1.13.2", 1))
	m.Store("jquery-ui-selectmenu", NewScript("jquery-ui-selectmenu", "/wp-includes/js/jquery/ui/selectmenu"+suffix+".js", []string{"jquery-ui-menu"}, "1.13.2", 1))
	m.Store("jquery-ui-slider", NewScript("jquery-ui-slider", "/wp-includes/js/jquery/ui/slider"+suffix+".js", []string{"jquery-ui-mouse"}, "1.13.2", 1))
	m.Store("jquery-ui-spinner", NewScript("jquery-ui-spinner", "/wp-includes/js/jquery/ui/spinner"+suffix+".js", []string{"jquery-ui-button"}, "1.13.2", 1))
	m.Store("jquery-ui-tabs", NewScript("jquery-ui-tabs", "/wp-includes/js/jquery/ui/tabs"+suffix+".js", []string{"jquery-ui-core"}, "1.13.2", 1))
	m.Store("jquery-ui-tooltip", NewScript("jquery-ui-tooltip", "/wp-includes/js/jquery/ui/tooltip"+suffix+".js", []string{"jquery-ui-core"}, "1.13.2", 1))
	m.Store("jquery-ui-checkboxradio", NewScript("jquery-ui-checkboxradio", "/wp-includes/js/jquery/ui/checkboxradio"+suffix+".js", []string{"jquery-ui-core"}, "1.13.2", 1))
	m.Store("jquery-ui-controlgroup", NewScript("jquery-ui-controlgroup", "/wp-includes/js/jquery/ui/controlgroup"+suffix+".js", []string{"jquery-ui-core"}, "1.13.2", 1))
	m.Store("jquery-ui-draggable", NewScript("jquery-ui-draggable", "/wp-includes/js/jquery/ui/draggable"+suffix+".js", []string{"jquery-ui-mouse"}, "1.13.2", 1))
	m.Store("jquery-ui-droppable", NewScript("jquery-ui-droppable", "/wp-includes/js/jquery/ui/droppable"+suffix+".js", []string{"jquery-ui-draggable"}, "1.13.2", 1))
	m.Store("jquery-ui-resizable", NewScript("jquery-ui-resizable", "/wp-includes/js/jquery/ui/resizable"+suffix+".js", []string{"jquery-ui-mouse"}, "1.13.2", 1))
	m.Store("jquery-ui-selectable", NewScript("jquery-ui-selectable", "/wp-includes/js/jquery/ui/selectable"+suffix+".js", []string{"jquery-ui-mouse"}, "1.13.2", 1))
	m.Store("jquery-ui-sortable", NewScript("jquery-ui-sortable", "/wp-includes/js/jquery/ui/sortable"+suffix+".js", []string{"jquery-ui-mouse"}, "1.13.2", 1))
	m.Store("jquery-ui-position", NewScript("jquery-ui-position", ""+suffix+".js", []string{"jquery-ui-core"}, "1.13.2", 1))
	m.Store("jquery-ui-widget", NewScript("jquery-ui-widget", ""+suffix+".js", []string{"jquery-ui-core"}, "1.13.2", 1))
	m.Store("jquery-form", NewScript("jquery-form", "/wp-includes/js/jquery/jquery.form"+suffix+".js", []string{"jquery"}, "4.3.0", 1))
	m.Store("jquery-color", NewScript("jquery-color", "/wp-includes/js/jquery/jquery.color"+suffix+".js", []string{"jquery"}, "2.2.0", 1))
	m.Store("schedule", NewScript("schedule", "/wp-includes/js/jquery/jquery.schedule.js"+suffix+".js", []string{"jquery"}, "20m", 1))
	m.Store("jquery-query", NewScript("jquery-query", "/wp-includes/js/jquery/jquery.query.js"+suffix+".js", []string{"jquery"}, "2.2.3", 1))
	m.Store("jquery-serialize-object", NewScript("jquery-serialize-object", "/wp-includes/js/jquery/jquery.serialize-object.js"+suffix+".js", []string{"jquery"}, "0.2-wp", 1))
	m.Store("jquery-hotkeys", NewScript("jquery-hotkeys", "/wp-includes/js/jquery/jquery.hotkeys"+suffix+".js", []string{"jquery"}, "0.0.2m", 1))
	m.Store("jquery-table-hotkeys", NewScript("jquery-table-hotkeys", "/wp-includes/js/jquery/jquery.table-hotkeys"+suffix+".js", []string{"jquery", "jquery-hotkeys"}, "", 1))
	m.Store("jquery-touch-punch", NewScript("jquery-touch-punch", "/wp-includes/js/jquery/jquery.ui.touch-punch.js"+suffix+".js", []string{"jquery-ui-core", "jquery-ui-mouse"}, "0.2.2", 1))
	m.Store("suggest", NewScript("suggest", "/wp-includes/js/jquery/suggest"+suffix+".js", []string{"jquery"}, "1.1-20110113", 1))
	m.Store("imagesloaded", NewScript("imagesloaded", "/wp-includes/js/imagesloaded"+suffix+".js", nil, "4.1.4", 1))
	m.Store("masonry", NewScript("masonry", "/wp-includes/js/masonry"+suffix+".js", []string{"imagesloaded"}, "4.2.2", 1))
	m.Store("jquery-masonry", NewScript("jquery-masonry", "/wp-includes/js/jquery/jquery.masonry"+suffix+".js", []string{"jquery", "masonry"}, "3.1.2b", 1))
	m.Store("thickbox", NewScript("thickbox", "/wp-includes/js/thickbox/thickbox.js"+suffix+".js", []string{"jquery"}, "3.1-20121105", 1))
	m.Store("jcrop", NewScript("jcrop", "/wp-includes/js/jcrop/jquery.Jcrop"+suffix+".js", []string{"jquery"}, "0.9.15", nil))
	m.Store("swfobject", NewScript("swfobject", "/wp-includes/js/swfobject.js"+suffix+".js", nil, "2.2-20120417", nil))
	m.Store("moxiejs", NewScript("moxiejs", "/wp-includes/js/plupload/moxie"+suffix+".js", nil, "1.3.5", nil))
	m.Store("plupload", NewScript("plupload", "/wp-includes/js/plupload/plupload"+suffix+".js", []string{"moxiejs"}, "2.1.9", nil))
	m.Store("plupload-all", NewScript("plupload-all", ""+suffix+".js", []string{"plupload"}, "2.1.1", nil))
	m.Store("plupload-html5", NewScript("plupload-html5", ""+suffix+".js", []string{"plupload"}, "2.1.1", nil))
	m.Store("plupload-flash", NewScript("plupload-flash", ""+suffix+".js", []string{"plupload"}, "2.1.1", nil))
	m.Store("plupload-silverlight", NewScript("plupload-silverlight", ""+suffix+".js", []string{"plupload"}, "2.1.1", nil))
	m.Store("plupload-html4", NewScript("plupload-html4", ""+suffix+".js", []string{"plupload"}, "2.1.1", nil))
	m.Store("plupload-handlers", NewScript("plupload-handlers", "/wp-includes/js/plupload/handlers"+suffix+".js", []string{"clipboard", "jquery", "plupload", "underscore", "wp-a11y", "wp-i18n"}, "", nil))
	m.Store("wp-plupload", NewScript("wp-plupload", "/wp-includes/js/plupload/wp-plupload"+suffix+".js", []string{"plupload", "jquery", "json2", "media-models"}, "", 1))
	m.Store("swfupload", NewScript("swfupload", "/wp-includes/js/swfupload/swfupload.js"+suffix+".js", nil, "2201-20110113", nil))
	m.Store("swfupload-all", NewScript("swfupload-all", ""+suffix+".js", []string{"swfupload"}, "2201", nil))
	m.Store("swfupload-handlers", NewScript("swfupload-handlers", "/wp-includes/js/swfupload/handlers"+suffix+".js", []string{"swfupload-all", "jquery"}, "2201-20110524", nil))
	m.Store("comment-reply", NewScript("comment-reply", "/wp-includes/js/comment-reply"+suffix+".js", nil, "", 1))
	m.Store("json2", NewScript("json2", "/wp-includes/js/json2"+suffix+".js", nil, "2015-05-03", nil))
	m.Store("underscore", NewScript("underscore", "/wp-includes/js/underscore"+suffix+".js", nil, "1.13.4", 1))
	m.Store("backbone", NewScript("backbone", "/wp-includes/js/backbone"+suffix+".js", []string{"underscore", "jquery"}, "1.4.1", 1))
	m.Store("wp-util", NewScript("wp-util", "/wp-includes/js/wp-util"+suffix+".js", []string{"underscore", "jquery"}, "", 1))
	m.Store("wp-backbone", NewScript("wp-backbone", "/wp-includes/js/wp-backbone"+suffix+".js", []string{"backbone", "wp-util"}, "", 1))
	m.Store("revisions", NewScript("revisions", "/wp-admin/js/revisions"+suffix+".js", []string{"wp-backbone", "jquery-ui-slider", "hoverIntent"}, "", 1))
	m.Store("imgareaselect", NewScript("imgareaselect", "/wp-includes/js/imgareaselect/jquery.imgareaselect"+suffix+".js", []string{"jquery"}, "", 1))
	m.Store("mediaelement", NewScript("mediaelement", ""+suffix+".js", []string{"jquery", "mediaelement-core", "mediaelement-migrate"}, "4.2.17", 1))
	m.Store("mediaelement-core", NewScript("mediaelement-core", "/wp-includes/js/mediaelement/mediaelement-and-player"+suffix+".js", nil, "4.2.17", 1))
	m.Store("mediaelement-migrate", NewScript("mediaelement-migrate", "/wp-includes/js/mediaelement/mediaelement-migrate"+suffix+".js", nil, "", 1))
	m.Store("mediaelement-vimeo", NewScript("mediaelement-vimeo", "/wp-includes/js/mediaelement/renderers/vimeo"+suffix+".js", []string{"mediaelement"}, "4.2.17", 1))
	m.Store("wp-mediaelement", NewScript("wp-mediaelement", "/wp-includes/js/mediaelement/wp-mediaelement"+suffix+".js", []string{"mediaelement"}, "", 1))
	m.Store("wp-codemirror", NewScript("wp-codemirror", "/wp-includes/js/codemirror/codemirror"+suffix+".js", nil, "5.29.1-alpha-ee20357", nil))
	m.Store("csslint", NewScript("csslint", "/wp-includes/js/codemirror/csslint.js"+suffix+".js", nil, "1.0.5", nil))
	m.Store("esprima", NewScript("esprima", "/wp-includes/js/codemirror/esprima.js"+suffix+".js", nil, "4.0.0", nil))
	m.Store("jshint", NewScript("jshint", "/wp-includes/js/codemirror/fakejshint.js"+suffix+".js", []string{"esprima"}, "2.9.5", nil))
	m.Store("jsonlint", NewScript("jsonlint", "/wp-includes/js/codemirror/jsonlint.js"+suffix+".js", nil, "1.6.2", nil))
	m.Store("htmlhint", NewScript("htmlhint", "/wp-includes/js/codemirror/htmlhint.js"+suffix+".js", nil, "0.9.14-xwp", nil))
	m.Store("htmlhint-kses", NewScript("htmlhint-kses", "/wp-includes/js/codemirror/htmlhint-kses.js"+suffix+".js", []string{"htmlhint"}, "", nil))
	m.Store("code-editor", NewScript("code-editor", "/wp-admin/js/code-editor"+suffix+".js", []string{"jquery", "wp-codemirror", "underscore"}, "", nil))
	m.Store("wp-theme-plugin-editor", NewScript("wp-theme-plugin-editor", "/wp-admin/js/theme-plugin-editor"+suffix+".js", []string{"common", "wp-util", "wp-sanitize", "jquery", "jquery-ui-core", "wp-a11y", "underscore", "wp-i18n"}, "", 1))
	m.Store("wp-playlist", NewScript("wp-playlist", "/wp-includes/js/mediaelement/wp-playlist"+suffix+".js", []string{"wp-util", "backbone", "mediaelement"}, "", 1))
	m.Store("zxcvbn-async", NewScript("zxcvbn-async", "/wp-includes/js/zxcvbn-async"+suffix+".js", nil, "1.0", nil))
	m.Store("password-strength-meter", NewScript("password-strength-meter", "/wp-admin/js/password-strength-meter"+suffix+".js", []string{"jquery", "zxcvbn-async", "wp-i18n"}, "", 1))
	m.Store("application-passwords", NewScript("application-passwords", "/wp-admin/js/application-passwords"+suffix+".js", []string{"jquery", "wp-util", "wp-api-request", "wp-date", "wp-i18n", "wp-hooks"}, "", 1))
	m.Store("auth-app", NewScript("auth-app", "/wp-admin/js/auth-app"+suffix+".js", []string{"jquery", "wp-api-request", "wp-i18n", "wp-hooks"}, "", 1))
	m.Store("user-profile", NewScript("user-profile", "/wp-admin/js/user-profile"+suffix+".js", []string{"jquery", "password-strength-meter", "wp-util", "wp-i18n"}, "", 1))
	m.Store("language-chooser", NewScript("language-chooser", "/wp-admin/js/language-chooser"+suffix+".js", []string{"jquery"}, "", 1))
	m.Store("user-suggest", NewScript("user-suggest", "/wp-admin/js/user-suggest"+suffix+".js", []string{"jquery-ui-autocomplete"}, "", 1))
	m.Store("admin-bar", NewScript("admin-bar", "/wp-includes/js/admin-bar"+suffix+".js", []string{"hoverintent-js"}, "", 1))
	m.Store("wplink", NewScript("wplink", "/wp-includes/js/wplink"+suffix+".js", []string{"jquery", "wp-a11y"}, "", 1))
	m.Store("wpdialogs", NewScript("wpdialogs", "/wp-includes/js/wpdialog"+suffix+".js", []string{"jquery-ui-dialog"}, "", 1))
	m.Store("word-count", NewScript("word-count", "/wp-admin/js/word-count"+suffix+".js", nil, "", 1))
	m.Store("media-upload", NewScript("media-upload", "/wp-admin/js/media-upload"+suffix+".js", []string{"thickbox", "shortcode"}, "", 1))
	m.Store("hoverIntent", NewScript("hoverIntent", "/wp-includes/js/hoverIntent"+suffix+".js", []string{"jquery"}, "1.10.2", 1))
	m.Store("hoverintent-js", NewScript("hoverintent-js", "/wp-includes/js/hoverintent-js"+suffix+".js", nil, "2.2.1", 1))
	m.Store("customize-base", NewScript("customize-base", "/wp-includes/js/customize-base"+suffix+".js", []string{"jquery", "json2", "underscore"}, "", 1))
	m.Store("customize-loader", NewScript("customize-loader", "/wp-includes/js/customize-loader"+suffix+".js", []string{"customize-base"}, "", 1))
	m.Store("customize-preview", NewScript("customize-preview", "/wp-includes/js/customize-preview"+suffix+".js", []string{"wp-a11y", "customize-base"}, "", 1))
	m.Store("customize-models", NewScript("customize-models", "/wp-includes/js/customize-models.js"+suffix+".js", []string{"underscore", "backbone"}, "", 1))
	m.Store("customize-views", NewScript("customize-views", "/wp-includes/js/customize-views.js"+suffix+".js", []string{"jquery", "underscore", "imgareaselect", "customize-models", "media-editor", "media-views"}, "", 1))
	m.Store("customize-controls", NewScript("customize-controls", "/wp-admin/js/customize-controls"+suffix+".js", []string{"customize-base", "wp-a11y", "wp-util", "jquery-ui-core"}, "", 1))
	m.Store("customize-selective-refresh", NewScript("customize-selective-refresh", "/wp-includes/js/customize-selective-refresh"+suffix+".js", []string{"jquery", "wp-util", "customize-preview"}, "", 1))
	m.Store("customize-widgets", NewScript("customize-widgets", "/wp-admin/js/customize-widgets"+suffix+".js", []string{"jquery", "jquery-ui-sortable", "jquery-ui-droppable", "wp-backbone", "customize-controls"}, "", 1))
	m.Store("customize-preview-widgets", NewScript("customize-preview-widgets", "/wp-includes/js/customize-preview-widgets"+suffix+".js", []string{"jquery", "wp-util", "customize-preview", "customize-selective-refresh"}, "", 1))
	m.Store("customize-nav-menus", NewScript("customize-nav-menus", "/wp-admin/js/customize-nav-menus"+suffix+".js", []string{"jquery", "wp-backbone", "customize-controls", "accordion", "nav-menu", "wp-sanitize"}, "", 1))
	m.Store("customize-preview-nav-menus", NewScript("customize-preview-nav-menus", "/wp-includes/js/customize-preview-nav-menus"+suffix+".js", []string{"jquery", "wp-util", "customize-preview", "customize-selective-refresh"}, "", 1))
	m.Store("wp-custom-header", NewScript("wp-custom-header", "/wp-includes/js/wp-custom-header"+suffix+".js", []string{"wp-a11y"}, "", 1))
	m.Store("accordion", NewScript("accordion", "/wp-admin/js/accordion"+suffix+".js", []string{"jquery"}, "", 1))
	m.Store("shortcode", NewScript("shortcode", "/wp-includes/js/shortcode"+suffix+".js", []string{"underscore"}, "", 1))
	m.Store("media-models", NewScript("media-models", "/wp-includes/js/media-models"+suffix+".js", []string{"wp-backbone"}, "", 1))
	m.Store("wp-embed", NewScript("wp-embed", "/wp-includes/js/wp-embed"+suffix+".js", nil, "", 1))
	m.Store("media-views", NewScript("media-views", "/wp-includes/js/media-views"+suffix+".js", []string{"utils", "media-models", "wp-plupload", "jquery-ui-sortable", "wp-mediaelement", "wp-api-request", "wp-a11y", "clipboard", "wp-i18n"}, "", 1))
	m.Store("media-editor", NewScript("media-editor", "/wp-includes/js/media-editor"+suffix+".js", []string{"shortcode", "media-views", "wp-i18n"}, "", 1))
	m.Store("media-audiovideo", NewScript("media-audiovideo", "/wp-includes/js/media-audiovideo"+suffix+".js", []string{"media-editor"}, "", 1))
	m.Store("mce-view", NewScript("mce-view", "/wp-includes/js/mce-view"+suffix+".js", []string{"shortcode", "jquery", "media-views", "media-audiovideo"}, "", 1))
	m.Store("wp-api", NewScript("wp-api", "/wp-includes/js/wp-api"+suffix+".js", []string{"jquery", "backbone", "underscore", "wp-api-request"}, "", 1))
	m.Store("react", NewScript("react", "/wp-includes/js/dist/vendor/react"+suffix+".js", []string{"wp-polyfill"}, "18.2.0", 1))
	m.Store("react-dom", NewScript("react-dom", "/wp-includes/js/dist/vendor/react-dom"+suffix+".js", []string{"react"}, "18.2.0", 1))
	m.Store("regenerator-runtime", NewScript("regenerator-runtime", "/wp-includes/js/dist/vendor/regenerator-runtime"+suffix+".js", nil, "0.13.11", 1))
	m.Store("moment", NewScript("moment", "/wp-includes/js/dist/vendor/moment"+suffix+".js", nil, "2.29.4", 1))
	m.Store("lodash", NewScript("lodash", "/wp-includes/js/dist/vendor/lodash"+suffix+".js", nil, "4.17.19", 1))
	m.Store("wp-polyfill-fetch", NewScript("wp-polyfill-fetch", "/wp-includes/js/dist/vendor/wp-polyfill-fetch"+suffix+".js", nil, "3.6.2", 1))
	m.Store("wp-polyfill-formdata", NewScript("wp-polyfill-formdata", "/wp-includes/js/dist/vendor/wp-polyfill-formdata"+suffix+".js", nil, "4.0.10", 1))
	m.Store("wp-polyfill-node-contains", NewScript("wp-polyfill-node-contains", "/wp-includes/js/dist/vendor/wp-polyfill-node-contains"+suffix+".js", nil, "4.6.0", 1))
	m.Store("wp-polyfill-url", NewScript("wp-polyfill-url", "/wp-includes/js/dist/vendor/wp-polyfill-url"+suffix+".js", nil, "3.6.4", 1))
	m.Store("wp-polyfill-dom-rect", NewScript("wp-polyfill-dom-rect", "/wp-includes/js/dist/vendor/wp-polyfill-dom-rect"+suffix+".js", nil, "4.6.0", 1))
	m.Store("wp-polyfill-element-closest", NewScript("wp-polyfill-element-closest", "/wp-includes/js/dist/vendor/wp-polyfill-element-closest"+suffix+".js", nil, "3.0.2", 1))
	m.Store("wp-polyfill-object-fit", NewScript("wp-polyfill-object-fit", "/wp-includes/js/dist/vendor/wp-polyfill-object-fit"+suffix+".js", nil, "2.3.5", 1))
	m.Store("wp-polyfill-inert", NewScript("wp-polyfill-inert", "/wp-includes/js/dist/vendor/wp-polyfill-inert"+suffix+".js", nil, "3.1.2", 1))
	m.Store("wp-polyfill", NewScript("wp-polyfill", "/wp-includes/js/dist/vendor/wp-polyfill"+suffix+".js", []string{"wp-polyfill-inert", "regenerator-runtime"}, "3.15.0", 1))
	m.Store("wp-tinymce-root", NewScript("wp-tinymce-root", "http://wp.test/wp-includes/js/tinymce/tinymce"+suffix+".js", nil, "49110-20201110", nil))
	m.Store("wp-tinymce", NewScript("wp-tinymce", "http://wp.test/wp-includes/js/tinymce/plugins/compat3x/plugin"+suffix+".js", []string{"wp-tinymce-root"}, "49110-20201110", nil))
	m.Store("wp-tinymce-lists", NewScript("wp-tinymce-lists", "http://wp.test/wp-includes/js/tinymce/plugins/lists/plugin"+suffix+".js", []string{"wp-tinymce"}, "49110-20201110", nil))
	m.Store("wp-a11y", NewScript("wp-a11y", "/wp-includes/js/dist/a11y"+suffix+".js", []string{"wp-dom-ready", "wp-i18n", "wp-polyfill"}, "ecce20f002eda4c19664", 1))
	m.Store("wp-annotations", NewScript("wp-annotations", "/wp-includes/js/dist/annotations"+suffix+".js", []string{"wp-data", "wp-hooks", "wp-i18n", "wp-polyfill", "wp-rich-text"}, "1720fc5d5c76f53a1740", 1))
	m.Store("wp-api-fetch", NewScript("wp-api-fetch", "/wp-includes/js/dist/api-fetch"+suffix+".js", []string{"wp-i18n", "wp-polyfill", "wp-url"}, "bc0029ca2c943aec5311", 1))
	m.Store("wp-autop", NewScript("wp-autop", "/wp-includes/js/dist/autop"+suffix+".js", []string{"wp-polyfill"}, "43197d709df445ccf849", 1))
	m.Store("wp-blob", NewScript("wp-blob", "/wp-includes/js/dist/blob"+suffix+".js", []string{"wp-polyfill"}, "e7b4ea96175a89b263e2", 1))
	m.Store("wp-block-directory", NewScript("wp-block-directory", "/wp-includes/js/dist/block-directory"+suffix+".js", []string{"wp-a11y", "wp-api-fetch", "wp-block-editor", "wp-blocks", "wp-components", "wp-compose", "wp-core-data", "wp-data", "wp-editor", "wp-element", "wp-hooks", "wp-html-entities", "wp-i18n", "wp-notices", "wp-plugins", "wp-polyfill", "wp-primitives", "wp-url"}, "9c45b8d28fc867ceed45", 1))
	m.Store("wp-block-editor", NewScript("wp-block-editor", "/wp-includes/js/dist/block-editor"+suffix+".js", []string{"lodash", "react", "react-dom", "wp-a11y", "wp-api-fetch", "wp-blob", "wp-blocks", "wp-components", "wp-compose", "wp-data", "wp-date", "wp-deprecated", "wp-dom", "wp-element", "wp-escape-html", "wp-hooks", "wp-html-entities", "wp-i18n", "wp-is-shallow-equal", "wp-keyboard-shortcuts", "wp-keycodes", "wp-notices", "wp-polyfill", "wp-preferences", "wp-primitives", "wp-private-apis", "wp-rich-text", "wp-shortcode", "wp-style-engine", "wp-token-list", "wp-url", "wp-warning", "wp-wordcount"}, "43e40e04f77d598ede94", 1))
	m.Store("wp-block-library", NewScript("wp-block-library", "/wp-includes/js/dist/block-library"+suffix+".js", []string{"lodash", "wp-a11y", "wp-api-fetch", "wp-autop", "wp-blob", "wp-block-editor", "wp-blocks", "wp-components", "wp-compose", "wp-core-data", "wp-data", "wp-date", "wp-deprecated", "wp-dom", "wp-element", "wp-escape-html", "wp-hooks", "wp-html-entities", "wp-i18n", "wp-keycodes", "wp-notices", "wp-polyfill", "wp-primitives", "wp-private-apis", "wp-reusable-blocks", "wp-rich-text", "wp-server-side-render", "wp-url", "wp-viewport", "editor"}, "3115f0b5551a55bb6d3b", 1))
	m.Store("wp-block-serialization-default-parser", NewScript("wp-block-serialization-default-parser", "/wp-includes/js/dist/block-serialization-default-parser"+suffix+".js", []string{"wp-polyfill"}, "30ffd7e7e199f10b2a6d", 1))
	m.Store("wp-blocks", NewScript("wp-blocks", "/wp-includes/js/dist/blocks"+suffix+".js", []string{"lodash", "wp-autop", "wp-blob", "wp-block-serialization-default-parser", "wp-compose", "wp-data", "wp-deprecated", "wp-dom", "wp-element", "wp-hooks", "wp-html-entities", "wp-i18n", "wp-is-shallow-equal", "wp-polyfill", "wp-shortcode"}, "639e14271099dc3d85bf", 1))
	m.Store("wp-components", NewScript("wp-components", "/wp-includes/js/dist/components"+suffix+".js", []string{"lodash", "react", "react-dom", "wp-a11y", "wp-compose", "wp-date", "wp-deprecated", "wp-dom", "wp-element", "wp-escape-html", "wp-hooks", "wp-html-entities", "wp-i18n", "wp-is-shallow-equal", "wp-keycodes", "wp-polyfill", "wp-primitives", "wp-private-apis", "wp-rich-text", "wp-warning"}, "bf6e0ec3089253604b52", 1))
	m.Store("wp-compose", NewScript("wp-compose", "/wp-includes/js/dist/compose"+suffix+".js", []string{"react", "wp-deprecated", "wp-dom", "wp-element", "wp-is-shallow-equal", "wp-keycodes", "wp-polyfill", "wp-priority-queue"}, "7d5916e3b2ef0ea01400", 1))
	m.Store("wp-core-data", NewScript("wp-core-data", "/wp-includes/js/dist/core-data"+suffix+".js", []string{"lodash", "wp-api-fetch", "wp-blocks", "wp-compose", "wp-data", "wp-deprecated", "wp-element", "wp-html-entities", "wp-i18n", "wp-is-shallow-equal", "wp-polyfill", "wp-url"}, "fc0de6bb17aa25caf698", 1))
	m.Store("wp-customize-widgets", NewScript("wp-customize-widgets", "/wp-includes/js/dist/customize-widgets"+suffix+".js", []string{"wp-block-editor", "wp-block-library", "wp-blocks", "wp-components", "wp-compose", "wp-core-data", "wp-data", "wp-deprecated", "wp-dom", "wp-element", "wp-hooks", "wp-i18n", "wp-is-shallow-equal", "wp-keyboard-shortcuts", "wp-keycodes", "wp-media-utils", "wp-polyfill", "wp-preferences", "wp-primitives", "wp-private-apis", "wp-widgets"}, "7ae69cc350436c0cf301", 1))
	m.Store("wp-data", NewScript("wp-data", "/wp-includes/js/dist/data"+suffix+".js", []string{"lodash", "wp-compose", "wp-deprecated", "wp-element", "wp-is-shallow-equal", "wp-polyfill", "wp-priority-queue", "wp-private-apis", "wp-redux-routine"}, "90cebfec01d1a3f0368e", 1))
	m.Store("wp-data-controls", NewScript("wp-data-controls", "/wp-includes/js/dist/data-controls"+suffix+".js", []string{"wp-api-fetch", "wp-data", "wp-deprecated", "wp-polyfill"}, "e10d473d392daa8501e8", 1))
	m.Store("wp-date", NewScript("wp-date", "/wp-includes/js/dist/date"+suffix+".js", []string{"moment", "wp-deprecated", "wp-polyfill"}, "f8550b1212d715fbf745", 1))
	m.Store("wp-deprecated", NewScript("wp-deprecated", "/wp-includes/js/dist/deprecated"+suffix+".js", []string{"wp-hooks", "wp-polyfill"}, "6c963cb9494ba26b77eb", 1))
	m.Store("wp-dom", NewScript("wp-dom", "/wp-includes/js/dist/dom"+suffix+".js", []string{"wp-deprecated", "wp-polyfill"}, "e03c89e1dd68aee1cb3a", 1))
	m.Store("wp-dom-ready", NewScript("wp-dom-ready", "/wp-includes/js/dist/dom-ready"+suffix+".js", []string{"wp-polyfill"}, "392bdd43726760d1f3ca", 1))
	m.Store("wp-edit-post", NewScript("wp-edit-post", "/wp-includes/js/dist/edit-post"+suffix+".js", []string{"lodash", "wp-a11y", "wp-api-fetch", "wp-block-editor", "wp-block-library", "wp-blocks", "wp-components", "wp-compose", "wp-core-data", "wp-data", "wp-deprecated", "wp-editor", "wp-element", "wp-hooks", "wp-i18n", "wp-keyboard-shortcuts", "wp-keycodes", "wp-media-utils", "wp-notices", "wp-plugins", "wp-polyfill", "wp-preferences", "wp-primitives", "wp-private-apis", "wp-url", "wp-viewport", "wp-warning", "wp-widgets", "media-models", "media-views", "postbox", "wp-dom-ready"}, "d098b8ee5bdffa238c03", 1))
	m.Store("wp-edit-site", NewScript("wp-edit-site", "/wp-includes/js/dist/edit-site"+suffix+".js", []string{"lodash", "react", "wp-a11y", "wp-api-fetch", "wp-block-editor", "wp-block-library", "wp-blocks", "wp-components", "wp-compose", "wp-core-data", "wp-data", "wp-deprecated", "wp-editor", "wp-element", "wp-hooks", "wp-html-entities", "wp-i18n", "wp-keyboard-shortcuts", "wp-keycodes", "wp-media-utils", "wp-notices", "wp-plugins", "wp-polyfill", "wp-preferences", "wp-primitives", "wp-private-apis", "wp-reusable-blocks", "wp-url", "wp-viewport", "wp-widgets"}, "fcf81e803ab1af60d4f8", 1))
	m.Store("wp-edit-widgets", NewScript("wp-edit-widgets", "/wp-includes/js/dist/edit-widgets"+suffix+".js", []string{"wp-api-fetch", "wp-block-editor", "wp-block-library", "wp-blocks", "wp-components", "wp-compose", "wp-core-data", "wp-data", "wp-deprecated", "wp-dom", "wp-element", "wp-hooks", "wp-i18n", "wp-keyboard-shortcuts", "wp-keycodes", "wp-media-utils", "wp-notices", "wp-plugins", "wp-polyfill", "wp-preferences", "wp-primitives", "wp-private-apis", "wp-reusable-blocks", "wp-url", "wp-viewport", "wp-widgets"}, "d683d5fc75e655fdf974", 1))
	m.Store("wp-editor", NewScript("wp-editor", "/wp-includes/js/dist/editor"+suffix+".js", []string{"lodash", "react", "wp-a11y", "wp-api-fetch", "wp-blob", "wp-block-editor", "wp-blocks", "wp-components", "wp-compose", "wp-core-data", "wp-data", "wp-date", "wp-deprecated", "wp-dom", "wp-element", "wp-hooks", "wp-html-entities", "wp-i18n", "wp-keyboard-shortcuts", "wp-keycodes", "wp-media-utils", "wp-notices", "wp-polyfill", "wp-preferences", "wp-primitives", "wp-private-apis", "wp-reusable-blocks", "wp-rich-text", "wp-server-side-render", "wp-url", "wp-wordcount"}, "1fb5fcf129627da4939e", 1))
	m.Store("wp-element", NewScript("wp-element", "/wp-includes/js/dist/element"+suffix+".js", []string{"react", "react-dom", "wp-escape-html", "wp-polyfill"}, "b3bda690cfc516378771", 1))
	m.Store("wp-escape-html", NewScript("wp-escape-html", "/wp-includes/js/dist/escape-html"+suffix+".js", []string{"wp-polyfill"}, "03e27a7b6ae14f7afaa6", 1))
	m.Store("wp-format-library", NewScript("wp-format-library", "/wp-includes/js/dist/format-library"+suffix+".js", []string{"wp-a11y", "wp-block-editor", "wp-components", "wp-data", "wp-element", "wp-html-entities", "wp-i18n", "wp-polyfill", "wp-primitives", "wp-rich-text", "wp-url"}, "cd4a10ec005e2f001978", 1))
	m.Store("wp-hooks", NewScript("wp-hooks", "/wp-includes/js/dist/hooks"+suffix+".js", []string{"wp-polyfill"}, "4169d3cf8e8d95a3d6d5", 1))
	m.Store("wp-html-entities", NewScript("wp-html-entities", "/wp-includes/js/dist/html-entities"+suffix+".js", []string{"wp-polyfill"}, "36a4a255da7dd2e1bf8e", 1))
	m.Store("wp-i18n", NewScript("wp-i18n", "/wp-includes/js/dist/i18n"+suffix+".js", []string{"wp-hooks", "wp-polyfill"}, "9e794f35a71bb98672ae", 1))
	m.Store("wp-is-shallow-equal", NewScript("wp-is-shallow-equal", "/wp-includes/js/dist/is-shallow-equal"+suffix+".js", []string{"wp-polyfill"}, "20c2b06ecf04afb14fee", 1))
	m.Store("wp-keyboard-shortcuts", NewScript("wp-keyboard-shortcuts", "/wp-includes/js/dist/keyboard-shortcuts"+suffix+".js", []string{"wp-data", "wp-element", "wp-keycodes", "wp-polyfill"}, "b696c16720133edfc065", 1))
	m.Store("wp-keycodes", NewScript("wp-keycodes", "/wp-includes/js/dist/keycodes"+suffix+".js", []string{"wp-i18n", "wp-polyfill"}, "184b321fa2d3bc7fd173", 1))
	m.Store("wp-list-reusable-blocks", NewScript("wp-list-reusable-blocks", "/wp-includes/js/dist/list-reusable-blocks"+suffix+".js", []string{"wp-api-fetch", "wp-components", "wp-compose", "wp-element", "wp-i18n", "wp-polyfill"}, "6ba78be26d660b6af113", 1))
	m.Store("wp-media-utils", NewScript("wp-media-utils", "/wp-includes/js/dist/media-utils"+suffix+".js", []string{"wp-api-fetch", "wp-blob", "wp-element", "wp-i18n", "wp-polyfill"}, "f837b6298c83612cd6f6", 1))
	m.Store("wp-notices", NewScript("wp-notices", "/wp-includes/js/dist/notices"+suffix+".js", []string{"wp-data", "wp-polyfill"}, "9c1575b7a31659f45a45", 1))
	m.Store("wp-nux", NewScript("wp-nux", "/wp-includes/js/dist/nux"+suffix+".js", []string{"wp-components", "wp-compose", "wp-data", "wp-deprecated", "wp-element", "wp-i18n", "wp-polyfill", "wp-primitives"}, "038c48e26a91639ae8ab", 1))
	m.Store("wp-plugins", NewScript("wp-plugins", "/wp-includes/js/dist/plugins"+suffix+".js", []string{"wp-compose", "wp-element", "wp-hooks", "wp-polyfill", "wp-primitives"}, "0d1b90278bae7df6ecf9", 1))
	m.Store("wp-preferences", NewScript("wp-preferences", "/wp-includes/js/dist/preferences"+suffix+".js", []string{"wp-a11y", "wp-components", "wp-data", "wp-element", "wp-i18n", "wp-polyfill", "wp-primitives", "wp-preferences-persistence"}, "c66e137a7e588dab54c3", 1))
	m.Store("wp-preferences-persistence", NewScript("wp-preferences-persistence", "/wp-includes/js/dist/preferences-persistence"+suffix+".js", []string{"wp-api-fetch", "wp-polyfill"}, "c5543628aa7ff5bd5be4", 1))
	m.Store("wp-primitives", NewScript("wp-primitives", "/wp-includes/js/dist/primitives"+suffix+".js", []string{"wp-element", "wp-polyfill"}, "dfac1545e52734396640", 1))
	m.Store("wp-priority-queue", NewScript("wp-priority-queue", "/wp-includes/js/dist/priority-queue"+suffix+".js", []string{"wp-polyfill"}, "422e19e9d48b269c5219", 1))
	m.Store("wp-private-apis", NewScript("wp-private-apis", "/wp-includes/js/dist/private-apis"+suffix+".js", []string{"wp-polyfill"}, "6f247ed2bc3571743bba", 1))
	m.Store("wp-redux-routine", NewScript("wp-redux-routine", "/wp-includes/js/dist/redux-routine"+suffix+".js", []string{"wp-polyfill"}, "d86e7e9f062d7582f76b", 1))
	m.Store("wp-reusable-blocks", NewScript("wp-reusable-blocks", "/wp-includes/js/dist/reusable-blocks"+suffix+".js", []string{"wp-block-editor", "wp-blocks", "wp-components", "wp-core-data", "wp-data", "wp-element", "wp-i18n", "wp-notices", "wp-polyfill", "wp-primitives", "wp-url"}, "a7367a6154c724b51b31", 1))
	m.Store("wp-rich-text", NewScript("wp-rich-text", "/wp-includes/js/dist/rich-text"+suffix+".js", []string{"wp-a11y", "wp-compose", "wp-data", "wp-deprecated", "wp-element", "wp-escape-html", "wp-i18n", "wp-keycodes", "wp-polyfill"}, "9307ec04c67d79b6e813", 1))
	m.Store("wp-server-side-render", NewScript("wp-server-side-render", "/wp-includes/js/dist/server-side-render"+suffix+".js", []string{"wp-api-fetch", "wp-blocks", "wp-components", "wp-compose", "wp-data", "wp-element", "wp-i18n", "wp-polyfill", "wp-url"}, "d1bc93277666143a3f5e", 1))
	m.Store("wp-shortcode", NewScript("wp-shortcode", "/wp-includes/js/dist/shortcode"+suffix+".js", []string{"wp-polyfill"}, "7539044b04e6bca57f2e", 1))
	m.Store("wp-style-engine", NewScript("wp-style-engine", "/wp-includes/js/dist/style-engine"+suffix+".js", []string{"lodash", "wp-polyfill"}, "528e6cf281ffc9b7bd3c", 1))
	m.Store("wp-token-list", NewScript("wp-token-list", "/wp-includes/js/dist/token-list"+suffix+".js", []string{"wp-polyfill"}, "f2cf0bb3ae80de227e43", 1))
	m.Store("wp-url", NewScript("wp-url", "/wp-includes/js/dist/url"+suffix+".js", []string{"wp-polyfill"}, "16185fce2fb043a0cfed", 1))
	m.Store("wp-viewport", NewScript("wp-viewport", "/wp-includes/js/dist/viewport"+suffix+".js", []string{"wp-compose", "wp-data", "wp-element", "wp-polyfill"}, "4f6bd168b2b8b45c8a6b", 1))
	m.Store("wp-warning", NewScript("wp-warning", "/wp-includes/js/dist/warning"+suffix+".js", []string{"wp-polyfill"}, "4acee5fc2fd9a24cefc2", 1))
	m.Store("wp-widgets", NewScript("wp-widgets", "/wp-includes/js/dist/widgets"+suffix+".js", []string{"wp-api-fetch", "wp-block-editor", "wp-blocks", "wp-components", "wp-compose", "wp-core-data", "wp-data", "wp-element", "wp-i18n", "wp-notices", "wp-polyfill", "wp-primitives"}, "040ac8be5e0cfc4b52df", 1))
	m.Store("wp-wordcount", NewScript("wp-wordcount", "/wp-includes/js/dist/wordcount"+suffix+".js", []string{"wp-polyfill"}, "feb9569307aec24292f2", 1))
}

func defaultLocalize() {
	/*AddDynamicLocalize(h, "utils", "userSettings", map[string]any{
		"url":    h.C.Request.RequestURI,
		"uid":    "0",
		"time":   number.IntToString(time.Now().Unix()),
		"secure": h.IsHttps(),
	})*/

	AddStaticLocalize("wp-ajax-response", "wpAjax", map[string]any{
		"noPerm": `抱歉，您不能这么做。`,
		"broken": `出现了问题。`,
	})
	AddStaticLocalize("wp-api-request", "wpApiSettings", map[string]any{
		"root":          `/wp-json/`,
		"nonce":         `9a3ce0a5d7`,
		"versionString": `wp/v2/`,
	})
	AddStaticLocalize("jquery-ui-autocomplete", "uiAutocompleteL10n", map[string]any{
		"noResults":    `未找到结果。`,
		"oneResult":    `找到1个结果。使用上下方向键来导航。`,
		"manyResults":  `找到%d个结果。使用上下方向键来导航。`,
		"itemSelected": `已选择项目。`,
	})
	AddStaticLocalize("thickbox", "thickboxL10n", map[string]any{
		"next":             `下一页 &gt;`,
		"prev":             `&lt; 上一页`,
		"image":            `图片`,
		"of":               `/`,
		"close":            `关闭`,
		"noiframes":        `这个功能需要iframe的支持。您可能禁止了iframe的显示，或您的浏览器不支持此功能。`,
		"loadingAnimation": `/wp-includes/js/thickbox/loadingAnimation.gif`,
	})

	AddStaticLocalize("mediaelement", "_wpmejsSettings", map[string]any{
		"pluginPath":            `/wp-includes/js/mediaelement/`,
		"classPrefix":           `mejs-`,
		"stretching":            `responsive`,
		"audioShortcodeLibrary": `mediaelement`,
		"videoShortcodeLibrary": `mediaelement`,
	})
	AddStaticLocalize("zxcvbn-async", "_zxcvbnSettings", map[string]any{
		"src": `/wp-includes/js/zxcvbn.min.js`,
	})

	AddStaticLocalize("mce-view", "mceViewL10n", map[string]any{
		"shortcodes": []string{"wp_caption", "caption", "gallery", "playlist", "audio", "video", "embed"},
	})
	AddStaticLocalize("word-count", "wordCountL10n", map[string]any{
		"type":       `characters_excluding_spaces`,
		"shortcodes": []string{"wp_caption", "caption", "gallery", "playlist", "audio", "video", "embed"},
	})

}

func defaultTranslate() {
	SetTranslation("common", "default", "")
	SetTranslation("wp-pointer", "default", "")
	SetTranslation("wp-auth-check", "default", "")
	SetTranslation("wp-theme-plugin-editor", "default", "")
	SetTranslation("password-strength-meter", "default", "")
	SetTranslation("application-passwords", "default", "")
	SetTranslation("auth-app", "default", "")
	SetTranslation("user-profile", "default", "")
	SetTranslation("media-views", "default", "")
	SetTranslation("media-editor", "default", "")
	SetTranslation("wp-a11y", "default", "")
	SetTranslation("wp-annotations", "default", "")
	SetTranslation("wp-api-fetch", "default", "")
	SetTranslation("wp-block-directory", "default", "")
	SetTranslation("wp-block-editor", "default", "")
	SetTranslation("wp-block-library", "default", "")
	SetTranslation("wp-blocks", "default", "")
	SetTranslation("wp-components", "default", "")
	SetTranslation("wp-core-data", "default", "")
	SetTranslation("wp-customize-widgets", "default", "")
	SetTranslation("wp-edit-post", "default", "")
	SetTranslation("wp-edit-site", "default", "")
	SetTranslation("wp-edit-widgets", "default", "")
	SetTranslation("wp-editor", "default", "")
	SetTranslation("wp-format-library", "default", "")
	SetTranslation("wp-keycodes", "default", "")
	SetTranslation("wp-list-reusable-blocks", "default", "")
	SetTranslation("wp-media-utils", "default", "")
	SetTranslation("wp-nux", "default", "")
	SetTranslation("wp-preferences", "default", "")
	SetTranslation("wp-reusable-blocks", "default", "")
	SetTranslation("wp-rich-text", "default", "")
	SetTranslation("wp-server-side-render", "default", "")
	SetTranslation("wp-widgets", "default", "")
	SetTranslation("common", "default", "")
	SetTranslation("wp-pointer", "default", "")
	SetTranslation("wp-auth-check", "default", "")
	SetTranslation("wp-theme-plugin-editor", "default", "")
	SetTranslation("password-strength-meter", "default", "")
	SetTranslation("application-passwords", "default", "")
	SetTranslation("auth-app", "default", "")
	SetTranslation("user-profile", "default", "")
	SetTranslation("media-views", "default", "")
	SetTranslation("media-editor", "default", "")
	SetTranslation("wp-a11y", "default", "")
	SetTranslation("wp-annotations", "default", "")
	SetTranslation("wp-api-fetch", "default", "")
	SetTranslation("wp-block-directory", "default", "")
	SetTranslation("wp-block-editor", "default", "")
	SetTranslation("wp-block-library", "default", "")
	SetTranslation("wp-blocks", "default", "")
	SetTranslation("wp-components", "default", "")
	SetTranslation("wp-core-data", "default", "")
	SetTranslation("wp-customize-widgets", "default", "")
	SetTranslation("wp-edit-post", "default", "")
	SetTranslation("wp-edit-site", "default", "")
	SetTranslation("wp-edit-widgets", "default", "")
	SetTranslation("wp-editor", "default", "")
	SetTranslation("wp-format-library", "default", "")
	SetTranslation("wp-keycodes", "default", "")
	SetTranslation("wp-list-reusable-blocks", "default", "")
	SetTranslation("wp-media-utils", "default", "")
	SetTranslation("wp-nux", "default", "")
	SetTranslation("wp-preferences", "default", "")
	SetTranslation("wp-reusable-blocks", "default", "")
	SetTranslation("wp-rich-text", "default", "")
	SetTranslation("wp-server-side-render", "default", "")
	SetTranslation("wp-widgets", "default", "")
}

func defaultAddData() {
	AddData("json2", "conditional", `lt IE 8`)
	AddData("wp-embed-template-ie", "conditional", `lte IE 8`)
	AddData("wp-block-library-theme", "path", `wp-includes/css/dist/block-library/theme.min.css`)
	AddData("wp-block-editor", "path", `/wp-includes/css/dist/block-editor/style.min.css`)
	AddData("wp-block-library", "path", `/wp-includes/css/dist/block-library/style.min.css`)
	AddData("wp-block-directory", "path", `/wp-includes/css/dist/block-directory/style.min.css`)
	AddData("wp-components", "path", `/wp-includes/css/dist/components/style.min.css`)
	AddData("wp-edit-post", "path", `/wp-includes/css/dist/edit-post/style.min.css`)
	AddData("wp-editor", "path", `/wp-includes/css/dist/editor/style.min.css`)
	AddData("wp-format-library", "path", `/wp-includes/css/dist/format-library/style.min.css`)
	AddData("wp-list-reusable-blocks", "path", `/wp-includes/css/dist/list-reusable-blocks/style.min.css`)
	AddData("wp-reusable-blocks", "path", `/wp-includes/css/dist/reusable-blocks/style.min.css`)
	AddData("wp-nux", "path", `/wp-includes/css/dist/nux/style.min.css`)
	AddData("wp-widgets", "path", `/wp-includes/css/dist/widgets/style.min.css`)
	AddData("wp-edit-widgets", "path", `/wp-includes/css/dist/edit-widgets/style.min.css`)
	AddData("wp-customize-widgets", "path", `/wp-includes/css/dist/customize-widgets/style.min.css`)
	AddData("wp-edit-site", "path", `/wp-includes/css/dist/edit-site/style.min.css`)
	AddData("common", "rtl", `replace`)
	AddData("common", "suffix", `.min`)
	AddData("forms", "rtl", `replace`)
	AddData("forms", "suffix", `.min`)
	AddData("admin-menu", "rtl", `replace`)
	AddData("admin-menu", "suffix", `.min`)
	AddData("dashboard", "rtl", `replace`)
	AddData("dashboard", "suffix", `.min`)
	AddData("list-tables", "rtl", `replace`)
	AddData("list-tables", "suffix", `.min`)
	AddData("edit", "rtl", `replace`)
	AddData("edit", "suffix", `.min`)
	AddData("revisions", "rtl", `replace`)
	AddData("revisions", "suffix", `.min`)
	AddData("media", "rtl", `replace`)
	AddData("media", "suffix", `.min`)
	AddData("themes", "rtl", `replace`)
	AddData("themes", "suffix", `.min`)
	AddData("about", "rtl", `replace`)
	AddData("about", "suffix", `.min`)
	AddData("nav-menus", "rtl", `replace`)
	AddData("nav-menus", "suffix", `.min`)
	AddData("widgets", "rtl", `replace`)
	AddData("widgets", "suffix", `.min`)
	AddData("site-icon", "rtl", `replace`)
	AddData("site-icon", "suffix", `.min`)
	AddData("l10n", "rtl", `replace`)
	AddData("l10n", "suffix", `.min`)
	AddData("install", "rtl", `replace`)
	AddData("install", "suffix", `.min`)
	AddData("wp-color-picker", "rtl", `replace`)
	AddData("wp-color-picker", "suffix", `.min`)
	AddData("customize-controls", "rtl", `replace`)
	AddData("customize-controls", "suffix", `.min`)
	AddData("customize-widgets", "rtl", `replace`)
	AddData("customize-widgets", "suffix", `.min`)
	AddData("customize-nav-menus", "rtl", `replace`)
	AddData("customize-nav-menus", "suffix", `.min`)
	AddData("customize-preview", "rtl", `replace`)
	AddData("customize-preview", "suffix", `.min`)
	AddData("login", "rtl", `replace`)
	AddData("login", "suffix", `.min`)
	AddData("site-health", "rtl", `replace`)
	AddData("site-health", "suffix", `.min`)
	AddData("buttons", "rtl", `replace`)
	AddData("buttons", "suffix", `.min`)
	AddData("admin-bar", "rtl", `replace`)
	AddData("admin-bar", "suffix", `.min`)
	AddData("wp-auth-check", "rtl", `replace`)
	AddData("wp-auth-check", "suffix", `.min`)
	AddData("editor-buttons", "rtl", `replace`)
	AddData("editor-buttons", "suffix", `.min`)
	AddData("media-views", "rtl", `replace`)
	AddData("media-views", "suffix", `.min`)
	AddData("wp-pointer", "rtl", `replace`)
	AddData("wp-pointer", "suffix", `.min`)
	AddData("wp-jquery-ui-dialog", "rtl", `replace`)
	AddData("wp-jquery-ui-dialog", "suffix", `.min`)
	AddData("wp-reset-editor-styles", "rtl", `replace`)
	AddData("wp-reset-editor-styles", "suffix", `.min`)
	AddData("wp-editor-classic-layout-styles", "rtl", `replace`)
	AddData("wp-editor-classic-layout-styles", "suffix", `.min`)
	AddData("wp-block-library-theme", "rtl", `replace`)
	AddData("wp-block-library-theme", "suffix", `.min`)
	AddData("wp-edit-blocks", "rtl", `replace`)
	AddData("wp-edit-blocks", "suffix", `.min`)
	AddData("wp-block-editor", "rtl", `replace`)
	AddData("wp-block-editor", "suffix", `.min`)
	AddData("wp-block-library", "rtl", `replace`)
	AddData("wp-block-library", "suffix", `.min`)
	AddData("wp-block-directory", "rtl", `replace`)
	AddData("wp-block-directory", "suffix", `.min`)
	AddData("wp-components", "rtl", `replace`)
	AddData("wp-components", "suffix", `.min`)
	AddData("wp-customize-widgets", "rtl", `replace`)
	AddData("wp-customize-widgets", "suffix", `.min`)
	AddData("wp-edit-post", "rtl", `replace`)
	AddData("wp-edit-post", "suffix", `.min`)
	AddData("wp-edit-site", "rtl", `replace`)
	AddData("wp-edit-site", "suffix", `.min`)
	AddData("wp-edit-widgets", "rtl", `replace`)
	AddData("wp-edit-widgets", "suffix", `.min`)
	AddData("wp-editor", "rtl", `replace`)
	AddData("wp-editor", "suffix", `.min`)
	AddData("wp-format-library", "rtl", `replace`)
	AddData("wp-format-library", "suffix", `.min`)
	AddData("wp-list-reusable-blocks", "rtl", `replace`)
	AddData("wp-list-reusable-blocks", "suffix", `.min`)
	AddData("wp-reusable-blocks", "rtl", `replace`)
	AddData("wp-reusable-blocks", "suffix", `.min`)
	AddData("wp-nux", "rtl", `replace`)
	AddData("wp-nux", "suffix", `.min`)
	AddData("wp-widgets", "rtl", `replace`)
	AddData("wp-widgets", "suffix", `.min`)
	AddData("deprecated-media", "rtl", `replace`)
	AddData("deprecated-media", "suffix", `.min`)
	AddData("farbtastic", "rtl", `replace`)
	AddData("farbtastic", "suffix", `.min`)

}

func defaultAddInLineScript() {
	AddInlineScript("mediaelement-core", `var mejsL10n = {"language":"zh","strings":{"mejs.download-file":"\u4e0b\u8f7d\u6587\u4ef6","mejs.install-flash":"\u60a8\u6b63\u5728\u4f7f\u7528\u7684\u6d4f\u89c8\u5668\u672a\u5b89\u88c5\u6216\u542f\u7528Flash\u64ad\u653e\u5668\uff0c\u8bf7\u542f\u7528\u60a8\u7684Flash\u64ad\u653e\u5668\u63d2\u4ef6\uff0c\u6216\u4ece https:\/\/get.adobe.com\/flashplayer\/ \u4e0b\u8f7d\u6700\u65b0\u7248\u3002","mejs.fullscreen":"\u5168\u5c4f","mejs.play":"\u64ad\u653e","mejs.pause":"\u6682\u505c","mejs.time-slider":"\u65f6\u95f4\u8f74","mejs.time-help-text":"\u4f7f\u7528\u5de6\/\u53f3\u7bad\u5934\u952e\u6765\u524d\u8fdb\u4e00\u79d2\uff0c\u4e0a\/\u4e0b\u7bad\u5934\u952e\u6765\u524d\u8fdb\u5341\u79d2\u3002","mejs.live-broadcast":"\u73b0\u573a\u76f4\u64ad","mejs.volume-help-text":"\u4f7f\u7528\u4e0a\/\u4e0b\u7bad\u5934\u952e\u6765\u589e\u9ad8\u6216\u964d\u4f4e\u97f3\u91cf\u3002","mejs.unmute":"\u53d6\u6d88\u9759\u97f3","mejs.mute":"\u9759\u97f3","mejs.volume-slider":"\u97f3\u91cf","mejs.video-player":"\u89c6\u9891\u64ad\u653e\u5668","mejs.audio-player":"\u97f3\u9891\u64ad\u653e\u5668","mejs.captions-subtitles":"\u8bf4\u660e\u6587\u5b57\u6216\u5b57\u5e55","mejs.captions-chapters":"\u7ae0\u8282","mejs.none":"\u65e0","mejs.afrikaans":"\u5357\u975e\u8377\u5170\u8bed","mejs.albanian":"\u963f\u5c14\u5df4\u5c3c\u4e9a\u8bed","mejs.arabic":"\u963f\u62c9\u4f2f\u8bed","mejs.belarusian":"\u767d\u4fc4\u7f57\u65af\u8bed","mejs.bulgarian":"\u4fdd\u52a0\u5229\u4e9a\u8bed","mejs.catalan":"\u52a0\u6cf0\u7f57\u5c3c\u4e9a\u8bed","mejs.chinese":"\u4e2d\u6587","mejs.chinese-simplified":"\u4e2d\u6587\uff08\u7b80\u4f53\uff09","mejs.chinese-traditional":"\u4e2d\u6587(\uff08\u7e41\u4f53\uff09","mejs.croatian":"\u514b\u7f57\u5730\u4e9a\u8bed","mejs.czech":"\u6377\u514b\u8bed","mejs.danish":"\u4e39\u9ea6\u8bed","mejs.dutch":"\u8377\u5170\u8bed","mejs.english":"\u82f1\u8bed","mejs.estonian":"\u7231\u6c99\u5c3c\u4e9a\u8bed","mejs.filipino":"\u83f2\u5f8b\u5bbe\u8bed","mejs.finnish":"\u82ac\u5170\u8bed","mejs.french":"\u6cd5\u8bed","mejs.galician":"\u52a0\u5229\u897f\u4e9a\u8bed","mejs.german":"\u5fb7\u8bed","mejs.greek":"\u5e0c\u814a\u8bed","mejs.haitian-creole":"\u6d77\u5730\u514b\u91cc\u5965\u5c14\u8bed","mejs.hebrew":"\u5e0c\u4f2f\u6765\u8bed","mejs.hindi":"\u5370\u5730\u8bed","mejs.hungarian":"\u5308\u7259\u5229\u8bed","mejs.icelandic":"\u51b0\u5c9b\u8bed","mejs.indonesian":"\u5370\u5ea6\u5c3c\u897f\u4e9a\u8bed","mejs.irish":"\u7231\u5c14\u5170\u8bed","mejs.italian":"\u610f\u5927\u5229\u8bed","mejs.japanese":"\u65e5\u8bed","mejs.korean":"\u97e9\u8bed","mejs.latvian":"\u62c9\u8131\u7ef4\u4e9a\u8bed","mejs.lithuanian":"\u7acb\u9676\u5b9b\u8bed","mejs.macedonian":"\u9a6c\u5176\u987f\u8bed","mejs.malay":"\u9a6c\u6765\u8bed","mejs.maltese":"\u9a6c\u8033\u4ed6\u8bed","mejs.norwegian":"\u632a\u5a01\u8bed","mejs.persian":"\u6ce2\u65af\u8bed","mejs.polish":"\u6ce2\u5170\u8bed","mejs.portuguese":"\u8461\u8404\u7259\u8bed","mejs.romanian":"\u7f57\u9a6c\u5c3c\u4e9a\u8bed","mejs.russian":"\u4fc4\u8bed","mejs.serbian":"\u585e\u5c14\u7ef4\u4e9a\u8bed","mejs.slovak":"\u65af\u6d1b\u4f10\u514b\u8bed","mejs.slovenian":"\u65af\u6d1b\u6587\u5c3c\u4e9a\u8bed","mejs.spanish":"\u897f\u73ed\u7259\u8bed","mejs.swahili":"\u65af\u74e6\u5e0c\u91cc\u8bed","mejs.swedish":"\u745e\u5178\u8bed","mejs.tagalog":"\u4ed6\u52a0\u7984\u8bed","mejs.thai":"\u6cf0\u8bed","mejs.turkish":"\u571f\u8033\u5176\u8bed","mejs.ukrainian":"\u4e4c\u514b\u5170\u8bed","mejs.vietnamese":"\u8d8a\u5357\u8bed","mejs.welsh":"\u5a01\u5c14\u58eb\u8bed","mejs.yiddish":"\u610f\u7b2c\u7eea\u8bed"}};`, "before")
	AddInlineScript("lodash", `window.lodash = _.noConflict();`, "after")
	AddInlineScript("moment", `moment.updateLocale( 'zh_CN', {"months":["1\u6708","2\u6708","3\u6708","4\u6708","5\u6708","6\u6708","7\u6708","8\u6708","9\u6708","10\u6708","11\u6708","12\u6708"],"monthsShort":["1\u6708","2\u6708","3\u6708","4\u6708","5\u6708","6\u6708","7\u6708","8\u6708","9\u6708","10\u6708","11\u6708","12\u6708"],"weekdays":["\u661f\u671f\u65e5","\u661f\u671f\u4e00","\u661f\u671f\u4e8c","\u661f\u671f\u4e09","\u661f\u671f\u56db","\u661f\u671f\u4e94","\u661f\u671f\u516d"],"weekdaysShort":["\u5468\u65e5","\u5468\u4e00","\u5468\u4e8c","\u5468\u4e09","\u5468\u56db","\u5468\u4e94","\u5468\u516d"],"week":{"dow":1},"longDateFormat":{"LT":"ag:i","LTS":null,"L":null,"LL":"Y\u5e74n\u6708j\u65e5","LLL":"Y\u5e74n\u6708j\u65e5a g:i","LLLL":null}} );`, "after")
	AddInlineScript("wp-i18n", `wp.i18n.setLocaleData( { 'text direction\u0004ltr': [ 'ltr' ] } );`, "after")
	AddInlineScript("text-widgets", `wp.textWidgets.idBases.push( "text" );`, "after")
	AddInlineScript("custom-html-widgets", `wp.customHtmlWidgets.idBases.push( "custom_html" );`, "after")
}

func defaultAddInLineStyle() {
	AddInlineStyle("global-styles", GetGlobalStyletSheet())
	AddInlineStyle("global-styles", `.wp-block-navigation a:where(:not(.wp-element-button)){color: inherit;}`)
	AddInlineStyle("global-styles", `:where(.wp-block-columns.is-layout-flex){gap: 2em;}`)
	AddInlineStyle("global-styles", `.wp-block-pullquote{font-size: 1.5em;line-height: 1.6;}`)
}

var re = regexp.MustCompile(`(?is:\([A-Za-z0-9'.:\-/, ]+\))`)
var rea = regexp.MustCompile(`array\(array\(.*?\)\)`)

func InitDefaultScriptSetting() {
	initThemeJson()
	defaultLocalize()
	defaultTranslate()
	defaultAddData()
	defaultAddInLineScript()
	defaultAddInLineStyle()
}

func initThemeJson() {
	blocksData := __blocksData()
	path := config.GetConfig().WpDir
	f, err := os.ReadFile(filepath.Join(path, "wp-includes/theme.json"))
	if err != nil {
		logs.Error(err, "can't open theme json", path)
		return
	}

	var j map[string]any
	err = json.Unmarshal(f, &j)
	if err != nil {
		logs.Error(err, "can't parse theme json")
		return
	}
	t := ThemeJson{blocksData, j}
	setThemeJson(j)
	setSpacingSizes(t)
	__themeJson.Store(t)
}

func setThemeJson(m map[string]any) {
	blocks, ok := maps.GetStrAnyVal[map[string]any](m, "settings.blocks")
	if !ok {
		return
	}
	var paths = [][]string{{"settings"}}
	for name := range blocks {
		paths = append(paths, []string{"settings", "blocks", name})
	}
	for _, path := range paths {
		for _, metadatum := range presetsMetadata {
			pathx := append(path, metadatum.path...)
			key := strings.Join(pathx, ".")
			preset, ok := maps.GetStrAnyVal[[]any](m, key)
			if !ok || len(preset) < 1 {
				continue
			}
			var presets []map[string]string
			for _, a := range preset {
				v, ok := a.(map[string]any)
				if !ok {
					continue
				}
				mm := map[string]string{}
				for kk, vv := range v {
					val, ok := vv.(string)
					if !ok {
						continue
					}
					mm[kk] = val
				}
				presets = append(presets, mm)
			}
			maps.SetStrAnyVal(m, key, map[string]any{
				"default": presets,
			})
		}
	}
}

var __propertyMappings = map[string]string{
	"apiVersion":      "api_version",
	"title":           "title",
	"category":        "category",
	"parent":          "parent",
	"ancestor":        "ancestor",
	"icon":            "icon",
	"description":     "description",
	"keywords":        "keywords",
	"attributes":      "attributes",
	"providesContext": "provides_context",
	"usesContext":     "uses_context",
	"supports":        "supports",
	"styles":          "styles",
	"variations":      "variations",
	"example":         "example",
}

func __propertyMap(m map[string]any) {
	for k, mappedKey := range __propertyMappings {
		vv, ok := m[k]
		if ok {
			m[mappedKey] = vv
		}
	}
}

func __blocksData() map[string]any {
	path := config.GetConfig().WpDir
	//path := "/var/www/html/wordpress"
	b, err := os.ReadFile(filepath.Join(path, "wp-includes/blocks/blocks-json.php"))
	if err != nil {
		logs.Error(err, "can't open block json", path)
		return nil
	}
	bb := re.ReplaceAllStringFunc(string(b), func(s string) string {
		return str.Replace(s, map[string]string{
			"(": "[",
			")": "]",
		})
	})
	bb = strings.ReplaceAll(bb, "\"", `\"`)
	bb = rea.ReplaceAllStringFunc(bb, func(s string) string {
		s = strings.ReplaceAll(s, "array(array", "[")
		ss := []rune(s)
		ss[len(ss)-1] = ']'
		s = string(ss)
		return s
	})
	bb = str.Replace(bb, map[string]string{
		"<?php":  "",
		"return": "",
		"array":  "",
		"()":     "[]",
		"(":      "{",
		")":      "}",
		"=>":     ":",
		";":      "",
		"'":      `"`,
	})

	var blocks map[string]any
	err = json.Unmarshal([]byte(bb), &blocks)
	if err != nil {
		logs.Error(err, "can't parse block json")
		return nil
	}
	c := map[string]any{
		"version": int64(2),
	}
	for k, v := range blocks {
		m, ok := v.(map[string]any)
		if !ok {
			continue
		}
		_, ok = m["style"]
		if !ok {
			m["style"] = str.Join("wp-block-", k)
		}
		_, ok = m["editorStyle"]
		if !ok {
			m["editorStyle"] = str.Join("wp-block-", k, "-editor")
		}
		__propertyMap(m)
		name := maps.GetStrAnyValWithDefaults(m, "name", str.Join("core/", k))
		blocks[name] = v
		if name != k {
			delete(blocks, k)
		}
		__blockSelectors(m)
		s, ok := maps.GetStrAnyVal[map[string]any](m, "supports.__experimentalStyle")
		if ok {
			__removeComment(s)
			maps.SetStrAnyVal(c, str.Join("styles.blocks.", name), s)
		}
		_, ok = maps.GetStrAnyVal[string](m, "supports.spacing.blockGap.__experimentalDefault")
		if ok {
			key := str.Join("styles.blocks.", name, ".spacing.blockGap")
			_, ok := maps.GetStrAnyVal[string](c, key)
			if !ok {
				maps.SetStrAnyVal[map[string]any](c, key, nil)
			}
		}
	}
	return map[string]any{
		"blocks_metadata": blocks,
		"theme_json":      c,
	}
}

func __blockSelectors(m map[string]any) {
	selector, ok := maps.GetStrAnyVal[string](m, "supports.__experimentalSelector")
	if !ok {
		vv, _ := maps.GetStrAnyVal[string](m, "name")
		selector = str.Join(".wp-block-", str.Replaces(vv,
			[]string{"core/", ""},
			[]string{"/", "-"},
		))
	}
	var features = map[string]string{}
	maps.SetStrAnyVal(m, "supports.selector", selector)
	color, ok := maps.GetStrAnyVal[string](m, "supports.color.__experimentalDuotone")
	if ok {
		maps.SetStrAnyVal(m, "duotone", color)
	}
	for k, v := range blockSupportFeatureLevelSelectors {
		key := str.Join("supports.", k, ".__experimentalSelector")
		vv, ok := maps.GetStrAnyVal[string](m, key)
		if ok && vv != "" {
			features[v] = scopeSelector(selector, vv)
		}
	}
	if len(features) > 0 {
		m["features"] = features
	}
	blockSelector := strings.Split(selector, ",")
	for name, selor := range __elements {
		var els []string
		for _, s := range blockSelector {
			if s == selor {
				els = append(els, selor)
				break
			}
			els = append(els, appendToSelector(selor, str.Join(s, " "), "left"))
		}
		maps.SetStrAnyVal(m, str.Join("elements.", name), strings.Join(els, ","))
	}
	styles, ok := maps.GetStrAnyVal[[]any](m, "styles")
	if ok {
		var styleSelectors = map[string]string{}
		for _, ss := range styles {
			s, ok := ss.(map[string]any)
			if !ok {
				continue
			}
			name, ok := maps.GetStrAnyVal[string](s, "name")
			if !ok {
				continue
			}
			styleSelectors[name] = appendToSelector(str.Join(".is-style-", name, ".is-style-", name), selector, "")
		}
		m["styleVariations"] = styleSelectors
	}
}
