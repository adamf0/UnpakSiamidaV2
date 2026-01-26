package infrastructure

import (
	create "UnpakSiamida/modules/mataprogram/application/CreateMataProgram"
	delete "UnpakSiamida/modules/mataprogram/application/DeleteMataProgram"
	getAll "UnpakSiamida/modules/mataprogram/application/GetAllMataPrograms"
	get "UnpakSiamida/modules/mataprogram/application/GetMataProgram"
	setupUuid "UnpakSiamida/modules/mataprogram/application/SetupUuidMataProgram"
	update "UnpakSiamida/modules/mataprogram/application/UpdateMataProgram"
	domainMataProgram "UnpakSiamida/modules/mataprogram/domain"

	infraTahunProker "UnpakSiamida/modules/tahunproker/infrastructure"

	"github.com/mehdihadeli/go-mediatr"

	// "gorm.io/driver/mysql"
	"gorm.io/gorm"
	// "fmt"
)

func RegisterModuleMataProgram(db *gorm.DB) error {
	// dsn := "root:@tcp(127.0.0.1:3306)/unpak_sijamu_server?charset=utf8mb4&parseTime=true&loc=Local"

	// db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	// if err != nil {
	// 	return fmt.Errorf("Indikator Renstra DB connection failed: %w", err)
	// 	// panic(err)
	// }

	repoMataProgram := NewMataProgramRepository(db)
	repoTahunProker := infraTahunProker.NewTahunProkerRepository(db)
	// if err := db.AutoMigrate(&domainMataProgram.MataProgram{}); err != nil {
	// 	panic(err)
	// }

	// Pipeline behavior
	// mediatr.RegisterRequestPipelineBehaviors(NewValidationBehaviorMataProgram())

	// Register request handler
	mediatr.RegisterRequestHandler[
		create.CreateMataProgramCommand,
		string,
	](&create.CreateMataProgramCommandHandler{
		Repo:            repoMataProgram,
		RepoTahunProker: repoTahunProker,
	})

	mediatr.RegisterRequestHandler[
		update.UpdateMataProgramCommand,
		string,
	](&update.UpdateMataProgramCommandHandler{
		Repo:            repoMataProgram,
		RepoTahunProker: repoTahunProker,
	})

	mediatr.RegisterRequestHandler[
		delete.DeleteMataProgramCommand,
		string,
	](&delete.DeleteMataProgramCommandHandler{
		Repo: repoMataProgram,
	})

	mediatr.RegisterRequestHandler[
		get.GetMataProgramByUuidQuery,
		*domainMataProgram.MataProgram,
	](&get.GetMataProgramByUuidQueryHandler{
		Repo: repoMataProgram,
	})

	mediatr.RegisterRequestHandler[
		getAll.GetAllMataProgramsQuery,
		domainMataProgram.PagedMataPrograms,
	](&getAll.GetAllMataProgramsQueryHandler{
		Repo: repoMataProgram,
	})

	mediatr.RegisterRequestHandler[
		setupUuid.SetupUuidMataProgramCommand,
		string,
	](&setupUuid.SetupUuidMataProgramCommandHandler{
		Repo: repoMataProgram,
	})

	return nil
}
