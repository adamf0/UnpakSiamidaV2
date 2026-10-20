package UpdateTimelineV3

import (
	domain "UnpakSiamida/modules/prokerv3/domain"
	"context"
	"errors"
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
		return nil, errors.New("Invalid timeline UUID")
	}

	tl, err := h.Repo.GetTimelineByUuid(ctx, tlUUID)
	if err != nil {
		return nil, err
	}
	if tl == nil {
		return nil, errors.New("Timeline aktivitas tidak ditemukan")
	}

	if cmd.NamaAktivitas != "" {
		tl.NamaAktivitas = strings.TrimSpace(cmd.NamaAktivitas)
	}
	if cmd.PIC != "" {
		tl.PIC = strings.TrimSpace(cmd.PIC)
	}
	if cmd.TitikKontrol != "" {
		tl.TitikKontrol = strings.TrimSpace(cmd.TitikKontrol)
	}
	if cmd.NoAktivitas > 0 {
		tl.NoAktivitas = cmd.NoAktivitas
	}
	if cmd.JadwalRk != "" {
		tl.JadwalRk = strings.TrimSpace(cmd.JadwalRk)
	}
	if cmd.JadwalR != "" {
		tl.JadwalR = strings.TrimSpace(cmd.JadwalR)
	}
	if cmd.Capaian != "" {
		tl.Capaian = strings.TrimSpace(cmd.Capaian)
	}
	if cmd.Masalah != "" {
		tl.Masalah = strings.TrimSpace(cmd.Masalah)
	}
	if cmd.PemecahanMasalah != "" {
		tl.PemecahanMasalah = strings.TrimSpace(cmd.PemecahanMasalah)
	}
	if cmd.Status != "" {
		tl.Status = strings.TrimSpace(cmd.Status)
	}
	if t := parseDateFlexible(cmd.TanggalRkAwal); t != nil {
		tl.TanggalRkAwal = *t
	}
	if t := parseDateFlexible(cmd.TanggalRkAkhir); t != nil {
		tl.TanggalRkAkhir = *t
	}
	if t := parseDateFlexible(cmd.TanggalRAwal); t != nil {
		tl.TanggalRAwal = t
	}
	if t := parseDateFlexible(cmd.TanggalRAkhir); t != nil {
		tl.TanggalRAkhir = t
	}

	if err := h.Repo.UpdateTimeline(ctx, tl); err != nil {
		return nil, err
	}

	return tl, nil
}
