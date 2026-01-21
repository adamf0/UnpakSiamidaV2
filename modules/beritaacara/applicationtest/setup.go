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

func setupBeritaAcaraMySQL(t *testing.T) (*gorm.DB, func()) {
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
		DROP TABLE IF EXISTS users;
        CREATE TABLE users (
            id bigint(20) UNSIGNED PRIMARY KEY AUTO_INCREMENT NOT NULL,
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
            id int(11) PRIMARY KEY AUTO_INCREMENT NOT NULL,
            uuid varchar(36) DEFAULT NULL,
            kode_fakultas char(9) DEFAULT NULL,
            kode_prodi char(10) DEFAULT NULL,
            nama varchar(100) DEFAULT NULL,
            id_m_prodi int(11) DEFAULT NULL,
            standalone tinyint(4) DEFAULT 0
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

		DROP TABLE IF EXISTS berita_acara;
		CREATE TABLE berita_acara (
			id int(11) PRIMARY KEY AUTO_INCREMENT NOT NULL,
			uuid varchar(36) DEFAULT NULL,
			tahun year(4) NOT NULL,
			fakultas_unit_old int(11) DEFAULT NULL,
			fakultas_unit int(11) DEFAULT NULL,
			tanggal date NOT NULL,
			auditee bigint(20) UNSIGNED DEFAULT NULL,
			auditor1 bigint(20) UNSIGNED DEFAULT NULL,
			auditor2 bigint(20) UNSIGNED DEFAULT NULL,
			created_at timestamp NULL DEFAULT NULL,
			updated_at timestamp NULL DEFAULT NULL
		);
    `).Error

	if err != nil {
		t.Fatalf("migration failed: %v", err)
	}

	seedAllBeritaAcara(t, gdb)

	cleanup := func() {
		sqlDB, _ := gdb.DB()
		resetDBBeritaAcara(t, gdb)
		sqlDB.Close()
		// mysqlC.Terminate(ctx)
	}

	return gdb, cleanup
}

func resetDBBeritaAcara(t *testing.T, gdb *gorm.DB) {
	gdb.Exec("SET FOREIGN_KEY_CHECKS=0")

	tables := []string{
		"m_fakultas",
		"m_program_studi",
		"sijamu_fakultas_unit",
		"users",
		"berita_acara",
	}

	for _, tbl := range tables {
		gdb.Exec("TRUNCATE TABLE " + tbl)
	}

	gdb.Exec("SET FOREIGN_KEY_CHECKS=1")

	seedAllBeritaAcara(t, gdb)
}

func resetDBOnlyBeritaAcara(t *testing.T, gdb *gorm.DB) {
	gdb.Exec("SET FOREIGN_KEY_CHECKS=0")

	tables := []string{
		"m_fakultas",
		"m_program_studi",
		"sijamu_fakultas_unit",
		"users",
		"berita_acara",
	}

	for _, tbl := range tables {
		gdb.Exec("TRUNCATE TABLE " + tbl)
	}

	gdb.Exec("SET FOREIGN_KEY_CHECKS=1")
}

func seedAllBeritaAcara(t *testing.T, gdb *gorm.DB) {
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

		INSERT INTO berita_acara (id, uuid, tahun, fakultas_unit_old, fakultas_unit, tanggal, auditee, auditor1, auditor2, created_at, updated_at) VALUES
		(16, '14212231-792f-4935-bb1c-9a38695a4b6b', '2023', NULL, 9, '2023-10-28', 271, 376, 151, '2023-10-27 13:42:59', '2023-10-27 13:42:59'),
		(19, 'ba9f4f89-41bd-4737-acf4-48554ed9210e', '2023', NULL, 2, '2023-10-28', 247, 118, 151, '2023-10-27 13:58:26', '2023-10-27 13:58:26'),
		(22, 'ccab56fa-c290-4497-8248-657937802ae0', '2023', NULL, 3, '2023-12-20', 250, 70, 211, '2023-10-28 03:43:34', '2023-12-20 03:40:48'),
		(28, '62039c56-77c4-4532-afa2-665be559a966', '2023', NULL, 12, '2023-12-18', 286, 67, 208, '2023-12-18 03:16:16', '2023-12-18 03:16:16'),
		(31, '10f63b40-4235-45e5-9333-2917624b7210', '2023', NULL, 100, '2023-12-18', 274, 46, 154, '2023-12-18 04:30:41', '2023-12-18 04:30:41'),
		(37, '19eac91d-8c6b-4514-a62d-10da1a6e1c49', '2023', NULL, 13, '2023-12-18', 283, 124, 217, '2023-12-18 09:38:49', '2023-12-18 09:38:49'),
		(40, 'ab8e52ab-53c4-4fd2-b0f8-0bfe8c34f780', '2023', NULL, 112, '2023-12-19', 355, 55, 163, '2023-12-18 21:30:58', '2023-12-18 21:30:58'),
		(43, '67f2fe94-95e5-4fd9-896c-abf38900716b', '2023', NULL, 25, '2023-12-19', 319, 46, 151, '2023-12-18 23:27:38', '2023-12-18 23:27:38'),
		(46, '7990715d-247e-4e35-b567-c8a0091aa9c5', '2023', NULL, 91, '2023-12-20', 232, 82, 190, '2023-12-19 02:21:49', '2023-12-19 02:21:49'),
		(49, '11b56c9a-3236-4db3-b954-fab8f4b38472', '2023', NULL, 36, '2023-12-19', 370, 43, 202, '2023-12-19 02:31:04', '2023-12-19 02:31:04'),
		(52, '243d50a5-9edf-4bcc-bd6f-5eb2eb8bc4b1', '2023', NULL, 2, '2023-12-20', 247, 97, 202, '2023-12-19 08:09:31', '2023-12-19 08:09:31'),
		(55, 'e2225d01-e834-4370-ae02-12186c90dd9d', '2023', NULL, 40, '2023-12-19', 358, 76, 154, '2023-12-19 14:46:12', '2023-12-19 14:46:12'),
		(58, '505c5699-354b-425d-9193-de3073cea616', '2023', NULL, 4, '2023-12-20', 244, 121, 181, '2023-12-20 08:02:37', '2023-12-20 08:02:37'),
		(61, 'a736389e-b8ee-4bc4-8015-bfd2b4ced3a3', '2023', NULL, 1, '2023-12-20', 235, 70, 166, '2023-12-20 08:18:22', '2023-12-20 08:18:22'),
		(64, '645de2aa-1436-45f4-981a-825a6a80b87b', '2023', NULL, 14, '2023-12-18', 277, 70, 196, '2023-12-20 08:18:58', '2023-12-20 08:18:58'),
		(67, '091c480a-7a01-43b5-8519-c9961b0d4f78', '2023', NULL, 94, '2023-12-20', 241, 76, 184, '2023-12-20 12:13:49', '2023-12-20 12:13:49'),
		(70, 'ad9b8941-2d8f-4923-9e72-35ca2f96f12f', '2023', NULL, 22, '2023-12-21', 310, 124, 175, '2023-12-21 01:35:51', '2023-12-21 01:35:51'),
		(73, '36963565-2ede-4969-8d81-2e1334e31997', '2023', NULL, 24, '2023-12-21', 322, 61, 226, '2023-12-21 02:09:01', '2023-12-21 02:09:01'),
		(76, '6030733b-13d5-45ae-9fa4-6efe14a11bf9', '2023', NULL, 21, '2023-12-21', 316, 73, 208, '2023-12-21 04:43:41', '2023-12-21 04:43:41'),
		(79, '4133efe2-f8a3-46d2-9c8d-def75ac1e4b4', '2023', NULL, 106, '2023-12-21', 307, 76, 184, '2023-12-22 08:23:53', '2023-12-22 08:23:53'),
		(82, '046b042f-17ef-4a61-b3ec-40e54f31777f', '2023', NULL, 32, '2023-12-27', 334, 91, 181, '2023-12-23 01:09:54', '2023-12-27 08:46:29'),
		(85, '37ca9fbe-c430-4965-9193-4c45cfe578e4', '2023', NULL, 97, '2023-12-23', 253, 67, 175, '2023-12-23 02:28:40', '2023-12-23 02:28:40'),
		(88, '1e714579-d8c5-463f-9681-8a6b622e48a8', '2023', NULL, 11, '2023-12-23', 259, 58, 220, '2023-12-23 03:53:32', '2023-12-23 03:53:32'),
		(91, 'f67e5d42-3b12-470d-a2c5-726ef44e5cac', '2023', NULL, 6, '2023-12-23', 268, 106, 175, '2023-12-23 06:32:29', '2023-12-23 06:32:29'),
		(94, '7675f40b-4148-424b-b6e9-23f30b7ada85', '2023', NULL, 39, '2023-12-27', 367, 97, 172, '2023-12-25 01:18:25', '2023-12-25 01:18:25'),
		(97, '18f94ec8-ae1e-4d61-b82d-794bfc00ba85', '2023', NULL, 109, '2023-12-27', 325, 61, 169, '2023-12-27 02:40:53', '2023-12-27 02:45:06'),
		(100, 'ffe855dd-7a31-4911-9ac3-ce404e1660b8', '2023', NULL, 20, '2023-12-27', 295, 121, 217, '2023-12-27 02:44:51', '2023-12-27 02:44:51'),
		(103, '5a0f899d-01de-462b-93da-295c708265b3', '2023', NULL, 28, '2023-12-28', 337, 46, 175, '2023-12-27 02:54:26', '2023-12-27 02:54:26'),
		(106, 'f799d4aa-d655-443f-aae3-e6047757bee9', '2023', NULL, 28, '2023-12-27', 337, 43, 175, '2023-12-27 02:56:24', '2023-12-27 02:56:24'),
		(118, '4b4e3049-2210-48dd-a5e7-b58b8558c2b8', '2023', NULL, 15, '2023-12-18', 280, 115, 160, '2023-12-27 03:55:10', '2023-12-27 03:55:10'),
		(121, '950807bc-b278-4135-9e22-3033dd964d97', '2023', NULL, 26, '2023-12-27', 331, 43, 187, '2023-12-27 06:54:21', '2023-12-27 06:54:21'),
		(124, '59642265-b7e5-4f47-b76a-57f806cb33f4', '2023', NULL, 17, '2023-12-28', 298, 127, 172, '2023-12-27 08:12:53', '2023-12-31 08:18:51'),
		(127, '9d18c470-01c0-4aab-a453-55a1b601caae', '2023', NULL, 33, '2023-12-27', 352, 43, 178, '2023-12-27 08:55:28', '2023-12-27 08:55:28'),
		(130, 'e553cf70-7060-4c9a-a2b1-3cd07ae063ce', '2023', NULL, 103, '2023-12-28', 289, 70, NULL, '2023-12-28 01:08:54', '2023-12-28 01:08:54'),
		(133, '16cc7f98-4e56-49de-ae60-47c0033bd53b', '2023', NULL, 103, '2023-12-28', 289, 70, NULL, '2023-12-28 01:10:00', '2023-12-28 01:10:00'),
		(136, 'e5d846aa-1612-4c9c-a437-d3626ae3a8f1', '2023', NULL, 25, '2023-12-28', 319, 43, 208, '2023-12-28 03:03:02', '2023-12-28 03:03:02'),
		(139, '6f09ab15-ec72-496c-b3af-525f342f932f', '2023', NULL, 103, '2023-12-28', 289, 70, 178, '2023-12-28 06:01:39', '2023-12-28 06:01:39'),
		(142, '3ed6a5a2-69df-4e4f-8e3f-9cae77b15456', '2023', NULL, 16, '2023-12-28', 301, 85, 178, '2023-12-28 08:26:03', '2023-12-28 08:26:03'),
		(145, 'b549076f-1016-4cf5-a9c8-4c9952f62308', '2023', NULL, 18, '2023-12-28', 304, 127, 166, '2023-12-28 09:07:06', '2023-12-29 06:32:05'),
		(148, '49fecfdc-d1a7-40e1-bdd8-93e29929f7af', '2023', NULL, 10, '2023-12-28', 256, 79, 223, '2023-12-28 09:44:34', '2023-12-28 09:44:34'),
		(151, 'cd6abee2-2203-498b-b60e-b4d75203e9f9', '2023', NULL, 10, '2023-12-28', 256, 79, 223, '2023-12-28 09:45:00', '2023-12-28 09:45:00'),
		(154, 'afec5583-0dd3-4018-a747-1a52ede238d9', '2023', NULL, 19, '2023-12-29', 292, 85, 175, '2023-12-29 02:38:26', '2023-12-29 02:38:26'),
		(157, 'bef81a79-da46-40b0-8a8e-a3c40f7a972e', '2023', NULL, 38, '2023-12-19', 364, 73, 211, '2023-12-30 03:38:07', '2023-12-30 03:38:07'),
		(160, '57f87db6-14dd-49e8-b568-c55773e7a111', '2023', NULL, 23, '2024-01-24', 313, 82, 217, '2024-01-01 08:06:11', '2024-01-01 08:09:09'),
		(163, '8a6d823e-bbcc-405c-9b9d-679bcf1e43c6', '2023', NULL, 7, '2023-12-29', 265, 82, 187, '2024-01-01 08:11:10', '2024-01-16 04:53:25'),
		(166, '5ee0381c-3381-4b7d-826d-97f18d8093b2', '2023', NULL, 37, '2023-12-19', 373, 91, 151, '2024-01-02 01:25:56', '2024-01-02 01:46:02'),
		(167, 'e514ab9b-0f0b-436e-8ef0-5efcd8b9f9c7', '2023', NULL, 35, '2023-12-27', 349, 64, 157, '2024-01-30 04:03:02', '2024-01-30 04:03:02'),
		(170, 'a1e04ee0-d2ac-4a4f-ac4d-f5fd455a9a88', '2023', NULL, 27, '2023-12-28', 343, 106, 181, '2024-01-30 04:05:27', '2024-01-30 04:05:27'),
		(173, '90487851-c3e6-4b58-a336-127813d1b6de', '2023', NULL, 5, '2023-12-23', 262, 91, 178, '2024-01-30 04:08:02', '2024-01-30 04:08:02'),
		(176, '24d933cf-1fe9-409b-ad92-8939922f0030', '2023', NULL, 29, '2023-12-27', 340, 91, 178, '2024-01-30 04:12:57', '2024-01-30 04:12:57'),
		(179, 'e6a59c18-e703-4411-b92a-ea483230726b', '2023', NULL, 30, '2023-12-28', 346, 88, 160, '2024-01-30 05:09:19', '2024-01-30 05:09:19'),
		(182, '465951fc-af67-4d88-88aa-8177942a8f0a', '2023', NULL, 31, '2023-12-27', 328, 91, 178, '2024-01-30 05:10:34', '2024-01-30 05:10:34'),
		(184, '71f3c5d7-7fb6-4d10-b918-5ef0197846e4', '2024', NULL, 68, '2024-02-27', 403, NULL, 415, '2024-02-27 01:13:18', '2024-02-27 01:13:18'),
		(187, '84951ba9-05e4-41b4-b634-c2753853e591', '2024', NULL, 41, '2024-02-27', 406, NULL, 175, '2024-02-27 01:32:48', '2024-02-27 01:32:48'),
		(190, 'b8c7fd1d-2d9b-4938-b9ec-ba655c11a0db', '2023', NULL, 113, '2024-02-27', 400, NULL, 151, '2024-02-27 03:52:15', '2024-03-06 03:23:19'),
		(193, '89dd525d-e259-4d6e-9ba6-b8399ffe2864', '2024', 89, NULL, '2024-02-27', 397, NULL, 175, '2024-02-27 04:16:18', '2024-02-27 04:16:18'),
		(196, 'e6e59e2e-8df6-4d4a-90cc-727889c9a8f8', '2024', NULL, 71, '2024-02-27', 391, NULL, 415, '2024-02-27 06:40:21', '2024-02-27 06:40:21'),
		(199, '008a1317-1e23-48b6-82dd-ea51faed958f', '2024', NULL, 72, '2024-02-27', 388, NULL, 178, '2024-02-27 06:54:35', '2024-02-27 06:54:35'),
		(202, 'eb8da2c5-491b-465a-9843-c3f8f339191a', '2024', NULL, 69, '2024-02-29', 413, NULL, 163, '2024-02-29 03:40:52', '2024-02-29 03:40:52'),
		(205, 'a54d7a8e-b78a-49a4-8ce2-f5f896671f34', '2024', NULL, 119, '2024-02-29', 407, NULL, 178, '2024-02-29 06:41:38', '2024-02-29 06:41:38'),
		(208, '83ef63f1-9f22-43de-95ed-2002a57c4aa7', '2024', NULL, 70, '2024-02-29', 382, NULL, 151, '2024-03-01 08:26:25', '2024-03-01 08:26:25'),
		(211, '1b9b9219-fe5c-466a-a717-f51ed1e2b7ba', '2024', 88, 88, '2024-02-27', 394, NULL, 178, '2024-03-04 06:29:10', '2024-03-04 06:29:10'),
		(214, 'ac82f6dd-4329-4761-a4c2-324007d6584e', '2023', NULL, 116, '2024-02-29', 410, NULL, 151, '2024-03-04 08:17:49', '2024-03-04 08:17:49'),
		(217, 'a11ab2ea-8dc3-4dc4-9345-8a3ae118a41f', '2024', NULL, 75, '2024-03-05', 385, NULL, 175, '2024-03-05 06:36:18', '2024-03-05 06:36:18'),
		(220, 'b9bfe8ff-9709-4976-bf58-81ffddd4adf0', '2024', NULL, 41, '2024-09-03', 406, NULL, NULL, '2024-09-03 04:12:51', '2024-09-03 04:12:51'),
		(223, 'c5b509df-7a59-47a2-bcd9-3fd2563fb145', '2025', NULL, 41, '2024-09-08', 406, NULL, 226, '2024-09-03 04:14:32', '2024-09-03 04:14:32'),
		(226, '3ad65d2c-50c5-4855-aea6-e343329a695f', '2024', NULL, 1, '2024-09-09', 232, 40, 154, '2024-09-24 07:41:06', '2024-09-24 07:41:06'),
		(229, '7c9fdc9e-5a68-4a29-b322-5c02a8d922cb', '2024', NULL, 91, '2024-09-25', 232, 43, 415, '2024-09-25 03:33:55', '2024-09-25 03:33:55'),
		(236, 'a601a2a6-1de5-4961-a84c-d05b96569907', '2024', NULL, 25, '2024-12-17', 319, 43, 217, '2024-10-11 03:42:58', '2024-12-17 05:09:40'),
		(239, 'bcfd04b8-9a32-4720-b7b4-fa31aeb56887', '2024', NULL, 41, '2024-10-12', 425, 428, 431, '2024-10-12 04:03:45', '2024-10-12 04:03:45'),
		(242, '3e30aa21-9198-46ec-abf1-757e9bb8b422', '2024', NULL, 3, '2024-12-14', 250, 43, 211, '2024-10-15 07:04:01', '2024-12-14 05:31:10'),
		(248, 'f02b4607-ac0d-4703-ac7c-3ebdd5602194', '2024', NULL, 14, '2024-12-14', 277, 70, 199, '2024-10-18 02:50:56', '2024-12-16 07:53:42'),
		(250, 'b50dff9b-2a8c-4dd6-97b9-ec27173a4ef8', '2024', NULL, 91, '2024-11-22', 425, 428, 431, '2024-11-22 03:25:51', '2024-11-22 03:25:51'),
		(253, '8a12ad25-ed2d-47b3-87f4-b491804bc157', '2024', NULL, 12, '2024-12-14', 286, 70, 187, '2024-12-14 06:59:15', '2024-12-14 06:59:31'),
		(256, '8d974e90-7792-4710-82ae-30bc181a0d29', '2024', NULL, 4, '2024-12-14', 244, 70, 415, '2024-12-14 09:29:29', '2024-12-14 09:29:29'),
		(259, '3ef55003-dec0-434b-9a6a-5539fc2a0888', '2024', NULL, 1, '2024-12-16', 235, 70, 415, '2024-12-16 07:53:30', '2024-12-16 07:53:30'),
		(262, '06879764-a841-4561-9efc-497ff53273e6', '2024', NULL, 37, '2024-12-16', 373, 106, 199, '2024-12-16 08:26:31', '2024-12-16 08:26:31'),
		(265, '7a2c9e3f-93e4-4ff5-bbc3-d07b37a6dafb', '2024', NULL, 15, '2024-12-17', 280, 52, 217, '2024-12-17 03:52:14', '2024-12-17 03:52:14'),
		(268, '238ad9d5-00f0-4d1f-ac9c-b429915f58b0', '2024', NULL, 36, '2024-12-13', 370, 115, 415, '2024-12-17 05:13:44', '2024-12-17 05:13:44'),
		(271, '43dba6bc-a87a-4ee9-b4dd-7b43d5c47289', '2024', NULL, 21, '2024-12-17', 316, 67, 178, '2024-12-17 06:53:35', '2024-12-17 06:53:35'),
		(274, '63a781bc-05af-4782-ae4b-8c5c236d98ee', '2024', NULL, 23, '2024-12-17', 313, 76, 160, '2024-12-17 06:59:51', '2024-12-18 09:44:27'),
		(277, '68038646-35e5-44fe-8fae-91c6ca5c303e', '2024', NULL, 13, '2024-12-18', 283, 43, 421, '2024-12-18 03:01:33', '2024-12-24 02:07:00'),
		(280, 'b39ec3f1-6538-457c-a26d-44d54fa2b278', '2024', NULL, 11, '2024-12-18', 259, 76, 415, '2024-12-18 05:29:56', '2024-12-18 09:45:12'),
		(283, '324115f3-91c9-4af2-84e3-64e08188dcc5', '2024', NULL, 5, '2024-12-18', 262, 70, 175, '2024-12-18 08:57:36', '2024-12-18 08:59:39'),
		(286, '2dcbf9bc-11d9-454e-9efc-215fd4627b6a', '2024', NULL, 6, '2024-12-18', 268, 70, 151, '2024-12-18 08:58:12', '2024-12-18 08:59:48'),
		(292, 'ebc41dfd-88ee-4945-9c17-eb1abfafefd8', '2024', NULL, 9, '2024-12-18', 271, 91, 421, '2024-12-18 22:03:30', '2024-12-23 08:27:53'),
		(295, '8f63e794-bfe3-4d87-8bda-b83f6385bbd2', '2024', NULL, 109, '2024-12-19', 325, 61, 151, '2024-12-19 02:33:20', '2024-12-19 02:33:20'),
		(298, '0e61044a-d015-4ab8-b63c-b891ef2bf7b8', '2024', NULL, 27, '2024-12-19', 343, 70, 421, '2024-12-19 02:50:34', '2024-12-19 02:50:34'),
		(301, '579f4473-5f96-44a8-991c-081c85cc5fc3', '2024', NULL, 28, '2024-12-19', 337, 43, 217, '2024-12-19 04:18:29', '2024-12-19 04:18:29'),
		(304, '4e81462f-b258-4aa5-8399-2dcfe0df6751', '2024', NULL, 29, '2024-12-19', 340, 70, 223, '2024-12-19 05:30:19', '2024-12-19 05:30:19'),
		(307, '85ec16bf-ead7-42f2-a401-12efa407f55f', '2024', NULL, 31, '2024-12-19', 328, 91, 178, '2024-12-19 07:03:11', '2024-12-23 08:14:29'),
		(310, '01d1c545-4466-4c74-86d2-b3b5936ac5e2', '2024', NULL, 26, '2024-12-19', 331, 43, 226, '2024-12-19 07:56:10', '2024-12-19 08:14:33'),
		(313, '20d82302-6578-448b-bd8a-0ff27f02151f', '0000', NULL, 30, '2024-12-19', 346, 88, 169, '2024-12-19 07:58:07', '2024-12-19 07:58:07'),
		(316, 'd6e3ffce-291f-4721-9ad8-d6a01518e8af', '2024', NULL, 33, '2024-12-19', 352, 70, 415, '2024-12-19 09:16:48', '2024-12-19 09:16:48'),
		(319, 'a3720210-863c-4445-8e7f-bfff97e658e7', '2024', NULL, 24, '2024-12-17', 322, 88, 226, '2024-12-19 17:25:27', '2024-12-19 17:25:27'),
		(325, 'd1ec541c-6f53-4512-b6e8-953308c09195', '2024', NULL, 30, '2024-12-19', 346, 88, 169, '2024-12-19 17:29:05', '2024-12-19 17:29:05'),
		(328, 'a1a6604e-a118-46c7-a4cc-8382715f4e32', '2024', NULL, 18, '2024-12-20', 304, 43, 229, '2024-12-20 01:49:09', '2024-12-20 01:49:09'),
		(331, '1c7129da-9477-44a0-9c54-6d4439838e73', '2024', NULL, 103, '2024-12-20', 289, 85, 178, '2024-12-20 03:01:51', '2024-12-20 03:01:51'),
		(334, '5dde3262-c637-4bf8-b71b-c8808178138c', '2024', NULL, 94, '2024-12-20', 241, 43, 166, '2024-12-20 03:14:15', '2024-12-20 03:14:15'),
		(337, 'cd3b967e-b470-4602-ad0d-d7fa86fc81fe', '2024', NULL, 10, '2024-12-20', 256, 79, 217, '2024-12-20 04:12:32', '2024-12-23 08:35:49'),
		(340, 'c91cd4fb-ef3b-43b5-ab45-54353e32bab9', '2024', NULL, 20, '2024-12-20', 295, 70, 223, '2024-12-20 04:54:22', '2024-12-20 04:54:22'),
		(343, '89d26a1f-033d-414a-a960-6e59fd66d401', '2024', NULL, 19, '2024-12-20', 292, 85, 196, '2024-12-20 05:13:32', '2024-12-20 05:13:32'),
		(346, '0e78f5f8-96ec-4bc8-91bc-220c134d4934', '2024', NULL, 22, '2024-12-20', 310, 82, 151, '2024-12-20 07:37:05', '2024-12-20 07:37:05'),
		(349, 'b5c28bd8-7017-4835-aa9c-2cc1b87318af', '2024', NULL, 39, '2024-12-20', 367, 70, 223, '2024-12-20 07:58:28', '2024-12-20 07:58:28'),
		(352, 'bddcfb81-af99-404c-897d-d8a73370b087', '2024', NULL, 16, '2024-12-20', 301, 85, 181, '2024-12-20 22:21:31', '2024-12-20 22:21:31'),
		(355, 'ee698fee-94b6-4b3b-b5eb-a013b67ff8c9', '2024', NULL, 91, '2024-12-23', 232, 70, 415, '2024-12-23 01:38:01', '2024-12-23 01:38:01'),
		(358, '976110b1-33ee-4507-92f9-7934b8bb5029', '2024', NULL, 32, '2024-12-23', 334, 58, 421, '2024-12-23 08:52:39', '2024-12-23 08:52:39'),
		(361, '7887fc3a-fd50-4d6e-9baf-7311d56696af', '2024', NULL, 7, '2024-12-20', 265, 82, 415, '2024-12-23 13:31:05', '2025-01-06 02:48:32'),
		(364, 'bc9dc775-9dfe-43b9-baa2-baed69c41545', '2024', NULL, 106, '2024-12-17', 307, 43, 169, '2024-12-24 02:05:54', '2024-12-24 02:05:54'),
		(367, 'd8ca61e6-e894-40f7-b09c-4fb38a9eb522', '2024', NULL, 40, '2024-12-21', 358, 43, 178, '2024-12-24 02:08:42', '2024-12-24 02:08:42'),
		(370, '971eda05-b10c-4893-b101-f2a6e5b792d6', '2024', NULL, 2, '2024-12-14', 247, 64, 151, '2024-12-24 07:50:59', '2024-12-24 07:50:59'),
		(373, 'a6fb3381-3ab8-434f-9f25-0e5459edcb5f', '2024', NULL, 17, '2024-12-20', 298, 64, 229, '2024-12-24 07:52:07', '2024-12-24 07:52:07'),
		(376, 'bedd2009-9db2-4353-83c3-47db72b45d3a', '2024', NULL, 35, '2024-12-19', 349, 64, 175, '2024-12-24 07:53:22', '2024-12-24 07:53:22'),
		(379, '67677b6d-b69b-4672-a52c-20f1d74f38df', '2024', NULL, 97, '2024-12-18', 253, 64, 166, '2024-12-24 07:54:49', '2024-12-24 07:54:49'),
		(382, 'dc8ad910-e53e-48ee-82c3-4b83a59700af', '2024', NULL, 38, '2024-12-19', 364, 73, 421, '2024-12-24 09:58:55', '2024-12-24 09:58:55'),
		(385, 'ffee7f6e-7e85-4b40-bcb4-b57d5cf58f49', '2024', NULL, 100, '2024-12-21', 274, 49, 223, '2025-01-07 15:57:53', '2025-01-07 15:57:53'),
		(388, 'c40f505e-1917-47f5-bfd7-24d0d2c0113f', '2024', NULL, 112, '2024-12-16', 355, 49, 421, '2025-01-08 01:11:32', '2025-01-20 04:20:32'),
		(389, '61bc77bc-0a5c-4ea0-9320-b7a2b250fb12', '2024', NULL, 41, '2025-02-05', 406, 43, 178, '2025-02-05 01:53:28', '2025-02-05 01:53:28'),
		(392, '3dc63fe3-7f5c-43f7-8e24-c375046f69fa', '2024', NULL, 125, '2025-02-27', 400, 452, 415, '2025-02-27 02:37:52', '2025-02-27 02:37:52'),
		(394, '49291b6c-b074-48e8-a1ae-e7bbc4481553', '2025', NULL, 72, '2025-03-04', 388, 70, NULL, '2025-03-04 06:11:13', '2025-03-04 06:11:13'),
		(397, 'c9f4cdf6-b68b-4175-ad42-02c05e53ee31', '2025', NULL, 68, '2025-03-04', 403, 70, NULL, '2025-03-04 06:11:39', '2025-03-04 06:11:39'),
		(400, '274a0047-e9fd-4621-9995-3397030328c3', '2024', NULL, 71, '2025-03-04', 391, 452, 415, '2025-03-04 08:08:51', '2025-03-04 08:08:51'),
		(403, '5a8c472f-7070-4ecf-9194-0493be4d22ca', '2025', NULL, 41, '2025-03-06', 406, 70, NULL, '2025-03-06 07:05:50', '2025-03-06 07:05:50'),
		(406, '4df8c0bd-d12a-4936-86f3-f7e1ecc5c4ec', '2024', NULL, 69, '2025-02-24', 413, 55, 163, '2025-03-20 04:37:07', '2025-03-20 04:37:07'),
		(409, 'f4d81e72-d021-4cba-939b-d8d22ff8a47f', '2024', NULL, 70, '2025-03-04', 382, 43, NULL, '2025-03-20 06:46:01', '2025-03-20 06:46:01'),
		(412, '3b455867-8373-4e4a-8ffd-dbd4b4555b8f', '2024', 89, NULL, '2025-03-04', 394, 43, NULL, '2025-03-20 06:46:25', '2025-03-20 06:46:25'),
		(415, '76572ff5-662f-4205-b3d1-40cc79ce2d7e', '2024', NULL, 75, '2025-03-04', 385, 455, 421, '2025-03-21 06:34:36', '2025-03-21 06:34:36'),
		(418, 'c79249d7-e25f-4d33-a981-0cc7edd69bc9', '2025', 88, 88, '2025-02-24', 394, 70, NULL, '2025-05-14 01:42:25', '2025-05-14 01:42:25'),
		(419, '1625be7d-286f-483c-a0b6-92c4c29a6e1b', '2025', NULL, 91, '2025-11-29', 232, 428, 431, '2025-11-21 02:25:53', '2025-11-21 02:25:53');
	`).Error

	if err != nil {
		t.Fatalf("migration failed: %v", err)
	}
}
