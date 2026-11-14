package UpdateJadwalProkerV3

import (
	domain "UnpakSiamida/modules/prokerv3/domain"
	"context"
	"time"
)

type UpdateJadwalProkerV3CommandHandler struct {
	Repo domain.IProkerV3Repository
}

func (h *UpdateJadwalProkerV3CommandHandler) Handle(ctx context.Context, cmd UpdateJadwalProkerV3Command) (*UpdateJadwalProkerV3Result, error) {
	existingDTO, err := h.Repo.GetJadwalProkerV3(ctx, cmd.UUID)
	if err != nil || existingDTO == nil {
		return nil, domain.JadwalNotFound()
	}

	tMulai, _ := time.Parse("2006-01-02", existingDTO.TanggalMulai)
	tAkhir, _ := time.Parse("2006-01-02", existingDTO.TanggalAkhir)

	tahun := existingDTO.Tahun
	if cmd.Tahun != "" {
		tahun = cmd.Tahun
	}
	if cmd.TanggalMulai != nil {
		tMulai = *cmd.TanggalMulai
	}
	if cmd.TanggalAkhir != nil {
		tAkhir = *cmd.TanggalAkhir
	}

	fuID := existingDTO.IdFakultasUnit
	if cmd.HasFuKey {
		fuID = cmd.IdFakultasUnit
	}

	if tAkhir.Before(tMulai) {
		return nil, domain.TanggalAkhirInvalid()
	}

	item := domain.JadwalProkerV3{
		ID:             existingDTO.ID,
		UUID:           cmd.UUID,
		Tahun:          tahun,
		TanggalMulai:   tMulai,
		TanggalAkhir:   tAkhir,
		IdFakultasUnit: fuID,
		UpdatedAt:      time.Now(),
	}

	if err := h.Repo.UpdateJadwalProkerV3(ctx, cmd.UUID, &item); err != nil {
		return nil, err
	}

	return &UpdateJadwalProkerV3Result{
		Message: "Jadwal program kerja berhasil diperbarui",
		UUID:    cmd.UUID.String(),
		Data:    item,
	}, nil
}
