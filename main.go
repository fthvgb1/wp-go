package main

import (
	"fmt"
	"github/fthvgb1/wp-go/db"
	"github/fthvgb1/wp-go/models"
	"github/fthvgb1/wp-go/vars"
)

func init() {
	err := vars.InitDbConfig()
	if err != nil {
		panic(err)
	}
	err = db.InitDb()
	if err != nil {
		panic(err)
	}
}

func main() {
	T, total, err := models.WpPostsM.SimplePagination(nil, "*", 1, 10, models.SqlBuilder{{"ID", "desc"}})
	if err != nil {
		return
	}
	fmt.Println(T, total)
}
