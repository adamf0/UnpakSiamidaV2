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

func setupDokumenTambahanMySQL(t *testing.T) (*gorm.DB, func()) {
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

        DROP TABLE IF EXISTS dokumen_tambahan;
        CREATE TABLE dokumen_tambahan (
            id int(11) NOT NULL,
            uuid varchar(36) DEFAULT NULL,
            id_renstra_old int(11) DEFAULT NULL,
            id_renstra int(11) NOT NULL,
            id_template_dokumen_tambahan int(11) NOT NULL,
            file text DEFAULT NULL,
            capaian_auditor varchar(100) DEFAULT NULL,
            catatan_auditor text DEFAULT NULL,
            tugas varchar(100) NOT NULL DEFAULT 'auditor2',
            created_at datetime DEFAULT NULL,
            updated_at datetime DEFAULT NULL
        );

        ALTER TABLE dokumen_tambahan
        ADD PRIMARY KEY (id);

        DROP TABLE IF EXISTS jenis_file_renstra;
        CREATE TABLE jenis_file_renstra (
            id bigint(20) UNSIGNED NOT NULL,
            uuid varchar(36) DEFAULT NULL,
            nama text NOT NULL,
            created_at datetime DEFAULT NULL,
            updated_at datetime DEFAULT NULL
        );

        ALTER TABLE jenis_file_renstra
        ADD PRIMARY KEY (id);

        DROP TABLE IF EXISTS renstra;
        CREATE TABLE renstra (
            id int(11) NOT NULL,
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

        ALTER TABLE renstra
        ADD PRIMARY KEY (id);
        
        DROP TABLE IF EXISTS template_dokumen_tambahan;
        CREATE TABLE template_dokumen_tambahan (
            id int(11) NOT NULL,
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
        ADD PRIMARY KEY (id),
        ADD UNIQUE KEY uq_template_dokumen (tahun,jenis_file,fakultas_prodi_unit);

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
		"dokumen_tambahan",
		"jenis_file_renstra",
		"renstra",
		"template_dokumen_tambahan",
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
		"dokumen_tambahan",
		"jenis_file_renstra",
		"renstra",
		"template_dokumen_tambahan",
	}

	for _, tbl := range tables {
		gdb.Exec("TRUNCATE TABLE " + tbl)
	}

	gdb.Exec("SET FOREIGN_KEY_CHECKS=1")
}

func seedAllDokumenTambahan(t *testing.T, gdb *gorm.DB) {
	err := gdb.Exec(`
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

        INSERT INTO dokumen_tambahan (id, uuid, id_renstra_old, id_renstra, id_template_dokumen_tambahan, file, capaian_auditor, catatan_auditor, tugas, created_at, updated_at) VALUES
        (1271, 'c836800f-8c09-4e04-ba16-e0ca027ca571', NULL, 368, 28, 'https://www.linkedin.com/in/adamfurqon175901204', NULL, 'entah lah', 'auditor2', NULL, '2024-12-23 07:53:29'),
        (1274, 'db252319-9c5b-4884-ba10-0a3ad9bacee9', NULL, 368, 34, 'https://www.linkedin.com/in/adamfurqon175901204', NULL, NULL, 'auditor2', NULL, '2024-12-23 07:53:40'),
        (1277, '14722f04-575e-457d-ab25-fe355d038889', NULL, 368, 40, 'SOP KINERJA STRUKTURAL.pdf', NULL, NULL, 'auditor2', NULL, '2024-12-23 07:55:04'),
        (1280, '9ab98493-b1a5-4247-985a-7fe87221f990', NULL, 368, 46, 'TEMUAN AMI 2023 YANG SUDAH DI CLOSED.pdf', '1', NULL, 'auditor2', NULL, '2024-12-23 07:55:58'),
        (1283, 'bd340cc2-48cb-4f52-812e-b5caf6b7a99e', NULL, 371, 30, 'PROKER 2024 FAKULTAS HUKUM UNPAK FINAL 31 JANUARI 2024.xlsx', '1', NULL, 'auditor2', NULL, '2024-11-26 18:17:58'),
        (1286, 'd2373f1a-9149-4831-9749-35d2e059847e', NULL, 371, 36, 'STRUKTUR ORGANISASI DAN JOB DESCRIPTION JOB SPECIFICATION FH-UNPAK.pdf', '1', NULL, 'auditor2', NULL, '2024-11-26 18:19:24'),
        (1289, 'c1b23f35-e16e-4922-879f-7c6de2cb6c81', NULL, 371, 42, 'SOP KINERJA STRUKTURAL.pdf', '1', NULL, 'auditor2', NULL, '2024-11-26 18:19:53'),
        (1292, '3cf11241-83fa-4ba9-a49e-e6b96bb935fe', NULL, 371, 48, 'TEMUAN AMI 2023 YANG SUDAH DI CLOSED.pdf', '1', NULL, 'auditor2', NULL, '2024-11-26 18:20:42'),
        (1295, '802b0732-e5b5-4852-a770-a834a8b70746', NULL, 374, 28, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (1298, '39b5f87d-2cbe-4699-8be1-bf04f174590f', NULL, 374, 34, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (1301, 'b1c690f6-381b-45ba-8675-2c1ceff5ad68', NULL, 374, 40, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (1304, 'daddd871-85d7-48fc-a950-cab013442cbd', NULL, 374, 46, 'Berita Acara Audit 2023.pdf', NULL, NULL, 'auditor2', NULL, '2024-11-14 13:38:13'),
        (1307, '7f40fa57-3fa7-4545-b734-9a3999766d50', NULL, 377, 30, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (1310, '6e2b08e5-2db3-4a12-965c-ae74f54d032b', NULL, 377, 36, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (1313, '0a1b2b32-a9d1-434a-9b37-7cc8b7d53072', NULL, 377, 42, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (1316, '0581b311-a023-46ca-a522-ae338d68df4b', NULL, 377, 48, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (1319, 'c3272370-0756-42ad-9064-197802107886', NULL, 380, 30, NULL, '1', NULL, 'auditor2', NULL, '2024-12-14 09:36:50'),
        (1322, '8ab1d23b-e72f-41b8-bdf5-cb19ec9aa44c', NULL, 380, 36, 'Struktur Organisasi, JOB DESC dan JOB SPEC Prodi AKT.pdf', '1', NULL, 'auditor2', NULL, '2024-11-26 14:23:56'),
        (1325, 'dc48fa18-9e85-4d32-bef1-45118d51616e', NULL, 380, 42, NULL, '1', NULL, 'auditor2', NULL, '2024-12-14 09:41:16'),
        (1328, '90e8e1eb-f833-4ebb-92d0-3c4f0aaec47c', NULL, 380, 48, 'berita_acara.pdf', '1', NULL, 'auditor2', NULL, '2024-12-14 10:00:28'),
        (1331, '8c38a001-64bb-4c24-a4de-2d0f6a1dc12a', NULL, 383, 30, '4. Proker Bisnis Digital 2024 (11062024).xlsx', '0', 'Pada Proker 2024: Perlu penjelasan IKU/IKT PS. Terdapat ketidaksesuaian program dan indikator.', 'auditor2', NULL, '2024-12-14 11:21:31'),
        (1334, '54889824-1eb1-45b1-8403-fbd95bcec950', NULL, 383, 36, 'BUKU JOB DESCRIPTION.pdf', '1', NULL, 'auditor2', NULL, '2024-12-11 02:05:47'),
        (1337, 'c753fc38-c99a-4b01-80bd-d660168b23b1', NULL, 383, 42, 'SOP - PROKER.pdf', '1', NULL, 'auditor2', NULL, '2024-12-11 02:13:29'),
        (1340, 'c7053b08-134d-4ed2-9b47-79c36f777c4c', NULL, 383, 48, 'Closing audit.pdf', '1', NULL, 'auditor2', NULL, '2024-12-11 02:15:50'),
        (1343, '864285ba-3b78-4aaa-bbb3-02b162af12a6', NULL, 386, 28, 'Monitoring Proker FKIP 2024.xlsx', NULL, NULL, 'auditor2', NULL, '2024-11-22 21:35:26'),
        (1346, '9d9b2b9e-4487-4096-b828-ef2bd9b14886', NULL, 386, 34, '1 Struktur Organisasi dan Deskripsi Kerja FKIP-pages.pdf', NULL, NULL, 'auditor2', NULL, '2024-11-22 21:49:51'),
        (1349, '530342cd-5b9e-4437-9bb9-b3d569a5fff9', NULL, 386, 40, '1 Struktur Organisasi dan Deskripsi Kerja FKIP-pages.pdf', NULL, NULL, 'auditor2', NULL, '2024-11-22 21:50:01'),
        (1352, '3a3447d6-2dd8-431b-9a95-4e53b5d04d6a', NULL, 386, 46, 'BA+hasil audit 2023.pdf', NULL, NULL, 'auditor2', NULL, '2024-11-22 21:50:44'),
        (1367, 'e75254d2-a641-41f1-b1af-5a9d3403f9e4', NULL, 392, 30, 'PROKER DAN MONITORING 2024.pdf', '0', 'ada 1 program kerja belum tercapai sertfikasi penulis dan editor', 'auditor2', NULL, '2024-11-29 09:38:18'),
        (1370, '7829af6d-f54a-4277-bb64-c2eb68395ccd', NULL, 392, 36, 'SK, Struktur Organisasi, dan Job Des.pdf', '1', NULL, 'auditor2', NULL, '2024-12-20 08:28:14'),
        (1373, '704da2fd-434f-4039-aab1-499e1cc8ba8b', NULL, 392, 42, 'SOP JOBDES dan SOP Pembagian Tugas.pdf', '1', NULL, 'auditor2', NULL, '2024-11-29 09:36:19'),
        (1376, '2277f90f-3a09-46d3-afc2-b2aac263f244', NULL, 392, 48, 'kts (6).pdf', '1', NULL, 'auditor2', NULL, '2024-11-29 09:38:52'),
        (1379, '6af712ec-dd60-43b4-ad63-1a2e76564398', NULL, 395, 30, 'Evaluasi Proker 2024.xlsx', '1', NULL, 'auditor2', NULL, '2024-11-26 18:29:01'),
        (1382, 'd9ecd06f-3432-4c92-b35b-a3acc1ed4703', NULL, 395, 36, 'Job Desc and Specs _SK 2024.pdf', '1', NULL, 'auditor2', NULL, '2024-11-26 18:29:52'),
        (1385, '5990a6f7-6e81-462f-8b81-908dd5088d27', NULL, 395, 42, 'KUMPULAN SOP.pdf', '1', NULL, 'auditor2', NULL, '2024-11-26 18:30:17'),
        (1388, '7f50e862-9330-4051-ab9d-6466ce56d417', NULL, 395, 48, 'kts (6).pdf', '1', NULL, 'auditor2', NULL, '2024-11-26 18:30:34'),
        (1391, '8260b8e5-590d-4427-88be-bdf7a7547021', NULL, 398, 30, 'Monitoring Proker 2024.docx', '1', NULL, 'auditor2', NULL, '2024-12-02 09:06:58'),
        (1394, '6a13b09e-2f7b-4880-8bb9-a7ed18255335', NULL, 398, 36, '1 Struktur Organisasi dan Deskripsi Kerja FKIP.pdf', '1', NULL, 'auditor2', NULL, '2024-12-09 10:12:40'),
        (1397, '72f9dcdf-c068-4237-b485-2cb9d86f3679', NULL, 398, 42, 'SOP FKIP Size kecil.pdf', '1', NULL, 'auditor2', NULL, '2024-12-02 09:07:33'),
        (1400, '2432e79e-d741-4add-9eb2-75f15b3840b3', NULL, 398, 48, 'kts (2).pdf', '1', NULL, 'auditor2', NULL, '2024-12-09 10:13:17'),
        (1403, 'f65f42af-b5a1-4f63-816a-4150774e9f8a', NULL, 401, 30, 'PROKER 2024 dan Monitoring (Lengkap).pdf', '1', NULL, 'auditor2', NULL, '2024-11-26 17:43:02'),
        (1406, 'bb2d9850-34a5-419d-b2ea-6b628334ac83', NULL, 401, 36, 'PGSD. Struktur Organisasi Prodi PGSD Job Desc & Job Spec (2020-2025).pdf', '1', NULL, 'auditor2', NULL, '2024-11-26 17:55:54'),
        (1409, 'da58b1cf-e8d7-43df-9c8e-5076b890ada5', NULL, 401, 42, 'PGSD. Struktur Organisasi Prodi PGSD Job Desc & Job Spec (2020-2025).pdf', '0', 'Dokumen yang diunggah bukan SOP namun Jobdesc dan jobspec', 'auditor2', NULL, '2024-11-26 17:57:49'),
        (1412, '7d453f87-acb6-46e1-9eb6-6dc508f2107d', NULL, 401, 48, 'Closing PGSD AMI 2023 (kts 7-12)pdf.pdf', '1', NULL, 'auditor2', NULL, '2024-11-26 17:58:09'),
        (1415, '887686f7-4fca-414d-8a68-71aa8b634c3b', NULL, 404, 30, 'Program kerja Pend. IPA 2024.xlsx', '0', 'belum terlihat ada monitoringnya', 'auditor2', NULL, '2024-12-17 20:47:47'),
        (1418, '47433430-a4b3-4ec6-85c6-16bcc54b70c6', NULL, 404, 36, 'Job Des Organisasi Prodi 2024.pdf', '1', NULL, 'auditor2', NULL, '2024-12-17 20:48:31'),
        (1421, 'edb36953-e5eb-4934-b063-811c8d72f216', NULL, 404, 42, 'SOP TUPOKSI TATA PAMONG PRODI PENDIDIKAN IPA (S1).pdf', '1', NULL, 'auditor2', NULL, '2024-12-17 20:50:00'),
        (1424, 'ea652868-d54b-43b0-b4b5-c27306139db7', NULL, 404, 48, 'CLOSED AUDIT 2023.pdf', '1', NULL, 'auditor2', NULL, '2024-12-17 20:50:29'),
        (1427, 'daa25df8-9d9d-4c45-a313-14d833fba847', NULL, 407, 30, 'Proker PPG 2024.xlsx', '0', '1. ada beberapa program belum ada evaluasinya \r\n2. jadwak proker Realisasi belum di centang', 'auditor2', NULL, '2024-12-03 12:40:24'),
        (1430, '7c1e82e2-d104-4c3f-8f90-ee20f21efd41', NULL, 407, 36, 'LEMBAR PENGESAHAN STRUKTUR ORGANISASI 2024.pdf', '0', 'Lembar Pengesahan belum ditanda tangan', 'auditor2', NULL, '2024-12-03 13:52:30'),
        (1433, '8164b574-d9ca-449b-b640-f0dfd9debf4d', NULL, 407, 42, 'SOP PROKER.docx.pdf', '1', NULL, 'auditor2', NULL, '2024-12-03 13:53:19'),
        (1436, 'dd3bda91-0f45-4eba-a4b7-4484d1d82c24', NULL, 407, 48, 'berita_acara (1).pdf', '0', 'KTS belum closed', 'auditor2', NULL, '2024-12-04 12:37:20'),
        (1439, '2e54eddd-12a8-4104-b576-bf3fc8458856', NULL, 410, 28, 'Proker FISIB 2023-2024-Monev.xlsx', '1', NULL, 'auditor2', NULL, '2024-12-21 14:53:06'),
        (1442, '746c7e58-4426-4ec4-aaaa-682c09356f28', NULL, 410, 34, 'SK DEKAN STRUKTUR ORGANISASI & TUPOKSI FISIB MASA BAKTI 2020-2025.pdf', '1', NULL, 'auditor2', NULL, '2024-12-21 14:53:13'),
        (1445, '3f5897f0-930f-4cc7-9d4b-e3a3ca0ec1c8', NULL, 410, 40, 'SOP Struktur Organisasi FISIB .pdf', '1', NULL, 'auditor2', NULL, '2024-12-21 14:53:21'),
        (1448, 'b6d9cde2-63c0-41e6-95a4-28f3513e0e6e', NULL, 410, 46, 'Berita Acara Audit 2023.pdf', '1', NULL, 'auditor2', NULL, '2024-12-21 14:54:19'),
        (1451, 'b9558948-6395-4146-b1da-9159398710df', NULL, 413, 30, 'PROKER 2024_Sastra Indonesia-2.xlsx', '1', NULL, 'auditor2', NULL, '2024-12-03 12:41:47'),
        (1454, '7a0e6405-98e1-4679-9c8b-1c5677f2380e', NULL, 413, 36, 'Struktur organisasi Sastra Indonesia.docx', '1', NULL, 'auditor2', NULL, '2024-12-03 13:56:26'),
        (1457, 'c913f73f-a10c-4cb0-9b3f-a37395c7c8f3', NULL, 413, 42, 'KUMPULAN SOP FISIB.pdf', '1', NULL, 'auditor2', NULL, '2024-12-18 09:58:15'),
        (1460, '0de5c66d-bd0b-4458-8216-a8660d968215', NULL, 413, 48, 'Closing KTS.pdf', '1', NULL, 'auditor2', NULL, '2024-12-03 13:58:00'),
        (1463, '8e0576f3-c698-4748-840d-d2668539af8b', NULL, 416, 30, 'Proker 2024_Sastra Inggris.xlsx', '1', NULL, 'auditor2', NULL, '2024-12-09 19:52:52'),
        (1466, 'fe79b329-c037-42ea-aa43-6d11b9c12a40', NULL, 416, 36, 'Struktur Organisasi, Job Desc, dan Job Spec.pdf', '1', NULL, 'auditor2', NULL, '2024-12-09 19:54:43'),
        (1469, 'ff5cfa27-2abe-4391-bdf1-d5ab2d27041b', NULL, 416, 42, 'KUMPULAN SOP FISIB.pdf', '0', 'KUMPULAN SOP Tidak mecakup semua aktivitas sesuai  job desk dan sob spek yang ada dilingkungan FISIB ,. Mohon dibuat SOP semua aktivitas yang berjalan di FISIB dalam buku yang lengkap dengan kata pengantar, pengeshaan, daftar isi serta  hal lainnya dari kumpulan SOP tersebut', 'auditor2', NULL, '2024-12-09 19:59:31'),
        (1472, 'be3649ae-8098-4000-b3bd-a9ad0cff699a', NULL, 416, 48, 'KTS 2023 Sastra Inggris.pdf', '1', NULL, 'auditor2', NULL, '2024-12-09 20:00:04'),
        (1475, '63ae8d4b-074e-4710-9f15-2b3a9a6d8f09', NULL, 419, 30, 'Proker 2024 Sastra Jepang.xlsx', '0', 'Proker No. 8 Penguatan Mata Kuliah Berbahasa Jepang untuk Dosen\r\nProker No. 15 Peningkatan Kompetensi Dosen (JLPT)', 'auditor2', NULL, '2024-12-19 08:26:14'),
        (1478, 'a47a5c3b-dd88-4576-8a5d-71b0a7645524', NULL, 419, 36, 'Struktur Organisasi + SK.pdf', '1', NULL, 'auditor2', NULL, '2024-12-16 08:43:24'),
        (1481, '094ab4ab-94ca-4ca4-9ac1-45484cf275e2', NULL, 419, 42, 'SOP Struktur Organisasi FISIB .pdf', '1', NULL, 'auditor2', NULL, '2024-12-19 08:38:29'),
        (1484, 'e4a2b7cf-cea2-48ca-850c-e160c04694bc', NULL, 419, 48, 'bukti closing KTS.pdf', '1', NULL, 'auditor2', NULL, '2024-12-16 08:45:13'),
        (1487, '7551b08c-ed61-4275-a15c-cf0a9f843410', NULL, 422, 30, 'Proker 2024_Ilmu Komunikasi.xlsx', NULL, NULL, 'auditor2', NULL, '2024-11-22 16:31:55'),
        (1490, '9415de6f-1e35-47da-9a2a-3cef4dfd5426', NULL, 422, 36, 'SK_Struktur Organisasi_Job Description_Job Spesification.pdf', NULL, NULL, 'auditor2', NULL, '2024-11-22 21:24:14'),
        (1493, '3203fadd-7e72-4bf8-b59e-95b371d8d63b', NULL, 422, 42, 'SOP Struktur Organisasi Prodi Ilmu Komunikasi .pdf', NULL, NULL, 'auditor2', NULL, '2024-11-22 19:35:36'),
        (1496, '41b98daa-792e-4328-905e-37275c18ef22', NULL, 422, 48, 'KTS.pdf', NULL, NULL, 'auditor2', NULL, '2024-11-22 19:03:53'),
        (1499, 'dce5b6b3-0029-45c8-8a9d-969c5e1c0713', NULL, 425, 28, 'PROKER PENGELOLA 1 FT 2024_Rev 15112023.xlsx', NULL, NULL, 'auditor2', NULL, '2024-11-22 15:29:42'),
        (1502, 'ada4ec87-a13b-4079-b4d1-fd4a1fd3e5bb', NULL, 425, 34, 'Struktur Organisasi, Job Desk, Job Spec dan SK.pdf', '1', NULL, 'auditor2', NULL, '2024-12-20 08:00:21'),
        (1505, '41b08b45-e1d0-4b3c-ae68-41e7560462fa', NULL, 425, 40, 'SOP Tupoksi.pdf', '1', NULL, 'auditor2', NULL, '2024-12-20 08:00:37'),
        (1508, '5793411a-60a1-4f1b-ac92-b5f36a2d54b4', NULL, 425, 46, 'Surat Pemberitahuan kepada ketua LPM Unpak.pdf', NULL, NULL, 'auditor2', NULL, '2024-11-22 20:01:49'),
        (1511, '5e8bf208-ef67-4bed-98ae-452a43f6a023', NULL, 428, 30, 'Proker Elektro 2024-lengkap analisa baru.xlsx', '1', NULL, 'auditor2', NULL, '2024-12-17 05:01:57'),
        (1514, 'd82e80ff-394e-459e-8a88-af1c25f25dcf', NULL, 428, 36, 'Struktur Organisasi, Job Desk, Job Spec dan SK.pdf', '1', NULL, 'auditor2', NULL, '2024-12-17 05:33:56'),
        (1517, '0bf943fc-2b6d-4ae7-96b5-3373c71b7b66', NULL, 428, 42, 'SOP Elektro 2023_ok.pdf', '1', NULL, 'auditor2', NULL, '2024-12-17 06:24:39'),
        (1520, 'eac9ee34-770f-4e16-a272-1f61e76a14c6', NULL, 428, 48, 'close audit 2023_BA.pdf', '1', NULL, 'auditor2', NULL, '2024-12-17 06:25:26'),
        (1523, '61b1e560-dd4c-4a0e-9b06-1957c1a51a51', NULL, 431, 30, NULL, '0', 'sertifikasi laboran belum berjalan\r\nProkre belum di Monev sesuai template', 'auditor2', NULL, '2024-12-20 14:05:09'),
        (1526, '11ac773c-5df5-48d9-928e-45e2c89702a9', NULL, 431, 36, NULL, '1', NULL, 'auditor2', NULL, '2024-12-20 14:06:29'),
        (1529, '1897da5d-1be2-4bae-b095-3cc7076ffdec', NULL, 431, 42, NULL, '1', NULL, 'auditor2', NULL, '2024-12-20 14:07:53'),
        (1532, '2452dc6f-1bf5-4b0d-9c6f-71a27eaae532', NULL, 431, 48, NULL, '0', NULL, 'auditor2', NULL, '2024-12-20 14:08:16'),
        (1535, '8a61b782-9b82-4d20-994c-bbcbde203a0a', NULL, 434, 30, 'Laporan Proker Geodesi 2024_28102024.xlsx', '1', NULL, 'auditor2', NULL, '2024-12-11 11:24:26'),
        (1538, '8fbbe0a6-95b3-403a-8b6b-079ad1bc7454', NULL, 434, 36, 'Struktur Organisasi, Job Desk, Job Spec dan SK.pdf', '1', NULL, 'auditor2', NULL, '2024-12-11 11:25:08'),
        (1541, '5bc26836-9b49-4869-a2af-28e80c025e02', NULL, 434, 42, 'SOP Tupoksi.pdf', '1', NULL, 'auditor2', NULL, '2024-12-11 11:25:26'),
        (1544, '5101356d-036a-473a-9d53-97949910bdff', NULL, 434, 48, 'close KTS 1-9.pdf', '1', NULL, 'auditor2', NULL, '2024-12-11 11:25:44'),
        (1547, '65a2754c-9431-48f7-b8bb-7f1144a52ebd', NULL, 437, 30, 'Proker Prodi TEKNIK GEOLOGI 2024_perbaikan_used.xlsx', '1', NULL, 'auditor2', NULL, '2024-12-20 10:31:28'),
        (1550, 'f2b948fa-6685-4632-b6f0-80c441fbf2a3', NULL, 437, 36, '4. Struktur Organisasi, Job Desk, Job Spec dan SK.pdf', '1', NULL, 'auditor2', NULL, '2024-12-20 01:00:31'),
        (1553, 'ba2feaa4-2944-415b-94ab-ed0a29b4cd6e', NULL, 437, 42, '3. SOP Tupoksi.pdf', '1', NULL, 'auditor2', NULL, '2024-12-20 01:01:51'),
        (1556, '5b335438-f094-45ad-a65a-42040a668fca', NULL, 437, 48, NULL, '0', 'Belum ada bukti yang disampaikan', 'auditor2', NULL, '2024-12-20 01:02:44'),
        (1559, '4db8f67c-9296-44aa-954b-e7ed5bc35fdf', NULL, 440, 30, 'PROKER PWK 2024.xlsx', '1', NULL, 'auditor2', NULL, '2024-12-20 11:30:18'),
        (1562, 'e6dc1144-1030-46ca-a068-987d0d23fad5', NULL, 440, 36, 'Struktur Organisasi, Job Desk, Job Spec dan SK.pdf', '1', NULL, 'auditor2', NULL, '2024-12-20 11:32:03'),
        (1565, '32da1ba9-4ae6-4afc-a6b7-4b6be56af827', NULL, 440, 42, 'SOP Tupoksi.pdf', '1', NULL, 'auditor2', NULL, '2024-12-20 11:32:14'),
        (1568, '2a4a9b03-5c22-4f44-82ff-5749cf4fbb72', NULL, 440, 48, 'berita acara dan closing temuan AMI 2023.pdf', '1', NULL, 'auditor2', NULL, '2024-12-20 11:33:45'),
        (1571, 'a8902d98-6a34-42b2-9d1a-4f06e0b325fb', NULL, 443, 28, '0 Proker Fakultas 2024 rev191223 (3).pdf', '1', NULL, 'auditor2', NULL, '2024-12-09 10:24:15'),
        (1574, '53dc0563-fc58-408e-ac1c-87411a18439f', NULL, 443, 34, NULL, '1', NULL, 'auditor2', NULL, '2024-12-17 10:50:58'),
        (1577, 'a7ee2f13-4399-4ec9-9e28-69282f2ae483', NULL, 443, 40, NULL, '1', NULL, 'auditor2', NULL, '2024-12-17 10:51:16'),
        (1580, '60286211-eeb7-4f8d-abbd-5ea47c6b9090', NULL, 443, 46, NULL, '1', NULL, 'auditor2', NULL, '2024-12-17 10:35:00'),
        (1583, 'b1ffff6b-2746-4db6-aaab-57e6b572bb48', NULL, 446, 30, 'Template Proker 2024 Biologi.xlsx', '1', NULL, 'auditor2', NULL, '2024-12-20 12:20:26'),
        (1586, '8015c99f-964b-4260-8d2a-aa51f90f2f61', NULL, 446, 36, 'Struktur Organisasi FMIPA Job desk dan job spek _compressed (1).pdf', '1', NULL, 'auditor2', NULL, '2024-12-20 12:21:04'),
        (1589, '008b4521-b608-4f7d-961c-ed960a02c1c4', NULL, 446, 42, 'SOP Pelaksanaan Tugas Sesuai Dengan Job Description Dan Job Spesification (2).pdf', '1', NULL, 'auditor2', NULL, '2024-12-20 12:22:06'),
        (1592, 'c6f25123-12ed-41a3-a5dc-00205828e875', NULL, 446, 48, 'kts closed.pdf', '1', NULL, 'auditor2', NULL, '2024-12-20 12:22:40'),
        (1595, '3b78e3f6-39ea-48eb-befe-ac2f61672908', NULL, 449, 30, 'Proker PS Kimia 2024-Rev 181123.xlsx', '1', NULL, 'auditor2', NULL, '2024-12-05 15:53:12'),
        (1598, '633608e5-4397-4c5d-8bb3-0f44394011ba', NULL, 449, 36, 'STRUKTUR ORGANISASI JOB DES JOB SPEK PS KIMIA 2023 TTD Pak Dekan.pdf', '1', NULL, 'auditor2', NULL, '2024-12-05 15:57:46'),
        (1601, '51b8fc59-c7a1-4334-9969-e61a072d4c00', NULL, 449, 42, 'SOP Fakultas 211124.pdf', '1', NULL, 'auditor2', NULL, '2024-12-05 15:58:15'),
        (1604, '51fc1306-dbde-4b6f-8c11-ff84382898be', NULL, 449, 48, 'laporan audit 2023.pdf', '1', NULL, 'auditor2', NULL, '2024-12-05 15:58:46'),
        (1607, '49e2b09d-4f7f-4301-b3ca-6559ad58e42d', NULL, 452, 30, 'Rev Proker 2024 Matematika - Fitria Virgantari.xlsx', '1', NULL, 'auditor2', NULL, '2024-12-16 15:02:27'),
        (1610, '1ba3133f-a362-49fe-95eb-5097408a0de4', NULL, 452, 36, 'Jobdesc jobspec prodi.pdf', '1', NULL, 'auditor2', NULL, '2024-12-16 15:03:06'),
        (1613, '1496b466-6beb-46a7-b4fa-2e42995fd515', NULL, 452, 42, 'Daftar SOP NEW.pdf', '1', NULL, 'auditor2', NULL, '2024-12-16 15:04:04'),
        (1616, '9503cb6e-0ead-4772-9d20-97c8ff11f638', NULL, 452, 48, 'bukti closing.pdf', '1', NULL, 'auditor2', NULL, '2024-12-16 15:03:37'),
        (1619, '786e12cb-d61a-437c-8a33-8efd82bd1453', NULL, 455, 30, 'Proker Ilmu Komputer 2025 Edit.xlsx', '0', 'student excange dan tracer studi belum terlaksana', 'auditor2', NULL, '2024-11-29 09:52:20'),
        (1622, 'd8f59281-914a-499e-9fb0-52e61bfad2e3', NULL, 455, 36, 'Jobdes dan Job Spec Ilkom.pdf', '1', NULL, 'auditor2', NULL, '2024-11-29 09:53:26'),
        (1625, '6e54e3f5-a64a-4a54-9c04-c810a3ba614d', NULL, 455, 42, 'Jobdes dan Job Spec Ilkom.pdf', '1', NULL, 'auditor2', NULL, '2024-11-29 09:53:48'),
        (1628, 'eee68227-ded7-4647-a54f-e1fa15b9fea5', NULL, 455, 48, 'Closing audit Ilkom 2023.pdf', '1', NULL, 'auditor2', NULL, '2024-11-29 09:54:08'),
        (1631, '990dd457-16d2-4120-9534-1eaa29ab8dca', NULL, 458, 30, 'Program Kerja Farmasi 2024.xlsx', '1', NULL, 'auditor2', NULL, '2024-12-09 11:39:06'),
        (1634, '0f6dadd3-e877-47fd-af84-89bc02f0b63c', NULL, 458, 36, 'LINK STRUKTUR ORGANISASI, JOBDESK DAN JOB SPEC FARMASI UNIVERSITAS PAKUAN.docx', '1', NULL, 'auditor2', NULL, '2024-12-09 11:43:53'),
        (1637, 'bf8fd01a-d3cb-4b3f-a07b-ef6a15714c66', NULL, 458, 42, 'LINK SOP FARMASI UNIVERSITAS PAKUAN.docx', '1', NULL, 'auditor2', NULL, '2024-12-17 11:24:18'),
        (1640, '810e457c-c77e-452d-96f6-85ef563f4cf4', NULL, 458, 48, 'KTS 2023 Closed.pdf', '1', NULL, 'auditor2', NULL, '2024-12-09 11:55:57'),
        (1643, '3e24299d-1500-4677-b836-8eb517d84420', NULL, 461, 28, '2024 PROKER PASCASARJANA UPDATE 16 DESEMBER 2023.xlsx', '0', 'Proker belum terlihat dimonitor', 'auditor2', NULL, '2024-12-18 20:31:33'),
        (1646, 'be066704-5805-4607-95b2-4fbc28745415', NULL, 461, 34, 'PEDOMAN STRUKTUR ORGANISASI, JOB DES dan JOB SPEC 2024(2).pdf', '1', NULL, 'auditor2', NULL, '2024-12-19 09:18:50'),
        (1649, '823c41ee-d2b7-4706-b37f-422f6b845ad7', NULL, 461, 40, 'DAFTAR SOP SPS.pdf', '1', NULL, 'auditor2', NULL, '2024-12-18 20:34:55'),
        (1652, '3ac1bd53-64cd-4525-8141-70f4244d9579', NULL, 461, 46, 'berita_acara.pdf', '0', 'Bukti yang diberikan adalah berita acara audit akuntansi S1', 'auditor2', NULL, '2024-12-18 20:35:59'),
        (1655, '76477093-b1ae-4e62-9d13-21a468c0c1e6', NULL, 464, 30, '2024 PROKER MANAJEMEN PENDIDIKAN S3.xlsx', '1', NULL, 'auditor2', NULL, '2024-12-18 16:10:02'),
        (1658, 'fa9b198d-34e4-4325-8f13-f7721c98f9b5', NULL, 464, 36, '1. STRUKTUR ORGANISASI, TUGAS & TANGGUNGJWB.docx', '1', NULL, 'auditor2', NULL, '2024-12-19 14:09:41'),
        (1661, '769e4bb2-bb41-4506-9dab-d107a383f117', NULL, 464, 42, '6. SOP PENGELOLAAN UNIT PENGELOLA.pdf', '1', NULL, 'auditor2', NULL, '2024-12-18 16:11:31'),
        (1664, 'e9fedc0e-b2a0-4c9e-a33c-91e11e45bc44', NULL, 464, 48, 'kts-AUDIT INTERNAL 27-12-2023.pdf', '1', NULL, 'auditor2', NULL, '2024-12-19 14:08:00'),
        (1667, '06682e14-72ba-4f86-bacd-3e1d5b3545ba', NULL, 467, 30, 'Proker Prodi IM 2025 (1).xlsx', '0', 'Belum Lengkap dan tidak dilakukan Monev\r\nSiproker tidak diisi', 'auditor2', NULL, '2024-12-19 14:45:16'),
        (1670, '8e770859-f42b-4a32-8839-bf4d2e92dead', NULL, 467, 36, NULL, '1', NULL, 'auditor2', NULL, '2024-12-19 14:47:35'),
        (1673, 'fe8970d6-9434-4fb9-90d3-190e2ead4c8f', NULL, 467, 42, NULL, '1', NULL, 'auditor2', NULL, '2024-12-19 14:50:43'),
        (1676, 'e1bc92e5-e6ae-4fe7-bb05-f62e81c10127', NULL, 467, 48, NULL, '1', NULL, 'auditor2', NULL, '2024-12-19 14:51:43'),
        (1679, '351f5601-e7ce-4ffa-9e4c-7ce7950f2806', NULL, 470, 30, '2024 PROKER PRODI ADM. PENDIDIKAN S2.xlsx', '1', NULL, 'auditor2', NULL, '2024-12-03 13:41:51'),
        (1682, '1781e46a-4714-4ef8-81e2-e1bac1991740', NULL, 470, 36, 'PEDOMAN STRUKTUR ORGANISASI, JOB DES dan JOB SPEC.pdf', '1', NULL, 'auditor2', NULL, '2024-12-03 13:42:58'),
        (1685, 'febf6615-9351-4430-ac6b-55bd230ab480', NULL, 470, 42, 'SOP TATA USAHA SPs UNPAK 2023.pdf', '1', NULL, 'auditor2', NULL, '2024-12-03 13:49:09'),
        (1688, '980f967e-f37f-4bd2-b17d-9ea7a6d2fadf', NULL, 470, 48, 'kts(1)_kts(2)_kts_merged.pdf', '1', NULL, 'auditor2', NULL, '2024-12-03 13:50:44'),
        (1691, 'af51114d-55af-49e2-a900-6c5042be916a', NULL, 473, 30, 'Proker MNL 2024 Revisi 181123.xlsx', '0', 'Dua proker belum tercapai \r\n1. Pengembangan RPS\r\n2. Produk Inovasi Dosen dan Dosen', 'auditor2', NULL, '2024-12-19 10:35:05'),
        (1694, 'b1a8182a-afb8-47fd-a0a2-3d961a880431', NULL, 473, 36, NULL, '1', NULL, 'auditor2', NULL, '2024-12-19 10:36:52'),
        (1697, '1ee4b475-6379-4c6d-965d-7ec8be0a86d9', NULL, 473, 42, NULL, '1', NULL, 'auditor2', NULL, '2024-12-19 10:39:51'),
        (1700, '8bcce07a-d12d-4e94-bfb8-cbe972acfe75', NULL, 473, 48, NULL, '1', NULL, 'auditor2', NULL, '2024-12-19 10:40:38'),
        (1703, '33262158-808c-4414-ab0c-bf5213f3ee4f', NULL, 476, 30, NULL, '0', 'Siproker tidak di isi jadi tidak ada pelaporan yg sudah dilaksanakan atau belum Proker nya di tahun 2024.\r\nUntuk selanjutanya agar rutin di isi dan di laporkan', 'auditor2', NULL, '2024-12-19 11:58:51'),
        (1706, 'ac77aa62-547b-4921-8778-c4ed94ff3cdb', NULL, 476, 36, NULL, '1', NULL, 'auditor2', NULL, '2024-12-19 11:59:00'),
        (1709, 'e4341046-78f6-4177-acdc-ceed32938b52', NULL, 476, 42, NULL, '1', NULL, 'auditor2', NULL, '2024-12-19 11:59:08'),
        (1712, '97cc476d-d647-4c16-95b4-48faf096c858', NULL, 476, 48, NULL, '1', NULL, 'auditor2', NULL, '2024-12-19 11:59:18'),
        (1715, 'c05c87af-d30b-4784-8b62-be1537de2d25', NULL, 479, 30, NULL, '1', NULL, 'auditor2', NULL, '2024-12-19 09:14:28'),
        (1718, 'c796564a-a455-43ca-89e5-f9d8daab7cf6', NULL, 479, 36, NULL, '1', NULL, 'auditor2', NULL, '2024-12-19 09:08:54'),
        (1721, '68080757-fb3a-46da-8b9d-faba9c3c0b96', NULL, 479, 42, NULL, '1', NULL, 'auditor2', NULL, '2024-12-19 09:24:55'),
        (1724, '3ed6f2a6-9a41-4bd5-ac80-2a4bbe2363f7', NULL, 479, 48, NULL, '1', NULL, 'auditor2', NULL, '2024-12-19 09:13:32'),
        (1727, '14383f83-e98c-48c9-84c2-a67861c9aa1f', NULL, 482, 30, '2024 PROKER PENDIDIKAN IPA_(fix).xlsx', '1', NULL, 'auditor2', NULL, '2024-12-09 10:12:20'),
        (1730, 'cf1dbaad-46f2-496c-93a8-9ab5afa049e7', NULL, 482, 36, 'PEDOMAN STRUKTUR ORGANISASI, JOB DES dan JOB SPEC 2024.pdf', '1', NULL, 'auditor2', NULL, '2024-12-09 10:12:52'),
        (1733, '67e89cea-d09e-44d8-b99e-c7d27f1da083', NULL, 482, 42, 'SOP.pdf', '1', NULL, 'auditor2', NULL, '2024-12-09 10:13:05'),
        (1736, '0a89f066-a929-4641-a4e4-60fdb6c579a3', NULL, 482, 48, 'Hasil AMI 2023.pdf', '1', NULL, 'auditor2', NULL, '2024-12-09 10:14:24'),
        (1739, '61dace19-7372-4a68-8c15-a581353204bd', NULL, 485, 30, 'PROKER PWK S2 2025 (21112024).xlsx', '1', NULL, 'auditor2', NULL, '2024-12-09 10:18:26'),
        (1742, 'e6ae3eb8-ddf1-46d3-ae0a-e03b331be217', NULL, 485, 36, '03-STRUKTUR ORGANISASI MPWK UNPAK-2023 - REVISI[1].pdf', '1', NULL, 'auditor2', NULL, '2024-12-09 10:18:38'),
        (1745, '674a69be-229c-4f8e-b895-5b942859b8c2', NULL, 485, 42, '1. MATRIK KOMPETENSI[1].docx', '1', NULL, 'auditor2', NULL, '2024-12-09 10:18:48'),
        (1748, '787b7297-cc3a-48fd-9ab0-30820d90d99a', NULL, 485, 48, NULL, '0', 'Belum ada bukti dokumen penyelesaian KTS audit sebelumnya', 'auditor2', NULL, '2024-12-09 10:19:26'),
        (1751, '1a830b09-188e-407d-befb-5092eaf2f137', NULL, 488, 30, 'PROKER PENDAS 2024 REVISI.xls', '1', NULL, 'auditor2', NULL, '2024-11-26 18:41:57'),
        (1754, 'f0b90d43-9cb9-4525-a8b6-74a70d3843c3', NULL, 488, 36, NULL, '1', NULL, 'auditor2', NULL, '2024-12-19 16:03:17'),
        (1757, 'd464b068-c611-45ac-a5ad-5150c6f166a7', NULL, 488, 42, NULL, '0', 'belum ada SOP yang mendukung Job desk dan job spek', 'auditor2', NULL, '2024-12-19 16:08:09'),
        (1760, 'bf68b2a0-d8fa-4037-933b-e1e5b33a59bb', NULL, 488, 48, NULL, '1', NULL, 'auditor2', NULL, '2024-12-19 16:08:47'),
        (1763, 'dfebf2cb-26b3-4a3b-9025-2ccc6e6a06c8', NULL, 491, 28, 'PROKER SEKOLAH VOKASI 2024.xlsx', '1', NULL, 'auditor2', NULL, '2024-12-04 20:17:59'),
        (1766, 'dffb7e58-635e-4f44-a9d0-6492e5a5d88d', NULL, 491, 34, 'New SK SV- Struktur Organisasi Jobdesk dan Jobspec.pdf', '1', NULL, 'auditor2', NULL, '2024-12-04 20:18:18'),
        (1769, '74da83f9-3ac8-4e34-9778-ec0314d5292c', NULL, 491, 40, 'SOP SV-merged-2024.pdf', '1', NULL, 'auditor2', NULL, '2024-12-04 20:19:10'),
        (1772, 'ab626061-25e6-4946-b9f1-d069e334efa2', NULL, 491, 46, 'BA - SV 2023.pdf', '1', NULL, 'auditor2', NULL, '2024-12-16 13:28:04'),
        (1775, 'fa1340a2-ac60-4279-a47e-fe018c877b95', NULL, 494, 30, 'Program Kerja 2024 Beserta Monev.pdf', '1', NULL, 'auditor2', NULL, '2024-12-21 08:44:23'),
        (1778, 'a544ea85-eadd-41ac-9a5e-60c2c69fcc2c', NULL, 494, 36, 'SK dan Struktur Organisasi, Job Desc dan Job Spec 2024.pdf', '1', NULL, 'auditor2', NULL, '2024-12-21 08:44:33'),
        (1781, 'b431bb46-44c8-4e65-9727-c9600af24d2e', NULL, 494, 42, 'SOP Pelaksanaan Tugas, Job Desc dan Job Spec.pdf', '1', NULL, 'auditor2', NULL, '2024-12-21 08:45:03'),
        (1784, 'daa0552c-c857-4877-9cf7-73d3f5d41c98', NULL, 494, 48, 'Berita Acara Closing Audit Desember 2023.pdf', '1', NULL, 'auditor2', NULL, '2024-12-21 08:45:22'),
        (1787, '6525e4ff-f75e-428a-9c95-0bb7a8242646', NULL, 497, 30, 'Proker 24-Asli.xlsx', '1', NULL, 'auditor2', NULL, '2024-12-04 20:03:33'),
        (1790, 'e95adeec-d464-4f73-9c39-49ce1a8ea2ec', NULL, 497, 36, '06-SK Struktur Organisasi Jobdesk dan Jobspec.pdf', '1', NULL, 'auditor2', NULL, '2024-12-04 20:05:00'),
        (1793, '871effad-e200-4fcc-9cbc-6ae3ae4ee78e', NULL, 497, 42, '2. Review Pelaksanaan Perkuliahan.docx', '1', NULL, 'auditor2', NULL, '2024-12-19 11:38:52'),
        (1796, '60de0d3f-c35b-4034-b0d2-b79e23d089bb', NULL, 497, 48, 'kts.pdf', '1', NULL, 'auditor2', NULL, '2024-12-19 11:37:02'),
        (1799, '0d0366af-a7e9-4360-9fbe-fd390414f504', NULL, 500, 30, 'Proker Prodi MKP_2024.xlsx', '1', NULL, 'auditor2', NULL, '2024-12-20 14:56:40'),
        (1802, '61162894-e157-4f3b-8b19-d23e5d4b4b82', NULL, 500, 36, '06-SK Struktur Organisasi Jobdesk dan Jobspek_New.pdf', '1', NULL, 'auditor2', NULL, '2024-12-20 14:56:46'),
        (1805, '851344d5-5035-4222-8b35-d1c403646ae1', NULL, 500, 42, 'SOP SV-merged-2024.pdf', '1', NULL, 'auditor2', NULL, '2024-12-20 14:56:54'),
        (1808, '99e19209-b1cb-45b1-b69c-e9aa6701d1c3', NULL, 500, 48, 'Berita Acara Audit 2023.pdf', '1', NULL, 'auditor2', NULL, '2024-12-20 14:57:01'),
        (1811, '6cf71c9a-0228-4d68-a524-ed5b7356e4fa', NULL, 503, 30, '1.Edit  proker d3 teknik komputer-20240213 (1) (1).xlsx', '1', NULL, 'auditor2', NULL, '2024-11-26 18:49:03'),
        (1814, '44e42e64-4851-45a1-88a1-942042b4292f', NULL, 503, 36, 'SK_Struktur Organisasi_Jobdesk_Job Specification_TK.pdf', '1', NULL, 'auditor2', NULL, '2024-11-26 18:50:13'),
        (1817, '70d17efe-3182-4a0b-8d39-a07359868af0', NULL, 503, 42, '1. Pemilihaan struktural.pdf', '1', NULL, 'auditor2', NULL, '2024-11-26 18:51:06'),
        (1820, 'e4e29d84-c931-4a1c-ae42-42d462274441', NULL, 503, 48, 'kts_Clsing Dokumen Utama_Dokumen Tambahan.pdf', '1', NULL, 'auditor2', NULL, '2024-11-26 18:51:35'),
        (1823, '178f22a5-10ba-4326-b480-6ca83b31df95', NULL, 506, 30, 'Proker MI 2024 (05022024) (1).xlsx', '1', NULL, 'auditor2', NULL, '2024-12-16 10:26:58'),
        (1826, '574c2ee2-b13c-466a-bd8b-092129e7be69', NULL, 506, 36, 'SK_Struktur Organisasi_Jobdesk_Job Specification_MI.pdf', '0', 'tidak ada struktur organisasi level prodi, beserta jobdes dan jobspekny misalnya akademik alumni, penelitian, pengabidan  hingga  koordinator bidang keilmuan', 'auditor2', NULL, '2024-12-09 20:19:28'),
        (1829, '11606a40-cd3d-4035-bce4-a03c0e670442', NULL, 506, 42, '1. SOP Penyusunan Struktur Organisasi Prodi.pdf', '1', NULL, 'auditor2', NULL, '2024-12-16 10:27:17'),
        (1832, 'fc4ca230-c566-4714-b668-19ab582c406c', NULL, 506, 48, 'kts_Closing Dokumen Utama_C2 C3 C9.pdf', '1', NULL, 'auditor2', NULL, '2024-12-09 20:21:46'),
        (1916, '95dcf9cc-4e29-4698-987d-42ecbec514b4', 524, 524, 32, '2024_Proker Pusat Inovasi.xlsx', '0', 'Proker tidak di Monev', 'auditor2', NULL, '2025-02-24 08:44:44'),
        (1919, '48a311b3-2d32-4d1b-87bd-0cb427936e7d', 524, 524, 38, 'Struktur Org Job SK.pdf', '1', NULL, 'auditor2', NULL, '2025-02-24 07:34:26'),
        (1922, 'fc753fa1-9a2d-45f6-ac8f-d4b8e24992ea', 524, 524, 44, 'Gabung SOP Job SK.pdf', '1', NULL, 'auditor2', NULL, '2025-02-24 07:34:43'),
        (1925, '8453923e-cd95-446e-a517-255a3682e27e', 524, 524, 50, 'Perbaikan Struk+JobDes+JobSpek+SK 2023.pdf', '1', NULL, 'auditor2', NULL, '2025-02-24 09:32:45'),
        (1940, 'd38d820c-16cc-464f-bd70-110976fccfe2', NULL, 548, 32, 'Proker KKPKT_2024_12OKT2023.xlsx', '1', NULL, 'auditor2', NULL, '2025-03-04 09:13:56'),
        (1943, '1bcb7fef-7548-4519-8610-ce832c47eaed', NULL, 548, 38, NULL, '0', 'SO akan disesuaikan dgn Statuta baru', 'auditor2', NULL, '2025-03-04 09:30:08'),
        (1946, '2b20e936-fbfc-4c23-9ac3-64284abab8f1', NULL, 548, 44, NULL, '1', NULL, 'auditor2', NULL, '2025-03-04 08:05:29'),
        (1949, '4f210c22-8567-4945-ac64-d4f6c91d66df', NULL, 548, 50, NULL, '1', NULL, 'auditor2', NULL, '2025-03-04 08:03:54'),
        (1952, 'ac906681-226a-4f2a-bf7e-bec8e829fa70', NULL, 551, 32, 'Talkshow dengan Dosen Fatani University Thailand Islahuddin, S.S., M.A..pdf', '0', 'data proker tidak tepat', 'auditor2', NULL, '2025-03-04 09:11:37'),
        (1955, '9bb473e1-2176-4336-80c5-33d241a93be4', NULL, 551, 38, 'SKRUKTRUR ORGANISASI, JOBDESK, JOBSPEK.pdf', '1', NULL, 'auditor2', NULL, '2025-03-04 09:26:39'),
        (1958, '72f96613-8f26-42bf-b174-fcf9cbf749c8', NULL, 551, 44, 'SOP Perpustakaan Pusat.pdf', '1', NULL, 'auditor2', NULL, '2025-03-04 09:27:03'),
        (1961, 'da44ec6a-3ebe-4f86-b692-46cb75ce5465', NULL, 551, 50, 'CamScanner 14-02-2025 13.20.pdf', '1', NULL, 'auditor2', NULL, '2025-03-04 09:46:13'),
        (1964, '1f6a8d0a-e15d-47d8-912d-e1f9e4af4e43', NULL, 530, 32, 'PROGRAM KERJA LPPM 2024.xlsx', NULL, NULL, 'auditor2', NULL, '2025-02-09 16:01:40'),
        (1967, 'df4a2f70-6d6a-4dc0-b222-5b775ecd7a9e', NULL, 530, 38, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (1970, 'c628a1e8-f1cc-4a3f-81fe-0e33fbdd7cbf', NULL, 530, 44, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (1973, '493c246b-63b8-4cb8-8be1-e95a6d0bcb5f', NULL, 530, 50, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (1976, 'e8997f72-fa46-47ac-9bee-02b981e6af0b', NULL, 533, 32, NULL, '0', 'Proker belum di Monev', 'auditor2', NULL, '2025-03-04 10:06:04'),
        (1979, '08dbbef8-1cf2-4a3c-93b9-830331112df9', NULL, 533, 38, NULL, '1', NULL, 'auditor2', NULL, '2025-03-04 13:34:11'),
        (1982, 'e1ee5476-38e8-40f8-afb0-143e5e2be5b4', NULL, 533, 44, NULL, '1', NULL, 'auditor2', NULL, '2025-03-04 10:14:46'),
        (1985, '35f9c8e4-6828-4cfc-952e-5b38983eb983', NULL, 533, 50, NULL, '1', NULL, 'auditor2', NULL, '2025-03-04 10:18:08'),
        (1988, '97a92cff-dd32-42e2-80c5-39163251b879', NULL, 536, 32, 'Proker LPM 2024.xlsx', '1', NULL, 'auditor2', NULL, '2025-02-24 11:05:58'),
        (1991, 'b23b60ca-6fbf-4b40-b233-c6662a42e25c', NULL, 536, 38, 'softcopy buku profil LPM 2024 UPDATE OKT.pdf', '1', NULL, 'auditor2', NULL, '2025-02-24 11:06:03'),
        (1994, '6bdd803b-b918-499d-9600-468dcc095b35', NULL, 536, 44, 'softcopy buku profil LPM 2024 UPDATE OKT.pdf', '1', NULL, 'auditor2', NULL, '2025-02-24 11:06:06'),
        (1997, '0a3f59c7-898f-43bf-9be6-73c82448197b', NULL, 536, 50, 'log stasus LPM 24.pdf', '1', NULL, 'auditor2', NULL, '2025-02-24 11:06:10'),
        (2000, 'd933c8aa-c3af-44ca-8b2c-f0f8e479cff3', 539, 539, 32, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2003, '249bb97a-d79a-42db-ad74-161ed03b6c7c', 539, 539, 38, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2006, '4c193fff-51cf-46e5-bbb1-670e0c1db74a', 539, 539, 44, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2009, '402ef044-81f9-46f9-bb7f-392ca44f5067', 539, 539, 50, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2024, 'fb0fd5f4-ae1e-46cd-b5b9-75926d6161d4', NULL, 545, 32, 'Proker UPP 2024_TERBARU.xlsx', '0', NULL, 'auditor2', NULL, '2025-03-04 10:13:53'),
        (2027, 'b0642454-4c2c-43d0-8fb3-eaa20d09311e', NULL, 545, 38, 'STRUKTUR ORGANISASI & JOBDES UNPAK PRESS_2024.pdf', '1', NULL, 'auditor2', NULL, '2025-02-20 10:19:46'),
        (2030, '4f951c86-65ff-4b4d-8873-026cf6d9b4db', NULL, 545, 44, 'SOP PROGRAM KERJA & ADMIN UNPAK PRESS_2024.pdf', '1', NULL, 'auditor2', NULL, '2025-02-20 10:27:07'),
        (2033, '1d70cd3c-f2cc-4dfb-a65f-cff8164c1f0c', NULL, 545, 50, 'Berita Acara_Audit Internal 2023.pdf', '1', NULL, 'auditor2', NULL, '2025-02-20 10:29:42'),
        (2036, '071d0e51-b037-4475-8152-ef0f1603ee4f', NULL, 554, 32, 'LINK DOKUMEN SIAMIDA.pdf', '0', 'Dokumen Siproker belum dilengkapi dan Proker belum dimonev', 'auditor2', NULL, '2025-03-06 13:56:21'),
        (2039, '3bc439e5-dfba-4a01-9f3f-8da4249d4e7e', NULL, 554, 38, 'Link Siamida Struktur Organisasi Disertai Dengan Job Description.pdf', '0', 'ada Statuta baru sehingga akan disesuaikan', 'auditor2', NULL, '2025-03-06 13:38:48'),
        (2042, '7a5e15b2-e67f-464d-81c2-46d73d3e70ce', NULL, 554, 44, 'Link Apakah Memiliki SOP Pelaksanaan Tugas Sesuai Dengan Job Description Dan Job Spesification.pdf', '1', NULL, 'auditor2', NULL, '2025-03-06 13:40:43'),
        (2045, 'a4014ded-7fef-4c2d-80c1-2400c64b6ca1', NULL, 554, 50, 'Berita Acara Kegiatan Audit Internal 2023.pdf', '1', NULL, 'auditor2', NULL, '2025-03-06 13:41:39'),
        (2046, '17a321f1-d726-4575-8e7a-60374c6f7823', NULL, 557, 56, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2047, 'f8b92da5-fbf6-4eee-8d3c-e521724cde8e', NULL, 557, 59, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2048, '2774c2ae-67d8-45a5-8ea9-fd23e6de3bd2', NULL, 557, 62, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2049, '213cde41-674f-48ac-88da-bc81b4522301', NULL, 557, 65, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2050, '602661cd-f1bf-46fb-9b88-d77d5373cea7', NULL, 556, 56, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2051, '83a4d4d9-d394-4a71-aa90-1c0c3a7e4f65', NULL, 556, 59, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2052, 'b0f313c4-fb1b-4165-aaa5-f90903ad7568', NULL, 556, 62, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2053, '7d959915-fa74-42c9-aa65-5c72edae2cd0', NULL, 556, 65, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2054, '92c64c58-370e-445c-b3a5-3958de24f5db', NULL, 556, 56, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2055, 'a29175cb-e148-4183-89a4-841d096cc549', NULL, 556, 59, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2056, '86045b69-927d-4a99-8c0f-5774a59778f7', NULL, 556, 62, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2057, '56d0d874-bc62-4ab9-9e6b-52f24b8efda9', NULL, 556, 65, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2058, 'f9d7239f-9488-4732-b2bb-b7ac7e3bf9a2', NULL, 606, 57, 'Proker KBR vokasi - 2025.xlsx', NULL, NULL, 'auditor2', NULL, '2025-11-12 10:42:44'),
        (2059, '8a2e108b-02b1-4b98-ba68-a6497559266c', NULL, 606, 60, 'Jobdesk dan Jobspek 1 (2)-4-43.pdf', NULL, NULL, 'auditor2', NULL, '2025-11-28 09:39:01'),
        (2060, '4c7ce34d-e6a1-481a-a202-40252ffe5298', NULL, 606, 63, 'SOP SV-merged.pdf', NULL, NULL, 'auditor2', NULL, '2025-11-14 14:13:00'),
        (2061, 'f25ac2aa-fa25-4138-bd44-9724abbe4e4e', NULL, 606, 66, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2062, '32926153-38ba-49ad-b648-914f920e45ee', NULL, 564, 57, 'PROKER_2025_FH UNPAK_TERBARU (REVISI 20-1-2025) (2).xlsx', NULL, NULL, 'auditor2', NULL, '2025-11-28 11:26:53'),
        (2063, 'cb6224ba-9bff-4f1b-a9f5-589ad6577564', NULL, 564, 60, NULL, NULL, NULL, 'auditor2', NULL, '2025-10-15 10:17:52'),
        (2064, '2134d0c5-3918-4b7b-bb93-460c19907c67', NULL, 564, 63, 'SOP.pdf', NULL, NULL, 'auditor2', NULL, '2025-11-28 12:20:05'),
        (2065, '60dd677e-aa30-4bf1-b64f-8c6149e6e763', NULL, 564, 66, NULL, NULL, NULL, 'auditor2', NULL, '2025-10-15 10:18:02'),
        (2066, '94afc574-4092-4cec-a67c-98d780b7ad56', NULL, 608, 58, 'kts_renstra.pdf', '1', NULL, 'auditor2', NULL, '2025-10-15 10:19:34'),
        (2067, '7a5fc297-6bde-47cd-a6b7-f47f1920c8f3', NULL, 608, 61, 'kts_renstra (4).pdf', '1', NULL, 'auditor2', NULL, '2025-10-15 10:19:39'),
        (2068, 'd1bf7918-901b-44d9-a36b-cc0fa3de832b', NULL, 608, 64, 'Test.pdf', '1', NULL, 'auditor2', NULL, '2025-10-15 10:19:43'),
        (2069, '164ef30c-306c-4139-aa96-e580468098c7', NULL, 608, 67, 'kts_renstra (2).pdf', '1', NULL, 'auditor2', NULL, '2025-10-15 10:19:46'),
        (2070, 'cdd1dac7-204c-4fac-ad1c-3ed154c67d43', NULL, 609, 57, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2071, '0e6f03c8-7b47-4d4c-b29f-a6dd0ebd06e8', NULL, 609, 60, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2072, '95633a13-2897-4a04-b17a-55717ce221be', NULL, 609, 63, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2073, 'e725dbf2-19e9-4f29-bde6-21ef58e75717', NULL, 609, 66, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2074, '2c4850c0-4984-4193-8f68-daedba07868c', NULL, 610, 57, 'Proker S2-KOM 2025.xlsx', NULL, NULL, 'auditor2', NULL, '2025-11-28 12:49:59'),
        (2075, '392e9d88-6f0e-4874-a2d7-9d4538ead71f', NULL, 610, 60, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2076, 'bf01d14b-1484-47b9-b225-d50d6db94b44', NULL, 610, 63, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2077, '9f5e8a5c-277b-466b-9dfc-dfcdd43a7bd2', NULL, 610, 66, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2078, '13ee528e-2790-46d9-a647-6e48ac1d1014', NULL, 593, 57, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2079, 'e2a1f9bf-73de-4678-9aef-c4d5e2c2b23f', NULL, 593, 60, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2080, '631c4a75-92e4-4ee4-a9b9-41c490eb544a', NULL, 593, 63, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2081, 'a291a363-14e4-446b-8cc3-432935745a37', NULL, 593, 66, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2082, 'da2f1108-aaf1-4c43-9578-37a23e4b4cbc', NULL, 558, 56, 'PROKER 2025.pdf', NULL, NULL, 'auditor2', NULL, '2025-11-28 15:10:52'),
        (2083, 'dd66ae75-56cb-43c6-a182-358d91c8ca77', NULL, 558, 59, 'ORGANISASI DAN TUPOKSI PIMPINAN FKIP UNPAK_compressed.pdf', NULL, NULL, 'auditor2', NULL, '2025-11-26 14:18:27'),
        (2084, 'f2483839-ff83-4bed-9d54-7226f490a19a', NULL, 558, 62, 'SOP FKIP 2022.pdf', NULL, NULL, 'auditor2', NULL, '2025-11-27 10:55:49'),
        (2085, '3160516c-8759-4b2e-b657-6ed85e02ea62', NULL, 558, 65, 'kts_rensta_02_merged.pdf', NULL, NULL, 'auditor2', NULL, '2025-11-28 02:59:28'),
        (2086, 'b41d27c6-8733-4c36-a955-f9889b899ab9', NULL, 559, 56, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2087, '0bd33142-da0a-4703-b744-823297434dff', NULL, 559, 59, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2088, '57d5b554-6ded-4f06-bcf3-b13dfb6ddc10', NULL, 559, 62, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2089, '9ebb771d-8c33-4e0a-9d07-1a06d9067e10', NULL, 559, 65, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2090, 'ea935b70-ea86-496d-8a9b-4e0f9e4c9ac4', NULL, 560, 56, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2091, '5dcbb8db-717a-43b3-95c8-a6443bc918d1', NULL, 560, 59, 'SK  STRUKTUR ORGANISASI JOBDESC DAN JOBSPEC FAKULTAS TEKNIK_compressed (1).pdf', NULL, NULL, 'auditor2', NULL, '2025-11-26 14:50:24'),
        (2092, '8cd88b03-f425-4aa7-a2e5-518bd38579c9', NULL, 560, 62, 'SOP Tupoksi (1).pdf', NULL, NULL, 'auditor2', NULL, '2025-11-26 14:51:11'),
        (2093, '6925eb21-e4b3-47d6-9563-e4a2a4b005f1', NULL, 560, 65, 'Closing AMI 2024.pdf', NULL, NULL, 'auditor2', NULL, '2025-11-28 08:56:40'),
        (2094, '340cdcf5-6cd3-41b2-9020-3add0e6917c6', NULL, 561, 56, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2095, '12889d64-f94d-45db-b138-b64fb69fa5dc', NULL, 561, 59, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2096, 'd31505ae-43ce-4ce7-8af7-fac84712713d', NULL, 561, 62, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2097, '7b51315c-1f9c-47c8-86c6-cf718f2ffef4', NULL, 561, 65, 'CLOSED AUDIT FMIPA 2024.pdf', NULL, NULL, 'auditor2', NULL, '2025-11-28 12:00:47'),
        (2098, '9f5c8f90-4c38-4607-8a3a-6f92fe38154f', NULL, 562, 56, '2025 PROKER PASCASARJANA FULL.xlsx', NULL, NULL, 'auditor2', NULL, '2025-11-25 13:45:26'),
        (2099, '425d079a-9c2f-4987-ae62-5e1ee1a2d3e6', NULL, 562, 59, 'PEDOMAN STRUKTUR ORGANISASI, JOB DES dan JOB SPEC 2024.pdf', NULL, NULL, 'auditor2', NULL, '2025-11-25 13:46:06'),
        (2100, '6e4355ed-6504-4488-b55b-9025685f1f59', NULL, 562, 62, 'DAFTAR SOP TERBARU.pdf', NULL, NULL, 'auditor2', NULL, '2025-11-25 13:47:12'),
        (2101, '99652d5c-b030-48f6-b0dd-1a383d740a0a', NULL, 562, 65, 'berita_acara (1).pdf', NULL, NULL, 'auditor2', NULL, '2025-11-25 13:49:05'),
        (2102, '2ce1adf2-ec3d-446a-9f9e-205f863ec6d7', NULL, 563, 56, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2103, 'f0d0947a-f706-4d82-9730-cfad05690811', NULL, 563, 59, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2104, '73db4239-cfae-4673-9d42-4b216ffc6828', NULL, 563, 62, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2105, 'e02c1175-592e-44d0-8359-97848d3fa35e', NULL, 563, 65, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2106, '147cc179-467e-4377-8387-d9f0cd2cd2fe', NULL, 565, 57, 'Mjn-Proker 2025.xlsx', NULL, NULL, 'auditor2', NULL, '2025-11-28 14:47:48'),
        (2107, '7b1373a9-8225-4016-867a-590b33c08036', NULL, 565, 60, 'BUKU JOB DESCRIPTION-2025 per 5 Nov.docx', NULL, NULL, 'auditor2', NULL, '2025-11-28 14:58:51'),
        (2108, '357d6743-58a6-47fb-be64-f37cfcd14459', NULL, 565, 63, 'BUKU JOB DESCRIPTION-2025 per 5 Nov.docx', NULL, NULL, 'auditor2', NULL, '2025-11-28 14:59:37'),
        (2109, 'd31dcd89-aba6-4c5a-9ae4-6c529253c880', NULL, 565, 66, 'Temuan Hasil Audit Tahun 2024.pdf', NULL, NULL, 'auditor2', NULL, '2025-11-28 14:59:26'),
        (2110, '5d77ad5b-a5dc-4698-a5e9-10d9fa95dd0e', NULL, 566, 57, 'proker bdi-2025.pdf', NULL, NULL, 'auditor2', NULL, '2025-11-28 14:49:42'),
        (2111, '70d91842-0df8-4739-b1f7-764e58964122', NULL, 566, 60, 'BUKU JOB DESCRIPTION-2025 per Agustus 2025.docx.pdf', NULL, NULL, 'auditor2', NULL, '2025-11-28 14:33:44'),
        (2112, 'bbdc1944-02f1-4284-8a34-c709a536264d', NULL, 566, 63, 'BUKU JOB DESCRIPTION-2025 per Agustus 2025.docx.pdf', NULL, NULL, 'auditor2', NULL, '2025-11-28 14:49:31'),
        (2113, '044a9068-cca5-478e-b516-3cad271881b2', NULL, 566, 66, 'Close BA AMI 2024.pdf', NULL, NULL, 'auditor2', NULL, '2025-11-28 13:53:09'),
        (2114, 'a991cfdf-0be3-4c10-8627-3cb528df5a77', NULL, 567, 57, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2115, '9ba81bb2-3b67-4381-a0cb-6eee8fed5f08', NULL, 567, 60, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2116, 'cdb7bdc0-8d84-43b8-8c35-b7a714915544', NULL, 567, 63, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2117, 'd312e6ca-94af-4198-84fa-9ace0ca3ea2f', NULL, 567, 66, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2118, '6db1df5d-d28f-4e34-99fa-6afeeaef4536', NULL, 572, 57, 'Laporan PROKER 2025.pdf', NULL, NULL, 'auditor2', NULL, '2025-11-27 10:36:27'),
        (2119, 'acd8571a-7fea-4f50-9394-ba5e991a0198', NULL, 572, 60, 'Struktur Organisasi & Job Des 2025.pdf', NULL, NULL, 'auditor2', NULL, '2025-11-27 12:56:17'),
        (2120, 'f8c1c2f0-5a41-4078-bd86-b8c1ab6cb827', NULL, 572, 63, 'SOP JOBDES & Pelaksanaan Tugas 2025.pdf', NULL, NULL, 'auditor2', NULL, '2025-11-27 13:11:35'),
        (2121, '82cf8618-2722-4bbf-bd03-b715d5d821cc', NULL, 572, 66, 'kts_renstra AMI 2024.pdf', NULL, NULL, 'auditor2', NULL, '2025-11-27 12:24:17'),
        (2122, 'd57e11d0-a9ec-4602-a214-29b0035e7477', NULL, 573, 57, 'Proker 2025 Pendidikan Bahasa Inggris.pdf', NULL, NULL, 'auditor2', NULL, '2025-11-28 14:28:24'),
        (2123, '03ae7935-d779-476d-9b4c-e5ab7797f331', NULL, 573, 60, 'SK, Struktur Organisasi, Job Desc.pdf', NULL, NULL, 'auditor2', NULL, '2025-11-27 20:58:53'),
        (2124, 'bf63b25d-23f1-44b3-ae09-62a7085841b4', NULL, 573, 63, 'SOP JOBDES.pdf', NULL, NULL, 'auditor2', NULL, '2025-11-27 14:45:53'),
        (2125, '42de45b1-5132-4090-8bbe-5a4b43f694e0', NULL, 573, 66, 'Closing AMI 30 Agustus 2024.pdf', NULL, NULL, 'auditor2', NULL, '2025-11-27 15:01:41'),
        (2126, '659f3af6-18b9-49dc-9ee2-a04cfd028279', NULL, 574, 57, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2127, '0c3decff-39b8-4dcb-864e-231c743d1704', NULL, 574, 60, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2128, '7e79740b-60ce-403c-87b2-1abcaa0f89f3', NULL, 574, 63, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2129, 'cb8f0006-8d1d-4b1f-a614-15adc237b97f', NULL, 574, 66, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2130, '7c42a81d-443d-4393-a5be-0c2b8e3b61cb', NULL, 575, 57, '27112025 PROKER PGSD 2025 OK.xlsx', NULL, NULL, 'auditor2', NULL, '2025-11-27 15:30:51'),
        (2131, 'c04a0776-833b-459e-9e1f-eb2ae857627d', NULL, 575, 60, '2025-2030 JOBDESK DAN JOBSPE STRUKTUR ORGANISASI PRODI PGSD.pdf', NULL, NULL, 'auditor2', NULL, '2025-11-28 11:04:37'),
        (2132, 'f64fa33d-ae85-464d-8dd6-99c08d80336c', NULL, 575, 63, 'PGSD SOP Jobdesk dan Jobspek (2020-2027).pdf', NULL, NULL, 'auditor2', NULL, '2025-11-27 15:49:51'),
        (2133, 'ecfd3b65-7201-4a7e-abd7-f0fe7930e70b', NULL, 575, 66, 'Dokumen Closeing AMI KTS Tahun 2024.pdf', NULL, NULL, 'auditor2', NULL, '2025-11-27 15:23:15'),
        (2134, '882ac61e-f33b-4ba5-87fe-bddc33e5d1ab', NULL, 576, 57, 'Proker PPG 2025.xlsx', NULL, NULL, 'auditor2', NULL, '2025-11-20 08:33:13'),
        (2135, 'f61ea7d3-1ab8-499a-8d3c-168f4c8e30ef', NULL, 576, 60, 'Struktur Organisasi dan Job Deskripsi.pdf', NULL, NULL, 'auditor2', NULL, '2025-11-18 10:37:29'),
        (2136, 'c6de6fff-d785-46f1-b652-13c013222254', NULL, 576, 63, 'SOP.pdf', NULL, NULL, 'auditor2', NULL, '2025-11-18 10:40:57'),
        (2137, 'd966b64f-8d81-42c1-be69-96cd909e4e1a', NULL, 576, 66, 'berita_acara (6).pdf', NULL, NULL, 'auditor2', NULL, '2025-11-18 10:30:13'),
        (2138, 'bfccc789-0422-4537-be45-d57e89ff5c23', NULL, 578, 57, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2139, '906f0484-a84f-47cc-9cd5-d8c3cc634b7d', NULL, 578, 60, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2140, '64304fb3-4627-482d-b15e-b32ca3d4711a', NULL, 578, 63, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2141, '40aec6d7-89b7-4671-8a00-55eb459c8755', NULL, 578, 66, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2142, '07b2221a-27b7-448a-b312-9c97fe7bde7c', NULL, 579, 57, 'FORM PROKER TAHUN 2026.xlsx', NULL, NULL, 'auditor2', NULL, '2025-11-27 19:53:05'),
        (2143, 'a75a4231-7698-40f8-b81e-4a0dd918c5d4', NULL, 579, 60, 'Struktur Organisasi Program Studi Sastra Jepang FISIB Universitas Pakuan 2025.docx', NULL, NULL, 'auditor2', NULL, '2025-11-23 22:17:30'),
        (2144, '99090911-0f57-40d2-a596-8e6547b918a6', NULL, 579, 63, 'SOP Pelaksanaan Tugas Sesuai Dengan Job Description Dan Job Specification.docx', NULL, NULL, 'auditor2', NULL, '2025-11-23 22:46:48'),
        (2145, '9e18370e-e289-4b24-8f68-084099d81f27', NULL, 579, 66, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2146, '4de347a2-a33f-47f6-bef4-e60fe5b75755', NULL, 580, 57, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2147, '815d182d-be4d-4c53-905a-520994023c26', NULL, 580, 60, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2148, '88eb9630-81b3-405e-81a5-af9ac78e7762', NULL, 580, 63, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2149, '0ac7f652-9a3f-4951-81c1-3589ec67c16f', NULL, 580, 66, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2150, '1a705111-1c23-4d61-96ea-6d0865666eb6', NULL, 581, 57, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2151, '169859d9-4102-4085-8eb6-ea3814f5b042', NULL, 581, 60, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2152, 'eccb41e0-41e1-49c8-9844-42f5f8c54c1b', NULL, 581, 63, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2153, 'bb4bc5e7-32f0-4b2a-8fac-e899083aeaa2', NULL, 581, 66, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2154, '93353375-07ad-4606-b41a-b1a8241509d4', NULL, 582, 57, 'Proker Prodi Teknik Geologi_2025.xlsx', NULL, NULL, 'auditor2', NULL, '2025-11-25 13:58:53'),
        (2155, '1521aef4-078b-4185-9905-73c960bfe968', NULL, 582, 60, 'SK  STRUKTUR ORGANISASI JOBDESC DAN JOBSPEC FAKULTAS TEKNIK 2025_compressed.pdf', NULL, NULL, 'auditor2', NULL, '2025-11-26 17:11:59'),
        (2156, '70aadf30-37d4-417e-a92d-6d857600037f', NULL, 582, 63, '3. SOP Tupoksi.pdf', NULL, NULL, 'auditor2', NULL, '2025-11-26 17:14:49'),
        (2157, '9e147d91-d978-47a4-b988-74744463ae08', NULL, 582, 66, 'Closing KTS 2024 Teknik Geologi.pdf', NULL, NULL, 'auditor2', NULL, '2025-11-28 12:54:38'),
        (2158, '3066a2a2-6c76-4e3c-ab6c-597800b7621f', NULL, 583, 57, 'PROKER PRODI PWK 2025.xlsx', NULL, NULL, 'auditor2', NULL, '2025-11-20 11:28:01'),
        (2159, '616e312a-485c-4c9b-8e93-1bfc8125d174', NULL, 583, 60, '2. SK STRUKTUR ORGANISASI JOBDESC DAN JOBSPEC FAKULTAS TEKNIK-PWK_compressed.pdf', NULL, NULL, 'auditor2', NULL, '2025-11-26 11:31:21'),
        (2160, '8ef54937-1c1a-4915-93d3-bcbd27aae186', NULL, 583, 63, 'SOP Tupoksi.pdf', NULL, NULL, 'auditor2', NULL, '2025-11-18 06:48:13'),
        (2161, 'acb5db3c-661c-4dd6-a6fd-f166bb867624', NULL, 583, 66, 'Bukti Closing.pdf', NULL, NULL, 'auditor2', NULL, '2025-11-20 09:42:58'),
        (2162, 'c4c77e82-e997-4cb5-a506-49234e20d799', NULL, 584, 57, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2163, 'f7858092-a6d4-49f2-9a0d-5f64af699033', NULL, 584, 60, 'Struktur Organisasi, Job Description, Job Spesification, Disahkan Dengan SK_compressed.pdf', NULL, NULL, 'auditor2', NULL, '2025-11-26 14:06:29'),
        (2164, '429f53e9-33f3-474a-83f1-d83da5d5813b', NULL, 584, 63, 'SOP Tupoksi.pdf', NULL, NULL, 'auditor2', NULL, '2025-11-28 04:23:34'),
        (2165, '40b9919a-c0e2-4b4e-83f3-76385d6f33d0', NULL, 584, 66, 'Bukti Closing KTS 2024.pdf', NULL, NULL, 'auditor2', NULL, '2025-11-28 12:36:32'),
        (2166, '799b6e0d-02e5-46ee-a92a-ffaf69ef42e8', NULL, 585, 57, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2167, 'b82c4da0-5dc6-4376-afaa-dbd9ca35323d', NULL, 585, 60, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2168, '7c49ff1e-fdbe-4ad9-92c6-79f7d72e88f2', NULL, 585, 63, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2169, 'aa0ec288-3cf7-40a2-81bb-26b37bed9275', NULL, 585, 66, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2170, 'd3865b9e-5215-48b7-8950-539a10c47588', NULL, 586, 57, 'Proker Prodi Geodesi 2026.xlsx', NULL, NULL, 'auditor2', NULL, '2025-11-26 12:59:45'),
        (2171, '866b9a04-23b0-4396-b403-2f2f961795e0', NULL, 586, 60, 'SK  STRUKTUR ORGANISASI JOBDESC DAN JOBSPEC FAKULTAS T_compressed (1).pdf', NULL, NULL, 'auditor2', NULL, '2025-11-26 12:50:35'),
        (2172, '9cf74d07-aac2-4a66-b4ac-598e149dcdab', NULL, 586, 63, 'SOP Tupoksi.pdf', NULL, NULL, 'auditor2', NULL, '2025-11-26 12:58:23');
        INSERT INTO dokumen_tambahan (id, uuid, id_renstra_old, id_renstra, id_template_dokumen_tambahan, file, capaian_auditor, catatan_auditor, tugas, created_at, updated_at) VALUES
        (2173, '959cb292-aa65-4fb1-8ea7-b1bb83261c5f', NULL, 586, 66, 'berita acara+close kts 2024.pdf', NULL, NULL, 'auditor2', NULL, '2025-11-26 12:51:57'),
        (2174, 'c0748363-da08-4842-b243-c24e27a3e267', NULL, 587, 57, 'Template Proker 2025 Biologi.xlsx', NULL, NULL, 'auditor2', NULL, '2025-11-23 01:17:14'),
        (2175, '12eb1c40-1a85-4d1b-9661-8663b861deeb', NULL, 587, 60, 'Job Desk Prodi Biologi 2025.pdf', NULL, NULL, 'auditor2', NULL, '2025-11-26 14:11:16'),
        (2176, '5a640a1d-547d-4c0f-858b-b43a252a3a19', NULL, 587, 63, 'SOP Fakultas MIPA.pdf', NULL, NULL, 'auditor2', NULL, '2025-11-26 01:35:28'),
        (2177, 'bb6a3813-c78f-43e5-aaed-186bc3bac1cf', NULL, 587, 66, 'Dokumen Closed 2924.pdf', NULL, NULL, 'auditor2', NULL, '2025-11-23 01:21:40'),
        (2178, '88f628b0-f9ba-46f8-956f-ec305d80bec7', NULL, 588, 57, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2179, '993be3ed-1e13-4388-a737-d923cfa305d7', NULL, 588, 60, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2180, 'bb444442-2337-48d6-a875-e231e1b1047b', NULL, 588, 63, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2181, '21661a80-8088-4c7a-a3ec-7393ddb8585e', NULL, 588, 66, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2182, 'deac3f0c-6617-4aeb-a150-0deb7ad0207d', NULL, 589, 57, 'Evaluasi dan Jadwal Proker Matematika 2025.xlsx', NULL, NULL, 'auditor2', NULL, '2025-11-24 09:40:34'),
        (2183, '2d52cc2f-1f43-4469-874d-c16066239630', NULL, 589, 60, 'Jobdes dan spec Prodi 2025.pdf', NULL, NULL, 'auditor2', NULL, '2025-11-27 20:20:46'),
        (2184, 'eae35396-4e17-4805-b6f1-eabf8e032ac3', NULL, 589, 63, 'SOP Pelaksanaan Tugas Sesuai Dengan Job Description Dan Job Spesification----.pdf', NULL, NULL, 'auditor2', NULL, '2025-11-24 13:14:21'),
        (2185, 'af725a2e-ce7d-4bc3-844b-d0e1a4141302', NULL, 589, 66, 'Closing audit2024_merged.pdf', NULL, NULL, 'auditor2', NULL, '2025-11-24 12:31:11'),
        (2186, '598b8975-7e23-499d-82c1-f602f74c3d06', NULL, 590, 57, 'Proker_Ilmu_Komputer_2025_Edit_revisi_program 6 Des 24.xlsx', NULL, NULL, 'auditor2', NULL, '2025-11-19 21:27:19'),
        (2187, '005b8381-4e6e-4077-b654-bce174ddad39', NULL, 590, 60, 'Jobdes dan Job Spec Ilkom.pdf', NULL, NULL, 'auditor2', NULL, '2025-11-19 21:04:07'),
        (2188, '251710cf-0a86-4512-89ab-9f9270df8ff9', NULL, 590, 63, 'SOP Pengurus Divisi.pdf', NULL, NULL, 'auditor2', NULL, '2025-11-19 21:04:43'),
        (2189, 'b16918ed-d0bc-48a8-b43c-70ccfea58e98', NULL, 590, 66, 'kts_renstra (1).pdf', NULL, NULL, 'auditor2', NULL, '2025-11-19 21:07:13'),
        (2190, '47bee754-3584-477e-8b43-1ddc7f77a4f4', NULL, 591, 57, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2191, '5200e328-5b98-4021-9bbd-077aece093c1', NULL, 591, 60, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2192, '8ad23eb4-5e41-42a5-bce2-0a92c8318e85', NULL, 591, 63, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2193, 'd9820216-01dd-4712-9ee2-755746b510a6', NULL, 591, 66, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2194, 'cade62d5-3946-4342-b5a8-b743b7dac9d1', NULL, 592, 57, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2195, '80eab23b-3031-40dd-b901-37e8f6566041', NULL, 592, 60, '(2025) MP-S3. Struktur Organisasi. Job Desc. Job Spec_compressed.pdf', NULL, NULL, 'auditor2', NULL, '2025-11-28 11:40:56'),
        (2196, '18660953-5188-43ca-bad6-fa4225b8d127', NULL, 592, 63, 'Pascasarjana SOP Terbaru.pdf', NULL, NULL, 'auditor2', NULL, '2025-11-27 21:03:10'),
        (2197, '4a2e93a0-5641-4cbb-aa66-9043a65188c4', NULL, 592, 66, '(2024) MP-S3. KTS sudah di closed.pdf', NULL, NULL, 'auditor2', NULL, '2025-11-27 22:47:14'),
        (2198, '189bab18-30b4-4255-9ea6-2663986de1a1', NULL, 594, 57, 'Proker MP S2(2024-2025).xlsx', NULL, NULL, 'auditor2', NULL, '2025-11-20 15:27:06'),
        (2199, '09e69b71-eaf9-4daa-aea2-65515cc996f9', NULL, 594, 60, 'JOB DESC dan JOB DESPEC.pdf', NULL, NULL, 'auditor2', NULL, '2025-11-26 09:00:00'),
        (2200, 'e628c8b0-0723-4cbb-a181-3ce6cd463203', NULL, 594, 63, 'DAFTAR SOP TERBARU.pdf', NULL, NULL, 'auditor2', NULL, '2025-11-28 09:25:48'),
        (2201, '79826aaa-7bcb-4afe-8e89-970e0e2265ba', NULL, 594, 66, 'kts_renstra.pdf', NULL, NULL, 'auditor2', NULL, '2025-11-26 11:19:32'),
        (2202, '7668ea3a-d97c-492d-969a-c57ce3dc2f51', NULL, 595, 57, 'Laporan Monev Proker 2025.pdf', NULL, NULL, 'auditor2', NULL, '2025-11-26 15:59:47'),
        (2203, '0e2f22cb-2d07-4f68-9ef8-d7a4c194e2c3', NULL, 595, 60, 'B2.1. SK Dekan ttg Struktur Organisasi dan Job Desk Prodi ML.pdf', NULL, NULL, 'auditor2', NULL, '2025-11-13 14:15:55'),
        (2204, 'f5360022-7d46-4741-a217-3ed0ba24baa8', NULL, 595, 63, 'PEDOMAN STRUKTUR ORGANISASI, JOB DES dan JOB SPEC 2024.pdf', NULL, NULL, 'auditor2', NULL, '2025-11-13 14:17:22'),
        (2205, 'fdef4e0f-bc13-4ecd-b212-5d2261d67920', NULL, 595, 66, 'Closing Audit Mutu Internal 202420251119_13260832.pdf', NULL, NULL, 'auditor2', NULL, '2025-11-26 16:01:03'),
        (2206, '5a842bf5-806b-41ba-8b86-bdb690e3b348', NULL, 596, 57, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2207, 'c19d1e4c-e256-4cc8-b637-ff63354489eb', NULL, 596, 60, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2208, '7956df8c-21f8-45ec-91c2-b06974a76d9c', NULL, 596, 63, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2209, '30bdf3ee-4b47-4613-98ae-876bed77daaf', NULL, 596, 66, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2210, '83029328-0473-4acd-a485-c1795aee65a2', NULL, 597, 57, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2211, 'ba6370be-9b63-47cd-99d2-ea69208303c1', NULL, 597, 60, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2212, 'd80dcd77-8e33-41eb-8d63-247304606b06', NULL, 597, 63, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2213, '7651d0ac-7f55-4505-8bf9-01e1fec6c44b', NULL, 597, 66, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2214, '804eebb3-5690-4169-8986-421c00b49c2a', NULL, 598, 57, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2215, '4fae3e95-6ca4-41a2-9c79-5ed42908d11e', NULL, 598, 60, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2216, 'f0082989-d19c-42f4-8e16-1d285e673914', NULL, 598, 63, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2217, '7ac26c8e-b1a1-4b99-bfd0-26f779aeb523', NULL, 598, 66, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2218, '60796b02-d180-409b-ace3-71d61506a2a8', NULL, 599, 57, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2219, 'c3469af9-97e8-4b93-9983-60b4076b6bc2', NULL, 599, 60, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2220, '9ed13279-8a45-4941-97fc-c49b860c1126', NULL, 599, 63, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2221, '09563f96-e326-42b6-ad7c-9983b2481642', NULL, 599, 66, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2222, '33e34aef-d490-403e-8b40-35ef4538a659', NULL, 600, 57, 'PROKER PENDAS 2025_3 Desember 2024.xlsx', NULL, NULL, 'auditor2', NULL, '2025-11-26 14:55:43'),
        (2223, '3ab302ad-e437-44c5-b76f-86787822c4d9', NULL, 600, 60, 'Struktur Organisasi,Job.Desc, SK Kaprodi PENDAS 2025.pdf', NULL, NULL, 'auditor2', NULL, '2025-11-21 16:42:29'),
        (2224, 'ec0ab6fd-f665-4e39-931b-c1a92a325958', NULL, 600, 63, 'DAFTAR SOP TERBARU_merged.pdf', NULL, NULL, 'auditor2', NULL, '2025-11-28 09:55:56'),
        (2225, 'dcf611cb-2d8b-4388-8852-eea1bcc8a4f7', NULL, 600, 66, 'KTS CLOSED.docx', NULL, NULL, 'auditor2', NULL, '2025-11-28 09:56:14'),
        (2226, '7d7090a1-55f3-478e-9fdf-0858f4bab563', NULL, 601, 57, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2227, '33053958-046f-4854-bf69-8f8d9d23b328', NULL, 601, 60, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2228, 'd81171f1-1571-4b7c-95b4-7de5c56eafaf', NULL, 601, 63, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2229, '5bc8af81-b065-4df0-83ce-b1b16704ab44', NULL, 601, 66, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2230, '2b54f66c-b754-4ae1-b020-6d51e4e1ac28', NULL, 602, 57, 'anggaran 2025 prodi Akuntansi - pajak - MKP.xlsx', NULL, NULL, 'auditor2', NULL, '2025-11-25 12:19:10'),
        (2231, '3a0ceb43-75c7-487c-ba08-2c6df8e2e110', NULL, 602, 60, 'Jobdesk dan Jobspek 1 (2)-4-43.docx', NULL, NULL, 'auditor2', NULL, '2025-11-28 11:25:40'),
        (2232, 'ac4c3e3c-7068-4d2e-a603-bfddac847652', NULL, 602, 63, 'SOP-SOP Dokumen Tambahan.pdf', NULL, NULL, 'auditor2', NULL, '2025-11-25 13:33:38'),
        (2233, 'b34537ba-4257-4889-b93f-565b4104e225', NULL, 602, 66, 'Hasil Audit Mutu Internal 2024.pdf', NULL, NULL, 'auditor2', NULL, '2025-11-25 13:41:22'),
        (2234, '79084fed-39d5-4165-a80a-7a6436b15d38', NULL, 603, 57, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2235, '10bdfbdf-a892-4d25-a7ee-8063433b24e6', NULL, 603, 60, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2236, '9e8a3015-52d7-422e-af6e-168cd3e66621', NULL, 603, 63, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2237, '76afc6c1-1827-46f7-9408-1026f61fa077', NULL, 603, 66, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2238, 'b608a64c-31a1-4868-b23c-b26f31db0d15', NULL, 604, 57, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2239, '1ab64de7-ff33-4c80-9250-d95ee53e3109', NULL, 604, 60, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2240, '086112bb-6c35-4f97-9d45-5e4ec21f1c33', NULL, 604, 63, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2241, 'a53fae7f-bb23-4a9d-b512-cdea31c4d080', NULL, 604, 66, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2242, '982afdaa-26f3-42dd-95e9-1fcc29138282', NULL, 605, 57, 'Revisian Proker MI 2025.xlsx', NULL, NULL, 'auditor2', NULL, '2025-11-26 12:40:38'),
        (2243, '404bf9b8-f831-4a05-b053-5c5e62d78138', NULL, 605, 60, 'Jobdesk dan Jobspek 1 (2)-4-43.pdf', NULL, NULL, 'auditor2', NULL, '2025-11-28 09:37:31'),
        (2244, '64987f2d-cfa5-4c4c-8cd9-31942fe123d6', NULL, 605, 63, '1. SOP Penyusunan Struktur Organisasi Prodi (1).pdf', NULL, NULL, 'auditor2', NULL, '2025-11-28 09:59:49'),
        (2245, '61549c47-a25b-41c5-bf9a-71f316190ca5', NULL, 605, 66, 'Closing Audit Internal_2024.pdf', NULL, NULL, 'auditor2', NULL, '2025-11-26 12:39:57'),
        (2246, '5023dc83-ad47-4252-9a68-9beb9d1ac6a6', NULL, 611, 57, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2247, '61566338-7a11-4cf1-a407-fe694a330285', NULL, 611, 60, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2248, 'e16f5775-eb89-4792-b216-9620062c20f0', NULL, 611, 63, NULL, NULL, NULL, 'auditor2', NULL, NULL),
        (2249, '206f78a9-5f1a-412a-bfea-d9d0add1e5bb', NULL, 611, 66, NULL, NULL, NULL, 'auditor2', NULL, NULL);

        INSERT INTO jenis_file_renstra (id, uuid, nama, created_at, updated_at) VALUES
        (1, '14212231-792f-4935-bb1c-9a38695a4b6b', 'Program Kerja Sesuai Dengan Template 2024 disertai Monev', NULL, '2024-10-08 13:35:36'),
        (2, '1a353e22-1111-4fc5-96c1-a2ed2877a6a4', 'Struktur Organisasi Disertai Dengan Job Description, Job Spesification, Dan Disahkan Dengan SK', NULL, '2024-10-08 13:36:55'),
        (3, '08a5e4cc-1a30-4080-95ad-127abf8819f5', 'SOP Pelaksanaan Tugas Sesuai Dengan Job Description Dan Job Spesification', NULL, '2024-10-08 13:37:14'),
        (4, '3523b28b-829a-4548-96e0-459bf4f14dea', 'Hasil/Catatan Audit/KTS Sebelumnya (2023) Telah Diselesaikan (Closed)', NULL, '2024-10-08 13:37:47'),
        (14, '6aee7cd5-da31-4735-9243-8c19aa7497c0', 'Program Kerja Sesuai Dengan Template 2025 disertai Monev', '2025-10-13 15:14:06', '2025-10-13 15:14:06'),
        (15, '84066942-1f2d-44b0-be66-8b87cdab6e91', 'Hasil/Catatan Audit/KTS Sebelumnya (2024) Telah Diselesaikan (Closed)', '2025-10-13 15:14:25', '2025-10-13 15:14:25');

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
}
