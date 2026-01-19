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

func setupTemplateDokumenTambahanMySQL(t *testing.T) (*gorm.DB, func()) {
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
       	DROP TABLE IF EXISTS jenis_file_renstra;
		CREATE TABLE jenis_file_renstra (
			id bigint(20) UNSIGNED PRIMARY KEY AUTO_INCREMENT NOT NULL,
			uuid varchar(36) DEFAULT NULL,
			nama text NOT NULL,
			created_at datetime DEFAULT NULL,
			updated_at datetime DEFAULT NULL
		);

		INSERT INTO jenis_file_renstra (id, uuid, nama, created_at, updated_at) VALUES
		(1, '14212231-792f-4935-bb1c-9a38695a4b6b', 'Program Kerja Sesuai Dengan Template 2024 disertai Monev', NULL, '2024-10-08 13:35:36'),
		(2, '1a353e22-1111-4fc5-96c1-a2ed2877a6a4', 'Struktur Organisasi Disertai Dengan Job Description, Job Spesification, Dan Disahkan Dengan SK', NULL, '2024-10-08 13:36:55'),
		(3, '08a5e4cc-1a30-4080-95ad-127abf8819f5', 'SOP Pelaksanaan Tugas Sesuai Dengan Job Description Dan Job Spesification', NULL, '2024-10-08 13:37:14'),
		(4, '3523b28b-829a-4548-96e0-459bf4f14dea', 'Hasil/Catatan Audit/KTS Sebelumnya (2023) Telah Diselesaikan (Closed)', NULL, '2024-10-08 13:37:47'),
		(14, '6aee7cd5-da31-4735-9243-8c19aa7497c0', 'Program Kerja Sesuai Dengan Template 2025 disertai Monev', '2025-10-13 15:14:06', '2025-10-13 15:14:06'),
		(15, '84066942-1f2d-44b0-be66-8b87cdab6e91', 'Hasil/Catatan Audit/KTS Sebelumnya (2024) Telah Diselesaikan (Closed)', '2025-10-13 15:14:25', '2025-10-13 15:14:25');

	   	DROP TABLE IF EXISTS template_dokumen_tambahan;
		CREATE TABLE template_dokumen_tambahan (
			id int(11) PRIMARY KEY AUTO_INCREMENT NOT NULL,
			uuid varchar(36) DEFAULT NULL,
			tahun varchar(100) NOT NULL,
			jenis_file bigint(11) UNSIGNED NOT NULL,
			pertanyaan text NOT NULL,
			klasifikasi varchar(100) NOT NULL,
			fakultas_prodi_unit varchar(100) NOT NULL,
			tugas varchar(100) NOT NULL DEFAULT 'auditor2',
			created_at datetime DEFAULT NULL,
			updated_at datetime DEFAULT NULL
		);
		ALTER TABLE template_dokumen_tambahan
		ADD UNIQUE KEY uq_template_dokumen (tahun,jenis_file,fakultas_prodi_unit);

		INSERT INTO template_dokumen_tambahan (id, uuid, tahun, jenis_file, pertanyaan, klasifikasi, fakultas_prodi_unit, tugas, created_at, updated_at) VALUES
		(28, '9b354f31-be71-4173-9e26-c319d163660d', '2024', 1, 'Apakah Sudah Lengkap Sesuai Dengan Template Proker 2024 Beserta Monevnya?', 'minor', 'fakultas#all', 'auditor2', NULL, '2025-02-04 21:59:46'),
		(30, '4bda1f5a-c31f-4db0-a70a-ce007f7fb6de', '2024', 1, 'Apakah Sudah Lengkap Sesuai Dengan Template Proker 2024 Beserta Monevnya?', 'minor', 'prodi#all', 'auditor2', NULL, '2025-02-04 21:59:46'),
		(32, '0552e387-c22c-4074-9edf-6bcac2f11851', '2024', 1, 'Apakah Sudah Lengkap Sesuai Dengan Template Proker 2024 Beserta Monevnya?', 'minor', 'unit#all', 'auditor2', NULL, '2025-02-04 21:59:46'),
		(34, '9cc38a80-d11b-44c9-9b81-f394816eaa96', '2024', 2, 'Apakah Struktur Organisasi Sudah Lengkap Disertai Dengan Job Description, Job Spesification, Dan Disahkan Dengan SK?', 'minor', 'fakultas#all', 'auditor2', NULL, '2024-10-08 13:36:04'),
		(36, 'e72af7fd-b9a8-4e4b-a8b7-c725ed3ae597', '2024', 2, 'Apakah Struktur Organisasi Sudah Lengkap Disertai Dengan Job Description, Job Spesification, Dan Disahkan Dengan SK?', 'minor', 'prodi#all', 'auditor2', NULL, '2024-10-08 13:36:04'),
		(38, 'a5b748e5-0f0d-407b-9858-06d7d4b3f914', '2024', 2, 'Apakah Struktur Organisasi Sudah Lengkap Disertai Dengan Job Description, Job Spesification, Dan Disahkan Dengan SK?', 'minor', 'unit#all', 'auditor2', NULL, '2024-10-08 13:36:04'),
		(40, '19542f2a-2c8c-4a63-8d23-1a61b9d583b3', '2024', 3, 'Apakah Memiliki SOP Pelaksanaan Tugas Sesuai Dengan Job Description Dan Job Spesification?', 'minor', 'fakultas#all', 'auditor2', NULL, NULL),
		(42, '320b5233-e228-448a-80dc-d28c33faad5f', '2024', 3, 'Apakah Memiliki SOP Pelaksanaan Tugas Sesuai Dengan Job Description Dan Job Spesification?', 'minor', 'prodi#all', 'auditor2', NULL, NULL),
		(44, '431fc697-6869-47e7-8011-184bec8eaf49', '2024', 3, 'Apakah Memiliki SOP Pelaksanaan Tugas Sesuai Dengan Job Description Dan Job Spesification?', 'minor', 'unit#all', 'auditor2', NULL, NULL),
		(46, '01dcdfa5-ce8d-4cd2-8921-c2e40f595052', '2024', 4, 'Apakah KTS Pada Audit Sebelumnya (2023) Telah Diselesaikan (Closed)?', 'major', 'fakultas#all', 'auditor2', NULL, '2025-10-13 15:21:14'),
		(48, '23b60f74-b023-46b5-9e05-8fa1715e4702', '2024', 4, 'Apakah KTS Pada Audit Sebelumnya (2023) Telah Diselesaikan (Closed)?', 'major', 'prodi#all', 'auditor2', NULL, '2025-10-13 15:21:14'),
		(50, '8cde2ebe-6abd-46b7-a942-5941c4154f46', '2024', 4, 'Apakah KTS Pada Audit Sebelumnya (2023) Telah Diselesaikan (Closed)?', 'major', 'unit#all', 'auditor2', NULL, '2025-10-13 15:21:14'),
		(56, '9080464e-3d3f-4040-b9c5-cfef8d3a0bba', '2025', 14, 'Apakah Sudah Lengkap Sesuai Dengan Template Proker 2025 Beserta Monevnya?', 'minor', 'fakultas#all', 'auditor2', NULL, NULL),
		(57, 'b92f38bb-f328-44e3-b106-bb3f23cf4e5f', '2025', 14, 'Apakah Sudah Lengkap Sesuai Dengan Template Proker 2025 Beserta Monevnya?', 'minor', 'prodi#all', 'auditor2', NULL, NULL),
		(58, '88cc484d-dc60-4ed4-bcf0-5d2c2d05dafb', '2025', 14, 'Apakah Sudah Lengkap Sesuai Dengan Template Proker 2025 Beserta Monevnya?', 'minor', 'unit#all', 'auditor2', NULL, NULL),
		(59, '3f0ee834-3a7d-40c9-b04c-803e4df48ba5', '2025', 2, 'Apakah Struktur Organisasi Sudah Lengkap Disertai Dengan Job Description, Job Spesification, Dan Disahkan Dengan SK?', 'minor', 'fakultas#all', 'auditor2', NULL, NULL),
		(60, '817e397a-755f-4d8f-81c6-be4a9c2f27c9', '2025', 2, 'Apakah Struktur Organisasi Sudah Lengkap Disertai Dengan Job Description, Job Spesification, Dan Disahkan Dengan SK?', 'minor', 'prodi#all', 'auditor2', NULL, NULL),
		(61, '4ace1d89-8e03-40a2-be31-c73f8c746476', '2025', 2, 'Apakah Struktur Organisasi Sudah Lengkap Disertai Dengan Job Description, Job Spesification, Dan Disahkan Dengan SK?', 'minor', 'unit#all', 'auditor2', NULL, NULL),
		(62, '2d026ffb-359d-4203-800c-0b3211770457', '2025', 3, 'Apakah Memiliki SOP Pelaksanaan Tugas Sesuai Dengan Job Description Dan Job Spesification?', 'minor', 'fakultas#all', 'auditor2', NULL, NULL),
		(63, 'f5880e2a-9859-4bd7-94e9-0d1d57fde80f', '2025', 3, 'Apakah Memiliki SOP Pelaksanaan Tugas Sesuai Dengan Job Description Dan Job Spesification?', 'minor', 'prodi#all', 'auditor2', NULL, NULL),
		(64, 'aed66f0f-c584-4857-a1a1-08324a3c964a', '2025', 3, 'Apakah Memiliki SOP Pelaksanaan Tugas Sesuai Dengan Job Description Dan Job Spesification?', 'minor', 'unit#all', 'auditor2', NULL, NULL),
		(65, '5f904d3f-bd98-48ab-80b5-89a6a5c71a71', '2025', 15, 'Apakah KTS Pada Audit Sebelumnya (2024) Telah Diselesaikan (Closed)?', 'major', 'fakultas#all', 'auditor2', NULL, NULL),
		(66, '0265bb82-0707-49f0-82e7-ad452072f4c1', '2025', 15, 'Apakah KTS Pada Audit Sebelumnya (2024) Telah Diselesaikan (Closed)?', 'major', 'prodi#all', 'auditor2', NULL, NULL),
		(67, 'a94a807a-cb72-40b1-b9a4-d62bfa1a1206', '2025', 15, 'Apakah KTS Pada Audit Sebelumnya (2024) Telah Diselesaikan (Closed)?', 'major', 'unit#all', 'auditor2', NULL, NULL),
		(71, 'f3740e32-c373-4833-90bc-69b5cc358a93', '2026', 15, 'beta tester', 'minor', 'unit#all', 'auditor2', NULL, NULL);
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
