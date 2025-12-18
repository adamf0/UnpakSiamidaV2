package infrastructure

import (
    // domaingeneraterenstra "UnpakSiamida/modules/generaterenstra/domain"
    generate "UnpakSiamida/modules/generaterenstra/application/GenerateRenstra"
    deletegenerate "UnpakSiamida/modules/generaterenstra/application/DeleteGenerateRenstra"
    infraRenstra "UnpakSiamida/modules/renstra/infrastructure"
    
    infraFakultasUnit "UnpakSiamida/modules/fakultasunit/infrastructure"
    infraTemplateRenstra "UnpakSiamida/modules/templaterenstra/infrastructure"
    infraTemplateDokumenTambahan "UnpakSiamida/modules/templatedokumentambahan/infrastructure"
    "github.com/mehdihadeli/go-mediatr"
    "gorm.io/driver/mysql"
	"gorm.io/gorm"
    "fmt"
)

func RegisterModuleGenerateRenstra() error{
    dsn := "root:@tcp(127.0.0.1:3306)/unpak_sijamu_server?charset=utf8mb4&parseTime=true&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("Indikator Renstra DB connection failed: %w", err)
		// panic(err)
	}

    repoGenerate := NewGenerateRenstraRepository(db)
    repoRenstra := infraRenstra.NewRenstraRepository(db)

    repoFakultasUnit := infraFakultasUnit.NewFakultasUnitRepository(db)
    repoTemplateRenstra := infraTemplateRenstra.NewTemplateRenstraRepository(db)
    repoTemplateDokumenTambahan := infraTemplateDokumenTambahan.NewTemplateDokumenTambahanRepository(db)
    // if err := db.AutoMigrate(&domainrenstra.Renstra{}); err != nil {
	// 	panic(err)
	// }

    // Pipeline behavior
    // mediatr.RegisterRequestPipelineBehaviors(NewValidationBehaviorRenstra())

    // Register request handler
    mediatr.RegisterRequestHandler[
        generate.GenerateRenstraCommand,
        string,
    ](&generate.GenerateRenstraCommandHandler{
        Repo:                           repoGenerate, //renstra_nilai & dokumen_tambahan
        RepoRenstra:                    repoRenstra,
        
        RepoFakultasUnit:               repoFakultasUnit,
        RepoTemplateRenstra:            repoTemplateRenstra,
        RepoTemplateDokumenTambahan:    repoTemplateDokumenTambahan,
    })

    mediatr.RegisterRequestHandler[
        deletegenerate.DeleteGenerateRenstraCommand,
        string,
    ](&deletegenerate.DeleteGenerateRenstraCommandHandler{
        Repo:                           repoGenerate, //renstra_nilai & dokumen_tambahan
        RepoRenstra:                    repoRenstra,
    })


    return nil
}
