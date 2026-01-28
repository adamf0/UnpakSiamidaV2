package application

import (
	"context"

	domainAktivitasProker "UnpakSiamida/modules/aktivitasproker/domain"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GetAktivitasProkerByUuidQueryHandler struct {
	Repo domainAktivitasProker.IAktivitasProkerRepository
}

func (h *GetAktivitasProkerByUuidQueryHandler) Handle(
	ctx context.Context,
	q GetAktivitasProkerByUuidQuery,
) (*domainAktivitasProker.AktivitasProker, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	parsed, err := uuid.Parse(q.Uuid)
	if err != nil {
		return nil, domainAktivitasProker.NotFound(q.Uuid)
	}

	inAktivitasProker, err := h.Repo.GetByUuid(ctx, parsed)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainAktivitasProker.NotFound(q.Uuid)
		}
		return nil, err
	}

	return inAktivitasProker, nil
}
