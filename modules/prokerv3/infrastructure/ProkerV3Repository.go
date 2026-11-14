package infrastructure

import (
	commondomain "UnpakSiamida/common/domain"
	commoninfra "UnpakSiamida/common/infrastructure"
	domain "UnpakSiamida/modules/prokerv3/domain"
	"context"
	"errors"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProkerV3Repository struct {
	db       *gorm.DB
	dbIku    *gorm.DB
	dbSiakun *gorm.DB
}

func NewProkerV3Repository(db *gorm.DB, dbs ...*gorm.DB) domain.IProkerV3Repository {
	repo := &ProkerV3Repository{db: db}
	if len(dbs) > 0 {
		repo.dbIku = dbs[0]
	}
	if len(dbs) > 1 {
		repo.dbSiakun = dbs[1]
	}
	return repo
}

func formatTanggalIndo(t time.Time) string {
	bulan := []string{
		"", "Januari", "Februari", "Maret", "April", "Mei", "Juni",
		"Juli", "Agustus", "September", "Oktober", "November", "Desember",
	}
	m := int(t.Month())
	if m >= 1 && m <= 12 {
		return fmt.Sprintf("%02d %s %d", t.Day(), bulan[m], t.Year())
	}
	return t.Format("02-01-2006")
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

func (r *ProkerV3Repository) populatePaguInfo(ctx context.Context, item *domain.ProkerV3Default) {
	if item == nil {
		return
	}
	fuId := 0
	if item.IdFakultasUnit != nil {
		fuId = *item.IdFakultasUnit
	}
	fuName := item.NamaFakultasUnit
	tahun := item.Tahun
	if tahun == "" {
		tahun = "2026"
	}

	var paguRow struct {
		Pagu float64
	}
	errPagu := r.db.WithContext(ctx).Table("pagu_anggaran").
		Select("pagu").
		Where("tahun = ? AND (id_fakultas_unit = ? OR (fakultas_unit IS NOT NULL AND fakultas_unit != '' AND (fakultas_unit = ? OR fakultas_unit = ?)))", tahun, fuId, strconv.Itoa(fuId), fuName).
		Order("id DESC").
		Limit(1).
		Scan(&paguRow).Error

	if errPagu == nil && paguRow.Pagu > 0 {
		nomPagu := paguRow.Pagu
		item.HasPagu = true
		item.PaguUnit = nomPagu

		// Hitung total anggaran akumulatif unit pada tahun tersebut (hanya non-yayasan dan tidak batal)
		// Prioritaskan bagian realisasi jika ada, fallback ke total rincian RAB jika belum ada realisasi
		var totalUnitAnggaran float64
		_ = r.db.WithContext(ctx).Table("aktivitasv2 a").
			Joins("JOIN matrix_proker mp ON mp.id = a.id_matrik_proker").
			Joins("LEFT JOIN mata_program p ON p.id = mp.id_mata_program").
			Joins("LEFT JOIN master_tahun mt ON mt.id = p.id_master_tahun").
			Where("mp.deleted_at IS NULL AND (mp.id_fakultas_unit = ? OR (SELECT fu2.nama_fak_prod_unit FROM v_fakultas_unit fu2 WHERE fu2.id = mp.id_fakultas_unit) = ?) AND COALESCE(mt.tahun, '2026') = ? AND COALESCE(mp.is_batal, 0) = 0 AND COALESCE(a.is_batal, 0) = 0 AND COALESCE(mp.beban_yayasan, 0) = 0", fuId, fuName, tahun).
			Select("COALESCE(SUM(CASE WHEN a.realisasi IS NOT NULL THEN a.realisasi ELSE COALESCE((SELECT SUM(COALESCE(r2.total, r2.qty_satuan * r2.qty_paket * r2.sbu, 0)) FROM aktivitasv2_rab r2 WHERE r2.id_aktivitasv2 = a.id), 0) END), 0)").
			Scan(&totalUnitAnggaran).Error

		item.TotalAnggaranUnit = totalUnitAnggaran
		if totalUnitAnggaran > nomPagu && !item.BebanYayasan && !item.IsBatal {
			item.IsMelebihiPagu = true
		} else {
			item.IsMelebihiPagu = false
		}
	} else {
		item.HasPagu = false
		item.PaguUnit = 0
		item.IsMelebihiPagu = false
	}
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
			'' as TargetUnpak,
			'' as CapaianUnitLalu,
			'' as TargetUnit,
			'' as Output,
			'' as Outcome,
			COALESCE(NULLIF(mp_t.pic, ''), (SELECT a.pic FROM aktivitasv2 a WHERE a.id_matrik_proker = mp_t.id ORDER BY a.id ASC LIMIT 1), '') as PIC,
			COALESCE(NULLIF(mp_t.status, ''), 'draft') as Status,
			COALESCE(mp_t.catatan_penolakan, '') as CatatanPenolakan,
			COALESCE(NULLIF(mp_t.status_verifikator, ''), 'draft') as StatusVerifikator,
			COALESCE(mp_t.catatan_verifikator, '') as CatatanVerifikator,
			COALESCE(NULLIF(mp_t.status_anggaran, ''), 'menunggu_verifikator') as StatusAnggaran,
			COALESCE(mp_t.catatan_anggaran, '') as CatatanAnggaran,
			COALESCE(mp_t.beban_yayasan, 0) as BebanYayasan,
			COALESCE(mp_t.is_batal, 0) as IsBatal,
			COALESCE(mp_t.catatan_batal, '') as CatatanBatal,
			mp_t.deleted_at as DeletedAt,
			mp_t.deleted_by as DeletedBy,
			COALESCE(mp_t.deleted_by_name, '') as DeletedByName,
			COALESCE(mp_t.catatan_hapus, '') as CatatanHapus,
			mp_t.created_at as CreatedAt,
			mp_t.updated_at as UpdatedAt,
			COALESCE((SELECT COUNT(*) FROM aktivitasv2 a WHERE a.id_matrik_proker = mp_t.id AND COALESCE(a.is_batal, 0) = 0), 0) as TotalAktivitas,
			COALESCE((SELECT SUM(COALESCE(r.total, r.qty_satuan * r.qty_paket * r.sbu, 0)) FROM aktivitasv2_rab r JOIN aktivitasv2 a ON a.id = r.id_aktivitasv2 WHERE a.id_matrik_proker = mp_t.id AND COALESCE(a.is_batal, 0) = 0), 0) as TotalAnggaran,
			COALESCE((SELECT SUM(CASE WHEN a.realisasi IS NOT NULL THEN a.realisasi ELSE COALESCE((SELECT SUM(COALESCE(r2.total, r2.qty_satuan * r2.qty_paket * r2.sbu, 0)) FROM aktivitasv2_rab r2 WHERE r2.id_aktivitasv2 = a.id), 0) END) FROM aktivitasv2 a WHERE a.id_matrik_proker = mp_t.id AND COALESCE(a.is_batal, 0) = 0), 0) as TotalRealisasi,
			COALESCE((SELECT COUNT(*) FROM aktivitasv2 a WHERE a.id_matrik_proker = mp_t.id AND a.realisasi IS NOT NULL AND COALESCE(a.is_batal, 0) = 0), 0) > 0 as HasRealisasi,
			COALESCE((SELECT COUNT(*) FROM matrix_proker_dokumen d WHERE d.id_matrik_proker = mp_t.id AND d.deleted_at IS NULL), 0) as TotalDokumen,
			COALESCE((SELECT COUNT(*) FROM risk_assessments ra WHERE ra.id_mata_program = mp_t.id_mata_program), 0) as TotalRisiko,
			COALESCE((SELECT ra.tingkat_risiko_awal FROM risk_assessments ra WHERE ra.id_mata_program = mp_t.id_mata_program ORDER BY FIELD(ra.tingkat_risiko_awal, 'T', 'M', 'R') ASC LIMIT 1), '') as TingkatRisiko,
			COALESCE((
				SELECT GROUP_CONCAT(iri.id_iku ORDER BY iri.id_iku SEPARATOR ',')
				FROM indikator_renstra_iku iri
				WHERE iri.id_indikator_renstra = mp_t.id_indikator
			), '') as IkuList
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

	r.populatePaguInfo(ctx, &row)

	var scopes []domain.ProkerScopeDefault
	_ = r.db.WithContext(ctx).Table("matrix_proker_scope mps").
		Select(`
			mps.id as id,
			mps.uuid as uuid,
			mps.id_matrix_proker as id_matrix_proker,
			mps.id_mata_program as id_mata_program,
			mps.id_indikator as id_indikator,
			COALESCE(mir.indikator, '') as indikator_kinerja,
			COALESCE(mir.uuid, '') as indikator_uuid,
			COALESCE(mir.id_master_standar, msr.id, 0) as id_master_standar,
			COALESCE(msr.nama, '') as standar_sub_standar,
			COALESCE(msr.uuid, '') as standar_uuid,
			COALESCE(mps.target_unpak, '') as target_unpak,
			COALESCE(mps.capaian_terdahulu, '') as capaian_unit_lalu,
			COALESCE(mps.target_unit, '') as target_unit,
			COALESCE(mps.output, '') as output,
			COALESCE(mps.outcome, '') as outcome,
			COALESCE((
				SELECT GROUP_CONCAT(iri.id_iku ORDER BY iri.id_iku SEPARATOR ',')
				FROM indikator_renstra_iku iri
				WHERE iri.id_indikator_renstra = mps.id_indikator
			), '') as iku_list
		`).
		Joins("LEFT JOIN master_indikator_renstra mir ON mps.id_indikator = mir.id").
		Joins("LEFT JOIN master_standar_renstra msr ON mir.id_master_standar = msr.id").
		Where("mps.id_matrix_proker = ?", row.ID).
		Order("mps.id ASC").
		Scan(&scopes).Error

	if scopes == nil {
		scopes = []domain.ProkerScopeDefault{}
	}

	ikuMap := commoninfra.GetMasterIkuMap(r.dbIku)
	row.IkuList = commoninfra.FormatIkuList(row.IkuList, ikuMap)
	for i := range scopes {
		scopes[i].IkuList = commoninfra.FormatIkuList(scopes[i].IkuList, ikuMap)
	}
	row.Scopes = scopes
	if len(scopes) > 0 {
		row.TargetUnpak = scopes[0].TargetUnpak
		row.CapaianUnitLalu = scopes[0].CapaianUnitLalu
		row.TargetUnit = scopes[0].TargetUnit
		row.Output = scopes[0].Output
		row.Outcome = scopes[0].Outcome
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
	"output":               "(SELECT mps.output FROM matrix_proker_scope mps WHERE mps.id_matrix_proker = mp_t.id ORDER BY mps.id ASC LIMIT 1)",
	"outcome":              "(SELECT mps.outcome FROM matrix_proker_scope mps WHERE mps.id_matrix_proker = mp_t.id ORDER BY mps.id ASC LIMIT 1)",
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
	"target_unpak":         "(SELECT mps.target_unpak FROM matrix_proker_scope mps WHERE mps.id_matrix_proker = mp_t.id ORDER BY mps.id ASC LIMIT 1)",
	"target_capaian":       "(SELECT mps.target_unpak FROM matrix_proker_scope mps WHERE mps.id_matrix_proker = mp_t.id ORDER BY mps.id ASC LIMIT 1)",
	"capaian_unit_lalu":    "(SELECT mps.capaian_terdahulu FROM matrix_proker_scope mps WHERE mps.id_matrix_proker = mp_t.id ORDER BY mps.id ASC LIMIT 1)",
	"capaian_terdahulu":    "(SELECT mps.capaian_terdahulu FROM matrix_proker_scope mps WHERE mps.id_matrix_proker = mp_t.id ORDER BY mps.id ASC LIMIT 1)",
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
	isBuangTab := false

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
			'' as TargetUnpak,
			'' as CapaianUnitLalu,
			'' as TargetUnit,
			'' as Output,
			'' as Outcome,
			COALESCE(NULLIF(mp_t.pic, ''), (SELECT a.pic FROM aktivitasv2 a WHERE a.id_matrik_proker = mp_t.id ORDER BY a.id ASC LIMIT 1), '') as PIC,
			COALESCE(NULLIF(mp_t.status, ''), 'draft') as Status,
			COALESCE(mp_t.catatan_penolakan, '') as CatatanPenolakan,
			COALESCE(NULLIF(mp_t.status_verifikator, ''), 'draft') as StatusVerifikator,
			COALESCE(mp_t.catatan_verifikator, '') as CatatanVerifikator,
			COALESCE(NULLIF(mp_t.status_anggaran, ''), 'menunggu_verifikator') as StatusAnggaran,
			COALESCE(mp_t.catatan_anggaran, '') as CatatanAnggaran,
			COALESCE(mp_t.beban_yayasan, 0) as BebanYayasan,
			COALESCE(mp_t.is_batal, 0) as IsBatal,
			COALESCE(mp_t.catatan_batal, '') as CatatanBatal,
			mp_t.deleted_at as DeletedAt,
			mp_t.deleted_by as DeletedBy,
			COALESCE(mp_t.deleted_by_name, '') as DeletedByName,
			COALESCE(mp_t.catatan_hapus, '') as CatatanHapus,
			mp_t.created_at as CreatedAt,
			mp_t.updated_at as UpdatedAt,
			COALESCE((SELECT COUNT(*) FROM aktivitasv2 a WHERE a.id_matrik_proker = mp_t.id AND COALESCE(a.is_batal, 0) = 0), 0) as TotalAktivitas,
			COALESCE((SELECT SUM(COALESCE(r.total, r.qty_satuan * r.qty_paket * r.sbu, 0)) FROM aktivitasv2_rab r JOIN aktivitasv2 a ON a.id = r.id_aktivitasv2 WHERE a.id_matrik_proker = mp_t.id AND COALESCE(a.is_batal, 0) = 0), 0) as TotalAnggaran,
			COALESCE((SELECT SUM(CASE WHEN a.realisasi IS NOT NULL THEN a.realisasi ELSE COALESCE((SELECT SUM(COALESCE(r2.total, r2.qty_satuan * r2.qty_paket * r2.sbu, 0)) FROM aktivitasv2_rab r2 WHERE r2.id_aktivitasv2 = a.id), 0) END) FROM aktivitasv2 a WHERE a.id_matrik_proker = mp_t.id AND COALESCE(a.is_batal, 0) = 0), 0) as TotalRealisasi,
			COALESCE((SELECT COUNT(*) FROM aktivitasv2 a WHERE a.id_matrik_proker = mp_t.id AND a.realisasi IS NOT NULL AND COALESCE(a.is_batal, 0) = 0), 0) > 0 as HasRealisasi,
			COALESCE((SELECT COUNT(*) FROM matrix_proker_dokumen d WHERE d.id_matrik_proker = mp_t.id AND d.deleted_at IS NULL), 0) as TotalDokumen,
			COALESCE((SELECT COUNT(*) FROM risk_assessments ra WHERE ra.id_mata_program = mp_t.id_mata_program), 0) as TotalRisiko,
			COALESCE((SELECT ra.tingkat_risiko_awal FROM risk_assessments ra WHERE ra.id_mata_program = mp_t.id_mata_program ORDER BY FIELD(ra.tingkat_risiko_awal, 'T', 'M', 'R') ASC LIMIT 1), '') as TingkatRisiko,
			COALESCE((
				SELECT GROUP_CONCAT(iri.id_iku ORDER BY iri.id_iku SEPARATOR ',')
				FROM indikator_renstra_iku iri
				WHERE iri.id_indikator_renstra = mp_t.id_indikator
			), '') as IkuList
		`).
		Joins("LEFT JOIN mata_program prog ON mp_t.id_mata_program = prog.id").
		Joins("LEFT JOIN master_tahun mt ON prog.id_master_tahun = mt.id").
		Joins("LEFT JOIN master_indikator_renstra mir ON mp_t.id_indikator = mir.id").
		Joins("LEFT JOIN master_standar_renstra msr ON mir.id_master_standar = msr.id").
		Joins("LEFT JOIN v_fakultas_unit fu ON mp_t.id_fakultas_unit = fu.id")

	// Apply Search Filters
	var activeRole string
	for _, f := range searchFilters {
		fField := strings.TrimSpace(strings.ToLower(f.Field))
		if fField == "active_role" || fField == "role" {
			if f.Value != nil {
				activeRole = strings.ToLower(strings.TrimSpace(*f.Value))
			}
		}
	}

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
					db = db.Where("EXISTS (SELECT 1 FROM matrix_proker_scope mps WHERE mps.id_matrix_proker = mp_t.id AND (LOWER(COALESCE(mps.target_unpak, '')) LIKE ? OR LOWER(COALESCE(mps.target_unit, '')) LIKE ? OR LOWER(COALESCE(mps.capaian_terdahulu, '')) LIKE ?))", likeVal, likeVal, likeVal)
					continue
				}

				if field == "output" || field == "outcome" {
					likeVal := "%" + strings.ToLower(val) + "%"
					db = db.Where("EXISTS (SELECT 1 FROM matrix_proker_scope mps WHERE mps.id_matrix_proker = mp_t.id AND (LOWER(COALESCE(mps.output, '')) LIKE ? OR LOWER(COALESCE(mps.outcome, '')) LIKE ?))", likeVal, likeVal)
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
					isStrictScoped := activeRole == "upml" || activeRole == "upms"
					if strings.Contains(val, ",") {
						parts := strings.Split(val, ",")
						var ids []int
						for _, p := range parts {
							if id, err := strconv.Atoi(strings.TrimSpace(p)); err == nil {
								ids = append(ids, id)
							}
						}
						if len(ids) > 0 {
							if isStrictScoped {
								db = db.Where("mp_t.id_fakultas_unit IN (?)", ids)
							} else {
								db = db.Where("(mp_t.id_fakultas_unit IN (?) OR (fu.fakultas IS NOT NULL AND fu.fakultas != '' AND LOWER(TRIM(fu.fakultas)) IN (SELECT LOWER(TRIM(COALESCE(nama_fak_prod_unit, ''))) FROM v_fakultas_unit WHERE id IN (?) AND type = 'fakultas')))", ids, ids)
							}
							continue
						}
					}
					if _, err := uuid.Parse(val); err == nil {
						if isStrictScoped {
							db = db.Where("(fu.uuid = ? OR mp_t.id_fakultas_unit IN (SELECT id FROM v_fakultas_unit WHERE uuid = ?))", val, val)
						} else {
							db = db.Where("(fu.uuid = ? OR mp_t.id_fakultas_unit IN (SELECT id FROM v_fakultas_unit WHERE uuid = ?) OR (fu.fakultas IS NOT NULL AND fu.fakultas != '' AND LOWER(TRIM(fu.fakultas)) = (SELECT LOWER(TRIM(COALESCE(nama_fak_prod_unit, ''))) FROM v_fakultas_unit WHERE uuid = ? AND type = 'fakultas' LIMIT 1)))", val, val, val)
						}
					} else if valInt, err := strconv.Atoi(val); err == nil {
						if valInt == -1 {
							db = db.Where("1 = 0")
						} else if isStrictScoped {
							db = db.Where("mp_t.id_fakultas_unit = ?", valInt)
						} else {
							db = db.Where("(mp_t.id_fakultas_unit = ? OR (fu.fakultas IS NOT NULL AND fu.fakultas != '' AND LOWER(TRIM(fu.fakultas)) = (SELECT LOWER(TRIM(COALESCE(nama_fak_prod_unit, ''))) FROM v_fakultas_unit WHERE id = ? AND type = 'fakultas' LIMIT 1)))", valInt, valInt)
						}
					} else {
						if val == "none" || val == "-1" {
							db = db.Where("1 = 0")
						} else if isStrictScoped {
							likeVal := "%" + strings.ToLower(val) + "%"
							db = db.Where("(mp_t.id_fakultas_unit = ? OR fu.uuid = ? OR LOWER(fu.nama_fak_prod_unit) LIKE ?)", val, val, likeVal)
						} else {
							likeVal := "%" + strings.ToLower(val) + "%"
							db = db.Where("(mp_t.id_fakultas_unit = ? OR fu.uuid = ? OR LOWER(fu.nama_fak_prod_unit) LIKE ? OR (fu.fakultas IS NOT NULL AND LOWER(COALESCE(fu.fakultas, '')) LIKE ? AND EXISTS (SELECT 1 FROM v_fakultas_unit fux WHERE (fux.uuid = ? OR LOWER(fux.nama_fak_prod_unit) LIKE ?) AND fux.type = 'fakultas')))", val, val, likeVal, likeVal, val, likeVal)
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

					if tab == "buang" || tab == "trash" || tab == "deleted" {
						isBuangTab = true
						continue
					}

					if role == "anggaran" || role == "keuangan" {
						switch tab {
						case "terima":
							db = db.Where("(mp_t.status_anggaran = 'terima_anggaran' OR LOWER(mp_t.status) IN ('terima_anggaran', 'terima anggaran'))")
						case "tolak":
							db = db.Where("mp_t.status_anggaran = 'ditolak'")
						case "revisi":
							db = db.Where("(mp_t.status = 'revisi' OR mp_t.status_anggaran = 'revisi' OR mp_t.status_verifikator = 'revisi')")
						case "menunggu":
							db = db.Where("COALESCE(mp_t.status, '') NOT IN ('draft', 'draf', 'revisi') AND COALESCE(mp_t.status_anggaran, '') != 'terima_anggaran' AND LOWER(COALESCE(mp_t.status, '')) NOT IN ('terima_anggaran', 'terima anggaran') AND COALESCE(mp_t.status_anggaran, '') != 'ditolak' AND COALESCE(mp_t.status_verifikator, '') != 'ditolak' AND COALESCE(mp_t.status_anggaran, '') != 'dibatalkan' AND COALESCE(mp_t.status_anggaran, '') != 'revisi' AND COALESCE(mp_t.status_verifikator, '') != 'revisi'")
						}
					} else {
						// Verifikator (upml, upmf, upms, evaluator, admin, user)
						switch tab {
						case "terima":
							db = db.Where("(mp_t.status_verifikator IN ('terima_verifikator', 'disetujui') OR mp_t.status IN ('disetujui', 'terima_anggaran', 'terima anggaran') OR mp_t.status_verifikator LIKE '%diverifikasi%' OR mp_t.status LIKE '%diverifikasi%')")
						case "tolak":
							db = db.Where("(mp_t.status_verifikator = 'ditolak' OR (mp_t.status LIKE '%ditolak%' AND (mp_t.status_verifikator = 'ditolak' OR (mp_t.catatan_verifikator IS NOT NULL AND mp_t.catatan_verifikator != '') OR (mp_t.catatan_penolakan IS NOT NULL AND mp_t.catatan_penolakan != ''))))")
						case "revisi":
							db = db.Where("(mp_t.status = 'revisi' OR mp_t.status_verifikator = 'revisi' OR mp_t.status_anggaran = 'revisi')")
						case "menunggu":
							db = db.Where("(mp_t.status = 'diajukan' OR mp_t.status_verifikator = 'diajukan' OR mp_t.status_verifikator = 'menunggu') AND COALESCE(mp_t.status_verifikator, '') NOT IN ('terima_verifikator', 'disetujui', 'ditolak', 'revisi') AND COALESCE(mp_t.status, '') NOT IN ('disetujui', 'terima_anggaran', 'terima anggaran', 'draft', 'revisi') AND COALESCE(mp_t.status_verifikator, '') NOT LIKE '%diverifikasi%' AND COALESCE(mp_t.status, '') NOT LIKE '%diverifikasi%'")
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
			EXISTS (SELECT 1 FROM matrix_proker_scope mps WHERE mps.id_matrix_proker = mp_t.id AND (mps.output LIKE ? OR mps.outcome LIKE ?))
		)`, like, like, like, like, like, like)
	}

	// Soft delete scope
	if isBuangTab {
		db = db.Where("mp_t.deleted_at IS NOT NULL")
	} else {
		db = db.Where("mp_t.deleted_at IS NULL")
	}

	// Count Total
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Order by deleted_at / id DESC
	if isBuangTab {
		db = db.Order("mp_t.deleted_at DESC, mp_t.id DESC")
	} else {
		db = db.Order("mp_t.id DESC")
	}

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

	if len(list) > 0 {
		var pIDs []uint
		for _, item := range list {
			pIDs = append(pIDs, item.ID)
		}
		var allScopes []domain.ProkerScopeDefault
		_ = r.db.WithContext(ctx).Table("matrix_proker_scope mps").
			Select(`
				mps.id as id,
				mps.uuid as uuid,
				mps.id_matrix_proker as id_matrix_proker,
				mps.id_mata_program as id_mata_program,
				mps.id_indikator as id_indikator,
				COALESCE(mir.indikator, '') as indikator_kinerja,
				COALESCE(mir.uuid, '') as indikator_uuid,
				COALESCE(mir.id_master_standar, msr.id, 0) as id_master_standar,
				COALESCE(msr.nama, '') as standar_sub_standar,
				COALESCE(msr.uuid, '') as standar_uuid,
				COALESCE(mps.target_unpak, '') as target_unpak,
				COALESCE(mps.capaian_terdahulu, '') as capaian_unit_lalu,
				COALESCE(mps.target_unit, '') as target_unit,
				COALESCE(mps.output, '') as output,
				COALESCE(mps.outcome, '') as outcome,
				COALESCE((
					SELECT GROUP_CONCAT(iri.id_iku ORDER BY iri.id_iku SEPARATOR ',')
					FROM indikator_renstra_iku iri
					WHERE iri.id_indikator_renstra = mps.id_indikator
				), '') as iku_list
			`).
			Joins("LEFT JOIN master_indikator_renstra mir ON mps.id_indikator = mir.id").
			Joins("LEFT JOIN master_standar_renstra msr ON mir.id_master_standar = msr.id").
			Where("mps.id_matrix_proker IN ?", pIDs).
			Order("mps.id ASC").
			Scan(&allScopes).Error

		scopeMap := make(map[uint][]domain.ProkerScopeDefault)
		for _, sc := range allScopes {
			scopeMap[sc.IdMatrixProker] = append(scopeMap[sc.IdMatrixProker], sc)
		}

		ikuMap := commoninfra.GetMasterIkuMap(r.dbIku)
		for i := range list {
			list[i].NoProker = int64(i + 1)
			list[i].IkuList = commoninfra.FormatIkuList(list[i].IkuList, ikuMap)
			if scs, ok := scopeMap[list[i].ID]; ok {
				for j := range scs {
					scs[j].IkuList = commoninfra.FormatIkuList(scs[j].IkuList, ikuMap)
				}
				list[i].Scopes = scs
				if len(scs) > 0 {
					list[i].TargetUnpak = scs[0].TargetUnpak
					list[i].CapaianUnitLalu = scs[0].CapaianUnitLalu
					list[i].TargetUnit = scs[0].TargetUnit
					list[i].Output = scs[0].Output
					list[i].Outcome = scs[0].Outcome
				}
			} else {
				list[i].Scopes = []domain.ProkerScopeDefault{}
			}
			r.populatePaguInfo(ctx, &list[i])
		}
		if !isBuangTab {
			r.populateDeletedSimilarity(ctx, list)
		}
	}

	return list, total, nil
}

func (r *ProkerV3Repository) populateDeletedSimilarity(ctx context.Context, list []domain.ProkerV3Default) {
	if len(list) == 0 {
		return
	}

	yearSet := make(map[string]struct{})
	var years []string
	for _, item := range list {
		y := strings.TrimSpace(item.Tahun)
		if y == "" {
			y = "2026"
		}
		if _, ok := yearSet[y]; !ok {
			yearSet[y] = struct{}{}
			years = append(years, y)
		}
	}
	if len(years) == 0 {
		return
	}

	type deletedCandidate struct {
		UUID           string     `gorm:"column:uuid"`
		IdFakultasUnit *int       `gorm:"column:id_fakultas_unit"`
		Tahun          string     `gorm:"column:tahun"`
		NamaProgram    string     `gorm:"column:nama_program"`
		CatatanHapus   string     `gorm:"column:catatan_hapus"`
		DeletedByName  string     `gorm:"column:deleted_by_name"`
		DeletedAt      *time.Time `gorm:"column:deleted_at"`
	}

	var candidates []deletedCandidate
	if err := r.db.WithContext(ctx).
		Table("matrix_proker mp_d").
		Select(`
			mp_d.uuid as uuid,
			mp_d.id_fakultas_unit as id_fakultas_unit,
			COALESCE(NULLIF(mt_d.tahun, ''), '2026') as tahun,
			COALESCE(prog_d.mata_program, '') as nama_program,
			COALESCE(mp_d.catatan_hapus, '') as catatan_hapus,
			COALESCE(mp_d.deleted_by_name, '') as deleted_by_name,
			mp_d.deleted_at as deleted_at
		`).
		Joins("LEFT JOIN mata_program prog_d ON mp_d.id_mata_program = prog_d.id").
		Joins("LEFT JOIN master_tahun mt_d ON prog_d.id_master_tahun = mt_d.id").
		Where("mp_d.deleted_at IS NOT NULL AND COALESCE(NULLIF(mt_d.tahun, ''), '2026') IN ?", years).
		Order("mp_d.deleted_at DESC, mp_d.id DESC").
		Scan(&candidates).Error; err != nil || len(candidates) == 0 {
		return
	}

	byYear := make(map[string][]deletedCandidate, len(years))
	for _, cand := range candidates {
		if strings.TrimSpace(cand.NamaProgram) == "" {
			continue
		}
		byYear[cand.Tahun] = append(byYear[cand.Tahun], cand)
	}

	for i := range list {
		title := strings.TrimSpace(list[i].NamaProgram)
		if title == "" {
			continue
		}
		y := strings.TrimSpace(list[i].Tahun)
		if y == "" {
			y = "2026"
		}
		cands := byYear[y]
		if len(cands) == 0 {
			continue
		}

		var matches []domain.DuplicateDeletedProker
		for _, cand := range cands {
			if list[i].IdFakultasUnit != nil && cand.IdFakultasUnit != nil && *list[i].IdFakultasUnit != *cand.IdFakultasUnit {
				continue
			}
			sim := commoninfra.ComputeTitleSimilarity(title, cand.NamaProgram)
			if sim >= commoninfra.TitleSimilarityThreshold {
				matches = append(matches, domain.DuplicateDeletedProker{
					UUID:          cand.UUID,
					NamaProgram:   cand.NamaProgram,
					CatatanHapus:  cand.CatatanHapus,
					DeletedByName: cand.DeletedByName,
					DeletedAt:     cand.DeletedAt,
					Similarity:    int(math.Round(sim * 100)),
				})
			}
		}

		if len(matches) > 0 {
			sort.SliceStable(matches, func(a, b int) bool {
				return matches[a].Similarity > matches[b].Similarity
			})
			list[i].DuplicateDeletedCount = int64(len(matches))
			list[i].DuplicateDeletedTitles = matches
		}
	}
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
// SOFT DELETE (matrix_proker)
// ------------------------
func (r *ProkerV3Repository) SoftDelete(ctx context.Context, params domain.SoftDeleteProkerParams) error {
	catatan := strings.TrimSpace(params.Catatan)
	if catatan == "" {
		return domain.DeleteReasonRequired()
	}
	role := strings.ToLower(strings.TrimSpace(params.Role))
	if role == "" {
		role = "user"
	}
	actorName := strings.TrimSpace(params.ActorName)
	if actorName == "" {
		actorName = strings.ToUpper(role)
	}
	actorUsername := strings.TrimSpace(params.ActorUsername)
	if actorUsername == "" {
		actorUsername = actorName
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var matrix domain.ProkerV3
		if err := tx.Unscoped().Where("uuid = ?", params.UUID).First(&matrix).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return domain.ProkerNotFound()
			}
			return err
		}
		if matrix.DeletedAt.Valid {
			return domain.AlreadyDeleted()
		}

		now := time.Now()
		updates := map[string]interface{}{
			"deleted_at":      now,
			"deleted_by":      params.ActorID,
			"deleted_by_name": actorName,
			"catatan_hapus":   catatan,
			"updated_at":      now,
		}
		if err := tx.Unscoped().Model(&domain.ProkerV3{}).Where("id = ?", matrix.ID).Updates(updates).Error; err != nil {
			return err
		}

		fromState := strings.TrimSpace(matrix.Status)
		if fromState == "" || strings.EqualFold(fromState, "draft") {
			fromState = "draf"
		}
		mpIDInt := int(matrix.ID)
		historyEntry := domain.HistoryState{
			UUID:           uuid.New(),
			IdMatrikProker: &mpIDInt,
			EntityType:     "proker",
			EntityUUID:     matrix.UUID.String(),
			FromState:      fromState,
			ToState:        "dihapus",
			Role:           role,
			Action:         "dihapus",
			ActorID:        params.ActorID,
			ActorName:      actorName,
			ActorUsername:  actorUsername,
			UserID:         params.ActorID,
			UserName:       actorName,
			Catatan:        catatan,
			CreatedAt:      now,
		}
		return tx.Create(&historyEntry).Error
	})
}

