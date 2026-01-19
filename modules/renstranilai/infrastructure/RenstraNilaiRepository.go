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
func (r *RenstraNilaiRepository) GetDefaultByUuid( //[pr] ini salah, ini tanpa kts
	ctx context.Context,
	id uuid.UUID,
) (*domainrenstranilai.RenstraNilaiDefault, error) {

	query := `
        SELECT 
			k.id AS ID,
			k.uuid AS UUID,
			r.id AS RenstraId,
			k.id_renstra_nilai AS RenstraNilai,
			k.id_dokumen_tambahan AS DokumenTambahan,
			k.status AS Status,
			r.tahun AS Tahun,
			r.fakultas_unit AS IdTarget,
			CASE 
				WHEN fu.type = 'prodi' 
					THEN CONCAT(fu.nama_fak_prod_unit, ' (', fu.jenjang, ')')
				ELSE fu.nama_fak_prod_unit
			END AS Target,
			dt.id_template_dokumen_tambahan AS TemplateDokumen,
			tdt.pertanyaan AS Pertanyaan,
			jf.nama AS JenisFile,
			rn.template_renstra AS TemplateRenstra,
			sr.nama AS Standar,
			mir.indikator AS Indikator,

			k.nomor_laporan                    AS NomorLaporan,
			k.tanggal_laporan                  AS TanggalLaporan,
			k.auditor                          AS Auditor,

			k.uraian_ketidaksesuaian_p         AS KetidaksesuaianP,
			k.uraian_ketidaksesuaian_l         AS KetidaksesuaianL,
			k.uraian_ketidaksesuaian_o         AS KetidaksesuaianO,
			k.uraian_ketidaksesuaian_r         AS KetidaksesuaianR,

			k.akar_masalah                     AS AkarMasalah,
			k.tindakan_koreksi                 AS TindakanKoreksi,
			k.acc_auditor                      AS AccAuditor,

			k.status_acc_auditee               AS StatusAccAuditee,
			k.acc_auditee                      AS AccAuditee,
			k.keterangan_tolak_auditee         AS KeteranganTolak,
			k.tindakan_perbaikan               AS TindakanPerbaikan,

			k.tanggal_penyelesaian             AS TanggalPenyelesaian,
			k.tinjauan_tindakan_perbaikan      AS TinjauanTindakanPerbaikan,
			k.tanggal_closing_auditee          AS TanggalClosing,

			k.acc_auditor_final                AS AccFinal,
			k.tanggal_closing                  AS TanggalClosingFinal,
			k.wmm_upmf_upmps                   AS WmmUpmfUpmps,
			k.closingBy                        AS ClosingBy
		FROM kts_renstra k
		LEFT JOIN renstra_nilai rn 
			ON k.id_renstra_nilai = rn.id
		LEFT JOIN dokumen_tambahan dt 
			ON k.id_dokumen_tambahan = dt.id
		JOIN renstra r 
			ON r.id = COALESCE(rn.id_renstra, dt.id_renstra)
		JOIN v_fakultas_unit fu 
			ON r.fakultas_unit = fu.id
		LEFT JOIN template_renstra tr 
			ON rn.template_renstra = tr.id
		LEFT JOIN master_indikator_renstra mir 
			ON tr.indikator = mir.id
		LEFT JOIN master_standar_renstra sr 
			ON mir.id_master_standar = sr.id
		LEFT JOIN template_dokumen_tambahan tdt 
			ON dt.id_template_dokumen_tambahan = tdt.id
		LEFT JOIN jenis_file_renstra jf 
			ON tdt.jenis_file = jf.id
        WHERE k.uuid = ?
        LIMIT 1
    `

	var rowData domainrenstranilai.RenstraNilaiDefault

	res := r.db.WithContext(ctx).Raw(query, id).Scan(&rowData)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, res.Error
	}

	if rowData.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &rowData, nil
}

var allowedSearchColumns = map[string]string{
	// key:param -> db column
	// "uuidrenstra":      "r.uuid",
}

