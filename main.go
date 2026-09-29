package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/redis/go-redis/v9"

	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	commondomain "UnpakSiamida/common/domain"

	tahunprokerInfrastructure "UnpakSiamida/modules/tahunproker/infrastructure"

	tahunprokerPresentation "UnpakSiamida/modules/tahunproker/presentation"

	mataprogramInfrastructure "UnpakSiamida/modules/mataprogram/infrastructure"

	mataprogramPresentation "UnpakSiamida/modules/mataprogram/presentation"

	jadwalprokerInfrastructure "UnpakSiamida/modules/jadwalproker/infrastructure"

	jadwalprokerPresentation "UnpakSiamida/modules/jadwalproker/presentation"

	aktivitasprokerInfrastructure "UnpakSiamida/modules/aktivitasproker/infrastructure"

	aktivitasprokerPresentation "UnpakSiamida/modules/aktivitasproker/presentation"

	laporanInfrastructure "UnpakSiamida/modules/laporan/infrastructure"

	laporanPresentation "UnpakSiamida/modules/laporan/presentation"

	beritaacaraInfrastructure "UnpakSiamida/modules/beritaacara/infrastructure"

	beritaacaraPresentation "UnpakSiamida/modules/beritaacara/presentation"

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

	/////////

	commoninfra "UnpakSiamida/common/infrastructure"

	commonpresentation "UnpakSiamida/common/presentation"

	//////////

	eventBeritaAcara "UnpakSiamida/modules/beritaacara/event"
	eventKts "UnpakSiamida/modules/kts/event"
	eventUser "UnpakSiamida/modules/user/event"

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
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  10 * time.Second,
	})
	// app.Use(recover())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders: "*",
	}))
	app.Use(helmet.New(helmet.Config{
		XSSProtection:             "1; mode=block",
		ContentTypeNosniff:        "nosniff",     // X-Content-Type-Options
		XFrameOptions:             "DENY",        // X-Frame-Options
		ReferrerPolicy:            "no-referrer", // Referrer-Policy
		ContentSecurityPolicy:     "default-src 'self'; script-src 'self'; object-src 'none'; base-uri 'none'",
		CrossOriginEmbedderPolicy: "require-corp",
		CrossOriginOpenerPolicy:   "same-origin",
		CrossOriginResourcePolicy: "same-origin",
	}))
	app.Use(commonpresentation.LoggerMiddleware)
	app.Use(commonpresentation.HeaderSecurityMiddleware(cfg))
	app.Use(func(c *fiber.Ctx) error {
		c.Response().Header.Del("X-Powered-By")
		return c.Next()
	})

	mediatr.RegisterRequestPipelineBehaviors(NewValidationBehavior())

	var db *gorm.DB
	var redis commondomain.IRedisStore
	mustStart("Database", func() error {
		var err error
		db, err = NewMySQL()
		return err
	})
	mustStart("Redis", func() error {
		var err error
		redis = NewRedisStore()
		return err
	})

	var tg commoninfra.TelegramSender
	modeTelegram := os.Getenv("TELEGRAM_MODE")

	mustStart("Telegram Service", func() error {
		factory := &commoninfra.DefaultTelegramFactory{
			UseFake: modeTelegram != "dev",
		}

		client, err := factory.Create()
		if err != nil {
			return err
		}

		tg = client
		return nil
	})

	//berlaku untuk startup bukan hot reload
	mustStart("User Module", func() error {
		return userInfrastructure.RegisterModuleUser(db, tg)
	})

	mustStart("Berita Acara Module", func() error {
		return beritaacaraInfrastructure.RegisterModuleBeritaAcara(db, &redis)
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
		return ktsInfrastructure.RegisterModuleKts(db, &redis, tg)
	})

	mustStart("Tahun Proker Module", func() error { //buat audit
		return tahunprokerInfrastructure.RegisterModuleTahunProker(db)
	})

	mustStart("Mata Program Module", func() error { //buat audit
		return mataprogramInfrastructure.RegisterModuleMataProgram(db)
	})

	mustStart("Jadwal Proker Module", func() error { //buat audit
		return jadwalprokerInfrastructure.RegisterModuleJadwalProker(db)
	})

	mustStart("Aktivitas Proker Module", func() error { //buat audit
		return aktivitasprokerInfrastructure.RegisterModuleAktivitasProker(db)
	})

	mustStart("Laporan Module", func() error { //buat audit
		return laporanInfrastructure.RegisterModuleLaporan(db)
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
	commoninfra.RegisterEvent[eventUser.UserCreatedEvent](dispatcher)
	commoninfra.RegisterEvent[eventUser.UserUpdatedEvent](dispatcher)
	commoninfra.RegisterEvent[eventBeritaAcara.BeritaAcaraPdfRequestedEvent](dispatcher)
	commoninfra.RegisterEvent[eventKts.KtsPdfRequestedEvent](dispatcher)

	beritaacaraPresentation.ModuleBeritaAcara(app)
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
	tahunprokerPresentation.ModuleTahunProker(app)
	mataprogramPresentation.ModuleMataProgram(app)
	jadwalprokerPresentation.ModuleJadwalProker(app)
	aktivitasprokerPresentation.ModuleAktivitasProker(app)
	laporanPresentation.ModuleMonitoringProker(app)

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

	if err := commoninfra.Validate(request); err != nil {
		return nil, err
	}

	return next(ctx)
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

func NewRedisStore() *commoninfra.RedisStore {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})

	return commoninfra.NewRedisStore(client)
}
