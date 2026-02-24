package application

import (
	commondomain "UnpakSiamida/common/domain"
	domaintemplatedokumentambahan "UnpakSiamida/modules/templatedokumentambahan/domain"
	"context"
	"time"
)

type GetAllTemplateDokumenTambahansQueryHandler struct {
	Repo domaintemplatedokumentambahan.ITemplateDokumenTambahanRepository
}

func (h *GetAllTemplateDokumenTambahansQueryHandler) Handle(
	ctx context.Context,
	q GetAllTemplateDokumenTambahansQuery,
) (commondomain.Paged[domaintemplatedokumentambahan.TemplateDokumenTambahanDefault], error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	templatedokumentambahans, total, err := h.Repo.GetAll(
		ctx,
		q.Search,
		q.SearchFilters,
		q.Page,
		q.Limit,
	)
	if err != nil {
		return commondomain.Paged[domaintemplatedokumentambahan.TemplateDokumenTambahanDefault]{}, err
	}

	currentPage := 1
	totalPages := 1

	if q.Page != nil {
		currentPage = *q.Page
	}
	if q.Limit != nil && *q.Limit > 0 {
		totalPages = int((total + int64(*q.Limit) - 1) / int64(*q.Limit))
	}

	return commondomain.Paged[domaintemplatedokumentambahan.TemplateDokumenTambahanDefault]{
		Data:        templatedokumentambahans,
		Total:       total,
		CurrentPage: currentPage,
		TotalPages:  totalPages,
	}, nil
}
