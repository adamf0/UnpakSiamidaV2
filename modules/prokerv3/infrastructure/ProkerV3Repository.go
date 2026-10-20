package infrastructure

import (
	commondomain "UnpakSiamida/common/domain"
	domain "UnpakSiamida/modules/prokerv3/domain"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProkerV3Repository struct {
	db *gorm.DB
}

func NewProkerV3Repository(db *gorm.DB) domain.IProkerV3Repository {
	return &ProkerV3Repository{db: db}
}

// ------------------------
// GET BY UUID (matrix_proker)
// ------------------------
func (r *ProkerV3Repository) GetByUuid(ctx context.Context, uid uuid.UUID) (*domain.ProkerV3, error) {
	var row domain.ProkerV3
	err := r.db.WithContext(ctx).
		Where("uuid = ?", uid).
		First(&row).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &row, nil
}

// ------------------------
// GET DEFAULT BY UUID (matrix_proker joined)
// ------------------------
func (r *ProkerV3Repository) GetDefaultByUuid(ctx context.Context, uid uuid.UUID) (*domain.ProkerV3Default, error) {
	var row domain.ProkerV3Default

	err := r.db.WithContext(ctx).
		Table("matrix_proker mp_t").
		Select(`
			mp_t.id as Id,
			mp_t.uuid as UUID,
			mp_t.id_fakultas_unit as IdFakultasUnit,
			COALESCE(fu.nama_fak_prod_unit, '') as NamaFakultasUnit,
			COALESCE(fu.uuid, '') as UuidFakultasUnit,
			COALESCE(fu.type, '') as TypeFakultasUnit,
			COALESCE(fu.fakultas, '') as ParentFakultas,
			mp_t.id_mata_program as IdMataProgram,
			COALESCE(prog.mata_program, '') as NamaProgram,
			COALESCE(mt.tahun, '2026') as Tahun,
			1 as NoProker,
			mp_t.id_indikator as IdIndikator,
			COALESCE(mir.indikator, '') as IndikatorKinerja,
			COALESCE(mir.uuid, '') as IndikatorUUID,
			COALESCE(mir.id_master_standar, msr.id, 0) as IdMasterStandar,
			COALESCE(msr.nama, '') as StandarSubStandar,
			COALESCE(msr.uuid, '') as StandarUUID,
			COALESCE(mp_t.target_unpak, '') as TargetUnpak,
			COALESCE(mp_t.capaian_terdahulu, '') as CapaianUnitLalu,
			COALESCE(mp_t.target_unit, '') as TargetUnit,
			COALESCE(mp_t.output, '') as Output,
			COALESCE(mp_t.outcome, '') as Outcome,
			COALESCE((SELECT a.pic FROM aktivitasv2 a WHERE a.id_matrik_proker = mp_t.id ORDER BY a.id ASC LIMIT 1), '') as PIC,
			COALESCE(NULLIF(mp_t.status, ''), 'draft') as Status,
			COALESCE(mp_t.catatan_penolakan, '') as CatatanPenolakan,
			COALESCE(NULLIF(mp_t.status_verifikator, ''), 'draft') as StatusVerifikator,
			COALESCE(mp_t.catatan_verifikator, '') as CatatanVerifikator,
			COALESCE(NULLIF(mp_t.status_anggaran, ''), 'menunggu_verifikator') as StatusAnggaran,
			COALESCE(mp_t.catatan_anggaran, '') as CatatanAnggaran,
			mp_t.created_at as CreatedAt,
			mp_t.updated_at as UpdatedAt,
			COALESCE((SELECT COUNT(*) FROM aktivitasv2 a WHERE a.id_matrik_proker = mp_t.id), 0) as TotalAktivitas,
			COALESCE((SELECT SUM(a.anggaran) FROM aktivitasv2 a WHERE a.id_matrik_proker = mp_t.id), 0) as TotalAnggaran,
			COALESCE((SELECT SUM(CASE WHEN a.realisasi IS NOT NULL AND a.realisasi > 0 THEN a.realisasi ELSE a.anggaran END) FROM aktivitasv2 a WHERE a.id_matrik_proker = mp_t.id), 0) as TotalRealisasi,
			COALESCE((SELECT COUNT(*) FROM aktivitasv2 a WHERE a.id_matrik_proker = mp_t.id AND a.realisasi IS NOT NULL AND a.realisasi > 0), 0) > 0 as HasRealisasi,
			COALESCE((SELECT COUNT(*) FROM matrix_proker_dokumen d WHERE d.id_matrik_proker = mp_t.id AND d.deleted_at IS NULL), 0) as TotalDokumen
		`).
		Joins("LEFT JOIN mata_program prog ON mp_t.id_mata_program = prog.id").
		Joins("LEFT JOIN master_tahun mt ON prog.id_master_tahun = mt.id").
		Joins("LEFT JOIN master_indikator_renstra mir ON mp_t.id_indikator = mir.id").
		Joins("LEFT JOIN master_standar_renstra msr ON mir.id_master_standar = msr.id").
		Joins("LEFT JOIN v_fakultas_unit fu ON mp_t.id_fakultas_unit = fu.id").
		Where("mp_t.uuid = ?", uid).
		First(&row).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &row, nil
}

var allowedSearchColumns = map[string]string{
	"tahun":                "mt.tahun",
	"tahunuuid":            "mt.uuid",
	"filter_tahun":         "mt.tahun",
	"nama_program":         "prog.mata_program",
	"mataprogram":          "prog.mata_program",
	"program_kerja":        "prog.mata_program",
	"standar_sub_standar":  "msr.nama",
	"standar":              "msr.nama",
	"indikator_kinerja":    "mir.indikator",
	"indikator":            "mir.indikator",
	"output":               "mp_t.output",
	"outcome":              "mp_t.outcome",
	"id_fakultas_unit":     "mp_t.id_fakultas_unit",
	"fakultas_unit":        "fu.uuid",
	"fakultasunituuid":     "fu.uuid",
	"filter_fakultas_unit": "mp_t.id_fakultas_unit",
	"filter_prodi":         "fu.id",
	"id_prodi":             "fu.id",
	"prodi":                "fu.nama_fak_prod_unit",
	"filter_tab":           "filter_tab",
	"tab":                  "tab",
}

// ------------------------
// GET ALL (matrix_proker)
// ------------------------
func (r *ProkerV3Repository) GetAll(
	ctx context.Context,
	search string,
	searchFilters []commondomain.SearchFilter,
	page, limit *int,
) ([]domain.ProkerV3Default, int64, error) {

	var list = make([]domain.ProkerV3Default, 0)
	var total int64

	db := r.db.WithContext(ctx).
		Table("matrix_proker mp_t").
		Select(`
			mp_t.id as Id,
			mp_t.uuid as UUID,
			mp_t.id_fakultas_unit as IdFakultasUnit,
			COALESCE(fu.nama_fak_prod_unit, '') as NamaFakultasUnit,
			COALESCE(fu.uuid, '') as UuidFakultasUnit,
			COALESCE(fu.type, '') as TypeFakultasUnit,
			COALESCE(fu.fakultas, '') as ParentFakultas,
			mp_t.id_mata_program as IdMataProgram,
			COALESCE(prog.mata_program, '') as NamaProgram,
			COALESCE(mt.tahun, '2026') as Tahun,
			1 as NoProker,
			mp_t.id_indikator as IdIndikator,
			COALESCE(mir.indikator, '') as IndikatorKinerja,
			COALESCE(mir.uuid, '') as IndikatorUUID,
			COALESCE(mir.id_master_standar, msr.id, 0) as IdMasterStandar,
			COALESCE(msr.nama, '') as StandarSubStandar,
			COALESCE(msr.uuid, '') as StandarUUID,
			COALESCE(mp_t.target_unpak, '') as TargetUnpak,
			COALESCE(mp_t.capaian_terdahulu, '') as CapaianUnitLalu,
			COALESCE(mp_t.target_unit, '') as TargetUnit,
			COALESCE(mp_t.output, '') as Output,
			COALESCE(mp_t.outcome, '') as Outcome,
			COALESCE((SELECT a.pic FROM aktivitasv2 a WHERE a.id_matrik_proker = mp_t.id ORDER BY a.id ASC LIMIT 1), '') as PIC,
			COALESCE(NULLIF(mp_t.status, ''), 'draft') as Status,
			COALESCE(mp_t.catatan_penolakan, '') as CatatanPenolakan,
			COALESCE(NULLIF(mp_t.status_verifikator, ''), 'draft') as StatusVerifikator,
			COALESCE(mp_t.catatan_verifikator, '') as CatatanVerifikator,
			COALESCE(NULLIF(mp_t.status_anggaran, ''), 'menunggu_verifikator') as StatusAnggaran,
			COALESCE(mp_t.catatan_anggaran, '') as CatatanAnggaran,
			mp_t.created_at as CreatedAt,
			mp_t.updated_at as UpdatedAt,
			COALESCE((SELECT COUNT(*) FROM aktivitasv2 a WHERE a.id_matrik_proker = mp_t.id), 0) as TotalAktivitas,
			COALESCE((SELECT SUM(a.anggaran) FROM aktivitasv2 a WHERE a.id_matrik_proker = mp_t.id), 0) as TotalAnggaran,
			COALESCE((SELECT SUM(CASE WHEN a.realisasi IS NOT NULL AND a.realisasi > 0 THEN a.realisasi ELSE a.anggaran END) FROM aktivitasv2 a WHERE a.id_matrik_proker = mp_t.id), 0) as TotalRealisasi,
			COALESCE((SELECT COUNT(*) FROM aktivitasv2 a WHERE a.id_matrik_proker = mp_t.id AND a.realisasi IS NOT NULL AND a.realisasi > 0), 0) > 0 as HasRealisasi,
			COALESCE((SELECT COUNT(*) FROM matrix_proker_dokumen d WHERE d.id_matrik_proker = mp_t.id AND d.deleted_at IS NULL), 0) as TotalDokumen
		`).
		Joins("LEFT JOIN mata_program prog ON mp_t.id_mata_program = prog.id").
		Joins("LEFT JOIN master_tahun mt ON prog.id_master_tahun = mt.id").
		Joins("LEFT JOIN master_indikator_renstra mir ON mp_t.id_indikator = mir.id").
		Joins("LEFT JOIN master_standar_renstra msr ON mir.id_master_standar = msr.id").
		Joins("LEFT JOIN v_fakultas_unit fu ON mp_t.id_fakultas_unit = fu.id")

	// Apply Search Filters
	if len(searchFilters) > 0 {
		for _, f := range searchFilters {
			field := strings.TrimSpace(strings.ToLower(f.Field))
			operator := strings.TrimSpace(strings.ToLower(f.Operator))
			var val string
			if f.Value != nil {
				val = strings.TrimSpace(*f.Value)
			}

			if col, ok := allowedSearchColumns[field]; ok {
				if field == "filter_prodi" || field == "id_prodi" || field == "prodi" {
					if _, err := uuid.Parse(val); err == nil {
						db = db.Where("fu.uuid = ?", val)
					} else if valInt, err := strconv.Atoi(val); err == nil {
						db = db.Where("mp_t.id_fakultas_unit = ?", valInt)
					} else {
						db = db.Where("(mp_t.id_fakultas_unit = ? OR fu.nama_fak_prod_unit LIKE ?)", val, "%"+val+"%")
					}
					continue
				}

				if field == "filter_fakultas_unit" || field == "fakultas_unit" || field == "fakultasunituuid" || field == "id_fakultas_unit" {
					if strings.Contains(val, ",") {
						parts := strings.Split(val, ",")
						var ids []int
						for _, p := range parts {
							if id, err := strconv.Atoi(strings.TrimSpace(p)); err == nil {
								ids = append(ids, id)
							}
						}
						if len(ids) > 0 {
							db = db.Where("mp_t.id_fakultas_unit IN (?)", ids)
							continue
						}
					}
					if _, err := uuid.Parse(val); err == nil {
						db = db.Where("(fu.uuid = ? OR (fu.kode_fakultas IS NOT NULL AND fu.kode_fakultas != '' AND fu.kode_fakultas = (SELECT kode_fakultas FROM v_fakultas_unit WHERE uuid = ? LIMIT 1)) OR fu.fakultas = (SELECT nama_fak_prod_unit FROM v_fakultas_unit WHERE uuid = ? LIMIT 1) OR fu.fakultas = (SELECT fakultas FROM v_fakultas_unit WHERE uuid = ? LIMIT 1))", val, val, val, val)
					} else if valInt, err := strconv.Atoi(val); err == nil {
						db = db.Where("(mp_t.id_fakultas_unit = ? OR (fu.kode_fakultas IS NOT NULL AND fu.kode_fakultas != '' AND fu.kode_fakultas = (SELECT kode_fakultas FROM v_fakultas_unit WHERE id = ? LIMIT 1)) OR fu.fakultas = (SELECT nama_fak_prod_unit FROM v_fakultas_unit WHERE id = ? LIMIT 1) OR fu.fakultas = (SELECT fakultas FROM v_fakultas_unit WHERE id = ? LIMIT 1))", valInt, valInt, valInt, valInt)
					} else {
						db = db.Where("(mp_t.id_fakultas_unit = ? OR fu.uuid = ? OR fu.kode_fakultas = ? OR fu.nama_fak_prod_unit LIKE ? OR fu.fakultas LIKE ?)", val, val, val, "%"+val+"%", "%"+val+"%")
					}
					continue
				}

				if field == "tahun" || field == "filter_tahun" {
					if val == "2026" {
						db = db.Where("(mt.tahun = ? OR mt.tahun IS NULL OR mt.tahun = '')", val)
					} else {
						db = db.Where("mt.tahun = ?", val)
					}
					continue
				}

				if field == "tahunuuid" {
					db = db.Where("mt.uuid = ?", val)
					continue
				}

				if field == "filter_tab" || field == "tab" {
					parts := strings.Split(val, ":")
					var role, tab string
					if len(parts) == 2 {
						role = strings.ToLower(strings.TrimSpace(parts[0]))
						tab = strings.ToLower(strings.TrimSpace(parts[1]))
					} else {
						tab = strings.ToLower(strings.TrimSpace(val))
					}

					if role == "anggaran" || role == "keuangan" {
						switch tab {
						case "terima":
							db = db.Where("(mp_t.status_anggaran IN ('terima_anggaran', 'disetujui') OR mp_t.status IN ('terima_anggaran', 'terima anggaran'))")
						case "tolak":
							db = db.Where("(mp_t.status_anggaran = 'ditolak' OR (mp_t.status LIKE '%ditolak%' AND (mp_t.status_anggaran = 'ditolak' OR (mp_t.catatan_anggaran IS NOT NULL AND mp_t.catatan_anggaran != ''))))")
						case "menunggu":
							db = db.Where("COALESCE(mp_t.status_anggaran, '') NOT IN ('terima_anggaran', 'disetujui', 'ditolak') AND COALESCE(mp_t.status, '') NOT IN ('terima_anggaran', 'terima anggaran')")
						}
					} else {
						// Verifikator (upml, upmf, upms, evaluator, admin)
						switch tab {
						case "terima":
							db = db.Where("(mp_t.status_verifikator IN ('terima_verifikator', 'disetujui') OR mp_t.status IN ('disetujui', 'terima_anggaran', 'terima anggaran') OR mp_t.status_verifikator LIKE '%diverifikasi%' OR mp_t.status LIKE '%diverifikasi%')")
						case "tolak":
							db = db.Where("(mp_t.status_verifikator = 'ditolak' OR (mp_t.status LIKE '%ditolak%' AND (mp_t.status_verifikator = 'ditolak' OR (mp_t.catatan_verifikator IS NOT NULL AND mp_t.catatan_verifikator != '') OR (mp_t.catatan_penolakan IS NOT NULL AND mp_t.catatan_penolakan != ''))))")
						case "menunggu":
							db = db.Where("COALESCE(mp_t.status_verifikator, '') NOT IN ('terima_verifikator', 'disetujui', 'ditolak') AND COALESCE(mp_t.status, '') NOT IN ('disetujui', 'terima_anggaran', 'terima anggaran') AND COALESCE(mp_t.status_verifikator, '') NOT LIKE '%diverifikasi%' AND COALESCE(mp_t.status, '') NOT LIKE '%diverifikasi%'")
						}
					}
					continue
				}

				switch operator {
				case "eq", "=":
					db = db.Where(fmt.Sprintf("%s = ?", col), val)
				case "like", "contains":
					db = db.Where(fmt.Sprintf("%s LIKE ?", col), "%"+val+"%")
				default:
					db = db.Where(fmt.Sprintf("%s = ?", col), val)
				}
			}
		}
	}

	// General Search
	if search = strings.TrimSpace(search); search != "" {
		like := "%" + search + "%"
		db = db.Where(`(
			prog.mata_program LIKE ? OR
			mir.indikator LIKE ? OR
			msr.nama LIKE ? OR
			fu.nama_fak_prod_unit LIKE ? OR
			mp_t.output LIKE ? OR
			mp_t.outcome LIKE ?
		)`, like, like, like, like, like, like)
	}

	// Count Total
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Order by id DESC
	db = db.Order("mp_t.id DESC")

	// Pagination
	if limit != nil && *limit > 0 {
		db = db.Limit(*limit)
		if page != nil && *page > 0 {
			offset := (*page - 1) * (*limit)
			db = db.Offset(offset)
		}
	}

	if err := db.Find(&list).Error; err != nil {
		return nil, 0, err
	}

	for i := range list {
		list[i].NoProker = int64(i + 1)
	}

	return list, total, nil
}

// ------------------------
// CREATE (matrix_proker)
// ------------------------
func (r *ProkerV3Repository) Create(ctx context.Context, data *domain.ProkerV3) error {
	if data.UUID == uuid.Nil {
		data.UUID = uuid.New()
	}
	now := time.Now()
	data.CreatedAt = now
	data.UpdatedAt = now
	return r.db.WithContext(ctx).Create(data).Error
}

// ------------------------
// UPDATE (matrix_proker)
// ------------------------
func (r *ProkerV3Repository) Update(ctx context.Context, data *domain.ProkerV3) error {
	data.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(data).Error
}

// ------------------------
// DELETE (matrix_proker)
// ------------------------
func (r *ProkerV3Repository) Delete(ctx context.Context, uid uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var matrix domain.ProkerV3
		if err := tx.Where("uuid = ?", uid).First(&matrix).Error; err != nil {
			return err
		}

		// Find all child aktivitas
		var acts []domain.AktivitasV2
		if err := tx.Where("id_matrik_proker = ?", matrix.ID).Find(&acts).Error; err == nil {
			for _, act := range acts {
				// Delete timeline files
				_ = tx.Unscoped().Where("id_aktivitasv2 = ?", act.ID).Delete(&domain.ProkerTimelineFileV3{}).Error
				// Delete timelines
				_ = tx.Unscoped().Where("id_aktivitasv2 = ?", act.ID).Delete(&domain.ProkerTimelineV3{}).Error
				// Delete aktivitas
				_ = tx.Unscoped().Delete(&act).Error
			}
		}

		// Delete matrix_proker
		return tx.Unscoped().Delete(&matrix).Error
	})
}

