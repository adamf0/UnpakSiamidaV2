package DeleteProkerV3

import (
	domain "UnpakSiamida/modules/prokerv3/domain"
	"context"
	"strings"

	"github.com/google/uuid"
)

type DeleteProkerV3CommandHandler struct {
	Repo domain.IProkerV3Repository
}

func (h *DeleteProkerV3CommandHandler) Handle(
	ctx context.Context,
	cmd DeleteProkerV3Command,
) (string, error) {
	if err := cmd.Validate(); err != nil {
		return "", err
	}

	prokerUUID, err := uuid.Parse(strings.TrimSpace(cmd.UUID))
	if err != nil {
		return "", domain.ProkerNotFound()
	}

	if cmd.SoftDelete {
		if err := h.Repo.SoftDelete(ctx, domain.SoftDeleteProkerParams{
			UUID:          prokerUUID,
			Catatan:       strings.TrimSpace(cmd.Catatan),
			Role:          strings.TrimSpace(cmd.Role),
			ActorID:       cmd.ActorID,
			ActorName:     strings.TrimSpace(cmd.ActorName),
			ActorUsername: strings.TrimSpace(cmd.ActorUsername),
		}); err != nil {
			return "", err
		}
		return cmd.UUID, nil
	}

	if err := h.Repo.Delete(ctx, prokerUUID); err != nil {
		return "", err
	}

	return cmd.UUID, nil
}
