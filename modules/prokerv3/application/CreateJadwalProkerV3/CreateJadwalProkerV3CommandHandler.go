package CreateJadwalProkerV3

import (
	domain "UnpakSiamida/modules/prokerv3/domain"
	"context"
	"time"

	"github.com/google/uuid"
)

type CreateJadwalProkerV3CommandHandler struct {
	Repo domain.IProkerV3Repository
}

func (h *CreateJadwalProkerV3CommandHandler) Handle(ctx context.Context, cmd CreateJadwalProkerV3Command) (*CreateJadwalProkerV3Result, error) {
	if cmd.Tahun == "" {
		return nil, domain.TahunRequired()
	}
	if cmd.TanggalMulai.IsZero() {
		return nil, domain.TanggalMulaiRequired()
	}
	if cmd.TanggalAkhir.IsZero() {
		return nil, domain.TanggalAkhirRequired()
	}
	if cmd.TanggalAkhir.Before(cmd.TanggalMulai) {
		return nil, domain.TanggalAkhirInvalid()
	}

	newJadwal := domain.JadwalProkerV3{
		UUID:           uuid.New(),
		Tahun:          cmd.Tahun,
		TanggalMulai:   cmd.TanggalMulai,
		TanggalAkhir:   cmd.TanggalAkhir,
		IdFakultasUnit: cmd.IdFakultasUnit,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := h.Repo.CreateJadwalProkerV3(ctx, &newJadwal); err != nil {
		return nil, err
	}

	return &CreateJadwalProkerV3Result{
		Message: "Jadwal program kerja berhasil ditambahkan",
		UUID:    newJadwal.UUID.String(),
		Data:    newJadwal,
	}, nil
}
