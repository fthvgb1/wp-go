package db

import (
	"github.com/fthvgb1/wp-go/internal/pkg/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

func InitDb() (*sqlx.DB, error) {
	c := config.GetConfig()
	dsn := c.Mysql.Dsn.GetDsn()
	var err error
	db, err = sqlx.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if c.Mysql.Pool.ConnMaxIdleTime != 0 {
		db.SetConnMaxIdleTime(c.Mysql.Pool.ConnMaxLifetime)
	}
	if c.Mysql.Pool.MaxIdleConn != 0 {
		db.SetMaxIdleConns(c.Mysql.Pool.MaxIdleConn)
	}
	if c.Mysql.Pool.MaxOpenConn != 0 {
		db.SetMaxOpenConns(c.Mysql.Pool.MaxOpenConn)
	}
	if c.Mysql.Pool.ConnMaxLifetime != 0 {
		db.SetConnMaxLifetime(c.Mysql.Pool.ConnMaxLifetime)
	}
	return db, err
}
