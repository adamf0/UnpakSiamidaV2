package infrastructure

import (
	commondomainaktivitasproker "UnpakSiamida/common/domain"
	domainaktivitasproker "UnpakSiamida/modules/aktivitasproker/domain"
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
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
	query := `
		SELECT 
			a.id as ID,
			a.uuid as UUID,
			a.id_fakultas_unit as FakultasUnitId,
            vfu.uuid as FakultasUnitUuid,
			vfu.nama_fak_prod_unit as FakultasUnit,
			vfu.jenjang as Jenjang,
			vfu.type as Type,
			vfu.fakultas as Fakultas,
            a.id_mata_program as MataProgramId,
            mp.uuid as MataProgramUUID,
            mp.mata_program as MataProgram,
            a.aktivitas as Aktivitas,
            a.PIC as PIC,
			a.target_rk_awal as TanggalRKAwal,
			a.target_rk_akhir as TanggalRKAkhir
		FROM aktivitas a
		JOIN v_fakultas_unit vfu ON a.id_fakultas_unit = vfu.id 
        JOIN mata_program mp ON a.id_mata_program = mp.id
		WHERE a.uuid = ?
		LIMIT 1
	`

	var rowData domainaktivitasproker.AktivitasProkerDefault

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
func (r *AktivitasProkerRepository) GetAll(
	ctx context.Context,
	search string,
	searchFilters []commondomainaktivitasproker.SearchFilter,
	page, limit *int,
) ([]domainaktivitasproker.AktivitasProkerDefault, int64, error) {

	var rows []domainaktivitasproker.AktivitasProkerDefault
	var total int64

	db := r.db.WithContext(ctx).
		Table("aktivitas a").
		Select(`
			a.id as ID,
			a.uuid as UUID,
			a.id_fakultas_unit as FakultasUnitId,
            vfu.uuid as FakultasUnitUuid,
			vfu.nama_fak_prod_unit as FakultasUnit,
			vfu.jenjang as Jenjang,
			vfu.type as Type,
			vfu.fakultas as Fakultas,
            a.id_mata_program as MataProgramId,
            mp.uuid as MataProgramUUID,
            mp.mata_program as MataProgram,
            a.aktivitas as Aktivitas,
            a.PIC as PIC,
			a.target_rk_awal as TanggalRKAwal,
			a.target_rk_akhir as TanggalRKAkhir
		`).
		Joins(`v_fakultas_unit vfu ON jp.id_fakultas_unit = vfu.id`).
		Joins(`mata_program mp ON a.id_mata_program = mp.id`)

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

func (r *AktivitasProkerRepository) SetupUuid(ctx context.Context) error {
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
