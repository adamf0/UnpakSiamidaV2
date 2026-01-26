package application

import (
	"context"

	domaintahunproker "UnpakSiamida/modules/tahunproker/domain"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GetTahunProkerByUuidQueryHandler struct {
	Repo domaintahunproker.ITahunProkerRepository
}

func (h *GetTahunProkerByUuidQueryHandler) Handle(
	ctx context.Context,
	q GetTahunProkerByUuidQuery,
) (*domaintahunproker.TahunProker, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	parsed, err := uuid.Parse(q.Uuid)
	if err != nil {
		return nil, domaintahunproker.NotFound(q.Uuid)
	}

	inTahunProker, err := h.Repo.GetByUuid(ctx, parsed)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domaintahunproker.NotFound(q.Uuid)
		}
		return nil, err
	}

	return inTahunProker, nil
}
