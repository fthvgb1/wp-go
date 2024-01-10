package cache

import (
	"context"
	"github.com/fthvgb1/wp-go/app/pkg/logs"
	"github.com/fthvgb1/wp-go/app/pkg/models"
	"github.com/fthvgb1/wp-go/cache/cachemanager"
	"time"
)

func GetUserByName(ctx context.Context, username string) (models.Users, error) {
	return cachemanager.GetBy[models.Users]("usernameMapToUserData", ctx, username, time.Second)
}

func GetAllUsername(ctx context.Context) (map[string]struct{}, error) {
	return cachemanager.GetVarVal[map[string]struct{}]("allUsername", ctx, time.Second)
}

func GetUserById(ctx context.Context, uid uint64) models.Users {
	r, err := cachemanager.GetBy[models.Users]("userData", ctx, uid, time.Second)
	logs.IfError(err, "get user", uid)
	return r
}
