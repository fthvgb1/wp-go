package main

import (
	"flag"
	"fmt"
	"github.com/fthvgb1/wp-go/app/cmd/cachemanager"
	"github.com/fthvgb1/wp-go/app/cmd/reload"
	"github.com/fthvgb1/wp-go/app/cmd/route"
	"github.com/fthvgb1/wp-go/app/mail"
	"github.com/fthvgb1/wp-go/app/pkg/cache"
	"github.com/fthvgb1/wp-go/app/pkg/config"
	"github.com/fthvgb1/wp-go/app/pkg/db"
	"github.com/fthvgb1/wp-go/app/pkg/logs"
	"github.com/fthvgb1/wp-go/app/plugins"
	"github.com/fthvgb1/wp-go/app/plugins/wphandle"
	"github.com/fthvgb1/wp-go/app/theme"
	"github.com/fthvgb1/wp-go/app/wpconfig"
	"github.com/fthvgb1/wp-go/model"
	"log"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
	"time"
)

var confPath string
var address string
var intReg = regexp.MustCompile(`^\d`)

func init() {
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

func flushCache() {
	defer func() {
		if r := recover(); r != nil {
			err := mail.SendMail([]string{config.GetConfig().Mail.User}, "清空缓存失败", fmt.Sprintf("err:[%s]", r))
			logs.IfError(err, "发邮件失败")
		}
	}()
	cachemanager.Flush()
	log.Println("all cache flushed")
}

func reloads() {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()
	err := config.InitConfig(confPath)
	logs.IfError(err, "获取配置文件失败", confPath)
	err = logs.InitLogger()
	logs.IfError(err, "日志配置错误")
	_, err = db.InitDb()
	logs.IfError(err, "重新读取db失败", config.GetConfig().Mysql)
	err = wpconfig.InitOptions()
	logs.IfError(err, "获取网站设置WpOption失败")
	err = wpconfig.InitTerms()
	logs.IfError(err, "获取WpTerms表失败")
	wphandle.LoadPlugins()
	reload.Reload()
	flushCache()
	log.Println("reload complete")
}

func signalNotify() {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGUSR1, syscall.SIGUSR2)
	for {
		switch <-c {
		case syscall.SIGUSR1:
			go reloads()
		case syscall.SIGUSR2:
			go flushCache()
		}
	}
}

func main() {
	go signalNotify()
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
