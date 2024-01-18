package db

import (
	"context"
	"github.com/fthvgb1/wp-go/app/pkg/config"
	"github.com/fthvgb1/wp-go/cache/reload"
	"github.com/fthvgb1/wp-go/model"
	"github.com/fthvgb1/wp-go/safety"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"log"
	"runtime"
)

var safeDb = safety.NewVar[*sqlx.DB](nil)
var showQuerySql func() bool

func InitDb() (*safety.Var[*sqlx.DB], error) {
	c := config.GetConfig()
	dsn := c.Mysql.Dsn.GetDsn()
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	preDb := safeDb.Load()
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
	if preDb != nil {
		_ = preDb.Close()
	}
	if showQuerySql == nil {
		showQuerySql = reload.BuildFnVal("showQuerySql", false, func() bool {
			return config.GetConfig().ShowQuerySql
		})
	}
	return safeDb, err
}

func QueryDb(db *safety.Var[*sqlx.DB]) *model.SqlxQuery {
	query := model.NewSqlxQuery(db, model.NewUniversalDb(
		nil,
		nil))
	model.SetSelect(query, func(ctx context.Context, a any, s string, args ...any) error {
		if showQuerySql() {
			_, f, l, _ := runtime.Caller(5)
			go func() {
				log.Printf("%v:%v %v\n", f, l, model.FormatSql(s, args...))
			}()
		}
		return query.Selects(ctx, a, s, args...)
	})
	model.SetGet(query, func(ctx context.Context, a any, s string, args ...any) error {
		if showQuerySql() {
			_, f, l, _ := runtime.Caller(5)
			go func() {
				log.Printf("%v:%v %v\n", f, l, model.FormatSql(s, args...))
			}()
		}
		return query.Gets(ctx, a, s, args...)
	})
	return query
}
