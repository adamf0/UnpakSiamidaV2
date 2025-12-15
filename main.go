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

	templaterenstraInfrastructure "UnpakSiamida/modules/templaterenstra/infrastructure"
	templaterenstraPresentation "UnpakSiamida/modules/templaterenstra/presentation"

	templatedokumentambahanInfrastructure "UnpakSiamida/modules/templatedokumentambahan/infrastructure"
	templatedokumentambahanPresentation "UnpakSiamida/modules/templatedokumentambahan/presentation"

	fakultasunitInfrastructure "UnpakSiamida/modules/fakultasunit/infrastructure"
	fakultasunitPresentation "UnpakSiamida/modules/fakultasunit/presentation"

	jenisfileInfrastructure "UnpakSiamida/modules/jenisfile/infrastructure"
	jenisfilePresentation "UnpakSiamida/modules/jenisfile/presentation"

	renstraInfrastructure "UnpakSiamida/modules/renstra/infrastructure"
	renstraPresentation "UnpakSiamida/modules/renstra/presentation"
	
	createUser "UnpakSiamida/modules/user/application/CreateUser"
	updateUser "UnpakSiamida/modules/user/application/UpdateUser"
	deleteUser "UnpakSiamida/modules/user/application/DeleteUser"

	createStandarRenstra "UnpakSiamida/modules/standarrenstra/application/CreateStandarRenstra"
	updateStandarRenstra "UnpakSiamida/modules/standarrenstra/application/UpdateStandarRenstra"
	deleteStandarRenstra "UnpakSiamida/modules/standarrenstra/application/DeleteStandarRenstra"

	createIndikatorRenstra "UnpakSiamida/modules/indikatorrenstra/application/CreateIndikatorRenstra"
	updateIndikatorRenstra "UnpakSiamida/modules/indikatorrenstra/application/UpdateIndikatorRenstra"
	deleteIndikatorRenstra "UnpakSiamida/modules/indikatorrenstra/application/DeleteIndikatorRenstra"

	createTemplateRenstra "UnpakSiamida/modules/templaterenstra/application/CreateTemplateRenstra"
	updateTemplateRenstra "UnpakSiamida/modules/templaterenstra/application/UpdateTemplateRenstra"
	deleteTemplateRenstra "UnpakSiamida/modules/templaterenstra/application/DeleteTemplateRenstra"

	createTemplateDokumenTambahan "UnpakSiamida/modules/templatedokumentambahan/application/CreateTemplateDokumenTambahan"
	updateTemplateDokumenTambahan "UnpakSiamida/modules/templatedokumentambahan/application/UpdateTemplateDokumenTambahan"
	deleteTemplateDokumenTambahan "UnpakSiamida/modules/templatedokumentambahan/application/DeleteTemplateDokumenTambahan"

	createRenstra "UnpakSiamida/modules/renstra/application/CreateRenstra"
	updateRenstra "UnpakSiamida/modules/renstra/application/UpdateRenstra"
	giveCodeRenstra "UnpakSiamida/modules/renstra/application/GiveCodeAccessRenstra"
	deleteRenstra "UnpakSiamida/modules/renstra/application/DeleteRenstra"

	"github.com/gofiber/fiber/v2"
	"context"
	"strings"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/go-ozzo/ozzo-validation/v4"
	commoninfra "UnpakSiamida/common/infrastructure"
	commonpresentation "UnpakSiamida/common/presentation"

	"github.com/gofiber/swagger"      // ⬅️ INI WAJIB
	_ "github.com/swaggo/files"       // ⬅️ swagger embed files
	_ "UnpakSiamida/docs"              // ⬅️ hasil swag init
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

// @title UnpakSiamidaV2 API
// @version 1.0
// @description All Module Siamida
// @host localhost:3000
// @BasePath /
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
	mustStart("Template Renstra Module", templaterenstraInfrastructure.RegisterModuleTemplateRenstra)
	mustStart("Template Dokumen Tambahan Module", templatedokumentambahanInfrastructure.RegisterModuleTemplateDokumenTambahan)
	mustStart("Fakultas Unit Module", fakultasunitInfrastructure.RegisterModuleFakultasUnit)
	mustStart("Jenis File Module", jenisfileInfrastructure.RegisterModuleJenisFile)
	mustStart("Renstra Module", renstraInfrastructure.RegisterModuleRenstra)

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
	templaterenstraPresentation.ModuleTemplateRenstra(app)
	templatedokumentambahanPresentation.ModuleTemplateDokumenTambahan(app)
	fakultasunitPresentation.ModuleFakultasUnit(app)
	jenisfilePresentation.ModuleJenisFile(app)
	renstraPresentation.ModuleRenstra(app)

	app.Get("/swagger/*", swagger.HandlerDefault)
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
		// === TemplateRenstra Commands ===
		case createTemplateRenstra.CreateTemplateRenstraCommand:
			if err := createTemplateRenstra.CreateTemplateRenstraCommandValidation(cmd); err != nil {
				return nil, wrapValidationError("TemplateRenstra.Validation", err)
			}
		case updateTemplateRenstra.UpdateTemplateRenstraCommand:
			if err := updateTemplateRenstra.UpdateTemplateRenstraCommandValidation(cmd); err != nil {
				return nil, wrapValidationError("TemplateRenstra.Validation", err)
			}
		case deleteTemplateRenstra.DeleteTemplateRenstraCommand:
			if err := deleteTemplateRenstra.DeleteTemplateRenstraCommandValidation(cmd); err != nil {
				return nil, wrapValidationError("TemplateRenstra.Validation", err)
			}
		// === TemplateDokumenTambahan Commands ===
		case createTemplateDokumenTambahan.CreateTemplateDokumenTambahanCommand:
			if err := createTemplateDokumenTambahan.CreateTemplateDokumenTambahanCommandValidation(cmd); err != nil {
				return nil, wrapValidationError("TemplateDokumenTambahan.Validation", err)
			}
		case updateTemplateDokumenTambahan.UpdateTemplateDokumenTambahanCommand:
			if err := updateTemplateDokumenTambahan.UpdateTemplateDokumenTambahanCommandValidation(cmd); err != nil {
				return nil, wrapValidationError("TemplateDokumenTambahan.Validation", err)
			}
		case deleteTemplateDokumenTambahan.DeleteTemplateDokumenTambahanCommand:
			if err := deleteTemplateDokumenTambahan.DeleteTemplateDokumenTambahanCommandValidation(cmd); err != nil {
				return nil, wrapValidationError("TemplateDokumenTambahan.Validation", err)
			}

		// === Renstra Commands ===
		case createRenstra.CreateRenstraCommand:
			if err := createRenstra.CreateRenstraCommandValidation(cmd); err != nil {
				return nil, wrapValidationError("Renstra.Validation", err)
			}
		case updateRenstra.UpdateRenstraCommand:
			if err := updateRenstra.UpdateRenstraCommandValidation(cmd); err != nil {
				return nil, wrapValidationError("Renstra.Validation", err)
			}
		case giveCodeRenstra.GiveCodeAccessRenstraCommand:
			if err := giveCodeRenstra.GiveCodeAccessRenstraCommandValidation(cmd); err != nil {
				return nil, wrapValidationError("Renstra.Validation", err)
			}
		case deleteRenstra.DeleteRenstraCommand:
			if err := deleteRenstra.DeleteRenstraCommandValidation(cmd); err != nil {
				return nil, wrapValidationError("Renstra.Validation", err)
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