// ------------------------
// DELETE (matrix_proker)
// ------------------------
func (r *ProkerV3Repository) Delete(ctx context.Context, uid uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var matrix domain.ProkerV3
		if err := tx.Unscoped().Where("uuid = ?", uid).First(&matrix).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return domain.ProkerNotFound()
			}
			return err
		}

		// Find all child aktivitas IDs
		var actIDs []uint
		_ = tx.Table("aktivitasv2").Where("id_matrik_proker = ?", matrix.ID).Pluck("id", &actIDs).Error
		if len(actIDs) > 0 {
			_ = tx.Table("aktivitasv2_timeline_file").Where("id_aktivitasv2 IN (?)", actIDs).Delete(nil).Error
			_ = tx.Table("aktivitasv2_timeline").Where("id_aktivitasv2 IN (?)", actIDs).Delete(nil).Error
			_ = tx.Table("aktivitasv2_rab").Where("id_aktivitasv2 IN (?)", actIDs).Delete(nil).Error
			_ = tx.Table("aktivitasv2").Where("id IN (?)", actIDs).Delete(nil).Error
		}

		// Delete matrix_proker_dokumen
		_ = tx.Unscoped().Where("id_matrik_proker = ?", matrix.ID).Delete(&domain.MatrixProkerDokumen{}).Error

		// Delete matrix_proker_scope
		_ = tx.Unscoped().Where("id_matrix_proker = ?", matrix.ID).Delete(&domain.MatrixProkerScope{}).Error

		// Delete matrix_proker
		return tx.Unscoped().Delete(&matrix).Error
	})
}

