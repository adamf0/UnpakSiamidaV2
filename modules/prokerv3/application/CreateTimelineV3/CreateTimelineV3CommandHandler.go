package CreateTimelineV3

import (
	domain "UnpakSiamida/modules/prokerv3/domain"
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
)

type CreateTimelineV3CommandHandler struct {
	Repo domain.IProkerV3Repository
}

func parseDateFlexible(val string) *time.Time {
	val = strings.TrimSpace(val)
	if val == "" {
		return nil
	}
	if t, err := time.Parse("2006-01-02", val); err == nil {
		return &t
	}
	if t, err := time.Parse(time.RFC3339, val); err == nil {
		return &t
	}
	return nil
}

func (h *CreateTimelineV3CommandHandler) Handle(ctx context.Context, cmd CreateTimelineV3Command) (*domain.ProkerTimelineV3, error) {
	actUUID, err := uuid.Parse(cmd.AktivitasUUID)
	if err != nil {
		return nil, domain.InvalidAktivitasUuid()
	}

	act, err := h.Repo.GetAktivitasEntityByUuid(ctx, actUUID)
	if err != nil {
		return nil, err
	}
	if act == nil {
		return nil, domain.AktivitasNotFound()
	}

	rkAwal := parseDateFlexible(cmd.TanggalRkAwal)
	if rkAwal == nil {
		def, _ := time.Parse("2006-01-02", "2026-01-01")
		rkAwal = &def
	}

	rkAkhir := parseDateFlexible(cmd.TanggalRkAkhir)
	if rkAkhir == nil {
		def, _ := time.Parse("2006-01-02", "2026-12-31")
		rkAkhir = &def
	}

	noAktivitas := cmd.NoAktivitas
	if noAktivitas <= 0 && h.Repo != nil {
		count, _ := h.Repo.CountTimelineByAktivitasId(ctx, act.ID)
		noAktivitas = int(count + 1)
	}

	tlUUID := uuid.New()
	if strings.TrimSpace(cmd.UUID) != "" {
		if parsed, err := uuid.Parse(cmd.UUID); err == nil {
			tlUUID = parsed
		}
	}

	result := domain.NewProkerTimelineV3(
		act,
		tlUUID,
		noAktivitas,
		cmd.NamaAktivitas,
		cmd.PIC,
		cmd.TitikKontrol,
		cmd.Anggaran,
		*rkAwal,
		*rkAkhir,
		parseDateFlexible(cmd.TanggalRAwal),
		parseDateFlexible(cmd.TanggalRAkhir),
		cmd.JadwalRk,
		cmd.JadwalR,
		cmd.Capaian,
		cmd.Masalah,
		cmd.PemecahanMasalah,
		cmd.Status,
		cmd.RAB,
	)

	if !result.IsSuccess {
		return nil, result.Error
	}

	newTL := result.Value
	if err := h.Repo.CreateTimeline(ctx, newTL); err != nil {
		return nil, err
	}

	return newTL, nil
}
