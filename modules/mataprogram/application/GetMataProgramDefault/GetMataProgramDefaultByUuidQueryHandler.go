package application

import (
	"context"

	domainmataprogram "UnpakSiamida/modules/mataprogram/domain"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GetMataProgramDefaultByUuidQueryHandler struct {
	Repo domainmataprogram.IMataProgramRepository
}

func (h *GetMataProgramDefaultByUuidQueryHandler) Handle(
	ctx context.Context,
	q GetMataProgramDefaultByUuidQuery,
) (*domainmataprogram.MataProgramDefault, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	parsed, err := uuid.Parse(q.Uuid)
	if err != nil {
		return nil, domainmataprogram.NotFound(q.Uuid)
	}

	MataProgram, err := h.Repo.GetDefaultByUuid(ctx, parsed)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainmataprogram.NotFound(q.Uuid)
		}
		return nil, err
	}

	return MataProgram, nil
}
