package infrastructure

import (
	commondomain "UnpakSiamida/common/domain"
	domainRisk "UnpakSiamida/modules/riskassessment/domain"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"strconv"
	"strings"
)

type RiskAssessmentRepository struct {
	db *gorm.DB
}

func NewRiskAssessmentRepository(db *gorm.DB) domainRisk.IRiskAssessmentRepository {
	return &RiskAssessmentRepository{db: db}
}

// ------------------------
// GET BY UUID
// ------------------------
func (r *RiskAssessmentRepository) GetByUuid(ctx context.Context, uid uuid.UUID) (*domainRisk.RiskAssessment, error) {
	var item domainRisk.RiskAssessment
	err := r.db.WithContext(ctx).Table("risk_assessments").Where("uuid = ?", uid).First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// ------------------------
// GET DEFAULT BY UUID
// ------------------------
func (r *RiskAssessmentRepository) GetDefaultByUuid(ctx context.Context, uid uuid.UUID) (*domainRisk.RiskAssessmentDefault, error) {
	var row domainRisk.RiskAssessmentDefault
	err := r.db.WithContext(ctx).
		Table("risk_assessments r").
		Select(`
			r.id as Id,
			r.uuid as UUID,
			r.id_fakultas_unit as IdFakultasUnit,
			CASE 
				WHEN fu.type = 'prodi' AND fu.jenjang IS NOT NULL AND fu.jenjang != '' 
				THEN CONCAT(fu.nama_fak_prod_unit, ' (', UPPER(fu.jenjang), ')') 
				ELSE fu.nama_fak_prod_unit 
			END as NamaFakultasUnit,
			fu.uuid as UuidFakultasUnit,
			r.id_master_tahun as IdMasterTahun,
			CONVERT(COALESCE(mt.tahun, r.tahun) USING utf8mb4) as Tahun,
			mt.uuid as TahunUuid,

			r.id_mata_program as IdMataProgram,
			COALESCE(mp.mata_program, '') as MataProgram,
			mp.uuid as MataProgramUuid,
			COALESCE(NULLIF(r.isu_internal_eksternal, ''), r.penyebab, '') as IsuInternalEksternal,
			COALESCE(NULLIF(r.deskripsi_risiko, ''), r.pernyataan_risiko, '') as DeskripsiRisiko,
			COALESCE(NULLIF(r.kemungkinan_awal, 0), r.kemungkinan, 0) as KemungkinanAwal,
			COALESCE(NULLIF(r.dampak_awal, 0), r.keparahan, 0) as DampakAwal,
			COALESCE(r.tingkat_risiko, (COALESCE(NULLIF(r.kemungkinan_awal, 0), r.kemungkinan, 0) * COALESCE(NULLIF(r.dampak_awal, 0), r.keparahan, 0))) as TotalSkorAwal,
			COALESCE(NULLIF(r.tingkat_risiko_awal, ''), r.level_risiko, '') as TingkatRisikoAwal,
			COALESCE(rl_awal.nama_level, r.level_risiko, '') as NamaTingkatAwal,
			COALESCE(NULLIF(r.tindakan_mitigasi, ''), r.rencana_mitigasi, '') as TindakanMitigasi,
			COALESCE(r.peluang, '') as Peluang,
			COALESCE(r.kenyataan_residu, '') as KenyataanResidu,
			r.kemungkinan_residu as KemungkinanResidu,
			r.dampak_residu as DampakResidu,
			(COALESCE(r.kemungkinan_residu, 0) * COALESCE(r.dampak_residu, 0)) as TotalSkorResidu,
			r.tingkat_risiko_residu as TingkatRisikoResidu,
			COALESCE(rl_res.nama_level, '') as NamaTingkatResidu,
			COALESCE(r.rencana_tindak_lanjut, '') as RencanaTindakLanjut,

			COALESCE(r.kode_risiko, '') as KodeRisiko,
			COALESCE(r.kategori, '') as Kategori,
			COALESCE(r.pernyataan_risiko, r.deskripsi_risiko, '') as PernyataanRisiko,
			COALESCE(r.penyebab, r.isu_internal_eksternal, '') as Penyebab,
			COALESCE(r.dampak, '') as Dampak,
			r.kemungkinan as Kemungkinan,
			r.keparahan as Keparahan,
			r.tingkat_risiko as TingkatRisiko,
			COALESCE(r.level_risiko, r.tingkat_risiko_awal, '') as LevelRisiko,
			COALESCE(r.rencana_mitigasi, r.tindakan_mitigasi, '') as RencanaMitigasi,
			COALESCE(r.pengendalian, '') as Pengendalian,
			COALESCE(r.penanggung_jawab, '') as PenanggungJawab,
			COALESCE(r.target_waktu, '') as TargetWaktu,
			COALESCE(r.status, 'Identifikasi') as Status,
			r.created_at as CreatedAt,
			r.updated_at as UpdatedAt
		`).
		Joins("LEFT JOIN v_fakultas_unit fu ON r.id_fakultas_unit = fu.id").
		Joins("LEFT JOIN master_tahun mt ON r.id_master_tahun = mt.id").
		Joins("LEFT JOIN mata_program mp ON r.id_mata_program = mp.id").
		Joins("LEFT JOIN risk_levels rl_awal ON r.tingkat_risiko_awal = rl_awal.kode").
		Joins("LEFT JOIN risk_levels rl_res ON r.tingkat_risiko_residu = rl_res.kode").
		Where("r.uuid = ?", uid).
		Take(&row).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}

	return &row, nil
}

