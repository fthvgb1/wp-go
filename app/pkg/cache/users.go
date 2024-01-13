package cache

import (
	"context"
	"github.com/fthvgb1/wp-go/app/pkg/logs"
	"github.com/fthvgb1/wp-go/app/pkg/models"
	"github.com/fthvgb1/wp-go/cache/cachemanager"
	"time"
)

// GetUserByName query func see dao.GetUserByName
func GetUserByName(ctx context.Context, username string) (models.Users, error) {
	return cachemanager.GetBy[models.Users]("usernameToUserData", ctx, username, time.Second)
}

// GetAllUsername query func see dao.AllUsername
func GetAllUsername(ctx context.Context) (map[string]uint64, error) {
	return cachemanager.GetVarVal[map[string]uint64]("allUsername", ctx, time.Second)
}

// GetUserById query func see dao.GetUserById
func GetUserById(ctx context.Context, uid uint64) models.Users {
	r, err := cachemanager.GetBy[models.Users]("userData", ctx, uid, time.Second)
	logs.IfError(err, "get user", uid)
	return r
}
