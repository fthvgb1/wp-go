package main

import (
	"github/fthvgb1/wp-go/actions/common"
	"github/fthvgb1/wp-go/db"
	"github/fthvgb1/wp-go/models"
	"github/fthvgb1/wp-go/plugins"
	"github/fthvgb1/wp-go/route"
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

	err = models.InitOptions()
	if err != nil {
		panic(err)
	}

	err = models.InitTerms()
	if err != nil {
		panic(err)
	}

	common.InitCache()
	plugins.InitDigest()
}

func main() {
	err := route.SetupRouter().Run(":8082")
	if err != nil {
		panic(err)
	}
}
