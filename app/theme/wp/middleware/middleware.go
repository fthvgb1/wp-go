package middleware

import (
	"github.com/fthvgb1/wp-go/app/pkg/cache"
	"github.com/fthvgb1/wp-go/app/pkg/constraints"
	"github.com/fthvgb1/wp-go/app/theme/wp"
	"github.com/fthvgb1/wp-go/app/theme/wp/components/widget"
	"github.com/fthvgb1/wp-go/cache/reload"
	"github.com/fthvgb1/wp-go/helper/maps"
	"github.com/fthvgb1/wp-go/helper/number"
	str "github.com/fthvgb1/wp-go/helper/strings"
	"net/http"
	"strconv"
	"time"
)

var plainRouteParam = reload.Vars([]Plain{
	{
		Action: "p",
		Param: map[string]string{
			"p":     "id",
			"cpage": "page",
		},
		Scene: constraints.Detail,
	},
	{
		Action: "s",
		Scene:  constraints.Search,
	},
	{
		Scene: constraints.Category,
		Fn: func(h *wp.Handle) bool {
			c, ok := widget.IsCategory(h)
			if !ok {
				return false
			}
			h.C.AddParam("category", c.Name)
			h.C.AddParam("page", h.C.Query("paged"))
			return true
		},
	},
	{
		Scene:  constraints.Tag,
		Action: "tag",
		Param: map[string]string{
			"tag":   "tag",
			"paged": "page",
		},
	},
	{
		Scene: constraints.Archive,
		Fn: func(h *wp.Handle) bool {
			m := h.C.Query("m")
			if m == "" {
				return false
			}
			t, err := time.Parse("200601", m)
			if err != nil {
				return false
			}
			h.C.AddParam("year", strconv.Itoa(t.Year()))
			h.C.AddParam("month", number.IntToString(t.Month()))
			h.C.AddParam("page", h.C.Query("paged"))
			return true
		},
	},
	{
		Scene: constraints.Author,
		Fn: func(h *wp.Handle) bool {
			u := h.C.Query("author")
			if u == "" {
				return false
			}
			users := reload.GetAnyValBys("usersIds", struct{}{},
				func(_ struct{}) (map[uint64]string, bool) {
					users, err := cache.GetAllUsername(h.C)
					if err != nil {
						return nil, true
					}
					return maps.Flip(users), true
				})
			name, ok := users[str.ToInteger[uint64](u, 0)]
			if !ok {
				return false
			}
			h.C.AddParam("author", name)
			h.C.AddParam("page", h.C.Query("paged"))
			return true
		},
	},
})

func SetExplainRouteParam(p []Plain) {
	plainRouteParam.Store(p)
}
func GetExplainRouteParam() []Plain {
	return plainRouteParam.Load()
}
func PushExplainRouteParam(explain ...Plain) {
	v := plainRouteParam.Load()
	v = append(v, explain...)
	plainRouteParam.Store(v)
}

type Plain struct {
	Action string
	Param  map[string]string
	Scene  string
	Fn     func(h *wp.Handle) bool
}

func MixWithPlain(h *wp.Handle) {
	for _, explain := range plainRouteParam.Load() {
		if explain.Action == "" && explain.Fn == nil {
			continue
		}
		if explain.Fn != nil {
			if !explain.Fn(h) {
				continue
			}
			if explain.Scene != "" {
				h.SetScene(explain.Scene)
			}
			wp.Run(h, nil)
			h.Abort()
			return
		}
		if explain.Scene == "" {
			continue
		}
		q := h.C.Query(explain.Action)
		if q == "" {
			continue
		}
		h.SetScene(explain.Scene)
		for query, param := range explain.Param {
			h.C.AddParam(param, h.C.Query(query))
		}
		wp.Run(h, nil)
		h.Abort()
		return
	}
}

func ShowPreComment(h *wp.Handle) {
	v, ok := cache.NewCommentCache().Get(h.C, h.C.Request.URL.RawQuery)
	if ok {
		h.C.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
		h.C.Writer.WriteHeader(http.StatusOK)
		_, _ = h.C.Writer.Write([]byte(v))
		h.Abort()
	}
}

func CommonMiddleware(h *wp.Handle) {
	h.PushHandler(constraints.PipeMiddleware, constraints.Home,
		wp.NewHandleFn(MixWithPlain, 100, "middleware.MixWithPlain"),
	)
	h.PushHandler(constraints.PipeMiddleware, constraints.Detail,
		wp.NewHandleFn(ShowPreComment, 100, "middleware.ShowPreComment"),
	)
}
