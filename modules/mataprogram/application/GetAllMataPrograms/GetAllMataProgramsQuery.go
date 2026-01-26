package application

import "UnpakSiamida/common/domain"

type GetAllMataProgramsQuery struct {
	Search        string
	SearchFilters []domain.SearchFilter
	Page          *int
	Limit         *int
}
