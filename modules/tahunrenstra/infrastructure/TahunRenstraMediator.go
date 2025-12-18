package infrastructure

import (
    domainTahunRenstra "UnpakSiamida/modules/tahunrenstra/domain"
    get "UnpakSiamida/modules/tahunrenstra/application/GetActiveTahunRenstra"
    getAll "UnpakSiamida/modules/tahunrenstra/application/GetAllTahunRenstras"
    "github.com/mehdihadeli/go-mediatr"
    "gorm.io/driver/mysql"
	"gorm.io/gorm"
    "fmt"
)

func RegisterModuleTahunRenstra() error{
    dsn := "root:@tcp(127.0.0.1:3306)/unpak_sijamu_server?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=true&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("Tahun Renstra DB connection failed: %w", err)
		// panic(err)
	}

    repoTahunRenstra := NewTahunRenstraRepository(db)
	// if err := db.AutoMigrate(&domainTahunRenstra.TahunRenstra{}); err != nil {
	// 	panic(err)
	// }

    // Pipeline behavior
    // mediatr.RegisterRequestPipelineBehaviors(NewValidationBehaviorTahunRenstra())

    // Register request handler
    mediatr.RegisterRequestHandler[
        get.GetActiveTahunRenstraQuery,
        *domainTahunRenstra.TahunRenstra,
    ](&get.GetActiveTahunRenstraQueryHandler{
        Repo: repoTahunRenstra,
    })

    mediatr.RegisterRequestHandler[
        getAll.GetAllTahunRenstrasQuery,
        domainTahunRenstra.PagedTahunRenstras,
    ](&getAll.GetAllTahunRenstrasQueryHandler{
        Repo: repoTahunRenstra,
    })

    return nil
}
