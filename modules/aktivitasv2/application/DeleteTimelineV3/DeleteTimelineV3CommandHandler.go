package DeleteTimelineV3

import (
	domain "UnpakSiamida/modules/aktivitasv2/domain"
	"context"

	"github.com/google/uuid"
)

type DeleteTimelineV3CommandHandler struct {
	Repo domain.IAktivitasV2Repository
}

func (h *DeleteTimelineV3CommandHandler) Handle(ctx context.Context, cmd DeleteTimelineV3Command) (string, error) {
	tlUUID, err := uuid.Parse(cmd.UUID)
	if err != nil {
		return "", domain.InvalidTimelineUuid()
	}

	if err := h.Repo.DeleteTimeline(ctx, tlUUID); err != nil {
		return "", err
	}

	return cmd.UUID, nil
}
