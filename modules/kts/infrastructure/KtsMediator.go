package infrastructure

import (
	getAll "UnpakSiamida/modules/kts/application/GetAllKtss"
	get "UnpakSiamida/modules/kts/application/GetKts"
	setupUuid "UnpakSiamida/modules/kts/application/SetupUuidKts"
	update "UnpakSiamida/modules/kts/application/UpdateKts"
	domainKts "UnpakSiamida/modules/kts/domain"
	eventKts "UnpakSiamida/modules/kts/event"

	commoninfra "UnpakSiamida/common/infrastructure"
	infraUser "UnpakSiamida/modules/user/infrastructure"

	"github.com/mehdihadeli/go-mediatr"
	"gorm.io/gorm"
)

func RegisterModuleKts(db *gorm.DB, tg commoninfra.TelegramSender) error {
	repoKts := NewKtsRepository(db)
	repoUser := infraUser.NewUserRepository(db)

	// =========================
	// Command Handlers
	// =========================

	mediatr.RegisterRequestHandler[
		update.UpdateKtsCommand,
		string,
	](&update.UpdateKtsCommandHandler{
		Repo:     repoKts,
		RepoUser: repoUser,
	})

	mediatr.RegisterRequestHandler[
		get.GetKtsByUuidQuery,
		*domainKts.Kts,
	](&get.GetKtsByUuidQueryHandler{
		Repo: repoKts,
	})

	mediatr.RegisterRequestHandler[
		getAll.GetAllKtssQuery,
		domainKts.PagedKtss,
	](&getAll.GetAllKtssQueryHandler{
		Repo: repoKts,
	})

	mediatr.RegisterRequestHandler[
		setupUuid.SetupUuidKtsCommand,
		string,
	](&setupUuid.SetupUuidKtsCommandHandler{
		Repo: repoKts,
	})

	// =========================
	// Domain Event Handler
	// =========================
	commoninfra.RegisterDomainEvent(&eventKts.KtsUpdatedEvent{})

	mediatr.RegisterNotificationHandler[eventKts.KtsUpdatedEvent](
		eventKts.NewKtsUpdatedEventHandler(tg),
	)

	commoninfra.RegisterDomainEvent(&eventKts.KtsCreatedEvent{})

	mediatr.RegisterNotificationHandler[eventKts.KtsCreatedEvent](
		eventKts.NewKtsCreatedEventHandler(tg),
	)

	return nil
}
