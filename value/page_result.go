package value

type PageResult[T any] struct {
	TotalCount uint64
	Items      []T
}
