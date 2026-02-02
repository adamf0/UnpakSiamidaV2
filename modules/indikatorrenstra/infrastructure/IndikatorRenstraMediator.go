package infrastructure

import (
	create "UnpakSiamida/modules/indikatorrenstra/application/CreateIndikatorRenstra"
	delete "UnpakSiamida/modules/indikatorrenstra/application/DeleteIndikatorRenstra"
	getAll "UnpakSiamida/modules/indikatorrenstra/application/GetAllIndikatorRenstras"
	get "UnpakSiamida/modules/indikatorrenstra/application/GetIndikatorRenstra"
	setupUuid "UnpakSiamida/modules/indikatorrenstra/application/SetupUuidIndikatorRenstra"
	getTree "UnpakSiamida/modules/indikatorrenstra/application/TreeIndikatorRenstra"
	update "UnpakSiamida/modules/indikatorrenstra/application/UpdateIndikatorRenstra"
	domainindikatorrenstra "UnpakSiamida/modules/indikatorrenstra/domain"

	infraStandarRenstra "UnpakSiamida/modules/standarrenstra/infrastructure"

	"github.com/mehdihadeli/go-mediatr"

	// "gorm.io/driver/mysql"
	"gorm.io/gorm"
	// "fmt"
)

func RegisterModuleIndikatorRenstra(db *gorm.DB) error {
	// dsn := "root:@tcp(127.0.0.1:3306)/unpak_sijamu_server?charset=utf8mb4&parseTime=true&loc=Local"

	// db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	// if err != nil {
	// 	return fmt.Errorf("Indikator Renstra DB connection failed: %w", err)
	// 	// panic(err)
	// }

	repoIndikatorRenstra := NewIndikatorRenstraRepository(db)
	repoStandarRenstra := infraStandarRenstra.NewStandarRenstraRepository(db)
	// if err := db.AutoMigrate(&domainindikatorrenstra.IndikatorRenstra{}); err != nil {
	// 	panic(err)
	// }

	// Register request handler
	mediatr.RegisterRequestHandler[
		create.CreateIndikatorRenstraCommand,
		string,
	](&create.CreateIndikatorRenstraCommandHandler{
		Repo:               repoIndikatorRenstra,
		RepoStandarRenstra: repoStandarRenstra,
	})

	mediatr.RegisterRequestHandler[
		update.UpdateIndikatorRenstraCommand,
		string,
	](&update.UpdateIndikatorRenstraCommandHandler{
		Repo:               repoIndikatorRenstra,
		RepoStandarRenstra: repoStandarRenstra,
	})

	mediatr.RegisterRequestHandler[
		delete.DeleteIndikatorRenstraCommand,
		string,
	](&delete.DeleteIndikatorRenstraCommandHandler{
		Repo: repoIndikatorRenstra,
	})

	mediatr.RegisterRequestHandler[
		get.GetIndikatorRenstraByUuidQuery,
		*domainindikatorrenstra.IndikatorRenstra,
	](&get.GetIndikatorRenstraByUuidQueryHandler{
		Repo: repoIndikatorRenstra,
	})

	mediatr.RegisterRequestHandler[
		getAll.GetAllIndikatorRenstrasQuery,
		domainindikatorrenstra.PagedIndikatorRenstras,
	](&getAll.GetAllIndikatorRenstrasQueryHandler{
		Repo: repoIndikatorRenstra,
	})

	mediatr.RegisterRequestHandler[
		getTree.TreeIndikatorRenstraQuery,
		[]domainindikatorrenstra.IndikatorTree,
	](&getTree.TreeIndikatorRenstraQueryHandler{
		Repo: repoIndikatorRenstra,
	})

	mediatr.RegisterRequestHandler[
		setupUuid.SetupUuidIndikatorRenstraCommand,
		string,
	](&setupUuid.SetupUuidIndikatorRenstraCommandHandler{
		Repo: repoIndikatorRenstra,
	})

	return nil
}
