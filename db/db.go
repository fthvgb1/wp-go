package db

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github/fthvgb1/wp-go/config"
)

var Db *sqlx.DB

func InitDb() error {
	dsn := config.Conf.Mysql.Dsn.GetDsn()
	var err error
	Db, err = sqlx.Open("mysql", dsn)
	if err != nil {
		return err
	}
	if config.Conf.Mysql.Pool.ConnMaxIdleTime != 0 {
		Db.SetConnMaxIdleTime(config.Conf.Mysql.Pool.ConnMaxLifetime)
	}
	if config.Conf.Mysql.Pool.MaxIdleConn != 0 {
		Db.SetMaxIdleConns(config.Conf.Mysql.Pool.MaxIdleConn)
	}
	if config.Conf.Mysql.Pool.MaxOpenConn != 0 {
		Db.SetMaxOpenConns(config.Conf.Mysql.Pool.MaxOpenConn)
	}
	if config.Conf.Mysql.Pool.ConnMaxLifetime != 0 {
		Db.SetConnMaxLifetime(config.Conf.Mysql.Pool.ConnMaxLifetime)
	}
	return err
}
