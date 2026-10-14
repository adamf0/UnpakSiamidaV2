package application

import (
	domainRisk "UnpakSiamida/modules/riskassessment/domain"
	"context"
	"time"
)

type GetRiskAssessmentByUuidQueryHandler struct {
	Repo domainRisk.IRiskAssessmentRepository
}

func (h *GetRiskAssessmentByUuidQueryHandler) Handle(
	ctx context.Context,
	q GetRiskAssessmentByUuidQuery,
) (*domainRisk.RiskAssessmentDefault, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	item, err := h.Repo.GetDefaultByUuid(ctx, q.UUID)
	if err != nil {
		return nil, domainRisk.NotFound(q.UUID.String())
	}

	return item, nil
}
