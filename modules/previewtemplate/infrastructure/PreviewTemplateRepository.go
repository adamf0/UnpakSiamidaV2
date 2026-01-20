package infrastructure

import (
	domainpreviewtemplate "UnpakSiamida/modules/previewtemplate/domain"
	"context"

	"gorm.io/gorm"
)

type PreviewTemplateRepository struct {
	db *gorm.DB
}

func NewPreviewTemplateRepository(db *gorm.DB) domainpreviewtemplate.IPreviewTemplateRepository {
	return &PreviewTemplateRepository{db: db}
}

func (r *PreviewTemplateRepository) GetByTahunFakultasUnit(
	ctx context.Context,
	tahun string,
	fakultasUnit string,
) ([]domainpreviewtemplate.PreviewTemplate, error) {

	var results []domainpreviewtemplate.PreviewTemplate

	sql := `
		SELECT
			tr.uuid as UUID,
			tr.tahun AS Tahun,
			i.id AS IndikatorId,
			i.indikator AS Indikator,
			i.tahun AS IndikatorTahun,
			tr.pertanyaan AS IsPertanyaan,
			i.parent as ParentIndikatorId,
			tr.fakultas_unit AS FakultasUnitId,
			fu.nama_fak_prod_unit AS FakultasUnit,
			fu.type AS FakultasUnitType,
			fu.fakultas AS Fakultas,
			tr.klasifikasi AS Klasifikasi,
			tr.satuan,
			(
				CASE
					WHEN tr.target IS NOT NULL
						THEN tr.target
					WHEN tr.target IS NULL
						THEN CONCAT(tr.target_min, '≥ nilai ≤', tr.target_max)
					ELSE 'ada yg salah pada datanya'
				END
			) AS Target,
			tr.kategori AS Kategori
		FROM template_renstra tr
		INNER JOIN master_indikator_renstra i
			ON i.id = tr.indikator
		INNER JOIN v_fakultas_unit fu
			ON fu.id = tr.fakultas_unit
		WHERE
			tr.tahun = ?
			AND tr.fakultas_unit = ?
		ORDER BY
			tr.tahun ASC,
			tr.fakultas_unit ASC,
			tr.indikator ASC
	`

	err := r.db.WithContext(ctx).
		Raw(sql, tahun, fakultasUnit).
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	return results, nil
}

func (r *PreviewTemplateRepository) GetByTahunTag(
	ctx context.Context,
	tahun string,
	tag string,
) ([]domainpreviewtemplate.PreviewTemplate, error) {

	var results []domainpreviewtemplate.PreviewTemplate

	sql := `
		SELECT
				dt.uuid as UUID,
				dt.tahun AS Tahun,
				NULL AS IndikatorId,
				dt.pertanyaan AS Indikator,
				NULL AS IndikatorTahun,
				1 AS IsPertanyaan,
				NULL AS ParentIndikatorId,
				null AS FakultasUnitId,
				null AS FakultasUnit,
				null AS FakultasUnitType,
				null AS Fakultas,
				dt.klasifikasi AS Klasifikasi,
				NULL AS Satuan,
				'ya' AS Target,
				NULL AS Kategori
			FROM template_dokumen_tambahan dt
			INNER JOIN jenis_file_renstra jf ON jf.id = dt.jenis_file
			WHERE
				dt.tahun = ?
				AND dt.fakultas_prodi_unit = ?
			ORDER BY
				dt.tahun ASC,
				dt.id ASC;
	`

	err := r.db.WithContext(ctx).
		Raw(sql, tahun, tag).
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	return results, nil
}
