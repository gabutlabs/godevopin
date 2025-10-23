package pkg

// Pagination structure for generic use
type Pagination struct {
	Page  int    `json:"page"`
	Limit int    `json:"limit"`
	Sort  string `json:"sort"`
	Total int64  `json:"total"`
}

// PaginatedResponse for generic use
type PaginatedResponse[T any] struct {
	Data       []T        `json:"data"`
	Pagination Pagination `json:"pagination"`
}

// ApplyPaginationToSlice applies pagination to a slice of any type by manually slicing the data
func ApplyPaginationToSlice[T any](data []T, page, limit int) []T {
	start := (page - 1) * limit
	if start >= len(data) {
		return []T{}
	}

	end := min(start+limit, len(data))

	return data[start:end]
}

// GetPaginationMetadata returns pagination metadata without performing database operations
func GetPaginationMetadata(total int64, page, limit int, sort string) Pagination {
	return Pagination{
		Page:  page,
		Limit: limit,
		Sort:  sort,
		Total: total,
	}
}
