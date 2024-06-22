package ossigns

import (
	"fmt"
	"github.com/fthvgb1/wp-go/app/mail"
	"github.com/fthvgb1/wp-go/app/pkg/config"
	"github.com/fthvgb1/wp-go/app/pkg/db"
	"github.com/fthvgb1/wp-go/app/pkg/logs"
	"github.com/fthvgb1/wp-go/app/plugins/wphandle"
	"github.com/fthvgb1/wp-go/app/wpconfig"
	"github.com/fthvgb1/wp-go/cache/cachemanager"
	"github.com/fthvgb1/wp-go/cache/reload"
	"github.com/fthvgb1/wp-go/signs"
	"log"
	"syscall"
)

var confPath string

func SetConfPath(path string) {
	confPath = path
}

func FlushCache() {
	defer func() {
		if r := recover(); r != nil {
			err := mail.SendMail([]string{config.GetConfig().Mail.User}, "清空缓存失败", fmt.Sprintf("err:[%s]", r))
			logs.IfError(err, "发邮件失败")
		}
	}()
	cachemanager.Flush()
	log.Println("all cache flushed")
}

func Reloads() {
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
	reload.Reloads("themeArgAndConfig")
	FlushCache()
	log.Println("reload complete")
}

func SignalNotify() {
	rel := func() bool {
		go Reloads()
		return true
	}
	flu := func() bool {
		go FlushCache()
		return true
	}
	signs.Install(syscall.SIGUSR1, rel, "reload")
	signs.Install(syscall.SIGUSR2, flu, "flush")
	signs.Wait()
}
