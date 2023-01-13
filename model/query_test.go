package model

import (
	"context"
	"database/sql"
	"github/fthvgb1/wp-go/helper"
	"github/fthvgb1/wp-go/internal/pkg/config"
	"github/fthvgb1/wp-go/internal/pkg/db"
	models2 "github/fthvgb1/wp-go/internal/pkg/models"
	"reflect"
	"testing"
)

var ctx = context.Background()

func init() {
	err := config.InitConfig("../config.yaml")
	if err != nil {
		panic(err)
	}
	err = db.InitDb()
	if err != nil {
		panic(err)
	}
	InitDB(db.NewSqlxDb(db.Db))
}
func TestFind(t *testing.T) {
	type args struct {
		where  ParseWhere
		fields string
		group  string
		order  SqlBuilder
		join   SqlBuilder
		having SqlBuilder
		limit  int
		in     [][]any
	}
	type posts struct {
		models2.Posts
		N int `db:"n"`
	}
	tests := []struct {
		name    string
		args    args
		wantR   []posts
		wantErr bool
	}{
		{
			name: "in,orderBy",
			args: args{
				where: SqlBuilder{{
					"post_status", "publish",
				}, {"ID", "in", ""}},
				fields: "*",
				group:  "",
				order:  SqlBuilder{{"ID", "desc"}},
				join:   nil,
				having: nil,
				limit:  0,
				in:     [][]any{{1, 2, 3, 4}},
			},
			wantR: func() []posts {
				r, err := Select[posts](ctx, "select * from "+posts{}.Table()+" where post_status='publish' and ID in (1,2,3,4) order by ID desc")
				if err != nil {
					panic(err)
				}
				return r
			}(),
			wantErr: false,
		},
		{
			name: "or",
			args: args{
				where: SqlBuilder{{
					"and", "ID", "=", "1", "int",
				}, {"or", "ID", "=", "2", "int"}},
				fields: "*",
				group:  "",
				order:  nil,
				join:   nil,
				having: nil,
				limit:  0,
				in:     nil,
			},
			wantR: func() []posts {
				r, err := Select[posts](ctx, "select * from "+posts{}.Table()+" where (ID=1 or ID=2)")
				if err != nil {
					panic(err)
				}
				return r
			}(),
		},
		{
			name: "group,having",
			args: args{
				where: SqlBuilder{
					{"ID", "<", "1000", "int"},
				},
				fields: "post_status,count(*) n",
				group:  "post_status",
				order:  nil,
				join:   nil,
				having: SqlBuilder{
					{"n", ">", "1"},
				},
				limit: 0,
				in:    nil,
			},
			wantR: func() []posts {
				r, err := Select[posts](ctx, "select post_status,count(*) n from "+models2.Posts{}.Table()+" where ID<1000 group by post_status having n>1")
				if err != nil {
					panic(err)
				}
				return r
			}(),
		},
		{
			name: "or、多个in",
			args: args{
				where: SqlBuilder{
					{"and", "ID", "in", "", "", "or", "ID", "in", "", ""},
					{"or", "post_status", "=", "publish", "", "and", "post_status", "=", "closed", ""},
				},
				fields: "*",
				group:  "",
				order:  nil,
				join:   nil,
				having: nil,
				limit:  0,
				in:     [][]any{{1, 2, 3}, {4, 5, 6}},
			},
			wantR: func() []posts {
				r, err := Select[posts](ctx, "select * from "+posts{}.Table()+" where (ID in (1,2,3) or ID in (4,5,6)) or (post_status='publish' and post_status='closed')")
				if err != nil {
					panic(err)
				}
				return r
			}(),
		},
		{
			name: "all",
			args: args{
				where: SqlBuilder{
					{"b.user_login", "in", ""},
					{"and", "a.post_type", "=", "post", "", "or", "a.post_type", "=", "page", ""},
					{"a.comment_count", ">", "0", "int"},
					{"a.post_status", "publish"},
					{"e.name", "in", ""},
					{"d.taxonomy", "category"},
				},
				fields: "post_author,count(*) n",
				group:  "a.post_author",
				order:  SqlBuilder{{"n", "desc"}},
				join: SqlBuilder{
					{"a", "left join", models2.Users{}.Table() + " b", "a.post_author=b.ID"},
					{"left join", "wp_term_relationships c", "a.Id=c.object_id"},
					{"left join", models2.TermTaxonomy{}.Table() + " d", "c.term_taxonomy_id=d.term_taxonomy_id"},
					{"left join", models2.Terms{}.Table() + " e", "d.term_id=e.term_id"},
				},
				having: SqlBuilder{{"n", ">", "0", "int"}},
				limit:  10,
				in:     [][]any{{"test", "test2"}, {"web", "golang", "php"}},
			},
			wantR: func() []posts {
				r, err := Select[posts](ctx, "select post_author,count(*) n from wp_posts a left join wp_users b on a.post_author=b.ID left join  wp_term_relationships c  on a.Id=c.object_id  left join  wp_term_taxonomy d  on c.term_taxonomy_id=d.term_taxonomy_id  left join  wp_terms e on d.term_id=e.term_id where b.user_login in ('test','test2') and b.user_status=0 and (a.post_type='post' or a.post_type='page') and a.comment_count>0 and a.post_status='publish' and e.name in ('web','golang','php') and d.taxonomy='category' group by post_author having n > 0 order by n desc limit 10")
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
			gotR, err := Find[posts](ctx, tt.args.where, tt.args.fields, tt.args.group, tt.args.order, tt.args.join, tt.args.having, tt.args.limit, tt.args.in...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Find() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("Find() gotR = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestFindOneById(t *testing.T) {
	type args struct {
		id int
	}

	tests := []struct {
		name    string
		args    args
		want    models2.Posts
		wantErr bool
	}{
		{
			name: "t1",
			args: args{
				1,
			},
			want: func() models2.Posts {
				r, err := Get[models2.Posts](ctx, "select * from "+models2.Posts{}.Table()+" where ID=?", 1)
				if err != nil && err != sql.ErrNoRows {
					panic(err)
				} else if err == sql.ErrNoRows {
					err = nil
				}
				return r
			}(),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FindOneById[models2.Posts](ctx, tt.args.id)
			if err == sql.ErrNoRows {
				err = nil
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("FindOneById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindOneById() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFirstOne(t *testing.T) {
	type args struct {
		where  ParseWhere
		fields string
		order  SqlBuilder
		in     [][]any
	}
	tests := []struct {
		name    string
		args    args
		want    models2.Posts
		wantErr bool
	}{
		{
			name: "t1",
			args: args{
				where:  SqlBuilder{{"post_status", "publish"}},
				fields: "*",
				order:  SqlBuilder{{"ID", "desc"}},
				in:     nil,
			},
			wantErr: false,
			want: func() models2.Posts {
				r, err := Get[models2.Posts](ctx, "select * from "+models2.Posts{}.Table()+" where post_status='publish' order by ID desc limit 1")
				if err != nil && err != sql.ErrNoRows {
					panic(err)
				} else if err == sql.ErrNoRows {
					err = nil
				}
				return r
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FirstOne[models2.Posts](ctx, tt.args.where, tt.args.fields, tt.args.order, tt.args.in...)
			if (err != nil) != tt.wantErr {
				t.Errorf("FirstOne() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FirstOne() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLastOne(t *testing.T) {
	type args struct {
		where  ParseWhere
		fields string
		in     [][]any
	}
	tests := []struct {
		name    string
		args    args
		want    models2.Posts
		wantErr bool
	}{
		{
			name: "t1",
			args: args{
				where: SqlBuilder{{
					"post_status", "publish",
				}},
				fields: "*",
				in:     nil,
			},
			want: func() models2.Posts {
				r, err := Get[models2.Posts](ctx, "select * from "+models2.Posts{}.Table()+" where post_status='publish' order by  "+models2.Posts{}.PrimaryKey()+" desc limit 1")
				if err != nil {
					panic(err)
				}
				return r
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LastOne[models2.Posts](ctx, tt.args.where, tt.args.fields, tt.args.in...)
			if (err != nil) != tt.wantErr {
				t.Errorf("LastOne() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LastOne() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSimpleFind(t *testing.T) {
	type args struct {
		where  ParseWhere
		fields string
		in     [][]any
	}
	tests := []struct {
		name    string
		args    args
		want    []models2.Posts
		wantErr bool
	}{
		{
			name: "t1",
			args: args{
				where: SqlBuilder{
					{"ID", "in", ""},
				},
				fields: "*",
				in:     [][]any{{1, 2}},
			},
			want: func() (r []models2.Posts) {
				r, err := Select[models2.Posts](ctx, "select * from "+models2.Posts{}.Table()+" where ID in (?,?)", 1, 2)
				if err != nil && err != sql.ErrNoRows {
					panic(err)
				} else if err == sql.ErrNoRows {
					err = nil
				}
				return
			}(),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SimpleFind[models2.Posts](ctx, tt.args.where, tt.args.fields, tt.args.in...)
			if (err != nil) != tt.wantErr {
				t.Errorf("SimpleFind() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SimpleFind() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSimplePagination(t *testing.T) {
	type args struct {
		where    ParseWhere
		fields   string
		group    string
		page     int
		pageSize int
		order    SqlBuilder
		join     SqlBuilder
		having   SqlBuilder
		in       [][]any
	}
	tests := []struct {
		name      string
		args      args
		wantR     []models2.Posts
		wantTotal int
		wantErr   bool
	}{
		{
			name: "t1",
			args: args{
				where: SqlBuilder{
					{"ID", "in", ""},
				},
				fields:   "*",
				group:    "",
				page:     1,
				pageSize: 5,
				order:    nil,
				join:     nil,
				having:   nil,
				in:       [][]any{helper.SliceMap[int, any](helper.RangeSlice(431, 440, 1), helper.ToAny[int])},
			},
			wantR: func() (r []models2.Posts) {
				r, err := Select[models2.Posts](ctx, "select * from "+models2.Posts{}.Table()+" where ID in (?,?,?,?,?)", helper.SliceMap[int, any](helper.RangeSlice(431, 435, 1), helper.ToAny[int])...)
				if err != nil && err != sql.ErrNoRows {
					panic(err)
				} else if err == sql.ErrNoRows {
					err = nil
				}
				return
			}(),
			wantTotal: 10,
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotR, gotTotal, err := SimplePagination[models2.Posts](ctx, tt.args.where, tt.args.fields, tt.args.group, tt.args.page, tt.args.pageSize, tt.args.order, tt.args.join, tt.args.having, tt.args.in...)
			if (err != nil) != tt.wantErr {
				t.Errorf("SimplePagination() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("SimplePagination() gotR = %v, want %v", gotR, tt.wantR)
			}
			if gotTotal != tt.wantTotal {
				t.Errorf("SimplePagination() gotTotal = %v, want %v", gotTotal, tt.wantTotal)
			}
		})
	}
}
