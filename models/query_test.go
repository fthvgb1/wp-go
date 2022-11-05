package models

import (
	"github/fthvgb1/wp-go/config"
	"github/fthvgb1/wp-go/db"
	"github/fthvgb1/wp-go/models/wp"
	"reflect"
	"testing"
)

func init() {
	err := config.InitConfig("../config.yaml")
	if err != nil {
		panic(err)
	}
	err = db.InitDb()
	if err != nil {
		panic(err)
	}
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
	tests := []struct {
		name    string
		args    args
		wantR   []wp.Posts
		wantErr bool
	}{
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotR, err := Find[wp.Posts](tt.args.where, tt.args.fields, tt.args.group, tt.args.order, tt.args.join, tt.args.having, tt.args.limit, tt.args.in...)
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
	r, err := Get[wp.Posts]("select * from "+wp.Posts{}.Table()+" where ID=?", 1)
	if err != nil {
		panic(err)
	}
	tests := []struct {
		name    string
		args    args
		want    wp.Posts
		wantErr bool
	}{
		{
			name: "t1",
			args: args{
				1,
			},
			want:    r,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FindOneById[wp.Posts](tt.args.id)
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
	r, err := Get[wp.Posts]("select * from " + wp.Posts{}.Table() + " where post_status='publish' order by ID desc")
	if err != nil {
		panic(err)
	}
	tests := []struct {
		name    string
		args    args
		want    wp.Posts
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
			want:    r,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FirstOne[wp.Posts](tt.args.where, tt.args.fields, tt.args.order, tt.args.in...)
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

func TestGet(t *testing.T) {
	type args struct {
		sql    string
		params []any
	}
	tests := []struct {
		name    string
		args    args
		wantR   wp.Posts
		wantErr bool
	}{
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotR, err := Get[wp.Posts](tt.args.sql, tt.args.params...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("Get() gotR = %v, want %v", gotR, tt.wantR)
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
		want    wp.Posts
		wantErr bool
	}{
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LastOne[wp.Posts](tt.args.where, tt.args.fields, tt.args.in...)
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

func TestSelect(t *testing.T) {
	type args struct {
		sql    string
		params []any
	}
	tests := []struct {
		name    string
		args    args
		want    []wp.Posts
		wantErr bool
	}{
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Select[wp.Posts](tt.args.sql, tt.args.params...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Select() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Select() got = %v, want %v", got, tt.want)
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
		want    []wp.Posts
		wantErr bool
	}{
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SimpleFind[wp.Posts](tt.args.where, tt.args.fields, tt.args.in...)
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
		wantR     []wp.Posts
		wantTotal int
		wantErr   bool
	}{
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotR, gotTotal, err := SimplePagination[wp.Posts](tt.args.where, tt.args.fields, tt.args.group, tt.args.page, tt.args.pageSize, tt.args.order, tt.args.join, tt.args.having, tt.args.in...)
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
