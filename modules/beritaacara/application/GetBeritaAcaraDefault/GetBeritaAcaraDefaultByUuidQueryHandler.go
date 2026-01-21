package application

import (
	"context"

	domainberitaacara "UnpakSiamida/modules/beritaacara/domain"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GetBeritaAcaraDefaultByUuidQueryHandler struct {
	Repo domainberitaacara.IBeritaAcaraRepository
}

func (h *GetBeritaAcaraDefaultByUuidQueryHandler) Handle(
	ctx context.Context,
	q GetBeritaAcaraDefaultByUuidQuery,
) (*domainberitaacara.BeritaAcaraDefault, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	parsed, err := uuid.Parse(q.Uuid)
	if err != nil {
		return nil, domainberitaacara.NotFound(q.Uuid)
	}

	BeritaAcara, err := h.Repo.GetDefaultByUuid(ctx, parsed)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainberitaacara.NotFound(q.Uuid)
		}
		return nil, err
	}

	return BeritaAcara, nil
}
