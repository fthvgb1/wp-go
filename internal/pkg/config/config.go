package config

import (
	"fmt"
	"github.com/fthvgb1/wp-go/safety"
	"gopkg.in/yaml.v2"
	"os"
	"time"
)

var config safety.Var[Config]

func GetConfig() Config {
	return config.Load()
}

type Config struct {
	Ssl                Ssl       `yaml:"ssl"`
	Mysql              Mysql     `yaml:"mysql"`
	Mail               Mail      `yaml:"mail"`
	CacheTime          CacheTime `yaml:"cacheTime"`
	DigestWordCount    int       `yaml:"digestWordCount"`
	MaxRequestSleepNum int64     `yaml:"maxRequestSleepNum"`
	MaxRequestNum      int64     `yaml:"maxRequestNum"`
	SingleIpSearchNum  int64     `yaml:"singleIpSearchNum"`
	Gzip               bool      `yaml:"gzip"`
	PostCommentUrl     string    `yaml:"postCommentUrl"`
	TrustIps           []string  `yaml:"trustIps"`
	TrustServerNames   []string  `yaml:"trustServerNames"`
	Theme              string    `yaml:"theme"`
	PostOrder          string    `yaml:"postOrder"`
	UploadDir          string    `yaml:"uploadDir"`
	Pprof              string    `yaml:"pprof"`
	ListPagePlugins    []string  `yaml:"listPagePlugins"`
	PaginationStep     int       `yaml:"paginationStep"`
	ShowQuerySql       bool      `yaml:"showQuerySql"`
	Plugins            []string  `yaml:"plugins"`
}

type CacheTime struct {
	CacheControl             time.Duration   `yaml:"cacheControl"`
	RecentPostCacheTime      time.Duration   `yaml:"recentPostCacheTime"`
	CategoryCacheTime        time.Duration   `yaml:"categoryCacheTime"`
	ArchiveCacheTime         time.Duration   `yaml:"archiveCacheTime"`
	ContextPostCacheTime     time.Duration   `yaml:"contextPostCacheTime"`
	RecentCommentsCacheTime  time.Duration   `yaml:"recentCommentsCacheTime"`
	DigestCacheTime          time.Duration   `yaml:"digestCacheTime"`
	PostListCacheTime        time.Duration   `yaml:"postListCacheTime"`
	SearchPostCacheTime      time.Duration   `yaml:"searchPostCacheTime"`
	MonthPostCacheTime       time.Duration   `yaml:"monthPostCacheTime"`
	PostDataCacheTime        time.Duration   `yaml:"postDataCacheTime"`
	PostCommentsCacheTime    time.Duration   `yaml:"postCommentsCacheTime"`
	CrontabClearCacheTime    time.Duration   `yaml:"crontabClearCacheTime"`
	MaxPostIdCacheTime       time.Duration   `yaml:"maxPostIdCacheTime"`
	UserInfoCacheTime        time.Duration   `yaml:"userInfoCacheTime"`
	CommentsCacheTime        time.Duration   `yaml:"commentsCacheTime"`
	ThemeHeaderImagCacheTime time.Duration   `yaml:"themeHeaderImagCacheTime"`
	SleepTime                []time.Duration `yaml:"sleepTime"`
}

type Ssl struct {
	Cert string `yaml:"cert"`
	Key  string `yaml:"key"`
}

type Mail struct {
	User  string `yaml:"user"`
	Alias string `yaml:"alias"`
	Pass  string `yaml:"pass"`
	Host  string `yaml:"host"`
	Port  int    `yaml:"port"`
	Ssl   bool   `yaml:"ssl"`
}

type Mysql struct {
	Dsn  Dsn  `yaml:"dsn"`
	Pool Pool `yaml:"pool"`
}

func InitConfig(conf string) error {
	if conf == "" {
		conf = "config.yaml"
	}
	file, err := os.ReadFile(conf)
	if err != nil {
		return err
	}
	var c Config
	err = yaml.Unmarshal(file, &c)
	if err != nil {
		return err
	}
	config.Store(c)
	return nil
}

type Dsn struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Db       string `yaml:"db"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Charset  string `yaml:"charset"`
}

func (m Dsn) GetDsn() string {
	if m.Charset == "" {
		m.Charset = "utf8"
	}
	t := "%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local"
	return fmt.Sprintf(t, m.User, m.Password, m.Host, m.Port, m.Db, m.Charset)
}

type Pool struct {
	ConnMaxIdleTime time.Duration `yaml:"connMaxIdleTime"`
	MaxOpenConn     int           `yaml:"maxOpenConn"`
	MaxIdleConn     int           `yaml:"maxIdleConn"`
	ConnMaxLifetime time.Duration `yaml:"connMaxLifetime"`
}
