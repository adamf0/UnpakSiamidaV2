package application

import (
	domainRisk "UnpakSiamida/modules/riskassessment/domain"
	"context"
	"time"

	"github.com/google/uuid"
)

type DeleteRiskAssessmentCommandHandler struct {
	Repo domainRisk.IRiskAssessmentRepository
}

func (h *DeleteRiskAssessmentCommandHandler) Handle(
	ctx context.Context,
	cmd DeleteRiskAssessmentCommand,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	parsedUuid, err := uuid.Parse(cmd.UUID)
	if err != nil {
		return "", domainRisk.NotFound(cmd.UUID)
	}

	if err := h.Repo.Delete(ctx, parsedUuid); err != nil {
		return "", err
	}

	return cmd.UUID, nil
}