// -------------------------------------------------------------
// AKTIVITAS V2 REPOSITORY METHODS
// -------------------------------------------------------------

func GenerateWeeklySlots(start, end time.Time) string {
	if start.IsZero() || end.IsZero() {
		return "[]"
	}
	months := []struct {
		key   string
		month time.Month
	}{
		{"JANUARI", time.January},
		{"PEBRUARI", time.February},
		{"MARET", time.March},
		{"APRIL", time.April},
		{"MEI", time.May},
		{"JUNI", time.June},
		{"JULI", time.July},
		{"AGUSTUS", time.August},
		{"SEPTEMBER", time.September},
		{"OKTOBER", time.October},
		{"NOVEMBER", time.November},
		{"DESEMBER", time.December},
	}

	weeks := []struct {
		week   string
		startD int
		endD   int
	}{
		{"M1", 1, 7},
		{"M2", 8, 14},
		{"M3", 15, 21},
		{"M4", 22, 31},
	}

	var slots []string
	year := start.Year()
	if year <= 1970 {
		year = 2026
	}

	for _, m := range months {
		for _, w := range weeks {
			wStart := time.Date(year, m.month, w.startD, 0, 0, 0, 0, time.Local)
			lastDayOfMonth := time.Date(year, m.month+1, 0, 23, 59, 59, 0, time.Local).Day()
			endDay := w.endD
			if endDay > lastDayOfMonth {
				endDay = lastDayOfMonth
			}
			wEnd := time.Date(year, m.month, endDay, 23, 59, 59, 0, time.Local)

			if wEnd.After(start) && wStart.Before(end) {
				slots = append(slots, fmt.Sprintf("%s_%s", m.key, w.week))
			}
		}
	}

	if len(slots) == 0 {
		return "[]"
	}
	b, err := json.Marshal(slots)
	if err != nil {
		return "[]"
	}
	return string(b)
}

