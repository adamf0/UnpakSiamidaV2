package application

import (
    "context"

    domainrenstranilai "UnpakSiamida/modules/renstranilai/domain"
    "github.com/google/uuid"
	"errors"
    "gorm.io/gorm"
	"time"
)

type GetRenstraNilaiByUuidQueryHandler struct {
    Repo domainrenstranilai.IRenstraNilaiRepository
}

func (h *GetRenstraNilaiByUuidQueryHandler) Handle(
    ctx context.Context,
    q GetRenstraNilaiByUuidQuery,
) (*domainrenstranilai.RenstraNilai, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

    parsed, err := uuid.Parse(q.Uuid)
	if err != nil {
		return nil, domainrenstranilai.NotFound(q.Uuid)
	}

    inrenstranilai, err := h.Repo.GetByUuid(ctx, parsed)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainrenstranilai.NotFound(q.Uuid)
		}
		return nil, err
	}

    return inrenstranilai, nil
}
