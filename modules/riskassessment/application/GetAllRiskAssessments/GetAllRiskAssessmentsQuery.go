package application

import (
	commondomain "UnpakSiamida/common/domain"
)

type GetAllRiskAssessmentsQuery struct {
	Search        string
	SearchFilters []commondomain.SearchFilter
	Page          *int
	Limit         *int
}
