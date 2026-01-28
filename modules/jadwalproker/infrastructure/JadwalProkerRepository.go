package infrastructure

import (
	commondomainJadwalProker "UnpakSiamida/common/domain"
	domainJadwalProker "UnpakSiamida/modules/jadwalproker/domain"
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type JadwalProkerRepository struct {
	db *gorm.DB
}

func NewJadwalProkerRepository(db *gorm.DB) domainJadwalProker.IJadwalProkerRepository {
	return &JadwalProkerRepository{db: db}
}

// ------------------------
// GET BY UUID
// ------------------------
func (r *JadwalProkerRepository) GetByUuid(ctx context.Context, uid uuid.UUID) (*domainJadwalProker.JadwalProker, error) {
	var JadwalProker domainJadwalProker.JadwalProker

	err := r.db.WithContext(ctx).
		Where("uuid = ?", uid).
		First(&JadwalProker).Error

	// if errors.Is(err, gorm.ErrRecordNotFound) {
	// 	return nil, nil
	// }

	if err != nil {
		return nil, err
	}

	return &JadwalProker, nil
}

// ------------------------
// GET DEFAULT BY UUID
// ------------------------
func (r *JadwalProkerRepository) GetDefaultByUuid(
	ctx context.Context,
	id uuid.UUID,
) (*domainJadwalProker.JadwalProkerDefault, error) {

	// Ambil hanya kolom yang benar-benar ada di struct JadwalProkerDefault
	query := `
		SELECT 
			jp.id as ID,
			jp.uuid as UUID,
			jp.id_fakultas_unit as FakultasUnitId,
			vfu.uuid as FakultasUnitUuid,
			vfu.nama_fak_prod_unit as FakultasUnit,
			vfu.jenjang as Jenjang,
			vfu.type as Type,
			vfu.fakultas as Fakultas,
			jp.tanggal_tutup as TanggalTutupEntry,
			jp.tanggal_tutup_dokumen as TanggalTutupDokumen
		FROM jadwal_proker jp
		JOIN v_fakultas_unit vfu ON jp.id_fakultas_unit = vfu.id
		WHERE jp.uuid = ?
		LIMIT 1
	`

	var rowData domainJadwalProker.JadwalProkerDefault

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
func (r *JadwalProkerRepository) GetAll(
	ctx context.Context,
	search string,
	searchFilters []commondomainJadwalProker.SearchFilter,
	page, limit *int,
) ([]domainJadwalProker.JadwalProkerDefault, int64, error) {

	var rows []domainJadwalProker.JadwalProkerDefault
	var total int64

	db := r.db.WithContext(ctx).
		Table("jadwal_proker jp").
		Select(`
			jp.id as ID,
			jp.uuid as UUID,
			jp.id_fakultas_unit as FakultasUnitId,
			vfu.uuid as FakultasUnitUuid,
			vfu.nama_fak_prod_unit as FakultasUnit,
			vfu.jenjang as Jenjang,
			vfu.type as Type,
			vfu.fakultas as Fakultas,
			jp.tanggal_tutup as TanggalTutupEntry,
			jp.tanggal_tutup_dokumen as TanggalTutupDokumen
		`).
		Joins(`v_fakultas_unit vfu ON jp.id_fakultas_unit = vfu.id`)

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
	db = db.Order("jp.id DESC")

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
func (r *JadwalProkerRepository) Create(ctx context.Context, jadwalproker *domainJadwalProker.JadwalProker) error {
	return r.db.WithContext(ctx).Create(jadwalproker).Error
}

// ------------------------
// UPDATE
// ------------------------
func (r *JadwalProkerRepository) Update(ctx context.Context, jadwalproker *domainJadwalProker.JadwalProker) error {
	return r.db.WithContext(ctx).Save(jadwalproker).Error
}

// ------------------------
// DELETE (by UUID)
// ------------------------
func (r *JadwalProkerRepository) Delete(ctx context.Context, uid uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("uuid = ?", uid).
		Delete(&domainJadwalProker.JadwalProker{}).Error
}

func (r *JadwalProkerRepository) SetupUuid(ctx context.Context) error {
	const chunkSize = 500

	var ids []uint
	if err := r.db.WithContext(ctx).
		Model(&domainJadwalProker.JadwalProker{}).
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
			"UPDATE jadwal_proker SET uuid = %s WHERE id IN (?)",
			caseSQL,
		)

		if err := r.db.WithContext(ctx).
			Exec(query, args...).Error; err != nil {
			return err
		}
	}

	return nil
}
