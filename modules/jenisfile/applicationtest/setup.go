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

func setupJenisFileMySQL(t *testing.T) (*gorm.DB, func()) {
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
        CREATE TABLE jenis_file_renstra (
            id bigint(20) AUTO_INCREMENT UNSIGNED NOT NULL,
            uuid varchar(36) DEFAULT NULL,
            nama text NOT NULL,
            created_at datetime DEFAULT NULL,
            updated_at datetime DEFAULT NULL
        );

        ALTER TABLE jenis_file_renstra
        ADD PRIMARY KEY (id);

        INSERT INTO jenis_file_renstra (id, uuid, nama, created_at, updated_at) VALUES
        (1, '14212231-792f-4935-bb1c-9a38695a4b6b', 'Program Kerja Sesuai Dengan Template 2024 disertai Monev', NULL, '2024-10-08 13:35:36'),
        (2, '1a353e22-1111-4fc5-96c1-a2ed2877a6a4', 'Struktur Organisasi Disertai Dengan Job Description, Job Spesification, Dan Disahkan Dengan SK', NULL, '2024-10-08 13:36:55'),
        (3, '08a5e4cc-1a30-4080-95ad-127abf8819f5', 'SOP Pelaksanaan Tugas Sesuai Dengan Job Description Dan Job Spesification', NULL, '2024-10-08 13:37:14'),
        (4, '3523b28b-829a-4548-96e0-459bf4f14dea', 'Hasil/Catatan Audit/KTS Sebelumnya (2023) Telah Diselesaikan (Closed)', NULL, '2024-10-08 13:37:47'),
        (14, '6aee7cd5-da31-4735-9243-8c19aa7497c0', 'Program Kerja Sesuai Dengan Template 2025 disertai Monev', '2025-10-13 15:14:06', '2025-10-13 15:14:06'),
        (15, '84066942-1f2d-44b0-be66-8b87cdab6e91', 'Hasil/Catatan Audit/KTS Sebelumnya (2024) Telah Diselesaikan (Closed)', '2025-10-13 15:14:25', '2025-10-13 15:14:25');
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
