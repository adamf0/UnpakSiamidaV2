package GetDashboardProkerMonitoring

import (
	domain "UnpakSiamida/modules/prokerv3/domain"
	"context"
)

type GetDashboardProkerMonitoringQueryHandler struct {
	Repo domain.IProkerV3Repository
}

func (h *GetDashboardProkerMonitoringQueryHandler) Handle(ctx context.Context, q GetDashboardProkerMonitoringQuery) (*domain.DashboardMonitoringResult, error) {
	return h.Repo.GetDashboardProkerMonitoring(ctx, q.Filter)
}
