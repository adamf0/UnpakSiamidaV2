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

func setupMySQL(t *testing.T) (*gorm.DB, func()) {
    ctx := context.Background()

    req := testcontainers.ContainerRequest{
        Image:        "mysql:8.0",
        Env: map[string]string{
            "MYSQL_ROOT_PASSWORD": "pass",
            "MYSQL_DATABASE":      "testdb",
        },
        ExposedPorts: []string{"3306/tcp"},
        WaitingFor: wait.ForListeningPort("3306/tcp").
            WithStartupTimeout(90 * time.Second),
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
            id int(11) NOT NULL AUTO_INCREMENT,
            uuid varchar(36) DEFAULT NULL,
            nama text NOT NULL,
            created_at datetime DEFAULT NULL,
            updated_at datetime DEFAULT NULL,
            PRIMARY KEY (id)
        );
        INSERT INTO master_standar_renstra VALUES 
            (1,UUID(),'Standar Kompetensi Lulusan','2024-08-09 09:35:07','2024-08-25 19:24:41'),
            (4,UUID(),'Standar Isi Pembelajaran','2024-08-25 19:22:07','2024-08-25 19:22:07'),
            (7,UUID(),'Standar Proses Pembelajaran','2024-08-25 19:22:14','2024-08-25 19:22:14'),
            (10,UUID(),'Standar Penilaian Pembelajaran','2024-08-25 19:22:20','2024-08-25 19:22:20'),
            (13,UUID(),'Standar Dosen dan Tenaga Kependidikan','2024-08-25 19:22:33','2024-08-25 19:22:33'),
            (25,UUID(),'Standar Sarana dan Prasarana Pembelajaran','2024-09-09 14:25:14','2024-09-09 14:25:14'),
            (28,UUID(),'Standar Pengelolaan Pembelajaran','2024-09-09 14:25:28','2024-09-09 14:25:28'),
            (118,UUID(),'Standar Pembiayaan Pembelajaran','2024-09-10 09:12:37','2024-09-10 09:12:37'),
            (121,UUID(),'Standar Hasil Penelitian','2024-09-10 09:12:42','2024-09-10 09:12:42'),
            (124,UUID(),'Standar Isi Penelitian','2024-09-10 09:12:49','2024-09-10 09:12:49'),
            (127,UUID(),'Standar Proses Penelitian','2024-09-10 09:12:54','2024-09-10 09:12:54'),
            (130,UUID(),'Standar Penilaian Penelitian','2024-09-10 09:13:10','2024-09-10 09:13:10'),
            (133,UUID(),'Standar Peneliti','2024-09-10 09:13:24','2024-09-10 09:13:24'),
            (136,UUID(),'Standar Sarana dan Prasarana Penelitian','2024-09-10 09:13:33','2024-09-10 09:13:33'),
            (139,UUID(),'Standar Pengelolaan Penelitian','2024-09-10 09:13:41','2024-09-10 09:13:41'),
            (142,UUID(),'Standar Pendanaan dan Pembiayaan Penelitian','2024-09-10 09:13:58','2024-09-10 09:13:58'),
            (145,UUID(),'Standar Hasil Pengabdian Kepada Masyarakat','2024-09-10 09:23:17','2024-09-10 09:23:17'),
            (148,UUID(),'Standar Isi Pengabdian Kepada Masyarakat','2024-09-10 09:23:27','2024-09-10 09:23:27'),
            (151,UUID(),'Standar Proses Pengabdian Kepada Masyarakat','2024-09-10 09:23:38','2024-09-10 09:23:38'),
            (154,UUID(),'Standar Penilaian Pengabdian Kepada Masyarakat','2024-09-10 09:23:46','2024-09-10 09:23:46'),
            (157,UUID(),'Standar Pelaksana Pengabdian Kepada Masyarakat','2024-09-10 09:24:01','2024-09-10 09:24:01'),
            (160,UUID(),'Standar Sarana dan Prasarana Pengabdian Kepada Masyarakat','2024-09-10 09:25:16','2024-09-10 09:25:16'),
            (163,UUID(),'Standar Pengelolaan Pengabdian Kepada Masyarakat','2024-09-10 09:25:35','2024-09-10 09:25:35'),
            (166,UUID(),'Standar Pendanaan dan Pembiayaan PkM','2024-09-10 09:25:55','2024-09-10 09:25:55'),
            (169,UUID(),'Standar Visi Misi (kriteria 1)','2024-09-10 09:26:02','2024-09-10 09:26:02'),
            (172,UUID(),'Standar Ketaatan Pada Peraturan Perundang-Undangan (kriteria 2)','2024-09-10 09:26:08','2024-09-10 09:26:08'),
            (175,UUID(),'Standar Pengelolaan Tata Pamong (kriteria 2)','2024-09-10 09:26:12','2024-09-10 09:26:12'),
            (178,UUID(),'Standar Penjanjian Kerjasama (kriteria 2)','2024-09-10 09:26:17','2024-09-10 09:26:17'),
            (181,UUID(),'Standar Pemeliharaan/Peningkatan Jumlah Peminat/Pendaftar (kriteria 3)','2024-09-10 09:26:29','2024-09-10 09:26:29'),
            (184,UUID(),'Standar Layanan Kemahasiswaan (kriteria 3)','2024-09-10 09:26:40','2024-09-10 09:26:40'),
            (187,UUID(),'Standar Prestasi Mahasiswa (kriteria 3 dan 9)','2024-09-10 09:26:44','2024-09-10 09:26:44'),
            (190,UUID(),'Standar Pengelolaan Keuangan (kriteria 5)','2024-09-10 09:27:13','2024-09-10 09:27:13'),
            (193,UUID(),'Standar Sarana Prasarana Umum (kriteria 5)','2024-09-10 09:27:18','2024-09-10 09:27:18'),
            (196,UUID(),'Standar Sistem Informasi (kriteria 5)','2024-09-10 09:27:24','2024-09-10 09:27:24'),
            (199,UUID(),'Standar Pembiayaan MBKM (kriteria 5 dan 6)','2024-09-10 09:27:32','2024-09-10 09:27:32'),
            (202,UUID(),'Standar Pelaksanaan MBKM (kriteria 6)','2024-09-10 09:27:37','2024-09-10 09:27:37'),
            (205,UUID(),'Standar Pemeliharaan/Peningkatan Jumlah Lulusan (kriteria 6 dan 9)','2024-09-10 09:27:48','2024-09-10 09:27:48'),
            (208,UUID(),'Standar Tracer Study (kriteria 9)','2024-09-10 09:27:55','2024-09-10 09:27:55'),
            (211,UUID(),'Standar Inovasi dan Inkubator Bisnis (kriteria 9)','2024-09-10 09:28:01','2024-09-10 09:28:01');

        CREATE TABLE master_indikator_renstra (
            id int(11) NOT NULL AUTO_INCREMENT,
            uuid varchar(36) DEFAULT NULL,
            id_master_standar int(11) DEFAULT NULL,
            indikator text NOT NULL,
            parent int(11) DEFAULT NULL,
            tahun varchar(4) NOT NULL,
            tipe_target text DEFAULT NULL,
            operator varchar(5) DEFAULT NULL,
            created_at datetime DEFAULT NULL,
            updated_at datetime DEFAULT NULL,
            PRIMARY KEY (id)
        ) ENGINE=InnoDB AUTO_INCREMENT=477 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

        INSERT INTO master_indikator_renstra (
            uuid,
            id_master_standar,
            indikator,
            parent,
            tahun,
            tipe_target,
            operator
        )
        VALUES
            (16,'186f2427-8bdd-42d9-a757-65808f364eeb',1,'Lulusan memiliki sertifikat kompetensi atau Bahasa asing',NULL,2024,'numerik',NULL,'2024-08-25 19:25:05','2024-09-21 07:18:42'),
            (79,UUID(),1,'Lulusan bekerja pada ruang lingkup nasional',NULL,2024,'numerik',NULL,'2024-09-08 10:45:01','2024-09-23 08:17:48'),
            (82,UUID(),4,'Jumlah mata kuliah yang mengintegrasikan hasil penelitian dan PKM ke dalam perkuliahan',NULL,2024,'numerik',NULL,'2024-09-08 10:45:21','2024-09-23 08:17:59'),
            (85,UUID(),4,'Persentase Prodi Vokasi dan S1 yang mengimplementasikan kurikulum MBKM',NULL,2024,'numerik',NULL,'2024-09-09 14:40:00','2024-10-18 12:35:42'),
            (15,UUID()7,7,'Proses pembelajaran berbasis PBL dan/atau PjBL tiap semester',NULL,2024,'numerik',NULL,'2024-09-10 09:31:20','2024-09-23 08:21:32'),
            (16,UUID()0,7,'Jumlah exchange lecture di setiap Prodi sejenis',NULL,2024,'numerik',NULL,'2024-09-10 09:31:38','2024-09-23 08:22:38'),
            (16,UUID()3,7,'Perkuliahan bilingual di setiap Prodi',NULL,2024,'numerik',NULL,'2024-09-10 09:32:47','2024-09-23 08:22:49'),
            (16,UUID()6,7,'Peningkatan jumlah kegiatan sebagai penunjang suasana akademik',NULL,2024,'numerik',NULL,'2024-09-10 09:33:13','2024-09-23 08:23:09'),
            (16,UUID()9,7,'Persentase pembelajaran dengan moda campuran',NULL,2024,'numerik',NULL,'2024-09-10 09:33:36','2024-09-23 08:23:49'),
            (17,UUID()2,7,'blended/flip learning (Prodi)',169,2024,'numerik',NULL,'2024-09-10 09:33:54','2024-09-23 08:23:55'),
            (17,UUID()5,7,'hybrid learning (Fakultas)',169,2024,'numerik',NULL,'2024-09-10 09:34:09','2024-09-23 08:24:09'),
            (17,UUID()8,10,'IPK rata-rata lulusan/tahun (S1 ≥ 3,00; S2 ≥ 3,25; dan S3 ≥ 3,50)',NULL,2024,'range','>=','2024-09-10 09:35:10','2024-09-23 08:24:59'),
            (18,UUID()1,10,'Portofolio mahasiswa melalui kurikulum OBE',NULL,2024,'numerik',NULL,'2024-09-10 09:35:32','2024-09-23 08:25:10'),
            (18,UUID()4,13,'Dosen mampu Berbahasa Inggris dengan sangat baik',NULL,2024,'numerik',NULL,'2024-09-10 09:44:07','2024-09-23 08:25:20'),
            (18,UUID()7,13,'Peningkatan jumlah Guru Besar tiap Fakultas',NULL,2024,'numerik',NULL,'2024-09-10 09:50:14','2024-09-23 08:25:33'),
            (19,UUID()0,13,'Peningkatan jumlah Tendik registered dan memiliki jabatan fungsional',NULL,2024,'numerik',NULL,'2024-09-10 09:52:19','2024-09-24 22:06:40'),
            (19,UUID()3,25,'Tersedia Smart Class',NULL,2024,'numerik',NULL,'2024-09-10 09:53:17','2024-09-23 08:28:00'),
            (19,UUID()6,28,'Penyusunan kurikulum berbasi OBE',NULL,2024,'kategori',NULL,'2024-09-10 09:54:06','2024-09-23 08:28:11'),
            (19,UUID()9,28,'Implementasi OBE',NULL,2024,'numerik',NULL,'2024-09-10 09:54:18','2024-09-23 08:28:57'),
            (20,UUID()2,28,'Penyusunan kurikulum MBKM',NULL,2024,'kategori',NULL,'2024-09-10 09:54:39','2024-09-23 08:29:09'),
            (20,UUID()5,28,'Implementasi MBKM',NULL,2024,'kategori',NULL,'2024-09-10 09:59:06','2024-09-23 08:29:20'),
            (20,UUID()8,118,'Audit keuangan internal',NULL,2024,'kategori',NULL,'2024-09-10 09:59:41','2024-09-23 08:30:59'),
            (21,UUID()1,121,'Peningkatan jumlah publikasi dosen pada jurnal nasional minimal Sinta 3 dan jurnal internasional bereputasi',NULL,2024,'numerik',NULL,'2024-09-10 09:59:57','2024-09-23 08:31:50'),
            (21,UUID()4,121,'Penelitian kolektif atau kolaboratif dosen melibatkan mahasiswa',NULL,2024,'numerik',NULL,'2024-09-10 10:00:24','2024-09-23 08:31:56'),
            (21,UUID()7,121,'Kolaborasi publikasi dengan institusi dalam dan luar negeri',NULL,2024,'numerik',NULL,'2024-09-10 10:00:44','2024-09-23 08:32:02'),
            (22,UUID()0,121,'Hilirisasi hasil penelitian',NULL,2024,'numerik',NULL,'2024-09-10 10:01:24','2024-09-23 08:32:08'),
            (22,UUID()3,121,'Hasil penelitian diadopsi oleh masyarakat atau Teknologi Tepat Guna (TTG)',NULL,2024,'numerik',NULL,'2024-09-10 10:01:37','2024-09-23 08:32:14'),
            (22,UUID()6,121,'Hasil penelitian terstandarisasi atau tersertifikasi',NULL,2024,'numerik',NULL,'2024-09-10 10:01:54','2024-09-23 08:32:27'),
            (22,UUID()9,124,'Terdapat bukti kesesuaian isi penelitian dengan roadmap penelitian dosen',NULL,2024,'kategori',NULL,'2024-09-10 10:02:21','2024-09-23 08:33:45'),
            (23,UUID()2,124,'Penelitian dosen dilakukan secara multi disiplin ilmu',NULL,2024,'numerik',NULL,'2024-09-10 10:02:33','2024-09-23 08:35:02'),
            (23,UUID()5,124,'Research group di level Prodi/Fakultas/Universitas dengan legalisasi Rektor/Dekan',NULL,2024,'numerik',NULL,'2024-09-10 10:02:47','2024-09-23 08:35:13'),
            (23,UUID()8,124,'Research collaboration dengan institusi dalam dan luar negeri',NULL,2024,'numerik',NULL,'2024-09-10 10:02:59','2024-09-23 08:35:23'),
            (24,UUID()1,127,'Kegiatan penelitian dosen sesuai dengan roadmap',NULL,2024,'numerik',NULL,'2024-09-10 10:03:57','2024-09-23 08:35:28'),
            (24,UUID()4,130,'Reviewer berlisensi/bersertifikat',NULL,2024,'numerik',NULL,'2024-09-10 10:04:15','2024-09-23 08:35:36'),
            (24,UUID()7,133,'Setiap peneliti memiliki roadmap penelitian',NULL,2024,'numerik',NULL,'2024-09-10 10:04:27','2024-09-23 08:35:43'),
            (25,UUID()0,136,'Standarisasi Sarpras penelitian (KAN/SNI/ISO)',NULL,2024,'kategori',NULL,'2024-09-10 10:04:47','2024-09-24 22:38:59'),
            (25,UUID()3,139,'Sistem pengelolaan penelitian dan PKM berbasis informasi',NULL,2024,'kategori',NULL,'2024-09-10 10:05:10','2024-09-24 22:39:15'),
            (25,UUID()6,142,'Dosen mendapatkan hibah penelitian eksternal',NULL,2024,'numerik',NULL,'2024-09-10 10:10:21','2024-09-23 08:40:44'),
            (25,UUID()9,145,'Kegiatan PKM dosen melibatkan mahasiswa',NULL,2024,'numerik',NULL,'2024-09-10 10:10:58','2024-09-23 08:40:50'),
            (26,UUID()2,145,'Kolaborasi publikasi dengan institusi dalam dan luar negeri',NULL,2024,'numerik',NULL,'2024-09-10 10:12:44','2024-09-23 08:40:56'),
            (26,UUID()5,148,'Kegiatan PKM dosen dilakukan secara multi disiplin ilmu',NULL,2024,'numerik',NULL,'2024-09-10 10:13:56','2024-09-23 08:41:02'),
            (26,UUID()8,148,'PKM group di level Prodi/Fakultas/Universitas dengan legalisasi Rektor/Dekan',NULL,2024,'numerik',NULL,'2024-09-10 10:14:13','2024-09-23 08:41:18'),
            (27,UUID()1,148,'PKM collaboration dengan institusi dalam dan luar negeri',NULL,2024,'numerik',NULL,'2024-09-10 10:14:26','2024-09-23 08:41:24'),
            (27,UUID()4,151,'Kegiatan PKM dosen sesuai dengan roadmap',NULL,2024,'numerik',NULL,'2024-09-10 10:15:08','2024-09-23 08:41:30'),
            (27,UUID()7,154,'Reviewer berlisensi/bersertifikat',NULL,2024,'numerik',NULL,'2024-09-10 10:15:30','2024-09-23 08:41:36'),
            (28,UUID()0,157,'Setiap pelaksana PKM memiliki roadmap',NULL,2024,'numerik',NULL,'2024-09-10 10:16:01','2024-09-23 08:41:41'),
            (28,UUID()3,160,'Terdapat laboratorium yang memfasilitasi kegiatan PKM',NULL,2024,'numerik',NULL,'2024-09-10 10:16:31','2024-09-23 08:41:46'),
            (28,UUID()6,163,'Sistem pengelolaan penelitian dan PKM berbasis informasi',NULL,2024,'kategori',NULL,'2024-09-10 10:16:59','2024-09-25 06:27:45'),
            (28,UUID()9,166,'Dosen mendapatkan hibah PkM eksternal',NULL,2024,'numerik',NULL,'2024-09-10 10:17:37','2024-09-23 08:41:54'),
            (29,UUID()2,169,'Dilakukan peninjauan VMTS tiap 5 tahun',NULL,2024,'kategori',NULL,'2024-09-10 10:18:05','2024-09-23 08:43:03'),
            (29,UUID()5,169,'Akreditasi Internasional',NULL,2024,'numerik',NULL,'2024-09-10 10:18:28','2025-02-05 14:42:44'),
            (29,UUID()8,169,'Tingkat keterpahamanan VMTS',NULL,2024,'numerik',NULL,'2024-09-10 10:18:52','2024-09-23 08:43:26'),
            (30,UUID()1,172,'Implementasi standar per tahun',NULL,2024,'kategori',NULL,'2024-09-10 10:24:51','2024-09-25 06:28:22'),
            (30,UUID()4,175,'Tersedia program retensi untuk dosen dan tenaga kependidikan',NULL,2024,'kategori',NULL,'2024-09-10 10:25:34','2024-09-25 06:28:41'),
            (30,UUID()7,175,'Terdapat peningkatan Prodi terakreditasi unggul',NULL,2024,'numerik',NULL,'2024-09-10 10:26:09','2025-02-05 14:44:41'),
            (31,UUID()0,175,'Automasi penjaminan mutu',NULL,2024,'numerik',NULL,'2024-09-10 10:26:31','2025-02-05 14:45:37'),
            (31,UUID()3,178,'Prodi mengimplementasikan kerjasama untuk proses tridharma',NULL,2024,'numerik',NULL,'2024-09-10 10:27:52','2024-09-23 08:43:44'),
            (31,UUID()6,178,'Jumlah kerjasama internasional dalam lingkup tridharma',NULL,2024,'numerik',NULL,'2024-09-10 10:28:11','2024-09-23 08:45:43'),
            (31,UUID()9,181,'Persentase kenaikan jumlah mahasiswa baru',NULL,2024,'kategori',NULL,'2024-09-10 10:28:39','2024-09-25 06:29:01'),
            (32,UUID()2,181,'Terdapat mekanisme penentuan daya tampung Prodi',NULL,2024,'kategori',NULL,'2024-09-10 10:29:03','2024-09-25 06:29:12'),
            (32,UUID()5,181,'Tersedia formasi khusus untuk mahasiswa disabilitas',NULL,2024,'numerik',NULL,'2024-09-10 10:29:26','2025-02-04 11:44:33'),
            (32,UUID()8,181,'Terdapat peningkatan jumlah mahasiswa asing',NULL,2024,'numerik',NULL,'2024-09-10 10:31:59','2024-09-23 08:46:08'),
            (33,UUID()1,184,'Jumlah mahasiswa yang mendapatkan beasiswa eksternal',NULL,2024,'kategori',NULL,'2024-09-10 10:37:33','2024-09-25 08:23:38'),
            (33,UUID()4,184,'Tersedianya Unit Layanan Disabilitas (ULD) bagi mahasiswa',NULL,2024,'numerik',NULL,'2024-09-10 10:38:04','2025-02-04 13:21:32'),
            (33,UUID()7,187,'Terdapat peningkatan jumlah prestasi mahasiswa baik akademik maupun non akademik dalam ruang lingkup nasional dan internasional',NULL,2024,'numerik',NULL,'2024-09-10 10:39:43','2024-09-23 08:46:23'),
            (34,UUID()0,187,'Jumlah penelitian mahasiswa yang mendapatkan HKI',NULL,2024,'numerik',NULL,'2024-09-10 10:49:40','2024-09-23 08:46:33'),
            (34,UUID()3,190,'Tersedia pedoman dan SOP pengelolaan keuangan dan dilaporkan setiap tahun anggaran',NULL,2024,'kategori',NULL,'2024-09-10 10:50:03','2024-09-25 08:23:57'),
            (34,UUID()6,190,'Tersedia mekanisme alokasi pendanaan tridharma',NULL,2024,'kategori',NULL,'2024-09-10 10:50:27','2024-09-25 08:24:04'),
            (34,UUID()9,190,'Tersedia kebijakan biaya pendidikan untuk mahasiswa berpotensi akademik tapi kurang mampu',NULL,2024,'kategori',NULL,'2024-09-10 10:50:43','2024-09-25 08:24:12'),
            (35,UUID()2,193,'Peningkatan jumlah sarana dan prasarana yang menjadi Income Generator Unit (IGU)',NULL,2024,'kategori',NULL,'2024-09-10 10:51:39','2024-09-25 08:24:30'),
            (35,UUID()5,193,'Implementasi green campus pada pengembangan sarana dan prasarana',NULL,2024,'kategori',NULL,'2024-09-10 10:51:55','2024-09-25 08:24:39'),
            (35,UUID()8,196,'Tersedia sistem informasi terintegrasi untuk akademik, SDM, keuangan, Sarpras, kemahasiswaan dan alumni, kerjasama, riset dan inovasi',NULL,2024,'kategori',NULL,'2024-09-10 10:52:10','2024-09-25 08:24:45'),
            (36,UUID()1,199,'Tersedia kebijakan biaya pembelajaran peserta MBKM untuk mahasiswa yang memiliki prestasi akademik tapi kurang mampu',NULL,2024,'numerik',NULL,'2024-09-10 10:52:43','2025-02-04 14:56:13'),
            (36,UUID()4,202,'Tersedia fasilitas kegiatan magang/praktik industri maupun kegiatan proyek desa yang sesuai',NULL,2024,'numerik',NULL,'2024-09-10 10:53:06','2024-09-23 08:47:05'),
            (36,UUID()7,205,'Percepatan masa studi',NULL,2024,'range','>=','2024-09-10 10:53:28','2024-10-10 14:12:42'),
            (37,UUID()0,208,'Terdapat peningkatan respon pengisian dari pengguna lulusan',NULL,2024,'numerik',NULL,'2024-09-10 10:53:54','2024-09-23 08:47:20'),
            (37,UUID()3,208,'Terdapat peningkatan respon pengisian dari jumlah lulusan',NULL,2024,'numerik',NULL,'2024-09-10 10:54:15','2024-09-23 08:47:26'),
            (37,UUID()6,211,'Terdapat peningkatan jumlah karya inovasi mahasiswa per tahun',NULL,2024,'numerik',NULL,'2024-09-10 10:54:35','2024-09-23 08:47:31'),
            (38,UUID()5,118,'Audit keuangan internal dan eksternal tiap triwulan',NULL,2024,'kategori',NULL,'2024-09-18 10:02:04','2024-09-23 08:29:46'),
            (39,UUID()2,10,'10.	IPK rata-rata lulusan/tahun (S1 ≥ 3,00; S2 ≥ 3,25; dan S3 ≥ 3,50) (Lembaga dan Unit)',NULL,2024,'kategori',NULL,'2025-02-03 14:49:48','2025-02-03 14:49:48'),
            (39,UUID()9,1,'Lulusan memiliki sertifikat kompetensi atau Bahasa asing',NULL,2025,'numerik',NULL,NULL,NULL),
            (40,UUID()0,1,'Lulusan bekerja pada ruang lingkup nasional',NULL,2025,'numerik',NULL,NULL,NULL),
            (40,UUID()1,4,'Jumlah mata kuliah yang mengintegrasikan hasil penelitian dan PKM ke dalam perkuliahan',NULL,2025,'numerik',NULL,NULL,NULL),
            (40,UUID()2,4,'Persentase Prodi Vokasi dan S1 yang mengimplementasikan kurikulum MBKM',NULL,2025,'numerik',NULL,NULL,NULL),
            (40,UUID()3,7,'Proses pembelajaran berbasis PBL dan/atau PjBL tiap semester',NULL,2025,'numerik',NULL,NULL,NULL),
            (40,UUID()4,7,'Jumlah exchange lecture di setiap Prodi sejenis',NULL,2025,'numerik',NULL,NULL,NULL),
            (40,UUID()5,7,'Perkuliahan bilingual di setiap Prodi',NULL,2025,'numerik',NULL,NULL,NULL),
            (40,UUID()6,7,'Peningkatan jumlah kegiatan sebagai penunjang suasana akademik',NULL,2025,'numerik',NULL,NULL,NULL),
            (40,UUID()7,7,'Persentase pembelajaran dengan moda campuran',NULL,2025,'numerik',NULL,NULL,NULL),
            (40,UUID()8,7,'blended/flip learning (Prodi)',407,2025,'numerik',NULL,NULL,NULL),
            (40,UUID()9,7,'hybrid learning (Fakultas)',407,2025,'numerik',NULL,NULL,NULL),
            (41,UUID()0,10,'IPK rata-rata lulusan/tahun (Diploma ≥ 3,00; S1 ≥ 3,00; S2 ≥ 3,25; dan S3 ≥ 3,50)',NULL,2025,'range','>=',NULL,'2025-10-23 09:19:49'),
            (41,UUID()1,10,'Portofolio mahasiswa melalui kurikulum OBE',NULL,2025,'numerik',NULL,NULL,NULL),
            (41,UUID()2,13,'Dosen mampu Berbahasa Inggris dengan sangat baik',NULL,2025,'numerik',NULL,NULL,NULL),
            (41,UUID()3,13,'Peningkatan jumlah Guru Besar tiap Fakultas',NULL,2025,'numerik',NULL,NULL,NULL),
            (41,UUID()4,13,'Peningkatan jumlah Tendik registered dan memiliki jabatan fungsional',NULL,2025,'numerik',NULL,NULL,NULL),
            (41,UUID()5,25,'Tersedia Smart Class',NULL,2025,'numerik',NULL,NULL,NULL),
            (41,UUID()6,28,'Penyusunan kurikulum berbasi OBE',NULL,2025,'kategori',NULL,NULL,NULL),
            (41,UUID()7,28,'Implementasi OBE',NULL,2025,'numerik',NULL,NULL,NULL),
            (41,UUID()8,28,'Penyusunan kurikulum MBKM',NULL,2025,'kategori',NULL,NULL,NULL),
            (41,UUID()9,28,'Implementasi MBKM',NULL,2025,'kategori',NULL,NULL,NULL),
            (42,UUID()0,121,'Peningkatan jumlah publikasi dosen pada jurnal nasional minimal Sinta 3 dan jurnal internasional bereputasi',NULL,2025,'numerik',NULL,NULL,NULL),
            (42,UUID()1,121,'Penelitian kolektif atau kolaboratif dosen melibatkan mahasiswa',NULL,2025,'numerik',NULL,NULL,NULL),
            (42,UUID()2,121,'Kolaborasi publikasi dengan institusi dalam dan luar negeri',NULL,2025,'numerik',NULL,NULL,NULL),
            (42,UUID()3,121,'Hilirisasi hasil penelitian',NULL,2025,'numerik',NULL,NULL,NULL),
            (42,UUID()4,121,'Hasil penelitian diadopsi oleh masyarakat atau Teknologi Tepat Guna (TTG)',NULL,2025,'numerik',NULL,NULL,NULL),
            (42,UUID()5,121,'Hasil penelitian terstandarisasi atau tersertifikasi',NULL,2025,'numerik',NULL,NULL,NULL),
            (42,UUID()6,124,'Terdapat bukti kesesuaian isi penelitian dengan roadmap penelitian dosen',NULL,2025,'kategori',NULL,NULL,NULL),
            (42,UUID()7,124,'Penelitian dosen dilakukan secara multi disiplin ilmu',NULL,2025,'numerik',NULL,NULL,NULL),
            (42,UUID()8,124,'Research group di level Prodi/Fakultas/Universitas dengan legalisasi Rektor/Dekan',NULL,2025,'numerik',NULL,NULL,NULL),
            (42,UUID()9,124,'Research collaboration dengan institusi dalam dan luar negeri',NULL,2025,'numerik',NULL,NULL,NULL),
            (43,UUID()0,127,'Kegiatan penelitian dosen sesuai dengan roadmap',NULL,2025,'numerik',NULL,NULL,NULL),
            (43,UUID()1,130,'Reviewer berlisensi/bersertifikat (Penelitian)',NULL,2025,'numerik',NULL,NULL,'2025-10-22 14:43:03'),
            (43,UUID()2,133,'Setiap peneliti memiliki roadmap penelitian',NULL,2025,'numerik',NULL,NULL,NULL),
            (43,UUID()3,136,'Standarisasi Sarpras penelitian (KAN/SNI/ISO)',NULL,2025,'kategori',NULL,NULL,NULL),
            (43,UUID()4,139,'Sistem pengelolaan penelitian dan PKM berbasis informasi',NULL,2025,'kategori',NULL,NULL,NULL),
            (43,UUID()5,142,'Dosen mendapatkan hibah penelitian eksternal',NULL,2025,'numerik',NULL,NULL,NULL),
            (43,UUID()6,145,'Kegiatan PKM dosen melibatkan mahasiswa',NULL,2025,'numerik',NULL,NULL,NULL),
            (43,UUID()7,145,'Kolaborasi publikasi dengan institusi dalam dan luar negeri',NULL,2025,'numerik',NULL,NULL,NULL),
            (43,UUID()8,148,'Kegiatan PKM dosen dilakukan secara multi disiplin ilmu',NULL,2025,'numerik',NULL,NULL,NULL),
            (43,UUID()9,148,'PKM group di level Prodi/Fakultas/Universitas dengan legalisasi Rektor/Dekan',NULL,2025,'numerik',NULL,NULL,NULL),
            (44,UUID()0,148,'PKM collaboration dengan institusi dalam dan luar negeri',NULL,2025,'numerik',NULL,NULL,NULL),
            (44,UUID()1,151,'Kegiatan PKM dosen sesuai dengan roadmap',NULL,2025,'numerik',NULL,NULL,NULL),
            (44,UUID()2,154,'Reviewer berlisensi/bersertifikat (PkM)',NULL,2025,'numerik',NULL,NULL,'2025-10-22 14:43:26'),
            (44,UUID()3,157,'Setiap pelaksana PKM memiliki roadmap',NULL,2025,'numerik',NULL,NULL,NULL),
            (44,UUID()4,160,'Terdapat laboratorium yang memfasilitasi kegiatan PKM',NULL,2025,'numerik',NULL,NULL,NULL),
            (44,UUID()5,163,'Sistem pengelolaan penelitian dan PKM berbasis informasi',NULL,2025,'kategori',NULL,NULL,NULL),
            (44,UUID()6,166,'Dosen mendapatkan hibah PkM eksternal',NULL,2025,'numerik',NULL,NULL,NULL),
            (44,UUID()7,169,'Dilakukan peninjauan VMTS tiap 5 tahun',NULL,2025,'kategori',NULL,NULL,NULL),
            (44,UUID()8,169,'Akreditasi Internasional',NULL,2025,'numerik',NULL,NULL,NULL),
            (44,UUID()9,169,'Tingkat keterpahamanan VMTS',NULL,2025,'numerik',NULL,NULL,NULL),
            (45,UUID()0,172,'Implementasi standar per tahun',NULL,2025,'kategori',NULL,NULL,NULL),
            (45,UUID()1,175,'Tersedia program retensi untuk dosen dan tenaga kependidikan',NULL,2025,'kategori',NULL,NULL,NULL),
            (45,UUID()2,175,'Terdapat peningkatan Prodi terakreditasi unggul',NULL,2025,'numerik',NULL,NULL,NULL),
            (45,UUID()3,175,'Automasi penjaminan mutu',NULL,2025,'numerik',NULL,NULL,NULL),
            (45,UUID()4,178,'Prodi mengimplementasikan kerjasama untuk proses tridharma',NULL,2025,'numerik',NULL,NULL,NULL),
            (45,UUID()5,178,'Jumlah kerjasama internasional dalam lingkup tridharma',NULL,2025,'numerik',NULL,NULL,NULL),
            (45,UUID()6,181,'Persentase kenaikan jumlah mahasiswa baru',NULL,2025,'kategori',NULL,NULL,NULL),
            (45,UUID()7,181,'Terdapat mekanisme penentuan daya tampung Prodi',NULL,2025,'kategori',NULL,NULL,NULL),
            (45,UUID()8,181,'Tersedia formasi khusus untuk mahasiswa disabilitas',NULL,2025,'numerik',NULL,NULL,NULL),
            (45,UUID()9,181,'Terdapat peningkatan jumlah mahasiswa asing',NULL,2025,'numerik',NULL,NULL,NULL),
            (46,UUID()0,184,'Jumlah mahasiswa yang mendapatkan beasiswa eksternal',NULL,2025,'kategori',NULL,NULL,NULL),
            (46,UUID()1,184,'Tersedianya Unit Layanan Disabilitas (ULD) bagi mahasiswa',NULL,2025,'numerik',NULL,NULL,NULL),
            (46,UUID()2,187,'Terdapat peningkatan jumlah prestasi mahasiswa baik akademik maupun non akademik dalam ruang lingkup nasional dan internasional',NULL,2025,'numerik',NULL,NULL,NULL),
            (46,UUID()3,187,'Jumlah penelitian mahasiswa yang mendapatkan HKI',NULL,2025,'numerik',NULL,NULL,NULL),
            (46,UUID()4,190,'Tersedia pedoman dan SOP pengelolaan keuangan dan dilaporkan setiap tahun anggaran',NULL,2025,'kategori',NULL,NULL,NULL),
            (46,UUID()5,190,'Tersedia mekanisme alokasi pendanaan tridharma',NULL,2025,'kategori',NULL,NULL,NULL),
            (46,UUID()6,190,'Tersedia kebijakan biaya pendidikan untuk mahasiswa berpotensi akademik tapi kurang mampu',NULL,2025,'kategori',NULL,NULL,NULL),
            (46,UUID()7,193,'Peningkatan jumlah sarana dan prasarana yang menjadi Income Generator Unit (IGU)',NULL,2025,'kategori',NULL,NULL,NULL),
            (46,UUID()8,193,'Implementasi green campus pada pengembangan sarana dan prasarana',NULL,2025,'kategori',NULL,NULL,NULL),
            (46,UUID()9,196,'Tersedia sistem informasi terintegrasi untuk akademik, SDM, keuangan, Sarpras, kemahasiswaan dan alumni, kerjasama, riset dan inovasi',NULL,2025,'kategori',NULL,NULL,NULL),
            (47,UUID()0,199,'Tersedia kebijakan biaya pembelajaran peserta MBKM untuk mahasiswa yang memiliki prestasi akademik tapi kurang mampu',NULL,2025,'numerik',NULL,NULL,NULL),
            (47,UUID()1,202,'Tersedia fasilitas kegiatan magang/praktik industri maupun kegiatan proyek desa yang sesuai',NULL,2025,'numerik',NULL,NULL,NULL),
            (47,UUID()2,205,'Percepatan masa studi',NULL,2025,'range','>=',NULL,NULL),
            (47,UUID()3,208,'Terdapat peningkatan respon pengisian dari pengguna lulusan',NULL,2025,'numerik',NULL,NULL,NULL),
            (47,UUID()4,208,'Terdapat peningkatan respon pengisian dari jumlah lulusan',NULL,2025,'numerik',NULL,NULL,NULL),
            (47,UUID()5,211,'Terdapat peningkatan jumlah karya inovasi mahasiswa per tahun',NULL,2025,'numerik',NULL,NULL,NULL),
            (47,UUID()6,118,'Audit keuangan internal dan eksternal tiap triwulan',NULL,2025,'kategori',NULL,NULL,NULL);
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
