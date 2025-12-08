package main

import (
	userInfrastructure "UnpakSiamida/modules/user/infrastructure"
	userPresentation "UnpakSiamida/modules/user/presentation"

	standarrenstraInfrastructure "UnpakSiamida/modules/standarrenstra/infrastructure"
	standarrenstraPresentation "UnpakSiamida/modules/standarrenstra/presentation"

	indikatorrenstraInfrastructure "UnpakSiamida/modules/indikatorrenstra/infrastructure"
	indikatorrenstraPresentation "UnpakSiamida/modules/indikatorrenstra/presentation"

	tahunrenstraInfrastructure "UnpakSiamida/modules/tahunrenstra/infrastructure"
	tahunrenstraPresentation "UnpakSiamida/modules/tahunrenstra/presentation"
	
	createUser "UnpakSiamida/modules/user/application/CreateUser"
	updateUser "UnpakSiamida/modules/user/application/UpdateUser"
	deleteUser "UnpakSiamida/modules/user/application/DeleteUser"

	createStandarRenstra "UnpakSiamida/modules/standarrenstra/application/CreateStandarRenstra"
	updateStandarRenstra "UnpakSiamida/modules/standarrenstra/application/UpdateStandarRenstra"
	deleteStandarRenstra "UnpakSiamida/modules/standarrenstra/application/DeleteStandarRenstra"

	createIndikatorRenstra "UnpakSiamida/modules/indikatorrenstra/application/CreateIndikatorRenstra"
	updateIndikatorRenstra "UnpakSiamida/modules/indikatorrenstra/application/UpdateIndikatorRenstra"
	deleteIndikatorRenstra "UnpakSiamida/modules/indikatorrenstra/application/DeleteIndikatorRenstra"

	"github.com/gofiber/fiber/v2"
	"context"
	"strings"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/go-ozzo/ozzo-validation/v4"
	commoninfra "UnpakSiamida/common/infrastructure"
	commonpresentation "UnpakSiamida/common/presentation"
)

var startupErrors []fiber.Map

func mustStart(name string, fn func() error) {
    if err := fn(); err != nil {
        startupErrors = append(startupErrors, fiber.Map{
            "module": name,
            "error":  err.Error(),
        })
    }
}

func main() {
	cfg := commonpresentation.DefaultHeaderSecurityConfig()
	cfg.ResolveAndCheck = false

	app := fiber.New(fiber.Config{
		// DisableStartupMessage: true,
		ReadBufferSize: 16 * 1024,
		Prefork:      true, // gunakan semua CPU cores
		ServerHeader: "Fiber",
		// ReadTimeout: 10 * time.Second,
		// WriteTimeout: 10 * time.Second,
		// IdleTimeout: 10 * time.Second
	})
	app.Use(commonpresentation.LoggerMiddleware)
	// app.Use(commonpresentation.HeaderSecurityMiddleware(cfg))

	mediatr.RegisterRequestPipelineBehaviors(NewValidationBehavior())

	//berlaku untuk startup bukan hot reload
	mustStart("User Module", userInfrastructure.RegisterModuleUser)
	mustStart("Standar Renstra Module", standarrenstraInfrastructure.RegisterModuleStandarRenstra)
	mustStart("Indikator Renstra Module", indikatorrenstraInfrastructure.RegisterModuleIndikatorRenstra)
	mustStart("Tahun Renstra Module", tahunrenstraInfrastructure.RegisterModuleTahunRenstra)

	if len(startupErrors) > 0 {
		app.Use(func(c *fiber.Ctx) error {
			return c.Status(500).JSON(fiber.Map{
					"Code":    "INTERNAL_SERVER_ERROR",
					"Message": "Startup module failed",
					"Trace": startupErrors,
			})
		})
	}
	
	userPresentation.ModuleUser(app)
	standarrenstraPresentation.ModuleStandarRenstra(app)
	indikatorrenstraPresentation.ModuleIndikatorRenstra(app)
	tahunrenstraPresentation.ModuleTahunRenstra(app)

	app.Listen(":3000")
}

type ValidationBehavior struct{}

func NewValidationBehavior() *ValidationBehavior {
	return &ValidationBehavior{}
}

func (b *ValidationBehavior) Handle(
	ctx context.Context,
	request interface{},
	next mediatr.RequestHandlerFunc,
) (interface{}, error) {

	switch cmd := request.(type) {
		// === User Commands ===
		case createUser.CreateUserCommand:
			if err := createUser.CreateUserCommandValidation(cmd); err != nil {
				return nil, wrapValidationError("User.Validation", err)
			}
		case updateUser.UpdateUserCommand:
			if err := updateUser.UpdateUserCommandValidation(cmd); err != nil {
				return nil, wrapValidationError("User.Validation", err)
			}
		case deleteUser.DeleteUserCommand:
			if err := deleteUser.DeleteUserCommandValidation(cmd); err != nil {
				return nil, wrapValidationError("User.Validation", err)
			}

		// === StandarRenstra Commands ===
		case createStandarRenstra.CreateStandarRenstraCommand:
			if err := createStandarRenstra.CreateStandarRenstraCommandValidation(cmd); err != nil {
				return nil, wrapValidationError("StandarRenstra.Validation", err)
			}
		case updateStandarRenstra.UpdateStandarRenstraCommand:
			if err := updateStandarRenstra.UpdateStandarRenstraCommandValidation(cmd); err != nil {
				return nil, wrapValidationError("StandarRenstra.Validation", err)
			}
		case deleteStandarRenstra.DeleteStandarRenstraCommand:
			if err := deleteStandarRenstra.DeleteStandarRenstraCommandValidation(cmd); err != nil {
				return nil, wrapValidationError("StandarRenstra.Validation", err)
			}
		// === IndikatorRenstra Commands ===
		case createIndikatorRenstra.CreateIndikatorRenstraCommand:
			if err := createIndikatorRenstra.CreateIndikatorRenstraCommandValidation(cmd); err != nil {
				return nil, wrapValidationError("IndikatorRenstra.Validation", err)
			}
		case updateIndikatorRenstra.UpdateIndikatorRenstraCommand:
			if err := updateIndikatorRenstra.UpdateIndikatorRenstraCommandValidation(cmd); err != nil {
				return nil, wrapValidationError("IndikatorRenstra.Validation", err)
			}
		case deleteIndikatorRenstra.DeleteIndikatorRenstraCommand:
			if err := deleteIndikatorRenstra.DeleteIndikatorRenstraCommandValidation(cmd); err != nil {
				return nil, wrapValidationError("IndikatorRenstra.Validation", err)
			}

		default:
			// request lain → skip validation
	}

	return next(ctx)
}

func wrapValidationError(code string, err error) error {
	if ve, ok := err.(validation.Errors); ok {
		msgs := make(map[string]string)
		for field, ferr := range ve {
			key := strings.ToLower(field)
			msgs[key] = ferr.Error()
		}
		return commoninfra.NewResponseError(code, msgs)
	}
	return commoninfra.NewResponseError(code, err.Error())
}