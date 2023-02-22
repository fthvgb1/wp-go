package model

import "context"

func finds[T Model](db dbQuery[T], ctx context.Context, q *QueryCondition) ([]T, error) {
	s, args, err := BuildQuerySql[T](q)
	if err != nil {
		return nil, err
	}
	return db.Select(ctx, s, args...)
}

func scanners[T Model](db dbQuery[T], ctx context.Context, q *QueryCondition) ([]T, error) {
	s, args, err := BuildQuerySql[T](q)
	if err != nil {
		return nil, err
	}
	ctx = context.WithValue(ctx, "handle=>", "scanner")
	var r []T
	ctx = context.WithValue(ctx, "fn", func(t T) {
		r = append(r, t)
	})
	_, err = db.Select(ctx, s, args...)
	return r, err
}
