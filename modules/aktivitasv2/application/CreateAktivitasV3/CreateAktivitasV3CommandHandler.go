package CreateAktivitasV3

import (
	domain "UnpakSiamida/modules/aktivitasv2/domain"
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/google/uuid"
)

type CreateAktivitasV3CommandHandler struct {
	Repo domain.IAktivitasV2Repository
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
		return nil, domain.InvalidMatrixUuid()
	}

	proker, err := h.Repo.GetMatrixProkerByUuid(ctx, matrixUUID)
	if err != nil {
		return nil, err
	}
	if proker == nil {
		return nil, domain.ProkerNotFound()
	}
	if proker.DeletedAt != nil {
		return nil, domain.ProkerDeleted()
	}

	namaAktivitas := strings.TrimSpace(cmd.Aktivitas)
	if namaAktivitas == "" {
		return nil, domain.NamaAktivitasRequired()
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
		RAB: func() string {
			if strings.TrimSpace(cmd.RAB) != "" && cmd.RAB != "[]" {
				return strings.TrimSpace(cmd.RAB)
			}
			if cmd.Anggaran > 0 {
				b, _ := json.Marshal([]map[string]interface{}{
					{
						"kategori":   "Umum",
						"rincian":    "Alokasi Anggaran Kegiatan",
						"keterangan": "Alokasi Anggaran Kegiatan",
						"qty_satuan": 1,
						"qty_paket":  1,
						"sbu":        cmd.Anggaran,
						"total":      cmd.Anggaran,
					},
				})
				return string(b)
			}
			return ""
		}(),
		Status:           statusStr,
		StatusAnggaran:    "menunggu_verifikator",
		StatusUpml:       "menunggu",
		StatusUpmf:       "menunggu",
		StatusUpms:       "menunggu",
		StatusLpm:        "menunggu",
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
