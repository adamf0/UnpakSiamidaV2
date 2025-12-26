package infrastructure

import (
    domaintemplatedokumentambahan "UnpakSiamida/modules/templatedokumentambahan/domain"
    create "UnpakSiamida/modules/templatedokumentambahan/application/CreateTemplateDokumenTambahan"
    update "UnpakSiamida/modules/templatedokumentambahan/application/UpdateTemplateDokumenTambahan"
    delete "UnpakSiamida/modules/templatedokumentambahan/application/DeleteTemplateDokumenTambahan"
    get "UnpakSiamida/modules/templatedokumentambahan/application/GetTemplateDokumenTambahan"
    getAll "UnpakSiamida/modules/templatedokumentambahan/application/GetAllTemplateDokumenTambahans"
    infraJenisFile "UnpakSiamida/modules/jenisfile/infrastructure"
    "github.com/mehdihadeli/go-mediatr"
    // "gorm.io/driver/mysql"
	"gorm.io/gorm"
    // "fmt"
)

func RegisterModuleTemplateDokumenTambahan(db *gorm.DB) error{
    // dsn := "root:@tcp(127.0.0.1:3306)/unpak_sijamu_server?charset=utf8mb4&parseTime=true&loc=Local"

	// db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	// if err != nil {
	// 	return fmt.Errorf("Standar Renstra DB connection failed: %w", err)
	// 	// panic(err)
	// }

    repoTemplateDokumenTambahan := NewTemplateDokumenTambahanRepository(db)
    repoJenisFileRepo := infraJenisFile.NewJenisFileRepository(db)
    // if err := db.AutoMigrate(&domaintemplatedokumentambahan.TemplateDokumenTambahan{}); err != nil {
	// 	panic(err)
	// }

    // Pipeline behavior
    // mediatr.RegisterRequestPipelineBehaviors(NewValidationBehaviorTemplateDokumenTambahan())

    // Register request handler
    mediatr.RegisterRequestHandler[
        create.CreateTemplateDokumenTambahanCommand,
        string,
    ](&create.CreateTemplateDokumenTambahanCommandHandler{
        Repo: repoTemplateDokumenTambahan,
        JenisFileRepo: repoJenisFileRepo,
    })

    mediatr.RegisterRequestHandler[
        update.UpdateTemplateDokumenTambahanCommand,
        string,
    ](&update.UpdateTemplateDokumenTambahanCommandHandler{
        Repo: repoTemplateDokumenTambahan,
        JenisFileRepo: repoJenisFileRepo,
    })

    mediatr.RegisterRequestHandler[
        delete.DeleteTemplateDokumenTambahanCommand,
        string,
    ](&delete.DeleteTemplateDokumenTambahanCommandHandler{
        Repo: repoTemplateDokumenTambahan,
    })

    mediatr.RegisterRequestHandler[
        get.GetTemplateDokumenTambahanByUuidQuery,
        *domaintemplatedokumentambahan.TemplateDokumenTambahan,
    ](&get.GetTemplateDokumenTambahanByUuidQueryHandler{
        Repo: repoTemplateDokumenTambahan,
    })

    mediatr.RegisterRequestHandler[
        getAll.GetAllTemplateDokumenTambahansQuery,
        domaintemplatedokumentambahan.PagedTemplateDokumenTambahans,
    ](&getAll.GetAllTemplateDokumenTambahansQueryHandler{
        Repo: repoTemplateDokumenTambahan,
    })

    return nil
}
