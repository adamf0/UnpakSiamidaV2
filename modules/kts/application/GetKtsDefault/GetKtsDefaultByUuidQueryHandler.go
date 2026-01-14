package application

import (
    "context"

    domainKts "UnpakSiamida/modules/kts/domain"
    "github.com/google/uuid"
	"errors"
    "gorm.io/gorm"
	"time"
)

type GetKtsDefaultByUuidQueryHandler struct {
    Repo domainKts.IKtsRepository
}

func (h *GetKtsDefaultByUuidQueryHandler) Handle(
    ctx context.Context,
    q GetKtsDefaultByUuidQuery,
) (*domainKts.KtsDefault, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

    parsed, err := uuid.Parse(q.Uuid)
	if err != nil {
		return nil, domainKts.NotFound(q.Uuid)
	}

    Kts, err := h.Repo.GetDefaultByUuid(ctx, parsed)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainKts.NotFound(q.Uuid)
		}
		return nil, err
	}

    return Kts, nil
}
