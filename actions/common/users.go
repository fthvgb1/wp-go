package common

import (
	"context"
	"github.com/gin-gonic/gin"
	"github/fthvgb1/wp-go/logs"
	"github/fthvgb1/wp-go/models"
	"github/fthvgb1/wp-go/models/wp"
	"time"
)

func getUsers(a ...any) (m map[uint64]wp.Users, err error) {
	m = make(map[uint64]wp.Users)
	ctx := a[0].(context.Context)
	r, err := models.SimpleFind[wp.Users](ctx, nil, "*")
	for _, user := range r {
		m[user.Id] = user
	}
	return
}

func GetUser(ctx *gin.Context, uid uint64) wp.Users {
	r, err := usersCache.GetCache(ctx, uid, time.Second, ctx, uid)
	logs.ErrPrintln(err, "get user", uid)
	return r
}
