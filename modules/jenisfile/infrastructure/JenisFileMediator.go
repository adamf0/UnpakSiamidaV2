package infrastructure

import (
    domainJenisFile "UnpakSiamida/modules/jenisfile/domain"
    create "UnpakSiamida/modules/jenisfile/application/CreateJenisFile"
    update "UnpakSiamida/modules/jenisfile/application/UpdateJenisFile"
    delete "UnpakSiamida/modules/jenisfile/application/DeleteJenisFile"
    get "UnpakSiamida/modules/jenisfile/application/GetJenisFile"
    getAll "UnpakSiamida/modules/jenisfile/application/GetAllJenisFiles"
    setupUuid "UnpakSiamida/modules/jenisfile/application/SetupUuidJenisFile"
    "github.com/mehdihadeli/go-mediatr"
    // "gorm.io/driver/mysql"
	"gorm.io/gorm"
    // "fmt"
)

func RegisterModuleJenisFile(db *gorm.DB) error{
    // dsn := "root:@tcp(127.0.0.1:3306)/unpak_sijamu_server?charset=utf8mb4&parseTime=true&loc=Local"

	// db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	// if err != nil {
	// 	return fmt.Errorf("Indikator Renstra DB connection failed: %w", err)
	// 	// panic(err)
	// }

    repoJenisFile := NewJenisFileRepository(db)
	// if err := db.AutoMigrate(&domainJenisFile.JenisFile{}); err != nil {
	// 	panic(err)
	// }

    // Pipeline behavior
    // mediatr.RegisterRequestPipelineBehaviors(NewValidationBehaviorJenisFile())

    // Register request handler
    mediatr.RegisterRequestHandler[
        create.CreateJenisFileCommand,
        string,
    ](&create.CreateJenisFileCommandHandler{
        Repo: repoJenisFile,
    })

    mediatr.RegisterRequestHandler[
        update.UpdateJenisFileCommand,
        string,
    ](&update.UpdateJenisFileCommandHandler{
        Repo: repoJenisFile,
    })

    mediatr.RegisterRequestHandler[
        delete.DeleteJenisFileCommand,
        string,
    ](&delete.DeleteJenisFileCommandHandler{
        Repo: repoJenisFile,
    })

    mediatr.RegisterRequestHandler[
        get.GetJenisFileByUuidQuery,
        *domainJenisFile.JenisFile,
    ](&get.GetJenisFileByUuidQueryHandler{
        Repo: repoJenisFile,
    })

    mediatr.RegisterRequestHandler[
        getAll.GetAllJenisFilesQuery,
        domainJenisFile.PagedJenisFiles,
    ](&getAll.GetAllJenisFilesQueryHandler{
        Repo: repoJenisFile,
    })

    mediatr.RegisterRequestHandler[
        setupUuid.SetupUuidJenisFileCommand,
        string,
    ](&setupUuid.SetupUuidJenisFileCommandHandler{
        Repo: repoJenisFile,
    })

    return nil
}
