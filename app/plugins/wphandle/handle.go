package wphandle

import (
	"errors"
	"github.com/fthvgb1/wp-go/app/pkg/config"
	"github.com/fthvgb1/wp-go/app/pkg/logs"
	"github.com/fthvgb1/wp-go/app/plugins/wphandle/enlightjs"
	"github.com/fthvgb1/wp-go/app/plugins/wphandle/hiddenlogin"
	"github.com/fthvgb1/wp-go/app/theme/wp"
	str "github.com/fthvgb1/wp-go/helper/strings"
	"github.com/fthvgb1/wp-go/safety"
	"path/filepath"
	"plugin"
)

var plugins = func() *safety.Map[string, func(*wp.Handle)] {
	m := safety.NewMap[string, func(*wp.Handle)]()
	m.Store("Enlightjs", enlightjs.EnlighterJS)
	m.Store("HiddenLogin", hiddenlogin.HiddenLogin)
	return m
}()

func RegisterPlugin(name string, fn func(*wp.Handle)) {
	plugins.Store(name, fn)
}

func UsePlugins(h *wp.Handle, calls ...string) {
	calls = append(calls, config.GetConfig().Plugins...)
	for _, call := range calls {
		call = str.FirstUpper(call)
		if fn, ok := plugins.Load(call); ok {
			fn(h)
		}
	}
}

func LoadPlugins() {
	dirPath := config.GetConfig().PluginPath
	if dirPath == "" {
		return
	}
	glob, err := filepath.Glob(filepath.Join(dirPath, "*.so"))
	if err != nil {
		logs.Error(err, "读取插件目录错误", dirPath)
		return
	}
	for _, entry := range glob {
		f := filepath.Join(dirPath, entry)
		p, err := plugin.Open(f)
		if err != nil {
			logs.Error(err, "读取插件错误", f)
			continue
		}
		name := filepath.Ext(entry)
		name = str.FirstUpper(entry[0 : len(entry)-len(name)])
		pl, err := p.Lookup(name)
		if err != nil {
			logs.Error(err, "插件lookup错误", f)
			continue
		}
		plu, ok := pl.(func(*wp.Handle))
		if !ok {
			logs.Error(errors.New("switch func(*wp.Handle) fail"), "插件转换错误", f)
			continue
		}
		RegisterPlugin(name, plu)
	}

}
