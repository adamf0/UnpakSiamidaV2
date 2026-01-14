package domain

import (
	"context"
	commoninfra "UnpakSiamida/common/infrastructure"
)

type IKtsRepositoryTx interface {
	IKtsRepository
	InsertOutbox(ctx context.Context, msg *commoninfra.OutboxMessage) error
}
