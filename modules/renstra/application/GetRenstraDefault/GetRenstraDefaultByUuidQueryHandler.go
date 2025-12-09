package application

import (
    "context"

    domainrenstra "UnpakSiamida/modules/renstra/domain"
    "github.com/google/uuid"
	"errors"
    "gorm.io/gorm"
)

type GetRenstraDefaultByUuidQueryHandler struct {
    Repo domainrenstra.IRenstraRepository
}

func (h *GetRenstraDefaultByUuidQueryHandler) Handle(
    ctx context.Context,
    q GetRenstraDefaultByUuidQuery,
) (*domainrenstra.RenstraDefault, error) {

    parsed, err := uuid.Parse(q.Uuid)
	if err != nil {
		return nil, domainrenstra.NotFound(q.Uuid)
	}

    inrenstra, err := h.Repo.GetDefaultByUuid(ctx, parsed)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainrenstra.NotFound(q.Uuid)
		}
		return nil, err
	}

    return inrenstra, nil
}
