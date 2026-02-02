package application

import (
	"context"
	"time"

	domainlaporan "UnpakSiamida/modules/laporan/domain"

	"github.com/google/uuid"
)

type GetMonitoringProkerByTargetTahunQueryHandler struct {
	Repo domainlaporan.ILaporanRepository
}

func (h *GetMonitoringProkerByTargetTahunQueryHandler) Handle(
	ctx context.Context,
	q GetMonitoringProkerByTargetTahunQuery,
) (domainlaporan.Paged[domainlaporan.MonitoringProker], error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	uuidTahun, err := uuid.Parse(q.TahunUuid)
	if err != nil {
		return domainlaporan.Paged[domainlaporan.MonitoringProker]{}, domainlaporan.InvalidTahun()
	}

	uuidTarget, err := uuid.Parse(q.TargetUuid)
	if err != nil {
		return domainlaporan.Paged[domainlaporan.MonitoringProker]{}, domainlaporan.InvalidTarget()
	}

	monitoringProker, total, err := h.Repo.GetMonitoringByTargetTahun(ctx, uuidTarget, uuidTahun, q.Page, q.Limit)
	if err != nil {
		return domainlaporan.Paged[domainlaporan.MonitoringProker]{}, err
	}

	currentPage := 1
	totalPages := 1
	if q.Page != nil {
		currentPage = *q.Page
	}
	if q.Limit != nil && *q.Limit > 0 {
		totalPages = int((total + int64(*q.Limit) - 1) / int64(*q.Limit))
	}

	return domainlaporan.Paged[domainlaporan.MonitoringProker]{
		Data:        monitoringProker,
		Total:       total,
		CurrentPage: currentPage,
		TotalPages:  totalPages,
	}, nil
}
