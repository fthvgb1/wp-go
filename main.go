package main

import (
	"github/fthvgb1/wp-go/actions"
	"github/fthvgb1/wp-go/actions/common"
	"github/fthvgb1/wp-go/config"
	"github/fthvgb1/wp-go/db"
	"github/fthvgb1/wp-go/models/wp"
	"github/fthvgb1/wp-go/plugins"
	"github/fthvgb1/wp-go/route"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
	err := config.InitConfig()
	if err != nil {
		panic(err)
	}

	err = db.InitDb()
	if err != nil {
		panic(err)
	}

	err = wp.InitOptions()
	if err != nil {
		panic(err)
	}

	err = wp.InitTerms()
	if err != nil {
		panic(err)
	}
	actions.InitFeed()
	common.InitActionsCommonCache()
	plugins.InitDigestCache()
	go cronClearCache()
}

func cronClearCache() {
	t := time.NewTicker(config.Conf.CrontabClearCacheTime)
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
