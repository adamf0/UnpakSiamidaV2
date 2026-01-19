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

func setupIndikatorRenstraMySQL(t *testing.T) (*gorm.DB, func()) {
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
        CREATE TABLE master_indikator_renstra (
            id int(11) PRIMARY KEY AUTO_INCREMENT NOT NULL,
            uuid varchar(36) DEFAULT NULL,
            id_master_standar int(11) DEFAULT NULL,
            indikator text NOT NULL,
            parent int(11) DEFAULT NULL,
            tahun year(4) NOT NULL,
            tipe_target text DEFAULT NULL,
            operator varchar(5) DEFAULT NULL,
            created_at datetime DEFAULT NULL,
            updated_at datetime DEFAULT NULL
        );

        INSERT INTO master_indikator_renstra (id, uuid, id_master_standar, indikator, parent, tahun, tipe_target, operator, created_at, updated_at) VALUES
        (16, 'b763b5b3-a18e-416c-9d0d-a0c23aa6076c', 1, 'Lulusan memiliki sertifikat kompetensi atau Bahasa asing', NULL, '2024', 'numerik', NULL, '2024-08-25 19:25:05', '2024-09-21 07:18:42'),
        (79, 'c9e5f716-57e4-4349-a95c-83e264db62a1', 1, 'Lulusan bekerja pada ruang lingkup nasional', NULL, '2024', 'numerik', NULL, '2024-09-08 10:45:01', '2024-09-23 08:17:48'),
        (82, '5281d378-4a6f-4b4f-8f95-5abbbc0037e6', 4, 'Jumlah mata kuliah yang mengintegrasikan hasil penelitian dan PKM ke dalam perkuliahan', NULL, '2024', 'numerik', NULL, '2024-09-08 10:45:21', '2024-09-23 08:17:59'),
        (85, '8a4ad016-bbc1-48d6-860d-3b171f4197b9', 4, 'Persentase Prodi Vokasi dan S1 yang mengimplementasikan kurikulum MBKM', NULL, '2024', 'numerik', NULL, '2024-09-09 14:40:00', '2024-10-18 12:35:42'),
        (157, '6f07fdf2-9645-476c-978e-65482731a53e', 7, 'Proses pembelajaran berbasis PBL dan/atau PjBL tiap semester', NULL, '2024', 'numerik', NULL, '2024-09-10 09:31:20', '2024-09-23 08:21:32'),
        (160, '022ca2fe-0725-4c30-9347-ae8bb5c02d32', 7, 'Jumlah exchange lecture di setiap Prodi sejenis', NULL, '2024', 'numerik', NULL, '2024-09-10 09:31:38', '2024-09-23 08:22:38'),
        (163, '16f32b7e-208f-43ad-b695-813f0861434c', 7, 'Perkuliahan bilingual di setiap Prodi', NULL, '2024', 'numerik', NULL, '2024-09-10 09:32:47', '2024-09-23 08:22:49'),
        (166, '86f97046-21ac-486f-9793-d16da7fb864e', 7, 'Peningkatan jumlah kegiatan sebagai penunjang suasana akademik', NULL, '2024', 'numerik', NULL, '2024-09-10 09:33:13', '2024-09-23 08:23:09'),
        (169, '19712e87-2b3c-447c-a8f9-914f11d4166e', 7, 'Persentase pembelajaran dengan moda campuran', NULL, '2024', 'numerik', NULL, '2024-09-10 09:33:36', '2024-09-23 08:23:49'),
        (172, 'd19f83d0-f1c4-4c7b-b9f0-a4a533b7c68e', 7, 'blended/flip learning (Prodi)', 169, '2024', 'numerik', NULL, '2024-09-10 09:33:54', '2024-09-23 08:23:55'),
        (175, 'ad4d1c6c-376b-4170-89d2-78d7d7855244', 7, 'hybrid learning (Fakultas)', 169, '2024', 'numerik', NULL, '2024-09-10 09:34:09', '2024-09-23 08:24:09'),
        (178, '80433cae-6b76-4453-a524-9d9395d300e1', 10, 'IPK rata-rata lulusan/tahun (S1 ≥ 3,00; S2 ≥ 3,25; dan S3 ≥ 3,50)', NULL, '2024', 'range', '>=', '2024-09-10 09:35:10', '2024-09-23 08:24:59'),
        (181, 'fe77df5f-3640-491f-a219-c2a0b26b4b4b', 10, 'Portofolio mahasiswa melalui kurikulum OBE', NULL, '2024', 'numerik', NULL, '2024-09-10 09:35:32', '2024-09-23 08:25:10'),
        (184, 'af1f5720-7e1b-42d7-bd53-59296008a408', 13, 'Dosen mampu Berbahasa Inggris dengan sangat baik', NULL, '2024', 'numerik', NULL, '2024-09-10 09:44:07', '2024-09-23 08:25:20'),
        (187, '9914ae1e-06e2-4d20-b144-928a9e100487', 13, 'Peningkatan jumlah Guru Besar tiap Fakultas', NULL, '2024', 'numerik', NULL, '2024-09-10 09:50:14', '2024-09-23 08:25:33'),
        (190, 'f6f3f919-983c-4c16-a9f3-cce6410f39f3', 13, 'Peningkatan jumlah Tendik registered dan memiliki jabatan fungsional', NULL, '2024', 'numerik', NULL, '2024-09-10 09:52:19', '2024-09-24 22:06:40'),
        (193, '63bd3587-435c-4aca-ab2e-66e7bcd8346f', 25, 'Tersedia Smart Class', NULL, '2024', 'numerik', NULL, '2024-09-10 09:53:17', '2024-09-23 08:28:00'),
        (196, '1eedf2d8-eb70-40ff-9dca-b4df57052945', 28, 'Penyusunan kurikulum berbasi OBE', NULL, '2024', 'kategori', NULL, '2024-09-10 09:54:06', '2024-09-23 08:28:11'),
        (199, '752b0b18-e94b-4e77-b1f2-64bf89923f20', 28, 'Implementasi OBE', NULL, '2024', 'numerik', NULL, '2024-09-10 09:54:18', '2024-09-23 08:28:57'),
        (202, 'aa710cf3-4a39-4053-a7a9-fe1158668667', 28, 'Penyusunan kurikulum MBKM', NULL, '2024', 'kategori', NULL, '2024-09-10 09:54:39', '2024-09-23 08:29:09'),
        (205, 'c2ea650b-38d8-463f-9503-6e1cc3297efe', 28, 'Implementasi MBKM', NULL, '2024', 'kategori', NULL, '2024-09-10 09:59:06', '2024-09-23 08:29:20'),
        (208, '35017e39-e824-4a63-83aa-5db5d4c3e27a', 118, 'Audit keuangan internal', NULL, '2024', 'kategori', NULL, '2024-09-10 09:59:41', '2024-09-23 08:30:59'),
        (211, '69bee102-11f9-49a2-a9ff-07624bf39208', 121, 'Peningkatan jumlah publikasi dosen pada jurnal nasional minimal Sinta 3 dan jurnal internasional bereputasi', NULL, '2024', 'numerik', NULL, '2024-09-10 09:59:57', '2024-09-23 08:31:50'),
        (214, '747c5091-addc-494d-b189-94dbe0b381c4', 121, 'Penelitian kolektif atau kolaboratif dosen melibatkan mahasiswa', NULL, '2024', 'numerik', NULL, '2024-09-10 10:00:24', '2024-09-23 08:31:56'),
        (217, '73e1645e-757b-4ede-b178-99de85e68ad0', 121, 'Kolaborasi publikasi dengan institusi dalam dan luar negeri', NULL, '2024', 'numerik', NULL, '2024-09-10 10:00:44', '2024-09-23 08:32:02'),
        (220, '4d164d41-a87d-429e-801c-d48c16a9d6b9', 121, 'Hilirisasi hasil penelitian', NULL, '2024', 'numerik', NULL, '2024-09-10 10:01:24', '2024-09-23 08:32:08'),
        (223, 'ef778001-d3e0-440b-9493-1f3ae0e1c206', 121, 'Hasil penelitian diadopsi oleh masyarakat atau Teknologi Tepat Guna (TTG)', NULL, '2024', 'numerik', NULL, '2024-09-10 10:01:37', '2024-09-23 08:32:14'),
        (226, 'fa95e2cf-1c25-4ff1-8f07-cd7e6f125bbe', 121, 'Hasil penelitian terstandarisasi atau tersertifikasi', NULL, '2024', 'numerik', NULL, '2024-09-10 10:01:54', '2024-09-23 08:32:27'),
        (229, '7ebd2243-134c-4856-9cab-e25dc8bea9cd', 124, 'Terdapat bukti kesesuaian isi penelitian dengan roadmap penelitian dosen', NULL, '2024', 'kategori', NULL, '2024-09-10 10:02:21', '2024-09-23 08:33:45'),
        (232, 'd8dc41b7-950a-47e3-9f1b-0cd138484ebc', 124, 'Penelitian dosen dilakukan secara multi disiplin ilmu', NULL, '2024', 'numerik', NULL, '2024-09-10 10:02:33', '2024-09-23 08:35:02'),
        (235, 'b6cfe6fb-0965-48aa-a09c-0fa6d64eba03', 124, 'Research group di level Prodi/Fakultas/Universitas dengan legalisasi Rektor/Dekan', NULL, '2024', 'numerik', NULL, '2024-09-10 10:02:47', '2024-09-23 08:35:13'),
        (238, 'de41aa44-21bb-4fd4-98e7-41e3ec470270', 124, 'Research collaboration dengan institusi dalam dan luar negeri', NULL, '2024', 'numerik', NULL, '2024-09-10 10:02:59', '2024-09-23 08:35:23'),
        (241, 'e9c51fa4-49f2-4cbf-b0bd-06bf846056ae', 127, 'Kegiatan penelitian dosen sesuai dengan roadmap', NULL, '2024', 'numerik', NULL, '2024-09-10 10:03:57', '2024-09-23 08:35:28'),
        (244, '21ef1d12-a17e-4566-a51a-5a371b4a1b12', 130, 'Reviewer berlisensi/bersertifikat', NULL, '2024', 'numerik', NULL, '2024-09-10 10:04:15', '2024-09-23 08:35:36'),
        (247, '1bf37fe3-05db-46ca-88db-32174f633ded', 133, 'Setiap peneliti memiliki roadmap penelitian', NULL, '2024', 'numerik', NULL, '2024-09-10 10:04:27', '2024-09-23 08:35:43'),
        (250, '45b6e929-97de-4618-9d1c-995725e0c5fb', 136, 'Standarisasi Sarpras penelitian (KAN/SNI/ISO)', NULL, '2024', 'kategori', NULL, '2024-09-10 10:04:47', '2024-09-24 22:38:59'),
        (253, '6dafa446-4da5-4d6b-a484-f9360655a607', 139, 'Sistem pengelolaan penelitian dan PKM berbasis informasi', NULL, '2024', 'kategori', NULL, '2024-09-10 10:05:10', '2024-09-24 22:39:15'),
        (256, '36428f6e-778f-43ce-8b99-4a056ed3d385', 142, 'Dosen mendapatkan hibah penelitian eksternal', NULL, '2024', 'numerik', NULL, '2024-09-10 10:10:21', '2024-09-23 08:40:44'),
        (259, 'ac54fe6e-b595-4377-b341-f95812c432e1', 145, 'Kegiatan PKM dosen melibatkan mahasiswa', NULL, '2024', 'numerik', NULL, '2024-09-10 10:10:58', '2024-09-23 08:40:50'),
        (262, 'f825f600-fa7d-4058-b113-441be5879028', 145, 'Kolaborasi publikasi dengan institusi dalam dan luar negeri', NULL, '2024', 'numerik', NULL, '2024-09-10 10:12:44', '2024-09-23 08:40:56'),
        (265, '281ac4e4-9da0-4dea-acd6-2c7e0de92aa4', 148, 'Kegiatan PKM dosen dilakukan secara multi disiplin ilmu', NULL, '2024', 'numerik', NULL, '2024-09-10 10:13:56', '2024-09-23 08:41:02'),
        (268, '29137aca-c90b-4815-aaa1-280cf0e8eb88', 148, 'PKM group di level Prodi/Fakultas/Universitas dengan legalisasi Rektor/Dekan', NULL, '2024', 'numerik', NULL, '2024-09-10 10:14:13', '2024-09-23 08:41:18'),
        (271, 'b33241e7-58a7-4362-bfe3-da48204d1c12', 148, 'PKM collaboration dengan institusi dalam dan luar negeri', NULL, '2024', 'numerik', NULL, '2024-09-10 10:14:26', '2024-09-23 08:41:24'),
        (274, 'ab5808fb-74fd-4006-b976-7fe342d97863', 151, 'Kegiatan PKM dosen sesuai dengan roadmap', NULL, '2024', 'numerik', NULL, '2024-09-10 10:15:08', '2024-09-23 08:41:30'),
        (277, '51bf3ff9-84b0-4102-bd63-f4682e3b9613', 154, 'Reviewer berlisensi/bersertifikat', NULL, '2024', 'numerik', NULL, '2024-09-10 10:15:30', '2024-09-23 08:41:36'),
        (280, 'a160ddb4-34f9-47de-b1de-593eb4fccaa2', 157, 'Setiap pelaksana PKM memiliki roadmap', NULL, '2024', 'numerik', NULL, '2024-09-10 10:16:01', '2024-09-23 08:41:41'),
        (283, '256f863b-50fd-4ae2-bc8f-85d19ec56c32', 160, 'Terdapat laboratorium yang memfasilitasi kegiatan PKM', NULL, '2024', 'numerik', NULL, '2024-09-10 10:16:31', '2024-09-23 08:41:46'),
        (286, '9e8171a8-3ee7-49a9-b0bb-b0544d8a5b0e', 163, 'Sistem pengelolaan penelitian dan PKM berbasis informasi', NULL, '2024', 'kategori', NULL, '2024-09-10 10:16:59', '2024-09-25 06:27:45'),
        (289, '493f85f0-a522-4ab5-8935-cc68b0cb8a00', 166, 'Dosen mendapatkan hibah PkM eksternal', NULL, '2024', 'numerik', NULL, '2024-09-10 10:17:37', '2024-09-23 08:41:54'),
        (292, 'd255387b-0fc5-428b-bf58-bd00d6861a15', 169, 'Dilakukan peninjauan VMTS tiap 5 tahun', NULL, '2024', 'kategori', NULL, '2024-09-10 10:18:05', '2024-09-23 08:43:03'),
        (295, 'b4f56da8-9a4a-4509-b5bb-074006a68e12', 169, 'Akreditasi Internasional', NULL, '2024', 'numerik', NULL, '2024-09-10 10:18:28', '2025-02-05 14:42:44'),
        (298, '404f2e8e-a49e-439a-9c6b-4ca069b60d64', 169, 'Tingkat keterpahamanan VMTS', NULL, '2024', 'numerik', NULL, '2024-09-10 10:18:52', '2024-09-23 08:43:26'),
        (301, '58218a5d-ffae-41b0-8183-10d34f7e698f', 172, 'Implementasi standar per tahun', NULL, '2024', 'kategori', NULL, '2024-09-10 10:24:51', '2024-09-25 06:28:22'),
        (304, 'ab8b730d-8128-4845-a5f7-4645650ca139', 175, 'Tersedia program retensi untuk dosen dan tenaga kependidikan', NULL, '2024', 'kategori', NULL, '2024-09-10 10:25:34', '2024-09-25 06:28:41'),
        (307, '14c9a482-5e90-457c-af2b-a94cf8c27fde', 175, 'Terdapat peningkatan Prodi terakreditasi unggul', NULL, '2024', 'numerik', NULL, '2024-09-10 10:26:09', '2025-02-05 14:44:41'),
        (310, '8a3d1f79-3950-4270-82f6-c7277c6c9aff', 175, 'Automasi penjaminan mutu', NULL, '2024', 'numerik', NULL, '2024-09-10 10:26:31', '2025-02-05 14:45:37'),
        (313, '3759d048-6d0b-4c72-a99c-8f70b8705d0b', 178, 'Prodi mengimplementasikan kerjasama untuk proses tridharma', NULL, '2024', 'numerik', NULL, '2024-09-10 10:27:52', '2024-09-23 08:43:44'),
        (316, 'bcabb315-720d-4c09-9d5b-dda3647cb38f', 178, 'Jumlah kerjasama internasional dalam lingkup tridharma', NULL, '2024', 'numerik', NULL, '2024-09-10 10:28:11', '2024-09-23 08:45:43'),
        (319, 'ff0f0506-22aa-4595-92ce-1d4553012b3d', 181, 'Persentase kenaikan jumlah mahasiswa baru', NULL, '2024', 'kategori', NULL, '2024-09-10 10:28:39', '2024-09-25 06:29:01'),
        (322, '380c2fdb-7343-433d-a6cc-5d7de5f1d4e3', 181, 'Terdapat mekanisme penentuan daya tampung Prodi', NULL, '2024', 'kategori', NULL, '2024-09-10 10:29:03', '2024-09-25 06:29:12'),
        (325, '4db2bf44-1d0b-4ee8-ab06-5574e2afc02f', 181, 'Tersedia formasi khusus untuk mahasiswa disabilitas', NULL, '2024', 'numerik', NULL, '2024-09-10 10:29:26', '2025-02-04 11:44:33'),
        (328, '7d110d75-6b10-43c9-9312-3d188c130b64', 181, 'Terdapat peningkatan jumlah mahasiswa asing', NULL, '2024', 'numerik', NULL, '2024-09-10 10:31:59', '2024-09-23 08:46:08'),
        (331, '686be35c-0013-4111-a984-6c90adac7a1e', 184, 'Jumlah mahasiswa yang mendapatkan beasiswa eksternal', NULL, '2024', 'kategori', NULL, '2024-09-10 10:37:33', '2024-09-25 08:23:38'),
        (334, '3ee2d4c2-628f-4c76-8fd6-7451d71f53fd', 184, 'Tersedianya Unit Layanan Disabilitas (ULD) bagi mahasiswa', NULL, '2024', 'numerik', NULL, '2024-09-10 10:38:04', '2025-02-04 13:21:32'),
        (337, '81673252-5f56-411e-b4cf-18860f45bf89', 187, 'Terdapat peningkatan jumlah prestasi mahasiswa baik akademik maupun non akademik dalam ruang lingkup nasional dan internasional', NULL, '2024', 'numerik', NULL, '2024-09-10 10:39:43', '2024-09-23 08:46:23'),
        (340, '22a6bb11-620e-477d-8b0b-18f124b39b2d', 187, 'Jumlah penelitian mahasiswa yang mendapatkan HKI', NULL, '2024', 'numerik', NULL, '2024-09-10 10:49:40', '2024-09-23 08:46:33'),
        (343, '21f7aea8-1382-4acb-89ed-4f31545330ff', 190, 'Tersedia pedoman dan SOP pengelolaan keuangan dan dilaporkan setiap tahun anggaran', NULL, '2024', 'kategori', NULL, '2024-09-10 10:50:03', '2024-09-25 08:23:57'),
        (346, '4c14f1e0-e389-412c-aea5-c6ecdcbd1a8f', 190, 'Tersedia mekanisme alokasi pendanaan tridharma', NULL, '2024', 'kategori', NULL, '2024-09-10 10:50:27', '2024-09-25 08:24:04'),
        (349, '3c6712cd-fe1a-4463-9400-1a1a95a67dad', 190, 'Tersedia kebijakan biaya pendidikan untuk mahasiswa berpotensi akademik tapi kurang mampu', NULL, '2024', 'kategori', NULL, '2024-09-10 10:50:43', '2024-09-25 08:24:12'),
        (352, 'dbe6fa12-c591-42fc-9baf-22d13a6dd0a5', 193, 'Peningkatan jumlah sarana dan prasarana yang menjadi Income Generator Unit (IGU)', NULL, '2024', 'kategori', NULL, '2024-09-10 10:51:39', '2024-09-25 08:24:30'),
        (355, '7a2538c6-a42b-4b1e-8f3c-94732cbdad6e', 193, 'Implementasi green campus pada pengembangan sarana dan prasarana', NULL, '2024', 'kategori', NULL, '2024-09-10 10:51:55', '2024-09-25 08:24:39'),
        (358, 'c2bd90f8-c5f6-4d64-911c-707ae8f6cf96', 196, 'Tersedia sistem informasi terintegrasi untuk akademik, SDM, keuangan, Sarpras, kemahasiswaan dan alumni, kerjasama, riset dan inovasi', NULL, '2024', 'kategori', NULL, '2024-09-10 10:52:10', '2024-09-25 08:24:45'),
        (361, 'a61328a2-8f2b-4abc-9f6c-31da56c105c2', 199, 'Tersedia kebijakan biaya pembelajaran peserta MBKM untuk mahasiswa yang memiliki prestasi akademik tapi kurang mampu', NULL, '2024', 'numerik', NULL, '2024-09-10 10:52:43', '2025-02-04 14:56:13'),
        (364, '68ff6284-32f8-4b0a-af69-7271b0e92772', 202, 'Tersedia fasilitas kegiatan magang/praktik industri maupun kegiatan proyek desa yang sesuai', NULL, '2024', 'numerik', NULL, '2024-09-10 10:53:06', '2024-09-23 08:47:05'),
        (367, 'e1982e94-17c3-421a-85fb-b2e66a06c73e', 205, 'Percepatan masa studi', NULL, '2024', 'range', '>=', '2024-09-10 10:53:28', '2024-10-10 14:12:42'),
        (370, '494666ad-5fce-4f79-af92-eaacb72f1266', 208, 'Terdapat peningkatan respon pengisian dari pengguna lulusan', NULL, '2024', 'numerik', NULL, '2024-09-10 10:53:54', '2024-09-23 08:47:20'),
        (373, 'acac104e-ed75-4edd-8552-5741e72ba861', 208, 'Terdapat peningkatan respon pengisian dari jumlah lulusan', NULL, '2024', 'numerik', NULL, '2024-09-10 10:54:15', '2024-09-23 08:47:26'),
        (376, '442c9738-43f4-4b29-a55a-8737db62d963', 211, 'Terdapat peningkatan jumlah karya inovasi mahasiswa per tahun', NULL, '2024', 'numerik', NULL, '2024-09-10 10:54:35', '2024-09-23 08:47:31'),
        (385, 'd5a63c9f-7734-44fe-a962-a9437e11f8cd', 118, 'Audit keuangan internal dan eksternal tiap triwulan', NULL, '2024', 'kategori', NULL, '2024-09-18 10:02:04', '2024-09-23 08:29:46'),
        (392, '3d4284b4-29fc-4a66-ab3f-b1629aaa710a', 10, '10.	IPK rata-rata lulusan/tahun (S1 ≥ 3,00; S2 ≥ 3,25; dan S3 ≥ 3,50) (Lembaga dan Unit)', NULL, '2024', 'kategori', NULL, '2025-02-03 14:49:48', '2025-02-03 14:49:48'),
        (399, 'b21c7096-a52f-4f2c-9300-c8523d033440', 1, 'Lulusan memiliki sertifikat kompetensi atau Bahasa asing', NULL, '2025', 'numerik', NULL, NULL, NULL),
        (400, 'e86717f4-606e-4782-83fe-4a6e77e1e570', 1, 'Lulusan bekerja pada ruang lingkup nasional', NULL, '2025', 'numerik', NULL, NULL, NULL),
        (401, '87a5e0d1-a7eb-4c5e-8925-ede9632e436b', 4, 'Jumlah mata kuliah yang mengintegrasikan hasil penelitian dan PKM ke dalam perkuliahan', NULL, '2025', 'numerik', NULL, NULL, NULL),
        (402, 'dde95320-b51d-4f7d-a72a-4d36d5f4f24e', 4, 'Persentase Prodi Vokasi dan S1 yang mengimplementasikan kurikulum MBKM', NULL, '2025', 'numerik', NULL, NULL, NULL),
        (403, '3a36b609-e02e-4de2-add8-3198662e7d24', 7, 'Proses pembelajaran berbasis PBL dan/atau PjBL tiap semester', NULL, '2025', 'numerik', NULL, NULL, NULL),
        (404, 'ab76cd68-7c30-4e1a-bf1a-0eec60c95d01', 7, 'Jumlah exchange lecture di setiap Prodi sejenis', NULL, '2025', 'numerik', NULL, NULL, NULL),
        (405, '38a86ab8-4a3f-4f89-8b8a-ec0385aadb80', 7, 'Perkuliahan bilingual di setiap Prodi', NULL, '2025', 'numerik', NULL, NULL, NULL),
        (406, '802ef1e2-38eb-46d9-963b-d210111ced71', 7, 'Peningkatan jumlah kegiatan sebagai penunjang suasana akademik', NULL, '2025', 'numerik', NULL, NULL, NULL),
        (407, 'f5a5736e-edd2-4b7a-916b-a44d27c54ebd', 7, 'Persentase pembelajaran dengan moda campuran', NULL, '2025', 'numerik', NULL, NULL, NULL),
        (408, '574f2339-cb83-4418-b915-2846fdf784a4', 7, 'blended/flip learning (Prodi)', 407, '2025', 'numerik', NULL, NULL, NULL),
        (409, '9763ce5a-ef2d-4317-b3d8-6604075fd383', 7, 'hybrid learning (Fakultas)', 407, '2025', 'numerik', NULL, NULL, NULL),
        (410, '14b43b72-98f8-4ab6-aa8f-78052030eb26', 10, 'IPK rata-rata lulusan/tahun (Diploma ≥ 3,00; S1 ≥ 3,00; S2 ≥ 3,25; dan S3 ≥ 3,50)', NULL, '2025', 'range', '>=', NULL, '2025-10-23 09:19:49'),
        (411, 'adcd895a-e4d6-428c-bba8-00aff7582641', 10, 'Portofolio mahasiswa melalui kurikulum OBE', NULL, '2025', 'numerik', NULL, NULL, NULL),
        (412, 'dc38214a-b036-423a-a1dd-be0bf1037030', 13, 'Dosen mampu Berbahasa Inggris dengan sangat baik', NULL, '2025', 'numerik', NULL, NULL, NULL),
        (413, 'dbf0159e-af1c-4591-924f-3f25e65953bd', 13, 'Peningkatan jumlah Guru Besar tiap Fakultas', NULL, '2025', 'numerik', NULL, NULL, NULL),
        (414, 'd7601d82-1948-4037-9ca0-fe6e2aa40b50', 13, 'Peningkatan jumlah Tendik registered dan memiliki jabatan fungsional', NULL, '2025', 'numerik', NULL, NULL, NULL),
        (415, '5aa7eed7-1b06-44dd-8240-f33ac9674f59', 25, 'Tersedia Smart Class', NULL, '2025', 'numerik', NULL, NULL, NULL),
        (416, '4b6f392e-646c-4dba-8c31-ec84de27666d', 28, 'Penyusunan kurikulum berbasi OBE', NULL, '2025', 'kategori', NULL, NULL, NULL),
        (417, 'e941ec5d-d224-4009-a86f-1f5754219ff1', 28, 'Implementasi OBE', NULL, '2025', 'numerik', NULL, NULL, NULL),
        (418, '013c2eb1-6fc8-4f14-9674-6f8731dc1dfe', 28, 'Penyusunan kurikulum MBKM', NULL, '2025', 'kategori', NULL, NULL, NULL),
        (419, '007599a9-5b45-4871-b863-e817d3620b31', 28, 'Implementasi MBKM', NULL, '2025', 'kategori', NULL, NULL, NULL),
        (420, 'be47be0e-93fb-49c5-b669-d84d1bbc13fc', 121, 'Peningkatan jumlah publikasi dosen pada jurnal nasional minimal Sinta 3 dan jurnal internasional bereputasi', NULL, '2025', 'numerik', NULL, NULL, NULL),
        (421, 'dd3808ad-eadf-4711-8884-7da86de724c7', 121, 'Penelitian kolektif atau kolaboratif dosen melibatkan mahasiswa', NULL, '2025', 'numerik', NULL, NULL, NULL),
        (422, '44ae022c-4061-49c0-8d74-e9221cbbf6b8', 121, 'Kolaborasi publikasi dengan institusi dalam dan luar negeri', NULL, '2025', 'numerik', NULL, NULL, NULL),
        (423, 'eeb439ea-b088-42ce-95a9-b0d86712cfe4', 121, 'Hilirisasi hasil penelitian', NULL, '2025', 'numerik', NULL, NULL, NULL),
        (424, '81a4df9c-fef4-4f22-a42e-4b01e3aeed7c', 121, 'Hasil penelitian diadopsi oleh masyarakat atau Teknologi Tepat Guna (TTG)', NULL, '2025', 'numerik', NULL, NULL, NULL),
        (425, '5a7a8364-cc04-410c-a2de-32ff425c0a59', 121, 'Hasil penelitian terstandarisasi atau tersertifikasi', NULL, '2025', 'numerik', NULL, NULL, NULL),
        (426, 'f36f6d94-3ed4-4d21-975b-be207f2eebba', 124, 'Terdapat bukti kesesuaian isi penelitian dengan roadmap penelitian dosen', NULL, '2025', 'kategori', NULL, NULL, NULL),
        (427, 'c95ad74f-bf9c-4244-ac16-f4521434965e', 124, 'Penelitian dosen dilakukan secara multi disiplin ilmu', NULL, '2025', 'numerik', NULL, NULL, NULL),
        (428, '01011950-1fe9-4f23-961d-e1a6b3ae947e', 124, 'Research group di level Prodi/Fakultas/Universitas dengan legalisasi Rektor/Dekan', NULL, '2025', 'numerik', NULL, NULL, NULL),
        (429, '25671ad3-1436-418e-9b9f-7f2afa981340', 124, 'Research collaboration dengan institusi dalam dan luar negeri', NULL, '2025', 'numerik', NULL, NULL, NULL),
        (430, '6b28e4c0-f665-44d6-8b32-9c9d765f8204', 127, 'Kegiatan penelitian dosen sesuai dengan roadmap', NULL, '2025', 'numerik', NULL, NULL, NULL),
        (431, '499dd4ed-62e8-4920-a180-fbf916ca713a', 130, 'Reviewer berlisensi/bersertifikat (Penelitian)', NULL, '2025', 'numerik', NULL, NULL, '2025-10-22 14:43:03'),
        (432, 'd41babfd-15eb-4d8c-9c65-dc225df161a6', 133, 'Setiap peneliti memiliki roadmap penelitian', NULL, '2025', 'numerik', NULL, NULL, NULL),
        (433, '4642be1c-9107-4593-8762-d83350f69243', 136, 'Standarisasi Sarpras penelitian (KAN/SNI/ISO)', NULL, '2025', 'kategori', NULL, NULL, NULL),
        (434, 'd781e6c1-6944-4bdf-aa5a-5bdca3b962db', 139, 'Sistem pengelolaan penelitian dan PKM berbasis informasi', NULL, '2025', 'kategori', NULL, NULL, NULL),
        (435, '65895a3f-323a-4612-9722-e6ca2f31bdb1', 142, 'Dosen mendapatkan hibah penelitian eksternal', NULL, '2025', 'numerik', NULL, NULL, NULL),
        (436, '4bd32763-9a51-453f-b53f-089d24bb39b5', 145, 'Kegiatan PKM dosen melibatkan mahasiswa', NULL, '2025', 'numerik', NULL, NULL, NULL),
        (437, '9bcafc98-c29e-4819-9928-0ed284249105', 145, 'Kolaborasi publikasi dengan institusi dalam dan luar negeri', NULL, '2025', 'numerik', NULL, NULL, NULL),
        (438, 'f9d43307-9171-445c-8282-a65be98f49d2', 148, 'Kegiatan PKM dosen dilakukan secara multi disiplin ilmu', NULL, '2025', 'numerik', NULL, NULL, NULL),
        (439, 'ce56f707-39c4-4eb9-bd86-3be674268121', 148, 'PKM group di level Prodi/Fakultas/Universitas dengan legalisasi Rektor/Dekan', NULL, '2025', 'numerik', NULL, NULL, NULL),
        (440, '9382e63b-41c5-4674-b414-d8d852e033cb', 148, 'PKM collaboration dengan institusi dalam dan luar negeri', NULL, '2025', 'numerik', NULL, NULL, NULL),
        (441, '811df41a-4c42-4639-966a-68d0a500f041', 151, 'Kegiatan PKM dosen sesuai dengan roadmap', NULL, '2025', 'numerik', NULL, NULL, NULL),
        (442, 'b9704759-f963-4556-9d6f-30709c0e20b9', 154, 'Reviewer berlisensi/bersertifikat (PkM)', NULL, '2025', 'numerik', NULL, NULL, '2025-10-22 14:43:26'),
        (443, 'ea2e6fb1-f626-4108-a677-cb58798e1a75', 157, 'Setiap pelaksana PKM memiliki roadmap', NULL, '2025', 'numerik', NULL, NULL, NULL),
        (444, 'da6858b2-398f-4557-b8d4-aac7b47cd649', 160, 'Terdapat laboratorium yang memfasilitasi kegiatan PKM', NULL, '2025', 'numerik', NULL, NULL, NULL),
        (445, '6d8f6ed7-842a-4a2d-8801-b75893a33932', 163, 'Sistem pengelolaan penelitian dan PKM berbasis informasi', NULL, '2025', 'kategori', NULL, NULL, NULL),
        (446, '3e654b43-78ac-448b-bd0e-15b6e67e6d51', 166, 'Dosen mendapatkan hibah PkM eksternal', NULL, '2025', 'numerik', NULL, NULL, NULL),
        (447, 'c55861f5-1828-41e6-97e4-257a5befce90', 169, 'Dilakukan peninjauan VMTS tiap 5 tahun', NULL, '2025', 'kategori', NULL, NULL, NULL),
        (448, '0ffaf09a-c2f9-4a93-888e-595eb3c8c6a9', 169, 'Akreditasi Internasional', NULL, '2025', 'numerik', NULL, NULL, NULL),
        (449, '2caafa23-244d-41a8-875d-d2f4b6e26d2d', 169, 'Tingkat keterpahamanan VMTS', NULL, '2025', 'numerik', NULL, NULL, NULL),
        (450, '83a80296-986c-4f13-b384-63e8d8b7c3d3', 172, 'Implementasi standar per tahun', NULL, '2025', 'kategori', NULL, NULL, NULL),
        (451, '0dfceff4-1ae4-4ca6-92d4-a884a058ae59', 175, 'Tersedia program retensi untuk dosen dan tenaga kependidikan', NULL, '2025', 'kategori', NULL, NULL, NULL),
        (452, 'c8dc4b44-39ab-4992-b652-0a5dffc9f227', 175, 'Terdapat peningkatan Prodi terakreditasi unggul', NULL, '2025', 'numerik', NULL, NULL, NULL),
        (453, '96015613-191c-4619-8b20-2de392b66410', 175, 'Automasi penjaminan mutu', NULL, '2025', 'numerik', NULL, NULL, NULL),
        (454, 'a43038c3-9690-4ac6-b7c4-ea8245bd9739', 178, 'Prodi mengimplementasikan kerjasama untuk proses tridharma', NULL, '2025', 'numerik', NULL, NULL, NULL),
        (455, 'ecc4c700-1cf9-4a00-9d41-825c54061f3f', 178, 'Jumlah kerjasama internasional dalam lingkup tridharma', NULL, '2025', 'numerik', NULL, NULL, NULL),
        (456, '618dc755-97f9-4200-bffe-9cea23f21c97', 181, 'Persentase kenaikan jumlah mahasiswa baru', NULL, '2025', 'kategori', NULL, NULL, NULL),
        (457, 'de1f068f-35b8-4480-bcfd-678cc34c33a5', 181, 'Terdapat mekanisme penentuan daya tampung Prodi', NULL, '2025', 'kategori', NULL, NULL, NULL),
        (458, 'b1c8b1ad-ca02-4d9c-b78c-eaacd3d4361c', 181, 'Tersedia formasi khusus untuk mahasiswa disabilitas', NULL, '2025', 'numerik', NULL, NULL, NULL),
        (459, 'cbce60ff-a411-45d3-ae01-1a7b695b9377', 181, 'Terdapat peningkatan jumlah mahasiswa asing', NULL, '2025', 'numerik', NULL, NULL, NULL),
        (460, '069dda69-4ea7-4076-a810-75506d0c5a4f', 184, 'Jumlah mahasiswa yang mendapatkan beasiswa eksternal', NULL, '2025', 'kategori', NULL, NULL, NULL),
        (461, '907109d0-2a09-45be-b052-e3232cd24ae5', 184, 'Tersedianya Unit Layanan Disabilitas (ULD) bagi mahasiswa', NULL, '2025', 'numerik', NULL, NULL, NULL),
        (462, '123125e6-5979-4214-be05-f8914af65671', 187, 'Terdapat peningkatan jumlah prestasi mahasiswa baik akademik maupun non akademik dalam ruang lingkup nasional dan internasional', NULL, '2025', 'numerik', NULL, NULL, NULL),
        (463, '283bb0fe-a1f5-48ac-907f-5ae0e69e0656', 187, 'Jumlah penelitian mahasiswa yang mendapatkan HKI', NULL, '2025', 'numerik', NULL, NULL, NULL),
        (464, '2176dc04-81a2-46e7-ba31-288146e2d25b', 190, 'Tersedia pedoman dan SOP pengelolaan keuangan dan dilaporkan setiap tahun anggaran', NULL, '2025', 'kategori', NULL, NULL, NULL),
        (465, 'b878083e-cae0-428f-8015-95bd331cec99', 190, 'Tersedia mekanisme alokasi pendanaan tridharma', NULL, '2025', 'kategori', NULL, NULL, NULL),
        (466, '25aae6a4-3edf-4b52-8261-40df6c82e6c6', 190, 'Tersedia kebijakan biaya pendidikan untuk mahasiswa berpotensi akademik tapi kurang mampu', NULL, '2025', 'kategori', NULL, NULL, NULL),
        (467, '96aeedcb-ecfc-4581-bdf1-90b32bbb906f', 193, 'Peningkatan jumlah sarana dan prasarana yang menjadi Income Generator Unit (IGU)', NULL, '2025', 'kategori', NULL, NULL, NULL),
        (468, 'c0dbd9a3-c173-498c-a1e7-bbca0253e6f6', 193, 'Implementasi green campus pada pengembangan sarana dan prasarana', NULL, '2025', 'kategori', NULL, NULL, NULL),
        (469, '1b4e2e3c-4ac2-45de-b0b8-bba10d6a1209', 196, 'Tersedia sistem informasi terintegrasi untuk akademik, SDM, keuangan, Sarpras, kemahasiswaan dan alumni, kerjasama, riset dan inovasi', NULL, '2025', 'kategori', NULL, NULL, NULL),
        (470, 'd8ed69df-62fa-4f5c-8007-fe1d62f4bdc1', 199, 'Tersedia kebijakan biaya pembelajaran peserta MBKM untuk mahasiswa yang memiliki prestasi akademik tapi kurang mampu', NULL, '2025', 'numerik', NULL, NULL, NULL),
        (471, '7a950eac-b93f-4001-b7f8-b5b6734e2e1f', 202, 'Tersedia fasilitas kegiatan magang/praktik industri maupun kegiatan proyek desa yang sesuai', NULL, '2025', 'numerik', NULL, NULL, NULL),
        (472, 'abeabef6-0340-4cc8-b119-91b0d98e18c8', 205, 'Percepatan masa studi', NULL, '2025', 'range', '>=', NULL, NULL),
        (473, '85298e18-68d8-4450-952c-dd58a364e133', 208, 'Terdapat peningkatan respon pengisian dari pengguna lulusan', NULL, '2025', 'numerik', NULL, NULL, NULL),
        (474, 'da43e364-345c-4957-a84e-9f7edd066938', 208, 'Terdapat peningkatan respon pengisian dari jumlah lulusan', NULL, '2025', 'numerik', NULL, NULL, NULL),
        (475, '40a20311-bc97-4c67-a10a-c71685bb3911', 211, 'Terdapat peningkatan jumlah karya inovasi mahasiswa per tahun', NULL, '2025', 'numerik', NULL, NULL, NULL),
        (476, '8dafad3e-b182-46a7-a678-8d23bd553f92', 118, 'Audit keuangan internal dan eksternal tiap triwulan', NULL, '2025', 'kategori', NULL, NULL, NULL),
        (479, '9c0bcb0d-2a6d-42c3-90aa-132f69d6ccd3', 205, 'uji coba', NULL, '2026', 'range', '>=', '2024-09-10 10:53:28', '2024-10-10 14:12:42'),
        (480, 'd0754056-4d55-4091-a13d-0fae624a7616', 211, 'tester', NULL, '2026', 'kategori', '>=', NULL, NULL);

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
