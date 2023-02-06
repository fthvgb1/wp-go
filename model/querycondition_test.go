package model

import (
	"context"
	"database/sql"
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
		q   *QueryCondition
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
		q        *QueryCondition
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
		q        *QueryCondition
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
		q   *QueryCondition
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
		q   *QueryCondition
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
