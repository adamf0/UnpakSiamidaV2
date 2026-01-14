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

func setupFakultasUnitMySQL(t *testing.T) (*gorm.DB, func()) {
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

        INSERT INTO m_fakultas (kode_fakultas, kode_pt, nama_fakultas, pejabat, jabatan, wakil_pejabat, wakil_pejabat_adm, logo) VALUES
        ('01', '041004', 'HUKUM', '0410067306', 'H', '0417086801', '0414106202', ''),
        ('02', '041004', 'EKONOMI DAN BISNIS', '0415056901', 'H', '0425097604', '0414109101', ''),
        ('03', '041004', 'KIP', '0416076701', 'H', '0413018604', '0415128202', ''),
        ('04', '041004', 'ISIB', '0416068002', 'H', '0401098708', '0413026701', ''),
        ('05', '041004', 'TEKNIK', '0428106901', 'H', '0424076601', '0417097601', ''),
        ('06', '041004', 'MIPA', '0406097101', 'H', '0404107405', '0407027501', ''),
        ('07', '041004', 'PASCASARJANA', '0403055801', 'H', '0427017501', '0025057515', ''),
        ('08', '041004', 'VOKASI', '0413117601', 'H', '0402047301', '0404047007', '');

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

        CREATE TABLE sijamu_fakultas_unit (
            id int(11) NOT NULL,
            uuid varchar(36) DEFAULT NULL,
            kode_fakultas char(9) DEFAULT NULL,
            kode_prodi char(10) DEFAULT NULL,
            nama varchar(100) DEFAULT NULL,
            id_m_prodi int(11) DEFAULT NULL,
            standalone tinyint(4) DEFAULT 0
        );

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

	cleanup := func() {
		sqlDB, _ := gdb.DB()
		sqlDB.Close()
		mysqlC.Terminate(ctx)
	}

	return gdb, cleanup
}
