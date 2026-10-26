package infrastructure

import (
	commondomaindokumenproker "UnpakSiamida/common/domain"
	domaindokumenproker "UnpakSiamida/modules/dokumenproker/domain"
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DokumenProkerRepository struct {
	db *gorm.DB
}

func NewDokumenProkerRepository(db *gorm.DB) domaindokumenproker.IDokumenProkerRepository {
	return &DokumenProkerRepository{db: db}
}

// ------------------------
// GET BY UUID
// ------------------------
func (r *DokumenProkerRepository) GetByUuid(ctx context.Context, uid uuid.UUID) (*domaindokumenproker.DokumenProker, error) {
	var DokumenProker domaindokumenproker.DokumenProker

	err := r.db.WithContext(ctx).
		Where("uuid = ?", uid).
		First(&DokumenProker).Error

	// if errors.Is(err, gorm.ErrRecordNotFound) {
	// 	return nil, nil
	// }

	if err != nil {
		return nil, err
	}

	return &DokumenProker, nil
}

// ------------------------
// GET DEFAULT BY UUID
// ------------------------
func (r *DokumenProkerRepository) GetDefaultByUuid(
	ctx context.Context,
	id uuid.UUID,
) (*domaindokumenproker.DokumenProkerDefault, error) {

	// Ambil hanya kolom yang benar-benar ada di struct DokumenProkerDefault
	var rowData domaindokumenproker.DokumenProkerDefault

	err := r.db.WithContext(ctx).
		Table("dokumen_realisasi_proker drp").
		Select(`
		drp.id as ID,
		drp.uuid as UUID,
		drp.id_fakultas_unit as FakultasUnitId,
		vfu.uuid as FakultasUnitUuid,
		CASE 
			WHEN vfu.type = 'prodi' AND vfu.jenjang IS NOT NULL AND vfu.jenjang != '' 
			THEN CONCAT(vfu.nama_fak_prod_unit, ' (', UPPER(vfu.jenjang), ')') 
			ELSE vfu.nama_fak_prod_unit 
		END as FakultasUnit,
		vfu.jenjang as Jenjang,
		vfu.type as Type,
		vfu.fakultas as Fakultas,
		drp.id_mata_program as MataProgramId,
		mp.uuid as MataProgramUUID,
		mp.mata_program as MataProgram,
		mp.id_master_tahun as TahunId,
		mt.uuid as TahunUUID,
		mt.tahun as Tahun,
		drp.jenis_dokumen as JenisDokumen,
		drp.file as File,
		drp.status_verifikasi as Status,
		drp.catatan as Catatan
	`).
		Joins("JOIN v_fakultas_unit vfu ON drp.id_fakultas_unit = vfu.id").
		Joins("JOIN mata_program mp ON drp.id_mata_program = mp.id").
		Joins("LEFT JOIN master_tahun mt ON mp.id_master_tahun = mt.id").
		Where("drp.uuid = ?", id).
		Take(&rowData).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}

	return &rowData, nil
}

var allowedSearchColumns = map[string]string{
	// key:param -> db column
	"target":            "vfu.nama_fak_prod_unit",
	"fakultasunit":      "vfu.nama_fak_prod_unit",
	"fakultasunituuid":  "vfu.uuid",
	"id_fakultas_unit":  "drp.id_fakultas_unit",
	"fakultasunitid":    "drp.id_fakultas_unit",
	"mataprogram":       "mp.mata_program",
	"mataprogramuuid":   "mp.uuid",
	"id_mata_program":   "drp.id_mata_program",
	"tahun":             "mt.tahun",
	"tahunuuid":         "mt.uuid",
	"tahunid":           "mp.id_master_tahun",
	"jenis_dokumen":     "drp.jenis_dokumen",
	"jenisdokumen":      "drp.jenis_dokumen",
	"status":            "drp.status_verifikasi",
	"status_verifikasi": "drp.status_verifikasi",
	"catatan":           "drp.catatan",
	"file":              "drp.file",
}

