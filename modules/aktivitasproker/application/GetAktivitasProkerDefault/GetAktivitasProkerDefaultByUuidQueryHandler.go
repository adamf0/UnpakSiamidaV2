package application

import (
	"context"

	domainAktivitasProker "UnpakSiamida/modules/aktivitasproker/domain"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GetAktivitasProkerDefaultByUuidQueryHandler struct {
	Repo domainAktivitasProker.IAktivitasProkerRepository
}

func (h *GetAktivitasProkerDefaultByUuidQueryHandler) Handle(
	ctx context.Context,
	q GetAktivitasProkerDefaultByUuidQuery,
) (*domainAktivitasProker.AktivitasProkerDefault, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	parsed, err := uuid.Parse(q.Uuid)
	if err != nil {
		return nil, domainAktivitasProker.NotFound(q.Uuid)
	}

	AktivitasProker, err := h.Repo.GetDefaultByUuid(ctx, parsed)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainAktivitasProker.NotFound(q.Uuid)
		}
		return nil, err
	}

	return AktivitasProker, nil
}
