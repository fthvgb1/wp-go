package common

import (
	"github.com/gin-gonic/gin"
	"github/fthvgb1/wp-go/logs"
	"github/fthvgb1/wp-go/models"
	"github/fthvgb1/wp-go/models/wp"
	"time"
)

func getUsers(...any) (m map[uint64]wp.WpUsers, err error) {
	m = make(map[uint64]wp.WpUsers)
	r, err := models.SimpleFind[wp.WpUsers](nil, "*")
	for _, user := range r {
		m[user.Id] = user
	}
	return
}

func GetUser(ctx *gin.Context, uid uint64) wp.WpUsers {
	r, err := usersCache.GetCache(ctx, uid, time.Second, uid)
	logs.ErrPrintln(err, "get user", uid)
	return r
}
