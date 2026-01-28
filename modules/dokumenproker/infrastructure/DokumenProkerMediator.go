package infrastructure

import (
	create "UnpakSiamida/modules/dokumenproker/application/CreateDokumenProker"
	delete "UnpakSiamida/modules/dokumenproker/application/DeleteDokumenProker"
	getAll "UnpakSiamida/modules/dokumenproker/application/GetAllDokumenProkers"
	get "UnpakSiamida/modules/dokumenproker/application/GetDokumenProker"
	setupUuid "UnpakSiamida/modules/dokumenproker/application/SetupUuidDokumenProker"
	update "UnpakSiamida/modules/dokumenproker/application/UpdateDokumenProker"
	domainDokumenProker "UnpakSiamida/modules/dokumenproker/domain"
	infraFakultas "UnpakSiamida/modules/fakultasunit/infrastructure"
	infraMataProgram "UnpakSiamida/modules/mataprogram/infrastructure"

	"github.com/mehdihadeli/go-mediatr"

	// "gorm.io/driver/mysql"
	"gorm.io/gorm"
	// "fmt"
)

func RegisterModuleDokumenProker(db *gorm.DB) error {
	// dsn := "root:@tcp(127.0.0.1:3306)/unpak_sijamu_server?charset=utf8mb4&parseTime=true&loc=Local"

	// db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	// if err != nil {
	// 	return fmt.Errorf("Indikator Renstra DB connection failed: %w", err)
	// 	// panic(err)
	// }

	repoDokumenProker := NewDokumenProkerRepository(db)
	repoFakultasUnit := infraFakultas.NewFakultasUnitRepository(db)
	repoMataProgram := infraMataProgram.NewMataProgramRepository(db)
	// if err := db.AutoMigrate(&domainDokumenProker.DokumenProker{}); err != nil {
	// 	panic(err)
	// }

	// Pipeline behavior
	// mediatr.RegisterRequestPipelineBehaviors(NewValidationBehaviorDokumenProker())

	// Register request handler
	mediatr.RegisterRequestHandler[
		create.CreateDokumenProkerCommand,
		string,
	](&create.CreateDokumenProkerCommandHandler{
		Repo:             repoDokumenProker,
		RepoFakultasUnit: repoFakultasUnit,
		RepoMataProgram:  repoMataProgram,
	})

	mediatr.RegisterRequestHandler[
		update.UpdateDokumenProkerCommand,
		string,
	](&update.UpdateDokumenProkerCommandHandler{
		Repo:             repoDokumenProker,
		RepoFakultasUnit: repoFakultasUnit,
		RepoMataProgram:  repoMataProgram,
	})

	mediatr.RegisterRequestHandler[
		delete.DeleteDokumenProkerCommand,
		string,
	](&delete.DeleteDokumenProkerCommandHandler{
		Repo: repoDokumenProker,
	})

	mediatr.RegisterRequestHandler[
		get.GetDokumenProkerByUuidQuery,
		*domainDokumenProker.DokumenProker,
	](&get.GetDokumenProkerByUuidQueryHandler{
		Repo: repoDokumenProker,
	})

	mediatr.RegisterRequestHandler[
		getAll.GetAllDokumenProkersQuery,
		domainDokumenProker.PagedDokumenProkers,
	](&getAll.GetAllDokumenProkersQueryHandler{
		Repo: repoDokumenProker,
	})

	mediatr.RegisterRequestHandler[
		setupUuid.SetupUuidDokumenProkerCommand,
		string,
	](&setupUuid.SetupUuidDokumenProkerCommandHandler{
		Repo: repoDokumenProker,
	})

	return nil
}
