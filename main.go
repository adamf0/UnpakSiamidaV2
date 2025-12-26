package main

import (
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

import (
	"github.com/gofiber/fiber/v2/middleware/cors"

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

	generaterenstraInfrastructure "UnpakSiamida/modules/generaterenstra/infrastructure"
	generaterenstraPresentation "UnpakSiamida/modules/generaterenstra/presentation"

	previewtemplateInfrastructure "UnpakSiamida/modules/previewtemplate/infrastructure"
	previewtemplatePresentation "UnpakSiamida/modules/previewtemplate/presentation"
	
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

	generateRenstra "UnpakSiamida/modules/generaterenstra/application/GenerateRenstra"
	deleteGenerateRenstra "UnpakSiamida/modules/generaterenstra/application/DeleteGenerateRenstra"

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
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders: "*",
	}))
	app.Use(commonpresentation.LoggerMiddleware)
	// app.Use(commonpresentation.HeaderSecurityMiddleware(cfg))

	mediatr.RegisterRequestPipelineBehaviors(NewValidationBehavior())

	var db *gorm.DB
	mustStart("Database", func() error {
		var err error
		db, err = NewMySQL()
		return err
	})

	//berlaku untuk startup bukan hot reload
	mustStart("User Module", func() error {
		return userInfrastructure.RegisterModuleUser(db)
	})

	mustStart("Standar Renstra Module", func() error {
		return standarrenstraInfrastructure.RegisterModuleStandarRenstra(db)
	})

	mustStart("Indikator Renstra Module", func() error {
		return indikatorrenstraInfrastructure.RegisterModuleIndikatorRenstra(db)
	})

	mustStart("Tahun Renstra Module", func() error {
		return tahunrenstraInfrastructure.RegisterModuleTahunRenstra(db)
	})

	mustStart("Template Renstra Module", func() error {
		return templaterenstraInfrastructure.RegisterModuleTemplateRenstra(db)
	})

	mustStart("Template Dokumen Tambahan Module", func() error {
		return templatedokumentambahanInfrastructure.RegisterModuleTemplateDokumenTambahan(db)
	})

	mustStart("Fakultas Unit Module", func() error {
		return fakultasunitInfrastructure.RegisterModuleFakultasUnit(db)
	})

	mustStart("Jenis File Module", func() error {
		return jenisfileInfrastructure.RegisterModuleJenisFile(db)
	})

	mustStart("Renstra Module", func() error {
		return renstraInfrastructure.RegisterModuleRenstra(db)
	})

	mustStart("Generate Renstra Module", func() error {
		return generaterenstraInfrastructure.RegisterModuleGenerateRenstra(db)
	})

	mustStart("Preview Template Module", func() error {
		return previewtemplateInfrastructure.RegisterModulePreviewTemplate(db)
	})

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
	generaterenstraPresentation.ModuleGenerateRenstra(app)
	previewtemplatePresentation.ModulePreviewTemplate(app)

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
				return nil, wrapValidationError("UserCreate.Validation", err)
			}
		case updateUser.UpdateUserCommand:
			if err := updateUser.UpdateUserCommandValidation(cmd); err != nil {
				return nil, wrapValidationError("UseUpdate.Validation", err)
			}
		case deleteUser.DeleteUserCommand:
			if err := deleteUser.DeleteUserCommandValidation(cmd); err != nil {
				return nil, wrapValidationError("UserDelete.Validation", err)
			}

		// === StandarRenstra Commands ===
		case createStandarRenstra.CreateStandarRenstraCommand:
			if err := createStandarRenstra.CreateStandarRenstraCommandValidation(cmd); err != nil {
				return nil, wrapValidationError("StandarRenstraCreate.Validation", err)
			}
		case updateStandarRenstra.UpdateStandarRenstraCommand:
			if err := updateStandarRenstra.UpdateStandarRenstraCommandValidation(cmd); err != nil {
				return nil, wrapValidationError("StandarRenstraUpdate.Validation", err)
			}
		case deleteStandarRenstra.DeleteStandarRenstraCommand:
			if err := deleteStandarRenstra.DeleteStandarRenstraCommandValidation(cmd); err != nil {
				return nil, wrapValidationError("StandarRenstraDelete.Validation", err)
			}
		// === IndikatorRenstra Commands ===
		case createIndikatorRenstra.CreateIndikatorRenstraCommand:
			if err := createIndikatorRenstra.CreateIndikatorRenstraCommandValidation(cmd); err != nil {
				return nil, wrapValidationError("IndikatorRenstraCreate.Validation", err)
			}
		case updateIndikatorRenstra.UpdateIndikatorRenstraCommand:
			if err := updateIndikatorRenstra.UpdateIndikatorRenstraCommandValidation(cmd); err != nil {
				return nil, wrapValidationError("IndikatorRenstraUpdate.Validation", err)
			}
		case deleteIndikatorRenstra.DeleteIndikatorRenstraCommand:
			if err := deleteIndikatorRenstra.DeleteIndikatorRenstraCommandValidation(cmd); err != nil {
				return nil, wrapValidationError("IndikatorRenstraDelete.Validation", err)
			}
		// === TemplateRenstra Commands ===
		case createTemplateRenstra.CreateTemplateRenstraCommand:
			if err := createTemplateRenstra.CreateTemplateRenstraCommandValidation(cmd); err != nil {
				return nil, wrapValidationError("TemplateRenstraCreate.Validation", err)
			}
		case updateTemplateRenstra.UpdateTemplateRenstraCommand:
			if err := updateTemplateRenstra.UpdateTemplateRenstraCommandValidation(cmd); err != nil {
				return nil, wrapValidationError("TemplateRenstraUpdate.Validation", err)
			}
		case deleteTemplateRenstra.DeleteTemplateRenstraCommand:
			if err := deleteTemplateRenstra.DeleteTemplateRenstraCommandValidation(cmd); err != nil {
				return nil, wrapValidationError("TemplateRenstraDelete.Validation", err)
			}
		// === TemplateDokumenTambahan Commands ===
		case createTemplateDokumenTambahan.CreateTemplateDokumenTambahanCommand:
			if err := createTemplateDokumenTambahan.CreateTemplateDokumenTambahanCommandValidation(cmd); err != nil {
				return nil, wrapValidationError("TemplateDokumenTambahanCreate.Validation", err)
			}
		case updateTemplateDokumenTambahan.UpdateTemplateDokumenTambahanCommand:
			if err := updateTemplateDokumenTambahan.UpdateTemplateDokumenTambahanCommandValidation(cmd); err != nil {
				return nil, wrapValidationError("TemplateDokumenTambahanUpdate.Validation", err)
			}
		case deleteTemplateDokumenTambahan.DeleteTemplateDokumenTambahanCommand:
			if err := deleteTemplateDokumenTambahan.DeleteTemplateDokumenTambahanCommandValidation(cmd); err != nil {
				return nil, wrapValidationError("TemplateDokumenTambahanDelete.Validation", err)
			}

		// === Renstra Commands ===
		case createRenstra.CreateRenstraCommand:
			if err := createRenstra.CreateRenstraCommandValidation(cmd); err != nil {
				return nil, wrapValidationError("RenstraCreate.Validation", err)
			}
		case updateRenstra.UpdateRenstraCommand:
			if err := updateRenstra.UpdateRenstraCommandValidation(cmd); err != nil {
				return nil, wrapValidationError("RenstraUpdate.Validation", err)
			}
		case giveCodeRenstra.GiveCodeAccessRenstraCommand:
			if err := giveCodeRenstra.GiveCodeAccessRenstraCommandValidation(cmd); err != nil {
				return nil, wrapValidationError("RenstraGiveCode.Validation", err)
			}
		case deleteRenstra.DeleteRenstraCommand:
			if err := deleteRenstra.DeleteRenstraCommandValidation(cmd); err != nil {
				return nil, wrapValidationError("RenstraDelete.Validation", err)
			}

		// === Generate Renstra Commands ===
		case generateRenstra.GenerateRenstraCommand:
			if err := generateRenstra.GenerateRenstraCommandValidation(cmd); err != nil {
				return nil, wrapValidationError("GenerateRenstra.Validation", err)
			}

		case deleteGenerateRenstra.DeleteGenerateRenstraCommand:
			if err := deleteGenerateRenstra.DeleteGenerateRenstraCommandValidation(cmd); err != nil {
				return nil, wrapValidationError("DeleteGenerateRenstra.Validation", err)
			}

		// === Preview Template Commands ===
		// case previewTemplate.GetPreviewTemplateCommand:
		// 	if err := previewTemplate.GetPreviewTemplateCommandValidation(cmd); err != nil {
		// 		return nil, wrapValidationError("PreviewTemplate.Validation", err)
		// 	}

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

var (
	db   *gorm.DB
	once sync.Once
)

func NewMySQL() (*gorm.DB, error) {
	var err error

	once.Do(func() {
		dsn := "root:@tcp(127.0.0.1:3306)/unpak_sijamu_server?charset=utf8mb4&parseTime=true&loc=Local"

		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			return
		}

		sqlDB, _ := db.DB()

		sqlDB.SetMaxOpenConns(20)
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetConnMaxLifetime(60 * time.Minute)
		sqlDB.SetConnMaxIdleTime(10 * time.Minute)
	})

	return db, err
}
