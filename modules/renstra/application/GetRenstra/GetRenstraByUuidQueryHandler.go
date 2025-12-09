package application

import (
    "context"

    domainrenstra "UnpakSiamida/modules/renstra/domain"
    "github.com/google/uuid"
	"errors"
    "gorm.io/gorm"
)

type GetRenstraByUuidQueryHandler struct {
    Repo domainrenstra.IRenstraRepository
}

func (h *GetRenstraByUuidQueryHandler) Handle(
    ctx context.Context,
    q GetRenstraByUuidQuery,
) (*domainrenstra.Renstra, error) {

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
