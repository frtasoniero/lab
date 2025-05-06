package pagination

type PaginationResultContext struct {
	TotalPages  int   `json:"totalPages"`
	HasNextPage bool  `json:"hasNextPage"`
	TotalItems  int   `json:"totalItems"`
	Err         error `json:"Err,omitempty"`
}
