package application

import (
	"context"

	domainJadwalProker "UnpakSiamida/modules/jadwalproker/domain"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GetJadwalProkerDefaultByUuidQueryHandler struct {
	Repo domainJadwalProker.IJadwalProkerRepository
}

func (h *GetJadwalProkerDefaultByUuidQueryHandler) Handle(
	ctx context.Context,
	q GetJadwalProkerDefaultByUuidQuery,
) (*domainJadwalProker.JadwalProkerDefault, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	parsed, err := uuid.Parse(q.Uuid)
	if err != nil {
		return nil, domainJadwalProker.NotFound(q.Uuid)
	}

	JadwalProker, err := h.Repo.GetDefaultByUuid(ctx, parsed)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainJadwalProker.NotFound(q.Uuid)
		}
		return nil, err
	}

	return JadwalProker, nil
}
