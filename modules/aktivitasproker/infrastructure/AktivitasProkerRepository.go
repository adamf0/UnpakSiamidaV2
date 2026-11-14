package infrastructure

import (
	commondomainaktivitasproker "UnpakSiamida/common/domain"
	domainaktivitasproker "UnpakSiamida/modules/aktivitasproker/domain"
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AktivitasProkerRepository struct {
	db *gorm.DB
}

func NewAktivitasProkerRepository(db *gorm.DB) domainaktivitasproker.IAktivitasProkerRepository {
	return &AktivitasProkerRepository{db: db}
}

// ------------------------
// GET BY UUID
// ------------------------
func (r *AktivitasProkerRepository) GetByUuid(ctx context.Context, uid uuid.UUID) (*domainaktivitasproker.AktivitasProker, error) {
	var AktivitasProker domainaktivitasproker.AktivitasProker

	err := r.db.WithContext(ctx).
		Where("uuid = ?", uid).
		First(&AktivitasProker).Error

	// if errors.Is(err, gorm.ErrRecordNotFound) {
	// 	return nil, nil
	// }

	if err != nil {
		return nil, err
	}

	return &AktivitasProker, nil
}

// ------------------------
// GET DEFAULT BY UUID
// ------------------------
func (r *AktivitasProkerRepository) GetDefaultByUuid(
	ctx context.Context,
	id uuid.UUID,
) (*domainaktivitasproker.AktivitasProkerDefault, error) {

	// Ambil hanya kolom yang benar-benar ada di struct AktivitasProkerDefault
	var rowData domainaktivitasproker.AktivitasProkerDefault

	err := r.db.WithContext(ctx).
		Table("aktivitas a").
		Select(`
		a.id as ID,
		a.uuid as UUID,
		a.id_fakultas_unit as FakultasUnitId,
		vfu.uuid as FakultasUnitUuid,
		CASE 
			WHEN vfu.type = 'prodi' AND vfu.jenjang IS NOT NULL AND vfu.jenjang != '' 
			THEN CONCAT(vfu.nama_fak_prod_unit, ' (', UPPER(vfu.jenjang), ')') 
			ELSE vfu.nama_fak_prod_unit 
		END as FakultasUnit,
		vfu.jenjang as Jenjang,
		vfu.type as Type,
		vfu.fakultas as Fakultas,
		a.id_mata_program as MataProgramId,
		mp.uuid as MataProgramUUID,
		mp.mata_program as MataProgram,
		mp.id_master_tahun as TahunId,
		mt.uuid as TahunUUID,
		mt.tahun as Tahun,
		a.aktivitas as Aktivitas,
		a.PIC as PIC,
		a.target_rk_awal as TanggalRKAwal,
		a.target_rk_akhir as TanggalRKAkhir,
		a.status_realisasi as StatusRealisasi,
		a.target_r_awal as TanggalRAwal,
		a.target_r_akhir as TanggalRAkhir,
		a.detail as Detail
	`).
		Joins("JOIN v_fakultas_unit vfu ON a.id_fakultas_unit = vfu.id").
		Joins("JOIN mata_program mp ON a.id_mata_program = mp.id").
		Joins("LEFT JOIN master_tahun mt ON mp.id_master_tahun = mt.id").
		Where("a.uuid = ?", id).
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
	"target":           "vfu.nama_fak_prod_unit",
	"fakultasunit":     "vfu.nama_fak_prod_unit",
	"fakultasunituuid": "vfu.uuid",
	"id_fakultas_unit": "a.id_fakultas_unit",
	"fakultasunitid":   "a.id_fakultas_unit",
	"mataprogram":      "mp.mata_program",
	"mataprogramuuid":  "mp.uuid",
	"aktivitas":        "a.aktivitas",
	"pic":              "a.PIC",
	"tahun":            "mt.tahun",
	"tahunuuid":        "mt.uuid",
	"status_realisasi": "a.status_realisasi",
	"statusrealisasi":  "a.status_realisasi",
	"detail":           "a.detail",
}

// ------------------------
// GET ALL
// ------------------------
func (r *AktivitasProkerRepository) GetAll(
	ctx context.Context,
	search string,
	searchFilters []commondomainaktivitasproker.SearchFilter,
	page, limit *int,
) ([]domainaktivitasproker.AktivitasProkerDefault, int64, error) {

	var rows = make([]domainaktivitasproker.AktivitasProkerDefault, 0)
	var total int64

	db := r.db.WithContext(ctx).
		Table("aktivitas a").
		Select(`
		a.id as ID,
		a.uuid as UUID,
		a.id_fakultas_unit as FakultasUnitId,
		vfu.uuid as FakultasUnitUuid,
		CASE 
			WHEN vfu.type = 'prodi' AND vfu.jenjang IS NOT NULL AND vfu.jenjang != '' 
			THEN CONCAT(vfu.nama_fak_prod_unit, ' (', UPPER(vfu.jenjang), ')') 
			ELSE vfu.nama_fak_prod_unit 
		END as FakultasUnit,
		vfu.jenjang as Jenjang,
		vfu.type as Type,
		vfu.fakultas as Fakultas,
		a.id_mata_program as MataProgramId,
		mp.uuid as MataProgramUUID,
		mp.mata_program as MataProgram,
		mp.id_master_tahun as TahunId,
		mt.uuid as TahunUUID,
		mt.tahun as Tahun,
		a.aktivitas as Aktivitas,
		a.PIC as PIC,
		a.target_rk_awal as TanggalRKAwal,
		a.target_rk_akhir as TanggalRKAkhir,
		a.status_realisasi as StatusRealisasi,
		a.target_r_awal as TanggalRAwal,
		a.target_r_akhir as TanggalRAkhir,
		a.detail as Detail
	`).
		Joins("JOIN v_fakultas_unit vfu ON a.id_fakultas_unit = vfu.id").
		Joins("JOIN mata_program mp ON a.id_mata_program = mp.id").
		Joins("LEFT JOIN master_tahun mt ON mp.id_master_tahun = mt.id")

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
						db = db.Where("a.id_fakultas_unit = ?", *userRow.FakultasUnit)
					} else {
						db = db.Where("1 = 0")
					}
				}
			}
			continue
		}

		if (fieldLower == "status_realisasi" || fieldLower == "status") && strings.ToLower(val) == "belum_dilaporkan" {
			db = db.Where("(a.status_realisasi IS NULL OR a.status_realisasi = '')")
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
	db = db.Order("a.id DESC")

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
func (r *AktivitasProkerRepository) Create(ctx context.Context, aktivitas *domainaktivitasproker.AktivitasProker) error {
	return r.db.WithContext(ctx).Create(aktivitas).Error
}

// ------------------------
// UPDATE
// ------------------------
func (r *AktivitasProkerRepository) Update(ctx context.Context, aktivitas *domainaktivitasproker.AktivitasProker) error {
	return r.db.WithContext(ctx).Save(aktivitas).Error
}

// ------------------------
// DELETE (by UUID)
// ------------------------
func (r *AktivitasProkerRepository) Delete(ctx context.Context, uid uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("uuid = ?", uid).
		Delete(&domainaktivitasproker.AktivitasProker{}).Error
}

func (r *AktivitasProkerRepository) ensureUuidColumn(ctx context.Context, tableName string) error {
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

func (r *AktivitasProkerRepository) SetupUuid(ctx context.Context) error {
	if err := r.ensureUuidColumn(ctx, "aktivitas"); err != nil {
		return err
	}

	const chunkSize = 500

	var ids []uint
	if err := r.db.WithContext(ctx).
		Model(&domainaktivitasproker.AktivitasProker{}).
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
			"UPDATE aktivitas SET uuid = %s WHERE id IN (?)",
			caseSQL,
		)

		if err := r.db.WithContext(ctx).
			Exec(query, args...).Error; err != nil {
			return err
		}
	}

	return nil
}
