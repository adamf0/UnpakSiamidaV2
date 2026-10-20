package UpdateAktivitasV3

import (
	domain "UnpakSiamida/modules/prokerv3/domain"
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
)

type UpdateAktivitasV3CommandHandler struct {
	Repo domain.IProkerV3Repository
}

func (h *UpdateAktivitasV3CommandHandler) Handle(ctx context.Context, cmd UpdateAktivitasV3Command) (*domain.AktivitasV2, error) {
	actUUID, err := uuid.Parse(cmd.UUID)
	if err != nil {
		return nil, errors.New("Invalid aktivitas UUID")
	}

	act, err := h.Repo.GetAktivitasEntityByUuid(ctx, actUUID)
	if err != nil {
		return nil, err
	}
	if act == nil {
		return nil, errors.New("Aktivitas tidak ditemukan")
	}

	if cmd.Aktivitas != "" {
		act.Aktivitas = strings.TrimSpace(cmd.Aktivitas)
	}
	if cmd.PIC != "" {
		act.PIC = strings.TrimSpace(cmd.PIC)
	}
	if cmd.Anggaran > 0 {
		act.Anggaran = cmd.Anggaran
	}
	if cmd.Realisasi != nil {
		act.Realisasi = cmd.Realisasi
	}
	if cmd.Status != "" {
		act.Status = strings.TrimSpace(cmd.Status)
	}
	if cmd.Capaian != "" {
		act.Capaian = strings.TrimSpace(cmd.Capaian)
	}
	if cmd.Masalah != "" {
		act.Masalah = strings.TrimSpace(cmd.Masalah)
	}
	if cmd.PemecahanMasalah != "" {
		act.PemecahanMasalah = strings.TrimSpace(cmd.PemecahanMasalah)
	}

	if err := h.Repo.UpdateAktivitas(ctx, act); err != nil {
		return nil, err
	}

	return act, nil
}
