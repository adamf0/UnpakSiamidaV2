package ReorderTimelineV3

import (
	domain "UnpakSiamida/modules/prokerv3/domain"
	"context"
)

type ReorderTimelineV3CommandHandler struct {
	Repo domain.IProkerV3Repository
}

func (h *ReorderTimelineV3CommandHandler) Handle(ctx context.Context, cmd ReorderTimelineV3Command) (int, error) {
	if err := h.Repo.ReorderTimeline(ctx, cmd.Items); err != nil {
		return 0, err
	}
	return len(cmd.Items), nil
}
