package vars

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"time"
)

var Conf Config

type Config struct {
	Mysql                   Mysql           `yaml:"mysql"`
	Mail                    Mail            `yaml:"mail"`
	RecentPostCacheTime     time.Duration   `yaml:"recentPostCacheTime"`
	CategoryCacheTime       time.Duration   `yaml:"categoryCacheTime"`
	ArchiveCacheTime        time.Duration   `yaml:"archiveCacheTime"`
	ContextPostCacheTime    time.Duration   `yaml:"contextPostCacheTime"`
	RecentCommentsCacheTime time.Duration   `yaml:"recentCommentsCacheTime"`
	DigestCacheTime         time.Duration   `yaml:"digestCacheTime"`
	DigestWordCount         int             `yaml:"digestWordCount"`
	PostListCacheTime       time.Duration   `yaml:"postListCacheTime"`
	SearchPostCacheTime     time.Duration   `yaml:"searchPostCacheTime"`
	MonthPostCacheTime      time.Duration   `yaml:"monthPostCacheTime"`
	PostDataCacheTime       time.Duration   `yaml:"postDataCacheTime"`
	PostCommentsCacheTime   time.Duration   `yaml:"postCommentsCacheTime"`
	CrontabClearCacheTime   time.Duration   `yaml:"crontabClearCacheTime"`
	MaxRequestSleepNum      int64           `yaml:"maxRequestSleepNum"`
	SleepTime               []time.Duration `yaml:"sleepTime"`
	MaxRequestNum           int64           `yaml:"maxRequestNum"`
	SingleIpSearchNum       int64           `yaml:"singleIpSearchNum"`
	MaxPostIdCacheTime      time.Duration   `yaml:"maxPostIdCacheTime"`
	UserInfoCacheTime       time.Duration   `yaml:"userInfoCacheTime"`
	CommentsCacheTime       time.Duration   `yaml:"commentsCacheTime"`
	Gzip                    bool            `yaml:"gzip"`
	PostCommentUrl          string          `yaml:"postCommentUrl"`
	TrustIps                []string        `yaml:"trustIps"`
	TrustServerNames        []string        `yaml:"trustServerNames"`
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

func InitConfig() error {
	file, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(file, &Conf)
	if err != nil {
		return err
	}
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

func (m *Dsn) GetDsn() string {
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
