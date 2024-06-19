package main

import (
	"flag"
	"fmt"
	"github.com/fthvgb1/wp-go/app/ossigns"
	"github.com/fthvgb1/wp-go/app/pkg/cache"
	"github.com/fthvgb1/wp-go/app/pkg/config"
	"github.com/fthvgb1/wp-go/app/pkg/db"
	"github.com/fthvgb1/wp-go/app/pkg/logs"
	"github.com/fthvgb1/wp-go/app/plugins"
	"github.com/fthvgb1/wp-go/app/plugins/wphandle"
	"github.com/fthvgb1/wp-go/app/route"
	"github.com/fthvgb1/wp-go/app/theme"
	"github.com/fthvgb1/wp-go/app/wpconfig"
	"github.com/fthvgb1/wp-go/cache/cachemanager"
	"github.com/fthvgb1/wp-go/model"
	"os"
	"regexp"
	"strings"
	"time"
)

var confPath string
var address string
var intReg = regexp.MustCompile(`^\d`)

func inits() {
	flag.StringVar(&confPath, "c", "config.yaml", "config file support json,yaml or url")
	flag.StringVar(&address, "p", "", "listen address and port")
	flag.Parse()
	if address == "" && os.Getenv("PORT") == "" {
		address = "80"
	}
	if intReg.MatchString(address) && !strings.Contains(address, ":") {
		address = ":" + address
	}
	err := initConf(confPath)
	if err != nil {
		panic(err)
	}
	cache.InitActionsCommonCache()
	plugins.InitDigestCache()
	theme.InitTheme()
	go cronClearCache()
}

func initConf(c string) (err error) {
	err = config.InitConfig(c)
	if err != nil {
		return
	}
	err = config.InitTrans()
	if err != nil {
		return err
	}
	err = logs.InitLogger()
	if err != nil {
		return err
	}
	database, err := db.InitDb()
	if err != nil {
		return
	}
	model.InitDB(db.QueryDb(database))
	err = wpconfig.InitOptions()
	if err != nil {
		return
	}
	err = wpconfig.InitTerms()
	if err != nil {
		return
	}
	wphandle.LoadPlugins()
	return
}

func cronClearCache() {
	t := time.NewTicker(config.GetConfig().CacheTime.CrontabClearCacheTime)
	for {
		select {
		case <-t.C:
			cachemanager.ClearExpired()
		}
	}
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
			os.Exit(1)
		}
	}()
	inits()
	ossigns.SetConfPath(confPath)
	go ossigns.SignalNotify()
	Gin := route.SetupRouter()
	c := config.GetConfig()
	if c.Ssl.Key != "" && c.Ssl.Cert != "" {
		err := Gin.RunTLS(address, c.Ssl.Cert, c.Ssl.Key)
		if err != nil {
			panic(err)
		}
		return
	}
	err := Gin.Run(address)
	if err != nil {
		panic(err)
	}
}
