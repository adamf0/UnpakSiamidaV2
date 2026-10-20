package CreateTimelineV3

import (
	domain "UnpakSiamida/modules/prokerv3/domain"
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CreateTimelineV3CommandHandler struct {
	Repo domain.IProkerV3Repository
	DB   *gorm.DB
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
		return nil, errors.New("Invalid aktivitas UUID")
	}

	act, err := h.Repo.GetAktivitasEntityByUuid(ctx, actUUID)
	if err != nil {
		return nil, err
	}
	if act == nil {
		return nil, errors.New("Aktivitas tidak ditemukan")
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
	if noAktivitas <= 0 && h.DB != nil {
		var count int64
		_ = h.DB.WithContext(ctx).Table("aktivitasv2_timeline").Where("id_aktivitasv2 = ?", act.ID).Count(&count)
		noAktivitas = int(count + 1)
	}

	titikKontrol := strings.TrimSpace(cmd.TitikKontrol)
	if titikKontrol == "" {
		titikKontrol = "RK"
	}

	status := strings.TrimSpace(cmd.Status)
	if status == "" {
		status = "Belum Terlaksana"
	}

	namaAkt := strings.TrimSpace(cmd.NamaAktivitas)
	if namaAkt == "" {
		namaAkt = act.Aktivitas
	}

	pic := strings.TrimSpace(cmd.PIC)
	if pic == "" {
		pic = act.PIC
	}

	tlUUID := uuid.New()
	if strings.TrimSpace(cmd.UUID) != "" {
		if parsed, err := uuid.Parse(cmd.UUID); err == nil {
			tlUUID = parsed
		}
	}

	newTL := domain.ProkerTimelineV3{
		UUID:             tlUUID,
		IdAktivitasV2:    act.ID,
		NoAktivitas:      noAktivitas,
		NamaAktivitas:    namaAkt,
		PIC:              pic,
		TitikKontrol:     titikKontrol,
		TanggalRkAwal:    *rkAwal,
		TanggalRkAkhir:   *rkAkhir,
		TanggalRAwal:     parseDateFlexible(cmd.TanggalRAwal),
		TanggalRAkhir:    parseDateFlexible(cmd.TanggalRAkhir),
		JadwalRk:         strings.TrimSpace(cmd.JadwalRk),
		JadwalR:          strings.TrimSpace(cmd.JadwalR),
		Capaian:          strings.TrimSpace(cmd.Capaian),
		Masalah:          strings.TrimSpace(cmd.Masalah),
		PemecahanMasalah: strings.TrimSpace(cmd.PemecahanMasalah),
		Status:           status,
	}

	if err := h.Repo.CreateTimeline(ctx, &newTL); err != nil {
		return nil, err
	}

	return &newTL, nil
}
