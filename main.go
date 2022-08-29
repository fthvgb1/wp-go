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
	T, t, err := models.WpPostsM.SimplePagination(nil, "wp_posts.ID,b.meta_id post_author", 4, 2, nil, models.SqlBuilder{{
		"left join", "wp_postmeta b", "b.post_id=wp_posts.ID",
	}})
	if err != nil {
		return
	}
	fmt.Println(T, t)
}
