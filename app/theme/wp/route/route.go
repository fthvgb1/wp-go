package route

import (
	"github.com/fthvgb1/wp-go/app/theme/wp"
	"github.com/fthvgb1/wp-go/cache/reload"
	"github.com/fthvgb1/wp-go/helper/slice"
	"github.com/fthvgb1/wp-go/safety"
	"regexp"
)

type Route struct {
	Path   string
	Scene  string
	Method []string
	Type   string
}

var routeHook []func(Route) (Route, bool)

var regRoutes *safety.Map[string, *regexp.Regexp]
var routes = func() *safety.Map[string, Route] {
	r := safety.NewMap[string, Route]()
	reload.Push(func() {
		r.Flush()
		regRoutes.Flush()
	})
	regRoutes = safety.NewMap[string, *regexp.Regexp]()
	return r
}()

// PushRoute path can be const or regex string
//
//	eg: `(?P<control>\w+)/(?P<method>\w+)`, route.Route{
//			Path:   `(?P<control>\w+)/(?P<method>\w+)`,
//			Scene:  constraints.Home,
//			Method: []string{"GET"},
//			Type:   "reg",
//		}
func PushRoute(path string, route Route) error {
	if route.Type == "const" {
		routes.Store(path, route)
		return nil
	}
	re, err := regexp.Compile(route.Path)
	if err != nil {
		return err
	}
	regRoutes.Store(path, re)
	routes.Store(path, route)
	return err
}

func Delete(path string) {
	routeHook = append(routeHook, func(route Route) (Route, bool) {
		return route, route.Path != path
	})
}

func Replace(path string, route Route) {
	routeHook = append(routeHook, func(r Route) (Route, bool) {
		return route, path == route.Path
	})
}

func Hook(path string, fn func(Route) Route) {
	routeHook = append(routeHook, func(r Route) (Route, bool) {
		if path == r.Path {
			r = fn(r)
		}
		return r, path == r.Path
	})
}
func ResolveRoute(h *wp.Handle) {
	requestURI := h.C.Request.RequestURI
	rs, rrs := reload.GetAnyValBys("route",
		struct{}{},
		func(_ struct{}) (func() (map[string]Route, map[string]*regexp.Regexp), bool) {
			m := map[string]Route{}
			rrs := map[string]*regexp.Regexp{}
			routes.Range(func(key string, value Route) bool {
				vv, _ := regRoutes.Load(key)
				if len(routeHook) > 0 {
					for _, fn := range routeHook {
						v, ok := fn(value)
						if !ok {
							continue
						}
						m[v.Path] = v
						if v.Type != "reg" {
							continue
						}
						if v.Path != key {
							vvv, err := regexp.Compile(v.Path)
							if err != nil {
								panic(err)
							}
							vv = vvv
						}
						rrs[v.Path] = vv
					}
				} else {
					m[key] = value
					rrs[key] = vv
				}

				return true
			})
			return func() (map[string]Route, map[string]*regexp.Regexp) {
				return m, rrs
			}, true
		})()
	v, ok := rs[requestURI]
	if ok && slice.IsContained(v.Method, h.C.Request.Method) {
		h.SetScene(v.Scene)
		wp.Run(h, nil)
		return
	}
	for path, reg := range rrs {
		r := reg.FindAllStringSubmatch(requestURI, -1)
		if len(r) < 1 {
			return
		}
		rr := rs[path]
		if slice.IsContained(rr.Method, h.C.Request.Method) {
			h.SetScene(rr.Scene)
			h.C.Set("route", r)
			wp.Run(h, nil)
			return
		}
	}
}
