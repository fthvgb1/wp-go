package main

import (
	"fmt"
	"github/fthvgb1/wp-go/db"
	"github/fthvgb1/wp-go/models"
	"github/fthvgb1/wp-go/vars"
)

func init() {
	err := vars.InitConfig()
	if err != nil {
		panic(err)
	}
	err = db.InitDb()
	if err != nil {
		panic(err)
	}
}

func main() {
	T, err := models.WpPostsM.FindOneById(1)
	if err != nil {
		return
	}
	fmt.Println(T)
}
