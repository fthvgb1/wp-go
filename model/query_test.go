package model

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/fthvgb1/wp-go/safety"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"reflect"
	"testing"
	"time"
)

type post struct {
	Id                  uint64    `gorm:"column:ID" db:"ID" json:"ID" form:"ID"`
	PostAuthor          uint64    `gorm:"column:post_author" db:"post_author" json:"post_author" form:"post_author"`
	PostDate            time.Time `gorm:"column:post_date" db:"post_date" json:"post_date" form:"post_date"`
	PostDateGmt         time.Time `gorm:"column:post_date_gmt" db:"post_date_gmt" json:"post_date_gmt" form:"post_date_gmt"`
	PostContent         string    `gorm:"column:post_content" db:"post_content" json:"post_content" form:"post_content"`
	PostTitle           string    `gorm:"column:post_title" db:"post_title" json:"post_title" form:"post_title"`
	PostExcerpt         string    `gorm:"column:post_excerpt" db:"post_excerpt" json:"post_excerpt" form:"post_excerpt"`
	PostStatus          string    `gorm:"column:post_status" db:"post_status" json:"post_status" form:"post_status"`
	CommentStatus       string    `gorm:"column:comment_status" db:"comment_status" json:"comment_status" form:"comment_status"`
	PingStatus          string    `gorm:"column:ping_status" db:"ping_status" json:"ping_status" form:"ping_status"`
	PostPassword        string    `gorm:"column:post_password" db:"post_password" json:"post_password" form:"post_password"`
	PostName            string    `gorm:"column:post_name" db:"post_name" json:"post_name" form:"post_name"`
	ToPing              string    `gorm:"column:to_ping" db:"to_ping" json:"to_ping" form:"to_ping"`
	Pinged              string    `gorm:"column:pinged" db:"pinged" json:"pinged" form:"pinged"`
	PostModified        time.Time `gorm:"column:post_modified" db:"post_modified" json:"post_modified" form:"post_modified"`
	PostModifiedGmt     time.Time `gorm:"column:post_modified_gmt" db:"post_modified_gmt" json:"post_modified_gmt" form:"post_modified_gmt"`
	PostContentFiltered string    `gorm:"column:post_content_filtered" db:"post_content_filtered" json:"post_content_filtered" form:"post_content_filtered"`
	PostParent          uint64    `gorm:"column:post_parent" db:"post_parent" json:"post_parent" form:"post_parent"`
	Guid                string    `gorm:"column:guid" db:"guid" json:"guid" form:"guid"`
	MenuOrder           int       `gorm:"column:menu_order" db:"menu_order" json:"menu_order" form:"menu_order"`
	PostType            string    `gorm:"column:post_type" db:"post_type" json:"post_type" form:"post_type"`
	PostMimeType        string    `gorm:"column:post_mime_type" db:"post_mime_type" json:"post_mime_type" form:"post_mime_type"`
	CommentCount        int64     `gorm:"column:comment_count" db:"comment_count" json:"comment_count" form:"comment_count"`
}

type user struct {
	Id                uint64    `gorm:"column:ID" db:"ID" json:"ID"`
	UserLogin         string    `gorm:"column:user_login" db:"user_login" json:"user_login"`
	UserPass          string    `gorm:"column:user_pass" db:"user_pass" json:"user_pass"`
	UserNicename      string    `gorm:"column:user_nicename" db:"user_nicename" json:"user_nicename"`
	UserEmail         string    `gorm:"column:user_email" db:"user_email" json:"user_email"`
	UserUrl           string    `gorm:"column:user_url" db:"user_url" json:"user_url"`
	UserRegistered    time.Time `gorm:"column:user_registered" db:"user_registered" json:"user_registered"`
	UserActivationKey string    `gorm:"column:user_activation_key" db:"user_activation_key" json:"user_activation_key"`
	UserStatus        int       `gorm:"column:user_status" db:"user_status" json:"user_status"`
	DisplayName       string    `gorm:"column:display_name" db:"display_name" json:"display_name"`
}