var allowedSearchColumns = map[string]string{
	"tahun":                  "CONVERT(COALESCE(mt.tahun, r.tahun) USING utf8mb4)",
	"tahunuuid":              "mt.uuid",
	"id_mata_program":        "r.id_mata_program",
	"mataprogramuuid":        "mp.uuid",
	"mataprogram":            "mp.mata_program",
	"mata_program":           "mp.mata_program",
	"isu_internal_eksternal": "COALESCE(NULLIF(r.isu_internal_eksternal, ''), r.penyebab)",
	"deskripsi_risiko":       "COALESCE(NULLIF(r.deskripsi_risiko, ''), r.pernyataan_risiko)",
	"tingkat_risiko_awal":    "COALESCE(NULLIF(r.tingkat_risiko_awal, ''), r.level_risiko)",
	"tingkat_risiko_residu":  "r.tingkat_risiko_residu",
	"level_risiko":           "COALESCE(NULLIF(r.tingkat_risiko_awal, ''), r.level_risiko)",
	"tindakan_mitigasi":      "COALESCE(NULLIF(r.tindakan_mitigasi, ''), r.rencana_mitigasi)",
	"peluang":                "r.peluang",
	"kenyataan_residu":       "r.kenyataan_residu",
	"rencana_tindak_lanjut":  "r.rencana_tindak_lanjut",
	"id_fakultas_unit":       "r.id_fakultas_unit",
	"fakultas_unit":          "CASE WHEN fu.type = 'prodi' AND fu.jenjang IS NOT NULL AND fu.jenjang != '' THEN CONCAT(fu.nama_fak_prod_unit, ' (', UPPER(fu.jenjang), ')') ELSE fu.nama_fak_prod_unit END",
	"fakultasunituuid":       "fu.uuid",
	"kode_risiko":            "r.kode_risiko",
	"kategori":               "r.kategori",
	"status":                 "r.status",
}