func (r *ProkerV3Repository) GetAktivitasByMatrikId(ctx context.Context, matrikID uint) ([]domain.AktivitasV2Default, error) {
	var acts []domain.AktivitasV2
	err := r.db.WithContext(ctx).
		Where("id_matrik_proker = ?", matrikID).
		Order("id ASC").
		Find(&acts).Error
	if err != nil {
		return nil, err
	}

	result := make([]domain.AktivitasV2Default, 0, len(acts))
	for _, a := range acts {
		item := domain.AktivitasV2Default{
			ID:               a.ID,
			UUID:             a.UUID,
			IdMatrikProker:   a.IdMatrikProker,
			IdFakultasUnit:   a.IdFakultasUnit,
			Aktivitas:        a.Aktivitas,
			PIC:              a.PIC,
			Capaian:          a.Capaian,
			Masalah:          a.Masalah,
			Anggaran:         a.Anggaran,
			Realisasi:        a.Realisasi,
			Status:             a.Status,
			CatatanPenolakan:   a.CatatanPenolakan,
			StatusVerifikator:  a.StatusVerifikator,
			CatatanVerifikator: a.CatatanVerifikator,
			StatusAnggaran:     a.StatusAnggaran,
			CatatanAnggaran:    a.CatatanAnggaran,
			StatusUpml:         a.StatusUpml,
			CatatanTolakUpml:   a.CatatanTolakUpml,
			StatusUpmf:         a.StatusUpmf,
			CatatanTolakUpmf:   a.CatatanTolakUpmf,
			StatusUpms:         a.StatusUpms,
			CatatanTolakUpms:   a.CatatanTolakUpms,
			CreatedAt:          a.CreatedAt,
			UpdatedAt:          a.UpdatedAt,
		}

		// Load timelines
		var tls []domain.ProkerTimelineV3
		_ = r.db.WithContext(ctx).
			Where("id_aktivitasv2 = ?", a.ID).
			Order("no_aktivitas ASC, id ASC").
			Find(&tls).Error

		item.TimelinesCount = len(tls)
		item.Timelines = make([]domain.ProkerTimelineDefault, 0, len(tls))

		for _, tl := range tls {
			var rAwal, rAkhir *string
			if tl.TanggalRAwal != nil {
				s := tl.TanggalRAwal.Format("2006-01-02")
				rAwal = &s
			}
			if tl.TanggalRAkhir != nil {
				s := tl.TanggalRAkhir.Format("2006-01-02")
				rAkhir = &s
			}

			// Load files for timeline
			var rawFiles []domain.ProkerTimelineFileV3
			_ = r.db.WithContext(ctx).
				Where("id_aktivitasv2_timeline = ?", tl.ID).
				Find(&rawFiles).Error

			filesDefault := make([]domain.ProkerTimelineFileDefault, 0, len(rawFiles))
			for _, rf := range rawFiles {
				filesDefault = append(filesDefault, domain.ProkerTimelineFileDefault{
					ID:                    rf.ID,
					UUID:                  rf.UUID,
					IdAktivitasV2Timeline: rf.IdAktivitasV2Timeline,
					IdAktivitasV2:         rf.IdAktivitasV2,
					File:                  rf.File,
					NamaFile:              rf.NamaFile,
					FileURL:               rf.FileURL,
					FileSize:              rf.FileSize,
					FileType:              rf.FileType,
					Keterangan:            rf.Keterangan,
					CreatedAt:             rf.CreatedAt,
					UpdatedAt:             rf.UpdatedAt,
				})
			}

			jRk := strings.TrimSpace(tl.JadwalRk)
			if jRk == "" || jRk == "[]" {
				jRk = GenerateWeeklySlots(tl.TanggalRkAwal, tl.TanggalRkAkhir)
			}

			jR := strings.TrimSpace(tl.JadwalR)
			if (jR == "" || jR == "[]") && tl.TanggalRAwal != nil && tl.TanggalRAkhir != nil {
				jR = GenerateWeeklySlots(*tl.TanggalRAwal, *tl.TanggalRAkhir)
			}

			item.Timelines = append(item.Timelines, domain.ProkerTimelineDefault{
				ID:               tl.ID,
				UUID:             tl.UUID,
				IdAktivitasV2:    tl.IdAktivitasV2,
				TanggalRkAwal:    tl.TanggalRkAwal.Format("2006-01-02"),
				TanggalRkAkhir:   tl.TanggalRkAkhir.Format("2006-01-02"),
				TanggalRAwal:     rAwal,
				TanggalRAkhir:    rAkhir,
				NoAktivitas:      tl.NoAktivitas,
				NamaAktivitas:    tl.NamaAktivitas,
				PIC:              tl.PIC,
				TitikKontrol:     tl.TitikKontrol,
				JadwalRk:         jRk,
				JadwalR:          jR,
				Capaian:          tl.Capaian,
				Masalah:          tl.Masalah,
				PemecahanMasalah: tl.PemecahanMasalah,
				Status:           tl.Status,
				CreatedAt:        tl.CreatedAt,
				UpdatedAt:        tl.UpdatedAt,
				Files:            filesDefault,
			})
		}

		// Load Dokumen Referensi files (files with id_aktivitasv2_timeline is NULL)
		var docFiles []domain.ProkerTimelineFileV3
		_ = r.db.WithContext(ctx).
			Where("id_aktivitasv2 = ? AND (id_aktivitasv2_timeline IS NULL OR id_aktivitasv2_timeline = 0)", a.ID).
			Find(&docFiles).Error

		for _, df := range docFiles {
			item.DokumenReferensiFiles = append(item.DokumenReferensiFiles, domain.ProkerTimelineFileDefault{
				ID:                    df.ID,
				UUID:                  df.UUID,
				IdAktivitasV2Timeline: df.IdAktivitasV2Timeline,
				IdAktivitasV2:         df.IdAktivitasV2,
				File:                  df.File,
				NamaFile:              df.NamaFile,
				FileURL:               df.FileURL,
				FileSize:              df.FileSize,
				FileType:              df.FileType,
				Keterangan:            df.Keterangan,
				CreatedAt:             df.CreatedAt,
				UpdatedAt:             df.UpdatedAt,
			})
		}

		result = append(result, item)
	}

	return result, nil
}

