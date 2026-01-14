package application

import (
    "context"

    domainfakultasunit "UnpakSiamida/modules/fakultasunit/domain"
    "github.com/google/uuid"
	"errors"
    "gorm.io/gorm"
	"time"
)

type GetFakultasUnitByUuidQueryHandler struct {
    Repo domainfakultasunit.IFakultasUnitRepository
}

func (h *GetFakultasUnitByUuidQueryHandler) Handle(
    ctx context.Context,
    q GetFakultasUnitByUuidQuery,
) (*domainfakultasunit.FakultasUnit, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

    parsed, err := uuid.Parse(q.Uuid)
	if err != nil {
		return nil, domainfakultasunit.NotFound(q.Uuid)
	}

    fakultasunit, err := h.Repo.GetDefaultByUuid(ctx, parsed)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainfakultasunit.NotFound(q.Uuid)
		}
		return nil, err
	}

    return fakultasunit, nil
}
