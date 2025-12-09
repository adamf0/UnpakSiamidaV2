package application

import "UnpakSiamida/common/domain"

type GetAllTemplateDokumenTambahansQuery struct{
	Search       string
    SearchFilters []domain.SearchFilter
    Page         *int
    Limit        *int
}
