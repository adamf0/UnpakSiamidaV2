package infrastructure

import (
	commondomaindokumenproker "UnpakSiamida/common/domain"
	domaindokumenproker "UnpakSiamida/modules/dokumenproker/domain"
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
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
	query := `
		SELECT 
			drp.id as ID,
			drp.uuid as UUID,
			drp.id_fakultas_unit as FakultasUnitId,
            vfu.uuid as FakultasUnitUuid,
			vfu.nama_fak_prod_unit as FakultasUnit,
			vfu.jenjang as Jenjang,
			vfu.type as Type,
			vfu.fakultas as Fakultas,
            drp.id_mata_program as MataProgramId,
            mp.uuid as MataProgramUUID,
            mp.mata_program as MataProgram,
            drp.jenis_dokumen as JenisDokumen,
            drp.file as File,
            drp.status_verifikasi as Status,
            drp.catatan as Catatan
		FROM dokumen_realisasi_proker drp
		JOIN v_fakultas_unit vfu ON drp.id_fakultas_unit = vfu.id 
        JOIN mata_program mp ON drp.id_mata_program = mp.id
		WHERE a.uuid = ?
		LIMIT 1
	`

	var rowData domaindokumenproker.DokumenProkerDefault

	err := r.db.WithContext(ctx).Raw(query, id).Scan(&rowData).Error
	if err != nil {
		return nil, err
	}

	// Jika tidak ada row → struct kosong → anggap record not found
	if rowData.Id == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &rowData, nil
}

var allowedSearchColumns = map[string]string{
	// key:param -> db column
	"target": "vfu.nama_fak_prod_unit",
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

	var rows []domaindokumenproker.DokumenProkerDefault
	var total int64

	db := r.db.WithContext(ctx).
		Table("dokumen_realisasi_proker drp").
		Select(`
			drp.id as ID,
			drp.uuid as UUID,
			drp.id_fakultas_unit as FakultasUnitId,
            vfu.uuid as FakultasUnitUuid,
			vfu.nama_fak_prod_unit as FakultasUnit,
			vfu.jenjang as Jenjang,
			vfu.type as Type,
			vfu.fakultas as Fakultas,
            drp.id_mata_program as MataProgramId,
            mp.uuid as MataProgramUUID,
            mp.mata_program as MataProgram,
            drp.jenis_dokumen as JenisDokumen,
            drp.file as File,
            drp.status_verifikasi as Status,
            drp.catatan as Catatan
		`).
		Joins(`v_fakultas_unit vfu ON jp.id_fakultas_unit = vfu.id`).
		Joins(`mata_program mp ON drp.id_mata_program = mp.id`)

	// -----------------------------------
	// ADVANCED FILTERS
	// -----------------------------------
	for _, f := range searchFilters {
		col, ok := allowedSearchColumns[strings.ToLower(f.Field)]
		if !ok {
			continue
		}

		val := ""
		if f.Value != nil {
			val = strings.TrimSpace(*f.Value)
		}

		switch strings.ToLower(f.Operator) {
		case "eq":
			db = db.Where(col+" = ?", val)
		case "neq":
			db = db.Where(col+" <> ?", val)
		case "like":
			db = db.Where(col+" LIKE ?", "%"+val+"%")
		case "in":
			db = db.Where(col+" IN ?", strings.Split(val, ","))
		}
	}

	// -----------------------------------
	// GLOBAL SEARCH
	// -----------------------------------
	if strings.TrimSpace(search) != "" {
		like := "%" + search + "%"
		var or []string
		var args []interface{}

		for _, col := range allowedSearchColumns {
			or = append(or, col+" LIKE ?")
			args = append(args, like)
		}

		db = db.Where("("+strings.Join(or, " OR ")+")", args...)
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

func (r *DokumenProkerRepository) SetupUuid(ctx context.Context) error {
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