func (r *RenstraNilaiRepository) GetAll(
	ctx context.Context,
	search string,
	searchFilters []commondomain.SearchFilter,
	page, limit *int,
) ([]domainrenstranilai.RenstraNilaiDefault, int64, error) {

	var (
		result     []domainrenstranilai.RenstraNilaiDefault
		total      int64
		conditions []string
		args       []interface{}
	)

	// =====================================================
	// BASE FROM + JOIN (WAJIB PAKAI INI)
	// =====================================================
	baseFrom := `
		FROM kts_renstra k
		LEFT JOIN renstra_nilai rn 
			ON k.id_renstra_nilai = rn.id
		LEFT JOIN dokumen_tambahan dt 
			ON k.id_dokumen_tambahan = dt.id
		JOIN renstra r 
			ON r.id = COALESCE(rn.id_renstra, dt.id_renstra)
		JOIN v_fakultas_unit fu 
			ON r.fakultas_unit = fu.id
		LEFT JOIN template_renstra tr 
			ON rn.template_renstra = tr.id
		LEFT JOIN master_indikator_renstra mir 
			ON tr.indikator = mir.id
		LEFT JOIN master_standar_renstra sr 
			ON mir.id_master_standar = sr.id
		LEFT JOIN template_dokumen_tambahan tdt 
			ON dt.id_template_dokumen_tambahan = tdt.id
		LEFT JOIN jenis_file_renstra jf 
			ON tdt.jenis_file = jf.id
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

	orderBy := " ORDER BY r.tahun DESC, k.id DESC"

	// =====================================================
	// SELECT QUERY (INI YANG KAMU MINTA)
	// =====================================================
	selectQuery := `
		SELECT
			k.id AS ID,
			k.uuid AS UUID,
			r.id AS RenstraId,
			k.id_renstra_nilai AS RenstraNilai,
			k.id_dokumen_tambahan AS DokumenTambahan,
			k.status AS Status,
			r.tahun AS Tahun,
			r.fakultas_unit AS IdTarget,
			CASE 
				WHEN fu.type = 'prodi' 
					THEN CONCAT(fu.nama_fak_prod_unit, ' (', fu.jenjang, ')')
				ELSE fu.nama_fak_prod_unit
			END AS Target,
			dt.id_template_dokumen_tambahan AS TemplateDokumen,
			tdt.pertanyaan AS Pertanyaan,
			jf.nama AS JenisFile,
			rn.template_renstra AS TemplateRenstra,
			sr.nama AS Standar,
			mir.indikator AS Indikator,

			k.nomor_laporan                    AS NomorLaporan,
			k.tanggal_laporan                  AS TanggalLaporan,
			k.auditor                          AS Auditor,

			k.uraian_ketidaksesuaian_p         AS KetidaksesuaianP,
			k.uraian_ketidaksesuaian_l         AS KetidaksesuaianL,
			k.uraian_ketidaksesuaian_o         AS KetidaksesuaianO,
			k.uraian_ketidaksesuaian_r         AS KetidaksesuaianR,

			k.akar_masalah                     AS AkarMasalah,
			k.tindakan_koreksi                 AS TindakanKoreksi,
			k.acc_auditor                      AS AccAuditor,

			k.status_acc_auditee               AS StatusAccAuditee,
			k.acc_auditee                      AS AccAuditee,
			k.keterangan_tolak_auditee         AS KeteranganTolak,
			k.tindakan_perbaikan               AS TindakanPerbaikan,

			k.tanggal_penyelesaian             AS TanggalPenyelesaian,
			k.tinjauan_tindakan_perbaikan      AS TinjauanTindakanPerbaikan,
			k.tanggal_closing_auditee          AS TanggalClosing,

			k.acc_auditor_final                AS AccFinal,
			k.tanggal_closing                  AS TanggalClosingFinal,
			k.wmm_upmf_upmps                   AS WmmUpmfUpmps,
			k.closingBy                        AS ClosingBy
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
