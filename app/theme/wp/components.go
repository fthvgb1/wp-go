package wp

import (
	"github.com/fthvgb1/wp-go/app/cmd/reload"
	"github.com/fthvgb1/wp-go/app/pkg/constraints"
	"github.com/fthvgb1/wp-go/helper/maps"
	"github.com/fthvgb1/wp-go/helper/slice"
	str "github.com/fthvgb1/wp-go/helper/strings"
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
	componentss := reload.GetAnyValMapBy("scene-components", str.Join("allScene-", h.scene), h, func(h *Handle) map[string][]Components[string] {
		return maps.MergeBy(func(k string, v1, v2 []Components[string]) ([]Components[string], bool) {
			vv := append(v1, v2...)
			return vv, vv != nil
		}, nil, h.components[h.scene], h.components[constraints.AllScene])
	})
	for k, components := range componentss {
		key := str.Join("calComponents-", h.scene, "-", k)
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

func (h *Handle) PushComponents(scene, componentType string, components ...Components[string]) {
	c, ok := h.components[scene]
	if !ok {
		c = make(map[string][]Components[string])
		h.components[scene] = c
	}
	c[componentType] = append(c[componentType], components...)
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

func SetComponentsArgs(h *Handle, key string, value any) {
	h.componentsArgs[key] = value
}

func (h *Handle) ComponentFilterFn(name string) ([]func(*Handle, string, ...any) string, bool) {
	fn, ok := h.componentFilterFn[name]
	return fn, ok
}

func (h *Handle) AddActionFilter(name string, fns ...func(*Handle, string, ...any) string) {
	h.componentFilterFn[name] = append(h.componentFilterFn[name], fns...)
}
func (h *Handle) DoActionFilter(name, s string, args ...any) string {
	calls, ok := h.componentFilterFn[name]
	if ok {
		return slice.Reduce(calls, func(fn func(*Handle, string, ...any) string, r string) string {
			return fn(h, r, args...)
		}, s)
	}
	return s
}
