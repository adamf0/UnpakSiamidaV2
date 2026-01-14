package application

import "UnpakSiamida/common/domain"

type GetAllRenstraNilaisQuery struct{
	Search       string
    SearchFilters []domain.SearchFilter
    Page         *int
    Limit        *int
}
