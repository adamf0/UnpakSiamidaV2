package domain

import (
	"context"
	// commonDomain "UnpakSiamida/common/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IGenerateRenstraRepository interface {
	BeginTx(ctx context.Context) (*gorm.DB, error)
	Commit() error
	Rollback()

	// GetAllByTahunFakUnit(ctx context.Context, tahun string, fakultasUnit uint) ([]GenerateRenstra, error)
	GetAllRenstraNilaiByTahunFakUnitDefault(ctx context.Context, tahun string, fakultasUnit uint) ([]GenerateRenstraDefault, error)
	GetAllDokumenTambahanByTahunFakUnitDefault(ctx context.Context, tahun string, fakultasUnit uint) ([]GenerateDokumenTambahanDefault, error)
	
	CreateRenstraNilai(ctx context.Context, tx *gorm.DB, generateRenstra *GenerateRenstra) error
	DeleteRenstraNilai(ctx context.Context, tx *gorm.DB, generateRenstra *GenerateRenstra) error
	ForceDeleteRenstraNilai(ctx context.Context, uid uuid.UUID, renstra uint) error

	CreateDokumenTambahan(ctx context.Context, tx *gorm.DB, generateDokumenTamabahn *GenerateDokumenTambahan) error
	DeleteDokumenTambahan(ctx context.Context, tx *gorm.DB, generateDokumenTamabahn *GenerateDokumenTambahan) error
	ForceDeleteDokumenTambahan(ctx context.Context, uid uuid.UUID, renstra uint) error
}
