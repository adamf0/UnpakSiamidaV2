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

func setupAccountMySQL(t *testing.T) (*gorm.DB, func()) {
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