// ------------------------
// GET ALL
// ------------------------
func (r *DokumenProkerRepository) GetAll(
	ctx context.Context,
	search string,
	searchFilters []commondomaindokumenproker.SearchFilter,
	page, limit *int,
) ([]domaindokumenproker.DokumenProkerDefault, int64, error) {

	var rows = make([]domaindokumenproker.DokumenProkerDefault, 0)
	var total int64

	// Cek apakah filter tahun mengarah ke tahun >= 2027 (V3)
	var filterTahunVal string
	for _, f := range searchFilters {
		fieldLower := strings.ToLower(strings.TrimSpace(f.Field))
		if fieldLower == "tahun" || fieldLower == "tahunuuid" || fieldLower == "tahunid" || fieldLower == "filter_tahun" {
			if f.Value != nil {
				filterTahunVal = strings.TrimSpace(*f.Value)
				break
			}
		}
	}

	isV3 := false
	if filterTahunVal != "" {
		var mtCheck struct {
			Tahun string
		}
		if _, err := uuid.Parse(filterTahunVal); err == nil {
			_ = r.db.WithContext(ctx).Table("master_tahun").Select("tahun").Where("uuid = ?", filterTahunVal).Take(&mtCheck).Error
		} else {
			_ = r.db.WithContext(ctx).Table("master_tahun").Select("tahun").Where("tahun = ?", filterTahunVal).Take(&mtCheck).Error
		}
		tNum, _ := strconv.Atoi(strings.TrimSpace(mtCheck.Tahun))
		if tNum >= 2027 {
			isV3 = true
		}
	}

	if isV3 {
		// V3: Mengambil data dokumen dari matrix_proker_dokumen
		db := r.db.WithContext(ctx).
			Table("matrix_proker_dokumen mpd").
			Select(`
				mpd.id as Id,
				mpd.uuid as UUID,
				COALESCE(mpr.id_fakultas_unit, mp.id_fakultas_unit, 0) as FakultasUnitId,
				vfu.uuid as FakultasUnitUuid,
				CASE 
					WHEN vfu.type = 'prodi' AND vfu.jenjang IS NOT NULL AND vfu.jenjang != '' 
					THEN CONCAT(vfu.nama_fak_prod_unit, ' (', UPPER(vfu.jenjang), ')') 
					ELSE vfu.nama_fak_prod_unit 
				END as FakultasUnit,
				vfu.jenjang as Jenjang,
				vfu.type as Type,
				vfu.fakultas as Fakultas,
				COALESCE(mp.id, 0) as MataProgramId,
				mp.uuid as MataProgramUUID,
				COALESCE(mp.mata_program, '') as MataProgram,
				COALESCE(mt.id, 0) as TahunId,
				mt.uuid as TahunUUID,
				COALESCE(mt.tahun, '') as Tahun,
				mpd.jenis_dokumen as JenisDokumen,
				COALESCE(mpd.file_url, mpd.file, '') as File,
				mpd.status_verifikasi as Status,
				mpd.catatan as Catatan,
				mpd.nama_file as NamaFile,
				mpd.file_url as FileURL,
				mpd.file_size as FileSize,
				mpd.file_type as FileType,
				mpd.created_at as CreatedAt
			`).
			Joins("JOIN matrix_proker mpr ON mpd.id_matrik_proker = mpr.id").
			Joins("LEFT JOIN mata_program mp ON mpr.id_mata_program = mp.id").
			Joins("LEFT JOIN master_tahun mt ON mp.id_master_tahun = mt.id").
			Joins("LEFT JOIN v_fakultas_unit vfu ON COALESCE(mpr.id_fakultas_unit, mp.id_fakultas_unit) = vfu.id").
			Where("mpd.deleted_at IS NULL")

		allowedSearchColumnsV3 := map[string]string{
			"target":            "vfu.nama_fak_prod_unit",
			"fakultasunit":      "vfu.nama_fak_prod_unit",
			"fakultasunituuid":  "vfu.uuid",
			"id_fakultas_unit":  "mpr.id_fakultas_unit",
			"fakultasunitid":    "mpr.id_fakultas_unit",
			"mataprogram":       "mp.mata_program",
			"mataprogramuuid":   "mp.uuid",
			"id_mata_program":   "mpr.id_mata_program",
			"tahun":             "mt.tahun",
			"tahunuuid":         "mt.uuid",
			"tahunid":           "mp.id_master_tahun",
			"jenis_dokumen":     "mpd.jenis_dokumen",
			"jenisdokumen":      "mpd.jenis_dokumen",
			"status":            "mpd.status_verifikasi",
			"status_verifikasi": "mpd.status_verifikasi",
			"catatan":           "mpd.catatan",
			"file":              "mpd.file",
			"nama_file":         "mpd.nama_file",
		}

		for _, f := range searchFilters {
			fieldLower := strings.ToLower(strings.TrimSpace(f.Field))
			val := ""
			if f.Value != nil {
				val = strings.TrimSpace(*f.Value)
			}
			if val == "" {
				continue
			}

			if fieldLower == "sid" {
				var userRow struct {
					Level        string `gorm:"column:level"`
					FakultasUnit *uint  `gorm:"column:fakultas_unit"`
				}
				if err := r.db.WithContext(ctx).
					Table("users").
					Select("level, fakultas_unit").
					Where("uuid = ?", val).
					Take(&userRow).Error; err == nil {
					if strings.ToLower(strings.TrimSpace(userRow.Level)) != "admin" {
						if userRow.FakultasUnit != nil && *userRow.FakultasUnit > 0 {
							db = db.Where("mpr.id_fakultas_unit = ?", *userRow.FakultasUnit)
						} else {
							db = db.Where("1 = 0")
						}
					}
				}
				continue
			}

			if (fieldLower == "status_verifikasi" || fieldLower == "status") && strings.ToLower(val) == "belum_terverifikasi" {
				db = db.Where("(mpd.status_verifikasi = 'belum_terverifikasi' OR mpd.status_verifikasi IS NULL OR mpd.status_verifikasi = '')")
				continue
			}

			col, ok := allowedSearchColumnsV3[fieldLower]
			if !ok {
				continue
			}

			switch strings.ToLower(f.Operator) {
			case "eq":
				db = db.Where(clause.Eq{Column: col, Value: val})
			case "neq":
				db = db.Where(clause.Neq{Column: col, Value: val})
			case "like":
				db = db.Where(clause.Like{Column: col, Value: "%" + val + "%"})
			case "in":
				rawVals := strings.Split(val, ",")
				vals := make([]interface{}, 0, len(rawVals))
				for _, v := range rawVals {
					v = strings.TrimSpace(v)
					if v != "" {
						vals = append(vals, v)
					}
				}
				if len(vals) > 0 {
					db = db.Where(clause.IN{Column: col, Values: vals})
				}
			}
		}

		if strings.TrimSpace(search) != "" {
			like := "%" + search + "%"
			var conditions []clause.Expression
			for _, col := range allowedSearchColumnsV3 {
				conditions = append(conditions, clause.Like{Column: col, Value: like})
			}
			if len(conditions) > 0 {
				db = db.Where(clause.Or(conditions...))
			}
		}

		if err := db.Count(&total).Error; err != nil {
			return nil, 0, err
		}

		db = db.Order("mpd.id DESC")

		if page != nil && limit != nil && *limit > 0 {
			offset := (*page - 1) * (*limit)
			db = db.Offset(offset).Limit(*limit)
		}

		if err := db.Find(&rows).Error; err != nil {
			return nil, 0, err
		}

		return rows, total, nil
	}

	// Default V2: dokumen_realisasi_proker
	db := r.db.WithContext(ctx).
		Table("dokumen_realisasi_proker drp").
		Select(`
			drp.id as ID,
			drp.uuid as UUID,
			drp.id_fakultas_unit as FakultasUnitId,
			vfu.uuid as FakultasUnitUuid,
			CASE 
				WHEN vfu.type = 'prodi' AND vfu.jenjang IS NOT NULL AND vfu.jenjang != '' 
				THEN CONCAT(vfu.nama_fak_prod_unit, ' (', UPPER(vfu.jenjang), ')') 
				ELSE vfu.nama_fak_prod_unit 
			END as FakultasUnit,
			vfu.jenjang as Jenjang,
			vfu.type as Type,
			vfu.fakultas as Fakultas,
			drp.id_mata_program as MataProgramId,
			mp.uuid as MataProgramUUID,
			mp.mata_program as MataProgram,
			mp.id_master_tahun as TahunId,
			mt.uuid as TahunUUID,
			mt.tahun as Tahun,
			drp.jenis_dokumen as JenisDokumen,
			drp.file as File,
			drp.status_verifikasi as Status,
			drp.catatan as Catatan
		`).
		Joins(`JOIN v_fakultas_unit vfu ON drp.id_fakultas_unit = vfu.id`).
		Joins(`JOIN mata_program mp ON drp.id_mata_program = mp.id`).
		Joins(`LEFT JOIN master_tahun mt ON mp.id_master_tahun = mt.id`)

	// -----------------------------------
	// ADVANCED FILTERS
	// -----------------------------------
	for _, f := range searchFilters {
		fieldLower := strings.ToLower(strings.TrimSpace(f.Field))
		val := ""
		if f.Value != nil {
			val = strings.TrimSpace(*f.Value)
		}
		if val == "" {
			continue
		}

		if fieldLower == "sid" {
			var userRow struct {
				Level        string `gorm:"column:level"`
				FakultasUnit *uint  `gorm:"column:fakultas_unit"`
			}
			if err := r.db.WithContext(ctx).
				Table("users").
				Select("level, fakultas_unit").
				Where("uuid = ?", val).
				Take(&userRow).Error; err == nil {
				if strings.ToLower(strings.TrimSpace(userRow.Level)) != "admin" {
					if userRow.FakultasUnit != nil && *userRow.FakultasUnit > 0 {
						db = db.Where("drp.id_fakultas_unit = ?", *userRow.FakultasUnit)
					} else {
						db = db.Where("1 = 0")
					}
				}
			}
			continue
		}

		if (fieldLower == "status_verifikasi" || fieldLower == "status") && strings.ToLower(val) == "belum_terverifikasi" {
			db = db.Where("(drp.status_verifikasi = 'belum_terverifikasi' OR drp.status_verifikasi IS NULL OR drp.status_verifikasi = '')")
			continue
		}

		col, ok := allowedSearchColumns[fieldLower]
		if !ok {
			continue
		}

		switch strings.ToLower(f.Operator) {
		case "eq":
			db = db.Where(clause.Eq{
				Column: col,
				Value:  val,
			})
		case "neq":
			db = db.Where(clause.Neq{
				Column: col,
				Value:  val,
			})
		case "like":
			db = db.Where(clause.Like{
				Column: col,
				Value:  "%" + val + "%",
			})
		case "in":
			rawVals := strings.Split(val, ",")
			vals := make([]interface{}, 0, len(rawVals))

			for _, v := range rawVals {
				v = strings.TrimSpace(v)
				if v != "" {
					vals = append(vals, v)
				}
			}

			if len(vals) > 0 {
				db = db.Where(clause.IN{
					Column: col,
					Values: vals,
				})
			}
		}
	}

	// -----------------------------------
	// GLOBAL SEARCH
	// -----------------------------------
	if strings.TrimSpace(search) != "" {
		like := "%" + search + "%"
		var conditions []clause.Expression

		for _, col := range allowedSearchColumns {
			conditions = append(conditions, clause.Like{
				Column: col,
				Value:  like,
			})
		}

		if len(conditions) > 0 {
			db = db.Where(clause.Or(conditions...))
		}
	}

	// -----------------------------------
	// COUNT (AMAN)
	// -----------------------------------
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// -----------------------------------
	// ORDER + PAGINATION
	// -----------------------------------
	db = db.Order("drp.id DESC")

	if page != nil && limit != nil && *limit > 0 {
		offset := (*page - 1) * (*limit)
		db = db.Offset(offset).Limit(*limit)
	}

	// -----------------------------------
	// EXECUTE
	// -----------------------------------
	if err := db.Find(&rows).Error; err != nil {
		return nil, 0, err
	}

	return rows, total, nil
}

