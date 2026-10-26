package CreateAktivitasV3

import (
	domain "UnpakSiamida/modules/prokerv3/domain"
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

type CreateAktivitasV3CommandHandler struct {
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

func (h *CreateAktivitasV3CommandHandler) Handle(ctx context.Context, cmd CreateAktivitasV3Command) (*domain.AktivitasV2, error) {
	matrixUUID, err := uuid.Parse(cmd.MatrixUUID)
	if err != nil {
		return nil, errors.New("Invalid matrix UUID")
	}

	proker, err := h.Repo.GetDefaultByUuid(ctx, matrixUUID)
	if err != nil {
		return nil, err
	}
	if proker == nil {
		return nil, errors.New("Program kerja tidak ditemukan")
	}

	namaAktivitas := strings.TrimSpace(cmd.Aktivitas)
	if namaAktivitas == "" {
		return nil, errors.New("Nama aktivitas / kegiatan wajib diisi")
	}

	statusStr := strings.ToLower(strings.TrimSpace(cmd.Status))
	if statusStr == "" {
		statusStr = "draft"
	}

	var fuIDUint64 *uint64
	if proker.IdFakultasUnit != nil {
		val := uint64(*proker.IdFakultasUnit)
		fuIDUint64 = &val
	}
	matrikIDUint64 := uint64(proker.ID)

	act := &domain.AktivitasV2{
		UUID:             uuid.New(),
		IdMatrikProker:   &matrikIDUint64,
		IdFakultasUnit:   fuIDUint64,
		Aktivitas:        namaAktivitas,
		PIC: func() string {
			if p := strings.TrimSpace(cmd.PIC); p != "" {
				return p
			}
			return strings.TrimSpace(proker.PIC)
		}(),
		Capaian:          strings.TrimSpace(cmd.Capaian),
		Masalah:          strings.TrimSpace(cmd.Masalah),
		PemecahanMasalah: strings.TrimSpace(cmd.PemecahanMasalah),
		Anggaran:         cmd.Anggaran,
		Status:           statusStr,
		StatusVerifikator: "draft",
		StatusAnggaran:    "menunggu_verifikator",
		StatusUpml:       "menunggu",
		StatusUpmf:       "menunggu",
		StatusUpms:       "menunggu",
	}

	timelines := make([]domain.ProkerTimelineV3, 0, len(cmd.Timelines))
	for i, tInput := range cmd.Timelines {
		tlUUID := uuid.New()
		if strings.TrimSpace(tInput.UUID) != "" {
			if parsed, err := uuid.Parse(tInput.UUID); err == nil {
				tlUUID = parsed
			}
		}

		noAktivitas := tInput.NoAktivitas
		if noAktivitas <= 0 {
			noAktivitas = i + 1
		}

		rkAwal := parseDateFlexible(tInput.TanggalRkAwal)
		if rkAwal == nil {
			def, _ := time.Parse("2006-01-02", "2026-01-01")
			rkAwal = &def
		}

		rkAkhir := parseDateFlexible(tInput.TanggalRkAkhir)
		if rkAkhir == nil {
			def, _ := time.Parse("2006-01-02", "2026-12-31")
			rkAkhir = &def
		}

		rAwal := parseDateFlexible(tInput.TanggalRAwal)
		rAkhir := parseDateFlexible(tInput.TanggalRAkhir)

		titikKontrol := strings.TrimSpace(tInput.TitikKontrol)
		if titikKontrol == "" {
			titikKontrol = "RK"
		}

		tlStatus := strings.TrimSpace(tInput.Status)
		if tlStatus == "" {
			tlStatus = "Belum Terlaksana"
		}

		tNama := strings.TrimSpace(tInput.NamaAktivitas)
		if tNama == "" {
			tNama = namaAktivitas
		}

		timelines = append(timelines, domain.ProkerTimelineV3{
			UUID:             tlUUID,
			NoAktivitas:      noAktivitas,
			NamaAktivitas:    tNama,
			PIC:              strings.TrimSpace(tInput.PIC),
			TitikKontrol:     titikKontrol,
			TanggalRkAwal:    *rkAwal,
			TanggalRkAkhir:   *rkAkhir,
			TanggalRAwal:     rAwal,
			TanggalRAkhir:    rAkhir,
			JadwalRk:         strings.TrimSpace(tInput.JadwalRk),
			JadwalR:          strings.TrimSpace(tInput.JadwalR),
			Capaian:          strings.TrimSpace(tInput.Capaian),
			Masalah:          strings.TrimSpace(tInput.Masalah),
			PemecahanMasalah: strings.TrimSpace(tInput.PemecahanMasalah),
			Status:           tlStatus,
		})
	}

	if err := h.Repo.CreateAktivitas(ctx, act, timelines); err != nil {
		return nil, err
	}

	return act, nil
}
