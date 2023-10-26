package dao

import (
	"context"
	"github.com/fthvgb1/wp-go/app/pkg/models"
	"github.com/fthvgb1/wp-go/helper/slice"
	"github.com/fthvgb1/wp-go/model"
)

func GetUserById(ctx context.Context, uid uint64, _ ...any) (r models.Users, err error) {
	r, err = model.FindOneById[models.Users](ctx, uid)
	return
}

func AllUsername(a ...any) (map[string]struct{}, error) {
	ctx := a[0].(context.Context)
	r, err := model.SimpleFind[models.Users](ctx, model.SqlBuilder{
		{"user_status", "=", "0", "int"},
	}, "user_login")
	if err != nil {
		return nil, err
	}
	return slice.ToMap(r, func(t models.Users) (string, struct{}) {
		return t.UserLogin, struct{}{}
	}, true), nil
}

func GetUserByName(ctx context.Context, u string, _ ...any) (r models.Users, err error) {
	r, err = model.FirstOne[models.Users](ctx, model.SqlBuilder{{
		"user_login", u,
	}}, "*", nil)
	return
}