// ------------------------
// CREATE
// ------------------------
func (r *DokumenProkerRepository) Create(ctx context.Context, aktivitas *domaindokumenproker.DokumenProker) error {
	return r.db.WithContext(ctx).Create(aktivitas).Error
}

// ------------------------
// UPDATE
// ------------------------
func (r *DokumenProkerRepository) Update(ctx context.Context, aktivitas *domaindokumenproker.DokumenProker) error {
	return r.db.WithContext(ctx).Save(aktivitas).Error
}

// ------------------------
// DELETE (by UUID)
// ------------------------
func (r *DokumenProkerRepository) Delete(ctx context.Context, uid uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("uuid = ?", uid).
		Delete(&domaindokumenproker.DokumenProker{}).Error
}

func (r *DokumenProkerRepository) ensureUuidColumn(ctx context.Context, tableName string) error {
	var colCount int64
	checkSQL := `
		SELECT COUNT(*)
		FROM INFORMATION_SCHEMA.COLUMNS
		WHERE TABLE_SCHEMA = DATABASE()
		  AND TABLE_NAME = ?
		  AND COLUMN_NAME = 'uuid'
	`
	if err := r.db.WithContext(ctx).Raw(checkSQL, tableName).Scan(&colCount).Error; err != nil {
		return err
	}

	if colCount == 0 {
		alterSQL := fmt.Sprintf("ALTER TABLE `%s` ADD COLUMN `uuid` CHAR(36) NULL AFTER `id`", tableName)
		if err := r.db.WithContext(ctx).Exec(alterSQL).Error; err != nil {
			return err
		}
	}

	return nil
}

