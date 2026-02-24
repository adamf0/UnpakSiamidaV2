package application

import (
	"context"
	"time"

	commondomain "UnpakSiamida/common/domain"
	domainlaporan "UnpakSiamida/modules/laporan/domain"

	"github.com/google/uuid"
)

type GetMonitoringIndikatorByIndikatorTahunQueryHandler struct {
	Repo domainlaporan.ILaporanRepository
}

func (h *GetMonitoringIndikatorByIndikatorTahunQueryHandler) Handle(
	ctx context.Context,
	q GetMonitoringIndikatorByIndikatorTahunQuery,
) (commondomain.Paged[domainlaporan.MonitoringIndikator], error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	uuidIndikator, err := uuid.Parse(q.IndikatorUuid)
	if err != nil {
		return commondomain.Paged[domainlaporan.MonitoringIndikator]{}, domainlaporan.InvalidIndikator()
	}

	monitoringIndikator, total, err := h.Repo.GetMonitoringIndikatorByIndikatorTahun(ctx, uuidIndikator, q.Tahun, q.Page, q.Limit)
	if err != nil {
		return commondomain.Paged[domainlaporan.MonitoringIndikator]{}, err
	}

	currentPage := 1
	totalPages := 1
	if q.Page != nil {
		currentPage = *q.Page
	}
	if q.Limit != nil && *q.Limit > 0 {
		totalPages = int((total + int64(*q.Limit) - 1) / int64(*q.Limit))
	}

	return commondomain.Paged[domainlaporan.MonitoringIndikator]{
		Data:        monitoringIndikator,
		Total:       total,
		CurrentPage: currentPage,
		TotalPages:  totalPages,
	}, nil
}
