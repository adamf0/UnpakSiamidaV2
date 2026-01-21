package application

import (
	"context"

	domainberitaacara "UnpakSiamida/modules/beritaacara/domain"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GetBeritaAcaraByUuidQueryHandler struct {
	Repo domainberitaacara.IBeritaAcaraRepository
}

func (h *GetBeritaAcaraByUuidQueryHandler) Handle(
	ctx context.Context,
	q GetBeritaAcaraByUuidQuery,
) (*domainberitaacara.BeritaAcara, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	parsed, err := uuid.Parse(q.Uuid)
	if err != nil {
		return nil, domainberitaacara.NotFound(q.Uuid)
	}

	inBeritaAcara, err := h.Repo.GetByUuid(ctx, parsed)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainberitaacara.NotFound(q.Uuid)
		}
		return nil, err
	}

	return inBeritaAcara, nil
}
