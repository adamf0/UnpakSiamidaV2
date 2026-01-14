package application

import (
    "context"

    domainKts "UnpakSiamida/modules/kts/domain"
    "github.com/google/uuid"
	"errors"
    "gorm.io/gorm"
	"time"
)

type GetKtsByUuidQueryHandler struct {
    Repo domainKts.IKtsRepository
}

func (h *GetKtsByUuidQueryHandler) Handle(
    ctx context.Context,
    q GetKtsByUuidQuery,
) (*domainKts.Kts, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

    parsed, err := uuid.Parse(q.Uuid)
	if err != nil {
		return nil, domainKts.NotFound(q.Uuid)
	}

    inKts, err := h.Repo.GetByUuid(ctx, parsed)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainKts.NotFound(q.Uuid)
		}
		return nil, err
	}

    return inKts, nil
}
