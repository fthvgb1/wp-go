package wp

import (
	"github.com/fthvgb1/wp-go/app/pkg/constraints"
	"github.com/fthvgb1/wp-go/cache/reload"
	"github.com/fthvgb1/wp-go/helper/maps"
	"github.com/fthvgb1/wp-go/helper/slice"
	str "github.com/fthvgb1/wp-go/helper/strings"
	"github.com/fthvgb1/wp-go/safety"
	"strings"
)

var handleComponents = safety.NewMap[string, map[string][]Components[string]]()
var handleComponentHook = safety.NewMap[string, map[string][]func(Components[string]) (Components[string], bool)]()

var componentsArgs = safety.NewMap[string, any]()
var componentFilterFns = safety.NewMap[string, []func(*Handle, string, ...any) string]()

func (h *Handle) DeleteComponents(scene, componentKey, componentName string) {
	v, ok := handleComponentHook.Load(scene)
	if !ok {
		v = make(map[string][]func(Components[string]) (Components[string], bool))
	}

	v[componentKey] = append(v[componentKey], func(c Components[string]) (Components[string], bool) {
		return c, c.Name != componentName
	})
	handleComponentHook.Store(scene, v)
}
func (h *Handle) ReplaceComponents(scene, componentKey, componentName string, components Components[string]) {
	v, ok := handleComponentHook.Load(scene)
	if !ok {
		v = make(map[string][]func(Components[string]) (Components[string], bool))
	}
	v[componentKey] = append(v[componentKey], func(c Components[string]) (Components[string], bool) {
		if c.Name == componentName {
			c = components
		}
		return c, true
	})
	handleComponentHook.Store(scene, v)
}
func (h *Handle) PushComponentHooks(scene, componentKey string, fn func(Components[string]) (Components[string], bool)) {
	v, ok := handleComponentHook.Load(scene)
	if !ok {
		v = make(map[string][]func(Components[string]) (Components[string], bool))
	}
	v[componentKey] = append(v[componentKey], fn)
	handleComponentHook.Store(scene, v)
}

var GetAndHookComponents = reload.BuildMapFnWithAnyParams[string]("calComponents", HookComponent)

func HookComponent(a ...any) []Components[string] {
	componentKey := a[0].(string)
	scene := a[1].(string)
	mut := reload.GetGlobeMutex()
	mut.Lock()
	defer mut.Unlock()
	components := GetComponents(scene, componentKey)
	r := slice.FilterAndMap(components, func(component Components[string]) (Components[string], bool) {
		keyHooks, ok := handleComponentHook.Load(scene)
		if !ok {
			return component, true
		}
		hooks := keyHooks[componentKey]
		for _, fn := range hooks {
			hookedComponent, ok := fn(component)
			if !ok { // DeleteComponents fn
				return hookedComponent, false
			}
			component = hookedComponent // ReplaceComponents fn
		}
		return component, true
	})
	slice.SimpleSort(r, slice.DESC, func(t Components[string]) float64 {
		return t.Order
	})
	return r
}

func GetComponents(scene, key string) (r []Components[string]) {
	sceneComponents, _ := handleComponents.Load(scene)
	allSceneComponents, _ := handleComponents.Load(constraints.AllScene)
	r = append(sceneComponents[key], allSceneComponents[key]...)
	return
}

var CacheComponent = reload.BuildMapFnWithAnyParams[string]("cacheComponents", cacheComponentFn)

func cacheComponentFn(a ...any) string {
	return a[0].(Components[string]).Fn(a[1].(*Handle))
}

func CalComponent(h *Handle) func(string) string {
	return func(componentKey string) string {
		cacheKey := str.Join("get-hook-Components-", h.scene, "-", componentKey)
		hookedComponents := GetAndHookComponents(cacheKey, componentKey, h.scene)
		var s = make([]string, 0, len(hookedComponents))
		for _, component := range hookedComponents {
			if component.Val != "" {
				s = append(s, component.Val)
				continue
			}
			if component.Fn != nil {
				v := ""
				if component.Cached {
					key := str.Join(h.scene, "-", componentKey, "-", component.Name)
					v = CacheComponent(key, component, h)
				} else {
					v = component.Fn(h)
				}
				if v != "" {
					s = append(s, v)
				}
			}
		}
		return strings.Join(s, "\n")
	}
}

func (h *Handle) PushComponents(scene, componentType string, components ...Components[string]) {
	c, ok := handleComponents.Load(scene)
	if !ok {
		c = make(map[string][]Components[string])
	}
	c[componentType] = append(c[componentType], components...)
	handleComponents.Store(scene, c)
}

