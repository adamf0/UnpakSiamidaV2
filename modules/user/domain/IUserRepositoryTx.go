package domain

import (
	commoninfra "UnpakSiamida/common/infrastructure"
	"context"
)

type IUserRepositoryTx interface {
	IUserRepository
	InsertOutbox(ctx context.Context, msg *commoninfra.OutboxMessage) error
	InsertHistoryEmails(ctx context.Context, histories []HistoryEmail) error
}
