package application

import (
	"context"

	domaintahunproker "UnpakSiamida/modules/tahunproker/domain"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GetTahunProkerDefaultByUuidQueryHandler struct {
	Repo domaintahunproker.ITahunProkerRepository
}

func (h *GetTahunProkerDefaultByUuidQueryHandler) Handle(
	ctx context.Context,
	q GetTahunProkerDefaultByUuidQuery,
) (*domaintahunproker.TahunProkerDefault, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	parsed, err := uuid.Parse(q.Uuid)
	if err != nil {
		return nil, domaintahunproker.NotFound(q.Uuid)
	}

	TahunProker, err := h.Repo.GetDefaultByUuid(ctx, parsed)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domaintahunproker.NotFound(q.Uuid)
		}
		return nil, err
	}

	return TahunProker, nil
}
