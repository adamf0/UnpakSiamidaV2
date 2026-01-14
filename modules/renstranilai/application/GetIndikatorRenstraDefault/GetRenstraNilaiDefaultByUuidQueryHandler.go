package application

import (
    "context"

    domainrenstranilai "UnpakSiamida/modules/renstranilai/domain"
    "github.com/google/uuid"
	"errors"
    "gorm.io/gorm"
	"time"
)

type GetRenstraNilaiDefaultByUuidQueryHandler struct {
    Repo domainrenstranilai.IRenstraNilaiRepository
}

func (h *GetRenstraNilaiDefaultByUuidQueryHandler) Handle(
    ctx context.Context,
    q GetRenstraNilaiDefaultByUuidQuery,
) (*domainrenstranilai.RenstraNilaiDefault, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

    parsed, err := uuid.Parse(q.Uuid)
	if err != nil {
		return nil, domainrenstranilai.NotFound(q.Uuid)
	}

    inrenstranilai, err := h.Repo.GetDefaultByUuid(ctx, parsed)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainrenstranilai.NotFound(q.Uuid)
		}
		return nil, err
	}

    return inrenstranilai, nil
}
