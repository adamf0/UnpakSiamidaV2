package infrastructure

import (
	commondomain "UnpakSiamida/common/domain"
	commoninfra "UnpakSiamida/common/infrastructure"
	infraFakultas "UnpakSiamida/modules/fakultasunit/infrastructure"
	create "UnpakSiamida/modules/jadwalproker/application/CreateJadwalProker"
	delete "UnpakSiamida/modules/jadwalproker/application/DeleteJadwalProker"
	getAll "UnpakSiamida/modules/jadwalproker/application/GetAllJadwalProkers"
	get "UnpakSiamida/modules/jadwalproker/application/GetJadwalProker"
	setupUuid "UnpakSiamida/modules/jadwalproker/application/SetupUuidJadwalProker"
	update "UnpakSiamida/modules/jadwalproker/application/UpdateJadwalProker"
	domainJadwalProker "UnpakSiamida/modules/jadwalproker/domain"

	"github.com/mehdihadeli/go-mediatr"

	// "gorm.io/driver/mysql"
	"gorm.io/gorm"
	// "fmt"
)

func RegisterModuleJadwalProker(db *gorm.DB) error {
	// dsn := "root:@tcp(127.0.0.1:3306)/unpak_sijamu_server?charset=utf8mb4&parseTime=true&loc=Local"

	// db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	// if err != nil {
	// 	return fmt.Errorf("Indikator Renstra DB connection failed: %w", err)
	// 	// panic(err)
	// }

	repoJadwalProker := NewJadwalProkerRepository(db)
	repoFakultasUnit := infraFakultas.NewFakultasUnitRepository(db)
	// if err := db.AutoMigrate(&domainJadwalProker.JadwalProker{}); err != nil {
	// 	panic(err)
	// }

	// Pipeline behavior
	// mediatr.RegisterRequestPipelineBehaviors(NewValidationBehaviorJadwalProker())

	// Register request handler
	mediatr.RegisterRequestHandler[
		create.CreateJadwalProkerCommand,
		string,
	](&create.CreateJadwalProkerCommandHandler{
		Repo:             repoJadwalProker,
		RepoFakultasUnit: repoFakultasUnit,
	})

	mediatr.RegisterRequestHandler[
		update.UpdateJadwalProkerCommand,
		string,
	](&update.UpdateJadwalProkerCommandHandler{
		Repo:             repoJadwalProker,
		RepoFakultasUnit: repoFakultasUnit,
	})

	mediatr.RegisterRequestHandler[
		delete.DeleteJadwalProkerCommand,
		string,
	](&delete.DeleteJadwalProkerCommandHandler{
		Repo: repoJadwalProker,
	})

	mediatr.RegisterRequestHandler[
		get.GetJadwalProkerByUuidQuery,
		*domainJadwalProker.JadwalProker,
	](&get.GetJadwalProkerByUuidQueryHandler{
		Repo: repoJadwalProker,
	})

	mediatr.RegisterRequestHandler[
		getAll.GetAllJadwalProkersQuery,
		commondomain.Paged[domainJadwalProker.JadwalProkerDefault],
	](&getAll.GetAllJadwalProkersQueryHandler{
		Repo: repoJadwalProker,
	})

	mediatr.RegisterRequestHandler[
		setupUuid.SetupUuidJadwalProkerCommand,
		string,
	](&setupUuid.SetupUuidJadwalProkerCommandHandler{
		Repo: repoJadwalProker,
	})

	commoninfra.RegisterValidation(create.CreateJadwalProkerCommandValidation, "JadwalProkerCreate.Validation")
	commoninfra.RegisterValidation(update.UpdateJadwalProkerCommandValidation, "JadwalProkerUpdate.Validation")
	commoninfra.RegisterValidation(delete.DeleteJadwalProkerCommandValidation, "JadwalProkerDelete.Validation")

	return nil
}
