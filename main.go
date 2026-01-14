package main

import (
	"context"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mehdihadeli/go-mediatr"

	"github.com/gofiber/fiber/v2/middleware/cors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

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

	accountInfrastructure "UnpakSiamida/modules/account/infrastructure"

	accountPresentation "UnpakSiamida/modules/account/presentation"

	renstranilaiInfrastructure "UnpakSiamida/modules/renstranilai/infrastructure"

	renstranilaiPresentation "UnpakSiamida/modules/renstranilai/presentation"

	dokumentambahanInfrastructure "UnpakSiamida/modules/dokumentambahan/infrastructure"

	dokumentambahanPresentation "UnpakSiamida/modules/dokumentambahan/presentation"

	ktsInfrastructure "UnpakSiamida/modules/kts/infrastructure"

	ktsPresentation "UnpakSiamida/modules/kts/presentation"

	login "UnpakSiamida/modules/account/application/Login"

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

	updateRenstraNilai "UnpakSiamida/modules/renstranilai/application/UpdateRenstraNilai"

	deleteRenstraNilai "UnpakSiamida/modules/renstranilai/application/DeleteRenstraNilai"

	updateDokumenTambahan "UnpakSiamida/modules/dokumentambahan/application/UpdateDokumenTambahan"

	deleteDokumenTambahan "UnpakSiamida/modules/dokumentambahan/application/DeleteDokumenTambahan"

	updateKts "UnpakSiamida/modules/kts/application/UpdateKts"

	generateRenstra "UnpakSiamida/modules/generaterenstra/application/GenerateRenstra"

	deleteGenerateRenstra "UnpakSiamida/modules/generaterenstra/application/DeleteGenerateRenstra"

	/////////

	validation "github.com/go-ozzo/ozzo-validation/v4"

	commoninfra "UnpakSiamida/common/infrastructure"

	commonpresentation "UnpakSiamida/common/presentation"

	//////////

	eventKts "UnpakSiamida/modules/kts/event"

	_ "UnpakSiamida/docs"

	"github.com/gofiber/swagger"
	_ "github.com/swaggo/files"
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
		// Prefork:        true, // gunakan semua CPU cores
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

	mustStart("Account Module", func() error {
		return accountInfrastructure.RegisterModuleAccount(db)
	})

	mustStart("Renstra Nilai Module", func() error { //buat audit
		return renstranilaiInfrastructure.RegisterModuleRenstraNilai(db)
	})

	mustStart("Dokumen Tambahan Module", func() error { //buat audit
		return dokumentambahanInfrastructure.RegisterModuleDokumenTambahan(db)
	})

	mustStart("Kts Module", func() error { //buat audit
		tg := commoninfra.NewTelegramClient(
			// os.Getenv("TELEGRAM_BOT_TOKEN"),
			"8598253529:AAHx8G4uAIi0Ws5kGjYLKDFhBeKwrO9kbaQ",
			// os.Getenv("TELEGRAM_CHAT_ID"),
			"6327475613",
		)

		return ktsInfrastructure.RegisterModuleKts(db, tg)
	})

	if len(startupErrors) > 0 {
		app.Use(func(c *fiber.Ctx) error {
			return c.Status(500).JSON(fiber.Map{
				"Code":    "INTERNAL_SERVER_ERROR",
				"Message": "Startup module failed",
				"Trace":   startupErrors,
			})
		})
	}

	dispatcher := commoninfra.NewEventDispatcher()
	commoninfra.RegisterEvent[eventKts.KtsCreatedEvent](dispatcher)
	commoninfra.RegisterEvent[eventKts.KtsUpdatedEvent](dispatcher)

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
	accountPresentation.ModuleAccount(app)
	renstranilaiPresentation.ModuleRenstraNilai(app)
	dokumentambahanPresentation.ModuleDokumenTambahan(app)
	ktsPresentation.ModuleKts(app)

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	outboxProcessor := &commoninfra.OutboxProcessor{
		DB:         db,
		Dispatcher: dispatcher,
	}

	app.Get("/swagger/*", swagger.HandlerDefault)
	go commoninfra.StartOutboxWorker(ctx, outboxProcessor)
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

	// === Generate Login Commands ===
	case login.LoginCommand:
		if err := login.LoginCommandValidation(cmd); err != nil {
			return nil, wrapValidationError("Login.Validation", err)
		}

	// === Generate Renstra Nilai Commands ===
	case updateRenstraNilai.UpdateRenstraNilaiCommand:
		if err := updateRenstraNilai.UpdateRenstraNilaiCommandValidation(cmd); err != nil {
			return nil, wrapValidationError("RenstraNilaiUpdate.Validation", err)
		}
	case deleteRenstraNilai.DeleteRenstraNilaiCommand:
		if err := deleteRenstraNilai.DeleteRenstraNilaiCommandValidation(cmd); err != nil {
			return nil, wrapValidationError("RenstraNilaiDelete.Validation", err)
		}

	// === Generate Dokumen Tambahan Commands ===
	case updateDokumenTambahan.UpdateDokumenTambahanCommand:
		if err := updateDokumenTambahan.UpdateDokumenTambahanCommandValidation(cmd); err != nil {
			return nil, wrapValidationError("DokumenTambahanUpdate.Validation", err)
		}
	case deleteDokumenTambahan.DeleteDokumenTambahanCommand:
		if err := deleteDokumenTambahan.DeleteDokumenTambahanCommandValidation(cmd); err != nil {
			return nil, wrapValidationError("DokumenTambahanDelete.Validation", err)
		}

	// === Generate KTS Commands ===
	case updateKts.UpdateKtsCommand:
		if err := updateKts.UpdateKtsCommandValidation(cmd); err != nil {
			return nil, wrapValidationError("KtsUpdate.Validation", err)
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
		sqlDB.SetConnMaxLifetime(10 * time.Minute)
		sqlDB.SetConnMaxIdleTime(2 * time.Minute)
	})

	return db, err
}
