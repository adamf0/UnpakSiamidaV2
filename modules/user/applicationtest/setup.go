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

func setupUserMySQL(t *testing.T) (*gorm.DB, func()) {
    ctx := context.Background()

    req := testcontainers.ContainerRequest{
        Image:        "mysql:8.0",
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
            id bigint(20) unsigned NOT NULL AUTO_INCREMENT,
            uuid varchar(36) DEFAULT NULL,
            nidn_username varchar(255) NOT NULL,
            password varchar(255) NOT NULL,
            level enum('','admin','auditee','auditor1','auditor2','fakultas') NOT NULL,
            name varchar(255) NOT NULL,
            email varchar(255) DEFAULT NULL,
            fakultas_unit int(11) DEFAULT NULL,
            foto text DEFAULT NULL,
            email_verified_at timestamp NULL DEFAULT NULL,
            remember_token varchar(100) DEFAULT NULL,
            reset_password int(11) DEFAULT 0,
            created_at timestamp NULL DEFAULT NULL,
            updated_at timestamp NULL DEFAULT NULL,
            PRIMARY KEY (id)
        );
        INSERT INTO users VALUES 
            (2, '186f2427-8bdd-42d9-a757-65808f364eeb','admin','$2y$10$OKzK1M/XuKAuUCtfQ6FvBeZpyxjQYkMQnj8QcySVzY/cQ7Xk8s1hW','admin','Admin',NULL,NULL,'Tangkapan Layar 2025-10-17 pukul 09.48.28.png',NULL,NULL,0,'2023-05-18 18:45:52','2025-10-22 06:02:49'),
            (40, UUID(),'Didik NotoSudjono','$2y$10$uPFePuHtBrbp0FGwz82O/u.imhEDzr6C6ndQPHzSOfYyWo/e9d37u','auditor1','Prof. Dr. Ir. rer. pol. Didik NotoSudjono, M.Sc','didiknotosudjono@unpak.ac.id',NULL,NULL,NULL,NULL,NULL,'2023-10-02 02:45:45','2024-11-30 08:45:29'),
            (43, UUID(),'Istiqlaliah','$2y$10$Ld0Ws48jrqWfONzmsu0bweSDl0WqCwmhr7SVwkyeUfVlVVA4bqdcO','auditor1','Dr. Istiqlaliah Nurul  Hidayati, M.Pd','istiqlaliah@unpak.ac.id',NULL,'aku2.jpg',NULL,NULL,NULL,'2023-10-02 02:48:45','2024-12-30 07:04:23'),
            (46, UUID(),'Eri Sarimanah','$2y$10$jsegv/6gGVmrex3hxHnwo.vIav6hfnwYshQUJwk.OkJKbg1jkjPje','auditor1','Prof. Dr. Eri Sarimanah, M.Pd','erisarimanah@unpak.ac.id',NULL,NULL,NULL,NULL,NULL,'2023-10-02 02:49:22','2024-11-30 08:45:40'),
            (49, UUID(),'Yuary Farradia','$2y$10$ingyBMcI42jlDlKunW/HO.5Ni4y1.nlonFjsMyL4dGjXTbnDnFYn6','auditor1','Dr. Ir. Yuari Farradia, M.Sc','yuary.farradia@unpak.ac.id',NULL,NULL,NULL,NULL,NULL,'2023-10-02 02:52:36','2025-01-15 06:11:09'),
            (52, UUID(),'Andi Chaerunnas','$2y$10$QhSKATxTznNh9Dl2IBGlN.DyG7H2xvsIHJE8beWChxfrHP7WJQeoS','auditor1','Dr. Andi Chairunnas, M.Pd,.M.Kom','andichairunnas@unpak.ac.id',NULL,NULL,NULL,NULL,NULL,'2023-10-02 02:53:15','2024-11-30 08:46:20'),
            (55, UUID(),'Griet Helena','$2y$10$jyCXqbg1.JPqFF4EBnXOcOYvFa8WuRaiFUfQEFOdXonQgsSNYrPOa','auditor1','Dr. Griet Helena Laihad, M.Pd','grihela@unpak.ac.id',NULL,'FOTO GHL.jpeg',NULL,NULL,NULL,'2023-10-02 02:54:26','2025-02-23 13:27:47'),
            (58, UUID(),'Agung Fajar','$2y$10$lICelZarONq09QzEZVRQjeWMOvlm/7UMMcCjPpDQA4Nap8d4k6APC','auditor1','Dr. Agung Fajar Ilmiyono, SE.,M.Ak.,AWP.,CFA.,CAP','agung.fajar@unpak.ac.id',NULL,'IMG_Foto Agung.jpg',NULL,NULL,NULL,'2023-10-02 02:55:01','2024-12-23 03:41:52'),
            (61, UUID(),'Edi Rohaedi','$2y$10$uC7MyfUQijl0Yv.xTlILBeSOWQBd1IFTzWr3gxRw/XG8V3U1Qa2Ja','auditor1','Edi Rohaedi, SH.,MH','edi.rohaedi@unpak.ac.id',NULL,NULL,NULL,NULL,1,'2023-10-02 02:55:25','2024-11-29 12:53:03'),
            (64, UUID(),'Prihastuti Harsani','$2y$10$Ttz.xCyEQTV.ggSVcmIIueGRacOqfzTUdTyfdbhGnft2kymwcf/7S','auditor1','Dr. Prihastuti Harsani, M.Si','prihastuti.harsani@unpak.ac.id',NULL,'Foto.jpg',NULL,NULL,NULL,'2023-10-02 02:56:02','2024-12-16 09:38:10'),
            (67, UUID(),'Herman','$2y$10$EZshWaUkDra0ZbhWvNSlIuu5w2IYxPku3JM9VBncU9Ihh6kqOxYWG','auditor1','Dr. Herman, M.M','herman_fhz@unpak.ac.id',NULL,NULL,NULL,NULL,NULL,'2023-10-02 02:56:38','2025-07-26 02:19:09'),
            (70, UUID(),'Cantika Zaddana','$2y$10$61lLrJTmJjov7ys7n8N8V.s/CIk8TlTHpjUxnOKqmtgOp5sFHb/KS','auditor1','Cantika Zaddana, S.Gz, M.Si','cantika.zaddana@unpak.ac.id',NULL,'the newest.jpg',NULL,NULL,NULL,'2023-10-02 02:57:04','2025-08-11 08:03:05'),
            (73, UUID(),'Indri Yani','$2y$10$fsxG71.B7HYSPwnw4U7wVOu4YTZtKHUDJrOog1x59/0TGYmpWHyue','auditor1','Dr. Indri Yani, M.Pd','indri@unpak.ac.id',NULL,NULL,NULL,NULL,1,'2023-10-02 02:57:34','2024-11-29 12:20:50'),
            (76, UUID(),'Indarini','$2y$10$dhJ8GjS/iFeTo2wi3HDiy.PPzqLY4.ljrEwE4TieK9bLOIM/3UPRy','auditor1','Prof. Dr. Indarini Dwi Pursitasari, M.Si','indarini.dp@unpak.ac.id',NULL,NULL,NULL,NULL,1,'2023-10-02 02:59:10','2024-12-02 08:33:30'),
            (79, UUID(),'Dolly Priatna','$2y$10$MLDSXqLMI/3Zc7f8kMxN6.l3o4pzxNMPcPeqSouwUX8EKj5r4JN7O','auditor1','Dr. Dolly Priatna, M.Si','dollypriatna@unpak.ac.id',NULL,NULL,NULL,NULL,NULL,'2023-10-02 02:59:27','2025-11-18 03:31:54'),
            (82, UUID(),'Anna Permanasari','$2y$10$mCGcvCnOJt9W1hQOOt7qQ.TeqpF8ZxSJe4eqcT0fnigudZeocXjyq','auditor1','Prof. Dr. Anna Permanasari, M.Si','anna.permanasari@unpak.ac.id',NULL,NULL,NULL,NULL,1,'2023-10-02 02:59:44','2024-11-30 02:44:03'),
            (85, UUID(),'Indarti','$2y$10$jMZ2l6IgzjFgcmlpmOmGJOJv6xyEb4V9ujDQXMoeERNWesKM7h8Mq','auditor1','Dr.  Ir. Indarti Komala Dewi. M.Si','indarti@unpak.ac.id',NULL,NULL,NULL,NULL,NULL,'2023-10-02 03:00:12','2024-12-20 03:29:29'),
            (88, UUID(),'Herdiyana','$2y$10$Ojsaoli8D365aPwa8YDf8uGpDTMtu6gfmd/Sx6NY1gRExmVn9S0Ne','auditor1','Dr. Herdiyana, MM','herdiyana@unpak.ac.id',NULL,NULL,NULL,NULL,NULL,'2023-10-02 03:01:29','2025-01-02 03:36:35'),
            (91, UUID(),'Ade Heri','$2y$10$B0ruVIuyfYNLd5Rxuben6OFaUPeTydzUPv3wV985J5TlohB6XAkjO','auditor1','Dr. Ade Heri Mulyati, M.Si','adeheri.mulyati@unpak.ac.id',NULL,NULL,NULL,NULL,1,'2023-10-02 03:01:51','2024-12-10 01:56:13'),
            (94, UUID(),'Iwan Darmawan','$2y$10$FPwzLhPlJk57Nvzmz62MK.rND1N2m5tGKoFkWwAj1lKncgcNdHjOK','auditor1','Dr. Iwan Darmawan, MH','iwan.darmawan@unpak.ac.id',NULL,NULL,NULL,NULL,NULL,'2023-10-02 03:02:21','2024-11-30 08:54:15'),
            (97, UUID(),'Irvan Permana','$2y$10$NBKVlwCp9Rb9ZRbXvcNXGeDabp1wB70aRCWGvTiMGbhnDfgzPK.ze','auditor1','Dr. Irvan Permana, M.Pd','irvanpermana@unpak.ac.id',NULL,NULL,NULL,NULL,NULL,'2023-10-02 03:02:36','2024-11-30 08:54:32'),
            (103, UUID(),'Helen Susanti','$2y$10$wM6644lf0gqLwBSYnk1vFOk89QCD8uu8vtiuG0LOhywEiY/3MbIua','auditor1','Helen Susanti, M.Si','helen.susanti@unpak.ac.id',NULL,NULL,NULL,NULL,NULL,'2023-10-02 03:03:18','2024-11-30 08:54:50'),
            (106, UUID(),'Haqni Wijayanti','$2y$10$FZg.TYn0tbAKcsL.9Ml/tufTzlep7VSSkWyEqkFWPRx/P.vsiapbi','auditor1','Haqni Wijayanti, M.Si','hagnijantix@unpak.ac.id',NULL,NULL,NULL,NULL,1,'2023-10-02 03:04:00','2024-12-08 13:43:16'),
            (112, UUID(),'Rita Retnowati','$2y$10$ncRLgGe9LLPM9P.e3NqP3.evq6mxWjWdepMjLerwS1lznwyPLX1Rq','auditor1','Dr. Rita Retnowati, MS','ritaretnowati@unpak.ac.id',NULL,NULL,NULL,NULL,NULL,'2023-10-10 02:17:53','2024-11-30 08:55:34'),
            (115, UUID(),'Hari Muharam','$2y$10$gr5Sua95GJRIlvtGecwamexnPlKrRBGK2bNYNSBkIZbnaziH5iYQi','auditor1','Dr. Hari Muharam, S.E.,M.M','hari.muharam@unpak.ac.id',NULL,NULL,NULL,NULL,NULL,'2023-10-10 02:20:37','2025-10-02 06:01:03'),
            (118, UUID(),'Ellyn Octavianty','$2y$10$Ft37tJ9cQiOCJKDASllrgeGuYnJCXeo8iy0WawYxeYv/nxnfsyhoK','auditor1','Ellyn Octavianty, SE,. MM','ellynoctavianty@unpak.ac.id',NULL,NULL,NULL,NULL,NULL,'2023-10-10 02:29:19','2024-11-30 08:56:03'),
            (121, UUID(),'Heny Purwanti','$2y$10$4poGlC3wQ1h63/Y7bDFFZuivGASsoFEaT8ALhztETZqZnj4mEWtDO','auditor1','Heny Purwanti, M.T','henypurwanti@unpak.ac.id',NULL,NULL,NULL,NULL,NULL,'2023-10-10 02:29:49','2024-11-30 08:56:18'),
            (124, UUID(),'Patar Simamora','$2y$10$VHqgr8ldkVIo.oZf0JMcAOB2NurxmJndMLM4VXgQLf83yJri7PHyK','auditor1','Patar Simamora, S.E.,M.M','patar.simamora@unpak.ac.id',NULL,NULL,NULL,NULL,NULL,'2023-10-10 02:30:14','2024-11-30 08:56:35'),
            (127, UUID(),'Solihin','$2y$10$NcQ7MK7TuQNFtjICSv/TF.Ieur1xoNK5/aON8Y.fDuoi5X5So8fuG','','Solihin, M.T','solihin@unpak.ac.id',NULL,NULL,NULL,NULL,1,'2023-10-10 02:30:42','2025-10-21 08:44:48'),
            (151, UUID(),'Istiqlaliah','$2y$10$eGhZQKpMT0zDrJE0pkVwMO8T0NacgAN4A3gC8Q9ZZj4fPFVrJlcnK','auditor2','Dr. Istiqlaliah Nurul Hidayati, M.Pd','istiqlaliah@unpak.ac.id',NULL,NULL,NULL,NULL,NULL,'2023-10-16 02:13:49','2025-01-17 04:06:43'),
            (154, UUID(),'Eri Sarimanah','$2y$10$Gy37kYsYRjPjU9EmtS9z8OORQmSByxiHY8k.dK6XXPjwxg4SXWQJ.','auditor2','Prof. Dr. Eri Sarimanah, M.Pd','erisarimanah@unpak.ac.id',NULL,NULL,NULL,NULL,NULL,'2023-10-16 02:15:00','2024-12-02 01:33:21'),
            (157, UUID(),'Yuary Farradia','$2y$10$p.DsD1UpH7MSYK8RUeIDr.0ZE4rBzDpEL.dHp4reU3eWIOtV/GdHm','auditor2','Dr. Ir. Yuari Farradia, M.Sc','yuary.farradia@unpak.ac.id',NULL,NULL,NULL,NULL,NULL,'2023-10-16 02:15:18','2024-11-30 08:46:02'),
            (160, UUID(),'Andi Chaerunnas','$2y$10$NnPLz8wDZqsl6eeFrdadeuokvdoaU74uAIRGHqDfS0aXhrJK7vQ5O','auditor2','Dr. Andi Chairunnas, M.Pd,.M.Kom','andichairunnas@unpak.ac.id',NULL,NULL,NULL,NULL,NULL,'2023-10-16 02:15:33','2024-11-30 08:46:27'),
            (163, UUID(),'Griet Helena','$2y$10$Lpm3B5YTXfPxzEy7fSywKu8J78baQKHMHJp8Xn10UhyFy2sO9OUzi','auditor2','Dr. Griet Helena Laihad, M.Pd','grihela@unpak.ac.id',NULL,'FOTO GHL.jpeg',NULL,NULL,1,'2023-10-16 02:15:47','2024-12-15 21:55:11'),
            (166, UUID(),'Agung Fajar','$2y$10$2UczRfpQIwbSWzAKQDvSBeSjk1W2nCxvTqbgMFCOiR2NBU4HZc6Wm','auditor2','Dr. Agung Fajar Ilmiyono, SE., M.Ak.,AWP.,C.F.A.,CAP','agung.fajar@unpak.ac.id',NULL,'IMG_Foto Agung.jpg',NULL,NULL,NULL,'2023-10-16 02:16:07','2024-12-20 02:49:32'),
            (169, UUID(),'Edy Rohaedi','$2y$10$NWQ8P.vWVdlPIblw9gEcZeN4P/DkUP8hJR5/1HGLginlIkSH1NzhW','auditor2','Edy Rohaedi, SH.,MH','edi.rohaedi@unpak.ac.id',NULL,NULL,NULL,NULL,NULL,'2023-10-16 02:16:20','2024-11-30 08:51:52'),
            (172, UUID(),'Prihastuti Harsani','$2y$10$9HmAtDW7aNCQcWIEzfbp1Of6Iv6vkOdRycUupQIr/sEZUpZwokkN.','auditor2','Dr. Prihastuti Harsani, M.Si','prihastuti.harsani@unpak.ac.id',NULL,NULL,NULL,NULL,NULL,'2023-10-16 02:16:47','2024-11-30 08:52:13'),
            (175, UUID(),'Herman','$2y$10$I5yrPAYhJ.TuMKC074XRfuh7JzS0ZVA1oZW/ababPS0crsWJSyZMa','auditor2','Dr. Herman, M.M','herman_fhz@unpak.ac.id',NULL,'Herman_Foto.jpeg',NULL,NULL,NULL,'2023-10-16 02:16:58','2025-07-28 04:28:46'),
            (178, UUID(),'Cantika Zaddana','$2y$10$YdU4yWZks3STeHfZ1ATLEeTl1qDVuM9oCbnOg8zVpmNOhEiHlvUb6','auditor2','Cantika Zaddana, S.Gz, M.Si','cantika.zaddana@unpak.ac.id',NULL,'the newest.jpg',NULL,NULL,NULL,'2023-10-16 02:17:09','2025-07-28 04:19:06'),
            (181, UUID(),'Indri Yani','$2y$10$s/H7aWFfQtt4zKigwfOiCum32qMZNif0OtgYMQOxDBAyEKINny/JS','auditor2','Dr. Indri Yani, M.Pd','indri@unpak.ac.id',NULL,NULL,NULL,NULL,1,'2023-10-16 02:17:18','2024-12-01 00:43:10'),
            (184, UUID(),'Indarini','$2y$10$ahihMYC9oWJKy06gWNuIFuO0/0GjgbOAlSi8jHIQdoOM3kIJOS.4y','auditor2','Prof. Dr. Indarini Dwi Pursitasari, M.Si','indarini.dp@unpak.ac.id',NULL,NULL,NULL,NULL,NULL,'2023-10-16 02:17:29','2024-11-30 08:52:34'),
            (187, UUID(),'Dolly Priatna','$2y$10$M1ss6YWCHT6liuvVH5mar.KnjGvXHDpInv11OdZ9/lVqiqXBNOc.q','auditor2','Dr. Dolly Priatna, M.Si','dollypriatna@unpak.ac.id',NULL,'Dolly Priatna Foto Cop21 paris_2015.JPG',NULL,NULL,NULL,'2023-10-16 02:17:39','2025-04-11 04:19:59'),
            (190, UUID(),'Anna Permanasari','$2y$10$IDydlEycrG3m5V/WmJon4OoujMG3fi5DYrVIyYUPz0vczZDdRWJZ6','auditor2','Prof. Dr. Anna Permanasari, M.Si','anna.permanasari@unpak.ac.id',NULL,NULL,NULL,NULL,NULL,'2023-10-16 02:17:53','2024-11-30 08:53:13'),
            (193, UUID(),'Indarti','$2y$10$xuM1b3lGEVsGp02vg7mnMOlyIJe5TxOAMptBuek4D9K83uqRiA66q','auditor2','Dr. Indarti Komala Dewi. M.Si','indarti@unpak.ac.id',NULL,NULL,NULL,NULL,NULL,'2023-10-16 02:18:10','2024-11-30 08:53:30'),
            (196, UUID(),'Herdiyana','$2y$10$lnejzMvghFjgSD7aPe6pB.kthh2v1OISmCrD4yrxThmXdBykaM7oy','auditor2','Dr. Herdiyana, MM','herdiyana@unpak.ac.id',NULL,NULL,NULL,NULL,NULL,'2023-10-16 02:18:20','2024-11-30 08:53:47'),
            (199, UUID(),'Ade Heri','$2y$10$KUsCg4hycazNPNzbOr.E5eaCHbVIV3LFfB3/HU23y/U58ZXUjJORC','auditor2','Dr. Ade Heri Mulyati, M.Si','adeheri.mulyati@unpak.ac.id',NULL,'FOTO JAS ADE HERI.jpeg',NULL,NULL,1,'2023-10-16 02:18:30','2024-12-16 11:32:40'),
            (202, UUID(),'Iwan Darmawan','$2y$10$kVEi2awkz21XXS5qm6Ce9OFxGaZjWNYquvpVK52I8XHvmiYCUNELW','auditor2','Dr. Iwan Darmawan, MH','iwan.darmawan@unpak.ac.id',NULL,NULL,NULL,NULL,NULL,'2023-10-16 02:18:42','2024-11-30 08:54:22'),
            (205, UUID(),'Irvan Permana','$2y$10$cPPacELnWOcUh7M/a7N66OsaQhNdzeiVJyK9S5RfI711gog8Z9THm','auditor2','Dr. Irvan Permana, M.Pd','irvanpermana@unpak.ac.id',NULL,NULL,NULL,NULL,NULL,'2023-10-16 02:18:53','2024-11-30 08:54:39'),
            (208, UUID(),'Helen Susanti','$2y$10$pNTeg.CubrZz7v8rDf.queIQob5NLKWu203jcZJHXqMKWOpQUR2NK','auditor2','Helen Susanti, M.Si','helen.susanti@unpak.ac.id',NULL,NULL,NULL,NULL,NULL,'2023-10-16 02:19:24','2024-11-30 08:54:56'),
            (211, UUID(),'Haqni Wijayanti','$2y$10$3ZS0kz08B2DXsKcQVX9.h.JtNgmrwBnPGbMwJylKPsod.JMHAuDe.','auditor2','Haqni Wijayanti, M.Si','hagnijantix@unpak.ac.id',NULL,NULL,NULL,NULL,NULL,'2023-10-16 02:19:35','2025-03-20 04:14:47'),
            (214, UUID(),'Rita Retnowati','$2y$10$ylZt7LWLGuVqVnsgLueEQ..CjTJVB1hGasxdN7Ae4lpC/UlYiEGai','auditor2','Rita Retnowati','ritaretnowati@unpak.ac.id',NULL,NULL,NULL,NULL,NULL,'2023-10-16 02:19:46','2024-11-30 08:55:39'),
            (217, UUID(),'Hari Muharam','$2y$10$262w7uUNxbI3BzehL6hXV.50/k/Z3b2fCebHrgI7pNTLqmuREVtp6','auditor2','Dr. Hari Muharam, SE.,MM.,CIHCM','hari.muharam@unpak.ac.id',NULL,NULL,NULL,NULL,NULL,'2023-10-16 02:20:00','2025-10-02 06:01:11'),
            (220, UUID(),'Ellyn Octavianty','$2y$10$9eF012HeKIveFdXo7wibwenBg4FY3Mwua.p6Y9q.GysWdWG9UXfa6','auditor2','Ellyn Octavianty, SE., MM','ellynoctavianty@unpak.ac.id',NULL,NULL,NULL,NULL,NULL,'2023-10-16 02:20:10','2024-11-30 08:56:09'),
            (223, UUID(),'Heny Purwanti','$2y$10$j./Fj3eS3ZiaxF8C09Zxs.olJ6r6Arcys2Nlm0j7Jqk74RhF9pNCW','auditor2','Heny Purwanti, M.T','henypurwanti@unpak.ac.id',NULL,'FOTO.jpg',NULL,NULL,NULL,'2023-10-16 02:20:20','2024-12-16 12:21:31'),
            (226, UUID(),'Patar Simamora','$2y$10$k9Bn3veF6Wj7dCxjUDGf3OryaWi9FvzTHNooQmcYI4EWNplJb2q6u','auditor2','Patar Simamora, SE., M.Si.','patar.simamora@unpak.ac.id',NULL,NULL,NULL,NULL,NULL,'2023-10-16 02:20:33','2024-12-30 07:05:10'),
            (229, UUID(),'Solihin','$2y$10$5onF1l43ZlfS3xK/8h.km.qfp.MNO0zJD8ch4zk3IFXY2FwvgGXey','auditor2','Solihin, M.T','solihin@unpak.ac.id',NULL,'Solihin.jpg',NULL,NULL,1,'2023-10-16 02:20:50','2024-12-19 04:31:48'),
            (232, UUID(),'Fakultas Hukum','$2y$10$9LPF7aie0Yc1F2kAp8yXRu9XyYJpmon1BfMm9ZrUPngDX7j3QvxC6','auditee','Dr. Eka Ardianto Iskandar, S.H., M.H. (Dekan FH)','fakultashukum@unpak.ac.id',91,'OIP (1).jfif',NULL,NULL,NULL,'2023-10-16 02:25:53','2025-10-29 06:08:05'),
            (235, UUID(),'Ilmu Hukum','$2y$10$gnvb6sNXoPUrkkpCHiwS9uLLY0xQm4LQ2L9osFjJ8j2b7Ke1IKw26','auditee','Ari Wuisang, SH., MH. ( Kaprodi Ilmu Hukum S1)','fakultashukum@unpak.ac.id',1,'OIP (1).jfif',NULL,NULL,NULL,'2023-10-16 02:31:24','2025-10-29 06:52:40'),
            (241, UUID(),'FEB','$2y$10$JYKed7BSLvrppl5/U91AUOYQk.C1B5yj1Jbe94GxAsH9AKKalr5yW','auditee','Towaf Totok Irawan, S.E., M.E., Ph.D (Dekan FEB)','dummy@gmail.com',94,NULL,NULL,NULL,NULL,'2023-10-17 03:15:51','2025-10-14 07:23:52'),
            (244, UUID(),'Akuntansi','$2y$10$N6VxM0ovQHpgE7OHIAgs5ehuc8peG.RJLsKVd22Iy6UPb2YuNgfJu','auditee','Dr. Heru Satria Rukmana, S.E., Ak., M.M (Kaprodi Akuntansi)','akuntansi@unpak.ac.id',4,NULL,NULL,NULL,NULL,'2023-10-17 03:24:31','2025-10-11 03:15:10'),
            (247, UUID(),'Manajemen','$2y$10$j0C06V86F4NBXEITynbaBupk8ZVFtvFkLKPN3lKQA2/d0wri5vjMS','auditee','Prof.Dr. Yohanes Indrayono, Ak.,MM (Kaprodi Manajemen S1)','upmps.s1mjn@unpak.ac.id',2,'Logo Mjn New.jpg',NULL,NULL,NULL,'2023-10-17 03:26:11','2025-10-04 05:18:10'),
            (250, UUID(),'Bisnis Digital','$2y$10$ENlyFXwPHOVL.TBwxIlOLOnxS7sFQxG/stsGPZY2YFCQ7rZQN3iLO','auditee','Dr. Abel Gandhy, S.Pi., MM (Kaprodi Bisnis Digital)','bisnisdigital@unpak.ac.id',3,'bdi.png',NULL,NULL,NULL,'2023-10-17 03:32:18','2025-10-28 03:36:42'),
            (253, UUID(),'FKIP','$2y$10$99t67kLzuqg7FiQ.ypVGgOL6fBxEqEzS3UtQ10mEussugvrJgVdpy','auditee','Dr. H. Eka Suhardi, M.Si. (Dekan FKIP)','fkip@unpak.ac.id',97,NULL,NULL,NULL,NULL,'2023-10-17 03:33:57','2025-10-21 02:35:47'),
            (256, UUID(),'PBSI','$2y$10$I1YOEDUqUD//yo1WVz8rKOSiUEhwzW5um6PH4IPN3G9YkYHjmIxwe','auditee','Stella Talitha, M.Pd. (Kaprodi PBSI)','fkip.indo@unpak.ac.id',10,'Stella Talitha.jpeg',NULL,NULL,NULL,'2023-10-17 03:36:08','2025-10-30 01:32:00'),
            (259, UUID(),'PBI','$2y$10$dmqnLUUrlZEHkf/hRQbqXuXKaCm8VPDbIPlOlSUU7zG6Y3LvAc.46','auditee','Abdul Rosyid, M.Pd. (Kaprodi PBI)','englishedu.fkip@unpak.ac.id',11,'Abdul Rosyid_Foto.jpg',NULL,NULL,NULL,'2023-10-18 02:23:48','2025-11-24 09:21:05'),
            (262, UUID(),'Pen.Biologi','$2y$10$EIqS5RmbzDgGsw1H4FJ.r..Y3ZHUe/Lj4RLtdOUfnCyLYED99GOWq','auditee','Lufty Hari Susanto, M.Pd. (Kaprodi P.Biologi)','pendbiologifkip@unpak.ac.id',5,'WhatsApp Image 2025-06-20 at 15.33.04_1bad4e61.jpg',NULL,NULL,NULL,'2023-10-18 02:24:44','2025-11-18 01:42:16'),
            (265, UUID(),'PGSD','$2y$10$9jLvr.47Rl.4TDbTmCkNBunHYa/x96PfG3dTHg9uv7dD.WOKpAIuy','auditee','Dr. Nita Karmila, M.Pd. (Kaprodi PGSD)','dummy@gmail.com',7,NULL,NULL,NULL,NULL,'2023-10-18 02:25:27','2025-10-04 05:20:57'),
            (268, UUID(),'Pen.IPA','$2y$10$L7V4B2hifItRUhPM0OkMLuo8PykCHdzAtHMDUKJCCJkbYpaEGKrza','auditee','Lilis Supratman, M.Si. (Kaprodi P.IPA)',NULL,NULL,NULL,NULL,NULL,0,'2023-10-18 02:26:04','2023-10-18 02:26:04'),
            (271, UUID(),'PPG','$2y$10$.PoPnL6lGPDuGtwqfyv5ouvwVpbNdlWl9LhkYd830BX./4y1tXTTu','auditee','Dr. Indri Yani, M.Pd (Kaprodi PPG)','ppg@unpak.ac.id',9,NULL,NULL,NULL,NULL,'2023-10-18 02:27:00','2025-10-30 03:52:57'),
            (274, UUID(),'Fisib','$2y$10$XaYcZ/d9YAjreFKlyAH52.ikqE8jbiTahU2xC188gNt1FpvO7d./i','auditee','Dr. Muslim, M.Si (Dekan Fisib)','dekanfisib.2020@unpak.ac.id',100,NULL,NULL,NULL,NULL,'2023-10-18 02:28:09','2025-10-04 05:21:36'),
            (277, UUID(),'Sastra Inggris','$2y$10$a8Tmi80DmTPnE.1zvNmvj./xe3RRr.wAA/nNbLxYFfMY0uMaJ3wWK','auditee','Dyah Kristyowati,S.S.,M.Hum. (Kaprodi Sas.ing)','dummy@gmail.com',14,NULL,NULL,NULL,NULL,'2023-10-18 02:33:50','2025-10-04 05:21:46'),
            (280, UUID(),'Sastra Jepang','$2y$10$XbWkZ2dchRgIlPIMqCuAruTSJ36LPZcxlQvmc9u4CrPckZzLTkhVa','auditee','Mugiyanti, M.Si (Kapordi Sas.Jep)','fisib.sasjep@unpak.ac.id',15,NULL,NULL,NULL,NULL,'2023-10-18 02:35:50','2025-10-04 05:21:54'),
            (283, UUID(),'Sastra Indonesia','$2y$10$PPNip.mCc6n1vnfePsJT5uiulw7YjJBwmwv34Y7yV5GHEyLWiH3fe','auditee','Drs. Sasongko Suharto Putro, M.M. (Kaprodi Sas.In)','prodisastraindonesiaunpak@gmail.com',13,NULL,NULL,NULL,NULL,'2023-10-18 02:37:41','2025-10-22 07:29:19'),
            (286, UUID(),'I.kom','$2y$10$i4Cdm/vLB.13TpgBnQcrYOlb.7bDkahySy2YS6U/DQBozvK.NX/km','auditee','Ratih Siti Aminah., M.Si. (Kaprodi I.Kom)','rinifirdaus@unpak.ac.id',12,'Wajah 2.jpg',NULL,NULL,NULL,'2023-10-18 02:45:04','2025-10-04 05:22:19'),
            (289, UUID(),'FTeknik','$2y$10$VE6X7dMwQLGMzNhbzxFMUuZzedcAv1sxL.oAQRr.4BpU33Pf7Xyqe','auditee','Dr. Ir. Lilis Sri Mulyawati, M.Si (Dekan FT)','mutuft@unpak.ac.id',103,NULL,NULL,NULL,NULL,'2023-10-18 03:10:28','2025-10-17 03:35:40'),
            (292, UUID(),'T.Geologi','$2y$10$bXTc1BuOZjVjM4SNj/e.8Ot4ZNHZQJpVlzsfV9lGcV4qf9cQPX.pq','auditee','Helmi Setia Ritma P., ST., M.Si. (Kaprodi Geologi)','solihin@unpak.ac.id',19,'WhatsApp Image 2025-11-06 at 14.00.27_f6c8fe96.jpg',NULL,NULL,NULL,'2023-10-18 03:11:39','2025-11-06 07:01:50'),
            (295, UUID(),'PWK S1','$2y$10$gykHif0r19MQZJFaxiXXOueXwyHFAdf3sBkac6k5dA48aHrXFj0fS','auditee','Dr. Mujio, S.Pi., M.Si (Kaprodi PWK S1)','prodipwk@unpak.ac.id',20,'Mujio.png',NULL,NULL,NULL,'2023-10-18 03:15:07','2025-11-06 13:37:58'),
            (298, UUID(),'T.Sipil','$2y$10$AJQRgqp9nxPMQz5ruhNb8eqO7FyLUmBoR9NfkQuanq9edyAc.NxI2','auditee','Ir. Wahyu Gendam P, STP., M.Si (Prodi T.Sipil)','mutuft@unpak.ac.id',17,'Screenshot 2025-10-23 193738.png',NULL,NULL,NULL,'2023-10-19 02:01:22','2025-10-23 12:37:52'),
            (301, UUID(),'T.Elektro','$2y$10$WTSJE/Ew78A6V4V.UQVxi.1PO9hadLo7qob3BKOVvYd6Y1kCth/4C','auditee','Ir. Yamato, M.T (Kaprodi T.Elektro)','waryani@unpak.ac.id',16,'Capture.JPG',NULL,NULL,NULL,'2023-10-19 02:03:06','2025-10-04 05:23:18'),
            (304, UUID(),'T.Geodesi','$2y$10$Nl6GnpPtKwt3yJGaAd9X0uh15g7lnYg5NAZ4v5gNjrZq7UFDNTNQi','auditee','Mohamad Mahfudz, ST., MT. (Kaprodi T.Geodesi)','prodi_geodesi@unpak.ac.id',18,'kap_gd.png',NULL,NULL,NULL,'2023-10-19 02:04:27','2025-11-07 03:38:56'),
            (307, UUID(),'FMIPA','$2y$10$Wa.T9Ktxf2zH/vlscTqUlOGGP/Bhb9qwNNQFwxPftJ0ipixjq8IWy','auditee','Asep Denih, S.kom., M.Sc., Ph.D. (Dekan Fmipa)','asep.denih@unpak.ac.id',106,NULL,NULL,NULL,NULL,'2023-10-19 02:06:37','2025-10-04 05:23:52'),
            (310, UUID(),'Biologi','$2y$10$xwPB65hlQIHQyh6bEv4A2Os.OvSphMt8U0heQUZKOmvNy3Y.FXtki','auditee','Dra. Triastinurmiatiningsih, M.Si','dummy@gmail.com',22,NULL,NULL,NULL,NULL,'2023-10-19 02:08:31','2025-10-20 06:16:21'),
            (313, UUID(),'Kimia','$2y$10$Jw3Ed9Cgt5iFBr.dge8kCOwVbBhN2ZMFTlELc38rAmxvZajguVIwG','auditee','Dr. Uswatun Hasanah, S.Si., M.Si. (Kaprodi Kimia)','kimia@unpak.ac.id',23,'Screenshot 2025-10-29 at 10.58.07.png',NULL,NULL,NULL,'2023-10-19 03:21:12','2025-10-29 03:58:54'),
            (316, UUID(),'Matematika','$2y$10$7Fc3q5JKhR.dwDY1JqX.C.DgO1niTN/94AFHuUSUPoEVqgdMqPT2a','auditee','Dr. Embay Rohaeti, S.Si., M.Si. (Kaprodi MTK)','matematika@unpak.ac.id',21,NULL,NULL,NULL,NULL,'2023-10-19 03:22:07','2025-10-04 05:24:27'),
            (319, UUID(),'Ilkom','$2y$10$b0mxxIzVkTQgLrAoTbVpquxFUi0pxGKLfzNyNrrUbLOXB5oU7p6wm','auditee','Dr. Fajar Delli Wihartiko, S.Si., M.M., M.Kom (Kaprodi Ilkom)','akreilkom@unpak.ac.id',25,NULL,NULL,NULL,NULL,'2023-10-19 03:53:48','2025-10-04 05:24:41'),
            (322, UUID(),'Farmasi','$2y$10$inii6GEuCjDK2CUIhygo6uT94ZyhLZZJy.6TblqAvluYYU6qqLkB.','auditee','apt. Emy Oktaviani, S.Farm., M.Clin.Pharm. (Kaprodi Farmasi)','cyntiawahyuningrum@unpak.ac.id',24,NULL,NULL,NULL,NULL,'2023-10-19 03:54:17','2025-10-04 05:24:50'),
            (325, UUID(),'Pascasarjana','$2y$10$QToSV66aAzlnSvVsd/b0FOnHFUnrRisVgOxV8il/qp0oWiDx4IsF2','auditee','Prof. Dr. Sri Setyaningsih, M.Si (Dekan Pascasarjana)','pasca@unpak.ac.id',109,NULL,NULL,NULL,NULL,'2023-10-20 04:15:35','2025-10-04 05:25:08'),
            (328, UUID(),'MP S3','$2y$10$h8hDBbEPtnm6uwU6iL9iEe3tDxFPWTvR9sLUra876BkRwsSB292uy','auditee','Dr. Suhendra, M.Pd. (Kaprodi MP S3)','dummy@gmail.com',31,NULL,NULL,NULL,NULL,'2023-10-20 04:17:29','2025-10-04 05:25:20'),
            (331, UUID(),'Manajemen S3','$2y$10$IpWslWUPWI.qH1sfkAJ4y.3d7l3MrUE63RmhAVcObRosentB1zPtu','auditee','Dr. Nancy Yusnita, SE., MM  (Kaprodi Manajemen S3)','dummy@gmail.com',26,'Nanc.jpg',NULL,NULL,NULL,'2023-10-20 04:18:26','2025-10-22 07:00:56'),
            (334, UUID(),'MP S2','$2y$10$ZlbC6RizafqHg5lvcK8u/ur5hPtOBisX0qysnpBcfWo2tbuJKWJjK','auditee','Dr. Lina Novita, M.Pd. (Kaprodi MP S2)','dummy@gmail.com',32,NULL,NULL,NULL,NULL,'2023-10-20 04:19:54','2025-10-04 05:25:40'),
            (337, UUID(),'ML S2','$2y$10$UrexCNIJ/8brJw3Sl7luSOdH1.mVKWkaDF18lKcG097DEns9AXZfG','auditee','Dr. Rosadi, SP, MM (Kaprodi ML S2)','dummy@gmail.com',28,NULL,NULL,NULL,NULL,'2023-10-20 04:20:39','2025-10-04 05:25:47'),
            (340, UUID(),'Ilmu Hukum S2','$2y$10$H4FsiRmvZJQi4gMHDug7WuL/rO903Ubl2y4fWCPK7OwsHDe/YaZ3e','auditee','Dr. Iwan Darmawan, SH., MH. (Kaprodi Ilmu Hukum S2)','dummy@gmail.com',29,NULL,NULL,NULL,NULL,'2023-10-20 04:21:23','2025-10-22 06:46:26'),
            (343, UUID(),'Manajemen S2','$2y$10$IAMMrq7xpwm4CjTxwRs1zuMbrOY0GNIotCqTWGj1qOh6sVzhvQI5u','auditee','Dr. Agus Setyo Pranowo, MM., SE. (Kaprodi Manajemen S2)','dummy@gmail.com',27,NULL,NULL,NULL,NULL,'2023-10-20 04:22:12','2025-10-04 05:26:07'),
            (346, UUID(),'IPA S2','$2y$10$glsT8Xpsw1xWSE7bYTcuNuh8vuWuqPY5ljoYWLqFhNZpfKBQmmkzu','auditee','Dr. Didit Ardianto, M.Pd. (Kaprodi IPA S2)','dummy@gmail.com',30,NULL,NULL,NULL,NULL,'2023-10-20 04:23:01','2025-10-04 05:26:15'),
            (349, UUID(),'PWK S2','$2y$10$YjmFMhfBAPwHBiYXNTIwIOZ77/tLoB/c5RjgVOcoQO6s/7nNShWQS','auditee','Dr. Ir. Anugrah, M.Si. (Kaprodi PWK S2)','dummy@gmail.com',35,NULL,NULL,NULL,NULL,'2023-10-20 04:23:32','2025-10-04 05:26:27'),
            (352, UUID(),'PENDAS S2','$2y$10$WbTeNjiAcealTWgejX5bxOfIDJAkpvl5wWpz8ILtkmMWHc9bcTpSy','auditee','Dr. Tustiyana Windiyani, M.Pd. (Kaprodi PENDAS S2)','tustiyana@unpak.ac.id',33,'Foto mamah windi 4 x 6 (1).jpeg',NULL,NULL,NULL,'2023-10-20 04:24:17','2025-10-04 05:26:35'),
            (355, UUID(),'Svokasi','$2y$10$1NtL4tJs1/w4r9CRtF5q4.JydIXxQJ0Y4natXVPBIwSwd3u6qUhCG','auditee','Dr. Lia Dahlia Iryani, S.E., M.Si (Dekan SV)','sonniadarmasih123@gmail.com',112,'hehe.php',NULL,NULL,NULL,'2023-10-21 01:49:01','2025-10-04 05:27:45'),
            (358, UUID(),'0413117601','$2y$10$zO51umOGlCIrOD9QS61yu.g0snHAbt33BZd2H0BwlffgDNFuqbNsO','auditee','Dr. Lia Dahlia Iryani, SE., M.Si ( KaProdi Akuntansi D3)','dahlia.iryani@unpak.ac.id',NULL,'logo vokasi.jpg',NULL,NULL,1,'2023-10-21 01:50:13','2024-12-21 01:55:43'),
            (364, UUID(),'perpajakan D3','$2y$10$R1uZF6/WbF0Gmra7nuDsDupgRJyWt8tZhMMbwAFJi2plRL0EK278m','auditee','Chandra Pribadi, Ak., M.Si., CPA. (Kaprodi Perpajakan D3)','dummy@gmail.com',38,NULL,NULL,NULL,NULL,'2023-10-27 03:40:58','2025-10-04 05:28:34'),
            (367, UUID(),'MPK D3','$2y$10$WVvjKmREhipH3qnw.ZQqMu5uS/gJX6eczgOrKv35E19uv6j..qGTe','auditee','Djoko Hardjanto, S.Pt, M.Si (Kaprodi MKP D3)','d3.mkp@unpak.ac.id',NULL,NULL,NULL,NULL,1,'2023-10-27 03:43:34','2024-12-13 08:43:20'),
            (370, UUID(),'T.Kom D3','$2y$10$8WZTMqtZzOlcO/PCbFDOQupDUvECWvneNQgjzRt4dHAX7I0NMcMM6','auditee','Akbar Sugih Miftahul Huda, M. Kom. (Kaprodi Teknik Komputer D3)','dummy@gmail.com',36,NULL,NULL,NULL,NULL,'2023-10-27 03:44:50','2025-10-22 03:37:06'),
            (373, UUID(),'S.Informasi','$2y$10$vAQBpJ./w9h0cL4IZsdXMOCaRL.1ngKGKNiw.d/HJb/WlZFASj9ku','auditee','Dian Kartika Utami, M.Kom (Kaprodi Manajemen Informatika d3)','sv_si@unpak.ac.id',37,'C83C3368 crop.JPG',NULL,NULL,NULL,'2023-10-27 03:46:13','2025-11-13 04:19:46'),
            (376, UUID(),'aries','$2y$10$o.6Hly1CHcBQ4hS0.hg7q.0FLDx9kC4jLvUuEX9s4CORZjsHTDyE2','auditor1','Aries Maesya','adamilkom00@gmail.com',NULL,NULL,NULL,NULL,1,'2023-10-27 11:39:16','2024-11-29 11:22:59'),
            (379, UUID(),'KKPKT','$2y$10$nFz1VGY/wBbBSkOkKff4E.Msn6xJf1Mn5K0oSxmI/9nPZQTkR5BGu','auditee','KKPKT',NULL,NULL,NULL,NULL,NULL,0,'2023-11-03 01:55:23','2023-11-03 01:55:23'),
            (382, UUID(),'LPPM','$2y$10$.r1mSpN2ZkMnFoSSRWMkA.rB9R8kNiScU5Zk.YMop1MjX2XjFxlDm','auditee','Dr. Dolly Priyatna, M.Si','dollypriatna@unpak.ac.id',70,NULL,NULL,NULL,NULL,'2024-01-08 06:06:20','2025-09-01 06:09:27'),
            (385, UUID(),'UNPAK PRESS','$2y$10$iTd7C7cPIFxjrMyJAp9ppetWAkLQS.TIL42.uU54SJCOlNBZ9nFfe','auditee','Nina Agustina ,S.E.,M.E','dummy@gmail.com',75,NULL,NULL,NULL,NULL,'2024-01-08 06:07:07','2025-09-01 06:08:56'),
            (388, UUID(),'HUMAS','$2y$10$9bJVYOSArUTdOQMRY573ru4F.ZkzaPQkIejN8rUr9KinYsPmeDkuO','auditee','Aditya Prima Yudha, S.Pi.,M.M.','dummy@gmail.com',72,NULL,NULL,NULL,NULL,'2024-01-08 06:07:36','2025-09-01 06:11:02'),
            (391, UUID(),'PERPUS_PUSAT','$2y$10$gk0H/5NydAKvXBacuSlKmu0cXdwOzhtoYuzbRWR1fCcBrCcTNvfyS','auditee','Wildan Fauzi Mubarock M. Pd','dummy@gmail.com',71,NULL,NULL,NULL,NULL,'2024-01-08 06:08:04','2025-09-01 06:17:38'),
            (394, UUID(),'INOVASI','$2y$10$eckTCXLOu/RzN5FX/IkW9uEZMCxw3TcJB.i744CtYhfbXRP65SxRe','auditee','Asep Saepulrohman, M.Si','asepspl42@gmail.com',88,NULL,NULL,NULL,NULL,'2024-01-08 06:08:30','2025-10-13 02:55:15'),
            (397, UUID(),'INKUBATOR','$2y$10$u6rmC98w0a.xOCtMkoEsGufnb1rdFrKfIgF3LQYfy24.wkjCQr7Ei','auditee','Asep Saepulrohman, M.Si','asepspl42@gmail.com',89,NULL,NULL,NULL,NULL,'2024-01-08 06:09:01','2025-10-13 02:55:25'),
            (400, UUID(),'KEMITRAAN','$2y$10$T8J5cfA3e2kEhakjA7EKYesqCr4shHNVCB/PDKxfdMYLjA7zllDnC','auditee','Cucu Mariam, M.Pd','dummy@gmail.com',125,NULL,NULL,NULL,NULL,'2024-01-08 06:09:32','2025-09-01 06:12:04'),
            (403, UUID(),'KARIR','$2y$10$s3EjXRoNPAO1qcG4flXc7.kpQbnD8XtSfpNbwcfGae/i2S7Qxflgu','auditee','Dr. Herman, SE., MM.','dummy@gmail.com',68,NULL,NULL,NULL,NULL,'2024-01-08 06:09:55','2025-09-01 06:12:58'),
            (406, UUID(),'BPSI','$2y$10$5H2cyGbrhV0NDNS1u9GjvO8avGTOTze069.U3Y5tQlXaz0/hqR5iK','auditee','Aries Maesya, M.Kom','putik@unpak.ac.id',41,NULL,NULL,'$2y$10$7EZZYmgnfQ2oKb4QD1ZyIOPmL8pg3f4Mma1vnzr8dXgfLWLCgmY5O',NULL,'2024-01-11 07:28:00','2025-09-12 03:18:28'),
            (407, UUID(),'BAAK','$2y$10$/5Q9ilEgbJ5iGuZ3rABJJe419Iptj1n5kpOzyItwHFIWVXxGnHyva','auditee','Dr. Eka Ardianto Iskandar, S.H., M.H.','unpak@gmail.com',91,NULL,NULL,NULL,NULL,'2024-01-26 04:23:35','2025-10-04 06:32:15'),
            (410, UUID(),'BAUM','$2y$10$lomy6xBGmv0nIDkpUToxjO0m3uWPLoedAj.hCJJ/PhFK3RkaCMCQK','auditee','Wijaya Kusumah, S.E. (Kepala BAUM)',NULL,NULL,NULL,NULL,NULL,0,'2024-01-26 04:23:48','2024-02-02 01:42:32'),
            (413, UUID(),'LPM','$2y$10$wF76ZOkxWSY42maN5licUu8BKbJsPb3BEaV6zNNEtH7QRYp0ixcx6','auditee','Dr. Diana Widiastuti, M.Phil','lpm@unpak.ac.id',69,NULL,NULL,NULL,NULL,'2024-02-02 02:20:14','2025-08-29 04:40:11'),
            (415, UUID(),'Feri Ferdinan','$2y$10$Xi2k9oUCDaYWjSR0KX0GceYECXGvDMslSJ09x1.K4yTf8x/VYN4qq','auditor2','Dr. Feri Ferdinan Alamsyah, M.I.Kom.','feriferdinan@unpak.ac.id',NULL,'003.jpeg',NULL,NULL,1,'2024-02-15 05:53:15','2024-12-14 01:22:45'),
            (418, UUID(),'aries2','$2y$10$/sEjEdpfcHfl06vV3VBoRepVXIS9OvAImHmALbBltllp/5iqhdzX.','auditor2','aries2','adamilkom00@gmail.com',NULL,NULL,NULL,NULL,0,'2024-02-15 05:53:15','2024-11-29 07:00:28'),
            (421, UUID(),'Mahfudz','$2y$10$x6g4bNeJWmxOsGM14CzOIO3BgcgoULrtLvb4StI36EyBVa1VVa9Pe','auditor2','M. Mahfudz, M.T','mohamadmahfudz@unpak.ac.id',NULL,NULL,NULL,NULL,NULL,'2024-09-18 02:17:37','2025-11-26 03:08:11'),
            (425, UUID(),'testauditee','$2y$10$/rcTO5xNMGIecshQKDabFe2zSUZ2hQGJrI37aZneWlIoVFMriupsO','auditee','testauditee','testauditee@gmail.com',91,NULL,NULL,NULL,NULL,'2024-10-12 02:52:23','2025-11-20 14:34:53'),
            (428, UUID(),'testauditor1','$2y$10$VOPdhlRek2vzKfkBg5wYqeDCl4IyS8biHoZRcPNmc4QPBXxyyFb3a','auditor1','testauditor1','testauditor1@gmail.com',NULL,NULL,NULL,NULL,NULL,'2024-10-12 02:52:55','2025-11-20 14:35:21'),
            (431, UUID(),'testauditor2','$2y$10$yXuEfDS8owHNLilpdm9WfOF3wwlhkmqWaCV8XKoH0N.bE7ZPF0Pba','auditor2','testauditor2','testauditor2@gmail.com',NULL,NULL,NULL,NULL,NULL,'2024-10-12 02:53:04','2025-11-20 14:35:46'),
            (452, UUID(),'Feri Ferdinan','$2y$10$4k03/IVfdC.SUFwoMATG3.2Zw8SLGHI2rUQHZvl1AJ/ovfbWYPFmO','auditor1','Dr. Feri Ferdinan Alamsyah, M.I.Kom.','feriferdinan@unpak.ac.id',NULL,NULL,NULL,NULL,0,'2025-02-04 08:20:48','2025-02-04 08:20:48'),
            (455, UUID(),'Mahfudz','$2y$10$QFZ2xuoSDdI7CXwfLsariuX8GgqHs22dVsL4tFhvOTa7f7eWvnLlS','auditor1','M. Mahfudz, M.T','mohamadmahfudz@unpak.ac.id',NULL,NULL,NULL,NULL,NULL,'2025-02-04 08:21:20','2025-11-26 03:07:51'),
            (462, UUID(),'Diana Widiastuti','$2y$10$0jfbqSIIwAA61HQ29L2MAeN7T3YwEELFP/90i.kMUv/MXLea9D0Jy','auditor1','Dr. Diana Widiastuti, S.Si, M.Phil','diana@gmail.com',NULL,NULL,NULL,NULL,NULL,'2025-08-22 03:03:41','2025-11-18 03:07:55'),
            (463, UUID(),'Diana Widiastuti','$2y$10$MX0iu.T.hZp6SPzxowbo5.DPH4HXnSe0pGsz0UdrnR5YQmOqUf9UO','auditor2','Dr. Diana Widiastuti, S.Si, M.Phil','diana@gmail.com',NULL,NULL,NULL,NULL,NULL,'2025-08-22 03:04:18','2025-11-18 03:21:38'),
            (464, UUID(),'Muhammad Fathurrahman','$2y$10$XP40Y3XnMkbXin4kKQ4Q9.EpfTHJFPiyp91ZfTMZmaHg7S4UY1Ct.','auditor1','Dr. Muhammad Fathurrahman, S.Pd., M.Si.','fathur110590@unpak.ac.id',NULL,NULL,NULL,NULL,NULL,'2025-08-22 03:05:03','2025-11-18 03:25:41'),
            (465, UUID(),'Muhammad Fathurrahman','$2y$10$94/7Jt7BWU.oB7jOo5cGzeXJVfHfTDmIky7a7XVfOm7sRde66yjTC','auditor2','Dr. Muhammad Fathurrahman, S.Pd., M.Si.','fathur110590@unpak.ac.id',NULL,NULL,NULL,NULL,NULL,'2025-08-22 03:05:27','2025-11-18 03:07:31'),
            (466, UUID(),'Kecerdasan Buatan dan Robotika','$2y$10$AAwGvQ2IsbuA69LQNoGGUeHM3F8QSx/aOdYWiXnAIdB.T2eZM5xxu','auditee','Fikri Adzikri, S.T., M.T. (Ka Prodi Kecerdasan Buatan dan Robotika)','dummy@gmail.com',121,NULL,NULL,NULL,NULL,'2025-08-25 04:54:19','2025-10-04 05:30:02'),
            (467, UUID(),'Ilmu Komputer S2','$2y$10$l.FSqUg2sHKY7KaP4VM2ie1sx2Dan9vqCkbevSCfwYs9eARBAgdjy','auditee','Dr. Eneng Tita Tosida, M.Si. (Ka Prodi Ilmu Komputer S2)','dummy@gmail.com',126,NULL,NULL,NULL,NULL,'2025-08-26 04:25:15','2025-10-04 05:26:53'),
            (468, UUID(),'PPA','$2y$10$XeNkzrKmAW9ICXIbtz0xZessMKo7caCyyjTcxZedCTAjVJtEz8PZq','auditee','Dr. apt. Bina Lohitasari, M.Pd., M.Farm','binalohitasari@unpak.ac.id',127,'3x4.jpg',NULL,NULL,NULL,'2025-08-26 04:27:15','2025-11-12 08:55:50'),
            (469, UUID(),'E-Journal FKIP','$2y$10$fKdarrKnt4UXV.vBZaqdrOm61XC16dN/UffGZLlHKb2O/Ra0FmOiW','auditee','E-Journal (Unit FKIP)','dummy@gmail.com',78,NULL,NULL,NULL,NULL,'2025-10-07 03:12:29','2025-10-07 03:13:52'),
            (470, UUID(),'KKN FKIP','$2y$10$ct8G8WdqsA6icA2A0Gz0O.bIeav1synwx2AcrSBZD7i1XEHa3A3/.','auditee','KKN (Unit FKIP)','dummy@gmail.com',76,NULL,NULL,NULL,NULL,'2025-10-07 03:13:30','2025-10-07 03:14:18'),
            (471, UUID(),'Bimbingan Konseling FKIP','$2y$10$uLNPCLZw5YrPiiuAxUeGL.0SFLusOnWboq411VmZrc7305fSl/3pW','auditee','Bimbingan Konseling (Unit Fkip)','dummy@gmail.com',87,NULL,NULL,NULL,0,'2025-10-07 03:18:41','2025-10-07 03:18:41'),
            (472, UUID(),'Laboratorium Seni (FKIP)','$2y$10$Z93xpkC/R2v4/7CPFKwdFOBegjbuVfN4GvsM6p9CkWLruABNwsAiO','auditee','Laboratorium Seni (Unit FKIP)','dummy@gmail.com',128,NULL,NULL,NULL,NULL,'2025-10-07 04:05:52','2025-10-11 02:53:19'),
            (473, UUID(),'Laboratorium Bahasa Inggris (FKIP)','$2y$10$JSmtsvgQqQEi218JqrcT8ezb8VPQLJ5WSVjKrnt.H4Crj2Ih6/NeW','auditee','Laboratorium Bahasa Inggris (Unit FKIP)','dummy@gmail.com',129,NULL,NULL,NULL,0,'2025-10-11 02:54:30','2025-10-11 02:54:30'),
            (474, UUID(),'Laboratorium Microteaching (FKIP)','$2y$10$WBnKtk4qhQ8z0z0.35HwKOsAe6J3BPWFRJTKggQCwnd1d6Dj5oUGe','auditee','Laboratorium Microteaching (Unit FKIP)','dummy@gmail.com',130,NULL,NULL,NULL,0,'2025-10-11 02:55:21','2025-10-11 02:55:21'),
            (475, UUID(),'Laboratorium Komputer (FKIP)','$2y$10$niw9zkI3DVJ9B1WrPd9PuupZjaZeGSikPlZGdzkiWgNHSHUqyZ5IS','auditee','Laboratorium Komputer (Unit FKIP)','dummy@gmail.com',131,NULL,NULL,NULL,0,'2025-10-11 02:56:48','2025-10-11 02:56:48'),
            (476, UUID(),'AktStudio','$2y$10$zjpiSGILR4dFkAkqS3bmDuPAq40exCd66xKTTchVip4ilvY/KyQe.','admin','Admin','adamilkom00@gmail.com',NULL,'carbon (1).png',NULL,NULL,NULL,'2023-05-18 18:45:52','2025-10-13 07:28:13'),
            (478, UUID(),'Verifikator FH','$2y$10$om1jTBXX9ZXW9JOSVwzGBOpB4uYNaS9NhyXpZb.nPkgaKdFbFVPUO','fakultas','Verifikator FH','dummy@gmail.com',91,NULL,NULL,NULL,NULL,'2025-10-15 01:07:12','2025-10-15 01:08:54'),
            (479, UUID(),'Verifikator FEB','$2y$10$FfRY9P4cO3YdfaLNH.Gp.OcIXkD7R25Iiv0H4azIfInRwMQoLf4P2','fakultas','Verifikator FEB','dummy@gmail.com',94,NULL,NULL,NULL,0,'2025-10-15 01:07:59','2025-10-15 01:07:59'),
            (480, UUID(),'Verifikator FKIP','$2y$10$prYvCgwMXmk1rZpfePGtm.ssqkqNt6.3xBV4WJcirI95Bpfc4IF0O','fakultas','Verifikator FKIP','dummy@gmail.com',97,NULL,NULL,NULL,NULL,'2025-10-15 01:16:28','2025-10-18 01:53:07'),
            (481, UUID(),'Verifikator FISIB','$2y$10$7kPzGY7QrUF5UVLaM0.MmOQ/unwDwnPvb.5/BjFxn63EY5miPeNzu','fakultas','Verifikator FISIB','dummy@gmail.com',100,NULL,NULL,NULL,0,'2025-10-15 01:27:55','2025-10-15 01:27:55'),
            (482, UUID(),'Verifikator FT','$2y$10$BMoP02xUWj32iEgluRaiZOTwe3L.CCfmWnEBNln0cqWQSTs.DR6Fi','fakultas','Verifikator FT','dummy@gmail.com',103,NULL,NULL,NULL,0,'2025-10-15 02:04:20','2025-10-15 02:04:20'),
            (483, UUID(),'Verifikator FMIPA','$2y$10$Xv3wNbohilhZ14qRANgI.uj6UphgqYj4xoOVBwWGxkzBe7gXSG64u','fakultas','Verifikator FMIPA','dummy@gmail.com',106,NULL,NULL,NULL,0,'2025-10-15 02:04:42','2025-10-15 02:04:42'),
            (484, UUID(),'Verifikator Pascasarjana','$2y$10$OtQpv.Na9bPlmC4fv.9T0.ZoQ0MoadpwKVfQ.exHadWuVfpIbBy..','fakultas','Verifikator Pascasarjana','dummy@gmail.com',109,NULL,NULL,NULL,NULL,'2025-10-15 02:06:04','2025-10-15 02:06:22'),
            (485, UUID(),'Verifikator Vokasi','$2y$10$D.S5m6co.G8Qf.7z4UFATexnxvZKtNcxQ2i9Vf5RpsUa3wGnhLXLS','fakultas','Verifikator Vokasi','dummy@gmail.com',112,NULL,NULL,NULL,0,'2025-10-15 02:06:45','2025-10-15 02:06:45'),
            (486, UUID(),'Singgih Irianto','$2y$10$og1dbBHJn4cSP9f0PdQpTugmEPBEfzDZYU9PRPCZnhy4oelLnx77K','auditor1','Dr. Singgih Irianto, MT','dummy@gmail.com',NULL,NULL,NULL,NULL,NULL,'2025-11-18 02:39:20','2025-11-18 03:08:03'),
            (487, UUID(),'Singgih Irianto','$2y$10$4hShGb.9/B0m09SayZq2hOXlQQ5s1V/BqdUw0l3f3HVIuVURHStNi','auditor2','Dr. Singgih Irianto, MT','dummy@gmail.com',NULL,NULL,NULL,NULL,0,'2025-11-18 02:53:45','2025-11-18 02:53:45'),
            (488, UUID(),'Siti Warnasih','$2y$10$y8YKA148yLICYrW7aljF3exTov/RhvrLjo/OPTiR4i9rgNXTBvwEK','','Siti Warnasih, M.Si','siti.warnasih@unpak.ac.id',NULL,NULL,NULL,NULL,1,'2025-11-18 02:56:59','2025-11-28 08:16:33'),
            (489, UUID(),'Kotim Subandi','$2y$10$Mp9gHKhzY4ZD8Rg63lWDI.3glqu.7JfUcvjFMjD3E0VymH7yS6Taa','auditor2','Kotim Subandi, S.Kom., M.T','dummy@gmail.com',NULL,NULL,NULL,NULL,0,'2025-11-18 02:59:29','2025-11-18 02:59:29'),
            (490, UUID(),'Oktori Kiswati Zaini','$2y$10$6VbLUp.5pzxIlj7Lh9UYruy/dCnEOOse/FtWya3eZoFcww9tbRpzq','auditor2','Oktori Kiswati Zaini, S.E., M.M','dummy@gmail.com',NULL,NULL,NULL,NULL,0,'2025-11-18 03:04:30','2025-11-18 03:04:30'),
            (491, UUID(),'Delta Hadi Purnama','$2y$10$56ho7lpNkmeG30LYvUqJrOyUu0lUyUxzzTCjsA4Ui7Fdx3FC6AHXa','auditor2','Delta Hadi Purnama, S.E., M.E.Sy','dummy@gmail.com',NULL,NULL,NULL,NULL,0,'2025-11-18 03:06:05','2025-11-18 03:06:05'),
            (492, UUID(),'Mahifal','$2y$10$ar.TiK3rWqZiaURdnlEKEeJxCgxd6vq9mp.6a60AcrMtfpQP/IiLm','auditor2','Dr. Mahifal, S.H., M.H','dummy@gmail.com',NULL,NULL,NULL,NULL,0,'2025-11-18 03:12:21','2025-11-18 03:12:21'),
            (493, UUID(),'Atti Herawati','$2y$10$6Cynt4Y.62D5dIvRSezL1edR3YykT8bPqrMUng2v/lbjERHMGykfO','auditor2','Dr. Atti Herawati, M.Pd','dummy@gmail.com',NULL,NULL,NULL,NULL,0,'2025-11-18 03:13:55','2025-11-18 03:13:55'),
            (494, UUID(),'Anwar Sulaiman','$2y$10$yGj/ws9BnataxAOZLphVTuTJ/SX2wmjitKxjOKWMV.EJ3ndRlwjzG','auditor2','Dr. Anwar Sulaiman, S.E., M.M. S.H','dummy@gmail.com',NULL,NULL,NULL,NULL,0,'2025-11-18 03:17:24','2025-11-18 03:17:24'),
            (495, UUID(),'Desti Herawati','$2y$10$P6u5YQJ9UaOo0DDRPMdme.TgmHwP/po4NS3oFfQSEKie5a6pPdDwS','auditor2','Desti Herawati, M.Pd','dummy@gmail.com',NULL,NULL,NULL,NULL,0,'2025-11-18 03:18:18','2025-11-18 03:18:18'),
            (496, UUID(),'Diana Amaliasari','$2y$10$63M4MS2cG377YotjbdcFu.Kc339XTYLvZuJBDbsFV6jGj0a8s0Wie','auditor2','Diana Amaliasari, M.Si','dummy@gmail.com',NULL,NULL,NULL,NULL,0,'2025-11-18 03:20:13','2025-11-18 03:20:13'),
            (497, UUID(),'Siti Maimunah','$2y$10$ULGXj8hRcyQx91N6oJZXbe5F2pRCV..we3hxUUCWM5FeOQNtVbZZ2','auditor2','Dr. Siti Maimunah, S.E., M.Si','dummy@gmail.com',NULL,NULL,NULL,NULL,0,'2025-11-18 03:22:57','2025-11-18 03:22:57'),
            (498, UUID(),'Abdul Kohar','$2y$10$DcPpqO6SFkB7xmIJ0DWtLO3bJifXhgxmrtXfrFaTD/544My0wFugy','auditor2','Dr. Abdul Kohar, S.E., M.Ak','dummy@gmail.com',NULL,NULL,NULL,NULL,0,'2025-11-18 03:24:00','2025-11-18 03:24:00'),
            (499, UUID(),'Henny Suharyati','$2y$10$QVWvoz9zHb8OMkcAE2U6S.kFZCHpJO04tnNRb97v8zRWa4J4CP2HS','auditor1','Prof. Henny Suharyati, M.Pd','dummy@gmail.com',NULL,NULL,NULL,NULL,0,'2025-11-18 03:27:10','2025-11-18 03:27:10'),
            (500, UUID(),'Dewi Taurusyanti','$2y$10$v3GajUxO109oUgxEmNfN4eFbaTwE5wPfHOJUmrSlTO1gDRe.1xINm','auditor2','Dr.  Dewi Taurusyanti, S.E., M.M','dummy@gmail.com',NULL,NULL,NULL,NULL,0,'2025-11-18 03:28:26','2025-11-18 03:28:26'),
            (501, UUID(),'Yenny Febrianty','$2y$10$zcmAyyqEOzxwmZeTwXk6tuEGXiC8BFoQQtTopOM90pdfRDZGj47L6','auditor2','Dr. Yenny Febrianty, S.H., M.Hum., M.Kn','dummy@gmail.com',NULL,NULL,NULL,NULL,0,'2025-11-18 03:32:45','2025-11-18 03:32:45'),
            (502, UUID(),'Hasrul','$2y$10$52TbWvyZ0AMNlGUJRx1Nqe1p9SIPNgOvRUd2VvonqGSQuRCjQU.uG','auditor2','Dr. Hasrul, MM','dummy@gmail.com',NULL,NULL,NULL,NULL,0,'2025-11-18 03:33:31','2025-11-18 03:33:31'),
            (503, UUID(),'Elly Sukmanasa','$2y$10$67Q51YyGYL3v4MCf0pB.QOfB8NnlNwMXfBzl7KQ0hvGvSdalTXAEu','auditor2','Dr. Elly Sukmanasa, M .Pd','dummy@gmail.com',NULL,NULL,NULL,NULL,0,'2025-11-18 03:37:11','2025-11-18 03:37:11'),
            (504, UUID(),'Yunita Rahma','$2y$10$ocjML4tg1woR/rcNzhuBEOMyxiLwzRZT7c3clNwZwHVwbYeKTECv6','auditor2','Yunita Rahma, M.Kom','dummy@gmail.com',NULL,NULL,NULL,NULL,0,'2025-11-18 03:39:50','2025-11-18 03:39:50'),
            (505, UUID(),'Muhammad Reza','$2y$10$rq6fwxWrjZOG3ncTOllOp.xr953gDUjzyS3nn5KiEoOzt9/CNIRSu','auditor2','Muhammad Reza, M.Si','dummy@gmail.com',NULL,NULL,NULL,NULL,0,'2025-11-18 03:40:33','2025-11-18 03:40:33'),
            (506, UUID(),'Evyta Wismiana','$2y$10$X6bO.FXvxbHkmVFdGf0uCOcw6LsDsLYToIKiDwG7b6HoDwagis1tm','auditor2','Evyta Wismiana, S.T., M.T','dummy@gmail.com',NULL,NULL,NULL,NULL,0,'2025-11-18 03:42:35','2025-11-18 03:42:35'),
            (507, UUID(),'Restiawan Permana','$2y$10$gcPXVufgqxrSfp.ECT4Npu6.ideN4MX9k8pgfLFlXinJ4wIkFOm5a','auditor2','Restiawan Permana, M.Si','dummy@gmail.com',NULL,NULL,NULL,NULL,NULL,'2025-11-18 03:44:36','2025-11-18 03:44:53'),
            (508, UUID(),'Helmi Setia Pamungkas','$2y$10$qDsFf7n1sdqDrL8RRYzDEeu6CHMmrzXYUPie26/fIA1s.fY3OgwaC','auditor2','Helmi Setia Pamungkas, S.T., M.Si','dummy@gmail.com',NULL,NULL,NULL,NULL,0,'2025-11-18 04:09:45','2025-11-18 04:09:45'),
            (509, UUID(),'Lia Amelia Megawati','$2y$10$X9sk8HKDTu86lfIqUzU2DuZO54LECqpCV4HLv8/esbwaIO.4kYTtu','auditor2','Lia Amelia Megawati, S.pd., M.T','dummy@gmail.com',NULL,NULL,NULL,NULL,0,'2025-11-18 04:10:28','2025-11-18 04:10:28'),
            (510, UUID(),'Yan Noviar Nasution','$2y$10$WWcNPMusu9qjamfpWdnoPekx0MpHbMXd1FXd07fda72xWOtrZaaJC','auditor1','Dr. Yan Noviar Nasution, S.E., M.M., CA','dummy@gmail.com',NULL,NULL,NULL,NULL,0,'2025-11-18 04:11:45','2025-11-18 04:11:45'),
            (511, UUID(),'Komarudin','$2y$10$Dw653hB21hTX4N8M.050c.39b5AOxb.ZmJMLyRCxyUnTOUG5uXcAW','auditor2','Komarudin, M.H','dummy@gmail.com',NULL,NULL,NULL,NULL,0,'2025-11-18 04:12:43','2025-11-18 04:12:43'),
            (512, UUID(),'Dwi Rini Sovia Firdaus','$2y$10$q4kKllXwk6z1ZBn4zmVT3uX29/rK3af.I4qyr1eGZBgaU2IL1rpPW','auditor2','Dr. Dwi Rini Sovia Firdaus, M. Comn','dummy@gmail.com',NULL,NULL,NULL,NULL,0,'2025-11-18 04:13:41','2025-11-18 04:13:41'),
            (513, UUID(),'Lusi Agus Setiani','$2y$10$.RaYo/aSF/YIwKrr/gpVLOy2I6YOlz8mIo/RwKZiyGgYSbuS6shB2','auditor2','Dr. apt. Lusi Agus Setiani, M.Farm','dummy@gmail.com',NULL,NULL,NULL,NULL,0,'2025-11-18 04:14:31','2025-11-18 04:14:31'),
            (514, UUID(),'May Mulyaningsih','$2y$10$QsJBJO72hiuH1RSRh/xLn.93krIkX/e4tEEpZBauWH.5BfbTMr9dO','auditor2','May Mulyaningsih, S.E., M.Ak','dummy@gmail.com',NULL,NULL,NULL,NULL,0,'2025-11-18 04:17:18','2025-11-18 04:17:18'),
            (515, UUID(),'Sata Yoshida Srie Rahayu','$2y$10$BtXCae0raD2waLa987JayeYxCJRm6I0Lq95vARNaWiJhka/C/.9R2','auditor2','Prof. Dr. Sata Yoshida Srie Rahayu, M.Si','dummy@gmail.com',NULL,NULL,NULL,NULL,0,'2025-11-18 04:18:12','2025-11-18 04:18:12');;
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