func (h *Handle) PushGroupComponentStr(scene, componentType, name string, order float64, strs ...string) {
	var component = Components[string]{
		Val: strings.Join(slice.FilterAndMap(strs, func(t string) (string, bool) {
			t = strings.Trim(t, " \n\r\t\v\x00")
			if t == "" {
				return "", false
			}
			return t, true
		}), "\n"),
		Order: order,
		Name:  name,
	}
	h.PushComponents(scene, componentType, component)
}

func (h *Handle) PushCacheGroupHeadScript(scene, name string, order float64, fns ...func(*Handle) string) {
	h.PushGroupCacheComponentFn(scene, constraints.HeadScript, name, order, fns...)
}

func (h *Handle) PushFooterScript(scene string, components ...Components[string]) {
	h.PushComponents(scene, constraints.FooterScript, components...)
}

func (h *Handle) PushGroupFooterScript(scene, name string, order float64, strs ...string) {
	h.PushGroupComponentStr(scene, constraints.FooterScript, name, order, strs...)
}

func (h *Handle) PushCacheGroupFooterScript(scene, name string, order float64, fns ...func(*Handle) string) {
	h.PushGroupCacheComponentFn(scene, constraints.FooterScript, name, order, fns...)
}
func (h *Handle) PushGroupCacheComponentFn(scene, componentType, name string, order float64, fns ...func(*Handle) string) {
	h.PushComponents(scene, componentType, NewComponent(name, "", true, order, func(h *Handle) string {
		return strings.Join(slice.Map(fns, func(t func(*Handle) string) string {
			return t(h)
		}), "\n")
	}))
}

func NewComponent(name, val string, cached bool, order float64, fn func(handle *Handle) string) Components[string] {
	return Components[string]{Fn: fn, Name: name, Cached: cached, Order: order, Val: val}
}

func (h *Handle) AddCacheComponent(scene, componentType, name string, order float64, fn func(*Handle) string) {
	h.PushComponents(scene, componentType, NewComponent(name, "", true, order, fn))
}

func (h *Handle) PushHeadScript(scene string, components ...Components[string]) {
	h.PushComponents(scene, constraints.HeadScript, components...)
}
func (h *Handle) PushGroupHeadScript(scene, name string, order float64, str ...string) {
	h.PushGroupComponentStr(scene, constraints.HeadScript, name, order, str...)
}

func GetComponentsArgs[T any](k string, defaults T) T {
	v, ok := componentsArgs.Load(k)
	if ok {
		vv, ok := v.(T)
		if ok {
			return vv
		}
	}
	return defaults
}

func PushComponentsArgsForSlice[T any](name string, v ...T) {
	val, ok := componentsArgs.Load(name)
	if !ok {
		var vv []T
		vv = append(vv, v...)
		componentsArgs.Store(name, vv)
		return
	}
	vv, ok := val.([]T)
	if ok {
		vv = append(vv, v...)
		componentsArgs.Store(name, vv)
	}
}
func SetComponentsArgsForMap[K comparable, V any](name string, key K, v V) {
	val, ok := componentsArgs.Load(name)
	if !ok {
		vv := make(map[K]V)
		vv[key] = v
		componentsArgs.Store(name, vv)
		return
	}
	vv, ok := val.(map[K]V)
	if ok {
		vv[key] = v
		componentsArgs.Store(name, vv)
	}
}
func MergeComponentsArgsForMap[K comparable, V any](name string, m map[K]V) {
	val, ok := componentsArgs.Load(name)
	if !ok {
		componentsArgs.Store(name, m)
		return
	}
	vv, ok := val.(map[K]V)
	if ok {
		componentsArgs.Store(name, maps.Merge(vv, m))
	}
}

func SetComponentsArgs(key string, value any) {
	componentsArgs.Store(key, value)
}

func (h *Handle) GetComponentFilterFn(name string) ([]func(*Handle, string, ...any) string, bool) {
	return componentFilterFns.Load(name)
}

func (h *Handle) AddActionFilter(name string, fns ...func(*Handle, string, ...any) string) {
	v, _ := componentFilterFns.Load(name)
	v = append(v, fns...)
	componentFilterFns.Store(name, v)
}
func (h *Handle) DoActionFilter(name, s string, args ...any) string {
	calls, ok := componentFilterFns.Load(name)
	if ok {
		return slice.Reduce(calls, func(fn func(*Handle, string, ...any) string, r string) string {
			return fn(h, r, args...)
		}, s)
	}
	return s
}
