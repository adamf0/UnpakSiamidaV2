package application

import "UnpakSiamida/common/domain"

type GetAllKtssQuery struct{
	Search       string
    SearchFilters []domain.SearchFilter
    Page         *int
    Limit        *int
}
