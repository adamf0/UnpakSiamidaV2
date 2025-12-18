package infrastructure

import (
    domainpreviewtemplate "UnpakSiamida/modules/previewtemplate/domain"
    get "UnpakSiamida/modules/previewtemplate/application/GetPreviewTemplate"
    infraFakultasUnit "UnpakSiamida/modules/fakultasunit/infrastructure"
    "github.com/mehdihadeli/go-mediatr"
    "gorm.io/driver/mysql"
	"gorm.io/gorm"
    "fmt"
)

func RegisterModulePreviewTemplate() error{
    dsn := "root:@tcp(127.0.0.1:3306)/unpak_sijamu_server?charset=utf8mb4&parseTime=true&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("Standar Renstra DB connection failed: %w", err)
		// panic(err)
	}

    repoPreviewTemplate := NewPreviewTemplateRepository(db)
    repoFakultasUnit := infraFakultasUnit.NewFakultasUnitRepository(db)
	// if err := db.AutoMigrate(&domainpreviewtemplate.PreviewTemplate{}); err != nil {
	// 	panic(err)
	// }

    // Pipeline behavior
    // mediatr.RegisterRequestPipelineBehaviors(NewValidationBehaviorPreviewTemplate())

    // Register request handler
    mediatr.RegisterRequestHandler[
        get.GetPreviewTemplateByTahunFakultasUnitQuery,
        []domainpreviewtemplate.PreviewTemplate,
    ](&get.GetPreviewTemplateByTahunFakultasUnitQueryHandler{
        Repo: repoPreviewTemplate,
        RepoFakultasUnit: repoFakultasUnit,
    })

    return nil
}