// ------------------------
// GET ALL
// ------------------------
func (r *RiskAssessmentRepository) GetAll(
	ctx context.Context,
	search string,
	searchFilters []commondomain.SearchFilter,
	page, limit *int,
) ([]domainRisk.RiskAssessmentDefault, int64, error) {

	var riskList = make([]domainRisk.RiskAssessmentDefault, 0)
	var total int64

	db := r.db.WithContext(ctx).
		Table("risk_assessments r").
		Select(`
			r.id as Id,
			r.uuid as UUID,
			r.id_fakultas_unit as IdFakultasUnit,
			CASE WHEN fu.type = 'prodi' AND fu.jenjang IS NOT NULL AND fu.jenjang != '' THEN CONCAT(fu.nama_fak_prod_unit, ' (', UPPER(fu.jenjang), ')') ELSE fu.nama_fak_prod_unit END as NamaFakultasUnit,
			fu.uuid as UuidFakultasUnit,
			r.id_master_tahun as IdMasterTahun,
			CONVERT(COALESCE(mt.tahun, r.tahun) USING utf8mb4) as Tahun,
			mt.uuid as TahunUuid,

			r.id_mata_program as IdMataProgram,
			COALESCE(mp.mata_program, '') as MataProgram,
			mp.uuid as MataProgramUuid,
			COALESCE(NULLIF(r.isu_internal_eksternal, ''), r.penyebab, '') as IsuInternalEksternal,
			COALESCE(NULLIF(r.deskripsi_risiko, ''), r.pernyataan_risiko, '') as DeskripsiRisiko,
			COALESCE(NULLIF(r.kemungkinan_awal, 0), r.kemungkinan, 0) as KemungkinanAwal,
			COALESCE(NULLIF(r.dampak_awal, 0), r.keparahan, 0) as DampakAwal,
			COALESCE(r.tingkat_risiko, (COALESCE(NULLIF(r.kemungkinan_awal, 0), r.kemungkinan, 0) * COALESCE(NULLIF(r.dampak_awal, 0), r.keparahan, 0))) as TotalSkorAwal,
			COALESCE(NULLIF(r.tingkat_risiko_awal, ''), r.level_risiko, '') as TingkatRisikoAwal,
			COALESCE(rl_awal.nama_level, r.level_risiko, '') as NamaTingkatAwal,
			COALESCE(NULLIF(r.tindakan_mitigasi, ''), r.rencana_mitigasi, '') as TindakanMitigasi,
			COALESCE(r.peluang, '') as Peluang,
			COALESCE(r.kenyataan_residu, '') as KenyataanResidu,
			r.kemungkinan_residu as KemungkinanResidu,
			r.dampak_residu as DampakResidu,
			(COALESCE(r.kemungkinan_residu, 0) * COALESCE(r.dampak_residu, 0)) as TotalSkorResidu,
			r.tingkat_risiko_residu as TingkatRisikoResidu,
			COALESCE(rl_res.nama_level, '') as NamaTingkatResidu,
			COALESCE(r.rencana_tindak_lanjut, '') as RencanaTindakLanjut,

			COALESCE(r.kode_risiko, '') as KodeRisiko,
			COALESCE(r.kategori, '') as Kategori,
			COALESCE(r.pernyataan_risiko, r.deskripsi_risiko, '') as PernyataanRisiko,
			COALESCE(r.penyebab, r.isu_internal_eksternal, '') as Penyebab,
			COALESCE(r.dampak, '') as Dampak,
			r.kemungkinan as Kemungkinan,
			r.keparahan as Keparahan,
			r.tingkat_risiko as TingkatRisiko,
			COALESCE(r.level_risiko, r.tingkat_risiko_awal, '') as LevelRisiko,
			COALESCE(r.rencana_mitigasi, r.tindakan_mitigasi, '') as RencanaMitigasi,
			COALESCE(r.pengendalian, '') as Pengendalian,
			COALESCE(r.penanggung_jawab, '') as PenanggungJawab,
			COALESCE(r.target_waktu, '') as TargetWaktu,
			COALESCE(r.status, 'Identifikasi') as Status,
			r.created_at as CreatedAt,
			r.updated_at as UpdatedAt
		`).
		Joins("LEFT JOIN v_fakultas_unit fu ON r.id_fakultas_unit = fu.id").
		Joins("LEFT JOIN master_tahun mt ON r.id_master_tahun = mt.id").
		Joins("LEFT JOIN mata_program mp ON r.id_mata_program = mp.id").
		Joins("LEFT JOIN risk_levels rl_awal ON r.tingkat_risiko_awal = rl_awal.kode").
		Joins("LEFT JOIN risk_levels rl_res ON r.tingkat_risiko_residu = rl_res.kode")

	// Global Search
	if search != "" {
		s := "%" + search + "%"
		db = db.Where(`(
			mp.mata_program LIKE ? OR
			r.isu_internal_eksternal LIKE ? OR
			r.deskripsi_risiko LIKE ? OR
			r.tindakan_mitigasi LIKE ? OR
			r.peluang LIKE ? OR
			r.kenyataan_residu LIKE ? OR
			r.rencana_tindak_lanjut LIKE ? OR
			fu.nama_fak_prod_unit LIKE ? OR
			CONVERT(COALESCE(mt.tahun, r.tahun) USING utf8mb4) LIKE ?
		)`, s, s, s, s, s, s, s, s, s)
	}

	// Filter Columns
	for _, f := range searchFilters {
		col, ok := allowedSearchColumns[f.Field]
		if !ok {
			continue
		}
		if f.Value == nil || *f.Value == "" {
			continue
		}

		val := *f.Value
		switch f.Operator {
		case "eq":
			db = db.Where(fmt.Sprintf("%s = ?", col), val)
		case "neq":
			db = db.Where(fmt.Sprintf("%s != ?", col), val)
		case "like":
			db = db.Where(fmt.Sprintf("%s LIKE ?", col), "%"+val+"%")
		default:
			db = db.Where(fmt.Sprintf("%s = ?", col), val)
		}
	}

	// Count
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Sorting & Pagination
	db = db.Order("CONVERT(COALESCE(mt.tahun, r.tahun) USING utf8mb4) DESC, r.id DESC")
	if page != nil && limit != nil && *limit > 0 {
		offset := (*page - 1) * (*limit)
		db = db.Offset(offset).Limit(*limit)
	}

	if err := db.Find(&riskList).Error; err != nil {
		return nil, 0, err
	}

	return riskList, total, nil
}

