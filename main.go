package main

import (
	"flag"
	"fmt"
	"github/fthvgb1/wp-go/actions"
	"github/fthvgb1/wp-go/actions/common"
	"github/fthvgb1/wp-go/config"
	"github/fthvgb1/wp-go/db"
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
	err = config.InitOptions()
	if err != nil {
		panic(err)
	}
	err = config.InitTerms()
	if err != nil {
		panic(err)
	}
	actions.InitFeed()
	common.InitActionsCommonCache()
	plugins.InitDigestCache()
	go cronClearCache()
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

func signalNotify() {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGUSR1, syscall.SIGUSR2)
	conf := config.Conf.Load()
	for {
		switch <-c {
		case syscall.SIGUSR1:
		//todo 更新配置
		case syscall.SIGUSR2:
			go func() {
				defer func() {
					if r := recover(); r != nil {
						mail.SendMail([]string{conf.Mail.User}, "清空缓存失败", fmt.Sprintf("err:[%s]", r))
					}
				}()
				common.FlushCache()
				plugins.FlushCache()
				actions.FlushCache()
				log.Println("清除缓存成功")
			}()
		}
	}
}

func main() {
	c := config.Conf.Load()
	if c.Port == "" {
		c.Port = "80"
		config.Conf.Store(c)
	}
	go signalNotify()
	err := route.SetupRouter().Run(c.Port)
	if err != nil {
		panic(err)
	}
}
