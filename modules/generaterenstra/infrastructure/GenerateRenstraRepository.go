package infrastructure

import (
	"context"
	// commondomaingeneraterenstra "UnpakSiamida/common/domain"
	domaingeneraterenstra "UnpakSiamida/modules/generaterenstra/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"

	// "strings"
	"errors"
)

type GenerateRenstraRepository struct {
	db *gorm.DB
	tx *gorm.DB
}

func NewGenerateRenstraRepository(db *gorm.DB) domaingeneraterenstra.IGenerateRenstraRepository {
	return &GenerateRenstraRepository{
		db: db,
		tx: nil,
	}
}

func (r *GenerateRenstraRepository) BeginTx(
	ctx context.Context,
) (*gorm.DB, error) {
	r.tx = r.db.WithContext(ctx).Begin()
	if r.tx.Error != nil {
		return nil, r.tx.Error
	}
	return r.tx, nil
}

func (r *GenerateRenstraRepository) Commit() error {
	if r.tx == nil {
		return errors.New("transaction not started")
	}

	err := r.tx.Commit().Error
	r.tx = nil
	return err
}

func (r *GenerateRenstraRepository) Rollback() {
	if r.tx != nil {
		_ = r.tx.Rollback()
		r.tx = nil
	}
}

// ------------------------
// GET ALL BY Tahun & FakultasUnit
// ------------------------
// func (r *GenerateRenstraRepository) GetAllByTahunFakUnit(
// 	ctx context.Context,
// 	tahun string,
// 	fakultasUnit uint,
// ) ([]domaingeneraterenstra.GenerateRenstra, error) {

// 	var generaterenstras []domaingeneraterenstra.GenerateRenstra

// 	err := r.db.WithContext(ctx).
// 		Where("tahun = ? AND fakultas_unit = ?", tahun, fakultasUnit).
// 		Find(&generaterenstras).Error

// 	if err != nil {
// 		return nil, err
// 	}

// 	return generaterenstras, nil
// }

// ------------------------
// GET ALL Renstra Nilai BY Tahun & FakultasUnit (default)
// ------------------------
func (r *GenerateRenstraRepository) GetAllRenstraNilaiByTahunFakUnitDefault(
	ctx context.Context,
	tahun string,
	fakultasUnit uint,
) ([]domaingeneraterenstra.GenerateRenstraDefault, error) {

	var results = make([]domaingeneraterenstra.GenerateRenstraDefault, 0)

	err := r.db.WithContext(ctx).
		Table("renstra_nilai gr").
		Select(`
		gr.id AS ID,
		gr.uuid AS UUID,
		gr.id_renstra AS RenstraId,
		tr.uuid AS TemplateRenstraUuid,
		gr.template_renstra AS TemplateRenstra,
		i.indikator AS Indikator,
		gr.tugas AS Tugas,
		r.tahun AS RenstraTahun,
		tr.tahun AS TemplateTahun
	`).
		Joins("INNER JOIN renstra r ON r.id = gr.id_renstra").
		Joins("INNER JOIN template_renstra tr ON tr.id = gr.template_renstra").
		Joins("INNER JOIN master_indikator_renstra i ON tr.indikator = i.id").
		Where("r.tahun = ? AND r.fakultas_unit = ?", tahun, fakultasUnit).
		Find(&results).Error

	if err != nil {
		return nil, err
	}

	return results, nil

}

// ------------------------
// GET ALL Dokumen Tambahan BY Tahun & FakultasUnit (default)
// ------------------------
func (r *GenerateRenstraRepository) GetAllDokumenTambahanByTahunFakUnitDefault(
	ctx context.Context,
	tahun string,
	fakultasUnit uint,
) ([]domaingeneraterenstra.GenerateDokumenTambahanDefault, error) {

	var results = make([]domaingeneraterenstra.GenerateDokumenTambahanDefault, 0)

	err := r.db.WithContext(ctx).
		Table("dokumen_tambahan d").
		Select(`
		d.id AS ID,
		d.uuid AS UUID,
		r.id AS RenstraId,
		dt.uuid AS TemplateDokumenTambahanUuid,
		dt.id AS TemplateDokumenTambahan,
		jf.id AS JenisFileId,
		jf.nama AS JenisFile,
		dt.pertanyaan AS Pertanyaan,
		dt.tugas AS Tugas,
		r.fakultas_unit AS FakultasUnit,
		r.tahun AS RenstraTahun,
		dt.tahun AS TemplateTahun
	`).
		Joins("INNER JOIN renstra r ON r.id = d.id_renstra").
		Joins("INNER JOIN template_dokumen_tambahan dt ON dt.id = d.id_template_dokumen_tambahan").
		Joins("INNER JOIN jenis_file jf ON jf.id = dt.jenis_file").
		Where("r.tahun = ? AND r.fakultas_unit = ?", tahun, fakultasUnit).
		Find(&results).Error

	if err != nil {
		return nil, err
	}

	return results, nil

}

// ------------------------
// GET DEFAULT BY UUID
// ------------------------
// func (r *GenerateRenstraRepository) GetDefaultByUuid(
// 	ctx context.Context,
// 	id uuid.UUID,
// ) (*domaingeneraterenstra.RenstraDefault, error) {

