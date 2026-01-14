package application

import "UnpakSiamida/common/domain"

type GetAllDokumenTambahansQuery struct{
	Search       string
    SearchFilters []domain.SearchFilter
    Page         *int
    Limit        *int
}
