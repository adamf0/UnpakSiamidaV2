package application

import (
	"context"

	domainDokumenProker "UnpakSiamida/modules/dokumenproker/domain"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GetDokumenProkerDefaultByUuidQueryHandler struct {
	Repo domainDokumenProker.IDokumenProkerRepository
}

func (h *GetDokumenProkerDefaultByUuidQueryHandler) Handle(
	ctx context.Context,
	q GetDokumenProkerDefaultByUuidQuery,
) (*domainDokumenProker.DokumenProkerDefault, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	parsed, err := uuid.Parse(q.Uuid)
	if err != nil {
		return nil, domainDokumenProker.NotFound(q.Uuid)
	}

	DokumenProker, err := h.Repo.GetDefaultByUuid(ctx, parsed)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainDokumenProker.NotFound(q.Uuid)
		}
		return nil, err
	}

	return DokumenProker, nil
}
