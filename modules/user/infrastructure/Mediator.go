package infrastructure

import (
    domain "UnpakSiamida/modules/user/domain"
    create "UnpakSiamida/modules/user/application/CreateUser"
    update "UnpakSiamida/modules/user/application/UpdateUser"
    delete "UnpakSiamida/modules/user/application/DeleteUser"
    get "UnpakSiamida/modules/user/application/GetUser"
    getAll "UnpakSiamida/modules/user/application/GetAllUsers"
    "github.com/mehdihadeli/go-mediatr"
    // "UnpakSiamida/modules/user/domain"
    "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func RegisterModule() {
    dsn := "root:@tcp(127.0.0.1:3306)/unpak_sijamu_server?charset=utf8mb4&parseTime=true&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

    repoUser := NewUserRepository(db)
	// if err := db.AutoMigrate(&domain.User{}); err != nil {
	// 	panic(err)
	// }

	NewUserRepository(db)

    // Pipeline behavior
    mediatr.RegisterRequestPipelineBehaviors(NewValidationBehavior())

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
        *domain.User,
    ](&get.GetUserByUuidQueryHandler{
        Repo: repoUser,
    })

    mediatr.RegisterRequestHandler[
        getAll.GetAllUsersQuery,
        domain.PagedUsers,
    ](&getAll.GetAllUsersQueryHandler{
        Repo: repoUser,
    })


}
