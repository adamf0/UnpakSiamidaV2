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
        CREATE TABLE v_fakultas_unit (
            id BIGINT AUTO_INCREMENT PRIMARY KEY,
            uuid VARCHAR(36) NULL,
            nama_fak_prod_unit VARCHAR(255) NULL,
            kode_jenjang VARCHAR(1) NULL,
            jenjang VARCHAR(5) NULL,
            jenjang_int VARCHAR(1) NULL,
            type VARCHAR(100) NULL, 
            fakultas VARCHAR(255) NULL
        );

        INSERT INTO v_fakultas_unit (uuid, nama_fak_prod_unit, kode_jenjang, jenjang, jenjang_int, type, fakultas)
        VALUES
        (null, 'ERRROR_1', null, null, null, null, null),
        (null, 'ERRROR_2', null, null, null, 'Unit', null),
        (e76447c9-097a-4a1f-8c85-066058e0c299, 'PUTIK', null, null, null, 'Unit', null),
        (UUID(), 'Teknik', 'C', 'S1', '1', 'Prodi', 'VOKASI'),
        (UUID(), 'Ekonomi', 'A', 'S1', '1', 'Fakultas', 'EKONOMI DAN BISNIS');
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
