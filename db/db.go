package db

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github/fthvgb1/wp-go/vars"
)

var Db *sqlx.DB

func InitDb() error {
	dsn := vars.Conf.Mysql.Dsn.GetDsn()
	var err error
	Db, err = sqlx.Open("mysql", dsn)
	if err != nil {
		return err
	}
	if vars.Conf.Mysql.Pool.ConnMaxIdleTime != 0 {
		Db.SetConnMaxIdleTime(vars.Conf.Mysql.Pool.ConnMaxLifetime)
	}
	if vars.Conf.Mysql.Pool.MaxIdleConn != 0 {
		Db.SetMaxIdleConns(vars.Conf.Mysql.Pool.MaxIdleConn)
	}
	if vars.Conf.Mysql.Pool.MaxOpenConn != 0 {
		Db.SetMaxOpenConns(vars.Conf.Mysql.Pool.MaxOpenConn)
	}
	if vars.Conf.Mysql.Pool.ConnMaxLifetime != 0 {
		Db.SetConnMaxLifetime(vars.Conf.Mysql.Pool.ConnMaxLifetime)
	}
	return err
}
