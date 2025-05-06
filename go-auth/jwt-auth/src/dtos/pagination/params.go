package pagination

type PaginationParams struct {
	Field       string
	Value       any
	Result      interface{}
	Skip        int
	Limit       int
	SearchField string
	SearchValue string
	SortField   string
	SortOrder   int
}
