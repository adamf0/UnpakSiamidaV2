package application

import (
    "context"

    domainJenisFile "UnpakSiamida/modules/jenisfile/domain"
    "github.com/google/uuid"
	"errors"
    "gorm.io/gorm"
)

type GetJenisFileDefaultByUuidQueryHandler struct {
    Repo domainJenisFile.IJenisFileRepository
}

func (h *GetJenisFileDefaultByUuidQueryHandler) Handle(
    ctx context.Context,
    q GetJenisFileDefaultByUuidQuery,
) (*domainJenisFile.JenisFileDefault, error) {

    parsed, err := uuid.Parse(q.Uuid)
	if err != nil {
		return nil, domainJenisFile.NotFound(q.Uuid)
	}

    JenisFile, err := h.Repo.GetDefaultByUuid(ctx, parsed)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainJenisFile.NotFound(q.Uuid)
		}
		return nil, err
	}

    return JenisFile, nil
}
