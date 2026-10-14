package infrastructure

import (
	commondomain "UnpakSiamida/common/domain"
	commoninfra "UnpakSiamida/common/infrastructure"
	infraFakultasUnit "UnpakSiamida/modules/fakultasunit/infrastructure"
	create "UnpakSiamida/modules/swot/application/CreateSwot"
	delete "UnpakSiamida/modules/swot/application/DeleteSwot"
	getAll "UnpakSiamida/modules/swot/application/GetAllSwots"
	get "UnpakSiamida/modules/swot/application/GetSwot"
	update "UnpakSiamida/modules/swot/application/UpdateSwot"
	domainSwot "UnpakSiamida/modules/swot/domain"

	"github.com/mehdihadeli/go-mediatr"
	"gorm.io/gorm"
)

func RegisterModuleSwot(db *gorm.DB) error {
	// if err := db.AutoMigrate(&domainSwot.Swot{}); err != nil {
	// 	return err
	// }

	repoSwot := NewSwotRepository(db)
	repoFakultasUnit := infraFakultasUnit.NewFakultasUnitRepository(db)

	mediatr.RegisterRequestHandler[
		create.CreateSwotCommand,
		string,
	](&create.CreateSwotCommandHandler{
		Repo:             repoSwot,
		RepoFakultasUnit: repoFakultasUnit,
		DB:               db,
	})

	mediatr.RegisterRequestHandler[
		update.UpdateSwotCommand,
		string,
	](&update.UpdateSwotCommandHandler{
		Repo:             repoSwot,
		RepoFakultasUnit: repoFakultasUnit,
		DB:               db,
	})

	mediatr.RegisterRequestHandler[
		delete.DeleteSwotCommand,
		string,
	](&delete.DeleteSwotCommandHandler{
		Repo: repoSwot,
	})

	mediatr.RegisterRequestHandler[
		get.GetSwotByUuidQuery,
		*domainSwot.SwotDefault,
	](&get.GetSwotByUuidQueryHandler{
		Repo: repoSwot,
	})

	mediatr.RegisterRequestHandler[
		getAll.GetAllSwotsQuery,
		commondomain.Paged[domainSwot.SwotDefault],
	](&getAll.GetAllSwotsQueryHandler{
		Repo: repoSwot,
	})

	commoninfra.RegisterValidation(create.CreateSwotCommandValidation, "SwotCreate.Validation")
	commoninfra.RegisterValidation(update.UpdateSwotCommandValidation, "SwotUpdate.Validation")
	commoninfra.RegisterValidation(delete.DeleteSwotCommandValidation, "SwotDelete.Validation")

	return nil
}
