package cache

import (
	"context"
	"github/fthvgb1/wp-go/internal/logs"
	"github/fthvgb1/wp-go/internal/models"
	"github/fthvgb1/wp-go/model"
	"time"
)

func getUserById(a ...any) (r models.Users, err error) {
	ctx := a[0].(context.Context)
	uid := a[1].(uint64)
	r, err = model.FindOneById[models.Users](ctx, uid)
	return
}

func GetUserByName(ctx context.Context, username string) (models.Users, error) {
	return usersNameCache.GetCache(ctx, username, time.Second, ctx, username)
}

func GetUserById(ctx context.Context, uid uint64) models.Users {
	r, err := usersCache.GetCache(ctx, uid, time.Second, ctx, uid)
	logs.ErrPrintln(err, "get user", uid)
	return r
}
