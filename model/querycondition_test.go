package model

import (
	"context"
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
