package application

import (
	"context"

	domainberitaacara "UnpakSiamida/modules/beritaacara/domain"
	"UnpakSiamida/modules/beritaacara/event"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/mehdihadeli/go-mediatr"
	"gorm.io/gorm"
)

type PublishBeritaAcaraCommandHandler struct {
	Repo domainberitaacara.IBeritaAcaraRepository
}

func (h *PublishBeritaAcaraCommandHandler) Handle(
	ctx context.Context,
	cmd PublishBeritaAcaraCommand,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	parsed, err := uuid.Parse(cmd.Uuid)
	if err != nil {
		return "", domainberitaacara.NotFound(cmd.Uuid)
	}

	beritaacara, err := h.Repo.GetByUuid(ctx, parsed)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", domainberitaacara.NotFound(cmd.Uuid)
		}
		return "", err
	}

	err = mediatr.Publish(ctx, event.BeritaAcaraPdfRequestedEvent{
		EventID:         uuid.New(),
		BeritaAcaraUUID: beritaacara.UUID,
		Token:           cmd.Token,
		OccurredOn:      time.Now().UTC(),
	})
	if err != nil {
		return "", err
	}

	beritaacara.ClearDomainEvents()

	return beritaacara.UUID.String(), nil
}
