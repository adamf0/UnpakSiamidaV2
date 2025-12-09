package application

import (
    "context"

    domainJenisFile "UnpakSiamida/modules/jenisfile/domain"
    "github.com/google/uuid"
	"errors"
    "gorm.io/gorm"
)

type GetJenisFileByUuidQueryHandler struct {
    Repo domainJenisFile.IJenisFileRepository
}

func (h *GetJenisFileByUuidQueryHandler) Handle(
    ctx context.Context,
    q GetJenisFileByUuidQuery,
) (*domainJenisFile.JenisFile, error) {

    parsed, err := uuid.Parse(q.Uuid)
	if err != nil {
		return nil, domainJenisFile.NotFound(q.Uuid)
	}

    inJenisFile, err := h.Repo.GetByUuid(ctx, parsed)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainJenisFile.NotFound(q.Uuid)
		}
		return nil, err
	}

    return inJenisFile, nil
}
