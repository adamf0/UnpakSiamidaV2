package infrastructure

import (
	commondomainKts "UnpakSiamida/common/domain"
	commoninfra "UnpakSiamida/common/infrastructure"
	domainKts "UnpakSiamida/modules/kts/domain"
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type KtsRepository struct {
	db  *gorm.DB
	uow *commoninfra.UnitOfWork
}

func NewKtsRepository(db *gorm.DB) domainKts.IKtsRepository {
	return &KtsRepository{db: db, uow: commoninfra.NewUnitOfWork(db)}
}

// ------------------------
// GET BY UUID
// ------------------------
func (r *KtsRepository) GetByUuid(ctx context.Context, uid uuid.UUID) (*domainKts.Kts, error) {
	var Kts domainKts.Kts

	err := r.db.WithContext(ctx).
		Where("uuid = ?", uid).
		First(&Kts).Error

	// if errors.Is(err, gorm.ErrRecordNotFound) {
	// 	return nil, nil
	// }

	if err != nil {
		return nil, err
	}

	return &Kts, nil
}

// ------------------------
// GET DEFAULT BY UUID
// ------------------------
func (r *KtsRepository) GetDefaultByUuid(
	ctx context.Context,
	id uuid.UUID,
) (*domainKts.KtsDefault, error) {

	// Ambil hanya kolom yang benar-benar ada di struct KtsDefault
	var rowData domainKts.KtsDefault

	err := r.db.WithContext(ctx).
		Table("kts_renstra k").
		Select(`
		k.id AS ID,
		k.uuid AS UUID,

		r.id AS RenstraId,
		k.id_renstra_nilai AS RenstraNilai,
		k.id_dokumen_tambahan AS DokumenTambahan,
		k.status AS Status,

		dt.id_template_dokumen_tambahan AS TemplateDokumen,
		tdt.pertanyaan AS Pertanyaan,
		jf.id AS JenisFileId,
		jf.nama AS JenisFile,

		rn.template_renstra AS TemplateRenstra,
		sr.id AS StandarId,
		sr.nama AS Standar,
		mir.id AS IndikatorId,
		mir.indikator AS Indikator,

		r.tahun AS Tahun,
		r.fakultas_unit AS IdTarget,
		CASE 
			WHEN fu.type COLLATE utf8mb4_unicode_ci = 'prodi'
			THEN CONCAT(fu.nama_fak_prod_unit, ' (', fu.jenjang, ')')
			ELSE fu.nama_fak_prod_unit
		END AS Target,

		k.nomor_laporan as NomorLaporan,
		k.tanggal_laporan as TanggalLaporan,
		
		k.auditor as Auditor,
		u2.name as NamaAuditor,
		u2.uuid as AuditorUuid,
		r.auditee as Auditee,
		u1.name as NamaAuditee,
		u1.uuid as AuditeeUuid,

		k.uraian_ketidaksesuaian_p as KetidaksesuaianP,
		k.uraian_ketidaksesuaian_l as KetidaksesuaianL,
		k.uraian_ketidaksesuaian_o as KetidaksesuaianO,
		k.uraian_ketidaksesuaian_r as KetidaksesuaianR,
		k.referensi as Referensi,
		k.hasil_temuan as HasilTemuan,
		k.akar_masalah as AkarMasalah,
		k.tindakan_koreksi as TindakanKoreksi,
		k.acc_auditor as AccAuditor,

		k.status_acc_auditee as StatusAccAuditee,
		k.acc_auditee as AccAuditee,
		k.keterangan_tolak_auditee as KeteranganTolak,
		k.tindakan_perbaikan as TindakanPerbaikan,

		k.tanggal_penyelesaian as TanggalPenyelesaian,

		k.tinjauan_tindakan_perbaikan as TinjauanTindakanPerbaikan,
		k.tanggal_closing_auditee as TanggalClosing,
		k.acc_auditor_final as AccFinal,

		k.tanggal_closing as TanggalClosingFinal,
		k.wmm_upmf_upmps as WmmUpmfUpmps,
		k.closingBy as ClosingBy
	`).
		Joins(`LEFT JOIN renstra_nilai rn ON k.id_renstra_nilai = rn.id`).
		Joins(`LEFT JOIN dokumen_tambahan dt ON k.id_dokumen_tambahan = dt.id`).
		Joins(`JOIN renstra r ON r.id = COALESCE(rn.id_renstra, dt.id_renstra)`).
		Joins(`JOIN v_fakultas_unit fu ON r.fakultas_unit = fu.id`).
		Joins(`LEFT JOIN template_renstra tr ON rn.template_renstra = tr.id`).
		Joins(`LEFT JOIN master_indikator_renstra mir ON tr.indikator = mir.id`).
		Joins(`LEFT JOIN master_standar_renstra sr ON mir.id_master_standar = sr.id`).
		Joins(`LEFT JOIN template_dokumen_tambahan tdt ON dt.id_template_dokumen_tambahan = tdt.id`).
		Joins(`LEFT JOIN jenis_file_renstra jf ON tdt.jenis_file = jf.id`).
		Joins(`LEFT JOIN users u1 ON r.auditee = u1.id`).
		Joins(`LEFT JOIN users u2 ON k.auditor = u2.id`).
		Where("k.uuid = ?", id).
		Order("k.id_dokumen_tambahan DESC").
		Take(&rowData).Error // Take otomatis LIMIT 1

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &rowData, nil
}

var allowedSearchColumns = map[string]string{
	// key:param -> db column
	"status":     "k.status",
	"pertanyaan": "tdt.pertanyaan",
	"jenisfile":  "jf.nama",
	"standar":    "sr.nama",
	"indikator":  "mir.indikator",
	"target":     "CASE WHEN fu.type COLLATE utf8mb4_unicode_ci = 'prodi' THEN CONCAT(fu.nama_fak_prod_unit, ' (', fu.jenjang, ')') ELSE fu.nama_fak_prod_unit END",
	"targetuuid": "fu.uuid",
	"tahun":      "r.tahun",
}

// ------------------------
// GET ALL
// ------------------------
func (r *KtsRepository) GetAll(
	ctx context.Context,
	search string,
	searchFilters []commondomainKts.SearchFilter,
	page, limit *int,
) ([]domainKts.KtsDefault, int64, error) {

	var rows = make([]domainKts.KtsDefault, 0)
	var total int64

	db := r.db.WithContext(ctx).
		Table("kts_renstra k").
		Select(`
			k.id AS ID,
			k.uuid AS UUID,

			r.id AS RenstraId,
			k.id_renstra_nilai AS RenstraNilai,
			k.id_dokumen_tambahan AS DokumenTambahan,
			k.status AS Status,

			dt.id_template_dokumen_tambahan AS TemplateDokumen,
			tdt.pertanyaan AS Pertanyaan,
			jf.id AS JenisFileId,
			jf.nama AS JenisFile,

			rn.template_renstra AS TemplateRenstra,
			sr.id AS StandarId,
			sr.nama AS Standar,
			mir.id AS IndikatorId,
			mir.indikator AS Indikator,

			r.tahun AS Tahun,
			r.fakultas_unit AS IdTarget,
			CASE 
				WHEN fu.type COLLATE utf8mb4_unicode_ci = 'prodi'
				THEN CONCAT(fu.nama_fak_prod_unit, ' (', fu.jenjang, ')')
				ELSE fu.nama_fak_prod_unit
			END AS Target,

			k.nomor_laporan as NomorLaporan,
			k.tanggal_laporan as TanggalLaporan,
			
			k.auditor as Auditor,
			u2.name as NamaAuditor,
			u2.uuid as AuditorUuid,
			r.auditee as Auditee,
			u1.name as NamaAuditee,
			u1.uuid as AuditeeUuid,

			k.uraian_ketidaksesuaian_p as KetidaksesuaianP,
			k.uraian_ketidaksesuaian_l as KetidaksesuaianL,
			k.uraian_ketidaksesuaian_o as KetidaksesuaianO,
			k.uraian_ketidaksesuaian_r as KetidaksesuaianR,
			k.referensi as Referensi,
			k.hasil_temuan as HasilTemuan,
			k.akar_masalah as AkarMasalah,
			k.tindakan_koreksi as TindakanKoreksi,
			k.acc_auditor as AccAuditor,

			k.status_acc_auditee as StatusAccAuditee,
			k.acc_auditee as AccAuditee,
			k.keterangan_tolak_auditee as KeteranganTolak,
			k.tindakan_perbaikan as TindakanPerbaikan,

			k.tanggal_penyelesaian as TanggalPenyelesaian,

			k.tinjauan_tindakan_perbaikan as TinjauanTindakanPerbaikan,
			k.tanggal_closing_auditee as TanggalClosing,
			k.acc_auditor_final as AccFinal,

			k.tanggal_closing as TanggalClosingFinal,
			k.wmm_upmf_upmps as WmmUpmfUpmps,
			k.closingBy as ClosingBy
		`).
		Joins(`
			LEFT JOIN renstra_nilai rn 
				ON k.id_renstra_nilai = rn.id
		`).
		Joins(`
			LEFT JOIN dokumen_tambahan dt 
				ON k.id_dokumen_tambahan = dt.id
		`).
		Joins(`
			JOIN renstra r 
				ON r.id = COALESCE(rn.id_renstra, dt.id_renstra)
		`).
		Joins(`
			JOIN v_fakultas_unit fu 
				ON r.fakultas_unit = fu.id
		`).
		Joins(`
			LEFT JOIN template_renstra tr 
				ON rn.template_renstra = tr.id
		`).
		Joins(`
			LEFT JOIN master_indikator_renstra mir 
				ON tr.indikator = mir.id
		`).
		Joins(`
			LEFT JOIN master_standar_renstra sr 
				ON mir.id_master_standar = sr.id
		`).
		Joins(`
			LEFT JOIN template_dokumen_tambahan tdt 
				ON dt.id_template_dokumen_tambahan = tdt.id
		`).
		Joins(`
			LEFT JOIN jenis_file_renstra jf 
				ON tdt.jenis_file = jf.id
		`).
		Joins(`LEFT JOIN users u1 ON r.auditee = u1.id`).
		Joins(`LEFT JOIN users u2 ON k.auditor = u2.id`)

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
		if val == "" {
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
	// COUNT
	// -----------------------------------
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// -----------------------------------
	// ORDER + PAGINATION
	// -----------------------------------
	db = db.Order("k.id_dokumen_tambahan DESC")

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
func (r *KtsRepository) Create(ctx context.Context, kts *domainKts.Kts) error {
	return r.db.WithContext(ctx).Create(kts).Error
}

// ------------------------
// UPDATE
// ------------------------
func (r *KtsRepository) Update(ctx context.Context, kts *domainKts.Kts) error {
	// return r.db.WithContext(ctx).Save(kts).Error
	if err := r.db.WithContext(ctx).Save(kts).Error; err != nil {
		return err
	}

	return r.uow.Save(&kts.Entity)
}

// ------------------------
// DELETE (by UUID)
// ------------------------
func (r *KtsRepository) Delete(ctx context.Context, uid uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("uuid = ?", uid).
		Delete(&domainKts.Kts{}).Error
}

func (r *KtsRepository) SetupUuid(ctx context.Context) error {
	const chunkSize = 500

	var ids []uint
	if err := r.db.WithContext(ctx).
		Model(&domainKts.Kts{}).
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

func (r *KtsRepository) WithTx(
	ctx context.Context,
	fn func(txRepo domainKts.IKtsRepositoryTx) error,
) error {

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txRepo := &KtsRepository{
			db:  tx,
			uow: commoninfra.NewUnitOfWork(tx),
		}
		return fn(txRepo)
	})
}

func (r *KtsRepository) InsertOutbox(
	ctx context.Context,
	msg *commoninfra.OutboxMessage,
) error {
	return r.db.WithContext(ctx).Create(msg).Error
}
