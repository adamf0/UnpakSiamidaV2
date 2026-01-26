package infrastructure

import (
	create "UnpakSiamida/modules/tahunproker/application/CreateTahunProker"
	delete "UnpakSiamida/modules/tahunproker/application/DeleteTahunProker"
	getAll "UnpakSiamida/modules/tahunproker/application/GetAllTahunProkers"
	get "UnpakSiamida/modules/tahunproker/application/GetTahunProker"
	setupUuid "UnpakSiamida/modules/tahunproker/application/SetupUuidTahunProker"
	update "UnpakSiamida/modules/tahunproker/application/UpdateTahunProker"
	domainTahunProker "UnpakSiamida/modules/tahunproker/domain"

	"github.com/mehdihadeli/go-mediatr"

	// "gorm.io/driver/mysql"
	"gorm.io/gorm"
	// "fmt"
)

func RegisterModuleTahunProker(db *gorm.DB) error {
	// dsn := "root:@tcp(127.0.0.1:3306)/unpak_sijamu_server?charset=utf8mb4&parseTime=true&loc=Local"

	// db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	// if err != nil {
	// 	return fmt.Errorf("Indikator Renstra DB connection failed: %w", err)
	// 	// panic(err)
	// }

	repoTahunProker := NewTahunProkerRepository(db)
	// if err := db.AutoMigrate(&domainTahunProker.TahunProker{}); err != nil {
	// 	panic(err)
	// }

	// Pipeline behavior
	// mediatr.RegisterRequestPipelineBehaviors(NewValidationBehaviorTahunProker())

	// Register request handler
	mediatr.RegisterRequestHandler[
		create.CreateTahunProkerCommand,
		string,
	](&create.CreateTahunProkerCommandHandler{
		Repo: repoTahunProker,
	})

	mediatr.RegisterRequestHandler[
		update.UpdateTahunProkerCommand,
		string,
	](&update.UpdateTahunProkerCommandHandler{
		Repo: repoTahunProker,
	})

	mediatr.RegisterRequestHandler[
		delete.DeleteTahunProkerCommand,
		string,
	](&delete.DeleteTahunProkerCommandHandler{
		Repo: repoTahunProker,
	})

	mediatr.RegisterRequestHandler[
		get.GetTahunProkerByUuidQuery,
		*domainTahunProker.TahunProker,
	](&get.GetTahunProkerByUuidQueryHandler{
		Repo: repoTahunProker,
	})

	mediatr.RegisterRequestHandler[
		getAll.GetAllTahunProkersQuery,
		domainTahunProker.PagedTahunProkers,
	](&getAll.GetAllTahunProkersQueryHandler{
		Repo: repoTahunProker,
	})

	mediatr.RegisterRequestHandler[
		setupUuid.SetupUuidTahunProkerCommand,
		string,
	](&setupUuid.SetupUuidTahunProkerCommandHandler{
		Repo: repoTahunProker,
	})

	return nil
}
