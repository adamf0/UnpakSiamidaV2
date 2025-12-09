package application

import "UnpakSiamida/common/domain"

type GetAllJenisFilesQuery struct{
	Search       string
    SearchFilters []domain.SearchFilter
    Page         *int
    Limit        *int
}