func (r *ProkerV3Repository) GetAktivitasByUuid(ctx context.Context, uid uuid.UUID) (*domain.AktivitasV2Default, error) {
	var a domain.AktivitasV2
	err := r.db.WithContext(ctx).
		Where("uuid = ?", uid).
		First(&a).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	item := &domain.AktivitasV2Default{
		ID:               a.ID,
		UUID:             a.UUID,
		IdMatrikProker:   a.IdMatrikProker,
		IdFakultasUnit:   a.IdFakultasUnit,
		Aktivitas:        a.Aktivitas,
		PIC:              a.PIC,
		Capaian:          a.Capaian,
		Masalah:          a.Masalah,
		Anggaran:         a.Anggaran,
		Realisasi:        a.Realisasi,
		Status:             a.Status,
		CatatanPenolakan:   a.CatatanPenolakan,
		StatusVerifikator:  a.StatusVerifikator,
		CatatanVerifikator: a.CatatanVerifikator,
		StatusAnggaran:     a.StatusAnggaran,
		CatatanAnggaran:    a.CatatanAnggaran,
		StatusUpml:         a.StatusUpml,
		CatatanTolakUpml:   a.CatatanTolakUpml,
		StatusUpmf:         a.StatusUpmf,
		CatatanTolakUpmf:   a.CatatanTolakUpmf,
		StatusUpms:         a.StatusUpms,
		CatatanTolakUpms:   a.CatatanTolakUpms,
		CreatedAt:          a.CreatedAt,
		UpdatedAt:          a.UpdatedAt,
	}

	// Timelines
	var tls []domain.ProkerTimelineV3
	_ = r.db.WithContext(ctx).
		Where("id_aktivitasv2 = ?", a.ID).
		Order("no_aktivitas ASC, id ASC").
		Find(&tls).Error

	item.TimelinesCount = len(tls)
	item.Timelines = make([]domain.ProkerTimelineDefault, 0, len(tls))

	for _, tl := range tls {
		var rAwal, rAkhir *string
		if tl.TanggalRAwal != nil {
			s := tl.TanggalRAwal.Format("2006-01-02")
			rAwal = &s
		}
		if tl.TanggalRAkhir != nil {
			s := tl.TanggalRAkhir.Format("2006-01-02")
			rAkhir = &s
		}

		var rawFiles []domain.ProkerTimelineFileV3
		_ = r.db.WithContext(ctx).
			Where("id_aktivitasv2_timeline = ?", tl.ID).
			Find(&rawFiles).Error

		filesDefault := make([]domain.ProkerTimelineFileDefault, 0, len(rawFiles))
		for _, rf := range rawFiles {
			filesDefault = append(filesDefault, domain.ProkerTimelineFileDefault{
				ID:                    rf.ID,
				UUID:                  rf.UUID,
				IdAktivitasV2Timeline: rf.IdAktivitasV2Timeline,
				IdAktivitasV2:         rf.IdAktivitasV2,
				File:                  rf.File,
				NamaFile:              rf.NamaFile,
				FileURL:               rf.FileURL,
				FileSize:              rf.FileSize,
				FileType:              rf.FileType,
				Keterangan:            rf.Keterangan,
				CreatedAt:             rf.CreatedAt,
				UpdatedAt:             rf.UpdatedAt,
			})
		}

		jRk := strings.TrimSpace(tl.JadwalRk)
		if jRk == "" || jRk == "[]" {
			jRk = GenerateWeeklySlots(tl.TanggalRkAwal, tl.TanggalRkAkhir)
		}

		jR := strings.TrimSpace(tl.JadwalR)
		if (jR == "" || jR == "[]") && tl.TanggalRAwal != nil && tl.TanggalRAkhir != nil {
			jR = GenerateWeeklySlots(*tl.TanggalRAwal, *tl.TanggalRAkhir)
		}

		item.Timelines = append(item.Timelines, domain.ProkerTimelineDefault{
			ID:               tl.ID,
			UUID:             tl.UUID,
			IdAktivitasV2:    tl.IdAktivitasV2,
			TanggalRkAwal:    tl.TanggalRkAwal.Format("2006-01-02"),
			TanggalRkAkhir:   tl.TanggalRkAkhir.Format("2006-01-02"),
			TanggalRAwal:     rAwal,
			TanggalRAkhir:    rAkhir,
			NoAktivitas:      tl.NoAktivitas,
			NamaAktivitas:    tl.NamaAktivitas,
			PIC:              tl.PIC,
			TitikKontrol:     tl.TitikKontrol,
			JadwalRk:         jRk,
			JadwalR:          jR,
			Capaian:          tl.Capaian,
			Masalah:          tl.Masalah,
			PemecahanMasalah: tl.PemecahanMasalah,
			Status:           tl.Status,
			CreatedAt:        tl.CreatedAt,
			UpdatedAt:        tl.UpdatedAt,
			Files:            filesDefault,
		})
	}

	return item, nil
}

