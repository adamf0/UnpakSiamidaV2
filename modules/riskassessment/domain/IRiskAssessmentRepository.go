package domain

import (
	commondomain "UnpakSiamida/common/domain"
	"context"

	"github.com/google/uuid"
)

type IRiskAssessmentRepository interface {
	GetByUuid(ctx context.Context, uid uuid.UUID) (*RiskAssessment, error)
	GetDefaultByUuid(ctx context.Context, uid uuid.UUID) (*RiskAssessmentDefault, error)
	GetAll(
		ctx context.Context,
		search string,
		searchFilters []commondomain.SearchFilter,
		page, limit *int,
	) ([]RiskAssessmentDefault, int64, error)
	Create(ctx context.Context, item *RiskAssessment) error
	Update(ctx context.Context, item *RiskAssessment) error
	Delete(ctx context.Context, uid uuid.UUID) error
}
