package pagination

type PaginationResult[T any] struct {
	Items       []T `json:"items"`
	PageCurrent int `json:"pageCurrent"`
	PaginationResultContext
}