type termTaxonomy struct {
	TermTaxonomyId uint64 `gorm:"column:term_taxonomy_id" db:"term_taxonomy_id" json:"term_taxonomy_id" form:"term_taxonomy_id"`
	TermId         uint64 `gorm:"column:term_id" db:"term_id" json:"term_id" form:"term_id"`
	Taxonomy       string `gorm:"column:taxonomy" db:"taxonomy" json:"taxonomy" form:"taxonomy"`
	Description    string `gorm:"column:description" db:"description" json:"description" form:"description"`
	Parent         uint64 `gorm:"column:parent" db:"parent" json:"parent" form:"parent"`
	Count          int64  `gorm:"column:count" db:"count" json:"count" form:"count"`
}

type terms struct {
	TermId    uint64 `gorm:"column:term_id" db:"term_id" json:"term_id" form:"term_id"`
	Name      string `gorm:"column:name" db:"name" json:"name" form:"name"`
	Slug      string `gorm:"column:slug" db:"slug" json:"slug" form:"slug"`
	TermGroup int64  `gorm:"column:term_group" db:"term_group" json:"term_group" form:"term_group"`
}

func (t terms) PrimaryKey() string {
	return "term_id"
}
func (t terms) Table() string {
	return "wp_terms"
}

func (w termTaxonomy) PrimaryKey() string {
	return "term_taxonomy_id"
}

func (w termTaxonomy) Table() string {
	return "wp_term_taxonomy"
}

func (u user) Table() string {
	return "wp_users"
}

func (u user) PrimaryKey() string {
	return "ID"
}

func (p post) PrimaryKey() string {
	return "ID"
}

func (p post) Table() string {
	return "wp_posts"
}

var ctx = context.Background()

var glob *SqlxQuery
var ddb *sqlx.DB

