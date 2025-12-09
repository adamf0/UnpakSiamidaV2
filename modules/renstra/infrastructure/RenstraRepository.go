package infrastructure

import (
	"context"
	commondomainrenstra "UnpakSiamida/common/domain"
	domainrenstra "UnpakSiamida/modules/renstra/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"strings"
	"fmt"
)

type RenstraRepository struct {
	db *gorm.DB
}

func NewRenstraRepository(db *gorm.DB) domainrenstra.IRenstraRepository {
	return &RenstraRepository{db: db}
}

func (r *RenstraRepository) IsUnique(
	ctx context.Context,
	fakultas_unit uint,
	tahun string,
) (bool, error) {

	var count int64

	err := r.db.WithContext(ctx).
		Model(&domainrenstra.Renstra{}).
		Where("fakultas_unit = ?", fakultas_unit).
		Where("tahun = ?", tahun).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count == 0, nil
}

// ------------------------
// GET BY UUID
// ------------------------
func (r *RenstraRepository) GetByUuid(ctx context.Context, uid uuid.UUID) (*domainrenstra.Renstra, error) {
	var renstra domainrenstra.Renstra

	err := r.db.WithContext(ctx).
		Where("uuid = ?", uid).
		First(&renstra).Error

	// if errors.Is(err, gorm.ErrRecordNotFound) {
	// 	return nil, nil
	// }

	if err != nil {
		return nil, err
	}

	return &renstra, nil
}

// ------------------------
// GET DEFAULT BY UUID
// ------------------------
func (r *RenstraRepository) GetDefaultByUuid(
	ctx context.Context,
	id uuid.UUID,
) (*domainrenstra.RenstraDefault, error) {

	query := `
		SELECT 
			r.id as Id,
			r.uuid as Uuid,
			r.tahun as Tahun,
			r.fakultas_unit as FakultasUnit,
			r.periode_upload_mulai as PeriodeUploadMulai,
			r.periode_upload_akhir as PeriodeUploadAkhir,
			r.periode_assesment_dokumen_mulai as PeriodeAssesmentDokumenMulai,
			r.periode_assesment_dokumen_akhir as PeriodeAssesmentDokumenAkhir,
			r.periode_assesment_lapangan_mulai as PeriodeAssesmentLapanganMulai,
			r.periode_assesment_lapangan_akhir as PeriodeAssesmentLapanganAkhir,
			r.auditee as Auditee,
			r.auditor1 as Auditor1,
			r.auditor2  as Auditor2,
			r.kodeAkses as KodeAkses,
			r.catatan as Catatan1,
			r.catatan2 as Catatan2,
			a0.name AS NamaAuditee,
			a1.name AS NamaAuditor1,
			a2.name AS NamaAuditor2,
			vfu.uuid AS UUIDFakultasUnit,
			vfu.nama_fak_prod_unit as NamaFakultasUnit,
			vfu.jenjang as Jenjang,
			vfu.type as Type,
			vfu.fakultas as Fakultas
		FROM renstra r
		JOIN v_fakultas_unit vfu ON r.fakultas_unit = vfu.id
		LEFT JOIN users a0 ON r.auditee = a0.id
		LEFT JOIN users a1 ON r.auditor1 = a1.id
		LEFT JOIN users a2 ON r.auditor2 = a2.id
		WHERE r.uuid = ?
		ORDER BY r.tahun DESC
		LIMIT 1
	`

	var rowData domainrenstra.RenstraDefault

	row := r.db.WithContext(ctx).Raw(query, id).Row()
	err := row.Scan(
		&rowData.ID,
		&rowData.UUID,
		&rowData.Tahun,
		&rowData.FakultasUnit,
		&rowData.PeriodeUploadMulai,
		&rowData.PeriodeUploadAkhir,
		&rowData.PeriodeAssesmentDokumenMulai,
		&rowData.PeriodeAssesmentDokumenAkhir,
		&rowData.PeriodeAssesmentLapanganMulai,
		&rowData.PeriodeAssesmentLapanganAkhir,
		&rowData.Auditee,
		&rowData.Auditor1,
		&rowData.Auditor2,
		&rowData.KodeAkses,
		&rowData.Catatan1,
		&rowData.Catatan2,
		&rowData.NamaAuditee,
		&rowData.NamaAuditor1,
		&rowData.NamaAuditor2,
		&rowData.UUIDFakultasUnit,
		&rowData.NamaFakultasUnit,
		&rowData.Jenjang,
		&rowData.Type,
		&rowData.Fakultas,
	)

	if err != nil {
		return nil, err
	}

	if rowData.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &rowData, nil
}

var allowedSearchColumns = map[string]string{
	"tahun":         "r.tahun",
	"fakultas_unit": "r.fakultas_unit",
	"auditee":       "a0.auditee",
	"auditor1":      "a1.auditor1",
	"auditor2":      "a2.auditor2",
	"kodeakses":     "r.kodeAkses",
	"jenjang":     	 "vfu.jenjang",
	"fakultas":      "vfu.fakultas",
}

