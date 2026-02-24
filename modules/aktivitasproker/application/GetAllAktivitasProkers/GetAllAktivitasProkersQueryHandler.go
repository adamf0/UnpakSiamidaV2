package application

import (
	commondomain "UnpakSiamida/common/domain"
	domainAktivitasProker "UnpakSiamida/modules/aktivitasproker/domain"
	"context"
	"time"
)

type GetAllAktivitasProkersQueryHandler struct {
	Repo domainAktivitasProker.IAktivitasProkerRepository
}

func (h *GetAllAktivitasProkersQueryHandler) Handle(
	ctx context.Context,
	q GetAllAktivitasProkersQuery,
) (commondomain.Paged[domainAktivitasProker.AktivitasProkerDefault], error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	AktivitasProkers, total, err := h.Repo.GetAll(
		ctx,
		q.Search,
		q.SearchFilters,
		q.Page,
		q.Limit,
	)
	if err != nil {
		return commondomain.Paged[domainAktivitasProker.AktivitasProkerDefault]{}, err
	}

	currentPage := 1
	totalPages := 1

	if q.Page != nil {
		currentPage = *q.Page
	}
	if q.Limit != nil && *q.Limit > 0 {
		totalPages = int((total + int64(*q.Limit) - 1) / int64(*q.Limit))
	}

	return commondomain.Paged[domainAktivitasProker.AktivitasProkerDefault]{
		Data:        AktivitasProkers,
		Total:       total,
		CurrentPage: currentPage,
		TotalPages:  totalPages,
	}, nil
}
