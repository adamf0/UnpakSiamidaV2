package application

import (
	"context"

	domainmataprogram "UnpakSiamida/modules/mataprogram/domain"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GetMataProgramByUuidQueryHandler struct {
	Repo domainmataprogram.IMataProgramRepository
}

func (h *GetMataProgramByUuidQueryHandler) Handle(
	ctx context.Context,
	q GetMataProgramByUuidQuery,
) (*domainmataprogram.MataProgram, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	parsed, err := uuid.Parse(q.Uuid)
	if err != nil {
		return nil, domainmataprogram.NotFound(q.Uuid)
	}

	inMataProgram, err := h.Repo.GetByUuid(ctx, parsed)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainmataprogram.NotFound(q.Uuid)
		}
		return nil, err
	}

	return inMataProgram, nil
}
