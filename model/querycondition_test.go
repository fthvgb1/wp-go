package model

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/fthvgb1/wp-go/helper/number"
	"github.com/fthvgb1/wp-go/helper/slice"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func TestFinds(t *testing.T) {
	type args struct {
		ctx context.Context
		q   QueryCondition
	}
	type testCase[T Model] struct {
		name    string
		args    args
		wantR   []T
		wantErr bool
	}
	tests := []testCase[post]{
		{
			name: "t1",
			args: args{
				ctx: context.Background(),
				q: Conditions(
					Where(SqlBuilder{
						{"post_status", "publish"}, {"ID", "in", ""}},
					),
					Order(SqlBuilder{{"ID", "desc"}}),
					Offset(10),
					Limit(10),
					In([][]any{slice.ToAnySlice(number.Range(1, 1000, 1))}...),
				),
			},
			wantR: func() []post {
				r, err := Select[post](ctx, "select * from "+post{}.Table()+" where post_status='publish' and ID in ("+strings.Join(slice.Map(number.Range(1, 1000, 1), strconv.Itoa), ",")+") order by ID desc limit 10 offset 10 ")
				if err != nil {
					panic(err)
				}
				return r
			}(),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotR, err := Finds[post](tt.args.ctx, tt.args.q)
			if (err != nil) != tt.wantErr {
				t.Errorf("Findx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("Findx() gotR = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestChunkFind(t *testing.T) {
	type args struct {
		ctx      context.Context
		perLimit int
		q        QueryCondition
	}
	type testCase[T Model] struct {
		name    string
		args    args
		wantR   []T
		wantErr bool
	}
	n := 500
	tests := []testCase[post]{
		{
			name: "in,orderBy",
			args: args{
				ctx: ctx,
				q: Conditions(
					Where(SqlBuilder{{
						"post_status", "publish",
					}, {"ID", "in", ""}}),
					Order(SqlBuilder{{"ID", "desc"}}),
					In([][]any{slice.ToAnySlice(number.Range(1, n, 1))}...),
				),
				perLimit: 20,
			},
			wantR: func() []post {
				r, err := Select[post](ctx, "select * from "+post{}.Table()+" where post_status='publish' and ID in ("+strings.Join(slice.Map(number.Range(1, n, 1), strconv.Itoa), ",")+")  order by ID desc")
				if err != nil {
					panic(err)
				}
				return r
			}(),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotR, err := ChunkFind[post](tt.args.ctx, tt.args.perLimit, tt.args.q)
			if (err != nil) != tt.wantErr {
				t.Errorf("ChunkFind() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("ChunkFind() gotR = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestChunk(t *testing.T) {
	type args[T Model, R any] struct {
		ctx      context.Context
		perLimit int
		fn       func(rows T) (R, bool)
		q        QueryCondition
	}
	type testCase[T Model, R any] struct {
		name    string
		args    args[T, R]
		wantR   []R
		wantErr bool
	}
	n := 500
	tests := []testCase[post, uint64]{
		{
			name: "t1",
			args: args[post, uint64]{
				ctx:      ctx,
				perLimit: 20,
				fn: func(t post) (uint64, bool) {
					if t.Id > 300 {
						return t.Id, true
					}
					return 0, false
				},
				q: Conditions(
					Where(SqlBuilder{{
						"post_status", "publish",
					}, {"ID", "in", ""}}),
					Order(SqlBuilder{{"ID", "desc"}}),
					In([][]any{slice.ToAnySlice(number.Range(1, n, 1))}...),
				),
			},
			wantR: func() []uint64 {
				r, err := Select[post](ctx, "select * from "+post{}.Table()+" where post_status='publish' and ID in ("+strings.Join(slice.Map(number.Range(1, n, 1), strconv.Itoa), ",")+")  order by ID desc")
				if err != nil {
					panic(err)
				}
				return slice.FilterAndMap(r, func(t post) (uint64, bool) {
					if t.Id <= 300 {
						return 0, false
					}
					return t.Id, true
				})
			}(),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotR, err := Chunk[post](tt.args.ctx, tt.args.perLimit, tt.args.fn, tt.args.q)
			if (err != nil) != tt.wantErr {
				t.Errorf("Chunk() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("Chunk() gotR = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestPagination(t *testing.T) {
	type args struct {
		ctx context.Context
		q   QueryCondition
	}
	type testCase[T Model] struct {
		name    string
		args    args
		want    []T
		want1   int
		wantErr bool
	}
	tests := []testCase[post]{
		{
			name: "t1",
			args: args{
				ctx: ctx,
				q: Conditions(
					Where(SqlBuilder{
						{"ID", "in", ""},
					}),
					Page(1),
					Limit(5),
					In([][]any{slice.ToAnySlice(number.Range(431, 440, 1))}...),
				),
			},
			want: func() (r []post) {
				r, err := Select[post](ctx, "select * from "+post{}.Table()+" where ID in (?,?,?,?,?)", slice.ToAnySlice(number.Range(431, 435, 1))...)
				if err != nil && err != sql.ErrNoRows {
					panic(err)
				} else if err == sql.ErrNoRows {
					err = nil
				}
				return
			}(),
			want1:   10,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := Pagination[post](tt.args.ctx, tt.args.q)
			if (err != nil) != tt.wantErr {
				t.Errorf("Pagination() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Pagination() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Pagination() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestColumn(t *testing.T) {
	type args[V Model, T any] struct {
		ctx context.Context
		fn  func(V) (T, bool)
		q   QueryCondition
	}
	type testCase[V Model, T any] struct {
		name    string
		args    args[V, T]
		wantR   []T
		wantErr bool
	}
	tests := []testCase[post, uint64]{
		{
			name: "t1",
			args: args[post, uint64]{
				ctx: ctx,
				fn: func(t post) (uint64, bool) {
					return t.Id, true
				},
				q: Conditions(
					Where(SqlBuilder{
						{"ID", "<", "200", "int"},
					}),
				),
			},
			wantR: []uint64{63, 64, 190, 193},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotR, err := Column[post](tt.args.ctx, tt.args.fn, tt.args.q)
			if (err != nil) != tt.wantErr {
				t.Errorf("Column() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("Column() gotR = %v, want %v", gotR, tt.wantR)
			}
		})
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

func Test_getField(t *testing.T) {
	{
		name := "string"
		db := glob
		field := "option_value"
		q := Conditions(Where(SqlBuilder{{"option_name", "blogname"}}))
		wantR := "记录并见证自己的成长"
		wantErr := false
		t.Run(name, func(t *testing.T) {
			gotR, err := getField[options](db, ctx, field, q)
			if (err != nil) != wantErr {
				t.Errorf("getField() error = %v, wantErr %v", err, wantErr)
				return
			}
			if !reflect.DeepEqual(gotR, wantR) {
				t.Errorf("getField() gotR = %v, want %v", gotR, wantR)
			}
		})
	}

	{
		name := "t2"
		db := glob
		field := "option_id"
		q := Conditions(Where(SqlBuilder{{"option_name", "blogname"}}))
		wantR := "3"
		wantErr := false
		t.Run(name, func(t *testing.T) {
			gotR, err := getField[options](db, ctx, field, q)
			if (err != nil) != wantErr {
				t.Errorf("getField() error = %v, wantErr %v", err, wantErr)
				return
			}
			if !reflect.DeepEqual(gotR, wantR) {
				t.Errorf("getField() gotR = %v, want %v", gotR, wantR)
			}
		})
	}
	{
		name := "count(*)"
		db := glob
		field := "count(*)"
		q := Conditions()
		wantR := "385"
		wantErr := false
		t.Run(name, func(t *testing.T) {
			gotR, err := getField[options](db, ctx, field, q)
			if (err != nil) != wantErr {
				t.Errorf("getField() error = %v, wantErr %v", err, wantErr)
				return
			}
			if !reflect.DeepEqual(gotR, wantR) {
				t.Errorf("getField() gotR = %v, want %v", gotR, wantR)
			}
		})
	}
}

func Test_getToStringMap(t *testing.T) {
	type args struct {
		db  dbQuery
		ctx context.Context
		q   QueryCondition
	}
	tests := []struct {
		name    string
		args    args
		wantR   map[string]string
		wantErr bool
	}{
		{
			name: "t1",
			args: args{
				db:  glob,
				ctx: ctx,
				q:   Conditions(Where(SqlBuilder{{"option_name", "users_can_register"}})),
			},
			wantR: map[string]string{
				"option_id":    "5",
				"option_value": "0",
				"option_name":  "users_can_register",
				"autoload":     "yes",
			},
		},
		{
			name: "t2",
			args: args{
				db:  glob,
				ctx: ctx,
				q: Conditions(
					Where(SqlBuilder{{"option_name", "users_can_register"}}),
					Fields("option_id id"),
				),
			},
			wantR: map[string]string{
				"id": "5",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotR, err := getToStringMap[options](tt.args.db, tt.args.ctx, tt.args.q)
			if (err != nil) != tt.wantErr {
				t.Errorf("getToStringMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("getToStringMap() gotR = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func Test_findToStringMap(t *testing.T) {
	type args struct {
		db  dbQuery
		ctx context.Context
		q   QueryCondition
	}
	tests := []struct {
		name    string
		args    args
		wantR   []map[string]string
		wantErr bool
	}{
		{
			name: "t1",
			args: args{
				db:  glob,
				ctx: ctx,
				q:   Conditions(Where(SqlBuilder{{"option_id", "5"}})),
			},
			wantR: []map[string]string{{
				"option_id":    "5",
				"option_value": "0",
				"option_name":  "users_can_register",
				"autoload":     "yes",
			}},
			wantErr: false,
		},
		{
			name: "t2",
			args: args{
				db:  glob,
				ctx: ctx,
				q: Conditions(
					Where(SqlBuilder{{"option_id", "5"}}),
					Fields("option_value,option_name"),
				),
			},
			wantR: []map[string]string{{
				"option_value": "0",
				"option_name":  "users_can_register",
			}},
			wantErr: false,
		},
		{
			name: "t3",
			args: args{
				db:  glob,
				ctx: ctx,
				q: Conditions(
					Where(SqlBuilder{{"option_id", "5"}}),
					Fields("option_value v,option_name k"),
				),
			},
			wantR: []map[string]string{{
				"v": "0",
				"k": "users_can_register",
			}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotR, err := findToStringMap[options](tt.args.db, tt.args.ctx, tt.args.q)
			if (err != nil) != tt.wantErr {
				t.Errorf("findToStringMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("findToStringMap() gotR = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func Test_findScanner(t *testing.T) {
	type args[T Model] struct {
		db  dbQuery
		ctx context.Context
		fn  func(T)
		q   QueryCondition
	}
	type testCase[T Model] struct {
		name    string
		args    args[T]
		wantErr bool
	}
	tests := []testCase[options]{
		{
			name: "t1",
			args: args[options]{glob, ctx, func(t options) {
				fmt.Println(t)
			}, Conditions(Where(SqlBuilder{{"option_id", "<", "10", "int"}}))},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := findScanner[options](tt.args.db, tt.args.ctx, tt.args.fn, tt.args.q); (err != nil) != tt.wantErr {
				t.Errorf("findScanner() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func BenchmarkSqlxQueryXX(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r []options
		err := ddb.Select(&r, "select * from wp_options where option_id<100 and option_id>50")
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkScannerXX(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r []options
		err := findScanner[options](glob, ctx, func(t options) {
			r = append(r, t)
			//fmt.Println(t)
		}, Conditions(Where(SqlBuilder{{"option_id<100"}, {"option_id>50"}})))
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkFindsXX(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := finds[options](glob, ctx, Conditions(
			Where(SqlBuilder{{"option_id<100"}, {"option_id>50"}})),
		)
		if err != nil {
			panic(err)
		}
	}
}

func Test_gets(t *testing.T) {
	type args struct {
		db  dbQuery
		ctx context.Context
		q   QueryCondition
	}
	type testCase[T Model] struct {
		name    string
		args    args
		wantR   T
		wantErr bool
	}
	tests := []testCase[options]{
		{
			name: "t1",
			args: args{
				db:  glob,
				ctx: ctx,
				q:   Conditions(Where(SqlBuilder{{"option_name", "blogname"}})),
			},
			wantR:   options{3, "blogname", "记录并见证自己的成长", "yes"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotR, err := gets[options](tt.args.db, tt.args.ctx, tt.args.q)
			if (err != nil) != tt.wantErr {
				t.Errorf("gets() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("gets() gotR = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func Test_finds(t *testing.T) {
	type args struct {
		db  dbQuery
		ctx context.Context
		q   QueryCondition
	}
	type testCase[T Model] struct {
		name    string
		args    args
		wantR   []T
		wantErr bool
	}
	var u user
	tests := []testCase[options]{
		{
			name: "sub query",
			args: args{db: glob, ctx: ctx, q: Conditions(
				From("(select * from wp_options where option_id <100) a"),
				Where(SqlBuilder{{"option_id", ">", "50", "int"}}),
			)},
			wantR: func() []options {
				r, err := Select[options](ctx, "select * from (select * from wp_options where option_id <100) a where option_id>50")
				if err != nil {
					panic(err)
				}
				return r
			}(),
			wantErr: false,
		},
		{
			name: "mixed query",
			args: args{db: glob, ctx: ctx, q: Conditions(
				From("(select * from wp_options where option_id <100) a"),
				Where(SqlBuilder{
					{"u.ID", "<", "50", "int"}}),
				Join(SqlBuilder{
					{"left join", user.Table(u) + " u", "a.option_id=u.ID"},
				}),
				Fields("u.user_login autoload,option_name,option_value"),
			)},
			wantR: func() []options {
				r, err := Select[options](ctx, "select u.user_login autoload,option_name,option_value from (select * from wp_options where option_id <100) a  left join  wp_users u  on a.option_id=u.ID   where `u`.`ID`<50")
				if err != nil {
					panic(err)
				}
				return r
			}(),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotR, err := finds[options](tt.args.db, tt.args.ctx, tt.args.q)
			if (err != nil) != tt.wantErr {
				t.Errorf("finds() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("finds() gotR = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}
