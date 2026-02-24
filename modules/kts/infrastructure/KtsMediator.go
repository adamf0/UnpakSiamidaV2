package infrastructure

import (
	export "UnpakSiamida/modules/kts/application/ExportKts"
	getAll "UnpakSiamida/modules/kts/application/GetAllKtss"
	get "UnpakSiamida/modules/kts/application/GetKts"
	setupUuid "UnpakSiamida/modules/kts/application/SetupUuidKts"
	update "UnpakSiamida/modules/kts/application/UpdateKts"
	domainKts "UnpakSiamida/modules/kts/domain"
	eventKts "UnpakSiamida/modules/kts/event"

	commonDomain "UnpakSiamida/common/domain"
	commoninfra "UnpakSiamida/common/infrastructure"
	infraUser "UnpakSiamida/modules/user/infrastructure"

	"github.com/mehdihadeli/go-mediatr"
	"gorm.io/gorm"
)

func RegisterModuleKts(db *gorm.DB, redis *commonDomain.IRedisStore, tg commoninfra.TelegramSender) error {
	repoKts := NewKtsRepository(db)
	repoUser := infraUser.NewUserRepository(db)

	// =========================
	// Command Handlers
	// =========================

	mediatr.RegisterRequestHandler[
		export.PublishKtsCommand,
		string,
	](&export.PublishKtsCommandHandler{
		Repo: repoKts,
	})

	mediatr.RegisterRequestHandler[
		export.ExportKtsCommand,
		[]byte,
	](&export.ExportKtsCommandHandler{
		Repo:  repoKts,
		Redis: *redis,
	})

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
		commonDomain.Paged[domainKts.KtsDefault],
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

	commoninfra.RegisterDomainEvent(&eventKts.KtsPdfRequestedEvent{})

	mediatr.RegisterNotificationHandler[eventKts.KtsPdfRequestedEvent](
		eventKts.NewKtsPdfRequestedEventHandler(*redis),
	)

	commoninfra.RegisterValidation(export.PublishKtsCommandValidation, "KtsPublish.Validation")
	commoninfra.RegisterValidation(update.UpdateKtsCommandValidation, "KtsUpdate.Validation")
	commoninfra.RegisterValidation(export.ExportKtsCommandValidation, "KtsExport.Validation")

	return nil
}
