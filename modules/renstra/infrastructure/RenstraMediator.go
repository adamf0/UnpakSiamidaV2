package infrastructure

import (
	infraFakultasUnit "UnpakSiamida/modules/fakultasunit/infrastructure"
	create "UnpakSiamida/modules/renstra/application/CreateRenstra"
	delete "UnpakSiamida/modules/renstra/application/DeleteRenstra"
	getAll "UnpakSiamida/modules/renstra/application/GetAllRenstras"
	getDefault "UnpakSiamida/modules/renstra/application/GetRenstraDefault"
	giveCode "UnpakSiamida/modules/renstra/application/GiveCodeAccessRenstra"
	setupUuid "UnpakSiamida/modules/renstra/application/SetupUuidRenstra"
	update "UnpakSiamida/modules/renstra/application/UpdateRenstra"
	domainrenstra "UnpakSiamida/modules/renstra/domain"
	infraUser "UnpakSiamida/modules/user/infrastructure"

	"github.com/mehdihadeli/go-mediatr"

	// "gorm.io/driver/mysql"
	"gorm.io/gorm"
	// "fmt"
)

func RegisterModuleRenstra(db *gorm.DB) error {
	// dsn := "root:@tcp(127.0.0.1:3306)/unpak_sijamu_server?charset=utf8mb4&parseTime=true&loc=Local"

	// db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	// if err != nil {
	// 	return fmt.Errorf("Indikator Renstra DB connection failed: %w", err)
	// 	// panic(err)
	// }

	repoRenstra := NewRenstraRepository(db)
	repoFakultasUnit := infraFakultasUnit.NewFakultasUnitRepository(db)
	repoUser := infraUser.NewUserRepository(db)
	// if err := db.AutoMigrate(&domainrenstra.Renstra{}); err != nil {
	// 	panic(err)
	// }

	// Pipeline behavior
	// mediatr.RegisterRequestPipelineBehaviors(NewValidationBehaviorRenstra())

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
		domainrenstra.PagedRenstras,
	](&getAll.GetAllRenstrasQueryHandler{
		Repo: repoRenstra,
	})

	mediatr.RegisterRequestHandler[
		setupUuid.SetupUuidRenstraCommand,
		string,
	](&setupUuid.SetupUuidRenstraCommandHandler{
		Repo: repoRenstra,
	})

	return nil
}
