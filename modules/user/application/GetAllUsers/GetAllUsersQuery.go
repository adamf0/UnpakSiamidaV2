package application

import "UnpakSiamida/common/domain"

type GetAllUsersQuery struct{
	Search       string
    SearchFilters []domain.SearchFilter
    Page         *int
    Limit        *int
}
