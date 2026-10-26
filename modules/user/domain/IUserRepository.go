package domain

import (
	commonDomain "UnpakSiamida/common/domain"
	"context"

	"github.com/google/uuid"
)

type IUserRepository interface {
	GetByUuid(ctx context.Context, uid uuid.UUID) (*User, error)
	GetAllStrict(ctx context.Context) ([]UserOptions, error)
	GetAll(
		ctx context.Context,
		search string,
		searchFilters []commonDomain.SearchFilter,
		page, limit *int,
	) ([]User, int64, error)
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, uid uuid.UUID) error
	SetupUuid(ctx context.Context) error
	MergeUser(ctx context.Context, primaryID uint, secondaryIDs []uint) error
	UnmergeUser(ctx context.Context, userID uint) error
	WithTx(ctx context.Context, fn func(txRepo IUserRepositoryTx) error) error

	// User ABAC
	GetAllUserAbac(ctx context.Context, tahun, status, role, search, userUuid, fuUuid string) ([]UserAbac, error)
	GetUserAbacByUuid(ctx context.Context, uid uuid.UUID) (*UserAbac, error)
	CreateUserAbac(ctx context.Context, abac *UserAbac) error
	UpdateUserAbac(ctx context.Context, id uint, updates map[string]interface{}) error
	DeleteUserAbac(ctx context.Context, uid uuid.UUID) error
	CheckUserAbacExists(ctx context.Context, userID uint, tahun string, fuID *uint) (bool, error)
	GetFakultasUnitInfo(ctx context.Context, id *uint, uidStr string) (idRes uint, typeRes string, namaRes string, err error)
}
