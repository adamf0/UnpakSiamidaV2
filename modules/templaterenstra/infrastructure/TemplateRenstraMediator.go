package infrastructure

import (
	commondomain "UnpakSiamida/common/domain"
	infraFakultasUnit "UnpakSiamida/modules/fakultasunit/infrastructure"
	infraIndikatorRenstra "UnpakSiamida/modules/indikatorrenstra/infrastructure"
	create "UnpakSiamida/modules/templaterenstra/application/CreateTemplateRenstra"
	delete "UnpakSiamida/modules/templaterenstra/application/DeleteTemplateRenstra"
	getAll "UnpakSiamida/modules/templaterenstra/application/GetAllTemplateRenstras"
	get "UnpakSiamida/modules/templaterenstra/application/GetTemplateRenstra"
	setupUuid "UnpakSiamida/modules/templaterenstra/application/SetupUuidTemplateRenstra"
	update "UnpakSiamida/modules/templaterenstra/application/UpdateTemplateRenstra"
	domaintemplaterenstra "UnpakSiamida/modules/templaterenstra/domain"

	"github.com/mehdihadeli/go-mediatr"

	// "gorm.io/driver/mysql"
	"gorm.io/gorm"
	// "fmt"
)

func RegisterModuleTemplateRenstra(db *gorm.DB) error {
	// dsn := "root:@tcp(127.0.0.1:3306)/unpak_sijamu_server?charset=utf8mb4&parseTime=true&loc=Local"

	// db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	// if err != nil {
	// 	return fmt.Errorf("Standar Renstra DB connection failed: %w", err)
	// 	// panic(err)
	// }

	repoTemplateRenstra := NewTemplateRenstraRepository(db)
	repoFakultasUnitRepo := infraFakultasUnit.NewFakultasUnitRepository(db)
	repoIndikatorRenstra := infraIndikatorRenstra.NewIndikatorRenstraRepository(db)
	// if err := db.AutoMigrate(&domaintemplaterenstra.TemplateRenstra{}); err != nil {
	// 	panic(err)
	// }

	// Pipeline behavior
	// mediatr.RegisterRequestPipelineBehaviors(NewValidationBehaviorTemplateRenstra())

	// Register request handler
	mediatr.RegisterRequestHandler[
		create.CreateTemplateRenstraCommand,
		string,
	](&create.CreateTemplateRenstraCommandHandler{
		Repo:                 repoTemplateRenstra,
		FakultasUnitRepo:     repoFakultasUnitRepo,
		IndikatorRenstraRepo: repoIndikatorRenstra,
	})

	mediatr.RegisterRequestHandler[
		update.UpdateTemplateRenstraCommand,
		string,
	](&update.UpdateTemplateRenstraCommandHandler{
		Repo:                 repoTemplateRenstra,
		FakultasUnitRepo:     repoFakultasUnitRepo,
		IndikatorRenstraRepo: repoIndikatorRenstra,
	})

	mediatr.RegisterRequestHandler[
		delete.DeleteTemplateRenstraCommand,
		string,
	](&delete.DeleteTemplateRenstraCommandHandler{
		Repo: repoTemplateRenstra,
	})

	mediatr.RegisterRequestHandler[
		get.GetTemplateRenstraByUuidQuery,
		*domaintemplaterenstra.TemplateRenstra,
	](&get.GetTemplateRenstraByUuidQueryHandler{
		Repo: repoTemplateRenstra,
	})

	mediatr.RegisterRequestHandler[
		getAll.GetAllTemplateRenstrasQuery,
		commondomain.Paged[domaintemplaterenstra.TemplateRenstraDefault],
	](&getAll.GetAllTemplateRenstrasQueryHandler{
		Repo: repoTemplateRenstra,
	})

	mediatr.RegisterRequestHandler[
		setupUuid.SetupUuidTemplateRenstraCommand,
		string,
	](&setupUuid.SetupUuidTemplateRenstraCommandHandler{
		Repo: repoTemplateRenstra,
	})

	return nil
}
