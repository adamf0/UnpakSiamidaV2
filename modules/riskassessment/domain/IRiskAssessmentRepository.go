package domain

import (
	commondomain "UnpakSiamida/common/domain"
	"context"

	"github.com/google/uuid"
)

type IRiskAssessmentRepository interface {
	GetByUuid(ctx context.Context, uid uuid.UUID) (*RiskAssessment, error)
	GetDefaultByUuid(ctx context.Context, uid uuid.UUID) (*RiskAssessmentDefault, error)
	GetAll(
		ctx context.Context,
		search string,
		searchFilters []commondomain.SearchFilter,
		page, limit *int,
	) ([]RiskAssessmentDefault, int64, error)
	Create(ctx context.Context, item *RiskAssessment) error
	Update(ctx context.Context, item *RiskAssessment) error
	Delete(ctx context.Context, uid uuid.UUID) error

	// Master Options
	GetCriteria(ctx context.Context, tipe string) ([]RiskCriteriaItem, error)
	GetLevels(ctx context.Context) ([]RiskLevelItem, error)
	GetCategories(ctx context.Context) ([]RiskCategoryItem, error)
	GetStatuses(ctx context.Context) ([]RiskStatusItem, error)
	GetPrograms(ctx context.Context, fu, thn, search string) ([]RiskMataProgramItem, error)
	GetOptions(ctx context.Context, fu, thn string) (*RiskOptionsResult, error)
}

type RiskCriteriaItem struct {
	ID        int    `json:"id" gorm:"column:id"`
	Tipe      string `json:"tipe" gorm:"column:tipe"`
	Skor      int    `json:"skor" gorm:"column:skor"`
	Label     string `json:"label" gorm:"column:label"`
	Deskripsi string `json:"deskripsi" gorm:"column:deskripsi"`
}

type RiskLevelItem struct {
	ID         int    `json:"id" gorm:"column:id"`
	Kode       string `json:"kode" gorm:"column:kode"`
	NamaLevel  string `json:"nama_level" gorm:"column:nama_level"`
	WarnaBadge string `json:"warna_badge" gorm:"column:warna_badge"`
	Keterangan string `json:"keterangan" gorm:"column:keterangan"`
	MinSkor    int    `json:"min_skor" gorm:"column:min_skor"`
	MaxSkor    int    `json:"max_skor" gorm:"column:max_skor"`
}

type RiskCategoryItem struct {
	ID        int    `json:"id" gorm:"column:id"`
	Nama      string `json:"nama" gorm:"column:nama"`
	Deskripsi string `json:"deskripsi" gorm:"column:deskripsi"`
}

type RiskStatusItem struct {
	ID        int    `json:"id" gorm:"column:id"`
	Nama      string `json:"nama" gorm:"column:nama"`
	Deskripsi string `json:"deskripsi" gorm:"column:deskripsi"`
	Urutan    int    `json:"urutan" gorm:"column:urutan"`
}

type RiskMataProgramItem struct {
	ID               int     `json:"id" gorm:"column:id"`
	UUID             string  `json:"uuid" gorm:"column:uuid"`
	MataProgram      string  `json:"mata_program" gorm:"column:mata_program"`
	IdFakultasUnit   *int    `json:"id_fakultas_unit" gorm:"column:id_fakultas_unit"`
	UuidFakultasUnit *string `json:"uuid_fakultas_unit" gorm:"column:uuid_fakultas_unit"`
	NamaFakultasUnit *string `json:"nama_fakultas_unit" gorm:"column:nama_fakultas_unit"`
	IdMasterTahun    *int    `json:"id_master_tahun" gorm:"column:id_master_tahun"`
	Tahun            *string `json:"tahun" gorm:"column:tahun"`
}

type RiskOptionsResult struct {
	Criteria    []RiskCriteriaItem    `json:"criteria"`
	Kemungkinan []RiskCriteriaItem    `json:"kemungkinan"`
	Dampak      []RiskCriteriaItem    `json:"dampak"`
	Levels      []RiskLevelItem       `json:"levels"`
	Categories  []RiskCategoryItem    `json:"categories"`
	Statuses    []RiskStatusItem      `json:"statuses"`
	Programs    []RiskMataProgramItem `json:"programs"`
}
