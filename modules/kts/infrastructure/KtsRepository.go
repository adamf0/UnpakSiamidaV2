package infrastructure

import (
	commondomainKts "UnpakSiamida/common/domain"
	commoninfra "UnpakSiamida/common/infrastructure"
	domainKts "UnpakSiamida/modules/kts/domain"
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

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

func (r *KtsRepository) ensureUuidColumn(ctx context.Context, tableName string) error {
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

func (r *KtsRepository) SetupUuid(ctx context.Context) error {
	if err := r.ensureUuidColumn(ctx, "kts_renstra"); err != nil {
		return err
	}

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

func (r *KtsRepository) GetTemuanKts(
	ctx context.Context,
	renstraUuid string,
	tugas string,
	tahun string,
) ([]domainKts.TemuanKts, error) {
	if strings.TrimSpace(tahun) != "" {
		if tNum, err := strconv.Atoi(strings.TrimSpace(tahun)); err == nil && tNum <= 2023 {
			return nil, commondomainKts.FailureError("Kts.MigrasiFitur", "masih migrasi fitur")
		}
	}

	results := make([]domainKts.TemuanKts, 0)

	normTugas := strings.ToLower(strings.TrimSpace(tugas))
	if normTugas == "utama" {
		normTugas = "auditor1"
	} else if normTugas == "tambahan" {
		normTugas = "auditor2"
	}

	// 1. Query Renstra Temuan (Auditor 1, Auditor 2, atau Semua)
	if normTugas == "" || normTugas == "auditor1" || normTugas == "auditor2" {
		var renstraRows []domainKts.TemuanKts

		renstraSQL := `
			SELECT 
				kts.id AS ID,
				kts.uuid AS UUID,
				rn.id AS RenstraNilaiID,
				rn.uuid AS RenstraNilaiUUID,
				NULL AS DokumenTambahanID,
				NULL AS DokumenTambahanUUID,
				'renstra' AS Source,
				r.auditor1 AS Auditor1,
				u1.name AS NamaAuditor1,
				r.auditor2 AS Auditor2,
				u2.name AS NamaAuditor2,
				r.auditee AS Auditee,
				u3.name AS NamaAuditee,
				(CASE WHEN kts.id IS NULL THEN 'belum_dibuat' ELSE COALESCE(kts.status, 'draf') END) AS Status,
				COALESCE(venr.nama_standar, sr.nama, '') AS NamaStandar,
				COALESCE(venr.indikator, mir.indikator, '') AS Indikator,
				COALESCE(venr.klasifikasi, tr.klasifikasi, '') AS Klasifikasi,
				COALESCE(venr.capaian, rn.capaian, '') AS Capaian,
				COALESCE(venr.capaian_auditor, rn.capaian_auditor, '') AS CapaianAuditor,
				COALESCE(kts.tindakan_perbaikan, '') AS TindakanPerbaikan,
				(CASE WHEN kts.status = 'tutup_kts' THEN 'Close' ELSE '' END) AS Keterangan,
				COALESCE(kts.keterangan_tolak_auditee, '') AS KeteranganTolakAuditee,
				COALESCE(kts.nomor_laporan, '') AS NomorLaporan,
				COALESCE(DATE_FORMAT(kts.tanggal_laporan, '%Y-%m-%d'), '') AS TanggalLaporan,
				COALESCE(kts.referensi, '') AS Referensi,
				COALESCE(kts.hasil_temuan, '') AS HasilTemuan,
				COALESCE(kts.uraian_ketidaksesuaian_p, '') AS KetidaksesuaianP,
				COALESCE(kts.uraian_ketidaksesuaian_l, '') AS KetidaksesuaianL,
				COALESCE(kts.uraian_ketidaksesuaian_o, '') AS KetidaksesuaianO,
				COALESCE(kts.uraian_ketidaksesuaian_r, '') AS KetidaksesuaianR,
				TRIM(CONCAT_WS(' ', NULLIF(kts.uraian_ketidaksesuaian_p, ''), NULLIF(kts.uraian_ketidaksesuaian_l, ''), NULLIF(kts.uraian_ketidaksesuaian_o, ''), NULLIF(kts.uraian_ketidaksesuaian_r, ''))) AS PLOR,
				COALESCE(kts.akar_masalah, '') AS AkarMasalah,
				COALESCE(kts.tindakan_koreksi, '') AS TindakanKoreksi,
				COALESCE(DATE_FORMAT(kts.tanggal_penyelesaian, '%Y-%m-%d'), '') AS TanggalPenyelesaian,
				COALESCE(kts.tinjauan_tindakan_perbaikan, '') AS TinjauanTindakanPerbaikan,
				COALESCE(DATE_FORMAT(kts.tanggal_closing_auditee, '%Y-%m-%d'), DATE_FORMAT(kts.tanggal_closing, '%Y-%m-%d'), '') AS TanggalClosing,
				LOWER(TRIM(COALESCE(venr.tugas, tr.tugas, 'auditor1'))) AS Tugas,
				r.tahun AS Tahun,
				r.fakultas_unit AS FakultasUnit,
				CASE 
					WHEN fu.type COLLATE utf8mb4_unicode_ci = 'prodi'
					THEN CONCAT(fu.nama_fak_prod_unit, ' (', fu.jenjang, ')')
					ELSE fu.nama_fak_prod_unit
				END AS Target,
				fu.uuid AS TargetUUID,
				r.uuid AS RenstraUUID,
				fu.jenjang AS Jenjang,
				fu.type AS Type,
				fu.fakultas AS Fakultas,
				COALESCE(CAST(venr.target_seharusnya AS CHAR), CAST(tr.target AS CHAR), '') AS TargetIndikator,
				COALESCE(CAST(venr.target_min AS CHAR), CAST(tr.target_min AS CHAR), '') AS TargetMin,
				COALESCE(CAST(venr.target_max AS CHAR), CAST(tr.target_max AS CHAR), '') AS TargetMax,
				(CASE 
					WHEN COALESCE(venr.target_min, tr.target_min) IS NOT NULL AND COALESCE(venr.target_max, tr.target_max) IS NOT NULL 
					THEN CONCAT(COALESCE(venr.target_min, tr.target_min), ' - ', COALESCE(venr.target_max, tr.target_max))
					WHEN COALESCE(venr.target_min, tr.target_min) IS NOT NULL 
					THEN CAST(COALESCE(venr.target_min, tr.target_min) AS CHAR)
					WHEN COALESCE(venr.target_max, tr.target_max) IS NOT NULL 
					THEN CAST(COALESCE(venr.target_max, tr.target_max) AS CHAR)
					ELSE NULL
				END) AS RangeTarget,
				COALESCE(venr.tipe_target, mir.tipe_target, '') AS TipeTarget,
				COALESCE(venr.operator, mir.operator, '') AS Operator,
				COALESCE(venr.satuan, tr.satuan, '') AS Satuan
			FROM renstra_nilai rn
			JOIN renstra r ON rn.id_renstra = r.id
			JOIN view_evaluasi_nilai_renstra venr ON rn.id = venr.id
			JOIN template_renstra tr ON rn.template_renstra = tr.id
			LEFT JOIN master_indikator_renstra mir ON tr.indikator = mir.id
			LEFT JOIN master_standar_renstra sr ON mir.id_master_standar = sr.id
			LEFT JOIN kts_renstra kts ON rn.id = kts.id_renstra_nilai
			LEFT JOIN users u1 ON r.auditor1 = u1.id
			LEFT JOIN users u2 ON r.auditor2 = u2.id
			LEFT JOIN users u3 ON r.auditee = u3.id
			LEFT JOIN v_fakultas_unit fu ON r.fakultas_unit = fu.id
			WHERE LOWER(TRIM(venr.klasifikasi)) != 'ok'
			  AND venr.klasifikasi IS NOT NULL
			  AND TRIM(venr.klasifikasi) != ''
		`
		var args []interface{}
		if strings.TrimSpace(renstraUuid) != "" && strings.TrimSpace(renstraUuid) != "all" {
			renstraSQL += " AND (r.uuid = ? OR CAST(r.id AS CHAR) = ?)"
			args = append(args, renstraUuid, renstraUuid)
		}
		if normTugas != "" {
			renstraSQL += " AND LOWER(TRIM(COALESCE(venr.tugas, tr.tugas))) = ?"
			args = append(args, normTugas)
		}
		if tahun != "" {
			renstraSQL += " AND r.tahun = ?"
			args = append(args, tahun)
		}
		renstraSQL += " ORDER BY rn.id ASC"

		if err := r.db.WithContext(ctx).Raw(renstraSQL, args...).Scan(&renstraRows).Error; err != nil {
			return nil, err
		}
		results = append(results, renstraRows...)
	}

	// 2. Query Dokumen Tambahan Temuan
	if normTugas == "" || normTugas == "auditor2" || normTugas == "auditor1" {
		var tambahanRows []domainKts.TemuanKts

		tambahanSQL := `
			SELECT 
				kts.id AS ID,
				kts.uuid AS UUID,
				NULL AS RenstraNilaiID,
				NULL AS RenstraNilaiUUID,
				dt.id AS DokumenTambahanID,
				dt.uuid AS DokumenTambahanUUID,
				'dokumen_tambahan' AS Source,
				r.auditor1 AS Auditor1,
				u1.name AS NamaAuditor1,
				r.auditor2 AS Auditor2,
				u2.name AS NamaAuditor2,
				r.auditee AS Auditee,
				u3.name AS NamaAuditee,
				(CASE WHEN kts.id IS NULL THEN 'belum_dibuat' ELSE COALESCE(kts.status, 'draf') END) AS Status,
				'' AS NamaStandar,
				COALESCE(tdt.pertanyaan, vedt.indikator, '') AS Indikator,
				COALESCE(vedt.klasifikasi, tdt.klasifikasi, '') AS Klasifikasi,
				'' AS Capaian,
				(CASE WHEN dt.capaian_auditor = '1' OR dt.capaian_auditor = 'Ya' THEN 'Ya' WHEN dt.capaian_auditor = '0' OR dt.capaian_auditor = 'Tidak' THEN 'Tidak' ELSE COALESCE(dt.capaian_auditor, '') END) AS CapaianAuditor,
				COALESCE(kts.tindakan_perbaikan, '') AS TindakanPerbaikan,
				(CASE WHEN kts.status = 'tutup_kts' THEN 'Close' ELSE '' END) AS Keterangan,
				COALESCE(kts.keterangan_tolak_auditee, '') AS KeteranganTolakAuditee,
				COALESCE(kts.nomor_laporan, '') AS NomorLaporan,
				COALESCE(DATE_FORMAT(kts.tanggal_laporan, '%Y-%m-%d'), '') AS TanggalLaporan,
				COALESCE(kts.referensi, '') AS Referensi,
				COALESCE(kts.hasil_temuan, '') AS HasilTemuan,
				COALESCE(kts.uraian_ketidaksesuaian_p, '') AS KetidaksesuaianP,
				COALESCE(kts.uraian_ketidaksesuaian_l, '') AS KetidaksesuaianL,
				COALESCE(kts.uraian_ketidaksesuaian_o, '') AS KetidaksesuaianO,
				COALESCE(kts.uraian_ketidaksesuaian_r, '') AS KetidaksesuaianR,
				TRIM(CONCAT_WS(' ', NULLIF(kts.uraian_ketidaksesuaian_p, ''), NULLIF(kts.uraian_ketidaksesuaian_l, ''), NULLIF(kts.uraian_ketidaksesuaian_o, ''), NULLIF(kts.uraian_ketidaksesuaian_r, ''))) AS PLOR,
				COALESCE(kts.akar_masalah, '') AS AkarMasalah,
				COALESCE(kts.tindakan_koreksi, '') AS TindakanKoreksi,
				COALESCE(DATE_FORMAT(kts.tanggal_penyelesaian, '%Y-%m-%d'), '') AS TanggalPenyelesaian,
				COALESCE(kts.tinjauan_tindakan_perbaikan, '') AS TinjauanTindakanPerbaikan,
				COALESCE(DATE_FORMAT(kts.tanggal_closing_auditee, '%Y-%m-%d'), DATE_FORMAT(kts.tanggal_closing, '%Y-%m-%d'), '') AS TanggalClosing,
				COALESCE(NULLIF(LOWER(TRIM(COALESCE(vedt.tugas, tdt.tugas))), ''), 'auditor2') AS Tugas,
				r.tahun AS Tahun,
				r.fakultas_unit AS FakultasUnit,
				CASE 
					WHEN fu.type COLLATE utf8mb4_unicode_ci = 'prodi'
					THEN CONCAT(fu.nama_fak_prod_unit, ' (', fu.jenjang, ')')
					ELSE fu.nama_fak_prod_unit
				END AS Target,
				fu.uuid AS TargetUUID,
				r.uuid AS RenstraUUID,
				fu.jenjang AS Jenjang,
				fu.type AS Type,
				fu.fakultas AS Fakultas,
				'Terlaksana / Sesuai Template' AS TargetIndikator,
				'' AS TargetMin,
				'' AS TargetMax,
				NULL AS RangeTarget,
				'kualitatif' AS TipeTarget,
				'' AS Operator,
				'' AS Satuan
			FROM dokumen_tambahan dt
			JOIN renstra r ON dt.id_renstra = r.id
			JOIN view_evaluasi_dokumen_tambahan vedt ON dt.id = vedt.id
			LEFT JOIN template_dokumen_tambahan tdt ON dt.id_template_dokumen_tambahan = tdt.id
			LEFT JOIN kts_renstra kts ON dt.id = kts.id_dokumen_tambahan
			LEFT JOIN users u1 ON r.auditor1 = u1.id
			LEFT JOIN users u2 ON r.auditor2 = u2.id
			LEFT JOIN users u3 ON r.auditee = u3.id
			LEFT JOIN v_fakultas_unit fu ON r.fakultas_unit = fu.id
			WHERE LOWER(TRIM(vedt.klasifikasi)) != 'ok'
			  AND vedt.klasifikasi IS NOT NULL
			  AND TRIM(vedt.klasifikasi) != ''
		`
		var tbArgs []interface{}
		if strings.TrimSpace(renstraUuid) != "" && strings.TrimSpace(renstraUuid) != "all" {
			tambahanSQL += " AND (r.uuid = ? OR CAST(r.id AS CHAR) = ?)"
			tbArgs = append(tbArgs, renstraUuid, renstraUuid)
		}
		if normTugas == "auditor1" {
			tambahanSQL += " AND LOWER(TRIM(COALESCE(vedt.tugas, tdt.tugas))) = 'auditor1'"
		} else if normTugas == "auditor2" {
			tambahanSQL += " AND (COALESCE(vedt.tugas, tdt.tugas) IS NULL OR TRIM(COALESCE(vedt.tugas, tdt.tugas)) = '' OR LOWER(TRIM(COALESCE(vedt.tugas, tdt.tugas))) = 'auditor2')"
		}
		if tahun != "" {
			tambahanSQL += " AND r.tahun = ?"
			tbArgs = append(tbArgs, tahun)
		}
		tambahanSQL += " ORDER BY dt.id ASC"

		if err := r.db.WithContext(ctx).Raw(tambahanSQL, tbArgs...).Scan(&tambahanRows).Error; err != nil {
			return nil, err
		}
		results = append(results, tambahanRows...)
	}

	return results, nil
}

func (r *KtsRepository) GenerateKtsFromTemuan(
	ctx context.Context,
	renstraNilaiId *uint,
	renstraNilaiUuid *string,
	dokumenTambahanId *uint,
	dokumenTambahanUuid *string,
	auditorId string,
) (*uuid.UUID, error) {
	// 1. Resolve renstraNilaiId if UUID is provided
	if (renstraNilaiId == nil || *renstraNilaiId == 0) && renstraNilaiUuid != nil && *renstraNilaiUuid != "" {
		var rnRow struct {
			ID uint
		}
		if err := r.db.WithContext(ctx).Table("renstra_nilai").Select("id").Where("uuid = ?", *renstraNilaiUuid).First(&rnRow).Error; err == nil && rnRow.ID > 0 {
			renstraNilaiId = &rnRow.ID
		}
	}

	// 2. Resolve dokumenTambahanId if UUID is provided
	if (dokumenTambahanId == nil || *dokumenTambahanId == 0) && dokumenTambahanUuid != nil && *dokumenTambahanUuid != "" {
		var dtRow struct {
			ID uint
		}
		if err := r.db.WithContext(ctx).Table("dokumen_tambahan").Select("id").Where("uuid = ?", *dokumenTambahanUuid).First(&dtRow).Error; err == nil && dtRow.ID > 0 {
			dokumenTambahanId = &dtRow.ID
		}
	}

	if (renstraNilaiId == nil || *renstraNilaiId == 0) && (dokumenTambahanId == nil || *dokumenTambahanId == 0) {
		return nil, errors.New("renstraNilaiId/renstraNilaiUuid atau dokumenTambahanId/dokumenTambahanUuid wajib diisi dan valid")
	}

	// 3. Check existing KTS
	var existing domainKts.Kts
	var query *gorm.DB
	if renstraNilaiId != nil && *renstraNilaiId > 0 {
		query = r.db.WithContext(ctx).Where("id_renstra_nilai = ?", *renstraNilaiId)
	} else {
		query = r.db.WithContext(ctx).Where("id_dokumen_tambahan = ?", *dokumenTambahanId)
	}

	if err := query.First(&existing).Error; err == nil {
		return &existing.UUID, nil
	}

	// 4. Resolve auditor to integer string (kts_renstra.auditor column is INT in MySQL)
	var auditorIntStr *string
	trimmedAuditor := strings.TrimSpace(auditorId)
	if trimmedAuditor != "" {
		if _, err := strconv.ParseUint(trimmedAuditor, 10, 32); err == nil {
			auditorIntStr = &trimmedAuditor
		} else {
			// Query users table by uuid
			var userRow struct {
				ID uint
			}
			err := r.db.WithContext(ctx).Table("users").Select("id").
				Where("uuid = ?", trimmedAuditor).
				First(&userRow).Error
			if err == nil && userRow.ID > 0 {
				strID := strconv.FormatUint(uint64(userRow.ID), 10)
				auditorIntStr = &strID
			}
		}
	}

	// Fallback to assigned auditor from renstra if auditorIntStr is still nil
	if auditorIntStr == nil {
		if renstraNilaiId != nil && *renstraNilaiId > 0 {
			var fallbackAuditor struct {
				Auditor1 *string
				Auditor2 *string
				Tugas    string
			}
			_ = r.db.WithContext(ctx).Table("renstra_nilai rn").
				Select("r.auditor1 AS Auditor1, r.auditor2 AS Auditor2, tr.tugas AS Tugas").
				Joins("JOIN renstra r ON rn.id_renstra = r.id").
				Joins("JOIN template_renstra tr ON rn.template_renstra = tr.id").
				Where("rn.id = ?", *renstraNilaiId).Scan(&fallbackAuditor)
			if strings.ToLower(strings.TrimSpace(fallbackAuditor.Tugas)) == "auditor2" && fallbackAuditor.Auditor2 != nil && *fallbackAuditor.Auditor2 != "" {
				auditorIntStr = fallbackAuditor.Auditor2
			} else if fallbackAuditor.Auditor1 != nil && *fallbackAuditor.Auditor1 != "" {
				auditorIntStr = fallbackAuditor.Auditor1
			}
		} else if dokumenTambahanId != nil && *dokumenTambahanId > 0 {
			var fallbackAuditor struct {
				Auditor1 *string
				Auditor2 *string
			}
			_ = r.db.WithContext(ctx).Table("dokumen_tambahan dt").
				Select("r.auditor1 AS Auditor1, r.auditor2 AS Auditor2").
				Joins("JOIN renstra r ON dt.id_renstra = r.id").
				Where("dt.id = ?", *dokumenTambahanId).Scan(&fallbackAuditor)
			if fallbackAuditor.Auditor2 != nil && *fallbackAuditor.Auditor2 != "" {
				auditorIntStr = fallbackAuditor.Auditor2
			} else if fallbackAuditor.Auditor1 != nil && *fallbackAuditor.Auditor1 != "" {
				auditorIntStr = fallbackAuditor.Auditor1
			}
		}
	}

	newUUID := uuid.New()
	now := time.Now()

	var referensi string
	var hasilTemuan string

	if renstraNilaiId != nil && *renstraNilaiId > 0 {
		type RnInfo struct {
			Indikator   string
			Klasifikasi string
		}
		var info RnInfo
		_ = r.db.WithContext(ctx).Table("renstra_nilai rn").
			Select("COALESCE(venr.indikator, mir.indikator, '') AS Indikator, COALESCE(venr.klasifikasi, tr.klasifikasi, '') AS Klasifikasi").
			Joins("JOIN view_evaluasi_nilai_renstra venr ON rn.id = venr.id").
			Joins("JOIN template_renstra tr ON rn.template_renstra = tr.id").
			Joins("LEFT JOIN master_indikator_renstra mir ON tr.indikator = mir.id").
			Where("rn.id = ?", *renstraNilaiId).Scan(&info)
		referensi = info.Indikator
		hasilTemuan = info.Klasifikasi
	} else if dokumenTambahanId != nil && *dokumenTambahanId > 0 {
		type DtInfo struct {
			Pertanyaan  string
			Klasifikasi string
		}
		var info DtInfo
		_ = r.db.WithContext(ctx).Table("dokumen_tambahan dt").
			Select("COALESCE(tdt.pertanyaan, vedt.indikator, '') AS Pertanyaan, COALESCE(vedt.klasifikasi, tdt.klasifikasi, '') AS Klasifikasi").
			Joins("JOIN view_evaluasi_dokumen_tambahan vedt ON dt.id = vedt.id").
			Joins("LEFT JOIN template_dokumen_tambahan tdt ON dt.id_template_dokumen_tambahan = tdt.id").
			Where("dt.id = ?", *dokumenTambahanId).Scan(&info)
		referensi = info.Pertanyaan
		hasilTemuan = info.Klasifikasi
	}

	newKts := domainKts.Kts{
		UUID:            newUUID,
		RenstraNilai:    renstraNilaiId,
		DokumenTambahan: dokumenTambahanId,
		Status:          "draf",
		TanggalLaporan:  &now,
		Auditor:         auditorIntStr,
		Referensi:       &referensi,
		HasilTemuan:     &hasilTemuan,
	}

	if err := r.db.WithContext(ctx).Create(&newKts).Error; err != nil {
		return nil, err
	}

	return &newUUID, nil
}

