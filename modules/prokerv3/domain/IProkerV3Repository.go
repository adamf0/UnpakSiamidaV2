package domain

import (
	commondomain "UnpakSiamida/common/domain"
	"context"
	"time"

	"github.com/google/uuid"
)

type SoftDeleteProkerParams struct {
	UUID          uuid.UUID
	Catatan       string
	Role          string
	ActorID       *int
	ActorName     string
	ActorUsername string
}

type IProkerV3Repository interface {
	GetByUuid(ctx context.Context, uid uuid.UUID) (*ProkerV3, error)
	GetDefaultByUuid(ctx context.Context, uid uuid.UUID) (*ProkerV3Default, error)
	GetAll(ctx context.Context, search string, searchFilters []commondomain.SearchFilter, page, limit *int) ([]ProkerV3Default, int64, error)
	Create(ctx context.Context, data *ProkerV3) error
	Update(ctx context.Context, data *ProkerV3) error
	Delete(ctx context.Context, uid uuid.UUID) error
	SoftDelete(ctx context.Context, params SoftDeleteProkerParams) error

	// Beban Yayasan & Batal
	UpdateBebanYayasan(ctx context.Context, uid uuid.UUID, beban bool) (*ProkerV3, error)
	UpdateBatalProker(ctx context.Context, uid uuid.UUID, isBatal bool, catatan string) (*ProkerV3, error)

	// Options & Target & COA
	GetOptions(ctx context.Context) (*ProkerOptionsResult, error)
	GetTargetUnpak(ctx context.Context, tahun, indikatorParam, unitParam string) (*TargetUnpakResult, error)
	GetCoa(ctx context.Context, search string) ([]CoaItem, error)

	// Status & History
	UpdateStatus(ctx context.Context, req UpdateStatusRequest) (*UpdateStatusResult, error)
	GetHistoryState(ctx context.Context, entityUuid string) ([]HistoryState, error)
	RecordHistoryState(ctx context.Context, h *HistoryState) error

	// Dokumen & Aktivitas Options
	GetDokumenByMatrikId(ctx context.Context, matrikID uint) ([]MatrixProkerDokumen, error)
	GetAktivitasOptionsByMatrikId(ctx context.Context, matrikID uint) ([]AktivitasDokumenOption, error)
	GetDokumenByUuid(ctx context.Context, uid uuid.UUID) (*MatrixProkerDokumen, error)
	CreateDokumen(ctx context.Context, doc *MatrixProkerDokumen) error
	DeleteDokumen(ctx context.Context, uid uuid.UUID) error
	UpdateVerifikasiDokumen(ctx context.Context, uid uuid.UUID, status, catatan, catatanPenolakan string) (*MatrixProkerDokumen, error)

	// App Settings & Proker Batal
	EnsureAppSettingsTable(ctx context.Context) error
	GetProkerBatalSettingStatus(ctx context.Context, now time.Time, targetUnitID *int) (*ProkerBatalSettingStatus, error)
	GetProkerBatalSettingRaw(ctx context.Context) (string, error)
	SaveProkerBatalSetting(ctx context.Context, serialized string, updatedBy *int) error
	GetFakultasUnitNamesMap(ctx context.Context, unitIDs []int) (map[int]string, error)
	GetUserIDByUUID(ctx context.Context, uuid string) (*int, error)
	GetFakultasUnitIDByUUIDOrName(ctx context.Context, key string) (*int, error)

	// Jadwal Proker V3
	GetAllJadwalProkerV3(ctx context.Context, filter JadwalFilter) ([]JadwalProkerV3DTO, error)
	GetJadwalProkerV3(ctx context.Context, uid uuid.UUID) (*JadwalProkerV3DTO, error)
	CheckJadwalProkerV3(ctx context.Context, req CheckJadwalRequest) (*CheckJadwalResult, error)
	CreateJadwalProkerV3(ctx context.Context, data *JadwalProkerV3) error
	UpdateJadwalProkerV3(ctx context.Context, uid uuid.UUID, data *JadwalProkerV3) error
	DeleteJadwalProkerV3(ctx context.Context, uid uuid.UUID) error

	// Pagu Anggaran
	GetAllPaguAnggaran(ctx context.Context, tahun string, idFu string) ([]PaguAnggaran, error)
	SavePaguAnggaran(ctx context.Context, req SavePaguAnggaranRequest) (*PaguAnggaran, error)
	DeletePaguAnggaran(ctx context.Context, id string) error
	GetUnitPaguSummary(ctx context.Context, tahun string, reqFu, userFu string) (*UnitPaguSummaryResult, error)

	// Dashboard Monitoring & Sankey
	GetDashboardProkerMonitoring(ctx context.Context, filter DashboardMonitoringFilter) (*DashboardMonitoringResult, error)
	GetSankeyData(ctx context.Context, tahun string) (*SankeyDataResult, error)
}