// ------------------------
// GET ALL
// ------------------------
func (r *RenstraRepository) GetAll(
	ctx context.Context,
	search string,
	searchFilters []commondomainrenstra.SearchFilter,
	page, limit *int,
) ([]domainrenstra.RenstraDefault, int64, error) {

	var renstras []domainrenstra.RenstraDefault
	var total int64

	db := r.db.WithContext(ctx).Table("renstra r").
		Select(`
			r.id as Id,
			r.uuid as Uuid,
			r.tahun as Tahun,
			r.fakultas_unit as FakultasUnit,
			r.periode_upload_mulai as PeriodeUploadMulai,
			r.periode_upload_akhir as PeriodeUploadAkhir,
			r.periode_assesment_dokumen_mulai as PeriodeAssesmentDokumenMulai,
			r.periode_assesment_dokumen_akhir as PeriodeAssesmentDokumenAkhir,
			r.periode_assesment_lapangan_mulai as PeriodeAssesmentLapanganMulai,
			r.periode_assesment_lapangan_akhir as PeriodeAssesmentLapanganAkhir,
			r.auditee as Auditee,
			r.auditor1 as Auditor1,
			r.auditor2  as Auditor2,
			r.kodeAkses as KodeAkses,
			r.catatan as Catatan1,
			r.catatan2 as Catatan2,
			a0.name AS NamaAuditee,
			a1.name AS NamaAuditor1,
			a2.name AS NamaAuditor2,
			vfu.uuid AS UUIDFakultasUnit,
			vfu.nama_fak_prod_unit as NamaFakultasUnit,
			vfu.jenjang as Jenjang,
			vfu.type as Type,
			vfu.fakultas as Fakultas
		`).
		Joins("JOIN v_fakultas_unit vfu ON r.fakultas_unit = vfu.id").
		Joins("LEFT JOIN users a0 ON r.auditee = a0.id").
		Joins("LEFT JOIN users a1 ON r.auditor1 = a1.id").
		Joins("LEFT JOIN users a2 ON r.auditor2 = a2.id")

	// -------------------------------
	// SEARCH FILTERS (ADVANCED)
	// -------------------------------
	if len(searchFilters) > 0 {
		for _, f := range searchFilters {
			field := strings.TrimSpace(strings.ToLower(f.Field))
			operator := strings.TrimSpace(strings.ToLower(f.Operator))
			
			var value string
			if f.Value != nil {
				value = strings.TrimSpace(*f.Value)
			} else {
				value = "" // nil dianggap kosong
			}

			// if value == "" {
			// 	continue
			// }

			// Validate allowed column
			col, ok := allowedSearchColumns[field]
			if !ok {
				continue // skip unknown field
			}

			switch operator {
			case "eq":
				db = db.Where(fmt.Sprintf("%s = ?", col), value)
			case "neq":
				db = db.Where(fmt.Sprintf("%s <> ?", col), value)
			case "like":
				db = db.Where(fmt.Sprintf("%s LIKE ?", col), "%"+value+"%")
			case "gt":
				db = db.Where(fmt.Sprintf("%s > ?", col), value)
			case "gte":
				db = db.Where(fmt.Sprintf("%s >= ?", col), value)
			case "lt":
				db = db.Where(fmt.Sprintf("%s < ?", col), value)
			case "lte":
				db = db.Where(fmt.Sprintf("%s <= ?", col), value)
			case "in":
				db = db.Where(fmt.Sprintf("%s IN (?)", col), strings.Split(value, ","))
			default:
				// default fallback → LIKE
				db = db.Where(fmt.Sprintf("%s LIKE ?", col), "%"+value+"%")
			}
		}

	}
	if strings.TrimSpace(search) != "" {

		// -------------------------------
		// GLOBAL SEARCH
		// -------------------------------
		like := "%" + search + "%"
		var orParts []string
		var params []interface{}

		for _, col := range allowedSearchColumns {
			orParts = append(orParts, fmt.Sprintf("%s LIKE ?", col))
			params = append(params, like)
		}

		db = db.Where("(" + strings.Join(orParts, " OR ") + ")", params...)
	}

	// -------------------------------
	// COUNT
	// -------------------------------
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// -------------------------------
	// PAGINATION
	// -------------------------------
	if page != nil && limit != nil && *limit > 0 {
		p := *page
		l := *limit

		if p < 1 {
			p = 1
		}

		offset := (p - 1) * l
		db = db.Offset(offset).Limit(l)
	}

	// -------------------------------
	// EXECUTE QUERY
	// -------------------------------
	if err := db.Order("r.tahun DESC").Find(&renstras).Error; err != nil {
		return nil, 0, err
	}

	return renstras, total, nil
}


// ------------------------
// CREATE
// ------------------------
func (r *RenstraRepository) Create(ctx context.Context, renstra *domainrenstra.Renstra) error {
	return r.db.WithContext(ctx).Create(renstra).Error
}

// ------------------------
// UPDATE
// ------------------------
func (r *RenstraRepository) Update(ctx context.Context, renstra *domainrenstra.Renstra) error {
	return r.db.WithContext(ctx).Save(renstra).Error
}

// ------------------------
// DELETE (by UUID)
// ------------------------
func (r *RenstraRepository) Delete(ctx context.Context, uid uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("uuid = ?", uid).
		Delete(&domainrenstra.Renstra{}).Error
}
