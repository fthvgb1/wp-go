package main

import (
	"github/fthvgb1/wp-go/actions"
	"github/fthvgb1/wp-go/actions/common"
	"github/fthvgb1/wp-go/db"
	"github/fthvgb1/wp-go/models"
	"github/fthvgb1/wp-go/plugins"
	"github/fthvgb1/wp-go/route"
	"github/fthvgb1/wp-go/vars"
	"time"
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
	actions.InitFeed()
	common.InitActionsCommonCache()
	plugins.InitDigestCache()
	go cronClearCache()
}

func cronClearCache() {
	t := time.NewTicker(vars.Conf.CrontabClearCacheTime)
	for {
		select {
		case <-t.C:
			common.ClearCache()
			plugins.ClearDigestCache()
			actions.ClearCache()
		}
	}
}

func main() {
	err := route.SetupRouter().Run(":8082")
	if err != nil {
		panic(err)
	}
}
