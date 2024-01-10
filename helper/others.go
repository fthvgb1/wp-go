package helper

type PaginationData[T any] struct {
	Data     []T
	TotalRaw int
}
