package infrastructure

import (
	domainlaporan "UnpakSiamida/modules/laporan/domain"
	"context"

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
	var result = make([]domainlaporan.MonitoringIndikator, 0)
	var total int64

	db := r.db.WithContext(ctx).
		Table("renstra_nilai rn").
		Select(`
		r.tahun as Tahun,
		r.fakultas_unit as FakultasUnitId,
		vfu.uuid as FakultasUnitUuid,
		vfu.nama_fak_prod_unit as FakultasUnit,
		vfu.jenjang as Jenjang,
		vfu.type as Type,
		vfu.fakultas as Fakultas,
		mir.id as IndikatorId,
		mir.uuid as IndikatorUuid,
		mir.indikator as Indikator,
		mir.tipe_target as TipeTarget,
		rn.capaian as Capaian,
		rn.capaian_auditor as CapaianAuditor
	`).
		Joins("JOIN renstra r ON rn.id_renstra = r.id").
		Joins("JOIN template_renstra tr ON rn.template_renstra = tr.id").
		Joins("JOIN master_indikator_renstra mir ON tr.indikator = mir.id").
		Joins("JOIN v_fakultas_unit vfu ON r.fakultas_unit = vfu.id").
		Where("mir.uuid = ? AND r.tahun = ?", uuidIndikator.String(), tahun).
		Order("vfu.id ASC")

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
