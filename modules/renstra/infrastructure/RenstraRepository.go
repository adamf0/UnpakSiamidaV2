package infrastructure

import (
	commondomainrenstra "UnpakSiamida/common/domain"
	domainrenstra "UnpakSiamida/modules/renstra/domain"
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
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
			r.id as ID,
			r.uuid as UUID,
			r.tahun as Tahun,

			r.fakultas_unit as FakultasUnitId,
			vfu.uuid as FakultasUnitUuid,
			vfu.nama_fak_prod_unit as FakultasUnit,

			r.periode_upload_mulai as PeriodeUploadMulai,
			r.periode_upload_akhir as PeriodeUploadAkhir,
			r.periode_assesment_dokumen_mulai as PeriodeAssesmentDokumenMulai,
			r.periode_assesment_dokumen_akhir as PeriodeAssesmentDokumenAkhir,
			r.periode_assesment_lapangan_mulai as PeriodeAssesmentLapanganMulai,
			r.periode_assesment_lapangan_akhir as PeriodeAssesmentLapanganAkhir,

			r.auditee as AuditeeId,
			r.auditor1 as Auditor1Id,
			r.auditor2 as Auditor2Id,

			r.kodeAkses as KodeAkses,
			r.catatan as Catatan1,
			r.catatan2 as Catatan2,

			a0.name as Auditee,
			a1.name as Auditor1,
			a2.name as Auditor2,

			a0.uuid as AuditeeUuid,
			a1.uuid as Auditor1Uuid,
			a2.uuid as Auditor2Uuid,

			vfu.jenjang as Jenjang,
			vfu.type as Type,
			vfu.fakultas as Fakultas,

			COUNT(DISTINCT rn.id) AS TotalRenstra,
			SUM(CASE WHEN rn.capaian IS NOT NULL THEN 1 ELSE 0 END) AS TotalRenstraAuditee,
			SUM(CASE WHEN rn.capaian_auditor IS NOT NULL THEN 1 ELSE 0 END) AS TotalRenstraAuditor,

			COUNT(DISTINCT dt.id) AS TotalDokumen,
			SUM(CASE WHEN dt.file IS NOT NULL THEN 1 ELSE 0 END) AS TotalDokumenAuditee,
			SUM(CASE WHEN dt.capaian_auditor IS NOT NULL THEN 1 ELSE 0 END) AS TotalDokumenAuditor
		FROM renstra r
		JOIN v_fakultas_unit vfu ON r.fakultas_unit = vfu.id
		LEFT JOIN users a0 ON r.auditee = a0.id
		LEFT JOIN users a1 ON r.auditor1 = a1.id
		LEFT JOIN users a2 ON r.auditor2 = a2.id
		LEFT JOIN renstra_nilai rn  ON rn.id_renstra = r.id
		LEFT JOIN dokumen_tambahan dt ON dt.id_renstra = r.id
		WHERE r.uuid = ?
		LIMIT 1
	`

	var result domainrenstra.RenstraDefault

	err := r.db.WithContext(ctx).Raw(query, id).Scan(&result).Error
	if err != nil {
		return nil, err
	}

	if result.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &result, nil
}

var allowedSearchColumns = map[string]string{
	"tahun":         "r.tahun",
	"fakultas_unit": "vfu.nama_fak_prod_unit",
	"auditee":       "a0.name",
	"auditor1":      "a1.name",
	"auditor2":      "a2.name",
	"uuidauditee":   "a0.uuid",
	"uuidauditor1":  "a1.uuid",
	"uuidauditor2":  "a2.uuid",
	"kodeakses":     "r.kodeAkses",
	"jenjang":       "vfu.jenjang",
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
	scope string,
) ([]domainrenstra.RenstraDefault, int64, error) {

	var renstras []domainrenstra.RenstraDefault
	var total int64

	db := r.db.WithContext(ctx).Table("renstra r").
		Select(`
			r.id as ID,
			r.uuid as UUID,
			r.tahun as Tahun,

			r.fakultas_unit as FakultasUnitId,
			vfu.uuid as FakultasUnitUuid,
			vfu.nama_fak_prod_unit as FakultasUnit,

			r.periode_upload_mulai as PeriodeUploadMulai,
			r.periode_upload_akhir as PeriodeUploadAkhir,
			r.periode_assesment_dokumen_mulai as PeriodeAssesmentDokumenMulai,
			r.periode_assesment_dokumen_akhir as PeriodeAssesmentDokumenAkhir,
			r.periode_assesment_lapangan_mulai as PeriodeAssesmentLapanganMulai,
			r.periode_assesment_lapangan_akhir as PeriodeAssesmentLapanganAkhir,

			r.auditee as AuditeeId,
			r.auditor1 as Auditor1Id,
			r.auditor2 as Auditor2Id,

			r.kodeAkses as KodeAkses,
			r.catatan as Catatan1,
			r.catatan2 as Catatan2,

			a0.name as Auditee,
			a1.name as Auditor1,
			a2.name as Auditor2,

			a0.uuid as AuditeeUuid,
			a1.uuid as Auditor1Uuid,
			a2.uuid as Auditor2Uuid,

			vfu.jenjang as Jenjang,
			vfu.type as Type,
			vfu.fakultas as Fakultas,

			COUNT(DISTINCT rn.id) AS TotalRenstra,
			SUM(CASE WHEN rn.capaian IS NOT NULL THEN 1 ELSE 0 END) AS TotalRenstraAuditee,
			SUM(CASE WHEN rn.capaian_auditor IS NOT NULL THEN 1 ELSE 0 END) AS TotalRenstraAuditor,

			COUNT(DISTINCT dt.id) AS TotalDokumen,
			SUM(CASE WHEN dt.file IS NOT NULL THEN 1 ELSE 0 END) AS TotalDokumenAuditee,
			SUM(CASE WHEN dt.capaian_auditor IS NOT NULL THEN 1 ELSE 0 END) AS TotalDokumenAuditor
		`).
		Joins("JOIN v_fakultas_unit vfu ON r.fakultas_unit = vfu.id").
		Joins("LEFT JOIN users a0 ON r.auditee = a0.id").
		Joins("LEFT JOIN users a1 ON r.auditor1 = a1.id").
		Joins("LEFT JOIN renstra_nilai rn  ON rn.id_renstra = r.id").
		Joins("LEFT JOIN dokumen_tambahan dt ON dt.id_renstra = r.id")

	// -------------------------------
	// SEARCH FILTERS (ADVANCED)
	// -------------------------------
	var auditOrParts []string
	var auditParams []interface{}

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

			if scope == "audit" &&
				(col == "a0.uuid" || col == "a1.uuid" || col == "a2.uuid") {

				switch operator {
				case "eq":
					auditOrParts = append(auditOrParts, fmt.Sprintf("%s = ?", col))
					auditParams = append(auditParams, value)
				case "neq":
					auditOrParts = append(auditOrParts, fmt.Sprintf("%s <> ?", col))
					auditParams = append(auditParams, value)
				default:
					auditOrParts = append(auditOrParts, fmt.Sprintf("%s LIKE ?", col))
					auditParams = append(auditParams, "%"+value+"%")
				}
				continue
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

	if scope == "audit" && len(auditOrParts) > 0 {
		db = db.Where(
			"("+strings.Join(auditOrParts, " OR ")+")",
			auditParams...,
		)
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

		db = db.Where("("+strings.Join(orParts, " OR ")+")", params...)
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
	if err := db.Order("r.tahun DESC").Order("r.id DESC").Find(&renstras).Error; err != nil {
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

func (r *RenstraRepository) SetupUuid(ctx context.Context) error {
	const chunkSize = 500

	var ids []uint
	if err := r.db.WithContext(ctx).
		Model(&domainrenstra.Renstra{}).
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
			"UPDATE renstra SET uuid = %s WHERE id IN (?)",
			caseSQL,
		)

		if err := r.db.WithContext(ctx).
			Exec(query, args...).Error; err != nil {
			return err
		}
	}

	return nil
}
