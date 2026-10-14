package application

import (
	commondomain "UnpakSiamida/common/domain"
	domainRisk "UnpakSiamida/modules/riskassessment/domain"
	"context"
	"time"
)

type GetAllRiskAssessmentsQueryHandler struct {
	Repo domainRisk.IRiskAssessmentRepository
}

func (h *GetAllRiskAssessmentsQueryHandler) Handle(
	ctx context.Context,
	q GetAllRiskAssessmentsQuery,
) (commondomain.Paged[domainRisk.RiskAssessmentDefault], error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	data, total, err := h.Repo.GetAll(ctx, q.Search, q.SearchFilters, q.Page, q.Limit)
	if err != nil {
		return commondomain.Paged[domainRisk.RiskAssessmentDefault]{}, err
	}

	return commondomain.Paged[domainRisk.RiskAssessmentDefault]{
		Data:  data,
		Total: total,
	}, nil
}