// ------------------------
// CREATE
// ------------------------
func (r *RiskAssessmentRepository) Create(ctx context.Context, item *domainRisk.RiskAssessment) error {
	return r.db.WithContext(ctx).Table("risk_assessments").Create(item).Error
}

// ------------------------
// SAVE
// ------------------------
func (r *RiskAssessmentRepository) Save(ctx context.Context, item *domainRisk.RiskAssessment) error {
	return r.db.WithContext(ctx).Table("risk_assessments").Create(item).Error
}

// ------------------------
// UPDATE
// ------------------------
func (r *RiskAssessmentRepository) Update(ctx context.Context, item *domainRisk.RiskAssessment) error {
	return r.db.WithContext(ctx).Table("risk_assessments").Where("uuid = ?", item.UUID).Updates(item).Error
}

// ------------------------
// DELETE
// ------------------------
func (r *RiskAssessmentRepository) Delete(ctx context.Context, uid uuid.UUID) error {
	return r.db.WithContext(ctx).Table("risk_assessments").Where("uuid = ?", uid).Delete(&domainRisk.RiskAssessment{}).Error
}

// ------------------------
// MASTER OPTIONS
// ------------------------
func (r *RiskAssessmentRepository) GetCriteria(ctx context.Context, tipe string) ([]domainRisk.RiskCriteriaItem, error) {
	var items []domainRisk.RiskCriteriaItem
	q := r.db.WithContext(ctx).Table("risk_criteria").Order("tipe ASC, skor ASC")
	if tipe != "" {
		q = q.Where("tipe = ?", strings.ToUpper(tipe))
	}
	if err := q.Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *RiskAssessmentRepository) GetLevels(ctx context.Context) ([]domainRisk.RiskLevelItem, error) {
	var items []domainRisk.RiskLevelItem
	if err := r.db.WithContext(ctx).Table("risk_levels").Order("min_skor ASC, id ASC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *RiskAssessmentRepository) GetCategories(ctx context.Context) ([]domainRisk.RiskCategoryItem, error) {
	var items []domainRisk.RiskCategoryItem
	if err := r.db.WithContext(ctx).Table("risk_categories").Order("id ASC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *RiskAssessmentRepository) GetStatuses(ctx context.Context) ([]domainRisk.RiskStatusItem, error) {
	var items []domainRisk.RiskStatusItem
	if err := r.db.WithContext(ctx).Table("risk_statuses").Order("urutan ASC, id ASC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *RiskAssessmentRepository) GetPrograms(ctx context.Context, fu, thn, search string) ([]domainRisk.RiskMataProgramItem, error) {
	var programs []domainRisk.RiskMataProgramItem
	q := r.db.WithContext(ctx).Table("mata_program mp").
		Select(`
			mp.id,
			mp.uuid,
			mp.mata_program,
			mp.id_fakultas_unit,
			vfu.uuid as uuid_fakultas_unit,
			vfu.nama_fak_prod_unit as nama_fakultas_unit,
			mp.id_master_tahun,
			mt.tahun
		`).
		Joins("LEFT JOIN master_tahun mt ON mp.id_master_tahun = mt.id").
		Joins("LEFT JOIN v_fakultas_unit vfu ON mp.id_fakultas_unit = vfu.id").
		Order("mp.mata_program ASC")

	if fu != "" {
		if id, err := strconv.Atoi(fu); err == nil && id > 0 {
			q = q.Where("mp.id_fakultas_unit = ?", id)
		} else {
			q = q.Where("vfu.uuid = ?", fu)
		}
	}
	if thn != "" {
		if id, err := strconv.Atoi(thn); err == nil && id > 0 {
			q = q.Where("mp.id_master_tahun = ? OR mt.tahun = ?", id, thn)
		} else {
			q = q.Where("mt.uuid = ? OR mt.tahun = ?", thn, thn)
		}
	}
	if search != "" {
		q = q.Where("mp.mata_program LIKE ?", "%"+search+"%")
	}

	if err := q.Find(&programs).Error; err != nil {
		return nil, err
	}
	return programs, nil
}

func (r *RiskAssessmentRepository) GetOptions(ctx context.Context, fu, thn string) (*domainRisk.RiskOptionsResult, error) {
	criteria, err := r.GetCriteria(ctx, "")
	if err != nil {
		return nil, err
	}
	levels, err := r.GetLevels(ctx)
	if err != nil {
		return nil, err
	}
	categories, err := r.GetCategories(ctx)
	if err != nil {
		return nil, err
	}
	statuses, err := r.GetStatuses(ctx)
	if err != nil {
		return nil, err
	}
	programs, err := r.GetPrograms(ctx, fu, thn, "")
	if err != nil {
		return nil, err
	}

	var kemungkinan []domainRisk.RiskCriteriaItem
	var dampak []domainRisk.RiskCriteriaItem
	for _, cr := range criteria {
		if strings.EqualFold(cr.Tipe, "KEMUNGKINAN") {
			kemungkinan = append(kemungkinan, cr)
		} else if strings.EqualFold(cr.Tipe, "DAMPAK") {
			dampak = append(dampak, cr)
		}
	}

	return &domainRisk.RiskOptionsResult{
		Criteria:    criteria,
		Kemungkinan: kemungkinan,
		Dampak:      dampak,
		Levels:      levels,
		Categories:  categories,
		Statuses:    statuses,
		Programs:    programs,
	}, nil
}
