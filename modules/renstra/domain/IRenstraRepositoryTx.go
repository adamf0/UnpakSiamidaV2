package domain

import (
	commoninfra "UnpakSiamida/common/infrastructure"
	"context"
)

type IRenstraRepositoryTx interface {
	IRenstraRepository
	InsertOutbox(ctx context.Context, msg *commoninfra.OutboxMessage) error
}
