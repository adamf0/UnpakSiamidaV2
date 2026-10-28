package UpdateTimelineV3

import (
	domain "UnpakSiamida/modules/prokerv3/domain"
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
)

type UpdateTimelineV3CommandHandler struct {
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

func (h *UpdateTimelineV3CommandHandler) Handle(ctx context.Context, cmd UpdateTimelineV3Command) (*domain.ProkerTimelineV3, error) {
	tlUUID, err := uuid.Parse(cmd.UUID)
	if err != nil {
		return nil, domain.InvalidTimelineUuid()
	}

	tl, err := h.Repo.GetTimelineByUuid(ctx, tlUUID)
	if err != nil {
		return nil, err
	}
	if tl == nil {
		return nil, domain.TimelineNotFound()
	}

	var act *domain.AktivitasV2
	if cmd.Anggaran != nil {
		act, _ = h.Repo.GetAktivitasEntityById(ctx, tl.IdAktivitasV2)
	}

	result := domain.UpdateProkerTimelineV3(
		tl,
		act,
		cmd.Anggaran,
		cmd.NamaAktivitas,
		cmd.PIC,
		cmd.TitikKontrol,
		cmd.NoAktivitas,
		parseDateFlexible(cmd.TanggalRkAwal),
		parseDateFlexible(cmd.TanggalRkAkhir),
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

	updatedTL := result.Value
	if err := h.Repo.UpdateTimeline(ctx, updatedTL); err != nil {
		return nil, err
	}

	return updatedTL, nil
}