// -------------------------------------------------------------
// DOKUMEN & AKTIVITAS OPTIONS
// -------------------------------------------------------------

func (r *ProkerV3Repository) GetAktivitasOptionsByMatrikId(ctx context.Context, matrikID uint) ([]domain.AktivitasDokumenOption, error) {
	var list []domain.AktivitasDokumenOption
	err := r.db.WithContext(ctx).Table("aktivitasv2").
		Select("id, uuid, aktivitas").
		Where("id_matrik_proker = ?", matrikID).
		Order("id ASC").
		Scan(&list).Error
	return list, err
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
	_ = r.db.WithContext(ctx).Unscoped().Where("uuid = ?", rawUuid).First(&mp).Error
	var act struct {
		ID               uint
		IdMatrikProker   *uint64   `gorm:"column:id_matrik_proker"`
		Status           string    `gorm:"column:status"`
		StatusLpm        string    `gorm:"column:status_lpm"`
		StatusUpml       string    `gorm:"column:status_upml"`
		StatusUpmf       string    `gorm:"column:status_upmf"`
		StatusUpms       string    `gorm:"column:status_upms"`
		StatusAnggaran   string    `gorm:"column:status_anggaran"`
		CreatedAt        time.Time `gorm:"column:created_at"`
		CatatanPenolakan string    `gorm:"column:catatan_penolakan"`
		CatatanAnggaran  string    `gorm:"column:catatan_anggaran"`
	}
	if mp.ID == 0 {
		_ = r.db.WithContext(ctx).Table("aktivitasv2").Where("uuid = ?", rawUuid).First(&act).Error
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
			if act.StatusLpm == "disetujui" {
				currentStatusVerifikator = "terima_verifikator"
			} else if act.StatusLpm == "ditolak" || act.StatusUpml == "ditolak" || act.StatusUpmf == "ditolak" || act.StatusUpms == "ditolak" {
				currentStatusVerifikator = "ditolak"
			} else {
				currentStatusVerifikator = "menunggu"
			}
			currentStatusAnggaran = act.StatusAnggaran
			createdAt = act.CreatedAt
			catatan = act.CatatanPenolakan
			if catatan == "" {
				catatan = act.CatatanAnggaran
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

		if mp.ID > 0 && mp.DeletedAt.Valid {
			delActor := strings.TrimSpace(mp.DeletedByName)
			if delActor == "" {
				delActor = "Unit Pengusul"
			}
			fromSt := strings.TrimSpace(currentStatus)
			if fromSt == "" || strings.EqualFold(fromSt, "draft") {
				fromSt = "draf"
			}
			history = append(history, domain.HistoryState{
				UUID:          uuid.New(),
				EntityUUID:    rawUuid,
				FromState:     fromSt,
				ToState:       "dihapus",
				Role:          "user",
				Action:        "dihapus",
				ActorID:       mp.DeletedBy,
				ActorName:     delActor,
				ActorUsername: delActor,
				UserID:        mp.DeletedBy,
				UserName:      delActor,
				Catatan:       mp.CatatanHapus,
				CreatedAt:     mp.DeletedAt.Time,
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
	if level == "fakultas" {
		level = "upmf"
	}
	isVerifier := level == "upml" || level == "upmf" || level == "upms" || level == "admin" || level == "evaluator"
	isAnggaran := level == "anggaran" || level == "keuangan" || level == "admin"

	userName := req.UserName
	if userName == "" {
		userName = strings.ToUpper(level)
	}

	var mp domain.ProkerV3
	if errMp := r.db.WithContext(ctx).Where("uuid = ?", parsedUuid).First(&mp).Error; errMp != nil {
		return nil, errors.New("Data program kerja tidak ditemukan")
	}

	if !isVerifier && !isAnggaran && level != "admin" {
		if newStatus != "draft" && newStatus != "diajukan" && newStatus != "revisi" {
			return nil, errors.New("Sebagai auditee, status hanya boleh Draf, Diajukan, atau Revisi")
		}
	}

	// Validasi Jadwal Proker V3 ketika mengajukan program kerja
	if newStatus == "diajukan" && level != "admin" {
		var prokerYear string
		if mp.IdMataProgram != nil {
			var mt struct {
				Tahun string `gorm:"column:tahun"`
			}
			_ = r.db.WithContext(ctx).Table("mata_program prog").
				Select("COALESCE(mt.tahun, '') as tahun").
				Joins("LEFT JOIN master_tahun mt ON prog.id_master_tahun = mt.id").
				Where("prog.id = ?", *mp.IdMataProgram).
				First(&mt).Error
			prokerYear = mt.Tahun
		}

		if prokerYear != "" {
			var jw domain.JadwalProkerV3
			foundSchedule := false

			// 1. Cek jadwal spesifik unit
			if mp.IdFakultasUnit != nil && *mp.IdFakultasUnit > 0 {
				if errJw := r.db.WithContext(ctx).
					Where("tahun = ? AND id_fakultas_unit = ?", prokerYear, *mp.IdFakultasUnit).
					Order("id DESC").First(&jw).Error; errJw == nil {
					foundSchedule = true
				} else {
					// Jika unit adalah prodi, cek jadwal parent fakultas
					var fu struct {
						Type     string
						Fakultas string
					}
					if errFu := r.db.WithContext(ctx).Table("v_fakultas_unit").
						Select("type, fakultas").Where("id = ?", *mp.IdFakultasUnit).Take(&fu).Error; errFu == nil {
						if strings.EqualFold(fu.Type, "prodi") && fu.Fakultas != "" {
							var parentFu struct{ ID uint64 }
							if errP := r.db.WithContext(ctx).Table("v_fakultas_unit").
								Select("id").
								Where("LOWER(TRIM(nama_fak_prod_unit)) = ? AND LOWER(type) = 'fakultas'", strings.ToLower(strings.TrimSpace(fu.Fakultas))).
								Take(&parentFu).Error; errP == nil {
								if errJwP := r.db.WithContext(ctx).
									Where("tahun = ? AND id_fakultas_unit = ?", prokerYear, parentFu.ID).
									Order("id DESC").First(&jw).Error; errJwP == nil {
									foundSchedule = true
								}
							}
						}
					}
				}
			}

			// 2. Fallback ke jadwal global (semua unit)
			if !foundSchedule {
				if errJw := r.db.WithContext(ctx).
					Where("tahun = ? AND (id_fakultas_unit IS NULL OR id_fakultas_unit = 0)", prokerYear).
					Order("id DESC").First(&jw).Error; errJw == nil {
					foundSchedule = true
				}
			}

			if foundSchedule {
				now := time.Now()
				startOfDay := time.Date(jw.TanggalMulai.Year(), jw.TanggalMulai.Month(), jw.TanggalMulai.Day(), 0, 0, 0, 0, jw.TanggalMulai.Location())
				endOfDay := time.Date(jw.TanggalAkhir.Year(), jw.TanggalAkhir.Month(), jw.TanggalAkhir.Day(), 23, 59, 59, 999999999, jw.TanggalAkhir.Location())

				if now.Before(startOfDay) {
					return nil, fmt.Errorf("Pengajuan program kerja tahun %s belum dibuka. Jadwal pengajuan: %s s/d %s", prokerYear, formatTanggalIndo(jw.TanggalMulai), formatTanggalIndo(jw.TanggalAkhir))
				}
				if now.After(endOfDay) {
					return nil, fmt.Errorf("Pengajuan program kerja tahun %s telah ditutup pada %s", prokerYear, formatTanggalIndo(jw.TanggalAkhir))
				}
			}
		}
	}

	if isAnggaran && !isVerifier && level != "admin" {
		verifOk := mp.StatusVerifikator == "terima_verifikator" ||
			mp.StatusVerifikator == "disetujui" ||
			strings.Contains(mp.StatusVerifikator, "diverifikasi") ||
			strings.Contains(mp.Status, "diverifikasi") ||
			mp.Status == "disetujui" ||
			mp.Status == "terima_anggaran"
		if !verifOk && newStatus != "draft" && newStatus != "draf" && newStatus != "revisi" {
			return nil, errors.New("Program kerja belum diverifikasi oleh verifikator")
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
			updatesMp["status_anggaran"] = "draft"
			updatesMp["status"] = "draft"
			updatesAct["status_anggaran"] = "draft"
			updatesAct["status"] = "draft"
			fromStateHistory = "terima verifikator"
			toStateHistory = "draf anggaran"
		} else if newStatus == "revisi" || strings.Contains(newStatus, "tolak") {
			updatesMp["status_anggaran"] = "ditolak"
			updatesMp["status"] = "revisi"
			if catatanPenolakan != "" {
				updatesMp["catatan_anggaran"] = catatanPenolakan
				updatesMp["catatan_penolakan"] = catatanPenolakan
			}
			updatesAct["status_anggaran"] = "ditolak"
			updatesAct["status"] = "revisi"
			if catatanPenolakan != "" {
				updatesAct["catatan_penolakan"] = catatanPenolakan
			}
			fromStateHistory = "terima verifikator"
			toStateHistory = "revisi anggaran"
		} else {
			updatesMp["status_anggaran"] = "terima_anggaran"
			updatesMp["catatan_anggaran"] = ""
			updatesMp["status"] = "terima_anggaran"
			updatesAct["status_anggaran"] = "terima_anggaran"
			updatesAct["status"] = "terima_anggaran"
			fromStateHistory = "terima verifikator"
			toStateHistory = "terima anggaran"
		}
	} else if isVerifier && level != "admin" {
		fromStateHistory = mp.Status
		if fromStateHistory == "" {
			fromStateHistory = "diajukan"
		}
		if strings.Contains(newStatus, "tolak") {
			updatesMp["status_verifikator"] = "ditolak"
			updatesMp["catatan_verifikator"] = catatanPenolakan
			updatesMp["status_anggaran"] = "dibatalkan"
			updatesMp["status"] = "ditolak"
			updatesMp["catatan_penolakan"] = catatanPenolakan

			updatesAct["status_anggaran"] = "dibatalkan"
			updatesAct["status"] = "ditolak"
			updatesAct["catatan_penolakan"] = catatanPenolakan
			toStateHistory = "ditolak verifikator"
		} else if newStatus == "revisi" || strings.Contains(newStatus, "revisi") {
			updatesMp["status_verifikator"] = "revisi"
			updatesMp["catatan_verifikator"] = catatanPenolakan
			updatesMp["status_anggaran"] = "menunggu_verifikator"
			updatesMp["status"] = "revisi"
			updatesMp["catatan_penolakan"] = catatanPenolakan

			updatesAct["status_anggaran"] = "menunggu_verifikator"
			updatesAct["status"] = "revisi"
			updatesAct["catatan_penolakan"] = catatanPenolakan
			toStateHistory = "revisi verifikator"
		} else if newStatus == "draft" || newStatus == "draf" {
			updatesMp["status_verifikator"] = "draft"
			updatesMp["catatan_verifikator"] = catatanPenolakan
			updatesMp["status_anggaran"] = "menunggu_verifikator"
			updatesMp["status"] = "draft"
			updatesMp["catatan_penolakan"] = catatanPenolakan

			updatesAct["status_anggaran"] = "menunggu_verifikator"
			updatesAct["status"] = "draft"
			updatesAct["catatan_penolakan"] = catatanPenolakan
			toStateHistory = "draf verifikator"
		} else {
			updatesMp["status_verifikator"] = "terima_verifikator"
			updatesMp["catatan_verifikator"] = ""
			updatesMp["status_anggaran"] = "menunggu_anggaran"
			updatesMp["status"] = newStatus

			updatesAct["status_anggaran"] = "menunggu_anggaran"
			updatesAct["status"] = newStatus
			toStateHistory = "terima verifikator"

			if level == "upml" {
				updatesMp["verified_upml_at"] = time.Now()
			} else if level == "upmf" {
				updatesMp["verified_upmf_at"] = time.Now()
			} else if level == "upms" || level == "evaluator" {
				updatesMp["verified_upms_at"] = time.Now()
			}
		}
	} else {
		if newStatus == "diajukan" {
			updatesMp["status_verifikator"] = "diajukan"
			updatesMp["status_anggaran"] = "menunggu_verifikator"
			updatesMp["status"] = "diajukan"
			updatesAct["status_anggaran"] = "menunggu_verifikator"
			updatesAct["status"] = "diajukan"
			fromStateHistory = mp.Status
			if fromStateHistory == "" {
				fromStateHistory = "draf"
			}
			toStateHistory = "diajukan"
		} else if newStatus == "draft" {
			updatesMp["status_verifikator"] = "draft"
			updatesMp["status_anggaran"] = "menunggu_verifikator"
			updatesMp["status"] = "draft"
			updatesAct["status_anggaran"] = "menunggu_verifikator"
			updatesAct["status"] = "draft"
			fromStateHistory = "diajukan"
			toStateHistory = "draf"
		} else if newStatus == "revisi" {
			fromStateHistory = mp.Status
			if fromStateHistory == "" {
				fromStateHistory = "diajukan"
			}
			updatesMp["status_verifikator"] = "revisi"
			updatesMp["status_anggaran"] = "menunggu_verifikator"
			updatesMp["status"] = "revisi"
			if catatanPenolakan != "" {
				updatesMp["catatan_verifikator"] = catatanPenolakan
				updatesMp["catatan_penolakan"] = catatanPenolakan
			}
			updatesAct["status_anggaran"] = "menunggu_verifikator"
			updatesAct["status"] = "revisi"
			if catatanPenolakan != "" {
				updatesAct["catatan_penolakan"] = catatanPenolakan
			}
			toStateHistory = "revisi"
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
		_ = r.db.WithContext(ctx).Table("aktivitasv2").Where("id_matrik_proker = ?", mp.ID).Updates(updatesAct).Error
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

	return &domain.UpdateStatusResult{
		Message:          "Status program kerja berhasil diubah",
		Status:           newStatus,
		CatatanPenolakan: catatanPenolakan,
		UUID:             req.UUID,
	}, nil
}

func (r *ProkerV3Repository) EnsureAppSettingsTable(ctx context.Context) error {
	if r.db == nil {
		return nil
	}
	return r.db.WithContext(ctx).Exec(`
		CREATE TABLE IF NOT EXISTS app_settings (
			id bigint(20) unsigned NOT NULL AUTO_INCREMENT,
			setting_key varchar(100) NOT NULL,
			setting_value text DEFAULT NULL,
			updated_by int(11) DEFAULT NULL,
			created_at timestamp NULL DEFAULT current_timestamp(),
			updated_at timestamp NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
			PRIMARY KEY (id),
			UNIQUE KEY uk_app_settings_key (setting_key)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
	`).Error
}

func (r *ProkerV3Repository) GetProkerBatalSettingStatus(ctx context.Context, now time.Time, targetUnitID *int) (*domain.ProkerBatalSettingStatus, error) {
	if r.db == nil {
		return nil, errors.New("Database tidak tersedia")
	}

	var row domain.AppSetting
	err := r.db.WithContext(ctx).
		Where("setting_key = ?", domain.ProkerBatalAktifTanggalKey).
		Take(&row).Error

	val := ""
	if err == nil && row.SettingValue != nil {
		val = *row.SettingValue
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	status := domain.EvalProkerBatalSettingWithUnit(val, now, targetUnitID)
	if err == nil && !row.UpdatedAt.IsZero() {
		status.UpdatedAt = row.UpdatedAt.Format("2006-01-02 15:04:05")
	}

	if targetUnitID != nil && *targetUnitID > 0 {
		var uUUID string
		_ = r.db.WithContext(ctx).Table("v_fakultas_unit").Select("uuid").Where("id = ?", *targetUnitID).Take(&uUUID).Error
		status.TargetUnitUUID = uUUID
	}

	if len(status.Units) > 0 {
		unitIDs := make([]int, len(status.Units))
		for i, u := range status.Units {
			unitIDs[i] = u.IdFakultasUnit
		}
		var fuRows []struct {
			ID   int    `gorm:"column:id"`
			UUID string `gorm:"column:uuid"`
		}
		if err := r.db.WithContext(ctx).Table("v_fakultas_unit").Select("id, uuid").Where("id IN (?)", unitIDs).Find(&fuRows).Error; err == nil {
			uuidMap := make(map[int]string)
			for _, rRow := range fuRows {
				uuidMap[rRow.ID] = rRow.UUID
			}
			for i := range status.Units {
				status.Units[i].UuidFakultasUnit = uuidMap[status.Units[i].IdFakultasUnit]
			}
		}
	}

	return &status, nil
}

func (r *ProkerV3Repository) GetProkerBatalSettingRaw(ctx context.Context) (string, error) {
	if r.db == nil {
		return "", errors.New("Database tidak tersedia")
	}
	var row domain.AppSetting
	err := r.db.WithContext(ctx).
		Where("setting_key = ?", domain.ProkerBatalAktifTanggalKey).
		Take(&row).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil
		}
		return "", err
	}
	if row.SettingValue != nil {
		return *row.SettingValue, nil
	}
	return "", nil
}

func (r *ProkerV3Repository) SaveProkerBatalSetting(ctx context.Context, serialized string, updatedBy *int) error {
	if r.db == nil {
		return errors.New("Database tidak tersedia")
	}
	return r.db.WithContext(ctx).Exec(`
		INSERT INTO app_settings (setting_key, setting_value, updated_by, created_at, updated_at)
		VALUES (?, ?, ?, NOW(), NOW())
		ON DUPLICATE KEY UPDATE
			setting_value = VALUES(setting_value),
			updated_by = VALUES(updated_by),
			updated_at = NOW()
	`, domain.ProkerBatalAktifTanggalKey, &serialized, updatedBy).Error
}

func (r *ProkerV3Repository) GetFakultasUnitNamesMap(ctx context.Context, unitIDs []int) (map[int]string, error) {
	namesMap := make(map[int]string)
	if len(unitIDs) == 0 || r.db == nil {
		return namesMap, nil
	}
	var fuRows []struct {
		ID   int    `gorm:"column:id"`
		Nama string `gorm:"column:nama_fak_prod_unit"`
	}
	err := r.db.WithContext(ctx).Table("v_fakultas_unit").
		Where("id IN (?)", unitIDs).
		Select("id, nama_fak_prod_unit").
		Find(&fuRows).Error
	if err != nil {
		return namesMap, err
	}
	for _, row := range fuRows {
		namesMap[row.ID] = row.Nama
	}
	return namesMap, nil
}

func (r *ProkerV3Repository) GetUserIDByUUID(ctx context.Context, uuid string) (*int, error) {
	if r.db == nil || strings.TrimSpace(uuid) == "" {
		return nil, nil
	}
	var uid int
	err := r.db.WithContext(ctx).Table("usersv3").Select("id").Where("uuid = ?", strings.TrimSpace(uuid)).Take(&uid).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &uid, nil
}

func (r *ProkerV3Repository) GetFakultasUnitIDByUUIDOrName(ctx context.Context, key string) (*int, error) {
	if r.db == nil || strings.TrimSpace(key) == "" {
		return nil, nil
	}
	var uid int
	err := r.db.WithContext(ctx).Table("v_fakultas_unit").Select("id").
		Where("uuid = ? OR nama_fak_prod_unit = ?", key, key).
		Take(&uid).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &uid, nil
}

// =======================================================
// Beban Yayasan & Batal Proker
// =======================================================

func (r *ProkerV3Repository) UpdateBebanYayasan(ctx context.Context, uid uuid.UUID, beban bool) (*domain.ProkerV3, error) {
	proker, err := r.GetByUuid(ctx, uid)
	if err != nil || proker == nil {
		return nil, domain.ProkerNotFound()
	}
	proker.BebanYayasan = beban
	proker.UpdatedAt = time.Now()
	if err := r.db.WithContext(ctx).Model(&domain.ProkerV3{}).Where("id = ?", proker.ID).Updates(map[string]interface{}{
		"beban_yayasan": beban,
		"updated_at":    proker.UpdatedAt,
	}).Error; err != nil {
		return nil, err
	}
	return proker, nil
}

func (r *ProkerV3Repository) UpdateBatalProker(ctx context.Context, uid uuid.UUID, isBatal bool, catatan string) (*domain.ProkerV3, error) {
	proker, err := r.GetByUuid(ctx, uid)
	if err != nil || proker == nil {
		return nil, domain.ProkerNotFound()
	}
	proker.IsBatal = isBatal
	proker.CatatanBatal = catatan
	proker.UpdatedAt = time.Now()

	err = r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&domain.ProkerV3{}).Where("id = ?", proker.ID).Updates(map[string]interface{}{
			"is_batal":      isBatal,
			"catatan_batal": catatan,
			"updated_at":    proker.UpdatedAt,
		}).Error; err != nil {
			return err
		}

		actUpdates := map[string]interface{}{
			"is_batal":   isBatal,
			"updated_at": time.Now(),
		}
		if isBatal {
			actUpdates["catatan_batal"] = catatan
		} else {
			actUpdates["catatan_batal"] = ""
		}
		return tx.Table("aktivitasv2").Where("id_matrik_proker = ?", proker.ID).Updates(actUpdates).Error
	})
	if err != nil {
		return nil, err
	}
	return proker, nil
}

// =======================================================
// COA
// =======================================================

func (r *ProkerV3Repository) GetCoa(ctx context.Context, search string) ([]domain.CoaItem, error) {
	targetDB := r.dbSiakun
	tableName := "m_coa"
	if targetDB == nil {
		targetDB = r.db
		tableName = "yps_siakun.m_coa"
	}
	if targetDB == nil {
		return nil, errors.New("Database not initialized")
	}

	combinedNoAkun := "CONCAT(TRIM(no_akun), IF(kode_unit IS NOT NULL AND kode_unit > 0, kode_unit, ''))"

	query := targetDB.WithContext(ctx).Table(tableName).
		Select("DISTINCT " + combinedNoAkun + " as no_akun, TRIM(nama_akun_id) as nama_akun").
		Where("LENGTH(TRIM(no_akun)) >= 6").
		Where("nama_akun_id IS NOT NULL AND TRIM(nama_akun_id) != ''")

	if search != "" {
		s := "%" + strings.ToLower(search) + "%"
		query = query.Where("("+combinedNoAkun+" LIKE ? OR LOWER(nama_akun_id) LIKE ?)", s, s)
	}

	var items []domain.CoaItem
	err := query.Order("no_akun ASC").Find(&items).Error
	return items, err
}

// =======================================================
// Jadwal Proker V3
// =======================================================

func (r *ProkerV3Repository) GetAllJadwalProkerV3(ctx context.Context, filter domain.JadwalFilter) ([]domain.JadwalProkerV3DTO, error) {
	if r.db == nil {
		return nil, errors.New("Database tidak tersedia")
	}
	var list []domain.JadwalProkerV3Row
	query := r.db.WithContext(ctx).
		Table("jadwal_proker_v3 j").
		Select("j.*, COALESCE(fu.nama_fak_prod_unit, '') as nama_fakultas_unit, COALESCE(fu.type, '') as type_fakultas_unit").
		Joins("LEFT JOIN v_fakultas_unit fu ON j.id_fakultas_unit = fu.id")

	if filter.Tahun != "" && filter.Tahun != "all" {
		query = query.Where("j.tahun = ?", filter.Tahun)
	}
	if filter.FilterUnit == "global" {
		query = query.Where("j.id_fakultas_unit IS NULL OR j.id_fakultas_unit = 0")
	} else if filter.FilterUnit == "specific" {
		query = query.Where("j.id_fakultas_unit IS NOT NULL AND j.id_fakultas_unit > 0")
	} else if filter.FilterUnit != "" && filter.FilterUnit != "all" {
		query = query.Where("j.id_fakultas_unit = ?", filter.FilterUnit)
	}
	if filter.Search != "" {
		likeSearch := "%" + filter.Search + "%"
		query = query.Where("(j.tahun LIKE ? OR fu.nama_fak_prod_unit LIKE ?)", likeSearch, likeSearch)
	}

	if err := query.Order("j.tahun DESC, (j.id_fakultas_unit IS NULL) DESC, j.tanggal_mulai ASC").Find(&list).Error; err != nil {
		return nil, err
	}

	now := time.Now()
	res := make([]domain.JadwalProkerV3DTO, 0, len(list))
	for _, item := range list {
		startOfDay := time.Date(item.TanggalMulai.Year(), item.TanggalMulai.Month(), item.TanggalMulai.Day(), 0, 0, 0, 0, item.TanggalMulai.Location())
		endOfDay := time.Date(item.TanggalAkhir.Year(), item.TanggalAkhir.Month(), item.TanggalAkhir.Day(), 23, 59, 59, 999999999, item.TanggalAkhir.Location())

		status := "berjalan"
		statusLabel := "Sedang Dibuka"
		if now.Before(startOfDay) {
			status = "belum_dibuka"
			statusLabel = "Belum Dibuka"
		} else if now.After(endOfDay) {
			status = "ditutup"
			statusLabel = "Ditutup"
		}

		namaUnit := item.NamaFakultasUnit
		typeUnit := item.TypeFakultasUnit
		if item.IdFakultasUnit == nil || *item.IdFakultasUnit == 0 {
			namaUnit = "Semua Unit / Fakultas (Global)"
			typeUnit = "all"
		}

		res = append(res, domain.JadwalProkerV3DTO{
			ID:               item.ID,
			UUID:             item.UUID.String(),
			Tahun:            item.Tahun,
			TanggalMulai:     item.TanggalMulai.Format("2006-01-02"),
			TanggalAkhir:     item.TanggalAkhir.Format("2006-01-02"),
			IdFakultasUnit:   item.IdFakultasUnit,
			NamaFakultasUnit: namaUnit,
			TypeFakultasUnit: typeUnit,
			Status:           status,
			StatusLabel:      statusLabel,
			CreatedAt:        item.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:        item.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	return res, nil
}

func (r *ProkerV3Repository) GetJadwalProkerV3(ctx context.Context, uid uuid.UUID) (*domain.JadwalProkerV3DTO, error) {
	if r.db == nil {
		return nil, errors.New("Database tidak tersedia")
	}
	var item domain.JadwalProkerV3Row
	if err := r.db.WithContext(ctx).
		Table("jadwal_proker_v3 j").
		Select("j.*, COALESCE(fu.nama_fak_prod_unit, '') as nama_fakultas_unit, COALESCE(fu.type, '') as type_fakultas_unit").
		Joins("LEFT JOIN v_fakultas_unit fu ON j.id_fakultas_unit = fu.id").
		Where("j.uuid = ?", uid).First(&item).Error; err != nil {
		return nil, err
	}

	now := time.Now()
	startOfDay := time.Date(item.TanggalMulai.Year(), item.TanggalMulai.Month(), item.TanggalMulai.Day(), 0, 0, 0, 0, item.TanggalMulai.Location())
	endOfDay := time.Date(item.TanggalAkhir.Year(), item.TanggalAkhir.Month(), item.TanggalAkhir.Day(), 23, 59, 59, 999999999, item.TanggalAkhir.Location())

	status := "berjalan"
	statusLabel := "Sedang Dibuka"
	if now.Before(startOfDay) {
		status = "belum_dibuka"
		statusLabel = "Belum Dibuka"
	} else if now.After(endOfDay) {
		status = "ditutup"
		statusLabel = "Ditutup"
	}

	namaUnit := item.NamaFakultasUnit
	typeUnit := item.TypeFakultasUnit
	if item.IdFakultasUnit == nil || *item.IdFakultasUnit == 0 {
		namaUnit = "Semua Unit / Fakultas (Global)"
		typeUnit = "all"
	}

	return &domain.JadwalProkerV3DTO{
		ID:               item.ID,
		UUID:             item.UUID.String(),
		Tahun:            item.Tahun,
		TanggalMulai:     item.TanggalMulai.Format("2006-01-02"),
		TanggalAkhir:     item.TanggalAkhir.Format("2006-01-02"),
		IdFakultasUnit:   item.IdFakultasUnit,
		NamaFakultasUnit: namaUnit,
		TypeFakultasUnit: typeUnit,
		Status:           status,
		StatusLabel:      statusLabel,
		CreatedAt:        item.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:        item.UpdatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

func (r *ProkerV3Repository) CheckJadwalProkerV3(ctx context.Context, req domain.CheckJadwalRequest) (*domain.CheckJadwalResult, error) {
	if r.db == nil {
		return nil, errors.New("Database tidak tersedia")
	}
	tahun := strings.TrimSpace(req.Tahun)
	if tahun == "" {
		tahun = "2027"
	}
	unitParam := strings.TrimSpace(req.UnitParam)

	var targetUnitID uint64
	var targetUnitName string
	var parentFakultasID uint64

	if unitParam != "" && unitParam != "all" && unitParam != "0" {
		var fu struct {
			ID       uint64 `gorm:"column:id"`
			Nama     string `gorm:"column:nama_fak_prod_unit"`
			Type     string `gorm:"column:type"`
			Fakultas string `gorm:"column:fakultas"`
		}
		var qFu *gorm.DB
		if _, errU := uuid.Parse(unitParam); errU == nil {
			qFu = r.db.WithContext(ctx).Table("v_fakultas_unit").Where("uuid = ?", unitParam)
		} else {
			qFu = r.db.WithContext(ctx).Table("v_fakultas_unit").Where("id = ?", unitParam)
		}
		if errTake := qFu.Take(&fu).Error; errTake == nil {
			targetUnitID = fu.ID
			targetUnitName = fu.Nama

			if strings.EqualFold(fu.Type, "prodi") && fu.Fakultas != "" {
				var parentFu struct {
					ID uint64 `gorm:"column:id"`
				}
				if errP := r.db.WithContext(ctx).Table("v_fakultas_unit").
					Select("id").
					Where("LOWER(TRIM(nama_fak_prod_unit)) = ? AND LOWER(type) = 'fakultas'", strings.ToLower(strings.TrimSpace(fu.Fakultas))).
					Take(&parentFu).Error; errP == nil {
					parentFakultasID = parentFu.ID
				}
			}
		}
	}

	var item domain.JadwalProkerV3
	found := false
	isUnitSpecific := false

	if targetUnitID > 0 {
		if errFind := r.db.WithContext(ctx).
			Where("tahun = ? AND id_fakultas_unit = ?", tahun, targetUnitID).
			Order("id DESC").First(&item).Error; errFind == nil {
			found = true
			isUnitSpecific = true
		}
	}

	if !found && parentFakultasID > 0 {
		if errFind := r.db.WithContext(ctx).
			Where("tahun = ? AND id_fakultas_unit = ?", tahun, parentFakultasID).
			Order("id DESC").First(&item).Error; errFind == nil {
			found = true
			isUnitSpecific = true
		}
	}

	if !found {
		if errFind := r.db.WithContext(ctx).
			Where("tahun = ? AND (id_fakultas_unit IS NULL OR id_fakultas_unit = 0)", tahun).
			Order("id DESC").First(&item).Error; errFind == nil {
			found = true
			isUnitSpecific = false
		}
	}

	if !found {
		return &domain.CheckJadwalResult{
			HasSchedule:    false,
			IsOpen:         true,
			Status:         "tidak_ada_jadwal",
			StatusLabel:    "Tidak Ada Pembatasan Jadwal",
			Tahun:          tahun,
			IsUnitSpecific: false,
			Message:        fmt.Sprintf("Belum ada jadwal khusus pengajuan program kerja untuk tahun %s", tahun),
		}, nil
	}

	now := time.Now()
	startOfDay := time.Date(item.TanggalMulai.Year(), item.TanggalMulai.Month(), item.TanggalMulai.Day(), 0, 0, 0, 0, item.TanggalMulai.Location())
	endOfDay := time.Date(item.TanggalAkhir.Year(), item.TanggalAkhir.Month(), item.TanggalAkhir.Day(), 23, 59, 59, 999999999, item.TanggalAkhir.Location())

	isOpen := true
	status := "berjalan"
	statusLabel := "Sedang Dibuka"

	unitSuffix := ""
	if isUnitSpecific && targetUnitName != "" {
		unitSuffix = fmt.Sprintf(" (%s)", targetUnitName)
	}

	message := fmt.Sprintf("Pengajuan program kerja tahun %s%s dibuka dari %s s/d %s", item.Tahun, unitSuffix, formatTanggalIndo(item.TanggalMulai), formatTanggalIndo(item.TanggalAkhir))

	if now.Before(startOfDay) {
		isOpen = false
		status = "belum_dibuka"
		statusLabel = "Belum Dibuka"
		message = fmt.Sprintf("Pengajuan program kerja tahun %s%s belum dibuka. Jadwal pengajuan: %s s/d %s", item.Tahun, unitSuffix, formatTanggalIndo(item.TanggalMulai), formatTanggalIndo(item.TanggalAkhir))
	} else if now.After(endOfDay) {
		isOpen = false
		status = "ditutup"
		statusLabel = "Ditutup"
		message = fmt.Sprintf("Pengajuan program kerja tahun %s%s telah ditutup pada %s", item.Tahun, unitSuffix, formatTanggalIndo(item.TanggalAkhir))
	}

	return &domain.CheckJadwalResult{
		HasSchedule:      true,
		IsOpen:           isOpen,
		Status:           status,
		StatusLabel:      statusLabel,
		Tahun:            item.Tahun,
		TanggalMulai:     item.TanggalMulai.Format("2006-01-02"),
		TanggalAkhir:     item.TanggalAkhir.Format("2006-01-02"),
		IdFakultasUnit:   item.IdFakultasUnit,
		IsUnitSpecific:   isUnitSpecific,
		NamaFakultasUnit: targetUnitName,
		Message:          message,
	}, nil
}

func (r *ProkerV3Repository) CreateJadwalProkerV3(ctx context.Context, data *domain.JadwalProkerV3) error {
	if r.db == nil {
		return errors.New("Database tidak tersedia")
	}
	return r.db.WithContext(ctx).Create(data).Error
}

func (r *ProkerV3Repository) UpdateJadwalProkerV3(ctx context.Context, uid uuid.UUID, data *domain.JadwalProkerV3) error {
	if r.db == nil {
		return errors.New("Database tidak tersedia")
	}
	var existing domain.JadwalProkerV3
	if err := r.db.WithContext(ctx).Where("uuid = ?", uid).First(&existing).Error; err != nil {
		return err
	}
	data.ID = existing.ID
	data.UUID = existing.UUID
	data.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(data).Error
}

func (r *ProkerV3Repository) DeleteJadwalProkerV3(ctx context.Context, uid uuid.UUID) error {
	if r.db == nil {
		return errors.New("Database tidak tersedia")
	}
	return r.db.WithContext(ctx).Where("uuid = ?", uid).Delete(&domain.JadwalProkerV3{}).Error
}

// =======================================================
// Pagu Anggaran
// =======================================================

func (r *ProkerV3Repository) GetAllPaguAnggaran(ctx context.Context, tahun string, idFu string) ([]domain.PaguAnggaran, error) {
	if r.db == nil {
		return nil, errors.New("Database tidak tersedia")
	}
	var list []domain.PaguAnggaran
	query := r.db.WithContext(ctx).Table("pagu_anggaran")
	if tahun != "" {
		query = query.Where("tahun = ?", tahun)
	}
	if idFu != "" && idFu != "0" {
		query = query.Where("id_fakultas_unit = ?", idFu)
	}
	err := query.Order("id DESC").Find(&list).Error
	return list, err
}

func (r *ProkerV3Repository) SavePaguAnggaran(ctx context.Context, req domain.SavePaguAnggaranRequest) (*domain.PaguAnggaran, error) {
	if r.db == nil {
		return nil, errors.New("Database tidak tersedia")
	}
	if req.Tahun == "" {
		return nil, errors.New("Tahun wajib diisi")
	}

	nom := float64(0)
	if req.Pagu != nil {
		nom = *req.Pagu
	}

	fakultasUnit := req.FakultasUnit
	if req.IdFakultasUnit != nil && *req.IdFakultasUnit > 0 && fakultasUnit == "" {
		var fuName string
		_ = r.db.WithContext(ctx).Table("v_fakultas_unit").
			Select("nama_fak_prod_unit").
			Where("id = ?", *req.IdFakultasUnit).
			Scan(&fuName).Error
		fakultasUnit = fuName
	}

	var item domain.PaguAnggaran
	if req.ID > 0 {
		if err := r.db.WithContext(ctx).Where("id = ?", req.ID).First(&item).Error; err != nil {
			return nil, err
		}
		item.IdFakultasUnit = req.IdFakultasUnit
		item.FakultasUnit = fakultasUnit
		item.Tahun = req.Tahun
		item.Pagu = nom
		item.UpdatedAt = time.Now()
		if err := r.db.WithContext(ctx).Save(&item).Error; err != nil {
			return nil, err
		}
	} else {
		var existing domain.PaguAnggaran
		checkQuery := r.db.WithContext(ctx).Where("tahun = ?", req.Tahun)
		if req.IdFakultasUnit != nil {
			checkQuery = checkQuery.Where("id_fakultas_unit = ?", *req.IdFakultasUnit)
		} else if fakultasUnit != "" {
			checkQuery = checkQuery.Where("fakultas_unit = ?", fakultasUnit)
		}
		if err := checkQuery.First(&existing).Error; err == nil {
			existing.Pagu = nom
			existing.UpdatedAt = time.Now()
			if fakultasUnit != "" {
				existing.FakultasUnit = fakultasUnit
			}
			if req.IdFakultasUnit != nil {
				existing.IdFakultasUnit = req.IdFakultasUnit
			}
			if err := r.db.WithContext(ctx).Save(&existing).Error; err != nil {
				return nil, err
			}
			item = existing
		} else {
			item = domain.PaguAnggaran{
				IdFakultasUnit: req.IdFakultasUnit,
				FakultasUnit:   fakultasUnit,
				Tahun:          req.Tahun,
				Pagu:           nom,
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			}
			if err := r.db.WithContext(ctx).Create(&item).Error; err != nil {
				return nil, err
			}
		}
	}
	return &item, nil
}

func (r *ProkerV3Repository) DeletePaguAnggaran(ctx context.Context, id string) error {
	if r.db == nil {
		return errors.New("Database tidak tersedia")
	}
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.PaguAnggaran{}).Error
}

func (r *ProkerV3Repository) GetUnitPaguSummary(ctx context.Context, tahun string, reqFu, userFu string) (*domain.UnitPaguSummaryResult, error) {
	if r.db == nil {
		return nil, errors.New("Database tidak tersedia")
	}
	if tahun == "" {
		tahun = "2027"
	}
	if reqFu == "" || reqFu == "0" || reqFu == "all" {
		reqFu = userFu
	}

	var fuId int
	var fuName string

	if reqFu != "" {
		var fuRow struct {
			ID              int    `gorm:"column:id"`
			NamaFakProdUnit string `gorm:"column:nama_fak_prod_unit"`
		}
		_ = r.db.WithContext(ctx).Table("v_fakultas_unit").
			Where("id = ? OR nama_fak_prod_unit = ?", reqFu, reqFu).
			First(&fuRow).Error

		if fuRow.ID > 0 {
			fuId = fuRow.ID
			fuName = fuRow.NamaFakProdUnit
		} else if parsed, err := strconv.Atoi(reqFu); err == nil {
			fuId = parsed
		}
	}

	if fuId == 0 && userFu != "" {
		if parsed, err := strconv.Atoi(userFu); err == nil {
			fuId = parsed
		}
	}

	var paguRow struct {
		Pagu float64
	}
	errPagu := r.db.WithContext(ctx).Table("pagu_anggaran").
		Select("pagu").
		Where("tahun = ? AND (id_fakultas_unit = ? OR (fakultas_unit IS NOT NULL AND fakultas_unit != '' AND (fakultas_unit = ? OR fakultas_unit = ?)))", tahun, fuId, strconv.Itoa(fuId), fuName).
		Order("id DESC").
		Limit(1).
		Scan(&paguRow).Error

	hasPagu := errPagu == nil && paguRow.Pagu > 0
	paguVal := float64(0)
	if hasPagu {
		paguVal = paguRow.Pagu
	}

	var prokers []struct {
		ID             uint    `gorm:"column:id"`
		TotAnggaran    float64 `gorm:"column:tot_anggaran"`
		TotRealisasi   float64 `gorm:"column:tot_realisasi"`
		CountRealisasi int     `gorm:"column:count_realisasi"`
	}

	_ = r.db.WithContext(ctx).Table("matrix_proker mp").
		Select(`
			mp.id,
			COALESCE((SELECT SUM(COALESCE(r.total, r.qty_satuan * r.qty_paket * r.sbu, 0)) FROM aktivitasv2_rab r JOIN aktivitasv2 a ON a.id = r.id_aktivitasv2 WHERE a.id_matrik_proker = mp.id AND COALESCE(a.is_batal, 0) = 0), 0) as tot_anggaran,
			COALESCE((SELECT SUM(CASE WHEN a.realisasi IS NOT NULL THEN a.realisasi ELSE COALESCE((SELECT SUM(COALESCE(r2.total, r2.qty_satuan * r2.qty_paket * r2.sbu, 0)) FROM aktivitasv2_rab r2 WHERE r2.id_aktivitasv2 = a.id), 0) END) FROM aktivitasv2 a WHERE a.id_matrik_proker = mp.id AND COALESCE(a.is_batal, 0) = 0), 0) as tot_realisasi,
			COALESCE((SELECT COUNT(*) FROM aktivitasv2 a WHERE a.id_matrik_proker = mp.id AND a.realisasi IS NOT NULL AND COALESCE(a.is_batal, 0) = 0), 0) as count_realisasi
		`).
		Joins("LEFT JOIN mata_program p ON p.id = mp.id_mata_program").
		Joins("LEFT JOIN master_tahun mt ON mt.id = p.id_master_tahun").
		Where("mp.deleted_at IS NULL AND (mp.id_fakultas_unit = ? OR (SELECT fu2.nama_fak_prod_unit FROM v_fakultas_unit fu2 WHERE fu2.id = mp.id_fakultas_unit) = ?) AND COALESCE(mt.tahun, '') = ? AND COALESCE(mp.beban_yayasan, 0) = 0 AND COALESCE(mp.is_batal, 0) = 0", fuId, fuName, tahun).
		Scan(&prokers).Error

	totalAnggaran := float64(0)
	totalRealisasi := float64(0)
	hasAnyRealisasi := false

	for _, p := range prokers {
		totalAnggaran += p.TotAnggaran
		totalRealisasi += p.TotRealisasi
		if p.CountRealisasi > 0 {
			hasAnyRealisasi = true
		}
	}

	hasRealDiff := hasAnyRealisasi && (totalRealisasi != totalAnggaran)
	effectiveTotal := totalRealisasi
	if !hasAnyRealisasi {
		effectiveTotal = totalAnggaran
	}
	isMelebihiPagu := hasPagu && (effectiveTotal > paguVal)
	sisaPagu := float64(0)
	if hasPagu {
		sisaPagu = paguVal - effectiveTotal
	}

	return &domain.UnitPaguSummaryResult{
		HasPagu:          hasPagu,
		Pagu:             paguVal,
		TotalAnggaran:    totalAnggaran,
		TotalRealisasi:   totalRealisasi,
		HasRealisasi:     hasAnyRealisasi,
		HasRealDiff:      hasRealDiff,
		IsMelebihiPagu:   isMelebihiPagu,
		SisaPagu:         sisaPagu,
		NamaFakultasUnit: fuName,
		IdFakultasUnit:   fuId,
		Tahun:            tahun,
	}, nil
}

// =======================================================
// Dashboard Monitoring & Sankey Data
// =======================================================

func (r *ProkerV3Repository) GetDashboardProkerMonitoring(ctx context.Context, filter domain.DashboardMonitoringFilter) (*domain.DashboardMonitoringResult, error) {
	if r.db == nil {
		return nil, errors.New("Database tidak tersedia")
	}

	var userRow struct {
		ID           uint   `gorm:"column:id"`
		UUID         string `gorm:"column:uuid"`
		NidnUsername string `gorm:"column:nidn_username"`
		Level        string `gorm:"column:level"`
		Name         string `gorm:"column:name"`
		FakultasUnit string `gorm:"column:fakultas_unit"`
	}

	if filter.SID != "" {
		_ = r.db.WithContext(ctx).Table("usersv3").Where("uuid = ?", filter.SID).First(&userRow).Error
	}
	if userRow.ID == 0 && filter.UserFU != "" {
		_ = r.db.WithContext(ctx).Table("usersv3").Where("fakultas_unit = ? AND level = ?", filter.UserFU, filter.Level).First(&userRow).Error
	}

	userName := strings.TrimSpace(userRow.NidnUsername)
	userLevel := strings.ToLower(strings.TrimSpace(userRow.Level))
	if filter.Level != "" {
		userLevel = strings.ToLower(strings.TrimSpace(filter.Level))
	}
	fuID := strings.TrimSpace(userRow.FakultasUnit)
	if fuID == "" {
		fuID = strings.TrimSpace(filter.UserFU)
	}

	isAdmin := userLevel == "admin" || strings.EqualFold(userName, "admin")
	isWarek := strings.Contains(strings.ToLower(userName), "wakil rektor") ||
		strings.Contains(strings.ToLower(userName), "warek") ||
		strings.Contains(userLevel, "wakil rektor") ||
		strings.Contains(userLevel, "warek") ||
		fuID == "65" || fuID == "66" || fuID == "67" || fuID == "122"

	var fuInfo struct {
		ID              int    `gorm:"column:id"`
		NamaFakProdUnit string `gorm:"column:nama_fak_prod_unit"`
		Type            string `gorm:"column:type"`
		KodeFakultas    string `gorm:"column:kode_fakultas"`
	}
	if fuID != "" {
		_ = r.db.WithContext(ctx).Table("v_fakultas_unit").
			Where("id = ? OR nama_fak_prod_unit = ?", fuID, fuID).
			First(&fuInfo).Error
	}

	isFakultas := fuInfo.Type == "fakultas" || userLevel == "fakultas"

	if !isAdmin && !isWarek && !isFakultas {
		return nil, errors.New("Fitur monitoring ini hanya dapat diakses oleh Wakil Rektor 1-4 dan Tingkat Fakultas")
	}

	userRoleType := "warek"
	if isFakultas && !isWarek && !isAdmin {
		userRoleType = "fakultas"
	} else if isAdmin {
		userRoleType = "admin"
	}

	var unitOptions []domain.DashboardMonitoringUnitOption
	if isAdmin || isWarek {
		unitOptions = append(unitOptions, domain.DashboardMonitoringUnitOption{
			ID:   "all",
			Nama: "Semua Fakultas & Unit",
			Type: "all",
		})
		var allUnits []struct {
			ID              int    `gorm:"column:id"`
			NamaFakProdUnit string `gorm:"column:nama_fak_prod_unit"`
			Type            string `gorm:"column:type"`
			Jenjang         string `gorm:"column:jenjang"`
		}
		_ = r.db.WithContext(ctx).Table("v_fakultas_unit").
			Order("FIELD(type, 'fakultas', 'prodi', 'unit'), nama_fak_prod_unit ASC").
			Scan(&allUnits).Error

		for _, u := range allUnits {
			label := u.NamaFakProdUnit
			if u.Type == "prodi" && u.Jenjang != "" {
				label = fmt.Sprintf("%s (%s)", u.NamaFakProdUnit, strings.ToUpper(u.Jenjang))
			}
			unitOptions = append(unitOptions, domain.DashboardMonitoringUnitOption{
				ID:   u.ID,
				Nama: label,
				Type: u.Type,
			})
		}
	} else if isFakultas {
		facultyName := fuInfo.NamaFakProdUnit
		if facultyName == "" {
			facultyName = "Fakultas"
		}
		unitOptions = append(unitOptions, domain.DashboardMonitoringUnitOption{
			ID:   "all",
			Nama: fmt.Sprintf("Semua Prodi (%s)", facultyName),
			Type: "all",
		})
		var facUnits []struct {
			ID              int    `gorm:"column:id"`
			NamaFakProdUnit string `gorm:"column:nama_fak_prod_unit"`
			Type            string `gorm:"column:type"`
			Jenjang         string `gorm:"column:jenjang"`
		}
		_ = r.db.WithContext(ctx).Table("v_fakultas_unit").
			Where("kode_fakultas = ?", fuInfo.KodeFakultas).
			Order("FIELD(type, 'fakultas', 'prodi'), nama_fak_prod_unit ASC").
			Scan(&facUnits).Error

		for _, u := range facUnits {
			label := u.NamaFakProdUnit
			if u.Type == "prodi" && u.Jenjang != "" {
				label = fmt.Sprintf("%s (%s)", u.NamaFakProdUnit, strings.ToUpper(u.Jenjang))
			}
			unitOptions = append(unitOptions, domain.DashboardMonitoringUnitOption{
				ID:   u.ID,
				Nama: label,
				Type: u.Type,
			})
		}
	}

	tahun := strings.TrimSpace(filter.Tahun)
	if tahun == "" {
		var activeTahun string
		_ = r.db.WithContext(ctx).Table("master_tahun").
			Where("status = 'active'").
			Order("tahun DESC").
			Pluck("tahun", &activeTahun).Error
		if activeTahun != "" {
			tahun = activeTahun
		} else {
			tahun = "2027"
		}
	}
	selectedUnit := strings.TrimSpace(filter.FilterUnit)
	search := strings.TrimSpace(filter.Search)

	query := r.db.WithContext(ctx).Table("matrix_proker mp").
		Select(`
			mp.id as id,
			mp.uuid as uuid,
			mp.id_fakultas_unit as id_fakultas_unit,
			COALESCE(fu.nama_fak_prod_unit, '') as nama_fakultas_unit,
			COALESCE(fu.type, '') as type_fakultas_unit,
			COALESCE(fu.jenjang, '') as jenjang,
			COALESCE(prog.mata_program, '') as nama_program,
			COALESCE(mt.tahun, '') as tahun,
			COALESCE(mir.indikator, '') as indikator_kinerja,
			COALESCE((SELECT mps.target_unpak FROM matrix_proker_scope mps WHERE mps.id_matrix_proker = mp.id ORDER BY mps.id ASC LIMIT 1), '') as target_unpak,
			COALESCE((SELECT mps.target_unit FROM matrix_proker_scope mps WHERE mps.id_matrix_proker = mp.id ORDER BY mps.id ASC LIMIT 1), '') as target_unit,
			COALESCE((SELECT mps.output FROM matrix_proker_scope mps WHERE mps.id_matrix_proker = mp.id ORDER BY mps.id ASC LIMIT 1), '') as output,
			COALESCE((SELECT mps.outcome FROM matrix_proker_scope mps WHERE mps.id_matrix_proker = mp.id ORDER BY mps.id ASC LIMIT 1), '') as outcome,
			COALESCE(NULLIF(mp.pic, ''), '') as pic,
			COALESCE(NULLIF(mp.status, ''), 'draft') as status,
			COALESCE(NULLIF(mp.status_verifikator, ''), 'draft') as status_verifikator,
			COALESCE(NULLIF(mp.status_anggaran, ''), 'menunggu_verifikator') as status_anggaran,
			COALESCE(mp.beban_yayasan, 0) as beban_yayasan,
			COALESCE(mp.is_batal, 0) as is_batal,
			COALESCE(mp.catatan_batal, '') as catatan_batal
		`).
		Joins("LEFT JOIN mata_program prog ON mp.id_mata_program = prog.id").
		Joins("LEFT JOIN master_tahun mt ON prog.id_master_tahun = mt.id").
		Joins("LEFT JOIN master_indikator_renstra mir ON mp.id_indikator = mir.id").
		Joins("LEFT JOIN v_fakultas_unit fu ON mp.id_fakultas_unit = fu.id").
		Where("mp.deleted_at IS NULL AND COALESCE(mt.tahun, '') = ?", tahun)

	if selectedUnit != "" && selectedUnit != "all" {
		query = query.Where("mp.id_fakultas_unit = ?", selectedUnit)
	} else if isFakultas && !isAdmin && !isWarek && fuInfo.KodeFakultas != "" {
		query = query.Where("mp.id_fakultas_unit IN (SELECT id FROM v_fakultas_unit WHERE kode_fakultas = ?)", fuInfo.KodeFakultas)
	}

	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where("(prog.mata_program LIKE ? OR fu.nama_fak_prod_unit LIKE ? OR mir.indikator LIKE ? OR mp.pic LIKE ?)", searchPattern, searchPattern, searchPattern, searchPattern)
	}

	var prokerRows []struct {
		Id                uint   `gorm:"column:id"`
		UUID              string `gorm:"column:uuid"`
		IdFakultasUnit    *int   `gorm:"column:id_fakultas_unit"`
		NamaFakultasUnit  string `gorm:"column:nama_fakultas_unit"`
		TypeFakultasUnit  string `gorm:"column:type_fakultas_unit"`
		Jenjang           string `gorm:"column:jenjang"`
		NamaProgram       string `gorm:"column:nama_program"`
		Tahun             string `gorm:"column:tahun"`
		IndikatorKinerja  string `gorm:"column:indikator_kinerja"`
		TargetUnpak       string `gorm:"column:target_unpak"`
		TargetUnit        string `gorm:"column:target_unit"`
		Output            string `gorm:"column:output"`
		Outcome           string `gorm:"column:outcome"`
		PIC               string `gorm:"column:pic"`
		Status            string `gorm:"column:status"`
		StatusVerifikator string `gorm:"column:status_verifikator"`
		StatusAnggaran    string `gorm:"column:status_anggaran"`
		BebanYayasan      int    `gorm:"column:beban_yayasan"`
		IsBatal           int    `gorm:"column:is_batal"`
		CatatanBatal      string `gorm:"column:catatan_batal"`
	}

	if err := query.Order("mp.id DESC").Scan(&prokerRows).Error; err != nil {
		return nil, err
	}

	var prokerIDs []uint
	for _, p := range prokerRows {
		prokerIDs = append(prokerIDs, p.Id)
	}

	var rawAktivitas []struct {
		Id                    uint     `gorm:"column:id"`
		UUID                  string   `gorm:"column:uuid"`
		IdMatrikProker        uint     `gorm:"column:id_matrik_proker"`
		Aktivitas             string   `gorm:"column:aktivitas"`
		PIC                   string   `gorm:"column:pic"`
		Capaian               string   `gorm:"column:capaian"`
		Realisasi             *float64 `gorm:"column:realisasi"`
		CatatanRevisiAnggaran string   `gorm:"column:catatan_revisi_anggaran"`
		Pencairan             *float64 `gorm:"column:pencairan"`
		Status                string   `gorm:"column:status"`
		StatusAnggaran        string   `gorm:"column:status_anggaran"`
		StatusUpml            string   `gorm:"column:status_upml"`
		StatusUpmf            string   `gorm:"column:status_upmf"`
		StatusUpms            string   `gorm:"column:status_upms"`
		StatusLpm             string   `gorm:"column:status_lpm"`
		IsBatal               int      `gorm:"column:is_batal"`
		CatatanBatal          string   `gorm:"column:catatan_batal"`
	}

	var aktivitasIDs []uint
	if len(prokerIDs) > 0 {
		_ = r.db.WithContext(ctx).Table("aktivitasv2").
			Select("id, uuid, id_matrik_proker, aktivitas, pic, capaian, realisasi, COALESCE(catatan_revisi_anggaran, '') as catatan_revisi_anggaran, pencairan, status, status_anggaran, status_upml, status_upmf, status_upms, status_lpm, is_batal, COALESCE(catatan_batal, '') as catatan_batal").
			Where("id_matrik_proker IN (?)", prokerIDs).
			Order("id ASC").
			Scan(&rawAktivitas).Error

		for _, a := range rawAktivitas {
			aktivitasIDs = append(aktivitasIDs, a.Id)
		}
	}

	var rawRAB []struct {
		Id            uint    `gorm:"column:id"`
		UUID          string  `gorm:"column:uuid"`
		IdAktivitasv2 uint    `gorm:"column:id_aktivitasv2"`
		NoAkun        string  `gorm:"column:no_akun"`
		NamaAkun      string  `gorm:"column:nama_akun"`
		Kategori      string  `gorm:"column:kategori"`
		Rincian       string  `gorm:"column:rincian"`
		Keterangan    string  `gorm:"column:keterangan"`
		QtySatuan     float64 `gorm:"column:qty_satuan"`
		QtyPaket      float64 `gorm:"column:qty_paket"`
		Sbu           float64 `gorm:"column:sbu"`
		Total         float64 `gorm:"column:total"`
	}

	if len(aktivitasIDs) > 0 {
		_ = r.db.WithContext(ctx).Table("aktivitasv2_rab").
			Where("id_aktivitasv2 IN (?)", aktivitasIDs).
			Order("id ASC").
			Scan(&rawRAB).Error
	}

	var rawTimeline []struct {
		Id               uint       `gorm:"column:id"`
		UUID             string     `gorm:"column:uuid"`
		IdAktivitasv2    uint       `gorm:"column:id_aktivitasv2"`
		NoAktivitas      int        `gorm:"column:no_aktivitas"`
		NamaAktivitas    string     `gorm:"column:nama_aktivitas"`
		TanggalRkAwal    *time.Time `gorm:"column:tanggal_rk_awal"`
		TanggalRkAkhir   *time.Time `gorm:"column:tanggal_rk_akhir"`
		TanggalRAwal     *time.Time `gorm:"column:tanggal_r_awal"`
		TanggalRAkhir    *time.Time `gorm:"column:tanggal_r_akhir"`
		Capaian          string     `gorm:"column:capaian"`
		Masalah          string     `gorm:"column:masalah"`
		PemecahanMasalah string     `gorm:"column:pemecahan_masalah"`
		Status           string     `gorm:"column:status"`
	}

	if len(aktivitasIDs) > 0 {
		_ = r.db.WithContext(ctx).Table("aktivitasv2_timeline").
			Where("id_aktivitasv2 IN (?)", aktivitasIDs).
			Order("no_aktivitas ASC, id ASC").
			Scan(&rawTimeline).Error
	}

	timelineMap := make(map[uint][]domain.DashboardMonitoringTimelineItem)
	for _, t := range rawTimeline {
		var rkAwalStr, rkAkhirStr string
		if t.TanggalRkAwal != nil && !t.TanggalRkAwal.IsZero() {
			rkAwalStr = t.TanggalRkAwal.Format("2006-01-02")
		}
		if t.TanggalRkAkhir != nil && !t.TanggalRkAkhir.IsZero() {
			rkAkhirStr = t.TanggalRkAkhir.Format("2006-01-02")
		}
		var rAwalStr, rAkhirStr *string
		if t.TanggalRAwal != nil && !t.TanggalRAwal.IsZero() {
			s := t.TanggalRAwal.Format("2006-01-02")
			rAwalStr = &s
		}
		if t.TanggalRAkhir != nil && !t.TanggalRAkhir.IsZero() {
			s := t.TanggalRAkhir.Format("2006-01-02")
			rAkhirStr = &s
		}

		item := domain.DashboardMonitoringTimelineItem{
			ID:               t.Id,
			UUID:             t.UUID,
			IdAktivitasv2:    t.IdAktivitasv2,
			NoAktivitas:      t.NoAktivitas,
			NamaAktivitas:    t.NamaAktivitas,
			TanggalRkAwal:    rkAwalStr,
			TanggalRkAkhir:   rkAkhirStr,
			TanggalRAwal:     rAwalStr,
			TanggalRAkhir:    rAkhirStr,
			Capaian:          t.Capaian,
			Masalah:          t.Masalah,
			PemecahanMasalah: t.PemecahanMasalah,
			Status:           t.Status,
		}
		timelineMap[t.IdAktivitasv2] = append(timelineMap[t.IdAktivitasv2], item)
	}

	rabMap := make(map[uint][]domain.DashboardMonitoringRabItem)
	rabTotalMap := make(map[uint]float64)
	for _, rItem := range rawRAB {
		tot := rItem.Total
		if tot == 0 && rItem.QtySatuan > 0 && rItem.QtyPaket > 0 && rItem.Sbu > 0 {
			tot = rItem.QtySatuan * rItem.QtyPaket * rItem.Sbu
		}
		item := domain.DashboardMonitoringRabItem{
			ID:         rItem.Id,
			UUID:       rItem.UUID,
			NoAkun:     rItem.NoAkun,
			NamaAkun:   rItem.NamaAkun,
			Kategori:   rItem.Kategori,
			Rincian:    rItem.Rincian,
			Keterangan: rItem.Keterangan,
			QtySatuan:  rItem.QtySatuan,
			QtyPaket:   rItem.QtyPaket,
			Sbu:        rItem.Sbu,
			Total:      tot,
		}
		rabMap[rItem.IdAktivitasv2] = append(rabMap[rItem.IdAktivitasv2], item)
		rabTotalMap[rItem.IdAktivitasv2] += tot
	}

	aktMap := make(map[uint][]domain.DashboardMonitoringAktivitas)
	prokerAnggaranMap := make(map[uint]float64)
	prokerAnggaranPengajuanMap := make(map[uint]float64)
	prokerIsRevisedMap := make(map[uint]bool)
	prokerPencairanMap := make(map[uint]float64)

	for _, a := range rawAktivitas {
		rabItems := rabMap[a.Id]
		if rabItems == nil {
			rabItems = []domain.DashboardMonitoringRabItem{}
		}
		totRab := rabTotalMap[a.Id]

		targetBanding := totRab
		if a.Realisasi != nil && *a.Realisasi >= 0 {
			targetBanding = *a.Realisasi
		}

		statusPencairan := "belum_dicairkan"
		if a.Pencairan != nil && *a.Pencairan > 0 {
			pencairanVal := *a.Pencairan
			if pencairanVal == targetBanding {
				statusPencairan = "cair"
			} else if pencairanVal < targetBanding {
				statusPencairan = "tidak_penuh"
			} else {
				statusPencairan = "cair_kelebihan"
			}
			if a.IsBatal == 0 {
				prokerPencairanMap[a.IdMatrikProker] += pencairanVal
			}
		}

		actApprovedBudget := totRab
		if a.Realisasi != nil && *a.Realisasi >= 0 {
			actApprovedBudget = *a.Realisasi
		}

		if a.IsBatal == 0 {
			prokerAnggaranMap[a.IdMatrikProker] += actApprovedBudget
			prokerAnggaranPengajuanMap[a.IdMatrikProker] += totRab
			if totRab != actApprovedBudget {
				prokerIsRevisedMap[a.IdMatrikProker] = true
			}
		}

		timelineItems := timelineMap[a.Id]
		if timelineItems == nil {
			timelineItems = []domain.DashboardMonitoringTimelineItem{}
		}

		aktItem := domain.DashboardMonitoringAktivitas{
			ID:                    a.Id,
			UUID:                  a.UUID,
			Aktivitas:             a.Aktivitas,
			PIC:                   a.PIC,
			Capaian:               a.Capaian,
			Realisasi:             a.Realisasi,
			CatatanRevisiAnggaran: a.CatatanRevisiAnggaran,
			Pencairan:             a.Pencairan,
			Status:                a.Status,
			StatusAnggaran:        a.StatusAnggaran,
			StatusUpml:            a.StatusUpml,
			StatusUpmf:            a.StatusUpmf,
			StatusUpms:            a.StatusUpms,
			StatusLpm:             a.StatusLpm,
			IsBatal:               a.IsBatal == 1,
			CatatanBatal:          a.CatatanBatal,
			TotalRab:              totRab,
			StatusPencairan:       statusPencairan,
			Rab:                   rabItems,
			Timelines:             timelineItems,
		}
		aktMap[a.IdMatrikProker] = append(aktMap[a.IdMatrikProker], aktItem)
	}

	var prokerList []domain.DashboardMonitoringProker
	var globalTotalProker int
	var globalTotalAktivitas int
	var globalTotalAnggaran float64
	var globalTotalPencairan float64

	for _, p := range prokerRows {
		activities := aktMap[p.Id]
		if activities == nil {
			activities = []domain.DashboardMonitoringAktivitas{}
		}
		totAnggaran := prokerAnggaranMap[p.Id]
		totPengajuan := prokerAnggaranPengajuanMap[p.Id]
		isRevised := prokerIsRevisedMap[p.Id]
		totPencairan := prokerPencairanMap[p.Id]

		unitLabel := p.NamaFakultasUnit
		if p.TypeFakultasUnit == "prodi" && p.Jenjang != "" {
			unitLabel = fmt.Sprintf("%s (%s)", p.NamaFakultasUnit, strings.ToUpper(p.Jenjang))
		}

		prokerItem := domain.DashboardMonitoringProker{
			ID:                     p.Id,
			UUID:                   p.UUID,
			IdFakultasUnit:         p.IdFakultasUnit,
			NamaFakultasUnit:       unitLabel,
			TypeFakultasUnit:       p.TypeFakultasUnit,
			Jenjang:                p.Jenjang,
			NamaProgram:            p.NamaProgram,
			Tahun:                  p.Tahun,
			IndikatorKinerja:       p.IndikatorKinerja,
			TargetUnpak:            p.TargetUnpak,
			TargetUnit:             p.TargetUnit,
			Output:                 p.Output,
			Outcome:                p.Outcome,
			PIC:                    p.PIC,
			Status:                 p.Status,
			StatusVerifikator:      p.StatusVerifikator,
			StatusAnggaran:         p.StatusAnggaran,
			BebanYayasan:           p.BebanYayasan == 1,
			IsBatal:                p.IsBatal == 1,
			CatatanBatal:           p.CatatanBatal,
			TotalAktivitas:         len(activities),
			TotalAnggaran:          totAnggaran,
			TotalAnggaranPengajuan: totPengajuan,
			IsAnggaranRevised:      isRevised,
			TotalPencairan:         totPencairan,
			Aktivitas:              activities,
		}

		prokerList = append(prokerList, prokerItem)
		globalTotalProker++
		globalTotalAktivitas += len(activities)
		if p.IsBatal == 0 {
			globalTotalAnggaran += totAnggaran
			globalTotalPencairan += totPencairan
		}
	}

	if prokerList == nil {
		prokerList = []domain.DashboardMonitoringProker{}
	}
	if unitOptions == nil {
		unitOptions = []domain.DashboardMonitoringUnitOption{}
	}

	return &domain.DashboardMonitoringResult{
		Meta: map[string]interface{}{
			"allowed":      true,
			"role_type":    userRoleType,
			"user_name":    userName,
			"unit_options": unitOptions,
		},
		Summary: map[string]interface{}{
			"total_proker":    globalTotalProker,
			"total_aktivitas": globalTotalAktivitas,
			"total_anggaran":  globalTotalAnggaran,
			"total_pencairan": globalTotalPencairan,
		},
		Data: prokerList,
	}, nil
}

func (r *ProkerV3Repository) GetSankeyData(ctx context.Context, tahun string) (*domain.SankeyDataResult, error) {
	if r.db == nil {
		return nil, errors.New("Database tidak tersedia")
	}

	if tahun == "" {
		_ = r.db.WithContext(ctx).Table("master_tahun").
			Where("status = 'active'").
			Order("tahun DESC").
			Pluck("tahun", &tahun).Error
		if tahun == "" {
			tahun = "2027"
		}
	}

	var masterIkus []map[string]interface{}
	_ = r.db.WithContext(ctx).Table("master_iku").Order("id asc").Scan(&masterIkus).Error

	var allUnits []map[string]interface{}
	_ = r.db.WithContext(ctx).Table("v_fakultas_unit").
		Order("FIELD(type, 'fakultas', 'prodi', 'unit'), nama_fak_prod_unit ASC").
		Scan(&allUnits).Error

	var prokerRows []struct {
		ID               uint   `gorm:"column:id"`
		UUID             string `gorm:"column:uuid"`
		IdFakultasUnit   *int   `gorm:"column:id_fakultas_unit"`
		IdIndikator      *int   `gorm:"column:id_indikator"`
		NamaFakultasUnit string `gorm:"column:nama_fakultas_unit"`
		TypeFakultasUnit string `gorm:"column:type_fakultas_unit"`
		Jenjang          string `gorm:"column:jenjang"`
		NamaProgram      string `gorm:"column:nama_program"`
		Tahun            string `gorm:"column:tahun"`
		IndikatorKinerja string `gorm:"column:indikator_kinerja"`
		PIC              string `gorm:"column:pic"`
		Status           string `gorm:"column:status"`
		IsBatal          int    `gorm:"column:is_batal"`
	}

	_ = r.db.WithContext(ctx).Table("matrix_proker mp").
		Select(`
			mp.id, mp.uuid, mp.id_fakultas_unit, mp.id_indikator,
			COALESCE(fu.nama_fak_prod_unit, '') as nama_fakultas_unit,
			COALESCE(fu.type, '') as type_fakultas_unit,
			COALESCE(fu.jenjang, '') as jenjang,
			COALESCE(prog.mata_program, '') as nama_program,
			COALESCE(mt.tahun, '') as tahun,
			COALESCE(mir.indikator, '') as indikator_kinerja,
			COALESCE(NULLIF(mp.pic, ''), '') as pic,
			COALESCE(NULLIF(mp.status, ''), 'draft') as status,
			COALESCE(mp.is_batal, 0) as is_batal
		`).
		Joins("LEFT JOIN mata_program prog ON mp.id_mata_program = prog.id").
		Joins("LEFT JOIN master_tahun mt ON prog.id_master_tahun = mt.id").
		Joins("LEFT JOIN master_indikator_renstra mir ON mp.id_indikator = mir.id").
		Joins("LEFT JOIN v_fakultas_unit fu ON mp.id_fakultas_unit = fu.id").
		Where("mp.deleted_at IS NULL AND COALESCE(mt.tahun, '') = ?", tahun).
		Order("mp.id DESC").
		Scan(&prokerRows).Error

	var prokerIDs []uint
	for _, p := range prokerRows {
		prokerIDs = append(prokerIDs, p.ID)
	}

	var rawAktivitas []struct {
		ID                    uint     `gorm:"column:id"`
		UUID                  string   `gorm:"column:uuid"`
		IdMatrikProker        uint     `gorm:"column:id_matrik_proker"`
		Aktivitas             string   `gorm:"column:aktivitas"`
		PIC                   string   `gorm:"column:pic"`
		Realisasi             *float64 `gorm:"column:realisasi"`
		Pencairan             *float64 `gorm:"column:pencairan"`
		CatatanRevisiAnggaran string   `gorm:"column:catatan_revisi_anggaran"`
		Status                string   `gorm:"column:status"`
		IsBatal               int      `gorm:"column:is_batal"`
		CatatanBatal          string   `gorm:"column:catatan_batal"`
	}
	var aktivitasIDs []uint
	if len(prokerIDs) > 0 {
		_ = r.db.WithContext(ctx).Table("aktivitasv2").
			Select("id, uuid, id_matrik_proker, aktivitas, pic, realisasi, pencairan, COALESCE(catatan_revisi_anggaran, '') as catatan_revisi_anggaran, status, is_batal, COALESCE(catatan_batal, '') as catatan_batal").
			Where("id_matrik_proker IN (?)", prokerIDs).
			Order("id ASC").
			Scan(&rawAktivitas).Error

		for _, a := range rawAktivitas {
			aktivitasIDs = append(aktivitasIDs, a.ID)
		}
	}

	var rawRAB []struct {
		ID            uint    `gorm:"column:id"`
		UUID          string  `gorm:"column:uuid"`
		IdAktivitasV2 uint    `gorm:"column:id_aktivitasv2"`
		NoAkun        string  `gorm:"column:no_akun"`
		NamaAkun      string  `gorm:"column:nama_akun"`
		Rincian       string  `gorm:"column:rincian"`
		QtySatuan     float64 `gorm:"column:qty_satuan"`
		QtyPaket      float64 `gorm:"column:qty_paket"`
		Sbu           float64 `gorm:"column:sbu"`
		Total         float64 `gorm:"column:total"`
	}
	if len(aktivitasIDs) > 0 {
		_ = r.db.WithContext(ctx).Table("aktivitasv2_rab").
			Where("id_aktivitasv2 IN (?)", aktivitasIDs).
			Order("id ASC").
			Scan(&rawRAB).Error
	}

	rabMap := make(map[uint][]map[string]interface{})
	for _, rItem := range rawRAB {
		tot := rItem.Total
		if tot == 0 && rItem.QtySatuan > 0 && rItem.QtyPaket > 0 && rItem.Sbu > 0 {
			tot = rItem.QtySatuan * rItem.QtyPaket * rItem.Sbu
		}
		rabMap[rItem.IdAktivitasV2] = append(rabMap[rItem.IdAktivitasV2], map[string]interface{}{
			"id":         rItem.ID,
			"uuid":       rItem.UUID,
			"no_akun":    rItem.NoAkun,
			"nama_akun":  rItem.NamaAkun,
			"rincian":    rItem.Rincian,
			"qty_satuan": rItem.QtySatuan,
			"qty_paket":  rItem.QtyPaket,
			"sbu":        rItem.Sbu,
			"total":      tot,
		})
	}

	aktMap := make(map[uint][]map[string]interface{})
	for _, a := range rawAktivitas {
		rabs := rabMap[a.ID]
		if rabs == nil {
			rabs = []map[string]interface{}{}
		}
		var totRab float64
		for _, rItem := range rabs {
			if t, ok := rItem["total"].(float64); ok {
				totRab += t
			}
		}
		aktMap[a.IdMatrikProker] = append(aktMap[a.IdMatrikProker], map[string]interface{}{
			"id":                      a.ID,
			"uuid":                    a.UUID,
			"aktivitas":               a.Aktivitas,
			"pic":                     a.PIC,
			"realisasi":               a.Realisasi,
			"pencairan":               a.Pencairan,
			"status":                  a.Status,
			"catatan_revisi_anggaran": a.CatatanRevisiAnggaran,
			"catatan_batal":           a.CatatanBatal,
			"is_batal":                a.IsBatal == 1,
			"total_rab":               totRab,
			"rab":                     rabs,
		})
	}

	var monitoringList []map[string]interface{}
	for _, p := range prokerRows {
		acts := aktMap[p.ID]
		if acts == nil {
			acts = []map[string]interface{}{}
		}
		var totAnggaran, totPencairan float64
		for _, act := range acts {
			if rabTot, ok := act["total_rab"].(float64); ok {
				totAnggaran += rabTot
			}
			if pcr, ok := act["pencairan"].(*float64); ok && pcr != nil {
				totPencairan += *pcr
			} else if act["status"] == "terima_anggaran" || act["status"] == "terlaksana" {
				if rabTot, ok := act["total_rab"].(float64); ok {
					totPencairan += rabTot
				}
			}
		}

		monitoringList = append(monitoringList, map[string]interface{}{
			"id":                 p.ID,
			"uuid":               p.UUID,
			"id_fakultas_unit":   p.IdFakultasUnit,
			"id_indikator":       p.IdIndikator,
			"nama_fakultas_unit": p.NamaFakultasUnit,
			"type_fakultas_unit": p.TypeFakultasUnit,
			"jenjang":            p.Jenjang,
			"nama_program":       p.NamaProgram,
			"tahun":              p.Tahun,
			"indikator_kinerja":  p.IndikatorKinerja,
			"pic":                p.PIC,
			"status":             p.Status,
			"is_batal":           p.IsBatal == 1,
			"total_anggaran":     totAnggaran,
			"total_pencairan":    totPencairan,
			"aktivitas":          acts,
		})
	}

	return &domain.SankeyDataResult{
		Tahun:        tahun,
		Monitoring:   monitoringList,
		MasterIkus:   masterIkus,
		FakultasUnit: allUnits,
	}, nil
}
