package infrastructure

import (
	domainlaporan "UnpakSiamida/modules/laporan/domain"
	"context"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LaporanRepository struct {
	db *gorm.DB
}

func NewLaporanRepository(db *gorm.DB) domainlaporan.ILaporanRepository {
	return &LaporanRepository{db: db}
}

// ------------------------
// GET MONITORING BY TARGET & TAHUN
// ------------------------
func (r *LaporanRepository) GetMonitoringByTargetTahun(
	ctx context.Context,
	uuidFakultasUnit uuid.UUID,
	uuidTahun uuid.UUID,
	page, limit *int,
) ([]domainlaporan.MonitoringProker, int64, error) {
	var result = make([]domainlaporan.MonitoringProker, 0)
	var total int64

	// Cek tahun dari master_tahun
	var mtRow struct {
		ID    uint      `gorm:"column:id"`
		UUID  uuid.UUID `gorm:"column:uuid"`
		Tahun string    `gorm:"column:tahun"`
	}
	_ = r.db.WithContext(ctx).Table("master_tahun").Where("uuid = ?", uuidTahun.String()).Take(&mtRow).Error

	tahunNum, errAtoi := strconv.Atoi(strings.TrimSpace(mtRow.Tahun))
	isV3 := errAtoi == nil && tahunNum >= 2027

	if isV3 {
		// V3 (Tahun >= 2027): Mengambil data dokumen dari matrix_proker_dokumen
		// berelasi dengan matrix_proker dan mata_program dengan tahun >= 2027
		db := r.db.WithContext(ctx).Table("mata_program mp").
			Select(`
				vfu.id as FakultasUnitId,
				vfu.uuid as FakultasUnitUuid,
				vfu.nama_fak_prod_unit as FakultasUnit,
				vfu.jenjang as Jenjang,
				vfu.type as Type,
				vfu.fakultas as Fakultas,

				mp.id as MataProgramId,
				mp.uuid as MataProgramUUID,
				mp.mata_program as MataProgram,

				mt.id as TahunId,
				mt.uuid as TahunUUID,
				mt.tahun as Tahun,

				COALESCE(COUNT(CASE WHEN UPPER(TRIM(mpd.jenis_dokumen)) = 'SK' THEN 1 END), 0) as Sk,
				COALESCE(COUNT(CASE WHEN UPPER(TRIM(mpd.jenis_dokumen)) = 'SK' AND (mpd.status_verifikasi IS NULL OR TRIM(mpd.status_verifikasi) = '' OR LOWER(TRIM(mpd.status_verifikasi)) IN ('r0', 'draft')) THEN 1 END), 0) as SkR0,
				COALESCE(COUNT(CASE WHEN UPPER(TRIM(mpd.jenis_dokumen)) = 'SK' AND LOWER(TRIM(mpd.status_verifikasi)) IN ('belum_terverifikasi', 'belum terverifikasi', 'proses', 'r1', 'menunggu') THEN 1 END), 0) as SkR1,
				COALESCE(COUNT(CASE WHEN UPPER(TRIM(mpd.jenis_dokumen)) = 'SK' AND LOWER(TRIM(mpd.status_verifikasi)) IN ('terverifikasi', 'disetujui', 'approved', 'r2') THEN 1 END), 0) as SkR2,
				COALESCE(COUNT(CASE WHEN UPPER(TRIM(mpd.jenis_dokumen)) = 'SK' AND LOWER(TRIM(mpd.status_verifikasi)) IN ('gagal_terverifikasi', 'gagal terverifikasi', 'ditolak', 'rejected', 'revisi', 'r3') THEN 1 END), 0) as SkR3,

				COALESCE(COUNT(CASE WHEN UPPER(TRIM(mpd.jenis_dokumen)) = 'SOP' THEN 1 END), 0) as Sop,
				COALESCE(COUNT(CASE WHEN UPPER(TRIM(mpd.jenis_dokumen)) = 'SOP' AND (mpd.status_verifikasi IS NULL OR TRIM(mpd.status_verifikasi) = '' OR LOWER(TRIM(mpd.status_verifikasi)) IN ('r0', 'draft')) THEN 1 END), 0) as SopR0,
				COALESCE(COUNT(CASE WHEN UPPER(TRIM(mpd.jenis_dokumen)) = 'SOP' AND LOWER(TRIM(mpd.status_verifikasi)) IN ('belum_terverifikasi', 'belum terverifikasi', 'proses', 'r1', 'menunggu') THEN 1 END), 0) as SopR1,
				COALESCE(COUNT(CASE WHEN UPPER(TRIM(mpd.jenis_dokumen)) = 'SOP' AND LOWER(TRIM(mpd.status_verifikasi)) IN ('terverifikasi', 'disetujui', 'approved', 'r2') THEN 1 END), 0) as SopR2,
				COALESCE(COUNT(CASE WHEN UPPER(TRIM(mpd.jenis_dokumen)) = 'SOP' AND LOWER(TRIM(mpd.status_verifikasi)) IN ('gagal_terverifikasi', 'gagal terverifikasi', 'ditolak', 'rejected', 'revisi', 'r3') THEN 1 END), 0) as SopR3,

				COALESCE(COUNT(CASE WHEN UPPER(TRIM(mpd.jenis_dokumen)) IN ('PROPOSAL/TOR', 'PROPOSAL', 'TOR', 'PROPOSAL / TOR') THEN 1 END), 0) as ProposalTor,
				COALESCE(COUNT(CASE WHEN UPPER(TRIM(mpd.jenis_dokumen)) IN ('PROPOSAL/TOR', 'PROPOSAL', 'TOR', 'PROPOSAL / TOR') AND (mpd.status_verifikasi IS NULL OR TRIM(mpd.status_verifikasi) = '' OR LOWER(TRIM(mpd.status_verifikasi)) IN ('r0', 'draft')) THEN 1 END), 0) as ProposalTorR0,
				COALESCE(COUNT(CASE WHEN UPPER(TRIM(mpd.jenis_dokumen)) IN ('PROPOSAL/TOR', 'PROPOSAL', 'TOR', 'PROPOSAL / TOR') AND LOWER(TRIM(mpd.status_verifikasi)) IN ('belum_terverifikasi', 'belum terverifikasi', 'proses', 'r1', 'menunggu') THEN 1 END), 0) as ProposalTorR1,
				COALESCE(COUNT(CASE WHEN UPPER(TRIM(mpd.jenis_dokumen)) IN ('PROPOSAL/TOR', 'PROPOSAL', 'TOR', 'PROPOSAL / TOR') AND LOWER(TRIM(mpd.status_verifikasi)) IN ('terverifikasi', 'disetujui', 'approved', 'r2') THEN 1 END), 0) as ProposalTorR2,
				COALESCE(COUNT(CASE WHEN UPPER(TRIM(mpd.jenis_dokumen)) IN ('PROPOSAL/TOR', 'PROPOSAL', 'TOR', 'PROPOSAL / TOR') AND LOWER(TRIM(mpd.status_verifikasi)) IN ('gagal_terverifikasi', 'gagal terverifikasi', 'ditolak', 'rejected', 'revisi', 'r3') THEN 1 END), 0) as ProposalTorR3,

				COALESCE(COUNT(CASE WHEN UPPER(TRIM(mpd.jenis_dokumen)) = 'LAPORAN' THEN 1 END), 0) as Laporan,
				COALESCE(COUNT(CASE WHEN UPPER(TRIM(mpd.jenis_dokumen)) = 'LAPORAN' AND (mpd.status_verifikasi IS NULL OR TRIM(mpd.status_verifikasi) = '' OR LOWER(TRIM(mpd.status_verifikasi)) IN ('r0', 'draft')) THEN 1 END), 0) as LaporanR0,
				COALESCE(COUNT(CASE WHEN UPPER(TRIM(mpd.jenis_dokumen)) = 'LAPORAN' AND LOWER(TRIM(mpd.status_verifikasi)) IN ('belum_terverifikasi', 'belum terverifikasi', 'proses', 'r1', 'menunggu') THEN 1 END), 0) as LaporanR1,
				COALESCE(COUNT(CASE WHEN UPPER(TRIM(mpd.jenis_dokumen)) = 'LAPORAN' AND LOWER(TRIM(mpd.status_verifikasi)) IN ('terverifikasi', 'disetujui', 'approved', 'r2') THEN 1 END), 0) as LaporanR2,
				COALESCE(COUNT(CASE WHEN UPPER(TRIM(mpd.jenis_dokumen)) = 'LAPORAN' AND LOWER(TRIM(mpd.status_verifikasi)) IN ('gagal_terverifikasi', 'gagal terverifikasi', 'ditolak', 'rejected', 'revisi', 'r3') THEN 1 END), 0) as LaporanR3,

				COALESCE(COUNT(CASE WHEN UPPER(TRIM(mpd.jenis_dokumen)) IN ('DOCUMENT PENDUKUNG', 'DOKUMEN PENDUKUNG') THEN 1 END), 0) as DokumenPendukung,
				COALESCE(COUNT(CASE WHEN UPPER(TRIM(mpd.jenis_dokumen)) IN ('DOCUMENT PENDUKUNG', 'DOKUMEN PENDUKUNG') AND (mpd.status_verifikasi IS NULL OR TRIM(mpd.status_verifikasi) = '' OR LOWER(TRIM(mpd.status_verifikasi)) IN ('r0', 'draft')) THEN 1 END), 0) as DokumenPendukungR0,
				COALESCE(COUNT(CASE WHEN UPPER(TRIM(mpd.jenis_dokumen)) IN ('DOCUMENT PENDUKUNG', 'DOKUMEN PENDUKUNG') AND LOWER(TRIM(mpd.status_verifikasi)) IN ('belum_terverifikasi', 'belum terverifikasi', 'proses', 'r1', 'menunggu') THEN 1 END), 0) as DokumenPendukungR1,
				COALESCE(COUNT(CASE WHEN UPPER(TRIM(mpd.jenis_dokumen)) IN ('DOCUMENT PENDUKUNG', 'DOKUMEN PENDUKUNG') AND LOWER(TRIM(mpd.status_verifikasi)) IN ('terverifikasi', 'disetujui', 'approved', 'r2') THEN 1 END), 0) as DokumenPendukungR2,
				COALESCE(COUNT(CASE WHEN UPPER(TRIM(mpd.jenis_dokumen)) IN ('DOCUMENT PENDUKUNG', 'DOKUMEN PENDUKUNG') AND LOWER(TRIM(mpd.status_verifikasi)) IN ('gagal_terverifikasi', 'gagal terverifikasi', 'ditolak', 'rejected', 'revisi', 'r3') THEN 1 END), 0) as DokumenPendukungR3
			`).
			Joins("JOIN master_tahun mt ON mp.id_master_tahun = mt.id").
			Joins("JOIN v_fakultas_unit vfu ON vfu.uuid = ?", uuidFakultasUnit.String()).
			Joins("LEFT JOIN matrix_proker mpr ON mpr.id_mata_program = mp.id AND mpr.id_fakultas_unit = vfu.id").
			Joins("LEFT JOIN matrix_proker_dokumen mpd ON mpd.id_matrik_proker = mpr.id AND mpd.deleted_at IS NULL").
			Where("mt.uuid = ? AND (mpr.id IS NOT NULL OR mp.id_fakultas_unit = vfu.id OR (mp.id_fakultas_unit IS NULL AND NOT EXISTS (SELECT 1 FROM matrix_proker other_mpr WHERE other_mpr.id_mata_program = mp.id AND other_mpr.id_fakultas_unit != vfu.id)))", uuidTahun.String()).
			Group("vfu.id, mt.id, mp.id").
			Order("mp.id ASC")

		// Hitung total dengan subquery agar akurat pada query GROUP BY
		countSubQuery := r.db.WithContext(ctx).Table("mata_program mp").
			Select("mp.id").
			Joins("JOIN master_tahun mt ON mp.id_master_tahun = mt.id").
			Joins("JOIN v_fakultas_unit vfu ON vfu.uuid = ?", uuidFakultasUnit.String()).
			Joins("LEFT JOIN matrix_proker mpr ON mpr.id_mata_program = mp.id AND mpr.id_fakultas_unit = vfu.id").
			Where("mt.uuid = ? AND (mpr.id IS NOT NULL OR mp.id_fakultas_unit = vfu.id OR (mp.id_fakultas_unit IS NULL AND NOT EXISTS (SELECT 1 FROM matrix_proker other_mpr WHERE other_mpr.id_mata_program = mp.id AND other_mpr.id_fakultas_unit != vfu.id)))", uuidTahun.String()).
			Group("mp.id")

		if err := r.db.WithContext(ctx).Table("(?) as sub", countSubQuery).Count(&total).Error; err != nil {
			return nil, 0, err
		}

		// Pagination
		if page != nil && limit != nil && *limit > 0 {
			p := *page
			l := *limit
			if p < 1 {
				p = 1
			}
			offset := (p - 1) * l
			db = db.Limit(l).Offset(offset)
		}

		// Order & eksekusi query
		if err := db.Scan(&result).Error; err != nil {
			return nil, 0, err
		}

		return result, total, nil
	}

	// Default V2: data dari view v_monitoring_siproker
	db := r.db.WithContext(ctx).Table("v_monitoring_siproker ms").
		Select(`
			ms.id_fakultas_unit as FakultasUnitId,
			vfu.uuid as FakultasUnitUuid,
			vfu.nama_fak_prod_unit as FakultasUnit,
			vfu.jenjang as Jenjang,
			vfu.type as Type,
			vfu.fakultas as Fakultas,

			ms.id_mata_program as MataProgramId,
			mp.uuid as MataProgramUUID,
			mp.mata_program as MataProgram,

			ms.id_tahun as TahunId,
			mt.uuid as TahunUUID,
			mt.tahun as Tahun,
			
			ms.sk as Sk,
			ms.sk_r0 as SkR0,
			ms.sk_r1 as SkR1,
			ms.sk_r2 as SkR2,
			ms.sk_r3 as SkR3,
			
			ms.sop as Sop,
			ms.sop_r0 as SopR0,
			ms.sop_r1 as SopR1,
			ms.sop_r2 as SopR2,
			ms.sop_r3 as SopR3,
			
			ms.proposal_tor as ProposalTor,
			ms.proposal_tor_r0 as ProposalTorR0,
			ms.proposal_tor_r1 as ProposalTorR1,
			ms.proposal_tor_r2 as ProposalTorR2,
			ms.proposal_tor_r3 as ProposalTorR3,
			
			ms.laporan as Laporan,
			ms.laporan_r0 as LaporanR0,
			ms.laporan_r1 as LaporanR1,
			ms.laporan_r2 as LaporanR2,
			ms.laporan_r3 as LaporanR3,
			
			ms.dokumen_pendukung as DokumenPendukung,
			ms.dokumen_pendukung_r0 as DokumenPendukungR0,
			ms.dokumen_pendukung_r1 as DokumenPendukungR1,
			ms.dokumen_pendukung_r2 as DokumenPendukungR2,
			ms.dokumen_pendukung_r3 as DokumenPendukungR3
		`).
		Joins("JOIN v_fakultas_unit vfu ON ms.id_fakultas_unit = vfu.id").
		Joins("JOIN mata_program mp ON ms.id_mata_program = mp.id").
		Joins("JOIN master_tahun mt ON ms.id_tahun = mt.id").
		Order("ms.id_mata_program ASC").
		Where("vfu.uuid = ? AND mt.uuid = ?", uuidFakultasUnit.String(), uuidTahun.String())

	// Hitung total
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Pagination
	if page != nil && limit != nil && *limit > 0 {
		p := *page
		l := *limit
		if p < 1 {
			p = 1
		}
		offset := (p - 1) * l
		db = db.Limit(l).Offset(offset)
	}

	// Order & eksekusi query
	if err := db.Scan(&result).Error; err != nil {
		return nil, 0, err
	}

	return result, total, nil
}

// ------------------------
// GET MONITORING BY TARGET & TAHUN
// ------------------------
func (r *LaporanRepository) GetMonitoringIndikatorByIndikatorTahun(
	ctx context.Context,
	uuidIndikator uuid.UUID,
	tahun string,
	page, limit *int,
) ([]domainlaporan.MonitoringIndikator, int64, error) {
	adapter := GetLaporanAdapter(r.db, tahun)
	return adapter.GetMonitoringIndikator(ctx, uuidIndikator, tahun, page, limit)
}

