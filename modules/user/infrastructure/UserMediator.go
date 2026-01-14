package infrastructure

import (
    domainuser "UnpakSiamida/modules/user/domain"
    create "UnpakSiamida/modules/user/application/CreateUser"
    update "UnpakSiamida/modules/user/application/UpdateUser"
    delete "UnpakSiamida/modules/user/application/DeleteUser"
    get "UnpakSiamida/modules/user/application/GetUser"
    getAll "UnpakSiamida/modules/user/application/GetAllUsers"
    setupUuid "UnpakSiamida/modules/user/application/SetupUuidUser"
    "github.com/mehdihadeli/go-mediatr"
    // "gorm.io/driver/mysql"
	"gorm.io/gorm"
    // "fmt"
)

func RegisterModuleUser(db *gorm.DB) error {
    // dsn := "root:@tcp(127.0.0.1:3306)/unpak_sijamu_server?charset=utf8mb4&parseTime=true&loc=Local"

	// db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	// if err != nil {
    //     return fmt.Errorf("User DB connection failed: %w", err)
	// 	// panic(err)
	// }

    repoUser := NewUserRepository(db)
	// if err := db.AutoMigrate(&domainuser.User{}); err != nil {
	// 	panic(err)
	// }

    // Pipeline behavior
    // mediatr.RegisterRequestPipelineBehaviors(NewValidationBehaviorUser())

    // Register request handler
    mediatr.RegisterRequestHandler[
        create.CreateUserCommand,
        string,
    ](&create.CreateUserCommandHandler{
        Repo: repoUser,
    })

    mediatr.RegisterRequestHandler[
        update.UpdateUserCommand,
        string,
    ](&update.UpdateUserCommandHandler{
        Repo: repoUser,
    })

    mediatr.RegisterRequestHandler[
        delete.DeleteUserCommand,
        string,
    ](&delete.DeleteUserCommandHandler{
        Repo: repoUser,
    })

    mediatr.RegisterRequestHandler[
        get.GetUserByUuidQuery,
        *domainuser.User,
    ](&get.GetUserByUuidQueryHandler{
        Repo: repoUser,
    })

    mediatr.RegisterRequestHandler[
        getAll.GetAllUsersQuery,
        domainuser.PagedUsers,
    ](&getAll.GetAllUsersQueryHandler{
        Repo: repoUser,
    })

    mediatr.RegisterRequestHandler[
        setupUuid.SetupUuidUserCommand,
        string,
    ](&setupUuid.SetupUuidUserCommandHandler{
        Repo: repoUser,
    })

    return nil
}
