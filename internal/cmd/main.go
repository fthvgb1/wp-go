package main

import (
	"flag"
	"fmt"
	"github.com/fthvgb1/wp-go/internal/cmd/route"
	"github.com/fthvgb1/wp-go/internal/mail"
	"github.com/fthvgb1/wp-go/internal/pkg/cache"
	"github.com/fthvgb1/wp-go/internal/pkg/config"
	"github.com/fthvgb1/wp-go/internal/pkg/db"
	"github.com/fthvgb1/wp-go/internal/pkg/logs"
	"github.com/fthvgb1/wp-go/internal/plugins"
	"github.com/fthvgb1/wp-go/internal/theme"
	"github.com/fthvgb1/wp-go/internal/wpconfig"
	"github.com/fthvgb1/wp-go/model"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
	"time"
)

var confPath string
var address string
var middleWareReloadFn func()
var intReg = regexp.MustCompile(`^\d`)

func init() {
	flag.StringVar(&confPath, "c", "config.yaml", "config file")
	flag.StringVar(&address, "p", "", "listen address(port)")
	flag.Parse()
	if address == "" && os.Getenv("PORT") == "" {
		address = "80"
	}
	if intReg.MatchString(address) && !strings.Contains(address, ":") {
		address = ":" + address
	}
	rand.Seed(time.Now().UnixNano())
	err := initConf(confPath)
	if err != nil {
		panic(err)
	}
	cache.InitActionsCommonCache()
	plugins.InitDigestCache()
	theme.InitThemeAndTemplateFuncMap()
	go cronClearCache()
}

func initConf(c string) (err error) {
	err = config.InitConfig(c)
	if err != nil {
		return
	}

	database, err := db.InitDb()
	if err != nil {
		return
	}
	model.InitDB(model.NewSqlxQuery(database))
	err = wpconfig.InitOptions()
	if err != nil {
		return
	}
	err = wpconfig.InitTerms()
	if err != nil {
		return
	}
	return
}

func cronClearCache() {
	t := time.NewTicker(config.GetConfig().CacheTime.CrontabClearCacheTime)
	for {
		select {
		case <-t.C:
			cache.ClearCache()
			plugins.ClearDigestCache()
		}
	}
}

func flushCache() {
	defer func() {
		if r := recover(); r != nil {
			err := mail.SendMail([]string{config.GetConfig().Mail.User}, "清空缓存失败", fmt.Sprintf("err:[%s]", r))
			logs.ErrPrintln(err, "发邮件失败")
		}
	}()
	cache.FlushCache()
	plugins.FlushCache()
	log.Println("all cache flushed")
}

func reload() {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()
	err := config.InitConfig(confPath)
	logs.ErrPrintln(err, "获取配置文件失败", confPath)
	err = wpconfig.InitOptions()
	logs.ErrPrintln(err, "获取网站设置WpOption失败")
	err = wpconfig.InitTerms()
	logs.ErrPrintln(err, "获取WpTerms表失败")
	if middleWareReloadFn != nil {
		middleWareReloadFn()
	}
	flushCache()
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

func main() {
	go signalNotify()
	Gin, reloadFn := route.SetupRouter()
	c := config.GetConfig()
	middleWareReloadFn = reloadFn
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
