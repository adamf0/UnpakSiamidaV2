package application

import (
	"context"

	domainDokumenProker "UnpakSiamida/modules/dokumenproker/domain"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GetDokumenProkerByUuidQueryHandler struct {
	Repo domainDokumenProker.IDokumenProkerRepository
}

func (h *GetDokumenProkerByUuidQueryHandler) Handle(
	ctx context.Context,
	q GetDokumenProkerByUuidQuery,
) (*domainDokumenProker.DokumenProker, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	parsed, err := uuid.Parse(q.Uuid)
	if err != nil {
		return nil, domainDokumenProker.NotFound(q.Uuid)
	}

	inDokumenProker, err := h.Repo.GetByUuid(ctx, parsed)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainDokumenProker.NotFound(q.Uuid)
		}
		return nil, err
	}

	return inDokumenProker, nil
}
