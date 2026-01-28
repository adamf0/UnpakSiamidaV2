package application

import (
	"context"

	domainJadwalProker "UnpakSiamida/modules/jadwalproker/domain"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GetJadwalProkerByUuidQueryHandler struct {
	Repo domainJadwalProker.IJadwalProkerRepository
}

func (h *GetJadwalProkerByUuidQueryHandler) Handle(
	ctx context.Context,
	q GetJadwalProkerByUuidQuery,
) (*domainJadwalProker.JadwalProker, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	parsed, err := uuid.Parse(q.Uuid)
	if err != nil {
		return nil, domainJadwalProker.NotFound(q.Uuid)
	}

	inJadwalProker, err := h.Repo.GetByUuid(ctx, parsed)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainJadwalProker.NotFound(q.Uuid)
		}
		return nil, err
	}

	return inJadwalProker, nil
}
