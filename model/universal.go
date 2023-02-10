package model

import "context"

type UniversalDb struct {
	selects QueryFn
	gets    QueryFn
}

func NewUniversalDb(selects QueryFn, gets QueryFn) UniversalDb {
	return UniversalDb{selects: selects, gets: gets}
}

func (u UniversalDb) Select(ctx context.Context, a any, s string, args ...any) error {
	return u.selects(ctx, a, s, args...)
}

func (u UniversalDb) Get(ctx context.Context, a any, s string, args ...any) error {
	return u.gets(ctx, a, s, args...)
}
