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

func setupRenstraMySQL(t *testing.T) (*gorm.DB, func()) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image: "mysql:8.0",
		Tmpfs: map[string]string{
			"/var/lib/mysql": "rw",
		},
		Cmd: []string{
			"--innodb_flush_log_at_trx_commit=2",
			"--sync_binlog=0",
		},
		Env: map[string]string{
			"MYSQL_ROOT_PASSWORD": "pass",
			"MYSQL_DATABASE":      "testdb",
		},
		Labels: map[string]string{
			"testcontainers.sessionId": "tahunrenstra",
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
        DROP TABLE IF EXISTS users;
        CREATE TABLE users (
            id bigint(20) UNSIGNED NOT NULL,
            uuid varchar(36) DEFAULT NULL,
            nidn_username varchar(255) NOT NULL,
            password varchar(255) NOT NULL,
            level enum('','admin','user','auditee','auditor1','auditor2','fakultas') NOT NULL,
            name varchar(255) NOT NULL,
            email varchar(255) DEFAULT NULL,
            fakultas_unit int(11) DEFAULT NULL,
            foto text DEFAULT NULL,
            email_verified_at timestamp NULL DEFAULT NULL,
            remember_token varchar(100) DEFAULT NULL,
            reset_password int(11) DEFAULT 0,
            created_at timestamp NULL DEFAULT NULL,
            updated_at timestamp NULL DEFAULT NULL
        );

        DROP TABLE IF EXISTS m_fakultas;
        CREATE TABLE m_fakultas (
            kode_fakultas char(9) NOT NULL,
            kode_pt char(10) NOT NULL,
            nama_fakultas varchar(100) NOT NULL DEFAULT '',
            pejabat varchar(50) DEFAULT NULL,
            jabatan char(1) DEFAULT NULL,
            wakil_pejabat varchar(50) DEFAULT NULL,
            wakil_pejabat_adm varchar(50) DEFAULT NULL,
            logo varchar(50) DEFAULT NULL
        );

        DROP TABLE IF EXISTS m_program_studi;
        CREATE TABLE m_program_studi (
            kode_prodi char(10) NOT NULL,
            kode_pt char(10) NOT NULL,
            kode_fak char(9) DEFAULT NULL,
            kode_jenjang varchar(1) DEFAULT NULL,
            kode_jurusan char(5) NOT NULL,
            nama_prodi varchar(50) DEFAULT NULL,
            alamat varchar(100) DEFAULT NULL,
            kode_kabupaten int(10) DEFAULT NULL,
            kode_propinsi int(10) DEFAULT NULL,
            kode_negara int(10) DEFAULT NULL,
            kode_pos varchar(10) DEFAULT NULL,
            telepon varchar(20) DEFAULT NULL,
            fax varchar(20) DEFAULT NULL,
            email varchar(50) DEFAULT NULL,
            website varchar(50) DEFAULT NULL,
            sks_lulus int(11) DEFAULT NULL,
            status_prodi varchar(1) DEFAULT NULL,
            tgl_awal_berdiri date DEFAULT NULL,
            semester_awal varchar(5) DEFAULT NULL,
            mulai_semester varchar(5) DEFAULT NULL,
            no_sk_dikti varchar(40) DEFAULT NULL,
            tgl_sk_dikti date DEFAULT NULL,
            tgl_akhir_sk_dikti date DEFAULT NULL,
            no_sk_ban varchar(40) DEFAULT NULL,
            tgl_sk_ban date DEFAULT NULL,
            tgl_akhir_sk_ban date DEFAULT NULL,
            kode_akreditasi varchar(1) DEFAULT NULL,
            frekuensi_kurikulum varchar(1) DEFAULT NULL,
            pelaksanaan_kurikulum varchar(1) DEFAULT NULL,
            idd_ketua_prodi varchar(50) DEFAULT NULL,
            hp_ketua varchar(20) DEFAULT NULL,
            idd_nama_operator varchar(50) DEFAULT NULL,
            telepon_operator varchar(20) DEFAULT NULL,
            nama_sesi varchar(20) NOT NULL,
            jumlah_sesi int(11) NOT NULL,
            batas_sesi int(11) NOT NULL,
            gelar varchar(20) NOT NULL,
            gelar_panjang varchar(200) NOT NULL,
            no_sk_ban_lama varchar(40) NOT NULL,
            logo varchar(50) DEFAULT NULL,
            nama_prodi_ing varchar(50) DEFAULT NULL
        );

        DROP TABLE IF EXISTS sijamu_fakultas_unit;
        CREATE TABLE sijamu_fakultas_unit (
            id int(11) NOT NULL,
            uuid varchar(36) DEFAULT NULL,
            kode_fakultas char(9) DEFAULT NULL,
            kode_prodi char(10) DEFAULT NULL,
            nama varchar(100) DEFAULT NULL,
            id_m_prodi int(11) DEFAULT NULL,
            standalone tinyint(4) DEFAULT 0
        );

        DROP TABLE IF EXISTS renstra;
        CREATE TABLE renstra (
            id int(11) PRIMARY KEY AUTO_INCREMENT NOT NULL,
            uuid varchar(36) DEFAULT NULL,
            tahun year(4) NOT NULL,
            fakultas_unit_old int(11) DEFAULT NULL,
            fakultas_unit int(11) DEFAULT NULL,
            periode_upload_mulai datetime NOT NULL,
            periode_upload_akhir datetime NOT NULL,
            periode_assesment_dokumen_mulai datetime NOT NULL,
            periode_assesment_dokumen_akhir datetime NOT NULL,
            periode_assesment_lapangan_mulai datetime NOT NULL,
            periode_assesment_lapangan_akhir datetime NOT NULL,
            kodeAkses varchar(255) DEFAULT NULL,
            auditee bigint(20) UNSIGNED DEFAULT NULL,
            auditor1 bigint(20) UNSIGNED DEFAULT NULL,
            auditor2 bigint(20) UNSIGNED DEFAULT NULL,
            catatan text DEFAULT NULL,
            catatan2 text DEFAULT NULL,
            created_at datetime DEFAULT NULL,
            updated_at datetime DEFAULT NULL
        );

        CREATE OR REPLACE VIEW v_fakultas_unit AS
        SELECT
            sfu.id AS id,
            sfu.uuid AS uuid,
            CASE
                WHEN sfu.standalone THEN fak.nama_fakultas
                WHEN sfu.kode_fakultas IS NOT NULL
                    AND sfu.kode_prodi IS NOT NULL
                    AND sfu.nama IS NULL THEN prod.nama_prodi
                WHEN sfu.kode_fakultas IS NOT NULL
                    AND sfu.kode_prodi IS NOT NULL
                    AND sfu.nama IS NOT NULL THEN sfu.nama
                WHEN sfu.kode_fakultas IS NOT NULL
                    AND sfu.kode_prodi IS NULL
                    AND sfu.nama IS NOT NULL THEN sfu.nama
                ELSE 'tidak diketahui'
            END AS nama_fak_prod_unit,
            prod.kode_jenjang AS kode_jenjang,
            CASE
                WHEN prod.kode_jenjang = 'C' THEN 's1'
                WHEN prod.kode_jenjang = 'B' THEN 's2'
                WHEN prod.kode_jenjang = 'A' THEN 's3'
                WHEN prod.kode_jenjang = 'E' THEN 'd3'
                WHEN prod.kode_jenjang = 'D' THEN 'd4'
                WHEN prod.kode_jenjang = 'J' THEN 'profesi'
                ELSE NULL
            END AS jenjang,
            CASE
                WHEN prod.kode_jenjang = 'C' THEN '1'
                WHEN prod.kode_jenjang = 'B' THEN '2'
                WHEN prod.kode_jenjang = 'A' THEN '3'
                WHEN prod.kode_jenjang = 'E' THEN '4'
                WHEN prod.kode_jenjang = 'D' THEN '5'
                WHEN prod.kode_jenjang = 'J' THEN '6'
                ELSE '7'
            END AS jenjang_int,
            CASE
                WHEN sfu.standalone THEN 'fakultas'
                WHEN sfu.kode_fakultas IS NOT NULL
                    AND sfu.kode_prodi IS NOT NULL
                    AND sfu.nama IS NULL THEN 'prodi'
                WHEN sfu.kode_fakultas IS NOT NULL
                    AND sfu.kode_prodi IS NOT NULL
                    AND sfu.nama IS NOT NULL THEN 'prodi'
                WHEN sfu.kode_fakultas IS NOT NULL
                    AND sfu.kode_prodi IS NULL
                    AND sfu.nama IS NOT NULL THEN 'unit'
                ELSE NULL
            END AS type,
            CASE
                WHEN sfu.standalone THEN fak.nama_fakultas
                WHEN sfu.kode_fakultas IS NOT NULL
                    AND sfu.kode_prodi IS NOT NULL THEN fak.nama_fakultas
                ELSE NULL
            END AS fakultas
        FROM sijamu_fakultas_unit sfu
        LEFT JOIN m_fakultas fak
            ON sfu.kode_fakultas = fak.kode_fakultas
        LEFT JOIN m_program_studi prod
            ON sfu.kode_prodi = prod.kode_prodi;
    `).Error

	if err != nil {
		t.Fatalf("migration failed: %v", err)
	}
	seedAllDokumenTambahan(t, gdb)

	cleanup := func() {
		sqlDB, _ := gdb.DB()
		resetDBDokumenTambahan(t, gdb)
		sqlDB.Close()
		// mysqlC.Terminate(ctx)
	}

	return gdb, cleanup
}

func resetDBDokumenTambahan(t *testing.T, gdb *gorm.DB) {
	gdb.Exec("SET FOREIGN_KEY_CHECKS=0")

	tables := []string{
		"m_fakultas",
		"m_program_studi",
		"sijamu_fakultas_unit",
		"users",
		"renstra",
	}

	for _, tbl := range tables {
		gdb.Exec("TRUNCATE TABLE " + tbl)
	}

	gdb.Exec("SET FOREIGN_KEY_CHECKS=1")

	seedAllDokumenTambahan(t, gdb)
}

func resetDBOnlyDokumenTambahan(t *testing.T, gdb *gorm.DB) {
	gdb.Exec("SET FOREIGN_KEY_CHECKS=0")

	tables := []string{
		"m_fakultas",
		"m_program_studi",
		"sijamu_fakultas_unit",
		"users",
		"renstra",
	}

	for _, tbl := range tables {
		gdb.Exec("TRUNCATE TABLE " + tbl)
	}

	gdb.Exec("SET FOREIGN_KEY_CHECKS=1")
}

func seedAllDokumenTambahan(t *testing.T, gdb *gorm.DB) {
	err := gdb.Exec(`
        INSERT INTO users (id, uuid, nidn_username, password, level, name, email, fakultas_unit, foto, email_verified_at, remember_token, reset_password, created_at, updated_at) VALUES
        (2, 'f524cbfd-b5aa-41d9-9a94-d9d5065918b4', 'admin', '123', 'admin', 'Admin', NULL, NULL, 'Tangkapan Layar 2025-10-17 pukul 09.48.28.png', NULL, '$2y$10$OKzK1M/XuKAuUCtfQ6FvBeZpyxjQYkMQnj8QcySVzY/cQ7Xk8s1hW', 0, '2023-05-18 18:45:52', '2025-10-22 06:02:49'),
        (40, '495fe283-3e42-4323-a172-c110036b0c60', 'Didik NotoSudjono', '$2y$10$uPFePuHtBrbp0FGwz82O/u.imhEDzr6C6ndQPHzSOfYyWo/e9d37u', 'auditor1', 'Prof. Dr. Ir. rer. pol. Didik NotoSudjono, M.Sc', 'didiknotosudjono@unpak.ac.id', NULL, NULL, NULL, NULL, NULL, '2023-10-02 02:45:45', '2024-11-30 08:45:29'),
        (43, 'd3d2b976-49c5-4fc8-8a78-a92484a97189', 'Istiqlaliah', '$2y$10$Ld0Ws48jrqWfONzmsu0bweSDl0WqCwmhr7SVwkyeUfVlVVA4bqdcO', 'auditor1', 'Dr. Istiqlaliah Nurul  Hidayati, M.Pd', 'istiqlaliah@unpak.ac.id', NULL, 'aku2.jpg', NULL, NULL, NULL, '2023-10-02 02:48:45', '2024-12-30 07:04:23'),
        (46, '496a2940-70eb-429a-88e6-e210babb323e', 'Eri Sarimanah', '$2y$10$jsegv/6gGVmrex3hxHnwo.vIav6hfnwYshQUJwk.OkJKbg1jkjPje', 'auditor1', 'Prof. Dr. Eri Sarimanah, M.Pd', 'erisarimanah@unpak.ac.id', NULL, NULL, NULL, NULL, NULL, '2023-10-02 02:49:22', '2024-11-30 08:45:40'),
        (49, 'd726a8ff-45b5-4be5-b76f-2edfe754881b', 'Yuary Farradia', '$2y$10$ingyBMcI42jlDlKunW/HO.5Ni4y1.nlonFjsMyL4dGjXTbnDnFYn6', 'auditor1', 'Dr. Ir. Yuari Farradia, M.Sc', 'yuary.farradia@unpak.ac.id', NULL, NULL, NULL, NULL, NULL, '2023-10-02 02:52:36', '2025-01-15 06:11:09'),
        (52, 'b09581fc-14c8-4b6d-8748-ceeb044288dd', 'Andi Chaerunnas', '$2y$10$QhSKATxTznNh9Dl2IBGlN.DyG7H2xvsIHJE8beWChxfrHP7WJQeoS', 'auditor1', 'Dr. Andi Chairunnas, M.Pd,.M.Kom', 'andichairunnas@unpak.ac.id', NULL, NULL, NULL, NULL, NULL, '2023-10-02 02:53:15', '2024-11-30 08:46:20'),
        (55, 'f3789a1f-3c9e-4d1f-8cf9-359184d93c14', 'Griet Helena', '$2y$10$jyCXqbg1.JPqFF4EBnXOcOYvFa8WuRaiFUfQEFOdXonQgsSNYrPOa', 'auditor1', 'Dr. Griet Helena Laihad, M.Pd', 'grihela@unpak.ac.id', NULL, 'FOTO GHL.jpeg', NULL, NULL, NULL, '2023-10-02 02:54:26', '2025-02-23 13:27:47'),
        (58, 'd0df6a80-6f66-4bd3-b844-b7f163cc0130', 'Agung Fajar', '$2y$10$lICelZarONq09QzEZVRQjeWMOvlm/7UMMcCjPpDQA4Nap8d4k6APC', 'auditor1', 'Dr. Agung Fajar Ilmiyono, SE.,M.Ak.,AWP.,CFA.,CAP', 'agung.fajar@unpak.ac.id', NULL, 'IMG_Foto Agung.jpg', NULL, NULL, NULL, '2023-10-02 02:55:01', '2024-12-23 03:41:52'),
        (61, '6052aa7c-31e4-4b17-bbd6-afe555b119c7', 'Edi Rohaedi', '$2y$10$uC7MyfUQijl0Yv.xTlILBeSOWQBd1IFTzWr3gxRw/XG8V3U1Qa2Ja', 'auditor1', 'Edi Rohaedi, SH.,MH', 'edi.rohaedi@unpak.ac.id', NULL, NULL, NULL, NULL, 1, '2023-10-02 02:55:25', '2024-11-29 12:53:03'),
        (64, '7c69c39e-8e70-460c-97d1-0aaa17d7f430', 'Prihastuti Harsani', '$2y$10$Ttz.xCyEQTV.ggSVcmIIueGRacOqfzTUdTyfdbhGnft2kymwcf/7S', 'auditor1', 'Dr. Prihastuti Harsani, M.Si', 'prihastuti.harsani@unpak.ac.id', NULL, 'Foto.jpg', NULL, NULL, NULL, '2023-10-02 02:56:02', '2024-12-16 09:38:10'),
        (67, 'e478505d-3b81-4a67-b924-7e5bda8bff1a', 'Herman', '$2y$10$EZshWaUkDra0ZbhWvNSlIuu5w2IYxPku3JM9VBncU9Ihh6kqOxYWG', 'auditor1', 'Dr. Herman, M.M', 'herman_fhz@unpak.ac.id', NULL, NULL, NULL, NULL, NULL, '2023-10-02 02:56:38', '2025-07-26 02:19:09'),
        (70, '56ce6c95-e23f-463b-bcf6-80fa4bea2a1e', 'Cantika Zaddana', '123', 'auditor1', 'Cantika Zaddana, S.Gz, M.Si', 'cantika.zaddana@unpak.ac.id', NULL, 'the newest.jpg', NULL, '$2y$10$61lLrJTmJjov7ys7n8N8V.s/CIk8TlTHpjUxnOKqmtgOp5sFHb/KS', NULL, '2023-10-02 02:57:04', '2025-08-11 08:03:05'),
        (73, 'b15f4b66-c696-40a4-a047-c51d2be63d4b', 'Indri Yani', '$2y$10$fsxG71.B7HYSPwnw4U7wVOu4YTZtKHUDJrOog1x59/0TGYmpWHyue', 'auditor1', 'Dr. Indri Yani, M.Pd', 'indri@unpak.ac.id', NULL, NULL, NULL, NULL, 1, '2023-10-02 02:57:34', '2024-11-29 12:20:50'),
        (76, '2107dc19-942e-41ed-ab73-27496f2cb72a', 'Indarini', '$2y$10$dhJ8GjS/iFeTo2wi3HDiy.PPzqLY4.ljrEwE4TieK9bLOIM/3UPRy', 'auditor1', 'Prof. Dr. Indarini Dwi Pursitasari, M.Si', 'indarini.dp@unpak.ac.id', NULL, NULL, NULL, NULL, 1, '2023-10-02 02:59:10', '2024-12-02 08:33:30'),
        (79, '9e7ed60b-75a1-4851-8490-fe9c46a9674d', 'Dolly Priatna', '$2y$10$MLDSXqLMI/3Zc7f8kMxN6.l3o4pzxNMPcPeqSouwUX8EKj5r4JN7O', 'auditor1', 'Dr. Dolly Priatna, M.Si', 'dollypriatna@unpak.ac.id', NULL, NULL, NULL, NULL, NULL, '2023-10-02 02:59:27', '2025-11-18 03:31:54'),
        (82, 'f74d3f0a-ed59-4f35-bdf9-ea0f0e784c35', 'Anna Permanasari', '$2y$10$mCGcvCnOJt9W1hQOOt7qQ.TeqpF8ZxSJe4eqcT0fnigudZeocXjyq', 'auditor1', 'Prof. Dr. Anna Permanasari, M.Si', 'anna.permanasari@unpak.ac.id', NULL, NULL, NULL, NULL, 1, '2023-10-02 02:59:44', '2024-11-30 02:44:03'),
        (85, '934c41f3-b31e-408f-b867-c919c4edbda2', 'Indarti', '$2y$10$jMZ2l6IgzjFgcmlpmOmGJOJv6xyEb4V9ujDQXMoeERNWesKM7h8Mq', 'auditor1', 'Dr.  Ir. Indarti Komala Dewi. M.Si', 'indarti@unpak.ac.id', NULL, NULL, NULL, NULL, NULL, '2023-10-02 03:00:12', '2024-12-20 03:29:29'),
        (88, '0a0e63d2-bb5a-4e34-9a74-5b8df78b51e0', 'Herdiyana', '$2y$10$Ojsaoli8D365aPwa8YDf8uGpDTMtu6gfmd/Sx6NY1gRExmVn9S0Ne', 'auditor1', 'Dr. Herdiyana, MM', 'herdiyana@unpak.ac.id', NULL, NULL, NULL, NULL, NULL, '2023-10-02 03:01:29', '2025-01-02 03:36:35'),
        (91, 'e290844e-5f5d-460f-b4d3-dcb9e4a8e962', 'Ade Heri', '$2y$10$B0ruVIuyfYNLd5Rxuben6OFaUPeTydzUPv3wV985J5TlohB6XAkjO', 'auditor1', 'Dr. Ade Heri Mulyati, M.Si', 'adeheri.mulyati@unpak.ac.id', NULL, NULL, NULL, NULL, 1, '2023-10-02 03:01:51', '2024-12-10 01:56:13'),
        (94, '3a0cbdb7-4d08-4d96-950c-be7341207bd3', 'Iwan Darmawan', '$2y$10$FPwzLhPlJk57Nvzmz62MK.rND1N2m5tGKoFkWwAj1lKncgcNdHjOK', 'auditor1', 'Dr. Iwan Darmawan, MH', 'iwan.darmawan@unpak.ac.id', NULL, NULL, NULL, NULL, NULL, '2023-10-02 03:02:21', '2024-11-30 08:54:15'),
        (97, 'bc46bf93-2899-4654-bee4-1f738da8d094', 'Irvan Permana', '$2y$10$NBKVlwCp9Rb9ZRbXvcNXGeDabp1wB70aRCWGvTiMGbhnDfgzPK.ze', 'auditor1', 'Dr. Irvan Permana, M.Pd', 'irvanpermana@unpak.ac.id', NULL, NULL, NULL, NULL, NULL, '2023-10-02 03:02:36', '2024-11-30 08:54:32'),
        (103, '0dd8b333-0a9c-462a-a84e-d34095d70916', 'Helen Susanti', '$2y$10$wM6644lf0gqLwBSYnk1vFOk89QCD8uu8vtiuG0LOhywEiY/3MbIua', 'auditor1', 'Helen Susanti, M.Si', 'helen.susanti@unpak.ac.id', NULL, NULL, NULL, NULL, NULL, '2023-10-02 03:03:18', '2024-11-30 08:54:50'),
        (106, 'a201b9b2-c94b-4ced-b886-81d7344cc789', 'Haqni Wijayanti', '$2y$10$FZg.TYn0tbAKcsL.9Ml/tufTzlep7VSSkWyEqkFWPRx/P.vsiapbi', 'auditor1', 'Haqni Wijayanti, M.Si', 'hagnijantix@unpak.ac.id', NULL, NULL, NULL, NULL, 1, '2023-10-02 03:04:00', '2024-12-08 13:43:16'),
        (112, '1cabd1d7-8a71-449b-b3a3-ee33a2d6f720', 'Rita Retnowati', '$2y$10$ncRLgGe9LLPM9P.e3NqP3.evq6mxWjWdepMjLerwS1lznwyPLX1Rq', 'auditor1', 'Dr. Rita Retnowati, MS', 'ritaretnowati@unpak.ac.id', NULL, NULL, NULL, NULL, NULL, '2023-10-10 02:17:53', '2024-11-30 08:55:34'),
        (115, '8d1a3f17-e1a5-4be2-bbf3-7bb991afc231', 'Hari Muharam', '$2y$10$gr5Sua95GJRIlvtGecwamexnPlKrRBGK2bNYNSBkIZbnaziH5iYQi', 'auditor1', 'Dr. Hari Muharam, S.E.,M.M', 'hari.muharam@unpak.ac.id', NULL, NULL, NULL, NULL, NULL, '2023-10-10 02:20:37', '2025-10-02 06:01:03'),
        (118, '668b7480-d2df-4592-83e5-65e5574d4344', 'Ellyn Octavianty', '$2y$10$Ft37tJ9cQiOCJKDASllrgeGuYnJCXeo8iy0WawYxeYv/nxnfsyhoK', 'auditor1', 'Ellyn Octavianty, SE,. MM', 'ellynoctavianty@unpak.ac.id', NULL, NULL, NULL, NULL, NULL, '2023-10-10 02:29:19', '2024-11-30 08:56:03'),
        (121, 'b5ed2f47-29e3-48be-97c0-359fb5eabaa6', 'Heny Purwanti', '$2y$10$4poGlC3wQ1h63/Y7bDFFZuivGASsoFEaT8ALhztETZqZnj4mEWtDO', 'auditor1', 'Heny Purwanti, M.T', 'henypurwanti@unpak.ac.id', NULL, NULL, NULL, NULL, NULL, '2023-10-10 02:29:49', '2024-11-30 08:56:18'),
        (124, '013fb4be-98d5-4f1b-a863-fb37aeb63cc9', 'Patar Simamora', '$2y$10$VHqgr8ldkVIo.oZf0JMcAOB2NurxmJndMLM4VXgQLf83yJri7PHyK', 'auditor1', 'Patar Simamora, S.E.,M.M', 'patar.simamora@unpak.ac.id', NULL, NULL, NULL, NULL, NULL, '2023-10-10 02:30:14', '2024-11-30 08:56:35'),
        (127, '7b69be28-5ba1-4493-8ee9-8f908be14869', 'Solihin', '$2y$10$NcQ7MK7TuQNFtjICSv/TF.Ieur1xoNK5/aON8Y.fDuoi5X5So8fuG', '', 'Solihin, M.T', 'solihin@unpak.ac.id', NULL, NULL, NULL, NULL, 1, '2023-10-10 02:30:42', '2025-10-21 08:44:48'),
        (151, '6c10337d-227f-4349-94ce-a998ae58d591', 'Istiqlaliah', '$2y$10$eGhZQKpMT0zDrJE0pkVwMO8T0NacgAN4A3gC8Q9ZZj4fPFVrJlcnK', 'auditor2', 'Dr. Istiqlaliah Nurul Hidayati, M.Pd', 'istiqlaliah@unpak.ac.id', NULL, NULL, NULL, NULL, NULL, '2023-10-16 02:13:49', '2025-01-17 04:06:43'),
        (154, '5971b74b-7a19-443d-97bd-f9c6fb34f4c3', 'Eri Sarimanah', '$2y$10$Gy37kYsYRjPjU9EmtS9z8OORQmSByxiHY8k.dK6XXPjwxg4SXWQJ.', 'auditor2', 'Prof. Dr. Eri Sarimanah, M.Pd', 'erisarimanah@unpak.ac.id', NULL, NULL, NULL, NULL, NULL, '2023-10-16 02:15:00', '2024-12-02 01:33:21'),
        (157, '931e98fe-b64f-4694-bab3-3f39aa342bdf', 'Yuary Farradia', '$2y$10$p.DsD1UpH7MSYK8RUeIDr.0ZE4rBzDpEL.dHp4reU3eWIOtV/GdHm', 'auditor2', 'Dr. Ir. Yuari Farradia, M.Sc', 'yuary.farradia@unpak.ac.id', NULL, NULL, NULL, NULL, NULL, '2023-10-16 02:15:18', '2024-11-30 08:46:02'),
        (160, '027f79d0-4071-4f41-be79-7fa8e7786b8b', 'Andi Chaerunnas', '$2y$10$NnPLz8wDZqsl6eeFrdadeuokvdoaU74uAIRGHqDfS0aXhrJK7vQ5O', 'auditor2', 'Dr. Andi Chairunnas, M.Pd,.M.Kom', 'andichairunnas@unpak.ac.id', NULL, NULL, NULL, NULL, NULL, '2023-10-16 02:15:33', '2024-11-30 08:46:27'),
        (163, 'd23dae62-24aa-4b3a-b6bc-b8abd8324401', 'Griet Helena', '$2y$10$Lpm3B5YTXfPxzEy7fSywKu8J78baQKHMHJp8Xn10UhyFy2sO9OUzi', 'auditor2', 'Dr. Griet Helena Laihad, M.Pd', 'grihela@unpak.ac.id', NULL, 'FOTO GHL.jpeg', NULL, NULL, 1, '2023-10-16 02:15:47', '2024-12-15 21:55:11'),
        (166, 'bc6f9d3f-8683-42af-ad16-16a52c8122f6', 'Agung Fajar', '$2y$10$2UczRfpQIwbSWzAKQDvSBeSjk1W2nCxvTqbgMFCOiR2NBU4HZc6Wm', 'auditor2', 'Dr. Agung Fajar Ilmiyono, SE., M.Ak.,AWP.,C.F.A.,CAP', 'agung.fajar@unpak.ac.id', NULL, 'IMG_Foto Agung.jpg', NULL, NULL, NULL, '2023-10-16 02:16:07', '2024-12-20 02:49:32'),
        (169, '4d6f64ab-b774-4c97-8c52-2a8bced2cd87', 'Edy Rohaedi', '$2y$10$NWQ8P.vWVdlPIblw9gEcZeN4P/DkUP8hJR5/1HGLginlIkSH1NzhW', 'auditor2', 'Edy Rohaedi, SH.,MH', 'edi.rohaedi@unpak.ac.id', NULL, NULL, NULL, NULL, NULL, '2023-10-16 02:16:20', '2024-11-30 08:51:52'),
        (172, '9d48fa00-1218-47fa-82a0-4841d107f36c', 'Prihastuti Harsani', '$2y$10$9HmAtDW7aNCQcWIEzfbp1Of6Iv6vkOdRycUupQIr/sEZUpZwokkN.', 'auditor2', 'Dr. Prihastuti Harsani, M.Si', 'prihastuti.harsani@unpak.ac.id', NULL, NULL, NULL, NULL, NULL, '2023-10-16 02:16:47', '2024-11-30 08:52:13'),
        (175, 'a410c93a-963d-4451-80df-add6cd8c3e86', 'Herman', '$2y$10$I5yrPAYhJ.TuMKC074XRfuh7JzS0ZVA1oZW/ababPS0crsWJSyZMa', 'auditor2', 'Dr. Herman, M.M', 'herman_fhz@unpak.ac.id', NULL, 'Herman_Foto.jpeg', NULL, NULL, NULL, '2023-10-16 02:16:58', '2025-07-28 04:28:46'),
        (178, 'ddd5af5d-db80-4723-b5ac-9c7f4547a65a', 'Cantika Zaddana', '$2y$10$YdU4yWZks3STeHfZ1ATLEeTl1qDVuM9oCbnOg8zVpmNOhEiHlvUb6', 'auditor2', 'Cantika Zaddana, S.Gz, M.Si', 'cantika.zaddana@unpak.ac.id', NULL, 'the newest.jpg', NULL, NULL, NULL, '2023-10-16 02:17:09', '2025-07-28 04:19:06'),
        (181, '4e8f88c9-c182-4f07-9655-10d60ced702e', 'Indri Yani', '$2y$10$s/H7aWFfQtt4zKigwfOiCum32qMZNif0OtgYMQOxDBAyEKINny/JS', 'auditor2', 'Dr. Indri Yani, M.Pd', 'indri@unpak.ac.id', NULL, NULL, NULL, NULL, 1, '2023-10-16 02:17:18', '2024-12-01 00:43:10'),
        (184, 'cc5c047b-cdd8-4029-8fb0-73bf684a100a', 'Indarini', '$2y$10$ahihMYC9oWJKy06gWNuIFuO0/0GjgbOAlSi8jHIQdoOM3kIJOS.4y', 'auditor2', 'Prof. Dr. Indarini Dwi Pursitasari, M.Si', 'indarini.dp@unpak.ac.id', NULL, NULL, NULL, NULL, NULL, '2023-10-16 02:17:29', '2024-11-30 08:52:34'),
        (187, 'e25d94f1-831d-4cfc-b0a6-2ec688f356c3', 'Dolly Priatna', '$2y$10$M1ss6YWCHT6liuvVH5mar.KnjGvXHDpInv11OdZ9/lVqiqXBNOc.q', 'auditor2', 'Dr. Dolly Priatna, M.Si', 'dollypriatna@unpak.ac.id', NULL, 'Dolly Priatna Foto Cop21 paris_2015.JPG', NULL, NULL, NULL, '2023-10-16 02:17:39', '2025-04-11 04:19:59'),
        (190, '323ab6c1-f993-4a1d-8800-123d1054b5f3', 'Anna Permanasari', '$2y$10$IDydlEycrG3m5V/WmJon4OoujMG3fi5DYrVIyYUPz0vczZDdRWJZ6', 'auditor2', 'Prof. Dr. Anna Permanasari, M.Si', 'anna.permanasari@unpak.ac.id', NULL, NULL, NULL, NULL, NULL, '2023-10-16 02:17:53', '2024-11-30 08:53:13'),
        (193, '4965c979-40c8-46f4-a443-f41e80db2326', 'Indarti', '$2y$10$xuM1b3lGEVsGp02vg7mnMOlyIJe5TxOAMptBuek4D9K83uqRiA66q', 'auditor2', 'Dr. Indarti Komala Dewi. M.Si', 'indarti@unpak.ac.id', NULL, NULL, NULL, NULL, NULL, '2023-10-16 02:18:10', '2024-11-30 08:53:30'),
        (196, '5cc94b9b-de6a-44b2-8998-7b25e7f2acf8', 'Herdiyana', '$2y$10$lnejzMvghFjgSD7aPe6pB.kthh2v1OISmCrD4yrxThmXdBykaM7oy', 'auditor2', 'Dr. Herdiyana, MM', 'herdiyana@unpak.ac.id', NULL, NULL, NULL, NULL, NULL, '2023-10-16 02:18:20', '2024-11-30 08:53:47'),
        (199, 'a7477de5-606d-442c-b90b-aa4405926e23', 'Ade Heri', '$2y$10$KUsCg4hycazNPNzbOr.E5eaCHbVIV3LFfB3/HU23y/U58ZXUjJORC', 'auditor2', 'Dr. Ade Heri Mulyati, M.Si', 'adeheri.mulyati@unpak.ac.id', NULL, 'FOTO JAS ADE HERI.jpeg', NULL, NULL, 1, '2023-10-16 02:18:30', '2024-12-16 11:32:40'),
        (202, '006e6a99-cc8c-4a94-879f-134a338e5dab', 'Iwan Darmawan', '$2y$10$kVEi2awkz21XXS5qm6Ce9OFxGaZjWNYquvpVK52I8XHvmiYCUNELW', 'auditor2', 'Dr. Iwan Darmawan, MH', 'iwan.darmawan@unpak.ac.id', NULL, NULL, NULL, NULL, NULL, '2023-10-16 02:18:42', '2024-11-30 08:54:22'),
        (205, 'c1158c21-88ac-453e-9eda-707ba6a4169b', 'Irvan Permana', '$2y$10$cPPacELnWOcUh7M/a7N66OsaQhNdzeiVJyK9S5RfI711gog8Z9THm', 'auditor2', 'Dr. Irvan Permana, M.Pd', 'irvanpermana@unpak.ac.id', NULL, NULL, NULL, NULL, NULL, '2023-10-16 02:18:53', '2024-11-30 08:54:39'),
        (208, '45d43d61-a7a4-450f-b9ea-f5ea4e42f7d6', 'Helen Susanti', '$2y$10$pNTeg.CubrZz7v8rDf.queIQob5NLKWu203jcZJHXqMKWOpQUR2NK', 'auditor2', 'Helen Susanti, M.Si', 'helen.susanti@unpak.ac.id', NULL, NULL, NULL, NULL, NULL, '2023-10-16 02:19:24', '2024-11-30 08:54:56'),
        (211, '6ca1c3c3-7f89-4b94-89cf-79d192ddf916', 'Haqni Wijayanti', '$2y$10$3ZS0kz08B2DXsKcQVX9.h.JtNgmrwBnPGbMwJylKPsod.JMHAuDe.', 'auditor2', 'Haqni Wijayanti, M.Si', 'hagnijantix@unpak.ac.id', NULL, NULL, NULL, NULL, NULL, '2023-10-16 02:19:35', '2025-03-20 04:14:47'),
        (214, '267df8fe-975b-4650-a9d7-619ea48d6a27', 'Rita Retnowati', '$2y$10$ylZt7LWLGuVqVnsgLueEQ..CjTJVB1hGasxdN7Ae4lpC/UlYiEGai', 'auditor2', 'Rita Retnowati', 'ritaretnowati@unpak.ac.id', NULL, NULL, NULL, NULL, NULL, '2023-10-16 02:19:46', '2024-11-30 08:55:39'),
        (217, 'a8885055-424c-4929-b3da-b7a02e755e57', 'Hari Muharam', '$2y$10$262w7uUNxbI3BzehL6hXV.50/k/Z3b2fCebHrgI7pNTLqmuREVtp6', 'auditor2', 'Dr. Hari Muharam, SE.,MM.,CIHCM', 'hari.muharam@unpak.ac.id', NULL, NULL, NULL, NULL, NULL, '2023-10-16 02:20:00', '2025-10-02 06:01:11'),
        (220, '9a1dff6c-056f-42b6-96f8-eaf963313e9f', 'Ellyn Octavianty', '$2y$10$9eF012HeKIveFdXo7wibwenBg4FY3Mwua.p6Y9q.GysWdWG9UXfa6', 'auditor2', 'Ellyn Octavianty, SE., MM', 'ellynoctavianty@unpak.ac.id', NULL, NULL, NULL, NULL, NULL, '2023-10-16 02:20:10', '2024-11-30 08:56:09'),
        (223, '91f53476-368a-49be-ac24-beb914af73f9', 'Heny Purwanti', '$2y$10$j./Fj3eS3ZiaxF8C09Zxs.olJ6r6Arcys2Nlm0j7Jqk74RhF9pNCW', 'auditor2', 'Heny Purwanti, M.T', 'henypurwanti@unpak.ac.id', NULL, 'FOTO.jpg', NULL, NULL, NULL, '2023-10-16 02:20:20', '2024-12-16 12:21:31'),
        (226, 'ce3f2a64-c9bb-453b-a036-ffe6fef7c809', 'Patar Simamora', '$2y$10$k9Bn3veF6Wj7dCxjUDGf3OryaWi9FvzTHNooQmcYI4EWNplJb2q6u', 'auditor2', 'Patar Simamora, SE., M.Si.', 'patar.simamora@unpak.ac.id', NULL, NULL, NULL, NULL, NULL, '2023-10-16 02:20:33', '2024-12-30 07:05:10'),
        (229, '00968a1d-992d-46e2-9ae4-6dae1a27ac0a', 'Solihin', '$2y$10$5onF1l43ZlfS3xK/8h.km.qfp.MNO0zJD8ch4zk3IFXY2FwvgGXey', 'auditor2', 'Solihin, M.T', 'solihin@unpak.ac.id', NULL, 'Solihin.jpg', NULL, NULL, 1, '2023-10-16 02:20:50', '2024-12-19 04:31:48'),
        (232, 'c7fd1d83-2d34-42a7-9cfe-38fa5f813188', 'Fakultas Hukum', '123', 'auditee', 'Dr. Eka Ardianto Iskandar, S.H., M.H. (Dekan FH)', 'fakultashukum@unpak.ac.id', 91, 'OIP (1).jfif', NULL, '$2y$10$9LPF7aie0Yc1F2kAp8yXRu9XyYJpmon1BfMm9ZrUPngDX7j3QvxC6', NULL, '2023-10-16 02:25:53', '2025-10-29 06:08:05'),
        (235, 'f033fb62-4a60-41f5-8909-ba3ca3df62e5', 'Ilmu Hukum', '123', 'auditee', 'Ari Wuisang, SH., MH. ( Kaprodi Ilmu Hukum S1)', 'fakultashukum@unpak.ac.id', 1, 'OIP (1).jfif', NULL, '$2y$10$gnvb6sNXoPUrkkpCHiwS9uLLY0xQm4LQ2L9osFjJ8j2b7Ke1IKw26', NULL, '2023-10-16 02:31:24', '2025-10-29 06:52:40'),
        (241, '94b19655-7c7c-444a-ae18-2479d0190519', 'FEB', '$2y$10$JYKed7BSLvrppl5/U91AUOYQk.C1B5yj1Jbe94GxAsH9AKKalr5yW', 'auditee', 'Towaf Totok Irawan, S.E., M.E., Ph.D (Dekan FEB)', 'dummy@gmail.com', 94, NULL, NULL, NULL, NULL, '2023-10-17 03:15:51', '2025-10-14 07:23:52'),
        (244, '05f396c7-663d-4449-a06b-ea01b48878d3', 'Akuntansi', '$2y$10$N6VxM0ovQHpgE7OHIAgs5ehuc8peG.RJLsKVd22Iy6UPb2YuNgfJu', 'auditee', 'Dr. Heru Satria Rukmana, S.E., Ak., M.M (Kaprodi Akuntansi)', 'akuntansi@unpak.ac.id', 4, NULL, NULL, NULL, NULL, '2023-10-17 03:24:31', '2025-10-11 03:15:10'),
        (247, '02f1b6d6-aff7-4ade-802a-52d9b4f29c2f', 'Manajemen', '$2y$10$j0C06V86F4NBXEITynbaBupk8ZVFtvFkLKPN3lKQA2/d0wri5vjMS', 'auditee', 'Prof.Dr. Yohanes Indrayono, Ak.,MM (Kaprodi Manajemen S1)', 'upmps.s1mjn@unpak.ac.id', 2, 'Logo Mjn New.jpg', NULL, NULL, NULL, '2023-10-17 03:26:11', '2025-10-04 05:18:10'),
        (250, '61aed0a7-a62e-423f-ba69-fed60bd3d0bb', 'Bisnis Digital', '$2y$10$ENlyFXwPHOVL.TBwxIlOLOnxS7sFQxG/stsGPZY2YFCQ7rZQN3iLO', 'auditee', 'Dr. Abel Gandhy, S.Pi., MM (Kaprodi Bisnis Digital)', 'bisnisdigital@unpak.ac.id', 3, 'bdi.png', NULL, NULL, NULL, '2023-10-17 03:32:18', '2025-10-28 03:36:42'),
        (253, '0a853a5f-0475-4b95-aa55-f9009b165771', 'FKIP', '$2y$10$99t67kLzuqg7FiQ.ypVGgOL6fBxEqEzS3UtQ10mEussugvrJgVdpy', 'auditee', 'Dr. H. Eka Suhardi, M.Si. (Dekan FKIP)', 'fkip@unpak.ac.id', 97, NULL, NULL, NULL, NULL, '2023-10-17 03:33:57', '2025-10-21 02:35:47'),
        (256, '5a663291-69fc-4694-9c02-f31a68879219', 'PBSI', '$2y$10$I1YOEDUqUD//yo1WVz8rKOSiUEhwzW5um6PH4IPN3G9YkYHjmIxwe', 'auditee', 'Stella Talitha, M.Pd. (Kaprodi PBSI)', 'fkip.indo@unpak.ac.id', 10, 'Stella Talitha.jpeg', NULL, NULL, NULL, '2023-10-17 03:36:08', '2025-10-30 01:32:00'),
        (259, 'f33b0300-e0f6-4251-bdec-995e178dcc9a', 'PBI', '$2y$10$dmqnLUUrlZEHkf/hRQbqXuXKaCm8VPDbIPlOlSUU7zG6Y3LvAc.46', 'auditee', 'Abdul Rosyid, M.Pd. (Kaprodi PBI)', 'englishedu.fkip@unpak.ac.id', 11, 'Abdul Rosyid_Foto.jpg', NULL, NULL, NULL, '2023-10-18 02:23:48', '2025-11-24 09:21:05'),
        (262, '9dc043cb-75a1-46f3-b210-f93fe5b91845', 'Pen.Biologi', '$2y$10$EIqS5RmbzDgGsw1H4FJ.r..Y3ZHUe/Lj4RLtdOUfnCyLYED99GOWq', 'auditee', 'Lufty Hari Susanto, M.Pd. (Kaprodi P.Biologi)', 'pendbiologifkip@unpak.ac.id', 5, 'WhatsApp Image 2025-06-20 at 15.33.04_1bad4e61.jpg', NULL, NULL, NULL, '2023-10-18 02:24:44', '2025-11-18 01:42:16'),
        (265, '191ae003-060c-42d3-8ac7-f0e6d5b110f3', 'PGSD', '$2y$10$9jLvr.47Rl.4TDbTmCkNBunHYa/x96PfG3dTHg9uv7dD.WOKpAIuy', 'auditee', 'Dr. Nita Karmila, M.Pd. (Kaprodi PGSD)', 'dummy@gmail.com', 7, NULL, NULL, NULL, NULL, '2023-10-18 02:25:27', '2025-10-04 05:20:57'),
        (268, '7af91608-8625-461e-b996-b027bd17fc55', 'Pen.IPA', '$2y$10$L7V4B2hifItRUhPM0OkMLuo8PykCHdzAtHMDUKJCCJkbYpaEGKrza', 'auditee', 'Lilis Supratman, M.Si. (Kaprodi P.IPA)', NULL, NULL, NULL, NULL, NULL, 0, '2023-10-18 02:26:04', '2023-10-18 02:26:04'),
        (271, '3c60ff48-2560-472b-b64d-7b93df32e452', 'PPG', '$2y$10$.PoPnL6lGPDuGtwqfyv5ouvwVpbNdlWl9LhkYd830BX./4y1tXTTu', 'auditee', 'Dr. Indri Yani, M.Pd (Kaprodi PPG)', 'ppg@unpak.ac.id', 9, NULL, NULL, NULL, NULL, '2023-10-18 02:27:00', '2025-10-30 03:52:57'),
        (274, '5d197aa4-39c8-4bf5-88de-aa64d7fdd8ae', 'Fisib', '$2y$10$XaYcZ/d9YAjreFKlyAH52.ikqE8jbiTahU2xC188gNt1FpvO7d./i', 'auditee', 'Dr. Muslim, M.Si (Dekan Fisib)', 'dekanfisib.2020@unpak.ac.id', 100, NULL, NULL, NULL, NULL, '2023-10-18 02:28:09', '2025-10-04 05:21:36'),
        (277, 'f340f1ab-c41d-4c88-ba00-7d9c427c5e7e', 'Sastra Inggris', '$2y$10$a8Tmi80DmTPnE.1zvNmvj./xe3RRr.wAA/nNbLxYFfMY0uMaJ3wWK', 'auditee', 'Dyah Kristyowati,S.S.,M.Hum. (Kaprodi Sas.ing)', 'dummy@gmail.com', 14, NULL, NULL, NULL, NULL, '2023-10-18 02:33:50', '2025-10-04 05:21:46'),
        (280, '3d73601b-6fb1-45e3-bf41-377e598259f1', 'Sastra Jepang', '$2y$10$XbWkZ2dchRgIlPIMqCuAruTSJ36LPZcxlQvmc9u4CrPckZzLTkhVa', 'auditee', 'Mugiyanti, M.Si (Kapordi Sas.Jep)', 'fisib.sasjep@unpak.ac.id', 15, NULL, NULL, NULL, NULL, '2023-10-18 02:35:50', '2025-10-04 05:21:54'),
        (283, 'c884277d-d314-4eee-a9fe-7e9e867155be', 'Sastra Indonesia', '$2y$10$PPNip.mCc6n1vnfePsJT5uiulw7YjJBwmwv34Y7yV5GHEyLWiH3fe', 'auditee', 'Drs. Sasongko Suharto Putro, M.M. (Kaprodi Sas.In)', 'prodisastraindonesiaunpak@gmail.com', 13, NULL, NULL, NULL, NULL, '2023-10-18 02:37:41', '2025-10-22 07:29:19'),
        (286, '4e50c6de-d74a-4ecf-91fe-f18ead93ec70', 'I.kom', '$2y$10$i4Cdm/vLB.13TpgBnQcrYOlb.7bDkahySy2YS6U/DQBozvK.NX/km', 'auditee', 'Ratih Siti Aminah., M.Si. (Kaprodi I.Kom)', 'rinifirdaus@unpak.ac.id', 12, 'Wajah 2.jpg', NULL, NULL, NULL, '2023-10-18 02:45:04', '2025-10-04 05:22:19'),
        (289, '37e90b8b-ac2e-40d7-81f3-1d60103f12f1', 'FTeknik', '$2y$10$VE6X7dMwQLGMzNhbzxFMUuZzedcAv1sxL.oAQRr.4BpU33Pf7Xyqe', 'auditee', 'Dr. Ir. Lilis Sri Mulyawati, M.Si (Dekan FT)', 'mutuft@unpak.ac.id', 103, NULL, NULL, NULL, NULL, '2023-10-18 03:10:28', '2025-10-17 03:35:40'),
        (292, '490f1f08-a28e-44cc-9f13-01e65d39a48f', 'T.Geologi', '$2y$10$bXTc1BuOZjVjM4SNj/e.8Ot4ZNHZQJpVlzsfV9lGcV4qf9cQPX.pq', 'auditee', 'Helmi Setia Ritma P., ST., M.Si. (Kaprodi Geologi)', 'solihin@unpak.ac.id', 19, 'WhatsApp Image 2025-11-06 at 14.00.27_f6c8fe96.jpg', NULL, NULL, NULL, '2023-10-18 03:11:39', '2025-11-06 07:01:50'),
        (295, 'c67062ae-3cf8-46a6-832f-fecd2522cafa', 'PWK S1', '$2y$10$gykHif0r19MQZJFaxiXXOueXwyHFAdf3sBkac6k5dA48aHrXFj0fS', 'auditee', 'Dr. Mujio, S.Pi., M.Si (Kaprodi PWK S1)', 'prodipwk@unpak.ac.id', 20, 'Mujio.png', NULL, NULL, NULL, '2023-10-18 03:15:07', '2025-11-06 13:37:58'),
        (298, '2a14747d-f8c7-4693-9d95-868f17c135b6', 'T.Sipil', '$2y$10$AJQRgqp9nxPMQz5ruhNb8eqO7FyLUmBoR9NfkQuanq9edyAc.NxI2', 'auditee', 'Ir. Wahyu Gendam P, STP., M.Si (Prodi T.Sipil)', 'mutuft@unpak.ac.id', 17, 'Screenshot 2025-10-23 193738.png', NULL, NULL, NULL, '2023-10-19 02:01:22', '2025-10-23 12:37:52'),
        (301, 'bfd1f6c6-c468-4d98-b414-3269d3c4ba66', 'T.Elektro', '$2y$10$WTSJE/Ew78A6V4V.UQVxi.1PO9hadLo7qob3BKOVvYd6Y1kCth/4C', 'auditee', 'Ir. Yamato, M.T (Kaprodi T.Elektro)', 'waryani@unpak.ac.id', 16, 'Capture.JPG', NULL, NULL, NULL, '2023-10-19 02:03:06', '2025-10-04 05:23:18'),
        (304, '74db7760-2b94-4cf1-befb-e32a7215969b', 'T.Geodesi', '$2y$10$Nl6GnpPtKwt3yJGaAd9X0uh15g7lnYg5NAZ4v5gNjrZq7UFDNTNQi', 'auditee', 'Mohamad Mahfudz, ST., MT. (Kaprodi T.Geodesi)', 'prodi_geodesi@unpak.ac.id', 18, 'kap_gd.png', NULL, NULL, NULL, '2023-10-19 02:04:27', '2025-11-07 03:38:56'),
        (307, '2e24b19d-d7a2-4018-8182-e0e973875d20', 'FMIPA', '$2y$10$Wa.T9Ktxf2zH/vlscTqUlOGGP/Bhb9qwNNQFwxPftJ0ipixjq8IWy', 'auditee', 'Asep Denih, S.kom., M.Sc., Ph.D. (Dekan Fmipa)', 'asep.denih@unpak.ac.id', 106, NULL, NULL, NULL, NULL, '2023-10-19 02:06:37', '2025-10-04 05:23:52'),
        (310, '9eaec77f-0fa0-4d34-9da6-18b28ae0e410', 'Biologi', '$2y$10$xwPB65hlQIHQyh6bEv4A2Os.OvSphMt8U0heQUZKOmvNy3Y.FXtki', 'auditee', 'Dra. Triastinurmiatiningsih, M.Si', 'dummy@gmail.com', 22, NULL, NULL, NULL, NULL, '2023-10-19 02:08:31', '2025-10-20 06:16:21'),
        (313, 'b6310844-f193-4da0-9fd9-641f9978119c', 'Kimia', '$2y$10$Jw3Ed9Cgt5iFBr.dge8kCOwVbBhN2ZMFTlELc38rAmxvZajguVIwG', 'auditee', 'Dr. Uswatun Hasanah, S.Si., M.Si. (Kaprodi Kimia)', 'kimia@unpak.ac.id', 23, 'Screenshot 2025-10-29 at 10.58.07.png', NULL, NULL, NULL, '2023-10-19 03:21:12', '2025-10-29 03:58:54'),
        (316, '010081f6-93a5-4e43-8ab3-484aa1edc0f8', 'Matematika', '$2y$10$7Fc3q5JKhR.dwDY1JqX.C.DgO1niTN/94AFHuUSUPoEVqgdMqPT2a', 'auditee', 'Dr. Embay Rohaeti, S.Si., M.Si. (Kaprodi MTK)', 'matematika@unpak.ac.id', 21, NULL, NULL, NULL, NULL, '2023-10-19 03:22:07', '2025-10-04 05:24:27'),
        (319, 'faac64ea-d23d-4544-9b5d-4a21eb3fb83a', 'Ilkom', '$2y$10$b0mxxIzVkTQgLrAoTbVpquxFUi0pxGKLfzNyNrrUbLOXB5oU7p6wm', 'auditee', 'Dr. Fajar Delli Wihartiko, S.Si., M.M., M.Kom (Kaprodi Ilkom)', 'akreilkom@unpak.ac.id', 25, NULL, NULL, NULL, NULL, '2023-10-19 03:53:48', '2025-10-04 05:24:41'),
        (322, 'b211988d-9e37-4d16-a836-4ee19866906d', 'Farmasi', '$2y$10$inii6GEuCjDK2CUIhygo6uT94ZyhLZZJy.6TblqAvluYYU6qqLkB.', 'auditee', 'apt. Emy Oktaviani, S.Farm., M.Clin.Pharm. (Kaprodi Farmasi)', 'cyntiawahyuningrum@unpak.ac.id', 24, NULL, NULL, NULL, NULL, '2023-10-19 03:54:17', '2025-10-04 05:24:50'),
        (325, 'a178d078-07e2-4b44-8d72-06bea0b08473', 'Pascasarjana', '$2y$10$QToSV66aAzlnSvVsd/b0FOnHFUnrRisVgOxV8il/qp0oWiDx4IsF2', 'auditee', 'Prof. Dr. Sri Setyaningsih, M.Si (Dekan Pascasarjana)', 'pasca@unpak.ac.id', 109, NULL, NULL, NULL, NULL, '2023-10-20 04:15:35', '2025-10-04 05:25:08'),
        (328, '74c7e91f-e480-4c73-8133-43e09bc7c78b', 'MP S3', '$2y$10$h8hDBbEPtnm6uwU6iL9iEe3tDxFPWTvR9sLUra876BkRwsSB292uy', 'auditee', 'Dr. Suhendra, M.Pd. (Kaprodi MP S3)', 'dummy@gmail.com', 31, NULL, NULL, NULL, NULL, '2023-10-20 04:17:29', '2025-10-04 05:25:20'),
        (331, '2b5c8e47-e551-40f3-bd80-122f029d1e61', 'Manajemen S3', '$2y$10$IpWslWUPWI.qH1sfkAJ4y.3d7l3MrUE63RmhAVcObRosentB1zPtu', 'auditee', 'Dr. Nancy Yusnita, SE., MM  (Kaprodi Manajemen S3)', 'dummy@gmail.com', 26, 'Nanc.jpg', NULL, NULL, NULL, '2023-10-20 04:18:26', '2025-10-22 07:00:56'),
        (334, '8a970297-6756-42d1-8cca-691b795a2607', 'MP S2', '$2y$10$ZlbC6RizafqHg5lvcK8u/ur5hPtOBisX0qysnpBcfWo2tbuJKWJjK', 'auditee', 'Dr. Lina Novita, M.Pd. (Kaprodi MP S2)', 'dummy@gmail.com', 32, NULL, NULL, NULL, NULL, '2023-10-20 04:19:54', '2025-10-04 05:25:40'),
        (337, 'fbacc3fb-ce9b-4eea-afcb-c0755a0a8f69', 'ML S2', '$2y$10$UrexCNIJ/8brJw3Sl7luSOdH1.mVKWkaDF18lKcG097DEns9AXZfG', 'auditee', 'Dr. Rosadi, SP, MM (Kaprodi ML S2)', 'dummy@gmail.com', 28, NULL, NULL, NULL, NULL, '2023-10-20 04:20:39', '2025-10-04 05:25:47'),
        (340, 'cb0ab124-244a-4378-9fc6-c3ebcfa96ae3', 'Ilmu Hukum S2', '$2y$10$H4FsiRmvZJQi4gMHDug7WuL/rO903Ubl2y4fWCPK7OwsHDe/YaZ3e', 'auditee', 'Dr. Iwan Darmawan, SH., MH. (Kaprodi Ilmu Hukum S2)', 'dummy@gmail.com', 29, NULL, NULL, NULL, NULL, '2023-10-20 04:21:23', '2025-10-22 06:46:26'),
        (343, 'dcde3652-22ae-439b-b2bc-2dd1907c26e3', 'Manajemen S2', '$2y$10$IAMMrq7xpwm4CjTxwRs1zuMbrOY0GNIotCqTWGj1qOh6sVzhvQI5u', 'auditee', 'Dr. Agus Setyo Pranowo, MM., SE. (Kaprodi Manajemen S2)', 'dummy@gmail.com', 27, NULL, NULL, NULL, NULL, '2023-10-20 04:22:12', '2025-10-04 05:26:07'),
        (346, '569e5c96-8bf8-42cd-805d-8895545c16bd', 'IPA S2', '$2y$10$glsT8Xpsw1xWSE7bYTcuNuh8vuWuqPY5ljoYWLqFhNZpfKBQmmkzu', 'auditee', 'Dr. Didit Ardianto, M.Pd. (Kaprodi IPA S2)', 'dummy@gmail.com', 30, NULL, NULL, NULL, NULL, '2023-10-20 04:23:01', '2025-10-04 05:26:15'),
        (349, '3951ef54-a3ec-4d19-aac0-f522e276bf6c', 'PWK S2', '$2y$10$YjmFMhfBAPwHBiYXNTIwIOZ77/tLoB/c5RjgVOcoQO6s/7nNShWQS', 'auditee', 'Dr. Ir. Anugrah, M.Si. (Kaprodi PWK S2)', 'dummy@gmail.com', 35, NULL, NULL, NULL, NULL, '2023-10-20 04:23:32', '2025-10-04 05:26:27'),
        (352, '13bb7063-49a4-485c-8b5f-d2d5e29c6e23', 'PENDAS S2', '$2y$10$WbTeNjiAcealTWgejX5bxOfIDJAkpvl5wWpz8ILtkmMWHc9bcTpSy', 'auditee', 'Dr. Tustiyana Windiyani, M.Pd. (Kaprodi PENDAS S2)', 'tustiyana@unpak.ac.id', 33, 'Foto mamah windi 4 x 6 (1).jpeg', NULL, NULL, NULL, '2023-10-20 04:24:17', '2025-10-04 05:26:35'),
        (355, '0c843eae-d26f-4d66-abd2-8f27fd72e759', 'Svokasi', '$2y$10$1NtL4tJs1/w4r9CRtF5q4.JydIXxQJ0Y4natXVPBIwSwd3u6qUhCG', 'auditee', 'Dr. Lia Dahlia Iryani, S.E., M.Si (Dekan SV)', 'sonniadarmasih123@gmail.com', 112, 'hehe.php', NULL, NULL, NULL, '2023-10-21 01:49:01', '2025-10-04 05:27:45'),
        (358, 'f91ca73e-3b11-42dd-b38e-ebeb42e084ae', '0413117601', '$2y$10$zO51umOGlCIrOD9QS61yu.g0snHAbt33BZd2H0BwlffgDNFuqbNsO', 'auditee', 'Dr. Lia Dahlia Iryani, SE., M.Si ( KaProdi Akuntansi D3)', 'dahlia.iryani@unpak.ac.id', NULL, 'logo vokasi.jpg', NULL, NULL, 1, '2023-10-21 01:50:13', '2024-12-21 01:55:43'),
        (364, 'e95cb68a-f522-4370-8dfd-50154559b91a', 'perpajakan D3', '$2y$10$R1uZF6/WbF0Gmra7nuDsDupgRJyWt8tZhMMbwAFJi2plRL0EK278m', 'auditee', 'Chandra Pribadi, Ak., M.Si., CPA. (Kaprodi Perpajakan D3)', 'dummy@gmail.com', 38, NULL, NULL, NULL, NULL, '2023-10-27 03:40:58', '2025-10-04 05:28:34'),
        (367, '2f7747e3-7780-4801-8c6a-a8d50cbfafaf', 'MPK D3', '$2y$10$WVvjKmREhipH3qnw.ZQqMu5uS/gJX6eczgOrKv35E19uv6j..qGTe', 'auditee', 'Djoko Hardjanto, S.Pt, M.Si (Kaprodi MKP D3)', 'd3.mkp@unpak.ac.id', NULL, NULL, NULL, NULL, 1, '2023-10-27 03:43:34', '2024-12-13 08:43:20'),
        (370, '16acae5a-b38a-4752-a28e-6a05c5bcd5cc', 'T.Kom D3', '$2y$10$8WZTMqtZzOlcO/PCbFDOQupDUvECWvneNQgjzRt4dHAX7I0NMcMM6', 'auditee', 'Akbar Sugih Miftahul Huda, M. Kom. (Kaprodi Teknik Komputer D3)', 'dummy@gmail.com', 36, NULL, NULL, NULL, NULL, '2023-10-27 03:44:50', '2025-10-22 03:37:06'),
        (373, '7678c084-2d6d-4c49-b774-3d14f0032845', 'S.Informasi', '$2y$10$vAQBpJ./w9h0cL4IZsdXMOCaRL.1ngKGKNiw.d/HJb/WlZFASj9ku', 'auditee', 'Dian Kartika Utami, M.Kom (Kaprodi Manajemen Informatika d3)', 'sv_si@unpak.ac.id', 37, 'C83C3368 crop.JPG', NULL, NULL, NULL, '2023-10-27 03:46:13', '2025-11-13 04:19:46'),
        (376, '1f358061-1244-4d20-871b-76bee78249fe', 'aries', '$2y$10$o.6Hly1CHcBQ4hS0.hg7q.0FLDx9kC4jLvUuEX9s4CORZjsHTDyE2', 'auditor1', 'Aries Maesya', 'adamilkom00@gmail.com', NULL, NULL, NULL, NULL, 1, '2023-10-27 11:39:16', '2024-11-29 11:22:59'),
        (379, 'be0785c7-e97b-4e77-b102-ee82242f7811', 'KKPKT', '$2y$10$nFz1VGY/wBbBSkOkKff4E.Msn6xJf1Mn5K0oSxmI/9nPZQTkR5BGu', 'auditee', 'KKPKT', NULL, NULL, NULL, NULL, NULL, 0, '2023-11-03 01:55:23', '2023-11-03 01:55:23'),
        (382, '51347d7b-a70e-4228-bf67-d0feed84a8d0', 'LPPM', '$2y$10$.r1mSpN2ZkMnFoSSRWMkA.rB9R8kNiScU5Zk.YMop1MjX2XjFxlDm', 'auditee', 'Dr. Dolly Priyatna, M.Si', 'dollypriatna@unpak.ac.id', 70, NULL, NULL, NULL, NULL, '2024-01-08 06:06:20', '2025-09-01 06:09:27'),
        (385, '2036e3fe-5b80-4daa-a584-4b56e430aac1', 'UNPAK PRESS', '$2y$10$iTd7C7cPIFxjrMyJAp9ppetWAkLQS.TIL42.uU54SJCOlNBZ9nFfe', 'auditee', 'Nina Agustina ,S.E.,M.E', 'dummy@gmail.com', 75, NULL, NULL, NULL, NULL, '2024-01-08 06:07:07', '2025-09-01 06:08:56'),
        (388, '368df968-365d-4b35-ad48-85670141b932', 'HUMAS', '$2y$10$9bJVYOSArUTdOQMRY573ru4F.ZkzaPQkIejN8rUr9KinYsPmeDkuO', 'auditee', 'Aditya Prima Yudha, S.Pi.,M.M.', 'dummy@gmail.com', 72, NULL, NULL, NULL, NULL, '2024-01-08 06:07:36', '2025-09-01 06:11:02'),
        (391, '38629163-55fe-41b4-892a-65319b7474d6', 'PERPUS_PUSAT', '$2y$10$gk0H/5NydAKvXBacuSlKmu0cXdwOzhtoYuzbRWR1fCcBrCcTNvfyS', 'auditee', 'Wildan Fauzi Mubarock M. Pd', 'dummy@gmail.com', 71, NULL, NULL, NULL, NULL, '2024-01-08 06:08:04', '2025-09-01 06:17:38'),
        (394, 'f922e886-a566-43bf-b0ea-676f5436b20a', 'INOVASI', '$2y$10$eckTCXLOu/RzN5FX/IkW9uEZMCxw3TcJB.i744CtYhfbXRP65SxRe', 'auditee', 'Asep Saepulrohman, M.Si', 'asepspl42@gmail.com', 88, NULL, NULL, NULL, NULL, '2024-01-08 06:08:30', '2025-10-13 02:55:15'),
        (397, '3d74f809-eb48-4930-9745-9a842f0bf7d7', 'INKUBATOR', '$2y$10$u6rmC98w0a.xOCtMkoEsGufnb1rdFrKfIgF3LQYfy24.wkjCQr7Ei', 'auditee', 'Asep Saepulrohman, M.Si', 'asepspl42@gmail.com', 89, NULL, NULL, NULL, NULL, '2024-01-08 06:09:01', '2025-10-13 02:55:25'),
        (400, '23b40c60-2e73-4042-881e-3b4b09e721c9', 'KEMITRAAN', '$2y$10$T8J5cfA3e2kEhakjA7EKYesqCr4shHNVCB/PDKxfdMYLjA7zllDnC', 'auditee', 'Cucu Mariam, M.Pd', 'dummy@gmail.com', 125, NULL, NULL, NULL, NULL, '2024-01-08 06:09:32', '2025-09-01 06:12:04'),
        (403, '9144f43b-d633-43c1-be6e-a9bb5ca495a9', 'KARIR', '$2y$10$s3EjXRoNPAO1qcG4flXc7.kpQbnD8XtSfpNbwcfGae/i2S7Qxflgu', 'auditee', 'Dr. Herman, SE., MM.', 'dummy@gmail.com', 68, NULL, NULL, NULL, NULL, '2024-01-08 06:09:55', '2025-09-01 06:12:58'),
        (406, 'f2c8ec89-d7dd-4fdd-84bb-954bb46a1d26', 'BPSI', '$2y$10$5H2cyGbrhV0NDNS1u9GjvO8avGTOTze069.U3Y5tQlXaz0/hqR5iK', 'auditee', 'Aries Maesya, M.Kom', 'putik@unpak.ac.id', 41, NULL, NULL, '$2y$10$7EZZYmgnfQ2oKb4QD1ZyIOPmL8pg3f4Mma1vnzr8dXgfLWLCgmY5O', NULL, '2024-01-11 07:28:00', '2025-09-12 03:18:28'),
        (407, 'f45b26d8-abd8-4ca9-891f-32c8e0347a1c', 'BAAK', '$2y$10$/5Q9ilEgbJ5iGuZ3rABJJe419Iptj1n5kpOzyItwHFIWVXxGnHyva', 'auditee', 'Dr. Eka Ardianto Iskandar, S.H., M.H.', 'unpak@gmail.com', 91, NULL, NULL, NULL, NULL, '2024-01-26 04:23:35', '2025-10-04 06:32:15'),
        (410, 'eb827018-2137-432a-86b8-4626fc9ba9e8', 'BAUM', '$2y$10$lomy6xBGmv0nIDkpUToxjO0m3uWPLoedAj.hCJJ/PhFK3RkaCMCQK', 'auditee', 'Wijaya Kusumah, S.E. (Kepala BAUM)', NULL, NULL, NULL, NULL, NULL, 0, '2024-01-26 04:23:48', '2024-02-02 01:42:32'),
        (413, '832a3472-e91c-4196-8392-0544653b65ce', 'LPM', '$2y$10$wF76ZOkxWSY42maN5licUu8BKbJsPb3BEaV6zNNEtH7QRYp0ixcx6', 'auditee', 'Dr. Diana Widiastuti, M.Phil', 'lpm@unpak.ac.id', 69, NULL, NULL, NULL, NULL, '2024-02-02 02:20:14', '2025-08-29 04:40:11'),
        (415, '63b1c4b2-5e13-407f-a9fc-a8c775d9ecaa', 'Feri Ferdinan', '123', 'auditor2', 'Dr. Feri Ferdinan Alamsyah, M.I.Kom.', 'feriferdinan@unpak.ac.id', NULL, '003.jpeg', NULL, '$2y$10$Xi2k9oUCDaYWjSR0KX0GceYECXGvDMslSJ09x1.K4yTf8x/VYN4qq', 1, '2024-02-15 05:53:15', '2024-12-14 01:22:45'),
        (418, '160f722d-33db-4ead-909b-097076c482d7', 'aries2', '$2y$10$/sEjEdpfcHfl06vV3VBoRepVXIS9OvAImHmALbBltllp/5iqhdzX.', 'auditor2', 'aries2', 'adamilkom00@gmail.com', NULL, NULL, NULL, NULL, 0, '2024-02-15 05:53:15', '2024-11-29 07:00:28'),
        (421, 'f818e8f1-690c-4b66-8b5a-3a567ff78fb6', 'Mahfudz', '$2y$10$x6g4bNeJWmxOsGM14CzOIO3BgcgoULrtLvb4StI36EyBVa1VVa9Pe', 'auditor2', 'M. Mahfudz, M.T', 'mohamadmahfudz@unpak.ac.id', NULL, NULL, NULL, NULL, NULL, '2024-09-18 02:17:37', '2025-11-26 03:08:11'),
        (425, 'b6aa8939-ed31-4ff6-982b-ee7d7998de2f', 'testauditee', '$2y$10$/rcTO5xNMGIecshQKDabFe2zSUZ2hQGJrI37aZneWlIoVFMriupsO', 'auditee', 'testauditee', 'testauditee@gmail.com', 91, NULL, NULL, NULL, NULL, '2024-10-12 02:52:23', '2025-11-20 14:34:53'),
        (428, 'de804d20-95a9-4766-939a-86260e0b987c', 'testauditor1', '$2y$10$VOPdhlRek2vzKfkBg5wYqeDCl4IyS8biHoZRcPNmc4QPBXxyyFb3a', 'auditor1', 'testauditor1', 'testauditor1@gmail.com', NULL, NULL, NULL, NULL, NULL, '2024-10-12 02:52:55', '2025-11-20 14:35:21'),
        (431, '77fc3092-3131-4881-8e9f-84b974c0c724', 'testauditor2', '$2y$10$yXuEfDS8owHNLilpdm9WfOF3wwlhkmqWaCV8XKoH0N.bE7ZPF0Pba', 'auditor2', 'testauditor2', 'testauditor2@gmail.com', NULL, NULL, NULL, NULL, NULL, '2024-10-12 02:53:04', '2025-11-20 14:35:46'),
        (452, 'ff73ce69-5545-41ba-aac8-4567dee6fd78', 'Feri Ferdinan', '$2y$10$4k03/IVfdC.SUFwoMATG3.2Zw8SLGHI2rUQHZvl1AJ/ovfbWYPFmO', 'auditor1', 'Dr. Feri Ferdinan Alamsyah, M.I.Kom.', 'feriferdinan@unpak.ac.id', NULL, NULL, NULL, NULL, 0, '2025-02-04 08:20:48', '2025-02-04 08:20:48'),
        (455, 'e8603e5f-dad4-4c23-82e1-235e458f6e10', 'Mahfudz', '$2y$10$QFZ2xuoSDdI7CXwfLsariuX8GgqHs22dVsL4tFhvOTa7f7eWvnLlS', 'auditor1', 'M. Mahfudz, M.T', 'mohamadmahfudz@unpak.ac.id', NULL, NULL, NULL, NULL, NULL, '2025-02-04 08:21:20', '2025-11-26 03:07:51'),
        (462, 'd38a231f-3e58-47f6-94c9-d3f05a9c021c', 'Diana Widiastuti', '$2y$10$0jfbqSIIwAA61HQ29L2MAeN7T3YwEELFP/90i.kMUv/MXLea9D0Jy', 'auditor1', 'Dr. Diana Widiastuti, S.Si, M.Phil', 'diana@gmail.com', NULL, NULL, NULL, NULL, NULL, '2025-08-22 03:03:41', '2025-11-18 03:07:55'),
        (463, '6de81337-43aa-43da-be46-48efbdc73061', 'Diana Widiastuti', '$2y$10$MX0iu.T.hZp6SPzxowbo5.DPH4HXnSe0pGsz0UdrnR5YQmOqUf9UO', 'auditor2', 'Dr. Diana Widiastuti, S.Si, M.Phil', 'diana@gmail.com', NULL, NULL, NULL, NULL, NULL, '2025-08-22 03:04:18', '2025-11-18 03:21:38'),
        (464, 'de549dd1-68fa-40e2-ab26-7a67e4afcaf9', 'Muhammad Fathurrahman', '$2y$10$XP40Y3XnMkbXin4kKQ4Q9.EpfTHJFPiyp91ZfTMZmaHg7S4UY1Ct.', 'auditor1', 'Dr. Muhammad Fathurrahman, S.Pd., M.Si.', 'fathur110590@unpak.ac.id', NULL, NULL, NULL, NULL, NULL, '2025-08-22 03:05:03', '2025-11-18 03:25:41'),
        (465, 'f712e038-6974-4cca-a994-c7e8d013d024', 'Muhammad Fathurrahman', '$2y$10$94/7Jt7BWU.oB7jOo5cGzeXJVfHfTDmIky7a7XVfOm7sRde66yjTC', 'auditor2', 'Dr. Muhammad Fathurrahman, S.Pd., M.Si.', 'fathur110590@unpak.ac.id', NULL, NULL, NULL, NULL, NULL, '2025-08-22 03:05:27', '2025-11-18 03:07:31'),
        (466, 'b9a9dd25-5c69-4229-ae7c-834f14254e9d', 'Kecerdasan Buatan dan Robotika', '$2y$10$AAwGvQ2IsbuA69LQNoGGUeHM3F8QSx/aOdYWiXnAIdB.T2eZM5xxu', 'auditee', 'Fikri Adzikri, S.T., M.T. (Ka Prodi Kecerdasan Buatan dan Robotika)', 'dummy@gmail.com', 121, NULL, NULL, NULL, NULL, '2025-08-25 04:54:19', '2025-10-04 05:30:02'),
        (467, '9a74cdf9-e614-41e3-bcc2-e61c6b902e23', 'Ilmu Komputer S2', '$2y$10$l.FSqUg2sHKY7KaP4VM2ie1sx2Dan9vqCkbevSCfwYs9eARBAgdjy', 'auditee', 'Dr. Eneng Tita Tosida, M.Si. (Ka Prodi Ilmu Komputer S2)', 'dummy@gmail.com', 126, NULL, NULL, NULL, NULL, '2025-08-26 04:25:15', '2025-10-04 05:26:53'),
        (468, '70575659-ce43-4759-b2aa-de281ef4641b', 'PPA', '$2y$10$XeNkzrKmAW9ICXIbtz0xZessMKo7caCyyjTcxZedCTAjVJtEz8PZq', 'auditee', 'Dr. apt. Bina Lohitasari, M.Pd., M.Farm', 'binalohitasari@unpak.ac.id', 127, '3x4.jpg', NULL, NULL, NULL, '2025-08-26 04:27:15', '2025-11-12 08:55:50'),
        (469, 'aa3e8268-ff45-470c-836b-f39a62389548', 'E-Journal FKIP', '$2y$10$fKdarrKnt4UXV.vBZaqdrOm61XC16dN/UffGZLlHKb2O/Ra0FmOiW', 'auditee', 'E-Journal (Unit FKIP)', 'dummy@gmail.com', 78, NULL, NULL, NULL, NULL, '2025-10-07 03:12:29', '2025-10-07 03:13:52'),
        (470, 'd15c3722-bc53-4457-82c2-83be6fa986cb', 'KKN FKIP', '$2y$10$ct8G8WdqsA6icA2A0Gz0O.bIeav1synwx2AcrSBZD7i1XEHa3A3/.', 'auditee', 'KKN (Unit FKIP)', 'dummy@gmail.com', 76, NULL, NULL, NULL, NULL, '2025-10-07 03:13:30', '2025-10-07 03:14:18'),
        (471, '7ffc1393-73a4-49cb-aca4-b53c975838a2', 'Bimbingan Konseling FKIP', '$2y$10$uLNPCLZw5YrPiiuAxUeGL.0SFLusOnWboq411VmZrc7305fSl/3pW', 'auditee', 'Bimbingan Konseling (Unit Fkip)', 'dummy@gmail.com', 87, NULL, NULL, NULL, 0, '2025-10-07 03:18:41', '2025-10-07 03:18:41'),
        (472, 'b818d342-de7e-430a-b1c9-c03809504942', 'Laboratorium Seni (FKIP)', '$2y$10$Z93xpkC/R2v4/7CPFKwdFOBegjbuVfN4GvsM6p9CkWLruABNwsAiO', 'auditee', 'Laboratorium Seni (Unit FKIP)', 'dummy@gmail.com', 128, NULL, NULL, NULL, NULL, '2025-10-07 04:05:52', '2025-10-11 02:53:19'),
        (473, '595eec54-2275-46c9-b08b-77188931f209', 'Laboratorium Bahasa Inggris (FKIP)', '$2y$10$JSmtsvgQqQEi218JqrcT8ezb8VPQLJ5WSVjKrnt.H4Crj2Ih6/NeW', 'auditee', 'Laboratorium Bahasa Inggris (Unit FKIP)', 'dummy@gmail.com', 129, NULL, NULL, NULL, 0, '2025-10-11 02:54:30', '2025-10-11 02:54:30'),
        (474, 'd2d5ab70-76dc-4639-9967-fd4513bb6072', 'Laboratorium Microteaching (FKIP)', '$2y$10$WBnKtk4qhQ8z0z0.35HwKOsAe6J3BPWFRJTKggQCwnd1d6Dj5oUGe', 'auditee', 'Laboratorium Microteaching (Unit FKIP)', 'dummy@gmail.com', 130, NULL, NULL, NULL, 0, '2025-10-11 02:55:21', '2025-10-11 02:55:21'),
        (475, '1414be4a-ec1a-4166-a25a-7504ecbcaef6', 'Laboratorium Komputer (FKIP)', '$2y$10$niw9zkI3DVJ9B1WrPd9PuupZjaZeGSikPlZGdzkiWgNHSHUqyZ5IS', 'auditee', 'Laboratorium Komputer (Unit FKIP)', 'dummy@gmail.com', 131, NULL, NULL, NULL, 0, '2025-10-11 02:56:48', '2025-10-11 02:56:48'),
        (476, '4cd95290-9bba-4170-995e-11c5c89dc436', 'AktStudio', '$2y$10$zjpiSGILR4dFkAkqS3bmDuPAq40exCd66xKTTchVip4ilvY/KyQe.', 'admin', 'Admin', 'adamilkom00@gmail.com', NULL, 'carbon (1).png', NULL, NULL, NULL, '2023-05-18 18:45:52', '2025-10-13 07:28:13'),
        (478, 'b79df916-d2f7-4630-8b6a-b4b09305473d', 'Verifikator FH', '$2y$10$om1jTBXX9ZXW9JOSVwzGBOpB4uYNaS9NhyXpZb.nPkgaKdFbFVPUO', 'fakultas', 'Verifikator FH', 'dummy@gmail.com', 91, NULL, NULL, NULL, NULL, '2025-10-15 01:07:12', '2025-10-15 01:08:54'),
        (479, '119489af-deaa-4ceb-a181-0656602ebbf3', 'Verifikator FEB', '$2y$10$FfRY9P4cO3YdfaLNH.Gp.OcIXkD7R25Iiv0H4azIfInRwMQoLf4P2', 'fakultas', 'Verifikator FEB', 'dummy@gmail.com', 94, NULL, NULL, NULL, 0, '2025-10-15 01:07:59', '2025-10-15 01:07:59'),
        (480, '67680076-ab9d-4a79-a3f3-065f477ff287', 'Verifikator FKIP', '$2y$10$prYvCgwMXmk1rZpfePGtm.ssqkqNt6.3xBV4WJcirI95Bpfc4IF0O', 'fakultas', 'Verifikator FKIP', 'dummy@gmail.com', 97, NULL, NULL, NULL, NULL, '2025-10-15 01:16:28', '2025-10-18 01:53:07'),
        (481, '127d7368-51c7-42e5-9aff-37b5c2faf3b6', 'Verifikator FISIB', '$2y$10$7kPzGY7QrUF5UVLaM0.MmOQ/unwDwnPvb.5/BjFxn63EY5miPeNzu', 'fakultas', 'Verifikator FISIB', 'dummy@gmail.com', 100, NULL, NULL, NULL, 0, '2025-10-15 01:27:55', '2025-10-15 01:27:55'),
        (482, '4bc35427-6503-4b89-8c31-f1d01abc0a07', 'Verifikator FT', '$2y$10$BMoP02xUWj32iEgluRaiZOTwe3L.CCfmWnEBNln0cqWQSTs.DR6Fi', 'fakultas', 'Verifikator FT', 'dummy@gmail.com', 103, NULL, NULL, NULL, 0, '2025-10-15 02:04:20', '2025-10-15 02:04:20'),
        (483, 'a3c1d18e-96df-40c5-a89f-8c89457382f0', 'Verifikator FMIPA', '$2y$10$Xv3wNbohilhZ14qRANgI.uj6UphgqYj4xoOVBwWGxkzBe7gXSG64u', 'fakultas', 'Verifikator FMIPA', 'dummy@gmail.com', 106, NULL, NULL, NULL, 0, '2025-10-15 02:04:42', '2025-10-15 02:04:42'),
        (484, '9c4a3581-1f62-490d-a4b2-670ad95f614e', 'Verifikator Pascasarjana', '$2y$10$OtQpv.Na9bPlmC4fv.9T0.ZoQ0MoadpwKVfQ.exHadWuVfpIbBy..', 'fakultas', 'Verifikator Pascasarjana', 'dummy@gmail.com', 109, NULL, NULL, NULL, NULL, '2025-10-15 02:06:04', '2025-10-15 02:06:22'),
        (485, 'f2661b35-0597-4aa1-9883-f355587c44a9', 'Verifikator Vokasi', '$2y$10$D.S5m6co.G8Qf.7z4UFATexnxvZKtNcxQ2i9Vf5RpsUa3wGnhLXLS', 'fakultas', 'Verifikator Vokasi', 'dummy@gmail.com', 112, NULL, NULL, NULL, 0, '2025-10-15 02:06:45', '2025-10-15 02:06:45'),
        (486, 'c734bce0-27e7-44fd-ba3a-7e81912ae442', 'Singgih Irianto', '$2y$10$og1dbBHJn4cSP9f0PdQpTugmEPBEfzDZYU9PRPCZnhy4oelLnx77K', 'auditor1', 'Dr. Singgih Irianto, MT', 'dummy@gmail.com', NULL, NULL, NULL, NULL, NULL, '2025-11-18 02:39:20', '2025-11-18 03:08:03'),
        (487, 'eeb323d3-464a-4d8f-9a3d-51d4940785f5', 'Singgih Irianto', '$2y$10$4hShGb.9/B0m09SayZq2hOXlQQ5s1V/BqdUw0l3f3HVIuVURHStNi', 'auditor2', 'Dr. Singgih Irianto, MT', 'dummy@gmail.com', NULL, NULL, NULL, NULL, 0, '2025-11-18 02:53:45', '2025-11-18 02:53:45'),
        (488, 'e9cab6ae-076d-45a9-a457-23a7cfc0beb9', 'Siti Warnasih', '$2y$10$y8YKA148yLICYrW7aljF3exTov/RhvrLjo/OPTiR4i9rgNXTBvwEK', '', 'Siti Warnasih, M.Si', 'siti.warnasih@unpak.ac.id', NULL, NULL, NULL, NULL, 1, '2025-11-18 02:56:59', '2025-11-28 08:16:33'),
        (489, '5b5fbc0d-6677-4e45-a3b0-72a94686e82b', 'Kotim Subandi', '$2y$10$Mp9gHKhzY4ZD8Rg63lWDI.3glqu.7JfUcvjFMjD3E0VymH7yS6Taa', 'auditor2', 'Kotim Subandi, S.Kom., M.T', 'dummy@gmail.com', NULL, NULL, NULL, NULL, 0, '2025-11-18 02:59:29', '2025-11-18 02:59:29'),
        (490, '644fd797-4cfb-4a14-9221-dc9cd4da081b', 'Oktori Kiswati Zaini', '$2y$10$6VbLUp.5pzxIlj7Lh9UYruy/dCnEOOse/FtWya3eZoFcww9tbRpzq', 'auditor2', 'Oktori Kiswati Zaini, S.E., M.M', 'dummy@gmail.com', NULL, NULL, NULL, NULL, 0, '2025-11-18 03:04:30', '2025-11-18 03:04:30'),
        (491, 'c19305f3-d47b-4df4-93d0-518b9a708d8e', 'Delta Hadi Purnama', '$2y$10$56ho7lpNkmeG30LYvUqJrOyUu0lUyUxzzTCjsA4Ui7Fdx3FC6AHXa', 'auditor2', 'Delta Hadi Purnama, S.E., M.E.Sy', 'dummy@gmail.com', NULL, NULL, NULL, NULL, 0, '2025-11-18 03:06:05', '2025-11-18 03:06:05'),
        (492, '30d01f90-f308-4784-9fc4-b90486156208', 'Mahifal', '$2y$10$ar.TiK3rWqZiaURdnlEKEeJxCgxd6vq9mp.6a60AcrMtfpQP/IiLm', 'auditor2', 'Dr. Mahifal, S.H., M.H', 'dummy@gmail.com', NULL, NULL, NULL, NULL, 0, '2025-11-18 03:12:21', '2025-11-18 03:12:21'),
        (493, 'a75e962c-57d4-4706-8102-c3fdf655792b', 'Atti Herawati', '$2y$10$6Cynt4Y.62D5dIvRSezL1edR3YykT8bPqrMUng2v/lbjERHMGykfO', 'auditor2', 'Dr. Atti Herawati, M.Pd', 'dummy@gmail.com', NULL, NULL, NULL, NULL, 0, '2025-11-18 03:13:55', '2025-11-18 03:13:55'),
        (494, '44197945-dc0c-4688-b0e0-bda5aaa46e7e', 'Anwar Sulaiman', '$2y$10$yGj/ws9BnataxAOZLphVTuTJ/SX2wmjitKxjOKWMV.EJ3ndRlwjzG', 'auditor2', 'Dr. Anwar Sulaiman, S.E., M.M. S.H', 'dummy@gmail.com', NULL, NULL, NULL, NULL, 0, '2025-11-18 03:17:24', '2025-11-18 03:17:24'),
        (495, 'e0d65c1b-10db-40d6-a716-8dbdc28a3bd6', 'Desti Herawati', '$2y$10$P6u5YQJ9UaOo0DDRPMdme.TgmHwP/po4NS3oFfQSEKie5a6pPdDwS', 'auditor2', 'Desti Herawati, M.Pd', 'dummy@gmail.com', NULL, NULL, NULL, NULL, 0, '2025-11-18 03:18:18', '2025-11-18 03:18:18'),
        (496, 'e317e9c1-1f75-4ab4-adb3-f12157535ad8', 'Diana Amaliasari', '$2y$10$63M4MS2cG377YotjbdcFu.Kc339XTYLvZuJBDbsFV6jGj0a8s0Wie', 'auditor2', 'Diana Amaliasari, M.Si', 'dummy@gmail.com', NULL, NULL, NULL, NULL, 0, '2025-11-18 03:20:13', '2025-11-18 03:20:13'),
        (497, '207f3eae-2a75-4176-90fb-6192dfdc3343', 'Siti Maimunah', '$2y$10$ULGXj8hRcyQx91N6oJZXbe5F2pRCV..we3hxUUCWM5FeOQNtVbZZ2', 'auditor2', 'Dr. Siti Maimunah, S.E., M.Si', 'dummy@gmail.com', NULL, NULL, NULL, NULL, 0, '2025-11-18 03:22:57', '2025-11-18 03:22:57'),
        (498, 'abd8913e-ab4f-4898-9574-c4c579d6cb5c', 'Abdul Kohar', '$2y$10$DcPpqO6SFkB7xmIJ0DWtLO3bJifXhgxmrtXfrFaTD/544My0wFugy', 'auditor2', 'Dr. Abdul Kohar, S.E., M.Ak', 'dummy@gmail.com', NULL, NULL, NULL, NULL, 0, '2025-11-18 03:24:00', '2025-11-18 03:24:00'),
        (499, '5b0f9fa0-e8be-43e0-9d70-f1edb4a016af', 'Henny Suharyati', '$2y$10$QVWvoz9zHb8OMkcAE2U6S.kFZCHpJO04tnNRb97v8zRWa4J4CP2HS', 'auditor1', 'Prof. Henny Suharyati, M.Pd', 'dummy@gmail.com', NULL, NULL, NULL, NULL, 0, '2025-11-18 03:27:10', '2025-11-18 03:27:10'),
        (500, '20a96cb7-9778-4add-8ca5-598c33856894', 'Dewi Taurusyanti', '$2y$10$v3GajUxO109oUgxEmNfN4eFbaTwE5wPfHOJUmrSlTO1gDRe.1xINm', 'auditor2', 'Dr.  Dewi Taurusyanti, S.E., M.M', 'dummy@gmail.com', NULL, NULL, NULL, NULL, 0, '2025-11-18 03:28:26', '2025-11-18 03:28:26'),
        (501, '4f6fa3ff-dd40-4b97-850e-cc12765a7dde', 'Yenny Febrianty', '$2y$10$zcmAyyqEOzxwmZeTwXk6tuEGXiC8BFoQQtTopOM90pdfRDZGj47L6', 'auditor2', 'Dr. Yenny Febrianty, S.H., M.Hum., M.Kn', 'dummy@gmail.com', NULL, NULL, NULL, NULL, 0, '2025-11-18 03:32:45', '2025-11-18 03:32:45'),
        (502, '4f9cd28a-a248-47df-9bea-218bc98c72ed', 'Hasrul', '$2y$10$52TbWvyZ0AMNlGUJRx1Nqe1p9SIPNgOvRUd2VvonqGSQuRCjQU.uG', 'auditor2', 'Dr. Hasrul, MM', 'dummy@gmail.com', NULL, NULL, NULL, NULL, 0, '2025-11-18 03:33:31', '2025-11-18 03:33:31'),
        (503, '397ad431-eb20-44f5-b44c-57b4ae56edf9', 'Elly Sukmanasa', '$2y$10$67Q51YyGYL3v4MCf0pB.QOfB8NnlNwMXfBzl7KQ0hvGvSdalTXAEu', 'auditor2', 'Dr. Elly Sukmanasa, M .Pd', 'dummy@gmail.com', NULL, NULL, NULL, NULL, 0, '2025-11-18 03:37:11', '2025-11-18 03:37:11'),
        (504, '099c4f09-a89c-42ce-aa3d-e366b43399f7', 'Yunita Rahma', '$2y$10$ocjML4tg1woR/rcNzhuBEOMyxiLwzRZT7c3clNwZwHVwbYeKTECv6', 'auditor2', 'Yunita Rahma, M.Kom', 'dummy@gmail.com', NULL, NULL, NULL, NULL, 0, '2025-11-18 03:39:50', '2025-11-18 03:39:50'),
        (505, '7911f1eb-2305-4766-9355-2e9e2bf2965e', 'Muhammad Reza', '$2y$10$rq6fwxWrjZOG3ncTOllOp.xr953gDUjzyS3nn5KiEoOzt9/CNIRSu', 'auditor2', 'Muhammad Reza, M.Si', 'dummy@gmail.com', NULL, NULL, NULL, NULL, 0, '2025-11-18 03:40:33', '2025-11-18 03:40:33'),
        (506, '761c5f96-4dab-4aee-a455-b67215449009', 'Evyta Wismiana', '$2y$10$X6bO.FXvxbHkmVFdGf0uCOcw6LsDsLYToIKiDwG7b6HoDwagis1tm', 'auditor2', 'Evyta Wismiana, S.T., M.T', 'dummy@gmail.com', NULL, NULL, NULL, NULL, 0, '2025-11-18 03:42:35', '2025-11-18 03:42:35'),
        (507, '69058276-ebd4-44be-9cee-3d3e53d975ba', 'Restiawan Permana', '$2y$10$gcPXVufgqxrSfp.ECT4Npu6.ideN4MX9k8pgfLFlXinJ4wIkFOm5a', 'auditor2', 'Restiawan Permana, M.Si', 'dummy@gmail.com', NULL, NULL, NULL, NULL, NULL, '2025-11-18 03:44:36', '2025-11-18 03:44:53'),
        (508, 'd29bbdcb-252a-4adb-8f07-dc98ac8a24d3', 'Helmi Setia Pamungkas', '$2y$10$qDsFf7n1sdqDrL8RRYzDEeu6CHMmrzXYUPie26/fIA1s.fY3OgwaC', 'auditor2', 'Helmi Setia Pamungkas, S.T., M.Si', 'dummy@gmail.com', NULL, NULL, NULL, NULL, 0, '2025-11-18 04:09:45', '2025-11-18 04:09:45'),
        (509, 'c65eed10-3b49-42da-8ba7-d4230f5dfcd2', 'Lia Amelia Megawati', '$2y$10$X9sk8HKDTu86lfIqUzU2DuZO54LECqpCV4HLv8/esbwaIO.4kYTtu', 'auditor2', 'Lia Amelia Megawati, S.pd., M.T', 'dummy@gmail.com', NULL, NULL, NULL, NULL, 0, '2025-11-18 04:10:28', '2025-11-18 04:10:28'),
        (510, '40a79496-b215-4380-924a-056f08439752', 'Yan Noviar Nasution', '$2y$10$WWcNPMusu9qjamfpWdnoPekx0MpHbMXd1FXd07fda72xWOtrZaaJC', 'auditor1', 'Dr. Yan Noviar Nasution, S.E., M.M., CA', 'dummy@gmail.com', NULL, NULL, NULL, NULL, 0, '2025-11-18 04:11:45', '2025-11-18 04:11:45'),
        (511, 'acfcdb84-aabd-4ec7-98fb-0417e94a911b', 'Komarudin', '$2y$10$Dw653hB21hTX4N8M.050c.39b5AOxb.ZmJMLyRCxyUnTOUG5uXcAW', 'auditor2', 'Komarudin, M.H', 'dummy@gmail.com', NULL, NULL, NULL, NULL, 0, '2025-11-18 04:12:43', '2025-11-18 04:12:43'),
        (512, '270e5ff9-ffd2-453f-a971-8931b67935e2', 'Dwi Rini Sovia Firdaus', '$2y$10$q4kKllXwk6z1ZBn4zmVT3uX29/rK3af.I4qyr1eGZBgaU2IL1rpPW', 'auditor2', 'Dr. Dwi Rini Sovia Firdaus, M. Comn', 'dummy@gmail.com', NULL, NULL, NULL, NULL, 0, '2025-11-18 04:13:41', '2025-11-18 04:13:41'),
        (513, 'ff48688f-dda3-4e11-ae79-d5ddef31a260', 'Lusi Agus Setiani', '$2y$10$.RaYo/aSF/YIwKrr/gpVLOy2I6YOlz8mIo/RwKZiyGgYSbuS6shB2', 'auditor2', 'Dr. apt. Lusi Agus Setiani, M.Farm', 'dummy@gmail.com', NULL, NULL, NULL, NULL, 0, '2025-11-18 04:14:31', '2025-11-18 04:14:31'),
        (514, '061a3606-5ae9-4de1-ac75-8c306d1fe75e', 'May Mulyaningsih', '$2y$10$QsJBJO72hiuH1RSRh/xLn.93krIkX/e4tEEpZBauWH.5BfbTMr9dO', 'auditor2', 'May Mulyaningsih, S.E., M.Ak', 'dummy@gmail.com', NULL, NULL, NULL, NULL, 0, '2025-11-18 04:17:18', '2025-11-18 04:17:18'),
        (515, '0e340ca2-eb94-4880-a33b-0eac6fc86b37', 'Sata Yoshida Srie Rahayu', '$2y$10$BtXCae0raD2waLa987JayeYxCJRm6I0Lq95vARNaWiJhka/C/.9R2', 'auditor2', 'Prof. Dr. Sata Yoshida Srie Rahayu, M.Si', 'dummy@gmail.com', NULL, NULL, NULL, NULL, 0, '2025-11-18 04:18:12', '2025-11-18 04:18:12'),
        (516, 'f03ba8b3-918d-4fd2-867d-6943dc14a5ac', 'beta', '123', '', 'beta', 'a@unpak.ac.id', 6, NULL, NULL, NULL, 0, NULL, NULL);

        INSERT INTO m_fakultas (kode_fakultas, kode_pt, nama_fakultas, pejabat, jabatan, wakil_pejabat, wakil_pejabat_adm, logo) VALUES
        ('01', '041004', 'HUKUM', '0410067306', 'H', '0417086801', '0414106202', ''),
        ('02', '041004', 'EKONOMI DAN BISNIS', '0415056901', 'H', '0425097604', '0414109101', ''),
        ('03', '041004', 'KIP', '0416076701', 'H', '0413018604', '0415128202', ''),
        ('04', '041004', 'ISIB', '0416068002', 'H', '0401098708', '0413026701', ''),
        ('05', '041004', 'TEKNIK', '0428106901', 'H', '0424076601', '0417097601', ''),
        ('06', '041004', 'MIPA', '0406097101', 'H', '0404107405', '0407027501', ''),
        ('07', '041004', 'PASCASARJANA', '0403055801', 'H', '0427017501', '0025057515', ''),
        ('08', '041004', 'VOKASI', '0413117601', 'H', '0402047301', '0404047007', '');

        INSERT INTO m_program_studi (kode_prodi, kode_pt, kode_fak, kode_jenjang, kode_jurusan, nama_prodi, alamat, kode_kabupaten, kode_propinsi, kode_negara, kode_pos, telepon, fax, email, website, sks_lulus, status_prodi, tgl_awal_berdiri, semester_awal, mulai_semester, no_sk_dikti, tgl_sk_dikti, tgl_akhir_sk_dikti, no_sk_ban, tgl_sk_ban, tgl_akhir_sk_ban, kode_akreditasi, frekuensi_kurikulum, pelaksanaan_kurikulum, idd_ketua_prodi, hp_ketua, idd_nama_operator, telepon_operator, nama_sesi, jumlah_sesi, batas_sesi, gelar, gelar_panjang, no_sk_ban_lama, logo, nama_prodi_ing) VALUES
        ('20201', '041004', '05', 'C', '05', 'TEKNIK ELEKTRO', ' Jl. Pakuan PO Box 452 Bogor', 200, 2, NULL, '16143', '081382928035', '0521 8314332', 'fteknik@unpak.ac.id', 'http://ft.unpak.ac.id', 144, 'A', '1986-05-22', '20021', '20021', '8414/D/T/K-IV/2011', '2011-08-19', '2015-09-23', '1151/SK/BAN-PT/Akred/S/XI/2015', '2015-11-14', '2020-11-14', 'B', 'F', 'A', '0428106301', '081316626393', 'EVYTA WISMIANA', '0521 8314332', 'SEMESTER', 2, 8, 'S.T.', 'Sarjana Teknik', '-', '', NULL),
        ('22201', '041004', '05', 'C', '05', 'TEKNIK SIPIL', 'Jl. Pakuan', 200, 2, NULL, '16154', '0251314136', '0251314136', 'ft@unpak.ac.id', 'http://ft.unpak.ac.id', 100, 'A', '1985-12-30', '20021', '20021', '8415/D/T/K-IV/2011', '2011-08-19', '2015-09-23', '247/SK/BAN-PT/Ak-PPJ/S/I/2021', '2020-12-30', '2025-12-30', 'B', 'F', 'A', '0405048403', '', '', '', 'SEMESTER', 2, 8, 'S.T.', 'Sarjana Teknik', '-', '', '-'),
        ('29201', '041004', '05', 'C', '05', 'TEKNIK GEODESI', 'Jl. Pakuan PO Box 452  Bogor 16143 Jawa Barat Indonesia', 200, 2, NULL, '16150', '08174990457', '02518311007', 'geo_unpak@yahoo.com.', 'http://ft.unpak.ac.id', 4, 'A', '1980-10-28', '20021', '20021', '8417/D/TK-IV/2011', '2011-08-19', '2015-09-23', '2275/SK/BAN-PT/Akred/S/VII/2019', '2019-07-09', '2024-07-09', 'B', 'F', 'A', '0415037606', '', '', '', 'SEMESTER', 2, 8, 'S.T.', 'Sarjana Teknik', '0', '', '-'),
        ('34201', '041004', '05', 'C', '05', 'TEKNIK GEOLOGI', 'Jl. Ciheulet Pakuan No.1', 200, 2, NULL, '16143', '081931102259', '02518311007', 'geologi@unpak.ac.id', 'http://www.unpak.ac.id', 4, 'A', '1985-12-30', '20021', '20021', '8417/D/TK-IV/2011', '2011-08-19', '2015-09-23', '5381/SK/BAN-PT/Ak-PPJ/S/IX/2020', '2020-09-04', '2025-09-04', 'B', 'D', 'A', '0411098303', '', '', '', 'SEMESTER', 2, 8, 'S.T.', 'Sarjana Teknik', '0', '', '-'),
        ('35201', '041004', '05', 'C', '05', 'PERENCANAAN WILAYAH DAN KOTA', 'Jl. Pakuan', 200, 2, NULL, '16143', '08567220884', '02518311007', 'ft@unpak.ac.id', 'http://ft.unpak.ac.id', 130, 'A', '1980-10-28', '20021', '20021', '8417/D/TK-IV/2011', '2011-08-19', '2015-09-23', '10328/SK/BAN-PT/Akred/M/VIII/2021', '2021-08-24', '2026-08-24', 'B', 'F', 'A', '0414057303', '08128061453', 'SUPRIA', '02518311007', 'SEMESTER', 2, 8, 'S.P.W.K.', 'Sarjana Perencanaan Wilayah dan Kota', '-', '', '-'),
        ('44201', '041004', '06', 'C', '06', 'MATEMATIKA', 'Jl. Pakuan PO Box 452  Bogor 16143 Jawa Barat Indonesia\r\n', 200, 2, NULL, '16143', '02517122609', '0251375547', 'matematika@unpak.ac.id', 'http://math.unpak.ac.id', 90, 'A', '1983-12-22', '20021', '20021', '8064/D/T/K-IV/2011', '2011-07-31', '2015-09-23', '440/SK/BAN-PT/Ak-PPJ/S/I/2021', '2021-01-12', '2026-01-12', 'B', 'D', 'E', '0404017801', '08128104148', 'ANI ANDRIYATI', '0251375547', 'SEMESTER', 2, 8, 'S.Mat.', 'Sarjana Matematika', '0', '', '-'),
        ('46201', '041004', '06', 'C', '06', 'BIOLOGI', 'Jl. Pakuan No.1', 200, 2, NULL, '16151', '081311082832', '02518375547', 'fmipa@unpak.ac.id', 'http://www.unpak.ac.id', 100, 'A', '1981-11-08', '20021', '20021', '8062/D/T/K-IV/2011', '2011-07-31', '2015-09-23', '4782/SK/BAN-PT/Akred/S/VIII/2020', '2020-08-19', '2025-08-19', 'B', 'D', 'E', '0422016902', '081287794946', 'DRA. TRIASTINURMIATININGSIH,M.SI', '02518375547', 'SEMESTER', 2, 8, 'S.Si.', 'Sarjana Sains', '-', '', NULL),
        ('47201', '041004', '06', 'C', '06', 'KIMIA', 'Jl. Pakuan PO Box 452  Bogor 16143 Jawa Barat Indonesia\r\n', 200, 2, NULL, '16143', '08121835020', '0251-8375547', 'kimia@unpak.ac.id', 'http://kimia.fmipa.unpak.ac.id', 100, 'A', '1982-12-04', '20021', '20021', '8065/D/T/K-IV/2011', '2011-07-31', '2015-10-27', '2967/SK/BAN-PT/Akred/S/VIII/2019', '2019-08-07', '2024-08-07', 'A', 'D', 'B', '0408129001', '', '', '', 'SEMESTER', 2, 8, 'S.Si.', 'Sarjana Sains', '0', '', '-'),
        ('48201', '041004', '06', 'C', '06', 'FARMASI', 'Jl. Pakuan PO Box 452  Bogor 16143 Jawa Barat Indonesia', 200, 2, NULL, '16143', '02519162290', '0251-8375547', 'farmasi_unpak@yahoo.com', 'http://farmasi.fmipa.unpak.ac.id', 24, 'A', '2001-05-07', '20021', '20021', '10386/D/T/K-IV/2012', '2012-02-07', '2016-05-26', '0470/LAM-PTKes/Akr/Sar/VII/2018', '2018-07-07', '2023-07-07', 'B', 'D', 'A', '0425109101', '08129918454', 'SRI WARDATUN', '0251-8375547', 'SEMESTER', 2, 8, 'S.Farm.', 'Sarjana Farmasi', '0', '', '-'),
        ('48901', '041004', '06', 'J', '06', 'PROFESI APOTEKER', 'Jl. Pakuan Bogor - Jawa Barat', 200, 2, NULL, '16129', '0251-8312206', '0251-8356927', 'rektorat@unpak.ac.id', 'https://fmipa.unpak.ac.id/', 36, 'A', '2025-01-21', '20242', '20242', '-', '2025-01-21', '2029-12-12', '-', '2025-01-21', '2029-12-12', 'U', 'F', 'A', '0403086301', '', '', '', 'Semester', 2, 4, 'apt', 'Apoteker', '-', '', 'Pharmacy Profession'),
        ('55101', '041004', '07', 'B', '07', 'ILMU KOMPUTER', 'Jalan Pakuan No.1 Ciheuleut , Kec. Kota Bogor Tengah, Kota Bogor, Prov. Jawa Barat', 200, 2, NULL, '16129', '0251-8312206', '0251-8356927', 'pasca@unpak.ac.id', 'https://pasca.unpak.ac.id', 39, 'A', '2024-12-12', '20251', '20251', '57/A/O/2024', '2024-12-12', '2029-12-12', '57/A/O/2024', '2024-12-12', '2029-12-12', 'B', 'F', 'B', '0425087601', '', '', '', 'Semester', 2, 4, 'M.Kom', 'Magister Komputer', '57/A/O/2024', '', 'Magister Computer Science'),
        ('55201', '041004', '06', 'C', '06', 'ILMU KOMPUTER', 'Jl. Pakuan', 200, 2, NULL, '16154', '081318103370', '0251-8375547', 'ilkom@unpak.ac.id', 'http://www.unpak.ac.id', 17, 'A', '2005-06-24', '20021', '20021', '10994/D/T/K-IV/2012', '2012-03-12', '2015-08-24', '13704/SK/BAN-PT/Ak-PPJ/S/XII/2021', '2021-12-28', '2026-12-28', 'B', 'D', 'A', '0425038403', '', '', '', 'SEMESTER', 2, 8, 'S.Kom.', 'Sarjana Komputer', '-', '', '-'),
        ('56209', '041004', '08', 'D', '08', 'KECERDASAN BUATAN DAN ROBOTIK', 'Jl. Pakuan No.1 Kota Bogor', 200, 2, NULL, '16151', '0821', '0', 'vokasi@unpak.ac.id', 'https://vokasi.unpak.ac.id', 144, 'A', '2024-03-25', '20241', '20241', 'Nomor 106/D/O/2024', '2024-03-25', '2050-03-25', '-', '2024-03-25', '2024-03-25', 'B', 'F', 'A', '0422089401', '', '', '', 'SEMESTER', 8, 4, 'S.Tr', 'Sarjana Terapan', '-', '', 'Artificial Intelligence and Robotics'),
        ('56401', '041004', '08', 'E', '08', 'TEKNIK KOMPUTER', 'J. Pakuan', 200, 2, NULL, '16143', '087870663038', '0251-375547', 'sekretariatd3tk@unpak.ac.id', 'http://www.unpak.ac.id', 110, 'A', '2007-08-02', '20081', '20081', '8412/D/T/K-IV/2011', '2011-08-19', '2013-08-19', '4184/SK/BAN-PT/Akred/Dipl-III/XI/2017', '2017-11-07', '2022-11-07', 'B', 'C', 'A', '0418098906', '', '', '', 'SEMESTER', 2, 6, 'A.Md.T.', 'Ahli Madya Teknik Komputer', '-', '', '-'),
        ('57401', '041004', '08', 'E', '08', 'MANAJEMEN INFORMATIKA', 'Jl. Pakuan', 200, 2, NULL, '16143', '085719377063', '02518375547', 'sekretariatd3mi@unpak.ac.id', 'http://www.unpak.ac.id', 100, 'A', '2007-08-02', '20071', '20071', '3020/D/T/K-IV/2010', '2010-07-14', '2012-07-14', '13824/SK/BAN-PT/Ak-PPJ/Dipl-III/XII/2021', '2021-12-28', '2026-12-28', 'B', 'C', 'A', '0402108604', '08569093787', 'JAMILUDIN', '02518375547', 'SEMESTER', 2, 6, 'A.Md.Kom.', 'Ahli Madya Komputer', '-', '', '-'),
        ('61001', '041004', '07', 'A', '07', 'ILMU MANAJEMEN', 'Jl. Pakuan No. 1', 200, 2, NULL, '16143', '0251-8312206', '0251-8356927', 's3manajemen@unpak.ac.id', 'http://pasca.unpak.ac.id', 39, 'A', '2017-05-31', '20171', '20171', '287/KPT/I/2017', '2017-05-31', '2022-05-31', '2823/SK/BAN-PT/Akred/D/VII/2019', '2019-07-31', '2024-07-31', 'C', 'A', 'A', '0411035501', '', 'LANA FADILAH', '02518320123', 'SEMESTER', 2, 6, 'Dr.', 'Doktor', '0', '', 'management science'),
        ('61101', '041004', '07', 'B', '07', 'MANAJEMEN', 'Jl. Pakuan', 200, 2, NULL, '16154', '08111106652', '02518320123', 's2manajemen@unpak.ac.id', 'http://pasca.unpak.ac.id', 6, 'A', '2006-09-28', '20071', '20071', '3019/D/T/K-IV/2010', '2010-07-14', '2013-07-14', '2154/SK/BAN-PT/Akred/M/VI/2017', '2017-06-20', '2022-06-20', 'A', 'A', 'A', '0319087004', '', '', '', 'SEMESTER', 2, 4, 'M.M.', 'Magister Manajemen', '-', '', '-'),
        ('61105', '041004', '07', 'B', '07', 'MANAJEMEN LINGKUNGAN', 'Jl. Pakuan', 200, 2, NULL, '16143', '085218534206', '085218534206', 'pasca@unpak.ac.id', 'http://pasca.unpak.ac.id', 6, 'A', '2017-12-05', '20172', '20172', '677/KPT/I/2017', '2017-12-05', '2017-12-05', '5001/SK/BAN-PT/Ak-PNB/M/IX/2020', '2020-07-04', '2022-07-04', 'A', 'D', 'A', '0405077105', '081281639705', 'Fredy Herlambang', '085697770081', 'SEMESTER', 2, 8, 'M.Ling.', 'Magister Lingkungan', '0', '', NULL),
        ('61201', '041004', '02', 'C', '02', 'MANAJEMEN', 'Bogor', 200, 2, NULL, '16143', '085710581117', '0251-314918 (109)', 'fekonomi@unpak.ac.id', 'http://fekonomi.unpak.ac.id', 6, 'A', '1982-05-10', '20021', '20021', '2299/D/T/2007', '2020-08-25', '2025-08-25', '4869/SK/BAN-PT/Akred/S/VIII/2020', '2020-08-25', '2025-08-25', 'A', 'D', 'B', '8878801019', '', 'AHMAD SADIKIN', '0251-314918 (119)', 'SEMESTER', 2, 8, 'S.M.', 'Sarjana Manajemen', '-', '', NULL),
        ('61209', '041004', '02', 'C', '02', 'BISNIS DIGITAL', 'Jl. Pakuan No.1', 200, 2, NULL, '16143', '02518312206', '02518312206', 'bisnisdigital@unpak.ac.id', 'https://bisnisdigital-feb.unpak.ac.id', 144, 'A', '2021-10-04', '20211', '20211', '-', '2021-10-04', '2026-10-04', '-', '2021-10-04', '2026-10-04', 'C', 'F', 'B', '0409128303', '081282821076', 'Sintia Andriani', '08569800750', 'SEMESTER', 8, 14, 'S.Bs', 'Sarjana Bisnis', '-', '', 'Digital Business'),
        ('61403', '041004', '08', 'E', '08', 'MANAJEMEN PERPAJAKAN', 'Jl. Pakuan No.1', 200, 2, NULL, '16143', '02518314918', '0251-8314918 (109)', 'fekonomi@unpak.ac.id', 'http://fekonomi.unpak.ac.id', 6, 'A', '2000-04-05', '20021', '20021', '10385/D/T/K-IV/2012', '2012-02-06', '2015-02-06', '12265/SK/BAN-PT/Ak-PPJ/Dipl-III/XI/2021	', '2021-10-28', '2025-10-28', 'B', 'D', 'B', '0411046605', '', 'JEFRI HARIMUDIANTO', '0251-8314918', 'SEMESTER', 2, 6, 'A.Md.M.', 'Ahli Madya Manajemen Pajak', '-', '', '-'),
        ('61406', '041004', '08', 'E', '08', 'KEUANGAN DAN PERBANKAN', 'Jl. Pakuan', 200, 2, NULL, '16143', '081111', '02518314918', 'fekonomi@unpak.ac.id', 'http://fekonomi.unpak.ac.id', 6, 'A', '2006-08-23', '20071', '20071', '/2012', '2012-08-23', '2017-12-06', '12263/SK/BAN-PT/Ak-PPJ/Dipl-III/XI/2021', '2021-10-28', '2026-10-28', 'B', 'B', 'B', '000002055', '', 'JEFRI', '02518314918', 'SEMESTER', 2, 6, 'A.Md.M.', 'Ahli Madya Perbankan dan Keuangan', '-', '', '-'),
        ('62201', '041004', '02', 'C', '02', 'AKUNTANSI', 'Bogor', 200, 2, NULL, '16143', '0251 314918', '0251 314918 EXT 109', 'fekonomi@unpak.ac.id', 'http://fekonomi.unpak.ac.id', 6, 'A', '1982-05-10', '20021', '20021', '10990/D/T/K-IV/2012', '2012-03-12', '2015-09-23', '10988/SK/BAN-PT/Ak-PPJ/S/IX/2021', '2021-09-17', '2026-09-17', 'A', 'D', 'B', '0408127102', '02518314918', 'AHMAD SUCHIA', '0251 314918', 'SEMESTER', 2, 8, 'S.Ak.', 'Sarjana Akuntansi', '-', '', '-'),
        ('62401', '041004', '08', 'E', '08', 'AKUNTANSI', 'Jl. Pakuan', 200, 2, NULL, '16143', '02518341918', '0251-320373 (109)', 'fekonomi@unpak.ac.id', 'http://fekonomi.unpak.ac.id', 6, 'A', '2000-04-05', '20021', '20021', '10384/D/T/K-IV/2011', '2011-02-07', '2015-02-06', '10996/SK/BAN-PT/Ak-PPJ/Dipl-III/IX/2021', '2021-09-17', '2026-09-17', 'B', 'D', 'B', '0413117601', '081311128757', 'JEFRI HARIMUDIANTO', '0251-320373', 'SEMESTER', 2, 6, 'A.Md.Ak.', 'Ahli Madya Akuntansi', '-', '', '-'),
        ('70201', '041004', '04', 'C', '04', 'ILMU KOMUNIKASI', 'Jl. Pakuan', 200, 2, NULL, '16245', '011', '02518338650', 'fisib@unpak.ac.id', 'http://fisib.unpak.ac.id', 50, 'A', '2008-08-20', '20081', '20081', '11273/D/T/K-IV/2012', '2012-03-29', '2017-03-29', '150/SK/BAN-PT/Akred/S/I/2018', '2018-01-03', '2023-01-03', 'B', 'A', 'A', '0411037204', '', '', '', 'SEMESTER', 2, 8, 'S.I.Kom.', 'Sarjana Ilmu Komunikasi', '-', '', 'communication science'),
        ('74101', '041004', '07', 'B', '07', 'ILMU HUKUM', 'Jl. Pakuan No. 1', 200, 2, NULL, '16143', '085811211448', '(0251)320123', 'pasca@unpak.ac.id', 'http://pasca.unpak.ac.id', 6, 'A', '2004-08-04', '20041', '20091', '10955/D/T/K-IV/2012', '2012-03-12', '2014-07-02', '2116/SK/BAN-PT/Ak-PPJ/M/IV/2020', '2020-04-01', '2025-04-01', 'B', 'A', 'A', '0408076801', '08129884477', 'UMPS ILMU HUKUM', '(0251)8320123', 'SEMESTER', 2, 4, 'M.H.', 'Magister Hukum', '-', '', '-'),
        ('74201', '041004', '01', 'C', '01', 'ILMU HUKUM', 'Jl. Pakuan PO Box 452  Bogor 16143 Jawa Barat Indonesia', 200, 2, NULL, '16143', '081317567126', '0251373588', 'fhukum@unpak.ac.id', 'http://fhukum.unpak.ac.id', 100, 'A', '1982-05-10', '20021', '20021', '2301/D/T/2011', '2011-08-21', '2015-08-21', '2085/SK/BAN-PT/Ak-PPJ/S/IV/2020', '2020-04-01', '2025-04-01', 'A', 'B', 'A', '0408017802', '085219213405', 'MUCHYAR KAMALDI', '0251373588', 'SEMESTER', 2, 8, 'S.H.', 'Sarjana Hukum', '0', '', 'Law'),
        ('79201', '041004', '04', 'C', '04', 'SASTRA INDONESIA', 'Jl. Pakuan', 200, 2, NULL, '16451', '0251-338650', '0251-338650', 'fisib@unpak.ac.id', 'http://fisib.unpak.ac.id', 135, 'A', '2000-10-16', '20021', '20021', '10992/D/T/K-IV/2012', '2012-03-12', '2017-03-12', '0431/SK/BAN-PT/Akred/S/I/2017', '2017-01-26', '2022-01-26', 'B', 'B', 'A', '0423066701', '', '', '', 'SEMESTER', 2, 8, 'S.S.', 'Sarjana Sastra', '-', '', 'Indonesian literature'),
        ('79202', '041004', '04', 'C', '04', 'SASTRA INGGRIS', 'Jl. Pakuan', 200, 2, NULL, '16143', '0251-338650', '0251-338650', 'fisib@unpak.ac.id', 'http://fisib.unpak.ac.id', 135, 'A', '1982-05-10', '20021', '20021', '8959/D/T/K-IV/2011', '2011-09-30', '2016-09-30', '12403/SK/BAN-PT/Ak-PPJ/S/XI/2021', '2021-10-07', '2026-10-07', 'B', 'B', 'A', '0428118505', '081294909678', 'IMAN SUDIRMAN', '0251-338650', 'SEMESTER', 2, 8, 'S.S.', 'Sarjana Sastra', '-', '', '-'),
        ('79204', '041004', '04', 'C', '04', 'SASTRA JEPANG', 'Jl. Pakuan', 200, 2, NULL, '16452', '0251-338650', '0251-338650', 'fisib@unpak.ac.id', 'http://fisib.unpak.ac.id', 135, 'A', '1986-01-15', '20021', '20021', '8418/D/T/K-IV/2011', '2011-08-19', '2016-08-19', '13705/SK/BAN-PT/Ak-PPJ/S/XII/2021', '2021-12-28', '2026-12-28', 'B', 'A', 'A', '0313076603', '08129173030', 'IMAN SUDIRMAN', '0251-338650', 'SEMESTER', 2, 8, 'S.S.', 'Sarjana Sastra', '3317/SK/BAN-PT/Akred/S/XII/2016', '', '-'),
        ('84101', '041004', '07', 'B', '07', 'PENDIDIKAN IPA', 'Jl. Pakuan No. 1, Ciheuleut.', 200, 2, NULL, '16143', '0251 8320 123', '0251 8320 123', 'pasca@unpak.ac.id', 'http://pasca.unpak.ac.id', 39, 'A', '2013-12-04', '20141', '20141', 'SK', '2013-12-02', '2019-05-02', '4331/SK/BAN-PT/Akred/M/VII/2020', '2020-07-29', '2025-07-29', 'A', 'B', 'A', '0405028902', '', '', '', 'SEMESTER', 2, 4, 'M.Pd.', 'Magister Pendidikan', 'SK 11111', '', 'MASTER OF SCIENCE EDUCATION'),
        ('84205', '041004', '03', 'C', '03', 'PENDIDIKAN BIOLOGI', 'Jl. Pakuan', 200, 2, NULL, '16143', '081211111111', '0251-8375608', 'fkip@unpak.ac.id', 'http://fkip.unpak.ac.id', 6, 'A', '1982-05-10', '20021', '20021', '2297/D/T/2011', '2011-08-21', '2015-08-21', '4703/SK/BAN-PT/Akred/S/VIII/2020', '2020-08-18', '2025-08-18', 'B', 'D', 'A', '0430058702', '', '', '', 'SEMESTER', 2, 8, 'S.Pd.', 'Sarjana Pendidikan', '0', '', 'biology education'),
        ('84206', '041004', '03', 'C', '03', 'PENDIDIKAN IPA', 'Jl. Pakuan', 200, 2, NULL, '16143', '081211111111', '081211111111', 'fkip@unpak.ac.id', 'http://fkip.unpak.ac.id', 6, 'A', '2019-01-07', '20191', '20191', '2297/D/T/2019', '2019-01-06', '2030-01-06', '9452/SK/BAN-PT/Akred/S/VII/2021', '2021-07-13', '2021-07-13', 'B', 'D', 'A', '0409018403', '085693810037', 'HILDA', '0821', 'SEMESTER', 2, 8, 'S.Pd.', 'Sarjana Pendidikan', '0', '', '-'),
        ('86004', '041004', '07', 'A', '07', 'MANAJEMEN PENDIDIKAN', 'Jl. Pakuan', 200, 2, NULL, '16154', '1111', '02518320123', 'pasca@unpak.ac.id', 'http://pasca.unpak.ac.id', 35, 'A', '2009-02-09', '20091', '20091', '11274/D/T/K-IV/2012', '2012-03-29', '2014-02-09', '2562/SK/BAN-PT/Akred/D/VIII/2017', '2017-08-01', '2022-08-01', 'B', 'A', 'A', '0428047402', '', 'LANA FADILAH', '02518320123', 'SEMESTER', 2, 6, 'Dr.', 'Doktor', '-', '', ' Education Management'),
        ('86104', '041004', '07', 'B', '07', 'ADMINISTRASI PENDIDIKAN', 'Pakuan 452 Bogor', 200, 2, NULL, '11111', '082518534206', '0251320123', 'pasca@unpak.ac.id', 'http://pasca.unpak.ac.id', 6, 'N', '2000-09-22', '20021', '20021', '8063/D/T/K-IV/2011', '2011-07-31', '2013-10-27', '0315/SK/BAN-PT/Akred/M/I/2017', '2017-01-10', '2022-01-10', 'A', 'A', 'A', '0426067204', '081282921597', 'YUDHIE SUCHYADI', '0251320123', 'SEMESTER', 2, 4, 'M.Pd.', 'Magister Pendidikan', '8063/D/T/K-IV/2011', '', '-'),
        ('86122', '041004', '07', 'B', '07', 'PENDIDIKAN DASAR', 'Jl. Pakuan No.1', 200, 2, NULL, '16143', '0251', '0251', 's2pendas@unpak.ac.id', 'https://pendas-pasca.unpak.ac.id', 39, 'A', '2022-09-15', '20221', '20221', '067/E/O/2022', '2022-09-15', '2052-09-15', '-', '2022-09-15', '2022-09-15', 'B', 'F', 'B', '0403087102', '', 'Lana', '', 'SEMESTER', 4, 8, 'M.Pd', 'Magister Pendidikan', '067/E/O/2022', '', 'Basic Education'),
        ('86139', '041004', '07', 'B', '07', 'MANAJEMEN PENDIDIKAN', 'Jl. Pakuan No.1', 200, 2, NULL, '16143', '089638332732', '0251320123', 'pasca@unpak.ac.id', 'https://pasca.unpak.ac.id', 6, 'A', '2023-02-15', '20232', '20232', '199/E/O/2024', '2024-02-15', '2050-02-15', '2184/SK/BAN-PT/Ak.Ppj/M/VI/2023', '2023-06-06', '2028-06-06', 'U', 'A', 'A', '0412117404', '081282921597', 'LINGGAR', '089638332732', 'SEMESTER', 2, 4, 'M.Pd', 'Magister Pendidikan', '8063/D/T/K-IV/2011', '', 'Education Management'),
        ('86206', '041004', '03', 'C', '03', 'PENDIDIKAN GURU SEKOLAH DASAR', 'jl. Pakuan', 200, 2, NULL, '16143', '08111106652', '0251-8356927', 'pgsd@unpak.ac.id', 'http://pgsd.fkip.unpak.ac.id', 6, 'A', '2007-07-19', '20071', '20151', '1923/D/T/2011', '2011-08-19', '2015-08-19', '13797/SK/BAN-PT/Ak-PPJ/S/XII/2021', '2021-12-28', '2026-12-28', 'B', 'F', 'A', '0425128802', '', 'ELLA', '085811849444', 'SEMESTER', 2, 16, 'S.Pd.', 'Sarjana Pendidikan', '3318/SK/BAN-PT/Akred/S/XII/2016 ', '', '-'),
        ('86276', '041004', '03', 'C', '03', 'PSKGJ PENDIDIKAN GURU SEKOLAH DASAR (PGSD)', 'Jl. Pakuan No.1', 200, 2, NULL, '16143', '0251-8312206', '0251-8312206', 'fkip@unpak.ac.id', 'http://fkip.unpak.ac.id', 6, 'N', '2009-02-16', '20091', '20091', '015/P/2009', '2009-02-16', '2015-02-16', '-', '2012-10-10', '2012-10-10', 'C', 'D', 'A', '0405076901', '', 'FAJAR', '', 'SEMESTER', 2, 8, 'S.Pd.', 'Sarjana Pendidikan', '-', NULL, NULL),
        ('86904', '041004', '03', 'J', '03', 'PENDIDIKAN PROFESI GURU', 'Jl. Pakuan', 200, 2, NULL, '16143', '0251-8312206', '0251-8356927', 'fkip@unpak.ac.id', 'http://www.unpak.ac.id', 8, 'A', '2018-11-30', '20191', '20191', '-', '2018-11-30', '2018-11-30', '768/SK/BAN-PT/Ak-PKP/PP/II/2021', '2021-02-10', '2026-02-10', 'B', 'B', 'A', '0405128703', '', 'M. Ganeswara', '081389121786', 'SEMESTER', 2, 2, '-', '-', '-', '', '-'),
        ('88201', '041004', '03', 'C', '03', 'PENDIDIKAN BAHASA DAN SASTRA INDONESIA', 'Jl. Pakuan', 200, 2, NULL, '16456', '081311111111', '0251-375608', 'fkip@unpak.ac.id', 'http://www.unpak.ac.id', 6, 'A', '1982-05-10', '20021', '20021', '2304/D/T/2011', '2011-08-21', '2015-08-21', '2123/SK/BAN-PT/Ak-PPJ/S/IV/2020', '2020-04-01', '2025-04-01', 'B', 'D', 'A', '0417099101', '', '', '', 'SEMESTER', 2, 8, 'S.Pd.', 'Sarjana Pendidikan', '0', '', 'Indonesian language and literature education'),
        ('88203', '041004', '03', 'C', '03', 'PENDIDIKAN BAHASA INGGRIS', 'Pakuan 452', 200, 2, NULL, '111111', '0812000000000', '02518375608', 'fkip@unpak.ac.id', 'http://fkip.unpak.ac.id', 6, 'A', '1982-05-10', '20021', '20021', '2305/D/T/2011', '2011-08-21', '2015-08-21', '13796/SK/BAN-PT/Ak-PPJ/S/XII/2021', '2022-01-11', '2027-01-11', 'B', 'A', 'A', '0404038901', '', '', '', 'SEMESTER', 2, 8, 'S.Pd.', 'Sarjana Pendidikan', '0165/SK/BAN-PT/Akred/S/I/2017', '', '-'),
        ('95102', '041004', '07', 'B', '07', 'PKLH', 'Jl. Pakuan', 200, 2, NULL, '16145', '085218534206', '02518320123', 'pasca@unpak.ac.id', 'http://pasca.unpak.ac.id', 6, 'N', '2000-09-22', '20021', '20021', '1359/D/T/K-IV/2012', '2012-10-19', '2015-10-27', '2212/SK/BAN-PT/Akred/M/VII/2017', '2017-07-04', '2022-07-04', 'A', 'A', 'A', '0323016503', '', 'Fredy Herlambang', '02518320123', 'SEMESTER', 2, 4, 'M.Ling.', 'Magister Lingkungan', '-', '', '-'),
        ('95125', '041004', '07', 'B', '07', 'PERENCANAAN WILAYAH DAN KOTA', 'Jl. Pakuan', 200, 2, NULL, '16143', '02518311007', '02518311007', 'pasca@unpak.ac.id', 'https://pasca.unpak.ac.id', 5, 'A', '2019-01-10', '20191', '20191', '-', '2019-01-10', '2030-01-20', '-', '2019-01-20', '2025-01-20', 'C', 'D', 'A', '0024087504', '081355131023', '', '', 'SEMESTER', 2, 4, 'M.P.W.K.', 'Magister Perencanaan Wilayah dan Kota', '-', '', '-');

        INSERT INTO sijamu_fakultas_unit (id, uuid, kode_fakultas, kode_prodi, nama, id_m_prodi, standalone) VALUES
        (1, '0d2fa3f8-6df3-45b8-8985-654cb49d5d03', '01', '74201', NULL, 7, 0),
        (2, 'ce5459fa-8aa7-4efc-9cb0-3f067df561c9', '02', '61201', NULL, 8, 0),
        (3, '882561ac-8ff0-442e-a0c7-e03268fd4cf8', '02', '61209', NULL, 84, 0),
        (4, '1971e6af-c1f4-4084-ad8f-d500dd431d4f', '02', '62201', NULL, 12, 0),
        (5, '47ec32d5-3371-4352-a153-d7c90a9f3a2a', '03', '84205', NULL, 18, 0),
        (6, '130fd75c-20c7-4ec0-bd87-eeb9eaa0862e', '03', '84206', NULL, 73, 0),
        (7, '48184036-62ce-4a7c-8598-9ccef3e4cbd2', '03', '86206', NULL, 19, 0),
        (8, 'c81e8e79-c535-4b50-ad35-63eb5f5c1380', '03', '86276', 'PKSGJ', NULL, 0),
        (9, '2db0c45a-5c2a-4fc4-ae8e-5fddca464847', '03', '86904', 'PPG', NULL, 0),
        (10, '36428f6e-778f-43ce-8b99-4a056ed3d384', '03', '88201', NULL, 17, 0),
        (11, 'b3fad690-be32-4d30-a6a6-f7e9f5f83be8', '03', '88203', NULL, 16, 0),
        (12, 'e00f35c9-f679-4584-a0c5-8f4ae4e69b0f', '04', '70201', NULL, 23, 0),
        (13, '9f348e73-3d42-4fd6-ad34-127fd5f6acd6', '04', '79201', NULL, 21, 0),
        (14, 'd5aed764-36d1-44ef-ab7f-0109d6cce47e', '04', '79202', NULL, 20, 0),
        (15, 'b28a1af9-792d-46af-8811-86383f4a3a21', '04', '79204', NULL, 22, 0),
        (16, '53740d58-f981-4890-8f93-4ed2c0690ae0', '05', '20201', NULL, 27, 0),
        (17, '0d1f9211-8a20-4e63-8207-a66ec021d419', '05', '22201', NULL, 26, 0),
        (18, '669c148e-f16a-4dad-ad52-d121ef62af80', '05', '29201', NULL, 24, 0),
        (19, '1223561f-e0de-4809-80a1-3c2d596f5ef5', '05', '34201', NULL, 28, 0),
        (20, '001e973e-f150-46c5-bcd7-0528bbc6a590', '05', '35201', NULL, 25, 0),
        (21, '8a114de9-581f-4ca2-b47e-766afed63f9b', '06', '44201', NULL, 31, 0),
        (22, 'b343d749-a399-4d8c-9c4e-2be2f94c5dca', '06', '46201', NULL, 29, 0),
        (23, '820926e7-561d-4a3b-b9e7-a6ef48cb9e82', '06', '47201', NULL, 30, 0),
        (24, '89d432f2-d1a8-439e-8ebd-b3429eb2440c', '06', '48201', NULL, 33, 0),
        (25, '9f3c62a3-deea-49f7-b17e-c676cddc4d92', '06', '55201', NULL, 32, 0),
        (26, 'fc4bd287-7cb1-4138-b825-f9e7d18e6e3a', '07', '61001', NULL, 58, 0),
        (27, '547b4104-e95a-4e84-b68a-d0d19952340b', '07', '61101', NULL, 54, 0),
        (28, '8b2b8c4c-374c-45fe-add0-2d861a72431e', '07', '61105', NULL, 59, 0),
        (29, '246efbbb-aea5-4b00-b3d8-5183e417721e', '07', '74101', NULL, 55, 0),
        (30, '2b2af760-3d54-42df-8c75-043b1ad8b40c', '07', '84101', NULL, 56, 0),
        (31, '8b897ac6-84cb-4e69-896c-1d878f0e015d', '07', '86004', NULL, 57, 0),
        (32, '78fdd67c-0577-49f1-9ef6-0d474e58b049', '07', '86104', NULL, 53, 0),
        (33, '27d61ecb-a2d7-4cd0-a58e-5afe43b4f16b', '07', '86122', NULL, 87, 0),
        (34, 'e972fae8-56f4-4fdb-91d6-cf8e2d7b7896', '07', '95102', 'PKLH', NULL, 0),
        (35, 'd0fe41d9-c044-42bb-bb98-ba73324912ff', '07', '95125', NULL, 70, 0),
        (36, 'ea275d77-f3f4-4bf8-8f17-6f47e1bcae4a', '08', '56401', NULL, 34, 0),
        (37, '3e345fd0-c32d-4c0f-a76c-7992e42463c4', '08', '57401', NULL, 35, 0),
        (38, 'b4b771bf-3610-46be-bcb2-da9aa6d6178b', '08', '61403', NULL, 14, 0),
        (39, 'db9dc7ea-8e41-428d-8aa3-63cda1c40fec', '08', '61406', NULL, 15, 0),
        (40, 'f8980dee-f0a3-47fc-bb41-c7462436a120', '08', '62401', NULL, 13, 0),
        (41, '935b0af7-07f8-4e9a-a83b-a36d54bdc5e2', 'U1', NULL, 'PUTIK (PUSAT TEKNOLOGI INFORMASI DAN KOMUNIKASI)', NULL, 0),
        (64, 'b9bd4e7d-2cee-4426-9be0-20428f5a146d', 'U2', NULL, 'SEKRETARIS UNIVERSITAS', NULL, 0),
        (65, 'ac485198-4731-4903-800d-0bccb3b7adb1', 'U3', NULL, 'WAKIL REKTOR 1', NULL, 0),
        (66, '0c9601d4-81ec-4bf1-b29a-ff0df3065c3e', 'U4', NULL, 'WAKIL REKTOR 2', NULL, 0),
        (67, 'fe918795-1c4c-4e7d-8f6b-7a0dbe6bf73b', 'U5', NULL, 'WAKIL REKTOR 3', NULL, 0),
        (68, '6216408c-305e-410f-a1e6-c93d705f30ff', 'U6', NULL, 'Kantor Pengembangan Karir dan Tracer Study', NULL, 0),
        (69, '2e7d3374-a140-4b43-a74f-d4d4be98d0ed', 'U7', NULL, 'LPM', NULL, 0),
        (70, '8419016f-77a7-4616-a486-3bbadda923f5', 'U8', NULL, 'LPPM', NULL, 0),
        (71, 'ba9394e5-cfa1-480e-825f-6fc306ee45ee', 'U9', NULL, 'PERPUSTAKAAN PUSAT', NULL, 0),
        (72, '9e814dff-d3e7-4d8f-bdd1-93a37400a3ce', 'U10', NULL, 'HUMAS DAN PROMOSI', NULL, 0),
        (73, 'fd927be3-3c48-4e8e-8f5f-00655886d3b7', '06', NULL, 'PERPUSTAKAAN MIPA', 61, 0),
        (74, '326eb279-9026-4f2d-8245-45bd13c47e67', '06', NULL, 'DIVISI PENELITIAN DAN PENGABDIAN FMIPA', 62, 0),
        (75, '6e1dfc64-1788-4a2c-8450-b30c584e3858', 'U13', NULL, 'UNPAK PRES', NULL, 0),
        (76, '209e5085-af39-4d3b-a373-57dfc57d3582', '03', NULL, 'FKIP KKN', 65, 0),
        (77, '2be4fe70-b909-4e80-8ba3-dde2ae6338a8', '03', NULL, 'LABORATORIUM FKIP', 66, 0),
        (78, 'e50a32ac-4052-4c08-86bd-5e3abfbc9476', '03', NULL, 'FKIP E-JOURNAL', 68, 0),
        (79, '19183828-8671-416d-86ff-a7653b242465', '06', NULL, 'UPMF MIPA', 75, 0),
        (80, '546a5292-ca15-4a9c-9ef7-3fe319c45f85', '06', NULL, 'Laboratorium FMIPA ', 76, 0),
        (81, '87a30c72-b758-4009-b912-0daa02b90948', '06', NULL, 'Laboratorium Service ', 77, 0),
        (82, '48abebf3-2a32-426c-a88d-b12c3b240a04', '06', NULL, 'Apotek', 78, 0),
        (83, '5cbdf964-65ef-4a05-80c1-66c44e4e368b', '06', NULL, 'Unit Pelayanan Kesehatan', 79, 0),
        (84, '7bc0e072-1aad-42ed-8a59-c1719a7fec65', '06', NULL, 'Unit Publikasi dan jurnal', 80, 0),
        (85, '7b78a3e8-f611-4e84-88d4-30efe1fbb0a7', '06', NULL, 'Unit Comstract dan Data Sience Center', 81, 0),
        (86, '1a1ab873-346c-483a-bad5-1ae857b96da3', '03', NULL, 'Perpustakaan FKIP', 82, 0),
        (87, 'd3714886-6e0d-4e08-a696-b77f1e6bfc16', '03', NULL, 'Bimbingan Konseling FKIP', 83, 0),
        (88, '0eb236ae-5e4b-44ff-9cdc-753a99d8f5af', 'U26', NULL, 'Pusat Inovasi', NULL, 0),
        (89, 'db0279c5-8301-42fc-bb7b-ebe010672193', 'U27', NULL, 'Inkubator Bisnis', NULL, 0),
        (91, 'dea9a83f-70b3-4295-85ed-459eb1a9f6a0', '01', NULL, NULL, NULL, 1),
        (94, 'd6182e3e-b71d-449b-8ed8-947dbf98ab1b', '02', NULL, NULL, NULL, 1),
        (97, '3dc75fe5-81b3-407e-8a98-c7bf586acc03', '03', NULL, NULL, NULL, 1),
        (100, 'cd828022-a823-4c2d-916a-d6a0624414bc', '04', NULL, NULL, NULL, 1),
        (103, '7620b4d7-309a-46fc-a0f1-e5161ebbead5', '05', NULL, NULL, NULL, 1),
        (106, '58fef9e7-be02-4db9-8e70-0d4ea874708b', '06', NULL, NULL, NULL, 1),
        (109, 'db8dbc27-98b2-41e8-abcd-aaf130722f6b', '07', NULL, NULL, NULL, 1),
        (112, 'b831040a-e97d-49ff-ba96-1b316a3c24ad', '08', NULL, NULL, NULL, 1),
        (113, '5ebddaf0-8c63-4fd6-a1f9-436491a61693', 'U155', NULL, 'Kantor Kemitraan dan Hubungan Internasional', NULL, 0),
        (116, '785bb88d-5018-48db-81dd-3aa23c5f4fb8', 'U28', NULL, 'BAUM', NULL, 0),
        (119, 'eefd25ce-fa25-41fc-ab59-069b6b248a53', 'U29', NULL, 'BAAK', NULL, 0),
        (121, 'fb52bb92-b7db-4145-ae06-ee45239d2cda', '08', '56209', NULL, NULL, 0),
        (122, '9ea41450-d5d0-4732-a5aa-a2b723ad0dcd', 'U28', NULL, 'WAKIL REKTOR 4', NULL, 0),
        (125, '36428f6e-778f-43ce-8b99-4a056ed3d385', 'U29', NULL, 'Kemitraan', NULL, 0),
        (126, '5246f479-67af-4e67-83b1-b3aa3ba801bf', '07', '55101', NULL, NULL, 0),
        (127, 'f72658bf-c6c5-46a8-9612-52a6dde66079', '06', '48901', NULL, NULL, 0),
        (128, 'd57c622c-850b-4ce9-8756-27ce0a38322c', '03', NULL, 'LABORATORIUM SENI', 66, 0),
        (129, 'b13f35eb-91c9-4f5c-a79c-67542c7a195b', '03', NULL, 'LABORATORIUM BAHASA INGGRIS', 66, 0),
        (130, 'a85f99c4-bbe1-4f9c-827a-fea3b958ab58', '03', NULL, 'LABORATORIUM MICROTEACHING', 66, 0),
        (131, 'c031ad9c-5288-4343-bfee-0a9854c76d77', '03', NULL, 'LABORATORIUM KOMPUTER', 66, 0);

        INSERT INTO renstra (id, uuid, tahun, fakultas_unit_old, fakultas_unit, periode_upload_mulai, periode_upload_akhir, periode_assesment_dokumen_mulai, periode_assesment_dokumen_akhir, periode_assesment_lapangan_mulai, periode_assesment_lapangan_akhir, kodeAkses, auditee, auditor1, auditor2, catatan, catatan2, created_at, updated_at) VALUES
        (368, 'c67a37c3-7f25-43de-835d-e4bece0eb308', '2024', 0, 91, '2024-10-18 00:00:00', '2024-11-22 00:00:00', '2024-11-23 00:00:00', '2024-12-11 00:00:00', '2024-12-12 00:00:00', '2024-12-31 00:00:00', '567244', 232, 70, 415, NULL, NULL, '2024-10-18 13:41:48', '2024-12-31 10:37:12'),
        (371, '0025699d-d69b-41e4-b712-f437aa15d3b1', '2024', 0, 1, '2024-10-18 00:00:00', '2024-11-22 00:00:00', '2024-11-23 00:00:00', '2024-12-11 00:00:00', '2024-12-12 00:00:00', '2024-12-31 00:00:00', '567244', 235, 70, 415, NULL, 'o	mahasiswa asing \r\no	Fasilitas magang\r\no	Inovasi mahasiswa', '2024-10-18 13:42:09', '2025-09-11 15:11:59'),
        (374, '88f3f34e-ee4e-4ca0-8bbc-ddd4b726ca06', '2024', 0, 94, '2024-10-18 00:00:00', '2024-11-22 00:00:00', '2024-11-23 00:00:00', '2024-12-11 00:00:00', '2024-12-12 00:00:00', '2025-02-27 00:00:00', '895730', 241, 43, 166, NULL, NULL, '2024-10-18 13:44:31', '2025-02-27 10:06:05'),
        (377, '7ccb8a57-f39d-4218-98ea-a4242f281de1', '2024', 0, 2, '2024-10-18 00:00:00', '2024-11-22 00:00:00', '2024-11-23 00:00:00', '2024-12-11 00:00:00', '2024-12-12 00:00:00', '2025-02-27 00:00:00', '290347', 247, 64, 151, 'Terdapat beberapa data lapangan yang jauh lebih tinggi daripada data yang dilaporkan', NULL, '2024-10-18 13:45:46', '2025-02-27 10:04:17'),
        (380, '72d597d7-687a-49c1-bbd9-166e6688ed13', '2024', 0, 4, '2024-10-18 00:00:00', '2024-11-22 00:00:00', '2024-11-23 00:00:00', '2024-12-11 00:00:00', '2024-12-12 00:00:00', '2025-04-14 00:00:00', '567244', 244, 70, 415, NULL, 'Mahasiswa asing', '2024-10-18 13:47:09', '2025-04-14 08:16:20'),
        (383, 'f98cdbb9-a3f5-4b1b-be3e-0cd6e495cdc0', '2024', 0, 3, '2024-10-18 00:00:00', '2024-11-22 00:00:00', '2024-11-23 00:00:00', '2024-12-11 00:00:00', '2024-12-12 00:00:00', '2025-03-24 00:00:00', '470975', 250, 43, 211, 'ada beberapa temuan yang akan dikonfirmasi saat AL', 'Belum terdapat mahasiswa asing. Belum terdapat penelitian mahasiswa yang mendapatkan HKI', '2024-10-18 13:47:36', '2025-03-24 10:46:47'),
        (386, 'e9d1c872-4fb5-437f-ae95-a5de28f1d799', '2024', 0, 97, '2024-10-18 00:00:00', '2024-11-22 00:00:00', '2024-11-23 00:00:00', '2024-12-11 00:00:00', '2024-12-12 00:00:00', '2024-12-23 00:00:00', '616584', 253, 64, 166, NULL, NULL, '2024-10-18 13:48:09', '2024-12-23 08:03:22'),
        (392, 'e6276ad3-3450-4d4c-928b-f287b3244030', '2024', 0, 10, '2024-10-18 00:00:00', '2024-11-22 00:00:00', '2024-11-23 00:00:00', '2024-12-11 00:00:00', '2024-12-12 00:00:00', '2024-12-21 00:00:00', '605998', 256, 79, 217, 'Beberapa bukti mohon diupload sesuai kesepakatan', 'data ada yang tidak terbaca', '2024-10-18 13:48:46', '2025-01-02 11:12:47'),
        (395, '7184d9ed-bd51-4b16-880d-e41263c558cc', '2024', 0, 11, '2024-10-18 00:00:00', '2024-11-22 00:00:00', '2024-11-23 00:00:00', '2024-12-11 00:00:00', '2024-12-12 00:00:00', '2024-12-21 00:00:00', '937115', 259, 76, 415, 'Data dan kompetensi lulusan tidak ditemukan', 'mahasiswa asing\r\nprestasi dalam ruang lingkup internasional\r\nfasilitas magang pada link terdapat 5', '2024-10-18 13:49:18', '2024-12-06 22:20:02'),
        (398, '54409c75-ba7c-47a6-9d33-4d178b5ea9ec', '2024', 0, 5, '2024-10-18 00:00:00', '2024-11-22 00:00:00', '2024-11-23 00:00:00', '2024-12-11 00:00:00', '2024-12-12 00:00:00', '2024-12-27 00:00:00', '567244', 262, 70, 175, NULL, NULL, '2024-10-18 13:49:51', '2024-12-27 14:12:44'),
        (401, '0451c391-6d2e-465d-8f20-bff30fd4c221', '2024', 0, 7, '2024-10-18 00:00:00', '2024-11-22 00:00:00', '2024-11-23 00:00:00', '2024-12-11 00:00:00', '2024-12-12 00:00:00', '2024-12-23 00:00:00', '186377', 265, 82, 415, 'data sudah terkonfirmasi', 'o	Mahasiswa asing\r\no	Fasilitas kegiatan magang\r\no	Inovasi baru 1 dari 3 (abon', '2024-10-18 13:52:51', '2024-12-23 11:46:38'),
        (404, '88c8f743-2cee-470a-bbce-a2a6caa15a8c', '2024', 0, 6, '2024-10-18 00:00:00', '2024-11-22 00:00:00', '2024-11-23 00:00:00', '2024-12-11 00:00:00', '2024-12-12 00:00:00', '2025-02-27 00:00:00', '567244', 268, 70, 151, NULL, NULL, '2024-10-18 13:53:28', '2025-02-27 10:04:31'),
        (407, '9554e65e-5106-48cd-b006-934507a81ae0', '2024', 0, 9, '2024-10-18 00:00:00', '2024-11-22 00:00:00', '2024-11-23 00:00:00', '2024-12-11 00:00:00', '2024-12-12 00:00:00', '2025-03-24 00:00:00', '398944', 271, 91, 421, 'Belum ada Bukti laporan kegiatan dan publikasi semua dosen', '1. capaian yang diharapkan mahasiswa mendapatkan HKI, tapi yang dilampirkan HKI dosen', '2024-10-18 13:55:51', '2025-03-24 11:04:02'),
        (410, '3c4c7234-c6b2-48ff-a575-162eba7d59a8', '2024', 0, 100, '2024-10-18 00:00:00', '2024-11-22 00:00:00', '2024-11-23 00:00:00', '2024-12-11 00:00:00', '2024-12-12 00:00:00', '2025-01-09 00:00:00', '611830', 274, 49, 223, 'terdapat 15 indikator  yang belum mencapai 100% target.. meskipun sudah dilaksanakan yaitu indikator no :6,14,22,24,26,27,29,31,40,41,42,43,45,46,57.  terdapat 1 indicator GB yang sama sekali belum tercapai , yaitu indikator no 6', 'Sudah terlaksana baik akademik dan non akademik ada 24', '2024-10-18 13:56:35', '2025-01-07 22:47:48'),
        (413, '618170ae-489f-43aa-aa42-deeeb48fd8e1', '2024', 0, 13, '2024-10-18 00:00:00', '2024-11-22 00:00:00', '2024-11-23 00:00:00', '2024-12-11 00:00:00', '2024-12-12 00:00:00', '2025-02-27 00:00:00', '057735', 283, 43, 421, NULL, '1. fasilitas kegiatan magang/praktek industri linknya kosong\r\n2. Jumlah inovasi mahasiswa lampiran berbeda dengan capaian yang diharapkan', '2024-10-18 13:57:03', '2025-02-27 10:04:36'),
        (416, '845d9222-0511-4610-899e-e731e9b6cd5c', '2024', 0, 14, '2024-10-18 00:00:00', '2024-11-22 00:00:00', '2024-11-23 00:00:00', '2024-12-11 00:00:00', '2024-12-12 00:00:00', '2025-11-25 00:00:00', '567244', 277, 70, 199, NULL, 'belum ada dokumen karya inovasi mahasiswa', '2024-10-18 13:59:14', '2025-11-25 10:04:40'),
        (419, '1f5df06f-8672-4b51-9f50-f323ace4c6df', '2024', 0, 15, '2024-10-18 00:00:00', '2024-11-22 00:00:00', '2024-11-23 00:00:00', '2024-12-11 00:00:00', '2024-12-12 00:00:00', '2024-12-21 00:00:00', '110875', 280, 52, 217, NULL, '1. Standard Prestasi Mahasiswa baru 4 Kurang 1 dari 5\r\n2. Belum Ada Mahasiswa Asing (Kriteria 3)\r\n3. Belum ada Karya Inovasi Mahasiswa Per Tahun', '2024-10-18 13:59:36', '2024-12-16 08:36:59'),
        (422, '6730228c-4325-4b0d-9be4-d050fe422c90', '2024', 0, 12, '2024-10-18 00:00:00', '2024-11-22 00:00:00', '2024-11-23 00:00:00', '2024-12-11 00:00:00', '2024-12-12 00:00:00', '2025-04-11 00:00:00', '567244', 286, 70, 187, NULL, '1. Bukti mahasiswa asing belum otentik\r\n2. Beberapa capaian indikator perlu dibuat listnya dalam catatan auditee untuk memudahkan tracing', '2024-10-18 14:00:36', '2025-04-11 10:17:55'),
        (425, '422199b8-64db-45e1-b325-76fb5c3af6ea', '2024', 0, 103, '2024-10-18 00:00:00', '2025-07-05 00:00:00', '2025-07-07 00:00:00', '2025-07-08 00:00:00', '2025-07-09 00:00:00', '2025-07-10 00:00:00', '567244', 289, 85, 178, NULL, NULL, '2024-10-18 14:01:04', '2025-07-01 10:36:19'),
        (428, '91ba15dd-7501-478f-8ee9-48d60285b36e', '2024', 0, 16, '2024-10-18 00:00:00', '2025-07-05 00:00:00', '2025-07-07 00:00:00', '2025-07-08 00:00:00', '2025-07-09 00:00:00', '2025-07-10 00:00:00', '282803', 301, 85, 181, 'UNTUK PRODI : beberapa link tidak dapat dibuka; beberapa link tidak menunjukkan bukti yang diminta; Beberapa bukti perlu dilengkapi.\r\nUNTUK SIAMIDA: pertanyaan no 4 bukan untuk prodi, tapi Fakultas.', NULL, '2024-10-18 14:02:42', '2025-07-03 07:31:25'),
        (431, 'b7b087da-fa23-45f1-9418-928f2d1f4794', '2024', 0, 17, '2024-10-18 00:00:00', '2024-11-22 00:00:00', '2024-11-23 00:00:00', '2024-12-11 00:00:00', '2024-12-12 00:00:00', '2024-12-21 00:00:00', '658393', 298, 64, 229, NULL, 'Fasilitas kegiatan magang belum ada', '2024-10-18 14:03:05', '2024-12-20 13:54:36'),
        (434, '8d2614fe-3116-4aa7-bc33-4d900fc4e297', '2024', 0, 18, '2024-10-18 00:00:00', '2024-11-22 00:00:00', '2024-11-23 00:00:00', '2024-12-11 00:00:00', '2024-12-12 00:00:00', '2025-02-27 00:00:00', '267689', 304, 43, 229, NULL, 'Pada Tahun 2024, mahasiswa Teknik Geodesi membuat 4 inovasi.', '2024-10-18 14:03:37', '2025-02-27 10:04:43'),
        (437, '75950dcb-0dfd-464e-ae85-6d9b20373a93', '2024', 0, 19, '2024-10-18 00:00:00', '2025-07-05 00:00:00', '2025-07-07 00:00:00', '2025-07-08 00:00:00', '2025-07-09 00:00:00', '2025-07-10 00:00:00', '511782', 292, 85, 196, 'UNTUK PRODI: Beberapa bukti ada yang tidak sesuai dengan yang ditanyakan; Beberapa bukti ada yang perlu diperbaiki penyajiannya agar memudahkan auditor menganalisis capaian dari auditee.\r\nUNTUK SIAMIDA: pertanyaan no 4 bukan untuk prodi tapi untuk Fakultas; pertanyaan no 74  tidak selaras antara capaian indikator yang diharapkan dengan keterangan YA /TIDAK berdasarkan isian dari auditee', NULL, '2024-10-18 14:03:55', '2025-07-01 10:36:52'),
        (440, '7468537a-bd51-4272-96b1-a4864e6503f3', '2024', 0, 20, '2024-10-18 00:00:00', '2024-11-22 00:00:00', '2024-11-23 00:00:00', '2024-12-11 00:00:00', '2024-12-12 00:00:00', '2024-12-23 00:00:00', '567244', 295, 70, 223, NULL, 'Tercapai', '2024-10-18 14:04:26', '2024-12-22 17:55:45'),
        (443, '162686b4-6e71-4231-afc4-6b45933accc3', '2024', 0, 106, '2024-10-18 00:00:00', '2024-11-22 00:00:00', '2024-11-23 00:00:00', '2024-12-11 00:00:00', '2024-12-12 00:00:00', '2025-02-27 00:00:00', '253523', 307, 43, 169, NULL, '1. Tidak ada mahasiswa asing\r\n2. Jumlah karya inovasi mahasiswa hanya 6', '2024-10-18 14:05:07', '2025-02-27 10:06:11'),
        (446, '579b3a0a-52ed-46a8-9823-37613d3e6a69', '2024', 0, 22, '2024-10-18 00:00:00', '2024-11-22 00:00:00', '2024-11-23 00:00:00', '2024-12-11 00:00:00', '2024-12-12 00:00:00', '2025-02-27 00:00:00', '343923', 310, 82, 151, 'Bukti sertifikat kompetensi telah lengkap, 100%', NULL, '2024-10-18 14:05:54', '2025-02-27 08:51:12'),
        (449, '4c06716d-d023-4740-b965-d9260f51dba3', '2024', 0, 23, '2024-10-18 00:00:00', '2024-11-22 00:00:00', '2024-11-23 00:00:00', '2024-12-11 00:00:00', '2024-12-12 00:00:00', '2024-12-21 00:00:00', '460861', 313, 76, 160, NULL, 'tidak ada temuan', '2024-10-18 14:06:27', '2024-12-05 16:05:16'),
        (452, '71e67647-3629-4f64-9d1b-65cfff84d35f', '2024', 0, 21, '2024-10-18 00:00:00', '2024-11-22 00:00:00', '2024-11-23 00:00:00', '2024-12-11 00:00:00', '2024-12-12 00:00:00', '2025-05-20 00:00:00', '567244', 316, 67, 178, '1. Jumlah mata kuliah yang mengintegrasikan hasil penelitian dan PKM ke dalam perkuliahan hanya mencapai 19% dari yang seharusnya 40%\r\n2. Proses pembelajaran berbasis PBL/PJBL tiap semester hanya mencapai 65% dari seharusnya 85%', NULL, '2024-10-18 14:07:21', '2025-05-20 08:09:48'),
        (455, '1e680de6-b16f-44f1-a9c9-0d01800305c1', '2024', 0, 25, '2024-10-18 00:00:00', '2024-11-22 00:00:00', '2024-11-23 00:00:00', '2024-12-11 00:00:00', '2024-12-12 00:00:00', '2025-07-26 00:00:00', '325055', 319, 43, 217, NULL, 'data mahasiswa asing yang kuliah di ilkom belum ada', '2024-10-18 14:09:11', '2025-07-21 10:32:07'),
        (458, 'd77d5a21-01d6-40ed-8789-c64236bb9943', '2024', 0, 24, '2024-10-18 00:00:00', '2024-11-22 00:00:00', '2024-11-23 00:00:00', '2024-12-11 00:00:00', '2024-12-12 00:00:00', '2024-12-21 00:00:00', '704099', 322, 88, 226, 'Beberapa temuan perlu dikonfirmasi saat assesment lapangan', 'Terdapat ketidaksesuaian bukti yang harus diverifikasi saat Assesmen Lapangan', '2024-10-18 14:12:35', '2025-07-29 17:33:08'),
        (461, '0d92c660-8db8-4819-a6f3-d388c8e98bfc', '2024', 0, 109, '2024-10-18 00:00:00', '2024-11-22 00:00:00', '2024-11-23 00:00:00', '2024-12-11 00:00:00', '2024-12-12 00:00:00', '2025-02-27 00:00:00', '892955', 325, 61, 151, NULL, NULL, '2024-10-18 14:12:59', '2025-02-27 10:06:17'),
        (464, '79b6849f-c955-4335-a097-d09b97fd8b63', '2024', 0, 31, '2024-10-18 00:00:00', '2024-11-22 00:00:00', '2024-11-23 00:00:00', '2024-12-11 00:00:00', '2024-12-12 00:00:00', '2024-12-23 00:00:00', '567244', 328, 91, 178, NULL, NULL, '2024-10-18 14:13:40', '2024-12-22 17:56:18'),
        (467, 'c3d03c51-bf47-4a8f-bf7f-9adfc1045613', '2024', 0, 26, '2024-10-18 00:00:00', '2024-11-22 00:00:00', '2024-11-23 00:00:00', '2024-12-11 00:00:00', '2024-12-12 00:00:00', '2025-02-27 00:00:00', '908508', 331, 43, 226, NULL, NULL, '2024-10-18 14:16:20', '2025-02-27 10:04:55'),
        (470, 'bfd4629c-0e29-443b-b38f-648e3d45f578', '2024', 0, 32, '2024-10-18 00:00:00', '2024-11-22 00:00:00', '2024-11-23 00:00:00', '2024-12-11 00:00:00', '2024-12-12 00:00:00', '2025-01-02 00:00:00', '138373', 334, 58, 421, NULL, NULL, '2024-10-18 14:18:18', '2025-01-02 09:57:53'),
        (473, 'a4ba6674-a448-446e-a549-c458e1dffffb', '2024', 0, 28, '2024-10-18 00:00:00', '2024-11-22 00:00:00', '2024-11-23 00:00:00', '2024-12-11 00:00:00', '2024-12-12 00:00:00', '2024-12-27 00:00:00', '813700', 337, 43, 217, NULL, 'belum tercapai semua restra', '2024-10-18 14:19:40', '2025-03-24 11:06:23'),
        (476, '847840d6-1dc6-48b4-a98b-e3c2d65e56a1', '2024', 0, 29, '2024-10-18 00:00:00', '2024-11-22 00:00:00', '2024-11-23 00:00:00', '2024-12-11 00:00:00', '2024-12-12 00:00:00', '2025-03-15 00:00:00', '567244', 340, 70, 223, NULL, 'Tercapai', '2024-10-18 14:22:18', '2025-03-15 08:56:30'),
        (479, '134e0172-ed77-4dce-8f19-6b3cf12e3ba6', '2024', 0, 27, '2024-10-18 00:00:00', '2024-11-22 00:00:00', '2024-11-23 00:00:00', '2024-12-11 00:00:00', '2024-12-12 00:00:00', '2025-03-24 00:00:00', '567244', 343, 70, 421, NULL, '1. capaian mahasiswa asing tidak tercapai\r\n2. Prestasi mahasiswa tidak ada', '2024-10-18 14:23:06', '2025-03-24 11:06:35'),
        (482, '61881d94-1eab-4a42-bd53-1d5b772c046d', '2024', 0, 30, '2024-10-18 00:00:00', '2024-11-22 00:00:00', '2024-11-23 00:00:00', '2024-12-11 00:00:00', '2024-12-12 00:00:00', '2024-12-21 00:00:00', '408731', 346, 88, 169, NULL, 'Auditee melampirkan bukti fasilitas di Fakultas.', '2024-10-18 14:24:05', '2025-01-02 10:48:32'),
        (485, 'da4b253f-9b4d-4457-a6d3-9907d047f4d8', '2024', 0, 35, '2024-10-18 00:00:00', '2024-11-22 00:00:00', '2024-11-23 00:00:00', '2024-12-11 00:00:00', '2024-12-12 00:00:00', '2024-12-21 00:00:00', '414481', 349, 64, 175, NULL, NULL, '2024-10-18 14:24:32', '2025-09-12 10:12:54'),
        (488, '9ce7643f-72d6-478e-ae8e-65be8c7a5a62', '2024', 0, 33, '2024-10-18 00:00:00', '2024-11-22 00:00:00', '2024-11-23 00:00:00', '2024-12-11 00:00:00', '2024-12-12 00:00:00', '2025-03-24 00:00:00', '567244', 352, 70, 415, NULL, 'o	Mahasiswa asing\r\no	Fasilitas magang', '2024-10-18 14:25:06', '2025-03-24 11:05:18'),
        (491, '656b17df-ced2-44db-a9ec-748d16b411ce', '2024', 0, 112, '2024-10-18 00:00:00', '2024-11-22 00:00:00', '2024-11-23 00:00:00', '2024-12-11 00:00:00', '2024-12-12 00:00:00', '2025-01-21 00:00:00', '175735', 355, 49, 421, 'Indikator 6,7,13,15,25 sama sekali belum ada.  Sedangkan indikator yang sudah terlaksana namun belum mencapai target adalah :12,14,24,26,27,29,31,32,37,39,40,41,42,43,46,48,51,56', NULL, '2024-10-18 14:26:08', '2025-01-21 11:44:35'),
        (494, '41a48232-eb0c-4cda-a554-ac9819693656', '2024', 0, 40, '2024-10-18 00:00:00', '2024-11-22 00:00:00', '2024-11-23 00:00:00', '2024-12-11 00:00:00', '2024-12-12 00:00:00', '2025-02-27 00:00:00', '567244', 358, 43, 178, NULL, NULL, '2024-10-18 14:26:40', '2025-02-27 10:05:06'),
        (497, '293c8b02-23d8-4c21-ad6a-4f9379dcc22d', '2024', 0, 38, '2024-10-18 00:00:00', '2024-11-22 00:00:00', '2024-11-23 00:00:00', '2024-12-11 00:00:00', '2024-12-12 00:00:00', '2024-12-21 00:00:00', '760746', 364, 73, 421, 'Belum tersedia kurikulum MBKM', NULL, '2024-10-18 14:27:04', '2024-12-19 10:28:53'),
        (500, '4e5f4638-e4d3-4e7c-b0e8-726392ae9ecd', '2024', 0, 39, '2024-10-18 00:00:00', '2024-11-22 00:00:00', '2024-11-23 00:00:00', '2024-12-11 00:00:00', '2024-12-12 00:00:00', '2024-12-23 00:00:00', '567244', 367, 70, 223, NULL, NULL, '2024-10-18 14:27:49', '2024-12-22 17:57:21'),
        (503, 'f642c77d-7bdb-4318-b1c3-684b3749f72b', '2024', 0, 36, '2024-10-18 00:00:00', '2024-11-22 00:00:00', '2024-11-23 00:00:00', '2024-12-11 00:00:00', '2024-12-12 00:00:00', '2024-12-21 00:00:00', '884918', 370, 115, 415, '1. Proses pembelajaran PBL kurang 5 % dari Target\r\n2. Belum ada Exchange Lecturer\r\n3. Portofolio mahasiswa melalui Kurikulum OBE kurang 6 %\r\n4. Dosen belum mampu berbahasa inggris dengan baik\r\n5. PKM multi disiplin ilmu belum mencukupi dari 82 baru 40 %\r\n6. Belum ada kerjasama internasional untuk tridharma', 'o	Mahasiswa asing\r\no	Prestasi (belum ada internasional)\r\no	HKI\r\no	Fasilitas magang\r\no	Inovasi', '2024-10-18 14:28:26', '2024-12-13 14:10:28'),
        (506, '92acd180-fe14-432b-92f6-f6c26fa7f6fc', '2024', 0, 37, '2024-10-18 00:00:00', '2024-11-22 00:00:00', '2024-11-23 00:00:00', '2024-12-11 00:00:00', '2024-12-12 00:00:00', '2024-12-21 00:00:00', '247394', 373, 106, 199, 'Beberapa indikator belum terisi. Satuan capaian masih belum sesuai. Beberapa indikator tidak disertakan bukti.', 'tidak ada mahasiswa asing', '2024-10-18 14:28:49', '2024-12-16 08:51:00'),
        (524, 'c1886c9e-98d5-4d4c-a7ba-8ef7123622fd', '2024', 88, 88, '2025-02-05 00:00:00', '2025-02-17 00:00:00', '2025-02-18 00:00:00', '2025-02-23 00:00:00', '2025-02-24 00:00:00', '2025-02-27 00:00:00', '123456', 394, 70, 178, NULL, NULL, '2025-02-04 15:16:19', '2025-02-04 22:00:38'),
        (527, '57f66a0d-4099-45a7-9150-6835cfb3f05d', '2026', 0, 125, '2025-02-05 00:00:00', '2025-02-17 00:00:00', '2025-02-18 00:00:00', '2025-02-23 00:00:00', '2025-02-24 00:00:00', '2025-02-27 00:00:00', '654321', 400, 452, 415, 'untuk tautan atau link data, sebaiknya yang memudahkan diakses (pengembangan silakerma, untuk memudahkan kategorisasi)', NULL, '2025-02-04 15:23:26', '2025-02-27 09:36:39'),
        (530, 'ba442eb4-6ba9-485c-8cd9-f75699941ef6', '2024', 0, 70, '2025-02-05 00:00:00', '2025-02-17 00:00:00', '2025-02-18 00:00:00', '2025-02-23 00:00:00', '2025-02-24 00:00:00', '2025-03-04 00:00:00', '112233', 382, 43, 151, NULL, NULL, '2025-02-04 15:24:12', '2025-03-03 08:13:00'),
        (533, 'b3970fae-2fa3-44c8-9d19-c9562563ad28', '2024', 0, 72, '2025-02-05 00:00:00', '2025-02-17 00:00:00', '2025-02-18 00:00:00', '2025-02-23 00:00:00', '2025-02-24 00:00:00', '2025-03-04 00:00:00', '123456', 388, 70, 178, NULL, NULL, '2025-02-04 15:25:00', '2025-03-03 08:54:53'),
        (536, 'daf16b86-2997-4db4-8763-7fe3aaa6f40e', '2024', 0, 69, '2025-02-05 00:00:00', '2025-02-18 00:00:00', '2025-02-19 00:00:00', '2025-02-23 00:00:00', '2025-02-24 00:00:00', '2025-02-27 00:00:00', '123123', 413, 55, 163, NULL, NULL, '2025-02-04 15:26:02', '2025-02-18 09:29:30'),
        (539, 'cd7f75ea-4f72-484f-ad89-f1b08d397b16', '2024', 89, NULL, '2025-02-05 00:00:00', '2025-02-17 00:00:00', '2025-02-18 00:00:00', '2025-02-23 00:00:00', '2025-02-24 00:00:00', '2025-03-06 00:00:00', '112233', NULL, NULL, NULL, 'auditee: 397; auditor1: 43; auditor2: 151;', NULL, '2025-02-04 15:26:52', '2025-04-23 09:46:30'),
        (545, '72f04d36-4ede-4494-adb4-e59560232c7e', '2024', 0, 75, '2025-02-05 00:00:00', '2025-02-17 00:00:00', '2025-02-18 00:00:00', '2025-02-23 00:00:00', '2025-02-24 00:00:00', '2025-03-04 00:00:00', '332211', 385, 455, 421, NULL, NULL, '2025-02-04 15:30:44', '2025-03-07 11:19:52'),
        (548, '3dd3816f-e493-4ebe-9341-b629ba1a970f', '2024', 0, 68, '2025-02-05 00:00:00', '2025-02-17 00:00:00', '2025-02-18 00:00:00', '2025-02-23 00:00:00', '2025-02-24 00:00:00', '2025-03-04 00:00:00', '123456', 403, 70, 178, NULL, NULL, '2025-02-04 15:31:23', '2025-03-03 09:18:08'),
        (551, 'd9e3c6a8-999e-41ef-aaba-197ebdf96851', '2024', 0, 71, '2025-02-05 00:00:00', '2025-02-17 00:00:00', '2025-02-18 00:00:00', '2025-02-23 00:00:00', '2025-02-24 00:00:00', '2025-03-04 00:00:00', '654321', 391, 452, 415, 'data bukan capaian mengenai publikasi penelitian', NULL, '2025-02-04 15:34:23', '2025-03-04 15:02:47'),
        (554, '1e51d30d-05a6-4bdb-9c4f-c9024607e466', '2024', 0, 41, '2025-02-05 00:00:00', '2025-02-17 00:00:00', '2025-02-18 00:00:00', '2025-03-05 00:00:00', '2025-03-06 00:00:00', '2025-03-08 00:00:00', '123456', 406, 70, 178, NULL, NULL, '2025-02-05 11:36:41', '2025-03-06 13:23:14'),
        (556, 'eaa6c3a9-bd95-408a-ac8b-9171a2c0524b', '2025', 0, 91, '2025-10-21 00:00:00', '2025-11-28 00:00:00', '2025-11-29 00:00:00', '2025-12-07 00:00:00', '2025-12-08 00:00:00', '2025-12-20 00:00:00', '586299', 407, 64, 169, NULL, NULL, '2025-09-25 09:34:12', '2025-11-28 13:38:02'),
        (557, 'be3879b5-cbbe-495d-ad9a-ff05ffb21605', '2025', 0, 94, '2025-10-21 00:00:00', '2025-11-28 00:00:00', '2025-11-29 00:00:00', '2025-12-07 00:00:00', '2025-12-08 00:00:00', '2025-12-20 00:00:00', '778815', 241, 55, 421, NULL, NULL, '2025-09-25 09:35:19', '2025-11-28 13:34:23'),
        (558, '057b8aed-04bd-413a-a9e3-e5d7b7c3fdd6', '2025', 0, 97, '2025-10-21 00:00:00', '2025-11-28 00:00:00', '2025-11-29 00:00:00', '2025-12-07 00:00:00', '2025-12-08 00:00:00', '2025-12-20 00:00:00', '886063', 253, 58, 498, NULL, NULL, '2025-09-25 09:40:45', '2025-11-28 14:00:19'),
        (559, '4902d0ed-6587-4d30-ad80-1d1ab681fa6b', '2025', 0, 100, '2025-10-21 00:00:00', '2025-11-28 00:00:00', '2025-11-29 00:00:00', '2025-12-07 00:00:00', '2025-12-08 00:00:00', '2025-12-20 00:00:00', '292052', 274, 486, 465, NULL, NULL, '2025-09-25 09:42:07', '2025-11-28 14:13:50'),
        (560, '614d8eca-83e9-4a99-9711-fc4904a3ad15', '2025', 0, 103, '2025-10-21 00:00:00', '2025-11-28 00:00:00', '2025-11-29 00:00:00', '2025-12-07 00:00:00', '2025-12-08 00:00:00', '2025-12-20 00:00:00', '764147', 289, 46, 463, NULL, NULL, '2025-09-25 09:43:05', '2025-11-28 13:57:10'),
        (561, '3922cd16-9b4d-4890-8a29-e57ef81ac74c', '2025', 0, 106, '2025-10-21 00:00:00', '2025-11-28 00:00:00', '2025-11-29 00:00:00', '2025-12-07 00:00:00', '2025-12-08 00:00:00', '2025-12-20 00:00:00', '638638', 307, 79, 421, NULL, NULL, '2025-09-25 09:45:59', '2025-11-18 13:46:40'),
        (562, '6e84379c-9b4e-4afa-847e-d1ac5b5921eb', '2025', 0, 109, '2025-10-21 00:00:00', '2025-11-28 00:00:00', '2025-11-29 00:00:00', '2025-12-07 00:00:00', '2025-12-08 00:00:00', '2025-12-20 00:00:00', '604191', 325, 55, 463, NULL, NULL, '2025-09-25 09:47:08', '2025-11-18 13:35:22'),
        (563, '63469e0f-25e5-4efc-b901-1e3cf3f377f3', '2025', 0, 112, '2025-10-21 00:00:00', '2025-11-28 00:00:00', '2025-11-29 00:00:00', '2025-12-07 00:00:00', '2025-12-08 00:00:00', '2025-12-20 00:00:00', '206624', 355, 462, 465, NULL, NULL, '2025-09-25 09:48:00', '2025-11-18 13:26:42'),
        (564, '47c942dc-d7ea-43c2-afba-450af5204092', '2025', 0, 1, '2025-10-18 00:00:00', '2025-11-28 00:00:00', '2025-11-29 00:00:00', '2025-12-07 00:00:00', '2025-12-08 00:00:00', '2025-12-20 00:00:00', '184841', 235, 70, 505, NULL, 'OKE', '2025-09-25 10:53:18', '2025-11-18 13:51:43'),
        (565, '1fee4b5c-c76e-4690-a588-4c756c5a5e97', '2025', 0, 2, '2025-10-21 00:00:00', '2025-11-28 00:00:00', '2025-11-29 00:00:00', '2025-12-07 00:00:00', '2025-12-08 00:00:00', '2025-12-20 00:00:00', '125484', 247, 452, 502, NULL, NULL, '2025-09-25 10:54:54', '2025-11-18 13:55:36'),
        (566, '6ac7dffd-3ac9-44fa-82ef-15e51ab0c92f', '2025', 0, 3, '2025-10-21 00:00:00', '2025-11-28 00:00:00', '2025-11-29 00:00:00', '2025-12-07 00:00:00', '2025-12-08 00:00:00', '2025-12-20 00:00:00', '811131', 250, 462, 507, NULL, NULL, '2025-09-25 10:58:35', '2025-11-28 13:37:07'),
        (567, '803dc594-f2fa-4de9-8b5f-f4e3b8b1c980', '2025', 0, 4, '2025-10-21 00:00:00', '2025-11-28 00:00:00', '2025-11-29 00:00:00', '2025-12-07 00:00:00', '2025-12-08 00:00:00', '2025-12-20 00:00:00', '811131', 244, 91, 506, NULL, NULL, '2025-09-25 11:00:01', '2025-11-18 13:52:40'),
        (572, 'f706b941-f4db-4d4c-86bf-701b0860a588', '2025', 0, 10, '2025-10-21 00:00:00', '2025-11-28 00:00:00', '2025-11-29 00:00:00', '2025-12-07 00:00:00', '2025-12-08 00:00:00', '2025-12-20 00:00:00', '669114', 256, 82, 208, NULL, NULL, '2025-09-27 14:31:27', '2025-11-28 14:02:12'),
        (573, '4e239a31-b9bc-400e-80e1-58745bb9c6c6', '2025', 0, 11, '2025-10-21 00:00:00', '2025-11-28 00:00:00', '2025-11-29 00:00:00', '2025-12-07 00:00:00', '2025-12-08 00:00:00', '2025-12-20 00:00:00', '527234', 259, 76, 512, NULL, NULL, '2025-09-27 14:32:05', '2025-11-28 14:10:51'),
        (574, '3fbc7f4f-4c50-4c33-a64a-5e0655fa18f1', '2025', 0, 5, '2025-10-21 00:00:00', '2025-11-28 00:00:00', '2025-11-29 00:00:00', '2025-12-07 00:00:00', '2025-12-08 00:00:00', '2025-12-20 00:00:00', '245259', 262, 115, 509, NULL, NULL, '2025-09-27 14:32:36', '2025-11-18 13:57:09'),
        (575, '326f02e4-5c74-4456-b99f-a89e3b15fbe5', '2025', 0, 7, '2025-10-21 00:00:00', '2025-11-28 00:00:00', '2025-11-29 00:00:00', '2025-12-07 00:00:00', '2025-12-08 00:00:00', '2025-12-20 00:00:00', '303918', 265, 46, 178, NULL, NULL, '2025-09-27 14:33:36', '2025-11-28 14:07:21'),
        (576, '3a0013af-07f6-47f0-9a28-251b8c396abd', '2025', 0, 9, '2025-10-21 00:00:00', '2025-11-28 00:00:00', '2025-11-29 00:00:00', '2025-12-07 00:00:00', '2025-12-08 00:00:00', '2025-12-20 00:00:00', '577055', 271, 46, 465, NULL, NULL, '2025-09-27 14:34:30', '2025-11-28 14:13:26'),
        (578, '4246a864-c10b-40f3-a44c-431fb0aa49fb', '2025', 0, 14, '2025-10-21 00:00:00', '2025-11-28 00:00:00', '2025-11-29 00:00:00', '2025-12-07 00:00:00', '2025-12-08 00:00:00', '2025-12-20 00:00:00', '685677', 277, 91, 178, NULL, NULL, '2025-09-27 14:46:40', '2025-11-28 14:16:31'),
        (579, '1f4508e1-c526-47ec-85ef-76b9ad7fb30c', '2025', 0, 15, '2025-10-21 00:00:00', '2025-11-28 00:00:00', '2025-11-29 00:00:00', '2025-12-07 00:00:00', '2025-12-08 00:00:00', '2025-12-20 00:00:00', '931316', 280, 88, 488, NULL, NULL, '2025-09-27 14:47:08', '2025-11-28 14:15:58'),
        (580, 'e588c27c-b3f5-4d9c-a40d-7253aa3b961f', '2025', 0, 13, '2025-10-21 00:00:00', '2025-11-28 00:00:00', '2025-11-29 00:00:00', '2025-12-07 00:00:00', '2025-12-08 00:00:00', '2025-12-20 00:00:00', '823401', 283, 61, 488, NULL, NULL, '2025-09-27 14:47:41', '2025-11-18 13:16:06'),
        (581, 'e2f247cb-340e-4a8e-8cad-2dd8505d0d32', '2025', 0, 12, '2025-10-21 00:00:00', '2025-11-28 00:00:00', '2025-11-29 00:00:00', '2025-12-07 00:00:00', '2025-12-08 00:00:00', '2025-12-20 00:00:00', '374832', 286, 486, 490, NULL, NULL, '2025-09-27 14:48:14', '2025-11-28 14:15:34'),
        (582, 'a962c6e9-b07c-427c-b624-bc535899bf48', '2025', 0, 19, '2025-10-21 00:00:00', '2025-11-28 00:00:00', '2025-11-29 00:00:00', '2025-12-07 00:00:00', '2025-12-08 00:00:00', '2025-12-20 00:00:00', '306337', 292, 64, 463, NULL, NULL, '2025-09-30 10:23:23', '2025-11-28 13:58:06'),
        (583, 'd04e3ac2-b664-4c17-9332-63fab147db9c', '2025', 0, 20, '2025-10-21 00:00:00', '2025-11-28 00:00:00', '2025-11-29 00:00:00', '2025-12-07 00:00:00', '2025-12-08 00:00:00', '2025-12-20 00:00:00', '800823', 295, 70, 208, NULL, NULL, '2025-09-30 12:55:07', '2025-11-28 13:58:44'),
        (584, '0c7ea967-ae78-45c9-bbd3-ef74b4358d90', '2025', 0, 17, '2025-10-21 00:00:00', '2025-11-28 00:00:00', '2025-11-29 00:00:00', '2025-12-07 00:00:00', '2025-12-08 00:00:00', '2025-12-20 00:00:00', '587663', 298, 462, 514, NULL, NULL, '2025-09-30 12:55:54', '2025-11-18 14:01:18'),
        (585, 'f57d3249-1514-4441-ba77-0cb85932066f', '2025', 0, 16, '2025-10-21 00:00:00', '2025-11-28 00:00:00', '2025-11-29 00:00:00', '2025-12-07 00:00:00', '2025-12-08 00:00:00', '2025-12-20 00:00:00', '625066', 301, 64, 515, NULL, NULL, '2025-09-30 12:56:27', '2025-11-28 13:59:24'),
        (586, '1e75d5ec-5c4f-43fa-ba05-5b64be1e3986', '2025', 0, 18, '2025-10-21 00:00:00', '2025-11-28 00:00:00', '2025-11-29 00:00:00', '2025-12-07 00:00:00', '2025-12-08 00:00:00', '2025-12-20 00:00:00', '318469', 304, 76, 513, NULL, NULL, '2025-09-30 12:57:00', '2025-11-28 13:57:32'),
        (587, 'ea77c7e4-99f7-438f-a47c-7804bb8f5ccd', '2025', 0, 22, '2025-10-21 00:00:00', '2025-11-28 00:00:00', '2025-11-29 00:00:00', '2025-12-07 00:00:00', '2025-12-08 00:00:00', '2025-12-20 00:00:00', '741710', 310, 79, 501, NULL, NULL, '2025-09-30 12:58:21', '2025-11-18 13:47:01'),
        (588, '6bce3b3a-73d2-4ffc-b648-68f536f7479f', '2025', 0, 23, '2025-10-21 00:00:00', '2025-11-28 00:00:00', '2025-11-29 00:00:00', '2025-12-07 00:00:00', '2025-12-08 00:00:00', '2025-12-20 00:00:00', '473660', 313, 43, 504, NULL, NULL, '2025-09-30 12:58:52', '2025-11-28 13:55:54'),
        (589, 'bdaa0d91-ca60-4fed-84dc-772abe7dda1d', '2025', 0, 21, '2025-10-21 00:00:00', '2025-11-28 00:00:00', '2025-11-29 00:00:00', '2025-12-07 00:00:00', '2025-12-08 00:00:00', '2025-12-20 00:00:00', '698676', 316, 76, 503, NULL, NULL, '2025-09-30 12:59:17', '2025-11-28 13:54:51'),
        (590, '2389f748-e7fc-4649-9912-275591531b52', '2025', 0, 25, '2025-10-21 00:00:00', '2025-11-28 00:00:00', '2025-11-29 00:00:00', '2025-12-07 00:00:00', '2025-12-08 00:00:00', '2025-12-20 00:00:00', '255749', 319, 58, 498, NULL, NULL, '2025-09-30 12:59:59', '2025-11-28 13:55:36'),
        (591, '01f1956a-d632-4eec-8224-c7472e19c261', '2025', 0, 24, '2025-10-21 00:00:00', '2025-11-28 00:00:00', '2025-11-29 00:00:00', '2025-12-07 00:00:00', '2025-12-08 00:00:00', '2025-12-20 00:00:00', '646545', 322, 43, 502, NULL, NULL, '2025-09-30 13:00:23', '2025-11-18 13:48:24'),
        (592, '02638050-fb99-4a26-a49a-2a6c0d055602', '2025', 0, 31, '2025-10-21 00:00:00', '2025-11-28 00:00:00', '2025-11-29 00:00:00', '2025-12-07 00:00:00', '2025-12-08 00:00:00', '2025-12-20 00:00:00', '455737', 328, 462, 415, NULL, NULL, '2025-09-30 13:08:44', '2025-11-28 13:44:54'),
        (593, '37c92825-ddf1-473b-9b6b-432a3788cddf', '2025', 0, 26, '2025-10-21 00:00:00', '2025-11-28 00:00:00', '2025-11-29 00:00:00', '2025-12-07 00:00:00', '2025-12-08 00:00:00', '2025-12-20 00:00:00', '704515', 331, 76, 175, NULL, NULL, '2025-09-30 13:38:39', '2025-11-28 13:51:04'),
        (594, '4998d5e5-7c29-4044-9e42-51c23950723d', '2025', 0, 32, '2025-10-21 00:00:00', '2025-11-28 00:00:00', '2025-11-29 00:00:00', '2025-12-07 00:00:00', '2025-12-08 00:00:00', '2025-12-20 00:00:00', '556483', 334, 499, 199, NULL, NULL, '2025-09-30 13:40:29', '2025-11-18 13:38:52'),
        (595, 'af460a99-ec36-4db2-b238-2422c64cc2e6', '2025', 0, 28, '2025-10-21 00:00:00', '2025-11-28 00:00:00', '2025-11-29 00:00:00', '2025-12-07 00:00:00', '2025-12-08 00:00:00', '2025-12-20 00:00:00', '455035', 337, 82, 497, NULL, NULL, '2025-09-30 13:43:11', '2025-11-18 13:37:10'),
        (596, '5920bdd2-98dd-4a82-a0c8-022df83f6862', '2025', 0, 29, '2025-10-21 00:00:00', '2025-11-28 00:00:00', '2025-11-29 00:00:00', '2025-12-07 00:00:00', '2025-12-08 00:00:00', '2025-12-20 00:00:00', '205305', 340, 115, 196, NULL, NULL, '2025-09-30 13:43:48', '2025-11-18 13:36:16'),
        (597, '46878fb7-4e69-4706-875f-674c4d7cc7ab', '2025', 0, 27, '2025-10-21 00:00:00', '2025-11-28 00:00:00', '2025-11-29 00:00:00', '2025-12-07 00:00:00', '2025-12-08 00:00:00', '2025-12-20 00:00:00', '211027', 343, 55, 178, NULL, NULL, '2025-09-30 13:44:37', '2025-11-18 13:41:54'),
        (598, 'f6105e08-d4d2-43a7-80d9-c2bc2872a9f7', '2025', 0, 30, '2025-10-21 00:00:00', '2025-11-28 00:00:00', '2025-11-29 00:00:00', '2025-12-07 00:00:00', '2025-12-08 00:00:00', '2025-12-20 00:00:00', '599760', 346, 88, 415, NULL, NULL, '2025-09-30 13:46:24', '2025-11-28 13:45:46'),
        (599, 'dbdbd1bc-277d-4f45-80eb-128cba86ef28', '2025', 0, 35, '2025-10-21 00:00:00', '2025-11-28 00:00:00', '2025-11-29 00:00:00', '2025-12-07 00:00:00', '2025-12-08 00:00:00', '2025-12-20 00:00:00', '235117', 349, 43, 465, NULL, NULL, '2025-09-30 13:49:22', '2025-11-18 13:38:22'),
        (600, '9089e2f9-2f79-4dc0-aaae-8d58e9972b5f', '2025', 0, 33, '2025-10-21 00:00:00', '2025-11-28 00:00:00', '2025-11-29 00:00:00', '2025-12-07 00:00:00', '2025-12-08 00:00:00', '2025-12-20 00:00:00', '364807', 352, 462, 500, NULL, NULL, '2025-09-30 13:49:55', '2025-11-18 13:39:53'),
        (601, '5ac7884c-f136-4ec8-8ca8-56245d06a473', '2025', 0, 40, '2025-10-21 00:00:00', '2025-11-28 00:00:00', '2025-11-29 00:00:00', '2025-12-07 00:00:00', '2025-12-08 00:00:00', '2025-12-20 00:00:00', '505079', 358, 70, 492, NULL, NULL, '2025-09-30 13:50:48', '2025-11-18 13:27:12'),
        (602, '882a4ee6-7ed1-4dd5-b5c2-c90c777e4d60', '2025', 0, 38, '2025-10-21 00:00:00', '2025-11-28 00:00:00', '2025-11-29 00:00:00', '2025-12-07 00:00:00', '2025-12-08 00:00:00', '2025-12-20 00:00:00', '386326', 364, 61, 495, NULL, NULL, '2025-09-30 13:51:19', '2025-11-18 13:29:21'),
        (603, '6bdb065b-25b4-4b01-bab2-4bf1d6c2e035', '2025', 0, 39, '2025-10-21 00:00:00', '2025-11-28 00:00:00', '2025-11-29 00:00:00', '2025-12-07 00:00:00', '2025-12-08 00:00:00', '2025-12-20 00:00:00', '489581', 367, 70, 496, NULL, NULL, '2025-09-30 13:52:25', '2025-11-28 13:32:40'),
        (604, 'fd6d5b74-05bf-4d29-9316-fb8a0c1dbfe6', '2025', 0, 36, '2025-10-21 00:00:00', '2025-11-28 00:00:00', '2025-11-29 00:00:00', '2025-12-07 00:00:00', '2025-12-08 00:00:00', '2025-12-20 00:00:00', '006830', 370, 67, 493, NULL, NULL, '2025-09-30 13:52:55', '2025-11-18 13:27:35'),
        (605, '763afb3d-5aea-4701-b614-d6cc7c36a517', '2025', 0, 37, '2025-10-21 00:00:00', '2025-11-28 00:00:00', '2025-11-29 00:00:00', '2025-12-07 00:00:00', '2025-12-08 00:00:00', '2025-12-20 00:00:00', '091574', 373, 452, 493, NULL, NULL, '2025-09-30 13:53:30', '2025-11-18 13:32:40'),
        (606, 'e01126da-f031-452d-b6be-cf24a31ab371', '2025', 0, 121, '2025-10-21 00:00:00', '2025-11-28 00:00:00', '2025-11-29 00:00:00', '2025-12-07 00:00:00', '2025-12-08 00:00:00', '2025-12-20 00:00:00', '014163', 466, 82, 494, NULL, NULL, '2025-09-30 13:54:02', '2025-11-18 13:28:36'),
        (608, 'e99080a1-1005-4a3f-8d87-b9e49aea5d9d', '2025', 0, 41, '2025-10-01 00:00:00', '2025-10-01 00:00:00', '2025-11-01 00:00:00', '2025-11-01 00:00:00', '2025-11-02 00:00:00', '2025-11-30 00:00:00', '123456', 406, 70, 178, NULL, NULL, '2025-10-15 09:30:11', '2025-10-22 16:48:20'),
        (609, '2a3df19f-8bef-45fd-89a1-c64b0e2f2625', '2025', 0, 127, '2025-10-21 00:00:00', '2025-11-28 00:00:00', '2025-11-29 00:00:00', '2025-12-07 00:00:00', '2025-12-08 00:00:00', '2025-12-20 00:00:00', '111111', 468, 64, 501, NULL, NULL, '2025-10-21 08:17:06', '2025-11-18 13:50:52'),
        (610, 'bb687877-bb64-4945-b3e9-9012e69d787e', '2025', 0, 126, '2025-10-21 00:00:00', '2025-11-28 00:00:00', '2025-11-29 00:00:00', '2025-12-07 00:00:00', '2025-12-08 00:00:00', '2025-12-20 00:00:00', '123123', 467, 115, 508, NULL, NULL, '2025-10-22 14:04:33', '2025-11-28 13:54:03'),
        (611, '756e3cd5-d62f-42ee-a6f6-ac97c7ef3442', '2025', NULL, 6, '2025-10-18 00:00:00', '2025-11-28 00:00:00', '2025-11-29 00:00:00', '2025-12-07 00:00:00', '2025-12-08 00:00:00', '2025-12-20 00:00:00', '321123', 268, 82, 511, NULL, NULL, '2025-11-18 13:58:42', '2025-11-28 14:11:57'),
        (619, '2e2d4672-cf8f-4578-8a85-992b5bce5f9f', '2033', NULL, 10, '2025-12-01 00:00:00', '2025-12-01 00:00:00', '2025-12-02 00:00:00', '2025-12-02 00:00:00', '2025-12-03 00:00:00', '2025-12-03 00:00:00', NULL, 515, 514, 513, NULL, NULL, NULL, NULL);

    `).Error

	if err != nil {
		t.Fatalf("migration failed: %v", err)
	}
}
