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

	var result domainrenstra.RenstraDefault

	err := r.db.WithContext(ctx).
		Table("renstra r").
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

		(
			SELECT COUNT(*)
			FROM renstra_nilai
			WHERE renstra_nilai.id_renstra = r.id
		) AS TotalRenstra,
		(
			SELECT COUNT(*)
			FROM renstra_nilai
			WHERE renstra_nilai.id_renstra = r.id AND renstra_nilai.capaian IS NOT NULL
		) AS TotalRenstraAuditee,
		(
			SELECT COUNT(*)
			FROM renstra_nilai
			WHERE renstra_nilai.id_renstra = r.id AND renstra_nilai.capaian_auditor IS NOT NULL
		) AS TotalRenstraAuditor,
		(
			SELECT COUNT(*)
			FROM dokumen_tambahan
			WHERE dokumen_tambahan.id_renstra = r.id
		) AS TotalDokumen,
		(
			SELECT COUNT(*)
			FROM dokumen_tambahan
			WHERE dokumen_tambahan.id_renstra = r.id AND dokumen_tambahan.file IS NOT NULL
		) AS TotalDokumenAuditee,
		(
			SELECT COUNT(*)
			FROM dokumen_tambahan
			WHERE dokumen_tambahan.id_renstra = r.id AND dokumen_tambahan.capaian_auditor IS NOT NULL
		) AS TotalDokumenAuditor,

		(
			SELECT COUNT(*)
			FROM renstra_nilai rn
			JOIN template_renstra tr ON rn.template_renstra = tr.id
			WHERE rn.id_renstra = r.id
			  AND LOWER(TRIM(tr.tugas)) = 'auditor1'
			  AND LOWER(TRIM(tr.klasifikasi)) != 'ok'
			  AND tr.klasifikasi IS NOT NULL
			  AND TRIM(tr.klasifikasi) != ''
		) AS TotalTemuanAuditor1Renstra,
		(
			SELECT COUNT(*)
			FROM dokumen_tambahan dt
			JOIN template_dokumen_tambahan tdt ON dt.id_template_dokumen_tambahan = tdt.id
			WHERE dt.id_renstra = r.id
			  AND LOWER(TRIM(tdt.tugas)) = 'auditor1'
			  AND LOWER(TRIM(tdt.klasifikasi)) != 'ok'
			  AND tdt.klasifikasi IS NOT NULL
			  AND TRIM(tdt.klasifikasi) != ''
		) AS TotalTemuanAuditor1Dokumen,
		(
			(
				SELECT COUNT(*)
				FROM renstra_nilai rn
				JOIN template_renstra tr ON rn.template_renstra = tr.id
				WHERE rn.id_renstra = r.id
				  AND LOWER(TRIM(tr.tugas)) = 'auditor1'
				  AND LOWER(TRIM(tr.klasifikasi)) != 'ok'
				  AND tr.klasifikasi IS NOT NULL
				  AND TRIM(tr.klasifikasi) != ''
			) + (
				SELECT COUNT(*)
				FROM dokumen_tambahan dt
				JOIN template_dokumen_tambahan tdt ON dt.id_template_dokumen_tambahan = tdt.id
				WHERE dt.id_renstra = r.id
				  AND LOWER(TRIM(tdt.tugas)) = 'auditor1'
				  AND LOWER(TRIM(tdt.klasifikasi)) != 'ok'
				  AND tdt.klasifikasi IS NOT NULL
				  AND TRIM(tdt.klasifikasi) != ''
			)
		) AS TotalTemuanAuditor1,
		(
			SELECT COUNT(*)
			FROM renstra_nilai rn
			JOIN template_renstra tr ON rn.template_renstra = tr.id
			WHERE rn.id_renstra = r.id
			  AND LOWER(TRIM(tr.tugas)) = 'auditor2'
			  AND LOWER(TRIM(tr.klasifikasi)) != 'ok'
			  AND tr.klasifikasi IS NOT NULL
			  AND TRIM(tr.klasifikasi) != ''
		) AS TotalTemuanAuditor2Renstra,
		(
			SELECT COUNT(*)
			FROM dokumen_tambahan dt
			JOIN template_dokumen_tambahan tdt ON dt.id_template_dokumen_tambahan = tdt.id
			WHERE dt.id_renstra = r.id
			  AND (tdt.tugas IS NULL OR TRIM(tdt.tugas) = '' OR LOWER(TRIM(tdt.tugas)) = 'auditor2')
			  AND LOWER(TRIM(tdt.klasifikasi)) != 'ok'
			  AND tdt.klasifikasi IS NOT NULL
			  AND TRIM(tdt.klasifikasi) != ''
		) AS TotalTemuanAuditor2Dokumen,
		(
			(
				SELECT COUNT(*)
				FROM renstra_nilai rn
				JOIN template_renstra tr ON rn.template_renstra = tr.id
				WHERE rn.id_renstra = r.id
				  AND LOWER(TRIM(tr.tugas)) = 'auditor2'
				  AND LOWER(TRIM(tr.klasifikasi)) != 'ok'
				  AND tr.klasifikasi IS NOT NULL
				  AND TRIM(tr.klasifikasi) != ''
			) + (
				SELECT COUNT(*)
				FROM dokumen_tambahan dt
				JOIN template_dokumen_tambahan tdt ON dt.id_template_dokumen_tambahan = tdt.id
				WHERE dt.id_renstra = r.id
				  AND (tdt.tugas IS NULL OR TRIM(tdt.tugas) = '' OR LOWER(TRIM(tdt.tugas)) = 'auditor2')
				  AND LOWER(TRIM(tdt.klasifikasi)) != 'ok'
				  AND tdt.klasifikasi IS NOT NULL
				  AND TRIM(tdt.klasifikasi) != ''
			)
		) AS TotalTemuanAuditor2,
		(
			(
				SELECT COUNT(*)
				FROM renstra_nilai rn
				JOIN template_renstra tr ON rn.template_renstra = tr.id
				LEFT JOIN kts_renstra k ON rn.id = k.id_renstra_nilai
				WHERE rn.id_renstra = r.id
				  AND LOWER(TRIM(tr.klasifikasi)) != 'ok'
				  AND tr.klasifikasi IS NOT NULL
				  AND TRIM(tr.klasifikasi) != ''
				  AND k.id IS NULL
			) + (
				SELECT COUNT(*)
				FROM dokumen_tambahan dt
				JOIN template_dokumen_tambahan tdt ON dt.id_template_dokumen_tambahan = tdt.id
				LEFT JOIN kts_renstra k ON dt.id = k.id_dokumen_tambahan
				WHERE dt.id_renstra = r.id
				  AND LOWER(TRIM(tdt.klasifikasi)) != 'ok'
				  AND tdt.klasifikasi IS NOT NULL
				  AND TRIM(tdt.klasifikasi) != ''
				  AND k.id IS NULL
			)
		) AS TotalKtsBelumDibuat,
		(
			SELECT COUNT(*)
			FROM kts_renstra k
			LEFT JOIN renstra_nilai rn ON k.id_renstra_nilai = rn.id
			LEFT JOIN dokumen_tambahan dt ON k.id_dokumen_tambahan = dt.id
			WHERE (rn.id_renstra = r.id OR dt.id_renstra = r.id)
			  AND k.status = 'draf'
		) AS TotalKtsDraf,
		(
			SELECT COUNT(*)
			FROM kts_renstra k
			LEFT JOIN renstra_nilai rn ON k.id_renstra_nilai = rn.id
			LEFT JOIN dokumen_tambahan dt ON k.id_dokumen_tambahan = dt.id
			WHERE (rn.id_renstra = r.id OR dt.id_renstra = r.id)
			  AND k.status = 'menunggu_verif_auditee'
		) AS TotalKtsMenungguVerif,
		(
			SELECT COUNT(*)
			FROM kts_renstra k
			LEFT JOIN renstra_nilai rn ON k.id_renstra_nilai = rn.id
			LEFT JOIN dokumen_tambahan dt ON k.id_dokumen_tambahan = dt.id
			WHERE (rn.id_renstra = r.id OR dt.id_renstra = r.id)
			  AND k.status = 'terima_auditee'
		) AS TotalKtsTerimaAuditee,
		(
			SELECT COUNT(*)
			FROM kts_renstra k
			LEFT JOIN renstra_nilai rn ON k.id_renstra_nilai = rn.id
			LEFT JOIN dokumen_tambahan dt ON k.id_dokumen_tambahan = dt.id
			WHERE (rn.id_renstra = r.id OR dt.id_renstra = r.id)
			  AND k.status = 'tolak_auditee'
		) AS TotalKtsTolakAuditee,
		(
			SELECT COUNT(*)
			FROM kts_renstra k
			LEFT JOIN renstra_nilai rn ON k.id_renstra_nilai = rn.id
			LEFT JOIN dokumen_tambahan dt ON k.id_dokumen_tambahan = dt.id
			WHERE (rn.id_renstra = r.id OR dt.id_renstra = r.id)
			  AND k.status = 'menunggu_penyelesaian'
		) AS TotalKtsMenungguPenyelesaian,
		(
			SELECT COUNT(*)
			FROM kts_renstra k
			LEFT JOIN renstra_nilai rn ON k.id_renstra_nilai = rn.id
			LEFT JOIN dokumen_tambahan dt ON k.id_dokumen_tambahan = dt.id
			WHERE (rn.id_renstra = r.id OR dt.id_renstra = r.id)
			  AND k.status = 'tindakan_penyelesaian'
		) AS TotalKtsTindakanPenyelesaian,
		(
			SELECT COUNT(*)
			FROM kts_renstra k
			LEFT JOIN renstra_nilai rn ON k.id_renstra_nilai = rn.id
			LEFT JOIN dokumen_tambahan dt ON k.id_dokumen_tambahan = dt.id
			WHERE (rn.id_renstra = r.id OR dt.id_renstra = r.id)
			  AND k.status NOT IN ('draf', 'tutup_kts')
		) AS TotalKtsProses,
		(
			SELECT COUNT(*)
			FROM kts_renstra k
			LEFT JOIN renstra_nilai rn ON k.id_renstra_nilai = rn.id
			LEFT JOIN dokumen_tambahan dt ON k.id_dokumen_tambahan = dt.id
			WHERE (rn.id_renstra = r.id OR dt.id_renstra = r.id)
			  AND k.status = 'tutup_kts'
		) AS TotalKtsTutup,
		(
			SELECT COUNT(*)
			FROM kts_renstra k
			LEFT JOIN renstra_nilai rn ON k.id_renstra_nilai = rn.id
			LEFT JOIN dokumen_tambahan dt ON k.id_dokumen_tambahan = dt.id
			WHERE (rn.id_renstra = r.id OR dt.id_renstra = r.id)
		) AS TotalKts
	`).
		Joins("JOIN v_fakultas_unit vfu ON r.fakultas_unit = vfu.id").
		Joins("LEFT JOIN users a0 ON r.auditee = a0.id").
		Joins("LEFT JOIN users a1 ON r.auditor1 = a1.id").
		Joins("LEFT JOIN users a2 ON r.auditor2 = a2.id").
		Where("r.uuid = ?", id).
		Take(&result).Error // Take = LIMIT 1

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}

	return &result, nil
}

var allowedSearchColumns = map[string]string{
	"tahun":           "r.tahun",
	"fakultasunit":    "vfu.nama_fak_prod_unit",
	"fakultas_unit":   "vfu.nama_fak_prod_unit",
	"targetaudit":     "vfu.nama_fak_prod_unit",
	"target_audit":    "vfu.nama_fak_prod_unit",
	"namafakprodunit": "vfu.nama_fak_prod_unit",
	"type":            "vfu.type",
	"tipe":            "vfu.type",
	"jenjang":         "vfu.jenjang",
	"fakultas":        "vfu.fakultas",
	"auditee":         "a0.name",
	"auditor1":        "a1.name",
	"auditor2":        "a2.name",
	"uuidauditee":     "a0.uuid",
	"uuid_auditee":    "a0.uuid",
	"uuidauditor1":    "a1.uuid",
	"uuid_auditor1":   "a1.uuid",
	"uuidauditor2":    "a2.uuid",
	"uuid_auditor2":   "a2.uuid",
	"kodeakses":       "r.kodeAkses",
	"kode_akses":      "r.kodeAkses",
}

var globalSearchColumns = []string{
	"r.tahun",
	"vfu.nama_fak_prod_unit",
	"vfu.type",
	"vfu.jenjang",
	"vfu.fakultas",
	"a0.name",
	"a1.name",
	"a2.name",
	"r.kodeAkses",
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

	var renstras = make([]domainrenstra.RenstraDefault, 0)
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

			(
				SELECT
					COUNT(*)
				FROM
					renstra_nilai
				WHERE
					renstra_nilai.id_renstra = r.id
			) AS TotalRenstra,
			(
				SELECT
					COUNT(*)
				FROM
					renstra_nilai
				WHERE
					renstra_nilai.id_renstra = r.id AND renstra_nilai.capaian IS NOT NULL
			) AS TotalRenstraAuditee,
			(
				SELECT
					COUNT(*)
				FROM
					renstra_nilai
				WHERE
					renstra_nilai.id_renstra = r.id AND renstra_nilai.capaian_auditor IS NOT NULL
			) AS TotalRenstraAuditor,
			(
				SELECT
					COUNT(*)
				FROM
					dokumen_tambahan
				WHERE
					dokumen_tambahan.id_renstra = r.id
			) AS TotalDokumen,
			(
				SELECT
					COUNT(*)
				FROM
					dokumen_tambahan
				WHERE
					dokumen_tambahan.id_renstra = r.id AND dokumen_tambahan.file IS NOT NULL
			) AS TotalDokumenAuditee,
			(
				SELECT
					COUNT(*)
				FROM
					dokumen_tambahan
				WHERE
					dokumen_tambahan.id_renstra = r.id AND dokumen_tambahan.capaian_auditor IS NOT NULL
			) AS TotalDokumenAuditor,

			(
				SELECT COUNT(*)
				FROM renstra_nilai rn
				JOIN template_renstra tr ON rn.template_renstra = tr.id
				WHERE rn.id_renstra = r.id
				  AND LOWER(TRIM(tr.tugas)) = 'auditor1'
				  AND LOWER(TRIM(tr.klasifikasi)) != 'ok'
				  AND tr.klasifikasi IS NOT NULL
				  AND TRIM(tr.klasifikasi) != ''
			) AS TotalTemuanAuditor1Renstra,
			(
				SELECT COUNT(*)
				FROM dokumen_tambahan dt
				JOIN template_dokumen_tambahan tdt ON dt.id_template_dokumen_tambahan = tdt.id
				WHERE dt.id_renstra = r.id
				  AND LOWER(TRIM(tdt.tugas)) = 'auditor1'
				  AND LOWER(TRIM(tdt.klasifikasi)) != 'ok'
				  AND tdt.klasifikasi IS NOT NULL
				  AND TRIM(tdt.klasifikasi) != ''
			) AS TotalTemuanAuditor1Dokumen,
			(
				(
					SELECT COUNT(*)
					FROM renstra_nilai rn
					JOIN template_renstra tr ON rn.template_renstra = tr.id
					WHERE rn.id_renstra = r.id
					  AND LOWER(TRIM(tr.tugas)) = 'auditor1'
					  AND LOWER(TRIM(tr.klasifikasi)) != 'ok'
					  AND tr.klasifikasi IS NOT NULL
					  AND TRIM(tr.klasifikasi) != ''
				) + (
					SELECT COUNT(*)
					FROM dokumen_tambahan dt
					JOIN template_dokumen_tambahan tdt ON dt.id_template_dokumen_tambahan = tdt.id
					WHERE dt.id_renstra = r.id
					  AND LOWER(TRIM(tdt.tugas)) = 'auditor1'
					  AND LOWER(TRIM(tdt.klasifikasi)) != 'ok'
					  AND tdt.klasifikasi IS NOT NULL
					  AND TRIM(tdt.klasifikasi) != ''
				)
			) AS TotalTemuanAuditor1,
			(
				SELECT COUNT(*)
				FROM renstra_nilai rn
				JOIN template_renstra tr ON rn.template_renstra = tr.id
				WHERE rn.id_renstra = r.id
				  AND LOWER(TRIM(tr.tugas)) = 'auditor2'
				  AND LOWER(TRIM(tr.klasifikasi)) != 'ok'
				  AND tr.klasifikasi IS NOT NULL
				  AND TRIM(tr.klasifikasi) != ''
			) AS TotalTemuanAuditor2Renstra,
			(
				SELECT COUNT(*)
				FROM dokumen_tambahan dt
				JOIN template_dokumen_tambahan tdt ON dt.id_template_dokumen_tambahan = tdt.id
				WHERE dt.id_renstra = r.id
				  AND (tdt.tugas IS NULL OR TRIM(tdt.tugas) = '' OR LOWER(TRIM(tdt.tugas)) = 'auditor2')
				  AND LOWER(TRIM(tdt.klasifikasi)) != 'ok'
				  AND tdt.klasifikasi IS NOT NULL
				  AND TRIM(tdt.klasifikasi) != ''
			) AS TotalTemuanAuditor2Dokumen,
			(
				(
					SELECT COUNT(*)
					FROM renstra_nilai rn
					JOIN template_renstra tr ON rn.template_renstra = tr.id
					WHERE rn.id_renstra = r.id
					  AND LOWER(TRIM(tr.tugas)) = 'auditor2'
					  AND LOWER(TRIM(tr.klasifikasi)) != 'ok'
					  AND tr.klasifikasi IS NOT NULL
					  AND TRIM(tr.klasifikasi) != ''
				) + (
					SELECT COUNT(*)
					FROM dokumen_tambahan dt
					JOIN template_dokumen_tambahan tdt ON dt.id_template_dokumen_tambahan = tdt.id
					WHERE dt.id_renstra = r.id
					  AND (tdt.tugas IS NULL OR TRIM(tdt.tugas) = '' OR LOWER(TRIM(tdt.tugas)) = 'auditor2')
					  AND LOWER(TRIM(tdt.klasifikasi)) != 'ok'
					  AND tdt.klasifikasi IS NOT NULL
					  AND TRIM(tdt.klasifikasi) != ''
				)
			) AS TotalTemuanAuditor2,
			(
				(
					SELECT COUNT(*)
					FROM renstra_nilai rn
					JOIN template_renstra tr ON rn.template_renstra = tr.id
					LEFT JOIN kts_renstra k ON rn.id = k.id_renstra_nilai
					WHERE rn.id_renstra = r.id
					  AND LOWER(TRIM(tr.klasifikasi)) != 'ok'
					  AND tr.klasifikasi IS NOT NULL
					  AND TRIM(tr.klasifikasi) != ''
					  AND k.id IS NULL
				) + (
					SELECT COUNT(*)
					FROM dokumen_tambahan dt
					JOIN template_dokumen_tambahan tdt ON dt.id_template_dokumen_tambahan = tdt.id
					LEFT JOIN kts_renstra k ON dt.id = k.id_dokumen_tambahan
					WHERE dt.id_renstra = r.id
					  AND LOWER(TRIM(tdt.klasifikasi)) != 'ok'
					  AND tdt.klasifikasi IS NOT NULL
					  AND TRIM(tdt.klasifikasi) != ''
					  AND k.id IS NULL
				)
			) AS TotalKtsBelumDibuat,
			(
				SELECT COUNT(*)
				FROM kts_renstra k
				LEFT JOIN renstra_nilai rn ON k.id_renstra_nilai = rn.id
				LEFT JOIN dokumen_tambahan dt ON k.id_dokumen_tambahan = dt.id
				WHERE (rn.id_renstra = r.id OR dt.id_renstra = r.id)
				  AND k.status = 'draf'
			) AS TotalKtsDraf,
			(
				SELECT COUNT(*)
				FROM kts_renstra k
				LEFT JOIN renstra_nilai rn ON k.id_renstra_nilai = rn.id
				LEFT JOIN dokumen_tambahan dt ON k.id_dokumen_tambahan = dt.id
				WHERE (rn.id_renstra = r.id OR dt.id_renstra = r.id)
				  AND k.status = 'menunggu_verif_auditee'
			) AS TotalKtsMenungguVerif,
			(
				SELECT COUNT(*)
				FROM kts_renstra k
				LEFT JOIN renstra_nilai rn ON k.id_renstra_nilai = rn.id
				LEFT JOIN dokumen_tambahan dt ON k.id_dokumen_tambahan = dt.id
				WHERE (rn.id_renstra = r.id OR dt.id_renstra = r.id)
				  AND k.status = 'terima_auditee'
			) AS TotalKtsTerimaAuditee,
			(
				SELECT COUNT(*)
				FROM kts_renstra k
				LEFT JOIN renstra_nilai rn ON k.id_renstra_nilai = rn.id
				LEFT JOIN dokumen_tambahan dt ON k.id_dokumen_tambahan = dt.id
				WHERE (rn.id_renstra = r.id OR dt.id_renstra = r.id)
				  AND k.status = 'tolak_auditee'
			) AS TotalKtsTolakAuditee,
			(
				SELECT COUNT(*)
				FROM kts_renstra k
				LEFT JOIN renstra_nilai rn ON k.id_renstra_nilai = rn.id
				LEFT JOIN dokumen_tambahan dt ON k.id_dokumen_tambahan = dt.id
				WHERE (rn.id_renstra = r.id OR dt.id_renstra = r.id)
				  AND k.status = 'menunggu_penyelesaian'
			) AS TotalKtsMenungguPenyelesaian,
			(
				SELECT COUNT(*)
				FROM kts_renstra k
				LEFT JOIN renstra_nilai rn ON k.id_renstra_nilai = rn.id
				LEFT JOIN dokumen_tambahan dt ON k.id_dokumen_tambahan = dt.id
				WHERE (rn.id_renstra = r.id OR dt.id_renstra = r.id)
				  AND k.status = 'tindakan_penyelesaian'
			) AS TotalKtsTindakanPenyelesaian,
			(
				SELECT COUNT(*)
				FROM kts_renstra k
				LEFT JOIN renstra_nilai rn ON k.id_renstra_nilai = rn.id
				LEFT JOIN dokumen_tambahan dt ON k.id_dokumen_tambahan = dt.id
				WHERE (rn.id_renstra = r.id OR dt.id_renstra = r.id)
				  AND k.status NOT IN ('draf', 'tutup_kts')
			) AS TotalKtsProses,
			(
				SELECT COUNT(*)
				FROM kts_renstra k
				LEFT JOIN renstra_nilai rn ON k.id_renstra_nilai = rn.id
				LEFT JOIN dokumen_tambahan dt ON k.id_dokumen_tambahan = dt.id
				WHERE (rn.id_renstra = r.id OR dt.id_renstra = r.id)
				  AND k.status = 'tutup_kts'
			) AS TotalKtsTutup,
			(
				SELECT COUNT(*)
				FROM kts_renstra k
				LEFT JOIN renstra_nilai rn ON k.id_renstra_nilai = rn.id
				LEFT JOIN dokumen_tambahan dt ON k.id_dokumen_tambahan = dt.id
				WHERE (rn.id_renstra = r.id OR dt.id_renstra = r.id)
			) AS TotalKts
		`).
		Joins("JOIN v_fakultas_unit vfu ON r.fakultas_unit = vfu.id").
		Joins("LEFT JOIN users a0 ON r.auditee = a0.id").
		Joins("LEFT JOIN users a1 ON r.auditor1 = a1.id").
		Joins("LEFT JOIN users a2 ON r.auditor2 = a2.id")

	// -------------------------------
	// SEARCH FILTERS (ADVANCED)
	// -------------------------------
	var auditOrParts []string
	var auditParams []interface{}

	if len(searchFilters) > 0 {
		for _, f := range searchFilters {
			rawField := strings.TrimSpace(strings.ToLower(f.Field))
			operator := strings.TrimSpace(strings.ToLower(f.Operator))

			var value string
			if f.Value != nil {
				value = strings.TrimSpace(*f.Value)
			}

			if value == "" {
				continue
			}

			// Validate allowed column with fallback
			col, ok := allowedSearchColumns[rawField]
			if !ok {
				col, ok = allowedSearchColumns[strings.ReplaceAll(rawField, "_", "")]
			}
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
				parts := strings.Split(value, ",")
				var inVals []string
				for _, p := range parts {
					if trimmed := strings.TrimSpace(p); trimmed != "" {
						inVals = append(inVals, trimmed)
					}
				}
				if len(inVals) > 0 {
					db = db.Where(fmt.Sprintf("%s IN (?)", col), inVals)
				}
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

		for _, col := range globalSearchColumns {
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

func (r *RenstraRepository) ensureUuidColumn(ctx context.Context, tableName string) error {
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

func (r *RenstraRepository) SetupUuid(ctx context.Context) error {
	if err := r.ensureUuidColumn(ctx, "renstra"); err != nil {
		return err
	}

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
