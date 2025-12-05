package domain

import (
	"context"
	commonDomain "UnpakSiamida/common/domain"
	"github.com/google/uuid"
)

type IUserRepository interface {
	GetByUuid(ctx context.Context, uid uuid.UUID) (*User, error)
	GetAll(
        ctx context.Context,
        search string,
        searchFilters []commonDomain.SearchFilter,
        page, limit *int,
    ) ([]User, int64, error)
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, uid uuid.UUID) error
}
