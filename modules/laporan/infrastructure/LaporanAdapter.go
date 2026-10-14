package infrastructure

import (
	domainlaporan "UnpakSiamida/modules/laporan/domain"
	"context"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ILaporanAdapter interface {
	GetVersion() string
	GetMonitoringIndikator(
		ctx context.Context,
		uuidIndikator uuid.UUID,
		tahun string,
		page, limit *int,
	) ([]domainlaporan.MonitoringIndikator, int64, error)
}

// LaporanAdapterV1 handles legacy data (tahun <= 2023)
type LaporanAdapterV1 struct {
	tahun string
}

func NewLaporanAdapterV1(tahun string) *LaporanAdapterV1 {
	return &LaporanAdapterV1{tahun: tahun}
}

func (a *LaporanAdapterV1) GetVersion() string {
	return "v1"
}

func (a *LaporanAdapterV1) GetMonitoringIndikator(
	ctx context.Context,
	uuidIndikator uuid.UUID,
	tahun string,
	page, limit *int,
) ([]domainlaporan.MonitoringIndikator, int64, error) {
	// Masuk adapter v1 langsung reject: masih migrasi fitur
	return nil, 0, domainlaporan.MigrasiFitur()
}

// LaporanAdapterV2 handles current Renstra data (tahun >= 2024)
type LaporanAdapterV2 struct {
	db *gorm.DB
}

func NewLaporanAdapterV2(db *gorm.DB) *LaporanAdapterV2 {
	return &LaporanAdapterV2{db: db}
}

func (a *LaporanAdapterV2) GetVersion() string {
	return "v2"
}

func (a *LaporanAdapterV2) GetMonitoringIndikator(
	ctx context.Context,
	uuidIndikator uuid.UUID,
	tahun string,
	page, limit *int,
) ([]domainlaporan.MonitoringIndikator, int64, error) {
	var result = make([]domainlaporan.MonitoringIndikator, 0)
	var total int64

	db := a.db.WithContext(ctx).
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
		mir.uuid as UndikatorUuid,
		mir.indikator as Indikator,
		mir.tipe_target as TipeTarget,
		mir.operator as Operator,
		tr.satuan as Satuan,
		tr.target as Target,
		tr.target_min as TargetMin,
		tr.target_max as TargetMax,
		CASE 
			WHEN tr.target_min IS NOT NULL AND tr.target_max IS NOT NULL 
			THEN CONCAT(tr.target_min, ' - ', tr.target_max)
			WHEN tr.target_min IS NOT NULL 
			THEN tr.target_min
			WHEN tr.target_max IS NOT NULL 
			THEN tr.target_max
			ELSE NULL
		END as RangeTarget,
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

// GetLaporanAdapter factory function returning v1 or v2 adapter
func GetLaporanAdapter(db *gorm.DB, tahun string) ILaporanAdapter {
	trimmed := strings.TrimSpace(tahun)
	tahunNum, err := strconv.Atoi(trimmed)
	if err == nil && tahunNum <= 2023 {
		return NewLaporanAdapterV1(trimmed)
	}
	return NewLaporanAdapterV2(db)
}
