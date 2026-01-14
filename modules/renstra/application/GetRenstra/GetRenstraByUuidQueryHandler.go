package application

import (
    "context"

    domainrenstra "UnpakSiamida/modules/renstra/domain"
    "github.com/google/uuid"
	"errors"
    "gorm.io/gorm"
	"time"
)

type GetRenstraByUuidQueryHandler struct {
    Repo domainrenstra.IRenstraRepository
}

func (h *GetRenstraByUuidQueryHandler) Handle(
    ctx context.Context,
    q GetRenstraByUuidQuery,
) (*domainrenstra.Renstra, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

    parsed, err := uuid.Parse(q.Uuid)
	if err != nil {
		return nil, domainrenstra.NotFound(q.Uuid)
	}

    inrenstra, err := h.Repo.GetByUuid(ctx, parsed)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainrenstra.NotFound(q.Uuid)
		}
		return nil, err
	}

    return inrenstra, nil
}
