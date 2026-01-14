package domain

import (
	"context"
)

type IAccountRepository interface {
	Auth(ctx context.Context, username string, password string) (*Account, error)
	Get(ctx context.Context, uuid string) (*Account, error)
}
