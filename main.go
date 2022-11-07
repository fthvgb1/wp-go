package main

import (
	"flag"
	"github/fthvgb1/wp-go/actions"
	"github/fthvgb1/wp-go/actions/common"
	"github/fthvgb1/wp-go/config"
	"github/fthvgb1/wp-go/db"
	"github/fthvgb1/wp-go/models"
	"github/fthvgb1/wp-go/models/wp"
	"github/fthvgb1/wp-go/plugins"
	"github/fthvgb1/wp-go/route"
	"math/rand"
	"time"
)

func init() {
	var c string
	flag.StringVar(&c, "c", "config.yaml", "config file")
	flag.Parse()
	rand.Seed(time.Now().UnixNano())
	err := config.InitConfig(c)
	if err != nil {
		panic(err)
	}

	err = db.InitDb()
	if err != nil {
		panic(err)
	}
	models.InitDB(db.NewSqlxDb(db.Db))
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
	if config.Conf.Port == "" {
		config.Conf.Port = "80"
	}
	err := route.SetupRouter().Run(config.Conf.Port)
	if err != nil {
		panic(err)
	}
}