func (r *ProkerV3Repository) CreateAktivitas(ctx context.Context, data *domain.AktivitasV2, timelines []domain.ProkerTimelineV3) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if data.UUID == uuid.Nil {
			data.UUID = uuid.New()
		}
		if err := tx.Create(data).Error; err != nil {
			return err
		}

		for _, t := range timelines {
			t.IdAktivitasV2 = data.ID
			if t.UUID == uuid.Nil {
				t.UUID = uuid.New()
			}
			if err := tx.Create(&t).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *ProkerV3Repository) UpdateAktivitas(ctx context.Context, data *domain.AktivitasV2) error {
	data.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(data).Error
}

func (r *ProkerV3Repository) DeleteAktivitas(ctx context.Context, uid uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var a domain.AktivitasV2
		if err := tx.Where("uuid = ?", uid).First(&a).Error; err != nil {
			return err
		}

		// Delete timeline files
		_ = tx.Unscoped().Where("id_aktivitasv2 = ?", a.ID).Delete(&domain.ProkerTimelineFileV3{}).Error
		// Delete timelines
		_ = tx.Unscoped().Where("id_aktivitasv2 = ?", a.ID).Delete(&domain.ProkerTimelineV3{}).Error
		// Delete aktivitas
		return tx.Unscoped().Delete(&a).Error
	})
}

// -------------------------------------------------------------
// TIMELINE & FILES METHODS
// -------------------------------------------------------------

func (r *ProkerV3Repository) GetTimelineByUuid(ctx context.Context, uid uuid.UUID) (*domain.ProkerTimelineV3, error) {
	var tl domain.ProkerTimelineV3
	err := r.db.WithContext(ctx).
		Where("uuid = ?", uid).
		First(&tl).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &tl, nil
}

func (r *ProkerV3Repository) AddTimelineFiles(ctx context.Context, files []domain.ProkerTimelineFileV3) error {
	if len(files) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Create(&files).Error
}

func (r *ProkerV3Repository) DeleteTimelineFile(ctx context.Context, fileUid uuid.UUID) error {
	return r.db.WithContext(ctx).Unscoped().Where("uuid = ?", fileUid).Delete(&domain.ProkerTimelineFileV3{}).Error
}

func (r *ProkerV3Repository) GetTimelineFileByUuid(ctx context.Context, fileUid uuid.UUID) (*domain.ProkerTimelineFileV3, error) {
	var file domain.ProkerTimelineFileV3
	err := r.db.WithContext(ctx).Where("uuid = ?", fileUid).First(&file).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &file, nil
}
