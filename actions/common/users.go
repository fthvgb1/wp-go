package common

import (
	"github.com/gin-gonic/gin"
	"github/fthvgb1/wp-go/logs"
	"github/fthvgb1/wp-go/models"
	"time"
)

func getUsers(...any) (m map[uint64]models.WpUsers, err error) {
	m = make(map[uint64]models.WpUsers)
	r, err := models.Find[models.WpUsers](nil, "*", "", nil, nil, nil, 0)
	for _, user := range r {
		m[user.Id] = user
	}
	return
}

func GetUser(ctx *gin.Context, uid uint64) models.WpUsers {
	r, err := usersCache.GetCache(ctx, uid, time.Second, uid)
	logs.ErrPrintln(err, "get user", uid)
	return r
}
