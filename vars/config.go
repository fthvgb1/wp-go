package vars

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"time"
)

var Conf Config

type Config struct {
	Mysql                Mysql         `yaml:"mysql"`
	RecentPostCacheTime  time.Duration `yaml:"recentPostCacheTime"`
	CategoryCacheTime    time.Duration `yaml:"categoryCacheTime"`
	ArchiveCacheTime     time.Duration `yaml:"archiveCacheTime"`
	ContextPostCacheTime time.Duration `yaml:"contextPostCacheTime"`
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