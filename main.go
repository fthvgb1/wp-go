package main

import (
	"flag"
	"fmt"
	"github/fthvgb1/wp-go/actions"
	"github/fthvgb1/wp-go/actions/common"
	"github/fthvgb1/wp-go/config"
	"github/fthvgb1/wp-go/db"
	"github/fthvgb1/wp-go/logs"
	"github/fthvgb1/wp-go/mail"
	"github/fthvgb1/wp-go/models"
	"github/fthvgb1/wp-go/plugins"
	"github/fthvgb1/wp-go/route"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var confPath string

func init() {
	flag.StringVar(&confPath, "c", "config.yaml", "config file")
	flag.Parse()
	rand.Seed(time.Now().UnixNano())
	err := initConf(confPath)
	if err != nil {
		panic(err)
	}
	actions.InitFeed()
	common.InitActionsCommonCache()
	plugins.InitDigestCache()
	go cronClearCache()
}

func initConf(c string) (err error) {
	err = config.InitConfig(c)
	if err != nil {
		return
	}

	err = db.InitDb()
	if err != nil {
		return
	}
	models.InitDB(db.NewSqlxDb(db.Db))
	err = config.InitOptions()
	if err != nil {
		return
	}
	err = config.InitTerms()
	if err != nil {
		return
	}
	return
}

func cronClearCache() {
	t := time.NewTicker(config.Conf.Load().CrontabClearCacheTime)
	for {
		select {
		case <-t.C:
			common.ClearCache()
			plugins.ClearDigestCache()
			actions.ClearCache()
		}
	}
}

func flushCache() {
	defer func() {
		if r := recover(); r != nil {
			err := mail.SendMail([]string{config.Conf.Load().Mail.User}, "清空缓存失败", fmt.Sprintf("err:[%s]", r))
			logs.ErrPrintln(err, "发邮件失败")
		}
	}()
	common.FlushCache()
	plugins.FlushCache()
	actions.FlushCache()
	log.Println("清除缓存成功")
}

func reload() {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()
	err := config.InitConfig(confPath)
	logs.ErrPrintln(err, "获取配置文件失败", confPath)
	err = config.InitOptions()
	logs.ErrPrintln(err, "获取网站设置WpOption失败")
	err = config.InitTerms()
	logs.ErrPrintln(err, "获取WpTerms表失败")
	middleWareReloadFn()
	log.Println("reload complete")
}

func signalNotify() {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGUSR1, syscall.SIGUSR2)
	for {
		switch <-c {
		case syscall.SIGUSR1:
			go reload()
		case syscall.SIGUSR2:
			go flushCache()
		}
	}
}

var middleWareReloadFn func()

func main() {
	c := config.Conf.Load()
	if c.Port == "" {
		c.Port = "80"
		config.Conf.Store(c)
	}
	go signalNotify()
	Gin, reloadFn := route.SetupRouter()
	middleWareReloadFn = reloadFn
	err := Gin.Run(c.Port)
	if err != nil {
		panic(err)
	}
}
