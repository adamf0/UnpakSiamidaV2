package GetRiskOptions

import (
	"context"

	domainRisk "UnpakSiamida/modules/riskassessment/domain"
)

type GetRiskOptionsQueryHandler struct {
	Repo domainRisk.IRiskAssessmentRepository
}

func (h *GetRiskOptionsQueryHandler) Handle(ctx context.Context, q GetRiskOptionsQuery) (*domainRisk.RiskOptionsResult, error) {
	return h.Repo.GetOptions(ctx, q.FakultasUnit, q.Tahun)
}

type GetRiskCriteriaQueryHandler struct {
	Repo domainRisk.IRiskAssessmentRepository
}

func (h *GetRiskCriteriaQueryHandler) Handle(ctx context.Context, q GetRiskCriteriaQuery) ([]domainRisk.RiskCriteriaItem, error) {
	return h.Repo.GetCriteria(ctx, q.Tipe)
}

type GetRiskLevelsQueryHandler struct {
	Repo domainRisk.IRiskAssessmentRepository
}

func (h *GetRiskLevelsQueryHandler) Handle(ctx context.Context, q GetRiskLevelsQuery) ([]domainRisk.RiskLevelItem, error) {
	return h.Repo.GetLevels(ctx)
}

type GetRiskCategoriesQueryHandler struct {
	Repo domainRisk.IRiskAssessmentRepository
}

func (h *GetRiskCategoriesQueryHandler) Handle(ctx context.Context, q GetRiskCategoriesQuery) ([]domainRisk.RiskCategoryItem, error) {
	return h.Repo.GetCategories(ctx)
}

type GetRiskStatusesQueryHandler struct {
	Repo domainRisk.IRiskAssessmentRepository
}

func (h *GetRiskStatusesQueryHandler) Handle(ctx context.Context, q GetRiskStatusesQuery) ([]domainRisk.RiskStatusItem, error) {
	return h.Repo.GetStatuses(ctx)
}

type GetRiskProgramsQueryHandler struct {
	Repo domainRisk.IRiskAssessmentRepository
}

func (h *GetRiskProgramsQueryHandler) Handle(ctx context.Context, q GetRiskProgramsQuery) ([]domainRisk.RiskMataProgramItem, error) {
	return h.Repo.GetPrograms(ctx, q.FakultasUnit, q.Tahun, q.Search)
}
