package model

import (
	"context"
	"github.com/fthvgb1/wp-go/safety"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"reflect"
	"sync"
	"testing"
)

var ctx = context.Background()

var glob = safety.NewMap[string, dbQuery[Model]]()
var dbMap = sync.Map{}

var sq *sqlx.DB

func anyDb[T Model]() *SqlxQuery[T] {
	var a T
	db, ok := dbMap.Load(a.Table())
	if ok {
		return db.(*SqlxQuery[T])
	}
	dbb := NewSqlxQuery[T](sq, UniversalDb[T]{nil, nil})
	dbMap.Store(a.Table(), dbb)
	return dbb
}

func init() {
	db, err := sqlx.Open("mysql", "root:root@tcp(192.168.66.47:3306)/wordpress?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	sq = db
	//glob = NewSqlxQuery(db, NewUniversalDb(nil, nil))

}

func Test_selects(t *testing.T) {
	type args[T Model] struct {
		db  dbQuery[T]
		ctx context.Context
		q   *QueryCondition
	}
	type testCase[T Model] struct {
		name    string
		args    args[T]
		want    []T
		wantErr bool
	}
	tests := []testCase[options]{
		{
			name: "t1",
			args: args[options]{
				anyDb[options](),
				ctx,
				Conditions(Where(SqlBuilder{{"option_name", "blogname"}})),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := finds[options](tt.args.db, tt.args.ctx, tt.args.q)
			if (err != nil) != tt.wantErr {
				t.Errorf("finds() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("finds() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkSelectXX(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := finds[options](anyDb[options](), ctx, Conditions())
		if err != nil {
			panic(err)
		}
	}
}
func BenchmarkScannerXX(b *testing.B) {
	for i := 0; i < b.N; i++ {

		_, err := scanners[options](anyDb[options](), ctx, Conditions())
		if err != nil {
			panic(err)
		}
	}
}

type options struct {
	OptionId    uint64 `gorm:"column:option_id" db:"option_id" json:"option_id" form:"option_id"`
	OptionName  string `gorm:"column:option_name" db:"option_name" json:"option_name" form:"option_name"`
	OptionValue string `gorm:"column:option_value" db:"option_value" json:"option_value" form:"option_value"`
	Autoload    string `gorm:"column:autoload" db:"autoload" json:"autoload" form:"autoload"`
}

func (w options) PrimaryKey() string {
	return "option_id"
}

func (w options) Table() string {
	return "wp_options"
}

func Test_scanners(t *testing.T) {
	type args[T Model] struct {
		db  dbQuery[T]
		ctx context.Context
		q   *QueryCondition
	}
	type testCase[T Model] struct {
		name    string
		args    args[T]
		wantErr bool
	}
	tests := []testCase[options]{
		{
			name: "t1",
			args: args[options]{
				anyDb[options](),
				ctx,
				Conditions(Where(SqlBuilder{{"option_name", "blogname"}})),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := scanners[options](tt.args.db, tt.args.ctx, tt.args.q); (err != nil) != tt.wantErr {
				t.Errorf("scanners() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
