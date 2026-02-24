package infrastructure

import (
	commondomain "UnpakSiamida/common/domain"
	create "UnpakSiamida/modules/standarrenstra/application/CreateStandarRenstra"
	delete "UnpakSiamida/modules/standarrenstra/application/DeleteStandarRenstra"
	getAll "UnpakSiamida/modules/standarrenstra/application/GetAllStandarRenstras"
	get "UnpakSiamida/modules/standarrenstra/application/GetStandarRenstra"
	setupUuid "UnpakSiamida/modules/standarrenstra/application/SetupUuidStandarRenstra"
	update "UnpakSiamida/modules/standarrenstra/application/UpdateStandarRenstra"
	domainstandarrenstra "UnpakSiamida/modules/standarrenstra/domain"

	"github.com/mehdihadeli/go-mediatr"

	// "gorm.io/driver/mysql"
	"gorm.io/gorm"
	// "fmt"
)

func RegisterModuleStandarRenstra(db *gorm.DB) error {
	// dsn := "root:@tcp(127.0.0.1:3306)/unpak_sijamu_server?charset=utf8mb4&parseTime=true&loc=Local"

	// db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	// if err != nil {
	// 	return fmt.Errorf("Standar Renstra DB connection failed: %w", err)
	// 	// panic(err)
	// }

	repoStandarRenstra := NewStandarRenstraRepository(db)
	// if err := db.AutoMigrate(&domainstandarrenstra.StandarRenstra{}); err != nil {
	// 	panic(err)
	// }

	// Pipeline behavior
	// mediatr.RegisterRequestPipelineBehaviors(NewValidationBehaviorStandarRenstra())

	// Register request handler
	mediatr.RegisterRequestHandler[
		create.CreateStandarRenstraCommand,
		string,
	](&create.CreateStandarRenstraCommandHandler{
		Repo: repoStandarRenstra,
	})

	mediatr.RegisterRequestHandler[
		update.UpdateStandarRenstraCommand,
		string,
	](&update.UpdateStandarRenstraCommandHandler{
		Repo: repoStandarRenstra,
	})

	mediatr.RegisterRequestHandler[
		delete.DeleteStandarRenstraCommand,
		string,
	](&delete.DeleteStandarRenstraCommandHandler{
		Repo: repoStandarRenstra,
	})

	mediatr.RegisterRequestHandler[
		get.GetStandarRenstraByUuidQuery,
		*domainstandarrenstra.StandarRenstra,
	](&get.GetStandarRenstraByUuidQueryHandler{
		Repo: repoStandarRenstra,
	})

	mediatr.RegisterRequestHandler[
		getAll.GetAllStandarRenstrasQuery,
		commondomain.Paged[domainstandarrenstra.StandarRenstra],
	](&getAll.GetAllStandarRenstrasQueryHandler{
		Repo: repoStandarRenstra,
	})

	mediatr.RegisterRequestHandler[
		setupUuid.SetupUuidStandarRenstraCommand,
		string,
	](&setupUuid.SetupUuidStandarRenstraCommandHandler{
		Repo: repoStandarRenstra,
	})

	return nil
}
