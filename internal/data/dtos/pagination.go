package dtos

type Pagination struct {
	Skip  int `json:"skip"  validate:"omitempty,gte=0"`
	Limit int `json:"limit" validate:"omitempty,gte=1,lte=100"`
}

// WithDefaults returns a copy of Pagination where the defaults have been set
// if they are at their zero values.
func (p Pagination) WithDefaults() Pagination {
	if p.Skip == 0 {
		p.Skip = 0
	}
	if p.Limit == 0 {
		p.Limit = 100
	}
	return p
}

type PaginationResponse[T any] struct {
	Total int `json:"total"`
	Items []T `json:"items"`
}