// 	query := `
// 		SELECT
// 			r.id as Id,
// 			r.uuid as Uuid,
// 			r.tahun as Tahun,
// 			r.fakultas_unit as FakultasUnit,
// 			r.periode_upload_mulai as PeriodeUploadMulai,
// 			r.periode_upload_akhir as PeriodeUploadAkhir,
// 			r.periode_assesment_dokumen_mulai as PeriodeAssesmentDokumenMulai,
// 			r.periode_assesment_dokumen_akhir as PeriodeAssesmentDokumenAkhir,
// 			r.periode_assesment_lapangan_mulai as PeriodeAssesmentLapanganMulai,
// 			r.periode_assesment_lapangan_akhir as PeriodeAssesmentLapanganAkhir,
// 			r.auditee as Auditee,
// 			r.auditor1 as Auditor1,
// 			r.auditor2  as Auditor2,
// 			r.kodeAkses as KodeAkses,
// 			r.catatan as Catatan1,
// 			r.catatan2 as Catatan2,
// 			a0.name AS NamaAuditee,
// 			a1.name AS NamaAuditor1,
// 			a2.name AS NamaAuditor2,
// 			vfu.uuid AS UUIDFakultasUnit,
// 			vfu.nama_fak_prod_unit as NamaFakultasUnit,
// 			vfu.jenjang as Jenjang,
// 			vfu.type as Type,
// 			vfu.fakultas as Fakultas
// 		FROM renstra r
// 		JOIN v_fakultas_unit vfu ON r.fakultas_unit = vfu.id
// 		LEFT JOIN users a0 ON r.auditee = a0.id
// 		LEFT JOIN users a1 ON r.auditor1 = a1.id
// 		LEFT JOIN users a2 ON r.auditor2 = a2.id
// 		WHERE r.uuid = ?
// 		ORDER BY r.tahun DESC
// 		LIMIT 1
// 	`

// 	var rowData domaingeneraterenstra.RenstraDefault

// 	row := r.db.WithContext(ctx).Raw(query, id).Row()
// 	err := row.Scan(
// 		&rowData.ID,
// 		&rowData.UUID,
// 		&rowData.Tahun,
// 		&rowData.FakultasUnit,
// 		&rowData.PeriodeUploadMulai,
// 		&rowData.PeriodeUploadAkhir,
// 		&rowData.PeriodeAssesmentDokumenMulai,
// 		&rowData.PeriodeAssesmentDokumenAkhir,
// 		&rowData.PeriodeAssesmentLapanganMulai,
// 		&rowData.PeriodeAssesmentLapanganAkhir,
// 		&rowData.Auditee,
// 		&rowData.Auditor1,
// 		&rowData.Auditor2,
// 		&rowData.KodeAkses,
// 		&rowData.Catatan1,
// 		&rowData.Catatan2,
// 		&rowData.NamaAuditee,
// 		&rowData.NamaAuditor1,
// 		&rowData.NamaAuditor2,
// 		&rowData.UUIDFakultasUnit,
// 		&rowData.NamaFakultasUnit,
// 		&rowData.Jenjang,
// 		&rowData.Type,
// 		&rowData.Fakultas,
// 	)

// 	if err != nil {
// 		return nil, err
// 	}

// 	if rowData.ID == 0 {
// 		return nil, gorm.ErrRecordNotFound
// 	}

// 	return &rowData, nil
// }

// ------------------------
// CREATE Renstra Nilai
// ------------------------
func (r *GenerateRenstraRepository) CreateRenstraNilai(
	ctx context.Context,
	tx *gorm.DB,
	data *domaingeneraterenstra.GenerateRenstra,
) error {
	return tx.WithContext(ctx).Create(data).Error
}

// ------------------------
// DELETE Renstra Nilai
// ------------------------
func (r *GenerateRenstraRepository) DeleteRenstraNilai(
	ctx context.Context,
	tx *gorm.DB,
	data *domaingeneraterenstra.GenerateRenstra,
) error {
	return tx.WithContext(ctx).
		Where(`
			id_renstra = ?
			AND template_renstra = ?
			AND tugas = ?
		`,
			data.RenstraId,
			data.TemplateRenstra,
			data.Tugas,
		).
		Delete(&domaingeneraterenstra.GenerateRenstra{}).
		Error
}

// ------------------------
// CREATE Dokumen Tambahan
// ------------------------
func (r *GenerateRenstraRepository) CreateDokumenTambahan(
	ctx context.Context,
	tx *gorm.DB,
	data *domaingeneraterenstra.GenerateDokumenTambahan,
) error {
	return tx.WithContext(ctx).Create(data).Error
}

// ------------------------
// DELETE Dokumen Tambahan
// ------------------------
func (r *GenerateRenstraRepository) DeleteDokumenTambahan(
	ctx context.Context,
	tx *gorm.DB,
	data *domaingeneraterenstra.GenerateDokumenTambahan,
) error {
	return tx.WithContext(ctx).
		Where(`
			id_renstra = ?
			AND id_template_dokumen_tambahan = ?
			AND tugas = ?
		`,
			data.RenstraId,
			data.TemplateDokumenTambahan,
			data.Tugas,
		).
		Delete(&domaingeneraterenstra.GenerateDokumenTambahan{}).
		Error
}

// ------------------------
// DELETE Renstra Nilai (by UUID)
// ------------------------
func (r *GenerateRenstraRepository) ForceDeleteRenstraNilai(ctx context.Context, uid uuid.UUID, renstra uint) error {
	return r.db.WithContext(ctx).
		Where("uuid = ?", uid).
		Where("id_renstra = ?", renstra).
		Delete(&domaingeneraterenstra.GenerateRenstra{}).Error
}

// ------------------------
// DELETE Dokumen Tambahan (by UUID)
// ------------------------
func (r *GenerateRenstraRepository) ForceDeleteDokumenTambahan(ctx context.Context, uid uuid.UUID, renstra uint) error {
	return r.db.WithContext(ctx).
		Where("uuid = ?", uid).
		Where("id_renstra = ?", renstra).
		Delete(&domaingeneraterenstra.GenerateDokumenTambahan{}).Error
}
