package common

import (
	"context"
	"github/fthvgb1/wp-go/internal/wp"
	"github/fthvgb1/wp-go/logs"
	"github/fthvgb1/wp-go/models"
	"time"
)

func getUserById(a ...any) (r wp.Users, err error) {
	ctx := a[0].(context.Context)
	uid := a[1].(uint64)
	r, err = models.FindOneById[wp.Users](ctx, uid)
	return
}

func GetUserByName(ctx context.Context, username string) (wp.Users, error) {
	return usersNameCache.GetCache(ctx, username, time.Second, ctx, username)
}

func getUserByName(a ...any) (r wp.Users, err error) {
	u := a[1].(string)
	ctx := a[0].(context.Context)
	r, err = models.FirstOne[wp.Users](ctx, models.SqlBuilder{{
		"user_login", u,
	}}, "*", nil)
	return
}

func GetUserById(ctx context.Context, uid uint64) wp.Users {
	r, err := usersCache.GetCache(ctx, uid, time.Second, ctx, uid)
	logs.ErrPrintln(err, "get user", uid)
	return r
}
