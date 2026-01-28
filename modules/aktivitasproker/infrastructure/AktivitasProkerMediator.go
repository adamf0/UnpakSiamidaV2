package infrastructure

import (
	create "UnpakSiamida/modules/aktivitasproker/application/CreateAktivitasProker"
	delete "UnpakSiamida/modules/aktivitasproker/application/DeleteAktivitasProker"
	get "UnpakSiamida/modules/aktivitasproker/application/GetAktivitasProker"
	getAll "UnpakSiamida/modules/aktivitasproker/application/GetAllAktivitasProkers"
	setupUuid "UnpakSiamida/modules/aktivitasproker/application/SetupUuidAktivitasProker"
	update "UnpakSiamida/modules/aktivitasproker/application/UpdateAktivitasProker"
	domainAktivitasProker "UnpakSiamida/modules/aktivitasproker/domain"
	infraFakultas "UnpakSiamida/modules/fakultasunit/infrastructure"
	infraMataProgram "UnpakSiamida/modules/mataprogram/infrastructure"

	"github.com/mehdihadeli/go-mediatr"

	// "gorm.io/driver/mysql"
	"gorm.io/gorm"
	// "fmt"
)

func RegisterModuleAktivitasProker(db *gorm.DB) error {
	// dsn := "root:@tcp(127.0.0.1:3306)/unpak_sijamu_server?charset=utf8mb4&parseTime=true&loc=Local"

	// db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	// if err != nil {
	// 	return fmt.Errorf("Indikator Renstra DB connection failed: %w", err)
	// 	// panic(err)
	// }

	repoAktivitasProker := NewAktivitasProkerRepository(db)
	repoFakultasUnit := infraFakultas.NewFakultasUnitRepository(db)
	repoMataProgram := infraMataProgram.NewMataProgramRepository(db)
	// if err := db.AutoMigrate(&domainAktivitasProker.AktivitasProker{}); err != nil {
	// 	panic(err)
	// }

	// Pipeline behavior
	// mediatr.RegisterRequestPipelineBehaviors(NewValidationBehaviorAktivitasProker())

	// Register request handler
	mediatr.RegisterRequestHandler[
		create.CreateAktivitasProkerCommand,
		string,
	](&create.CreateAktivitasProkerCommandHandler{
		Repo:             repoAktivitasProker,
		RepoFakultasUnit: repoFakultasUnit,
		RepoMataProgram:  repoMataProgram,
	})

	mediatr.RegisterRequestHandler[
		update.UpdateAktivitasProkerCommand,
		string,
	](&update.UpdateAktivitasProkerCommandHandler{
		Repo:             repoAktivitasProker,
		RepoFakultasUnit: repoFakultasUnit,
		RepoMataProgram:  repoMataProgram,
	})

	mediatr.RegisterRequestHandler[
		delete.DeleteAktivitasProkerCommand,
		string,
	](&delete.DeleteAktivitasProkerCommandHandler{
		Repo: repoAktivitasProker,
	})

	mediatr.RegisterRequestHandler[
		get.GetAktivitasProkerByUuidQuery,
		*domainAktivitasProker.AktivitasProker,
	](&get.GetAktivitasProkerByUuidQueryHandler{
		Repo: repoAktivitasProker,
	})

	mediatr.RegisterRequestHandler[
		getAll.GetAllAktivitasProkersQuery,
		domainAktivitasProker.PagedAktivitasProkers,
	](&getAll.GetAllAktivitasProkersQueryHandler{
		Repo: repoAktivitasProker,
	})

	mediatr.RegisterRequestHandler[
		setupUuid.SetupUuidAktivitasProkerCommand,
		string,
	](&setupUuid.SetupUuidAktivitasProkerCommandHandler{
		Repo: repoAktivitasProker,
	})

	return nil
}