func (r *DokumenProkerRepository) SetupUuid(ctx context.Context) error {
	if err := r.ensureUuidColumn(ctx, "dokumen_realisasi_proker"); err != nil {
		return err
	}

	const chunkSize = 500

	var ids []uint
	if err := r.db.WithContext(ctx).
		Model(&domaindokumenproker.DokumenProker{}).
		Where("uuid IS NULL OR uuid = ''").
		Pluck("id", &ids).Error; err != nil {
		return err
	}

	if len(ids) == 0 {
		return nil
	}

	for i := 0; i < len(ids); i += chunkSize {
		end := i + chunkSize
		if end > len(ids) {
			end = len(ids)
		}

		chunk := ids[i:end]

		caseSQL := "CASE id "
		args := make([]any, 0, len(chunk)*2+1)

		for _, id := range chunk {
			u := uuid.NewString()
			caseSQL += "WHEN ? THEN ? "
			args = append(args, id, u)
		}

		caseSQL += "END"
		args = append(args, chunk)

		query := fmt.Sprintf(
			"UPDATE dokumen_realisasi_proker SET uuid = %s WHERE id IN (?)",
			caseSQL,
		)

		if err := r.db.WithContext(ctx).
			Exec(query, args...).Error; err != nil {
			return err
		}
	}

	return nil
}
