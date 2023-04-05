package db

import (
	"context"
	"github.com/fthvgb1/wp-go/internal/pkg/config"
	"github.com/fthvgb1/wp-go/model"
	"github.com/fthvgb1/wp-go/safety"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"log"
)

var safeDb = safety.NewVar[*sqlx.DB](nil)

func InitDb() (*safety.Var[*sqlx.DB], error) {
	c := config.GetConfig()
	dsn := c.Mysql.Dsn.GetDsn()
	db, err := sqlx.Open("mysql", dsn)
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
	safeDb.Store(db)
	return safeDb, err
}

func QueryDb(db *safety.Var[*sqlx.DB]) *model.SqlxQuery {
	query := model.NewSqlxQuery(db, model.NewUniversalDb(
		nil,
		nil))
	model.SetSelect(query, func(ctx context.Context, a any, s string, args ...any) error {
		if config.GetConfig().ShowQuerySql {
			go log.Println(model.FormatSql(s, args...))
		}
		return query.Selects(ctx, a, s, args...)
	})
	model.SetGet(query, func(ctx context.Context, a any, s string, args ...any) error {
		if config.GetConfig().ShowQuerySql {
			go log.Println(model.FormatSql(s, args...))
		}
		return query.Gets(ctx, a, s, args...)
	})
	return query
}
