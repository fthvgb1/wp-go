package helper

import "context"

type Pagination[T any] interface {
	SetData(ctx context.Context, data []T)
	GetData(ctx context.Context) []T
	TotalRaws(ctx context.Context) int
}

type PaginationData[T any] struct {
	Data     []T
	TotalRaw int
}

func (p *PaginationData[T]) SetData(ctx context.Context, data []T) {
	p.Data = data
}

func (p *PaginationData[T]) GetData(ctx context.Context) []T {
	return p.Data
}

func (p *PaginationData[T]) TotalRaws(ctx context.Context) int {
	return p.TotalRaw
}
