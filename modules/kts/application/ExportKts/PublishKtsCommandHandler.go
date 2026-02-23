package application

import (
	"context"

	domainkts "UnpakSiamida/modules/kts/domain"
	"UnpakSiamida/modules/kts/event"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/mehdihadeli/go-mediatr"
	"gorm.io/gorm"
)

type PublishKtsCommandHandler struct {
	Repo domainkts.IKtsRepository
}

func (h *PublishKtsCommandHandler) Handle(
	ctx context.Context,
	cmd PublishKtsCommand,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	parsed, err := uuid.Parse(cmd.Uuid)
	if err != nil {
		return "", domainkts.NotFound(cmd.Uuid)
	}

	kts, err := h.Repo.GetByUuid(ctx, parsed)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", domainkts.NotFound(cmd.Uuid)
		}
		return "", err
	}

	err = mediatr.Publish(ctx, event.KtsPdfRequestedEvent{
		EventID:    uuid.New(),
		KtsUUID:    kts.UUID,
		Token:      cmd.Token,
		OccurredOn: time.Now().UTC(),
	})
	if err != nil {
		return "", err
	}

	kts.ClearDomainEvents()

	return kts.UUID.String(), nil
}
