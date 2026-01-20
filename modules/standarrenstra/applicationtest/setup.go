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

func setupStandarRenstraMySQL(t *testing.T) (*gorm.DB, func()) {
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
		CREATE TABLE master_standar_renstra (
			id int(11) PRIMARY KEY AUTO_INCREMENT NOT NULL,
			uuid varchar(36) DEFAULT NULL,
			nama text NOT NULL,
			created_at datetime DEFAULT NULL,
			updated_at datetime DEFAULT NULL
		);

        INSERT INTO master_standar_renstra (id, uuid, nama, created_at, updated_at) VALUES
		(1, '5fd713d0-adfe-4086-a000-21c948faf84d', 'Standar Kompetensi Lulusan', '2024-08-09 09:35:07', '2024-08-25 19:24:41'),
		(4, '2025f1ad-7afa-4609-9fd5-24e369f2f463', 'Standar Isi Pembelajaran', '2024-08-25 19:22:07', '2024-08-25 19:22:07'),
		(7, '2486a3ca-0898-4e5b-8bb0-ca5312ee2347', 'Standar Proses Pembelajaran', '2024-08-25 19:22:14', '2024-08-25 19:22:14'),
		(10, '9d16efe3-8d04-4c7b-9433-5f9041a94c58', 'Standar Penilaian Pembelajaran', '2024-08-25 19:22:20', '2024-08-25 19:22:20'),
		(13, 'fab663c7-45c1-46f6-8d39-fe47ecae6045', 'Standar Dosen dan Tenaga Kependidikan', '2024-08-25 19:22:33', '2024-08-25 19:22:33'),
		(25, 'c88af609-88e7-4b4b-92cf-8cd4f3c78c50', 'Standar Sarana dan Prasarana Pembelajaran', '2024-09-09 14:25:14', '2024-09-09 14:25:14'),
		(28, '8d7fa053-4c42-4b42-9b4a-7efd5500416f', 'Standar Pengelolaan Pembelajaran', '2024-09-09 14:25:28', '2024-09-09 14:25:28'),
		(118, 'f15a113a-9bcb-4cd4-9463-0aa704137004', 'Standar Pembiayaan Pembelajaran', '2024-09-10 09:12:37', '2024-09-10 09:12:37'),
		(121, '130aa46d-4a2a-419b-bdc9-05160165716c', 'Standar Hasil Penelitian', '2024-09-10 09:12:42', '2024-09-10 09:12:42'),
		(124, '34c5cea9-3ba9-4090-a969-6524d74d869f', 'Standar Isi Penelitian', '2024-09-10 09:12:49', '2024-09-10 09:12:49'),
		(127, 'e8691c3d-9e73-4919-a9e5-6193a1f51c20', 'Standar Proses Penelitian', '2024-09-10 09:12:54', '2024-09-10 09:12:54'),
		(130, 'b7b2bce1-2afd-40f6-bde9-d9be18f49e52', 'Standar Penilaian Penelitian', '2024-09-10 09:13:10', '2024-09-10 09:13:10'),
		(133, 'c0ef6d0d-57eb-4c62-b69c-ad845e4da773', 'Standar Peneliti', '2024-09-10 09:13:24', '2024-09-10 09:13:24'),
		(136, '01c577fd-36d1-4bf0-bfbc-eaacaadfa0f0', 'Standar Sarana dan Prasarana Penelitian', '2024-09-10 09:13:33', '2024-09-10 09:13:33'),
		(139, '2b5f2d39-6855-4993-9f10-af9fefddfea8', 'Standar Pengelolaan Penelitian', '2024-09-10 09:13:41', '2024-09-10 09:13:41'),
		(142, '3422abc5-5252-4e75-b6e4-2a2cc7e659dc', 'Standar Pendanaan dan Pembiayaan Penelitian', '2024-09-10 09:13:58', '2024-09-10 09:13:58'),
		(145, '18ebbe28-9c2f-44d0-8eb3-708ad347bd9d', 'Standar Hasil Pengabdian Kepada Masyarakat', '2024-09-10 09:23:17', '2024-09-10 09:23:17'),
		(148, 'f1c3a7cf-0504-4349-a116-6c25b502a0c1', 'Standar Isi Pengabdian Kepada Masyarakat', '2024-09-10 09:23:27', '2024-09-10 09:23:27'),
		(151, '8c15b52d-a5f6-4ca9-8cbd-cd0a4f72af58', 'Standar Proses Pengabdian Kepada Masyarakat', '2024-09-10 09:23:38', '2024-09-10 09:23:38'),
		(154, 'b8194497-a685-4d63-8316-8836e89676ea', 'Standar Penilaian Pengabdian Kepada Masyarakat', '2024-09-10 09:23:46', '2024-09-10 09:23:46'),
		(157, '88bff504-ec0f-42de-b81a-2cb42e97e849', 'Standar Pelaksana Pengabdian Kepada Masyarakat', '2024-09-10 09:24:01', '2024-09-10 09:24:01'),
		(160, '99f520c0-fb9c-499e-b79b-c4a8beefa80f', 'Standar Sarana dan Prasarana Pengabdian Kepada Masyarakat', '2024-09-10 09:25:16', '2024-09-10 09:25:16'),
		(163, 'f5604771-e003-4ca5-b97b-2386dcce8504', 'Standar Pengelolaan Pengabdian Kepada Masyarakat', '2024-09-10 09:25:35', '2024-09-10 09:25:35'),
		(166, '8d302d77-9a86-4260-b693-1652138586d8', 'Standar Pendanaan dan Pembiayaan PkM', '2024-09-10 09:25:55', '2024-09-10 09:25:55'),
		(169, '2f66cc47-28f7-4a6a-add0-dfe18140c458', 'Standar Visi Misi (kriteria 1)', '2024-09-10 09:26:02', '2024-09-10 09:26:02'),
		(172, '57ef4038-0add-491b-be45-d7112f38017e', 'Standar Ketaatan Pada Peraturan Perundang-Undangan (kriteria 2)', '2024-09-10 09:26:08', '2024-09-10 09:26:08'),
		(175, '2e2f66f6-9568-4930-a4e3-ec79acf5721c', 'Standar Pengelolaan Tata Pamong (kriteria 2)', '2024-09-10 09:26:12', '2024-09-10 09:26:12'),
		(178, '2ce1e09a-0dce-4fc1-b1aa-212dff1184d9', 'Standar Penjanjian Kerjasama (kriteria 2)', '2024-09-10 09:26:17', '2024-09-10 09:26:17'),
		(181, '075dc3cf-b05f-45e5-9e92-45f1cccd8d6d', 'Standar Pemeliharaan/Peningkatan Jumlah Peminat/Pendaftar (kriteria 3)', '2024-09-10 09:26:29', '2024-09-10 09:26:29'),
		(184, '7a64ef57-df41-4152-8100-06590e7a91e4', 'Standar Layanan Kemahasiswaan (kriteria 3)', '2024-09-10 09:26:40', '2024-09-10 09:26:40'),
		(187, '952af701-4589-4923-9ff2-e1babe4ccdd2', 'Standar Prestasi Mahasiswa (kriteria 3 dan 9)', '2024-09-10 09:26:44', '2024-09-10 09:26:44'),
		(190, 'ef6df7fb-fb39-4fce-b757-d6e823461981', 'Standar Pengelolaan Keuangan (kriteria 5)', '2024-09-10 09:27:13', '2024-09-10 09:27:13'),
		(193, '7b371ba9-053e-4477-80b4-b16efcc6836b', 'Standar Sarana Prasarana Umum (kriteria 5)', '2024-09-10 09:27:18', '2024-09-10 09:27:18'),
		(196, 'e3a92a8c-17e0-46b5-93a0-99900b95d848', 'Standar Sistem Informasi (kriteria 5)', '2024-09-10 09:27:24', '2024-09-10 09:27:24'),
		(199, '2561e7d1-dad5-4164-ba58-4a8c17493a2e', 'Standar Pembiayaan MBKM (kriteria 5 dan 6)', '2024-09-10 09:27:32', '2024-09-10 09:27:32'),
		(202, '990b94a6-61d5-40ad-818e-191f272cf1fe', 'Standar Pelaksanaan MBKM (kriteria 6)', '2024-09-10 09:27:37', '2024-09-10 09:27:37'),
		(205, '1fc4e63d-b5ee-4acf-bd19-8a477886c9f5', 'Standar Pemeliharaan/Peningkatan Jumlah Lulusan (kriteria 6 dan 9)', '2024-09-10 09:27:48', '2024-09-10 09:27:48'),
		(208, '0e3231d0-7e02-441c-923b-ca97bfbfe300', 'Standar Tracer Study (kriteria 9)', '2024-09-10 09:27:55', '2024-09-10 09:27:55'),
		(211, 'a5372639-4deb-4e62-8b8f-b6b183a58662', 'Standar Inovasi dan Inkubator Bisnis (kriteria 9)', '2024-09-10 09:28:01', '2024-09-10 09:28:01');

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
