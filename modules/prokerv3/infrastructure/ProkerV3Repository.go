package infrastructure

import (
	commondomain "UnpakSiamida/common/domain"
	domain "UnpakSiamida/modules/prokerv3/domain"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
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
			COALESCE(NULLIF(mp_t.pic, ''), (SELECT a.pic FROM aktivitasv2 a WHERE a.id_matrik_proker = mp_t.id ORDER BY a.id ASC LIMIT 1), '') as PIC,
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
	"pic":                  "mp_t.pic",
	"standar_sub_standar":  "msr.nama",
	"standar":              "msr.nama",
	"indikator_kinerja":    "mir.indikator",
	"indikator":            "mir.indikator",
	"output":               "mp_t.output",
	"outcome":              "mp_t.outcome",
	"id_fakultas_unit":     "mp_t.id_fakultas_unit",
	"fakultas_unit":        "fu.nama_fak_prod_unit",
	"nama_fakultas_unit":   "fu.nama_fak_prod_unit",
	"namafakultasunit":     "fu.nama_fak_prod_unit",
	"target_unit":          "fu.nama_fak_prod_unit",
	"target":               "fu.nama_fak_prod_unit",
	"fakultasunituuid":     "fu.uuid",
	"filter_fakultas_unit": "mp_t.id_fakultas_unit",
	"filter_prodi":         "fu.id",
	"id_prodi":             "fu.id",
	"prodi":                "fu.nama_fak_prod_unit",
	"target_unpak":         "mp_t.target_unpak",
	"target_capaian":       "mp_t.target_unpak",
	"capaian_unit_lalu":    "mp_t.capaian_terdahulu",
	"capaian_terdahulu":    "mp_t.capaian_terdahulu",
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

	// Pastikan tahun 2027 ada di master_tahun dan sinkronkan mata_program yang terkait 2027 (seperti id 15048)
	var count2027 int64
	_ = r.db.WithContext(ctx).Table("master_tahun").Where("tahun = ?", "2027").Count(&count2027).Error
	if count2027 == 0 {
		_ = r.db.WithContext(ctx).Exec("INSERT INTO master_tahun (uuid, tahun, status) VALUES (?, '2027', 'aktif')", uuid.New().String()).Error
	}
	_ = r.db.WithContext(ctx).Exec("UPDATE mata_program SET id_master_tahun = (SELECT id FROM master_tahun WHERE tahun = '2027' LIMIT 1) WHERE id = 15048").Error

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
			COALESCE(NULLIF(mp_t.pic, ''), (SELECT a.pic FROM aktivitasv2 a WHERE a.id_matrik_proker = mp_t.id ORDER BY a.id ASC LIMIT 1), '') as PIC,
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
				if field == "nama_program" || field == "program_kerja" || field == "mataprogram" {
					likeVal := "%" + strings.ToLower(val) + "%"
					db = db.Where("(LOWER(prog.mata_program) LIKE ? OR LOWER(COALESCE(mp_t.pic, '')) LIKE ?)", likeVal, likeVal)
					continue
				}

				if field == "pic" {
					db = db.Where("LOWER(COALESCE(mp_t.pic, '')) LIKE ?", "%"+strings.ToLower(val)+"%")
					continue
				}

				if field == "target_capaian" || field == "target_unpak" || field == "capaian_unit_lalu" || field == "capaian_terdahulu" {
					likeVal := "%" + strings.ToLower(val) + "%"
					db = db.Where("(LOWER(COALESCE(mp_t.target_unpak, '')) LIKE ? OR LOWER(COALESCE(mp_t.target_unit, '')) LIKE ? OR LOWER(COALESCE(mp_t.capaian_terdahulu, '')) LIKE ?)", likeVal, likeVal, likeVal)
					continue
				}

				if field == "filter_prodi" || field == "id_prodi" || field == "prodi" {
					if _, err := uuid.Parse(val); err == nil {
						db = db.Where("fu.uuid = ?", val)
					} else if valInt, err := strconv.Atoi(val); err == nil {
						db = db.Where("mp_t.id_fakultas_unit = ?", valInt)
					} else {
						db = db.Where("(mp_t.id_fakultas_unit = ? OR LOWER(fu.nama_fak_prod_unit) LIKE ?)", val, "%"+strings.ToLower(val)+"%")
					}
					continue
				}

				if field == "filter_fakultas_unit" || field == "fakultas_unit" || field == "fakultasunituuid" || field == "id_fakultas_unit" || field == "nama_fakultas_unit" || field == "namafakultasunit" || field == "target_unit" || field == "target" {
					if strings.Contains(val, ",") {
						parts := strings.Split(val, ",")
						var ids []int
						for _, p := range parts {
							if id, err := strconv.Atoi(strings.TrimSpace(p)); err == nil {
								ids = append(ids, id)
							}
						}
						if len(ids) > 0 {
							db = db.Where("(mp_t.id_fakultas_unit IN (?) OR (fu.fakultas IS NOT NULL AND fu.fakultas != '' AND LOWER(TRIM(fu.fakultas)) IN (SELECT LOWER(TRIM(COALESCE(nama_fak_prod_unit, ''))) FROM v_fakultas_unit WHERE id IN (?))))", ids, ids)
							continue
						}
					}
					if _, err := uuid.Parse(val); err == nil {
						db = db.Where("(fu.uuid = ? OR mp_t.id_fakultas_unit IN (SELECT id FROM v_fakultas_unit WHERE uuid = ?) OR (fu.fakultas IS NOT NULL AND fu.fakultas != '' AND LOWER(TRIM(fu.fakultas)) = (SELECT LOWER(TRIM(COALESCE(nama_fak_prod_unit, ''))) FROM v_fakultas_unit WHERE uuid = ? LIMIT 1)))", val, val, val)
					} else if valInt, err := strconv.Atoi(val); err == nil {
						if valInt == -1 {
							db = db.Where("1 = 0")
						} else {
							db = db.Where("(mp_t.id_fakultas_unit = ? OR (fu.fakultas IS NOT NULL AND fu.fakultas != '' AND LOWER(TRIM(fu.fakultas)) = (SELECT LOWER(TRIM(COALESCE(nama_fak_prod_unit, ''))) FROM v_fakultas_unit WHERE id = ? LIMIT 1)))", valInt, valInt)
						}
					} else {
						if val == "none" || val == "-1" {
							db = db.Where("1 = 0")
						} else {
							likeVal := "%" + strings.ToLower(val) + "%"
							db = db.Where("(mp_t.id_fakultas_unit = ? OR fu.uuid = ? OR LOWER(fu.nama_fak_prod_unit) LIKE ? OR LOWER(COALESCE(fu.fakultas, '')) LIKE ?)", val, val, likeVal, likeVal)
						}
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
					if len(parts) >= 2 {
						role = strings.ToLower(strings.TrimSpace(parts[0]))
						tab = strings.ToLower(strings.TrimSpace(parts[len(parts)-1]))
					} else {
						tab = strings.ToLower(strings.TrimSpace(val))
					}

					if role == "anggaran" || role == "keuangan" {
						switch tab {
						case "terima":
							db = db.Where("(mp_t.status_anggaran = 'terima_anggaran' OR LOWER(mp_t.status) IN ('terima_anggaran', 'terima anggaran'))")
						case "tolak":
							db = db.Where("mp_t.status_anggaran = 'ditolak'")
						case "menunggu":
							db = db.Where("COALESCE(mp_t.status, '') NOT IN ('draft', 'draf') AND COALESCE(mp_t.status_anggaran, '') != 'terima_anggaran' AND LOWER(COALESCE(mp_t.status, '')) NOT IN ('terima_anggaran', 'terima anggaran') AND COALESCE(mp_t.status_anggaran, '') != 'ditolak' AND COALESCE(mp_t.status_verifikator, '') != 'ditolak' AND COALESCE(mp_t.status_anggaran, '') != 'dibatalkan'")
						}
					} else {
						// Verifikator (upml, upmf, upms, evaluator, admin)
						switch tab {
						case "terima":
							db = db.Where("(mp_t.status_verifikator IN ('terima_verifikator', 'disetujui') OR mp_t.status IN ('disetujui', 'terima_anggaran', 'terima anggaran') OR mp_t.status_verifikator LIKE '%diverifikasi%' OR mp_t.status LIKE '%diverifikasi%')")
						case "tolak":
							db = db.Where("(mp_t.status_verifikator = 'ditolak' OR (mp_t.status LIKE '%ditolak%' AND (mp_t.status_verifikator = 'ditolak' OR (mp_t.catatan_verifikator IS NOT NULL AND mp_t.catatan_verifikator != '') OR (mp_t.catatan_penolakan IS NOT NULL AND mp_t.catatan_penolakan != ''))))")
						case "menunggu":
							db = db.Where("(mp_t.status = 'diajukan' OR mp_t.status_verifikator = 'diajukan' OR mp_t.status_verifikator = 'menunggu') AND COALESCE(mp_t.status_verifikator, '') NOT IN ('terima_verifikator', 'disetujui', 'ditolak') AND COALESCE(mp_t.status, '') NOT IN ('disetujui', 'terima_anggaran', 'terima anggaran', 'draft') AND COALESCE(mp_t.status_verifikator, '') NOT LIKE '%diverifikasi%' AND COALESCE(mp_t.status, '') NOT LIKE '%diverifikasi%'")
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
			ID:                 a.ID,
			UUID:               a.UUID,
			IdMatrikProker:     a.IdMatrikProker,
			IdFakultasUnit:     a.IdFakultasUnit,
			Aktivitas:          a.Aktivitas,
			PIC:                a.PIC,
			Capaian:            a.Capaian,
			Masalah:            a.Masalah,
			Anggaran:           a.Anggaran,
			Realisasi:          a.Realisasi,
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
				Anggaran:         tl.Anggaran,
				RAB:              tl.RAB,
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
		ID:                 a.ID,
		UUID:               a.UUID,
		IdMatrikProker:     a.IdMatrikProker,
		IdFakultasUnit:     a.IdFakultasUnit,
		Aktivitas:          a.Aktivitas,
		PIC:                a.PIC,
		Capaian:            a.Capaian,
		Masalah:            a.Masalah,
		Anggaran:           a.Anggaran,
		Realisasi:          a.Realisasi,
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
			Anggaran:         tl.Anggaran,
			RAB:              tl.RAB,
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
	err := r.db.WithContext(ctx).Where("uuid = ?", fileUid).First(&file)
	if err.Error != nil {
		if errors.Is(err.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err.Error
	}
	return &file, nil
}

func (r *ProkerV3Repository) GetAktivitasEntityByUuid(ctx context.Context, uid uuid.UUID) (*domain.AktivitasV2, error) {
	var a domain.AktivitasV2
	err := r.db.WithContext(ctx).Where("uuid = ?", uid).First(&a).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &a, nil
}

func (r *ProkerV3Repository) GetAktivitasEntityById(ctx context.Context, id uint) (*domain.AktivitasV2, error) {
	var a domain.AktivitasV2
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&a).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &a, nil
}

// ------------------------
// OPTIONS & TARGET UNPAK
// ------------------------

func (r *ProkerV3Repository) GetOptions(ctx context.Context) (*domain.ProkerOptionsResult, error) {
	var units []domain.ProkerUnitOption
	_ = r.db.WithContext(ctx).Table("v_fakultas_unit fu").
		Select("fu.id, fu.uuid, fu.nama_fak_prod_unit as nama, COALESCE(fu.type, '') as type, COALESCE(fu.fakultas, '') as fakultas, COALESCE(fu.jenjang, '') as jenjang, COALESCE(sfu.kode_fakultas, '') as kode_fakultas").
		Joins("LEFT JOIN sijamu_fakultas_unit sfu ON fu.id = sfu.id").
		Order("fu.nama_fak_prod_unit ASC").
		Scan(&units).Error

	var tahuns []domain.ProkerTahunOption
	_ = r.db.WithContext(ctx).Table("master_tahun").
		Select("id, uuid, tahun").
		Order("tahun DESC").
		Scan(&tahuns).Error

	var standars []domain.ProkerStandarOption
	_ = r.db.WithContext(ctx).Table("master_standar_renstra").
		Select("id, uuid, nama").
		Order("nama ASC").
		Scan(&standars).Error

	var indikators []domain.ProkerIndikatorOption
	_ = r.db.WithContext(ctx).Table("master_indikator_renstra").
		Select("id, uuid, id_master_standar, indikator, tahun").
		Order("indikator ASC").
		Scan(&indikators).Error

	return &domain.ProkerOptionsResult{
		Units:      units,
		Tahuns:     tahuns,
		Standars:   standars,
		Indikators: indikators,
	}, nil
}

func (r *ProkerV3Repository) GetTargetUnpak(ctx context.Context, tahun, indikatorParam, unitParam string) (*domain.TargetUnpakResult, error) {
	type TemplateRow struct {
		ID        int    `gorm:"column:id"`
		Target    string `gorm:"column:target"`
		TargetMin string `gorm:"column:target_min"`
		TargetMax string `gorm:"column:target_max"`
		Satuan    string `gorm:"column:satuan"`
	}

	formatTargetRow := func(t TemplateRow) string {
		var targetVal string
		target := strings.TrimSpace(t.Target)
		min := strings.TrimSpace(t.TargetMin)
		max := strings.TrimSpace(t.TargetMax)
		satuan := strings.TrimSpace(t.Satuan)

		if target != "" {
			targetVal = target
		} else if min != "" && max != "" {
			targetVal = fmt.Sprintf("%s - %s", min, max)
		} else if min != "" {
			targetVal = fmt.Sprintf(">= %s", min)
		} else if max != "" {
			targetVal = fmt.Sprintf("<= %s", max)
		}

		if targetVal != "" && satuan != "" {
			return fmt.Sprintf("%s %s", targetVal, satuan)
		}
		return targetVal
	}

	if indikatorParam == "" {
		return &domain.TargetUnpakResult{
			Found:       false,
			Target:      "-",
			TargetMin:   "",
			TargetMax:   "",
			Satuan:      "",
			TargetUnpak: "-",
		}, nil
	}

	indikatorID, err := strconv.Atoi(indikatorParam)
	if err != nil {
		var indRow struct{ ID int }
		if err := r.db.WithContext(ctx).Table("master_indikator_renstra").Select("id").Where("uuid = ?", indikatorParam).Take(&indRow).Error; err == nil {
			indikatorID = indRow.ID
		}
	}

	var unitID int
	if unitParam != "" {
		if id, err := strconv.Atoi(unitParam); err == nil {
			unitID = id
		} else {
			var uRow struct{ ID int }
			if err := r.db.WithContext(ctx).Table("v_fakultas_unit").Select("id").Where("uuid = ?", unitParam).Take(&uRow).Error; err == nil {
				unitID = uRow.ID
			}
		}
	}

	candidateIDs := []int{indikatorID}
	if indikatorID > 0 {
		var indRow struct {
			Indikator string `gorm:"column:indikator"`
		}
		if err := r.db.WithContext(ctx).Table("master_indikator_renstra").Select("indikator").Where("id = ?", indikatorID).First(&indRow).Error; err == nil && indRow.Indikator != "" {
			var sameIDs []int
			_ = r.db.WithContext(ctx).Table("master_indikator_renstra").Where("indikator = ?", indRow.Indikator).Pluck("id", &sameIDs).Error
			for _, sid := range sameIDs {
				if sid != indikatorID {
					candidateIDs = append(candidateIDs, sid)
				}
			}
		}
	}

	var tpl TemplateRow
	found := false

	if tahun != "" && unitID > 0 && len(candidateIDs) > 0 {
		if err := r.db.WithContext(ctx).Table("template_renstra").
			Where("indikator IN ? AND fakultas_unit = ? AND tahun = ?", candidateIDs, unitID, tahun).
			Order("id DESC").First(&tpl).Error; err == nil {
			if formatTargetRow(tpl) != "" {
				found = true
			}
		}
	}

	if !found && unitID > 0 && len(candidateIDs) > 0 {
		if err := r.db.WithContext(ctx).Table("template_renstra").
			Where("indikator IN ? AND fakultas_unit = ?", candidateIDs, unitID).
			Order("tahun DESC, id DESC").First(&tpl).Error; err == nil {
			if formatTargetRow(tpl) != "" {
				found = true
			}
		}
	}

	if !found && tahun != "" && len(candidateIDs) > 0 {
		if err := r.db.WithContext(ctx).Table("template_renstra").
			Where("indikator IN ? AND tahun = ?", candidateIDs, tahun).
			Order("id DESC").First(&tpl).Error; err == nil {
			if formatTargetRow(tpl) != "" {
				found = true
			}
		}
	}

	if !found && len(candidateIDs) > 0 {
		if err := r.db.WithContext(ctx).Table("template_renstra").
			Where("indikator IN ?", candidateIDs).
			Order("tahun DESC, id DESC").First(&tpl).Error; err == nil {
			if formatTargetRow(tpl) != "" {
				found = true
			}
		}
	}

	formatted := formatTargetRow(tpl)
	if formatted == "" {
		formatted = "-"
	}

	return &domain.TargetUnpakResult{
		Found:       found,
		Target:      tpl.Target,
		TargetMin:   tpl.TargetMin,
		TargetMax:   tpl.TargetMax,
		Satuan:      tpl.Satuan,
		TargetUnpak: formatted,
	}, nil
}

// ------------------------
// TIMELINE CRUD
// ------------------------

func (r *ProkerV3Repository) CountTimelineByAktivitasId(ctx context.Context, aktivitasId uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Table("aktivitasv2_timeline").Where("id_aktivitasv2 = ?", aktivitasId).Count(&count).Error
	return count, err
}

func (r *ProkerV3Repository) CreateTimeline(ctx context.Context, tl *domain.ProkerTimelineV3) error {
	if tl.UUID == uuid.Nil {
		tl.UUID = uuid.New()
	}
	return r.db.WithContext(ctx).Create(tl).Error
}

func (r *ProkerV3Repository) UpdateTimeline(ctx context.Context, tl *domain.ProkerTimelineV3) error {
	tl.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(tl).Error
}

func (r *ProkerV3Repository) DeleteTimeline(ctx context.Context, uid uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var tl domain.ProkerTimelineV3
		if err := tx.Where("uuid = ?", uid).First(&tl).Error; err != nil {
			return err
		}
		_ = tx.Unscoped().Where("id_aktivitasv2_timeline = ?", tl.ID).Delete(&domain.ProkerTimelineFileV3{}).Error
		return tx.Delete(&tl).Error
	})
}

func (r *ProkerV3Repository) ReorderTimeline(ctx context.Context, items []domain.ReorderTimelineItem) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, item := range items {
			if item.UUID != "" {
				if err := tx.Table("aktivitasv2_timeline").Where("uuid = ?", item.UUID).Update("no_aktivitas", item.NoAktivitas).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}

// ------------------------
// DOKUMEN CRUD
// ------------------------

func (r *ProkerV3Repository) GetDokumenByMatrikId(ctx context.Context, matrikID uint) ([]domain.MatrixProkerDokumen, error) {
	var docs []domain.MatrixProkerDokumen
	err := r.db.WithContext(ctx).Unscoped().Where("id_matrik_proker = ?", matrikID).Order("id DESC").Find(&docs).Error
	return docs, err
}

func (r *ProkerV3Repository) GetDokumenByUuid(ctx context.Context, uid uuid.UUID) (*domain.MatrixProkerDokumen, error) {
	var doc domain.MatrixProkerDokumen
	err := r.db.WithContext(ctx).Unscoped().Where("uuid = ?", uid).First(&doc).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &doc, nil
}

func (r *ProkerV3Repository) CreateDokumen(ctx context.Context, doc *domain.MatrixProkerDokumen) error {
	if doc.UUID == uuid.Nil {
		doc.UUID = uuid.New()
	}
	return r.db.WithContext(ctx).Create(doc).Error
}

func (r *ProkerV3Repository) DeleteDokumen(ctx context.Context, uid uuid.UUID) error {
	var doc domain.MatrixProkerDokumen
	if err := r.db.WithContext(ctx).Unscoped().Where("uuid = ?", uid).First(&doc).Error; err != nil {
		return err
	}
	filePath := doc.File
	if err := r.db.WithContext(ctx).Unscoped().Delete(&doc).Error; err != nil {
		return err
	}
	if filePath != "" {
		_ = os.Remove(filePath)
	}
	return nil
}

func (r *ProkerV3Repository) UpdateVerifikasiDokumen(ctx context.Context, uid uuid.UUID, status, catatan, catatanPenolakan string) (*domain.MatrixProkerDokumen, error) {
	var doc domain.MatrixProkerDokumen
	if err := r.db.WithContext(ctx).Where("uuid = ?", uid).First(&doc).Error; err != nil {
		return nil, err
	}

	status = strings.TrimSpace(status)
	if status == "" {
		status = "terverifikasi"
	}
	doc.StatusVerifikasi = status

	if status == "ditolak" || status == "gagal_terverifikasi" {
		alasan := strings.TrimSpace(catatanPenolakan)
		if alasan == "" {
			alasan = strings.TrimSpace(catatan)
		}
		doc.CatatanPenolakan = alasan
	} else if status == "terverifikasi" {
		doc.CatatanPenolakan = ""
	}
	doc.UpdatedAt = time.Now()

	if err := r.db.WithContext(ctx).Save(&doc).Error; err != nil {
		return nil, err
	}
	return &doc, nil
}

// ------------------------
// HISTORY STATE & STATUS
// ------------------------

func (r *ProkerV3Repository) RecordHistoryState(ctx context.Context, h *domain.HistoryState) error {
	if h.EntityUUID == "" || h.ToState == "" {
		return nil
	}
	if h.UUID == uuid.Nil {
		h.UUID = uuid.New()
	}
	if h.CreatedAt.IsZero() {
		h.CreatedAt = time.Now()
	}
	return r.db.WithContext(ctx).Create(h).Error
}

func (r *ProkerV3Repository) GetHistoryState(ctx context.Context, rawUuid string) ([]domain.HistoryState, error) {
	var mp domain.ProkerV3
	_ = r.db.WithContext(ctx).Where("uuid = ?", rawUuid).First(&mp).Error
	var act domain.AktivitasV2
	if mp.ID == 0 {
		_ = r.db.WithContext(ctx).Where("uuid = ?", rawUuid).First(&act).Error
	}

	var history []domain.HistoryState
	q := r.db.WithContext(ctx).Model(&domain.HistoryState{})
	if mp.ID > 0 {
		q = q.Where("entity_uuid = ? OR id_matrik_proker = ?", rawUuid, mp.ID)
	} else if act.ID > 0 {
		if act.IdMatrikProker != nil {
			q = q.Where("entity_uuid = ? OR id_aktivitas = ? OR id_matrik_proker = ?", rawUuid, act.ID, *act.IdMatrikProker)
		} else {
			q = q.Where("entity_uuid = ? OR id_aktivitas = ?", rawUuid, act.ID)
		}
	} else {
		q = q.Where("entity_uuid = ?", rawUuid)
	}
	_ = q.Order("created_at ASC, id ASC").Find(&history).Error

	if len(history) == 0 {
		currentStatus := "draft"
		currentStatusVerifikator := "draft"
		currentStatusAnggaran := "menunggu_verifikator"
		createdAt := time.Now()
		catatan := ""

		if mp.ID > 0 {
			currentStatus = mp.Status
			currentStatusVerifikator = mp.StatusVerifikator
			currentStatusAnggaran = mp.StatusAnggaran
			createdAt = mp.CreatedAt
			catatan = mp.CatatanPenolakan
			if catatan == "" {
				catatan = mp.CatatanAnggaran
			}
			if catatan == "" {
				catatan = mp.CatatanVerifikator
			}
		} else if act.ID > 0 {
			currentStatus = act.Status
			currentStatusVerifikator = act.StatusVerifikator
			currentStatusAnggaran = act.StatusAnggaran
			createdAt = act.CreatedAt
			catatan = act.CatatanPenolakan
			if catatan == "" {
				catatan = act.CatatanAnggaran
			}
			if catatan == "" {
				catatan = act.CatatanVerifikator
			}
		}

		history = append(history, domain.HistoryState{
			UUID:          uuid.New(),
			EntityUUID:    rawUuid,
			FromState:     "-",
			ToState:       "draf",
			Role:          "auditee",
			Action:        "draf",
			ActorName:     "Unit Pengusul",
			ActorUsername: "pengusul",
			UserName:      "Unit Pengusul",
			Catatan:       "Pembuatan draf awal program kerja",
			CreatedAt:     createdAt,
		})

		if currentStatus != "draft" && currentStatus != "" {
			history = append(history, domain.HistoryState{
				UUID:          uuid.New(),
				EntityUUID:    rawUuid,
				FromState:     "draf",
				ToState:       "diajukan",
				Role:          "auditee",
				Action:        "diajukan",
				ActorName:     "Unit Pengusul",
				ActorUsername: "pengusul",
				UserName:      "Unit Pengusul",
				Catatan:       "Pengajuan program kerja untuk diverifikasi",
				CreatedAt:     createdAt.Add(1 * time.Hour),
			})
		}

		if strings.Contains(currentStatusVerifikator, "tolak") || (strings.Contains(currentStatus, "tolak") && currentStatusAnggaran != "ditolak") {
			history = append(history, domain.HistoryState{
				UUID:          uuid.New(),
				EntityUUID:    rawUuid,
				FromState:     "diajukan",
				ToState:       "ditolak verifikator",
				Role:          "upml",
				Action:        "ditolak verifikator",
				ActorName:     "Verifikator UPM",
				ActorUsername: "verifikator",
				UserName:      "Verifikator UPM",
				Catatan:       catatan,
				CreatedAt:     createdAt.Add(2 * time.Hour),
			})
		} else if currentStatusVerifikator == "terima_verifikator" || currentStatusVerifikator == "disetujui" ||
			strings.Contains(currentStatus, "diverifikasi") || currentStatus == "disetujui" || currentStatus == "terima_anggaran" ||
			currentStatusAnggaran == "ditolak" || strings.Contains(currentStatusAnggaran, "tolak") {
			history = append(history, domain.HistoryState{
				UUID:          uuid.New(),
				EntityUUID:    rawUuid,
				FromState:     "diajukan",
				ToState:       "terima verifikator",
				Role:          "upml",
				Action:        "terima verifikator",
				ActorName:     "Verifikator UPM",
				ActorUsername: "verifikator",
				UserName:      "Verifikator UPM",
				Catatan:       "Diverifikasi dan disetujui oleh verifikator UPM",
				CreatedAt:     createdAt.Add(2 * time.Hour),
			})
		}

		if currentStatusAnggaran == "terima_anggaran" || currentStatus == "terima_anggaran" {
			history = append(history, domain.HistoryState{
				UUID:          uuid.New(),
				EntityUUID:    rawUuid,
				FromState:     "terima verifikator",
				ToState:       "terima anggaran",
				Role:          "anggaran",
				Action:        "terima anggaran",
				ActorName:     "Bagian Anggaran",
				ActorUsername: "anggaran",
				UserName:      "Bagian Anggaran",
				Catatan:       "Alokasi anggaran disetujui",
				CreatedAt:     createdAt.Add(3 * time.Hour),
			})
		} else if strings.Contains(currentStatusAnggaran, "tolak") {
			history = append(history, domain.HistoryState{
				UUID:          uuid.New(),
				EntityUUID:    rawUuid,
				FromState:     "terima verifikator",
				ToState:       "ditolak anggaran",
				Role:          "anggaran",
				Action:        "ditolak anggaran",
				ActorName:     "Bagian Anggaran",
				ActorUsername: "anggaran",
				UserName:      "Bagian Anggaran",
				Catatan:       catatan,
				CreatedAt:     createdAt.Add(3 * time.Hour),
			})
		}
	} else {
		for i := range history {
			if history[i].UserName == "" && history[i].ActorName != "" {
				history[i].UserName = history[i].ActorName
			}
			if history[i].ActorName == "" && history[i].UserName != "" {
				history[i].ActorName = history[i].UserName
			}
		}

		hasInitialDraft := false
		for _, h := range history {
			if strings.EqualFold(h.ToState, "draf") || strings.EqualFold(h.Action, "draf") {
				hasInitialDraft = true
				break
			}
		}
		if !hasInitialDraft {
			initCreatedAt := time.Now()
			if mp.ID > 0 && !mp.CreatedAt.IsZero() {
				initCreatedAt = mp.CreatedAt
			} else if act.ID > 0 && !act.CreatedAt.IsZero() {
				initCreatedAt = act.CreatedAt
			} else if len(history) > 0 && !history[0].CreatedAt.IsZero() {
				initCreatedAt = history[0].CreatedAt.Add(-30 * time.Minute)
			}
			initialDraft := domain.HistoryState{
				UUID:          uuid.New(),
				EntityUUID:    rawUuid,
				FromState:     "-",
				ToState:       "draf",
				Role:          "auditee",
				Action:        "draf",
				ActorName:     "Unit Pengusul",
				ActorUsername: "pengusul",
				UserName:      "Unit Pengusul",
				Catatan:       "Pembuatan draf awal program kerja",
				CreatedAt:     initCreatedAt,
			}
			history = append([]domain.HistoryState{initialDraft}, history...)
		}
	}
	return history, nil
}

func (r *ProkerV3Repository) UpdateStatus(ctx context.Context, req domain.UpdateStatusRequest) (*domain.UpdateStatusResult, error) {
	parsedUuid, err := uuid.Parse(req.UUID)
	if err != nil {
		return nil, errors.New("Invalid UUID")
	}

	newStatus := strings.ToLower(strings.TrimSpace(req.Status))
	if newStatus == "" {
		return nil, errors.New("Status tidak boleh kosong")
	}

	catatanPenolakan := strings.TrimSpace(req.CatatanPenolakan)
	level := strings.ToLower(strings.TrimSpace(req.Level))
	isVerifier := level == "upml" || level == "upmf" || level == "upms" || level == "admin" || level == "evaluator"
	isAnggaran := level == "anggaran" || level == "keuangan" || level == "admin"

	userName := req.UserName
	if userName == "" {
		userName = strings.ToUpper(level)
	}

	// 1. Try finding AktivitasV2
	var act domain.AktivitasV2
	errAct := r.db.WithContext(ctx).Where("uuid = ?", parsedUuid).First(&act).Error
	if errAct != nil {
		// 2. Try finding Matrix Proker
		var mp domain.ProkerV3
		if errMp := r.db.WithContext(ctx).Where("uuid = ?", parsedUuid).First(&mp).Error; errMp == nil {
			if !isVerifier && !isAnggaran && level != "admin" {
				if newStatus != "draft" && newStatus != "diajukan" {
					return nil, errors.New("Sebagai auditee, status hanya boleh Draf atau Diajukan")
				}
			}

			if isAnggaran && !isVerifier && level != "admin" {
				verifOk := mp.StatusVerifikator == "terima_verifikator" ||
					mp.StatusVerifikator == "disetujui" ||
					strings.Contains(mp.StatusVerifikator, "diverifikasi") ||
					strings.Contains(mp.Status, "diverifikasi") ||
					mp.Status == "disetujui" ||
					mp.Status == "terima_anggaran"
				if !verifOk && newStatus != "draft" && newStatus != "draf" {
					return nil, errors.New("Program kerja belum diverifikasi oleh verifikator (UPM)")
				}
			}

			updatesMp := map[string]interface{}{
				"updated_at": time.Now(),
			}
			updatesAct := map[string]interface{}{
				"updated_at": time.Now(),
			}
			var toStateHistory string
			var fromStateHistory string

			if (level == "anggaran" || level == "keuangan") && !isVerifier {
				if newStatus == "draft" || newStatus == "draf" {
					fromStateHistory = mp.Status
					if fromStateHistory == "" {
						fromStateHistory = "terima anggaran"
					}
					updatesMp["status_anggaran"] = "menunggu_verifikator"
					updatesMp["status_verifikator"] = "draft"
					updatesMp["status"] = "draft"
					if catatanPenolakan != "" {
						updatesMp["catatan_anggaran"] = catatanPenolakan
					}

					updatesAct["status_anggaran"] = "menunggu_verifikator"
					updatesAct["status_verifikator"] = "draft"
					updatesAct["status"] = "draft"
					if catatanPenolakan != "" {
						updatesAct["catatan_anggaran"] = catatanPenolakan
					}
					toStateHistory = "draf"
				} else if strings.Contains(newStatus, "tolak") {
					updatesMp["status_anggaran"] = "ditolak"
					updatesMp["catatan_anggaran"] = catatanPenolakan
					updatesMp["status"] = "ditolak"
					updatesMp["catatan_penolakan"] = catatanPenolakan

					updatesAct["status_anggaran"] = "ditolak"
					updatesAct["catatan_anggaran"] = catatanPenolakan
					updatesAct["status"] = "ditolak"
					updatesAct["catatan_penolakan"] = catatanPenolakan
					toStateHistory = "ditolak anggaran"
				} else {
					updatesMp["status_anggaran"] = "terima_anggaran"
					updatesMp["catatan_anggaran"] = ""
					updatesMp["status"] = "terima_anggaran"

					updatesAct["status_anggaran"] = "terima_anggaran"
					updatesAct["catatan_anggaran"] = ""
					updatesAct["status"] = "terima_anggaran"
					toStateHistory = "terima anggaran"
				}
			} else if isVerifier && (level == "upml" || level == "upmf" || level == "upms" || level == "evaluator") {
				if newStatus == "draft" || newStatus == "draf" {
					fromStateHistory = mp.Status
					if fromStateHistory == "" {
						fromStateHistory = "terima verifikator"
					}
					updatesMp["status_verifikator"] = "draft"
					updatesMp["status_anggaran"] = "menunggu_verifikator"
					updatesMp["status"] = "draft"
					if catatanPenolakan != "" {
						updatesMp["catatan_verifikator"] = catatanPenolakan
					}

					updatesAct["status_verifikator"] = "draft"
					updatesAct["status_anggaran"] = "menunggu_verifikator"
					updatesAct["status"] = "draft"
					if catatanPenolakan != "" {
						updatesAct["catatan_verifikator"] = catatanPenolakan
					}
					toStateHistory = "draf"
				} else if strings.Contains(newStatus, "tolak") {
					updatesMp["status_verifikator"] = "ditolak"
					updatesMp["catatan_verifikator"] = catatanPenolakan
					updatesMp["status_anggaran"] = "dibatalkan"
					updatesMp["status"] = "ditolak"
					updatesMp["catatan_penolakan"] = catatanPenolakan

					updatesAct["status_verifikator"] = "ditolak"
					updatesAct["catatan_verifikator"] = catatanPenolakan
					updatesAct["status_anggaran"] = "dibatalkan"
					updatesAct["status"] = "ditolak"
					updatesAct["catatan_penolakan"] = catatanPenolakan
					toStateHistory = "ditolak verifikator"
				} else {
					updatesMp["status_verifikator"] = "terima_verifikator"
					updatesMp["catatan_verifikator"] = ""
					updatesMp["status_anggaran"] = "menunggu_anggaran"
					updatesMp["status"] = newStatus

					updatesAct["status_verifikator"] = "terima_verifikator"
					updatesAct["catatan_verifikator"] = ""
					updatesAct["status_anggaran"] = "menunggu_anggaran"
					updatesAct["status"] = newStatus
					toStateHistory = "terima verifikator"
				}

				if level == "upml" {
					updatesMp["verified_upml_at"] = time.Now()
				} else if level == "upmf" {
					updatesMp["verified_upmf_at"] = time.Now()
				} else if level == "upms" || level == "evaluator" {
					updatesMp["verified_upms_at"] = time.Now()
				}
			} else {
				if newStatus == "diajukan" {
					updatesMp["status_verifikator"] = "diajukan"
					updatesMp["status_anggaran"] = "menunggu_verifikator"
					updatesMp["status"] = "diajukan"
					updatesAct["status_verifikator"] = "diajukan"
					updatesAct["status_anggaran"] = "menunggu_verifikator"
					updatesAct["status"] = "diajukan"
					fromStateHistory = "draf"
					toStateHistory = "diajukan"
				} else if newStatus == "draft" {
					updatesMp["status_verifikator"] = "draft"
					updatesMp["status_anggaran"] = "menunggu_verifikator"
					updatesMp["status"] = "draft"
					updatesAct["status_verifikator"] = "draft"
					updatesAct["status_anggaran"] = "menunggu_verifikator"
					updatesAct["status"] = "draft"
					fromStateHistory = "diajukan"
					toStateHistory = "draf"
				} else if newStatus == "terima_anggaran" {
					updatesMp["status_anggaran"] = "terima_anggaran"
					updatesMp["status"] = "terima_anggaran"
					updatesAct["status_anggaran"] = "terima_anggaran"
					updatesAct["status"] = "terima_anggaran"
					fromStateHistory = "terima verifikator"
					toStateHistory = "terima anggaran"
				} else if strings.Contains(newStatus, "diverifikasi") || newStatus == "disetujui" {
					updatesMp["status_verifikator"] = "terima_verifikator"
					updatesMp["status_anggaran"] = "menunggu_anggaran"
					updatesMp["status"] = newStatus
					updatesAct["status_verifikator"] = "terima_verifikator"
					updatesAct["status_anggaran"] = "menunggu_anggaran"
					updatesAct["status"] = newStatus
					fromStateHistory = "diajukan"
					toStateHistory = "terima verifikator"
				} else {
					updatesMp["status"] = newStatus
					updatesAct["status"] = newStatus
					toStateHistory = newStatus
				}
			}

			if err := r.db.WithContext(ctx).Model(&domain.ProkerV3{}).Where("id = ?", mp.ID).Updates(updatesMp).Error; err != nil {
				return nil, err
			}
			if len(updatesAct) > 1 {
				_ = r.db.WithContext(ctx).Model(&domain.AktivitasV2{}).Where("id_matrik_proker = ?", mp.ID).Updates(updatesAct).Error
			}

			mpIDInt := int(mp.ID)
			_ = r.RecordHistoryState(ctx, &domain.HistoryState{
				UUID:           uuid.New(),
				IdMatrikProker: &mpIDInt,
				EntityType:     "proker",
				EntityUUID:     req.UUID,
				FromState:      fromStateHistory,
				ToState:        toStateHistory,
				Role:           level,
				Action:         toStateHistory,
				ActorID:        req.UserID,
				ActorName:      userName,
				ActorUsername:  userName,
				UserID:         req.UserID,
				UserName:       userName,
				Catatan:        catatanPenolakan,
				CreatedAt:      time.Now(),
			})

			var count int64
			_ = r.db.WithContext(ctx).Model(&domain.AktivitasV2{}).Where("id_matrik_proker = ?", mp.ID).Count(&count)
			if count == 0 {
				var fuUint *uint64
				if mp.IdFakultasUnit != nil {
					v := uint64(*mp.IdFakultasUnit)
					fuUint = &v
				}
				mid := uint64(mp.ID)
				newAct := domain.AktivitasV2{
					UUID:              uuid.New(),
					IdMatrikProker:    &mid,
					IdFakultasUnit:    fuUint,
					Aktivitas:         "Pelaksanaan Program Kerja",
					Status:            newStatus,
					StatusVerifikator: "draft",
					StatusAnggaran:    "menunggu_verifikator",
					CatatanPenolakan:  catatanPenolakan,
					CreatedAt:         time.Now(),
					UpdatedAt:         time.Now(),
				}
				_ = r.db.WithContext(ctx).Create(&newAct)
			}

			return &domain.UpdateStatusResult{
				Message:          "Status program kerja berhasil diubah",
				Status:           newStatus,
				CatatanPenolakan: catatanPenolakan,
				UUID:             req.UUID,
			}, nil
		}
		return nil, errors.New("Data aktivitas/program kerja tidak ditemukan")
	}

	// Found AktivitasV2
	var fuID *int
	if act.IdFakultasUnit != nil {
		val := int(*act.IdFakultasUnit)
		fuID = &val
	} else if act.IdMatrikProker != nil {
		var parentMp domain.ProkerV3
		if errMp := r.db.WithContext(ctx).Select("id_fakultas_unit").Where("id = ?", *act.IdMatrikProker).First(&parentMp).Error; errMp == nil {
			fuID = parentMp.IdFakultasUnit
		}
	}

	if isVerifier && level != "admin" && act.IdMatrikProker != nil {
		var parentMp domain.ProkerV3
		if errMp := r.db.WithContext(ctx).Where("id = ?", *act.IdMatrikProker).First(&parentMp).Error; errMp == nil {
			isParentApproved := (parentMp.StatusAnggaran == "terima_anggaran" || parentMp.Status == "terima_anggaran" || parentMp.Status == "disetujui") &&
				(parentMp.StatusVerifikator == "terima_verifikator" || parentMp.StatusVerifikator == "disetujui" || strings.Contains(parentMp.StatusVerifikator, "diverifikasi") || parentMp.Status == "terima_anggaran")
			if !isParentApproved && newStatus != "draft" && newStatus != "draf" {
				return nil, errors.New("Aktivitas belum dapat diverifikasi: Program kerja harus disetujui oleh verifikator dan anggaran terlebih dahulu")
			}
		}
	}

	if !isVerifier && !isAnggaran && level != "admin" {
		if newStatus != "draft" && newStatus != "diajukan" {
			return nil, errors.New("Sebagai auditee, status hanya boleh Draf atau Diajukan")
		}
	}

	if isAnggaran && !isVerifier && level != "admin" {
		verifOk := act.StatusVerifikator == "terima_verifikator" ||
			act.StatusVerifikator == "disetujui" ||
			strings.Contains(act.StatusVerifikator, "diverifikasi") ||
			strings.Contains(act.Status, "diverifikasi") ||
			act.Status == "disetujui" ||
			act.Status == "terima_anggaran"
		if !verifOk && newStatus != "draft" && newStatus != "draf" {
			return nil, errors.New("Aktivitas belum diverifikasi oleh verifikator (UPM)")
		}
	}

	var isProdi bool
	if fuID != nil {
		var fu struct {
			Type string `gorm:"column:type"`
		}
		if err := r.db.WithContext(ctx).Table("v_fakultas_unit").Select("type").Where("id = ?", *fuID).First(&fu).Error; err == nil {
			isProdi = strings.ToLower(fu.Type) == "prodi"
		}
	}

	targetRole := strings.ToLower(strings.TrimSpace(req.Role))
	if targetRole == "" {
		targetRole = level
	}
	if targetRole == "evaluator" {
		targetRole = "upms"
	}
	if level == "admin" && (targetRole == "admin" || targetRole == "") {
		if act.StatusUpml != "disetujui" {
			targetRole = "upml"
		} else if isProdi && act.StatusUpmf != "disetujui" {
			targetRole = "upmf"
		} else {
			targetRole = "upms"
		}
	}

	actUpdates := map[string]interface{}{
		"updated_at": time.Now(),
	}
	var toStateHistoryAct string
	var fromStateHistoryAct string

	if (level == "anggaran" || level == "keuangan") && !isVerifier {
		fromStateHistoryAct = "terima verifikator"
		if newStatus == "draft" || newStatus == "draf" {
			actUpdates["status_anggaran"] = "menunggu_verifikator"
			actUpdates["status_verifikator"] = "draft"
			actUpdates["status"] = "draft"
			if catatanPenolakan != "" {
				actUpdates["catatan_anggaran"] = catatanPenolakan
			}
			toStateHistoryAct = "draf"
		} else if strings.Contains(newStatus, "tolak") {
			actUpdates["status_anggaran"] = "ditolak"
			actUpdates["catatan_anggaran"] = catatanPenolakan
			actUpdates["status"] = "ditolak"
			actUpdates["catatan_penolakan"] = catatanPenolakan
			toStateHistoryAct = "ditolak anggaran"
		} else {
			actUpdates["status_anggaran"] = "terima_anggaran"
			actUpdates["catatan_anggaran"] = ""
			actUpdates["status"] = "terima_anggaran"
			toStateHistoryAct = "terima anggaran"
		}
	} else if isVerifier && (targetRole == "upml" || targetRole == "upmf" || targetRole == "upms") {
		if level != "admin" && level != targetRole && !(level == "evaluator" && targetRole == "upms") {
			return nil, fmt.Errorf("Role %s tidak memiliki akses untuk verifikasi tahap %s", strings.ToUpper(level), strings.ToUpper(targetRole))
		}

		if newStatus == "draft" || newStatus == "draf" {
			fromStateHistoryAct = act.Status
			actUpdates["status_upml"] = "draft"
			actUpdates["status_upmf"] = "draft"
			actUpdates["status_upms"] = "draft"
			actUpdates["status_verifikator"] = "draft"
			actUpdates["status_anggaran"] = "menunggu_verifikator"
			actUpdates["status"] = "draft"
			if catatanPenolakan != "" {
				actUpdates["catatan_penolakan"] = catatanPenolakan
			}
			toStateHistoryAct = "draf"
		} else if targetRole == "upml" {
			fromStateHistoryAct = act.StatusUpml
			if fromStateHistoryAct == "" {
				fromStateHistoryAct = "menunggu"
			}
			if strings.Contains(newStatus, "tolak") {
				actUpdates["status_upml"] = "ditolak"
				actUpdates["catatan_tolak_upml"] = catatanPenolakan
				actUpdates["status"] = "ditolak"
				actUpdates["catatan_penolakan"] = catatanPenolakan
				toStateHistoryAct = "ditolak UPML"
			} else {
				actUpdates["status_upml"] = "disetujui"
				actUpdates["catatan_tolak_upml"] = ""
				if isProdi {
					actUpdates["status"] = "menunggu_upmf"
					if act.StatusUpmf != "disetujui" {
						actUpdates["status_upmf"] = "menunggu"
					}
				} else {
					actUpdates["status"] = "menunggu_upms"
					actUpdates["status_upmf"] = "-"
					if act.StatusUpms != "disetujui" {
						actUpdates["status_upms"] = "menunggu"
					}
				}
				actUpdates["catatan_penolakan"] = ""
				toStateHistoryAct = "disetujui UPML"
			}
		} else if targetRole == "upmf" {
			if !isProdi {
				return nil, errors.New("Aktivitas berasal dari Unit non-Prodi, verifikasi UPMF tidak berlaku (hanya UPML dan UPMS)")
			}
			if act.StatusUpml != "disetujui" {
				return nil, errors.New("Aktivitas belum disetujui oleh UPML. Verifikasi harus berurutan (UPML -> UPMF -> UPMS)")
			}
			fromStateHistoryAct = act.StatusUpmf
			if fromStateHistoryAct == "" {
				fromStateHistoryAct = "menunggu"
			}
			if strings.Contains(newStatus, "tolak") {
				actUpdates["status_upmf"] = "ditolak"
				actUpdates["catatan_tolak_upmf"] = catatanPenolakan
				actUpdates["status"] = "ditolak"
				actUpdates["catatan_penolakan"] = catatanPenolakan
				toStateHistoryAct = "ditolak UPMF"
			} else {
				actUpdates["status_upmf"] = "disetujui"
				actUpdates["catatan_tolak_upmf"] = ""
				actUpdates["status"] = "menunggu_upms"
				if act.StatusUpms != "disetujui" {
					actUpdates["status_upms"] = "menunggu"
				}
				actUpdates["catatan_penolakan"] = ""
				toStateHistoryAct = "disetujui UPMF"
			}
		} else if targetRole == "upms" {
			if act.StatusUpml != "disetujui" {
				return nil, errors.New("Aktivitas belum disetujui oleh UPML. Verifikasi harus berurutan")
			}
			if isProdi && act.StatusUpmf != "disetujui" {
				return nil, errors.New("Aktivitas prodi belum disetujui oleh UPMF. Verifikasi harus berurutan (UPML -> UPMF -> UPMS)")
			}
			fromStateHistoryAct = act.StatusUpms
			if fromStateHistoryAct == "" {
				fromStateHistoryAct = "menunggu"
			}
			if strings.Contains(newStatus, "tolak") {
				actUpdates["status_upms"] = "ditolak"
				actUpdates["catatan_tolak_upms"] = catatanPenolakan
				actUpdates["status"] = "ditolak"
				actUpdates["catatan_penolakan"] = catatanPenolakan
				toStateHistoryAct = "ditolak UPMS"
			} else {
				actUpdates["status_upms"] = "disetujui"
				actUpdates["catatan_tolak_upms"] = ""
				actUpdates["status"] = "disetujui"
				actUpdates["status_verifikator"] = "terima_verifikator"
				actUpdates["catatan_penolakan"] = ""
				toStateHistoryAct = "disetujui UPMS"
			}
		}
	} else {
		if newStatus == "diajukan" {
			actUpdates["status_verifikator"] = "diajukan"
			actUpdates["status_anggaran"] = "menunggu_verifikator"
			actUpdates["status"] = "diajukan"
			actUpdates["status_upml"] = "menunggu"
			actUpdates["catatan_tolak_upml"] = ""
			if isProdi {
				actUpdates["status_upmf"] = "menunggu"
			} else {
				actUpdates["status_upmf"] = "-"
			}
			actUpdates["catatan_tolak_upmf"] = ""
			actUpdates["status_upms"] = "menunggu"
			actUpdates["catatan_tolak_upms"] = ""
			actUpdates["catatan_penolakan"] = ""
			fromStateHistoryAct = act.Status
			toStateHistoryAct = "diajukan"
		} else if newStatus == "draft" {
			actUpdates["status_verifikator"] = "draft"
			actUpdates["status_anggaran"] = "menunggu_verifikator"
			actUpdates["status"] = "draft"
			fromStateHistoryAct = act.Status
			toStateHistoryAct = "draf"
		} else if newStatus == "terima_anggaran" {
			actUpdates["status_anggaran"] = "terima_anggaran"
			actUpdates["status"] = "terima_anggaran"
			fromStateHistoryAct = "terima verifikator"
			toStateHistoryAct = "terima anggaran"
		} else if strings.Contains(newStatus, "diverifikasi") || newStatus == "disetujui" {
			actUpdates["status_verifikator"] = "terima_verifikator"
			actUpdates["status_anggaran"] = "menunggu_anggaran"
			actUpdates["status"] = newStatus
			fromStateHistoryAct = "diajukan"
			toStateHistoryAct = "terima verifikator"
		} else {
			actUpdates["status"] = newStatus
			toStateHistoryAct = newStatus
		}
	}

	if err := r.db.WithContext(ctx).Model(&act).Updates(actUpdates).Error; err != nil {
		return nil, err
	}

	actIDInt := int(act.ID)
	var pMpID *int
	if act.IdMatrikProker != nil {
		val := int(*act.IdMatrikProker)
		pMpID = &val
	}
	_ = r.RecordHistoryState(ctx, &domain.HistoryState{
		UUID:           uuid.New(),
		IdMatrikProker: pMpID,
		IdAktivitas:    &actIDInt,
		EntityType:     "aktivitas",
		EntityUUID:     req.UUID,
		FromState:      fromStateHistoryAct,
		ToState:        toStateHistoryAct,
		Role:           level,
		Action:         toStateHistoryAct,
		ActorID:        req.UserID,
		ActorName:      userName,
		ActorUsername:  userName,
		UserID:         req.UserID,
		UserName:       userName,
		Catatan:        catatanPenolakan,
		CreatedAt:      time.Now(),
	})

	if act.IdMatrikProker != nil {
		mpUpdates := map[string]interface{}{
			"updated_at": time.Now(),
		}
		if targetRole == "upml" {
			mpUpdates["verified_upml_at"] = time.Now()
		} else if targetRole == "upmf" {
			mpUpdates["verified_upmf_at"] = time.Now()
		} else if targetRole == "upms" {
			mpUpdates["verified_upms_at"] = time.Now()
		}
		_ = r.db.WithContext(ctx).Model(&domain.ProkerV3{}).Where("id = ?", *act.IdMatrikProker).Updates(mpUpdates)
	}

	return &domain.UpdateStatusResult{
		Message:          "Status aktivitas berhasil diubah",
		Status:           newStatus,
		CatatanPenolakan: catatanPenolakan,
		UUID:             req.UUID,
	}, nil
}
