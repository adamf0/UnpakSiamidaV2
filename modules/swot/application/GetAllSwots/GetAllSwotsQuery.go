package application

import (
	commondomain "UnpakSiamida/common/domain"
)

type GetAllSwotsQuery struct {
	Search        string
	SearchFilters []commondomain.SearchFilter
	Page          *int
	Limit         *int
}
