package applicationtest

import (
	"context"
	"fmt"
	"testing"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupTahunProkerMySQL(t *testing.T) (*gorm.DB, func()) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image: "mysql:8.0",
		Env: map[string]string{
			"MYSQL_ROOT_PASSWORD": "pass",
			"MYSQL_DATABASE":      "testdb",
		},
		ExposedPorts: []string{"3306/tcp"},
		WaitingFor: wait.ForListeningPort("3306/tcp").
			WithStartupTimeout(10 * time.Minute),
	}

	mysqlC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("failed to start container: %v", err)
	}

	host, _ := mysqlC.Host(ctx)
	port, _ := mysqlC.MappedPort(ctx, "3306")

	dsn := fmt.Sprintf("root:pass@tcp(%s:%s)/testdb?parseTime=true&multiStatements=true&allowNativePasswords=true", host, port.Port())

	var gdb *gorm.DB
	// Retry connect 10x
	for i := 0; i < 10; i++ {
		gdb, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		t.Logf("retrying GORM connection: %v", err)
		time.Sleep(1 * time.Second)
	}
	if err != nil {
		t.Fatalf("cannot connect via GORM after retries: %v", err)
	}

	// Buat table & data contoh
	err = gdb.Exec(`
       CREATE TABLE master_tahun (
			id int(11) PRIMARY KEY AUTO_INCREMENT NOT NULL,
			uuid varchar(36) DEFAULT NULL,
			tahun varchar(100) DEFAULT NULL,
			status varchar(100) DEFAULT NULL,
			created_at timestamp NULL DEFAULT NULL,
			updated_at timestamp NULL DEFAULT NULL,
			UNIQUE KEY uk_master_tahun_tahun (tahun)
		);

        INSERT INTO master_tahun (id, uuid, tahun, status, created_at, updated_at) VALUES
		(7, '666a6b72-d2b4-481f-adb8-298d807e9e20', '2023', 'non-aktif', '2025-07-24 23:49:36', '2025-09-01 06:45:30'),
		(10, 'ea08bcc7-2333-4afe-8c7b-a1cd14d6114c', '2024', 'aktif', '2025-07-24 23:49:36', '2025-11-25 03:07:47'),
		(13, '1630ec58-6183-44a3-85c1-cf7da95dbd9e', '2025', 'aktif', '2025-07-24 23:49:36', '2025-08-22 03:34:11');
    `).Error

	if err != nil {
		t.Fatalf("migration failed: %v", err)
	}

	cleanup := func() {
		sqlDB, _ := gdb.DB()
		sqlDB.Close()
		mysqlC.Terminate(ctx)
	}

	return gdb, cleanup
}