func init() {
	db, err := sqlx.Open("mysql", "root:root@tcp(192.168.66.47:3306)/wordpress?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	glob = NewSqlxQuery(safety.NewVar(db), NewUniversalDb(func(ctx2 context.Context, a any, s string, a2 ...any) error {
		x := FormatSql(s, a2...)
		fmt.Println(x)
		return glob.Selects(ctx2, a, s, a2...)
	}, func(ctx2 context.Context, a any, s string, a2 ...any) error {
		x := FormatSql(s, a2...)
		fmt.Println(x)
		return glob.Gets(ctx2, a, s, a2...)
	}))
	ddb = db
	InitDB(glob)
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
		post
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
				r, err := Select[posts](ctx, "select post_status,count(*) n from "+post{}.Table()+" where ID<1000 group by post_status having n>1")
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
					{"a", "left join", user{}.Table() + " b", "a.post_author=b.ID"},
					{"left join", "wp_term_relationships c", "a.Id=c.object_id"},
					{"left join", termTaxonomy{}.Table() + " d", "c.term_taxonomy_id=d.term_taxonomy_id"},
					{"left join", terms{}.Table() + " e", "d.term_id=e.term_id"},
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
		want    post
		wantErr bool
	}{
		{
			name: "t1",
			args: args{
				1,
			},
			want: func() post {
				r, err := Get[post](ctx, "select * from "+post{}.Table()+" where ID=?", 1)
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
			got, err := FindOneById[post](ctx, tt.args.id)
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
		want    post
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
			want: func() post {
				r, err := Get[post](ctx, "select * from "+post{}.Table()+" where post_status='publish' order by ID desc limit 1")
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
			got, err := FirstOne[post](ctx, tt.args.where, tt.args.fields, tt.args.order, tt.args.in...)
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
		want    post
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
			want: func() post {
				r, err := Get[post](ctx, "select * from "+post{}.Table()+" where post_status='publish' order by  "+post{}.PrimaryKey()+" desc limit 1")
				if err != nil {
					panic(err)
				}
				return r
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LastOne[post](ctx, tt.args.where, tt.args.fields, tt.args.in...)
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
		want    []post
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
			want: func() (r []post) {
				r, err := Select[post](ctx, "select * from "+post{}.Table()+" where ID in (?,?)", 1, 2)
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
			got, err := SimpleFind[post](ctx, tt.args.where, tt.args.fields, tt.args.in...)
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

func Test_pagination(t *testing.T) {
	type args struct {
		db       dbQuery
		ctx      context.Context
		q        QueryCondition
		page     int
		pageSize int
	}
	type testCase[T Model] struct {
		name      string
		args      args
		wantR     []T
		wantTotal int
		wantErr   bool
	}
	tests := []testCase[post]{
		{
			name: "t1",
			args: args{
				db:  glob,
				ctx: ctx,
				q: QueryCondition{
					Fields: "post_type,count(*) ID",
					Group:  "post_type",
					Having: SqlBuilder{{"ID", ">", "1", "int"}},
				},
				page:     1,
				pageSize: 2,
			},
			wantR: func() (r []post) {

				err := glob.Selects(ctx, &r, "select post_type,count(*) ID from wp_posts group by post_type having `ID`> 1 limit 2")
				if err != nil {
					panic(err)
				}
				return r
			}(),
			wantTotal: 7,
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotR, gotTotal, err := pagination[post](tt.args.db, tt.args.ctx, tt.args.q, tt.args.page, tt.args.pageSize)
			if (err != nil) != tt.wantErr {
				t.Errorf("pagination() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("pagination() gotR = %v, want %v", gotR, tt.wantR)
			}
			if gotTotal != tt.wantTotal {
				t.Errorf("pagination() gotTotal = %v, want %v", gotTotal, tt.wantTotal)
			}
		})
	}
}

func Test_paginationToMap(t *testing.T) {
	type args struct {
		db       dbQuery
		ctx      context.Context
		q        QueryCondition
		page     int
		pageSize int
	}
	tests := []struct {
		name      string
		args      args
		wantR     []map[string]string
		wantTotal int
		wantErr   bool
	}{
		{
			name: "t1",
			args: args{
				db:  glob,
				ctx: ctx,
				q: QueryCondition{
					Fields: "ID",
					Where:  SqlBuilder{{"ID < 200"}},
				},
				page:     1,
				pageSize: 2,
			},
			wantR:     []map[string]string{{"ID": "63"}, {"ID": "64"}},
			wantTotal: 4,
		},
		{
			name: "t2",
			args: args{
				db:  glob,
				ctx: ctx,
				q: QueryCondition{
					Fields: "ID",
					Where:  SqlBuilder{{"ID < 200"}},
				},
				page:     2,
				pageSize: 2,
			},
			wantR:     []map[string]string{{"ID": "190"}, {"ID": "193"}},
			wantTotal: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotR, gotTotal, err := PaginationToMap[post](tt.args.ctx, tt.args.q, tt.args.page, tt.args.pageSize)
			fmt.Println(gotR, gotTotal, err)
			gotR, gotTotal, err = PaginationToMapFromDB[post](tt.args.db, tt.args.ctx, tt.args.q, tt.args.page, tt.args.pageSize)
			fmt.Println(gotR, gotTotal, err)
			gotR, gotTotal, err = paginationToMap[post](tt.args.db, tt.args.ctx, tt.args.q, tt.args.page, tt.args.pageSize)

			if (err != nil) != tt.wantErr {
				t.Errorf("paginationToMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("paginationToMap() gotR = %v, want %v", gotR, tt.wantR)
			}
			if gotTotal != tt.wantTotal {
				t.Errorf("paginationToMap() gotTotal = %v, want %v", gotTotal, tt.wantTotal)
			}
		})
	}
}
