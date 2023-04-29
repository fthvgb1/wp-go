package wp

import (
	"github.com/fthvgb1/wp-go/helper/maps"
	"github.com/fthvgb1/wp-go/helper/slice"
	str "github.com/fthvgb1/wp-go/helper/strings"
	"github.com/fthvgb1/wp-go/internal/cmd/reload"
	"github.com/fthvgb1/wp-go/internal/pkg/constraints"
	"strings"
)

func (h *Handle) DeleteComponents(scene, name string) {
	h.componentHook[scene] = append(h.componentHook[scene], func(c Components[string]) (Components[string], bool) {
		return c, c.Name != name
	})
}
func (h *Handle) ReplaceComponents(scene, name string, components Components[string]) {
	h.componentHook[scene] = append(h.componentHook[scene], func(c Components[string]) (Components[string], bool) {
		if c.Name == name {
			c = components
		}
		return c, true
	})
}
func (h *Handle) HookComponents(scene string, fn func(Components[string]) (Components[string], bool)) {
	h.componentHook[scene] = append(h.componentHook[scene], fn)
}

func CalComponents(h *Handle) {
	for k, components := range h.components {
		key := str.Join("calComponents-", k)
		key = h.ComponentFilterFnHook("calComponents", key, k)
		ss := reload.GetAnyValMapBy("calComponents", key, h, func(h *Handle) []Components[string] {
			r := slice.FilterAndMap(components, func(t Components[string]) (Components[string], bool) {
				fns, ok := h.componentHook[k]
				if !ok {
					return t, true
				}
				for _, fn := range fns {
					c, ok := fn(t)
					if !ok {
						return c, false
					}
					t = c
				}
				return t, true
			})
			slice.Sort(r, func(i, j Components[string]) bool {
				return i.Order > j.Order
			})
			return r
		})
		var s = make([]string, 0, len(ss))
		for _, component := range ss {
			if component.Val != "" {
				s = append(s, component.Val)
				continue
			}
			if component.Fn != nil {
				v := ""
				if component.Cached {
					v = reload.GetAnyValMapBy("cacheComponents", component.Name, h, component.Fn)
				} else {
					v = component.Fn(h)
				}
				if v != "" {
					s = append(s, v)
				}
			}
		}
		h.ginH[k] = strings.Join(s, "\n")
	}
}

func (h *Handle) PushComponents(name string, components ...Components[string]) {
	h.components[name] = append(h.components[name], components...)
}

func (h *Handle) PushGroupComponentStr(componentType, name string, order int, strs ...string) {
	var calls []Components[string]
	for _, val := range strs {
		calls = append(calls, Components[string]{
			Val:   val,
			Order: order,
			Name:  name,
		})
	}
	h.components[componentType] = append(h.components[componentType], calls...)
}

func (h *Handle) PushCacheGroupHeadScript(key string, order int, fns ...func(*Handle) string) {
	h.PushGroupCacheComponentFn(constraints.HeadScript, key, order, fns...)
}

func (h *Handle) PushFooterScript(components ...Components[string]) {
	h.PushComponents(constraints.FooterScript, components...)
}

func (h *Handle) PushGroupFooterScript(name string, order int, strs ...string) {
	h.PushGroupComponentStr(constraints.FooterScript, name, order, strs...)
}

func (h *Handle) PushCacheGroupFooterScript(name string, order int, fns ...func(*Handle) string) {
	h.PushGroupCacheComponentFn(constraints.FooterScript, name, order, fns...)
}
func (h *Handle) PushGroupCacheComponentFn(componentType, name string, order int, fns ...func(*Handle) string) {
	h.PushComponents(componentType, h.NewComponent(name, true, order, func(h *Handle) string {
		return strings.Join(slice.Map(fns, func(t func(*Handle) string) string {
			return t(h)
		}), "\n")
	}))
}

func (h *Handle) NewComponent(name string, cached bool, order int, fn func(handle *Handle) string) Components[string] {
	return Components[string]{Fn: fn, Name: name, Cached: cached, Order: order}
}

func (h *Handle) AddCacheComponent(componentType, name string, order int, fn func(*Handle) string) {
	h.components[componentType] = append(h.components[componentType], h.NewComponent(name, true, order, fn))
}

func (h *Handle) PushHeadScript(components ...Components[string]) {
	h.PushComponents(constraints.HeadScript, components...)
}
func (h *Handle) PushGroupHeadScript(name string, order int, str ...string) {
	h.PushGroupComponentStr(constraints.HeadScript, name, order, str...)
}

func GetComponentsArgs[T any](h *Handle, k string, defaults T) T {
	v, ok := h.componentsArgs[k]
	if ok {
		vv, ok := v.(T)
		if ok {
			return vv
		}
	}
	return defaults
}

func PushComponentsArgsForSlice[T any](h *Handle, name string, v ...T) {
	val, ok := h.componentsArgs[name]
	if !ok {
		var vv []T
		vv = append(vv, v...)
		h.componentsArgs[name] = vv
		return
	}
	vv, ok := val.([]T)
	if ok {
		vv = append(vv, v...)
		h.componentsArgs[name] = vv
	}
}
func SetComponentsArgsForMap[K comparable, V any](h *Handle, name string, key K, v V) {
	val, ok := h.componentsArgs[name]
	if !ok {
		vv := make(map[K]V)
		vv[key] = v
		h.componentsArgs[name] = vv
		return
	}
	vv, ok := val.(map[K]V)
	if ok {
		vv[key] = v
		h.componentsArgs[name] = vv
	}
}
func MergeComponentsArgsForMap[K comparable, V any](h *Handle, name string, m map[K]V) {
	val, ok := h.componentsArgs[name]
	if !ok {
		h.componentsArgs[name] = m
		return
	}
	vv, ok := val.(map[K]V)
	if ok {
		h.componentsArgs[name] = maps.Merge(vv, m)
	}
}

func (h *Handle) SetComponentsArgs(key string, value any) {
	h.componentsArgs[key] = value
}

func (h *Handle) ComponentFilterFn(name string) ([]func(*Handle, string, ...any) string, bool) {
	fn, ok := h.componentFilterFn[name]
	return fn, ok
}

func (h *Handle) PushComponentFilterFn(name string, fns ...func(*Handle, string, ...any) string) {
	h.componentFilterFn[name] = append(h.componentFilterFn[name], fns...)
}
func (h *Handle) ComponentFilterFnHook(name, s string, args ...any) string {
	calls, ok := h.componentFilterFn[name]
	if ok {
		return slice.Reduce(calls, func(fn func(*Handle, string, ...any) string, r string) string {
			return fn(h, r, args...)
		}, s)
	}
	return s
}
