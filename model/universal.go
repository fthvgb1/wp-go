package model

import "context"

type UniversalDb[T any] struct {
	selects QuerySelect[T]
	gets    QueryGet[T]
}

func (u *UniversalDb[T]) Select(ctx context.Context, s string, a ...any) ([]T, error) {
	return u.selects(ctx, s, a...)
}

func (u *UniversalDb[T]) Get(ctx context.Context, s string, a ...any) (T, error) {
	return u.gets(ctx, s, a...)
}
