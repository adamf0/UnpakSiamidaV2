package domain

import "context"

type IRedisStore interface {
	Exists(ctx context.Context, key string) (bool, error)
	Set(ctx context.Context, key string) error
	Delete(ctx context.Context, key string) error
}
