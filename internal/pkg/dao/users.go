package common

import (
	"context"
	"github.com/fthvgb1/wp-go/internal/pkg/models"
	"github.com/fthvgb1/wp-go/model"
)

func GetUserById(a ...any) (r models.Users, err error) {
	ctx := a[0].(context.Context)
	uid := a[1].(uint64)
	r, err = model.FindOneById[models.Users](ctx, uid)
	return
}

func GetUserByName(a ...any) (r models.Users, err error) {
	u := a[1].(string)
	ctx := a[0].(context.Context)
	r, err = model.FirstOne[models.Users](ctx, model.SqlBuilder{{
		"user_login", u,
	}}, "*", nil)
	return
}
