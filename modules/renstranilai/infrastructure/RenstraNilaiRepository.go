package infrastructure

import (
	commondomain "UnpakSiamida/common/domain"
	domainrenstranilai "UnpakSiamida/modules/renstranilai/domain"
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RenstraNilaiRepository struct {
	db *gorm.DB
}

func NewRenstraNilaiRepository(db *gorm.DB) domainrenstranilai.IRenstraNilaiRepository {
	return &RenstraNilaiRepository{db: db}
}

// ------------------------
// GET BY UUID
// ------------------------
func (r *RenstraNilaiRepository) GetByUuid(ctx context.Context, uid uuid.UUID) (*domainrenstranilai.RenstraNilai, error) {
	var renstranilai domainrenstranilai.RenstraNilai

	err := r.db.WithContext(ctx).
		Where("uuid = ?", uid).
		First(&renstranilai).Error

	// if errors.Is(err, gorm.ErrRecordNotFound) {
	// 	return nil, nil
	// }

	if err != nil {
		return nil, err
	}

	return &renstranilai, nil
}

// ------------------------
// GET DEFAULT BY UUID
// ------------------------
func (r *RenstraNilaiRepository) GetDefaultByUuid(
	ctx context.Context,
	id uuid.UUID,
) (*domainrenstranilai.RenstraNilaiDefault, error) {

	var rowData domainrenstranilai.RenstraNilaiDefault

	err := r.db.WithContext(ctx).
		Table("renstra_nilai rn").
		Select(`
		r.id as RenstraId,
		r.uuid as RenstraUUID,
		r.tahun as TahunRenstra,
		
		rn.id as ID,
		rn.uuid as UUID,
		fu.nama_fak_prod_unit AS TargetAudit,
		fu.jenjang AS Jenjang,
		fu.fakultas AS Fakultas,
		fu.type AS Type,
		
		sr.id as StandarId,
		sr.uuid as StandarUUID,
		sr.nama as NamaStandar,
		
		mir.id as IndikatorId,
		mir.uuid as IndikatorUUID,
		mir.indikator as NamaIndikator,
		mir.tahun as TahunIndikator,
		mir.tipe_target as TipeTarget,
		mir.operator as Operator,
		
		tr.id as TemplateRenstraId,
		tr.uuid as TemplateRenstraUUID,
		tr.satuan as Satuan,
		tr.target as Target,
		tr.target_min as TargetMin,
		tr.target_max as TargetMax,
		tr.tugas as TugasTemplate,
		tr.tahun as TahunTemplate,
		tr.pertanyaan as IsPertanyaan,
		
		rn.tugas as Tugas,
		rn.capaian as CapaianAuditee,
		rn.catatan as CatatanAuditee,
		rn.link_bukti as LinkBukti,
		rn.capaian_auditor as CapaianAuditor,
		rn.catatan_auditor as CatatanAuditor
	`).
		Joins("JOIN renstra r ON rn.id_renstra = r.id").
		Joins("JOIN template_renstra tr ON rn.template_renstra = tr.id").
		Joins("JOIN master_indikator_renstra mir ON tr.indikator = mir.id").
		Joins("JOIN master_standar_renstra sr ON mir.id_master_standar = sr.id").
		Joins("JOIN v_fakultas_unit fu ON r.fakultas_unit = fu.id").
		Where("rn.uuid = ?", id).
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
	"uuidrenstra":   "r.uuid",
	"uuidindikator": "mir.uuid",
}

func (r *RenstraNilaiRepository) GetAll(
	ctx context.Context,
	search string,
	searchFilters []commondomain.SearchFilter,
	page, limit *int,
) ([]domainrenstranilai.RenstraNilaiDefault, int64, error) {

	var (
		result     = make([]domainrenstranilai.RenstraNilaiDefault, 0)
		total      int64
		conditions []string
		args       []interface{}
	)

	// =====================================================
	// BASE FROM + JOIN (WAJIB PAKAI INI)
	// =====================================================
	//COALESCE(rn.id_renstra, dt.id_renstra)
	baseFrom := `
		FROM
			renstra_nilai rn
		JOIN renstra r ON
			rn.id_renstra = r.id
		JOIN template_renstra tr ON
			rn.template_renstra = tr.id
		JOIN master_indikator_renstra mir ON
			tr.indikator = mir.id
		join master_standar_renstra sr on mir.id_master_standar = sr.id
		JOIN v_fakultas_unit fu ON
			r.fakultas_unit = fu.id
	`

	// =====================================================
	// ADVANCED FILTERS
	// =====================================================
	for _, f := range searchFilters {
		field := strings.TrimSpace(strings.ToLower(f.Field))
		operator := strings.TrimSpace(strings.ToLower(f.Operator))

		col, ok := allowedSearchColumns[field] // harus sudah pakai alias tr./i./fu.
		if !ok {
			continue
		}

		value := ""
		if f.Value != nil {
			value = strings.TrimSpace(*f.Value)
		}

		switch operator {
		case "eq":
			conditions = append(conditions, fmt.Sprintf("%s = ?", col))
			args = append(args, value)

		case "neq":
			conditions = append(conditions, fmt.Sprintf("%s <> ?", col))
			args = append(args, value)

		case "like":
			conditions = append(conditions, fmt.Sprintf("%s LIKE ?", col))
			args = append(args, "%"+value+"%")

		case "gt":
			conditions = append(conditions, fmt.Sprintf("%s > ?", col))
			args = append(args, value)

		case "gte":
			conditions = append(conditions, fmt.Sprintf("%s >= ?", col))
			args = append(args, value)

		case "lt":
			conditions = append(conditions, fmt.Sprintf("%s < ?", col))
			args = append(args, value)

		case "lte":
			conditions = append(conditions, fmt.Sprintf("%s <= ?", col))
			args = append(args, value)

		case "in":
			values := strings.Split(value, ",")
			if len(values) > 0 {
				placeholders := strings.TrimRight(strings.Repeat("?,", len(values)), ",")
				conditions = append(
					conditions,
					fmt.Sprintf("%s IN (%s)", col, placeholders),
				)
				for _, v := range values {
					args = append(args, strings.TrimSpace(v))
				}
			}

		default:
			conditions = append(conditions, fmt.Sprintf("%s LIKE ?", col))
			args = append(args, "%"+value+"%")
		}
	}

	// =====================================================
	// GLOBAL SEARCH
	// =====================================================
	if strings.TrimSpace(search) != "" {
		var orParts []string
		like := "%" + search + "%"

		for _, col := range allowedSearchColumns {
			orParts = append(orParts, fmt.Sprintf("%s LIKE ?", col))
			args = append(args, like)
		}

		conditions = append(
			conditions,
			"("+strings.Join(orParts, " OR ")+")",
		)
	}

	// =====================================================
	// WHERE CLAUSE
	// =====================================================
	whereClause := ""
	if len(conditions) > 0 {
		whereClause = " WHERE " + strings.Join(conditions, " AND ")
	}

	// =====================================================
	// COUNT QUERY (HARUS PAKAI JOIN JUGA)
	// =====================================================
	countQuery := `
		SELECT COUNT(DISTINCT rn.id)
	` + baseFrom + whereClause

	if err := r.db.WithContext(ctx).
		Raw(countQuery, args...).
		Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	// =====================================================
	// PAGINATION
	// =====================================================
	pagination := ""
	if page != nil && limit != nil && *limit > 0 {
		p := *page
		l := *limit
		if p < 1 {
			p = 1
		}

		offset := (p - 1) * l
		pagination = " LIMIT ? OFFSET ?"
		args = append(args, l, offset)
	}

	orderBy := " ORDER BY r.tahun DESC, r.id DESC"

	// =====================================================
	// SELECT QUERY (INI YANG KAMU MINTA)
	// =====================================================
	selectQuery := `
		SELECT
			r.id as RenstraId,
			r.uuid as RenstraUUID,
			r.tahun as TahunRenstra,
			
			rn.id as ID,
			rn.uuid as UUID,
			fu.nama_fak_prod_unit AS TargetAudit,
			fu.jenjang AS Jenjang,
			fu.fakultas AS Fakultas,
			fu.type AS Type,
			
			sr.id as StandarId,
			sr.uuid as StandarUUID,
			sr.nama as NamaStandar,
			
			mir.id as IndikatorId,
			mir.uuid as IndikatorUUID,
			mir.indikator as NamaIndikator,
			mir.tahun as TahunIndikator,
			mir.tipe_target as TipeTarget,
			mir.operator as Operator,
			
			tr.id as TemplateRenstraId,
            tr.uuid as TemplateRenstraUUID,
			tr.satuan as Satuan,
			tr.target as Target,
			tr.target_min as TargetMin,
			tr.target_max as TargetMax,
			tr.tugas as TugasTemplate,
			tr.tahun as TahunTemplate,
			tr.pertanyaan as IsPertanyaan,
			
			rn.tugas as Tugas,
			rn.capaian as CapaianAuditee,
			rn.catatan as CatatanAuditee,
			rn.link_bukti as LinkBukti,
			rn.capaian_auditor as CapaianAuditor,
			rn.catatan_auditor as CatatanAuditor
	` + baseFrom + whereClause + orderBy + pagination

	if err := r.db.WithContext(ctx).
		Raw(selectQuery, args...).
		Scan(&result).Error; err != nil {
		return nil, 0, err
	}

	return result, total, nil
}

// ------------------------
// UPDATE
// ------------------------
func (r *RenstraNilaiRepository) Update(ctx context.Context, renstranilai *domainrenstranilai.RenstraNilai) error {
	return r.db.WithContext(ctx).Save(renstranilai).Error
}

// ------------------------
// DELETE (by UUID)
// ------------------------
func (r *RenstraNilaiRepository) Delete(ctx context.Context, uid uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("uuid = ?", uid).
		Delete(&domainrenstranilai.RenstraNilai{}).Error
}

func (r *RenstraNilaiRepository) SetupUuid(ctx context.Context) error {
	const chunkSize = 500

	var ids []uint
	if err := r.db.WithContext(ctx).
		Model(&domainrenstranilai.RenstraNilai{}).
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
			"UPDATE kts_renstra SET uuid = %s WHERE id IN (?)",
			caseSQL,
		)

		if err := r.db.WithContext(ctx).
			Exec(query, args...).Error; err != nil {
			return err
		}
	}

	return nil
}
