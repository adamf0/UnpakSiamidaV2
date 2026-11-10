package domain

import (
	"context"
)

type IAccountRepository interface {
	Auth(ctx context.Context, username string, password string) (*Account, error)
	Get(ctx context.Context, uuid string) (*Account, error)
	IsUsernameTakenByOther(ctx context.Context, username string, excludeUserUUID string) (bool, error)
	UpdateProfile(ctx context.Context, userUUID string, username string, name string, email *string, password *string) error
}
