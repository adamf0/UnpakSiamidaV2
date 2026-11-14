package ReorderTimelineV3

import (
	domain "UnpakSiamida/modules/aktivitasv2/domain"
	"context"
)

type ReorderTimelineV3CommandHandler struct {
	Repo domain.IAktivitasV2Repository
}

func (h *ReorderTimelineV3CommandHandler) Handle(ctx context.Context, cmd ReorderTimelineV3Command) (string, error) {
	if err := h.Repo.ReorderTimeline(ctx, cmd.Items); err != nil {
		return "", err
	}
	return "Urutan timeline berhasil diperbarui", nil
}
