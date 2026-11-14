package infrastructure

import (
	commondomain "UnpakSiamida/common/domain"
	commoninfra "UnpakSiamida/common/infrastructure"
	infraFakultasUnit "UnpakSiamida/modules/fakultasunit/infrastructure"
	create "UnpakSiamida/modules/renstra/application/CreateRenstra"
	delete "UnpakSiamida/modules/renstra/application/DeleteRenstra"
	getAll "UnpakSiamida/modules/renstra/application/GetAllRenstras"
	getDefault "UnpakSiamida/modules/renstra/application/GetRenstraDefault"
	giveCode "UnpakSiamida/modules/renstra/application/GiveCodeAccessRenstra"
	sendEmail "UnpakSiamida/modules/renstra/application/SendEmailRenstra"
	setupUuid "UnpakSiamida/modules/renstra/application/SetupUuidRenstra"
	update "UnpakSiamida/modules/renstra/application/UpdateRenstra"
	domainrenstra "UnpakSiamida/modules/renstra/domain"
	eventRenstra "UnpakSiamida/modules/renstra/event"
	infraUser "UnpakSiamida/modules/user/infrastructure"

	"github.com/mehdihadeli/go-mediatr"

	// "gorm.io/driver/mysql"
	"gorm.io/gorm"
	// "fmt"
)

func RegisterModuleRenstra(db *gorm.DB, tgs ...commoninfra.TelegramSender) error {
	var tg commoninfra.TelegramSender
	if len(tgs) > 0 {
		tg = tgs[0]
	}

	repoRenstra := NewRenstraRepository(db)
	repoFakultasUnit := infraFakultasUnit.NewFakultasUnitRepository(db)
	repoUser := infraUser.NewUserRepository(db)

	commoninfra.RegisterDomainEvent(&eventRenstra.RenstraSendEmailEvent{})

	mediatr.RegisterNotificationHandler[
		eventRenstra.RenstraSendEmailEvent,
	](eventRenstra.NewRenstraSendEmailEventHandler(tg))

	// Register request handler
	mediatr.RegisterRequestHandler[
		create.CreateRenstraCommand,
		string,
	](&create.CreateRenstraCommandHandler{
		Repo:             repoRenstra,
		FakultasUnitRepo: repoFakultasUnit,
		UserRepo:         repoUser,
	})

	mediatr.RegisterRequestHandler[
		update.UpdateRenstraCommand,
		string,
	](&update.UpdateRenstraCommandHandler{
		Repo:             repoRenstra,
		FakultasUnitRepo: repoFakultasUnit,
		UserRepo:         repoUser,
	})

	mediatr.RegisterRequestHandler[
		giveCode.GiveCodeAccessRenstraCommand,
		string,
	](&giveCode.GiveCodeAccessRenstraCommandHandler{
		Repo: repoRenstra,
	})

	mediatr.RegisterRequestHandler[
		delete.DeleteRenstraCommand,
		string,
	](&delete.DeleteRenstraCommandHandler{
		Repo: repoRenstra,
	})

	mediatr.RegisterRequestHandler[
		getDefault.GetRenstraDefaultByUuidQuery,
		*domainrenstra.RenstraDefault,
	](&getDefault.GetRenstraDefaultByUuidQueryHandler{
		Repo: repoRenstra,
	})

	mediatr.RegisterRequestHandler[
		getAll.GetAllRenstrasQuery,
		commondomain.Paged[domainrenstra.RenstraDefault],
	](&getAll.GetAllRenstrasQueryHandler{
		Repo: repoRenstra,
	})

	mediatr.RegisterRequestHandler[
		setupUuid.SetupUuidRenstraCommand,
		string,
	](&setupUuid.SetupUuidRenstraCommandHandler{
		Repo: repoRenstra,
	})

	mediatr.RegisterRequestHandler[
		sendEmail.SendEmailRenstraCommand,
		*sendEmail.SendEmailRenstraResult,
	](&sendEmail.SendEmailRenstraCommandHandler{
		Repo: repoRenstra,
	})

	commoninfra.RegisterValidation(create.CreateRenstraCommandValidation, "RenstraCreate.Validation")
	commoninfra.RegisterValidation(update.UpdateRenstraCommandValidation, "RenstraUpdate.Validation")
	commoninfra.RegisterValidation(delete.DeleteRenstraCommandValidation, "RenstraDelete.Validation")
	commoninfra.RegisterValidation(giveCode.GiveCodeAccessRenstraCommandValidation, "RenstraGiveCode.Validation")
	commoninfra.RegisterValidation(sendEmail.SendEmailRenstraCommandValidation, "RenstraSendEmail.Validation")

	return nil
}
