package config

import (
	"encoding/json"
	"fmt"
	"github.com/fthvgb1/wp-go/safety"
	"gopkg.in/yaml.v2"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

var config safety.Var[Config]

func GetConfig() Config {
	return config.Load()
}

type Config struct {
	Ssl                Ssl       `yaml:"ssl" json:"ssl"`
	Mysql              Mysql     `yaml:"mysql" json:"mysql"`
	Mail               Mail      `yaml:"mail" json:"mail"`
	CacheTime          CacheTime `yaml:"cacheTime" json:"cacheTime"`
	PluginPath         string    `yaml:"pluginPath" json:"pluginPath"`
	ExternScript       []string  `json:"externScript" yaml:"externScript"`
	DigestWordCount    int       `yaml:"digestWordCount" json:"digestWordCount,omitempty"`
	DigestAllowTag     string    `yaml:"digestAllowTag" json:"digestAllowTag"`
	MaxRequestSleepNum int64     `yaml:"maxRequestSleepNum" json:"maxRequestSleepNum,omitempty"`
	MaxRequestNum      int64     `yaml:"maxRequestNum" json:"maxRequestNum,omitempty"`
	SingleIpSearchNum  int64     `yaml:"singleIpSearchNum" json:"singleIpSearchNum,omitempty"`
	Gzip               bool      `yaml:"gzip" json:"gzip,omitempty"`
	PostCommentUrl     string    `yaml:"postCommentUrl" json:"postCommentUrl,omitempty"`
	TrustIps           []string  `yaml:"trustIps" json:"trustIps,omitempty"`
	TrustServerNames   []string  `yaml:"trustServerNames" json:"trustServerNames,omitempty"`
	Theme              string    `yaml:"theme" json:"theme,omitempty"`
	PostOrder          string    `yaml:"postOrder" json:"postOrder,omitempty"`
	UploadDir          string    `yaml:"uploadDir" json:"uploadDir,omitempty"`
	Pprof              string    `yaml:"pprof" json:"pprof,omitempty"`
	ListPagePlugins    []string  `yaml:"listPagePlugins" json:"listPagePlugins,omitempty"`
	PaginationStep     int       `yaml:"paginationStep" json:"paginationStep,omitempty"`
	ShowQuerySql       bool      `yaml:"showQuerySql" json:"showQuerySql,omitempty"`
	Plugins            []string  `yaml:"plugins" json:"plugins,omitempty"`
	LogOutput          string    `yaml:"logOutput" json:"logOutput,omitempty"`
	WpDir              string    `yaml:"wpDir" json:"wpDir"`
}

type CacheTime struct {
	CacheControl            time.Duration   `yaml:"cacheControl" json:"cacheControl,omitempty"`
	RecentPostCacheTime     time.Duration   `yaml:"recentPostCacheTime" json:"recentPostCacheTime,omitempty"`
	CategoryCacheTime       time.Duration   `yaml:"categoryCacheTime" json:"categoryCacheTime,omitempty"`
	ArchiveCacheTime        time.Duration   `yaml:"archiveCacheTime" json:"archiveCacheTime,omitempty"`
	ContextPostCacheTime    time.Duration   `yaml:"contextPostCacheTime" json:"contextPostCacheTime,omitempty"`
	RecentCommentsCacheTime time.Duration   `yaml:"recentCommentsCacheTime" json:"recentCommentsCacheTime,omitempty"`
	DigestCacheTime         time.Duration   `yaml:"digestCacheTime" json:"digestCacheTime,omitempty"`
	PostListCacheTime       time.Duration   `yaml:"postListCacheTime" json:"postListCacheTime,omitempty"`
	SearchPostCacheTime     time.Duration   `yaml:"searchPostCacheTime" json:"searchPostCacheTime,omitempty"`
	MonthPostCacheTime      time.Duration   `yaml:"monthPostCacheTime" json:"monthPostCacheTime,omitempty"`
	PostDataCacheTime       time.Duration   `yaml:"postDataCacheTime" json:"postDataCacheTime,omitempty"`
	PostCommentsCacheTime   time.Duration   `yaml:"postCommentsCacheTime" json:"postCommentsCacheTime,omitempty"`
	CrontabClearCacheTime   time.Duration   `yaml:"crontabClearCacheTime" json:"crontabClearCacheTime,omitempty"`
	MaxPostIdCacheTime      time.Duration   `yaml:"maxPostIdCacheTime" json:"maxPostIdCacheTime,omitempty"`
	UserInfoCacheTime       time.Duration   `yaml:"userInfoCacheTime" json:"userInfoCacheTime,omitempty"`
	CommentsCacheTime       time.Duration   `yaml:"commentsCacheTime" json:"commentsCacheTime,omitempty"`
	SleepTime               []time.Duration `yaml:"sleepTime" json:"sleepTime,omitempty"`
}

type Ssl struct {
	Cert string `yaml:"cert" json:"cert,omitempty"`
	Key  string `yaml:"key" json:"key,omitempty"`
}

type Mail struct {
	User               string `yaml:"user" json:"user,omitempty"`
	Alias              string `yaml:"alias" json:"alias,omitempty"`
	Pass               string `yaml:"pass" json:"pass,omitempty"`
	Host               string `yaml:"host" json:"host,omitempty"`
	Port               int    `yaml:"port" json:"port,omitempty"`
	InsecureSkipVerify bool   `yaml:"insecureSkipVerify" json:"insecureSkipVerify,omitempty"`
}

type Mysql struct {
	Dsn  Dsn  `yaml:"dsn" json:"dsn"`
	Pool Pool `yaml:"pool" json:"pool"`
}

func InitConfig(conf string) error {
	if conf == "" {
		conf = "config.yaml"
	}
	var file []byte
	var err error
	if strings.Contains(conf, "http") {
		get, err := http.Get(conf)
		if err != nil {
			return err
		}
		file, err = io.ReadAll(get.Body)
	} else {
		file, err = os.ReadFile(conf)
	}
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

func jsonToYaml[T any](b []byte, c T) error {
	var v map[string]any
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}
	marshal, er := yaml.Marshal(v)
	if er != nil {
		return er
	}
	err = yaml.Unmarshal(marshal, c)
	return err
}

type Dsn struct {
	Host     string `yaml:"host" json:"host,omitempty"`
	Port     string `yaml:"port" json:"port,omitempty"`
	Db       string `yaml:"db" json:"db,omitempty"`
	User     string `yaml:"user" json:"user,omitempty"`
	Password string `yaml:"password" json:"password,omitempty"`
	Charset  string `yaml:"charset" json:"charset,omitempty"`
}

func (m Dsn) GetDsn() string {
	if m.Charset == "" {
		m.Charset = "utf8"
	}
	t := "%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local"
	return fmt.Sprintf(t, m.User, m.Password, m.Host, m.Port, m.Db, m.Charset)
}

type Pool struct {
	ConnMaxIdleTime time.Duration `yaml:"connMaxIdleTime" json:"connMaxIdleTime,omitempty"`
	MaxOpenConn     int           `yaml:"maxOpenConn" json:"maxOpenConn,omitempty"`
	MaxIdleConn     int           `yaml:"maxIdleConn" json:"maxIdleConn,omitempty"`
	ConnMaxLifetime time.Duration `yaml:"connMaxLifetime" json:"connMaxLifetime,omitempty"`
}
