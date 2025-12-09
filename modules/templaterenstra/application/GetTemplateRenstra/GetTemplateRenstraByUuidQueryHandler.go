package application

import (
    "context"

    domaintemplaterenstra "UnpakSiamida/modules/templaterenstra/domain"
    "github.com/google/uuid"
	"errors"
    "gorm.io/gorm"
)

type GetTemplateRenstraByUuidQueryHandler struct {
    Repo domaintemplaterenstra.ITemplateRenstraRepository
}

func (h *GetTemplateRenstraByUuidQueryHandler) Handle(
    ctx context.Context,
    q GetTemplateRenstraByUuidQuery,
) (*domaintemplaterenstra.TemplateRenstra, error) {

    parsed, err := uuid.Parse(q.Uuid)
	if err != nil {
		return nil, domaintemplaterenstra.NotFound(q.Uuid)
	}

    templaterenstra, err := h.Repo.GetByUuid(ctx, parsed)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domaintemplaterenstra.NotFound(q.Uuid)
		}
		return nil, err
	}

    return templaterenstra, nil
}
