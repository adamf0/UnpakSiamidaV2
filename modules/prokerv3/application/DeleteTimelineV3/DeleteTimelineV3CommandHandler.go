package DeleteTimelineV3

import (
	domain "UnpakSiamida/modules/prokerv3/domain"
	"context"
	"errors"

	"github.com/google/uuid"
)

type DeleteTimelineV3CommandHandler struct {
	Repo domain.IProkerV3Repository
}

func (h *DeleteTimelineV3CommandHandler) Handle(ctx context.Context, cmd DeleteTimelineV3Command) (string, error) {
	tlUUID, err := uuid.Parse(cmd.UUID)
	if err != nil {
		return "", errors.New("Invalid timeline UUID")
	}

	if err := h.Repo.DeleteTimeline(ctx, tlUUID); err != nil {
		return "", err
	}

	return cmd.UUID, nil
}
